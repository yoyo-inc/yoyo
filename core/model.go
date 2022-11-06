package core

import (
	"strings"

	"github.com/google/uuid"
	"github.com/yoyo-inc/yoyo/common/datatypes"
	"gorm.io/gorm"
)

type Model struct {
	// 主键
	ID int `json:"id" gorm:"primarykey;autoIncreatment;comment:主键"`
	// 创建时间
	CreateTime *datatypes.LocalTime `json:"createTime" gorm:"type:timestamp;default:current_timestamp;<-:create;comment:创建时间"`
	// 更新时间
	ModifyTime *datatypes.LocalTime `json:"modifyTime" gorm:"type:timestamp;default:current_timestamp on update current_timestamp;comment:更新时间"`
}

type SModel struct {
	// 主键
	ID string `json:"id" gorm:"primarykey;size:32;comment:主键"`
	// 创建时间
	CreateTime *datatypes.LocalTime `json:"createTime" gorm:"type:timestamp;default:current_timestamp;<-:create;comment:创建时间"`
	// 更新时间
	ModifyTime *datatypes.LocalTime `json:"modifyTime" gorm:"type:timestamp;default:current_timestamp on update current_timestamp;comment:更新时间"`
}

func (m *SModel) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		id := uuid.NewString()
		id = strings.ReplaceAll(id, "-", "")
		m.ID = id
	}

	m.CreateTime = nil
	m.ModifyTime = nil
	return
}

func (m *SModel) BeforeUpdate(tx *gorm.DB) (err error) {
	m.CreateTime = nil
	m.ModifyTime = nil
	return
}
