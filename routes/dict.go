package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/yoyo-inc/yoyo/common/logger"
	"github.com/yoyo-inc/yoyo/core"
	"github.com/yoyo-inc/yoyo/errs"
	"github.com/yoyo-inc/yoyo/models"
)

type dictController struct{}

// QueryDicts
//	@Summary	查询字典表
//	@Tags		dict
//	@Accept		json
//	@Produce	json
//	@Param		query	query		models.Pagination	false	"参数"
//	@Success	200		{object}	core.Response{data=core.PaginatedData{list=models.Dict}}
//	@Security	JWT
//	@Router		/dicts [get]
func (*dictController) QueryDicts(c *gin.Context) {
	queries := core.GetPaginatedQuery(&models.Dict{})

	var dicts []models.Dict
	if res := queries[0].Scopes(core.Paginator(c), core.DateTimeRanger(c, "")).Find(&dicts); res.Error != nil {
		logger.Error(res.Error)
		c.Error(errs.ErrQueryDict)
		return
	}

	var count int64
	if res := queries[1].Count(&count); res.Error != nil {
		logger.Error(res.Error)
		c.Error(errs.ErrQueryDict)
		return
	}

	core.OK(c, core.Paginated(dicts, count))
}

func (dc *dictController) Setup(r *gin.RouterGroup) {
	r.GET("/dicts", dc.QueryDicts)
}
