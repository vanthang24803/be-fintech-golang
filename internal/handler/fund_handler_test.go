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

type stubFundService struct {
	createFn   func(context.Context, int64, *models.CreateFundRequest) (*models.Fund, error)
	getAllFn    func(context.Context, int64) ([]*models.Fund, error)
	getByIDFn  func(context.Context, int64, int64) (*models.Fund, error)
	updateFn   func(context.Context, int64, int64, *models.UpdateFundRequest) (*models.Fund, error)
	deleteFn   func(context.Context, int64, int64) error
	depositFn  func(context.Context, int64, int64, *models.FundTransactionRequest) (*models.Fund, error)
	withdrawFn func(context.Context, int64, int64, *models.FundTransactionRequest) (*models.Fund, error)
}

func (s *stubFundService) Create(ctx context.Context, userID int64, req *models.CreateFundRequest) (*models.Fund, error) {
	if s.createFn != nil {
		return s.createFn(ctx, userID, req)
	}
	return &models.Fund{ID: 1, UserID: userID, Name: req.Name}, nil
}

func (s *stubFundService) GetAll(ctx context.Context, userID int64) ([]*models.Fund, error) {
	if s.getAllFn != nil {
		return s.getAllFn(ctx, userID)
	}
	return []*models.Fund{{ID: 1, UserID: userID}}, nil
}

func (s *stubFundService) GetByID(ctx context.Context, id, userID int64) (*models.Fund, error) {
	if s.getByIDFn != nil {
		return s.getByIDFn(ctx, id, userID)
	}
	return &models.Fund{ID: id, UserID: userID}, nil
}

func (s *stubFundService) Update(ctx context.Context, id, userID int64, req *models.UpdateFundRequest) (*models.Fund, error) {
	if s.updateFn != nil {
		return s.updateFn(ctx, id, userID, req)
	}
	return &models.Fund{ID: id, Name: req.Name}, nil
}

func (s *stubFundService) Delete(ctx context.Context, id, userID int64) error {
	if s.deleteFn != nil {
		return s.deleteFn(ctx, id, userID)
	}
	return nil
}

func (s *stubFundService) Deposit(ctx context.Context, id, userID int64, req *models.FundTransactionRequest) (*models.Fund, error) {
	if s.depositFn != nil {
		return s.depositFn(ctx, id, userID, req)
	}
	return &models.Fund{ID: id, Balance: req.Amount}, nil
}

func (s *stubFundService) Withdraw(ctx context.Context, id, userID int64, req *models.FundTransactionRequest) (*models.Fund, error) {
	if s.withdrawFn != nil {
		return s.withdrawFn(ctx, id, userID, req)
	}
	return &models.Fund{ID: id, Balance: 0}, nil
}

func TestFundHandler_Create(t *testing.T) {
	t.Parallel()

	h := NewFundHandler(&stubFundService{})
	app := fiber.New()
	app.Post("/funds", withUser(h.Create))

	req := httptest.NewRequest(http.MethodPost, "/funds", strings.NewReader(`{"name":"Savings"}`))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}
}

func TestFundHandler_Create_Unauthorized(t *testing.T) {
	t.Parallel()

	h := NewFundHandler(&stubFundService{})
	app := fiber.New()
	app.Post("/funds", h.Create)

	req := httptest.NewRequest(http.MethodPost, "/funds", strings.NewReader(`{"name":"Savings"}`))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", res.StatusCode)
	}
}

func TestFundHandler_GetAll(t *testing.T) {
	t.Parallel()

	h := NewFundHandler(&stubFundService{})
	app := fiber.New()
	app.Get("/funds", withUser(h.GetAll))

	req := httptest.NewRequest(http.MethodGet, "/funds", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}
}

func TestFundHandler_GetByID(t *testing.T) {
	t.Parallel()

	h := NewFundHandler(&stubFundService{})
	app := fiber.New()
	app.Get("/funds/:id", withUser(h.GetByID))

	req := httptest.NewRequest(http.MethodGet, "/funds/1", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}
}

func TestFundHandler_GetByID_InvalidID(t *testing.T) {
	t.Parallel()

	h := NewFundHandler(&stubFundService{})
	app := fiber.New()
	app.Get("/funds/:id", withUser(h.GetByID))

	req := httptest.NewRequest(http.MethodGet, "/funds/abc", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", res.StatusCode)
	}
}

func TestFundHandler_Update(t *testing.T) {
	t.Parallel()

	h := NewFundHandler(&stubFundService{})
	app := fiber.New()
	app.Put("/funds/:id", withUser(h.Update))

	req := httptest.NewRequest(http.MethodPut, "/funds/1", strings.NewReader(`{"name":"Updated Fund"}`))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}
}

func TestFundHandler_Update_InvalidID(t *testing.T) {
	t.Parallel()

	h := NewFundHandler(&stubFundService{})
	app := fiber.New()
	app.Put("/funds/:id", withUser(h.Update))

	req := httptest.NewRequest(http.MethodPut, "/funds/abc", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", res.StatusCode)
	}
}

func TestFundHandler_Delete(t *testing.T) {
	t.Parallel()

	h := NewFundHandler(&stubFundService{})
	app := fiber.New()
	app.Delete("/funds/:id", withUser(h.Delete))

	req := httptest.NewRequest(http.MethodDelete, "/funds/1", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}
}

func TestFundHandler_Deposit(t *testing.T) {
	t.Parallel()

	h := NewFundHandler(&stubFundService{})
	app := fiber.New()
	app.Post("/funds/:id/deposit", withUser(h.Deposit))

	req := httptest.NewRequest(http.MethodPost, "/funds/1/deposit", strings.NewReader(`{"amount":500}`))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}
}

func TestFundHandler_Deposit_ValidationFail(t *testing.T) {
	t.Parallel()

	h := NewFundHandler(&stubFundService{})
	app := fiber.New()
	app.Post("/funds/:id/deposit", withUser(h.Deposit))

	// Amount = 0, should fail validation (gt=0)
	req := httptest.NewRequest(http.MethodPost, "/funds/1/deposit", strings.NewReader(`{"amount":0}`))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", res.StatusCode)
	}
}

func TestFundHandler_Withdraw(t *testing.T) {
	t.Parallel()

	h := NewFundHandler(&stubFundService{})
	app := fiber.New()
	app.Post("/funds/:id/withdraw", withUser(h.Withdraw))

	req := httptest.NewRequest(http.MethodPost, "/funds/1/withdraw", strings.NewReader(`{"amount":100}`))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}
}

func TestFundHandler_Withdraw_InvalidID(t *testing.T) {
	t.Parallel()

	h := NewFundHandler(&stubFundService{})
	app := fiber.New()
	app.Post("/funds/:id/withdraw", withUser(h.Withdraw))

	req := httptest.NewRequest(http.MethodPost, "/funds/abc/withdraw", strings.NewReader(`{"amount":100}`))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", res.StatusCode)
	}
}
