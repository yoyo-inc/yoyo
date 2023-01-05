package models

import (
	"github.com/yoyo-inc/yoyo/common/db"
	"github.com/yoyo-inc/yoyo/core"
)

type Permission struct {
	core.IModel
	Name        string `json:"name,omitempty" gorm:"size:255;unique;not null;comment:权限名称"` // 权限名称
	Description string `json:"description,omitempty" gorm:"comment:权限描述"`                   // 介绍
	ParentID    int    `json:"parentID,omitempty" gorm:"comment:父级权限"`                      // 父级权限
	Enable      bool   `json:"enable,omitempty" gorm:"default:true;comment:是否启用"`           // 是否开启
}

func init() {
	db.AddAutoMigrateModel(&Permission{})
}
