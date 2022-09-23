package core

import (
	"github.com/gin-gonic/gin"
	"github.com/yoyo-inc/yoyo/common/db"
	"gorm.io/gorm"
)

type Controller interface {
	Setup(r *gin.RouterGroup)
}

func GetPaginatedQuery(model interface{}) [2]*gorm.DB {
	modelQuery := db.Client.Model(model)
	totalQuery := db.Client.Model(model)
	return [2]*gorm.DB{modelQuery, totalQuery}
}
