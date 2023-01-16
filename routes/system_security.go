package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/yoyo-inc/yoyo/common/db"
	"github.com/yoyo-inc/yoyo/common/logger"
	"github.com/yoyo-inc/yoyo/core"
	"github.com/yoyo-inc/yoyo/models"
)

type systemSecurityController struct{}

func (*systemSecurityController) QuerySystemSecurity(c *gin.Context) {
	var systemSecurity models.SystemSecurity

	if res := db.Client.Model(&models.SystemSecurity{}).Find(&systemSecurity); res.Error != nil {
		logger.Error(res.Error)
		c.Error(res.Error)
		return
	}

	core.OK(c, systemSecurity)
}

func (*systemSecurityController) UpdateSystemSecurity() {}

func (ssc *systemSecurityController) Setup(r *gin.RouterGroup) {
}
