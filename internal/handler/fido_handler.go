package handler

import (
	"context"
	"strconv"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/gofiber/fiber/v2"
	"github.com/maynguyen24/sever/internal/models"
	"github.com/maynguyen24/sever/pkg/response"
)

// FIDOService defines the contract for the FIDO handler layer
type FIDOService interface {
	BeginEnrollment(ctx context.Context, userID, deviceID int64) (*protocol.CredentialCreation, error)
	FinishEnrollment(ctx context.Context, userID, deviceID int64, body []byte) (*models.Device, error)
	BeginAuthentication(ctx context.Context, credentialID string) (*protocol.CredentialAssertion, error)
	FinishAuthentication(ctx context.Context, body []byte) (string, error)
}

// FIDOHandler handles FIDO2 WebAuthn endpoints
type FIDOHandler struct {
	service FIDOService
}

func NewFIDOHandler(service FIDOService) *FIDOHandler {
	return &FIDOHandler{service: service}
}

// POST /api/v1/devices/biometric/enroll/:id/begin
func (h *FIDOHandler) BeginEnrollment(c *fiber.Ctx) error {
	userID, err := extractUserID(c)
	if err != nil {
		return err
	}

	deviceID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid device ID")
	}

	options, err := h.service.BeginEnrollment(c.Context(), userID, deviceID)
	if err != nil {
		return err
	}

	return response.Success(c, 2000, "Biometric enrollment started", options)
}

// POST /api/v1/devices/biometric/enroll/:id/finish
func (h *FIDOHandler) FinishEnrollment(c *fiber.Ctx) error {
	userID, err := extractUserID(c)
	if err != nil {
		return err
	}

	deviceID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid device ID")
	}

	device, err := h.service.FinishEnrollment(c.Context(), userID, deviceID, c.Body())
	if err != nil {
		return err
	}

	return response.Success(c, 2001, "Biometric enrolled successfully", device)
}

// POST /api/v1/auth/biometric/begin
func (h *FIDOHandler) BeginAuthentication(c *fiber.Ctx) error {
	var req struct {
		CredentialID string `json:"credential_id"`
	}
	if err := c.BodyParser(&req); err != nil || req.CredentialID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "credential_id is required")
	}

	options, err := h.service.BeginAuthentication(c.Context(), req.CredentialID)
	if err != nil {
		return err
	}

	return response.Success(c, 2000, "Biometric challenge issued", options)
}

// POST /api/v1/auth/biometric/finish
func (h *FIDOHandler) FinishAuthentication(c *fiber.Ctx) error {
	accessToken, err := h.service.FinishAuthentication(c.Context(), c.Body())
	if err != nil {
		return err
	}

	return response.Success(c, 2000, "Biometric verified", fiber.Map{
		"access_token": accessToken,
	})
}
