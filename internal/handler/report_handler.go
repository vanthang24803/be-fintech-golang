package handler

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/maynguyen24/sever/internal/models"
	"github.com/maynguyen24/sever/pkg/response"
)

// ReportService defines the contract for the handler layer
type ReportService interface {
	GetCategorySummary(ctx context.Context, userID int64, req *models.ReportRequest) ([]*models.CategorySummary, error)
	GetMonthlyTrend(ctx context.Context, userID int64, months int) ([]*models.MonthlySummary, error)
	GetIncomeCategoryBreakdown(ctx context.Context, userID int64, req *models.IncomeCategoryBreakdownRequest) (*models.IncomeCategoryBreakdownResponse, error)
	GetCategoryTrend(ctx context.Context, userID int64, req *models.CategoryTrendRequest) (*models.CategoryTrendResponse, error)
}

type ReportHandler struct {
	service ReportService
}

func NewReportHandler(service ReportService) *ReportHandler {
	return &ReportHandler{service: service}
}

// POST /api/v1/reports/category-summary
func (h *ReportHandler) GetCategorySummary(c *fiber.Ctx) error {
	userID, err := extractUserID(c)
	if err != nil {
		return err
	}

	var req models.ReportRequest
	if err := c.BodyParser(&req); err != nil {
		// Ignore error, use default range
	}

	summary, err := h.service.GetCategorySummary(c.Context(), userID, &req)
	if err != nil {
		return err
	}

	return response.Success(c, 2000, "Category summary fetched successfully", summary)
}

// POST /api/v1/reports/monthly-trend
func (h *ReportHandler) GetMonthlyTrend(c *fiber.Ctx) error {
	userID, err := extractUserID(c)
	if err != nil {
		return err
	}

	var req models.MonthlyTrendRequest
	if err := c.BodyParser(&req); err != nil {
		// Ignore error, use default count
	}

	trend, err := h.service.GetMonthlyTrend(c.Context(), userID, req.Months)
	if err != nil {
		return err
	}

	return response.Success(c, 2000, "Monthly trend fetched successfully", trend)
}

// POST /api/v1/reports/income-category-breakdown
func (h *ReportHandler) GetIncomeCategoryBreakdown(c *fiber.Ctx) error {
	userID, err := extractUserID(c)
	if err != nil {
		return err
	}

	var req models.IncomeCategoryBreakdownRequest
	if err := c.BodyParser(&req); err != nil {
		// Ignore error, use default range
	}

	result, err := h.service.GetIncomeCategoryBreakdown(c.Context(), userID, &req)
	if err != nil {
		return err
	}

	return response.Success(c, 2000, "Income category breakdown fetched successfully", result)
}

// POST /api/v1/reports/category-trend
func (h *ReportHandler) GetCategoryTrend(c *fiber.Ctx) error {
	userID, err := extractUserID(c)
	if err != nil {
		return err
	}

	var req models.CategoryTrendRequest
	if err := c.BodyParser(&req); err != nil {
		// Ignore error, use default range
	}

	result, err := h.service.GetCategoryTrend(c.Context(), userID, &req)
	if err != nil {
		return err
	}

	return response.Success(c, 2000, "Category trend fetched successfully", result)
}
