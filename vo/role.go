package vo

import "github.com/yoyo-inc/yoyo/models"

type RoleVO struct {
	models.Role
	Permissions []int `json:"permissions"`
}

type QueryRoleVO struct {
	// 角色名
	Name string `form:"name"`
	// 是否默认
	Default bool `form:"default"`
}
