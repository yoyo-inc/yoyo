package models

import (
	"github.com/ypli0629/yoyo/common/db"
	"github.com/ypli0629/yoyo/core"
)

type Organization struct {
	core.Model
	// 组织名
	Name     string `json:"name" gorm:"size:255;comment:组织名称"`
	ParentID string `json:"parentId"`
}

func init() {
	db.AddAutoMigrateModel(&Organization{})
}
