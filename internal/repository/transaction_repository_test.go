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

func TestTransactionRepository_CreateAndReads(t *testing.T) {
	t.Parallel()

	db, mock := newMockDB(t)
	repo := NewTransactionRepository(db)
	now := time.Now()
	desc := "salary"
	categoryID := int64(7)
	tx := &models.Transaction{
		UserID:          42,
		SourcePaymentID: 3,
		CategoryID:      &categoryID,
		Amount:          500,
		Type:            models.TransactionTypeIncome,
		Description:     &desc,
		TransactionDate: now,
	}

	mock.ExpectBegin()
	mock.ExpectQuery(quotedSQL(`
		INSERT INTO transactions (id, user_id, sourcepayment_id, category_id, amount, type, description, transaction_date)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING created_at, updated_at
	`)).
		WithArgs(sqlmock.AnyArg(), tx.UserID, tx.SourcePaymentID, tx.CategoryID, tx.Amount, tx.Type, tx.Description, tx.TransactionDate).
		WillReturnRows(timestampRows(now))
	mock.ExpectExec(quotedSQL(`
		UPDATE sourcepayment
		SET balance = balance + $1, updated_at = NOW()
		WHERE id = $2 AND user_id = $3
	`)).
		WithArgs(500.0, int64(3), int64(42)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()
	if err := repo.Create(context.Background(), tx); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	mock.ExpectQuery(quotedSQL(`
		SELECT
			t.id, t.user_id, t.sourcepayment_id, t.category_id,
			t.amount, t.type, t.description, t.transaction_date,
			t.created_at, t.updated_at,
			sp.name AS source_name,
			c.name  AS category_name
		FROM transactions t
		JOIN sourcepayment sp ON sp.id = t.sourcepayment_id
		LEFT JOIN categories c ON c.id = t.category_id
		WHERE t.user_id = $1
		 AND t.type = $2 AND t.category_id = $3 AND t.sourcepayment_id = $4 ORDER BY t.transaction_date DESC, t.created_at DESC
	`)).
		WithArgs(int64(42), models.TransactionTypeIncome, int64(7), int64(3)).
		WillReturnRows(sqlmock.NewRows(transactionDCols).
			AddRow(int64(1), int64(42), int64(3), int64(7), 500.0, models.TransactionTypeIncome, desc, now, now, now, "Wallet", "Salary"))
	list, err := repo.GetAllByUserID(context.Background(), 42, models.TransactionFilter{
		Type:            models.TransactionTypeIncome,
		CategoryID:      7,
		SourcePaymentID: 3,
	})
	if err != nil || len(list) != 1 || list[0].SourceName != "Wallet" {
		t.Fatalf("GetAllByUserID() = %+v, %v", list, err)
	}

	mock.ExpectQuery(quotedSQL(`
		SELECT
			t.id, t.user_id, t.sourcepayment_id, t.category_id,
			t.amount, t.type, t.description, t.transaction_date,
			t.created_at, t.updated_at,
			sp.name AS source_name,
			c.name  AS category_name
		FROM transactions t
		JOIN sourcepayment sp ON sp.id = t.sourcepayment_id
		LEFT JOIN categories c ON c.id = t.category_id
		WHERE t.id = $1 AND t.user_id = $2
		LIMIT 1
	`)).
		WithArgs(int64(1), int64(42)).
		WillReturnRows(sqlmock.NewRows(transactionDCols).
			AddRow(int64(1), int64(42), int64(3), int64(7), 500.0, models.TransactionTypeIncome, desc, now, now, now, "Wallet", "Salary"))
	detail, err := repo.GetByID(context.Background(), 1, 42)
	if err != nil || detail == nil || detail.ID != 1 {
		t.Fatalf("GetByID() = %+v, %v", detail, err)
	}

	mock.ExpectQuery(quotedSQL(`SELECT id, user_id, sourcepayment_id, category_id, amount, type, description, transaction_date, created_at, updated_at
		FROM transactions WHERE id = $1 AND user_id = $2 LIMIT 1`)).
		WithArgs(int64(1), int64(42)).
		WillReturnRows(sqlmock.NewRows(transactionCols).
			AddRow(int64(1), int64(42), int64(3), int64(7), 500.0, models.TransactionTypeIncome, desc, now, now, now))
	raw, err := repo.GetRawByID(context.Background(), 1, 42)
	if err != nil || raw == nil || raw.SourcePaymentID != 3 {
		t.Fatalf("GetRawByID() = %+v, %v", raw, err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("ExpectationsWereMet(): %v", err)
	}
}

func TestTransactionRepository_UpdateAndDelete(t *testing.T) {
	t.Parallel()

	db, mock := newMockDB(t)
	repo := NewTransactionRepository(db)
	now := time.Now()
	oldDesc := "old"
	newDesc := "new"
	oldCategoryID := int64(7)
	newCategoryID := int64(8)
	oldTx := &models.Transaction{
		ID:              1,
		UserID:          42,
		SourcePaymentID: 3,
		CategoryID:      &oldCategoryID,
		Amount:          100,
		Type:            models.TransactionTypeIncome,
		Description:     &oldDesc,
		TransactionDate: now,
	}
	newTx := &models.Transaction{
		ID:              1,
		UserID:          42,
		SourcePaymentID: 4,
		CategoryID:      &newCategoryID,
		Amount:          50,
		Type:            models.TransactionTypeExpense,
		Description:     &newDesc,
		TransactionDate: now,
	}

	mock.ExpectBegin()
	mock.ExpectExec(quotedSQL(`UPDATE sourcepayment SET balance = balance - $1, updated_at = NOW() WHERE id = $2 AND user_id = $3`)).
		WithArgs(100.0, int64(3), int64(42)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(quotedSQL(`UPDATE sourcepayment SET balance = balance + $1, updated_at = NOW() WHERE id = $2 AND user_id = $3`)).
		WithArgs(-50.0, int64(4), int64(42)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectQuery(quotedSQL(`
		UPDATE transactions
		SET sourcepayment_id = $1, category_id = $2, amount = $3, type = $4,
		    description = $5, transaction_date = $6, updated_at = NOW()
		WHERE id = $7 AND user_id = $8
		RETURNING updated_at
	`)).
		WithArgs(newTx.SourcePaymentID, newTx.CategoryID, newTx.Amount, newTx.Type, newTx.Description, newTx.TransactionDate, newTx.ID, newTx.UserID).
		WillReturnRows(sqlmock.NewRows([]string{"updated_at"}).AddRow(now))
	mock.ExpectCommit()
	if err := repo.Update(context.Background(), oldTx, newTx); err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	mock.ExpectQuery(quotedSQL(`SELECT id, user_id, sourcepayment_id, category_id, amount, type, description, transaction_date, created_at, updated_at
		FROM transactions WHERE id = $1 AND user_id = $2 LIMIT 1`)).
		WithArgs(int64(1), int64(42)).
		WillReturnRows(sqlmock.NewRows(transactionCols).
			AddRow(int64(1), int64(42), int64(4), int64(8), 50.0, models.TransactionTypeExpense, newDesc, now, now, now))
	mock.ExpectBegin()
	mock.ExpectExec(quotedSQL(`DELETE FROM transactions WHERE id = $1 AND user_id = $2`)).
		WithArgs(int64(1), int64(42)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(quotedSQL(`UPDATE sourcepayment SET balance = balance + $1, updated_at = NOW() WHERE id = $2 AND user_id = $3`)).
		WithArgs(50.0, int64(4), int64(42)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()
	if err := repo.Delete(context.Background(), 1, 42); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("ExpectationsWereMet(): %v", err)
	}
}

func TestTransactionRepository_ErrorBranches(t *testing.T) {
	t.Parallel()

	db, mock := newMockDB(t)
	repo := NewTransactionRepository(db)
	now := time.Now()
	desc := "expense"

	tx := &models.Transaction{
		UserID:          42,
		SourcePaymentID: 3,
		Amount:          50,
		Type:            models.TransactionTypeExpense,
		TransactionDate: now,
	}
	mock.ExpectBegin()
	mock.ExpectQuery(quotedSQL(`
		INSERT INTO transactions (id, user_id, sourcepayment_id, category_id, amount, type, description, transaction_date)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING created_at, updated_at
	`)).
		WithArgs(sqlmock.AnyArg(), tx.UserID, tx.SourcePaymentID, tx.CategoryID, tx.Amount, tx.Type, tx.Description, tx.TransactionDate).
		WillReturnRows(timestampRows(now))
	mock.ExpectExec(quotedSQL(`
		UPDATE sourcepayment
		SET balance = balance + $1, updated_at = NOW()
		WHERE id = $2 AND user_id = $3
	`)).
		WithArgs(-50.0, int64(3), int64(42)).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectRollback()
	if err := repo.Create(context.Background(), tx); err == nil {
		t.Fatal("expected create error when source payment update touches no rows")
	}

	mock.ExpectQuery(quotedSQL(`
		SELECT
			t.id, t.user_id, t.sourcepayment_id, t.category_id,
			t.amount, t.type, t.description, t.transaction_date,
			t.created_at, t.updated_at,
			sp.name AS source_name,
			c.name  AS category_name
		FROM transactions t
		JOIN sourcepayment sp ON sp.id = t.sourcepayment_id
		LEFT JOIN categories c ON c.id = t.category_id
		WHERE t.id = $1 AND t.user_id = $2
		LIMIT 1
	`)).
		WithArgs(int64(9), int64(42)).
		WillReturnError(sql.ErrNoRows)
	detail, err := repo.GetByID(context.Background(), 9, 42)
	if err != nil || detail != nil {
		t.Fatalf("expected nil detail on not found, got %+v err=%v", detail, err)
	}

	mock.ExpectQuery(quotedSQL(`SELECT id, user_id, sourcepayment_id, category_id, amount, type, description, transaction_date, created_at, updated_at
		FROM transactions WHERE id = $1 AND user_id = $2 LIMIT 1`)).
		WithArgs(int64(9), int64(42)).
		WillReturnError(sql.ErrNoRows)
	if err := repo.Delete(context.Background(), 9, 42); !errors.Is(err, apperr.ErrNotFound) {
		t.Fatalf("expected not found on delete, got %v", err)
	}

	mock.ExpectQuery(quotedSQL(`SELECT id, user_id, sourcepayment_id, category_id, amount, type, description, transaction_date, created_at, updated_at
		FROM transactions WHERE id = $1 AND user_id = $2 LIMIT 1`)).
		WithArgs(int64(1), int64(42)).
		WillReturnRows(sqlmock.NewRows(transactionCols).
			AddRow(int64(1), int64(42), int64(3), int64(7), 50.0, models.TransactionTypeIncome, desc, now, now, now))
	mock.ExpectBegin()
	mock.ExpectExec(quotedSQL(`DELETE FROM transactions WHERE id = $1 AND user_id = $2`)).
		WithArgs(int64(1), int64(42)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(quotedSQL(`UPDATE sourcepayment SET balance = balance + $1, updated_at = NOW() WHERE id = $2 AND user_id = $3`)).
		WithArgs(-50.0, int64(3), int64(42)).
		WillReturnError(sql.ErrConnDone)
	mock.ExpectRollback()
	if err := repo.Delete(context.Background(), 1, 42); err == nil {
		t.Fatal("expected delete reverse-balance error")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("ExpectationsWereMet(): %v", err)
	}
}
