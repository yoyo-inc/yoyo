package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/yoyo-inc/yoyo/common/db"
	"github.com/yoyo-inc/yoyo/common/logger"
	"github.com/yoyo-inc/yoyo/core"
	"github.com/yoyo-inc/yoyo/errs"
	"github.com/yoyo-inc/yoyo/models"
)

type roleController struct{}

// RetrieveRole returns role list
// @Summary  查询角色列表
// @Tags     role
// @Accept   json
// @Produce  json
// @Param    query  query     models.QueryRole  true  "参数"
// @Success  200    {object}  core.Response
// @Router   /role  [get]
func (role *roleController) RetrieveRole(c *gin.Context) {
	var query models.QueryRole
	if err := c.ShouldBindJSON(&query); err != nil {
		c.Error(core.NewParameterError(err.Error()))
		return
	}

	var roles []models.Role
	if result := db.Client.Model(&models.Role{}).Where(&query).Find(&roles); result.Error != nil {
		logger.Error(result.Error)
		c.Error(errs.ErrQueryRole)
		return
	}

	core.OK(c, roles)
}

func (role *roleController) CreateRole(c *gin.Context) {
	var query models.Role
	if err := c.ShouldBindJSON(&query); err != nil {
		c.Error(core.NewParameterError(err.Error()))
		return
	}

	if result := db.Client.Create(&query); result.Error != nil {
		logger.Error(result.Error)
		c.Error(errs.ErrCreateRole)
		return
	}

	core.OK(c, true)
}

func (role *roleController) Setup(r *gin.RouterGroup) {

}
