package routes

import (
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/yoyo-inc/yoyo/common/db"
	"github.com/yoyo-inc/yoyo/common/logger"
	"github.com/yoyo-inc/yoyo/core"
	"github.com/yoyo-inc/yoyo/errs"
	"github.com/yoyo-inc/yoyo/models"
	"github.com/yoyo-inc/yoyo/vo"
	"gorm.io/datatypes"
)

type systemSecurityController struct{}

// QuerySystemSecurity
// @Summary 查询系统安全设置
// @Tags    systemSecurity
// @Accept  json
// @Produce json
// @Success 200   {object} core.Response{data=models.SystemSecurity}
// @Security JWT
// @Router  /system_security [get]
func (*systemSecurityController) QuerySystemSecurity(c *gin.Context) {
	var systemSecurity models.SystemSecurity

	if res := db.Client.Model(&models.SystemSecurity{}).Find(&systemSecurity); res.Error != nil {
		logger.Error(res.Error)
		c.Error(errs.ErrQuerySystemSecurity)
		return
	}

	core.OK(c, systemSecurity)
}

// UpdateSystemSecurity
// @Summary 更新系统安全设置
// @Tags    systemSecurity
// @Accept  json
// @Produce json
// @Param   body body     vo.SystemSecurityVO true "参数"
// @Success 200   {object} core.Response{data=bool}
// @Security JWT
// @Router  /system_security [put]
func (*systemSecurityController) UpdateSystemSecurity(c *gin.Context) {
	var systemSecurityVO vo.SystemSecurityVO
	if err := c.ShouldBindJSON(&systemSecurityVO); err != nil {
		c.Error(core.NewParameterError(err))
		return
	}

	if systemSecurityVO.LoginIPWhitelist != nil {
		json := jsoniter.ConfigCompatibleWithStandardLibrary
		if b, err := json.Marshal(systemSecurityVO.LoginIPWhitelist); err != nil {
			systemSecurityVO.SystemSecurity.LoginIPWhitelist = (*datatypes.JSON)(&b)
		}
	}

	if res := db.Client.Model(&models.SystemSecurity{IModel: core.IModel{ID: systemSecurityVO.ID}}).Updates(systemSecurityVO.SystemSecurity); res.Error != nil {
		logger.Error(res.Error)
		c.Error(errs.ErrUpdateSystemSecurity)
		return
	}

	core.OK(c, true)
}

func (ssc *systemSecurityController) Setup(r *gin.RouterGroup) {
	r.GET("/system_security", ssc.QuerySystemSecurity)
}
