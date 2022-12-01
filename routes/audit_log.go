package routes

import (
	"github.com/duke-git/lancet/v2/slice"
	"github.com/gin-gonic/gin"
	"github.com/yoyo-inc/yoyo/common/db"
	"github.com/yoyo-inc/yoyo/common/logger"
	"github.com/yoyo-inc/yoyo/core"
	"github.com/yoyo-inc/yoyo/errs"
	"github.com/yoyo-inc/yoyo/models"
	"github.com/yoyo-inc/yoyo/vo"
)

type auditLogController struct {
}

// QueryAuditLog
// @Summary  查询审计日志列表
// @Tags     auditLog
// @Produce  json
// @Param    query query    models.Pagination false "参数"
// @Success  200   {object} core.Response{data=core.PaginatedData{list=[]models.AuditLog}}
// @Security JWT
// @Router   /audit_logs [get]
func (*auditLogController) QueryAuditLog(c *gin.Context) {
	var query vo.QueryAuditLogVO
	if err := c.ShouldBindQuery(&query); err != nil {
		logger.Error(err)
		c.Error(core.NewParameterError(err))
		return
	}

	var auditLogs []models.AuditLog
	queries := core.GetPaginatedQuery(&models.AuditLog{})
	if res := queries[0].Preload("User").Scopes(core.Paginator(c)).Where(query).Order("create_time desc").Find(&auditLogs); res.Error != nil {
		logger.Error(res.Error)
		c.Error(errs.ErrQueryAuditLog)
		return
	}
	var total int64
	if res := queries[1].Where(query).Count(&total); res.Error != nil {
		logger.Error(res.Error)
		c.Error(errs.ErrQueryAuditLog)
		return
	}
	core.OK(c, core.Paginated(auditLogs, total))
}

// QueryAuditLogModules
// @Summary  查询审计日志列表
// @Tags     auditLog
// @Produce  json
// @success  200 {object} core.Response{data=array,string}}
// @security jwt
// @Router   /audit_log/modules [get]
func (*auditLogController) QueryAuditLogModules(c *gin.Context) {
	var auditLogs []models.AuditLog
	if res := db.Client.Model(&models.AuditLog{}).Select("module").Distinct("module").Find(&auditLogs); res.Error != nil {
		logger.Error(res.Error)
		c.Error(errs.ErrQueryAuditLogModule)
		return
	}

	modules := slice.Map(auditLogs, func(_ int, item models.AuditLog) string {
		return item.Module
	})

	core.OK(c, modules)
}

func (alc *auditLogController) Setup(r *gin.RouterGroup) {
	r.GET("/audit_logs", alc.QueryAuditLog).GET("/audit_log/modules", alc.QueryAuditLogModules)
}
