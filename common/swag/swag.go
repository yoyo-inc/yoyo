package swag

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/yoyo-inc/yoyo/common/config"
	"github.com/yoyo-inc/yoyo/docs"
)

func Setup(r *gin.RouterGroup) {
	if config.GetString("name") != "" {
		docs.SwaggerInfo.Title = config.GetString("name")
	} else {
		docs.SwaggerInfo.Title = "yoyo"
	}
	docs.SwaggerInfo.Version = "1.0.0"
	docs.SwaggerInfo.BasePath = config.GetString("server.base_path")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler, ginSwagger.PersistAuthorization(true)))
}
