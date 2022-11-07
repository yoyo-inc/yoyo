package models

import (
	"github.com/yoyo-inc/yoyo/common/db"
	"github.com/yoyo-inc/yoyo/core"
	"github.com/yoyo-inc/yoyo/utils"
	"gorm.io/gorm"
)

// User represents user model
type User struct {
	core.Model
	Username       string       `json:"username" gorm:"size:15;unique;not null;index;comment:账户名" binding:"required,min=0,max=15"` // 账户名
	Nickname       string       `json:"nickname" gorm:"size:255;index;comment:昵称" binding:"max=255"`                               // 昵称
	Password       string       `json:"password,omitempty" gorm:"size:255;not null;comment:密码" binding:"max=255"`                  // 密码
	Email          string       `json:"email" gorm:"size:255;comment:邮箱" binding:"email,max=255"`                                  // 邮箱
	Phone          string       `json:"phone" gorm:"size:11;comment:手机号" binding:"len=11"`                                         // 手机号
	Avatar         string       `json:"avatar" gorm:"size:255;comment:头像" binding:"max=255"`                                       // 头像
	Sex            int          `json:"sex" gorm:"comment:性别"`                                                                     // 性别 0: 男 1: 女
	Age            int          `json:"age" gorm:"comment:年龄"`                                                                     // 年龄
	Roles          []Role       `json:"roles" gorm:"many2many:user_roles"`                                                         // 角色
	OrganizationID int          `json:"organizationID" gorm:"comment:组织ID"`                                                        // 组织ID
	Organization   Organization `json:"organization"`                                                                              // 组织
}

// UserRole represents the relations between user and role
type UserRole struct {
	UserID string `json:"userID"`
	RoleID string `json:"roleID"`
}

func init() {
	db.AddAutoMigrateModel(&User{})
}

// Check checks whether the passwords are the same
func (user User) Check(actual string) bool {
	return user.Password == utils.Encrypt(actual)
}

func (user *User) BeforeSave(tx *gorm.DB) (err error) {
	if user.Password != "" {
		user.Password = utils.Encrypt(user.Password)
	}
	return
}

func (user *User) BeforeUpdate(tx *gorm.DB) (err error) {
	if user.Password == "" {
		tx.Statement.Omit("password")
	}
	return
}
