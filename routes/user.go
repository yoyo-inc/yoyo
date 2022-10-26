package routes

import (
	"fmt"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/gin-gonic/gin"
	jwt "github.com/yoyo-inc/gin-jwt/v3"
	"github.com/yoyo-inc/yoyo/common/db"
	"github.com/yoyo-inc/yoyo/common/logger"
	"github.com/yoyo-inc/yoyo/core"
	"github.com/yoyo-inc/yoyo/errs"
	"github.com/yoyo-inc/yoyo/middlewares"
	"github.com/yoyo-inc/yoyo/models"
	"github.com/yoyo-inc/yoyo/services/audit_log"
	"github.com/yoyo-inc/yoyo/vo"
)

type userController struct{}

// QueryUsers
// @Summary  查询用户列表
// @Tags     user
// @Produce  json
// @Param    query  query     vo.QueryUserVO       false  "参数"
// @Param    query  query     models.Pagination  false  "参数"
// @Success  200    {object}  core.Response{data=core.PaginatedData{list=[]models.User}}
// @Router   /users [get]
func (*userController) QueryUsers(c *gin.Context) {
	var query vo.QueryUserVO
	if err := c.ShouldBindQuery(&query); err != nil {
		c.Error(core.NewParameterError(err.Error()))
		return
	}

	queries := core.GetPaginatedQuery(&models.User{})
	for i := range queries {
		if query.Username != "" {
			queries[i].Where("username like ?", "%"+query.Username+"%")
			query.Username = ""
		}
	}

	var users []models.User
	if res := queries[0].Scopes(core.Paginator(c)).Where(query).Find(&users); res.Error != nil {
		logger.Error(res.Error)
		c.Error(errs.ErrQueryUser)
		return
	}

	var total int64
	if res := queries[1].Where(query).Count(&total); res.Error != nil {
		logger.Error(res.Error)
		c.Error(errs.ErrQueryUser)
		return
	}

	core.OK(c, core.Paginated(users, total))
}

// CreateUser
// @Summary  创建用户
// @Tags     user
// @Accept   json
// @Produce  json
// @Param    query  body      vo.UserVO  true  "用户信息"
// @Success  200    {object}  core.Response{data=models.User}
// @Router   /user [post]
func (*userController) CreateUser(c *gin.Context) {
	var query vo.UserVO
	if err := c.ShouldBindJSON(&query); err != nil {
		c.Error(core.NewParameterError(err.Error()))
		return
	}

	if res := db.Client.Create(&query); res.Error != nil {
		logger.Error(res.Error)
		c.Error(errs.ErrCreateUser)
		audit_log.Fail(c, "用户", "新增", query.Username)
		return
	}

	audit_log.Success(c, "用户", "新增", query.Username)
	core.OK(c, query)
}

// UpdateUser
// @Summary  更新用户
// @Tags     user
// @Accept   json
// @Produce  json
// @Param    query  body      vo.UserVO  true  "用户信息"
// @Success  200    {object}  core.Response{data=bool}
// @Router   /user [put]
func (*userController) UpdateUser(c *gin.Context) {
	var query vo.UserVO
	if err := c.ShouldBindJSON(&query); err != nil {
		c.Error(core.NewParameterError(err.Error()))
		return
	}

	if res := db.Client.Save(&query.User); res.Error != nil {
		logger.Error(res.Error)
		c.Error(errs.ErrUpdateUser)
		audit_log.Fail(c, "用户", "更新", query.Username)
		return
	}

	// update user associations
	// update user roles
	if query.Roles != nil {
		roles := slice.Map(query.Roles, func(_ int, roleID string) models.Role {
			return models.Role{Model: core.Model{ID: roleID}}
		})
		if err := db.Client.Model(&models.User{}).Association("Roles").Replace(roles); err != nil {
			logger.Error(err)
			c.Error(errs.ErrUpdateUser)
			audit_log.Fail(c, "用户", "更新", query.Username)
		}
	}
	audit_log.Success(c, "用户", "更新", query.Username)
	core.OK(c, true)
}

// DeleteUser
// @Summary  删除用户
// @Tags     user
// @Produce  json
// @Param    userID  path      string  true  "用户ID"
// @Success  200     {object}  core.Response{data=bool}
// @Router   /user/{userID} [delete]
func (*userController) DeleteUser(c *gin.Context) {
	userID := c.Param("userID")

	var existUser models.User
	if res := db.Client.Model(&models.User{Model: core.Model{ID: userID}}).First(&existUser); res.Error != nil {
		logger.Error(res.Error)
		c.Error(errs.ErrUsernameNotExists)
		audit_log.Fail(c, "用户", "删除", fmt.Sprintf("用户ID(%s}不存在", userID))
		return
	}

	if res := db.Client.Delete(&models.User{}, "id = ?", userID); res.Error != nil {
		logger.Error(res.Error)
		c.Error(errs.ErrDeleteUser)
		audit_log.Fail(c, "用户", "删除", existUser.Username)
		return
	}
	audit_log.Success(c, "用户", "删除", existUser.Username)
	core.OK(c, true)
}

// QueryCurrentUser
// @Summary  查询当前用户信息
// @Tags     user
// @Accept   json
// @Produce  json
// @Success  200  {object}  core.Response{data=models.User}
// @Router   /user/current [get]
func (*userController) QueryCurrentUser(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	userID, _ := claims[middlewares.IdentityKey].(string)

	var user models.User
	if result := db.Client.Omit("password").Find(&user, userID); result.Error != nil {
		logger.Error(result.Error)
		c.Error(errs.ErrQueryCurrentUser)
		return
	}
	core.OK(c, user)
}

func (user *userController) Setup(r *gin.RouterGroup) {
	r.GET("/users", user.QueryUsers).
		POST("/user", user.CreateUser).
		PUT("/user", user.UpdateUser).
		DELETE("/user/:userID", user.DeleteUser).
		GET("/user/current", user.QueryCurrentUser)
}
