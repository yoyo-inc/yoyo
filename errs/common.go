package errs

import "github.com/yoyo-inc/yoyo/core"

var (
	ErrInvalidParameters = core.NewBusinessError("001", "无效的参数")
)
