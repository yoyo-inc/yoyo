package vo

import "github.com/yoyo-inc/yoyo/models"

type UpdateSystemSettingVO struct {
	models.SystemSetting
	ID int `json:"id" binding:"required"`
}
