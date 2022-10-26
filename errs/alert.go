package errs

import "github.com/yoyo-inc/yoyo/core"

var (
	ErrQueryAlert  = core.NewBusinessError("500", "查询告警失败")
	ErrCreateAlert = core.NewBusinessError("501", "创建告警失败")
	ErrUpdateAlert = core.NewBusinessError("502", "更新告警失败")
)
