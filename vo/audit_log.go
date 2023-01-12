package vo

type QueryAuditLogVO struct {
	Module string `form:"module"`
	Status int    `form:"status"`
}
