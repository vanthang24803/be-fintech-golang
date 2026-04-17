package service

import (
	"context"
	"fmt"
	"time"

	"github.com/maynguyen24/sever/internal/models"
	"github.com/maynguyen24/sever/pkg/apperr"
)

// ReportRepository defines the DB contract for this service
type ReportRepository interface {
	GetCategorySummary(ctx context.Context, userID int64, start, end time.Time) ([]*models.CategorySummary, error)
	GetMonthlyTrend(ctx context.Context, userID int64, since time.Time) ([]*models.MonthlySummary, error)
	GetIncomeCategoryBreakdown(ctx context.Context, userID int64, start, end time.Time, limit int) ([]*models.IncomeCategoryBreakdownItem, error)
	GetCategoryTrend(ctx context.Context, userID, categoryID int64, start, end time.Time, granularity string) ([]*models.CategoryTrendPoint, error)
}

type ReportService struct {
	repo ReportRepository
}

func NewReportService(repo ReportRepository) *ReportService {
	return &ReportService{repo: repo}
}

// GetCategorySummary provides expenses aggregated by category, including percentages
func (s *ReportService) GetCategorySummary(ctx context.Context, userID int64, req *models.ReportRequest) ([]*models.CategorySummary, error) {
	// Default to current month if not provided
	if req.StartDate.IsZero() {
		now := time.Now()
		req.StartDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		req.EndDate = now
	}

	summary, err := s.repo.GetCategorySummary(ctx, userID, req.StartDate, req.EndDate)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch category summary: %w", err)
	}

	// Calculate total amount to compute percentages
	var total float64
	for _, s := range summary {
		total += s.TotalAmount
	}

	for _, s := range summary {
		if total > 0 {
			s.Percentage = (s.TotalAmount / total) * 100
		}
	}

	return summary, nil
}

// GetMonthlyTrend provides month-by-month income/expense trends
func (s *ReportService) GetMonthlyTrend(ctx context.Context, userID int64, months int) ([]*models.MonthlySummary, error) {
	if months <= 0 {
		months = 6 // Default to 6 months
	}

	since := time.Now().AddDate(0, -months, 0)
	trend, err := s.repo.GetMonthlyTrend(ctx, userID, since)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch monthly trend: %w", err)
	}

	for _, m := range trend {
		m.NetProfit = m.Income - m.Expense
	}

	return trend, nil
}

func (s *ReportService) GetIncomeCategoryBreakdown(ctx context.Context, userID int64, req *models.IncomeCategoryBreakdownRequest) (*models.IncomeCategoryBreakdownResponse, error) {
	if req == nil {
		req = &models.IncomeCategoryBreakdownRequest{}
	}

	start, end, err := normalizeReportRange(req.StartDate, req.EndDate)
	if err != nil {
		return nil, err
	}

	items, err := s.repo.GetIncomeCategoryBreakdown(ctx, userID, start, end, req.Limit)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch income category breakdown: %w", err)
	}

	resp := &models.IncomeCategoryBreakdownResponse{
		RangeStart: start,
		RangeEnd:   end,
		Items:      items,
	}

	for _, item := range items {
		resp.TotalIncome += item.Amount
	}
	for _, item := range items {
		if resp.TotalIncome > 0 {
			item.Percentage = (item.Amount / resp.TotalIncome) * 100
		}
	}

	return resp, nil
}

func (s *ReportService) GetCategoryTrend(ctx context.Context, userID int64, req *models.CategoryTrendRequest) (*models.CategoryTrendResponse, error) {
	if req == nil {
		req = &models.CategoryTrendRequest{}
	}
	if req.CategoryID == 0 {
		return nil, fmt.Errorf("%w: category_id is required", apperr.ErrInvalidInput)
	}

	start, end, err := normalizeReportRange(req.StartDate, req.EndDate)
	if err != nil {
		return nil, err
	}

	granularity := req.Granularity
	if granularity == "" {
		granularity = "day"
	}
	if granularity != "day" {
		return nil, fmt.Errorf("%w: unsupported granularity", apperr.ErrInvalidInput)
	}

	points, err := s.repo.GetCategoryTrend(ctx, userID, req.CategoryID, start, end, granularity)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch category trend: %w", err)
	}

	return &models.CategoryTrendResponse{
		CategoryID:  req.CategoryID,
		Granularity: granularity,
		Points:      points,
	}, nil
}

func normalizeReportRange(start, end time.Time) (time.Time, time.Time, error) {
	if start.IsZero() && end.IsZero() {
		end = time.Now()
		start = end.AddDate(0, 0, -30)
	}
	if start.IsZero() != end.IsZero() {
		return time.Time{}, time.Time{}, fmt.Errorf("%w: start_date and end_date must be provided together", apperr.ErrInvalidInput)
	}

	start = start.UTC()
	end = end.UTC()
	if start.After(end) {
		return time.Time{}, time.Time{}, fmt.Errorf("%w: start_date must be before or equal to end_date", apperr.ErrInvalidInput)
	}
	if end.Sub(start) > 366*24*time.Hour {
		return time.Time{}, time.Time{}, fmt.Errorf("%w: date range must not exceed 366 days", apperr.ErrInvalidInput)
	}

	return start, end, nil
}
