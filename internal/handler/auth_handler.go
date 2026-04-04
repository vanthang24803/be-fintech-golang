package handler

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/maynguyen24/sever/internal/models"
	"github.com/maynguyen24/sever/pkg/response"
	"github.com/maynguyen24/sever/pkg/validator"
)

// AuthService interface defined where used
type AuthService interface {
	Login(ctx context.Context, req *models.LoginRequest) (*models.LoginResponse, error)
	RefreshToken(ctx context.Context, req *models.RefreshTokenRequest) (*models.TokenPair, error)
	Logout(ctx context.Context, req *models.LogoutRequest) error
	GetGoogleAuthURL(ctx context.Context) string
	HandleGoogleCallback(ctx context.Context, code string) (*models.LoginResponse, error)
}

type AuthHandler struct {
	service AuthService
}

func NewAuthHandler(service AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) GetGoogleAuthURL(c *fiber.Ctx) error {
	url := h.service.GetGoogleAuthURL(c.Context())
	return response.Success(c, 2000, "common.success", fiber.Map{"url": url})
}

func (h *AuthHandler) GoogleCallback(c *fiber.Ctx) error {
	code := c.Query("code")
	if code == "" {
		// Also try to read from body if Frontend sends it as POST
		type callbackReq struct {
			Code string `json:"code"`
		}
		var req callbackReq
		_ = c.BodyParser(&req)
		code = req.Code
	}

	if code == "" {
		return fiber.NewError(fiber.StatusBadRequest, "common.bad_request")
	}

	res, err := h.service.HandleGoogleCallback(c.Context(), code)
	if err != nil {
		return err
	}

	return response.Success(c, 2000, "auth.login_success", res)
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if err := validator.Validate(&req); err != nil {
		return err
	}

	res, err := h.service.Login(c.Context(), &req)
	if err != nil {
		return err // Let the central error handler manage this
	}

	return response.Success(c, 2000, "auth.login_success", res)
}

func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	var req models.RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "common.bad_request")
	}

	if err := validator.Validate(&req); err != nil {
		return err
	}

	res, err := h.service.RefreshToken(c.Context(), &req)
	if err != nil {
		return err
	}

	return response.Success(c, 2000, "auth.refresh_success", res)
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	var req models.LogoutRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "common.bad_request")
	}

	if err := validator.Validate(&req); err != nil {
		return err
	}

	if err := h.service.Logout(c.Context(), &req); err != nil {
		return err
	}

	return response.Success(c, 2000, "auth.logout_success", nil)
}
