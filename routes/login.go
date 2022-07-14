package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/yoyo-inc/yoyo/middlewares"
)

type loginController struct{}

func (login *loginController) Setup(r *gin.RouterGroup) {
	r.POST("/login", middlewares.SecurityMiddleware.LoginHandler).GET("/logout", middlewares.SecurityMiddleware.LogoutHandler)
}
