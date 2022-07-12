package core

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Model struct {
	// 主键
	ID string `json:"id" gorm:"primarykey;size:32;comment:主键"`
	// 创建时间
	CreateTime *time.Time `json:"createTime" gorm:"autoCreateTime;comment:创建时间"`
	// 更新时间
	ModifyTime *time.Time `json:"modifyTime" gorm:"autoUpdateTime;comment:更新时间"`
}

func (m *Model) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		id := uuid.NewString()
		id = strings.ReplaceAll(id, "-", "")
		m.ID = id
	}

	m.CreateTime = nil
	m.ModifyTime = nil
	return
}

func (m *Model) BeforeUpdate(tx *gorm.DB) (err error) {
	m.CreateTime = nil
	m.ModifyTime = nil
	return
}
