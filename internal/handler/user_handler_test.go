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

type stubUserService struct {
	registerFn      func(context.Context, *models.RegisterRequest) (*models.User, error)
	getProfileFn    func(context.Context, int64) (*models.ProfileResponse, error)
	updateProfileFn func(context.Context, int64, *models.UpdateProfileRequest) (*models.Profile, error)
}

func (s *stubUserService) RegisterUser(ctx context.Context, req *models.RegisterRequest) (*models.User, error) {
	if s.registerFn != nil {
		return s.registerFn(ctx, req)
	}
	return &models.User{ID: 1, Username: req.Username, Email: req.Email}, nil
}

func (s *stubUserService) GetProfile(ctx context.Context, userID int64) (*models.ProfileResponse, error) {
	if s.getProfileFn != nil {
		return s.getProfileFn(ctx, userID)
	}
	return &models.ProfileResponse{User: &models.User{ID: userID}, Profile: &models.Profile{ID: 1, UserID: userID}}, nil
}

func (s *stubUserService) UpdateProfile(ctx context.Context, userID int64, req *models.UpdateProfileRequest) (*models.Profile, error) {
	if s.updateProfileFn != nil {
		return s.updateProfileFn(ctx, userID, req)
	}
	return &models.Profile{ID: 1, UserID: userID}, nil
}

func TestUserHandler_Register(t *testing.T) {
	t.Parallel()

	h := NewUserHandler(&stubUserService{})
	app := fiber.New()
	app.Post("/register", h.Register)

	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(`{"username":"alice","email":"alice@example.com","password":"secret123"}`))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}
}

func TestUserHandler_Register_ValidationFail(t *testing.T) {
	t.Parallel()

	h := NewUserHandler(&stubUserService{})
	app := fiber.New()
	app.Post("/register", h.Register)

	// Missing email
	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(`{"username":"alice","password":"secret123"}`))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", res.StatusCode)
	}
}

func TestUserHandler_GetProfile(t *testing.T) {
	t.Parallel()

	h := NewUserHandler(&stubUserService{})
	app := fiber.New()
	app.Get("/profile", func(c *fiber.Ctx) error {
		c.Locals("user_id", int64(42))
		return h.GetProfile(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/profile", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}
}

func TestUserHandler_GetProfile_Unauthorized(t *testing.T) {
	t.Parallel()

	h := NewUserHandler(&stubUserService{})
	app := fiber.New()
	app.Get("/profile", h.GetProfile)

	req := httptest.NewRequest(http.MethodGet, "/profile", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", res.StatusCode)
	}
}

func TestUserHandler_UpdateProfile(t *testing.T) {
	t.Parallel()

	h := NewUserHandler(&stubUserService{})
	app := fiber.New()
	app.Put("/profile", func(c *fiber.Ctx) error {
		c.Locals("user_id", int64(42))
		return h.UpdateProfile(c)
	})

	req := httptest.NewRequest(http.MethodPut, "/profile", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}
}

func TestUserHandler_UpdateProfile_Unauthorized(t *testing.T) {
	t.Parallel()

	h := NewUserHandler(&stubUserService{})
	app := fiber.New()
	app.Put("/profile", h.UpdateProfile)

	req := httptest.NewRequest(http.MethodPut, "/profile", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", res.StatusCode)
	}
}
