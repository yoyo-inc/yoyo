package services

import (
	"errors"
	"github.com/yoyo-inc/yoyo/common/db"
	"github.com/yoyo-inc/yoyo/errs"
	"github.com/yoyo-inc/yoyo/models"
	"gorm.io/gorm"
)

// DoLogin does login
func DoLogin(username string, password string) (user models.User, err error) {

	if res := db.Client.Model(&models.User{}).Where("username = ?", username).First(&user); errors.Is(res.Error, gorm.ErrRecordNotFound) {
		err = errs.ErrUsernameNotExists
		return
	}
	if !user.Check(password) {
		err = errs.ErrUsernameOrPassword
	}
	return
}
