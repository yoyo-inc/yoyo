package core

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Paginator(c *gin.Context) func(db *gorm.DB) *gorm.DB {
	page, _ := strconv.Atoi(getParam(c, "current", "1"))
	pageSize, _ := strconv.Atoi(getParam(c, "pageSize", "10"))

	return func(db *gorm.DB) *gorm.DB {
		return db.Limit(pageSize).Offset((page - 1) * pageSize)
	}
}

func DateTimeRanger(c *gin.Context) func(db *gorm.DB) *gorm.DB {
	startTime := getParam(c, "startTime", "")
	endTime := getParam(c, "endTime", "")

	if startTime == "" || endTime == "" {
		return func(db *gorm.DB) *gorm.DB {
			return db
		}
	}

	return func(db *gorm.DB) *gorm.DB {
		return db.Where("create_time between ? and ?", startTime, endTime)
	}
}

func getParam(c *gin.Context, key string, defaultValue string) string {
	var result = defaultValue
	if param, exists := c.GetQuery(key); exists {
		result = param
	}
	return result
}
