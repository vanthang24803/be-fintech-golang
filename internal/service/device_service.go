package service

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/maynguyen24/sever/internal/models"
	"github.com/maynguyen24/sever/pkg/apperr"
)

// DeviceRepository defines the DB contract for this service
type DeviceRepository interface {
	Create(device *models.Device) error
	GetByFingerprint(fingerprint string) (*models.Device, error)
	GetByUserID(userID int64) ([]*models.Device, error)
	GetByID(id, userID int64) (*models.Device, error)
	Delete(id, userID int64) error
	UpdateLastUsed(id int64) error
}

type DeviceService struct {
	repo DeviceRepository
}

func NewDeviceService(repo DeviceRepository) *DeviceService {
	return &DeviceService{repo: repo}
}

// Register adds a new device for the user. 
// Enforces one-device-one-account: a device fingerprint can only be associated with one user.
func (s *DeviceService) Register(userID int64, req *models.RegisterDeviceRequest) (*models.Device, error) {
	existing, err := s.repo.GetByFingerprint(req.DeviceFingerprint)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing device: %w", err)
	}

	if existing != nil {
		// If belongs to a different user, reject (1 device = 1 account rule)
		if existing.UserID != userID {
			return nil, fiber.NewError(fiber.StatusConflict, "This device is already associated with another account")
		}
		if err := s.repo.UpdateLastUsed(existing.ID); err != nil {
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

	if err := s.repo.Create(device); err != nil {
		return nil, fmt.Errorf("failed to register device: %w", err)
	}

	return device, nil
}

func (s *DeviceService) GetList(userID int64) ([]*models.Device, error) {
	devices, err := s.repo.GetByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch devices: %w", err)
	}
	return devices, nil
}

func (s *DeviceService) Remove(userID int64, deviceID int64) error {
	if err := s.repo.Delete(deviceID, userID); err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "Device not found")
		}
		return err
	}
	return nil
}
