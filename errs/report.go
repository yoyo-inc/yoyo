package errs

import "github.com/yoyo-inc/yoyo/core"

var (
	ErrQueryReport        = core.NewBusinessError("900", "查询报告失败")
	ErrDeleteReport       = core.NewBusinessError("901", "删除报告失败")
	ErrPreviewReport      = core.NewBusinessError("902", "预览报告失败")
	ErrGenerateReport     = core.NewBusinessError("903", "生成报告失败")
	ErrQueryReportConfig  = core.NewBusinessError("904", "查询报告配置失败")
	ErrUpdateReportConfig = core.NewBusinessError("905", "更新报告配置失败")
	ErrQueryReportType    = core.NewBusinessError("906", "查询报表类型")
)
