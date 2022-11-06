package models

import (
	"github.com/yoyo-inc/yoyo/common/db"
	"github.com/yoyo-inc/yoyo/core"
)

type AuditLog struct {
	core.Model
	UserID    int    `json:"userID" gorm:"comment:用户ID"`
	User      User   `json:"user"`
	IP        string `json:"ip" gorm:"comment:ip"`
	Module    string `json:"module" gorm:"comment:模块"`
	Operation string `json:"operation" gorm:"comment:操作"`
	Status    int    `json:"status" gorm:"是否成功"`
	Detail    string `json:"detail" gorm:"详情"`
}

func init() {
	db.AddAutoMigrateModel(&AuditLog{})
}
