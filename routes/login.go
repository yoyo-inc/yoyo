package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ypli0629/yoyo/middlewares"
)

type LoginController struct{}

func (login *LoginController) Setup(r *gin.RouterGroup) {
	r.GET("/logout", middlewares.SecurityMiddleware.LogoutHandler)
}

func (login *LoginController) SetupNoSecurity(r *gin.RouterGroup) {
	r.POST("/login", middlewares.SecurityMiddleware.LoginHandler)
}
