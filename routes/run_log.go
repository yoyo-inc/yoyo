package routes

import (
	"github.com/gin-gonic/gin"
)

type runLog struct{}

// QueryRunLogVO
//	@Summary	查询运行日志列表
//	@Tags		runLog
//	@Produce	json
//	@Param		query	query		models.Pagination	false	"参数"
//	@Success	200		{object}	core.Response{data=core.PaginatedData{list=[]object}}
//	@Security	JWT
//	@Router		/runlogs [get]
func (*runLog) QueryRunLogs(c *gin.Context) {
}

func (*runLog) DownloadRunLog(c *gin.Context) {}

func (rl *runLog) Setup(r *gin.RouterGroup) {
	r.GET("/runlogs", rl.QueryRunLogs).GET("/runlog/download", rl.DownloadRunLog)
}
