# Response Format Guide

This document describes how to use the project's shared `response` package and how the global error handler intercepts errors across all Fiber routes.

---

## 1. Success Response

Use the `response.Success` helper whenever an operation completes successfully.

```go
import "github.com/maynguyen24/sever/pkg/response"

func MyHandler(c *fiber.Ctx) error {
    data := map[string]string{"foo": "bar"}
    // 2000 is a project-level business code (2001 = created)
    return response.Success(c, 2000, "Operation successful", data)
}
```

**Output:**
```json
{
  "code": 2000,
  "message": "Operation successful",
  "data": { "foo": "bar" }
}
```

---

## 2. Error Response (Global Error Interceptor)

Any `error` returned from a Fiber handler is automatically caught and formatted
by the global `ErrorHandler` middleware registered in `server.go`.

```go
func MyHandler(c *fiber.Ctx) error {
    result, err := someService.DoSomething()
    if err != nil {
        // Plain error → wrapped as 500 Internal Server Error
        return err

        // Fiber error → respects the provided HTTP status code
        // return fiber.NewError(fiber.StatusBadRequest, "Invalid parameter")
    }
    return response.Success(c, 2000, "Done", result)
}
```

**Output (example — 400 Bad Request):**
```json
{
  "code": 4000,
  "message": "Invalid parameter",
  "error": "Technical details or the original error message (for logs/debug)"
}
```

---

## 3. Handler Pattern Summary

```
Request → Handler → Service → Repository → Database
                ↓ error at any layer
         GlobalErrorHandler → JSON error response
```

| Layer | Responsibility |
|---|---|
| **Handler** | Parse input, call service, return `response.Success` |
| **Service** | Business logic, use `fiber.NewError` for client errors |
| **Repository** | Raw DB operations, return plain `error` |
| **ErrorHandler** | Intercept all errors, format into unified JSON |
