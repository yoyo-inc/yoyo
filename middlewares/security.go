package middlewares

import (
	"strconv"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/yoyo-inc/yoyo/common/logger"
	"github.com/yoyo-inc/yoyo/core"
	"github.com/yoyo-inc/yoyo/errs"
	"github.com/yoyo-inc/yoyo/models"
	"github.com/yoyo-inc/yoyo/services"
)

var (
	// SecurityMiddleware security middleware
	SecurityMiddleware *jwt.GinJWTMiddleware
	identityKey        = "userID"
)

type loginPayload struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// Setup setups security
func Setup() {
	var err error
	SecurityMiddleware, err = jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "yoyo",
		Key:         []byte{},
		Timeout:     30 * 24 * time.Hour,
		MaxRefresh:  30 * 24 * time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(models.User); ok {
				return jwt.MapClaims{
					identityKey: v.ID,
				}
			}

			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return claims[identityKey]
		},
		TokenLookup: "header: Authorization, query: token",
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var payload loginPayload
			if err := c.ShouldBindJSON(&payload); err != nil {
				return nil, errs.ErrUsernameOrPassword
			}

			user, err := services.DoLogin(payload.Username, payload.Password)
			if err != nil {
				return nil, err
			}

			return user, nil
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, core.FailedResponse(strconv.Itoa(code), message))
		},
		LoginResponse: func(c *gin.Context, code int, token string, expire time.Time) {
			core.OK(c, map[string]interface{}{
				"token":  token,
				"expire": expire.Format(time.RFC3339),
			})
		},
	})

	if err != nil {
		logger.Panicf("Failed to setup security: %s", err)
	}
}
