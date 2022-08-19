package models

import (
	"github.com/yoyo-inc/yoyo/common/db"
	"github.com/yoyo-inc/yoyo/core"
)

type AuditLog struct {
	core.Model
	UserID    string `json:"user_id"`
	User      User   `json:"user"`
	IP        string `json:"ip"`
	Module    string `json:"module"`
	Operation string `json:"operation"`
	Action    string `json:"action"`
	Detail    string `json:"detail"`
}

func init() {
	db.AddAutoMigrateModel(&AuditLog{})
}
