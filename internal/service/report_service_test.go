package service

import (
	"context"
	"testing"
	"time"

	"github.com/maynguyen24/sever/internal/models"
)

type stubReportRepository struct {
	breakdownStart time.Time
	breakdownEnd   time.Time
	trendStart     time.Time
	trendEnd       time.Time
	trendCategory  int64
	trendGran      string

	breakdownResp []*models.IncomeCategoryBreakdownItem
	trendResp     []*models.CategoryTrendPoint
}

func (s *stubReportRepository) GetCategorySummary(ctx context.Context, userID int64, start, end time.Time) ([]*models.CategorySummary, error) {
	return nil, nil
}

func (s *stubReportRepository) GetMonthlyTrend(ctx context.Context, userID int64, since time.Time) ([]*models.MonthlySummary, error) {
	return nil, nil
}

func (s *stubReportRepository) GetIncomeCategoryBreakdown(ctx context.Context, userID int64, start, end time.Time, limit int) ([]*models.IncomeCategoryBreakdownItem, error) {
	s.breakdownStart = start
	s.breakdownEnd = end
	return s.breakdownResp, nil
}

func (s *stubReportRepository) GetCategoryTrend(ctx context.Context, userID, categoryID int64, start, end time.Time, granularity string) ([]*models.CategoryTrendPoint, error) {
	s.trendCategory = categoryID
	s.trendStart = start
	s.trendEnd = end
	s.trendGran = granularity
	return s.trendResp, nil
}

func TestReportService_GetIncomeCategoryBreakdown_DefaultsLast30Days(t *testing.T) {
	repo := &stubReportRepository{}
	svc := NewReportService(repo)

	resp, err := svc.GetIncomeCategoryBreakdown(context.Background(), 99, &models.IncomeCategoryBreakdownRequest{})
	if err != nil {
		t.Fatalf("GetIncomeCategoryBreakdown() error = %v", err)
	}

	if resp == nil {
		t.Fatal("expected non-nil response")
	}

	if repo.breakdownStart.IsZero() || repo.breakdownEnd.IsZero() {
		t.Fatal("expected repository to receive default date range")
	}

	got := repo.breakdownEnd.Sub(repo.breakdownStart)
	if got < (29*24*time.Hour) || got > (31*24*time.Hour) {
		t.Fatalf("expected default range near 30 days, got %v", got)
	}
}

func TestReportService_GetIncomeCategoryBreakdown_ComputesPercentages(t *testing.T) {
	repo := &stubReportRepository{
		breakdownResp: []*models.IncomeCategoryBreakdownItem{
			{CategoryID: 1, CategoryName: "Salary", Amount: 80},
			{CategoryID: 2, CategoryName: "Bonus", Amount: 20},
		},
	}
	svc := NewReportService(repo)

	resp, err := svc.GetIncomeCategoryBreakdown(context.Background(), 99, &models.IncomeCategoryBreakdownRequest{
		StartDate: time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2026, 4, 30, 0, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("GetIncomeCategoryBreakdown() error = %v", err)
	}

	if resp.TotalIncome != 100 {
		t.Fatalf("expected total income 100, got %v", resp.TotalIncome)
	}
	if resp.Items[0].Percentage != 80 {
		t.Fatalf("expected first percentage 80, got %v", resp.Items[0].Percentage)
	}
	if resp.Items[1].Percentage != 20 {
		t.Fatalf("expected second percentage 20, got %v", resp.Items[1].Percentage)
	}
}

func TestReportService_GetIncomeCategoryBreakdown_RejectsInvalidRange(t *testing.T) {
	repo := &stubReportRepository{}
	svc := NewReportService(repo)

	_, err := svc.GetIncomeCategoryBreakdown(context.Background(), 99, &models.IncomeCategoryBreakdownRequest{
		StartDate: time.Date(2026, 4, 30, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC),
	})
	if err == nil {
		t.Fatal("expected error for invalid range")
	}
}

func TestReportService_GetCategoryTrend_DefaultsAndValidates(t *testing.T) {
	repo := &stubReportRepository{}
	svc := NewReportService(repo)

	resp, err := svc.GetCategoryTrend(context.Background(), 99, &models.CategoryTrendRequest{CategoryID: 12})
	if err != nil {
		t.Fatalf("GetCategoryTrend() error = %v", err)
	}
	if resp.Granularity != "day" {
		t.Fatalf("expected default granularity day, got %q", resp.Granularity)
	}
	if repo.trendCategory != 12 {
		t.Fatalf("expected category 12, got %d", repo.trendCategory)
	}
	if repo.trendStart.IsZero() || repo.trendEnd.IsZero() {
		t.Fatal("expected repository to receive default trend range")
	}

	_, err = svc.GetCategoryTrend(context.Background(), 99, &models.CategoryTrendRequest{
		CategoryID:   12,
		StartDate:    time.Date(2026, 4, 30, 0, 0, 0, 0, time.UTC),
		EndDate:      time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC),
		Granularity: "day",
	})
	if err == nil {
		t.Fatal("expected error for invalid trend range")
	}
}
