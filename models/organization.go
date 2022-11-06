package models

import (
	"github.com/yoyo-inc/yoyo/common/db"
	"github.com/yoyo-inc/yoyo/core"
)

type Organization struct {
	core.Model
	// 组织名
	Name     string `json:"name" gorm:"size:255;comment:组织名称"`
	ParentID int    `json:"parentId"`
}

func init() {
	db.AddAutoMigrateModel(&Organization{})
}
