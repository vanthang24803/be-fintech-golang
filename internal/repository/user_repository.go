package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/maynguyen24/sever/internal/models"
	"github.com/maynguyen24/sever/pkg/snowflake"
)

// UserRepository manages database operations for user entities.
// Following Go principles: We export the struct and do not define the interface here.
type UserRepository struct {
	db *sqlx.DB
}

// NewUserRepository creates a new repository instance.
func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetUserByEmailOrUsername(email, username string) (*models.User, error) {
	var user models.User
	query := "SELECT id, username, email, password_hash, google_id, created_at, updated_at FROM users WHERE email = $1 OR username = $2 LIMIT 1"
	err := r.db.Get(&user, query, email, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Not found
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	query := "SELECT id, username, email, password_hash, google_id, created_at, updated_at FROM users WHERE email = $1 LIMIT 1"
	err := r.db.Get(&user, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetUserByGoogleID(googleID string) (*models.User, error) {
	var user models.User
	query := "SELECT id, username, email, password_hash, google_id, created_at, updated_at FROM users WHERE google_id = $1 LIMIT 1"
	err := r.db.Get(&user, query, googleID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) CreateUser(user *models.User) error {
	// Begin transaction
	tx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// 1. Insert User
	userQuery := `
		INSERT INTO users (id, username, email, password_hash, google_id) 
		VALUES ($1, $2, $3, $4, $5) 
		RETURNING created_at, updated_at
	`
	err = tx.QueryRowx(userQuery, user.ID, user.Username, user.Email, user.PasswordHash, user.GoogleID).Scan(&user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}

	// 2. Insert Empty Profile
	profileQuery := `
		INSERT INTO profiles (id, user_id) 
		VALUES ($1, $2)
	`
	_, err = tx.Exec(profileQuery, snowflake.GenerateID(), user.ID)
	if err != nil {
		return fmt.Errorf("failed to insert empty profile: %w", err)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *UserRepository) GetProfileByUserID(userID int64) (*models.Profile, error) {
	var profile models.Profile
	query := "SELECT id, user_id, full_name, avatar_url, phone_number, date_of_birth, created_at, updated_at FROM profiles WHERE user_id = $1 LIMIT 1"
	err := r.db.Get(&profile, query, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Not found
		}
		return nil, err
	}
	return &profile, nil
}

func (r *UserRepository) GetUserByID(userID int64) (*models.User, error) {
	var user models.User
	query := "SELECT id, username, email, password_hash, google_id, created_at, updated_at FROM users WHERE id = $1 LIMIT 1"
	err := r.db.Get(&user, query, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Not found
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) LinkGoogleAccount(userID int64, googleID string) error {
	query := "UPDATE users SET google_id = $1, updated_at = NOW() WHERE id = $2"
	_, err := r.db.Exec(query, googleID, userID)
	return err
}

func (r *UserRepository) UpdateProfile(userID int64, req *models.UpdateProfileRequest) (*models.Profile, error) {
	var profile models.Profile
	query := `
		UPDATE profiles
		SET
			full_name    = COALESCE($1, full_name),
			avatar_url   = COALESCE($2, avatar_url),
			phone_number = COALESCE($3, phone_number),
			date_of_birth = COALESCE($4, date_of_birth),
			updated_at   = NOW()
		WHERE user_id = $5
		RETURNING id, user_id, full_name, avatar_url, phone_number, date_of_birth, created_at, updated_at
	`
	err := r.db.QueryRowx(query,
		req.FullName, req.AvatarURL, req.PhoneNumber, req.DateOfBirth, userID,
	).StructScan(&profile)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}
	return &profile, nil
}
