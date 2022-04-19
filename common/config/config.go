// Package config provide read custom configuration from file or environment variables
// Example:
//		server:
//			host: '127.0.0.1'
//			port:	'8080'
//			basepath: ''
//    db:
//			dns: root:123456@tcp(127.0.0.1:3306)/yoyo?charset=utf8mb4&parseTime=True&loc=Local
package config

import (
	"fmt"

	"github.com/spf13/viper"
	"github.com/ypli0629/yoyo/common/logger"
)

// Config contains user custom configuration
var Config Application

// DB contains database connection
type DB struct {
	DNS string
}

// Server contains server configuration
type Server struct {
	Host     string
	Port     string
	BasePath string `mapstructure:"base_path"`
}

// Application contains application configuration
type Application struct {
	Name   string
	Server Server
	DB     DB
}

// Setup setups config parser
func Setup() {
	// application config file name: application.yml
	viper.SetConfigName("application")
	viper.SetConfigType("yaml")
	// also read config from environment variables
	viper.AutomaticEnv()
	// lookup config file from current directory
	viper.AddConfigPath(".")
	// default config
	setDefaultConfig()

	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			logger.Warn("Config file not found")
			return
		}
		panic(fmt.Errorf("Fatal to read config file: %w", err))
	}

	err = viper.Unmarshal(&Config)
	if err != nil {
		panic(fmt.Errorf("Fatal to unmarshal config: %w", err))
	}

	logger.Info("Load config file successfully")
}

func setDefaultConfig() {
	viper.SetDefault("name", "default")
	viper.SetDefault("server.host", "127.0.0.1")
	viper.SetDefault("server.port", "8080")
}
