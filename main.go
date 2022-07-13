package main

import (
	"github.com/gin-gonic/gin"
	"github.com/yoyo-inc/yoyo/common/config"
	"github.com/yoyo-inc/yoyo/common/db"
	"github.com/yoyo-inc/yoyo/common/i18n"
	"github.com/yoyo-inc/yoyo/common/logger"
	"github.com/yoyo-inc/yoyo/common/swag"
	"github.com/yoyo-inc/yoyo/middlewares"
	"github.com/yoyo-inc/yoyo/routes"
)

func main() {
	// config
	config.Setup()
	// logger
	logger.Setup(logger.Options{
		Service: config.Config.Name,
	})
	// db
	db.Setup()
	// i18n
	i18n.Setup()

	engine := gin.New()
	// middlewares
	// logger
	engine.Use(middlewares.Logger())
	// onerror
	engine.Use(middlewares.OnError())
	// recovery
	engine.Use(gin.Recovery())
	// security
	// middlewares.Setup()

	r := &engine.RouterGroup
	if config.Config.Server.BasePath != "" {
		r = engine.Group(config.Config.Server.BasePath)
	}

	routes.SetupNoSecurity(r)
	// security
	// engine.Use(middlewares.SecurityMiddleware.MiddlewareFunc())
	routes.SetupSecurity(r)

	swag.Setup(r, config.Config)

	address := config.Config.Server.Host + ":" + config.Config.Server.Port
	logger.Infof("Listen and Serving HTTP on http://%s", address)

	// startup
	err := engine.Run(address)
	if err != nil {
		logger.Error(err)
		return
	}
}
