package core

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Model struct {
	ID         string     `json:"id" gorm:"primarykey;size:32"`
	CreateTime *time.Time `json:"createTime" gorm:"autoCreateTime;comment:创建时间"`
	ModifyTime *time.Time `json:"modifyTime" gorm:"autoUpdateTime;comment:更新时间"`
}

func (m *Model) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		id := uuid.NewString()
		id = strings.ReplaceAll(id, "-", "")
		m.ID = id
	}

	return
}
