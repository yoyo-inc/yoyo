package errs

import "github.com/yoyo-inc/yoyo/core"

var (
	ErrQueryResources   = core.NewBusinessError("800", "查询资源文件失败")
	ErrUploadResource   = core.NewBusinessError("801", "上传资源文件失败")
	ErrDeleteResource   = core.NewBusinessError("802", "删除资源文件失败")
	ErrNotExistResource = core.NewBusinessError("803", "资源文件不存在")
	ErrDownloadResource = core.NewBusinessError("804", "资源文件下载失败")
)
