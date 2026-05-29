# Backend Testing Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Achieve 70-90% test coverage on the Go backend through repository pattern refactoring, unit tests for pure logic, and HTTP integration tests with mocked database repositories.

**Architecture:** Introduce per-domain `Repository` interfaces (auth, hub, shipment) so handlers can be tested with `testify/mock` mocks instead of calling `database.DB` directly. Handlers become struct methods receiving the interface. Tracking and analytics handlers receive a `shipment.Repository`. Shipment handler also receives a `HubRepository` interface for status updates. GORM implementations remain for production use.

**Tech Stack:** Go 1.24, Fiber v2, GORM v2, testify (assert, require, mock), standard `testing` package, Fiber's `app.Test()`

---

## File Structure

### New files (19)

| File | Responsibility |
|------|---------------|
| `pkg/utils/response_test.go` | Test ComposeAddress, Success, Error, SuccessWithPagination |
| `pkg/utils/hash_test.go` | Test HashPassword, CheckPassword |
| `internal/config/config_test.go` | Test env loading and defaults |
| `internal/data/regions_test.go` | Test province-to-region mapping |
| `internal/models/shipment_test.go` | Test Shipment GeoPoint sync hooks |
| `internal/models/hub_test.go` | Test Hub GeoPoint sync hooks |
| `internal/middleware/ratelimit_test.go` | Test rateLimiter allow/deny/cleanup |
| `internal/middleware/auth_test.go` | Test AuthRequired JWT parsing, cookie fallback |
| `internal/auth/repository.go` | Auth Repository interface |
| `internal/auth/gorm_repository.go` | Auth GORM implementation |
| `internal/auth/handler_test.go` | Auth HTTP integration + helper unit tests |
| `internal/hub/repository.go` | Hub Repository interface |
| `internal/hub/gorm_repository.go` | Hub GORM implementation |
| `internal/hub/handler_test.go` | Hub HTTP integration tests |
| `internal/shipment/repository.go` | Shipment Repository interface + result types |
| `internal/shipment/gorm_repository.go` | Shipment GORM implementation |
| `internal/shipment/handler_test.go` | Shipment HTTP integration tests |
| `internal/tracking/handler_test.go` | Tracking HTTP integration tests |
| `internal/analytics/handler_test.go` | Analytics HTTP integration tests |

### Modified files (6)

| File | Change |
|------|--------|
| `internal/auth/handler.go` | Functions -> Handler methods, inject Repository |
| `internal/hub/handler.go` | Functions -> Handler methods, inject Repository |
| `internal/shipment/handler.go` | Functions -> Handler methods, inject Repository + HubRepository, remove ID generation |
| `internal/tracking/handler.go` | Functions -> Handler methods, inject shipment.Repository |
| `internal/analytics/handler.go` | Functions -> Handler methods, inject shipment.Repository |
| `cmd/server/main.go` | Instantiate GORM repos, pass to handler constructors |

### Deleted files (1)

| File | Reason |
|------|--------|
| `pkg/utils/idgen.go` | ID generation logic moves into GORM repository Create methods |

### go.mod change

- Promote `github.com/stretchr/testify` from indirect to direct dependency

---

## Tasks

### Task 1: Add testify dependency

- [ ] **Step 1: Promote testify to direct dependency**

```bash
go get github.com/stretchr/testify@v1.10.0
```

- [ ] **Step 2: Tidy**

```bash
go mod tidy
```

- [ ] **Step 3: Verify**

```bash
go build ./... && go vet ./...
```

Expected: clean output.

- [ ] **Step 4: Commit**

```bash
git add backend/go.mod backend/go.sum
git commit -m "chore: add testify dependency for testing"
```

---

### Task 2: Unit tests for utils, config, data, models

- [ ] **Step 1: Write `pkg/utils/hash_test.go`**

```go
package utils_test

import (
    "testing"
    "github.com/narinsak-u/backend/pkg/utils"
    "github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
    hash, err := utils.HashPassword("myPassword123")
    assert.NoError(t, err)
    assert.NotEmpty(t, hash)
    assert.NotEqual(t, "myPassword123", hash)
}

func TestCheckPassword(t *testing.T) {
    hash, _ := utils.HashPassword("correctPassword")
    assert.True(t, utils.CheckPassword("correctPassword", hash))
    assert.False(t, utils.CheckPassword("wrongPassword", hash))
}

func TestCheckPassword_WrongHash(t *testing.T) {
    assert.False(t, utils.CheckPassword("any", "invalidhash"))
}
```

- [ ] **Step 2: Write `pkg/utils/response_test.go`**

```go
package utils_test

import (
    "encoding/json"
    "testing"
    "github.com/gofiber/fiber/v2"
    "github.com/narinsak-u/backend/internal/models"
    "github.com/narinsak-u/backend/pkg/utils"
    "github.com/stretchr/testify/assert"
)

func TestComposeAddress(t *testing.T) {
    c := models.ContactInfo{SubDistrict: "S1", District: "D1", Province: "P1"}
    result := utils.ComposeAddress(c)
    assert.Equal(t, "S1, D1, P1", result)
}

func TestSuccess(t *testing.T) {
    app := fiber.New()
    c := app.AcquireCtx(&fiber.Ctx{})
    defer app.ReleaseCtx(c)

    err := utils.Success(c, fiber.Map{"key": "value"})
    assert.NoError(t, err)
    assert.Equal(t, 200, c.Response().StatusCode())

    var resp map[string]interface{}
    json.Unmarshal(c.Response().Body(), &resp)
    assert.True(t, resp["success"].(bool))
    assert.Equal(t, "value", resp["data"].(map[string]interface{})["key"])
}

func TestError(t *testing.T) {
    app := fiber.New()
    c := app.AcquireCtx(&fiber.Ctx{})
    defer app.ReleaseCtx(c)

    err := utils.Error(c, 400, "bad request")
    assert.NoError(t, err)
    assert.Equal(t, 400, c.Response().StatusCode())

    var resp map[string]interface{}
    json.Unmarshal(c.Response().Body(), &resp)
    assert.False(t, resp["success"].(bool))
    assert.Equal(t, "bad request", resp["error"].(string))
}

func TestSuccessWithPagination(t *testing.T) {
    app := fiber.New()
    c := app.AcquireCtx(&fiber.Ctx{})
    defer app.ReleaseCtx(c)

    data := []string{"a", "b"}
    err := utils.SuccessWithPagination(c, data, 1, 10, 25)
    assert.NoError(t, err)
    assert.Equal(t, 200, c.Response().StatusCode())

    var resp map[string]interface{}
    json.Unmarshal(c.Response().Body(), &resp)
    pagination := resp["pagination"].(map[string]interface{})
    assert.Equal(t, float64(1), pagination["page"])
    assert.Equal(t, float64(10), pagination["limit"])
    assert.Equal(t, float64(25), pagination["total"])
    assert.Equal(t, float64(3), pagination["totalPages"])
}

func TestSuccessWithPagination_ZeroTotal(t *testing.T) {
    app := fiber.New()
    c := app.AcquireCtx(&fiber.Ctx{})
    defer app.ReleaseCtx(c)

    utils.SuccessWithPagination(c, []interface{}{}, 1, 10, 0)
    var resp map[string]interface{}
    json.Unmarshal(c.Response().Body(), &resp)
    pagination := resp["pagination"].(map[string]interface{})
    assert.Equal(t, float64(1), pagination["totalPages"])
}
```

