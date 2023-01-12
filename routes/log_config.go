package routes

import (
	"errors"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/yoyo-inc/yoyo/common/config"
	"github.com/yoyo-inc/yoyo/common/db"
	"github.com/yoyo-inc/yoyo/common/logger"
	"github.com/yoyo-inc/yoyo/core"
	"github.com/yoyo-inc/yoyo/models"
	"github.com/yoyo-inc/yoyo/services"
	"gorm.io/gorm"
)

type logConfigController struct{}

// QueryLogConfig
//
//	@Summary	查询日志配置
//	@Tags		logConfig
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	core.Response{data=models.LogConfig}
//	@Security	JWT
//	@Router		/log_config [get]
func (*logConfigController) QueryLogConfig(c *gin.Context) {
	var config models.LogConfig

	if res := db.Client.Model(&models.LogConfig{}).Find(&config); res.Error != nil {
		logger.Error(res.Error)
		return
	}

	core.OK(c, config)
}

// SaveLogConfig
//	@Summary	保存日志配置
//	@Tags		logConfig
//	@Accept		json
//	@Produce	json
//	@Param		body	body		models.LogConfig	true	"参数"
//	@Success	200		{object}	core.Response{data=bool}
//	@Security	JWT
//	@Router		/log_config [post]
func (lcc *logConfigController) SaveLogConfig(c *gin.Context) {
	var logConfig models.LogConfig

	if logConfig.ID != "" {
		if res := db.Client.Omit("id").Where("id = ?", logConfig.ID).Updates(&logConfig); res.Error != nil {
			logger.Error(res.Error)
			return
		}
	} else {
		if res := db.Client.Create(&logConfig); res.Error != nil {
			logger.Error(res.Error)
			return
		}
	}

	lcc.registerLogSchedJob()

	core.OK(c, true)
}

func (*logConfigController) registerLogSchedJob() {
	var logConfig models.LogConfig

	if res := db.Client.Model(&models.LogConfig{}).First(&logConfig); res.Error != nil {
		if !errors.Is(res.Error, gorm.ErrRecordNotFound) {
			logger.Error(res.Error)
		}
		return
	}

	if logConfig.KeepTime != nil {
		services.AddSchedJob("deleteLog", "log", "定时清除日志文件", "0 1 * * *", func() error {
			srvs := config.GetMap("service_logs")
			for k, v := range srvs {
				p, err := filepath.Abs(v)
				if err != nil {
					logger.Errorf("清除%s(%s)失败: %s", k, v, err)
					continue
				}
				if err := services.DeleteLogByRecent(p, *logConfig.KeepTime); err != nil {
					logger.Errorf("清除%s(%s)失败: %s", k, v, err)
				}
			}
			return nil
		})
	}

	if logConfig.Archive {
		services.AddSchedJob("archiveLog", "log", "定时归档日志文件", "0 1 * * *", func() error {
			srvs := config.GetMap("service_logs")
			for k, v := range srvs {
				p, err := filepath.Abs(v)
				if err != nil {
					logger.Errorf("清除%s(%s)失败: %s", k, v, err)
					continue
				}
				if err := services.ArchiveLogByRecent(p, 1); err != nil {
					logger.Errorf("归档%s(%s)失败: %s", k, v, err)
				}
			}
			return nil
		})
	}
}

func (lcc *logConfigController) Setup(r *gin.RouterGroup) {
	r.GET("/log_config", lcc.QueryLogConfig).
		POST("/log_config", lcc.SaveLogConfig)

	lcc.registerLogSchedJob()
}
