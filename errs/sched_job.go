package errs

import "github.com/yoyo-inc/yoyo/core"

var (
	ErrQuerySchedJob      = core.NewBusinessError("1000", "查询定时任务失败")
	ErrStopSchedJob       = core.NewBusinessError("1001", "关闭定时任务失败")
	ErrStartSchedJob      = core.NewBusinessError("1002", "开启定时任务失败")
	ErrNotExistSchedJob   = core.NewBusinessError("1003", "定时任务不存在")
	ErrQuerySchedJobTypes = core.NewBusinessError("1004", "查询定时任务类型失败")
)
