package service

import (
	"context"
	"errors"
	"testing"

	"github.com/maynguyen24/sever/internal/models"
	"github.com/maynguyen24/sever/pkg/apperr"
)

type stubSourcePaymentRepo struct {
	createFn  func(context.Context, *models.SourcePayment) error
	getByIDFn func(context.Context, int64, int64) (*models.SourcePayment, error)
	updateFn  func(context.Context, *models.SourcePayment) error
	deleteFn  func(context.Context, int64, int64) error
}

func (s *stubSourcePaymentRepo) Create(ctx context.Context, source *models.SourcePayment) error {
	return s.createFn(ctx, source)
}

func (s *stubSourcePaymentRepo) GetAllByUserID(context.Context, int64) ([]*models.SourcePayment, error) {
	return nil, nil
}

func (s *stubSourcePaymentRepo) GetByID(ctx context.Context, id, userID int64) (*models.SourcePayment, error) {
	return s.getByIDFn(ctx, id, userID)
}

func (s *stubSourcePaymentRepo) Update(ctx context.Context, source *models.SourcePayment) error {
	return s.updateFn(ctx, source)
}

func (s *stubSourcePaymentRepo) Delete(ctx context.Context, id, userID int64) error {
	return s.deleteFn(ctx, id, userID)
}

func TestSourcePaymentService_Create_DefaultCurrencyAndValidation(t *testing.T) {
	t.Parallel()

	repo := &stubSourcePaymentRepo{
		createFn: func(ctx context.Context, source *models.SourcePayment) error {
			if source.Name != "Cash" || source.Type != "wallet" || source.Currency != "VND" {
				t.Fatalf("unexpected source passed to repo: %+v", source)
			}
			return nil
		},
	}
	svc := NewSourcePaymentService(repo)

	source, err := svc.Create(context.Background(), 99, &models.CreateSourcePaymentRequest{
		Name: "Cash",
		Type: "wallet",
	})
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if source.Currency != "VND" || source.UserID != 99 {
		t.Fatalf("unexpected source returned: %+v", source)
	}

	tests := []struct {
		name string
		req  *models.CreateSourcePaymentRequest
	}{
		{name: "missing name", req: &models.CreateSourcePaymentRequest{Type: "wallet"}},
		{name: "missing type", req: &models.CreateSourcePaymentRequest{Name: "Cash"}},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := svc.Create(context.Background(), 99, tt.req)
			if !errors.Is(err, apperr.ErrInvalidInput) {
				t.Fatalf("expected invalid input, got %v", err)
			}
		})
	}
}

func TestSourcePaymentService_GetByID_NotFound(t *testing.T) {
	t.Parallel()

	svc := NewSourcePaymentService(&stubSourcePaymentRepo{
		getByIDFn: func(context.Context, int64, int64) (*models.SourcePayment, error) {
			return nil, nil
		},
	})

	_, err := svc.GetByID(context.Background(), 1, 2)
	if !errors.Is(err, apperr.ErrNotFound) {
		t.Fatalf("expected not found, got %v", err)
	}
}

func TestSourcePaymentService_Update_And_Delete(t *testing.T) {
	t.Parallel()

	repo := &stubSourcePaymentRepo{}
	owned := &models.SourcePayment{ID: 1, UserID: 2, Name: "Cash", Type: "wallet", Currency: "VND"}
	repo.getByIDFn = func(context.Context, int64, int64) (*models.SourcePayment, error) {
		return owned, nil
	}
	repo.updateFn = func(ctx context.Context, source *models.SourcePayment) error {
		if source.Name != "Bank" || source.Type != "bank" || source.Currency != "USD" {
			t.Fatalf("unexpected source passed to update: %+v", source)
		}
		return nil
	}
	deleted := false
	repo.deleteFn = func(context.Context, int64, int64) error {
		deleted = true
		return nil
	}

	svc := NewSourcePaymentService(repo)

	source, err := svc.Update(context.Background(), 1, 2, &models.UpdateSourcePaymentRequest{
		Name:     "Bank",
		Type:     "bank",
		Currency: "USD",
	})
	if err != nil {
		t.Fatalf("Update returned error: %v", err)
	}
	if source.Name != "Bank" || source.Currency != "USD" {
		t.Fatalf("unexpected source returned: %+v", source)
	}

	if err := svc.Delete(context.Background(), 1, 2); err != nil {
		t.Fatalf("Delete returned error: %v", err)
	}
	if !deleted {
		t.Fatal("repo.Delete was not called")
	}

	repo.getByIDFn = func(context.Context, int64, int64) (*models.SourcePayment, error) {
		return nil, nil
	}
	_, err = svc.Update(context.Background(), 1, 2, &models.UpdateSourcePaymentRequest{Name: "Bank", Type: "bank"})
	if !errors.Is(err, apperr.ErrNotFound) {
		t.Fatalf("expected not found, got %v", err)
	}
	err = svc.Delete(context.Background(), 1, 2)
	if !errors.Is(err, apperr.ErrNotFound) {
		t.Fatalf("expected not found, got %v", err)
	}
}
