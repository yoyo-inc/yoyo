package models

type Pagination struct {
	// 页数
	Current int `json:"current"`
	// 每页大小
	PageSize int `json:"pageSize"`
}
