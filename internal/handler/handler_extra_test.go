package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
)

// Missing Unauthorized tests and additional coverage for various handlers

// --- Budget handler ---

func TestBudgetHandler_List_Unauthorized(t *testing.T) {
	t.Parallel()

	h := NewBudgetHandler(&stubBudgetService{})
	app := fiber.New()
	app.Post("/budgets/list", h.List)

	req := httptest.NewRequest(http.MethodPost, "/budgets/list", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", res.StatusCode)
	}
}

func TestBudgetHandler_Delete_Unauthorized(t *testing.T) {
	t.Parallel()

	h := NewBudgetHandler(&stubBudgetService{})
	app := fiber.New()
	app.Delete("/budgets/:id", h.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/budgets/1", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", res.StatusCode)
	}
}

func TestBudgetHandler_Create_ValidationFail(t *testing.T) {
	t.Parallel()

	h := NewBudgetHandler(&stubBudgetService{})
	app := fiber.New()
	app.Post("/budgets", withUser(h.Create))

	req := httptest.NewRequest(http.MethodPost, "/budgets", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", res.StatusCode)
	}
}

// --- Category handler ---

func TestCategoryHandler_GetAll_Unauthorized(t *testing.T) {
	t.Parallel()

	h := NewCategoryHandler(&stubCategoryService{})
	app := fiber.New()
	app.Get("/categories", h.GetAll)

	req := httptest.NewRequest(http.MethodGet, "/categories", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", res.StatusCode)
	}
}

func TestCategoryHandler_Update_ValidationFail(t *testing.T) {
	t.Parallel()

	h := NewCategoryHandler(&stubCategoryService{})
	app := fiber.New()
	app.Put("/categories/:id", withUser(h.Update))

	req := httptest.NewRequest(http.MethodPut, "/categories/1", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", res.StatusCode)
	}
}

func TestCategoryHandler_Update_InvalidID(t *testing.T) {
	t.Parallel()

	h := NewCategoryHandler(&stubCategoryService{})
	app := fiber.New()
	app.Put("/categories/:id", withUser(h.Update))

	req := httptest.NewRequest(http.MethodPut, "/categories/abc", strings.NewReader(`{"name":"Food","type":"expense"}`))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", res.StatusCode)
	}
}

func TestCategoryHandler_Delete_Unauthorized(t *testing.T) {
	t.Parallel()

	h := NewCategoryHandler(&stubCategoryService{})
	app := fiber.New()
	app.Delete("/categories/:id", h.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/categories/1", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", res.StatusCode)
	}
}

// --- Device handler ---

func TestDeviceHandler_List_Unauthorized(t *testing.T) {
	t.Parallel()

	h := NewDeviceHandler(&stubDeviceService{})
	app := fiber.New()
	app.Post("/devices/list", h.List)

	req := httptest.NewRequest(http.MethodPost, "/devices/list", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", res.StatusCode)
	}
}

func TestDeviceHandler_Delete_Unauthorized(t *testing.T) {
	t.Parallel()

	h := NewDeviceHandler(&stubDeviceService{})
	app := fiber.New()
	app.Post("/devices/:id", h.Delete)

	req := httptest.NewRequest(http.MethodPost, "/devices/1", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", res.StatusCode)
	}
}

// --- Fund handler ---

func TestFundHandler_GetAll_Unauthorized(t *testing.T) {
	t.Parallel()

	h := NewFundHandler(&stubFundService{})
	app := fiber.New()
	app.Get("/funds", h.GetAll)

	req := httptest.NewRequest(http.MethodGet, "/funds", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", res.StatusCode)
	}
}

func TestFundHandler_Delete_Unauthorized(t *testing.T) {
	t.Parallel()

	h := NewFundHandler(&stubFundService{})
	app := fiber.New()
	app.Delete("/funds/:id", h.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/funds/1", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", res.StatusCode)
	}
}

