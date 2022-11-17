package models

import (
	"github.com/yoyo-inc/yoyo/common/db"
	"github.com/yoyo-inc/yoyo/common/dt"
	"github.com/yoyo-inc/yoyo/core"
	"gorm.io/datatypes"
)

type Alert struct {
	core.Model
	StartAt        *dt.LocalTime `json:"startAt" gorm:"comment:告警产生时间"`
	Type           string        `json:"type" gorm:"size:20;index;comment:类型"`
	Level          string        `json:"level" gorm:"size:10;index;comment:级别"`
	Content        string        `json:"content" gorm:"size:2000;comment:告警内容"`
	Status         int           `json:"status" gorm:"comment:状态(0:未读 1:已读)"`
	ResolvedStatus int           `json:"resolvedStatus" gorm:"comment:处置状态(0:未处置 1:已手动处置 2:自动处置)"`
	Remark         string        `json:"remark" gorm:"size:2000;comment:备注"`
}

type AlertConfig struct {
	core.Model
	EmailEnable      bool           `json:"emailEnable" gorm:"comment:邮件通知是否开启"`
	SmtpServer       string         `json:"smtpServer" gorm:"size:255;comment:邮件服务地址"`
	SmtpSender       string         `json:"smtpSender" gorm:"size:255;comment:邮件发送方"`
	SmtpAuthUser     string         `json:"smtpAuthUser" gorm:"size:255;comment:邮件发送方账号"`
	SmtpAuthPassword string         `json:"smtpAuthPassword" gorm:"size:255;comment:邮件发送方密码"`
	SmtpReceivers    datatypes.JSON `json:"smtpReceivers" gorm:"type:json;default:(json_array());comment:邮件接收人" swaggertype:"array,object"`
}

type AlertAccess struct {
	core.Model
	AccessIP     string `json:"accessIP" gorm:"size:20;comment:接入IP" binding:"required"`
	Remark       string `json:"remark" gorm:"size:2000;comment:备注"`
	TypeField    string `json:"typeField" gorm:"size:20;comment:类型字段名"`
	LevelField   string `json:"levelField" gorm:"size:20;comment:级别字段名"`
	ContentField string `json:"contentField" gorm:"size:20;comment:告警内容字段名"`
}

func init() {
	db.AddAutoMigrateModel(&Alert{})
	db.AddAutoMigrateModel(&AlertConfig{})
	db.AddAutoMigrateModel(&AlertAccess{})
}
