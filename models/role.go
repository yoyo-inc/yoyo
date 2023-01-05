package models

import (
	"github.com/yoyo-inc/yoyo/common/db"
	"github.com/yoyo-inc/yoyo/core"
)

// Role represents user role
type Role struct {
	core.IModel
	Name        string       `json:"name,omitempty" gorm:"size:256;unique;not null;index;comment:角色名" binding:"required,max=250"`
	Remark      string       `json:"remark,omitempty" gorm:"size:256;comment:备注"`
	Permissions []Permission `json:"permissions,omitempty" gorm:"many2many:role_permissions;"`
	IsDefault   bool         `json:"is_default,omitempty" gorm:"comment:默认角色"`
}

// RolePermission represents the relations between role and permission
type RolePermission struct {
	RoleID       string `json:"roleID"`
	PermissionID string `json:"permissionID"`
}

func init() {
	db.AddAutoMigrateModel(&Role{})
}
