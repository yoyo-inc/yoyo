package routes

import (
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/yoyo-inc/yoyo/common/db"
	"github.com/yoyo-inc/yoyo/common/logger"
	"github.com/yoyo-inc/yoyo/core"
	"github.com/yoyo-inc/yoyo/errs"
	"github.com/yoyo-inc/yoyo/models"
	"github.com/yoyo-inc/yoyo/services/audit_log"
	"github.com/yoyo-inc/yoyo/vo"
	"gorm.io/datatypes"
	"gorm.io/gorm"
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
			systemSecurityVO.SystemSecurity.LoginIPWhitelist = datatypes.JSON(b)
		}
	}

	var res *gorm.DB
	if systemSecurityVO.ID == 0 {
		var count int64
		if res := db.Client.Model(&models.SystemSecurity{}).Count(&count); res.Error == nil {
			if count > 0 {
				db.Client.Delete(&models.SystemSecurity{})
			}
		}
		res = db.Client.Create(&(systemSecurityVO.SystemSecurity))
	} else {
		res = db.Client.Model(&models.SystemSecurity{IModel: core.IModel{ID: systemSecurityVO.ID}}).Select("*").Updates(systemSecurityVO.SystemSecurity)
	}

	if res.Error != nil {
		logger.Error(res.Error)
		c.Error(errs.ErrUpdateSystemSecurity)
		audit_log.Fail(c, "系统安全", "更新", "")
		return
	}

	audit_log.Success(c, "系统安全", "更新", "")
	core.OK(c, true)
}

func (ssc *systemSecurityController) Setup(r *gin.RouterGroup) {
	r.GET("/system_security", ssc.QuerySystemSecurity).PUT("/system_security", ssc.UpdateSystemSecurity)
}
