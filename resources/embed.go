package resources

import (
	"embed"
	_ "embed"
)

//go:embed sql/default.sql
var DefaultSql []byte

//go:embed prometheus
var PrometheusDir embed.FS

//go:embed alert-manager
var AlertmanagerDir embed.FS

//go:embed report
var InternalReportTplDir embed.FS
