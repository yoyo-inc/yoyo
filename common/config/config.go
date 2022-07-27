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
	"github.com/yoyo-inc/yoyo/common/logger"
)

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

	logger.Info("Load config file successfully")
}

func setDefaultConfig() {
	viper.SetDefault("name", "default")
	viper.SetDefault("server.host", "127.0.0.1")
	viper.SetDefault("server.port", "8080")
}

func Get(key string) interface{} {
	return viper.Get(key)
}

func GetString(key string) string {
	return viper.GetString(key)
}

func GetInt(key string) int {
	return viper.GetInt(key)
}

func GetSlice(key string) []interface{} {
	return GetSlice(key)
}

func GetIntSlice(key string) []int {
	return viper.GetIntSlice(key)
}

func GetStringSlice(key string) []string {
	return viper.GetStringSlice(key)
}

func GetMap(key string) map[string]string {
	return GetStringMapString(key)
}

func GetStringMapString(key string) map[string]string {
	return viper.GetStringMapString(key)
}

func GetStringMap(key string) map[string]interface{} {
	return viper.GetStringMap(key)
}
