package routes

import (
	"errors"
	"fmt"
	"net/http"
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
	"github.com/yoyo-inc/yoyo/services"
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
// @Success 200   {object} core.Response{data=core.PaginatedData{list=[]models.Resource}}
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
	if res := queries[0].Scopes(core.DateTimeRanger(c, "create_time"), core.Paginator(c)).Where(&query).Find(&resources); res.Error != nil {
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
// @Param resourceType path string true "资源类型"
// @Param   file formData     file true "参数"
// @Success 200   {object} core.Response{data=bool}
// @Security JWT
// @Router  /resource/:resourceType/upload [post]
func (*resourceController) UploadResource(c *gin.Context) {
	errUploadResource := errors.New("文件上传失败")

	file, err := c.FormFile("file")
	if err != nil {
		logger.Error(err)
		c.Error(errUploadResource)
		return
	}

	resourceType := c.Param("resourceType")
	resourceDir := services.GetTypedResourceDir(resourceType)
	if !fileutil.IsExist(resourceDir) {
		if err := fileutil.CreateDir(resourceDir); err != nil {
			logger.Error(err)
			c.Error(errUploadResource)
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
		c.Error(errUploadResource)
		audit_log.Fail(c, "资源管理", "上传", fmt.Sprintf("文件保存失败(文件名：%s，文件大小(kb)：%f)", file.Filename, filesize))
		return
	}

	var resource models.Resource
	resource.ResourceName = resourceName
	resource.ResourceType = resourceType
	resource.Filename = filename
	resource.FileType = strings.Replace(fileExt, ".", "", 1)
	resource.Filesize = filesize
	if res := db.Client.Create(&resource); res.Error != nil {
		logger.Error(res.Error)
		c.Error(errUploadResource)
		return
	}

	audit_log.Success(c, "资源管理", "上传", fmt.Sprintf("文件名：%s，文件大小(kb)：%f", file.Filename, filesize))
	core.OK(c, map[string]interface{}{
		"name": resource.Filename,
		"url":  config.GetString("server.base_path") + "/static/upload/" + resourceName,
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

	if err := services.DeleteResourceFile(c, id); err != nil {
		logger.Error(err)
		c.Error(errs.ErrDeleteResource)
		return
	}

	core.OK(c, true)
}

// DownloadResource
// @Summary 下载资源
// @Tags    resource
// @Accept  json
// @Produce octet-stream
// @Param   id  path   string   true "参数"
// @Security JWT
// @Router  /resource/download/:id [get]
func (*resourceController) DownloadResource(c *gin.Context) {
	id := c.Param("id")

	var resource models.Resource
	if res := db.Client.Model(&models.Resource{Model: core.Model{ID: id}}).First(&resource); res.Error != nil {
		logger.Error(res.Error)
		c.Error(errs.ErrDownloadResource)
		return
	}

	resourceFilepath := filepath.Join(services.GetTypedResourceDir(resource.ResourceType), resource.ResourceName)
	c.FileAttachment(resourceFilepath, resource.Filename)
}

func (rc *resourceController) Setup(r *gin.RouterGroup) {
	r.GET("/resources", rc.QueryResources).
		POST("/resource/:resourceType/upload", rc.UploadResource).
		DELETE("/resource/:id", rc.DeleteResource).
		StaticFS("/resource/upload", http.Dir(services.GetTypedResourceDir("upload"))).
		GET("/resource/download/:id", rc.DownloadResource)
}
