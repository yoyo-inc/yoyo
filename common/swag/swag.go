package swag

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/ypli0629/yoyo/common/config"
	"github.com/ypli0629/yoyo/docs"
)

func Setup(r *gin.RouterGroup, config config.Application) {
	if config.Name != "" {
		docs.SwaggerInfo.Title = config.Name
	} else {
		docs.SwaggerInfo.Title = "yoyo"
	}
	docs.SwaggerInfo.Version = "1.0.0"
	docs.SwaggerInfo.BasePath = config.Server.BasePath
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
