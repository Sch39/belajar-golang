// internal/storage/database/postgres_migrate.go
package database

import (
	"context"
	"embed"
	"fmt"
	"log"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed migrations/postgres/*.sql
var postgresMigrationFiles embed.FS

func MigratePostgres(pool *pgxpool.Pool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	const migrationPath = "migrations/postgres"
	const migrationTable = "migrations"

	// 1. ensure migration table
	createTableSQL := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			id BIGSERIAL PRIMARY KEY,
			filename TEXT NOT NULL UNIQUE,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
	`, migrationTable)

	if _, err := pool.Exec(ctx, createTableSQL); err != nil {
		return fmt.Errorf("create migration table failed: %w", err)
	}

	// 2. load applied migrations
	applied := map[string]bool{}

	rows, err := pool.Query(
		ctx,
		fmt.Sprintf(`SELECT filename FROM %s`, migrationTable),
	)
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

	if rows.Err() != nil {
		return rows.Err()
	}

	// 3. read migration files
	entries, err := postgresMigrationFiles.ReadDir(migrationPath)
	if err != nil {
		return err
	}

	var pending []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if strings.HasSuffix(entry.Name(), ".sql") && !applied[entry.Name()] {
			pending = append(pending, entry.Name())
		}
	}

	sort.Strings(pending)

	if len(pending) == 0 {
		log.Println("[migrate] postgres database is up to date")
		return nil
	}

	log.Printf("[migrate] found %d pending migrations\n", len(pending))

	// 4. apply migrations (transaction per file)
	for _, name := range pending {
		fullPath := path.Join(migrationPath, name)

		sqlBytes, err := postgresMigrationFiles.ReadFile(fullPath)
		if err != nil {
			return err
		}

		tx, err := pool.Begin(ctx)
		if err != nil {
			return err
		}

		if _, err := tx.Exec(ctx, string(sqlBytes)); err != nil {
			_ = tx.Rollback(ctx)
			return fmt.Errorf("migration %s failed: %w", name, err)
		}

		if _, err := tx.Exec(
			ctx,
			fmt.Sprintf(`INSERT INTO %s (filename) VALUES ($1)`, migrationTable),
			name,
		); err != nil {
			_ = tx.Rollback(ctx)
			return err
		}

		if err := tx.Commit(ctx); err != nil {
			return err
		}

		log.Printf("[migrate] ✅ %s applied\n", name)
	}

	log.Println("[migrate] postgres migration completed")
	return nil
}
