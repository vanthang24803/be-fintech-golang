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

func TestBudgetRepository_Operations(t *testing.T) {
	t.Parallel()

	db, mock := newMockDB(t)
	repo := NewBudgetRepository(db)
	now := time.Now()
	categoryID := int64(7)
	budget := &models.Budget{
		UserID:     42,
		CategoryID: &categoryID,
		Amount:     500,
		Period:     "monthly",
		StartDate:  now,
		EndDate:    now.AddDate(0, 1, 0),
		IsActive:   true,
	}

	mock.ExpectQuery(quotedSQL(`
		INSERT INTO budgets (id, user_id, category_id, amount, period, start_date, end_date, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING created_at, updated_at
	`)).
		WithArgs(sqlmock.AnyArg(), budget.UserID, budget.CategoryID, budget.Amount, budget.Period, budget.StartDate, budget.EndDate, budget.IsActive).
		WillReturnRows(timestampRows(now))
	if err := repo.Create(context.Background(), budget); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	mock.ExpectQuery(quotedSQL(`SELECT id, user_id, category_id, amount, period, start_date, end_date, is_active, created_at, updated_at
		FROM budgets WHERE user_id = $1 ORDER BY is_active DESC, end_date ASC`)).
		WithArgs(int64(42)).
		WillReturnRows(sqlmock.NewRows(budgetCols).
			AddRow(int64(1), int64(42), categoryID, 500.0, "monthly", now, now.AddDate(0, 1, 0), true, now, now))
	list, err := repo.GetByUserID(context.Background(), 42)
	if err != nil || len(list) != 1 {
		t.Fatalf("GetByUserID() = %+v, %v", list, err)
	}

	mock.ExpectQuery(quotedSQL(`SELECT id, user_id, category_id, amount, period, start_date, end_date, is_active, created_at, updated_at
		FROM budgets WHERE id = $1 AND user_id = $2 LIMIT 1`)).
		WithArgs(int64(1), int64(42)).
		WillReturnRows(sqlmock.NewRows(budgetCols).
			AddRow(int64(1), int64(42), categoryID, 500.0, "monthly", now, now.AddDate(0, 1, 0), true, now, now))
	got, err := repo.GetByID(context.Background(), 1, 42)
	if err != nil || got == nil || got.Amount != 500 {
		t.Fatalf("GetByID() = %+v, %v", got, err)
	}

	newAmount := 600.0
	isActive := false
	mock.ExpectQuery(quotedSQL(`
		UPDATE budgets
		SET
			amount    = COALESCE($1, amount),
			is_active = COALESCE($2, is_active),
			updated_at = NOW()
		WHERE id = $3 AND user_id = $4
		RETURNING id, user_id, category_id, amount, period, start_date, end_date, is_active, created_at, updated_at
	`)).
		WithArgs(&newAmount, &isActive, int64(1), int64(42)).
		WillReturnRows(sqlmock.NewRows(budgetCols).
			AddRow(int64(1), int64(42), categoryID, newAmount, "monthly", now, now.AddDate(0, 1, 0), false, now, now))
	got, err = repo.Update(context.Background(), 1, 42, &models.UpdateBudgetRequest{Amount: &newAmount, IsActive: &isActive})
	if err != nil || got == nil || got.Amount != newAmount || got.IsActive != false {
		t.Fatalf("Update() = %+v, %v", got, err)
	}

	mock.ExpectExec(quotedSQL(`DELETE FROM budgets WHERE id = $1 AND user_id = $2`)).
		WithArgs(int64(1), int64(42)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	if err := repo.Delete(context.Background(), 1, 42); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	mock.ExpectQuery(quotedSQL(`SELECT COALESCE(SUM(amount), 0) FROM transactions 
		WHERE user_id = $1 AND type = 'expense' AND created_at >= $2 AND created_at <= $3 AND category_id = $4`)).
		WithArgs(int64(42), budget.StartDate, budget.EndDate, categoryID).
		WillReturnRows(sqlmock.NewRows([]string{"coalesce"}).AddRow(150.0))
	total, err := repo.CalculateSpending(context.Background(), 42, &categoryID, budget.StartDate, budget.EndDate)
	if err != nil || total != 150 {
		t.Fatalf("CalculateSpending() = %v, %v", total, err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("ExpectationsWereMet(): %v", err)
	}
}

func TestBudgetRepository_ErrorBranches(t *testing.T) {
	t.Parallel()

	db, mock := newMockDB(t)
	repo := NewBudgetRepository(db)

	mock.ExpectQuery(quotedSQL(`SELECT id, user_id, category_id, amount, period, start_date, end_date, is_active, created_at, updated_at
		FROM budgets WHERE id = $1 AND user_id = $2 LIMIT 1`)).
		WithArgs(int64(9), int64(42)).
		WillReturnError(sql.ErrNoRows)
	got, err := repo.GetByID(context.Background(), 9, 42)
	if err != nil || got != nil {
		t.Fatalf("expected nil on missing budget, got %+v err=%v", got, err)
	}

	mock.ExpectQuery(quotedSQL(`
		UPDATE budgets
		SET
			amount    = COALESCE($1, amount),
			is_active = COALESCE($2, is_active),
			updated_at = NOW()
		WHERE id = $3 AND user_id = $4
		RETURNING id, user_id, category_id, amount, period, start_date, end_date, is_active, created_at, updated_at
	`)).
		WithArgs((*float64)(nil), (*bool)(nil), int64(9), int64(42)).
		WillReturnError(sql.ErrNoRows)
	got, err = repo.Update(context.Background(), 9, 42, &models.UpdateBudgetRequest{})
	if err != nil || got != nil {
		t.Fatalf("expected nil update result on missing budget, got %+v err=%v", got, err)
	}

	mock.ExpectExec(quotedSQL(`DELETE FROM budgets WHERE id = $1 AND user_id = $2`)).
		WithArgs(int64(9), int64(42)).
		WillReturnResult(sqlmock.NewResult(0, 0))
	if err := repo.Delete(context.Background(), 9, 42); !errors.Is(err, apperr.ErrNotFound) {
		t.Fatalf("expected not found, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("ExpectationsWereMet(): %v", err)
	}
}