- [ ] **Step 3: Write `internal/config/config_test.go`**

```go
package config_test

import (
    "os"
    "testing"
    "time"
    "github.com/narinsak-u/backend/internal/config"
    "github.com/stretchr/testify/assert"
)

func TestLoad_Defaults(t *testing.T) {
    os.Unsetenv("PORT")
    os.Unsetenv("DATABASE_URL")
    os.Unsetenv("JWT_SECRET")
    config.Load()
    assert.Equal(t, "8080", config.App.Port)
    assert.Equal(t, "postgres://user:pass@localhost:5432/shipments", config.App.DatabaseURL)
    assert.Equal(t, "change-me", config.App.JWTSecret)
    assert.Equal(t, 24*time.Hour, config.App.JWTTTL)
}

func TestLoad_FromEnv(t *testing.T) {
    os.Setenv("PORT", "9090")
    os.Setenv("DATABASE_URL", "postgres://test:test@localhost:9999/testdb")
    os.Setenv("JWT_SECRET", "my-secret-key")
    defer func() {
        os.Unsetenv("PORT"); os.Unsetenv("DATABASE_URL"); os.Unsetenv("JWT_SECRET")
    }()
    config.Load()
    assert.Equal(t, "9090", config.App.Port)
    assert.Equal(t, "postgres://test:test@localhost:9999/testdb", config.App.DatabaseURL)
    assert.Equal(t, "my-secret-key", config.App.JWTSecret)
}
```

- [ ] **Step 4: Write `internal/data/regions_test.go`**

```go
package data_test

import (
    "testing"
    "github.com/narinsak-u/backend/internal/data"
    "github.com/stretchr/testify/assert"
)

func TestThailandProvinceRegion_KnownProvinces(t *testing.T) {
    tests := []struct {
        province string
        expected string
    }{
        {"กรุงเทพมหานคร", "Central"},
        {"ชลบุรี", "East"},
        {"เชียงใหม่", "North"},
        {"กาญจนบุรี", "West"},
        {"ขอนแก่น", "North-east"},
        {"ภูเก็ต", "South"},
    }
    for _, tt := range tests {
        t.Run(tt.province, func(t *testing.T) {
            assert.Equal(t, tt.expected, data.ThailandProvinceRegion[tt.province])
        })
    }
}

func TestThailandProvinceRegion_UnknownProvince(t *testing.T) {
    assert.Empty(t, data.ThailandProvinceRegion["Unknown"])
}

func TestThailandProvinceRegion_NotEmpty(t *testing.T) {
    assert.Greater(t, len(data.ThailandProvinceRegion), 70)
}
```

- [ ] **Step 5: Write `internal/models/shipment_test.go`**

```go
package models_test

import (
    "testing"
    "github.com/narinsak-u/backend/internal/models"
    "github.com/stretchr/testify/assert"
)

func TestShipment_BeforeSave(t *testing.T) {
    s := &models.Shipment{
        Customer:      models.ContactInfo{Coords: models.GeoPoint{Lat: 10.0, Lng: 20.0}},
        Receiver:      models.ContactInfo{Coords: models.GeoPoint{Lat: 30.0, Lng: 40.0}},
        CurrentCoords: models.GeoPoint{Lat: 25.0, Lng: 35.0},
    }
    err := s.BeforeSave(nil)
    assert.NoError(t, err)
    assert.Equal(t, 10.0, s.CustomerLat)
    assert.Equal(t, 20.0, s.CustomerLng)
    assert.Equal(t, 30.0, s.ReceiverLat)
    assert.Equal(t, 40.0, s.ReceiverLng)
    assert.Equal(t, 25.0, s.CurrentLat)
    assert.Equal(t, 35.0, s.CurrentLng)
}

func TestShipment_AfterFind(t *testing.T) {
    s := &models.Shipment{
        CustomerLat: 10.0, CustomerLng: 20.0,
        ReceiverLat: 30.0, ReceiverLng: 40.0,
        CurrentLat: 25.0, CurrentLng: 35.0,
    }
    err := s.AfterFind(nil)
    assert.NoError(t, err)
    assert.Equal(t, 10.0, s.Customer.Coords.Lat)
    assert.Equal(t, 20.0, s.Customer.Coords.Lng)
    assert.Equal(t, 30.0, s.Receiver.Coords.Lat)
    assert.Equal(t, 40.0, s.Receiver.Coords.Lng)
    assert.Equal(t, 25.0, s.CurrentCoords.Lat)
    assert.Equal(t, 35.0, s.CurrentCoords.Lng)
}

func TestShipment_RoundTrip(t *testing.T) {
    s := &models.Shipment{
        Customer:      models.ContactInfo{Coords: models.GeoPoint{Lat: 10.0, Lng: 20.0}},
        Receiver:      models.ContactInfo{Coords: models.GeoPoint{Lat: 30.0, Lng: 40.0}},
        CurrentCoords: models.GeoPoint{Lat: 25.0, Lng: 35.0},
    }
    s.BeforeSave(nil)
    s.AfterFind(nil)
    assert.Equal(t, 10.0, s.Customer.Coords.Lat)
    assert.Equal(t, 20.0, s.Customer.Coords.Lng)
    assert.Equal(t, 30.0, s.Receiver.Coords.Lat)
    assert.Equal(t, 40.0, s.Receiver.Coords.Lng)
    assert.Equal(t, 25.0, s.CurrentCoords.Lat)
    assert.Equal(t, 35.0, s.CurrentCoords.Lng)
}
```

