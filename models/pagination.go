package models

type Pagination struct {
	Current  int `json:"current,omitempty"`  // 页数
	PageSize int `json:"pageSize,omitempty"` // 每页大小
}
