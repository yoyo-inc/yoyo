package models

import (
	"github.com/ypli0629/yoyo/core"
	"github.com/ypli0629/yoyo/utils"
)

// User represents user model
type User struct {
	core.Model
	Username string `json:"username" gorm:"size:15;unique;not null;index;comment:账户名"`
	Password string `json:"password" gorm:"size:256;not null;comment:密码"`
	Email    string `json:"email" gorm:"size:256;comment:邮箱"`
	Phone    string `json:"phone" gorm:"size:11;comment:手机号"`
	Avatar   string `json:"avatar" gorm:"size:256;comment:头像"`
	Sex      int    `json:"sex" gorm:"comment:性别"`
	Age      int    `json:"age" gorm:"comment:年龄"`
}

// Check checks whether the passwords are the same
func (user User) Check(actual string) bool {
	return user.Password == utils.Encrypt(actual)
}
