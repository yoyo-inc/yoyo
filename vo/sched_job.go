package vo

type QuerySchedJobVO struct {
	Type   string `form:"type"`
	Status *int   `form:"status"`
}

type UpdateSchedJobVO struct {
	ID int `json:"id" binding:"required"`
}
