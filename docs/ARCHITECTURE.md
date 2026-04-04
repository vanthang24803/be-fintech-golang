# System Architecture

A comprehensive overview of the project's layered architecture, entity relationships,
and current implementation status.

---

## 1. Clean Architecture Layers

The project strictly follows **Clean Architecture** with 4 dependency layers.
Each layer only depends on the layer below it — never the reverse.

```mermaid
flowchart TD
    Client(["🌐 HTTP Client"])

    subgraph Fiber["Fiber Framework"]
        MW["Middleware\n(JWT Auth)"]
        H["Handler Layer\n(HTTP I/O)"]
    end

    subgraph Core["Business Core"]
        S["Service Layer\n(Business Logic)"]
    end

    subgraph Data["Data Layer"]
        R["Repository Layer\n(SQL Queries)"]
        DB[("PostgreSQL\n(via sqlx)")]
    end

    Client -->|"HTTP Request"| MW
    MW -->|"Authenticated"| H
    H -->|"calls"| S
    S -->|"calls"| R
    R -->|"queries"| DB
    DB -->|"rows"| R
    R -->|"models"| S
    S -->|"models"| H
    H -->|"JSON response"| Client
```

---

## 2. Project Module Map

```mermaid
graph LR
    subgraph Modules["Feature Modules (current)"]
        AUTH["🔐 Auth\n/auth"]
        USER["👤 User\n/profile"]
        SRC["💳 Source Payment\n/sources"]
        CAT["🏷️ Category\n/categories"]
        TXN["💸 Transaction\n/transactions"]
        FUND["🪙 Fund\n/funds"]
        DEV["📱 Device\n/devices"]
        NOTIF["🔔 Notification\n/notifications"]
    end

    subgraph Infra["Infrastructure"]
        MW_JWT["JWT Middleware"]
        SNOW["Snowflake ID"]
        RESP["Response Helper"]
        DB2[("PostgreSQL")]
    end

    AUTH --> MW_JWT
    USER --> MW_JWT
    SRC --> MW_JWT
    CAT --> MW_JWT
    TXN --> MW_JWT
    FUND --> MW_JWT
    DEV --> MW_JWT
    NOTIF --> MW_JWT
    TXN -->|"updates balance"| SRC
    TXN --> CAT
    TXN -->|"emits"| NOTIF
    FUND -->|"emits"| NOTIF
    AUTH -->|"new device alert"| NOTIF
    AUTH -->|"registers"| DEV
    DEV -->|"FIDO2 biometric"| AUTH
    AUTH --> SNOW
    SRC --> SNOW
    CAT --> SNOW
    TXN --> SNOW
    FUND --> SNOW
    DEV --> SNOW
    NOTIF --> SNOW
    AUTH --> DB2
    USER --> DB2
    SRC --> DB2
    CAT --> DB2
    TXN --> DB2
    FUND --> DB2
    DEV --> DB2
    NOTIF --> DB2
```

---

## 3. Entity Relationship Diagram

```mermaid
erDiagram
    users {
        bigint id PK
        string username
        string email
        string password_hash
        timestamptz created_at
        timestamptz updated_at
    }

    sourcepayment {
        bigint id PK
        bigint user_id FK
        string name
        string type
        numeric balance
        string currency
        timestamptz created_at
        timestamptz updated_at
    }

    categories {
        bigint id PK
        bigint user_id FK
        string name
        string type
        string icon
        timestamptz created_at
        timestamptz updated_at
    }

    transactions {
        bigint id PK
        bigint user_id FK
        bigint sourcepayment_id FK
        bigint category_id FK
        numeric amount
        string type
        string description
        timestamptz transaction_date
        timestamptz created_at
        timestamptz updated_at
    }

    funds {
        bigint id PK
        bigint user_id FK
        string name
        text description
        numeric target_amount
        numeric balance
        string currency
        timestamptz created_at
        timestamptz updated_at
    }

    refresh_tokens {
        bigint id PK
        bigint user_id FK
        string token
        boolean revoked
        timestamptz expires_at
        timestamptz created_at
    }

    devices {
        bigint id PK
        bigint user_id FK
        string device_fingerprint
        string device_name
        string platform
        text push_token
        text fido_credential_id
        text fido_public_key
        bigint fido_sign_count
        boolean is_trusted
        boolean is_active
        timestamptz last_used_at
        timestamptz created_at
        timestamptz updated_at
    }

    notifications {
        bigint id PK
        bigint user_id FK
        string source
        bigint source_id
        string type
        string title
        text body
        jsonb metadata
        boolean is_read
        timestamptz read_at
        timestamptz created_at
    }

    users ||--o{ sourcepayment : "owns"
    users ||--o{ categories : "owns"
    users ||--o{ transactions : "makes"
    users ||--o{ funds : "holds"
    users ||--o{ refresh_tokens : "has"
    users ||--o{ devices : "registers"
    users ||--o{ notifications : "receives"
    sourcepayment ||--o{ transactions : "used in"
    categories ||--o{ transactions : "classifies"
```

