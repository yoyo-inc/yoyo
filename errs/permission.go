package errs

import "github.com/yoyo-inc/yoyo/core"

var (
  ErrQueryPermissions = core.NewBusinessError("400", "查询权限列表失败")
)
