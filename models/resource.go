package models

import (
	"github.com/yoyo-inc/yoyo/common/db"
	"github.com/yoyo-inc/yoyo/core"
)

type Resource struct {
	core.Model
	ResourceName string  `json:"resourceName,omitempty" gorm:"size:100;comment:资源名"`
	ResourceType string  `json:"resourceType,omitempty" gorm:"size:20;comment:资源类型"`
	Filename     string  `json:"filename,omitempty" gorm:"size:100;index;comment:原始文件名"`
	Filesize     float64 `json:"filesize,omitempty" gorm:"size:20;comment:文件大小"`
	FileType     string  `json:"filetype,omitempty" gorm:"size:20;comment:文件类型"`
}

func init() {
	db.AddAutoMigrateModel(&Resource{})
}
