package routes

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/gin-gonic/gin"
	"github.com/yoyo-inc/yoyo/common/config"
	"github.com/yoyo-inc/yoyo/common/logger"
	"github.com/yoyo-inc/yoyo/core"
	"github.com/yoyo-inc/yoyo/errs"
	"github.com/yoyo-inc/yoyo/services"
	"github.com/yoyo-inc/yoyo/vo"
)

type runLogController struct{}

// QueryRunLogVO
//
//	@Summary	查询运行日志列表
//	@Tags		runLog
//	@Produce	json
//	@Param		query	query		models.Pagination	false	"参数"
//	@Success	200		{object}	core.Response{data=core.PaginatedData{list=[]vo.RunLogVO}}
//	@Security	JWT
//	@Router		/runlogs [get]
func (*runLogController) QueryRunLogs(c *gin.Context) {
	srvs := config.GetMap("service_logs")

	runlogs := make([]vo.RunLogVO, 0)
	for s, v := range srvs {
		p, err := filepath.Abs(v)
		if err != nil {
			logger.Error(err)
			continue
		}

		stats, err := services.ScanLogByRecent(p, 0)
		if err != nil {
			logger.Error(err)
			continue
		}

		runlogs = append(runlogs, slice.Map(stats, func(index int, stat services.FileStat) vo.RunLogVO {
			return vo.RunLogVO{
				Service:  s,
				Filename: stat.Filename,
				Date:     stat.Date,
				Filesize: fmt.Sprintf("%.2f", float64(stat.Filesize/1024)),
			}
		})...)
	}

	err := slice.SortByField(runlogs, "Date", "desc")
	if err != nil {
		logger.Error(err)
		c.Error(errs.ErrQueryRunLog)
		return
	}

	page, _ := strconv.Atoi(core.GetParam(c, "current", "1"))
	pageSize, _ := strconv.Atoi(core.GetParam(c, "pageSize", "10"))

	start := (page - 1) * pageSize
	end := page * pageSize
	length := len(runlogs)

	res := make([]vo.RunLogVO, 0)
	if start <= length && length < end {
		res = runlogs[start:(length - start)]
	} else if length >= end {
		res = runlogs[start:end]
	}

	core.OK(c, core.Paginated(res, int64(len(runlogs))))
}

// DownloadRunLog
//	@Summary	下载运行日志
//	@Tags		runlog
//	@Accept		json
//	@Produce	octet-stream
//	@Param		query	query	vo.DownloadRunLogVO	false	"参数"
//	@Security	JWT
//	@Router		/runlog/download [get]
func (*runLogController) DownloadRunLog(c *gin.Context) {
	var query vo.DownloadRunLogVO
	if err := c.ShouldBindQuery(&query); err != nil {
		logger.Error(err)
		c.Error(core.NewParameterError(err))
		return
	}

	services := config.GetMap("service_logs")
	if v, ok := services[query.Service]; !ok {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		filePath, err := filepath.Abs(filepath.Join(v, query.Filename))
		if err != nil {
			c.AbortWithStatus(http.StatusNotFound)
		} else {
			c.FileAttachment(filePath, query.Filename)
		}
	}
}

func (rl *runLogController) Setup(r *gin.RouterGroup) {
	r.GET("/runlogs", rl.QueryRunLogs).GET("/runlog/download", rl.DownloadRunLog)
}
