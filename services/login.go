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

type loginPayload struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// DoLogin does login
func DoLogin(c *gin.Context) (user models.User, err error) {
	var systemSecurity models.SystemSecurity
	if res := db.Client.Model(&models.SystemSecurity{}).First(&systemSecurity); res.Error != nil {
		logger.Error(res.Error)
	}

	// ip whitelist
	if systemSecurity.LoginIPWhitelistEnable {
		var canPass bool
		clientIP := utils.GetClientIP(c)
		for _, item := range systemSecurity.GetLoginIPWhitelist() {
			if enable, ok := item["enable"].(bool); ok && enable {
				if ip, ok := item["ip"]; ok && ip == clientIP {
					canPass = true
					break
				}
			}
		}
		if !canPass {
			err = errs.ErrFailToLogin
			return
		}
	}

	var payload loginPayload
	if e := c.ShouldBindJSON(&payload); e != nil {
		err = errs.ErrUsernameOrPassword
		return
	}
	if res := db.Client.Model(&models.User{}).Where("username = ?", payload.Username).First(&user); errors.Is(res.Error, gorm.ErrRecordNotFound) {
		err = errs.ErrUsernameOrPassword
		audit_log.Fail(c, "用户", "登录", fmt.Sprintf("账号(%s)不存在", payload.Username))
		return
	}

	if systemSecurity.ForbidRepeatLogin == true {
		var count int64
		if res := db.Client.Model(&models.LoginSession{}).Where("username = ?", payload.Username).Count(&count); res.Error == nil {
			if count > 0 {
				err = errs.ErrRepeatLogin
				audit_log.Fail(c, "用户", "登录", fmt.Sprintf("帐号（%s）重复登录", payload.Username))
				return
			}
		}
	}

	if !user.Check(payload.Password) {
		err = errs.ErrUsernameOrPassword
		audit_log.Fail(c, "用户", "登录", fmt.Sprintf("账号(%s)或密码错误", payload.Username))
		return
	}
	audit_log.Success(c, "用户", "登录", fmt.Sprintf("账号(%s)登录成功", payload.Username))

	// record login session
	if res := db.Client.Create(&models.LoginSession{
		IP:       utils.GetClientIP(c),
		UserID:   user.ID,
		Username: payload.Username,
	}); res.Error != nil {
		logger.Errorf("Failed to record login session: %s", res.Error)
	}
	return
}
