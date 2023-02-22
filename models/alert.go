package models

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/yoyo-inc/yoyo/common/db"
	"github.com/yoyo-inc/yoyo/common/dt"
	"github.com/yoyo-inc/yoyo/core"
	"gorm.io/datatypes"
)

type Alert struct {
	core.IModel
	StartAt        *dt.LocalTime `json:"startAt,omitempty" gorm:"comment:告警产生时间"`
	Type           string        `json:"type,omitempty" form:"type" gorm:"size:20;index;comment:类型"`
	Level          string        `json:"level,omitempty" form:"level" gorm:"size:10;index;comment:级别"`
	Content        string        `json:"content,omitempty" gorm:"size:2000;comment:告警内容"`
	Status         int           `json:"status" gorm:"comment:状态(0:未读 1:已读)"`
	ResolvedStatus int           `json:"resolvedStatus" gorm:"comment:处置状态(0:未处置 1:已手动处置 2:自动处置 3:已忽略)"`
	Remark         string        `json:"remark,omitempty" gorm:"size:2000;comment:备注"`
	From           string        `json:"from,omitempty" form:"from" gorm:"size:255;default:本系统;comment:来源"`
}

type AlertConfig struct {
	core.IModel
	EmailEnable      bool           `json:"emailEnable,omitempty" gorm:"comment:邮件通知是否开启"`
	SmtpServer       string         `json:"smtpServer,omitempty" gorm:"size:255;comment:邮件服务地址"`
	SmtpSender       string         `json:"smtpSender,omitempty" gorm:"size:255;comment:邮件发送方"`
	SmtpAuthUser     string         `json:"smtpAuthUser,omitempty" gorm:"size:255;comment:邮件发送方账号"`
	SmtpAuthPassword string         `json:"smtpAuthPassword,omitempty" gorm:"size:255;comment:邮件发送方密码"`
	SmtpReceivers    datatypes.JSON `json:"smtpReceivers,omitempty" gorm:"type:json;default:(json_array());comment:邮件接收人" swaggertype:"array,object"`
}

type AlertAccess struct {
	core.IModel
	AccessIP     string `json:"accessIP,omitempty" gorm:"size:20;comment:接入IP" binding:"required"`
	Remark       string `json:"remark,omitempty" gorm:"size:2000;comment:备注"`
	StartAtField string `json:"startAtField,omitempty" gorm:"size:20;comment:告警时间字段名"`
	TypeField    string `json:"typeField,omitempty" gorm:"size:20;comment:类型字段名"`
	LevelField   string `json:"levelField,omitempty" gorm:"size:20;comment:级别字段名"`
	ContentField string `json:"contentField,omitempty" gorm:"size:20;comment:告警内容字段名"`
}

type AlertPush struct {
	core.IModel
	Enable        bool           `json:"enable,omitempty" gorm:"comment:是否启用"`
	Type          string         `json:"type,omitempty" gorm:"size:20;comment:推送类型"`
	SyslogNetwork string         `json:"syslogNetwork,omitempty" gorm:"size:20;comment:syslog网络类型(tcp|udp)"`
	SyslogAddr    string         `json:"syslogAddr,omitempty" gorm:"size:20;comment:syslog地址"`
	SyslogPort    int            `json:"syslogPort,omitempty" gorm:"comment:syslog端口"`
	SyslogTag     string         `json:"syslogTag,omitempty" gorm:"size:100;comment:syslog标志"`
	KafkaNetwork  string         `json:"kafkaNetwork,omitempty" gorm:"size:10;comment:kafka网络类型(tcp|udp)"`
	KafkaAddr     string         `json:"kafkaAddr,omitempty" gorm:"size:20;commnet:kafka服务地址"`
	KafkaPort     int            `json:"kafkaPort,omitempty" gorm:"comment:kafka端口"`
	KafkaTopic    string         `json:"kafkaTopic,omitempty" gorm:"comment:kafka topic"`
	AlertTypes    datatypes.JSON `json:"alertTypes,omitempty" gorm:"type:json;default:(json_array());comment:告警类型" swaggertype:"array,string"`
}

func (ap *AlertPush) GetAlertTypes() []string {
	if ap.AlertTypes != nil {
		return nil
	}

	var alertTypes []string
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	if err := json.Unmarshal(ap.AlertTypes, &alertTypes); err != nil {
		return nil
	}

	return alertTypes
}

func init() {
	db.AddAutoMigrateModel(&Alert{})
	db.AddAutoMigrateModel(&AlertConfig{})
	db.AddAutoMigrateModel(&AlertAccess{})
	db.AddAutoMigrateModel(&AlertPush{})
}
