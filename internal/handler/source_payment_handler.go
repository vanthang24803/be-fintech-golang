package handler

import (
	"context"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/maynguyen24/sever/internal/models"
	"github.com/maynguyen24/sever/pkg/response"
)

type SourcePaymentService interface {
	Create(ctx context.Context, userID int64, req *models.CreateSourcePaymentRequest) (*models.SourcePayment, error)
	GetAll(ctx context.Context, userID int64) ([]*models.SourcePayment, error)
	GetByID(ctx context.Context, id, userID int64) (*models.SourcePayment, error)
	Update(ctx context.Context, id, userID int64, req *models.UpdateSourcePaymentRequest) (*models.SourcePayment, error)
	Delete(ctx context.Context, id, userID int64) error
}

type SourcePaymentHandler struct {
	service SourcePaymentService
}

func NewSourcePaymentHandler(service SourcePaymentService) *SourcePaymentHandler {
	return &SourcePaymentHandler{service: service}
}

// extractUserID is a helper to get userID from Fiber context (set by AuthRequired middleware)
func extractUserID(c *fiber.Ctx) (int64, error) {
	val := c.Locals("user_id")
	if val == nil {
		return 0, fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}
	uid, ok := val.(int64)
	if !ok {
		return 0, fiber.NewError(fiber.StatusInternalServerError, "Invalid user ID in context")
	}
	return uid, nil
}

// POST /api/v1/sources
func (h *SourcePaymentHandler) Create(c *fiber.Ctx) error {
	userID, err := extractUserID(c)
	if err != nil {
		return err
	}

	var req models.CreateSourcePaymentRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	source, err := h.service.Create(c.Context(), userID, &req)
	if err != nil {
		return err
	}

	return response.Success(c, 2001, "Source payment created successfully", source)
}

// GET /api/v1/sources
func (h *SourcePaymentHandler) GetAll(c *fiber.Ctx) error {
	userID, err := extractUserID(c)
	if err != nil {
		return err
	}

	sources, err := h.service.GetAll(c.Context(), userID)
	if err != nil {
		return err
	}

	return response.Success(c, 2000, "Source payments fetched successfully", sources)
}

// GET /api/v1/sources/:id
func (h *SourcePaymentHandler) GetByID(c *fiber.Ctx) error {
	userID, err := extractUserID(c)
	if err != nil {
		return err
	}

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid source payment ID")
	}

	source, err := h.service.GetByID(c.Context(), id, userID)
	if err != nil {
		return err
	}

	return response.Success(c, 2000, "Source payment fetched successfully", source)
}

// PUT /api/v1/sources/:id
func (h *SourcePaymentHandler) Update(c *fiber.Ctx) error {
	userID, err := extractUserID(c)
	if err != nil {
		return err
	}

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid source payment ID")
	}

	var req models.UpdateSourcePaymentRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	source, err := h.service.Update(c.Context(), id, userID, &req)
	if err != nil {
		return err
	}

	return response.Success(c, 2000, "Source payment updated successfully", source)
}

// DELETE /api/v1/sources/:id
func (h *SourcePaymentHandler) Delete(c *fiber.Ctx) error {
	userID, err := extractUserID(c)
	if err != nil {
		return err
	}

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid source payment ID")
	}

	if err := h.service.Delete(c.Context(), id, userID); err != nil {
		return err
	}

	return response.Success(c, 2000, "Source payment deleted successfully", nil)
}
