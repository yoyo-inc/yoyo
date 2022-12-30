package models

import (
	"github.com/yoyo-inc/yoyo/common/db"
	"github.com/yoyo-inc/yoyo/core"
)

type Dict struct {
	core.IModel
	Type      string `json:"type" gorm:"comment:类型"`
	Label     string `json:"label" gorm:"comment:标签"`
	Value     string `json:"value" gorm:"comment:值"`
	ValueType string `json:"valueType" gorm:"comment:值类型"`
}

func init() {
	db.AddAutoMigrateModel(&Dict{})
}
