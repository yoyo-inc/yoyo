package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/yoyo-inc/yoyo/common/db"
	"github.com/yoyo-inc/yoyo/common/logger"
	"github.com/yoyo-inc/yoyo/core"
	"github.com/yoyo-inc/yoyo/errs"
	"github.com/yoyo-inc/yoyo/models"
	"github.com/yoyo-inc/yoyo/vo"
)

type alertController struct{}

func init() {
	AddService(&alertController{})
}

// QueryAlerts
// @Summary  查询告警列表
// @Tags     alert
// @Produce  json
// @Param    query  query     vo.QueryAlertVO    false  "参数"
// @Param    query  query     models.Pagination  false  "参数"
// @Success  200    {object}  core.Response{data=core.PaginatedData{list=[]models.Alert}}
// @Router   /alerts [get]
func (*alertController) QueryAlerts(c *gin.Context) {
	var query vo.QueryAlertVO
	if err := c.ShouldBindQuery(&query); err != nil {
		logger.Error(err)
		c.Error(core.NewParameterError(err.Error()))
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
// @Summary  更新告警
// @Tags     alert
// @Accept   json
// @Produce  json
// @Param    query  body      vo.UpdateAlertVO  true  "告警信息"
// @Success  200    {object}  core.Response{data=bool}
// @Router   /alert [put]
func (*alertController) UpdateAlert(c *gin.Context) {
	var query vo.UpdateAlertVO
	if err := c.ShouldBindJSON(&query); err != nil {
		logger.Error(err)
		c.Error(core.NewParameterError(err.Error()))
		return
	}

	if res := db.Client.Updates(&query); res.Error != nil {
		logger.Error(res.Error)
		c.Error(errs.ErrUpdateAlert)
		return
	}

	core.OK(c, true)
}

func (ac *alertController) Setup(r *gin.RouterGroup) {
	r.GET("/alerts", ac.QueryAlerts).PUT("/alert", ac.UpdateAlert)
}
