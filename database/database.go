package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func NewPostgres() (*sql.DB, error) {
	dsn := "host=pgdb port=5432 user=db_user password=pwd123 dbname=taskmanager_db sslmode=disable"

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return db, nil
}
