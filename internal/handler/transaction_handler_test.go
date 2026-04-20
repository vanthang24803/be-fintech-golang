package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/maynguyen24/sever/internal/models"
)

type stubTransactionService struct {
	createFn  func(context.Context, int64, *models.CreateTransactionRequest) (*models.Transaction, error)
	getAllFn   func(context.Context, int64, map[string]string) ([]*models.TransactionDetail, error)
	getByIDFn func(context.Context, int64, int64) (*models.TransactionDetail, error)
	updateFn  func(context.Context, int64, int64, *models.UpdateTransactionRequest) (*models.Transaction, error)
	deleteFn  func(context.Context, int64, int64) error
}

func (s *stubTransactionService) Create(ctx context.Context, userID int64, req *models.CreateTransactionRequest) (*models.Transaction, error) {
	if s.createFn != nil {
		return s.createFn(ctx, userID, req)
	}
	return &models.Transaction{ID: 1, UserID: userID, Amount: req.Amount, Type: req.Type}, nil
}

func (s *stubTransactionService) GetAll(ctx context.Context, userID int64, query map[string]string) ([]*models.TransactionDetail, error) {
	if s.getAllFn != nil {
		return s.getAllFn(ctx, userID, query)
	}
	return []*models.TransactionDetail{{Transaction: models.Transaction{ID: 1}}}, nil
}

func (s *stubTransactionService) GetByID(ctx context.Context, id, userID int64) (*models.TransactionDetail, error) {
	if s.getByIDFn != nil {
		return s.getByIDFn(ctx, id, userID)
	}
	return &models.TransactionDetail{Transaction: models.Transaction{ID: id}}, nil
}

func (s *stubTransactionService) Update(ctx context.Context, id, userID int64, req *models.UpdateTransactionRequest) (*models.Transaction, error) {
	if s.updateFn != nil {
		return s.updateFn(ctx, id, userID, req)
	}
	return &models.Transaction{ID: id, Amount: req.Amount, Type: req.Type}, nil
}

func (s *stubTransactionService) Delete(ctx context.Context, id, userID int64) error {
	if s.deleteFn != nil {
		return s.deleteFn(ctx, id, userID)
	}
	return nil
}

func TestTransactionHandler_Create(t *testing.T) {
	t.Parallel()

	h := NewTransactionHandler(&stubTransactionService{})
	app := fiber.New()
	app.Post("/transactions", withUser(h.Create))

	body := `{"source_payment_id":"1","amount":100.5,"type":"expense","transaction_date":"` + time.Now().Format(time.RFC3339) + `"}`
	req := httptest.NewRequest(http.MethodPost, "/transactions", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}
}

func TestTransactionHandler_Create_Unauthorized(t *testing.T) {
	t.Parallel()

	h := NewTransactionHandler(&stubTransactionService{})
	app := fiber.New()
	app.Post("/transactions", h.Create)

	req := httptest.NewRequest(http.MethodPost, "/transactions", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", res.StatusCode)
	}
}

func TestTransactionHandler_GetAll(t *testing.T) {
	t.Parallel()

	h := NewTransactionHandler(&stubTransactionService{})
	app := fiber.New()
	app.Post("/transactions/list", withUser(h.GetAll))

	req := httptest.NewRequest(http.MethodPost, "/transactions/list", strings.NewReader(`{"type":"expense"}`))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}
}

func TestTransactionHandler_GetByID(t *testing.T) {
	t.Parallel()

	h := NewTransactionHandler(&stubTransactionService{})
	app := fiber.New()
	app.Get("/transactions/:id", withUser(h.GetByID))

	req := httptest.NewRequest(http.MethodGet, "/transactions/1", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}
}

func TestTransactionHandler_GetByID_InvalidID(t *testing.T) {
	t.Parallel()

	h := NewTransactionHandler(&stubTransactionService{})
	app := fiber.New()
	app.Get("/transactions/:id", withUser(h.GetByID))

	req := httptest.NewRequest(http.MethodGet, "/transactions/abc", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", res.StatusCode)
	}
}

func TestTransactionHandler_Update(t *testing.T) {
	t.Parallel()

	h := NewTransactionHandler(&stubTransactionService{})
	app := fiber.New()
	app.Put("/transactions/:id", withUser(h.Update))

	body := `{"source_payment_id":"1","amount":200,"type":"income","transaction_date":"` + time.Now().Format(time.RFC3339) + `"}`
	req := httptest.NewRequest(http.MethodPut, "/transactions/1", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}
}

func TestTransactionHandler_Update_InvalidID(t *testing.T) {
	t.Parallel()

	h := NewTransactionHandler(&stubTransactionService{})
	app := fiber.New()
	app.Put("/transactions/:id", withUser(h.Update))

	req := httptest.NewRequest(http.MethodPut, "/transactions/abc", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", res.StatusCode)
	}
}

func TestTransactionHandler_Delete(t *testing.T) {
	t.Parallel()

	h := NewTransactionHandler(&stubTransactionService{})
	app := fiber.New()
	app.Delete("/transactions/:id", withUser(h.Delete))

	req := httptest.NewRequest(http.MethodDelete, "/transactions/1", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}
}

func TestTransactionHandler_Delete_InvalidID(t *testing.T) {
	t.Parallel()

	h := NewTransactionHandler(&stubTransactionService{})
	app := fiber.New()
	app.Delete("/transactions/:id", withUser(h.Delete))

	req := httptest.NewRequest(http.MethodDelete, "/transactions/abc", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", res.StatusCode)
	}
}
