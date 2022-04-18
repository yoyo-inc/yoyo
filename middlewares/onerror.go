package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ypli0629/yoyo/core"
)

func OnError() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		err := c.Errors.Last()
		if err != nil {
			return
		}

		if e, ok := err.Err.(core.BusinessError); ok {
			c.AbortWithStatusJSON(http.StatusOK, core.Fail(e.Code, e.Message))
		}
	}
}
