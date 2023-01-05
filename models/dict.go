package models

import (
	"github.com/yoyo-inc/yoyo/common/db"
	"github.com/yoyo-inc/yoyo/core"
)

type Dict struct {
	core.IModel
	Type      string `json:"type,omitempty" gorm:"comment:类型"`
	Desc      string `json:"desc,omitempty" gorm:"comment:描述"`
	Label     string `json:"label,omitempty" gorm:"comment:标签"`
	Value     string `json:"value,omitempty" gorm:"comment:值"`
	ValueType string `json:"valueType,omitempty" gorm:"comment:值类型"`
}

func init() {
	db.AddAutoMigrateModel(&Dict{})
}
