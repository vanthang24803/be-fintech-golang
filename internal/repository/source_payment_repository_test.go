package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/maynguyen24/sever/internal/models"
)

func TestSourcePaymentRepository_BasicOperations(t *testing.T) {
	t.Parallel()

	db, mock := newMockDB(t)
	repo := NewSourcePaymentRepository(db)
	now := time.Now()
	source := &models.SourcePayment{UserID: 42, Name: "Wallet", Type: "wallet", Balance: 100}

	mock.ExpectQuery(quotedSQL(`
		INSERT INTO sourcepayment (id, user_id, name, type, balance, currency)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at, updated_at
	`)).
		WithArgs(sqlmock.AnyArg(), source.UserID, source.Name, source.Type, source.Balance, "VND").
		WillReturnRows(timestampRows(now))
	if err := repo.Create(context.Background(), source); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	mock.ExpectQuery(quotedSQL(`SELECT id, user_id, name, type, balance, currency, created_at, updated_at
		FROM sourcepayment WHERE user_id = $1 ORDER BY created_at DESC`)).
		WithArgs(int64(42)).
		WillReturnRows(sqlmock.NewRows(sourceCols).
			AddRow(int64(1), int64(42), "Wallet", "wallet", 100.0, "VND", now, now))
	list, err := repo.GetAllByUserID(context.Background(), 42)
	if err != nil || len(list) != 1 {
		t.Fatalf("GetAllByUserID() = %+v, %v", list, err)
	}

	mock.ExpectQuery(quotedSQL(`SELECT id, user_id, name, type, balance, currency, created_at, updated_at
		FROM sourcepayment WHERE id = $1 AND user_id = $2 LIMIT 1`)).
		WithArgs(int64(1), int64(42)).
		WillReturnRows(sqlmock.NewRows(sourceCols).
			AddRow(int64(1), int64(42), "Wallet", "wallet", 100.0, "VND", now, now))
	got, err := repo.GetByID(context.Background(), 1, 42)
	if err != nil || got == nil || got.ID != 1 {
		t.Fatalf("GetByID() = %+v, %v", got, err)
	}

	got.Currency = "USD"
	mock.ExpectQuery(quotedSQL(`
		UPDATE sourcepayment
		SET name = $1, type = $2, currency = $3, updated_at = NOW()
		WHERE id = $4 AND user_id = $5
		RETURNING updated_at
	`)).
		WithArgs(got.Name, got.Type, got.Currency, got.ID, got.UserID).
		WillReturnRows(sqlmock.NewRows([]string{"updated_at"}).AddRow(now))
	if err := repo.Update(context.Background(), got); err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	mock.ExpectExec(quotedSQL(`DELETE FROM sourcepayment WHERE id = $1 AND user_id = $2`)).
		WithArgs(int64(1), int64(42)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	if err := repo.Delete(context.Background(), 1, 42); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	mock.ExpectQuery(quotedSQL(`SELECT id, user_id, name, type, balance, currency, created_at, updated_at
		FROM sourcepayment WHERE id = $1 AND user_id = $2 LIMIT 1`)).
		WithArgs(int64(9), int64(42)).
		WillReturnError(sql.ErrNoRows)
	got, err = repo.GetByID(context.Background(), 9, 42)
	if err != nil || got != nil {
		t.Fatalf("expected nil on not found, got %+v err=%v", got, err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("ExpectationsWereMet(): %v", err)
	}
}
