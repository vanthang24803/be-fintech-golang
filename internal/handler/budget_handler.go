package handler

import (
	"context"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/maynguyen24/sever/internal/models"
	"github.com/maynguyen24/sever/pkg/response"
	"github.com/maynguyen24/sever/pkg/validator"
)

// BudgetService defines the contract for the handler layer
type BudgetService interface {
	Create(ctx context.Context, userID int64, req *models.CreateBudgetRequest) (*models.Budget, error)
	GetList(ctx context.Context, userID int64) ([]*models.BudgetResponse, error)
	GetDetail(ctx context.Context, id, userID int64) (*models.BudgetResponse, error)
	Update(ctx context.Context, id, userID int64, req *models.UpdateBudgetRequest) (*models.Budget, error)
	Delete(ctx context.Context, id, userID int64) error
}

type BudgetHandler struct {
	service BudgetService
}

func NewBudgetHandler(service BudgetService) *BudgetHandler {
	return &BudgetHandler{service: service}
}

// POST /api/v1/budgets/create
func (h *BudgetHandler) Create(c *fiber.Ctx) error {
	userID, err := extractUserID(c)
	if err != nil {
		return err
	}

	var req models.CreateBudgetRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if err := validator.Validate(&req); err != nil {
		return err
	}

	budget, err := h.service.Create(c.Context(), userID, &req)
	if err != nil {
		return err
	}

	return response.Success(c, 2001, "Budget created successfully", budget)
}

// POST /api/v1/budgets/list
func (h *BudgetHandler) List(c *fiber.Ctx) error {
	userID, err := extractUserID(c)
	if err != nil {
		return err
	}

	budgets, err := h.service.GetList(c.Context(), userID)
	if err != nil {
		return err
	}

	return response.Success(c, 2000, "Budgets fetched successfully", budgets)
}

// POST /api/v1/budgets/detail/:id
func (h *BudgetHandler) GetDetail(c *fiber.Ctx) error {
	userID, err := extractUserID(c)
	if err != nil {
		return err
	}

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid budget ID")
	}

	budget, err := h.service.GetDetail(c.Context(), id, userID)
	if err != nil {
		return err
	}

	return response.Success(c, 2000, "Budget details fetched successfully", budget)
}

// POST /api/v1/budgets/update/:id
func (h *BudgetHandler) Update(c *fiber.Ctx) error {
	userID, err := extractUserID(c)
	if err != nil {
		return err
	}

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid budget ID")
	}

	var req models.UpdateBudgetRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if err := validator.Validate(&req); err != nil {
		return err
	}

	budget, err := h.service.Update(c.Context(), id, userID, &req)
	if err != nil {
		return err
	}

	return response.Success(c, 2000, "Budget updated successfully", budget)
}

// POST /api/v1/budgets/delete/:id
func (h *BudgetHandler) Delete(c *fiber.Ctx) error {
	userID, err := extractUserID(c)
	if err != nil {
		return err
	}

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid budget ID")
	}

	if err := h.service.Delete(c.Context(), id, userID); err != nil {
		return err
	}

	return response.Success(c, 2000, "Budget deleted successfully", nil)
}
