package models

import (
	"github.com/yoyo-inc/yoyo/common/db"
	"github.com/yoyo-inc/yoyo/core"
)

type LogConfig struct {
	core.Model
	KeepTime *int `json:"keepTime,omitempty" gorm:"default:30;comment:日志保留时间，单位天"`
	Archive  bool `json:"archive,omitempty" gorm:"default: true;comment:日志归档"`
}

func init() {
	db.AddAutoMigrateModel(&LogConfig{})
}
