package errs

import "github.com/yoyo-inc/yoyo/core"

var (
	ErrQueryAlert        = core.NewBusinessError("500", "查询告警失败")
	ErrCreateAlert       = core.NewBusinessError("501", "创建告警失败")
	ErrUpdateAlert       = core.NewBusinessError("502", "更新告警失败")
	ErrQueryAlertConfig  = core.NewBusinessError("503", "查询告警配置失败")
	ErrUpdateAlertConfig = core.NewBusinessError("504", "更新告警配置失败")
	ErrQueryAlertAccess  = core.NewBusinessError("505", "查询告警接入失败")
	ErrCreateAlertAccess = core.NewBusinessError("506", "创建告警接入失败")
	ErrUpdateAlertAccess = core.NewBusinessError("507", "更新告警接入失败")
	ErrDeleteAlertAccess = core.NewBusinessError("508", "删除告警接入失败")
)
