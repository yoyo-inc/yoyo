package models

import (
	"github.com/ypli0629/yoyo/common/db"
	"github.com/ypli0629/yoyo/core"
)

type Permission struct {
	core.Model
	// 权限名
	Name     string `json:"name" gorm:"size:255;unique;not null;comment:权限名称"`
	ParentID string `json:"parent_id" gorm:"size:32;comment:父级权限"`
}

func init() {
	db.AddAutoMigrateModel(&Permission{})
}
