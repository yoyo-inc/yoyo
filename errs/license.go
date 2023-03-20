package errs

import "github.com/yoyo-inc/yoyo/core"

var (
	ErrQueryLicense = core.NewBusinessError("", "查询激活信息失败")
)
