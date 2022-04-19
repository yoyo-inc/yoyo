package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ypli0629/yoyo/core"
)

var (
	SecurityControllers   []core.Controller
	NoSecurityControllers []core.NoSecurityController
)

func Init() {
	NoSecurityControllers = []core.NoSecurityController{
		&LoginController{},
	}
	SecurityControllers = []core.Controller{
		&LoginController{},
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