- [ ] **Step 6: Write `internal/models/hub_test.go`**

```go
package models_test

import (
    "testing"
    "github.com/narinsak-u/backend/internal/models"
    "github.com/stretchr/testify/assert"
)

func TestHub_BeforeSave(t *testing.T) {
    h := &models.Hub{Coords: models.GeoPoint{Lat: 13.0, Lng: 100.0}}
    err := h.BeforeSave(nil)
    assert.NoError(t, err)
    assert.Equal(t, 13.0, h.Lat)
    assert.Equal(t, 100.0, h.Lng)
}

func TestHub_AfterFind(t *testing.T) {
    h := &models.Hub{Lat: 13.0, Lng: 100.0}
    err := h.AfterFind(nil)
    assert.NoError(t, err)
    assert.Equal(t, 13.0, h.Coords.Lat)
    assert.Equal(t, 100.0, h.Coords.Lng)
}

func TestHub_RoundTrip(t *testing.T) {
    h := &models.Hub{Coords: models.GeoPoint{Lat: 10.0, Lng: 20.0}}
    h.BeforeSave(nil)
    h.AfterFind(nil)
    assert.Equal(t, 10.0, h.Coords.Lat)
    assert.Equal(t, 20.0, h.Coords.Lng)
}
```

- [ ] **Step 7: Run all pure unit tests**

```bash
go test ./pkg/utils/ ./internal/config/ ./internal/data/ ./internal/models/ -v -count=1
```

- [ ] **Step 8: Commit**

```bash
git add backend/pkg/utils/hash_test.go backend/pkg/utils/response_test.go backend/internal/config/config_test.go backend/internal/data/regions_test.go backend/internal/models/shipment_test.go backend/internal/models/hub_test.go
git commit -m "test: add unit tests for utils, config, data, and models"
```

---

### Task 3: Unit tests for middleware (rate limiter, auth)

- [ ] **Step 1: Write `internal/middleware/ratelimit_test.go`**

```go
package middleware

import (
    "sync"
    "testing"
    "time"
    "github.com/stretchr/testify/assert"
)

func TestRateLimiter_AllowsUpToLimit(t *testing.T) {
    rl := newRateLimiter(3, time.Minute)
    ip := "192.168.1.1"
    assert.True(t, rl.allow(ip))
    assert.True(t, rl.allow(ip))
    assert.True(t, rl.allow(ip))
    assert.False(t, rl.allow(ip))
}

func TestRateLimiter_DifferentIPsAreIndependent(t *testing.T) {
    rl := newRateLimiter(2, time.Minute)
    assert.True(t, rl.allow("10.0.0.1"))
    assert.True(t, rl.allow("10.0.0.1"))
    assert.False(t, rl.allow("10.0.0.1"))
    assert.True(t, rl.allow("10.0.0.2"))
    assert.True(t, rl.allow("10.0.0.2"))
    assert.False(t, rl.allow("10.0.0.2"))
}

func TestRateLimiter_WindowExpiry(t *testing.T) {
    rl := &rateLimiter{
        requests: make(map[string][]time.Time),
        limit:    1,
        window:   30 * time.Millisecond,
    }
    ip := "192.168.1.1"
    assert.True(t, rl.allow(ip))
    assert.False(t, rl.allow(ip))
    time.Sleep(40 * time.Millisecond)
    assert.True(t, rl.allow(ip))
}

func TestRateLimiter_ConcurrentAccess(t *testing.T) {
    rl := newRateLimiter(100, time.Minute)
    var wg sync.WaitGroup
    for i := 0; i < 50; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            rl.allow("192.168.1.1")
        }()
    }
    wg.Wait()
}
```

- [ ] **Step 2: Write `internal/middleware/auth_test.go`**

```go
package middleware

import (
    "net/http"
    "net/http/httptest"
    "testing"
    "time"
    "github.com/gofiber/fiber/v2"
    "github.com/golang-jwt/jwt/v5"
    "github.com/narinsak-u/backend/internal/config"
    "github.com/stretchr/testify/assert"
)

func setupConfig() {
    config.App = config.Config{
        JWTSecret: "test-secret",
        JWTTTL:    24 * time.Hour,
    }
}

func validToken() string {
    claims := jwt.MapClaims{
        "user_id": float64(1),
        "role":    "customer",
        "exp":     time.Now().Add(1 * time.Hour).Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    s, _ := token.SignedString([]byte(config.App.JWTSecret))
    return s
}

func expiredToken() string {
    claims := jwt.MapClaims{
        "user_id": float64(1),
        "role":    "customer",
        "exp":     time.Now().Add(-1 * time.Hour).Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    s, _ := token.SignedString([]byte(config.App.JWTSecret))
    return s
}

func testApp() *fiber.App {
    app := fiber.New()
    app.Use(AuthRequired())
    app.Get("/test", func(c *fiber.Ctx) error {
        return c.SendString("ok")
    })
    return app
}

func TestAuthRequired_ValidToken_Header(t *testing.T) {
    setupConfig()
    req := httptest.NewRequest("GET", "/test", nil)
    req.Header.Set("Authorization", "Bearer "+validToken())
    resp, _ := testApp().Test(req, 1000)
    assert.Equal(t, 200, resp.StatusCode)
}

func TestAuthRequired_NoToken(t *testing.T) {
    setupConfig()
    resp, _ := testApp().Test(httptest.NewRequest("GET", "/test", nil), 1000)
    assert.Equal(t, 401, resp.StatusCode)
}

func TestAuthRequired_ExpiredToken(t *testing.T) {
    setupConfig()
    req := httptest.NewRequest("GET", "/test", nil)
    req.Header.Set("Authorization", "Bearer "+expiredToken())
    resp, _ := testApp().Test(req, 1000)
    assert.Equal(t, 401, resp.StatusCode)
}

func TestAuthRequired_CookieFallback(t *testing.T) {
    setupConfig()
    req := httptest.NewRequest("GET", "/test", nil)
    req.AddCookie(&http.Cookie{Name: "jwt", Value: validToken()})
    resp, _ := testApp().Test(req, 1000)
    assert.Equal(t, 200, resp.StatusCode)
}

func TestAuthRequired_InvalidSignature(t *testing.T) {
    setupConfig()
    req := httptest.NewRequest("GET", "/test", nil)
    req.Header.Set("Authorization", "Bearer invalidtoken")
    resp, _ := testApp().Test(req, 1000)
    assert.Equal(t, 401, resp.StatusCode)
}
```

