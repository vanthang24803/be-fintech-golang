package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/maynguyen24/sever/internal/models"
	"github.com/maynguyen24/sever/pkg/apperr"
)

type stubBudgetServiceRepo struct {
	createFn            func(context.Context, *models.Budget) error
	getByUserIDFn       func(context.Context, int64) ([]*models.Budget, error)
	getByIDFn           func(context.Context, int64, int64) (*models.Budget, error)
	updateFn            func(context.Context, int64, int64, *models.UpdateBudgetRequest) (*models.Budget, error)
	deleteFn            func(context.Context, int64, int64) error
	calculateSpendingFn func(context.Context, int64, *int64, time.Time, time.Time) (float64, error)
}

func (s *stubBudgetServiceRepo) Create(ctx context.Context, budget *models.Budget) error {
	if s.createFn != nil {
		return s.createFn(ctx, budget)
	}
	return nil
}

func (s *stubBudgetServiceRepo) GetByUserID(ctx context.Context, userID int64) ([]*models.Budget, error) {
	if s.getByUserIDFn != nil {
		return s.getByUserIDFn(ctx, userID)
	}
	return nil, nil
}

func (s *stubBudgetServiceRepo) GetByID(ctx context.Context, id, userID int64) (*models.Budget, error) {
	if s.getByIDFn != nil {
		return s.getByIDFn(ctx, id, userID)
	}
	return nil, nil
}

func (s *stubBudgetServiceRepo) Update(ctx context.Context, id, userID int64, req *models.UpdateBudgetRequest) (*models.Budget, error) {
	if s.updateFn != nil {
		return s.updateFn(ctx, id, userID, req)
	}
	return nil, nil
}

func (s *stubBudgetServiceRepo) Delete(ctx context.Context, id, userID int64) error {
	if s.deleteFn != nil {
		return s.deleteFn(ctx, id, userID)
	}
	return nil
}

func (s *stubBudgetServiceRepo) CalculateSpending(ctx context.Context, userID int64, categoryID *int64, start, end time.Time) (float64, error) {
	if s.calculateSpendingFn != nil {
		return s.calculateSpendingFn(ctx, userID, categoryID, start, end)
	}
	return 0, nil
}

func TestBudgetService_Create_Monthly(t *testing.T) {
	t.Parallel()

	var created *models.Budget
	repo := &stubBudgetServiceRepo{
		createFn: func(ctx context.Context, b *models.Budget) error {
			created = b
			return nil
		},
	}

	catID := int64(5)
	svc := NewBudgetService(repo)
	got, err := svc.Create(context.Background(), 42, &models.CreateBudgetRequest{
		CategoryID: &catID,
		Amount:     1000,
		Period:     "monthly",
	})
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if got == nil || created == nil {
		t.Fatal("expected budget to be created")
	}
	if created.Amount != 1000 || created.Period != "monthly" || !created.IsActive {
		t.Fatalf("unexpected budget created: %+v", created)
	}
	if created.StartDate.Day() != 1 {
		t.Fatalf("expected start date to be first of month, got day %d", created.StartDate.Day())
	}
	if created.EndDate.Hour() != 23 || created.EndDate.Minute() != 59 {
		t.Fatalf("expected end date at 23:59, got %v", created.EndDate)
	}
}

func TestBudgetService_Create_Weekly(t *testing.T) {
	t.Parallel()

	var created *models.Budget
	repo := &stubBudgetServiceRepo{
		createFn: func(ctx context.Context, b *models.Budget) error {
			created = b
			return nil
		},
	}

	svc := NewBudgetService(repo)
	got, err := svc.Create(context.Background(), 42, &models.CreateBudgetRequest{
		Amount: 500,
		Period: "weekly",
	})
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if got == nil || created == nil {
		t.Fatal("expected budget to be created")
	}
	if created.Period != "weekly" {
		t.Fatalf("expected weekly period, got %s", created.Period)
	}
	diff := created.EndDate.Sub(created.StartDate)
	if diff < 5*24*time.Hour || diff > 7*24*time.Hour {
		t.Fatalf("unexpected weekly date range: %v", diff)
	}
}

func TestBudgetService_Create_ValidationErrors(t *testing.T) {
	t.Parallel()

	svc := NewBudgetService(&stubBudgetServiceRepo{})

	cases := []struct {
		name string
		req  *models.CreateBudgetRequest
	}{
		{"zero amount", &models.CreateBudgetRequest{Amount: 0, Period: "monthly"}},
		{"negative amount", &models.CreateBudgetRequest{Amount: -1, Period: "monthly"}},
		{"invalid period", &models.CreateBudgetRequest{Amount: 100, Period: "yearly"}},
	}
	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := svc.Create(context.Background(), 42, tt.req)
			if !errors.Is(err, apperr.ErrInvalidInput) {
				t.Fatalf("expected invalid input, got %v", err)
			}
		})
	}
}

func TestBudgetService_GetList(t *testing.T) {
	t.Parallel()

	catID := int64(1)
	budgets := []*models.Budget{
		{ID: 1, Amount: 1000, CategoryID: &catID, IsActive: true},
		{ID: 2, Amount: 500, IsActive: true},
	}

	repo := &stubBudgetServiceRepo{
		getByUserIDFn: func(ctx context.Context, userID int64) ([]*models.Budget, error) {
			return budgets, nil
		},
		calculateSpendingFn: func(ctx context.Context, userID int64, categoryID *int64, start, end time.Time) (float64, error) {
			return 200, nil
		},
	}

	svc := NewBudgetService(repo)
	resp, err := svc.GetList(context.Background(), 42)
	if err != nil {
		t.Fatalf("GetList returned error: %v", err)
	}
	if len(resp) != 2 {
		t.Fatalf("expected 2 budgets, got %d", len(resp))
	}
	if resp[0].ProgressPercent != 20 {
		t.Fatalf("expected progress 20, got %v", resp[0].ProgressPercent)
	}
	if resp[0].RemainingAmount != 800 {
		t.Fatalf("expected remaining 800, got %v", resp[0].RemainingAmount)
	}
}

