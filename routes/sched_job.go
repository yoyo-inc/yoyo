package routes

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/yoyo-inc/yoyo/common/db"
	"github.com/yoyo-inc/yoyo/common/logger"
	"github.com/yoyo-inc/yoyo/core"
	"github.com/yoyo-inc/yoyo/errs"
	"github.com/yoyo-inc/yoyo/models"
	"github.com/yoyo-inc/yoyo/services"
	"github.com/yoyo-inc/yoyo/services/audit_log"
	"github.com/yoyo-inc/yoyo/vo"
)

type schedJobController struct{}

// QuerySchedJobs
//
//	@Summary	查询定时任务
//	@Tags		schedJob
//	@Accept		json
//	@Produce	json
//	@Param		query	query		models.Pagination	false	"参数"
//	@Success	200		{object}	core.Response{data=core.PaginatedData{list=[]models.SchedJob}}
//	@Security	JWT
//	@Router		/schedjobs [get]
func (*schedJobController) QuerySchedJobs(c *gin.Context) {
	var query vo.QuerySchedJobVO
	if err := c.ShouldBindQuery(&query); err != nil {
		logger.Error(err)
		_ = c.Error(core.NewParameterError(err))
		return
	}
	queries := core.GetPaginatedQuery(&models.SchedJob{})

	var jobs []models.SchedJob
	if res := queries[0].Scopes(core.Paginator(c), core.DateTimeRanger(c, "")).Where(query).Order("create_time desc").Find(&jobs); res.Error != nil {
		logger.Error(res.Error)
		_ = c.Error(errs.ErrQuerySchedJob)
		return
	}

	var count int64
	if res := queries[1].Where(query).Count(&count); res.Error != nil {
		logger.Error(res.Error)
		_ = c.Error(errs.ErrQuerySchedJob)
		return
	}

	core.OK(c, core.Paginated(jobs, count))
}

// StartSchedJobs
//
//	@Summary	开启定时任务
//	@Tags		schedJob
//	@Accept		json
//	@Produce	json
//	@Param		body	body		vo.UpdateSchedJobVO	true	"参数"
//	@Success	200		{object}	core.Response{data=bool}
//	@Security	JWT
//	@Router		/schedjob/open [put]
func (*schedJobController) StartSchedJobs(c *gin.Context) {
	var query vo.UpdateSchedJobVO
	if err := c.ShouldBindJSON(&query); err != nil {
		logger.Error(err)
		_ = c.Error(core.NewParameterError(err))
		return
	}

	var job models.SchedJob
	if res := db.Client.Model(&models.SchedJob{}).Find(&job, "id = ?", query.ID); res.Error != nil {
		logger.Error(res.Error)
		_ = c.Error(errs.ErrNotExistSchedJob)
		return
	}

	if err := services.StartSchedJob(job.JobID); err != nil {
		audit_log.Fail(c, "定时任务", "开启", fmt.Sprintf("定时任务(%s)开启成功", job.Description))
		_ = c.Error(errs.ErrStartSchedJob)
		return
	}

	audit_log.Success(c, "定时任务", "开启", fmt.Sprintf("定时任务(%s)开启成功", job.Description))
	core.OK(c, true)
}

// CloseSchedJobs
//
//	@Summary	关闭定时任务
//	@Tags		schedJob
//	@Accept		json
//	@Produce	json
//	@Param		body	body		vo.UpdateSchedJobVO	true	"参数"
//	@Success	200		{object}	core.Response{data=bool}
//	@Security	JWT
//	@Router		/schedjob/close [put]
func (*schedJobController) CloseSchedJobs(c *gin.Context) {
	var query vo.UpdateSchedJobVO
	if err := c.ShouldBindJSON(&query); err != nil {
		logger.Error(err)
		_ = c.Error(core.NewParameterError(err))
		return
	}

	var job models.SchedJob
	if res := db.Client.Model(&models.SchedJob{}).Find(&job, "id = ?", query.ID); res.Error != nil {
		logger.Error(res.Error)
		_ = c.Error(errs.ErrNotExistSchedJob)
		return
	}

	if err := services.StopSchedJob(job.JobID); err != nil {
		audit_log.Fail(c, "定时任务", "关闭", fmt.Sprintf("定时任务(%s)关闭成功", job.Description))
		_ = c.Error(errs.ErrStopSchedJob)
		return
	}

	audit_log.Success(c, "定时任务", "关闭", fmt.Sprintf("定时任务(%s)关闭成功", job.Description))
	core.OK(c, true)
}

// QuerySchedJobTypes
//
//	@Summary	查询定时任务类型
//	@Tags		schedJob
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	core.Response{data=[]models.Dict}
//	@Security	JWT
//	@Router		/schedjob/types [get]
func (*schedJobController) QuerySchedJobTypes(c *gin.Context) {
	entries, err := services.GetEntriesByType("schedJob")
	if err != nil {
		logger.Error(err)
		_ = c.Error(errs.ErrQuerySchedJob)
		return
	}

	core.OK(c, entries)
}

func (sjc *schedJobController) Setup(r *gin.RouterGroup) {
	r.GET("/schedjobs", sjc.QuerySchedJobs).
		PUT("/schedjob/open", sjc.StartSchedJobs).
		PUT("/schedjob/close", sjc.CloseSchedJobs).GET("/schedjob/types", sjc.QuerySchedJobTypes)
}
