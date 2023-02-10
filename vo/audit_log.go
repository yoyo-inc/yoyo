package vo

type QueryAuditLogVO struct {
	Module string `form:"module"`
	Status *int   `form:"status"`
	IP     string `form:"ip"`
	UserID string `form:"user"`
}
