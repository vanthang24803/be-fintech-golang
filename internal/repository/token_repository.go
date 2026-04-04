package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/maynguyen24/sever/internal/models"
)

// TokenRepository manages database operations for auth tokens
type TokenRepository struct {
	db *sqlx.DB
}

func NewTokenRepository(db *sqlx.DB) *TokenRepository {
	return &TokenRepository{db: db}
}

func (r *TokenRepository) StoreRefreshToken(token *models.Token) error {
	query := `
		INSERT INTO tokens (id, user_id, token_string, expires_at)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.Exec(query, token.ID, token.UserID, token.TokenString, token.ExpiresAt)
	if err != nil {
		return fmt.Errorf("failed to store refresh token: %w", err)
	}
	return nil
}

func (r *TokenRepository) RevokeToken(tokenString string) error {
	query := `DELETE FROM tokens WHERE token_string = $1`
	_, err := r.db.Exec(query, tokenString)
	if err != nil {
		return fmt.Errorf("failed to revoke token: %w", err)
	}
	return nil
}

func (r *TokenRepository) GetToken(tokenString string) (*models.Token, error) {
	var token models.Token
	query := `SELECT id, user_id, token_string, expires_at, created_at FROM tokens WHERE token_string = $1 LIMIT 1`
	err := r.db.Get(&token, query, tokenString)
	if err != nil {
		return nil, err
	}
	return &token, nil
}
