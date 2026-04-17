package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/maynguyen24/sever/internal/models"
	"github.com/maynguyen24/sever/pkg/apperr"
	"github.com/maynguyen24/sever/pkg/queue"
)

type stubTransactionRepo struct {
	createFn  func(context.Context, *models.Transaction) error
	getAllFn  func(context.Context, int64, models.TransactionFilter) ([]*models.TransactionDetail, error)
	getByIDFn func(context.Context, int64, int64) (*models.TransactionDetail, error)
	getRawFn  func(context.Context, int64, int64) (*models.Transaction, error)
	updateFn  func(context.Context, *models.Transaction, *models.Transaction) error
	deleteFn  func(context.Context, int64, int64) error
}

func (s *stubTransactionRepo) Create(ctx context.Context, tx *models.Transaction) error {
	if s.createFn != nil {
		return s.createFn(ctx, tx)
	}
	return nil
}

func (s *stubTransactionRepo) GetAllByUserID(ctx context.Context, userID int64, filter models.TransactionFilter) ([]*models.TransactionDetail, error) {
	if s.getAllFn != nil {
		return s.getAllFn(ctx, userID, filter)
	}
	return nil, nil
}

func (s *stubTransactionRepo) GetByID(ctx context.Context, id, userID int64) (*models.TransactionDetail, error) {
	if s.getByIDFn != nil {
		return s.getByIDFn(ctx, id, userID)
	}
	return nil, nil
}

func (s *stubTransactionRepo) GetRawByID(ctx context.Context, id, userID int64) (*models.Transaction, error) {
	if s.getRawFn != nil {
		return s.getRawFn(ctx, id, userID)
	}
	return nil, nil
}

func (s *stubTransactionRepo) Update(ctx context.Context, old *models.Transaction, updated *models.Transaction) error {
	if s.updateFn != nil {
		return s.updateFn(ctx, old, updated)
	}
	return nil
}

func (s *stubTransactionRepo) Delete(ctx context.Context, id, userID int64) error {
	if s.deleteFn != nil {
		return s.deleteFn(ctx, id, userID)
	}
	return nil
}

type stubBudgetRepo struct {
	getByUserIDFn       func(context.Context, int64) ([]*models.Budget, error)
	calculateSpendingFn func(context.Context, int64, *int64, time.Time, time.Time) (float64, error)
}

func (s *stubBudgetRepo) GetByUserID(ctx context.Context, userID int64) ([]*models.Budget, error) {
	if s.getByUserIDFn != nil {
		return s.getByUserIDFn(ctx, userID)
	}
	return nil, nil
}

func (s *stubBudgetRepo) CalculateSpending(ctx context.Context, userID int64, categoryID *int64, start, end time.Time) (float64, error) {
	if s.calculateSpendingFn != nil {
		return s.calculateSpendingFn(ctx, userID, categoryID, start, end)
	}
	return 0, nil
}

type stubNotificationRepo struct {
	createFn func(context.Context, *models.Notification) error
}

func (s *stubNotificationRepo) Create(ctx context.Context, notif *models.Notification) error {
	if s.createFn != nil {
		return s.createFn(ctx, notif)
	}
	return nil
}

func TestTransactionService_Create(t *testing.T) {
	t.Parallel()

	desc := "coffee"
	now := time.Date(2026, 4, 18, 10, 0, 0, 0, time.UTC)
	userID := int64(99)
	categoryID := int64(7)
	repo := &stubTransactionRepo{
		createFn: func(ctx context.Context, tx *models.Transaction) error {
			if tx.UserID != userID || tx.Amount != 125.5 || tx.Type != models.TransactionTypeExpense || tx.Description == nil || *tx.Description != desc {
				t.Fatalf("unexpected transaction passed to repo: %+v", tx)
			}
			if tx.CategoryID == nil || *tx.CategoryID != categoryID {
				t.Fatalf("expected category id to be set: %+v", tx)
			}
			return nil
		},
	}
	svc := NewTransactionService(repo, &stubBudgetRepo{}, &stubNotificationRepo{}, nil)

	got, err := svc.Create(context.Background(), userID, &models.CreateTransactionRequest{
		SourcePaymentID: 3,
		CategoryID:      &categoryID,
		Amount:          125.5,
		Type:            models.TransactionTypeExpense,
		Description:     &desc,
		TransactionDate: now,
	})
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if got == nil || got.UserID != userID || got.Amount != 125.5 || got.Type != models.TransactionTypeExpense {
		t.Fatalf("unexpected transaction returned: %+v", got)
	}

	invalidCases := []struct {
		name string
		req  *models.CreateTransactionRequest
	}{
		{name: "zero amount", req: &models.CreateTransactionRequest{Amount: 0, Type: models.TransactionTypeIncome}},
		{name: "invalid type", req: &models.CreateTransactionRequest{Amount: 10, Type: "invalid"}},
	}
	for _, tt := range invalidCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := svc.Create(context.Background(), userID, tt.req)
			if !errors.Is(err, apperr.ErrInvalidInput) {
				t.Fatalf("expected invalid input, got %v", err)
			}
		})
	}
}

