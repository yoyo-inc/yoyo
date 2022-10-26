package models

import (
	"github.com/yoyo-inc/yoyo/common/db"
	"github.com/yoyo-inc/yoyo/core"
)

type Permission struct {
	core.Model
	// ID int `json:"id" gorm:"primaryKey;autoIncrement;comment:主键"`
	// 权限名称
	Name string `json:"name" gorm:"size:255;unique;not null;comment:权限名称"`
	// 介绍
	Description string `json:"description" gorm:"comment:权限描述"`
	// 父级权限
	ParentID string `json:"parentID" gorm:"size:32;comment:父级权限"`
	// 是否开启
	Enable bool `json:"enable" gorm:"comment:是否启用"`
}

func init() {
	db.AddAutoMigrateModel(&Permission{})
}
