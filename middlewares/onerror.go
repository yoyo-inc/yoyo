package middlewares

import (
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

		switch e := err.Err.(type) {
		case core.ParameterError:
			c.AbortWithStatusJSON(http.StatusBadRequest, core.NewFailedResponse(strconv.Itoa(http.StatusBadRequest), e.Error()))
		case core.BusinessError:
			c.AbortWithStatusJSON(http.StatusOK, core.NewFailedResponse(e.Code, e.Message))
		default:
			c.AbortWithError(http.StatusInternalServerError, e)
		}
	}
}
