package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/maynguyen24/sever/internal/models"
	"github.com/maynguyen24/sever/pkg/apperr"
)

// DeviceRepository defines the DB contract for this service
type DeviceRepository interface {
	Create(ctx context.Context, device *models.Device) error
	GetByFingerprint(ctx context.Context, fingerprint string) (*models.Device, error)
	GetByUserID(ctx context.Context, userID int64) ([]*models.Device, error)
	GetByID(ctx context.Context, id, userID int64) (*models.Device, error)
	Delete(ctx context.Context, id, userID int64) error
	UpdateLastUsed(ctx context.Context, id int64) error
}

type DeviceService struct {
	repo DeviceRepository
}

func NewDeviceService(repo DeviceRepository) *DeviceService {
	return &DeviceService{repo: repo}
}

// Register adds a new device for the user. 
// Enforces one-device-one-account: a device fingerprint can only be associated with one user.
func (s *DeviceService) Register(ctx context.Context, userID int64, req *models.RegisterDeviceRequest) (*models.Device, error) {
	existing, err := s.repo.GetByFingerprint(ctx, req.DeviceFingerprint)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing device: %w", err)
	}

	if existing != nil {
		// If belongs to a different user, reject (1 device = 1 account rule)
		if existing.UserID != userID {
			return nil, fmt.Errorf("%w: This device is already associated with another account", apperr.ErrConflict)
		}
		if err := s.repo.UpdateLastUsed(ctx, existing.ID); err != nil {
			return nil, fmt.Errorf("failed to update device status: %w", err)
		}
		existing.PushToken = req.PushToken
		return existing, nil
	}

	device := &models.Device{
		UserID:            userID,
		DeviceFingerprint: req.DeviceFingerprint,
		DeviceName:        req.DeviceName,
		Platform:          req.Platform,
		PushToken:         req.PushToken,
		IsActive:          true,
	}

	if err := s.repo.Create(ctx, device); err != nil {
		return nil, fmt.Errorf("failed to register device: %w", err)
	}

	return device, nil
}

func (s *DeviceService) GetList(ctx context.Context, userID int64) ([]*models.Device, error) {
	devices, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch devices: %w", err)
	}
	return devices, nil
}

func (s *DeviceService) Remove(ctx context.Context, userID int64, deviceID int64) error {
	if err := s.repo.Delete(ctx, deviceID, userID); err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			return fmt.Errorf("%w: Device not found", apperr.ErrNotFound)
		}
		return err
	}
	return nil
}