func TestFundHandler_Deposit_Unauthorized(t *testing.T) {
	t.Parallel()

	h := NewFundHandler(&stubFundService{})
	app := fiber.New()
	app.Post("/funds/:id/deposit", h.Deposit)

	req := httptest.NewRequest(http.MethodPost, "/funds/1/deposit", strings.NewReader(`{"amount":100}`))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", res.StatusCode)
	}
}

func TestFundHandler_Withdraw_Unauthorized(t *testing.T) {
	t.Parallel()

	h := NewFundHandler(&stubFundService{})
	app := fiber.New()
	app.Post("/funds/:id/withdraw", h.Withdraw)

	req := httptest.NewRequest(http.MethodPost, "/funds/1/withdraw", strings.NewReader(`{"amount":50}`))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", res.StatusCode)
	}
}

func TestFundHandler_Update_Unauthorized(t *testing.T) {
	t.Parallel()

	h := NewFundHandler(&stubFundService{})
	app := fiber.New()
	app.Put("/funds/:id", h.Update)

	req := httptest.NewRequest(http.MethodPut, "/funds/1", strings.NewReader(`{"name":"Test"}`))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", res.StatusCode)
	}
}

// --- Notification handler ---

func TestNotificationHandler_UnreadCount_Unauthorized(t *testing.T) {
	t.Parallel()

	h := NewNotificationHandler(&stubNotificationService{})
	app := fiber.New()
	app.Post("/notifications/unread-count", h.UnreadCount)

	req := httptest.NewRequest(http.MethodPost, "/notifications/unread-count", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", res.StatusCode)
	}
}

func TestNotificationHandler_MarkRead_Unauthorized(t *testing.T) {
	t.Parallel()

	h := NewNotificationHandler(&stubNotificationService{})
	app := fiber.New()
	app.Post("/notifications/mark-read", h.MarkRead)

	req := httptest.NewRequest(http.MethodPost, "/notifications/mark-read", strings.NewReader(`{"ids":[1,2]}`))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", res.StatusCode)
	}
}

func TestNotificationHandler_Delete_Unauthorized(t *testing.T) {
	t.Parallel()

	h := NewNotificationHandler(&stubNotificationService{})
	app := fiber.New()
	app.Post("/notifications/delete/:id", h.Delete)

	req := httptest.NewRequest(http.MethodPost, "/notifications/delete/1", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", res.StatusCode)
	}
}

// --- Savings goal handler ---

func TestSavingsGoalHandler_List_Unauthorized(t *testing.T) {
	t.Parallel()

	h := NewSavingsGoalHandler(&stubSavingsGoalService{})
	app := fiber.New()
	app.Get("/goals", h.List)

	req := httptest.NewRequest(http.MethodGet, "/goals", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", res.StatusCode)
	}
}

func TestSavingsGoalHandler_Withdraw_Unauthorized(t *testing.T) {
	t.Parallel()

	h := NewSavingsGoalHandler(&stubSavingsGoalService{})
	app := fiber.New()
	app.Post("/goals/withdraw", h.Withdraw)

	req := httptest.NewRequest(http.MethodPost, "/goals/withdraw", strings.NewReader(`{"goal_id":"1","fund_id":"2","amount":50}`))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", res.StatusCode)
	}
}

// --- Auth handler extra ---


// --- Report handler extra ---

func TestReportHandler_GetCategoryTrend(t *testing.T) {
	t.Parallel()

	svc := &stubReportService{}
	h := NewReportHandler(svc)
	app := fiber.New()
	app.Post("/reports/category-trend", func(c *fiber.Ctx) error {
		c.Locals("user_id", int64(42))
		return h.GetCategoryTrend(c)
	})

	req := httptest.NewRequest(http.MethodPost, "/reports/category-trend", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}
}

func TestReportHandler_GetCategoryTrend_Unauthorized(t *testing.T) {
	t.Parallel()

	h := NewReportHandler(&stubReportService{})
	app := fiber.New()
	app.Post("/reports/category-trend", h.GetCategoryTrend)

	req := httptest.NewRequest(http.MethodPost, "/reports/category-trend", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", res.StatusCode)
	}
}
