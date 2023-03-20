package models

import (
	"github.com/yoyo-inc/yoyo/common/db"
	"github.com/yoyo-inc/yoyo/core"
)

type License struct {
	core.Model
	MachineCode string `json:"machine_code" gorm:"comment:机器码"`
	Company     string `json:"company" gorm:"comment:公司名称"`
	IAT         string `json:"iat" gorm:"comment:签发时间"`
	EXP         string `json:"exp" gorm:"comment:过期时间"`
	ACC         string `json:"acc" gorm:"comment:激活码"`
}

func init() {
	db.AddAutoMigrateModel(&License{})
}
