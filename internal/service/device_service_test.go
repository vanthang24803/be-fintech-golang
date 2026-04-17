package service

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/maynguyen24/sever/internal/models"
	"github.com/maynguyen24/sever/pkg/apperr"
)

type stubDeviceRepository struct {
	getByFingerprintResp *models.Device
	getByFingerprintErr  error
	createErr            error
	getByUserIDResp      []*models.Device
	getByUserIDErr       error
	deleteErr            error

	createDevice     *models.Device
	deleteID         int64
	deleteUserID     int64
	getByUserIDArg   int64
	fingerprintArg   string
	updateLastUsedID int64
}

func (s *stubDeviceRepository) Create(ctx context.Context, device *models.Device) error {
	s.createDevice = device
	return s.createErr
}

func (s *stubDeviceRepository) GetByFingerprint(ctx context.Context, fingerprint string) (*models.Device, error) {
	s.fingerprintArg = fingerprint
	return s.getByFingerprintResp, s.getByFingerprintErr
}

func (s *stubDeviceRepository) GetByUserID(ctx context.Context, userID int64) ([]*models.Device, error) {
	s.getByUserIDArg = userID
	return s.getByUserIDResp, s.getByUserIDErr
}

func (s *stubDeviceRepository) GetByID(ctx context.Context, id, userID int64) (*models.Device, error) {
	return nil, nil
}

func (s *stubDeviceRepository) Delete(ctx context.Context, id, userID int64) error {
	s.deleteID = id
	s.deleteUserID = userID
	return s.deleteErr
}

func (s *stubDeviceRepository) UpdateLastUsed(ctx context.Context, id int64) error {
	s.updateLastUsedID = id
	return nil
}

func TestDeviceService_Register_CreatesTrustedActiveDevice(t *testing.T) {
	repo := &stubDeviceRepository{}
	svc := NewDeviceService(repo)
	name := "iPhone"
	push := "push-token"
	req := &models.RegisterDeviceRequest{
		DeviceFingerprint: "fingerprint-123456",
		DeviceName:        &name,
		Platform:          models.PlatformIOS,
		PushToken:         &push,
	}

	device, err := svc.Register(context.Background(), 99, req)
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	if repo.fingerprintArg != req.DeviceFingerprint {
		t.Fatalf("expected fingerprint %q, got %q", req.DeviceFingerprint, repo.fingerprintArg)
	}
	if repo.createDevice == nil {
		t.Fatal("expected repository Create to be called")
	}
	if device != repo.createDevice {
		t.Fatal("expected service to return created device pointer")
	}
	if got, want := device.UserID, int64(99); got != want {
		t.Fatalf("expected userID %d, got %d", want, got)
	}
	if got, want := device.DeviceFingerprint, req.DeviceFingerprint; got != want {
		t.Fatalf("expected fingerprint %q, got %q", want, got)
	}
	if got, want := device.DeviceName, req.DeviceName; got == nil || *got != *want {
		t.Fatalf("expected device name %q, got %v", *want, got)
	}
	if got, want := device.Platform, req.Platform; got != want {
		t.Fatalf("expected platform %q, got %q", want, got)
	}
	if got, want := device.PushToken, req.PushToken; got == nil || *got != *want {
		t.Fatalf("expected push token %q, got %v", *want, got)
	}
	if !device.IsActive {
		t.Fatal("expected device to be active")
	}
	if !device.IsTrusted {
		t.Fatal("expected device to be trusted")
	}
}

func TestDeviceService_Register_RejectsFingerprintOwnedByAnotherUser(t *testing.T) {
	repo := &stubDeviceRepository{
		getByFingerprintResp: &models.Device{ID: 7, UserID: 100},
	}
	svc := NewDeviceService(repo)

	_, err := svc.Register(context.Background(), 99, &models.RegisterDeviceRequest{
		DeviceFingerprint: "fingerprint-123456",
		Platform:          models.PlatformWeb,
	})
	if err == nil || !errors.Is(err, apperr.ErrConflict) || !strings.Contains(err.Error(), "already associated with another account") {
		t.Fatalf("expected conflict error, got %v", err)
	}
	if repo.createDevice != nil {
		t.Fatal("expected repository Create not to be called on conflict")
	}
}

func TestDeviceService_Register_WrapsRepositoryErrors(t *testing.T) {
	repoErr := errors.New("db unavailable")
	svc := NewDeviceService(&stubDeviceRepository{
		getByFingerprintErr: repoErr,
	})

	_, err := svc.Register(context.Background(), 99, &models.RegisterDeviceRequest{
		DeviceFingerprint: "fingerprint-123456",
		Platform:          models.PlatformAndroid,
	})
	if err == nil || !errors.Is(err, repoErr) || !strings.Contains(err.Error(), "failed to check existing device") {
		t.Fatalf("expected wrapped repository error, got %v", err)
	}
}

func TestDeviceService_GetList_PassesThroughResultsAndErrors(t *testing.T) {
	devices := []*models.Device{{ID: 1}, {ID: 2}}
	repo := &stubDeviceRepository{getByUserIDResp: devices}
	svc := NewDeviceService(repo)

	got, err := svc.GetList(context.Background(), 77)
	if err != nil {
		t.Fatalf("GetList() error = %v", err)
	}
	if repo.getByUserIDArg != 77 {
		t.Fatalf("expected userID 77, got %d", repo.getByUserIDArg)
	}
	if len(got) != len(devices) || got[0] != devices[0] || got[1] != devices[1] {
		t.Fatalf("expected repository results to pass through, got %#v", got)
	}

	listErr := errors.New("list failed")
	svc = NewDeviceService(&stubDeviceRepository{getByUserIDErr: listErr})
	_, err = svc.GetList(context.Background(), 77)
	if err == nil || !errors.Is(err, listErr) || !strings.Contains(err.Error(), "failed to fetch devices") {
		t.Fatalf("expected wrapped list error, got %v", err)
	}
}

func TestDeviceService_Remove_MapsNotFound(t *testing.T) {
	repo := &stubDeviceRepository{deleteErr: apperr.ErrNotFound}
	svc := NewDeviceService(repo)

	err := svc.Remove(context.Background(), 77, 15)
	if err == nil || !errors.Is(err, apperr.ErrNotFound) || !strings.Contains(err.Error(), "Device not found") {
		t.Fatalf("expected not found error, got %v", err)
	}
	if repo.deleteID != 15 || repo.deleteUserID != 77 {
		t.Fatalf("expected delete to be called with deviceID 15 and userID 77, got %d/%d", repo.deleteID, repo.deleteUserID)
	}
}
