# Shipment Code Flow (main.go → Handler → Repository → GORM)

This guide traces how a request travels through the codebase **end-to-end**. If you're new to Go, read this top-to-bottom — each section builds on the last.

---

## 1. The Big Picture

```
 Browser/CLI
     │
     ▼
  main.go  ←  registers routes, wires up dependencies
     │
     ▼
  Handler  ←  parses HTTP request, calls repository
     │
     ▼
  Repository (interface)  ←  contract ("what to do")
     │
     ▼
  GormRepository  ←  real implementation ("how to do it")
     │
     ▼
  Postgres  ←  the database
```

**Key insight:** The handler never touches the database directly. It talks to an **interface** (`Repository`). The real database code lives in `GormRepository`. This separation lets you test the handler with a mock database.

---

## 2. Step 1 — main.go Wires Everything Together

In `cmd/server/main.go`:

```go
// 1. Create the repository (real GORM implementation)
shipmentRepo := shipment.NewGormRepository(database.DB)

// 2. Create the hub repository (needed by UpdateStatus)
hubRepo := hub.NewGormRepository(database.DB)

// 3. Create the handler, passing the repository in
shipmentHandler := shipment.NewHandler(shipmentRepo, hubRepo)

// 4. Register routes — the handler methods handle requests
api.Get("/shipments", shipmentHandler.List)             // public
api.Get("/shipments/:orderId", shipmentHandler.GetByID)  // public
shipmentGroup := api.Group("/shipments", middleware.AuthRequired())
shipmentGroup.Post("/", shipmentHandler.Create)
```

This is called **dependency injection**: instead of the Handler creating its own database connection, main.go creates everything and "injects" the dependencies.

---

## 3. Step 2 — A Request Arrives at the Handler

When a client sends `POST /api/shipments`, Fiber routes it to `shipmentHandler.Create`.

In `handler.go`:

```go
func (h *Handler) Create(c *fiber.Ctx) error {
    // h.repo is a Repository (set by NewHandler in main.go)

    // Parse JSON from the request body
    var req CreateRequest
    if err := c.BodyParser(&req); err != nil {
        return utils.Error(c, 400, "invalid request body")
    }

    // Validate required fields
    if req.Customer.Name == "" || ... {
        return utils.Error(c, 400, "customer name is required")
    }

    // Build a Shipment model from the request
    shipment := models.Shipment{
        Customer:  req.Customer,
        Status:    "pending",
        Carrier:   req.Carrier,
        // ...
    }

    // Call the repository — at this point we don't know if
    // it's a real database or a test mock
    if err := h.repo.Create(&shipment); err != nil {
        return utils.Error(c, 500, "failed to create shipment")
    }

    return utils.Success(c, shipment)
}
```

**What the handler does:**
1. Parse the HTTP request (JSON body, URL params, query strings)
2. Validate the input
3. Call the repository
4. Return an HTTP response (success or error)

The handler does **not** know about Postgres, GORM, or SQL. It just calls `h.repo.Create()`.

---

## 4. Step 3 — The Repository Interface (The Contract)

In `repository.go`:

```go
type Repository interface {
    Create(shipment *models.Shipment) error
    FindByOrderID(orderID string) (*models.Shipment, error)
    Save(shipment *models.Shipment) error
    // ... more methods
}
```

An **interface** in Go is a **contract**: it says "any type that has these methods can be used as a Repository."

The Handler depends on the **interface**, not the concrete implementation. This is why the handler doesn't need to change if you swap the database — it just needs something that satisfies `Repository`.

Think of it like a USB port: your laptop (the handler) doesn't care what's on the other end (real hard drive vs. external SSD vs. thumb drive) as long as it fits the USB shape (the interface).

---

## 5. Step 4 — GormRepository Does the Actual Work

In `gorm_repository.go`:

```go
type GormRepository struct {
    db *gorm.DB  // real database connection
}

func (r *GormRepository) Create(shipment *models.Shipment) error {
    shipment.OrderID = r.generateOrderID()       // "ORD-10246"
    shipment.TrackingNumber = r.generateTrackingNumber() // "TH2026..."
    return r.db.Create(shipment).Error            // INSERT INTO Postgres
}
```

This is where the real database operations happen:
- `r.db` is a *gorm.DB connection set up in `database.ConnectPostgres()`
- `r.db.Create(shipment)` runs `INSERT INTO shipments ...`
- GORM converts the Go struct into SQL automatically

---

## 6. The Full Flow — Concrete Example

Let's trace `POST /api/shipments` with this JSON body:

```json
{
    "customer": { "name": "Somchai", "province": "Bangkok", "zipcode": "10100" },
    "receiver": { "name": "Somsri", "province": "Chonburi", "zipcode": "20000" },
    "carrier": "Thun-u-der Express",
    "weight": 2.5,
    "items": 1
}
```

### Step-by-step:

