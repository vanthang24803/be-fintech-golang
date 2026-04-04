package service

import (
	"context"
	"fmt"
	"time"

	"github.com/maynguyen24/sever/internal/models"
)

// ReportRepository defines the DB contract for this service
type ReportRepository interface {
	GetCategorySummary(ctx context.Context, userID int64, start, end time.Time) ([]*models.CategorySummary, error)
	GetMonthlyTrend(ctx context.Context, userID int64, since time.Time) ([]*models.MonthlySummary, error)
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
