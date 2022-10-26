package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/yoyo-inc/yoyo/common/logger"
	"github.com/yoyo-inc/yoyo/core"
	"github.com/yoyo-inc/yoyo/errs"
	"github.com/yoyo-inc/yoyo/models"
)

type auditLogController struct {
}

// QueryAuditLog
// @Summary  查询操作日志列表
// @Tags     auditLog
// @Produce  json
// @Param    query  query     models.Pagination  false  "参数"
// @Success  200    {object}  core.Response{data=core.PaginatedData{list=[]models.AuditLog}}
// @Router   /audit_logs [get]
func (*auditLogController) QueryAuditLog(c *gin.Context) {
	var auditLogs []models.AuditLog
	quires := core.GetPaginatedQuery(&models.AuditLog{})
	if res := quires[0].Scopes(core.Paginator(c)).Find(&auditLogs); res.Error != nil {
		logger.Error(res.Error)
		_ = c.Error(errs.ErrQueryAuditLog)
		return
	}
	var total int64
	if res := quires[1].Count(&total); res.Error != nil {
		logger.Error(res.Error)
		_ = c.Error(errs.ErrQueryAuditLog)
		return
	}
	core.OK(c, core.Paginated(quires, total))
}

func (alc *auditLogController) Setup(r *gin.RouterGroup) {
	r.GET("/audit_logs", alc.QueryAuditLog)
}
