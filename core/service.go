package core

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strconv"
)

func Paginator(c *gin.Context) func(db *gorm.DB) *gorm.DB {
	page, _ := strconv.Atoi(getParam(c, "current", "1"))
	pageSize, _ := strconv.Atoi(getParam(c, "pageSize", "10"))

	return func(db *gorm.DB) *gorm.DB {
		return db.Limit(pageSize).Offset((page - 1) * pageSize)
	}
}

func getParam(c *gin.Context, key string, defaultValue string) string {
	var result = defaultValue
	if param, exists := c.GetQuery(key); exists {
		result = param
	}
	return result
}
