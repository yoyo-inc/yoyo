package errs

import "github.com/yoyo-inc/yoyo/core"

var (
	ErrQueryReport  = core.NewBusinessError("900", "查询报告失败")
	ErrDeleteReport = core.NewBusinessError("901", "删除报告失败")
)
