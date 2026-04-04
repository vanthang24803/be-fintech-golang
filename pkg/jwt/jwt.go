package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/maynguyen24/sever/configs"
)

var (
	ErrInvalidToken = errors.New("invalid or expired token")
)

type TokenClaims struct {
	UserID       int64 `json:"user_id,string"`
	FIDOVerified bool  `json:"fido_verified"`
	jwt.RegisteredClaims
}

// GenerateTokenPair creates dual tokens (Access 15m, Refresh 30d)
func GenerateTokenPair(userID int64, fidoVerified bool, cfg *configs.Config) (string, string, error) {
	// 1. Access Token (15m)
	accessClaims := TokenClaims{
		UserID:       userID,
		FIDOVerified: fidoVerified,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		return "", "", err
	}

	// 2. Refresh Token (30d)
	refreshClaims := TokenClaims{
		UserID:       userID,
		FIDOVerified: fidoVerified,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(cfg.JWTRefreshSecret))
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}