---

## 4. Request Lifecycle

```mermaid
sequenceDiagram
    participant C as Client
    participant MW as JWT Middleware
    participant H as Handler
    participant S as Service
    participant R as Repository
    participant DB as PostgreSQL

    C->>MW: POST /api/v1/funds/:id/deposit
    MW->>MW: Validate Bearer token
    alt Token invalid
        MW-->>C: 401 Unauthorized
    end
    MW->>H: ctx + userID injected
    H->>H: Parse & validate body
    H->>S: Deposit(id, userID, req)
    S->>S: Business validation (amount > 0)
    S->>R: Deposit(id, userID, amount)
    R->>DB: UPDATE funds SET balance = balance + $1 ...
    DB-->>R: Updated row
    R-->>S: *Fund
    S-->>H: *Fund
    H-->>C: 200 { code: 2000, data: fund }
```

---

## 5. API Endpoint Reference

> **Convention:** All endpoints use `POST`. Action intent is expressed in the URL path, not the HTTP method.

### Public Endpoints (no auth required)

| Method | Path | Description |
|---|---|---|
| `POST` | `/api/v1/auth/register` | Register a new user |
| `POST` | `/api/v1/auth/login` | Log in, receive JWT pair |
| `POST` | `/api/v1/auth/refresh` | Refresh access token |
| `POST` | `/api/v1/auth/logout` | Revoke refresh token |
| `POST` | `/api/v1/auth/google/url` | Get Google login URL |
| `POST` | `/api/v1/auth/google/callback` | OAuth2 callback handler |

### Protected Endpoints (Bearer JWT required)

#### 👤 User
| Method | Path | Description |
|---|---|---|
| `POST` | `/api/v1/profile/me` | Get current user profile |

#### 💳 Source Payment
| Method | Path | Description |
|---|---|---|
| `POST` | `/api/v1/sources/create` | Create a payment source |
| `POST` | `/api/v1/sources/list` | List all sources |
| `POST` | `/api/v1/sources/detail/:id` | Get a source by ID |
| `POST` | `/api/v1/sources/update/:id` | Update a source |
| `POST` | `/api/v1/sources/delete/:id` | Delete a source |

#### 🏷️ Category
| Method | Path | Description |
|---|---|---|
| `POST` | `/api/v1/categories/create` | Create a category |
| `POST` | `/api/v1/categories/list` | List all categories |
| `POST` | `/api/v1/categories/detail/:id` | Get a category by ID |
| `POST` | `/api/v1/categories/update/:id` | Update a category |
| `POST` | `/api/v1/categories/delete/:id` | Delete a category |

#### 💸 Transaction
| Method | Path | Body Params | Description |
|---|---|---|---|
| `POST` | `/api/v1/transactions/create` | — | Create a transaction |
| `POST` | `/api/v1/transactions/list` | `type`, `category_id`, `source_id` | List with filters |
| `POST` | `/api/v1/transactions/detail/:id` | — | Get a transaction |
| `POST` | `/api/v1/transactions/update/:id` | — | Update a transaction |
| `POST` | `/api/v1/transactions/delete/:id` | — | Delete a transaction |

#### 🪙 Fund
| Method | Path | Description |
|---|---|---|
| `POST` | `/api/v1/funds/create` | Create a fund |
| `POST` | `/api/v1/funds/list` | List all funds |
| `POST` | `/api/v1/funds/detail/:id` | Get a fund by ID |
| `POST` | `/api/v1/funds/update/:id` | Update fund metadata |
| `POST` | `/api/v1/funds/delete/:id` | Delete a fund |
| `POST` | `/api/v1/funds/deposit/:id` | Deposit money into fund |
| `POST` | `/api/v1/funds/withdraw/:id` | Withdraw money from fund |

