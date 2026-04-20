package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/maynguyen24/sever/internal/models"
)

type stubNotificationService struct {
	getListFn        func(context.Context, int64, models.NotificationFilter) ([]*models.Notification, error)
	getUnreadCountFn func(context.Context, int64) (int, error)
	markReadFn       func(context.Context, int64, *models.MarkReadRequest) error
	deleteFn         func(context.Context, int64, int64) error
}

func (s *stubNotificationService) GetList(ctx context.Context, userID int64, filter models.NotificationFilter) ([]*models.Notification, error) {
	if s.getListFn != nil {
		return s.getListFn(ctx, userID, filter)
	}
	return []*models.Notification{{ID: 1, UserID: userID}}, nil
}

func (s *stubNotificationService) GetUnreadCount(ctx context.Context, userID int64) (int, error) {
	if s.getUnreadCountFn != nil {
		return s.getUnreadCountFn(ctx, userID)
	}
	return 3, nil
}

func (s *stubNotificationService) MarkRead(ctx context.Context, userID int64, req *models.MarkReadRequest) error {
	if s.markReadFn != nil {
		return s.markReadFn(ctx, userID, req)
	}
	return nil
}

func (s *stubNotificationService) Delete(ctx context.Context, userID int64, id int64) error {
	if s.deleteFn != nil {
		return s.deleteFn(ctx, userID, id)
	}
	return nil
}

func TestNotificationHandler_List(t *testing.T) {
	t.Parallel()

	h := NewNotificationHandler(&stubNotificationService{})
	app := fiber.New()
	app.Post("/notifications/list", withUser(h.List))

	req := httptest.NewRequest(http.MethodPost, "/notifications/list", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}
}

func TestNotificationHandler_List_Unauthorized(t *testing.T) {
	t.Parallel()

	h := NewNotificationHandler(&stubNotificationService{})
	app := fiber.New()
	app.Post("/notifications/list", h.List)

	req := httptest.NewRequest(http.MethodPost, "/notifications/list", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", res.StatusCode)
	}
}

func TestNotificationHandler_UnreadCount(t *testing.T) {
	t.Parallel()

	h := NewNotificationHandler(&stubNotificationService{})
	app := fiber.New()
	app.Post("/notifications/unread-count", withUser(h.UnreadCount))

	req := httptest.NewRequest(http.MethodPost, "/notifications/unread-count", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}
}

func TestNotificationHandler_MarkRead(t *testing.T) {
	t.Parallel()

	h := NewNotificationHandler(&stubNotificationService{})
	app := fiber.New()
	app.Post("/notifications/mark-read", withUser(h.MarkRead))

	req := httptest.NewRequest(http.MethodPost, "/notifications/mark-read", strings.NewReader(`{"ids":[1,2,3]}`))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}
}

func TestNotificationHandler_Delete(t *testing.T) {
	t.Parallel()

	h := NewNotificationHandler(&stubNotificationService{})
	app := fiber.New()
	app.Post("/notifications/:id", withUser(h.Delete))

	req := httptest.NewRequest(http.MethodPost, "/notifications/1", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}
}

func TestNotificationHandler_Delete_InvalidID(t *testing.T) {
	t.Parallel()

	h := NewNotificationHandler(&stubNotificationService{})
	app := fiber.New()
	app.Post("/notifications/:id", withUser(h.Delete))

	req := httptest.NewRequest(http.MethodPost, "/notifications/abc", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", res.StatusCode)
	}
}
