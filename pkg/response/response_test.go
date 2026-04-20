package response

import (
	"encoding/json"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/maynguyen24/sever/pkg/i18n"
)

func TestSuccessUsesLocalizedMessageAndDefaultLanguage(t *testing.T) {
	if err := i18n.LoadLocales(filepath.Join("..", "..", "locales")); err != nil {
		t.Fatalf("LoadLocales(): %v", err)
	}

	app := fiber.New()
	app.Get("/default", func(c *fiber.Ctx) error {
		return Success(c, 1000, "auth.login_success", fiber.Map{"ok": true})
	})
	app.Get("/vi", func(c *fiber.Ctx) error {
		c.Locals("lang", "vi")
		return Success(c, 1001, "auth.login_success", nil)
	})

	req := httptest.NewRequest(fiber.MethodGet, "/default", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test(default): %v", err)
	}
	defer resp.Body.Close()

	var body Response
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("Decode(default): %v", err)
	}
	if body.Code != 1000 || body.Message != "Login successful" {
		t.Fatalf("unexpected success response: %+v", body)
	}

	req = httptest.NewRequest(fiber.MethodGet, "/vi", nil)
	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("app.Test(vi): %v", err)
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("Decode(vi): %v", err)
	}
	if body.Message != "Đăng nhập thành công" {
		t.Fatalf("expected vietnamese message, got %+v", body)
	}
}

func TestErrorUsesLocalizedMessageAndStatus(t *testing.T) {
	if err := i18n.LoadLocales(filepath.Join("..", "..", "locales")); err != nil {
		t.Fatalf("LoadLocales(): %v", err)
	}

	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		c.Locals("lang", "vi")
		return Error(c, fiber.StatusBadRequest, 4001, "common.bad_request", "bad input")
	})

	req := httptest.NewRequest(fiber.MethodGet, "/", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test(): %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", fiber.StatusBadRequest, resp.StatusCode)
	}

	var body Response
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("Decode(): %v", err)
	}
	if body.Code != 4001 || body.Message != "Yêu cầu không hợp lệ" || body.Error != "bad input" {
		t.Fatalf("unexpected error response: %+v", body)
	}
}
