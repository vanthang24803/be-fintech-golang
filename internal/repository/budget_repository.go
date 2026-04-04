package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/maynguyen24/sever/internal/models"
	"github.com/maynguyen24/sever/pkg/snowflake"
)

// BudgetRepository handles database operations for budgets module
type BudgetRepository struct {
	db *sqlx.DB
}

func NewBudgetRepository(db *sqlx.DB) *BudgetRepository {
	return &BudgetRepository{db: db}
}

// Create inserts a new budget record
func (r *BudgetRepository) Create(budget *models.Budget) error {
	budget.ID = snowflake.GenerateID()

	query := `
		INSERT INTO budgets (id, user_id, category_id, amount, period, start_date, end_date, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING created_at, updated_at
	`
	return r.db.QueryRowx(query,
		budget.ID, budget.UserID, budget.CategoryID, budget.Amount,
		budget.Period, budget.StartDate, budget.EndDate, budget.IsActive,
	).Scan(&budget.CreatedAt, &budget.UpdatedAt)
}

// GetByUserID fetches all budgets for a user
func (r *BudgetRepository) GetByUserID(userID int64) ([]*models.Budget, error) {
	var budgets []*models.Budget
	query := `SELECT id, user_id, category_id, amount, period, start_date, end_date, is_active, created_at, updated_at
		FROM budgets WHERE user_id = $1 ORDER BY is_active DESC, end_date ASC`
	if err := r.db.Select(&budgets, query, userID); err != nil {
		return nil, err
	}
	return budgets, nil
}

// GetByID fetches a specific budget record
func (r *BudgetRepository) GetByID(id, userID int64) (*models.Budget, error) {
	var budget models.Budget
	query := `SELECT id, user_id, category_id, amount, period, start_date, end_date, is_active, created_at, updated_at
		FROM budgets WHERE id = $1 AND user_id = $2 LIMIT 1`
	err := r.db.Get(&budget, query, id, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &budget, nil
}

// Update modifies budget settings
func (r *BudgetRepository) Update(id, userID int64, req *models.UpdateBudgetRequest) (*models.Budget, error) {
	var budget models.Budget
	query := `
		UPDATE budgets
		SET
			amount    = COALESCE($1, amount),
			is_active = COALESCE($2, is_active),
			updated_at = NOW()
		WHERE id = $3 AND user_id = $4
		RETURNING id, user_id, category_id, amount, period, start_date, end_date, is_active, created_at, updated_at
	`
	err := r.db.QueryRowx(query, req.Amount, req.IsActive, id, userID).StructScan(&budget)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Not found or wrong user
		}
		return nil, err
	}
	return &budget, nil
}

// Delete removes a budget record
func (r *BudgetRepository) Delete(id, userID int64) error {
	result, err := r.db.Exec(`DELETE FROM budgets WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("budget not found")
	}
	return nil
}

// CalculateSpending returns total sum of expenses for a user and optional category within a date range
func (r *BudgetRepository) CalculateSpending(userID int64, categoryID *int64, start, end time.Time) (float64, error) {
	var total float64
	query := `SELECT COALESCE(SUM(amount), 0) FROM transactions 
		WHERE user_id = $1 AND type = 'expense' AND created_at >= $2 AND created_at <= $3`
	args := []interface{}{userID, start, end}

	if categoryID != nil {
		query += " AND category_id = $4"
		args = append(args, *categoryID)
	}

	if err := r.db.Get(&total, query, args...); err != nil {
		return 0, err
	}
	return total, nil
}
