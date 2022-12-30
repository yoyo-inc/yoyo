package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/yoyo-inc/yoyo/core"
	"github.com/yoyo-inc/yoyo/middlewares"
)

var controllers []core.Controller

func init() {
	controllers = []core.Controller{
		&loginController{},
		&userController{},
		&roleController{},
		&permissionController{},
		&notificationController{},
		&auditLogController{},
		&alertController{},
		&systemSettingController{},
		&resourceController{},
		&reportController{},
		&schedJobController{},
		&dictController{},
	}
}

func AddService(controller core.Controller) {
	controllers = append(controllers, controller)
}

func Setup(r *gin.RouterGroup) {
	middlewares.SecurityMiddleware.AddIgnore(
		"/api/login",
		"/api/swagger/*path",
		"/api/system/settings",
		"/api/resource/*path",
		"/api/report/preview/*path",
	)

	for _, sc := range controllers {
		sc.Setup(r)
	}
}
