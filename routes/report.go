package routes

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yoyo-inc/yoyo/common/db"
	"github.com/yoyo-inc/yoyo/common/logger"
	"github.com/yoyo-inc/yoyo/core"
	"github.com/yoyo-inc/yoyo/errs"
	"github.com/yoyo-inc/yoyo/models"
)

type reportController struct{}

// QueryReports
// @Summary 查询报告列表
// @Tags    report
// @Accept  json
// @Produce json
// @Success 200   {object} core.Response{data=core.PaginatedData{data=[]models.Report}}
// @Security JWT
// @Router  /reports [get]
func (*reportController) QueryReports(c *gin.Context) {
	queries := core.GetPaginatedQuery(&models.Report{})

	var reports []models.Report
	if res := queries[0].Find(&reports); res.Error != nil {
		logger.Error(res.Error)
		c.Error(errs.ErrQueryReport)
		return
	}

	var count int64
	if res := queries[1].Count(&count); res.Error != nil {
		logger.Error(res.Error)
		c.Error(errs.ErrQueryReport)
		return
	}

	core.OK(c, core.Paginated(reports, count))
}

// DeleteReport
// @Summary 删除报告
// @Tags    report
// @Accept  json
// @Produce json
// @Param   id path     string true "参数"
// @Success 200   {object} core.Response{data=bool}
// @Security JWT
// @Router  /report [delete]
func (*reportController) DeleteReport(c *gin.Context) {
	rawID := c.Param("id")
	id, err := strconv.Atoi(rawID)
	if err != nil {
		logger.Error(err)
		c.Error(errs.ErrInvalidParameters)
		return
	}

	if res := db.Client.Delete(&models.Report{IModel: core.IModel{ID: id}}); res.Error != nil {
		logger.Error(res.Error)
		c.Error(errs.ErrDeleteReport)
		return
	}

	core.OK(c, true)
}

func (rc *reportController) Setup(r *gin.RouterGroup) {
	r.GET("/reports", rc.QueryReports).DELETE("/report", rc.DeleteReport)
}
