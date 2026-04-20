package service

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/maynguyen24/sever/internal/models"
	"github.com/maynguyen24/sever/pkg/apperr"
)

func TestReportService_GetIncomeCategoryBreakdown_WrapsRepositoryError(t *testing.T) {
	t.Parallel()

	repoErr := errors.New("repo failed")
	svc := NewReportService(&stubReportRepositoryWithBreakdownError{err: repoErr})

	_, err := svc.GetIncomeCategoryBreakdown(context.Background(), 99, &models.IncomeCategoryBreakdownRequest{
		StartDate: time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2026, 4, 30, 0, 0, 0, 0, time.UTC),
	})
	if err == nil || !errors.Is(err, repoErr) || !strings.Contains(err.Error(), "failed to fetch income category breakdown") {
		t.Fatalf("expected wrapped breakdown error, got %v", err)
	}
}

type stubReportRepositoryWithBreakdownError struct {
	stubReportRepository
	err error
}

func (s *stubReportRepositoryWithBreakdownError) GetIncomeCategoryBreakdown(ctx context.Context, userID int64, start, end time.Time, limit int) ([]*models.IncomeCategoryBreakdownItem, error) {
	return nil, s.err
}

func TestReportService_GetCategoryTrend_RejectsMissingCategoryID(t *testing.T) {
	t.Parallel()

	svc := NewReportService(&stubReportRepository{})
	_, err := svc.GetCategoryTrend(context.Background(), 99, &models.CategoryTrendRequest{})
	if !errors.Is(err, apperr.ErrInvalidInput) {
		t.Fatalf("expected invalid input, got %v", err)
	}
}

func TestReportService_GetCategoryTrend_RejectsUnsupportedGranularity(t *testing.T) {
	t.Parallel()

	svc := NewReportService(&stubReportRepository{})
	_, err := svc.GetCategoryTrend(context.Background(), 99, &models.CategoryTrendRequest{
		CategoryID:  12,
		StartDate:   time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC),
		EndDate:     time.Date(2026, 4, 30, 0, 0, 0, 0, time.UTC),
		Granularity: "week",
	})
	if !errors.Is(err, apperr.ErrInvalidInput) {
		t.Fatalf("expected invalid input, got %v", err)
	}
}

func TestReportService_GetCategoryTrend_WrapsRepositoryError(t *testing.T) {
	t.Parallel()

	repoErr := errors.New("trend failed")
	svc := NewReportService(&stubReportRepositoryWithTrendError{err: repoErr})

	_, err := svc.GetCategoryTrend(context.Background(), 99, &models.CategoryTrendRequest{
		CategoryID:  12,
		StartDate:   time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC),
		EndDate:     time.Date(2026, 4, 30, 0, 0, 0, 0, time.UTC),
		Granularity: "day",
	})
	if err == nil || !errors.Is(err, repoErr) || !strings.Contains(err.Error(), "failed to fetch category trend") {
		t.Fatalf("expected wrapped trend error, got %v", err)
	}
}

type stubReportRepositoryWithTrendError struct {
	stubReportRepository
	err error
}

func (s *stubReportRepositoryWithTrendError) GetCategoryTrend(ctx context.Context, userID, categoryID int64, start, end time.Time, granularity string) ([]*models.CategoryTrendPoint, error) {
	return nil, s.err
}

func TestSavingsGoalService_GetDetail_ContributionError(t *testing.T) {
	t.Parallel()

	repoErr := errors.New("contribution lookup failed")
	repo := &stubSavingsGoalRepo{
		getGoalByIDFn: func(context.Context, int64) (*models.SavingsGoal, error) {
			return &models.SavingsGoal{ID: 1, UserID: 42, TargetAmount: 1000, CurrentAmount: 200}, nil
		},
		getContributionsByGoalFn: func(context.Context, int64) ([]models.GoalContribution, error) {
			return nil, repoErr
		},
	}

	svc := NewSavingsGoalService(repo, &stubSavingsFundRepo{}, &stubNotifier{}, nil)
	_, err := svc.GetDetail(context.Background(), 1, 42)
	if err == nil || !errors.Is(err, repoErr) || !strings.Contains(err.Error(), "failed to fetch contributions") {
		t.Fatalf("expected wrapped contributions error, got %v", err)
	}
}

func TestSavingsGoalService_Withdraw_DepositError(t *testing.T) {
	t.Parallel()

	depositErr := errors.New("deposit failed")
	repo := &stubSavingsGoalRepo{
		getGoalByIDFn: func(context.Context, int64) (*models.SavingsGoal, error) {
			return &models.SavingsGoal{ID: 1, UserID: 42, CurrentAmount: 500}, nil
		},
	}
	fundRepo := &stubSavingsFundRepo{
		depositFn: func(context.Context, int64, int64, float64) (*models.Fund, error) {
			return nil, depositErr
		},
	}

	svc := NewSavingsGoalService(repo, fundRepo, &stubNotifier{}, nil)
	_, err := svc.Withdraw(context.Background(), 42, &models.GoalWithdrawRequest{
		GoalID: 1,
		FundID: 2,
		Amount: 200,
	})
	if !errors.Is(err, depositErr) {
		t.Fatalf("expected deposit error, got %v", err)
	}
}

func TestSavingsGoalService_Withdraw_UpdateGoalAmountError(t *testing.T) {
	t.Parallel()

	updateErr := errors.New("update failed")
	repo := &stubSavingsGoalRepo{
		getGoalByIDFn: func(context.Context, int64) (*models.SavingsGoal, error) {
			return &models.SavingsGoal{ID: 1, UserID: 42, CurrentAmount: 500}, nil
		},
		updateGoalAmountFn: func(context.Context, int64, float64) error {
			return updateErr
		},
	}

	svc := NewSavingsGoalService(repo, &stubSavingsFundRepo{}, &stubNotifier{}, nil)
	_, err := svc.Withdraw(context.Background(), 42, &models.GoalWithdrawRequest{
		GoalID: 1,
		FundID: 2,
		Amount: 200,
	})
	if err == nil || !errors.Is(err, updateErr) || !strings.Contains(err.Error(), "failed to update goal balance") {
		t.Fatalf("expected wrapped update error, got %v", err)
	}
}

func TestSavingsGoalService_Withdraw_CreateContributionError(t *testing.T) {
	t.Parallel()

	contributionErr := errors.New("create contribution failed")
	repo := &stubSavingsGoalRepo{
		getGoalByIDFn: func(ctx context.Context, id int64) (*models.SavingsGoal, error) {
			if id != 1 {
				t.Fatalf("unexpected goal lookup id %d", id)
			}
			return &models.SavingsGoal{ID: 1, UserID: 42, CurrentAmount: 500}, nil
		},
		updateGoalAmountFn: func(context.Context, int64, float64) error {
			return nil
		},
		createContributionFn: func(context.Context, *models.GoalContribution) error {
			return contributionErr
		},
	}

	svc := NewSavingsGoalService(repo, &stubSavingsFundRepo{}, &stubNotifier{}, nil)
	_, err := svc.Withdraw(context.Background(), 42, &models.GoalWithdrawRequest{
		GoalID: 1,
		FundID: 2,
		Amount: 200,
	})
	if err == nil || !errors.Is(err, contributionErr) || !strings.Contains(err.Error(), "failed to log withdrawal") {
		t.Fatalf("expected wrapped contribution error, got %v", err)
	}
}
