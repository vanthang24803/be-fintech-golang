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
