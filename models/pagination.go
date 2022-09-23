package models

type Pagination struct {
	Current  int `json:"current"`  // 页数
	PageSize int `json:"pageSize"` // 每页大小
}
