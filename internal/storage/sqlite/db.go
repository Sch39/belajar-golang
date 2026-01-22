// internal\storage\sqlite\db.go
package sqlite

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func NewSQLiteDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	schema := `
	CREATE TABLE IF NOT EXISTS products (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		price INTEGER NOT NULL, 
		stock INTEGER NOT NULL
	)
	`
	_, err = db.Exec(schema)

	if err != nil {
		return nil, err
	}

	log.Println("Database initialized successfully")
	return db, err
}