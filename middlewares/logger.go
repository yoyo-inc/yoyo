package middlewares

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yoyo-inc/yoyo/common/logger"
)

// Logger middleware
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		// 请求路径
		path := c.Request.URL.Path
		// 客户端IP
		clientIP := c.ClientIP()
		// 请求方式
		method := c.Request.Method

		c.Next()

		// 请求处理时间
		latency := time.Now().Sub(start)
		// 响应状态吗
		status := c.Writer.Status()
		// 错误信息
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		logger.Info(fmt.Sprintf("%15s [%d] %s %#v %v %s", clientIP, status, method, path, latency, errorMessage))
	}
}
