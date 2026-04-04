package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/maynguyen24/sever/internal/models"
	"github.com/maynguyen24/sever/pkg/response"
	"github.com/maynguyen24/sever/pkg/validator"
)

// UserService interface is defined where it is used (handler layer)
type UserService interface {
	RegisterUser(req *models.RegisterRequest) (*models.User, error)
	GetProfile(userID int64) (*models.ProfileResponse, error)
	UpdateProfile(userID int64, req *models.UpdateProfileRequest) (*models.Profile, error)
}

type UserHandler struct {
	service UserService
}

func NewUserHandler(service UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) Register(c *fiber.Ctx) error {
	var req models.RegisterRequest

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if err := validator.Validate(&req); err != nil {
		return err
	}

	user, err := h.service.RegisterUser(&req)
	if err != nil {
		return err // The global ErrorHandler interceptor will catch this
	}

	res := models.RegisterResponse{User: user}
	return response.Success(c, 2001, "User registered successfully", res)
}

func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	userIDVal := c.Locals("user_id")
	if userIDVal == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	userID, ok := userIDVal.(int64)
	if !ok {
		if uidFloat, ok := userIDVal.(float64); ok {
			userID = int64(uidFloat)
		} else {
			return fiber.NewError(fiber.StatusInternalServerError, "Invalid user ID type in context")
		}
	}

	res, err := h.service.GetProfile(userID)
	if err != nil {
		return err // Global ErrorHandler will tackle this
	}

	return response.Success(c, 2000, "Profile fetched successfully", res)
}

// POST /api/v1/profile/update
func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
	userID, err := extractUserID(c)
	if err != nil {
		return err
	}

	var req models.UpdateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if err := validator.Validate(&req); err != nil {
		return err
	}

	profile, err := h.service.UpdateProfile(userID, &req)
	if err != nil {
		return err
	}

	return response.Success(c, 2000, "Profile updated successfully", profile)
}
