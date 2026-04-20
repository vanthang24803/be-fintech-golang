package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/maynguyen24/sever/internal/models"
)

func TestReportHandler_GetCategorySummary(t *testing.T) {
	t.Parallel()

	svc := &stubReportService{
		breakdownResp: &models.IncomeCategoryBreakdownResponse{},
	}
	h := NewReportHandler(svc)
	app := fiber.New()
	app.Post("/reports/category-summary", func(c *fiber.Ctx) error {
		c.Locals("user_id", int64(42))
		return h.GetCategorySummary(c)
	})

	req := httptest.NewRequest(http.MethodPost, "/reports/category-summary", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}
}

func TestReportHandler_GetCategorySummary_Unauthorized(t *testing.T) {
	t.Parallel()

	h := NewReportHandler(&stubReportService{})
	app := fiber.New()
	app.Post("/reports/category-summary", h.GetCategorySummary)

	req := httptest.NewRequest(http.MethodPost, "/reports/category-summary", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", res.StatusCode)
	}
}

func TestReportHandler_GetMonthlyTrend(t *testing.T) {
	t.Parallel()

	svc := &stubReportService{}
	h := NewReportHandler(svc)
	app := fiber.New()
	app.Post("/reports/monthly-trend", func(c *fiber.Ctx) error {
		c.Locals("user_id", int64(42))
		return h.GetMonthlyTrend(c)
	})

	req := httptest.NewRequest(http.MethodPost, "/reports/monthly-trend", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}
}

func TestReportHandler_GetMonthlyTrend_Unauthorized(t *testing.T) {
	t.Parallel()

	h := NewReportHandler(&stubReportService{})
	app := fiber.New()
	app.Post("/reports/monthly-trend", h.GetMonthlyTrend)

	req := httptest.NewRequest(http.MethodPost, "/reports/monthly-trend", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", res.StatusCode)
	}
}
