package routes

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/yoyo-inc/gin-jwt/v3"
	"github.com/yoyo-inc/yoyo/common/db"
	"github.com/yoyo-inc/yoyo/common/logger"
	"github.com/yoyo-inc/yoyo/middlewares"
	"github.com/yoyo-inc/yoyo/models"
	"github.com/yoyo-inc/yoyo/services/audit_log"
)

type loginController struct{}

func (*loginController) Logout(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	if userID, ok := claims["userID"]; ok {
		if res := db.Client.Where("user_id = ?", userID).Delete(&models.LoginSession{}); res.Error != nil {
			logger.Error(res.Error)
		}
	}

	middlewares.SecurityMiddleware.LogoutHandler(c)
	audit_log.Success(c, "用户", "登出", "退出登录成功")
}

func (login *loginController) Setup(r *gin.RouterGroup) {
	r.POST("/login", middlewares.SecurityMiddleware.LoginHandler).GET("/logout", login.Logout)
}
