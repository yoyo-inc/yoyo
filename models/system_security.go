package models

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/yoyo-inc/yoyo/common/db"
	"github.com/yoyo-inc/yoyo/core"
	"gorm.io/datatypes"
)

type SystemSecurity struct {
	core.IModel
	ForbidRepeatLogin      *bool           `json:"forbidRepeatLogin,omitempty" gorm:"comment:是否启用禁用重复登录"`
	LoginExpire            *bool           `json:"loginExpire,omitempty" gorm:"comment:是否启用登录过期"`
	LoginExpireTime        int             `json:"loginExpireTime,omitempty" gorm:"comment:登录过期时间，单位小时"`
	LoginIPWhitelistEnable *bool           `json:"loginIPWhitelistEnable,omitempty" gorm:"comment:登录IP限制"`
	LoginIPWhitelist       *datatypes.JSON `json:"loginIPWhitelist,omitempty" gorm:"comment:登录IP白名单" swaggertype:"array,object"`
}

func (ss *SystemSecurity) GetLoginIPWhitelist() []map[string]interface{} {
	if ss.LoginIPWhitelist == nil {
		return nil
	}

	json := jsoniter.ConfigCompatibleWithStandardLibrary

	var ipWhitelist []map[string]interface{}
	if err := json.Unmarshal(*ss.LoginIPWhitelist, &ipWhitelist); err != nil {
		return nil
	}

	return ipWhitelist
}

type LoginSession struct {
	core.Model
	IP       string
	UserID   int
	Username string
}

func init() {
	db.AddAutoMigrateModel(&SystemSecurity{})
	db.AddAutoMigrateModel(&LoginSession{})
}
