package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ypli0629/yoyo/common/db"
	"github.com/ypli0629/yoyo/core"
	"github.com/ypli0629/yoyo/errors"
	"github.com/ypli0629/yoyo/models"
)

type userController struct{}

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

func (user *userController) Setup(r *gin.RouterGroup) {
	r.POST("/user", user.CreateUser)
}