- [ ] **Step 3: Run middleware tests**

```bash
go test ./internal/middleware/ -v -count=1 -race
```

- [ ] **Step 4: Commit**

```bash
git add backend/internal/middleware/ratelimit_test.go backend/internal/middleware/auth_test.go
git commit -m "test: add unit tests for rate limiter and auth middleware"
```

---

### Task 4: Auth repository + handler refactoring + tests

- [ ] **Step 1: Create `internal/auth/repository.go`**

```go
package auth

import "github.com/narinsak-u/backend/internal/models"

type Repository interface {
    Create(user *models.User) error
    FindByEmail(email string) (*models.User, error)
    FindByID(id uint) (*models.User, error)
}
```

- [ ] **Step 2: Create `internal/auth/gorm_repository.go`**

```go
package auth

import (
    "github.com/narinsak-u/backend/internal/models"
    "gorm.io/gorm"
)

type GormRepository struct{ db *gorm.DB }

func NewGormRepository(db *gorm.DB) *GormRepository {
    return &GormRepository{db: db}
}

func (r *GormRepository) Create(user *models.User) error {
    return r.db.Create(user).Error
}

func (r *GormRepository) FindByEmail(email string) (*models.User, error) {
    var user models.User
    if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
        return nil, err
    }
    return &user, nil
}

func (r *GormRepository) FindByID(id uint) (*models.User, error) {
    var user models.User
    if err := r.db.First(&user, id).Error; err != nil {
        return nil, err
    }
    return &user, nil
}
```

- [ ] **Step 3: Refactor `internal/auth/handler.go`**

Add Handler struct + constructor after imports:
```go
type Handler struct {
    repo Repository
}

func NewHandler(repo Repository) *Handler {
    return &Handler{repo: repo}
}
```

Convert each function to method (`func Register(...)` -> `func (h *Handler) Register(...)`, same for Login, Me, Logout).

Replace database.DB calls:
- **Register**: `database.DB.Create(&user)` -> `h.repo.Create(&user)`
- **Login**: `database.DB.Where("email = ?", req.Email).First(&user)` -> `user, err := h.repo.FindByEmail(req.Email)` (delete `var user models.User`)
- **Me**: `database.DB.First(&user, userID)` -> `user, err := h.repo.FindByID(userID)` (delete `var user models.User`)
- Remove `"github.com/narinsak-u/backend/internal/database"` import

- [ ] **Step 4: Verify compilation**

```bash
go build ./internal/auth/
```

- [ ] **Step 5: Create `internal/auth/handler_test.go`**

