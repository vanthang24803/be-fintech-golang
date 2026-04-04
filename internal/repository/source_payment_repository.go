package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/maynguyen24/sever/internal/models"
	"github.com/maynguyen24/sever/pkg/snowflake"
)

// SourcePaymentRepository manages DB operations for payment sources
type SourcePaymentRepository struct {
	db *sqlx.DB
}

func NewSourcePaymentRepository(db *sqlx.DB) *SourcePaymentRepository {
	return &SourcePaymentRepository{db: db}
}

func (r *SourcePaymentRepository) Create(source *models.SourcePayment) error {
	source.ID = snowflake.GenerateID()
	if source.Currency == "" {
		source.Currency = "VND"
	}

	query := `
		INSERT INTO sourcepayment (id, user_id, name, type, balance, currency)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at, updated_at
	`
	return r.db.QueryRowx(query,
		source.ID, source.UserID, source.Name, source.Type, source.Balance, source.Currency,
	).Scan(&source.CreatedAt, &source.UpdatedAt)
}

func (r *SourcePaymentRepository) GetAllByUserID(userID int64) ([]*models.SourcePayment, error) {
	var sources []*models.SourcePayment
	query := `SELECT id, user_id, name, type, balance, currency, created_at, updated_at
		FROM sourcepayment WHERE user_id = $1 ORDER BY created_at DESC`
	if err := r.db.Select(&sources, query, userID); err != nil {
		return nil, err
	}
	return sources, nil
}

func (r *SourcePaymentRepository) GetByID(id, userID int64) (*models.SourcePayment, error) {
	var source models.SourcePayment
	query := `SELECT id, user_id, name, type, balance, currency, created_at, updated_at
		FROM sourcepayment WHERE id = $1 AND user_id = $2 LIMIT 1`
	err := r.db.Get(&source, query, id, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &source, nil
}

func (r *SourcePaymentRepository) Update(source *models.SourcePayment) error {
	query := `
		UPDATE sourcepayment
		SET name = $1, type = $2, currency = $3, updated_at = NOW()
		WHERE id = $4 AND user_id = $5
		RETURNING updated_at
	`
	err := r.db.QueryRowx(query,
		source.Name, source.Type, source.Currency, source.ID, source.UserID,
	).Scan(&source.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update source payment: %w", err)
	}
	return nil
}

func (r *SourcePaymentRepository) Delete(id, userID int64) error {
	query := `DELETE FROM sourcepayment WHERE id = $1 AND user_id = $2`
	result, err := r.db.Exec(query, id, userID)
	if err != nil {
		return fmt.Errorf("failed to delete source payment: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("source payment not found")
	}
	return nil
}
