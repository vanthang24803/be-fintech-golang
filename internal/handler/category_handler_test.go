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

type stubCategoryService struct {
	createFn  func(context.Context, int64, *models.CreateCategoryRequest) (*models.Category, error)
	getAllFn   func(context.Context, int64) ([]*models.Category, error)
	getByIDFn func(context.Context, int64, int64) (*models.Category, error)
	updateFn  func(context.Context, int64, int64, *models.UpdateCategoryRequest) (*models.Category, error)
	deleteFn  func(context.Context, int64, int64) error
}

func (s *stubCategoryService) Create(ctx context.Context, userID int64, req *models.CreateCategoryRequest) (*models.Category, error) {
	if s.createFn != nil {
		return s.createFn(ctx, userID, req)
	}
	return &models.Category{ID: 1, Name: req.Name, Type: req.Type}, nil
}

func (s *stubCategoryService) GetAll(ctx context.Context, userID int64) ([]*models.Category, error) {
	if s.getAllFn != nil {
		return s.getAllFn(ctx, userID)
	}
	return []*models.Category{{ID: 1, Name: "Food"}}, nil
}

func (s *stubCategoryService) GetByID(ctx context.Context, id, userID int64) (*models.Category, error) {
	if s.getByIDFn != nil {
		return s.getByIDFn(ctx, id, userID)
	}
	return &models.Category{ID: id, Name: "Food"}, nil
}

func (s *stubCategoryService) Update(ctx context.Context, id, userID int64, req *models.UpdateCategoryRequest) (*models.Category, error) {
	if s.updateFn != nil {
		return s.updateFn(ctx, id, userID, req)
	}
	return &models.Category{ID: id, Name: req.Name, Type: req.Type}, nil
}

func (s *stubCategoryService) Delete(ctx context.Context, id, userID int64) error {
	if s.deleteFn != nil {
		return s.deleteFn(ctx, id, userID)
	}
	return nil
}

func withUser(h func(*fiber.Ctx) error) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Locals("user_id", int64(42))
		return h(c)
	}
}

func TestCategoryHandler_Create(t *testing.T) {
	t.Parallel()

	h := NewCategoryHandler(&stubCategoryService{})
	app := fiber.New()
	app.Post("/categories", withUser(h.Create))

	req := httptest.NewRequest(http.MethodPost, "/categories", strings.NewReader(`{"name":"Food","type":"expense"}`))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}
}

func TestCategoryHandler_Create_Unauthorized(t *testing.T) {
	t.Parallel()

	h := NewCategoryHandler(&stubCategoryService{})
	app := fiber.New()
	app.Post("/categories", h.Create)

	req := httptest.NewRequest(http.MethodPost, "/categories", strings.NewReader(`{"name":"Food","type":"expense"}`))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", res.StatusCode)
	}
}

func TestCategoryHandler_Create_ValidationFail(t *testing.T) {
	t.Parallel()

	h := NewCategoryHandler(&stubCategoryService{})
	app := fiber.New()
	app.Post("/categories", withUser(h.Create))

	// Missing required type field
	req := httptest.NewRequest(http.MethodPost, "/categories", strings.NewReader(`{"name":"Food"}`))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", res.StatusCode)
	}
}

func TestCategoryHandler_GetAll(t *testing.T) {
	t.Parallel()

	h := NewCategoryHandler(&stubCategoryService{})
	app := fiber.New()
	app.Get("/categories", withUser(h.GetAll))

	req := httptest.NewRequest(http.MethodGet, "/categories", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}
}

func TestCategoryHandler_GetByID(t *testing.T) {
	t.Parallel()

	h := NewCategoryHandler(&stubCategoryService{})
	app := fiber.New()
	app.Get("/categories/:id", withUser(h.GetByID))

	req := httptest.NewRequest(http.MethodGet, "/categories/1", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}
}

func TestCategoryHandler_GetByID_InvalidID(t *testing.T) {
	t.Parallel()

	h := NewCategoryHandler(&stubCategoryService{})
	app := fiber.New()
	app.Get("/categories/:id", withUser(h.GetByID))

	req := httptest.NewRequest(http.MethodGet, "/categories/abc", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", res.StatusCode)
	}
}

func TestCategoryHandler_Update(t *testing.T) {
	t.Parallel()

	h := NewCategoryHandler(&stubCategoryService{})
	app := fiber.New()
	app.Put("/categories/:id", withUser(h.Update))

	req := httptest.NewRequest(http.MethodPut, "/categories/1", strings.NewReader(`{"name":"Food Updated","type":"expense"}`))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}
}

func TestCategoryHandler_Delete(t *testing.T) {
	t.Parallel()

	h := NewCategoryHandler(&stubCategoryService{})
	app := fiber.New()
	app.Delete("/categories/:id", withUser(h.Delete))

	req := httptest.NewRequest(http.MethodDelete, "/categories/1", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}
}

func TestCategoryHandler_Delete_InvalidID(t *testing.T) {
	t.Parallel()

	h := NewCategoryHandler(&stubCategoryService{})
	app := fiber.New()
	app.Delete("/categories/:id", withUser(h.Delete))

	req := httptest.NewRequest(http.MethodDelete, "/categories/abc", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", res.StatusCode)
	}
}
