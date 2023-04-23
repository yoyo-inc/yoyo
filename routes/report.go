package routes

import (
	"errors"
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
//
//	@Summary	查询报告列表
//	@Tags		report
//	@Accept		json
//	@Produce	json
//	@Param		query	query		models.Pagination	true	"参数"
//	@Param		query	query		vo.QueryReportVO	true	"参数"
//	@Success	200		{object}	core.Response{data=core.PaginatedData{data=[]models.Report}}
//	@Security	JWT
//	@Router		/reports [get]
func (*reportController) QueryReports(c *gin.Context) {
	var query vo.QueryReportVO
	if err := c.ShouldBindQuery(&query); err != nil {
		logger.Error(err)
		_ = c.Error(core.NewParameterError(err))
		return
	}
	queries := core.GetPaginatedQuery(&models.Report{})

	for i := range queries {
		if query.ReportName != "" {
			queries[i].Where("report_name like ?", "%"+query.ReportName+"%")
			query.ReportName = ""
		}
	}

	var reports []models.Report
	if res := queries[0].Preload("Resource").Scopes(core.Paginator(c), core.DateTimeRanger(c, ""), core.Orderer(c)).Where(&query).Find(&reports); res.Error != nil {
		logger.Error(res.Error)
		_ = c.Error(errs.ErrQueryReport)
		return
	}

	var count int64
	if res := queries[1].Where(&query).Count(&count); res.Error != nil {
		logger.Error(res.Error)
		_ = c.Error(errs.ErrQueryReport)
		return
	}

	core.OK(c, core.Paginated(reports, count))
}

// DeleteReport
//
//	@Summary	删除报告
//	@Tags		report
//	@Accept		json
//	@Produce	json
//	@Param		id	path		string	true	"参数"
//	@Success	200	{object}	core.Response{data=bool}
//	@Security	JWT
//	@Router		/report/:id [delete]
func (*reportController) DeleteReport(c *gin.Context) {
	id := c.Param("id")

	if err := services.DeleteResourceFile(c, id); err != nil {
		logger.Error(err)
		_ = c.Error(err)
		return
	}

	if res := db.Client.Delete(&models.Report{}, "id = ?", id); res.Error != nil {
		logger.Error(res.Error)
		_ = c.Error(errs.ErrDeleteReport)
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
			_ = c.Error(core.NewParameterError(err))
			return
		}

		page, err := services.RenderReport(services.GenerateReportOption{
			ReportType: reportType,
			StartTime:  query.StartTime,
			EndTime:    query.EndTime,
		})
		if err != nil {
			logger.Error(err)
			_ = c.Error(errs.ErrPreviewReport)
			return
		}

		c.Header("Content-Type", "text/html")
		_, _ = c.Writer.Write(page)
	} else {
		fsys, err := fs.Sub(resources.InternalReportTplDir, "report")
		if err != nil {
			logger.Error(err)
			_ = c.Error(errs.ErrPreviewReport)
			return
		}
		prefix := strings.Replace(c.Request.URL.Path, filepath.Join(reportType, file), "", 1)
		fileServer := http.StripPrefix(prefix, http.FileServer(http.FS(fsys)))
		fileServer.ServeHTTP(c.Writer, c.Request)
	}
}

