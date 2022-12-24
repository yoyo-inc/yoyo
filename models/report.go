package models

import (
	"github.com/yoyo-inc/yoyo/common/db"
	"github.com/yoyo-inc/yoyo/core"
	"gorm.io/datatypes"
)

type Report struct {
	core.IModel
	ReportName   string   `json:"reportName" gorm:"size:100;comment:报告名称"`
	ReportType   string   `json:"reportType" gorm:"size:20;comment:报告类型"`
	ReportStatus int      `json:"reportStatus" gorm:"comment:报告状态（0：生成中，1：成功，2：失败）"`
	ResourceID   string   `json:"resourceID" gorm:"comment:资源ID"`
	Resource     Resource `json:"resource" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type ReportConfig struct {
	core.Model
	Period     datatypes.JSON `json:"period" gorm:"type:json;default:(json_array());comment:周期" swaggertype:"array,string"`
	ReportType datatypes.JSON `json:"reportType" gorm:"type:json;default:(json_array());comment:报告类型"`
	ReportName string         `json:"reportName" gorm:"size:100;comment:报告名称"`
}

func init() {
	db.AddAutoMigrateModel(&Report{})
}
