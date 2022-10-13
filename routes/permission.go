package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/yoyo-inc/yoyo/common/db"
	"github.com/yoyo-inc/yoyo/common/logger"
	"github.com/yoyo-inc/yoyo/core"
	"github.com/yoyo-inc/yoyo/errs"
	"github.com/yoyo-inc/yoyo/models"
	"github.com/yoyo-inc/yoyo/vo"
)

type permissionController struct{}

// RetrievePermissions
// @Summary  查询列表
// @Tags     permissions
// @Produce  json
// @Success  200  {object}  core.Response{data=vo.PermissionVO}
// @Router   /permissions [get]
func (*permissionController) RetrievePermissions(c *gin.Context) {
	var permissions []models.Permission
	if res := db.Client.Model(&models.Permission{}).Find(&permissions); res.Error != nil {
		logger.Error(res.Error)
		c.Error(errs.ErrQueryPermissions)
		return
	}

	var rootPermission vo.PermissionVO
	var permissionMap = make(map[string]*vo.PermissionVO)
	for _, val := range permissions {
		permissionVO := vo.PermissionVO{
			Permission: val,
		}

		permissionMap[val.ID] = &permissionVO
	}

	for _, val := range permissions {
		var parent *vo.PermissionVO
		if val.ParentID == "" {
			parent = &rootPermission
		} else {
			parent = permissionMap[val.ParentID]
		}
		parent.Children = append(parent.Children, permissionMap[val.ID])
	}

	core.OK(c, rootPermission)
}

func (permission *permissionController) Setup(r *gin.RouterGroup) {
	r.GET("/permissions", permission.RetrievePermissions)
}
