package core

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strconv"
)

func Paginator(c *gin.Context) func(db *gorm.DB) *gorm.DB {
	page := getParam(c, "page", 1)
	pageSize := getParam(c, "pageSize", 10)

	return func(db *gorm.DB) *gorm.DB {
		return db.Limit(pageSize).Offset((page - 1) * pageSize)
	}
}

func getParam(c *gin.Context, key string, defaultValue int) int {
	var result = defaultValue
	if param, exists := c.GetQuery(key); exists {
		var err error
		result, err = strconv.Atoi(param)
		if err != nil {
			result = defaultValue
		}
	}

	return result
}
