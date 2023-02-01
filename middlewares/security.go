package middlewares

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	jwt "github.com/yoyo-inc/gin-jwt/v3"
	"github.com/yoyo-inc/yoyo/common/logger"
	"github.com/yoyo-inc/yoyo/core"
	"github.com/yoyo-inc/yoyo/errs"
	"github.com/yoyo-inc/yoyo/models"
	"github.com/yoyo-inc/yoyo/services"
)

var (
	// SecurityMiddleware security middleware
	SecurityMiddleware *jwt.GinJWTMiddleware
	IdentityKey        = "userID"
)

type loginPayload struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// Security setups security
func Security() func() gin.HandlerFunc {
	var err error
	SecurityMiddleware, err = jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "yoyo",
		Key:         []byte{},
		Timeout:     30 * 24 * time.Hour,
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
			var payload loginPayload
			if err := c.ShouldBindJSON(&payload); err != nil {
				logger.Error(err)
				return nil, errs.ErrUsernameOrPassword
			}

			user, err := services.DoLogin(c, payload.Username, payload.Password)
			if err != nil {
				logger.Errorf("%s: %s", err, payload.Username)
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
