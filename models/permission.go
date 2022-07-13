package models

import (
	"github.com/yoyo-inc/yoyo/common/db"
	"github.com/yoyo-inc/yoyo/core"
)

type Permission struct {
	core.Model
	// 权限名称
	Name string `json:"name" gorm:"size:255;unique;not null;comment:权限名称"`
	// 父级权限
	ParentID string `json:"parentID" gorm:"size:32;comment:父级权限"`
}

func init() {
	db.AddAutoMigrateModel(&Permission{})
}
