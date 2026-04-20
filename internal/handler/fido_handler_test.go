package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/gofiber/fiber/v2"
	"github.com/maynguyen24/sever/internal/models"
)

type stubFIDOService struct {
	beginEnrollmentFn    func(context.Context, int64, int64) (*protocol.CredentialCreation, error)
	finishEnrollmentFn   func(context.Context, int64, int64, []byte) (*models.Device, error)
	beginAuthFn          func(context.Context, string) (*protocol.CredentialAssertion, error)
	finishAuthFn         func(context.Context, []byte) (string, error)
}

func (s *stubFIDOService) BeginEnrollment(ctx context.Context, userID, deviceID int64) (*protocol.CredentialCreation, error) {
	if s.beginEnrollmentFn != nil {
		return s.beginEnrollmentFn(ctx, userID, deviceID)
	}
	return &protocol.CredentialCreation{}, nil
}

func (s *stubFIDOService) FinishEnrollment(ctx context.Context, userID, deviceID int64, body []byte) (*models.Device, error) {
	if s.finishEnrollmentFn != nil {
		return s.finishEnrollmentFn(ctx, userID, deviceID, body)
	}
	return &models.Device{ID: deviceID}, nil
}

func (s *stubFIDOService) BeginAuthentication(ctx context.Context, credentialID string) (*protocol.CredentialAssertion, error) {
	if s.beginAuthFn != nil {
		return s.beginAuthFn(ctx, credentialID)
	}
	return &protocol.CredentialAssertion{}, nil
}

func (s *stubFIDOService) FinishAuthentication(ctx context.Context, body []byte) (string, error) {
	if s.finishAuthFn != nil {
		return s.finishAuthFn(ctx, body)
	}
	return "access-token-abc", nil
}

func TestFIDOHandler_BeginEnrollment(t *testing.T) {
	t.Parallel()

	h := NewFIDOHandler(&stubFIDOService{})
	app := fiber.New()
	app.Post("/devices/biometric/enroll/:id/begin", withUser(h.BeginEnrollment))

	req := httptest.NewRequest(http.MethodPost, "/devices/biometric/enroll/1/begin", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}
}

func TestFIDOHandler_BeginEnrollment_Unauthorized(t *testing.T) {
	t.Parallel()

	h := NewFIDOHandler(&stubFIDOService{})
	app := fiber.New()
	app.Post("/devices/biometric/enroll/:id/begin", h.BeginEnrollment)

	req := httptest.NewRequest(http.MethodPost, "/devices/biometric/enroll/1/begin", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", res.StatusCode)
	}
}

func TestFIDOHandler_BeginEnrollment_InvalidID(t *testing.T) {
	t.Parallel()

	h := NewFIDOHandler(&stubFIDOService{})
	app := fiber.New()
	app.Post("/devices/biometric/enroll/:id/begin", withUser(h.BeginEnrollment))

	req := httptest.NewRequest(http.MethodPost, "/devices/biometric/enroll/abc/begin", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", res.StatusCode)
	}
}

func TestFIDOHandler_FinishEnrollment(t *testing.T) {
	t.Parallel()

	h := NewFIDOHandler(&stubFIDOService{})
	app := fiber.New()
	app.Post("/devices/biometric/enroll/:id/finish", withUser(h.FinishEnrollment))

	req := httptest.NewRequest(http.MethodPost, "/devices/biometric/enroll/1/finish", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}
}

func TestFIDOHandler_FinishEnrollment_InvalidID(t *testing.T) {
	t.Parallel()

	h := NewFIDOHandler(&stubFIDOService{})
	app := fiber.New()
	app.Post("/devices/biometric/enroll/:id/finish", withUser(h.FinishEnrollment))

	req := httptest.NewRequest(http.MethodPost, "/devices/biometric/enroll/abc/finish", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", res.StatusCode)
	}
}

func TestFIDOHandler_BeginAuthentication(t *testing.T) {
	t.Parallel()

	h := NewFIDOHandler(&stubFIDOService{})
	app := fiber.New()
	app.Post("/auth/biometric/begin", h.BeginAuthentication)

	req := httptest.NewRequest(http.MethodPost, "/auth/biometric/begin", strings.NewReader(`{"credential_id":"cred123"}`))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}
}

func TestFIDOHandler_BeginAuthentication_MissingCredentialID(t *testing.T) {
	t.Parallel()

	h := NewFIDOHandler(&stubFIDOService{})
	app := fiber.New()
	app.Post("/auth/biometric/begin", h.BeginAuthentication)

	req := httptest.NewRequest(http.MethodPost, "/auth/biometric/begin", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", res.StatusCode)
	}
}

func TestFIDOHandler_FinishAuthentication(t *testing.T) {
	t.Parallel()

	h := NewFIDOHandler(&stubFIDOService{})
	app := fiber.New()
	app.Post("/auth/biometric/finish", h.FinishAuthentication)

	req := httptest.NewRequest(http.MethodPost, "/auth/biometric/finish", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}
}
