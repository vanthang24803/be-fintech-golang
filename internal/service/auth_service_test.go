package service

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/maynguyen24/sever/configs"
	"github.com/maynguyen24/sever/internal/models"
	"github.com/maynguyen24/sever/pkg/apperr"
	jwtUtil "github.com/maynguyen24/sever/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

type stubAuthUserRepo struct {
	getUserByEmailOrUsernameFn func(context.Context, string, string) (*models.User, error)
	getUserByEmailFn           func(context.Context, string) (*models.User, error)
	getUserByGoogleIDFn        func(context.Context, string) (*models.User, error)
	linkGoogleAccountFn        func(context.Context, int64, string) error
	createUserFn               func(context.Context, *models.User) error
	updateProfileFn            func(context.Context, int64, *models.UpdateProfileRequest) (*models.Profile, error)
}

func (s *stubAuthUserRepo) CreateUser(ctx context.Context, user *models.User) error {
	if s.createUserFn != nil {
		return s.createUserFn(ctx, user)
	}
	return nil
}

func (s *stubAuthUserRepo) GetUserByEmailOrUsername(ctx context.Context, email, username string) (*models.User, error) {
	if s.getUserByEmailOrUsernameFn != nil {
		return s.getUserByEmailOrUsernameFn(ctx, email, username)
	}
	return nil, nil
}

func (s *stubAuthUserRepo) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	if s.getUserByEmailFn != nil {
		return s.getUserByEmailFn(ctx, email)
	}
	return nil, nil
}

func (s *stubAuthUserRepo) GetUserByGoogleID(ctx context.Context, googleID string) (*models.User, error) {
	if s.getUserByGoogleIDFn != nil {
		return s.getUserByGoogleIDFn(ctx, googleID)
	}
	return nil, nil
}

func (s *stubAuthUserRepo) GetUserByID(context.Context, int64) (*models.User, error) {
	return nil, nil
}

func (s *stubAuthUserRepo) GetProfileByUserID(context.Context, int64) (*models.Profile, error) {
	return nil, nil
}

func (s *stubAuthUserRepo) UpdateProfile(ctx context.Context, userID int64, req *models.UpdateProfileRequest) (*models.Profile, error) {
	if s.updateProfileFn != nil {
		return s.updateProfileFn(ctx, userID, req)
	}
	return nil, nil
}

func (s *stubAuthUserRepo) LinkGoogleAccount(ctx context.Context, userID int64, googleID string) error {
	if s.linkGoogleAccountFn != nil {
		return s.linkGoogleAccountFn(ctx, userID, googleID)
	}
	return nil
}

type stubTokenRepo struct {
	storeFn  func(context.Context, *models.Token) error
	revokeFn func(context.Context, string) error
	getFn    func(context.Context, string) (*models.Token, error)
}

func (s *stubTokenRepo) StoreRefreshToken(ctx context.Context, token *models.Token) error {
	if s.storeFn != nil {
		return s.storeFn(ctx, token)
	}
	return nil
}

func (s *stubTokenRepo) RevokeToken(ctx context.Context, tokenString string) error {
	if s.revokeFn != nil {
		return s.revokeFn(ctx, tokenString)
	}
	return nil
}

func (s *stubTokenRepo) GetToken(ctx context.Context, tokenString string) (*models.Token, error) {
	if s.getFn != nil {
		return s.getFn(ctx, tokenString)
	}
	return nil, nil
}

func testAuthConfig() *configs.Config {
	return &configs.Config{
		JWTSecret:         "test-jwt-secret",
		JWTRefreshSecret:  "test-refresh-secret",
		GoogleRedirectURL: "http://localhost/callback",
	}
}

func newPasswordHash(t *testing.T, password string) string {
	t.Helper()

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("GenerateFromPassword: %v", err)
	}
	return string(hash)
}

