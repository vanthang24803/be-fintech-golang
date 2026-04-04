package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/maynguyen24/sever/internal/models"
	"github.com/maynguyen24/sever/pkg/response"
)

type SavingsGoalService interface {
	Create(userID int64, req *models.CreateGoalRequest) (*models.SavingsGoal, error)
	List(userID int64) ([]models.SavingsGoal, error)
	GetDetail(id int64, userID int64) (*models.GoalResponse, error)
	Contribute(userID int64, req *models.GoalContributeRequest) (*models.SavingsGoal, error)
	Withdraw(userID int64, req *models.GoalWithdrawRequest) (*models.SavingsGoal, error)
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

	res, err := h.service.Create(userID, &req)
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
	res, err := h.service.List(userID)
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

	res, err := h.service.GetDetail(id, userID)
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

	res, err := h.service.Contribute(userID, &req)
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

	res, err := h.service.Withdraw(userID, &req)
	if err != nil {
		return err
	}

	return response.Success(c, 2000, "common.success", res)
}
