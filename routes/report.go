package routes

import (
	"fmt"
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-module/carbon/v2"
	jsoniter "github.com/json-iterator/go"
	"github.com/yoyo-inc/yoyo/common/db"
	"github.com/yoyo-inc/yoyo/common/logger"
	"github.com/yoyo-inc/yoyo/core"
	"github.com/yoyo-inc/yoyo/errs"
	"github.com/yoyo-inc/yoyo/models"
	"github.com/yoyo-inc/yoyo/resources"
	"github.com/yoyo-inc/yoyo/services"
	"github.com/yoyo-inc/yoyo/services/audit_log"
	"github.com/yoyo-inc/yoyo/vo"
	"gorm.io/gorm"
)

type reportController struct{}

// QueryReports
// @Summary  查询报告列表
// @Tags     report
// @Accept   json
// @Produce  json
// @Param    query query    models.Pagination false "参数"
// @Success  200 {object} core.Response{data=core.PaginatedData{data=[]models.Report}}
// @Security JWT
// @Router   /reports [get]
func (*reportController) QueryReports(c *gin.Context) {
	var query vo.QueryReportVO
	if err := c.ShouldBindQuery(&query); err != nil {
		logger.Error(err)
		c.Error(core.NewParameterError(err))
		return
	}
	queries := core.GetPaginatedQuery(&models.Report{})

	for i := range queries {
		if query.Filename != "" {
			queries[i].Where("filename like ?", "%"+query.Filename+"%")
			query.Filename = ""
		}
	}

	var reports []models.Report
	if res := queries[0].Preload("Resource").Scopes(core.Paginator(c), core.DateTimeRanger(c, ""), core.Orderer(c)).Where(&query).Find(&reports); res.Error != nil {
		logger.Error(res.Error)
		c.Error(errs.ErrQueryReport)
		return
	}

	var count int64
	if res := queries[1].Where(&query).Count(&count); res.Error != nil {
		logger.Error(res.Error)
		c.Error(errs.ErrQueryReport)
		return
	}

	core.OK(c, core.Paginated(reports, count))
}

// DeleteReport
// @Summary  删除报告
// @Tags     report
// @Accept   json
// @Produce  json
// @Param    id  path     string true "参数"
// @Success  200 {object} core.Response{data=bool}
// @Security JWT
// @Router   /report/:id [delete]
func (*reportController) DeleteReport(c *gin.Context) {
	id := c.Param("id")

	if err := services.DeleteResourceFile(c, id); err != nil {
		logger.Error(err)
		c.Error(err)
		return
	}

	if res := db.Client.Delete(&models.Report{}, "id = ?", id); res.Error != nil {
		logger.Error(res.Error)
		c.Error(errs.ErrDeleteReport)
		return
	}

	core.OK(c, true)
}

func (*reportController) PreviewReport(c *gin.Context) {
	file := c.Param("filepath")
	reportType := c.Param("reportType")
	if file == "/" || file == "index.html" {
		var query vo.GenerateReportVO
		if err := c.ShouldBindQuery(&query); err != nil {
			logger.Error(err)
			c.Error(core.NewParameterError(err))
			return
		}

		page, err := services.RenderReport(services.GenerateReportOption{
			ReportType: reportType,
			StartTime:  query.StartTime,
			EndTime:    query.EndTime,
		})
		if err != nil {
			logger.Error(err)
			c.Error(errs.ErrPreviewReport)
			return
		}

		c.Header("Content-Type", "text/html")
		_, _ = c.Writer.Write(page)
	} else {
		fsys, err := fs.Sub(resources.InternalReportTplDir, "report")
		if err != nil {
			logger.Error(err)
			c.Error(errs.ErrPreviewReport)
			return
		}
		prefix := strings.Replace(c.Request.URL.Path, filepath.Join(reportType, file), "", 1)
		fileServer := http.StripPrefix(prefix, http.FileServer(http.FS(fsys)))
		fileServer.ServeHTTP(c.Writer, c.Request)
	}
}

// GenerateReport
// @Summary  生成报告
// @Tags     report
// @Accept   json
// @Produce  json
// @Param    query      body     vo.GenerateReportVO ture "参数"
// @Param    reportType path     string              true "报告类型"
// @Success  200        {object} core.Response{data=bool}
// @Security JWT
// @Router   /report/generate/:reportType [post]
func (*reportController) GenerateReport(c *gin.Context) {
	var query vo.GenerateReportVO
	if err := c.ShouldBindJSON(&query); err != nil {
		logger.Error(err)
		c.Error(core.NewParameterError(err))
		return
	}

	reportType := c.Param("reportType")

	go func() {
		if err := services.GenerateReport(services.GenerateReportOption{
			ReportType: reportType,
			ReportName: query.ReportName,
			StartTime:  query.StartTime,
			EndTime:    query.EndTime,
		}); err != nil {
			logger.Error(err)
		}
	}()

	audit_log.Success(c, "报告", "任务创建", fmt.Sprintf("任务类型：%s,时间范围：%s至%s", reportType, query.StartTime, query.EndTime))
	core.OK(c, true)
}

