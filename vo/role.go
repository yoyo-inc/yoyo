package vo

import "github.com/yoyo-inc/yoyo/models"

type RoleVO struct {
	models.Role
	Permissions []string `json:"permissions"`
}

type QueryRoleVO struct {
	// 角色名
	Name string `json:"name"`
}
