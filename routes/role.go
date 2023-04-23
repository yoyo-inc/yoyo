package routes

import (
	"strconv"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/gin-gonic/gin"
	"github.com/yoyo-inc/yoyo/common/db"
	"github.com/yoyo-inc/yoyo/common/logger"
	"github.com/yoyo-inc/yoyo/core"
	"github.com/yoyo-inc/yoyo/errs"
	"github.com/yoyo-inc/yoyo/models"
	"github.com/yoyo-inc/yoyo/services/audit_log"
	"github.com/yoyo-inc/yoyo/vo"
	"gorm.io/gorm/clause"
)

type roleController struct{}

// QueryRoles
//
//	@Summary	查询角色列表
//	@Tags		role
//	@Accept		json
//	@Produce	json
//	@Param		query	query		models.Pagination	false	"参数"
//	@Param		query	query		vo.QueryRoleVO		true	"参数"
//	@Success	200		{object}	core.Response
//	@Security	JWT
//	@Router		/roles  [get]
func (*roleController) QueryRoles(c *gin.Context) {
	var query vo.QueryRoleVO
	if err := c.ShouldBindQuery(&query); err != nil {
		_ = c.Error(core.NewParameterError(err))
		return
	}

	queries := core.GetPaginatedQuery(&models.Role{})
	for i := range queries {
		if query.Name != "" {
			queries[i].Where("name like ?", "%"+query.Name+"%")
		}
	}
	query.Name = ""
	var roles []models.Role
	if res := queries[0].Preload("Permissions").Scopes(core.Paginator(c)).Where(&query).Find(&roles); res.Error != nil {
		logger.Error(res.Error)
		_ = c.Error(errs.ErrQueryRole)
		return
	}

	var count int64
	if res := queries[1].Where(&query).Count(&count); res.Error != nil {
		logger.Error(res.Error)
		_ = c.Error(errs.ErrQueryRole)
		return
	}

	core.OK(c, core.Paginated(roles, count))
}

// CreateRole
//
//	@Summary	创建角色
//	@Tags		role
//	@Accept		json
//	@Produce	json
//	@Param		query	body		vo.RoleVO	true	"参数"
//	@Success	200		{object}	core.Response{data=boolean}
//	@Security	JWT
//	@Router		/role  [post]
func (*roleController) CreateRole(c *gin.Context) {
	var query vo.RoleVO
	if err := c.ShouldBindJSON(&query); err != nil {
		_ = c.Error(core.NewParameterError(err))
		return
	}

	var count int64
	if res := db.Client.Model(&models.Role{}).Where("name = ?", query.Name).Count(&count); res.Error != nil {
		logger.Error(res.Error)
		_ = c.Error(errs.ErrCreateRole)
		return
	}
	if count > 0 {
		logger.Error(errs.ErrExistRole.Message)
		_ = c.Error(errs.ErrExistRole)
		return
	}

	// process role permissions
	role := query.Role
	if query.Permissions != nil {
		for _, val := range query.Permissions {
			role.Permissions = append(role.Permissions, models.Permission{IModel: core.IModel{ID: val}})
		}
	}
	if result := db.Client.Create(&role); result.Error != nil {
		logger.Error(result.Error)
		_ = c.Error(errs.ErrCreateRole)
		return
	}

	audit_log.Success(c, "角色", "新增", role.Name)
	core.OK(c, true)
}

// DeleteRole
//
//	@Summary	删除角色
//	@Tags		role
//	@Accept		json
//	@Produce	json
//	@Param		roleID	path		string	true	"参数"
//	@Success	200		{object}	core.Response{data=boolean}
//	@Security	JWT
//	@Router		/role/{roleID}  [delete]
func (*roleController) DeleteRole(c *gin.Context) {
	roleID := c.Param("roleID")
	id, _ := strconv.Atoi(roleID)

	var existRole models.Role
	if res := db.Client.Where("id = ?", id).First(&existRole); res.Error != nil {
		logger.Error(res.Error)
		_ = c.Error(errs.ErrNotExistRole)
		return
	}

	// delete associations
	if res := db.Client.Select(clause.Associations).Delete(&models.Role{IModel: core.IModel{ID: id}}); res.Error != nil {
		logger.Error(res.Error)
		_ = c.Error(errs.ErrDeleteRole)
		return
	}

	core.OK(c, true)
}

// UpdateRole
//
//	@Summary	更新角色
//	@Tags		role
//	@Accept		json
//	@Produce	json
//	@Param		query	body		vo.RoleVO	true	"参数"
//	@Success	200		{object}	core.Response{data=boolean}
//	@Security	JWT
//	@Router		/role  [put]
func (*roleController) UpdateRole(c *gin.Context) {
	var query vo.RoleVO

	if err := c.ShouldBindJSON(&query); err != nil {
		logger.Error(err)
		_ = c.Error(core.NewParameterError(err))
		return
	}

	if res := db.Client.Updates(&query.Role); res.Error != nil {
		logger.Error(res.Error)
		_ = c.Error(errs.ErrUpdateRole)
		return
	}

	// process role associations
	// process user permissions
	if query.Permissions != nil {
		permissions := slice.Map(query.Permissions, func(_ int, permissionID int) models.Permission {
			return models.Permission{IModel: core.IModel{ID: permissionID}}
		})

		if err := db.Client.Model(&query.Role).Association("Permissions").Replace(permissions); err != nil {
			logger.Error(err)
			_ = c.Error(errs.ErrUpdateRole)
			return
		}
	}

	core.OK(c, true)
}

func (role *roleController) Setup(r *gin.RouterGroup) {
	r.GET("/roles", role.QueryRoles).POST("/role", role.CreateRole).DELETE("/role/:roleID", role.DeleteRole).PUT("/role", role.UpdateRole)
}
