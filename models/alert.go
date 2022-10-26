package models

import (
	"github.com/yoyo-inc/yoyo/common/datatypes"
	"github.com/yoyo-inc/yoyo/common/db"
	"github.com/yoyo-inc/yoyo/core"
)

type Alert struct {
	core.Model
	StartAt        *datatypes.LocalTime `json:"startAt" gorm:"comment:告警产生时间"`
	Type           string               `json:"type" gorm:"comment:类型"`
	Level          string               `json:"level" gorm:"comment:级别"`
	Content        string               `json:"content" gorm:"comment:告警内容"`
	Status         int                  `json:"status" gorm:"comment:状态(0:未读 1:已读)"`
	ResolvedStatus int                  `json:"resolvedStatus" gorm:"comment:处置状态(0:未处置 1:已手动处置 2:自动处置)"`
	Remark         string               `json:"remark" gorm:"comment:备注"`
}

func init() {
	db.AddAutoMigrateModel(&Alert{})
}
