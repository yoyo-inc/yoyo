package errors

import "github.com/ypli0629/yoyo/core"

var (
	ErrInvalidParameters = core.NewBusinessError("001", "无效的参数")
)
