package routes

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/yoyo-inc/yoyo/common/db"
	"github.com/yoyo-inc/yoyo/common/logger"
	"github.com/yoyo-inc/yoyo/core"
	"github.com/yoyo-inc/yoyo/errs"
	"github.com/yoyo-inc/yoyo/models"
	"github.com/yoyo-inc/yoyo/services"
	"github.com/yoyo-inc/yoyo/services/audit_log"
	"github.com/yoyo-inc/yoyo/vo"
	"gorm.io/gorm"
)

type alertController struct{}

func init() {
	AddService(&alertController{})
}

// QueryAlerts
// @Summary 查询告警列表
// @Tags    alert
// @Produce json
// @Param   query query    vo.QueryAlertVO   false "参数"
// @Param   query query    models.Pagination false "参数"
// @Success 200   {object} core.Response{data=core.PaginatedData{list=[]models.Alert}}
// @Security JWT
// @Router  /alerts [get]
func (*alertController) QueryAlerts(c *gin.Context) {
	var query vo.QueryAlertVO
	if err := c.ShouldBindQuery(&query); err != nil {
		logger.Error(err)
		c.Error(core.NewParameterError(err))
		return
	}
	queries := core.GetPaginatedQuery(&models.Alert{})

	var alerts []models.Alert
	if res := queries[0].Scopes(core.Paginator(c)).Find(&alerts); res.Error != nil {
		logger.Error(res.Error)
		c.Error(errs.ErrQueryAlert)
		return
	}

	var count int64
	if res := queries[1].Count(&count); res.Error != nil {
		logger.Error(res.Error)
		c.Error(errs.ErrQueryAlert)
		return
	}

	core.OK(c, core.Paginated(alerts, count))
}

// UpdateAlert
// @Summary 更新告警
// @Tags    alert
// @Accept  json
// @Produce json
// @Param   query body     vo.UpdateAlertVO true "告警信息"
// @Success 200   {object} core.Response{data=bool}
// @Security JWT
// @Router  /alert [put]
func (*alertController) UpdateAlert(c *gin.Context) {
	var query vo.UpdateAlertVO
	if err := c.ShouldBindJSON(&query); err != nil {
		logger.Error(err)
		c.Error(core.NewParameterError(err))
		return
	}

	if res := db.Client.Updates(&query); res.Error != nil {
		logger.Error(res.Error)
		c.Error(errs.ErrUpdateAlert)
		return
	}

	audit_log.Success(c, "告警", "处置", "")

	core.OK(c, true)
}

// GetAlertConfig
// @Summary 查询告警配置
// @Tags    alert
// @Accept  json
// @Produce json
// @Success 200   {object} core.Response{data=models.AlertConfig}
// @Security JWT
// @Router  /alert/config [get]
func (*alertController) GetAlertConfig(c *gin.Context) {
	var config models.AlertConfig
	if res := db.Client.Model(&models.AlertConfig{}).First(&config); res.Error != nil {
		if !errors.Is(res.Error, gorm.ErrRecordNotFound) {
			logger.Error(res.Error)
			c.Error(errs.ErrQueryAlertConfig)
			return
		}
	}

	core.OK(c, config)
}

// UpdateAlertConfig
// @Summary 更新告警配置
// @Tags    alert
// @Accept  json
// @Produce json
// @Param body body models.AlertConfig true "参数"
// @Success 200   {object} core.Response{data=bool}
// @Security JWT
// @Router  /alert/config [put]
func (*alertController) UpdateAlertConfig(c *gin.Context) {
	var query vo.UpdateAlertConfigVO
	if err := c.ShouldBindJSON(&query); err != nil {
		c.Error(core.NewParameterError(err))
		return
	}

	json := jsoniter.ConfigCompatibleWithStandardLibrary
	if err := db.Client.Transaction(func(tx *gorm.DB) error {
		if err := tx.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.AlertConfig{}).Error; err != nil {
			return err
		}

		model := query.AlertConfig
		rawSmtpReceivers, err := json.Marshal(query.SmptReceivers)
		if err != nil {
			return err
		}
		model.SmtpReceivers = rawSmtpReceivers
		if err := tx.Create(&model).Error; err != nil {
			return err
		}

		if err := services.GenerateAlertManagerConfig(model); err != nil {
			return err
		}

		s, _ := json.MarshalToString(model)
		audit_log.Success(c, "告警配置", "更新", fmt.Sprintf("告警配置内容:%s", s))

		return nil
	}); err != nil {
		logger.Error(err)
		c.Error(errs.ErrUpdateAlertConfig)

		s, _ := json.MarshalToString(query.AlertConfig)
		audit_log.Fail(c, "告警配置", "更新", fmt.Sprintf("告警配置内容:%s", s))
		return
	}

	core.OK(c, true)
}