func TestTransactionService_GetAll(t *testing.T) {
	t.Parallel()

	expected := []*models.TransactionDetail{{Transaction: models.Transaction{ID: 1}}}
	var gotFilter models.TransactionFilter
	svc := NewTransactionService(&stubTransactionRepo{
		getAllFn: func(ctx context.Context, userID int64, filter models.TransactionFilter) ([]*models.TransactionDetail, error) {
			gotFilter = filter
			return expected, nil
		},
	}, &stubBudgetRepo{}, &stubNotificationRepo{}, nil)

	resp, err := svc.GetAll(context.Background(), 12, map[string]string{
		"type":        models.TransactionTypeExpense,
		"category_id": "7",
		"source_id":   "3",
	})
	if err != nil {
		t.Fatalf("GetAll returned error: %v", err)
	}
	if len(resp) != 1 || resp[0].ID != 1 {
		t.Fatalf("unexpected response: %+v", resp)
	}
	if gotFilter.Type != models.TransactionTypeExpense || gotFilter.CategoryID != 7 || gotFilter.SourcePaymentID != 3 {
		t.Fatalf("unexpected filter passed to repo: %+v", gotFilter)
	}
}

func TestTransactionService_GetByID_NotFound(t *testing.T) {
	t.Parallel()

	svc := NewTransactionService(&stubTransactionRepo{
		getByIDFn: func(context.Context, int64, int64) (*models.TransactionDetail, error) {
			return nil, nil
		},
	}, &stubBudgetRepo{}, &stubNotificationRepo{}, nil)

	_, err := svc.GetByID(context.Background(), 1, 2)
	if !errors.Is(err, apperr.ErrNotFound) {
		t.Fatalf("expected not found, got %v", err)
	}
}

func TestTransactionService_Update(t *testing.T) {
	t.Parallel()

	desc := "lunch"
	old := &models.Transaction{ID: 1, UserID: 9, Amount: 50, Type: models.TransactionTypeIncome}
	categoryID := int64(44)
	svc := NewTransactionService(&stubTransactionRepo{
		getRawFn: func(context.Context, int64, int64) (*models.Transaction, error) {
			return old, nil
		},
		updateFn: func(ctx context.Context, gotOld *models.Transaction, updated *models.Transaction) error {
			if gotOld != old {
				t.Fatalf("expected old tx to be passed through")
			}
			if updated.Amount != 88 || updated.Type != models.TransactionTypeExpense || updated.CategoryID == nil || *updated.CategoryID != categoryID || updated.Description == nil || *updated.Description != desc {
				t.Fatalf("unexpected updated tx: %+v", updated)
			}
			return nil
		},
	}, &stubBudgetRepo{}, &stubNotificationRepo{}, nil)

	got, err := svc.Update(context.Background(), 1, 9, &models.UpdateTransactionRequest{
		SourcePaymentID: 11,
		CategoryID:      &categoryID,
		Amount:          88,
		Type:            models.TransactionTypeExpense,
		Description:     &desc,
		TransactionDate: time.Date(2026, 4, 18, 10, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("Update returned error: %v", err)
	}
	if got == nil || got.Amount != 88 || got.Type != models.TransactionTypeExpense {
		t.Fatalf("unexpected update response: %+v", got)
	}

	invalidCases := []struct {
		name string
		req  *models.UpdateTransactionRequest
	}{
		{name: "zero amount", req: &models.UpdateTransactionRequest{Amount: 0, Type: models.TransactionTypeIncome}},
		{name: "invalid type", req: &models.UpdateTransactionRequest{Amount: 5, Type: "invalid"}},
	}
	for _, tt := range invalidCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := svc.Update(context.Background(), 1, 9, tt.req)
			if !errors.Is(err, apperr.ErrInvalidInput) {
				t.Fatalf("expected invalid input, got %v", err)
			}
		})
	}

	svc = NewTransactionService(&stubTransactionRepo{
		getRawFn: func(context.Context, int64, int64) (*models.Transaction, error) {
			return nil, nil
		},
	}, &stubBudgetRepo{}, &stubNotificationRepo{}, nil)
	_, err = svc.Update(context.Background(), 1, 9, &models.UpdateTransactionRequest{
		Amount:          5,
		Type:            models.TransactionTypeIncome,
		TransactionDate: time.Date(2026, 4, 18, 10, 0, 0, 0, time.UTC),
	})
	if !errors.Is(err, apperr.ErrNotFound) {
		t.Fatalf("expected not found, got %v", err)
	}
}

func TestTransactionService_Delete_NotFound(t *testing.T) {
	t.Parallel()

	svc := NewTransactionService(&stubTransactionRepo{
		deleteFn: func(context.Context, int64, int64) error {
			return apperr.ErrNotFound
		},
	}, &stubBudgetRepo{}, &stubNotificationRepo{}, nil)

	err := svc.Delete(context.Background(), 1, 2)
	if !errors.Is(err, apperr.ErrNotFound) {
		t.Fatalf("expected not found, got %v", err)
	}
}

func TestTransactionService_Delete_OK(t *testing.T) {
	t.Parallel()

	deleted := false
	svc := NewTransactionService(&stubTransactionRepo{
		deleteFn: func(context.Context, int64, int64) error {
			deleted = true
			return nil
		},
	}, &stubBudgetRepo{}, &stubNotificationRepo{}, nil)

	if err := svc.Delete(context.Background(), 1, 2); err != nil {
		t.Fatalf("Delete returned error: %v", err)
	}
	if !deleted {
		t.Fatal("expected delete to be called")
	}
}

func TestTransactionService_CreateExpenseWithNilQueue(t *testing.T) {
	t.Parallel()

	svc := NewTransactionService(&stubTransactionRepo{
		createFn: func(context.Context, *models.Transaction) error {
			return nil
		},
	}, &stubBudgetRepo{}, &stubNotificationRepo{}, nil)

	_, err := svc.Create(context.Background(), 1, &models.CreateTransactionRequest{
		Amount:          10,
		Type:            models.TransactionTypeExpense,
		TransactionDate: time.Date(2026, 4, 18, 10, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("Create returned error with nil queue: %v", err)
	}
}

var _ = queue.Client{}