```go
package auth

import (
    "encoding/json"
    "net/http/httptest"
    "strings"
    "testing"
    "time"
    "github.com/gofiber/fiber/v2"
    "github.com/golang-jwt/jwt/v5"
    "github.com/narinsak-u/backend/internal/config"
    "github.com/narinsak-u/backend/internal/models"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/gorm"
)

func init() {
    config.App = config.Config{
        JWTSecret: "test-secret",
        JWTTTL:    24 * time.Minute,
        Port:      "8080",
    }
}

type mockRepo struct{ mock.Mock }

func (m *mockRepo) Create(user *models.User) error {
    args := m.Called(user)
    return args.Error(0)
}

func (m *mockRepo) FindByEmail(email string) (*models.User, error) {
    args := m.Called(email)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*models.User), args.Error(1)
}

func (m *mockRepo) FindByID(id uint) (*models.User, error) {
    args := m.Called(id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*models.User), args.Error(1)
}

func TestRegister_Success(t *testing.T) {
    repo := new(mockRepo)
    repo.On("Create", mock.MatchedBy(func(u *models.User) bool {
        return u.Name == "John" && u.Email == "john@test.com"
    })).Return(nil)

    app := fiber.New()
    h := NewHandler(repo)
    app.Post("/api/auth/register", h.Register)

    body := `{"name":"John","email":"john@test.com","password":"pass123"}`
    req := httptest.NewRequest("POST", "/api/auth/register", strings.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    resp, _ := app.Test(req, 1000)

    assert.Equal(t, 200, resp.StatusCode)
    var result map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&result)
    assert.True(t, result["success"].(bool))
    repo.AssertExpectations(t)
}

func TestRegister_ValidationErrors(t *testing.T) {
    repo := new(mockRepo)
    app := fiber.New()
    h := NewHandler(repo)
    app.Post("/api/auth/register", h.Register)

    for _, tt := range []struct{ name, body string }{
        {"missing name",  `{"email":"e@t.com","password":"p123"}`},
        {"missing email", `{"name":"John","password":"p123"}`},
        {"missing pass",  `{"name":"John","email":"e@t.com"}`},
    } {
        t.Run(tt.name, func(t *testing.T) {
            req := httptest.NewRequest("POST", "/api/auth/register", strings.NewReader(tt.body))
            req.Header.Set("Content-Type", "application/json")
            resp, _ := app.Test(req, 1000)
            assert.Equal(t, 400, resp.StatusCode)
        })
    }
}

func TestRegister_DuplicateEmail(t *testing.T) {
    repo := new(mockRepo)
    repo.On("Create", mock.Anything).Return(gorm.ErrDuplicatedKey)

    app := fiber.New()
    h := NewHandler(repo)
    app.Post("/api/auth/register", h.Register)

    body := `{"name":"John","email":"dup@test.com","password":"pass123"}`
    req := httptest.NewRequest("POST", "/api/auth/register", strings.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    resp, _ := app.Test(req, 1000)
    assert.Equal(t, 409, resp.StatusCode)
}

func TestLogin_Success(t *testing.T) {
    repo := new(mockRepo)
    hash, _ := bcrypt.GenerateFromPassword([]byte("pass123"), bcrypt.DefaultCost)
    repo.On("FindByEmail", "john@test.com").Return(&models.User{
        ID: 1, Name: "John", Email: "john@test.com",
        Password: string(hash), Role: "customer",
    }, nil)

    app := fiber.New()
    h := NewHandler(repo)
    app.Post("/api/auth/login", h.Login)

    body := `{"email":"john@test.com","password":"pass123"}`
    req := httptest.NewRequest("POST", "/api/auth/login", strings.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    resp, _ := app.Test(req, 1000)

    assert.Equal(t, 200, resp.StatusCode)
    var result map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&result)
    data := result["data"].(map[string]interface{})
    assert.NotEmpty(t, data["token"])
    repo.AssertExpectations(t)
}

func TestLogin_WrongPassword(t *testing.T) {
    repo := new(mockRepo)
    hash, _ := bcrypt.GenerateFromPassword([]byte("correct"), bcrypt.DefaultCost)
    repo.On("FindByEmail", "john@test.com").Return(&models.User{
        ID: 1, Email: "john@test.com", Password: string(hash),
    }, nil)

    app := fiber.New()
    h := NewHandler(repo)
    app.Post("/api/auth/login", h.Login)

    body := `{"email":"john@test.com","password":"wrong"}`
    req := httptest.NewRequest("POST", "/api/auth/login", strings.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    resp, _ := app.Test(req, 1000)
    assert.Equal(t, 401, resp.StatusCode)
}

func TestLogin_UserNotFound(t *testing.T) {
    repo := new(mockRepo)
    repo.On("FindByEmail", "unknown@test.com").Return(nil, gorm.ErrRecordNotFound)

    app := fiber.New()
    h := NewHandler(repo)
    app.Post("/api/auth/login", h.Login)

    body := `{"email":"unknown@test.com","password":"pass123"}`
    req := httptest.NewRequest("POST", "/api/auth/login", strings.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    resp, _ := app.Test(req, 1000)
    assert.Equal(t, 401, resp.StatusCode)
}

func TestMe_Success(t *testing.T) {
    repo := new(mockRepo)
    repo.On("FindByID", uint(1)).Return(&models.User{
        ID: 1, Name: "John", Email: "john@test.com", Role: "customer",
    }, nil)

    app := fiber.New()
    h := NewHandler(repo)
    app.Get("/api/auth/me", func(c *fiber.Ctx) error {
        c.Locals("user_id", uint(1))
        return h.Me(c)
    })

    resp, _ := app.Test(httptest.NewRequest("GET", "/api/auth/me", nil), 1000)
    assert.Equal(t, 200, resp.StatusCode)
}

func TestMe_NotAuthenticated(t *testing.T) {
    repo := new(mockRepo)
    app := fiber.New()
    h := NewHandler(repo)
    app.Get("/api/auth/me", h.Me)

    resp, _ := app.Test(httptest.NewRequest("GET", "/api/auth/me", nil), 1000)
    assert.Equal(t, 401, resp.StatusCode)
}

func TestMe_UserNotFound(t *testing.T) {
    repo := new(mockRepo)
    repo.On("FindByID", uint(999)).Return(nil, gorm.ErrRecordNotFound)

    app := fiber.New()
    h := NewHandler(repo)
    app.Get("/api/auth/me", func(c *fiber.Ctx) error {
        c.Locals("user_id", uint(999))
        return h.Me(c)
    })

    resp, _ := app.Test(httptest.NewRequest("GET", "/api/auth/me", nil), 1000)
    assert.Equal(t, 404, resp.StatusCode)
}

func TestLogout(t *testing.T) {
    repo := new(mockRepo)
    app := fiber.New()
    h := NewHandler(repo)
    app.Post("/api/auth/logout", h.Logout)

    resp, _ := app.Test(httptest.NewRequest("POST", "/api/auth/logout", nil), 1000)
    assert.Equal(t, 200, resp.StatusCode)
}
```

- [ ] **Step 6: Run auth tests**

```bash
go test ./internal/auth/ -v -count=1 -race
```

- [ ] **Step 7: Commit**

```bash
git add backend/internal/auth/
git commit -m "refactor(auth): extract Repository interface, wire handler tests"
```

---

### Task 5: Hub repository + handler refactoring + tests

- [ ] **Step 1: Create `internal/hub/repository.go`**

```go
package hub

import "github.com/narinsak-u/backend/internal/models"

type Repository interface {
    FindAll() ([]models.Hub, error)
    FindByID(id string) (*models.Hub, error)
    Create(hub *models.Hub) error
    Save(hub *models.Hub) error
    Delete(id string) error
}
```

- [ ] **Step 2: Create `internal/hub/gorm_repository.go`** (see auth/GormRepository pattern — same structure, uses `r.db` for all GORM operations. `Create` calls `r.generateHubID()` if `hub.ID == ""`)

- [ ] **Step 3: Refactor `internal/hub/handler.go`**

Add Handler struct + constructor (same pattern as auth). Convert functions to methods. Replace `database.DB` with `h.repo` calls. Remove `database` import.

