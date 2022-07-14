package routes

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/yoyo-inc/gin-jwt/v3"
	"github.com/yoyo-inc/yoyo/common/db"
	"github.com/yoyo-inc/yoyo/common/logger"
	"github.com/yoyo-inc/yoyo/core"
	"github.com/yoyo-inc/yoyo/errs"
	"github.com/yoyo-inc/yoyo/middlewares"
	"github.com/yoyo-inc/yoyo/models"
)

type userController struct{}

// RetrieveUser
// @Summary  查询用户列表
// @Tags     user
// @Produce  json
// @Param    query  query    models.QueryUser   false  "参数"
// @Param    query  query    models.Pagination  false  "参数"
// @Success  200    {array}  core.Response{data=core.PaginatedData{list=[]models.User}}
// @Router   /user [get]
func (*userController) RetrieveUser(c *gin.Context) {
	var query models.QueryUser
	if err := c.ShouldBindQuery(&query); err != nil {
		c.Error(core.NewParameterError(err.Error()))
		return
	}

	var users []models.User
	if result := db.Client.Model(&models.User{}).Where(&query).Scopes(core.Paginator(c)).Find(&users); result.Error != nil {
		logger.Error(result.Error)
		c.Error(errs.ErrQueryUser)
		return
	}

	var total int64
	if result := db.Client.Model(&models.User{}).Where(&query).Count(&total); result.Error != nil {
		logger.Error(result.Error)
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
// @Param    query  body      models.User  true  "用户信息"
// @Success  200    {object}  core.Response{data=models.User}
// @Router   /user [post]
func (*userController) CreateUser(c *gin.Context) {
	var query models.User
	if err := c.ShouldBindJSON(&query); err != nil {
		c.Error(core.NewParameterError(err.Error()))
		return
	}
	if result := db.Client.Create(&query); result.Error != nil {
		logger.Error(result.Error)
		c.Error(errs.ErrCreateUser)
		return
	}
	core.OK(c, query)
}

// UpdateUser
// @Summary  更新用户
// @Tags     user
// @Accept   json
// @Produce  json
// @Param    query  body      models.User  true  "用户信息"
// @Success  200    {object}  core.Response{data=bool}
// @Router   /user [put]
func (*userController) UpdateUser(c *gin.Context) {
	var query models.User
	if err := c.ShouldBindJSON(&query); err != nil {
		c.Error(core.NewParameterError(err.Error()))
		return
	}

	if result := db.Client.Save(&query); result.Error != nil {
		logger.Error(result.Error)
		c.Error(errs.ErrUpdateUser)
		return
	}
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

	if result := db.Client.Delete(&models.User{}, userID); result.Error != nil {
		logger.Error(result.Error)
		c.Error(errs.ErrDeleteUser)
		return
	}
	core.OK(c, true)
}

// RetrieveCurrentUser
// @Summary      查询当前用户信息
// @Tags         user
// @Accept       json
// @Produce      json
// @Success      200  {object}  core.Response{data=models.User}
// @Router       /user/current [get]
func (*userController) RetrieveCurrentUser(c *gin.Context) {
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
	r.GET("/user", user.RetrieveUser).
		POST("/user", user.CreateUser).
		PUT("/user", user.UpdateUser).
		DELETE("/user/:userID", user.DeleteUser).
		GET("/user/current", user.RetrieveCurrentUser)
}
