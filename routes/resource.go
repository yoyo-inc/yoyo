package routes

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/duke-git/lancet/v2/fileutil"
	"github.com/gin-gonic/gin"
	"github.com/golang-module/carbon/v2"
	"github.com/yoyo-inc/yoyo/common/config"
	"github.com/yoyo-inc/yoyo/common/db"
	"github.com/yoyo-inc/yoyo/common/logger"
	"github.com/yoyo-inc/yoyo/core"
	"github.com/yoyo-inc/yoyo/errs"
	"github.com/yoyo-inc/yoyo/models"
	"github.com/yoyo-inc/yoyo/services/audit_log"
	"github.com/yoyo-inc/yoyo/vo"
)

type resourceController struct{}

// QueryResources
// @Summary 查询资源文件列表
// @Tags    resource
// @Accept  json
// @Produce json
// @Param query query vo.QueryResourceVO true "参数"
// @Success 200   {object} core.Response{data=array,models.Resource}
// @Security JWT
// @Router  /resources [get]
func (*resourceController) QueryResources(c *gin.Context) {
	var query vo.QueryResourceVO
	if err := c.ShouldBindQuery(&query); err != nil {
		logger.Error(err)
		c.Error(core.NewParameterError(err))
		return
	}

	queries := core.GetPaginatedQuery(&models.Resource{})

	for i := range queries {
		if query.Filename != "" {
			queries[i].Where("filename like ?", "%"+query.Filename+"%")
			query.Filename = ""
		}
	}

	var resources []models.Resource
	if res := queries[0].Scopes(core.DateTimeRanger(c, "create_time")).Where(&query).Find(&resources); res.Error != nil {
		logger.Error(res.Error)
		c.Error(errs.ErrQueryResources)
		return
	}

	var count int64
	if res := queries[1].Scopes(core.DateTimeRanger(c, "create_time")).Where(&query).Count(&count); res.Error != nil {
		logger.Error(res.Error)
		c.Error(errs.ErrQueryResources)
		return
	}

	core.OK(c, core.Paginated(resources, count))
}

// UploadResource
// @Summary 上传资源
// @Tags    resource
// @Accept  mpfd
// @Produce json
// @Param   file formData     file true "参数"
// @Success 200   {object} core.Response{data=bool}
// @Security JWT
// @Router  /resource [post]
func (*resourceController) UploadResource(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		logger.Error(err)
		c.Error(errs.ErrUploadResource)
		return
	}

	resourceDir, err := filepath.Abs(config.GetString("resource_dir"))
	if err != nil {
		logger.Error(err)
		c.Error(errs.ErrUploadResource)
		return
	}

	if !fileutil.IsExist(resourceDir) {
		if err := fileutil.CreateDir(resourceDir); err != nil {
			logger.Error(err)
			c.Error(errs.ErrUploadResource)
			return
		}
	}

	filename := file.Filename
	fileExt := filepath.Ext(filename)
	filenameWithoutExt := strings.Replace(file.Filename, fileExt, "", 1)
	resourceName := filenameWithoutExt + carbon.Now().ToShortDateTimeString() + fileExt
	resourcePath := filepath.Join(resourceDir, resourceName)
	filesize, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(file.Size)/float64(1024)), 64)

	if err := c.SaveUploadedFile(file, resourcePath); err != nil {
		logger.Error(err)
		c.Error(errs.ErrUploadResource)
		audit_log.Fail(c, "资源管理", "上传", fmt.Sprintf("文件保存失败(文件名：%s，文件大小(kb)：%f)", file.Filename, filesize))
		return
	}

	var resource models.Resource
	resource.ResourceName = resourceName
	resource.Filename = filename
	resource.FileType = strings.Replace(fileExt, ".", "", 1)
	resource.Filesize = filesize
	if res := db.Client.Create(&resource); res.Error != nil {
		logger.Error(res.Error)
		c.Error(errs.ErrUploadResource)
		return
	}

	audit_log.Success(c, "资源管理", "上传", fmt.Sprintf("文件名：%s，文件大小(kb)：%f", file.Filename, filesize))

	c.JSON(http.StatusOK, map[string]interface{}{
		"name": resource.Filename,
		"url":  config.GetString("server.base_path") + "/static" + resourceName,
		"id":   resource.ID,
	})
}

// DeleteResource
// @Summary 删除资源文件
// @Tags    resource
// @Accept  json
// @Produce json
// @Param   id path     string true "参数"
// @Success 200   {object} core.Response{data=bool}
// @Security JWT
// @Router  /resource/:id [delete]
func (rc *resourceController) DeleteResource(c *gin.Context) {
	id := c.Param("id")

	var resource models.Resource
	if res := db.Client.Model(&models.Resource{SModel: core.SModel{ID: id}}).Find(&resource); res.Error != nil {
		logger.Error(res.Error)
		c.Error(errs.ErrNotExistResource)
		return
	}

	resourceDir, err := filepath.Abs(config.GetString("resource_dir"))
	if err != nil {
		logger.Error(err)
		c.Error(errs.ErrDeleteResource)
		return
	}

	// delete resource file
	if err := fileutil.RemoveFile(filepath.Join(resourceDir, resource.ResourceName)); err != nil {
		logger.Error(err)
		if !errors.Is(err, os.ErrNotExist) {
			c.Error(errs.ErrDeleteResource)
			return
		}
		audit_log.Fail(c, "资源管理", "删除", fmt.Sprintf("资源文件不存在，资源文件名：%s", resource.Filename))
	}

	// delete resource record
	if res := db.Client.Delete(&models.Resource{SModel: core.SModel{ID: id}}); res.Error != nil {
		logger.Error(res.Error)
		c.Error(errs.ErrDeleteResource)
		audit_log.Fail(c, "资源管理", "删除", fmt.Sprintf("数据记录不存在，资源文件名：%s", resource.Filename))
		return
	}

	audit_log.Success(c, "资源管理", "删除", fmt.Sprintf("资源文件名：%s", resource.Filename))

	core.OK(c, true)
}

func (rc *resourceController) Setup(r *gin.RouterGroup) {
	resourceDir, _ := filepath.Abs(config.GetString("resource_dir"))
	r.GET("/resources", rc.QueryResources).
		POST("/resource", rc.UploadResource).
		DELETE("/resource/:id", rc.DeleteResource).
		StaticFS("/static", http.Dir(resourceDir))
}
