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

type stubSourcePaymentService struct {
	createFn  func(context.Context, int64, *models.CreateSourcePaymentRequest) (*models.SourcePayment, error)
	getAllFn   func(context.Context, int64) ([]*models.SourcePayment, error)
	getByIDFn func(context.Context, int64, int64) (*models.SourcePayment, error)
	updateFn  func(context.Context, int64, int64, *models.UpdateSourcePaymentRequest) (*models.SourcePayment, error)
	deleteFn  func(context.Context, int64, int64) error
}

func (s *stubSourcePaymentService) Create(ctx context.Context, userID int64, req *models.CreateSourcePaymentRequest) (*models.SourcePayment, error) {
	if s.createFn != nil {
		return s.createFn(ctx, userID, req)
	}
	return &models.SourcePayment{ID: 1, UserID: userID, Name: req.Name, Type: req.Type}, nil
}

func (s *stubSourcePaymentService) GetAll(ctx context.Context, userID int64) ([]*models.SourcePayment, error) {
	if s.getAllFn != nil {
		return s.getAllFn(ctx, userID)
	}
	return []*models.SourcePayment{{ID: 1, UserID: userID}}, nil
}

func (s *stubSourcePaymentService) GetByID(ctx context.Context, id, userID int64) (*models.SourcePayment, error) {
	if s.getByIDFn != nil {
		return s.getByIDFn(ctx, id, userID)
	}
	return &models.SourcePayment{ID: id, UserID: userID}, nil
}

func (s *stubSourcePaymentService) Update(ctx context.Context, id, userID int64, req *models.UpdateSourcePaymentRequest) (*models.SourcePayment, error) {
	if s.updateFn != nil {
		return s.updateFn(ctx, id, userID, req)
	}
	return &models.SourcePayment{ID: id, Name: req.Name}, nil
}

func (s *stubSourcePaymentService) Delete(ctx context.Context, id, userID int64) error {
	if s.deleteFn != nil {
		return s.deleteFn(ctx, id, userID)
	}
	return nil
}

func TestSourcePaymentHandler_Create(t *testing.T) {
	t.Parallel()

	h := NewSourcePaymentHandler(&stubSourcePaymentService{})
	app := fiber.New()
	app.Post("/sources", withUser(h.Create))

	req := httptest.NewRequest(http.MethodPost, "/sources", strings.NewReader(`{"name":"My Wallet","type":"wallet"}`))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}
}

func TestSourcePaymentHandler_Create_Unauthorized(t *testing.T) {
	t.Parallel()

	h := NewSourcePaymentHandler(&stubSourcePaymentService{})
	app := fiber.New()
	app.Post("/sources", h.Create)

	req := httptest.NewRequest(http.MethodPost, "/sources", strings.NewReader(`{"name":"Wallet","type":"wallet"}`))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", res.StatusCode)
	}
}

func TestSourcePaymentHandler_Create_ValidationFail(t *testing.T) {
	t.Parallel()

	h := NewSourcePaymentHandler(&stubSourcePaymentService{})
	app := fiber.New()
	app.Post("/sources", withUser(h.Create))

	// Missing required type
	req := httptest.NewRequest(http.MethodPost, "/sources", strings.NewReader(`{"name":"Wallet"}`))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", res.StatusCode)
	}
}

func TestSourcePaymentHandler_GetAll(t *testing.T) {
	t.Parallel()

	h := NewSourcePaymentHandler(&stubSourcePaymentService{})
	app := fiber.New()
	app.Get("/sources", withUser(h.GetAll))

	req := httptest.NewRequest(http.MethodGet, "/sources", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}
}

func TestSourcePaymentHandler_GetByID(t *testing.T) {
	t.Parallel()

	h := NewSourcePaymentHandler(&stubSourcePaymentService{})
	app := fiber.New()
	app.Get("/sources/:id", withUser(h.GetByID))

	req := httptest.NewRequest(http.MethodGet, "/sources/1", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}
}

func TestSourcePaymentHandler_GetByID_InvalidID(t *testing.T) {
	t.Parallel()

	h := NewSourcePaymentHandler(&stubSourcePaymentService{})
	app := fiber.New()
	app.Get("/sources/:id", withUser(h.GetByID))

	req := httptest.NewRequest(http.MethodGet, "/sources/abc", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", res.StatusCode)
	}
}

func TestSourcePaymentHandler_Update(t *testing.T) {
	t.Parallel()

	h := NewSourcePaymentHandler(&stubSourcePaymentService{})
	app := fiber.New()
	app.Put("/sources/:id", withUser(h.Update))

	req := httptest.NewRequest(http.MethodPut, "/sources/1", strings.NewReader(`{"name":"Updated Wallet","type":"wallet"}`))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}
}

func TestSourcePaymentHandler_Update_InvalidID(t *testing.T) {
	t.Parallel()

	h := NewSourcePaymentHandler(&stubSourcePaymentService{})
	app := fiber.New()
	app.Put("/sources/:id", withUser(h.Update))

	req := httptest.NewRequest(http.MethodPut, "/sources/abc", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", res.StatusCode)
	}
}

func TestSourcePaymentHandler_Delete(t *testing.T) {
	t.Parallel()

	h := NewSourcePaymentHandler(&stubSourcePaymentService{})
	app := fiber.New()
	app.Delete("/sources/:id", withUser(h.Delete))

	req := httptest.NewRequest(http.MethodDelete, "/sources/1", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}
}

func TestSourcePaymentHandler_Delete_InvalidID(t *testing.T) {
	t.Parallel()

	h := NewSourcePaymentHandler(&stubSourcePaymentService{})
	app := fiber.New()
	app.Delete("/sources/:id", withUser(h.Delete))

	req := httptest.NewRequest(http.MethodDelete, "/sources/abc", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", res.StatusCode)
	}
}