```
[1] Fiber receives POST /api/shipments
    │
    ├─ AuthRequired middleware checks JWT (valid? → continue)
    │
    ▼
[2] shipmentHandler.Create(c)  ← handler.go
    │
    ├─ c.BodyParser(&req) → fills CreateRequest struct
    ├─ Validate required fields (name, province, zipcode, weight, items, carrier)
    ├─ Build Shipment model:
    │     Customer      = { Somchai, Bangkok, ... }
    │     Receiver      = { Somsri, Chonburi, ... }
    │     Origin        = "Somchai, Bangkok, 10100"
    │     Destination   = "Somsri, Chonburi, 20000"
    │     Status        = "pending"
    │     Carrier       = "Thun-u-der Express"
    │     Weight        = 2.5
    │     Items         = 1
    │     EstimatedDelivery = now + 72h
    │     Progress      = 0
    │
    ├─ h.repo.Create(&shipment)  ← calls interface
    │
    ▼
[3] GormRepository.Create(shipment)  ← gorm_repository.go
    │
    ├─ shipment.OrderID        = generateOrderID()     → "ORD-10246"
    ├─ shipment.TrackingNumber = generateTrackingNumber() → "TH202612345"
    ├─ GORM BeforeSave hook fires → coords → flat columns
    ├─ r.db.Create(shipment)   → INSERT INTO shipments ...
    │
    ▼
[4] Postgres
    │
    ├─ New row inserted with all fields
    ├─ ID (auto-increment PK) generated by Postgres
    │
    ▼
[5] Back in handler.Create(c)
    │
    ├─ h.repo.CreateEvent(&event)  ← insert "Label Created" tracking event
    │
    ▼
[6] Return JSON response → 200 { shipment: { ... } }
```

---

## 7. Why Three Layers? (Handler vs Interface vs GORM)

This separation exists for one reason: **testability**.

### Without the interface

If the handler called GORM directly, you'd need a real Postgres database to test:

```go
// ❌ Hard to test — needs real database
func (h *Handler) Create(c *fiber.Ctx) error {
    db := gorm.Open(...)  // needs Postgres!
    db.Create(&shipment)
}
```

### With the interface

Tests can pass a **mock repository**:

```go
// ✅ Easy to test — no database needed
type mockRepo struct{}
func (m *mockRepo) Create(s *models.Shipment) error {
    s.OrderID = "ORD-TEST-001"  // fake it
    return nil
}

func TestCreate(t *testing.T) {
    handler := NewHandler(&mockRepo{}, nil)
    // ... test the handler logic without any database
}
```

---

## 8. Visual Summary

```
┌─────────────────────────────────────────────────────────────┐
│                     main.go                                  │
│                                                              │
│  shipmentRepo := NewGormRepository(database.DB)              │
│  hubRepo      := hub.NewGormRepository(database.DB)          │
│  handler      := NewHandler(shipmentRepo, hubRepo)           │
│                                                              │
│  app.Post("/api/shipments", handler.Create)                  │
└──────────────────────┬──────────────────────────────────────┘
                       │
                       ▼ request comes in
┌─────────────────────────────────────────────────────────────┐
│                  Handler (handler.go)                        │
│                                                              │
│  Fields:  repo  (Repository interface)                      │
│           hubRepo (HubRepository interface)                  │
│                                                              │
│  Job: Parse request → validate → repo.Create() → respond    │
│                                                              │
│  Knows: HTTP, JSON, validation rules                         │
│  Knows NOT: SQL, GORM, Postgres                              │
└──────────────────────┬──────────────────────────────────────┘
                       │
                       ▼ h.repo.Create(&shipment)
┌─────────────────────────────────────────────────────────────┐
│              Repository Interface (repository.go)             │
│                                                              │
│  type Repository interface {                                 │
│      Create(shipment *models.Shipment) error                 │
│      FindByOrderID(orderID string) (*models.Shipment, error) │
│      ...                                                     │
│  }                                                           │
│                                                              │
│  Job: Define the contract (what methods exist)               │
└──────────────────────┬──────────────────────────────────────┘
                       │
                       ▼ GormRepository satisfies the interface
┌─────────────────────────────────────────────────────────────┐
│              GormRepository (gorm_repository.go)              │
│                                                              │
│  Fields:  db *gorm.DB                                        │
│                                                              │
│  Job: Run actual SQL queries via GORM                        │
│                                                              │
│  Knows: SQL, Postgres, GORM hooks (BeforeSave, AfterFind)   │
└──────────────────────┬──────────────────────────────────────┘
                       │
                       ▼ r.db.Create(shipment)
┌─────────────────────────────────────────────────────────────┐
│                       Postgres                                │
│                                                              │
│  INSERT INTO shipments (...) VALUES (...)                     │
└─────────────────────────────────────────────────────────────┘
```

---

## 9. Key Go Patterns Used Here

| Pattern | Where | Why |
|---------|-------|-----|
| **Interface as contract** | `repository.go` | Decouples handler from database |
| **Dependency injection** | `main.go` → `NewHandler(repo, hubRepo)` | main.go decides what implementations to use |
| **Constructor function** | `NewGormRepository(db)` | Ensures the struct is properly initialized |
| **Pointer receiver** | `func (r *GormRepository) Create(...)` | Methods can modify the struct (e.g. generate IDs) |
| **Returning errors** | Every method returns `error` | Go's idiomatic way to handle failures — never exceptions |
| **Small interfaces** | Repository has ~14 methods | Still cohesive — all shipment-related operations |

### Reading the type annotations

```go
func (r *GormRepository) Create(shipment *models.Shipment) error
```

- `(r *GormRepository)` → this method belongs to `*GormRepository` (pointer to GormRepository)
- `Create` → the method name
- `(shipment *models.Shipment)` → takes a pointer to a Shipment (can modify it)
- `error` → returns an error (nil = success)

```go
func (h *Handler) Create(c *fiber.Ctx) error
```

- `(h *Handler)` → belongs to `*Handler`
- `Create` → the method name
- `(c *fiber.Ctx)` → takes the Fiber request context (has params, body, response methods)
- `error` → returns nil on success, or an error (Fiber handles the response)
