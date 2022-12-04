package models

import (
	"github.com/yoyo-inc/yoyo/common/db"
	"github.com/yoyo-inc/yoyo/core"
)

type Report struct {
	core.IModel
	Name string  `json:"name" gorm:"size:100;comment:报告名称"`
	Type string  `json:"type" gorm:"size:10;comment:报告类型"`
	Size float64 `json:"size" gorm:"comment:文件大小"`
}

type ReportConfig struct {
	core.IModel
}

type ReportJob struct {
	core.IModel
}

func init() {
	db.AddAutoMigrateModel(&Report{})
}
