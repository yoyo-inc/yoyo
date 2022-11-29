package models

import (
	"github.com/yoyo-inc/yoyo/common/db"
	"github.com/yoyo-inc/yoyo/core"
)

type Resource struct {
	core.Model
	ResourceName string `json:"resourceName" gorm:"size:100;comment:资源名"`
	Filename     string `json:"filename" gorm:"size:100;index;comment:原始文件名"`
	Filesize     string `json:"filesize" gorm:"size:20;comment:文件大小"`
	FileType     string `json:"filetype" gorm:"size:20;comment:文件类型"`
}

func init() {
	db.AddAutoMigrateModel(&Resource{})
}
