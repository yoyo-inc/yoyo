package models

import (
	"github.com/yoyo-inc/yoyo/common/db"
	"github.com/yoyo-inc/yoyo/core"
)

type Organization struct {
	core.IModel
	// 组织名
	Name     string `json:"name,omitempty" gorm:"size:255;comment:组织名称"`
	ParentID int    `json:"parentId,omitempty"`
}

func init() {
	db.AddAutoMigrateModel(&Organization{})
}