- [ ] **Step 4: Create `internal/hub/handler_test.go`** (same pattern as auth test: mockRepo struct implementing Repository, test each endpoint)

- [ ] **Step 5: Build and test**

```bash
go build ./internal/hub/ && go test ./internal/hub/ -v -count=1 -race
```

- [ ] **Step 6: Commit**

```bash
git add backend/internal/hub/
git commit -m "refactor(hub): extract Repository interface, wire handler tests"
```

---

### Task 6: Shipment repository + handler refactoring + tests

- [ ] **Step 1: Create `internal/shipment/repository.go`**

```go
package shipment

import "github.com/narinsak-u/backend/internal/models"

type ShipmentFilter struct {
    Page, Limit    int
    Search, Status, ExcludeStatus string
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
    List(filter ShipmentFilter) ([]models.Shipment, int64, error)
    FindByOrderID(orderID string) (*models.Shipment, error)
    FindByTrackingNumber(trackingNumber string) (*models.Shipment, error)
    Create(shipment *models.Shipment) error
    Save(shipment *models.Shipment) error
    Delete(shipment *models.Shipment) error
    CreateEvent(event *models.ShipmentEvent) error
    FindEventsByShipmentID(shipmentID uint) ([]models.ShipmentEvent, error)
    DeleteShipmentEvents(shipmentID uint) error
    Count() (int64, error)
    CountActive() (int64, error)
    CountByStatus() ([]StatusCountResult, error)
    CountByMonth() ([]MonthCountResult, error)
    CountByDayOfWeek() ([]DayCountResult, error)
}
```

- [ ] **Step 2: Create `internal/shipment/gorm_repository.go`**

Implement all Repository methods using `r.db` (*gorm.DB). Move `generateTrackingNumber()` and `generateOrderID()` (with mutex lock) into this file as private methods. `Create()` calls `r.generateOrderID()` and `r.generateTrackingNumber()` before `r.db.Create(shipment)`.

- [ ] **Step 3: Refactor `internal/shipment/handler.go`**

Add `HubRepository` interface and `Handler` struct at top:

```go
// HubRepository is the subset of hub operations needed for status updates.
type HubRepository interface {
    FindByID(id string) (*models.Hub, error)
}

type Handler struct {
    repo    Repository
    hubRepo HubRepository
}

func NewHandler(repo Repository, hubRepo HubRepository) *Handler {
    return &Handler{repo: repo, hubRepo: hubRepo}
}
```

Convert all functions to methods. Replace database.DB calls:
- **List**: `h.repo.List(filter)` instead of raw GORM query
- **Create**: `h.repo.Create(&shipment)` — OrderID/TrackingNumber are set by repo.Create
- **GetByID**: `h.repo.FindByOrderID(orderID)`
- **UpdateStatus**: `h.repo.FindByOrderID(orderID)`, `h.hubRepo.FindByID(body.HubID)`, `h.repo.Save(&shipment)`, `h.repo.CreateEvent(&event)`
- **Update**: `h.repo.FindByOrderID(orderID)`, `h.repo.Save(&shipment)`
- **Delete**: `h.repo.FindByOrderID(orderID)`, `h.repo.DeleteShipmentEvents(shipment.ID)`, `h.repo.Delete(&shipment)`

Remove `database` import and `generateTrackingNumber()` function from handler.go. Remove `utils.GenerateOrderID()` from Create body.

- [ ] **Step 4: Create `internal/shipment/handler_test.go`**

Create mockRepo (implements Repository) and mockHubRepo (implements HubRepository). Test each handler method (List, Create, GetByID, Update, UpdateStatus, Delete) with success and error cases. Same pattern as auth tests.

- [ ] **Step 5: Build and test**

```bash
go build ./internal/shipment/ && go test ./internal/shipment/ -v -count=1 -race
```

- [ ] **Step 6: Commit**

```bash
git add backend/internal/shipment/
git commit -m "refactor(shipment): extract Repository interface, wire handler tests"
```

---

### Task 7: Tracking handler refactoring + tests

- [ ] **Step 1: Refactor `internal/tracking/handler.go`**

Add Handler struct + constructor that accepts `shipment.Repository`:

```go
type Handler struct {
    repo shipment.Repository
}

func NewHandler(repo shipment.Repository) *Handler {
    return &Handler{repo: repo}
}
```

Convert `Track` to method. Replace:
```go
database.DB.Where("tracking_number = ?", trackingNumber).First(&shipment)
```
with:
```go
shipment, err := h.repo.FindByTrackingNumber(trackingNumber)
if err != nil { return utils.Error(c, 404, "shipment not found") }
```

Replace:
```go
database.DB.Where("shipment_id = ?", shipment.ID).Order("created_at asc").Find(&events)
```
with:
```go
events, err := h.repo.FindEventsByShipmentID(shipment.ID)
if err != nil { return utils.Error(c, 500, "failed to fetch events") }
```

Remove `database` import. Add `shipment` import.

- [ ] **Step 2: Create `internal/tracking/handler_test.go`**

