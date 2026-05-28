package sqldata

import "embed"

// FS embeds the SQL migration and seed files so the server binary is self-contained.
//
//go:embed migrations/*.sql seed.sql
var FS embed.FS
