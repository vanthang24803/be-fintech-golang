package repository

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/maynguyen24/sever/internal/models"
	"github.com/maynguyen24/sever/pkg/apperr"
)

func TestFundRepository_CRUDAndBalanceOperations(t *testing.T) {
	t.Parallel()

	db, mock := newMockDB(t)
	repo := NewFundRepository(db)
	now := time.Now()
	desc := "desc"
	fund := &models.Fund{UserID: 42, Name: "Trip", Description: &desc, TargetAmount: 1000, Balance: 250}

	mock.ExpectQuery(quotedSQL(`
		INSERT INTO funds (id, user_id, name, description, target_amount, balance, currency)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING created_at, updated_at
	`)).
		WithArgs(sqlmock.AnyArg(), fund.UserID, fund.Name, fund.Description, fund.TargetAmount, fund.Balance, "VND").
		WillReturnRows(timestampRows(now))
	if err := repo.Create(context.Background(), fund); err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if fund.Currency != "VND" {
		t.Fatalf("expected default currency VND, got %q", fund.Currency)
	}

	mock.ExpectQuery(quotedSQL(`
		SELECT id, user_id, name, description, target_amount, balance, currency, created_at, updated_at
		FROM funds
		WHERE user_id = $1
		ORDER BY created_at DESC
	`)).
		WithArgs(int64(42)).
		WillReturnRows(sqlmock.NewRows(fundCols).
			AddRow(int64(1), int64(42), "Trip", desc, 1000.0, 250.0, "VND", now, now))
	funds, err := repo.GetAllByUserID(context.Background(), 42)
	if err != nil || len(funds) != 1 {
		t.Fatalf("GetAllByUserID() = %+v, %v", funds, err)
	}

	mock.ExpectQuery(quotedSQL(`
		SELECT id, user_id, name, description, target_amount, balance, currency, created_at, updated_at
		FROM funds
		WHERE id = $1 AND user_id = $2
		LIMIT 1
	`)).
		WithArgs(int64(1), int64(42)).
		WillReturnRows(sqlmock.NewRows(fundCols).
			AddRow(int64(1), int64(42), "Trip", desc, 1000.0, 250.0, "VND", now, now))
	got, err := repo.GetByID(context.Background(), 1, 42)
	if err != nil || got == nil || got.ID != 1 {
		t.Fatalf("GetByID() = %+v, %v", got, err)
	}

	mock.ExpectQuery(quotedSQL(`
		UPDATE funds
		SET name = $1, description = $2, target_amount = $3, currency = $4, updated_at = NOW()
		WHERE id = $5 AND user_id = $6
		RETURNING updated_at
	`)).
		WithArgs("Trip", &desc, 1000.0, "USD", int64(1), int64(42)).
		WillReturnRows(sqlmock.NewRows([]string{"updated_at"}).AddRow(now))
	got.Currency = "USD"
	if err := repo.Update(context.Background(), got); err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	mock.ExpectExec(quotedSQL(`DELETE FROM funds WHERE id = $1 AND user_id = $2`)).
		WithArgs(int64(1), int64(42)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	if err := repo.Delete(context.Background(), 1, 42); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	mock.ExpectQuery(quotedSQL(`
		UPDATE funds
		SET balance = balance + $1, updated_at = NOW()
		WHERE id = $2 AND user_id = $3
		RETURNING id, user_id, name, description, target_amount, balance, currency, created_at, updated_at
	`)).
		WithArgs(50.0, int64(1), int64(42)).
		WillReturnRows(sqlmock.NewRows(fundCols).
			AddRow(int64(1), int64(42), "Trip", desc, 1000.0, 300.0, "USD", now, now))
	got, err = repo.Deposit(context.Background(), 1, 42, 50)
	if err != nil || got.Balance != 300 {
		t.Fatalf("Deposit() = %+v, %v", got, err)
	}

	mock.ExpectQuery(quotedSQL(`
		UPDATE funds
		SET balance = balance - $1, updated_at = NOW()
		WHERE id = $2 AND user_id = $3 AND balance >= $1
		RETURNING id, user_id, name, description, target_amount, balance, currency, created_at, updated_at
	`)).
		WithArgs(25.0, int64(1), int64(42)).
		WillReturnRows(sqlmock.NewRows(fundCols).
			AddRow(int64(1), int64(42), "Trip", desc, 1000.0, 275.0, "USD", now, now))
	got, err = repo.Withdraw(context.Background(), 1, 42, 25)
	if err != nil || got.Balance != 275 {
		t.Fatalf("Withdraw() = %+v, %v", got, err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("ExpectationsWereMet(): %v", err)
	}
}

func TestFundRepository_ErrorBranches(t *testing.T) {
	t.Parallel()

	db, mock := newMockDB(t)
	repo := NewFundRepository(db)

	mock.ExpectExec(quotedSQL(`DELETE FROM funds WHERE id = $1 AND user_id = $2`)).
		WithArgs(int64(9), int64(42)).
		WillReturnResult(sqlmock.NewResult(0, 0))
	if err := repo.Delete(context.Background(), 9, 42); !errors.Is(err, apperr.ErrNotFound) {
		t.Fatalf("expected not found, got %v", err)
	}

	mock.ExpectQuery(quotedSQL(`
		UPDATE funds
		SET balance = balance + $1, updated_at = NOW()
		WHERE id = $2 AND user_id = $3
		RETURNING id, user_id, name, description, target_amount, balance, currency, created_at, updated_at
	`)).
		WithArgs(50.0, int64(1), int64(42)).
		WillReturnError(sql.ErrNoRows)
	if _, err := repo.Deposit(context.Background(), 1, 42, 50); !errors.Is(err, apperr.ErrNotFound) {
		t.Fatalf("expected deposit not found, got %v", err)
	}

	mock.ExpectQuery(quotedSQL(`
		UPDATE funds
		SET balance = balance - $1, updated_at = NOW()
		WHERE id = $2 AND user_id = $3 AND balance >= $1
		RETURNING id, user_id, name, description, target_amount, balance, currency, created_at, updated_at
	`)).
		WithArgs(75.0, int64(1), int64(42)).
		WillReturnError(sql.ErrNoRows)
	mock.ExpectQuery(quotedSQL(`
		SELECT id, user_id, name, description, target_amount, balance, currency, created_at, updated_at
		FROM funds
		WHERE id = $1 AND user_id = $2
		LIMIT 1
	`)).
		WithArgs(int64(1), int64(42)).
		WillReturnRows(sqlmock.NewRows(fundCols).
			AddRow(int64(1), int64(42), "Trip", "desc", 1000.0, 10.0, "VND", time.Now(), time.Now()))
	if _, err := repo.Withdraw(context.Background(), 1, 42, 75); !errors.Is(err, apperr.ErrInsufficientBalance) {
		t.Fatalf("expected insufficient balance, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("ExpectationsWereMet(): %v", err)
	}
}
