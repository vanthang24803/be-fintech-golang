package middleware

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/maynguyen24/sever/pkg/logger"
	"go.uber.org/zap"
)

func TestRequestLogger(t *testing.T) {
	t.Parallel()

	logger.Log = zap.NewNop()

	app := fiber.New()
	app.Use(RequestLogger())
	app.Get("/ok", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusNoContent) })
	app.Get("/err", func(c *fiber.Ctx) error { return errors.New("boom") })

	resp, err := app.Test(httptest.NewRequest(fiber.MethodGet, "/ok", nil))
	if err != nil {
		t.Fatalf("app.Test(/ok): %v", err)
	}
	resp.Body.Close()

	resp, err = app.Test(httptest.NewRequest(fiber.MethodGet, "/err", nil))
	if err != nil {
		t.Fatalf("app.Test(/err): %v", err)
	}
	resp.Body.Close()
	if resp.StatusCode != fiber.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", resp.StatusCode)
	}
}
