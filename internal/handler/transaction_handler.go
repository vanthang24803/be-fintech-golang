package handler

import (
	"context"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/maynguyen24/sever/internal/models"
	"github.com/maynguyen24/sever/pkg/response"
	"github.com/maynguyen24/sever/pkg/validator"
)

// TransactionService defines the contract for the handler layer
type TransactionService interface {
	Create(ctx context.Context, userID int64, req *models.CreateTransactionRequest) (*models.Transaction, error)
	GetAll(ctx context.Context, userID int64, query map[string]string) ([]*models.TransactionDetail, error)
	GetByID(ctx context.Context, id, userID int64) (*models.TransactionDetail, error)
	Update(ctx context.Context, id, userID int64, req *models.UpdateTransactionRequest) (*models.Transaction, error)
	Delete(ctx context.Context, id, userID int64) error
}

type TransactionHandler struct {
	service TransactionService
}

func NewTransactionHandler(service TransactionService) *TransactionHandler {
	return &TransactionHandler{service: service}
}

// POST /api/v1/transactions
func (h *TransactionHandler) Create(c *fiber.Ctx) error {
	userID, err := extractUserID(c)
	if err != nil {
		return err
	}

	var req models.CreateTransactionRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if err := validator.Validate(&req); err != nil {
		return err
	}

	tx, err := h.service.Create(c.Context(), userID, &req)
	if err != nil {
		return err
	}

	return response.Success(c, 2001, "Transaction created successfully", tx)
}

// POST /api/v1/transactions/list
func (h *TransactionHandler) GetAll(c *fiber.Ctx) error {
	userID, err := extractUserID(c)
	if err != nil {
		return err
	}

	// Filters come from JSON body (all endpoints are POST)
	var body struct {
		Type            string `json:"type"`
		CategoryID      string `json:"category_id"`
		SourcePaymentID string `json:"source_id"`
	}
	// Body is optional — ignore parse error, default to no filter
	_ = c.BodyParser(&body)

	queryParams := map[string]string{
		"type":        body.Type,
		"category_id": body.CategoryID,
		"source_id":   body.SourcePaymentID,
	}

	txs, err := h.service.GetAll(c.Context(), userID, queryParams)
	if err != nil {
		return err
	}

	return response.Success(c, 2000, "Transactions fetched successfully", txs)
}

// GET /api/v1/transactions/:id
func (h *TransactionHandler) GetByID(c *fiber.Ctx) error {
	userID, err := extractUserID(c)
	if err != nil {
		return err
	}

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid transaction ID")
	}

	tx, err := h.service.GetByID(c.Context(), id, userID)
	if err != nil {
		return err
	}

	return response.Success(c, 2000, "Transaction fetched successfully", tx)
}

// PUT /api/v1/transactions/:id
func (h *TransactionHandler) Update(c *fiber.Ctx) error {
	userID, err := extractUserID(c)
	if err != nil {
		return err
	}

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid transaction ID")
	}

	var req models.UpdateTransactionRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if err := validator.Validate(&req); err != nil {
		return err
	}

	tx, err := h.service.Update(c.Context(), id, userID, &req)
	if err != nil {
		return err
	}

	return response.Success(c, 2000, "Transaction updated successfully", tx)
}

// DELETE /api/v1/transactions/:id
func (h *TransactionHandler) Delete(c *fiber.Ctx) error {
	userID, err := extractUserID(c)
	if err != nil {
		return err
	}

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid transaction ID")
	}

	if err := h.service.Delete(c.Context(), id, userID); err != nil {
		return err
	}

	return response.Success(c, 2000, "Transaction deleted successfully", nil)
}
