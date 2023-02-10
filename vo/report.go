package vo

import "github.com/yoyo-inc/yoyo/models"

type QueryReportVO struct {
	ReportName   string `form:"reportName"`
	ReportType   string `form:"reportType"`
	ReportStatus string `form:"reportStatus"`
}

type GenerateReportVO struct {
	ReportName string `form:"reportName" json:"reportName"`
	StartTime  string `form:"startTime" json:"startTime" binding:"required"`
	EndTime    string `form:"endTime" json:"endTime" binding:"required"`
}

type UpdateReportConfigVo struct {
	models.ReportConfig
	Period     []string `json:"period"`
	ReportType []string `json:"reportType"`
}
