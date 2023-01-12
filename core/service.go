package core

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Paginator(c *gin.Context) func(db *gorm.DB) *gorm.DB {
	page, _ := strconv.Atoi(GetParam(c, "current", "1"))
	pageSize, _ := strconv.Atoi(GetParam(c, "pageSize", "10"))

	return func(db *gorm.DB) *gorm.DB {
		return db.Limit(pageSize).Offset((page - 1) * pageSize)
	}
}

func DateTimeRanger(c *gin.Context, fieldName string) func(db *gorm.DB) *gorm.DB {
	startTime := GetParam(c, "startTime", "")
	endTime := GetParam(c, "endTime", "")

	if startTime == "" || endTime == "" {
		return func(db *gorm.DB) *gorm.DB {
			return db
		}
	}

	if fieldName == "" {
		fieldName = "create_time"
	}

	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("%s between ? and ?", fieldName), startTime, endTime)
	}
}

func Orderer(c *gin.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Order("create_time desc")
	}
}

func GetParam(c *gin.Context, key string, defaultValue string) string {
	result := defaultValue
	if param, exists := c.GetQuery(key); exists {
		result = param
	}
	return result
}