// QueryReportType
// @Summary  查询报告类型
// @Tags     report
// @Accept   json
// @Produce  json
// @Success  200 {object} core.Response{data=[]models.Dict}
// @Security JWT
// @Router   /report/types [get]
func (rc *reportController) QueryReportType(c *gin.Context) {
	types := []models.Dict{
		{
			Label: "默认",
			Value: "default",
		},
	}

	core.OK(c, types)
}

// QueryReportConfig
// @Summary 查询报告设置
// @Tags    report
// @Accept  json
// @Produce json
// @Success 200   {object} core.Response{data=models.ReportConfig}
// @Security JWT
// @Router  /report/configs [get]
func (*reportController) QueryReportConfig(c *gin.Context) {
	var config models.ReportConfig
	if res := db.Client.Model(&models.ReportConfig{}).Find(&config); res.Error != nil {
		logger.Error(res.Error)
		c.Error(errs.ErrQueryReportConfig)
		return
	}

	core.OK(c, config)
}

// UpdateReportConfig
// @Summary 更换报告设置
// @Tags    report
// @Accept  json
// @Produce json
// @Param   body body     vo.UpdateReportConfigVo true "参数"
// @Success 200   {object} core.Response{data=bool}
// @Security JWT
// @Router  /report/config [put]
func (*reportController) UpdateReportConfig(c *gin.Context) {
	var query vo.UpdateReportConfigVo
	if err := c.ShouldBindJSON(&query); err != nil {
		logger.Error(err)
		c.Error(core.NewParameterError(err))
		return
	}

	json := jsoniter.ConfigCompatibleWithStandardLibrary
	if period, err := json.Marshal(query.Period); err != nil {
		logger.Error(err)
		c.Error(errs.ErrUpdateReportConfig)
		return
	} else {
		query.ReportConfig.Period = period
	}

	err := db.Client.Transaction(func(tx *gorm.DB) error {
		if query.ID != "" {
			if res := db.Client.Delete(&models.ReportConfig{Model: core.Model{ID: query.ID}}); res.Error != nil {
				return res.Error
			}
		}
		if res := db.Client.Create(&query.ReportConfig); res.Error != nil {
			return res.Error
		}

		if len(query.Period) > 0 && len(query.ReportType) > 0 {
			for _, v := range query.Period {
				var spec string
				var description string
				switch v {
				case "day":
					spec = "0 0 * * *"
				case "week":
					spec = "0 0 * * 0"
				case "month":
					spec = "0 0 1 * *"
				}

				for _, t := range query.ReportType {
					jobID := t + v
					services.AddSchedJob(jobID, "report", description, spec, func() error {
						var startTime, endTime string
						switch v {
						case "day":
							startTime = carbon.Now().StartOfDay().ToDateTimeString()
							endTime = carbon.Now().EndOfDay().ToDateTimeString()
						case "week":
							startTime = carbon.Now().StartOfWeek().ToDateTimeString()
							endTime = carbon.Now().EndOfWeek().ToDateTimeString()
						case "month":
							startTime = carbon.Now().StartOfMonth().ToDateTimeString()
							endTime = carbon.Now().EndOfMonth().ToDateTimeString()
						}
						return services.GenerateReport(services.GenerateReportOption{
							ReportType: t,
							StartTime:  startTime,
							EndTime:    endTime,
						})
					})
				}
			}
		}

		return nil
	})
	if err != nil {
		logger.Error(err)
		c.Error(errs.ErrUpdateReportConfig)
		return
	}

	core.OK(c, true)
}

func AddSchedReportJob(period []string, reportName string) {}

func ReportDefaultCallback(startTime string, endTime string) services.ReportCallbackData {
	data := make(services.ReportCallbackData)

	return data
}

func (rc *reportController) Setup(r *gin.RouterGroup) {
	r.GET("/reports", rc.QueryReports).
		DELETE("/report/:id", rc.DeleteReport).
		GET("/report/types", rc.QueryReportType).
		GET("/report/preview/:reportType/*filepath", rc.PreviewReport).
		POST("/report/generate/:reportType", rc.GenerateReport).
		PUT("/report/config", rc.UpdateReportConfig)

	services.RegisterReportCallback("default", ReportDefaultCallback)
}
