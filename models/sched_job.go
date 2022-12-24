package models

import (
	"github.com/yoyo-inc/yoyo/common/db"
	"github.com/yoyo-inc/yoyo/core"
)

type SchedJob struct {
	core.IModel
	Type        string `json:"type" gorm:"size:20;comment:任务类型"`
	Spec        string `json:"spec" gorm:"size:100;comment:规则"`
	Description string `json:"description" gorm:"size:100;comment:描述"`
	Status      int    `json:"status" gorm:"comment:是否启用0：不启用，1：启用"`
	JobID       string `json:"jobID" gorm:"comment:任务唯一标识"`
	LastRunTime string `json:"lastRunTime" gorm:"comment:上次运行时间"`
}

func init() {
	db.AddAutoMigrateModel(&SchedJob{})
}
