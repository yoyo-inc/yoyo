package errs

import "github.com/yoyo-inc/yoyo/core"

var (
	ErrQueryAuditLog = core.NewBusinessError("300", "查询操作日志失败")
)
