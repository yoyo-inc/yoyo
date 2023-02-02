package services

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/yoyo-inc/yoyo/common/db"
	"github.com/yoyo-inc/yoyo/common/logger"
	"github.com/yoyo-inc/yoyo/errs"
	"github.com/yoyo-inc/yoyo/models"
	"github.com/yoyo-inc/yoyo/services/audit_log"
	"github.com/yoyo-inc/yoyo/utils"
	"gorm.io/gorm"
)

// DoLogin does login
func DoLogin(c *gin.Context, username string, password string) (user models.User, err error) {
	if res := db.Client.Model(&models.User{}).Where("username = ?", username).First(&user); errors.Is(res.Error, gorm.ErrRecordNotFound) {
		err = errs.ErrUsernameNotExists
		audit_log.Fail(c, "用户", "登录", fmt.Sprintf("账号(%s)不存在", username))
		return
	}

	var systemSecurity models.SystemSecurity
	if res := db.Client.Model(&models.SystemSecurity{}).First(&systemSecurity); res.Error == nil {
		if systemSecurity.ForbidRepeatLogin == true {
			var count int64
			if res := db.Client.Model(&models.LoginSession{}).Where("username = ?", username).Count(&count); res.Error == nil {
				if count > 1 {
					err = errs.ErrRepeatLogin
					audit_log.Fail(c, "用户", "登录", fmt.Sprintf("帐号（%s）重复登录", username))
					return
				}
			}
		}
	}

	if !user.Check(password) {
		err = errs.ErrUsernameOrPassword
		audit_log.Fail(c, "用户", "登录", fmt.Sprintf("账号(%s)或密码错误", username))
		return
	}
	audit_log.Success(c, "用户", "登录", fmt.Sprintf("账号(%s)登录成功", username))

	// record login session
	if res := db.Client.Create(&models.LoginSession{
		IP:       utils.GetClientIP(c),
		UserID:   user.ID,
		Username: username,
	}); res.Error != nil {
		logger.Error(res.Error)
	}
	return
}
