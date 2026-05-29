# Backend Overview

Go API server for the Thun-u-der Express shipment tracking platform. Built with **Fiber v2**, **GORM** (PostgreSQL), and **JWT** authentication.

---

## Tech Stack

| Component        | Library                                      |
| ---------------- | -------------------------------------------- |
| HTTP Framework   | `gofiber/fiber/v2` (FastHTTP)                |
| ORM              | `gorm.io/gorm` + `gorm.io/driver/postgres`   |
| Auth             | `golang-jwt/jwt/v5` (HS256)                  |
| Password Hashing | `golang.org/x/crypto` (bcrypt)               |
| Logging          | `rs/zerolog`                                 |
| Config           | `joho/godotenv`                              |
| Testing          | `stretchr/testify` (assert, require, mock)   |

---

## Directory Structure

```
backend/
├── cmd/server/main.go          # Entry point: config load, DB init, migrate, seed, route setup, server start
├── internal/
│   ├── config/config.go        # Global Config struct (Port, DatabaseURL, JWTSecret, JWTTTL)
│   ├── data/regions.go         # Static Thailand province-to-region mapping
│   ├── database/
│   │   └── postgres.go         # GORM connection (global DB var, 15x retry)
│   ├── seed/
│   │   ├── hubs.go             # SeedHubs — inserts 10 demo hubs (idempotent)
│   │   └── shipments.go        # SeedShipments — inserts 20 demo shipments + events (idempotent)
│   ├── models/
│   │   ├── user.go             # User model
│   │   ├── shipment.go         # Shipment + ContactInfo + GeoPoint (GORM hooks for coords)
│   │   ├── shipment_event.go   # ShipmentEvent + Location
│   │   └── hub.go              # Hub model
│   ├── middleware/
│   │   ├── auth.go             # JWT Bearer token / cookie verification
│   │   ├── cors.go             # Manual CORS headers (origin: localhost:5173)
│   │   ├── logger.go           # Request logging via zerolog
│   │   └── ratelimit.go        # In-memory sliding-window rate limiter (5 req/min on auth)
│   ├── auth/
│   │   ├── handler.go          # Handler struct (Register, Login, Me, Logout)
│   │   ├── repository.go       # Repository interface
│   │   └── gorm_repository.go  # GORM implementation
│   ├── shipment/
│   │   ├── handler.go          # Handler struct (List, Create, GetByID, Update, UpdateStatus, Delete)
│   │   ├── repository.go       # Repository interface + filter/result types
│   │   └── gorm_repository.go  # GORM implementation (CRUD, events, analytics queries)
│   ├── hub/
│   │   ├── handler.go          # Handler struct (List, GetByID, Create, Update, Delete)
│   │   ├── repository.go       # Repository interface
│   │   └── gorm_repository.go  # GORM implementation (CRUD + ID generation)
│   ├── tracking/handler.go     # Handler struct (Track) — uses shipment.Repository
│   └── analytics/handler.go    # Handler struct (Overview, TimeSeries) — uses shipment.Repository
├── pkg/utils/
│   ├── hash.go                 # bcrypt hash/compare helpers
│   └── response.go             # Standard JSON response writers (Success/Error/SuccessWithPagination)
├── Dockerfile                  # Multi-stage build (golang:1.24-alpine → alpine:3.19)
├── docker-compose.yml          # Postgres 16 + backend
├── .env.example
├── .dockerignore
├── go.mod / go.sum
├── README.md
└── docs/
    ├── OVERVIEW.md
    └── WORKFLOW.md
```

---

## API Endpoints

| Method | Path                              | Auth  | Handler                      | Description                          |
| ------ | --------------------------------- | ----- | ---------------------------- | ------------------------------------ |
| POST   | `/api/auth/register`              | No*   | `auth.Handler.Register`      | Register a new user                  |
| POST   | `/api/auth/login`                 | No*   | `auth.Handler.Login`         | Login, returns JWT + sets cookie     |
| GET    | `/api/auth/me`                    | JWT   | `auth.Handler.Me`            | Get current user profile             |
| POST   | `/api/auth/logout`                | No    | `auth.Handler.Logout`        | Clear auth cookie                    |
| GET    | `/api/shipments`                  | No    | `shipment.Handler.List`      | List shipments (pagination, search, filter) |
| GET    | `/api/shipments/:orderId`         | No    | `shipment.Handler.GetByID`   | Get shipment by order ID             |
| POST   | `/api/shipments`                  | JWT   | `shipment.Handler.Create`    | Create a new shipment                |
| PATCH  | `/api/shipments/:orderId/status`  | JWT   | `shipment.Handler.UpdateStatus` | Update shipment status + log event |
| PUT    | `/api/shipments/:orderId`         | JWT   | `shipment.Handler.Update`    | Update shipment fields               |
| DELETE | `/api/shipments/:orderId`         | JWT   | `shipment.Handler.Delete`    | Delete a shipment + its events       |
| GET    | `/api/track/:trackingNumber`      | No    | `tracking.Handler.Track`     | Public tracking lookup               |
| GET    | `/api/hubs`                       | No    | `hub.Handler.List`           | List all hubs                        |
| GET    | `/api/hubs/:id`                   | No    | `hub.Handler.GetByID`        | Get hub by ID                        |
| POST   | `/api/hubs`                       | JWT   | `hub.Handler.Create`         | Create a new hub                     |
| PUT    | `/api/hubs/:id`                   | JWT   | `hub.Handler.Update`         | Update hub fields                    |
| DELETE | `/api/hubs/:id`                   | JWT   | `hub.Handler.Delete`         | Delete a hub                         |
| GET    | `/api/analytics/overview`         | JWT   | `analytics.Handler.Overview` | Dashboard aggregate stats            |
| GET    | `/api/analytics/timeseries`       | JWT   | `analytics.Handler.TimeSeries` | Monthly + day-of-week trends       |

