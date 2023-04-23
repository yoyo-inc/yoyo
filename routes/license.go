package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/yoyo-inc/yoyo/common/db"
	"github.com/yoyo-inc/yoyo/common/license"
	"github.com/yoyo-inc/yoyo/common/logger"
	"github.com/yoyo-inc/yoyo/core"
	"github.com/yoyo-inc/yoyo/errs"
	"github.com/yoyo-inc/yoyo/models"
	"github.com/yoyo-inc/yoyo/vo"
)

type LicenseController struct {
}

// QueryLicense
//
//	@Summary	查询激活信息
//	@Tags		license
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	core.Response{data=models.License}
//	@Security	JWT
//	@Router		/license  [get]
func (*LicenseController) QueryLicense(c *gin.Context) {
	var license models.License
	if res := db.Client.Model(&models.License{}).First(license); res.Error != nil {
		logger.Error(res.Error)
		_ = c.Error(errs.ErrQueryLicense)
		return
	}

	core.OK(c, license)
}

// ActivateLicense
//
//	@Summary	激活
//	@Tags		license
//	@Accept		json
//	@Produce	json
//	@Param		body	body		vo.ActivateLicenseVO	true	"参数"
//	@Success	200		{object}	core.Response{data=bool}
//	@Security	JWT
//	@Router		/license [put]
func (*LicenseController) ActivateLicense(c *gin.Context) {
	var query vo.ActivateLicenseVO
	if err := c.ShouldBindJSON(&query); err != nil {
		logger.Error(err)
		_ = c.Error(core.NewParameterError(err))
		return
	}

	res := license.Service.Activate(query.License)
	core.OK(c, res)
}

func (lc *LicenseController) Setup(r *gin.RouterGroup) {
	r.GET("/license", lc.QueryLicense).PUT("/license", lc.ActivateLicense)
}
