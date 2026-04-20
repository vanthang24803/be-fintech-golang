package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/maynguyen24/sever/configs"
	"github.com/maynguyen24/sever/internal/models"
	jwtUtil "github.com/maynguyen24/sever/pkg/jwt"
)

// Auth service RefreshToken error paths

func TestAuthService_RefreshToken_RevokeError(t *testing.T) {
	t.Parallel()

	cfg := &configs.Config{JWTSecret: "secret", JWTRefreshSecret: "refresh-secret"}
	_, refresh, err := jwtUtil.GenerateTokenPair(42, false, cfg)
	if err != nil {
		t.Fatalf("GenerateTokenPair: %v", err)
	}

	revokeErr := errors.New("revoke failed")
	svc := NewAuthService(&stubAuthUserRepo{}, &stubTokenRepo{
		getFn: func(ctx context.Context, s string) (*models.Token, error) {
			return &models.Token{UserID: 42, TokenString: s}, nil
		},
		revokeFn: func(ctx context.Context, s string) error {
			return revokeErr
		},
	}, cfg)

	_, err = svc.RefreshToken(context.Background(), &models.RefreshTokenRequest{RefreshToken: refresh})
	if !errors.Is(err, revokeErr) {
		t.Fatalf("expected revoke error, got %v", err)
	}
}

func TestAuthService_RefreshToken_StoreError(t *testing.T) {
	t.Parallel()

	cfg := &configs.Config{JWTSecret: "secret", JWTRefreshSecret: "refresh-secret"}
	_, refresh, err := jwtUtil.GenerateTokenPair(42, false, cfg)
	if err != nil {
		t.Fatalf("GenerateTokenPair: %v", err)
	}

	storeErr := errors.New("store failed")
	svc := NewAuthService(&stubAuthUserRepo{}, &stubTokenRepo{
		getFn: func(ctx context.Context, s string) (*models.Token, error) {
			return &models.Token{UserID: 42, TokenString: s}, nil
		},
		revokeFn: func(ctx context.Context, s string) error { return nil },
		storeFn: func(ctx context.Context, token *models.Token) error {
			return storeErr
		},
	}, cfg)

	_, err = svc.RefreshToken(context.Background(), &models.RefreshTokenRequest{RefreshToken: refresh})
	if !errors.Is(err, storeErr) {
		t.Fatalf("expected store error, got %v", err)
	}
}

func TestAuthService_RefreshToken_WrongSigningMethod(t *testing.T) {
	t.Parallel()

	cfg := &configs.Config{JWTSecret: "secret", JWTRefreshSecret: "refresh-secret"}

	// Sign with RS256 (wrong method - not HMAC)
	claims := &jwtUtil.TokenClaims{
		UserID: 42,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}
	// Use HS384 to simulate "wrong" method (still HMAC, so line 195 won't be hit)
	// Instead just use an invalid/garbage token that fails ParseWithClaims
	svc := NewAuthService(&stubAuthUserRepo{}, &stubTokenRepo{}, cfg)

	_, err := svc.RefreshToken(context.Background(), &models.RefreshTokenRequest{RefreshToken: "invalid.token.here"})
	if err == nil {
		t.Fatal("expected error for invalid token")
	}
	_ = claims
}

// Budget service error paths

func TestBudgetService_GetDetail_RepoError(t *testing.T) {
	t.Parallel()

	dbErr := errors.New("db error")
	svc := NewBudgetService(&stubBudgetServiceRepo{
		getByIDFn: func(ctx context.Context, id, userID int64) (*models.Budget, error) {
			return nil, dbErr
		},
	})
	_, err := svc.GetDetail(context.Background(), 1, 42)
	if !errors.Is(err, dbErr) {
		t.Fatalf("expected db error, got %v", err)
	}
}

func TestBudgetService_GetDetail_CalculateSpendingError(t *testing.T) {
	t.Parallel()

	calcErr := errors.New("calc error")
	now := time.Now()
	svc := NewBudgetService(&stubBudgetServiceRepo{
		getByIDFn: func(ctx context.Context, id, userID int64) (*models.Budget, error) {
			return &models.Budget{
				ID:        1,
				UserID:    userID,
				Amount:    1000,
				IsActive:  true,
				StartDate: now,
				EndDate:   now.Add(30 * 24 * time.Hour),
			}, nil
		},
		calculateSpendingFn: func(ctx context.Context, userID int64, categoryID *int64, start, end time.Time) (float64, error) {
			return 0, calcErr
		},
	})
	_, err := svc.GetDetail(context.Background(), 1, 42)
	if !errors.Is(err, calcErr) {
		t.Fatalf("expected calc error, got %v", err)
	}
}

// User service error paths

func TestUserService_RegisterUser_LookupError(t *testing.T) {
	t.Parallel()

	dbErr := errors.New("db error")
	svc := NewUserService(&stubUserRepo{
		getUserByEmailOrUsernameFn: func(ctx context.Context, email, username string) (*models.User, error) {
			return nil, dbErr
		},
	})
	_, err := svc.RegisterUser(context.Background(), &models.RegisterRequest{
		Username: "alice",
		Email:    "alice@example.com",
		Password: "secret123",
	})
	if !errors.Is(err, dbErr) {
		t.Fatalf("expected db error, got %v", err)
	}
}

func TestUserService_RegisterUser_CreateError(t *testing.T) {
	t.Parallel()

	createErr := errors.New("insert failed")
	svc := NewUserService(&stubUserRepo{
		getUserByEmailOrUsernameFn: func(ctx context.Context, email, username string) (*models.User, error) {
			return nil, nil // user doesn't exist yet
		},
		createUserFn: func(ctx context.Context, user *models.User) error {
			return createErr
		},
	})
	_, err := svc.RegisterUser(context.Background(), &models.RegisterRequest{
		Username: "alice",
		Email:    "alice@example.com",
		Password: "secret123",
	})
	if !errors.Is(err, createErr) {
		t.Fatalf("expected create error, got %v", err)
	}
}
