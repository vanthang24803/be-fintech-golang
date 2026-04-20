package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/maynguyen24/sever/internal/models"
)

type stubAuthService struct {
	loginFn          func(context.Context, *models.LoginRequest) (*models.LoginResponse, error)
	refreshFn        func(context.Context, *models.RefreshTokenRequest) (*models.TokenPair, error)
	logoutFn         func(context.Context, *models.LogoutRequest) error
	getGoogleURLFn   func(context.Context) string
	googleCallbackFn func(context.Context, string) (*models.LoginResponse, error)
}

func (s *stubAuthService) Login(ctx context.Context, req *models.LoginRequest) (*models.LoginResponse, error) {
	if s.loginFn != nil {
		return s.loginFn(ctx, req)
	}
	return &models.LoginResponse{
		User:   &models.User{ID: 1, Username: "alice"},
		Tokens: &models.TokenPair{AccessToken: "at", RefreshToken: "rt"},
	}, nil
}

func (s *stubAuthService) RefreshToken(ctx context.Context, req *models.RefreshTokenRequest) (*models.TokenPair, error) {
	if s.refreshFn != nil {
		return s.refreshFn(ctx, req)
	}
	return &models.TokenPair{AccessToken: "new-at", RefreshToken: "new-rt"}, nil
}

func (s *stubAuthService) Logout(ctx context.Context, req *models.LogoutRequest) error {
	if s.logoutFn != nil {
		return s.logoutFn(ctx, req)
	}
	return nil
}

func (s *stubAuthService) GetGoogleAuthURL(ctx context.Context) string {
	if s.getGoogleURLFn != nil {
		return s.getGoogleURLFn(ctx)
	}
	return "https://accounts.google.com/o/oauth2/auth"
}

func (s *stubAuthService) HandleGoogleCallback(ctx context.Context, code string) (*models.LoginResponse, error) {
	if s.googleCallbackFn != nil {
		return s.googleCallbackFn(ctx, code)
	}
	return &models.LoginResponse{
		User:   &models.User{ID: 1},
		Tokens: &models.TokenPair{AccessToken: "at", RefreshToken: "rt"},
	}, nil
}

func TestAuthHandler_Login(t *testing.T) {
	t.Parallel()

	h := NewAuthHandler(&stubAuthService{})
	app := fiber.New()
	app.Post("/login", h.Login)

	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`{"identifier":"alice","password":"secret123"}`))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}
}

func TestAuthHandler_Login_ValidationFail(t *testing.T) {
	t.Parallel()

	h := NewAuthHandler(&stubAuthService{})
	app := fiber.New()
	app.Post("/login", h.Login)

	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`{"identifier":"alice"}`))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", res.StatusCode)
	}
}

func TestAuthHandler_Refresh(t *testing.T) {
	t.Parallel()

	h := NewAuthHandler(&stubAuthService{})
	app := fiber.New()
	app.Post("/refresh", h.Refresh)

	req := httptest.NewRequest(http.MethodPost, "/refresh", strings.NewReader(`{"refresh_token":"old-rt"}`))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}
}

func TestAuthHandler_Refresh_ValidationFail(t *testing.T) {
	t.Parallel()

	h := NewAuthHandler(&stubAuthService{})
	app := fiber.New()
	app.Post("/refresh", h.Refresh)

	req := httptest.NewRequest(http.MethodPost, "/refresh", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", res.StatusCode)
	}
}

func TestAuthHandler_Logout(t *testing.T) {
	t.Parallel()

	h := NewAuthHandler(&stubAuthService{})
	app := fiber.New()
	app.Post("/logout", h.Logout)

	req := httptest.NewRequest(http.MethodPost, "/logout", strings.NewReader(`{"refresh_token":"rt-123"}`))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}
}

func TestAuthHandler_Logout_ValidationFail(t *testing.T) {
	t.Parallel()

	h := NewAuthHandler(&stubAuthService{})
	app := fiber.New()
	app.Post("/logout", h.Logout)

	req := httptest.NewRequest(http.MethodPost, "/logout", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", res.StatusCode)
	}
}

func TestAuthHandler_GetGoogleAuthURL(t *testing.T) {
	t.Parallel()

	h := NewAuthHandler(&stubAuthService{})
	app := fiber.New()
	app.Get("/auth/google", h.GetGoogleAuthURL)

	req := httptest.NewRequest(http.MethodGet, "/auth/google", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}
}

func TestAuthHandler_GoogleCallback(t *testing.T) {
	t.Parallel()

	h := NewAuthHandler(&stubAuthService{})
	app := fiber.New()
	app.Post("/auth/callback", h.GoogleCallback)

	req := httptest.NewRequest(http.MethodPost, "/auth/callback", strings.NewReader(`{"code":"auth-code-123"}`))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}
}

func TestAuthHandler_GoogleCallback_MissingCode(t *testing.T) {
	t.Parallel()

	h := NewAuthHandler(&stubAuthService{})
	app := fiber.New()
	app.Post("/auth/callback", h.GoogleCallback)

	req := httptest.NewRequest(http.MethodPost, "/auth/callback", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", res.StatusCode)
	}
}

func TestAuthHandler_GoogleCallback_QueryParam(t *testing.T) {
	t.Parallel()

	h := NewAuthHandler(&stubAuthService{})
	app := fiber.New()
	app.Get("/auth/callback", h.GoogleCallback)

	req := httptest.NewRequest(http.MethodGet, "/auth/callback?code=my-code", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}
}
