package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/maynguyen24/sever/internal/models"
	"github.com/maynguyen24/sever/pkg/apperr"
	"github.com/maynguyen24/sever/pkg/snowflake"
)

// TransactionRepository manages DB operations for transactions
type TransactionRepository struct {
	db *sqlx.DB
}

func NewTransactionRepository(db *sqlx.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

// Create inserts a transaction and atomically updates the source payment balance
func (r *TransactionRepository) Create(tx *models.Transaction) error {
	tx.ID = snowflake.GenerateID()

	// Determine balance delta: income → +amount, expense → -amount
	var balanceDelta float64
	if tx.Type == models.TransactionTypeIncome {
		balanceDelta = tx.Amount
	} else {
		balanceDelta = -tx.Amount
	}

	dbTx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer dbTx.Rollback()

	// 1. Insert transaction record
	insertQuery := `
		INSERT INTO transactions (id, user_id, sourcepayment_id, category_id, amount, type, description, transaction_date)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING created_at, updated_at
	`
	err = dbTx.QueryRowx(insertQuery,
		tx.ID, tx.UserID, tx.SourcePaymentID, tx.CategoryID,
		tx.Amount, tx.Type, tx.Description, tx.TransactionDate,
	).Scan(&tx.CreatedAt, &tx.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to insert transaction: %w", err)
	}

	// 2. Atomically update balance on the source payment
	updateBalanceQuery := `
		UPDATE sourcepayment
		SET balance = balance + $1, updated_at = NOW()
		WHERE id = $2 AND user_id = $3
	`
	result, err := dbTx.Exec(updateBalanceQuery, balanceDelta, tx.SourcePaymentID, tx.UserID)
	if err != nil {
		return fmt.Errorf("failed to update balance: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("source payment not found or not owned by user")
	}

	return dbTx.Commit()
}

// GetAllByUserID lists all transactions for a user with optional filters
func (r *TransactionRepository) GetAllByUserID(userID int64, filter models.TransactionFilter) ([]*models.TransactionDetail, error) {
	query := `
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
	`
	args := []interface{}{userID}
	argIdx := 2

	if filter.Type != "" {
		query += fmt.Sprintf(" AND t.type = $%d", argIdx)
		args = append(args, filter.Type)
		argIdx++
	}
	if filter.CategoryID != 0 {
		query += fmt.Sprintf(" AND t.category_id = $%d", argIdx)
		args = append(args, filter.CategoryID)
		argIdx++
	}
	if filter.SourcePaymentID != 0 {
		query += fmt.Sprintf(" AND t.sourcepayment_id = $%d", argIdx)
		args = append(args, filter.SourcePaymentID)
		argIdx++
	}

	query += " ORDER BY t.transaction_date DESC, t.created_at DESC"

	var txs []*models.TransactionDetail
	if err := r.db.Select(&txs, query, args...); err != nil {
		return nil, err
	}
	return txs, nil
}

// GetByID fetches a single transaction detail for the authenticated user
func (r *TransactionRepository) GetByID(id, userID int64) (*models.TransactionDetail, error) {
	var tx models.TransactionDetail
	query := `
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
	`
	err := r.db.Get(&tx, query, id, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &tx, nil
}

// GetRawByID fetches a raw transaction (no join) for mutation operations
func (r *TransactionRepository) GetRawByID(id, userID int64) (*models.Transaction, error) {
	var tx models.Transaction
	query := `SELECT id, user_id, sourcepayment_id, category_id, amount, type, description, transaction_date, created_at, updated_at
		FROM transactions WHERE id = $1 AND user_id = $2 LIMIT 1`
	err := r.db.Get(&tx, query, id, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &tx, nil
}

// Update modifies a transaction and reverses + reapplies the balance delta atomically
func (r *TransactionRepository) Update(old *models.Transaction, updated *models.Transaction) error {
	dbTx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer dbTx.Rollback()

	// Calculate the net balance adjustment:
	// Reverse old effect, then apply new effect
	oldDelta := map[string]float64{models.TransactionTypeIncome: old.Amount, models.TransactionTypeExpense: -old.Amount}[old.Type]
	newDelta := map[string]float64{models.TransactionTypeIncome: updated.Amount, models.TransactionTypeExpense: -updated.Amount}[updated.Type]

	// 1. Reverse old balance on old source
	_, err = dbTx.Exec(
		`UPDATE sourcepayment SET balance = balance - $1, updated_at = NOW() WHERE id = $2 AND user_id = $3`,
		oldDelta, old.SourcePaymentID, old.UserID,
	)
	if err != nil {
		return fmt.Errorf("failed to reverse old balance: %w", err)
	}

	// 2. Apply new balance on new source (may be same or different)
	_, err = dbTx.Exec(
		`UPDATE sourcepayment SET balance = balance + $1, updated_at = NOW() WHERE id = $2 AND user_id = $3`,
		newDelta, updated.SourcePaymentID, updated.UserID,
	)
	if err != nil {
		return fmt.Errorf("failed to apply new balance: %w", err)
	}

	// 3. Update the transaction record
	updateQuery := `
		UPDATE transactions
		SET sourcepayment_id = $1, category_id = $2, amount = $3, type = $4,
		    description = $5, transaction_date = $6, updated_at = NOW()
		WHERE id = $7 AND user_id = $8
		RETURNING updated_at
	`
	err = dbTx.QueryRowx(updateQuery,
		updated.SourcePaymentID, updated.CategoryID, updated.Amount, updated.Type,
		updated.Description, updated.TransactionDate, updated.ID, updated.UserID,
	).Scan(&updated.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update transaction: %w", err)
	}

	return dbTx.Commit()
}

// Delete removes a transaction and reverses its balance impact atomically
func (r *TransactionRepository) Delete(id, userID int64) error {
	// Fetch the transaction first to know how to reverse balance
	raw, err := r.GetRawByID(id, userID)
	if err != nil {
		return err
	}
	if raw == nil {
		return apperr.ErrNotFound
	}

	// Reverse: income was +amount → subtract; expense was -amount → add back
	reverseDelta := map[string]float64{models.TransactionTypeIncome: -raw.Amount, models.TransactionTypeExpense: raw.Amount}[raw.Type]

	dbTx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer dbTx.Rollback()

	_, err = dbTx.Exec(
		`DELETE FROM transactions WHERE id = $1 AND user_id = $2`,
		id, userID,
	)
	if err != nil {
		return fmt.Errorf("failed to delete transaction: %w", err)
	}

	_, err = dbTx.Exec(
		`UPDATE sourcepayment SET balance = balance + $1, updated_at = NOW() WHERE id = $2 AND user_id = $3`,
		reverseDelta, raw.SourcePaymentID, userID,
	)
	if err != nil {
		return fmt.Errorf("failed to reverse balance on delete: %w", err)
	}

	return dbTx.Commit()
}
