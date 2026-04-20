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

type stubSavingsGoalService struct {
	createFn     func(context.Context, int64, *models.CreateGoalRequest) (*models.SavingsGoal, error)
	listFn       func(context.Context, int64) ([]models.SavingsGoal, error)
	getDetailFn  func(context.Context, int64, int64) (*models.GoalResponse, error)
	contributeFn func(context.Context, int64, *models.GoalContributeRequest) (*models.SavingsGoal, error)
	withdrawFn   func(context.Context, int64, *models.GoalWithdrawRequest) (*models.SavingsGoal, error)
}

func (s *stubSavingsGoalService) Create(ctx context.Context, userID int64, req *models.CreateGoalRequest) (*models.SavingsGoal, error) {
	if s.createFn != nil {
		return s.createFn(ctx, userID, req)
	}
	return &models.SavingsGoal{ID: 1, UserID: userID, Name: req.Name, TargetAmount: req.TargetAmount}, nil
}

func (s *stubSavingsGoalService) List(ctx context.Context, userID int64) ([]models.SavingsGoal, error) {
	if s.listFn != nil {
		return s.listFn(ctx, userID)
	}
	return []models.SavingsGoal{{ID: 1, UserID: userID}}, nil
}

func (s *stubSavingsGoalService) GetDetail(ctx context.Context, id int64, userID int64) (*models.GoalResponse, error) {
	if s.getDetailFn != nil {
		return s.getDetailFn(ctx, id, userID)
	}
	return &models.GoalResponse{Goal: &models.SavingsGoal{ID: id}}, nil
}

func (s *stubSavingsGoalService) Contribute(ctx context.Context, userID int64, req *models.GoalContributeRequest) (*models.SavingsGoal, error) {
	if s.contributeFn != nil {
		return s.contributeFn(ctx, userID, req)
	}
	return &models.SavingsGoal{ID: req.GoalID}, nil
}

func (s *stubSavingsGoalService) Withdraw(ctx context.Context, userID int64, req *models.GoalWithdrawRequest) (*models.SavingsGoal, error) {
	if s.withdrawFn != nil {
		return s.withdrawFn(ctx, userID, req)
	}
	return &models.SavingsGoal{ID: req.GoalID}, nil
}

func TestSavingsGoalHandler_Create(t *testing.T) {
	t.Parallel()

	h := NewSavingsGoalHandler(&stubSavingsGoalService{})
	app := fiber.New()
	app.Post("/goals", withUser(h.Create))

	req := httptest.NewRequest(http.MethodPost, "/goals", strings.NewReader(`{"name":"New Car","target_amount":5000}`))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}
}

func TestSavingsGoalHandler_Create_Unauthorized(t *testing.T) {
	t.Parallel()

	h := NewSavingsGoalHandler(&stubSavingsGoalService{})
	app := fiber.New()
	app.Post("/goals", h.Create)

	req := httptest.NewRequest(http.MethodPost, "/goals", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", res.StatusCode)
	}
}

func TestSavingsGoalHandler_Create_ValidationFail(t *testing.T) {
	t.Parallel()

	h := NewSavingsGoalHandler(&stubSavingsGoalService{})
	app := fiber.New()
	app.Post("/goals", withUser(h.Create))

	// Missing required name
	req := httptest.NewRequest(http.MethodPost, "/goals", strings.NewReader(`{"target_amount":5000}`))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", res.StatusCode)
	}
}

func TestSavingsGoalHandler_List(t *testing.T) {
	t.Parallel()

	h := NewSavingsGoalHandler(&stubSavingsGoalService{})
	app := fiber.New()
	app.Get("/goals", withUser(h.List))

	req := httptest.NewRequest(http.MethodGet, "/goals", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}
}

func TestSavingsGoalHandler_GetDetail(t *testing.T) {
	t.Parallel()

	h := NewSavingsGoalHandler(&stubSavingsGoalService{})
	app := fiber.New()
	app.Get("/goals/:id", withUser(h.GetDetail))

	req := httptest.NewRequest(http.MethodGet, "/goals/1", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}
}

func TestSavingsGoalHandler_GetDetail_InvalidID(t *testing.T) {
	t.Parallel()

	h := NewSavingsGoalHandler(&stubSavingsGoalService{})
	app := fiber.New()
	app.Get("/goals/:id", withUser(h.GetDetail))

	req := httptest.NewRequest(http.MethodGet, "/goals/abc", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", res.StatusCode)
	}
}

func TestSavingsGoalHandler_Contribute(t *testing.T) {
	t.Parallel()

	h := NewSavingsGoalHandler(&stubSavingsGoalService{})
	app := fiber.New()
	app.Post("/goals/contribute", withUser(h.Contribute))

	req := httptest.NewRequest(http.MethodPost, "/goals/contribute", strings.NewReader(`{"goal_id":"1","fund_id":"2","amount":100}`))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}
}

func TestSavingsGoalHandler_Contribute_ValidationFail(t *testing.T) {
	t.Parallel()

	h := NewSavingsGoalHandler(&stubSavingsGoalService{})
	app := fiber.New()
	app.Post("/goals/contribute", withUser(h.Contribute))

	// Missing required fields
	req := httptest.NewRequest(http.MethodPost, "/goals/contribute", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", res.StatusCode)
	}
}

func TestSavingsGoalHandler_Withdraw(t *testing.T) {
	t.Parallel()

	h := NewSavingsGoalHandler(&stubSavingsGoalService{})
	app := fiber.New()
	app.Post("/goals/withdraw", withUser(h.Withdraw))

	req := httptest.NewRequest(http.MethodPost, "/goals/withdraw", strings.NewReader(`{"goal_id":"1","fund_id":"2","amount":50}`))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}
}

func TestSavingsGoalHandler_Withdraw_ValidationFail(t *testing.T) {
	t.Parallel()

	h := NewSavingsGoalHandler(&stubSavingsGoalService{})
	app := fiber.New()
	app.Post("/goals/withdraw", withUser(h.Withdraw))

	// Missing required fields
	req := httptest.NewRequest(http.MethodPost, "/goals/withdraw", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", res.StatusCode)
	}
}
