package routes

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/gin-gonic/gin"
	"github.com/golang-module/carbon/v2"
	jsoniter "github.com/json-iterator/go"
	"github.com/prometheus/alertmanager/template"
	"github.com/yoyo-inc/yoyo/common/db"
	"github.com/yoyo-inc/yoyo/common/dt"
	"github.com/yoyo-inc/yoyo/common/hub"
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

// QueryAlerts
//
//	@Summary	查询告警列表
//	@Tags		alert
//	@Produce	json
//	@Param		query	query		vo.QueryAlertVO			false	"参数"
//	@Param		query	query		models.Pagination		false	"参数"
//	@Param		query	query		models.DateTimeRanger	false	"参数"
//	@Success	200		{object}	core.Response{data=core.PaginatedData{list=[]models.Alert}}
//	@Security	JWT
//	@Router		/alerts [get]
func (*alertController) QueryAlerts(c *gin.Context) {
	var query vo.QueryAlertVO
	if err := c.ShouldBindQuery(&query); err != nil {
		logger.Error(err)
		c.Error(core.NewParameterError(err))
		return
	}
	queries := core.GetPaginatedQuery(&models.Alert{})

	var alerts []models.Alert
	if res := queries[0].Scopes(core.Paginator(c), core.DateTimeRanger(c, "start_at")).
		Where(&query).Order("start_at desc").Find(&alerts); res.Error != nil {
		logger.Error(res.Error)
		c.Error(errs.ErrQueryAlert)
		return
	}

	var count int64
	if res := queries[1].Scopes(core.DateTimeRanger(c, "start_at")).Where(&query).Count(&count); res.Error != nil {
		logger.Error(res.Error)
		c.Error(errs.ErrQueryAlert)
		return
	}

	core.OK(c, core.Paginated(alerts, count))
}

// QueryAlertTypes
//
//	@Summary	查询告警类型
//	@Tags		alert
//	@Accept		json
//	@Produce	json
//	@Param		query	query		vo.QueryAlertTypeVO	true	"类型"
//	@Success	200		{object}	core.Response{data=array,vo.Record}
//	@Security	JWT
//	@Router		/alert/types [get]
func (*alertController) QueryAlertTypes(c *gin.Context) {
	var query vo.QueryAlertTypeVO
	if err := c.ShouldBindQuery(&query); err != nil {
		logger.Error(err)
		c.Error(core.NewParameterError(err))
		return
	}

	alertTypes, err := services.GetEntriesByType("alertType")
	if err != nil {
		logger.Error(err)
		c.Error(errs.ErrQueryAlertTypes)
		return
	}

	if query.Type == 1 {
		alertTypes = slice.Map(alertTypes, func(index int, item services.Entry) services.Entry {
			item.Value = item.Label
			return item
		})
		rows, err := db.Client.Raw("select distinct(type) as type from alerts").Rows()
		if err != nil {
			logger.Error(err)
			c.Error(errs.ErrQueryAlertTypes)
			return
		}

		for rows.Next() {
			var t string
			err := rows.Scan(&t)
			if err != nil {
				logger.Error(err)
				c.Error(errs.ErrQueryAlertTypes)
				return
			}
			if t != "" && !slice.Some(alertTypes, func(index int, item services.Entry) bool {
				return item.Label == t
			}) {
				alertTypes = append(alertTypes, services.Entry{Label: t, Value: t})
			}
		}
	}

	core.OK(c, alertTypes)
}

// ResolvedAlert
//
//	@Summary	处置告警
//	@Tags		alert
//	@Accept		json
//	@Produce	json
//	@Param		body	body		vo.ResolveAlertVO	true	"参数"
//	@Success	200		{object}	core.Response{data=bool}
//	@Security	JWT
//	@Router		/alert/resolve [put]
func (*alertController) ResolvedAlert(c *gin.Context) {
	var query vo.ResolveAlertVO
	if err := c.ShouldBindJSON(&query); err != nil {
		logger.Error(err)
		c.Error(core.NewParameterError(err))
		return
	}

	if res := db.Client.Model(&models.Alert{}).Where("id = ?", query.ID).Updates(map[string]interface{}{
		"remark":          query.Remark,
		"resolved_status": 1,
		"status":          1,
	}); res.Error != nil {
		logger.Error(res.Error)
		c.Error(errs.ErrResolveAlert)
		return
	}

	core.OK(c, true)
}

