package audit_log

import (
	"strconv"

	"github.com/gin-gonic/gin"
	jwt "github.com/yoyo-inc/gin-jwt/v3"
	"github.com/yoyo-inc/yoyo/common/db"
	"github.com/yoyo-inc/yoyo/common/logger"
	"github.com/yoyo-inc/yoyo/middlewares"
	"github.com/yoyo-inc/yoyo/models"
	"github.com/yoyo-inc/yoyo/utils"
)

// Log logs audit log
func Log(userID int, clientIP string, module string, operation string, status int, detail string) {
	auditLog := models.AuditLog{
		UserID:    userID,
		IP:        clientIP,
		Module:    module,
		Operation: operation,
		Status:    status,
		Detail:    detail,
	}
	if res := db.Client.Create(&auditLog); res.Error != nil {
		logger.Error(res.Error)
	}
}

// Success logs successful operation
func Success(ctx *gin.Context, module string, operation string, detail string) {
	userID := getUserID(ctx)
	clientIP := utils.GetClientIP(ctx)
	Log(userID, clientIP, module, operation, 1, detail)
}

// Fail logs failed operation
func Fail(ctx *gin.Context, module string, operation, detail string) {
	userID := getUserID(ctx)
	clientIP := utils.GetClientIP(ctx)
	Log(userID, clientIP, module, operation, 0, detail)
}

func getUserID(ctx *gin.Context) int {
	claims := jwt.ExtractClaims(ctx)
	userID := claims[middlewares.IdentityKey].(string)
	id, _ := strconv.Atoi(userID)
	return id
}
