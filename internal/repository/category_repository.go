package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/maynguyen24/sever/internal/models"
	"github.com/maynguyen24/sever/pkg/snowflake"
)

// CategoryRepository manages DB operations for categories
type CategoryRepository struct {
	db *sqlx.DB
}

func NewCategoryRepository(db *sqlx.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

// Create inserts a new user-owned category
func (r *CategoryRepository) Create(cat *models.Category) error {
	cat.ID = snowflake.GenerateID()
	query := `
		INSERT INTO categories (id, user_id, name, type, icon)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING created_at, updated_at
	`
	return r.db.QueryRowx(query,
		cat.ID, cat.UserID, cat.Name, cat.Type, cat.Icon,
	).Scan(&cat.CreatedAt, &cat.UpdatedAt)
}

// GetAllByUserID returns both system categories (user_id IS NULL) and user-owned ones
func (r *CategoryRepository) GetAllByUserID(userID int64) ([]*models.Category, error) {
	var cats []*models.Category
	query := `
		SELECT id, user_id, name, type, icon, created_at, updated_at
		FROM categories
		WHERE user_id = $1 OR user_id IS NULL
		ORDER BY user_id NULLS FIRST, name ASC
	`
	if err := r.db.Select(&cats, query, userID); err != nil {
		return nil, err
	}
	return cats, nil
}

// GetByID fetches a category only if it belongs to the user or is a system category
func (r *CategoryRepository) GetByID(id, userID int64) (*models.Category, error) {
	var cat models.Category
	query := `
		SELECT id, user_id, name, type, icon, created_at, updated_at
		FROM categories
		WHERE id = $1 AND (user_id = $2 OR user_id IS NULL)
		LIMIT 1
	`
	err := r.db.Get(&cat, query, id, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &cat, nil
}

// GetOwnedByID fetches a category only if explicitly owned by the user (for mutations)
func (r *CategoryRepository) GetOwnedByID(id, userID int64) (*models.Category, error) {
	var cat models.Category
	query := `
		SELECT id, user_id, name, type, icon, created_at, updated_at
		FROM categories
		WHERE id = $1 AND user_id = $2
		LIMIT 1
	`
	err := r.db.Get(&cat, query, id, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &cat, nil
}

// Update mutates an existing user-owned category
func (r *CategoryRepository) Update(cat *models.Category) error {
	query := `
		UPDATE categories
		SET name = $1, type = $2, icon = $3, updated_at = NOW()
		WHERE id = $4 AND user_id = $5
		RETURNING updated_at
	`
	err := r.db.QueryRowx(query,
		cat.Name, cat.Type, cat.Icon, cat.ID, cat.UserID,
	).Scan(&cat.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update category: %w", err)
	}
	return nil
}

// Delete removes a user-owned category
func (r *CategoryRepository) Delete(id, userID int64) error {
	query := `DELETE FROM categories WHERE id = $1 AND user_id = $2`
	result, err := r.db.Exec(query, id, userID)
	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("category not found or not owned by user")
	}
	return nil
}