\* Rate-limited to 5 requests/minute per IP.

**Auth notes:** Read operations (GET shipments, hubs, tracking) are **public**. Write operations (POST, PUT, PATCH, DELETE) require JWT.

**Pagination (shipments List):** Query params `page` (default 1), `limit` (default 10, use `-1` for all), `search` (ILIKE on order_id, tracking_number, customer_name, destination), `status` (filter by status), `exclude_status` (exclude by status).

---

## Data Model

```
User (standalone)
  - ID, Name, Email, Password (bcrypt, hidden from JSON), Role, CreatedAt

Shipment 1──N ShipmentEvent
  - Shipment: ID (uint PK), OrderID (unique, e.g. ORD-10245),
    TrackingNumber (unique, e.g. TH202612345), Customer (ContactInfo embedded),
    Receiver (ContactInfo embedded), Origin, Destination, CurrentCoords,
    Status, Carrier, HubID, Weight, Items, EstimatedDelivery,
    CreatedAt, Progress
  - ShipmentEvent: ID, ShipmentID (FK), Status, Location (embedded),
    Description, CreatedAt (returned as "timestamp")

Hub (standalone)
  - ID (string PK, e.g. HUB-001), Name, CarrierID (indexed), Address,
    Coords (GeoPoint), Capacity, CurrentUtilization, Status, CreatedAt
```

### Shared Types

| Type        | Fields                                      | Used In                                    |
| ----------- | ------------------------------------------- | ------------------------------------------ |
| `GeoPoint`  | `lat`, `lng`                                | `ContactInfo.Coords`, `Shipment.CurrentCoords`, `Hub.Coords` |
| `ContactInfo` | `name`, `zipcode`, `subDistrict`, `district`, `province`, `coords` | `Shipment.Customer`, `Shipment.Receiver` |
| `Location`  | `name`, `lat`, `lng`                        | `ShipmentEvent.Location`                   |

### Coords Storage Pattern

`GeoPoint` fields use `gorm:"-"` (not stored directly). Values are synced to flat `_lat`/`_lng` columns via GORM hooks:
- **`BeforeSave`** — copies `Coords.Lat/Lng` into flat columns before insert/update
- **`AfterFind`** — reconstructs `Coords` from flat columns after query

This gives clean JSON with nested objects while keeping normal DB columns for querying.

---

## Architecture Notes

- **Repository pattern:** Each domain package (auth, hub, shipment) defines a `Repository` interface and a GORM implementation. Handlers receive their repository via constructor injection (`NewHandler(repo)`), making them testable without a real database.
- **Dependency injection:** `main.go` wires GORM repositories → handlers → routes. No packages call `database.DB` directly except repositories and seed code.
- **Cross-domain dependency:** The shipment handler receives both a `shipment.Repository` and a `HubRepository` interface (defined locally in the shipment package to avoid circular imports) for status updates that need hub location data.
- **Standardized responses:** All endpoints return `{"success": true/false, "data": ...}` or `{"success": false, "error": "..."}`. Paginated endpoints include a `pagination` block.
- **JWT auth:** HS256 bearer tokens with `user_id` and `role` claims. Token TTL is hardcoded to 24h. Auth middleware checks `Authorization: Bearer` header first, then falls back to the `jwt` cookie. Auth endpoints are rate-limited to 5 requests/minute per IP via an in-memory sliding-window limiter.
- **Seed data:** On first startup, `main.go` runs `seed.SeedHubs()` and `seed.SeedShipments()` after AutoMigrate. Both skip if their table already has rows (idempotent).
- **Testing:** 19 test files, 64+ tests across 10 packages using testify (assert, require, mock). Handler tests use mocked repositories with Fiber's `app.Test()`. Pure-logic tests cover models (GORM hooks), middleware (rate limiter, JWT auth), utils, config, and data.
- **Middlewares run in order:** CORS → Logger → Rate limit (auth routes only) → Auth (protected routes only).

---

## Configuration

| Variable       | Default                                                      | Description                |
| -------------- | ------------------------------------------------------------ | -------------------------- |
| `PORT`         | `8080`                                                       | Server listen port         |
| `DATABASE_URL` | `postgres://user:pass@localhost:5432/shipments`              | PostgreSQL DSN             |
| `JWT_SECRET`   | `change-me`                                                  | HMAC secret for JWT        |

---

## Running Locally

```bash
# With Docker Compose (Postgres + backend)
cd backend
docker compose up

# Without Docker
cd backend
go run .            # requires separate Postgres
```

---

## Current State

| Feature | Status |
|---------|--------|
| Auth (register, login, me, logout) | Done |
| Shipment CRUD | Done (List with pagination/search/filter) |
| Shipment status updates with hub coords | Done (sets currentCoords to hub location) |
| Hub CRUD | Done |
| Public tracking lookup | Done |
| Analytics overview + timeseries | Done |
| Rate limiting on auth endpoints | Done (5 req/min per IP) |
| Repository pattern with DI | Done (auth, hub, shipment) |
| Tests | Done (~64 tests, 46% overall coverage, handler logic 70-100%) |
