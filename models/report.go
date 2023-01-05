package models

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/yoyo-inc/yoyo/common/db"
	"github.com/yoyo-inc/yoyo/core"
	"gorm.io/datatypes"
)

type Report struct {
	core.IModel
	ReportName   string   `json:"reportName,omitempty" gorm:"size:100;comment:报告名称"`
	ReportType   string   `json:"reportType,omitempty" gorm:"size:20;comment:报告类型"`
	ReportStatus int      `json:"reportStatus,omitempty" gorm:"comment:报告状态（0：生成中，1：成功，2：失败）"`
	ResourceID   string   `json:"resourceID,omitempty" gorm:"comment:资源ID"`
	Resource     Resource `json:"resource,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type ReportConfig struct {
	core.Model
	Period     datatypes.JSON `json:"period,omitempty" gorm:"type:json;default:(json_array());comment:周期" swaggertype:"array,string"`
	ReportType datatypes.JSON `json:"reportType,omitempty" gorm:"type:json;default:(json_array());comment:报告类型" swaggertype:"array,string"`
}

func (rc ReportConfig) GetPeriod() []string {
	json := jsoniter.ConfigCompatibleWithStandardLibrary

	var period []string
	if err := json.Unmarshal(rc.Period, &period); err != nil {
		return []string{}
	}

	return period
}

func (rc ReportConfig) GetReportType() []string {
	json := jsoniter.ConfigCompatibleWithStandardLibrary

	var reportType []string
	if err := json.Unmarshal(rc.ReportType, &reportType); err != nil {
		return []string{}
	}

	return reportType
}

func init() {
	db.AddAutoMigrateModel(&Report{})
	db.AddAutoMigrateModel(&ReportConfig{})
}