```go
package tracking

import (
    "encoding/json"
    "net/http/httptest"
    "testing"
    "github.com/gofiber/fiber/v2"
    "github.com/narinsak-u/backend/internal/models"
    "github.com/narinsak-u/backend/internal/shipment"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "gorm.io/gorm"
)

type mockRepo struct{ mock.Mock }

func (m *mockRepo) List(f shipment.ShipmentFilter) ([]models.Shipment, int64, error) {
    args := m.Called(f); return args.Get(0).([]models.Shipment), args.Get(1).(int64), args.Error(2)
}
func (m *mockRepo) FindByOrderID(id string) (*models.Shipment, error) {
    args := m.Called(id)
    if args.Get(0) == nil { return nil, args.Error(1) }
    return args.Get(0).(*models.Shipment), args.Error(1)
}
func (m *mockRepo) FindByTrackingNumber(tn string) (*models.Shipment, error) {
    args := m.Called(tn)
    if args.Get(0) == nil { return nil, args.Error(1) }
    return args.Get(0).(*models.Shipment), args.Error(1)
}
func (m *mockRepo) Create(*models.Shipment) error { return m.Called().Error(0) }
func (m *mockRepo) Save(*models.Shipment) error { return m.Called().Error(0) }
func (m *mockRepo) Delete(*models.Shipment) error { return m.Called().Error(0) }
func (m *mockRepo) CreateEvent(*models.ShipmentEvent) error { return m.Called().Error(0) }
func (m *mockRepo) FindEventsByShipmentID(id uint) ([]models.ShipmentEvent, error) {
    args := m.Called(id)
    if args.Get(0) == nil { return nil, args.Error(1) }
    return args.Get(0).([]models.ShipmentEvent), args.Error(1)
}
func (m *mockRepo) DeleteShipmentEvents(uint) error { return m.Called().Error(0) }
func (m *mockRepo) Count() (int64, error) { args := m.Called(); return args.Get(0).(int64), args.Error(1) }
func (m *mockRepo) CountActive() (int64, error) { args := m.Called(); return args.Get(0).(int64), args.Error(1) }
func (m *mockRepo) CountByStatus() ([]shipment.StatusCountResult, error) {
    args := m.Called()
    if args.Get(0) == nil { return nil, args.Error(1) }
    return args.Get(0).([]shipment.StatusCountResult), args.Error(1)
}
func (m *mockRepo) CountByMonth() ([]shipment.MonthCountResult, error) {
    args := m.Called()
    if args.Get(0) == nil { return nil, args.Error(1) }
    return args.Get(0).([]shipment.MonthCountResult), args.Error(1)
}
func (m *mockRepo) CountByDayOfWeek() ([]shipment.DayCountResult, error) {
    args := m.Called()
    if args.Get(0) == nil { return nil, args.Error(1) }
    return args.Get(0).([]shipment.DayCountResult), args.Error(1)
}

func TestTrack_Found(t *testing.T) {
    repo := new(mockRepo)
    repo.On("FindByTrackingNumber", "TH202600001").Return(&models.Shipment{
        OrderID: "ORD-10245", TrackingNumber: "TH202600001", Status: "in_transit",
    }, nil)
    repo.On("FindEventsByShipmentID", uint(0)).Return([]models.ShipmentEvent{
        {Status: "Label Created", Description: "Awaiting pickup."},
    }, nil)

    app := fiber.New()
    h := NewHandler(repo)
    app.Get("/api/track/:trackingNumber", h.Track)

    resp, _ := app.Test(httptest.NewRequest("GET", "/api/track/TH202600001", nil), 1000)
    assert.Equal(t, 200, resp.StatusCode)

    var result map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&result)
    assert.True(t, result["success"].(bool))
    data := result["data"].(map[string]interface{})
    assert.NotNil(t, data["shipment"])
    assert.NotNil(t, data["events"])
}

func TestTrack_NotFound(t *testing.T) {
    repo := new(mockRepo)
    repo.On("FindByTrackingNumber", "INVALID").Return(nil, gorm.ErrRecordNotFound)

    app := fiber.New()
    h := NewHandler(repo)
    app.Get("/api/track/:trackingNumber", h.Track)

    resp, _ := app.Test(httptest.NewRequest("GET", "/api/track/INVALID", nil), 1000)
    assert.Equal(t, 404, resp.StatusCode)
}
```

- [ ] **Step 3: Build and test**

```bash
go build ./internal/tracking/ && go test ./internal/tracking/ -v -count=1 -race
```

- [ ] **Step 4: Commit**

```bash
git add backend/internal/tracking/
git commit -m "refactor(tracking): inject shipment.Repository, add handler tests"
```

---

### Task 8: Analytics handler refactoring + tests

- [ ] **Step 1: Refactor `internal/analytics/handler.go`**

Add Handler struct + constructor that accepts `shipment.Repository`:

```go
type Handler struct {
    repo shipment.Repository
}

func NewHandler(repo shipment.Repository) *Handler {
    return &Handler{repo: repo}
}
```

Convert functions to methods. Replace all `database.DB` analytics queries with repo calls:
- `database.DB.Model(&models.Shipment{}).Count(&total)` -> `total, err := h.repo.Count()`
- `Where("status NOT IN ?", ...).Count(&active)` -> `active, err := h.repo.CountActive()`
- `Select("status, count(*) as count").Group("status").Scan(&byStatus)` -> `byStatus, err := h.repo.CountByStatus()` (the result type is already `[]shipment.StatusCountResult`)
- Same pattern for `CountByMonth` and `CountByDayOfWeek`

Remove `database` import. Add `shipment` and remove unused `models` import if possible.

Region mapping logic stays in the handler.

- [ ] **Step 2: Create `internal/analytics/handler_test.go`**

