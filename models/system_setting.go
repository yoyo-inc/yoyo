package models

import (
	"github.com/yoyo-inc/yoyo/common/db"
	"github.com/yoyo-inc/yoyo/core"
)

type SystemSetting struct {
	core.IModel
	Name        string    `json:"name" gorm:"size:100;comment:系统名称"`
	Description string    `json:"description" gorm:"size:1000;comment:系统描述"`
	ResourceID  string    `json:"resourceID" gorm:"comment:资源ID"`
	Logo        *Resource `json:"logo" gorm:"foreignKey:ResourceID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

func init() {
	db.AddAutoMigrateModel(&SystemSetting{})
}
