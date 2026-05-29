# Backend Testing Design

## Goal

Achieve 70-90% test coverage on the Go backend through unit tests for pure logic and HTTP integration tests with mocked database repositories. The work involves introducing repository interfaces, refactoring handlers from package-level functions to struct methods, and writing comprehensive test suites using `testify`.

---

## Architecture Change: Repository Pattern

### Current State

Handlers are package-level functions that call `database.DB` (global `*gorm.DB`) directly. No dependency injection. No interfaces.

### Target State

Each domain package defines a `Repository` interface and a GORM implementation. Handlers become methods on a `Handler` struct that receives the repository interface. Tests inject `testify/mock` mocks.

### Before / After

| Before | After |
|--------|-------|
| `auth.Register(c *fiber.Ctx)` | `auth.NewHandler(repo).Register(c)` |
| `shipment.List(c *fiber.Ctx)` | `shipment.NewHandler(repo).List(c)` |
| `hub.Create(c *fiber.Ctx)` | `hub.NewHandler(repo).Create(c)` |

In `cmd/server/main.go`, real GORM repos are wired; in tests, mocks are injected.

---

## Repository Interfaces

### auth.Repository

```go
type Repository interface {
    Create(user *models.User) error
    FindByEmail(email string) (*models.User, error)
    FindByID(id uint) (*models.User, error)
}
```

### shipment.Repository

```go
type ShipmentFilter struct {
    Page          int
    Limit         int
    Search        string
    Status        string
    ExcludeStatus string
}

type StatusCountResult struct {
    Status string `json:"status"`
    Count  int64  `json:"count"`
}

type MonthCountResult struct {
    Month string `json:"month"`
    Count int64  `json:"count"`
}

type DayCountResult struct {
    Day   string `json:"day"`
    Count int64  `json:"count"`
}

type Repository interface {
    // Shipment CRUD
    List(params ShipmentFilter) ([]models.Shipment, int64, error)
    FindByOrderID(orderID string) (*models.Shipment, error)
    FindByTrackingNumber(trackingNumber string) (*models.Shipment, error)
    Create(shipment *models.Shipment) error
    Save(shipment *models.Shipment) error
    Delete(shipment *models.Shipment) error
    // Events
    CreateEvent(event *models.ShipmentEvent) error
    FindEventsByShipmentID(shipmentID uint) ([]models.ShipmentEvent, error)
    DeleteShipmentEvents(shipmentID uint) error
    // Analytics
    Count() (int64, error)
    CountActive() (int64, error)
    CountByStatus() ([]StatusCountResult, error)
    CountByMonth() ([]MonthCountResult, error)
    CountByDayOfWeek() ([]DayCountResult, error)
}
```

### hub.Repository

```go
type Repository interface {
    FindAll() ([]models.Hub, error)
    FindByID(id string) (*models.Hub, error)
    Create(hub *models.Hub) error
    Save(hub *models.Hub) error
    Delete(id string) error
}
```

---

## GORM Repository Implementations

Each domain gets a `gorm_repository.go` file with the concrete implementation. ID generation (`GenerateOrderID`, `GenerateHubID`) moves from `pkg/utils/idgen.go` into the respective GORM repository's `Create` method, since both query the database to find the next available ID. The `pkg/utils/idgen.go` file is deleted.

---

## Handler Refactoring

Each handler package gets:

```go
type Handler struct {
    repo Repository
}

func NewHandler(repo Repository) *Handler {
    return &Handler{repo: repo}
}
```

Every exported handler function becomes a method on `Handler`. The `main.go` route registration changes from:

```go
authGroup.Post("/register", middleware.RateLimitAuth(), auth.Register)
```

to:

```go
authHandler := auth.NewHandler(authRepo)
authGroup.Post("/register", middleware.RateLimitAuth(), authHandler.Register)
```

### Cross-package dependencies

- **tracking handler**: accepts `shipment.Repository` (needs `FindByTrackingNumber` + `FindEventsByShipmentID`)
- **analytics handler**: accepts `shipment.Repository` (needs `Count`, `CountActive`, `CountByStatus`, `CountByMonth`, `CountByDayOfWeek`); the region-mapping logic stays in the handler itself using `data.ThailandProvinceRegion`

---

## Test Structure

### Unit Tests (pure logic, no DB)

