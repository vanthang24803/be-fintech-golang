package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/maynguyen24/sever/internal/models"
	"github.com/maynguyen24/sever/pkg/snowflake"
)

// FundRepository manages DB operations for funds
type FundRepository struct {
	db *sqlx.DB
}

func NewFundRepository(db *sqlx.DB) *FundRepository {
	return &FundRepository{db: db}
}

// Create inserts a new fund record
func (r *FundRepository) Create(fund *models.Fund) error {
	fund.ID = snowflake.GenerateID()
	if fund.Currency == "" {
		fund.Currency = "VND"
	}

	query := `
		INSERT INTO funds (id, user_id, name, description, target_amount, balance, currency)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING created_at, updated_at
	`
	return r.db.QueryRowx(query,
		fund.ID, fund.UserID, fund.Name, fund.Description,
		fund.TargetAmount, fund.Balance, fund.Currency,
	).Scan(&fund.CreatedAt, &fund.UpdatedAt)
}

// GetAllByUserID returns all funds owned by the given user
func (r *FundRepository) GetAllByUserID(userID int64) ([]*models.Fund, error) {
	var funds []*models.Fund
	query := `
		SELECT id, user_id, name, description, target_amount, balance, currency, created_at, updated_at
		FROM funds
		WHERE user_id = $1
		ORDER BY created_at DESC
	`
	if err := r.db.Select(&funds, query, userID); err != nil {
		return nil, err
	}
	return funds, nil
}

// GetByID fetches a single fund for the authenticated user
func (r *FundRepository) GetByID(id, userID int64) (*models.Fund, error) {
	var fund models.Fund
	query := `
		SELECT id, user_id, name, description, target_amount, balance, currency, created_at, updated_at
		FROM funds
		WHERE id = $1 AND user_id = $2
		LIMIT 1
	`
	err := r.db.Get(&fund, query, id, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &fund, nil
}

// Update modifies fund metadata (name, description, target, currency)
func (r *FundRepository) Update(fund *models.Fund) error {
	query := `
		UPDATE funds
		SET name = $1, description = $2, target_amount = $3, currency = $4, updated_at = NOW()
		WHERE id = $5 AND user_id = $6
		RETURNING updated_at
	`
	err := r.db.QueryRowx(query,
		fund.Name, fund.Description, fund.TargetAmount, fund.Currency,
		fund.ID, fund.UserID,
	).Scan(&fund.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update fund: %w", err)
	}
	return nil
}

// Delete removes a fund if it belongs to the user
func (r *FundRepository) Delete(id, userID int64) error {
	result, err := r.db.Exec(
		`DELETE FROM funds WHERE id = $1 AND user_id = $2`, id, userID,
	)
	if err != nil {
		return fmt.Errorf("failed to delete fund: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("fund not found")
	}
	return nil
}

// Deposit adds amount to the fund balance atomically
func (r *FundRepository) Deposit(id, userID int64, amount float64) (*models.Fund, error) {
	var fund models.Fund
	query := `
		UPDATE funds
		SET balance = balance + $1, updated_at = NOW()
		WHERE id = $2 AND user_id = $3
		RETURNING id, user_id, name, description, target_amount, balance, currency, created_at, updated_at
	`
	err := r.db.QueryRowx(query, amount, id, userID).StructScan(&fund)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("fund not found")
		}
		return nil, fmt.Errorf("failed to deposit: %w", err)
	}
	return &fund, nil
}

// Withdraw subtracts amount from the fund balance; returns error if balance insufficient
func (r *FundRepository) Withdraw(id, userID int64, amount float64) (*models.Fund, error) {
	var fund models.Fund
	query := `
		UPDATE funds
		SET balance = balance - $1, updated_at = NOW()
		WHERE id = $2 AND user_id = $3 AND balance >= $1
		RETURNING id, user_id, name, description, target_amount, balance, currency, created_at, updated_at
	`
	err := r.db.QueryRowx(query, amount, id, userID).StructScan(&fund)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Could be not found OR insufficient balance — check which
			existing, checkErr := r.GetByID(id, userID)
			if checkErr != nil {
				return nil, checkErr
			}
			if existing == nil {
				return nil, fmt.Errorf("fund not found")
			}
			return nil, fmt.Errorf("insufficient balance")
		}
		return nil, fmt.Errorf("failed to withdraw: %w", err)
	}
	return &fund, nil
}