// QueryAlertAccesses
// @Summary 查询告警接入
// @Tags    alert
// @Accept  json
// @Produce json
// @Param query query vo.QueryAlertAccessVO true "参数"
// @Success 200   {object} core.Response{data=core.PaginatedData{data=models.AlertAccess}}
// @Security JWT
// @Router  /alert/accesses [get]
func (*alertController) QueryAlertAccesses(c *gin.Context) {
	var query vo.QueryAlertAccessVO
	if err := c.ShouldBindQuery(&query); err != nil {
		logger.Error(err)
		c.Error(core.NewParameterError(err))
		return
	}

	queries := core.GetPaginatedQuery(&models.AlertAccess{})

	for i := range queries {
		if query.AccessIP != "" {
			queries[i].Where("access_ip like ?", "%"+query.AccessIP+"%")
		}
	}
	var alertAccesses []models.AlertAccess
	if res := queries[0].Scopes(core.Paginator(c), core.DateTimeRanger(c)).Find(&alertAccesses); res.Error != nil {
		logger.Error(res.Error)
		c.Error(errs.ErrQueryAlertAccess)
		return
	}
	var count int64
	if res := queries[1].Count(&count); res.Error != nil {
		logger.Error(res.Error)
		c.Error(errs.ErrQueryAlertAccess)
		return
	}

	core.OK(c, core.Paginated(alertAccesses, count))
}

// CreateAlertAccess
// @Summary 创建告警接入
// @Tags    alert
// @Accept  json
// @Produce json
// @Param   body body     models.AlertAccess true "参数"
// @Success 200   {object} core.Response{data=bool}
// @Security JWT
// @Router  /alert/access [post]
func (*alertController) CreateAlertAccess(c *gin.Context) {
	var access models.AlertAccess
	if err := c.ShouldBindJSON(&access); err != nil {
		logger.Error(err)
		c.Error(errs.ErrCreateAlertAccess)
		return
	}

	if res := db.Client.Create(&access); res.Error != nil {
		logger.Error(res.Error)
		c.Error(errs.ErrCreateAlertAccess)
		return
	}

	audit_log.Success(c, "告警接入", "新增", fmt.Sprintf("接入IP:%s", access.AccessIP))

	core.OK(c, true)
}

// UpdateAlertAccess
// @Summary 更新告警接入
// @Tags    alert
// @Accept  json
// @Produce json
// @Param   body body     vo.UpdateAlertAccessVO true "参数"
// @Success 200   {object} core.Response{data=bool}
// @Security JWT
// @Router  /alert/access [put]
func (*alertController) UpdateAlertAccess(c *gin.Context) {
	var query vo.UpdateAlertAccessVO
	if err := c.ShouldBindJSON(&query); err != nil {
		logger.Error(err)
		c.Error(core.NewParameterError(err))
		return
	}

	if res := db.Client.Where("id = ?", query.ID).Updates(&query.AlertAccess); res.Error != nil {
		logger.Error(res.Error)
		c.Error(errs.ErrUpdateAlertAccess)
		return
	}
	audit_log.Success(c, "告警接入", "更新", fmt.Sprintf("接入IP:%s", query.AccessIP))

	core.OK(c, true)
}

// DeleteAlertAccess
// @Summary 删除告警接入
// @Tags    alert
// @Accept  json
// @Produce json
// @Param   id path     string true "参数"
// @Success 200   {object} core.Response{data=bool}
// @Security JWT
// @Router  /alert/access/{id} [delete]
func (*alertController) DeleteAlertAccess(c *gin.Context) {
	rawID := c.Param("id")
	id, err := strconv.Atoi(rawID)
	if err != nil {
		logger.Error(err)
		c.Error(core.NewParameterError(err.Error()))
		return
	}

	var existAlertAccess models.AlertAccess
	if res := db.Client.Model(&models.AlertAccess{}).First(&existAlertAccess, "id = ?", id); res.Error != nil {
		logger.Error(res.Error)
		c.Error(errs.ErrDeleteAlertAccess)
		audit_log.Fail(c, "告警接入", "删除", fmt.Sprintf("ID不存在:%d", id))
		return
	}

	if res := db.Client.Delete(&models.AlertAccess{Model: core.Model{ID: id}}); res.Error != nil {
		logger.Error(res.Error)
		c.Error(errs.ErrDeleteAlertAccess)
		return
	}

	audit_log.Success(c, "告警接入", "删除", fmt.Sprintf("接入IP:%s", existAlertAccess.AccessIP))

	core.OK(c, true)
}

func (ac *alertController) Setup(r *gin.RouterGroup) {
	r.GET("/alerts", ac.QueryAlerts).
		PUT("/alert", ac.UpdateAlert).
		GET("/alert/config", ac.GetAlertConfig).
		PUT("/alert/config", ac.UpdateAlertConfig).
		GET("/alert/accesses", ac.QueryAlertAccesses).
		POST("/alert/access", ac.CreateAlertAccess).
		PUT("/alert/access", ac.UpdateAlertAccess).
		DELETE("/alert/access/:id", ac.DeleteAlertAccess)
}
