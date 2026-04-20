package service

import (
	"context"
	"errors"
	"testing"

	"github.com/maynguyen24/sever/configs"
	"github.com/maynguyen24/sever/internal/models"
	"github.com/maynguyen24/sever/pkg/apperr"
)

// Auth service extra coverage

func TestAuthService_GetGoogleAuthURL(t *testing.T) {
	svc := NewAuthService(&stubAuthUserRepo{}, &stubTokenRepo{}, &configs.Config{})
	url := svc.GetGoogleAuthURL(context.Background())
	if url == "" {
		t.Fatal("expected non-empty Google auth URL")
	}
}

// Device service extra coverage

func TestDeviceService_Register_ExistingSameUser(t *testing.T) {
	existing := &models.Device{ID: 10, UserID: 42}
	repo := &stubDeviceRepository{
		getByFingerprintResp: existing,
	}
	svc := NewDeviceService(repo)
	push := "new-token"
	got, err := svc.Register(context.Background(), 42, &models.RegisterDeviceRequest{
		DeviceFingerprint: "fp",
		PushToken:         &push,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got.ID != 10 {
		t.Fatalf("expected returned existing device ID 10, got %d", got.ID)
	}
	if repo.updateLastUsedID != 10 {
		t.Fatalf("expected UpdateLastUsed called with ID 10, got %d", repo.updateLastUsedID)
	}
}

func TestDeviceService_Remove_GenericError(t *testing.T) {
	dbErr := errors.New("connection lost")
	repo := &stubDeviceRepository{deleteErr: dbErr}
	svc := NewDeviceService(repo)
	err := svc.Remove(context.Background(), 42, 5)
	if !errors.Is(err, dbErr) {
		t.Fatalf("expected wrapped db error, got %v", err)
	}
}

// SourcePayment service extra coverage

func TestSourcePaymentService_GetByID_ReturnsNotFound(t *testing.T) {
	svc := NewSourcePaymentService(&stubSourcePaymentServiceRepo{
		getByIDFn: func(ctx context.Context, id, userID int64) (*models.SourcePayment, error) {
			return nil, nil
		},
	})
	_, err := svc.GetByID(context.Background(), 1, 42)
	if !errors.Is(err, apperr.ErrNotFound) {
		t.Fatalf("expected not found, got %v", err)
	}
}

func TestSourcePaymentService_GetByID_RepoError(t *testing.T) {
	dbErr := errors.New("db error")
	svc := NewSourcePaymentService(&stubSourcePaymentServiceRepo{
		getByIDFn: func(ctx context.Context, id, userID int64) (*models.SourcePayment, error) {
			return nil, dbErr
		},
	})
	_, err := svc.GetByID(context.Background(), 1, 42)
	if !errors.Is(err, dbErr) {
		t.Fatalf("expected db error, got %v", err)
	}
}

// Transaction service extra coverage

func TestTransactionService_GetByID_ReturnsNotFound(t *testing.T) {
	svc := NewTransactionService(
		&stubTransactionRepo{
			getByIDFn: func(ctx context.Context, id, userID int64) (*models.TransactionDetail, error) {
				return nil, nil
			},
		},
		&stubBudgetRepo{},
		&stubNotificationRepo{},
		nil,
	)
	_, err := svc.GetByID(context.Background(), 1, 42)
	if !errors.Is(err, apperr.ErrNotFound) {
		t.Fatalf("expected not found, got %v", err)
	}
}

func TestTransactionService_GetByID_ReturnsRepoError(t *testing.T) {
	dbErr := errors.New("db error")
	svc := NewTransactionService(
		&stubTransactionRepo{
			getByIDFn: func(ctx context.Context, id, userID int64) (*models.TransactionDetail, error) {
				return nil, dbErr
			},
		},
		&stubBudgetRepo{},
		&stubNotificationRepo{},
		nil,
	)
	_, err := svc.GetByID(context.Background(), 1, 42)
	if !errors.Is(err, dbErr) {
		t.Fatalf("expected db error, got %v", err)
	}
}

// Fund service extra coverage

func TestFundService_Delete_Success(t *testing.T) {
	deleted := false
	repo := &stubFundRepo{
		deleteFn: func(ctx context.Context, id, userID int64) error {
			deleted = true
			return nil
		},
	}
	svc := NewFundService(repo)
	if err := svc.Delete(context.Background(), 1, 42); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !deleted {
		t.Fatal("expected repo.Delete to be called")
	}
}

func TestFundService_GetByID_RepoError(t *testing.T) {
	dbErr := errors.New("db error")
	repo := &stubFundRepo{
		getByIDFn: func(ctx context.Context, id, userID int64) (*models.Fund, error) {
			return nil, dbErr
		},
	}
	svc := NewFundService(repo)
	_, err := svc.GetByID(context.Background(), 1, 42)
	if !errors.Is(err, dbErr) {
		t.Fatalf("expected db error, got %v", err)
	}
}

// Category service extra coverage

func TestCategoryService_GetByID_RepoError(t *testing.T) {
	dbErr := errors.New("db error")
	svc := NewCategoryService(&stubCategoryRepo{
		getByIDFn: func(ctx context.Context, id, userID int64) (*models.Category, error) {
			return nil, dbErr
		},
	})
	_, err := svc.GetByID(context.Background(), 1, 42)
	if !errors.Is(err, dbErr) {
		t.Fatalf("expected db error, got %v", err)
	}
}

// Notification service extra coverage

// SourcePayment service error coverage

func TestSourcePaymentService_GetAll_Error(t *testing.T) {
	t.Parallel()
	dbErr := errors.New("db error")
	svc := NewSourcePaymentService(&stubSourcePaymentServiceRepo{
		getAllByUserIDFn: func(ctx context.Context, userID int64) ([]*models.SourcePayment, error) {
			return nil, dbErr
		},
	})
	_, err := svc.GetAll(context.Background(), 42)
	if !errors.Is(err, dbErr) {
		t.Fatalf("expected db error, got %v", err)
	}
}

func TestSourcePaymentService_Delete_GetByIDError(t *testing.T) {
	t.Parallel()
	dbErr := errors.New("db error")
	svc := NewSourcePaymentService(&stubSourcePaymentServiceRepo{
		getByIDFn: func(ctx context.Context, id, userID int64) (*models.SourcePayment, error) {
			return nil, dbErr
		},
	})
	err := svc.Delete(context.Background(), 1, 42)
	if !errors.Is(err, dbErr) {
		t.Fatalf("expected db error, got %v", err)
	}
}

// Fund service error coverage

func TestFundService_GetAll_Error(t *testing.T) {
	t.Parallel()
	dbErr := errors.New("db error")
	repo := &stubFundServiceAllRepo{
		getAllByUserIDFn: func(ctx context.Context, userID int64) ([]*models.Fund, error) {
			return nil, dbErr
		},
	}
	svc := NewFundService(repo)
	_, err := svc.GetAll(context.Background(), 42)
	if !errors.Is(err, dbErr) {
		t.Fatalf("expected db error, got %v", err)
	}
}

func TestFundService_Deposit_GenericError(t *testing.T) {
	t.Parallel()
	dbErr := errors.New("deposit failed")
	repo := &stubFundRepo{
		depositFn: func(ctx context.Context, id, userID int64, amount float64) (*models.Fund, error) {
			return nil, dbErr
		},
	}
	svc := NewFundService(repo)
	_, err := svc.Deposit(context.Background(), 1, 42, &models.FundTransactionRequest{Amount: 100})
	if !errors.Is(err, dbErr) {
		t.Fatalf("expected db error, got %v", err)
	}
}

// Category service error coverage

func TestCategoryService_Delete_Error(t *testing.T) {
	t.Parallel()
	dbErr := errors.New("db error")
	owned := int64(42)
	svc := NewCategoryService(&stubCategoryRepo{
		getOwnedByIDFn: func(ctx context.Context, id, userID int64) (*models.Category, error) {
			return &models.Category{ID: 1, UserID: &owned}, nil
		},
		deleteFn: func(ctx context.Context, id, userID int64) error {
			return dbErr
		},
	})
	err := svc.Delete(context.Background(), 1, 42)
	if !errors.Is(err, dbErr) {
		t.Fatalf("expected db error, got %v", err)
	}
}

// NewFIDOService can be initialized with valid config

func TestNewFIDOService_ValidConfig(t *testing.T) {
	cfg := &configs.Config{
		WebAuthnRPID:   "localhost",
		WebAuthnRPName: "Test App",
		WebAuthnOrigin: "http://localhost",
	}
	_, err := NewFIDOService(nil, nil, cfg)
	if err != nil {
		t.Fatalf("expected no error creating FIDOService with valid config, got %v", err)
	}
}

// Notification service extra coverage

func TestNotificationService_Delete_Error(t *testing.T) {
	dbErr := errors.New("db error")
	svc := NewNotificationService(&stubFullNotifRepo{
		deleteNotifFn: func(ctx context.Context, userID, id int64) error {
			return dbErr
		},
	}, &stubNotifDeviceRepo{}, &stubPushClient{}, nil)
	err := svc.Delete(context.Background(), 42, 1)
	if !errors.Is(err, dbErr) {
		t.Fatalf("expected db error, got %v", err)
	}
}

func TestNotificationService_GetList_Error(t *testing.T) {
	dbErr := errors.New("db error")
	svc := NewNotificationService(&stubFullNotifRepo{
		getByUserIDFn: func(ctx context.Context, userID int64, filter models.NotificationFilter) ([]*models.Notification, error) {
			return nil, dbErr
		},
	}, &stubNotifDeviceRepo{}, &stubPushClient{}, nil)
	_, err := svc.GetList(context.Background(), 42, models.NotificationFilter{})
	if !errors.Is(err, dbErr) {
		t.Fatalf("expected db error, got %v", err)
	}
}

func TestNotificationService_GetUnreadCount_Error(t *testing.T) {
	dbErr := errors.New("db error")
	svc := NewNotificationService(&stubFullNotifRepo{
		getUnreadCountFn: func(ctx context.Context, userID int64) (int, error) {
			return 0, dbErr
		},
	}, &stubNotifDeviceRepo{}, &stubPushClient{}, nil)
	_, err := svc.GetUnreadCount(context.Background(), 42)
	if !errors.Is(err, dbErr) {
		t.Fatalf("expected db error, got %v", err)
	}
}