#### 📊 Budget
| Method | Path | Description |
|---|---|---|
| `POST` | `/api/v1/budgets/create` | Create a new budget limit |
| `POST` | `/api/v1/budgets/list` | List budgets with progress |
| `POST` | `/api/v1/budgets/detail/:id` | Get budget details & spending |
| `POST` | `/api/v1/budgets/update/:id` | Update budget amount/status |
| `POST` | `/api/v1/budgets/delete/:id` | Delete a budget |

#### 📈 Reports & Analytics
| Method | Path | Body Params | Description |
|---|---|---|---|
| `POST` | `/api/v1/reports/category-summary` | `start_date`, `end_date` | Spending breakdown by category |
| `POST` | `/api/v1/reports/monthly-trend` | `months` | Income vs Expense trend |

#### 📱 Device
| Method | Path | Description |
|---|---|---|
| `POST` | `/api/v1/devices/register` | Register a new device (first login) |
| `POST` | `/api/v1/devices/list` | List user's registered devices |
| `POST` | `/api/v1/devices/delete/:id` | Remove / untrust a device |
| `POST` | `/api/v1/devices/biometric/enroll/:id` | Enroll FIDO2 biometric credential (skeleton) |
| `POST` | `/api/v1/auth/biometric` | Authenticate with FIDO2 (public - skeleton) |

#### 🔔 Notification
| Method | Path | Body Params | Description |
|---|---|---|---|
| `POST` | `/api/v1/notifications/list` | `source`, `is_read` | List notifications with filters |
| `POST` | `/api/v1/notifications/unread-count` | — | Get count of unread notifications |
| `POST` | `/api/v1/notifications/mark-read` | `ids[]` | Bulk mark notifications as read |
| `POST` | `/api/v1/notifications/delete/:id` | — | Delete a notification |

---

## 6. Implementation Status

```mermaid
gantt
    title Feature Implementation Progress
    dateFormat  YYYY-MM-DD
    section Authentication
    Register / Login / Refresh / Logout   :done, 2026-03-20, 2026-03-26
    JWT Middleware                         :done, 2026-03-26, 2026-03-27
    section Core Entities
    Source Payment CRUD                    :done, 2026-03-26, 2026-03-27
    Category CRUD                          :done, 2026-03-26, 2026-03-27
    Transaction CRUD + Balance Sync        :done, 2026-03-27, 2026-03-28
    Fund CRUD + Deposit/Withdraw           :done, 2026-03-29, 2026-03-29
    section Security & Engagement
    Device Schema + One-Device-One-Account :done, 2026-03-29, 2026-03-29
    Notification Schema + Multi-source     :done, 2026-03-29, 2026-03-29
    Device CRUD + One-Device Enforcement   :done, 2026-03-29, 2026-03-29
    Notification CRUD Logic                :done, 2026-03-29, 2026-03-29
    Budget CRUD + Limit Enforcement        :done, 2026-03-29, 2026-03-29
    Reports & Analytics Aggregation        :done, 2026-03-29, 2026-03-29
    FIDO Middleware (Step-up Auth)         :done, 2026-03-29, 2026-03-29
    Google OAuth 2.0 Integration           :done, 2026-03-29, 2026-03-29
    Savings Goals Module Implementation    :done, 2026-03-29, 2026-03-29
    Push Delivery Integration (FCM)        :done, 2026-03-29, 2026-03-29
    FIDO2 WebAuthn Implementation          :active, 2026-04-01, 2026-04-05
    section Planned
    Recurring Transactions                 :2026-04-13, 2026-04-18
```

### Current Coverage

| Module | Model | Migration | Repository | Service | Handler | Routes | Status |
|---|:---:|:---:|:---:|:---:|:---:|:---:|---|
| Auth | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | **Done** |
| User | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | **Done** |
| Source Payment | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | **Done** |
| Category | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | **Done** |
| Transaction | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | **Done** |
| Fund | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | **Done** |
| **Device** | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | **Done** |
| **Notification** | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | **Done** |
| **Budget** | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | **Done** |
| **Reports** | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | **Done** |
| **Savings Goal** | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | **Done** |
<line_break_filler>
