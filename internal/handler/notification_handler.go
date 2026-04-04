package handler

import (
	"context"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/maynguyen24/sever/internal/models"
	"github.com/maynguyen24/sever/pkg/response"
)

// NotificationService defines the contract for the handler layer
type NotificationService interface {
	GetList(ctx context.Context, userID int64, filter models.NotificationFilter) ([]*models.Notification, error)
	GetUnreadCount(ctx context.Context, userID int64) (int, error)
	MarkRead(ctx context.Context, userID int64, req *models.MarkReadRequest) error
	Delete(ctx context.Context, userID int64, id int64) error
}

type NotificationHandler struct {
	service NotificationService
}

func NewNotificationHandler(service NotificationService) *NotificationHandler {
	return &NotificationHandler{service: service}
}

// POST /api/v1/notifications/list
func (h *NotificationHandler) List(c *fiber.Ctx) error {
	userID, err := extractUserID(c)
	if err != nil {
		return err
	}

	var filter models.NotificationFilter
	if err := c.BodyParser(&filter); err != nil {
		// Ignore error, use default filter
	}

	notifications, err := h.service.GetList(c.Context(), userID, filter)
	if err != nil {
		return err
	}

	return response.Success(c, 2000, "Notifications fetched successfully", notifications)
}

// POST /api/v1/notifications/unread-count
func (h *NotificationHandler) UnreadCount(c *fiber.Ctx) error {
	userID, err := extractUserID(c)
	if err != nil {
		return err
	}

	count, err := h.service.GetUnreadCount(c.Context(), userID)
	if err != nil {
		return err
	}

	return response.Success(c, 2000, "Unread count fetched successfully", fiber.Map{"unread_count": count})
}

// POST /api/v1/notifications/mark-read
func (h *NotificationHandler) MarkRead(c *fiber.Ctx) error {
	userID, err := extractUserID(c)
	if err != nil {
		return err
	}

	var req models.MarkReadRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if err := h.service.MarkRead(c.Context(), userID, &req); err != nil {
		return err
	}

	return response.Success(c, 2000, "Notifications marked as read", nil)
}

// POST /api/v1/notifications/delete/:id
func (h *NotificationHandler) Delete(c *fiber.Ctx) error {
	userID, err := extractUserID(c)
	if err != nil {
		return err
	}

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid notification ID")
	}

	if err := h.service.Delete(c.Context(), userID, id); err != nil {
		return err
	}

	return response.Success(c, 2000, "Notification deleted successfully", nil)
}
