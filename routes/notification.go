package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/yoyo-inc/yoyo/common/hub"
)

type notificationController struct {
}

func (notification *notificationController) PushNotification(c *gin.Context) {
	hub.Register(c.Request, c.Writer)
}

func (notification *notificationController) Setup(r *gin.RouterGroup) {
	r.GET("/ws/notification", notification.PushNotification)
}
