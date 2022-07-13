package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/yoyo-inc/yoyo/core"
)

var (
	SecurityControllers   []core.Controller
	NoSecurityControllers []core.NoSecurityController
)

func Init() {
	NoSecurityControllers = []core.NoSecurityController{
		&loginController{},
	}
	SecurityControllers = []core.Controller{
		&loginController{},
		&userController{},
	}
}

func SetupSecurity(r *gin.RouterGroup) {
	for _, sc := range SecurityControllers {
		sc.Setup(r)
	}
}

func SetupNoSecurity(r *gin.RouterGroup) {
	Init()

	for _, nsc := range NoSecurityControllers {
		nsc.SetupNoSecurity(r)
	}
}