func TestBudgetService_GetList_SpendingOverBudget(t *testing.T) {
	t.Parallel()

	budgets := []*models.Budget{{ID: 1, Amount: 100, IsActive: true}}
	repo := &stubBudgetServiceRepo{
		getByUserIDFn: func(ctx context.Context, userID int64) ([]*models.Budget, error) {
			return budgets, nil
		},
		calculateSpendingFn: func(ctx context.Context, userID int64, categoryID *int64, start, end time.Time) (float64, error) {
			return 150, nil
		},
	}

	svc := NewBudgetService(repo)
	resp, err := svc.GetList(context.Background(), 42)
	if err != nil {
		t.Fatalf("GetList returned error: %v", err)
	}
	if resp[0].RemainingAmount != 0 {
		t.Fatalf("expected remaining 0 when over budget, got %v", resp[0].RemainingAmount)
	}
	if resp[0].ProgressPercent != 100 {
		t.Fatalf("expected progress 100 when over budget, got %v", resp[0].ProgressPercent)
	}
}

func TestBudgetService_GetList_RepoError(t *testing.T) {
	t.Parallel()

	repo := &stubBudgetServiceRepo{
		getByUserIDFn: func(ctx context.Context, userID int64) ([]*models.Budget, error) {
			return nil, errors.New("db error")
		},
	}

	svc := NewBudgetService(repo)
	_, err := svc.GetList(context.Background(), 42)
	if err == nil {
		t.Fatal("expected error from repo")
	}
}

func TestBudgetService_GetDetail(t *testing.T) {
	t.Parallel()

	b := &models.Budget{ID: 1, UserID: 42, Amount: 500, IsActive: true}
	repo := &stubBudgetServiceRepo{
		getByIDFn: func(ctx context.Context, id, userID int64) (*models.Budget, error) {
			return b, nil
		},
		calculateSpendingFn: func(ctx context.Context, userID int64, categoryID *int64, start, end time.Time) (float64, error) {
			return 100, nil
		},
	}

	svc := NewBudgetService(repo)
	resp, err := svc.GetDetail(context.Background(), 1, 42)
	if err != nil {
		t.Fatalf("GetDetail returned error: %v", err)
	}
	if resp.Budget.ID != 1 {
		t.Fatalf("unexpected budget ID: %d", resp.Budget.ID)
	}
	if resp.CurrentSpending != 100 {
		t.Fatalf("expected spending 100, got %v", resp.CurrentSpending)
	}
}

func TestBudgetService_GetDetail_NotFound(t *testing.T) {
	t.Parallel()

	repo := &stubBudgetServiceRepo{
		getByIDFn: func(ctx context.Context, id, userID int64) (*models.Budget, error) {
			return nil, nil
		},
	}

	svc := NewBudgetService(repo)
	_, err := svc.GetDetail(context.Background(), 99, 42)
	if !errors.Is(err, apperr.ErrNotFound) {
		t.Fatalf("expected not found, got %v", err)
	}
}

func TestBudgetService_Update(t *testing.T) {
	t.Parallel()

	amount := float64(800)
	isActive := true
	b := &models.Budget{ID: 1, Amount: 800, IsActive: true}
	repo := &stubBudgetServiceRepo{
		updateFn: func(ctx context.Context, id, userID int64, req *models.UpdateBudgetRequest) (*models.Budget, error) {
			return b, nil
		},
	}

	svc := NewBudgetService(repo)
	got, err := svc.Update(context.Background(), 1, 42, &models.UpdateBudgetRequest{
		Amount: &amount, IsActive: &isActive,
	})
	if err != nil {
		t.Fatalf("Update returned error: %v", err)
	}
	if got.ID != 1 {
		t.Fatalf("unexpected budget returned: %+v", got)
	}
}

func TestBudgetService_Update_NotFound(t *testing.T) {
	t.Parallel()

	repo := &stubBudgetServiceRepo{
		updateFn: func(ctx context.Context, id, userID int64, req *models.UpdateBudgetRequest) (*models.Budget, error) {
			return nil, nil
		},
	}

	svc := NewBudgetService(repo)
	_, err := svc.Update(context.Background(), 1, 42, &models.UpdateBudgetRequest{})
	if !errors.Is(err, apperr.ErrNotFound) {
		t.Fatalf("expected not found, got %v", err)
	}
}

func TestBudgetService_Delete(t *testing.T) {
	t.Parallel()

	deleted := false
	repo := &stubBudgetServiceRepo{
		deleteFn: func(ctx context.Context, id, userID int64) error {
			deleted = true
			return nil
		},
	}

	svc := NewBudgetService(repo)
	if err := svc.Delete(context.Background(), 1, 42); err != nil {
		t.Fatalf("Delete returned error: %v", err)
	}
	if !deleted {
		t.Fatal("expected delete to be called")
	}
}

func TestBudgetService_Delete_Error(t *testing.T) {
	t.Parallel()

	repo := &stubBudgetServiceRepo{
		deleteFn: func(ctx context.Context, id, userID int64) error {
			return apperr.ErrNotFound
		},
	}

	svc := NewBudgetService(repo)
	err := svc.Delete(context.Background(), 1, 42)
	if !errors.Is(err, apperr.ErrNotFound) {
		t.Fatalf("expected not found, got %v", err)
	}
}
