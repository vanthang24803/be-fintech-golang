package database

import (
	"embed"
	"fmt"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// RunDBMigration reads sql files and runs the database migrations automatically
func RunDBMigration(dsn string) error {
	d, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("failed to load embedded migrations: %w", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, dsn)
	if err != nil {
		return fmt.Errorf("failed to initialize migrate instance: %w", err)
	}

	if err := m.Up(); err != nil {
		if strings.HasPrefix(err.Error(), "Dirty database") {
			// If dirty, try to force to the current version to clear the flag
			version, _, _ := m.Version()
			if version > 0 {
				fmt.Printf("Database is dirty at version %d, forcing to clean state...\n", version)
				if err := m.Force(int(version)); err != nil {
					return fmt.Errorf("failed to force version: %w", err)
				}
				// Retry Up after forcing
				if err := m.Up(); err != nil && err != migrate.ErrNoChange {
					return fmt.Errorf("failed to run migrate up after force: %w", err)
				}
			}
		} else if err != migrate.ErrNoChange {
			return fmt.Errorf("failed to run migrate up: %w", err)
		}
	}

	return nil
}
