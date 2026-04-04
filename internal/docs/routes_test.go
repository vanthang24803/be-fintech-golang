package docs

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestRegisterRoutesServesOpenAPISpec(t *testing.T) {
	app := fiber.New()
	RegisterRoutes(app)

	req := httptest.NewRequest(http.MethodGet, "/openapi.yaml", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	if got := resp.Header.Get("Content-Type"); !strings.Contains(got, "application/yaml") {
		t.Fatalf("content-type = %q, want application/yaml", got)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("io.ReadAll() error = %v", err)
	}

	if !strings.Contains(string(body), "openapi: 3.0.3") {
		t.Fatalf("body missing embedded OpenAPI version")
	}
}

func TestRegisterRoutesServesScalarReference(t *testing.T) {
	app := fiber.New()
	RegisterRoutes(app)

	req := httptest.NewRequest(http.MethodGet, "/docs", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	if got := resp.Header.Get("Content-Type"); !strings.Contains(got, "text/html") {
		t.Fatalf("content-type = %q, want text/html", got)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("io.ReadAll() error = %v", err)
	}

	content := string(body)
	if !strings.Contains(content, "Scalar.createApiReference") {
		t.Fatalf("body missing Scalar bootstrap")
	}
	if !strings.Contains(content, "\"url\":\"/openapi.yaml\"") {
		t.Fatalf("body missing local OpenAPI URL")
	}
}

func TestEmbeddedOpenAPISpecHasSingleComponentsBlock(t *testing.T) {
	content := string(openAPISpec)

	if count := strings.Count(content, "\ncomponents:"); count != 1 {
		t.Fatalf("components block count = %d, want 1", count)
	}

	if strings.Contains(content, "/auth/refresh#placeholder: {}") {
		t.Fatalf("embedded spec contains placeholder path hack")
	}
}
