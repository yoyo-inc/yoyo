package models

import (
	"github.com/ypli0629/yoyo/common/db"
	"github.com/ypli0629/yoyo/core"
)

// Role represents user role
type Role struct {
	core.Model
	Name   string `json:"name" gorm:"size:256;unique;not null;comment:角色名" binding:"required,max=250"`
	Remark string `json:"remark" gorm:"size:256;comment:备注"`
}

func init() {
	db.AddAutoMigrateModel(&Role{})
}
