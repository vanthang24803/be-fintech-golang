package middleware

import (
	"encoding/json"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/maynguyen24/sever/pkg/i18n"
	jwtUtil "github.com/maynguyen24/sever/pkg/jwt"
)

func TestFIDOMiddleware(t *testing.T) {
	t.Parallel()

	if err := i18n.LoadLocales(filepath.Join("..", "..", "locales")); err != nil {
		t.Fatalf("LoadLocales(): %v", err)
	}

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		mode := c.Query("mode")
		switch mode {
		case "verified":
			c.Locals("user", &jwtUtil.TokenClaims{UserID: 42, FIDOVerified: true})
		case "unverified":
			c.Locals("user", &jwtUtil.TokenClaims{UserID: 42, FIDOVerified: false})
		}
		return c.Next()
	})
	app.Use(FIDOMiddleware())
	app.Get("/", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusNoContent) })

	tests := []struct {
		name       string
		query      string
		wantStatus int
		wantCode   float64
	}{
		{name: "missing claims", query: "", wantStatus: fiber.StatusUnauthorized, wantCode: 4010},
		{name: "not verified", query: "?mode=unverified", wantStatus: fiber.StatusForbidden, wantCode: 4031},
		{name: "verified", query: "?mode=verified", wantStatus: fiber.StatusNoContent},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			resp, err := app.Test(httptest.NewRequest(fiber.MethodGet, "/"+tt.query, nil))
			if err != nil {
				t.Fatalf("app.Test(): %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.wantStatus {
				t.Fatalf("expected status %d, got %d", tt.wantStatus, resp.StatusCode)
			}
			if tt.wantCode == 0 {
				return
			}
			var body map[string]any
			if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
				t.Fatalf("Decode(): %v", err)
			}
			if body["code"] != tt.wantCode {
				t.Fatalf("expected code %v, got %v", tt.wantCode, body["code"])
			}
		})
	}
}