```go
package analytics

import (
    "encoding/json"
    "net/http/httptest"
    "testing"
    "github.com/gofiber/fiber/v2"
    "github.com/narinsak-u/backend/internal/models"
    "github.com/narinsak-u/backend/internal/shipment"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

type mockRepo struct{ mock.Mock }

func (m *mockRepo) List(f shipment.ShipmentFilter) ([]models.Shipment, int64, error) {
    args := m.Called(f); return args.Get(0).([]models.Shipment), args.Get(1).(int64), args.Error(2)
}
func (m *mockRepo) FindByOrderID(id string) (*models.Shipment, error) {
    args := m.Called(id)
    if args.Get(0) == nil { return nil, args.Error(1) }
    return args.Get(0).(*models.Shipment), args.Error(1)
}
func (m *mockRepo) FindByTrackingNumber(tn string) (*models.Shipment, error) {
    args := m.Called(tn)
    if args.Get(0) == nil { return nil, args.Error(1) }
    return args.Get(0).(*models.Shipment), args.Error(1)
}
func (m *mockRepo) Create(*models.Shipment) error { return m.Called().Error(0) }
func (m *mockRepo) Save(*models.Shipment) error { return m.Called().Error(0) }
func (m *mockRepo) Delete(*models.Shipment) error { return m.Called().Error(0) }
func (m *mockRepo) CreateEvent(*models.ShipmentEvent) error { return m.Called().Error(0) }
func (m *mockRepo) FindEventsByShipmentID(id uint) ([]models.ShipmentEvent, error) {
    args := m.Called(id)
    if args.Get(0) == nil { return nil, args.Error(1) }
    return args.Get(0).([]models.ShipmentEvent), args.Error(1)
}
func (m *mockRepo) DeleteShipmentEvents(uint) error { return m.Called().Error(0) }
func (m *mockRepo) Count() (int64, error) { args := m.Called(); return args.Get(0).(int64), args.Error(1) }
func (m *mockRepo) CountActive() (int64, error) { args := m.Called(); return args.Get(0).(int64), args.Error(1) }
func (m *mockRepo) CountByStatus() ([]shipment.StatusCountResult, error) {
    args := m.Called()
    if args.Get(0) == nil { return nil, args.Error(1) }
    return args.Get(0).([]shipment.StatusCountResult), args.Error(1)
}
func (m *mockRepo) CountByMonth() ([]shipment.MonthCountResult, error) {
    args := m.Called()
    if args.Get(0) == nil { return nil, args.Error(1) }
    return args.Get(0).([]shipment.MonthCountResult), args.Error(1)
}
func (m *mockRepo) CountByDayOfWeek() ([]shipment.DayCountResult, error) {
    args := m.Called()
    if args.Get(0) == nil { return nil, args.Error(1) }
    return args.Get(0).([]shipment.DayCountResult), args.Error(1)
}

func TestOverview(t *testing.T) {
    repo := new(mockRepo)
    repo.On("Count").Return(int64(10), nil)
    repo.On("CountActive").Return(int64(7), nil)
    repo.On("CountByStatus").Return([]shipment.StatusCountResult{
        {Status: "in_transit", Count: 4},
        {Status: "pending", Count: 3},
    }, nil)

    app := fiber.New()
    h := NewHandler(repo)
    app.Get("/api/analytics/overview", h.Overview)

    resp, _ := app.Test(httptest.NewRequest("GET", "/api/analytics/overview", nil), 1000)
    assert.Equal(t, 200, resp.StatusCode)

    var result map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&result)
    assert.True(t, result["success"].(bool))
    data := result["data"].(map[string]interface{})
    assert.Equal(t, float64(10), data["total"])
    assert.Equal(t, float64(7), data["active"])
}

func TestTimeSeries(t *testing.T) {
    repo := new(mockRepo)
    repo.On("CountByMonth").Return([]shipment.MonthCountResult{
        {Month: "2026-01", Count: 3},
        {Month: "2026-02", Count: 5},
    }, nil)
    repo.On("CountByDayOfWeek").Return([]shipment.DayCountResult{
        {Day: "Monday", Count: 2},
        {Day: "Wednesday", Count: 4},
    }, nil)

    app := fiber.New()
    h := NewHandler(repo)
    app.Get("/api/analytics/timeseries", h.TimeSeries)

    resp, _ := app.Test(httptest.NewRequest("GET", "/api/analytics/timeseries", nil), 1000)
    assert.Equal(t, 200, resp.StatusCode)

    var result map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&result)
    assert.True(t, result["success"].(bool))
}
```

- [ ] **Step 3: Build and test**

```bash
go build ./internal/analytics/ && go test ./internal/analytics/ -v -count=1 -race
```

- [ ] **Step 4: Commit**

```bash
git add backend/internal/analytics/
git commit -m "refactor(analytics): inject shipment.Repository, add handler tests"
```

---

### Task 9: Wire main.go and delete idgen.go

- [ ] **Step 1: Update `cmd/server/main.go`**

Replace direct handler references with constructor calls:

```go
// Before:
// auth.Register
// After:
authRepo := auth.NewGormRepository(database.DB)
authHandler := auth.NewHandler(authRepo)
authGroup.Post("/register", middleware.RateLimitAuth(), authHandler.Register)
authGroup.Post("/login", middleware.RateLimitAuth(), authHandler.Login)
authGroup.Get("/me", middleware.AuthRequired(), authHandler.Me)
authGroup.Post("/logout", authHandler.Logout)
```

Same pattern for shipment, hub, tracking, and analytics. Remove the `pkg/utils` import if no longer needed.

- [ ] **Step 2: Delete `pkg/utils/idgen.go`**

```bash
git rm backend/pkg/utils/idgen.go
```

- [ ] **Step 3: Build everything**

```bash
go build ./...
```

Expected: Clean build.

- [ ] **Step 4: Vet**

```bash
go vet ./...
```

- [ ] **Step 5: Run all tests**

```bash
go test ./... -count=1 -race
```

Expected: All tests pass (with possible timing-dependent rate limit test needing a small adjustment).

- [ ] **Step 6: Commit**

```bash
git add backend/cmd/server/main.go
git add backend/pkg/utils/idgen.go  # staged for deletion via git rm
git commit -m "feat: wire repository pattern in main.go, remove obsolete idgen.go"
```

---

### Task 10: Self-review and final verification

- [ ] **Step 1: Check spec coverage**

Verify the plan covers every requirement from the design spec:
- Repository interfaces per domain: tasks 4, 5, 6
- GORM implementations: tasks 4, 5, 6
- Handler struct refactoring: tasks 4, 5, 6, 7, 8
- Unit tests for pure logic: tasks 2, 3
- HTTP integration tests: tasks 4, 5, 6, 7, 8
- Main.go wiring + idgen.go deletion: task 9
- testify dependency: task 1

- [ ] **Step 2: Run final full test suite**

```bash
go test ./... -count=1 -race -v 2>&1 | tail -50
```

Check for any failures, panics, or race conditions.

- [ ] **Step 3: Check coverage**

```bash
go test ./... -count=1 -coverprofile=coverage.out && go tool cover -func=coverage.out | tail -20
```

Expected: 70-90% coverage across packages (lower in seed/database packages since they need a real DB).
