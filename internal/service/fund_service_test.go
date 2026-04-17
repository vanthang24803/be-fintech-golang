package service

import (
	"context"
	"errors"
	"testing"

	"github.com/maynguyen24/sever/internal/models"
	"github.com/maynguyen24/sever/pkg/apperr"
)

type stubFundRepo struct {
	createFn   func(context.Context, *models.Fund) error
	getByIDFn  func(context.Context, int64, int64) (*models.Fund, error)
	updateFn   func(context.Context, *models.Fund) error
	deleteFn   func(context.Context, int64, int64) error
	depositFn  func(context.Context, int64, int64, float64) (*models.Fund, error)
	withdrawFn func(context.Context, int64, int64, float64) (*models.Fund, error)
}

func (s *stubFundRepo) Create(ctx context.Context, fund *models.Fund) error {
	return s.createFn(ctx, fund)
}

func (s *stubFundRepo) GetAllByUserID(context.Context, int64) ([]*models.Fund, error) {
	return nil, nil
}

func (s *stubFundRepo) GetByID(ctx context.Context, id, userID int64) (*models.Fund, error) {
	return s.getByIDFn(ctx, id, userID)
}

func (s *stubFundRepo) Update(ctx context.Context, fund *models.Fund) error {
	return s.updateFn(ctx, fund)
}

func (s *stubFundRepo) Delete(ctx context.Context, id, userID int64) error {
	return s.deleteFn(ctx, id, userID)
}

func (s *stubFundRepo) Deposit(ctx context.Context, id, userID int64, amount float64) (*models.Fund, error) {
	return s.depositFn(ctx, id, userID, amount)
}

func (s *stubFundRepo) Withdraw(ctx context.Context, id, userID int64, amount float64) (*models.Fund, error) {
	return s.withdrawFn(ctx, id, userID, amount)
}

func TestFundService_Create(t *testing.T) {
	t.Parallel()

	desc := "desc"
	repo := &stubFundRepo{
		createFn: func(ctx context.Context, fund *models.Fund) error {
			if fund.UserID != 7 || fund.Name != "Trip" || fund.Description != &desc || fund.TargetAmount != 500 || fund.Balance != 100 || fund.Currency != "USD" {
				t.Fatalf("unexpected fund passed to repo: %+v", fund)
			}
			return nil
		},
	}
	svc := NewFundService(repo)

	fund, err := svc.Create(context.Background(), 7, &models.CreateFundRequest{
		Name:         "Trip",
		Description:  &desc,
		TargetAmount: 500,
		Balance:      100,
		Currency:     "USD",
	})
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if fund.UserID != 7 || fund.Name != "Trip" || fund.Currency != "USD" {
		t.Fatalf("unexpected fund returned: %+v", fund)
	}
}

func TestFundService_Create_Validation(t *testing.T) {
	t.Parallel()

	svc := NewFundService(&stubFundRepo{
		createFn: func(context.Context, *models.Fund) error {
			t.Fatal("repo.Create should not be called on invalid input")
			return nil
		},
	})

	tests := []struct {
		name string
		req  *models.CreateFundRequest
	}{
		{name: "negative balance", req: &models.CreateFundRequest{Balance: -1}},
		{name: "negative target", req: &models.CreateFundRequest{TargetAmount: -1}},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := svc.Create(context.Background(), 1, tt.req)
			if !errors.Is(err, apperr.ErrInvalidInput) {
				t.Fatalf("expected invalid input error, got %v", err)
			}
		})
	}
}

func TestFundService_GetByID_NotFound(t *testing.T) {
	t.Parallel()

	svc := NewFundService(&stubFundRepo{
		getByIDFn: func(context.Context, int64, int64) (*models.Fund, error) {
			return nil, nil
		},
	})

	_, err := svc.GetByID(context.Background(), 1, 2)
	if !errors.Is(err, apperr.ErrNotFound) {
		t.Fatalf("expected not found, got %v", err)
	}
}

