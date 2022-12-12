package services

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/duke-git/lancet/v2/fileutil"
	"github.com/gin-gonic/gin"
	"github.com/yoyo-inc/yoyo/common/config"
	"github.com/yoyo-inc/yoyo/common/db"
	"github.com/yoyo-inc/yoyo/core"
	"github.com/yoyo-inc/yoyo/errs"
	"github.com/yoyo-inc/yoyo/models"
	"github.com/yoyo-inc/yoyo/services/audit_log"
	"gorm.io/gorm"
)

// DeleteResourceFile delete resource file by resourceID
func DeleteResourceFile(c *gin.Context, resourceID string) error {
	var resource models.Resource
	if res := db.Client.Model(&models.Resource{Model: core.Model{ID: resourceID}}).First(&resource); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return errs.ErrNotExistResource
		} else {
			return errs.ErrDeleteResource
		}
	}

	resourcePath := filepath.Join(GetTypedResourceDir(resource.ResourceType), resource.ResourceName)
	if err := fileutil.RemoveFile(resourcePath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			audit_log.Fail(c, "资源管理", "删除", fmt.Sprintf("资源文件不存在，资源文件名：%s", resource.Filename))
		} else {
			audit_log.Fail(c, "资源管理", "删除", fmt.Sprintf("资源文件删除失败，资源文件名：%s", resource.Filename))
			return err
		}
	}

	if res := db.Client.Delete(&models.Resource{Model: core.Model{ID: resourceID}}); res.Error != nil {
		return res.Error
	}

	audit_log.Success(c, "资源管理", "删除", fmt.Sprintf("资源文件名：%s", resource.Filename))
	return nil
}

// GetTypedResourceDir returns directory path where stores the typed resources
func GetTypedResourceDir(t string) string {
	typepResourceDir, _ := filepath.Abs(filepath.Join(config.GetString("resource_dir"), t))
	return typepResourceDir
}
