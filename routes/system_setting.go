package routes

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/yoyo-inc/yoyo/common/db"
	"github.com/yoyo-inc/yoyo/common/logger"
	"github.com/yoyo-inc/yoyo/core"
	"github.com/yoyo-inc/yoyo/errs"
	"github.com/yoyo-inc/yoyo/models"
	"github.com/yoyo-inc/yoyo/services/audit_log"
	"github.com/yoyo-inc/yoyo/vo"
	"gorm.io/gorm"
)

type systemSettingController struct{}

// QuerySystemSettings
//	@Summary	查询系统设置
//	@Tags		system
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	core.Response{data=models.SystemSetting}
//	@Security	JWT
//	@Router		/system/settings [get]
func (*systemSettingController) QuerySystemSettings(c *gin.Context) {
	var systemSetting models.SystemSetting
	if res := db.Client.Preload("Logo").Model(&models.SystemSetting{}).First(&systemSetting); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			core.OK(c, map[string]string{})
			return
		} else {
			logger.Error(res.Error)
			c.Error(errs.ErrQuerySystemSetting)
			return
		}
	}

	core.OK(c, systemSetting)
}

// UpdateSystemSetting
//	@Summary	更新系统设置
//	@Tags		system
//	@Accept		json
//	@Produce	json
//	@Param		body	body		vo.UpdateSystemSettingVO	true	"参数"
//	@Success	200		{object}	core.Response{data=bool}
//	@Security	JWT
//	@Router		/system/setting [put]
func (*systemSettingController) UpdateSystemSetting(c *gin.Context) {
	var query vo.UpdateSystemSettingVO
	if err := c.ShouldBindJSON(&query); err != nil {
		logger.Error(err)
		c.Error(core.NewParameterError(err))
		return
	}

	if res := db.Client.Model(&models.SystemSetting{IModel: core.IModel{ID: query.ID}}).Updates(query.SystemSetting); res.Error != nil {
		logger.Error(res.Error)
		c.Error(errs.ErrUpdateSystemSetting)
		return
	}

	json := jsoniter.ConfigCompatibleWithStandardLibrary
	s, _ := json.MarshalToString(query)
	audit_log.Success(c, "系统设置", "更新", fmt.Sprintf("配置内容:%s", s))
	core.OK(c, true)
}

func (sc *systemSettingController) Setup(r *gin.RouterGroup) {
	r.GET("/system/settings", sc.QuerySystemSettings).PUT("/system/setting", sc.UpdateSystemSetting)
}
