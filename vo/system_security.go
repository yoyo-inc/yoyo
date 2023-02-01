package vo

import "github.com/yoyo-inc/yoyo/models"

type SystemSecurityVO struct {
	models.SystemSecurity
	LoginIPWhitelist []map[string]interface{} `json:"loginIPWhitelist,omitempty"`
}
