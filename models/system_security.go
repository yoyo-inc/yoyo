package models

import "github.com/yoyo-inc/yoyo/core"

type SystemSecurity struct {
	core.Model
	ForbidRepeatLogin *bool `json:"forbidRepeatLogin,omitempty" gorm:"comment:是否启用禁用重复登录"`
	LoginExpire       *bool `json:"loginExpire,omitempty" gorm:"comment:是否启用登录过期"`
	LoginExpireTime   int   `json:"loginExpireTime,omitempty" gorm:"comment:登录过期时间，单位小时"`
}