| Package | What's tested | File |
|---------|--------------|------|
| `pkg/utils` | `HashPassword`, `CheckPassword`, `ComposeAddress`, `Success`, `Error`, `SuccessWithPagination` | `response_test.go`, `hash_test.go` |
| `internal/middleware` | Rate limiter: allow/deny boundaries, cleanup loop, concurrent safety. Auth: valid/invalid/expired JWT, cookie fallback, missing token | `ratelimit_test.go`, `auth_test.go` |
| `internal/models` | Shipment.BeforeSave/AfterFind, Hub.BeforeSave/AfterFind | `shipment_test.go`, `hub_test.go` |
| `internal/config` | Env var loading, defaults | `config_test.go` |
| `internal/data` | Province-to-region mapping completeness | `regions_test.go` |
| `internal/auth` | `setAuthCookie` helper (cookie fields), token generation via `jwt.NewWithClaims` | `handler_test.go` |

### HTTP Integration Tests (Fiber app.Test with mocked repos)

Each test creates a minimal Fiber app with the handler under test and one mocked repository, then sends real HTTP requests and asserts response JSON.

#### auth/handler_test.go

| Test | Cases |
|------|-------|
| Register | success, missing name/email/password, duplicate email |
| Login | success, wrong password, non-existent email |
| Me | success with valid token, 401 with no token, 401 with expired token |
| Logout | clears cookie, returns success |

#### shipment/handler_test.go

| Test | Cases |
|------|-------|
| List | paginated results, search filter, status filter, exclude_status, empty results |
| Create | success, missing customer/receiver fields, zero weight, zero items, missing carrier |
| GetByID | found, not found |
| Update | success, not found |
| UpdateStatus | with hub (sets coords), without hub, invalid hub ID |
| Delete | success, not found |

#### hub/handler_test.go

| Test | Cases |
|------|-------|
| List | returns all hubs |
| GetByID | found, not found |
| Create | success, invalid body |
| Update | success, not found, invalid body |
| Delete | success, not found |

#### tracking/handler_test.go

| Test | Cases |
|------|-------|
| Track | found with events, not found |

#### analytics/handler_test.go

| Test | Cases |
|------|-------|
| Overview | returns aggregate stats |
| TimeSeries | returns monthly + day-of-week breakdowns |

---

## New and Modified Files

### New files (18)

| File | Purpose |
|------|---------|
| `internal/auth/repository.go` | Auth repository interface |
| `internal/auth/gorm_repository.go` | Auth GORM implementation |
| `internal/auth/handler_test.go` | Auth HTTP + unit tests |
| `internal/shipment/repository.go` | Shipment repository interface + shared result types |
| `internal/shipment/gorm_repository.go` | Shipment GORM implementation |
| `internal/shipment/handler_test.go` | Shipment HTTP integration tests |
| `internal/hub/repository.go` | Hub repository interface |
| `internal/hub/gorm_repository.go` | Hub GORM implementation |
| `internal/hub/handler_test.go` | Hub HTTP integration tests |
| `internal/tracking/handler_test.go` | Tracking HTTP integration tests |
| `internal/analytics/handler_test.go` | Analytics HTTP integration tests |
| `internal/middleware/ratelimit_test.go` | Rate limiter unit tests |
| `internal/middleware/auth_test.go` | Auth middleware unit tests |
| `internal/models/shipment_test.go` | Shipment GORM hook tests |
| `internal/models/hub_test.go` | Hub GORM hook tests |
| `internal/config/config_test.go` | Config loading tests |
| `internal/data/regions_test.go` | Province region mapping tests |
| `pkg/utils/response_test.go` | Response helper tests |
| `pkg/utils/hash_test.go` | Hash helper tests |

### Modified files (6)

| File | Change |
|------|--------|
| `internal/auth/handler.go` | Functions → `Handler` methods, remove `database.DB` dep |
| `internal/shipment/handler.go` | Functions → `Handler` methods, remove `database.DB` dep |
| `internal/hub/handler.go` | Functions → `Handler` methods, remove `database.DB` dep |
| `internal/tracking/handler.go` | Functions → `Handler` methods, accept `shipment.Repository` |
| `internal/analytics/handler.go` | Functions → `Handler` methods, accept `shipment.Repository` |
| `cmd/server/main.go` | Instantiate repos, pass to handlers via constructors |

### Deleted files (1)

| File | Reason |
|------|--------|
| `pkg/utils/idgen.go` | Logic moves into GORM repository Create methods |

### go.mod changes

- Promote `github.com/stretchr/testify` from indirect to direct dependency

---

## Testing Conventions

- **Test runner:** Go standard `testing` package
- **Assertions:** `github.com/stretchr/testify/assert` and `github.com/stretchr/testify/require`
- **Mocking:** `github.com/stretchr/testify/mock`
- **Mocks location:** In-package `_test.go` files define mock structs (no generated code, no mock package)
- **HTTP testing:** Fiber's `app.Test()` sends `*http.Request` directly, no server needed
- **Test DB:** none — all DB calls go through mocked repositories
- **Run:** `go test ./...` (or `go test -v -count=1 -race ./...`)
