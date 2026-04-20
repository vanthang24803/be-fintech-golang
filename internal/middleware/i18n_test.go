package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestI18nMiddleware(t *testing.T) {
	t.Parallel()

	app := fiber.New()
	app.Use(I18nMiddleware())
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString(c.Locals("lang").(string))
	})

	tests := []struct {
		name       string
		xLang      string
		acceptLang string
		want       string
	}{
		{name: "custom header wins", xLang: "vi", acceptLang: "en-US,en;q=0.9", want: "vi"},
		{name: "accept language normalized", acceptLang: "vi-VN,vi;q=0.9", want: "vi"},
		{name: "unsupported falls back", xLang: "fr", want: "en"},
		{name: "missing falls back", want: "en"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(fiber.MethodGet, "/", nil)
			if tt.xLang != "" {
				req.Header.Set("x-lang", tt.xLang)
			}
			if tt.acceptLang != "" {
				req.Header.Set("Accept-Language", tt.acceptLang)
			}

			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("app.Test(): %v", err)
			}
			defer resp.Body.Close()

			buf := make([]byte, 8)
			n, _ := resp.Body.Read(buf)
			if got := string(buf[:n]); got != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, got)
			}
		})
	}
}