func TestFundService_Update_NotFound(t *testing.T) {
	t.Parallel()

	svc := NewFundService(&stubFundRepo{
		getByIDFn: func(context.Context, int64, int64) (*models.Fund, error) {
			return nil, nil
		},
		updateFn: func(context.Context, *models.Fund) error {
			t.Fatal("repo.Update should not be called when fund is missing")
			return nil
		},
	})

	_, err := svc.Update(context.Background(), 1, 2, &models.UpdateFundRequest{Name: "x", TargetAmount: 10})
	if !errors.Is(err, apperr.ErrNotFound) {
		t.Fatalf("expected not found, got %v", err)
	}
}

func TestFundService_Delete_NotFound(t *testing.T) {
	t.Parallel()

	svc := NewFundService(&stubFundRepo{
		deleteFn: func(context.Context, int64, int64) error {
			return apperr.ErrNotFound
		},
	})

	err := svc.Delete(context.Background(), 1, 2)
	if !errors.Is(err, apperr.ErrNotFound) {
		t.Fatalf("expected not found, got %v", err)
	}
}

func TestFundService_Transactions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		setup   func(*stubFundRepo)
		method  func(*FundService) error
		wantErr error
	}{
		{
			name: "deposit invalid amount",
			setup: func(repo *stubFundRepo) {
				repo.depositFn = func(context.Context, int64, int64, float64) (*models.Fund, error) {
					t.Fatal("repo transaction should not be called for invalid amount")
					return nil, nil
				}
			},
			method: func(svc *FundService) error {
				_, err := svc.Deposit(context.Background(), 1, 2, &models.FundTransactionRequest{Amount: 0})
				return err
			},
			wantErr: apperr.ErrInvalidInput,
		},
		{
			name: "deposit not found",
			setup: func(repo *stubFundRepo) {
				repo.depositFn = func(context.Context, int64, int64, float64) (*models.Fund, error) {
					return nil, apperr.ErrNotFound
				}
			},
			method: func(svc *FundService) error {
				_, err := svc.Deposit(context.Background(), 1, 2, &models.FundTransactionRequest{Amount: 10})
				return err
			},
			wantErr: apperr.ErrNotFound,
		},
		{
			name: "withdraw invalid amount",
			setup: func(repo *stubFundRepo) {
				repo.withdrawFn = func(context.Context, int64, int64, float64) (*models.Fund, error) {
					t.Fatal("repo transaction should not be called for invalid amount")
					return nil, nil
				}
			},
			method: func(svc *FundService) error {
				_, err := svc.Withdraw(context.Background(), 1, 2, &models.FundTransactionRequest{Amount: -5})
				return err
			},
			wantErr: apperr.ErrInvalidInput,
		},
		{
			name: "withdraw not found",
			setup: func(repo *stubFundRepo) {
				repo.withdrawFn = func(context.Context, int64, int64, float64) (*models.Fund, error) {
					return nil, apperr.ErrNotFound
				}
			},
			method: func(svc *FundService) error {
				_, err := svc.Withdraw(context.Background(), 1, 2, &models.FundTransactionRequest{Amount: 10})
				return err
			},
			wantErr: apperr.ErrNotFound,
		},
		{
			name: "withdraw insufficient balance",
			setup: func(repo *stubFundRepo) {
				repo.withdrawFn = func(context.Context, int64, int64, float64) (*models.Fund, error) {
					return nil, apperr.ErrInsufficientBalance
				}
			},
			method: func(svc *FundService) error {
				_, err := svc.Withdraw(context.Background(), 1, 2, &models.FundTransactionRequest{Amount: 10})
				return err
			},
			wantErr: apperr.ErrInvalidInput,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := &stubFundRepo{}
			tt.setup(repo)

			svc := NewFundService(repo)
			err := tt.method(svc)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected %v, got %v", tt.wantErr, err)
			}
		})
	}
}
