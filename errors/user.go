package errors

import "github.com/ypli0629/yoyo/core"

var (
	FailedToLogin         = core.NewBusinessError("100", "登录失败")
	ErrUsernameNotExists  = core.NewBusinessError("101", "账号不存在")
	ErrUsernameOrPassword = core.NewBusinessError("102", "账号或密码错误")
)
