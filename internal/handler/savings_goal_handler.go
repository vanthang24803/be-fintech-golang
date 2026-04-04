package handler

import (
	"context"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/maynguyen24/sever/internal/models"
	"github.com/maynguyen24/sever/pkg/response"
)

type SavingsGoalService interface {
	Create(ctx context.Context, userID int64, req *models.CreateGoalRequest) (*models.SavingsGoal, error)
	List(ctx context.Context, userID int64) ([]models.SavingsGoal, error)
	GetDetail(ctx context.Context, id int64, userID int64) (*models.GoalResponse, error)
	Contribute(ctx context.Context, userID int64, req *models.GoalContributeRequest) (*models.SavingsGoal, error)
	Withdraw(ctx context.Context, userID int64, req *models.GoalWithdrawRequest) (*models.SavingsGoal, error)
}

type SavingsGoalHandler struct {
	service SavingsGoalService
}

func NewSavingsGoalHandler(service SavingsGoalService) *SavingsGoalHandler {
	return &SavingsGoalHandler{service: service}
}

func (h *SavingsGoalHandler) Create(c *fiber.Ctx) error {
	userID, err := extractUserID(c)
	if err != nil {
		return err
	}
	var req models.CreateGoalRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "common.bad_request")
	}

	res, err := h.service.Create(c.Context(), userID, &req)
	if err != nil {
		return err
	}

	return response.Success(c, 2000, "common.success", res)
}

func (h *SavingsGoalHandler) List(c *fiber.Ctx) error {
	userID, err := extractUserID(c)
	if err != nil {
		return err
	}
	res, err := h.service.List(c.Context(), userID)
	if err != nil {
		return err
	}

	return response.Success(c, 2000, "common.success", res)
}

func (h *SavingsGoalHandler) GetDetail(c *fiber.Ctx) error {
	userID, err := extractUserID(c)
	if err != nil {
		return err
	}
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid ID")
	}

	res, err := h.service.GetDetail(c.Context(), id, userID)
	if err != nil {
		return err
	}

	return response.Success(c, 2000, "common.success", res)
}

func (h *SavingsGoalHandler) Contribute(c *fiber.Ctx) error {
	userID, err := extractUserID(c)
	if err != nil {
		return err
	}
	var req models.GoalContributeRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "common.bad_request")
	}

	res, err := h.service.Contribute(c.Context(), userID, &req)
	if err != nil {
		return err
	}

	return response.Success(c, 2000, "common.success", res)
}

func (h *SavingsGoalHandler) Withdraw(c *fiber.Ctx) error {
	userID, err := extractUserID(c)
	if err != nil {
		return err
	}
	var req models.GoalWithdrawRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "common.bad_request")
	}

	res, err := h.service.Withdraw(c.Context(), userID, &req)
	if err != nil {
		return err
	}

	return response.Success(c, 2000, "common.success", res)
}
