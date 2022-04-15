package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ypli0629/yoyo/common/config"
	"github.com/ypli0629/yoyo/common/logger"
	"github.com/ypli0629/yoyo/middlewares"
)

func main() {
	// config
	config.Setup()
	// logger
	logger.Setup(logger.Options{
		Service: config.Config.Name,
	})

	r := gin.New()
	// middlewares
	// logger
	r.Use(middlewares.Logger())
	// recovery
	r.Use(gin.Recovery())

	r.Run()
}
