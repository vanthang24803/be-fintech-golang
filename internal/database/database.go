package database

import (
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

var DB *sqlx.DB

// Connect initializes the database connection using a provided DSN
func Connect(dsn string) error {
	if dsn == "" {
		return fmt.Errorf("DATABASE_URL is not provided")
	}

	// sqlx.Connect automatically calls Ping to ensure connection is valid
	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	DB = db

	// Auto-migrate schema on start
	if err := RunDBMigration(dsn); err != nil {
		return fmt.Errorf("auto-migration failed: %w", err)
	}

	return nil
}

// Close gracefully closes the database connection
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
