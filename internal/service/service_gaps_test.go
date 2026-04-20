package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/maynguyen24/sever/internal/models"
	"github.com/maynguyen24/sever/pkg/apperr"
)

// Additional coverage for category/fund/source_payment GetAll methods
// and transaction CheckBudgets

func TestCategoryService_GetAll(t *testing.T) {
	t.Parallel()

	cats := []*models.Category{{ID: 1, Name: "Food"}, {ID: 2, Name: "Salary"}}
	svc := NewCategoryService(&stubCategoryRepo{
		getAllByUserIDFn: func(ctx context.Context, userID int64) ([]*models.Category, error) {
			return cats, nil
		},
	})
	got, err := svc.GetAll(context.Background(), 42)
	if err != nil {
		t.Fatalf("GetAll returned error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 categories, got %d", len(got))
	}
}

func TestCategoryService_GetAll_Error(t *testing.T) {
	t.Parallel()

	svc := NewCategoryService(&stubCategoryRepo{
		getAllByUserIDFn: func(ctx context.Context, userID int64) ([]*models.Category, error) {
			return nil, errors.New("db error")
		},
	})
	_, err := svc.GetAll(context.Background(), 42)
	if err == nil {
		t.Fatal("expected error from repo")
	}
}

func TestFundService_GetAll(t *testing.T) {
	t.Parallel()

	funds := []*models.Fund{{ID: 1}, {ID: 2}}
	repo := &stubFundServiceAllRepo{
		getAllByUserIDFn: func(ctx context.Context, userID int64) ([]*models.Fund, error) {
			return funds, nil
		},
	}
	svc := NewFundService(repo)
	got, err := svc.GetAll(context.Background(), 42)
	if err != nil {
		t.Fatalf("GetAll returned error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 funds, got %d", len(got))
	}
}

type stubFundServiceAllRepo struct {
	getAllByUserIDFn func(context.Context, int64) ([]*models.Fund, error)
	getByIDFn       func(context.Context, int64, int64) (*models.Fund, error)
}

func (s *stubFundServiceAllRepo) Create(ctx context.Context, fund *models.Fund) error { return nil }
func (s *stubFundServiceAllRepo) GetAllByUserID(ctx context.Context, userID int64) ([]*models.Fund, error) {
	if s.getAllByUserIDFn != nil {
		return s.getAllByUserIDFn(ctx, userID)
	}
	return nil, nil
}
func (s *stubFundServiceAllRepo) GetByID(ctx context.Context, id, userID int64) (*models.Fund, error) {
	if s.getByIDFn != nil {
		return s.getByIDFn(ctx, id, userID)
	}
	return nil, nil
}
func (s *stubFundServiceAllRepo) Update(ctx context.Context, fund *models.Fund) error { return nil }
func (s *stubFundServiceAllRepo) Delete(ctx context.Context, id, userID int64) error  { return nil }
func (s *stubFundServiceAllRepo) Deposit(ctx context.Context, id, userID int64, amount float64) (*models.Fund, error) {
	return &models.Fund{ID: id}, nil
}
func (s *stubFundServiceAllRepo) Withdraw(ctx context.Context, id, userID int64, amount float64) (*models.Fund, error) {
	return &models.Fund{ID: id}, nil
}

func TestFundService_Update_Full(t *testing.T) {
	t.Parallel()

	existing := &models.Fund{ID: 1, UserID: 42, Name: "Old", Currency: "VND"}
	repo := &stubFundServiceAllRepo{
		getByIDFn: func(ctx context.Context, id, userID int64) (*models.Fund, error) {
			return existing, nil
		},
	}

	svc := NewFundService(repo)
	got, err := svc.Update(context.Background(), 1, 42, &models.UpdateFundRequest{
		Name:         "New Name",
		TargetAmount: 1000,
		Currency:     "USD",
	})
	if err != nil {
		t.Fatalf("Update returned error: %v", err)
	}
	if got.Name != "New Name" || got.Currency != "USD" {
		t.Fatalf("unexpected fund returned: %+v", got)
	}
}

func TestFundService_Update_InvalidTargetAmount(t *testing.T) {
	t.Parallel()

	svc := NewFundService(&stubFundServiceAllRepo{})
	_, err := svc.Update(context.Background(), 1, 42, &models.UpdateFundRequest{
		Name:         "Test",
		TargetAmount: -10,
	})
	if !errors.Is(err, apperr.ErrInvalidInput) {
		t.Fatalf("expected invalid input, got %v", err)
	}
}

func TestSourcePaymentService_GetAll(t *testing.T) {
	t.Parallel()

	sources := []*models.SourcePayment{{ID: 1}, {ID: 2}}
	svc := NewSourcePaymentService(&stubSourcePaymentServiceRepo{
		getAllByUserIDFn: func(ctx context.Context, userID int64) ([]*models.SourcePayment, error) {
			return sources, nil
		},
	})
	got, err := svc.GetAll(context.Background(), 42)
	if err != nil {
		t.Fatalf("GetAll returned error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 sources, got %d", len(got))
	}
}

