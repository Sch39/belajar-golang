// internal/storage/storage.go
package storage

import "embed"

//go:embed migrations/*.sql
var MigrationFiles embed.FS