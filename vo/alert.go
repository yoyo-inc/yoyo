package vo

import "github.com/yoyo-inc/yoyo/models"

type QueryAlertVO struct {
	models.Alert
	Status         *int `form:"status"`
	ResolvedStatus *int `form:"resolvedStatus"`
}

type ResolveAlertVO struct {
	models.Alert
	ID string `json:"id" binding:"required"`
}

type IgnoreAlertVO struct {
	ID string `json:"id" binding:"required"`
}

type SmtpReceiver struct {
	ID     int    `json:"id"`
	Email  string `json:"email"`
	Enable bool   `json:"enable"`
}

type UpdateAlertConfigVO struct {
	models.AlertConfig
}

type QueryAlertAccessVO struct {
	AccessIP string `form:"accessIP"`
}

type UpdateAlertAccessVO struct {
	models.AlertAccess
	ID int `json:"id" binding:"required"`
}

type QueryAlertCountVO struct {
	Status         int `form:"status"`
	ResolvedStatus int `form:"resolvedStatus"`
}
