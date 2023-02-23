package models

import (
	"github.com/yoyo-inc/yoyo/common/db"
	"github.com/yoyo-inc/yoyo/core"
)

type AuditLog struct {
	core.IModel
	UserID    *int   `json:"userID,omitempty" gorm:"comment:用户ID"`
	User      User   `json:"user,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	IP        string `json:"ip,omitempty" gorm:"comment:ip"`
	Module    string `json:"module,omitempty" gorm:"comment:模块"`
	Operation string `json:"operation,omitempty" gorm:"comment:操作"`
	Status    int    `json:"status" gorm:"是否成功"`
	Detail    string `json:"detail,omitempty" gorm:"详情"`
}

func init() {
	db.AddAutoMigrateModel(&AuditLog{})
}