// IgnoreAlert
//
//	@Summary	更新告警
//	@Tags		alert
//	@Accept		json
//	@Produce	json
//	@Param		body	body		vo.IgnoreAlertVO	true	"参数"
//	@Success	200		{object}	core.Response{data=bool}
//	@Security	JWT
//	@Router		/alert/ignore [put]
func (*alertController) IgnoreAlert(c *gin.Context) {
	var query vo.IgnoreAlertVO
	if err := c.ShouldBindJSON(&query); err != nil {
		logger.Error(err)
		c.Error(core.NewParameterError(err))
		return
	}

	if res := db.Client.Model(&models.Alert{}).Where("id = ?", query.ID).Updates(map[string]interface{}{
		"resolved_status": 3,
		"status":          1,
	}); res.Error != nil {
		logger.Error(res.Error)
		c.Error(errs.ErrIgnoreAlert)
		return
	}

	core.OK(c, true)
}

// GetAlertConfig
//
//	@Summary	查询告警配置
//	@Tags		alert
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	core.Response{data=models.AlertConfig}
//	@Security	JWT
//	@Router		/alert/config [get]
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
//
//	@Summary	更新告警配置
//	@Tags		alert
//	@Accept		json
//	@Produce	json
//	@Param		body	body		models.AlertConfig	true	"参数"
//	@Success	200		{object}	core.Response{data=bool}
//	@Security	JWT
//	@Router		/alert/config [put]
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
//
//	@Summary	查询告警接入
//	@Tags		alert
//	@Accept		json
//	@Produce	json
//	@Param		query	query		vo.QueryAlertAccessVO	false	"参数"
//	@Param		query	query		models.Pagination		false	"参数"
//	@Param		query	query		models.DateTimeRanger	false	"参数"
//	@Success	200		{object}	core.Response{data=core.PaginatedData{data=models.AlertAccess}}
//	@Security	JWT
//	@Router		/alert/accesses [get]
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
	if res := queries[0].Scopes(core.Paginator(c), core.DateTimeRanger(c, "")).Find(&alertAccesses); res.Error != nil {
		logger.Error(res.Error)
		c.Error(errs.ErrQueryAlertAccess)
		return
	}
	var count int64
	if res := queries[1].Scopes(core.DateTimeRanger(c, "")).Count(&count); res.Error != nil {
		logger.Error(res.Error)
		c.Error(errs.ErrQueryAlertAccess)
		return
	}

	core.OK(c, core.Paginated(alertAccesses, count))
}

// CreateAlertAccess
//
//	@Summary	创建告警接入
//	@Tags		alert
//	@Accept		json
//	@Produce	json
//	@Param		body	body		models.AlertAccess	true	"参数"
//	@Success	200		{object}	core.Response{data=bool}
//	@Security	JWT
//	@Router		/alert/access [post]
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
//
//	@Summary	更新告警接入
//	@Tags		alert
//	@Accept		json
//	@Produce	json
//	@Param		body	body		vo.UpdateAlertAccessVO	true	"参数"
//	@Success	200		{object}	core.Response{data=bool}
//	@Security	JWT
//	@Router		/alert/access [put]
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
//
//	@Summary	删除告警接入
//	@Tags		alert
//	@Accept		json
//	@Produce	json
//	@Param		id	path		string	true	"参数"
//	@Success	200	{object}	core.Response{data=bool}
//	@Security	JWT
//	@Router		/alert/access/{id} [delete]
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

	if res := db.Client.Delete(&models.AlertAccess{IModel: core.IModel{ID: id}}); res.Error != nil {
		logger.Error(res.Error)
		c.Error(errs.ErrDeleteAlertAccess)
		return
	}

	audit_log.Success(c, "告警接入", "删除", fmt.Sprintf("接入IP:%s", existAlertAccess.AccessIP))

	core.OK(c, true)
}

