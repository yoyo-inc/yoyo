package db

import (
	"github.com/yoyo-inc/yoyo/common/config"
	"github.com/yoyo-inc/yoyo/common/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	log "gorm.io/gorm/logger"
)

// Client provides the ability to manipulate database
var Client *gorm.DB

//AutoMigrateModels contains all models which to be migrated
var AutoMigrateModels []interface{}
var AutoMigrateMethods []func(db *gorm.DB)

// Setup setups database connection
func Setup() {
	dns := config.GetString("db.dns")
	var err error
	Client, err = gorm.Open(mysql.Open(dns), &gorm.Config{
		Logger: log.Default.LogMode(GetLoggerLevel(config.GetString("db.logger.level"))),
	})

	if err != nil {
		logger.Panicf("Failed to connect database: %s", err)
	} else {
		logger.Info("Connect database successfully")
	}

	logger.Info("Begin to autoMigrate")
	autoMigrate()
}

func GetLoggerLevel(level string) log.LogLevel {
	switch level {
	case "silent":
		return log.Silent
	case "error":
		return log.Error
	case "warn":
		return log.Warn
	case "info":
		return log.Info
	default:
		return log.Warn
	}
}

func AddAutoMigrateModel(model interface{}) {
	AutoMigrateModels = append(AutoMigrateModels, model)
}

func AddAutoMigrateMethods(method func(client *gorm.DB)) {
	AutoMigrateMethods = append(AutoMigrateMethods, method)
}

func MigrateModels(models []interface{}) {
	if err := Client.
		Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").
		AutoMigrate(
			models...,
		); err != nil {
		logger.Error(err)
	}
}

func MigrateMethods(methods []func(db *gorm.DB)) {
	for _, method := range methods {
		method(Client)
	}
}

func autoMigrate() {
	MigrateModels(AutoMigrateModels)
	MigrateMethods(AutoMigrateMethods)
}
