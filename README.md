# Expense Manager API

Backend service for a personal finance application built with Go, Fiber, PostgreSQL, and Redis. The API covers authentication, profile management, payment sources, categories, transactions, funds, budgets, savings goals, notifications, devices, reports, and FIDO2-based biometric step-up authentication.

## Highlights

- Clean architecture with handler, service, repository, and database layers
- JWT authentication with refresh tokens
- Google OAuth login flow
- FIDO2/WebAuthn step-up authentication for sensitive operations
- PostgreSQL persistence with embedded automatic migrations on startup
- Redis-backed async jobs for notifications and background processing
- OpenAPI spec with Scalar Docs served by the API

## Tech Stack

- Go `1.25`
- Fiber v2
- PostgreSQL
- Redis
- SQLX + PGX
- Asynq
- Zap logger

## Project Structure

```text
cmd/api/                 API entrypoint
configs/                 Environment-based configuration
internal/database/       DB connection and embedded migrations
internal/docs/           OpenAPI spec and Scalar Docs routes
internal/handler/        HTTP handlers
internal/middleware/     Fiber middleware
internal/models/         Domain/data models
internal/repository/     Data access layer
internal/router/         Route registration
internal/server/         Fiber server bootstrap
internal/service/        Business logic
internal/worker/         Background worker handlers
pkg/                     Shared packages
docs/                    Project documentation
```

## Prerequisites

- Go `1.25+`
- Docker and Docker Compose
- PostgreSQL `15+` if not using Docker
- Redis `7+` if not using Docker

## Environment Setup

1. Copy the example environment file:

```bash
cp .env.example .env
```

2. Set the required values in `.env`.

Minimum required:

- `DATABASE_URL`
- `JWT_SECRET`
- `JWT_REFRESH_SECRET`

Common optional values:

- `PORT` default: `8386`
- `GOOGLE_CLIENT_ID`
- `GOOGLE_CLIENT_SECRET`
- `GOOGLE_REDIRECT_URL`
- `FIREBASE_SERVICE_ACCOUNT_JSON`
- `REDIS_ADDR`
- `REDIS_PASSWORD`
- `WEBAUTHN_RP_ID`
- `WEBAUTHN_RP_NAME`
- `WEBAUTHN_ORIGIN`

Example PostgreSQL DSN:

```env
DATABASE_URL=postgres://postgres:postgres@localhost:5432/go_backend?sslmode=disable
```

## Run Locally

Start PostgreSQL and Redis:

```bash
docker compose up -d
```

Start the API:

```bash
go run ./cmd/api
```

Notes:

- Database migrations run automatically during startup from `internal/database/migrations/`.
- The server listens on `PORT`, which defaults to `8386`.
- If Firebase credentials are missing, the app falls back to a mock push client.

## API Docs

Once the server is running:

- Scalar Docs: `http://localhost:8386/docs`
- OpenAPI file: `http://localhost:8386/openapi.yaml`

If you change `PORT`, replace `8386` with your configured value.

## Development Workflow

Run tests:

```bash
go test ./...
```

Format Go code:

```bash
gofmt -w .
```

Useful docs already in the repo:

- [docs/ARCHITECTURE.md](/home/thang-nv/Workspace/GO-Proj/docs/ARCHITECTURE.md)
- [docs/STANDARDS.md](/home/thang-nv/Workspace/GO-Proj/docs/STANDARDS.md)
- [docs/RESPONSE_FORMAT.md](/home/thang-nv/Workspace/GO-Proj/docs/RESPONSE_FORMAT.md)

## Current API Areas

- Auth and token management
- Profile management
- Source payments
- Categories
- Transactions
- Funds
- Budgets
- Savings goals
- Notifications
- Devices and biometric enrollment
- Reports

## Notes

- Most API endpoints are currently modeled as `POST` operations, including list and detail actions.
- Sensitive flows such as profile updates, withdrawals, and some deletions require FIDO2 step-up authentication.
- Scalar Docs uses the checked-in OpenAPI file at `internal/docs/openapi.yaml` as the source of truth.