func (*alertController) AccessAlert(c *gin.Context) {
	var accesses []models.AlertAccess
	if res := db.Client.Model(&models.AlertAccess{}).Find(&accesses); res.Error != nil {
		logger.Error(res.Error)
		c.Error(errs.ErrQueryAlertAccess)
		return
	}

	var remoteIP string
	if ip := c.GetHeader("X-Real-IP"); ip != "" {
		remoteIP = ip
	} else {
		if proxy := c.GetHeader("X-Forwarded-For"); proxy != "" {
			ips := strings.Split(proxy, ",")
			if len(ips) > 0 {
				remoteIP = strings.TrimSpace(ips[0])
			}
		} else {
			remoteIP = c.RemoteIP()
		}
	}

	if access, ok := slice.Find(accesses, func(_ int, item models.AlertAccess) bool {
		return item.AccessIP == remoteIP
	}); ok {
		// has right
		var body interface{}
		json := jsoniter.ConfigCompatibleWithStandardLibrary
		if err := json.NewDecoder(c.Request.Body).Decode(&body); err != nil {
			logger.Error(err)
			c.Error(core.NewParameterError("参数有误"))
			return
		}

		assertStr := func(value reflect.Value, field string) (string, bool) {
			var s string
			if v := value.MapIndex(reflect.ValueOf(field)); v.IsValid() {
				s = (v.Interface()).(string)
				return s, true
			}
			return s, false
		}

		transform := func(value reflect.Value) error {
			var alert models.Alert
			if v, ok := assertStr(value, access.StartAtField); ok {
				t := carbon.Parse(v).Carbon2Time()
				alert.StartAt = (*dt.LocalTime)(&t)
			}
			if v, ok := assertStr(value, access.TypeField); ok {
				if v == "" {
					alert.Type = "其他"
				} else {
					alert.Type = v
				}
			}
			if v, ok := assertStr(value, access.LevelField); ok {
				alert.Level = v
			}
			if v, ok := assertStr(value, access.ContentField); ok {
				alert.Content = v
			}
			alert.From = access.AccessIP

			if res := db.Client.Create(&alert); res.Error != nil {
				logger.Error(res.Error)
				return res.Error
			}

			return nil
		}

		switch reflect.TypeOf(body).Kind() {
		case reflect.Slice:
			value := reflect.ValueOf(body)
			for i := 0; i < value.Len(); i++ {
				transform(value.Index(i))
			}
		case reflect.Map:
			transform(reflect.ValueOf(body))
		}
	} else {
		// have not right
		c.AbortWithStatus(http.StatusForbidden)
	}
}

// Webhook webhook
func (*alertController) Webhook(c *gin.Context) {
	var message template.Data
	if err := c.ShouldBindJSON(&message); err != nil {
		logger.Error(err)
		c.Error(errs.ErrReceiveAlertmanagerMessage)
		return
	}

	var pushes []models.AlertPush
	if res := db.Client.Model(&models.AlertPush{}).Find(&pushes); res.Error != nil {
		logger.Error(res.Error)
	}

	json := jsoniter.ConfigCompatibleWithStandardLibrary
	for _, alert := range message.Alerts {
		var modelAlert models.Alert
		modelAlert.From = "本系统"
		modelAlert.Level = alert.Labels["severity"]
		modelAlert.Type = services.GetLabelByValue("alertType", alert.Labels["group"])
		modelAlert.Content = alert.Annotations["summary"]
		modelAlert.StartAt = (*dt.LocalTime)(&alert.StartsAt)

		if res := db.Client.Create(&modelAlert); res.Error != nil {
			logger.Error(res.Error)
		}

		// broadcast alert to site station
		msg, _ := json.MarshalToString(&modelAlert)
		hub.SendMessage(msg)

		// broadcast alert to third party platform
		if pushes != nil && len(pushes) > 0 {
			for _, push := range pushes {
				if slice.Contain(push.GetAlertTypes(), alert.Labels["group"]) {
					if ok := services.PushAlert(push, msg); !ok {
						// failed to push message
						logger.Errorf("Failed to push %s to %v", msg, push)
					}
				}
			}
		}
	}

	core.OK(c, true)
}

// QueryAlertCount
//
//	@Summary	查询告警数量
//	@Tags		alert
//	@Accept		json
//	@Produce	json
//	@Param		query	query		vo.QueryAlertCountVO	false	"参数"
//	@Success	200		{object}	core.Response{data=boolean}
//	@Security	JWT
//	@Router		/alert/count [get]
func (*alertController) QueryAlertCount(c *gin.Context) {
	var query vo.QueryAlertCountVO
	if err := c.ShouldBindQuery(&query); err != nil {
		logger.Error(err)
		c.Error(core.NewParameterError(err))
		return
	}

	var count int64
	if res := db.Client.Model(&models.Alert{}).Where(&query).Count(&count); res.Error != nil {
		logger.Error(res.Error)
		c.Error(errs.ErrQueryAlertCount)
		return
	}

	core.OK(c, count)
}

// QueryAlertPush
//
//	@Summary	查询告警推送设置
//	@Tags		alert
//	@Accept		json
//	@Produce	json
//	@Param		query	query		models.DateTimeRanger	false	"参数"
//	@Success	200		{object}	core.Response{data=core.PaginatedData{data=models.AlertPush}}
//	@Security	JWT
//	@Router		/alert/push [get]
func (*alertController) QueryAlertPush(c *gin.Context) {
	queries := core.GetPaginatedQuery(&models.AlertPush{})
	var alertPushes []models.AlertPush
	if res := queries[0].Scopes(core.Paginator(c), core.DateTimeRanger(c, "")).Find(&alertPushes); res.Error != nil {
		logger.Error(res.Error)
		c.Error(errs.ErrQueryAlertPush)
		return
	}

	var count int64
	if res := queries[1].Scopes(core.DateTimeRanger(c, "")).Count(&count); res.Error != nil {
		logger.Error(res.Error)
		c.Error(errs.ErrQueryAlertPush)
		return
	}

	core.OK(c, core.Paginated(alertPushes, count))
}

