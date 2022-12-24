package db

import (
	"strings"

	"github.com/yoyo-inc/yoyo/common/config"
	"github.com/yoyo-inc/yoyo/common/dt"
	"github.com/yoyo-inc/yoyo/common/logger"
	"github.com/yoyo-inc/yoyo/resources"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	log "gorm.io/gorm/logger"
)

// Client provides the ability to manipulate database
var Client *gorm.DB

// AutoMigrateModels contains all models which to be migrated
var (
	AutoMigrateModels  []interface{}
	AutoMigrateMethods []func(db *gorm.DB)
)

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

func RunMigrateModels(models []interface{}) bool {
	if err := Client.
		Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").
		AutoMigrate(
			models...,
		); err != nil {
		logger.Error(err)
		return false
	}
	return true
}

func RunMigrateMethods(methods []func(db *gorm.DB)) {
	for _, method := range methods {
		method(Client)
	}
}

func MigrateFromFile() bool {
	sql := string(resources.DefaultSql)
	blocks := strings.Split(sql, ";")
	for _, block := range blocks {
		s := strings.Trim(block, "\n")
		s = strings.Trim(s, "\r\n")
		if s == "" {
			continue
		}

		if res := Client.Exec(s); res.Error != nil {
			logger.Error(res.Error)
			return false
		}
	}
	return true
}

type AutoMigrateModel struct {
	// 主键
	ID int `json:"id" gorm:"primarykey;autoIncreatment;comment:主键"`
	// 创建时间
	CreateTime *dt.LocalTime `json:"createTime" gorm:"type:timestamp;default:current_timestamp;<-:create;comment:创建时间"`
}

func (AutoMigrateModel) TableName() string {
	return "_automigrate"
}

// TODO:compare version
func canMigrate() bool {
	var migrateCount int64
	if res := Client.Model(&AutoMigrateModel{}).Count(&migrateCount); res.Error != nil {
		logger.Error(res.Error)
		return false
	}

	return migrateCount == 0
}

func recordMigrate() {
	if res := Client.Model(&AutoMigrateModel{}).Create(&AutoMigrateModel{}); res.Error != nil {
		logger.Error(res.Error)
	}
}

func autoMigrate() {
	RunMigrateModels([]interface{}{&AutoMigrateModel{}})

	logger.Info("Begin to autoMigrate")
	// migrate models
	if !RunMigrateModels(AutoMigrateModels) {
		return
	}
	if !canMigrate() {
		return
	}
	// migrate methods
	RunMigrateMethods(AutoMigrateMethods)
	// migrate from file
	if !MigrateFromFile() {
		return
	}

	recordMigrate()
}
