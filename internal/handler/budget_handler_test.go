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

type stubBudgetService struct {
	createFn    func(context.Context, int64, *models.CreateBudgetRequest) (*models.Budget, error)
	getListFn   func(context.Context, int64) ([]*models.BudgetResponse, error)
	getDetailFn func(context.Context, int64, int64) (*models.BudgetResponse, error)
	updateFn    func(context.Context, int64, int64, *models.UpdateBudgetRequest) (*models.Budget, error)
	deleteFn    func(context.Context, int64, int64) error
}

func (s *stubBudgetService) Create(ctx context.Context, userID int64, req *models.CreateBudgetRequest) (*models.Budget, error) {
	if s.createFn != nil {
		return s.createFn(ctx, userID, req)
	}
	return &models.Budget{ID: 1, UserID: userID, Amount: req.Amount, Period: req.Period}, nil
}

func (s *stubBudgetService) GetList(ctx context.Context, userID int64) ([]*models.BudgetResponse, error) {
	if s.getListFn != nil {
		return s.getListFn(ctx, userID)
	}
	return []*models.BudgetResponse{{Budget: models.Budget{ID: 1}}}, nil
}

func (s *stubBudgetService) GetDetail(ctx context.Context, id, userID int64) (*models.BudgetResponse, error) {
	if s.getDetailFn != nil {
		return s.getDetailFn(ctx, id, userID)
	}
	return &models.BudgetResponse{Budget: models.Budget{ID: id}}, nil
}

func (s *stubBudgetService) Update(ctx context.Context, id, userID int64, req *models.UpdateBudgetRequest) (*models.Budget, error) {
	if s.updateFn != nil {
		return s.updateFn(ctx, id, userID, req)
	}
	return &models.Budget{ID: id}, nil
}

func (s *stubBudgetService) Delete(ctx context.Context, id, userID int64) error {
	if s.deleteFn != nil {
		return s.deleteFn(ctx, id, userID)
	}
	return nil
}

func TestBudgetHandler_Create(t *testing.T) {
	t.Parallel()

	h := NewBudgetHandler(&stubBudgetService{})
	app := fiber.New()
	app.Post("/budgets", withUser(h.Create))

	req := httptest.NewRequest(http.MethodPost, "/budgets", strings.NewReader(`{"amount":1000,"period":"monthly"}`))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}
}

func TestBudgetHandler_Create_Unauthorized(t *testing.T) {
	t.Parallel()

	h := NewBudgetHandler(&stubBudgetService{})
	app := fiber.New()
	app.Post("/budgets", h.Create)

	req := httptest.NewRequest(http.MethodPost, "/budgets", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", res.StatusCode)
	}
}

func TestBudgetHandler_List(t *testing.T) {
	t.Parallel()

	h := NewBudgetHandler(&stubBudgetService{})
	app := fiber.New()
	app.Post("/budgets/list", withUser(h.List))

	req := httptest.NewRequest(http.MethodPost, "/budgets/list", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}
}

func TestBudgetHandler_GetDetail(t *testing.T) {
	t.Parallel()

	h := NewBudgetHandler(&stubBudgetService{})
	app := fiber.New()
	app.Post("/budgets/:id", withUser(h.GetDetail))

	req := httptest.NewRequest(http.MethodPost, "/budgets/1", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}
}

func TestBudgetHandler_GetDetail_InvalidID(t *testing.T) {
	t.Parallel()

	h := NewBudgetHandler(&stubBudgetService{})
	app := fiber.New()
	app.Post("/budgets/:id", withUser(h.GetDetail))

	req := httptest.NewRequest(http.MethodPost, "/budgets/abc", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", res.StatusCode)
	}
}

func TestBudgetHandler_Update(t *testing.T) {
	t.Parallel()

	h := NewBudgetHandler(&stubBudgetService{})
	app := fiber.New()
	app.Post("/budgets/update/:id", withUser(h.Update))

	req := httptest.NewRequest(http.MethodPost, "/budgets/update/1", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}
}

func TestBudgetHandler_Update_InvalidID(t *testing.T) {
	t.Parallel()

	h := NewBudgetHandler(&stubBudgetService{})
	app := fiber.New()
	app.Post("/budgets/update/:id", withUser(h.Update))

	req := httptest.NewRequest(http.MethodPost, "/budgets/update/abc", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", res.StatusCode)
	}
}

func TestBudgetHandler_Delete(t *testing.T) {
	t.Parallel()

	h := NewBudgetHandler(&stubBudgetService{})
	app := fiber.New()
	app.Delete("/budgets/:id", withUser(h.Delete))

	req := httptest.NewRequest(http.MethodDelete, "/budgets/1", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}
}

func TestBudgetHandler_Delete_InvalidID(t *testing.T) {
	t.Parallel()

	h := NewBudgetHandler(&stubBudgetService{})
	app := fiber.New()
	app.Delete("/budgets/:id", withUser(h.Delete))

	req := httptest.NewRequest(http.MethodDelete, "/budgets/abc", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", res.StatusCode)
	}
}
