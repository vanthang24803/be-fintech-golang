package handler

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/maynguyen24/sever/internal/models"
	"github.com/maynguyen24/sever/pkg/response"
	"github.com/maynguyen24/sever/pkg/validator"
)

type UserService interface {
	RegisterUser(ctx context.Context, req *models.RegisterRequest) (*models.User, error)
	GetProfile(ctx context.Context, userID int64) (*models.ProfileResponse, error)
	UpdateProfile(ctx context.Context, userID int64, req *models.UpdateProfileRequest) (*models.Profile, error)
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

	user, err := h.service.RegisterUser(c.Context(), &req)
	if err != nil {
		return err
	}

	res := models.RegisterResponse{User: user}
	return response.Success(c, 2001, "User registered successfully", res)
}

func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	userID, err := extractUserID(c)
	if err != nil {
		return err
	}

	res, err := h.service.GetProfile(c.Context(), userID)
	if err != nil {
		return err
	}

	return response.Success(c, 2000, "Profile fetched successfully", res)
}

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

	profile, err := h.service.UpdateProfile(c.Context(), userID, &req)
	if err != nil {
		return err
	}

	return response.Success(c, 2000, "Profile updated successfully", profile)
}
