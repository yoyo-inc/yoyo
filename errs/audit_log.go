package errs

import "github.com/yoyo-inc/yoyo/core"

var (
	ErrQueryAuditLog       = core.NewBusinessError("300", "查询审计日志失败")
	ErrQueryAuditLogModule = core.NewBusinessError("301", "查询审计日志模块失败")
)
