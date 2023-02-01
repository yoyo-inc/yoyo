package errs

import "github.com/yoyo-inc/yoyo/core"

var (
	ErrFailToLogin        = core.NewBusinessError("010", "登录失败")
	ErrUsernameNotExists  = core.NewBusinessError("011", "账号不存在")
	ErrUsernameOrPassword = core.NewBusinessError("012", "账号或密码错误")
	ErrRepeatLogin        = core.NewBusinessError("013", "该帐号已在其他地方登录")
)
