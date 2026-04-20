package middleware

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/maynguyen24/sever/configs"
	jwtUtil "github.com/maynguyen24/sever/pkg/jwt"
)

func testMiddlewareConfig() *configs.Config {
	return &configs.Config{
		JWTSecret:        "test-jwt-secret",
		JWTRefreshSecret: "test-refresh-secret",
	}
}

func TestAuthRequired(t *testing.T) {
	t.Parallel()

	cfg := testMiddlewareConfig()
	accessToken, _, err := jwtUtil.GenerateTokenPair(42, true, cfg)
	if err != nil {
		t.Fatalf("GenerateTokenPair(): %v", err)
	}

	app := fiber.New()
	app.Use(AuthRequired(cfg))
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"user_id": c.Locals("user_id"),
		})
	})

	tests := []struct {
		name         string
		authHeader   string
		wantStatus   int
		wantContains string
	}{
		{name: "missing header", wantStatus: fiber.StatusUnauthorized, wantContains: "Missing Authorization header"},
		{name: "invalid format", authHeader: "Token abc", wantStatus: fiber.StatusUnauthorized, wantContains: "Invalid Authorization format"},
		{name: "invalid token", authHeader: "Bearer invalid", wantStatus: fiber.StatusUnauthorized, wantContains: "Invalid or expired token"},
		{name: "success", authHeader: "Bearer " + accessToken, wantStatus: fiber.StatusOK},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(fiber.MethodGet, "/", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("app.Test(): %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.wantStatus {
				t.Fatalf("expected status %d, got %d", tt.wantStatus, resp.StatusCode)
			}
			if tt.wantContains != "" {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					t.Fatalf("ReadAll(): %v", err)
				}
				if string(body) != tt.wantContains {
					t.Fatalf("expected message %q, got %q", tt.wantContains, string(body))
				}
				return
			}

			var body map[string]any
			if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
				t.Fatalf("Decode(): %v", err)
			}
			if body["user_id"] != float64(42) {
				t.Fatalf("expected user_id 42, got %v", body["user_id"])
			}
		})
	}
}
