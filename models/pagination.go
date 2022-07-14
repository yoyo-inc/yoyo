package models

type Pagination struct {
	// 页数
	Page int `json:"page"`
	// 每页大小
	PageSize int `json:"pageSize"`
}
