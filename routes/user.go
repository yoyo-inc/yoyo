package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ypli0629/yoyo/common/db"
	"github.com/ypli0629/yoyo/core"
	"github.com/ypli0629/yoyo/errors"
	"github.com/ypli0629/yoyo/models"
)

type userController struct{}

// RetrieveUser returns user
// @Summary query user
// @Tags user
// @Produce json
// @Param query query models.QueryUser true "request param"
// @Success 200 {array} core.Response{data=models.User}
// @Router /user [get]
func (*userController) RetrieveUser(c *gin.Context) {
	var query models.QueryUser
	if err := c.ShouldBindQuery(&query); err != nil {
		c.Error(core.NewParameterError(err.Error()))
		return
	}

	var users []models.User
	if result := db.Client.Where("username", query.Username).Where("phone", query.Phone).Find(&users); result.Error != nil {
		c.Error(errors.ErrQueryUser)
		return
	}

	core.OK(c, users)
}

// CreateUser creates user
// @Summary create user
// @Tags user
// @Accept json
// @Produce json
// @Param query body models.User true "user information"
// @Success 200 {object} core.Response{data=models.User}
// @Router /user [post]
func (*userController) CreateUser(c *gin.Context) {
	var query models.User
	if err := c.ShouldBindJSON(&query); err != nil {
		c.Error(core.NewParameterError(err.Error()))
		return
	}
	if result := db.Client.Create(&query); result.Error != nil {
		c.Error(errors.ErrCreateUser)
		return
	}
	core.OK(c, query)
}

func (*userController) DeleteUser(c *gin.Context) {

}

func (user *userController) Setup(r *gin.RouterGroup) {
	r.GET("/user", user.RetrieveUser).POST("/user", user.CreateUser)
}
