package router

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/maynguyen24/sever/configs"
	"github.com/maynguyen24/sever/pkg/logger"
	"go.uber.org/zap"
)

func testRouterConfig() *configs.Config {
	return &configs.Config{
		JWTSecret:               "jwt-secret",
		JWTRefreshSecret:        "refresh-secret",
		GoogleRedirectURL:       "http://localhost/callback",
		WebAuthnRPID:            "localhost",
		WebAuthnRPName:          "Test",
		WebAuthnOrigin:          "http://localhost:3000",
		RedisAddr:               "127.0.0.1:0",
		RedisPassword:           "",
		FirebaseCredentialsPath: "",
	}
}

func TestSetupRoutesRegistersPublicAndProtectedEndpoints(t *testing.T) {
	t.Parallel()

	logger.Log = zap.NewNop()

	app := fiber.New()
	SetupRoutes(app, testRouterConfig(), nil)

	tests := []struct {
		name       string
		method     string
		target     string
		wantStatus int
	}{
		{name: "public login route registered", method: fiber.MethodPost, target: "/api/v1/auth/login", wantStatus: fiber.StatusBadRequest},
		{name: "google url route registered", method: fiber.MethodPost, target: "/api/v1/auth/google/url", wantStatus: fiber.StatusOK},
		{name: "protected funds route registered", method: fiber.MethodPost, target: "/api/v1/funds/list", wantStatus: fiber.StatusUnauthorized},
		{name: "protected goals route registered", method: fiber.MethodPost, target: "/api/v1/goals/list", wantStatus: fiber.StatusUnauthorized},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.target, nil)
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("app.Test(%s): %v", tt.target, err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.wantStatus {
				t.Fatalf("expected status %d for %s, got %d", tt.wantStatus, tt.target, resp.StatusCode)
			}
		})
	}
}
