package handler

import (
	"context"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/maynguyen24/sever/internal/models"
	"github.com/maynguyen24/sever/pkg/response"
)

// DeviceService defines the contract for the handler layer
type DeviceService interface {
	Register(ctx context.Context, userID int64, req *models.RegisterDeviceRequest) (*models.Device, error)
	GetList(ctx context.Context, userID int64) ([]*models.Device, error)
	Remove(ctx context.Context, userID int64, deviceID int64) error
}

type DeviceHandler struct {
	service DeviceService
}

func NewDeviceHandler(service DeviceService) *DeviceHandler {
	return &DeviceHandler{service: service}
}

// POST /api/v1/devices/register
func (h *DeviceHandler) Register(c *fiber.Ctx) error {
	userID, err := extractUserID(c)
	if err != nil {
		return err
	}

	var req models.RegisterDeviceRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	device, err := h.service.Register(c.Context(), userID, &req)
	if err != nil {
		return err
	}

	return response.Success(c, 2001, "Device registered successfully", device)
}

// POST /api/v1/devices/list
func (h *DeviceHandler) List(c *fiber.Ctx) error {
	userID, err := extractUserID(c)
	if err != nil {
		return err
	}

	devices, err := h.service.GetList(c.Context(), userID)
	if err != nil {
		return err
	}

	return response.Success(c, 2000, "Devices fetched successfully", devices)
}

// POST /api/v1/devices/delete/:id
func (h *DeviceHandler) Delete(c *fiber.Ctx) error {
	userID, err := extractUserID(c)
	if err != nil {
		return err
	}

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid device ID")
	}

	if err := h.service.Remove(c.Context(), userID, id); err != nil {
		return err
	}

	return response.Success(c, 2000, "Device deleted successfully", nil)
}
