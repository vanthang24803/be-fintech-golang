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

	dailyTrendSince time.Time
	dailyTrendErr   error
	dailyTrendResp  []*models.DailySummary

	breakdownStart time.Time
	breakdownEnd   time.Time
	breakdownResp  []*models.IncomeCategoryBreakdownItem

	trendStart    time.Time
	trendEnd      time.Time
	trendCategory int64
	trendGran     string
	trendResp     []*models.CategoryTrendPoint
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

func (s *stubReportRepository) GetDailyTrend(ctx context.Context, userID int64, since time.Time) ([]*models.DailySummary, error) {
	s.dailyTrendSince = since
	return s.dailyTrendResp, s.dailyTrendErr
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

func TestReportService_GetDailyTrend_DefaultsDays(t *testing.T) {
	repo := &stubReportRepository{
		dailyTrendResp: []*models.DailySummary{
			{Date: "2026-04-17", Income: 500, Expense: 150},
		},
	}
	svc := NewReportService(repo)

	fixedNow := time.Date(2026, 4, 18, 10, 30, 0, 0, time.UTC)
	prev := reportNow
	reportNow = func() time.Time { return fixedNow }
	defer func() { reportNow = prev }()

	resp, err := svc.GetDailyTrend(context.Background(), 42, 0)
	if err != nil {
		t.Fatalf("GetDailyTrend() error = %v", err)
	}

	if got, want := len(resp), 1; got != want {
		t.Fatalf("expected %d daily items, got %d", want, got)
	}
	if !repo.dailyTrendSince.Equal(time.Date(2026, 4, 11, 10, 30, 0, 0, time.UTC)) {
		t.Fatalf("expected default since to be 2026-04-11T10:30:00Z, got %s", repo.dailyTrendSince.Format(time.RFC3339Nano))
	}
}

func TestReportService_WrapsRepositoryErrors(t *testing.T) {
	categoryErr := errors.New("category summary failed")
	monthlyErr := errors.New("monthly trend failed")
	dailyErr := errors.New("daily trend failed")
	svc := NewReportService(&stubReportRepository{
		categorySummaryErr: categoryErr,
		monthlyTrendErr:    monthlyErr,
		dailyTrendErr:      dailyErr,
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

	_, err = svc.GetDailyTrend(context.Background(), 42, 7)
	if err == nil || !errors.Is(err, dailyErr) || !strings.Contains(err.Error(), "failed to fetch daily trend") {
		t.Fatalf("expected wrapped daily trend error, got %v", err)
	}
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
		CategoryID:  12,
		StartDate:   time.Date(2026, 4, 30, 0, 0, 0, 0, time.UTC),
		EndDate:     time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC),
		Granularity: "day",
	})
	if err == nil {
		t.Fatal("expected error for invalid trend range")
	}
}
