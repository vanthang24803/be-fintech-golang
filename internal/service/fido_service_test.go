package service

import (
	"context"
	"errors"
	"testing"

	"github.com/maynguyen24/sever/configs"
	"github.com/maynguyen24/sever/internal/models"
	"github.com/maynguyen24/sever/pkg/apperr"
)

type stubFIDODeviceRepo struct {
	getByIDFn          func(context.Context, int64, int64) (*models.Device, error)
	getByCredentialFn  func(context.Context, string) (*models.Device, error)
	updateFIDOFn       func(context.Context, int64, string, string, string, int64) error
	updateSignCountFn  func(context.Context, int64, int64) error
}

func (s *stubFIDODeviceRepo) GetByID(ctx context.Context, id, userID int64) (*models.Device, error) {
	if s.getByIDFn != nil {
		return s.getByIDFn(ctx, id, userID)
	}
	return nil, nil
}

func (s *stubFIDODeviceRepo) GetByCredentialID(ctx context.Context, credentialID string) (*models.Device, error) {
	if s.getByCredentialFn != nil {
		return s.getByCredentialFn(ctx, credentialID)
	}
	return nil, nil
}

func (s *stubFIDODeviceRepo) UpdateFIDOCredential(ctx context.Context, deviceID int64, credentialID, publicKey, aaguid string, signCount int64) error {
	if s.updateFIDOFn != nil {
		return s.updateFIDOFn(ctx, deviceID, credentialID, publicKey, aaguid, signCount)
	}
	return nil
}

func (s *stubFIDODeviceRepo) UpdateSignCount(ctx context.Context, deviceID, signCount int64) error {
	if s.updateSignCountFn != nil {
		return s.updateSignCountFn(ctx, deviceID, signCount)
	}
	return nil
}

func newTestFIDOService(t *testing.T) *FIDOService {
	t.Helper()
	svc, err := NewFIDOService(nil, nil, &configs.Config{
		WebAuthnRPID:   "localhost",
		WebAuthnRPName: "Test",
		WebAuthnOrigin: "http://localhost",
	})
	if err != nil {
		t.Fatalf("NewFIDOService: %v", err)
	}
	return svc
}

func TestFIDOService_BeginEnrollment_DeviceNotFound(t *testing.T) {
	t.Parallel()
	svc := newTestFIDOService(t)
	svc.deviceRepo = &stubFIDODeviceRepo{
		getByIDFn: func(ctx context.Context, id, userID int64) (*models.Device, error) {
			return nil, nil
		},
	}
	_, err := svc.BeginEnrollment(context.Background(), 42, 1)
	if !errors.Is(err, apperr.ErrNotFound) {
		t.Fatalf("expected not found, got %v", err)
	}
}

func TestFIDOService_BeginEnrollment_DeviceRepoError(t *testing.T) {
	t.Parallel()
	dbErr := errors.New("db error")
	svc := newTestFIDOService(t)
	svc.deviceRepo = &stubFIDODeviceRepo{
		getByIDFn: func(ctx context.Context, id, userID int64) (*models.Device, error) {
			return nil, dbErr
		},
	}
	_, err := svc.BeginEnrollment(context.Background(), 42, 1)
	if !errors.Is(err, dbErr) {
		t.Fatalf("expected db error, got %v", err)
	}
}

func TestFIDOService_FinishEnrollment_DeviceNotFound(t *testing.T) {
	t.Parallel()
	svc := newTestFIDOService(t)
	svc.deviceRepo = &stubFIDODeviceRepo{
		getByIDFn: func(ctx context.Context, id, userID int64) (*models.Device, error) {
			return nil, nil
		},
	}
	_, err := svc.FinishEnrollment(context.Background(), 42, 1, []byte(`{}`))
	if !errors.Is(err, apperr.ErrNotFound) {
		t.Fatalf("expected not found, got %v", err)
	}
}

func TestFIDOService_FinishEnrollment_DeviceRepoError(t *testing.T) {
	t.Parallel()
	dbErr := errors.New("db error")
	svc := newTestFIDOService(t)
	svc.deviceRepo = &stubFIDODeviceRepo{
		getByIDFn: func(ctx context.Context, id, userID int64) (*models.Device, error) {
			return nil, dbErr
		},
	}
	_, err := svc.FinishEnrollment(context.Background(), 42, 1, []byte(`{}`))
	if !errors.Is(err, dbErr) {
		t.Fatalf("expected db error, got %v", err)
	}
}

func TestFIDOService_BeginAuthentication_CredentialNotFound(t *testing.T) {
	t.Parallel()
	svc := newTestFIDOService(t)
	svc.deviceRepo = &stubFIDODeviceRepo{
		getByCredentialFn: func(ctx context.Context, credentialID string) (*models.Device, error) {
			return nil, nil
		},
	}
	_, err := svc.BeginAuthentication(context.Background(), "cred123")
	if !errors.Is(err, apperr.ErrNotFound) {
		t.Fatalf("expected not found, got %v", err)
	}
}

func TestFIDOService_BeginAuthentication_RepoError(t *testing.T) {
	t.Parallel()
	dbErr := errors.New("db error")
	svc := newTestFIDOService(t)
	svc.deviceRepo = &stubFIDODeviceRepo{
		getByCredentialFn: func(ctx context.Context, credentialID string) (*models.Device, error) {
			return nil, dbErr
		},
	}
	_, err := svc.BeginAuthentication(context.Background(), "cred123")
	if !errors.Is(err, dbErr) {
		t.Fatalf("expected db error, got %v", err)
	}
}

func TestFIDOService_FinishAuthentication_InvalidBody(t *testing.T) {
	t.Parallel()
	svc := newTestFIDOService(t)
	svc.deviceRepo = &stubFIDODeviceRepo{}

	// Invalid JSON
	_, err := svc.FinishAuthentication(context.Background(), []byte(`not json`))
	if !errors.Is(err, apperr.ErrInvalidInput) {
		t.Fatalf("expected invalid input, got %v", err)
	}

	// Valid JSON but empty id
	_, err = svc.FinishAuthentication(context.Background(), []byte(`{"id":""}`))
	if !errors.Is(err, apperr.ErrInvalidInput) {
		t.Fatalf("expected invalid input for empty id, got %v", err)
	}
}

func TestFIDOService_FinishAuthentication_DeviceNotFound(t *testing.T) {
	t.Parallel()
	svc := newTestFIDOService(t)
	svc.deviceRepo = &stubFIDODeviceRepo{
		getByCredentialFn: func(ctx context.Context, credentialID string) (*models.Device, error) {
			return nil, nil
		},
	}
	_, err := svc.FinishAuthentication(context.Background(), []byte(`{"id":"cred123"}`))
	if !errors.Is(err, apperr.ErrNotFound) {
		t.Fatalf("expected not found, got %v", err)
	}
}

func TestFIDOService_FinishAuthentication_RepoError(t *testing.T) {
	t.Parallel()
	dbErr := errors.New("db error")
	svc := newTestFIDOService(t)
	svc.deviceRepo = &stubFIDODeviceRepo{
		getByCredentialFn: func(ctx context.Context, credentialID string) (*models.Device, error) {
			return nil, dbErr
		},
	}
	_, err := svc.FinishAuthentication(context.Background(), []byte(`{"id":"cred123"}`))
	if !errors.Is(err, dbErr) {
		t.Fatalf("expected db error, got %v", err)
	}
}
