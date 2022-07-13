package middlewares

import (
	"github.com/yoyo-inc/yoyo/common/logger"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yoyo-inc/yoyo/core"
)

// OnError handles runtime error
func OnError() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		err := c.Errors.Last()
		if err == nil {
			return
		}

		if e, ok := err.Err.(core.ParameterError); ok {
			logger.Error(e)
			c.AbortWithStatusJSON(http.StatusBadRequest, core.FailedResponse(strconv.Itoa(http.StatusBadRequest), e.Error()))
		}
		if e, ok := err.Err.(core.BusinessError); ok {
			c.AbortWithStatusJSON(http.StatusOK, core.FailedResponse(e.Code, e.Message))
		}
	}
}
