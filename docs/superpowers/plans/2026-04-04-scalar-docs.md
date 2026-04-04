# Scalar Docs Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Serve Scalar API documentation from the existing Fiber app using `internal/docs/openapi.yaml` as the source of truth.

**Architecture:** Add a small `internal/docs` package that embeds the OpenAPI YAML and serves two public routes: `/openapi.yaml` for the raw spec and `/docs` for a Scalar HTML shell loaded from the official CDN. Wire the docs routes into the existing server bootstrap so the API and docs live in the same process.

**Tech Stack:** Go 1.25, Fiber v2, `embed`, `net/http/httptest`

---

### Task 1: Add docs-serving package

**Files:**
- Modify: `internal/docs/`
- Test: `internal/docs/routes_test.go`

- [ ] **Step 1: Write the failing route test**

```go
func TestRegisterRoutesServesDocsAndSpec(t *testing.T) {
	app := fiber.New()
	RegisterRoutes(app)
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/docs`
Expected: FAIL because `RegisterRoutes` does not exist yet

- [ ] **Step 3: Write minimal implementation**

```go
//go:embed openapi.yaml
var openAPISpec []byte

func RegisterRoutes(app *fiber.App) {
	app.Get("/openapi.yaml", serveOpenAPI)
	app.Get("/docs", serveScalarReference)
}
```

- [ ] **Step 4: Expand the test to assert response shape**

```go
req := httptest.NewRequest(http.MethodGet, "/docs", nil)
resp, _ := app.Test(req)
require.Equal(t, http.StatusOK, resp.StatusCode)
```

- [ ] **Step 5: Run the package tests**

Run: `go test ./internal/docs`
Expected: PASS

### Task 2: Mount docs routes in the server

**Files:**
- Modify: `internal/server/server.go`

- [ ] **Step 1: Import the docs package**

```go
import "github.com/maynguyen24/sever/internal/docs"
```

- [ ] **Step 2: Register docs routes before API routes**

```go
docs.RegisterRoutes(app)
router.SetupRoutes(app, s.cfg)
```

- [ ] **Step 3: Run a targeted build**

Run: `go test ./internal/server ./internal/docs`
Expected: PASS

### Task 3: Verify app-level integration

**Files:**
- Verify: `internal/docs/openapi.yaml`

- [ ] **Step 1: Run the focused package tests**

Run: `go test ./internal/docs ./internal/server`
Expected: PASS

- [ ] **Step 2: Run the full test suite if it exists**

Run: `go test ./...`
Expected: PASS, or document unrelated failures if present

- [ ] **Step 3: Review the diff**

Run: `git diff -- internal/docs internal/server/server.go docs/superpowers/plans/2026-04-04-scalar-docs.md`
Expected: only the Scalar Docs implementation and plan changes appear
