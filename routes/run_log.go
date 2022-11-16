package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/yoyo-inc/yoyo/common/logger"
	"github.com/yoyo-inc/yoyo/core"
	"github.com/yoyo-inc/yoyo/errs"
	"github.com/yoyo-inc/yoyo/models"
	"github.com/yoyo-inc/yoyo/vo"
)

type runLog struct{}

// QueryRunLogVO
// @Summary  查询运行日志列表
// @Tags     runLog
// @Produce  json
// @Param    query query    vo.QueryRunLogVO  false "参数"
// @Param    query query    models.Pagination false "参数"
// @Success  200   {object} core.Response{data=core.PaginatedData{list=[]models.RunLog}}
// @Security JWT
// @Router   /runlogs [get]
func (*runLog) QueryRunLogs(c *gin.Context) {
	var query vo.QueryRunLogVO
	if err := c.ShouldBindQuery(&query); err != nil {
		logger.Error(err)
		c.Error(core.NewParameterError(err))
		return
	}

	queries := core.GetPaginatedQuery(&models.RunLog{})
	for i := range queries {
		if query.Filename != "" {
			queries[i].Where("filename like ?", "%"+query.Filename+"%")
		}
	}

	var runLogs []models.RunLog
	if res := queries[0].Scopes(core.Paginator(c)).Find(&runLogs); res.Error != nil {
		logger.Error(res.Error)
		c.Error(errs.ErrQueryRunLog)
		return
	}

	var count int64
	if res := queries[1].Count(&count); res.Error != nil {
		logger.Error(res.Error)
		c.Error(errs.ErrQueryRunLog)
		return
	}

	core.OK(c, core.Paginated(runLogs, count))
}

func (*runLog) DownloadRunLog(c *gin.Context) {}

func (rl *runLog) Setup(r *gin.RouterGroup) {
	r.GET("/runlogs", rl.QueryRunLogs).GET("/runlog/download", rl.DownloadRunLog)
}
