package errs

import "github.com/yoyo-inc/yoyo/core"

var (
	ErrQuerySystemSecurity  = core.NewBusinessError("1200", "查询系统安全设置失败")
	ErrUpdateSystemSecurity = core.NewBusinessError("1201", "更新系统安全设置失败")
)
