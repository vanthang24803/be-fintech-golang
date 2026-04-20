package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/maynguyen24/sever/internal/models"
)

type stubDeviceService struct {
	registerFn func(context.Context, int64, *models.RegisterDeviceRequest) (*models.Device, error)
	getListFn  func(context.Context, int64) ([]*models.Device, error)
	removeFn   func(context.Context, int64, int64) error
}

func (s *stubDeviceService) Register(ctx context.Context, userID int64, req *models.RegisterDeviceRequest) (*models.Device, error) {
	if s.registerFn != nil {
		return s.registerFn(ctx, userID, req)
	}
	return &models.Device{ID: 1, UserID: userID, DeviceFingerprint: req.DeviceFingerprint}, nil
}

func (s *stubDeviceService) GetList(ctx context.Context, userID int64) ([]*models.Device, error) {
	if s.getListFn != nil {
		return s.getListFn(ctx, userID)
	}
	return []*models.Device{{ID: 1, UserID: userID}}, nil
}

func (s *stubDeviceService) Remove(ctx context.Context, userID int64, deviceID int64) error {
	if s.removeFn != nil {
		return s.removeFn(ctx, userID, deviceID)
	}
	return nil
}

func TestDeviceHandler_Register(t *testing.T) {
	t.Parallel()

	h := NewDeviceHandler(&stubDeviceService{})
	app := fiber.New()
	app.Post("/devices/register", withUser(h.Register))

	body := `{"device_fingerprint":"abcdef1234567890abcdef1234567890","platform":"android"}`
	req := httptest.NewRequest(http.MethodPost, "/devices/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}
}

func TestDeviceHandler_Register_Unauthorized(t *testing.T) {
	t.Parallel()

	h := NewDeviceHandler(&stubDeviceService{})
	app := fiber.New()
	app.Post("/devices/register", h.Register)

	req := httptest.NewRequest(http.MethodPost, "/devices/register", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", res.StatusCode)
	}
}

func TestDeviceHandler_List(t *testing.T) {
	t.Parallel()

	h := NewDeviceHandler(&stubDeviceService{})
	app := fiber.New()
	app.Post("/devices/list", withUser(h.List))

	req := httptest.NewRequest(http.MethodPost, "/devices/list", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}
}

func TestDeviceHandler_Delete(t *testing.T) {
	t.Parallel()

	h := NewDeviceHandler(&stubDeviceService{})
	app := fiber.New()
	app.Post("/devices/:id", withUser(h.Delete))

	req := httptest.NewRequest(http.MethodPost, "/devices/1", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}
}

func TestDeviceHandler_Delete_InvalidID(t *testing.T) {
	t.Parallel()

	h := NewDeviceHandler(&stubDeviceService{})
	app := fiber.New()
	app.Post("/devices/:id", withUser(h.Delete))

	req := httptest.NewRequest(http.MethodPost, "/devices/abc", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", res.StatusCode)
	}
}
