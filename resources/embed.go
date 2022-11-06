package resources

import (
	_ "embed"
)

//go:embed sql/default.sql
var DefaultSql []byte
