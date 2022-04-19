package core

import "github.com/gin-gonic/gin"

type Controller interface {
	Setup(r *gin.RouterGroup)
}

type NoSecurityController interface {
	SetupNoSecurity(r *gin.RouterGroup)
}
