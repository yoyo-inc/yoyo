package errs

import "github.com/yoyo-inc/yoyo/core"

var (
	ErrQuerySystemSetting  = core.NewBusinessError("700", "查询系统设置失败")
	ErrUpdateSystemSetting = core.NewBusinessError("701", "更新系统设置失败")
)
