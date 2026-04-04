package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/maynguyen24/sever/internal/models"
	"github.com/maynguyen24/sever/pkg/response"
	"github.com/maynguyen24/sever/pkg/validator"
)

// CategoryService defines the contract for the handler layer
type CategoryService interface {
	Create(userID int64, req *models.CreateCategoryRequest) (*models.Category, error)
	GetAll(userID int64) ([]*models.Category, error)
	GetByID(id, userID int64) (*models.Category, error)
	Update(id, userID int64, req *models.UpdateCategoryRequest) (*models.Category, error)
	Delete(id, userID int64) error
}

type CategoryHandler struct {
	service CategoryService
}

func NewCategoryHandler(service CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

// POST /api/v1/categories
func (h *CategoryHandler) Create(c *fiber.Ctx) error {
	userID, err := extractUserID(c)
	if err != nil {
		return err
	}

	var req models.CreateCategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if err := validator.Validate(&req); err != nil {
		return err
	}

	cat, err := h.service.Create(userID, &req)
	if err != nil {
		return err
	}

	return response.Success(c, 2001, "Category created successfully", cat)
}

// GET /api/v1/categories
func (h *CategoryHandler) GetAll(c *fiber.Ctx) error {
	userID, err := extractUserID(c)
	if err != nil {
		return err
	}

	cats, err := h.service.GetAll(userID)
	if err != nil {
		return err
	}

	return response.Success(c, 2000, "Categories fetched successfully", cats)
}

// GET /api/v1/categories/:id
func (h *CategoryHandler) GetByID(c *fiber.Ctx) error {
	userID, err := extractUserID(c)
	if err != nil {
		return err
	}

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid category ID")
	}

	cat, err := h.service.GetByID(id, userID)
	if err != nil {
		return err
	}

	return response.Success(c, 2000, "Category fetched successfully", cat)
}

// PUT /api/v1/categories/:id
func (h *CategoryHandler) Update(c *fiber.Ctx) error {
	userID, err := extractUserID(c)
	if err != nil {
		return err
	}

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid category ID")
	}

	var req models.UpdateCategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if err := validator.Validate(&req); err != nil {
		return err
	}

	cat, err := h.service.Update(id, userID, &req)
	if err != nil {
		return err
	}

	return response.Success(c, 2000, "Category updated successfully", cat)
}

// DELETE /api/v1/categories/:id
func (h *CategoryHandler) Delete(c *fiber.Ctx) error {
	userID, err := extractUserID(c)
	if err != nil {
		return err
	}

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid category ID")
	}

	if err := h.service.Delete(id, userID); err != nil {
		return err
	}

	return response.Success(c, 2000, "Category deleted successfully", nil)
}
