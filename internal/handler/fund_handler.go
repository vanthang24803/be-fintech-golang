package handler

import (
	"context"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/maynguyen24/sever/internal/models"
	"github.com/maynguyen24/sever/pkg/response"
)

// FundService defines the contract for the handler layer
type FundService interface {
	Create(ctx context.Context, userID int64, req *models.CreateFundRequest) (*models.Fund, error)
	GetAll(ctx context.Context, userID int64) ([]*models.Fund, error)
	GetByID(ctx context.Context, id, userID int64) (*models.Fund, error)
	Update(ctx context.Context, id, userID int64, req *models.UpdateFundRequest) (*models.Fund, error)
	Delete(ctx context.Context, id, userID int64) error
	Deposit(ctx context.Context, id, userID int64, req *models.FundTransactionRequest) (*models.Fund, error)
	Withdraw(ctx context.Context, id, userID int64, req *models.FundTransactionRequest) (*models.Fund, error)
}

type FundHandler struct {
	service FundService
}

func NewFundHandler(service FundService) *FundHandler {
	return &FundHandler{service: service}
}

// POST /api/v1/funds
func (h *FundHandler) Create(c *fiber.Ctx) error {
	userID, err := extractUserID(c)
	if err != nil {
		return err
	}

	var req models.CreateFundRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	fund, err := h.service.Create(c.Context(), userID, &req)
	if err != nil {
		return err
	}

	return response.Success(c, 2001, "Fund created successfully", fund)
}

// GET /api/v1/funds
func (h *FundHandler) GetAll(c *fiber.Ctx) error {
	userID, err := extractUserID(c)
	if err != nil {
		return err
	}

	funds, err := h.service.GetAll(c.Context(), userID)
	if err != nil {
		return err
	}

	return response.Success(c, 2000, "Funds fetched successfully", funds)
}

// GET /api/v1/funds/:id
func (h *FundHandler) GetByID(c *fiber.Ctx) error {
	userID, err := extractUserID(c)
	if err != nil {
		return err
	}

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid fund ID")
	}

	fund, err := h.service.GetByID(c.Context(), id, userID)
	if err != nil {
		return err
	}

	return response.Success(c, 2000, "Fund fetched successfully", fund)
}

// PUT /api/v1/funds/:id
func (h *FundHandler) Update(c *fiber.Ctx) error {
	userID, err := extractUserID(c)
	if err != nil {
		return err
	}

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid fund ID")
	}

	var req models.UpdateFundRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	fund, err := h.service.Update(c.Context(), id, userID, &req)
	if err != nil {
		return err
	}

	return response.Success(c, 2000, "Fund updated successfully", fund)
}

// DELETE /api/v1/funds/:id
func (h *FundHandler) Delete(c *fiber.Ctx) error {
	userID, err := extractUserID(c)
	if err != nil {
		return err
	}

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid fund ID")
	}

	if err := h.service.Delete(c.Context(), id, userID); err != nil {
		return err
	}

	return response.Success(c, 2000, "Fund deleted successfully", nil)
}

// POST /api/v1/funds/:id/deposit
func (h *FundHandler) Deposit(c *fiber.Ctx) error {
	userID, err := extractUserID(c)
	if err != nil {
		return err
	}

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid fund ID")
	}

	var req models.FundTransactionRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	fund, err := h.service.Deposit(c.Context(), id, userID, &req)
	if err != nil {
		return err
	}

	return response.Success(c, 2000, "Deposit successful", fund)
}

// POST /api/v1/funds/:id/withdraw
func (h *FundHandler) Withdraw(c *fiber.Ctx) error {
	userID, err := extractUserID(c)
	if err != nil {
		return err
	}

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid fund ID")
	}

	var req models.FundTransactionRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	fund, err := h.service.Withdraw(c.Context(), id, userID, &req)
	if err != nil {
		return err
	}

	return response.Success(c, 2000, "Withdrawal successful", fund)
}
