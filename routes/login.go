package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ypli0629/yoyo/middlewares"
)

type loginController struct{}

func (login *loginController) Setup(r *gin.RouterGroup) {
	r.GET("/logout", middlewares.SecurityMiddleware.LogoutHandler)
}

func (login *loginController) SetupNoSecurity(r *gin.RouterGroup) {
	r.POST("/login", middlewares.SecurityMiddleware.LoginHandler)
}