func TestAuthService_Login(t *testing.T) {
	t.Parallel()

	cfg := testAuthConfig()
	user := &models.User{ID: 42, Email: "alice@example.com", Username: "alice", PasswordHash: newPasswordHash(t, "secret123")}

	tests := []struct {
		name      string
		repo      *stubAuthUserRepo
		req       *models.LoginRequest
		wantErrIs error
		wantStore bool
	}{
		{
			name: "lookup error",
			repo: &stubAuthUserRepo{
				getUserByEmailOrUsernameFn: func(context.Context, string, string) (*models.User, error) {
					return nil, errors.New("db down")
				},
			},
			req:       &models.LoginRequest{Identifier: "alice", Password: "secret123"},
			wantErrIs: apperr.ErrUnauthorized,
		},
		{
			name: "password mismatch",
			repo: &stubAuthUserRepo{
				getUserByEmailOrUsernameFn: func(context.Context, string, string) (*models.User, error) {
					return user, nil
				},
			},
			req:       &models.LoginRequest{Identifier: "alice", Password: "wrong"},
			wantErrIs: apperr.ErrUnauthorized,
		},
		{
			name: "success",
			repo: &stubAuthUserRepo{
				getUserByEmailOrUsernameFn: func(context.Context, string, string) (*models.User, error) {
					return user, nil
				},
			},
			req:       &models.LoginRequest{Identifier: "alice", Password: "secret123"},
			wantStore: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			stored := false
			tokenRepo := &stubTokenRepo{
				storeFn: func(ctx context.Context, token *models.Token) error {
					stored = true
					if token == nil || token.UserID != user.ID || token.TokenString == "" {
						t.Fatalf("unexpected token stored: %+v", token)
					}
					return nil
				},
			}

			svc := NewAuthService(tt.repo, tokenRepo, cfg)
			resp, err := svc.Login(context.Background(), tt.req)

			if tt.wantErrIs != nil {
				if !errors.Is(err, tt.wantErrIs) {
					t.Fatalf("expected %v, got %v", tt.wantErrIs, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("Login returned error: %v", err)
			}
			if resp == nil || resp.User == nil || resp.Tokens == nil {
				t.Fatalf("unexpected login response: %+v", resp)
			}
			if resp.Tokens.AccessToken == "" || resp.Tokens.RefreshToken == "" {
				t.Fatalf("expected non-empty token pair: %+v", resp.Tokens)
			}
			if !stored && tt.wantStore {
				t.Fatal("expected refresh token to be stored")
			}
		})
	}
}

func TestAuthService_RefreshToken(t *testing.T) {
	t.Parallel()

	cfg := testAuthConfig()
	userID := int64(42)
	_, oldRefreshToken, err := jwtUtil.GenerateTokenPair(userID, false, cfg)
	if err != nil {
		t.Fatalf("GenerateTokenPair: %v", err)
	}

	tests := []struct {
		name       string
		refresh    string
		getFn      func(context.Context, string) (*models.Token, error)
		revokeFn   func(context.Context, string) error
		wantErrIs  error
		wantResult bool
	}{
		{
			name:    "malformed token",
			refresh: "not-a-jwt",
			getFn: func(context.Context, string) (*models.Token, error) {
				t.Fatal("GetToken should not be called for malformed token")
				return nil, nil
			},
			revokeFn: func(context.Context, string) error {
				t.Fatal("RevokeToken should not be called for malformed token")
				return nil
			},
			wantErrIs: apperr.ErrUnauthorized,
		},
		{
			name:    "revoked session",
			refresh: oldRefreshToken,
			getFn: func(context.Context, string) (*models.Token, error) {
				return nil, errors.New("missing")
			},
			revokeFn: func(context.Context, string) error {
				t.Fatal("RevokeToken should not be called when session is missing")
				return nil
			},
			wantErrIs: apperr.ErrUnauthorized,
		},
		{
			name:    "success",
			refresh: oldRefreshToken,
			getFn: func(context.Context, string) (*models.Token, error) {
				return &models.Token{UserID: userID, TokenString: oldRefreshToken}, nil
			},
			wantResult: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			var revoked []string
			var stored []*models.Token
			tokenRepo := &stubTokenRepo{
				getFn: tt.getFn,
				revokeFn: func(ctx context.Context, tokenString string) error {
					revoked = append(revoked, tokenString)
					if tt.revokeFn != nil {
						return tt.revokeFn(ctx, tokenString)
					}
					return nil
				},
				storeFn: func(ctx context.Context, token *models.Token) error {
					stored = append(stored, token)
					return nil
				},
			}

			svc := NewAuthService(&stubAuthUserRepo{}, tokenRepo, cfg)
			got, err := svc.RefreshToken(context.Background(), &models.RefreshTokenRequest{RefreshToken: tt.refresh})

			if tt.wantErrIs != nil {
				if !errors.Is(err, tt.wantErrIs) {
					t.Fatalf("expected %v, got %v", tt.wantErrIs, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("RefreshToken returned error: %v", err)
			}
			if got == nil || got.AccessToken == "" || got.RefreshToken == "" {
				t.Fatalf("unexpected token pair: %+v", got)
			}
			if len(revoked) != 1 || revoked[0] != tt.refresh {
				t.Fatalf("expected old token to be revoked once, got %v", revoked)
			}
			if len(stored) != 1 || stored[0].TokenString == "" {
				t.Fatalf("expected new refresh token to be stored, got %+v", stored)
			}
		})
	}
}

func TestAuthService_RefreshToken_ExpiredClaimRejected(t *testing.T) {
	t.Parallel()

	cfg := testAuthConfig()
	claims := jwtUtil.TokenClaims{
		UserID:       42,
		FIDOVerified: false,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Minute)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refreshToken, err := token.SignedString([]byte(cfg.JWTRefreshSecret))
	if err != nil {
		t.Fatalf("SignedString: %v", err)
	}

	svc := NewAuthService(&stubAuthUserRepo{}, &stubTokenRepo{
		getFn: func(context.Context, string) (*models.Token, error) {
			t.Fatal("GetToken should not be called for expired refresh token")
			return nil, nil
		},
	}, cfg)

	_, err = svc.RefreshToken(context.Background(), &models.RefreshTokenRequest{RefreshToken: refreshToken})
	if !errors.Is(err, apperr.ErrUnauthorized) {
		t.Fatalf("expected unauthorized, got %v", err)
	}
	if !strings.Contains(err.Error(), "invalid or expired refresh token") {
		t.Fatalf("expected invalid/expired message, got %v", err)
	}
}

func TestAuthService_Logout(t *testing.T) {
	t.Parallel()

	revokeErr := errors.New("db failed")
	var revoked string
	svc := NewAuthService(&stubAuthUserRepo{}, &stubTokenRepo{
		revokeFn: func(ctx context.Context, tokenString string) error {
			revoked = tokenString
			return nil
		},
	}, testAuthConfig())

	if err := svc.Logout(context.Background(), &models.LogoutRequest{RefreshToken: "rt-123"}); err != nil {
		t.Fatalf("Logout returned error: %v", err)
	}
	if revoked != "rt-123" {
		t.Fatalf("expected token to be revoked, got %q", revoked)
	}

	svc = NewAuthService(&stubAuthUserRepo{}, &stubTokenRepo{
		revokeFn: func(context.Context, string) error {
			return revokeErr
		},
	}, testAuthConfig())

	err := svc.Logout(context.Background(), &models.LogoutRequest{RefreshToken: "rt-123"})
	if err == nil || !strings.Contains(err.Error(), "failed to revoke token during logout") {
		t.Fatalf("expected wrapped logout error, got %v", err)
	}
}