type stubSourcePaymentServiceRepo struct {
	createFn        func(context.Context, *models.SourcePayment) error
	getAllByUserIDFn func(context.Context, int64) ([]*models.SourcePayment, error)
	getByIDFn       func(context.Context, int64, int64) (*models.SourcePayment, error)
	updateFn        func(context.Context, *models.SourcePayment) error
	deleteFn        func(context.Context, int64, int64) error
}

func (s *stubSourcePaymentServiceRepo) Create(ctx context.Context, source *models.SourcePayment) error {
	if s.createFn != nil {
		return s.createFn(ctx, source)
	}
	return nil
}

func (s *stubSourcePaymentServiceRepo) GetAllByUserID(ctx context.Context, userID int64) ([]*models.SourcePayment, error) {
	if s.getAllByUserIDFn != nil {
		return s.getAllByUserIDFn(ctx, userID)
	}
	return nil, nil
}

func (s *stubSourcePaymentServiceRepo) GetByID(ctx context.Context, id, userID int64) (*models.SourcePayment, error) {
	if s.getByIDFn != nil {
		return s.getByIDFn(ctx, id, userID)
	}
	return nil, nil
}

func (s *stubSourcePaymentServiceRepo) Update(ctx context.Context, source *models.SourcePayment) error {
	if s.updateFn != nil {
		return s.updateFn(ctx, source)
	}
	return nil
}

func (s *stubSourcePaymentServiceRepo) Delete(ctx context.Context, id, userID int64) error {
	if s.deleteFn != nil {
		return s.deleteFn(ctx, id, userID)
	}
	return nil
}

func TestTransactionService_CheckBudgets(t *testing.T) {
	t.Parallel()

	catID := int64(5)
	budgets := []*models.Budget{
		{ID: 1, UserID: 42, CategoryID: &catID, Amount: 1000, IsActive: true},
		{ID: 2, UserID: 42, CategoryID: nil, Amount: 500, IsActive: true}, // global budget
		{ID: 3, UserID: 42, Amount: 200, IsActive: false},                 // inactive
	}

	var notifications []*models.Notification
	svc := NewTransactionService(
		&stubTransactionRepo{},
		&stubBudgetRepo{
			getByUserIDFn: func(ctx context.Context, userID int64) ([]*models.Budget, error) {
				return budgets, nil
			},
			calculateSpendingFn: func(ctx context.Context, userID int64, categoryID *int64, start, end time.Time) (float64, error) {
				return 900, nil // 90% of budget -> warning
			},
		},
		&stubNotificationRepo{
			createFn: func(ctx context.Context, notif *models.Notification) error {
				notifications = append(notifications, notif)
				return nil
			},
		},
		nil,
	)

	err := svc.CheckBudgets(context.Background(), 42, &catID)
	if err != nil {
		t.Fatalf("CheckBudgets returned error: %v", err)
	}
	if len(notifications) == 0 {
		t.Fatal("expected at least one notification")
	}
}

func TestTransactionService_CheckBudgets_Exceeded(t *testing.T) {
	t.Parallel()

	budgets := []*models.Budget{
		{ID: 1, UserID: 42, Amount: 100, IsActive: true}, // global budget
	}

	var notifications []*models.Notification
	svc := NewTransactionService(
		&stubTransactionRepo{},
		&stubBudgetRepo{
			getByUserIDFn: func(ctx context.Context, userID int64) ([]*models.Budget, error) {
				return budgets, nil
			},
			calculateSpendingFn: func(ctx context.Context, userID int64, categoryID *int64, start, end time.Time) (float64, error) {
				return 110, nil // 110% -> exceeded
			},
		},
		&stubNotificationRepo{
			createFn: func(ctx context.Context, notif *models.Notification) error {
				notifications = append(notifications, notif)
				return nil
			},
		},
		nil,
	)

	err := svc.CheckBudgets(context.Background(), 42, nil)
	if err != nil {
		t.Fatalf("CheckBudgets returned error: %v", err)
	}
	if len(notifications) == 0 {
		t.Fatal("expected budget exceeded notification")
	}
}

func TestTransactionService_CheckBudgets_NoBudgets(t *testing.T) {
	t.Parallel()

	svc := NewTransactionService(
		&stubTransactionRepo{},
		&stubBudgetRepo{
			getByUserIDFn: func(ctx context.Context, userID int64) ([]*models.Budget, error) {
				return nil, nil
			},
		},
		&stubNotificationRepo{},
		nil,
	)

	if err := svc.CheckBudgets(context.Background(), 42, nil); err != nil {
		t.Fatalf("CheckBudgets returned error: %v", err)
	}
}

func TestTransactionService_CheckBudgets_RepoError(t *testing.T) {
	t.Parallel()

	svc := NewTransactionService(
		&stubTransactionRepo{},
		&stubBudgetRepo{
			getByUserIDFn: func(ctx context.Context, userID int64) ([]*models.Budget, error) {
				return nil, errors.New("db error")
			},
		},
		&stubNotificationRepo{},
		nil,
	)

	err := svc.CheckBudgets(context.Background(), 42, nil)
	if err == nil {
		t.Fatal("expected error from repo")
	}
}
