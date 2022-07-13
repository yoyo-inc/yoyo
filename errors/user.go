package errors

import "github.com/yoyo-inc/yoyo/core"

var (
	ErrFailedToLogin      = core.NewBusinessError("100", "登录失败")
	ErrUsernameNotExists  = core.NewBusinessError("101", "账号不存在")
	ErrUsernameOrPassword = core.NewBusinessError("102", "账号或密码错误")
	ErrCreateUser         = core.NewBusinessError("103", "用户创建失败")
	ErrQueryUser          = core.NewBusinessError("104", "用户查询失败")
	ErrDeleteUser         = core.NewBusinessError("105", "用户删除失败")
	ErrUpdateUser         = core.NewBusinessError("106", "用户更新失败")
)
