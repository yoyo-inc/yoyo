package models

import "github.com/yoyo-inc/yoyo/core"

type RunLog struct {
	core.Model
	Filename string `json:"filename"`
	Filesize string `json:"filesize"`
}
