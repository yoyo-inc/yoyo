package vo

import "github.com/yoyo-inc/yoyo/models"

type PermissionVO struct {
	models.Permission
	Children []*PermissionVO `json:"children"`
}
