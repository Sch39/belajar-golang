// internal/storage/database/postgres_migrate.go

package database

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"path"
	"sort"
	"strings"
)

//go:embed migrations/postgres/*.sql
var postgresMigrationFiles embed.FS

func MigratePostgres(db *sql.DB) error {
	const migrationPath = "migrations/postgres"
	const migrationTable = "migrations"

	createMigrationTableQuery := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			id SERIAL PRIMARY KEY,
			filename TEXT NOT NULL UNIQUE,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
	`, migrationTable)

	if _, err := db.Exec(createMigrationTableQuery); err != nil {
		return err
	}

	applied := make(map[string]bool)
	rows, err := db.Query(fmt.Sprintf(`SELECT filename FROM %s`, migrationTable))
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var filename string
		if err := rows.Scan(&filename); err != nil {
			return err
		}
		applied[filename] = true
	}

	entries, err := postgresMigrationFiles.ReadDir(migrationPath)
	if err != nil {
		return err
	}

	var pending []string
	for _, entry := range entries {
		name := entry.Name()
		if !entry.IsDir() && strings.HasSuffix(name, ".sql") {
			if !applied[name] {
				pending = append(pending, name)
			}
		}
	}

	sort.Strings(pending)

	if len(pending) == 0 {
		log.Println("[migrate] postgres database is up to date")
		return nil
	}

	log.Printf("[migrate] found %d pending migrations\n", len(pending))

	for _, name := range pending {
		fullPath := path.Join(migrationPath, name)

		sqlBytes, err := postgresMigrationFiles.ReadFile(fullPath)
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

		if _, err := tx.Exec(
			fmt.Sprintf(`INSERT INTO %s (filename) VALUES ($1)`, migrationTable),
			name,
		); err != nil {
			tx.Rollback()
			return err
		}

		if err := tx.Commit(); err != nil {
			return err
		}

		log.Printf("[migrate] ✅ %s applied successfully\n", name)
	}

	log.Println("[migrate] postgres migration completed")
	return nil
}
