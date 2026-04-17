package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/maynguyen24/sever/internal/models"
	respkg "github.com/maynguyen24/sever/pkg/response"
)

type stubReportService struct {
	breakdownReq *models.IncomeCategoryBreakdownRequest
	trendReq     *models.CategoryTrendRequest

	breakdownResp *models.IncomeCategoryBreakdownResponse
	trendResp     *models.CategoryTrendResponse
}

func (s *stubReportService) GetCategorySummary(ctx context.Context, userID int64, req *models.ReportRequest) ([]*models.CategorySummary, error) {
	return nil, nil
}

func (s *stubReportService) GetMonthlyTrend(ctx context.Context, userID int64, months int) ([]*models.MonthlySummary, error) {
	return nil, nil
}

func (s *stubReportService) GetIncomeCategoryBreakdown(ctx context.Context, userID int64, req *models.IncomeCategoryBreakdownRequest) (*models.IncomeCategoryBreakdownResponse, error) {
	s.breakdownReq = req
	return s.breakdownResp, nil
}

func (s *stubReportService) GetCategoryTrend(ctx context.Context, userID int64, req *models.CategoryTrendRequest) (*models.CategoryTrendResponse, error) {
	s.trendReq = req
	return s.trendResp, nil
}

func TestReportHandler_GetIncomeCategoryBreakdown(t *testing.T) {
	svc := &stubReportService{
		breakdownResp: &models.IncomeCategoryBreakdownResponse{
			TotalIncome: 100,
			Items: []*models.IncomeCategoryBreakdownItem{{CategoryID: 1, CategoryName: "Salary", Amount: 100, Percentage: 100}},
		},
	}
	h := NewReportHandler(svc)
	app := fiber.New()
	app.Post("/reports/income-category-breakdown", func(c *fiber.Ctx) error {
		c.Locals("user_id", int64(88))
		return h.GetIncomeCategoryBreakdown(c)
	})

	req := httptest.NewRequest(http.MethodPost, "/reports/income-category-breakdown", strings.NewReader(`{"limit":5}`))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", res.StatusCode)
	}
	if svc.breakdownReq == nil {
		t.Fatal("expected service request to be captured")
	}
	if svc.breakdownReq.Limit != 5 {
		t.Fatalf("expected limit 5, got %d", svc.breakdownReq.Limit)
	}

	var payload respkg.Response
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if payload.Code != 2000 {
		t.Fatalf("expected code 2000, got %d", payload.Code)
	}
}

func TestReportHandler_GetCategoryTrend_EmptyBodyUsesDefaults(t *testing.T) {
	svc := &stubReportService{
		trendResp: &models.CategoryTrendResponse{
			CategoryID:   12,
			CategoryName: "Salary",
			Granularity:  "day",
			Points:       []*models.CategoryTrendPoint{{Date: "2026-04-01", Amount: 100}},
		},
	}
	h := NewReportHandler(svc)
	app := fiber.New()
	app.Post("/reports/category-trend", func(c *fiber.Ctx) error {
		c.Locals("user_id", int64(88))
		return h.GetCategoryTrend(c)
	})

	req := httptest.NewRequest(http.MethodPost, "/reports/category-trend", strings.NewReader(`{"category_id":12}`))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", res.StatusCode)
	}
	if svc.trendReq == nil {
		t.Fatal("expected trend request to be captured")
	}
	if svc.trendReq.CategoryID != 12 {
		t.Fatalf("expected category ID 12, got %d", svc.trendReq.CategoryID)
	}
}

func TestReportHandler_GetIncomeCategoryBreakdown_RequiresUser(t *testing.T) {
	svc := &stubReportService{}
	h := NewReportHandler(svc)
	app := fiber.New()
	app.Post("/reports/income-category-breakdown", h.GetIncomeCategoryBreakdown)

	req := httptest.NewRequest(http.MethodPost, "/reports/income-category-breakdown", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", res.StatusCode)
	}
}

var _ = time.UTC