// GenerateReport
//
//	@Summary	生成报告
//	@Tags		report
//	@Accept		json
//	@Produce	json
//	@Param		query		body		vo.GenerateReportVO	ture	"参数"
//	@Param		reportType	path		string				true	"报告类型"
//	@Success	200			{object}	core.Response{data=bool}
//	@Security	JWT
//	@Router		/report/generate/:reportType [post]
func (*reportController) GenerateReport(c *gin.Context) {
	var query vo.GenerateReportVO
	if err := c.ShouldBindJSON(&query); err != nil {
		logger.Error(err)
		_ = c.Error(core.NewParameterError(err))
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
//
//	@Summary	查询报告类型
//	@Tags		report
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	core.Response{data=[]models.Dict}
//	@Security	JWT
//	@Router		/report/types [get]
func (rc *reportController) QueryReportType(c *gin.Context) {
	entries, err := services.GetEntriesByType("report")
	if err != nil {
		logger.Error(err)
		_ = c.Error(errs.ErrQueryReportType)
		return
	}
	core.OK(c, entries)
}

// QueryReportConfig
//
//	@Summary	查询报告设置
//	@Tags		report
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	core.Response{data=models.ReportConfig}
//	@Security	JWT
//	@Router		/report/config [get]
func (*reportController) QueryReportConfig(c *gin.Context) {
	var config models.ReportConfig
	if res := db.Client.Model(&models.ReportConfig{}).First(&config); res.Error != nil {
		logger.Error(res.Error)
		_ = c.Error(errs.ErrQueryReportConfig)
		return
	}

	core.OK(c, config)
}

// UpdateReportConfig
//
//	@Summary	更换报告设置
//	@Tags		report
//	@Accept		json
//	@Produce	json
//	@Param		body	body		vo.UpdateReportConfigVo	true	"参数"
//	@Success	200		{object}	core.Response{data=bool}
//	@Security	JWT
//	@Router		/report/config [put]
func (rc *reportController) UpdateReportConfig(c *gin.Context) {
	var query vo.UpdateReportConfigVo
	if err := c.ShouldBindJSON(&query); err != nil {
		logger.Error(err)
		_ = c.Error(core.NewParameterError(err))
		return
	}

	json := jsoniter.ConfigCompatibleWithStandardLibrary
	if query.Period != nil {
		if period, err := json.Marshal(query.Period); err != nil {
			logger.Error(err)
			_ = c.Error(errs.ErrUpdateReportConfig)
			return
		} else {
			query.ReportConfig.Period = period
		}
	}

	if query.ReportType != nil {
		if reportType, err := json.Marshal(query.ReportType); err != nil {
			logger.Error(err)
			_ = c.Error(errs.ErrUpdateReportConfig)
			return
		} else {
			query.ReportConfig.ReportType = reportType
		}
	}

	err := db.Client.Transaction(func(tx *gorm.DB) error {
		if err := services.RemoveSchedJobByType("report"); err != nil {
			return err
		}
		if res := db.Client.Where("1 = 1").Delete(&models.ReportConfig{}); res.Error != nil {
			return res.Error
		}
		if res := db.Client.Create(&query.ReportConfig); res.Error != nil {
			return res.Error
		}

		if err := rc.registerReportSchedJob(); err != nil {
			return err
		}

		return nil
	})

	cs, _ := json.MarshalToString(query.ReportConfig)

	if err != nil {
		logger.Error(err)
		_ = c.Error(errs.ErrUpdateReportConfig)
		audit_log.Fail(c, "报告设置", "更新", fmt.Sprintf("报告设置内容%s", cs))
		return
	}

	audit_log.Success(c, "报告设置", "更新", fmt.Sprintf("报告设置内容%s", cs))

	core.OK(c, true)
}

func (*reportController) registerReportSchedJob() error {
	var config models.ReportConfig
	if res := db.Client.Model(&models.ReportConfig{}).First(&config); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil
		}
		return res.Error
	}

	if len(config.GetPeriod()) > 0 && len(config.GetReportType()) > 0 {
		for _, v := range config.GetPeriod() {
			var spec string
			var description string
			switch v {
			case "day":
				spec = "0 22 * * *"
				description = "每日报告"
			case "week":
				spec = "0 22 * * 0"
				description = "每周报告"
			case "month":
				spec = "0 22 31 * *"
				description = "每月报告"
			case "year":
				spec = "0 22 31 12 *"
				description = "每年报告"
			}

			for _, t := range config.GetReportType() {
				jobID := t + v
				// remove exists sched job
				if err := services.RemoveSchedJob(jobID); err != nil {
					return err
				}
				if err := services.AddSchedJob(jobID, "report", fmt.Sprintf("%s%s", services.GetLabelByValue("report", t), description), spec, func() error {
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
					case "year":
						startTime = carbon.Now().StartOfYear().ToDateTimeString()
						endTime = carbon.Now().EndOfYear().ToDateTimeString()
					}
					return services.GenerateReport(services.GenerateReportOption{
						ReportType: t,
						StartTime:  startTime,
						EndTime:    endTime,
					})
				}); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func ReportDefaultCallback(startTime string, endTime string) (services.ReportCallbackData, error) {
	data := make(services.ReportCallbackData)

	return data, nil
}

func (rc *reportController) Setup(r *gin.RouterGroup) {
	r.GET("/reports", rc.QueryReports).
		DELETE("/report/:id", rc.DeleteReport).
		GET("/report/types", rc.QueryReportType).
		GET("/report/preview/:reportType/*filepath", rc.PreviewReport).
		POST("/report/generate/:reportType", rc.GenerateReport).
		GET("/report/config", rc.QueryReportConfig).
		PUT("/report/config", rc.UpdateReportConfig)

	services.RegisterReportCallback("default", ReportDefaultCallback)

	if err := rc.registerReportSchedJob(); err != nil {
		logger.Error(err)
	}
}
