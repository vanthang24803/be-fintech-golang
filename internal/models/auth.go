package models

import "time"

// LoginRequest is the JSON payload for logging in
type LoginRequest struct {
	Identifier string `json:"identifier" validate:"required"` // Can be email or username
	Password   string `json:"password" validate:"required,min=6"`
}

// TokenPair represents the access and refresh token returned to client
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// LoginResponse wraps user profile and their token pair
type LoginResponse struct {
	User   *User      `json:"user"`
	Tokens *TokenPair `json:"tokens"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// LogoutRequest represents payload to log out a user (revoke session)
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// Token represents the refresh token stored in the database
type Token struct {
	ID          int64     `db:"id"`
	UserID      int64     `db:"user_id"`
	TokenString string    `db:"token_string"`
	ExpiresAt   time.Time `db:"expires_at"`
	CreatedAt   time.Time `db:"created_at"`
}