// CreateAlertPush
//
//	@Summary	创建告警推送设置
//	@Tags		alert
//	@Accept		json
//	@Produce	json
//	@Param		body	body		models.AlertPush	true	"参数"
//	@Success	200		{object}	core.Response{data=bool}
//	@Security	JWT
//	@Router		/alert/push [post]
func (*alertController) CreateAlertPush(c *gin.Context) {
	var query models.AlertPush
	if err := c.ShouldBindJSON(&query); err != nil {
		logger.Error(err)
		c.Error(core.NewParameterError(err))
		return
	}

	if res := db.Client.Create(&query); res.Error != nil {
		logger.Error(res.Error)
		c.Error(errs.ErrCreateAlertPush)
		return
	}

	core.OK(c, true)
}

// UpdateAlertPush
//
//	@Summary	更新告警推送
//	@Tags		alert
//	@Accept		json
//	@Produce	json
//	@Param		body	body		models.AlertPush	true	"参数"
//	@Success	200		{object}	core.Response{data=bool}
//	@Security	JWT
//	@Router		/alert/push [put]
func (*alertController) UpdateAlertPush(c *gin.Context) {
	var query models.AlertPush
	if err := c.ShouldBindJSON(&query); err != nil {
		logger.Error(err)
		c.Error(core.NewParameterError(err))
		return
	}

	if res := db.Client.Model(&models.AlertPush{IModel: core.IModel{ID: query.ID}}).Updates(&query); res.Error != nil {
		logger.Error(res.Error)
		c.Error(errs.ErrUpdateAlertPush)
		return
	}

	core.OK(c, true)
}

// DeleteAlertPush
//
//	@Summary	删除告警推送
//	@Tags		alert
//	@Accept		json
//	@Produce	json
//	@Param		id	path		string	true	"参数"
//	@Success	200	{object}	core.Response{data=bool}
//	@Security	JWT
//	@Router		/alert/push/{id} [delete]
func (*alertController) DeleteAlertPush(c *gin.Context) {
	rawID := c.Param("id")
	id, err := strconv.Atoi(rawID)
	if err != nil {
		logger.Error(err)
		c.Error(core.NewParameterError(err.Error()))
		return
	}

	if res := db.Client.Delete(&models.AlertPush{IModel: core.IModel{ID: id}}); res.Error != nil {
		logger.Error(res.Error)
		c.Error(errs.ErrDeleteAlertPush)
		return
	}

	core.OK(c, true)
}

func (ac *alertController) Setup(r *gin.RouterGroup) {
	r.GET("/alerts", ac.QueryAlerts).
		GET("/alert/types", ac.QueryAlertTypes).
		GET("/alert/config", ac.GetAlertConfig).
		PUT("/alert/config", ac.UpdateAlertConfig).
		GET("/alert/accesses", ac.QueryAlertAccesses).
		POST("/alert/access", ac.CreateAlertAccess).
		PUT("/alert/access", ac.UpdateAlertAccess).
		DELETE("/alert/access/:id", ac.DeleteAlertAccess).
		POST("/access/alert", ac.AccessAlert).
		PUT("/alert/resolve", ac.ResolvedAlert).
		PUT("/alert/ignore", ac.IgnoreAlert).
		POST("/alert/webhook", ac.Webhook).
		GET("/alert/count", ac.QueryAlertCount).
		GET("/alert/push", ac.QueryAlertPush).
		POST("/alert/push", ac.CreateAlertPush).
		PUT("/alert/push", ac.UpdateAlertPush).
		DELETE("/alert/push/:id", ac.DeleteAlertPush)
	if err := services.GeneratePrometheusConfig(); err != nil {
		logger.Error(err)
	}

	// generate prometheus and alertmanager rules
	var alertConfig models.AlertConfig
	if res := db.Client.Model(&models.AlertConfig{}).First(&alertConfig); res.Error != nil {
		if !errors.Is(res.Error, gorm.ErrRecordNotFound) {
			logger.Debug(res.Error)
		}
	}
	if err := services.GenerateAlertManagerConfig(alertConfig); err != nil {
		logger.Error(err)
	}
}
