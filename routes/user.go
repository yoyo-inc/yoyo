package routes

import (
	"fmt"
	"strconv"

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
	"gorm.io/gorm/clause"
)

type userController struct{}

// QueryUsers
//
//	@Summary	查询用户列表
//	@Tags		user
//	@Produce	json
//	@Param		query	query		vo.QueryUserVO		false	"参数"
//	@Param		query	query		models.Pagination	false	"参数"
//	@Success	200		{object}	core.Response{data=core.PaginatedData{list=[]models.User}}
//	@Security	JWT
//	@Router		/users [get]
func (*userController) QueryUsers(c *gin.Context) {
	var query vo.QueryUserVO
	if err := c.ShouldBindQuery(&query); err != nil {
		_ = c.Error(core.NewParameterError(err))
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
	if res := queries[0].Preload("Roles").Omit("password").Scopes(core.Paginator(c), core.DateTimeRanger(c, "")).
		Where(query).Order("create_time desc").Find(&users); res.Error != nil {
		logger.Error(res.Error)
		_ = c.Error(errs.ErrQueryUser)
		return
	}

	var total int64
	if res := queries[1].Scopes(core.DateTimeRanger(c, "")).Where(query).Count(&total); res.Error != nil {
		logger.Error(res.Error)
		_ = c.Error(errs.ErrQueryUser)
		return
	}

	core.OK(c, core.Paginated(users, total))
}

// CreateUser
//
//	@Summary	创建用户
//	@Tags		user
//	@Accept		json
//	@Produce	json
//	@Param		query	body		vo.UserVO	true	"用户信息"
//	@Success	200		{object}	core.Response{data=models.User}
//	@Security	JWT
//	@Router		/user [post]
func (*userController) CreateUser(c *gin.Context) {
	var query vo.UserVO
	if err := c.ShouldBindJSON(&query); err != nil {
		_ = c.Error(core.NewParameterError(err))
		return
	}

	// judge whether user exists
	var count int64
	if res := db.Client.Model(&models.User{}).Where("username = ?", query.Username).Count(&count); res.Error != nil {
		logger.Error(res.Error)
		_ = c.Error(errs.ErrCreateUser)
		audit_log.Fail(c, "用户", "新增", "用户名："+query.Username)
		return
	}
	if count > 0 {
		_ = c.Error(errs.ErrExistUser)
		audit_log.Fail(c, "用户", "新增", "用户已存在，用户名："+query.Username)
		return
	}

	user := query.User
	user.Roles = slice.Map(query.Roles, func(_ int, roleID int) models.Role {
		return models.Role{IModel: core.IModel{ID: roleID}}
	})

	if user.OrganizationID == 0 {
		user.OrganizationID = 1
	}
	if res := db.Client.Create(&user); res.Error != nil {
		logger.Error(res.Error)
		_ = c.Error(errs.ErrCreateUser)
		audit_log.Fail(c, "用户", "新增", "用户名："+query.Username)
		return
	}

	audit_log.Success(c, "用户", "新增", "用户名："+query.Username)
	core.OK(c, query)
}

// UpdateUser
//
//	@Summary	更新用户
//	@Tags		user
//	@Accept		json
//	@Produce	json
//	@Param		query	body		vo.UserVO	true	"用户信息"
//	@Success	200		{object}	core.Response{data=bool}
//	@Security	JWT
//	@Router		/user [put]
func (*userController) UpdateUser(c *gin.Context) {
	var query vo.UserVO
	if err := c.ShouldBindJSON(&query); err != nil {
		_ = c.Error(core.NewParameterError(err))
		return
	}

	if res := db.Client.Where("id = ?", query.ID).Updates(&query.User); res.Error != nil {
		logger.Error(res.Error)
		_ = c.Error(errs.ErrUpdateUser)
		audit_log.Fail(c, "用户", "更新", "用户名："+query.Username)
		return
	}

	// update user associations
	// update user roles
	if query.Roles != nil {
		roles := slice.Map(query.Roles, func(_ int, roleID int) models.Role {
			return models.Role{IModel: core.IModel{ID: roleID}}
		})
		if err := db.Client.Model(&models.User{IModel: core.IModel{ID: query.ID}}).Association("Roles").Replace(roles); err != nil {
			logger.Error(err)
			_ = c.Error(errs.ErrUpdateUser)
			audit_log.Fail(c, "用户", "更新", "用户名："+query.Username)
			return
		}
	}
	audit_log.Success(c, "用户", "更新", "用户名："+query.Username)
	core.OK(c, true)
}

// DeleteUser
//
//	@Summary	删除用户
//	@Tags		user
//	@Produce	json
//	@Param		userID	path		string	true	"用户ID"
//	@Success	200		{object}	core.Response{data=bool}
//	@Security	JWT
//	@Router		/user/{userID} [delete]
func (*userController) DeleteUser(c *gin.Context) {
	userID := c.Param("userID")
	id, _ := strconv.Atoi(userID)

	var existUser models.User
	if res := db.Client.Model(&models.User{IModel: core.IModel{ID: id}}).First(&existUser); res.Error != nil {
		logger.Error(res.Error)
		_ = c.Error(errs.ErrUsernameNotExists)
		audit_log.Fail(c, "用户", "删除", fmt.Sprintf("用户ID(%s)不存在", userID))
		return
	}

	if res := db.Client.Select(clause.Associations).Delete(&models.User{IModel: core.IModel{ID: id}}); res.Error != nil {
		logger.Error(res.Error)
		_ = c.Error(errs.ErrDeleteUser)
		audit_log.Fail(c, "用户", "删除", "用户名："+existUser.Username)
		return
	}
	audit_log.Success(c, "用户", "删除", "用户名："+existUser.Username)
	core.OK(c, true)
}

// QueryCurrentUser
//
//	@Summary	查询当前用户信息
//	@Tags		user
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	core.Response{data=models.User}
//	@Security	JWT
//	@Router		/user/current [get]
func (*userController) QueryCurrentUser(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	userID, _ := claims[middlewares.IdentityKey].(string)

	var user models.User
	if result := db.Client.Omit("password").Preload("Roles").Find(&user, userID); result.Error != nil {
		logger.Error(result.Error)
		_ = c.Error(errs.ErrQueryCurrentUser)
		return
	}
	core.OK(c, user)
}

// QueryCurrentUserPermissions
//
//	@Summary	查询当前用户权限信息
//	@Tags		user
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	core.Response{data=array,string}
//	@Security	JWT
//	@Router		/user/current/permissions [get]
func (*userController) QueryCurrentUserPermissions(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	userID, _ := claims[middlewares.IdentityKey].(string)
	var user models.User
	if result := db.Client.Preload("Roles.Permissions").Find(&user, userID); result.Error != nil {
		logger.Error(result.Error)
		_ = c.Error(errs.ErrQueryCurrentUser)
		return
	}

	permissions := make([]string, 0)
	for _, role := range user.Roles {
		for _, permission := range role.Permissions {
			permissions = append(permissions, permission.Name)
		}
	}

	core.OK(c, permissions)
}

func (user *userController) Setup(r *gin.RouterGroup) {
	r.GET("/users", user.QueryUsers).
		POST("/user", user.CreateUser).
		PUT("/user", user.UpdateUser).
		DELETE("/user/:userID", user.DeleteUser).
		GET("/user/current", user.QueryCurrentUser).
		GET("/user/current/permissions", user.QueryCurrentUserPermissions)
}
