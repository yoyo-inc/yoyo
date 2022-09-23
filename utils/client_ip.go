package utils

import "github.com/gin-gonic/gin"

func GetClientIP(c *gin.Context) string {
	var ip string
	xRealIP := c.GetHeader("X-Real-IP")
	if xRealIP != "" {
		ip = xRealIP
	} else {
		ip = c.ClientIP()
	}
	return ip
}
