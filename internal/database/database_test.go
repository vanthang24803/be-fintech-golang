package database

import (
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func TestConnectRequiresDSN(t *testing.T) {
	t.Parallel()

	db, err := Connect("")
	if err == nil || db != nil {
		t.Fatalf("expected missing dsn error, got db=%v err=%v", db, err)
	}
}

func TestRunDBMigrationInvalidDSN(t *testing.T) {
	t.Parallel()

	err := RunDBMigration("invalid-dsn")
	if err == nil || !strings.Contains(err.Error(), "failed to initialize migrate instance") {
		t.Fatalf("expected migration init error, got %v", err)
	}
}

func TestCloseWithNilDB(t *testing.T) {
	t.Parallel()

	DB = nil
	if err := Close(); err != nil {
		t.Fatalf("Close() with nil DB error = %v", err)
	}
}

func TestCloseWithOpenDB(t *testing.T) {
	t.Parallel()

	rawDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New(): %v", err)
	}
	rawDB2, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New(): %v", err)
	}

	_ = rawDB.Close()
	mock.ExpectClose()
	DB = sqlx.NewDb(rawDB2, "sqlmock")
	if err := Close(); err != nil {
		t.Fatalf("Close() with open DB error = %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("ExpectationsWereMet() = %v", err)
	}
	DB = nil
}
