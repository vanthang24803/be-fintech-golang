# Logging & API Standards

This document defines the logging and API response standards for the project.
All team members must follow these conventions to maintain consistency across modules.

---

## 1. System Logging

Use `go.uber.org/zap` (or `slog`) with structured logging.

| Environment | Format | Notes |
|---|---|---|
| **Development** | Human-readable, colorized | Easy to read in terminal |
| **Production** | JSON | Includes `timestamp`, `caller`, `level`, `message` |

### Example — Development log
```
2026-03-29T14:00:00Z INFO  server started {"port": 8080}
```

### Example — Production log (JSON)
```json
{
  "level": "info",
  "ts": "2026-03-29T14:00:00.000Z",
  "caller": "server/server.go:42",
  "msg": "server started",
  "port": 8080
}
```

---

## 2. API Response Format

All responses returned to the client follow a **unified JSON envelope**.

### Success

```json
{
  "code": 2000,
  "message": "Operation successful",
  "data": { "..." }
}
```

### Error

```json
{
  "code": 4000,
  "message": "User-facing error message",
  "error": "Technical details or original error message (for debug/logging)"
}
```

### Business Code Convention

| Range | Meaning |
|---|---|
| `2000–2099` | Success (2001 = Created) |
| `4000–4099` | Client error (Bad Request, Not Found…) |
| `5000–5099` | Server / internal error |
