package services

import (
	"github.com/ypli0629/yoyo/common/db"
	"github.com/ypli0629/yoyo/errors"
	"github.com/ypli0629/yoyo/models"
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
