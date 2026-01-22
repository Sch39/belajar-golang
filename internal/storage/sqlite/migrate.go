// internal/storage/sqlite/migrate.go
package sqlite

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"path"
	"sort"
	"strings"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

func Migrate(db *sql.DB) error {
	const migrationPath = "migrations"
	const migrationTable = "migrations"

	createMigrationTableQuery := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		filename TEXT NOT NULL UNIQUE,
		applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`, migrationTable)

	_, err := db.Exec(createMigrationTableQuery)
	if err != nil {
		return err
	}

	appliedMap := make(map[string]bool)
	selectMigrationsQuery := fmt.Sprintf(`SELECT filename FROM %s`, migrationTable)
	rows, err := db.Query(selectMigrationsQuery)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var filename string
		if err := rows.Scan(&filename); err != nil {
			return err
		}
		appliedMap[filename] = true
	}

	entries, err := migrationFiles.ReadDir(migrationPath)
	if err != nil {
		return err
	}

	var pendingFiles []string
	for _, entry := range entries {
		name := entry.Name()
		if !entry.IsDir() && strings.HasSuffix(name, ".sql") {
			if !appliedMap[name] {
				pendingFiles = append(pendingFiles, name)
			}
		}
	}
	sort.Strings(pendingFiles)

	if len(pendingFiles) == 0 {
		log.Println("[migrate] database is up to date")
		return nil
	}

	log.Printf("[migrate] found %d pending migrations", len(pendingFiles))

	for _, name := range pendingFiles {
		fullPath := path.Join(migrationPath, name)

		sqlBytes, err := migrationFiles.ReadFile(fullPath)
		if err != nil {
			return err
		}
tx, err := db.Begin()
		if err != nil {
			return err
		}

		if _, err := tx.Exec(string(sqlBytes)); err != nil {
			tx.Rollback()
			log.Printf("[migrate] ❌ %s failed: %v\n", name, err)
			return err
		}

		insertMigrationQuery := fmt.Sprintf(`INSERT INTO %s (filename) VALUES (?)`, migrationTable)
		if _, err := tx.Exec(insertMigrationQuery, name); err != nil {
			tx.Rollback()
			return err
		}

		if err := tx.Commit(); err != nil {
			return err
		}

		log.Printf("[migrate] ✅ %s applied successfully\n", name)
	}

	log.Println("[migrate] done")
	return nil
}
