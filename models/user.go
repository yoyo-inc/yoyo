package models

import (
	"github.com/ypli0629/yoyo/core"
	"github.com/ypli0629/yoyo/utils"
)

// User represents user model
type User struct {
	core.Model
	// Username 账户名
	Username string `json:"username" gorm:"size:15;unique;not null;index;comment:账户名" binding:"required,min=0,max=15"`
	// Password 密码
	Password string `json:"password,omitempty" gorm:"size:256;not null;comment:密码" binding:"required,max=250"`
	// 邮箱
	Email string `json:"email" gorm:"size:256;comment:邮箱" binding:"email"`
	// 手机号
	Phone string `json:"phone" gorm:"size:11;comment:手机号" bining:"len=11"`
	// 头像
	Avatar string `json:"avatar" gorm:"size:256;comment:头像" binding:"max=250"`
	// 性别
	Sex int `json:"sex" gorm:"comment:性别"`
	// 年龄
	Age int `json:"age" gorm:"comment:年龄"`
}

// Check checks whether the passwords are the same
func (user User) Check(actual string) bool {
	return user.Password == utils.Encrypt(actual)
}
