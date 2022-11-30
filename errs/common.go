package errs

import "github.com/yoyo-inc/yoyo/core"

var (
	ErrInvalidParameters = core.NewParameterError("无效的参数")
)
