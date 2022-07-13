package services

import (
	"github.com/yoyo-inc/yoyo/common/db"
	"github.com/yoyo-inc/yoyo/errors"
	"github.com/yoyo-inc/yoyo/models"
)

// DoLogin does login
func DoLogin(username string, password string) (user models.User, err error) {
	res := db.Client.Model(&models.User{}).Where("username = ?", username).Find(&user)
	if res.Error != nil {
		err = errors.ErrUsernameNotExists
		return
	}
	if !user.Check(password) {
		err = errors.ErrUsernameOrPassword
	}
	return
}
