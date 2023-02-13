package errs

import "github.com/yoyo-inc/yoyo/core"

var (
	ErrQueryAlert                 = core.NewBusinessError("500", "查询告警失败")
	ErrCreateAlert                = core.NewBusinessError("501", "创建告警失败")
	ErrUpdateAlert                = core.NewBusinessError("502", "更新告警失败")
	ErrQueryAlertConfig           = core.NewBusinessError("503", "查询告警配置失败")
	ErrUpdateAlertConfig          = core.NewBusinessError("504", "更新告警配置失败")
	ErrQueryAlertAccess           = core.NewBusinessError("505", "查询告警接入失败")
	ErrCreateAlertAccess          = core.NewBusinessError("506", "创建告警接入失败")
	ErrUpdateAlertAccess          = core.NewBusinessError("507", "更新告警接入失败")
	ErrDeleteAlertAccess          = core.NewBusinessError("508", "删除告警接入失败")
	ErrResolveAlert               = core.NewBusinessError("509", "处置告警失败")
	ErrReadAlert                  = core.NewBusinessError("510", "已读告警失败")
	ErrIgnoreAlert                = core.NewBusinessError("511", "忽略告警失败")
	ErrReceiveAlertmanagerMessage = core.NewBusinessError("512", "接受信息失败")
	ErrQueryAlertCount            = core.NewBusinessError("513", "查询告警总数失败")
	ErrQueryAlertPush             = core.NewBusinessError("514", "查询告警推送配置失败")
	ErrCreateAlertPush            = core.NewBusinessError("515", "创建告警推送配置失败")
	ErrUpdateAlertPush            = core.NewBusinessError("516", "更新告警推送配置失败")
	ErrDeleteAlertPush            = core.NewBusinessError("517", "删除告警推送配置失败")
	ErrQueryAlertTypes            = core.NewBusinessError("518", "查询告警类型失败")
)
