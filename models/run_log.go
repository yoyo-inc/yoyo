package models

import "github.com/yoyo-inc/yoyo/core"

type RunLog struct {
	core.IModel
	Filename string `json:"filename"`
	Filesize string `json:"filesize"`
}
