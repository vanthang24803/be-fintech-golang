package server

import (
	"encoding/json"
	"errors"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/maynguyen24/sever/configs"
	"github.com/maynguyen24/sever/pkg/apperr"
	"github.com/maynguyen24/sever/pkg/i18n"
	"github.com/maynguyen24/sever/pkg/logger"
	"go.uber.org/zap"
)

func testServerConfig() *configs.Config {
	return &configs.Config{
		Port:                    "9999",
		JWTSecret:               "jwt-secret",
		JWTRefreshSecret:        "refresh-secret",
		GoogleRedirectURL:       "http://localhost/callback",
		WebAuthnRPID:            "localhost",
		WebAuthnRPName:          "Test",
		WebAuthnOrigin:          "http://localhost:3000",
		RedisAddr:               "127.0.0.1:0",
		RedisPassword:           "",
		FirebaseCredentialsPath: "",
	}
}

func TestCustomErrorHandler(t *testing.T) {
	t.Parallel()

	if err := i18n.LoadLocales(filepath.Join("..", "..", "locales")); err != nil {
		t.Fatalf("LoadLocales(): %v", err)
	}
	logger.Log = zap.NewNop()

	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantCode   float64
	}{
		{name: "not found", err: apperr.ErrNotFound, wantStatus: fiber.StatusNotFound, wantCode: 4040},
		{name: "conflict", err: apperr.ErrConflict, wantStatus: fiber.StatusConflict, wantCode: 4090},
		{name: "invalid input", err: apperr.ErrInvalidInput, wantStatus: fiber.StatusBadRequest, wantCode: 4000},
		{name: "unauthorized", err: apperr.ErrUnauthorized, wantStatus: fiber.StatusUnauthorized, wantCode: 4010},
		{name: "fiber error", err: fiber.NewError(fiber.StatusTeapot, "short and stout"), wantStatus: fiber.StatusTeapot, wantCode: 4180},
		{name: "generic error", err: errors.New("boom"), wantStatus: fiber.StatusInternalServerError, wantCode: 5000},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New(fiber.Config{ErrorHandler: customErrorHandler})
			app.Get("/", func(c *fiber.Ctx) error { return tt.err })

			resp, err := app.Test(httptest.NewRequest(fiber.MethodGet, "/", nil))
			if err != nil {
				t.Fatalf("app.Test(): %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.wantStatus {
				t.Fatalf("expected status %d, got %d", tt.wantStatus, resp.StatusCode)
			}
			var body map[string]any
			if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
				t.Fatalf("Decode(): %v", err)
			}
			if body["code"] != tt.wantCode {
				t.Fatalf("expected business code %v, got %v", tt.wantCode, body["code"])
			}
		})
	}
}

func TestServerOptionsAndRoutes(t *testing.T) {
	t.Parallel()

	logger.Log = zap.NewNop()
	cfg := testServerConfig()
	db := &sqlx.DB{}

	srv := NewServer(WithPort("9090"), WithConfig(cfg), WithDB(db))
	if srv.port != "9090" || srv.cfg != cfg || srv.db != db {
		t.Fatalf("unexpected server options: %+v", srv)
	}
	if srv.App == nil {
		t.Fatal("expected fiber app to be initialized")
	}

	resp, err := srv.App.Test(httptest.NewRequest(fiber.MethodGet, "/", nil))
	if err != nil {
		t.Fatalf("app.Test(/): %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected root status 200, got %d", resp.StatusCode)
	}

	resp, err = srv.App.Test(httptest.NewRequest(fiber.MethodGet, "/error", nil))
	if err != nil {
		t.Fatalf("app.Test(/error): %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("expected /error status 400, got %d", resp.StatusCode)
	}
}
