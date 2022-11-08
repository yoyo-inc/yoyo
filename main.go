package main

import (
	"github.com/gin-gonic/gin"
	"github.com/yoyo-inc/yoyo/common/config"
	"github.com/yoyo-inc/yoyo/common/db"
	"github.com/yoyo-inc/yoyo/common/hub"
	"github.com/yoyo-inc/yoyo/common/i18n"
	"github.com/yoyo-inc/yoyo/common/logger"
	"github.com/yoyo-inc/yoyo/common/swag"
	"github.com/yoyo-inc/yoyo/middlewares"
	"github.com/yoyo-inc/yoyo/routes"
)

func ReadLoggerConfig() logger.Options {
	var options logger.Options
	options.Service = config.GetString("name")

	options.Level = config.GetString("logger.level")

	// logger
	// syslog
	if config.Get("logger.syslog") != nil {
		var syslog logger.SyslogOption
		err := config.UnmarshalKey("logger.syslog", &syslog)
		if err != nil {
			logger.Error(err)
		}
		options.Syslog = &syslog
	}
	// rollingFile
	if config.Get("logger.rolling_file") != nil {
		var rollingFileOption logger.RollingFileOption
		err := config.UnmarshalKey("logger.rolling_file", &rollingFileOption)
		if err != nil {
			logger.Error(err)
		}
		options.RollingFile = &rollingFileOption
	}

	return options
}

// @Title yoyo
// @Description yoyo
// @Host 127.0.0.1:8080
// @BasePath /API
// @SecurityDefinitions.apikey JWT
// @In header
// @Name Authorization
func main() {
	// config
	config.Setup()
	// logger
	logger.Setup(ReadLoggerConfig())
	// db
	db.Setup()
	// i18n
	i18n.Setup()
	// hub
	hub.Setup()

	engine := gin.New()
	// middlewares
	// logger
	engine.Use(middlewares.Logger())
	// onerror
	engine.Use(middlewares.OnError())
	// recovery
	engine.Use(gin.Recovery())
	// security
	engine.Use(middlewares.Security()())

	// base_path
	r := &engine.RouterGroup
	if config.GetString("server.base_path") != "" {
		r = engine.Group(config.GetString("server.base_path"))
	}

	swag.Setup(r)
	routes.Setup(r)

	address := config.GetString("server.host") + ":" + config.GetString("server.port")
	logger.Infof("Listen and Serving HTTP on http://%s", address)

	// startup
	err := engine.Run(address)
	if err != nil {
		logger.Error(err)
		return
	}
}
