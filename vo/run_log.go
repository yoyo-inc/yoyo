package vo

type RunLogVO struct {
	Service  string `json:"service,omitempty"`
	Filename string `json:"filename,omitempty"`
	Date     string `json:"date,omitempty"`
	Filesize string `json:"filesize,omitempty"`
}

type DownloadRunLogVO struct {
	Service  string `form:"service" binding:"required"`
	Filename string `form:"filename" binding:"required"`
}
