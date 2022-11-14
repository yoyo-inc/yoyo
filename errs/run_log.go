package errs

import "github.com/yoyo-inc/yoyo/core"

var (
	ErrQueryRunLog = core.NewBusinessError("600", "查询运行日志失败")
)
