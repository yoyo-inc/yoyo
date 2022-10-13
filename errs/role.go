package errs

import "github.com/yoyo-inc/yoyo/core"

var (
	ErrCreateRole = core.NewBusinessError("200", "角色创建失败")
	ErrUpdateRole = core.NewBusinessError("201", "角色更新失败")
	ErrDeleteRole = core.NewBusinessError("202", "角色删除失败")
	ErrQueryRole  = core.NewBusinessError("203", "角色查询失败")
	ErrNotExistRole = core.NewBusinessError("204", "角色不存在")
)
