package service

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/maynguyen24/sever/internal/models"
)

type stubReportRepository struct {
	categorySummaryStart time.Time
	categorySummaryEnd   time.Time
	categorySummaryErr   error
	categorySummaryResp  []*models.CategorySummary

	monthlyTrendSince time.Time
	monthlyTrendErr   error
	monthlyTrendResp  []*models.MonthlySummary
}

func (s *stubReportRepository) GetCategorySummary(ctx context.Context, userID int64, start, end time.Time) ([]*models.CategorySummary, error) {
	s.categorySummaryStart = start
	s.categorySummaryEnd = end
	return s.categorySummaryResp, s.categorySummaryErr
}

func (s *stubReportRepository) GetMonthlyTrend(ctx context.Context, userID int64, since time.Time) ([]*models.MonthlySummary, error) {
	s.monthlyTrendSince = since
	return s.monthlyTrendResp, s.monthlyTrendErr
}

func TestReportService_GetCategorySummary_DefaultsCurrentMonthAndComputesPercentages(t *testing.T) {
	repo := &stubReportRepository{
		categorySummaryResp: []*models.CategorySummary{
			{CategoryID: 1, CategoryName: "Food", TotalAmount: 25},
			{CategoryID: 2, CategoryName: "Rent", TotalAmount: 75},
		},
	}
	svc := NewReportService(repo)

	fixedNow := time.Date(2026, 4, 18, 10, 30, 0, 0, time.UTC)
	prev := reportNow
	reportNow = func() time.Time { return fixedNow }
	defer func() { reportNow = prev }()

	resp, err := svc.GetCategorySummary(context.Background(), 42, &models.ReportRequest{})
	if err != nil {
		t.Fatalf("GetCategorySummary() error = %v", err)
	}

	if got, want := len(resp), 2; got != want {
		t.Fatalf("expected %d summary items, got %d", want, got)
	}
	if repo.categorySummaryStart.IsZero() || repo.categorySummaryEnd.IsZero() {
		t.Fatal("expected repository to receive default date range")
	}

	if got := repo.categorySummaryStart.Day(); got != 1 {
		t.Fatalf("expected default start date to be first day of month, got day %d", got)
	}
	if !repo.categorySummaryStart.Equal(time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC)) {
		t.Fatalf("expected default start date to be 2026-04-01, got %s", repo.categorySummaryStart.Format(time.RFC3339Nano))
	}
	if !repo.categorySummaryEnd.Equal(fixedNow) {
		t.Fatalf("expected default end date to match fixed clock, got %s", repo.categorySummaryEnd.Format(time.RFC3339Nano))
	}
	if got, want := resp[0].Percentage, 25.0; got != want {
		t.Fatalf("expected first percentage %v, got %v", want, got)
	}
	if got, want := resp[1].Percentage, 75.0; got != want {
		t.Fatalf("expected second percentage %v, got %v", want, got)
	}
}

func TestReportService_GetMonthlyTrend_DefaultsMonthsAndComputesNetProfit(t *testing.T) {
	repo := &stubReportRepository{
		monthlyTrendResp: []*models.MonthlySummary{
			{Month: "2026-03", Income: 500, Expense: 125},
			{Month: "2026-04", Income: 100, Expense: 150},
		},
	}
	svc := NewReportService(repo)

	fixedNow := time.Date(2026, 4, 18, 10, 30, 0, 0, time.UTC)
	prev := reportNow
	reportNow = func() time.Time { return fixedNow }
	defer func() { reportNow = prev }()

	resp, err := svc.GetMonthlyTrend(context.Background(), 42, 0)
	if err != nil {
		t.Fatalf("GetMonthlyTrend() error = %v", err)
	}

	if got, want := len(resp), 2; got != want {
		t.Fatalf("expected %d monthly items, got %d", want, got)
	}
	if !repo.monthlyTrendSince.Equal(time.Date(2025, 10, 18, 10, 30, 0, 0, time.UTC)) {
		t.Fatalf("expected default since to be 2025-10-18T10:30:00Z, got %s", repo.monthlyTrendSince.Format(time.RFC3339Nano))
	}
	if got, want := resp[0].NetProfit, 375.0; got != want {
		t.Fatalf("expected first net profit %v, got %v", want, got)
	}
	if got, want := resp[1].NetProfit, -50.0; got != want {
		t.Fatalf("expected second net profit %v, got %v", want, got)
	}
}

func TestReportService_WrapsRepositoryErrors(t *testing.T) {
	categoryErr := errors.New("category summary failed")
	monthlyErr := errors.New("monthly trend failed")
	svc := NewReportService(&stubReportRepository{
		categorySummaryErr: categoryErr,
		monthlyTrendErr:    monthlyErr,
	})

	_, err := svc.GetCategorySummary(context.Background(), 42, &models.ReportRequest{
		StartDate: time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2026, 4, 30, 0, 0, 0, 0, time.UTC),
	})
	if err == nil || !errors.Is(err, categoryErr) || !strings.Contains(err.Error(), "failed to fetch category summary") {
		t.Fatalf("expected wrapped category summary error, got %v", err)
	}

	_, err = svc.GetMonthlyTrend(context.Background(), 42, 3)
	if err == nil || !errors.Is(err, monthlyErr) || !strings.Contains(err.Error(), "failed to fetch monthly trend") {
		t.Fatalf("expected wrapped monthly trend error, got %v", err)
	}
}
