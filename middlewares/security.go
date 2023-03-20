package middlewares

import (
	"errors"
	"gorm.io/gorm"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	jwt "github.com/yoyo-inc/gin-jwt/v3"
	"github.com/yoyo-inc/yoyo/common/db"
	"github.com/yoyo-inc/yoyo/common/logger"
	"github.com/yoyo-inc/yoyo/core"
	"github.com/yoyo-inc/yoyo/models"
	"github.com/yoyo-inc/yoyo/services"
)

var (
	// SecurityMiddleware security middleware
	SecurityMiddleware *jwt.GinJWTMiddleware
	IdentityKey        = "userID"
	DefaultTimeout     = 1440 * time.Minute
)

// Security setups security
func Security() func() gin.HandlerFunc {
	var err error
	SecurityMiddleware, err = jwt.New(&jwt.GinJWTMiddleware{
		Realm:   "yoyo",
		Key:     []byte{},
		Timeout: DefaultTimeout,
		GetTimeout: func() time.Duration {
			// get token timeout from system setting
			var systemSecurity models.SystemSecurity
			if res := db.Client.Model(&models.SystemSecurity{}).First(&systemSecurity); res.Error != nil {
				if !errors.Is(res.Error, gorm.ErrRecordNotFound) {
					logger.Error(res.Error)
				}
				return DefaultTimeout
			}

			if systemSecurity.LoginExpireEnable {
				return time.Duration(systemSecurity.LoginExpireTime) * time.Minute
			} else {
				return DefaultTimeout
			}
		},
		MaxRefresh:  30 * 24 * time.Hour,
		IdentityKey: IdentityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(models.User); ok {
				return jwt.MapClaims{
					IdentityKey: strconv.Itoa(v.ID),
				}
			}

			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return claims[IdentityKey]
		},
		TokenLookup: "header: Authorization, query: token",
		Authenticator: func(c *gin.Context) (interface{}, error) {
			user, err := services.DoLogin(c)
			if err != nil {
				logger.Errorf("%s: %s", err, user.Username)
				return nil, err
			}

			return user, nil
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, core.NewFailedResponse(strconv.Itoa(code), message))
		},
		LoginResponse: func(c *gin.Context, code int, token string, expire time.Time) {
			core.OK(c, map[string]interface{}{
				"token":  token,
				"expire": expire.Format(time.RFC3339),
			})
		},
		LogoutResponse: func(c *gin.Context, code int) {
			core.OK(c, true)
		},
	})

	if err != nil {
		logger.Panicf("Failed to setup security: %s", err)
	}

	return SecurityMiddleware.MiddlewareFunc
}
