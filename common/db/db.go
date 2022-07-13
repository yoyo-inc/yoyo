package db

import (
	"github.com/yoyo-inc/yoyo/common/config"
	"github.com/yoyo-inc/yoyo/common/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Client provides the ability to manipulate database
var Client *gorm.DB

//AutoMigrateModels contains all models which to be migrated
var AutoMigrateModels []interface{}
var AutoMigrateMethods []func(db *gorm.DB)

// Setup setups database connection
func Setup() {
	dns := config.Config.DB.DNS
	var err error
	Client, err = gorm.Open(mysql.Open(dns), &gorm.Config{})

	if err != nil {
		logger.Panicf("Failed to connect database: %s", err)
	} else {
		logger.Info("Connect database successfully")
	}

	logger.Info("Begin to automigrate")
	automigrate()
}

func AddAutoMigrateModel(model interface{}) {
	AutoMigrateModels = append(AutoMigrateModels, model)
}

func AddAutoMigrateMethods(method func(client *gorm.DB)) {
	AutoMigrateMethods = append(AutoMigrateMethods, method)
}

func automigrate() {
	if err := Client.
		Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").
		AutoMigrate(
			AutoMigrateModels...,
		); err != nil {
		logger.Error(err)
	}

	for _, method := range AutoMigrateMethods {
		method(Client)
	}
}
