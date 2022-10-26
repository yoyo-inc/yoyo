package vo

import "github.com/yoyo-inc/yoyo/models"

type UpdateAlertVO struct {
	models.Alert
	ID string `json:"id" binding:"required"`
}

type QueryAlertVO struct{}
