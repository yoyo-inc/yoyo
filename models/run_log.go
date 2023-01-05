package models

import "github.com/yoyo-inc/yoyo/core"

type RunLog struct {
	core.IModel
	Filename string `json:"filename,omitempty"`
	Filesize string `json:"filesize,omitempty"`
}
