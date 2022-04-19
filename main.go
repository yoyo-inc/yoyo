package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ypli0629/yoyo/common/config"
	"github.com/ypli0629/yoyo/common/db"
	"github.com/ypli0629/yoyo/common/i18n"
	"github.com/ypli0629/yoyo/common/logger"
	"github.com/ypli0629/yoyo/middlewares"
	"github.com/ypli0629/yoyo/routes"
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
	middlewares.Setup()

	r := &engine.RouterGroup
	if config.Config.Server.BasePath != "" {
		r = engine.Group(config.Config.Server.BasePath)
	}

	routes.SetupNoSecurity(r)
	// security
	engine.Use(middlewares.SecurityMiddleware.MiddlewareFunc())
	routes.SetupSecurity(r)

	address := config.Config.Server.Host + ":" + config.Config.Server.Port
	logger.Infof("Listen and Serving HTTP on %s", address)

	// startup
	engine.Run(address)
}
