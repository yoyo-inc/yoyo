package errs

import "github.com/yoyo-inc/yoyo/core"

var (
	ErrCreateUser       = core.NewBusinessError("101", "用户创建失败")
	ErrQueryUser        = core.NewBusinessError("102", "用户查询失败")
	ErrDeleteUser       = core.NewBusinessError("103", "用户删除失败")
	ErrUpdateUser       = core.NewBusinessError("104", "用户更新失败")
	ErrQueryCurrentUser = core.NewBusinessError("105", "当前用户查询失败")
	ErrExistUser        = core.NewBusinessError("106", "用户已存在")
)
