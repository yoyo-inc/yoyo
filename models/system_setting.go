package models

import "github.com/yoyo-inc/yoyo/core"

type SystemSetting struct {
	core.Model
	Name        string `json:"name" gorm:"size:100;comment:系统名称"`
	Description string `json:"description" gorm:"size:1000;comment:系统描述"`
	Logo        string `json:"logo" gorm:"size:100;comment:系统logo"`
}
