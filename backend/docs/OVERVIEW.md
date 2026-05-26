# Backend Overview

Go API server for the Thun-u-der Express shipment tracking platform. Built with **Fiber v2**, **GORM** (PostgreSQL), **go-redis**, and **JWT** authentication.

---

## Tech Stack

| Component        | Library                                      |
| ---------------- | -------------------------------------------- |
| HTTP Framework   | `gofiber/fiber/v2` (FastHTTP)                |
| ORM              | `gorm.io/gorm` + `gorm.io/driver/postgres`   |
| Cache            | `redis/go-redis/v9` (initialized, not yet used) |
| Auth             | `golang-jwt/jwt/v5` (HS256)                  |
| Password Hashing | `golang.org/x/crypto` (bcrypt)               |
| WebSocket        | `gofiber/contrib/websocket` (gorilla/websocket) |
| Logging          | `rs/zerolog`                                 |
| Config           | `joho/godotenv`                              |

---

## Directory Structure

```
backend/
├── cmd/server/main.go          # Entry point: config load, DB init, migrate, seed, route setup, server start
├── internal/
│   ├── config/config.go        # Global Config struct (Port, DatabaseURL, RedisURL, JWTSecret, JWTTTL)
│   ├── database/
│   │   ├── postgres.go         # GORM connection (global DB var, 15x retry)
│   │   └── redis.go            # go-redis client (global Redis var)
│   ├── seed/
│   │   ├── hubs.go             # SeedHubs — inserts 6 demo hubs (idempotent)
│   │   └── shipments.go        # SeedShipments — inserts 12 demo shipments + events (idempotent)
│   ├── models/
│   │   ├── user.go             # User model
│   │   ├── shipment.go         # Shipment + ContactInfo + GeoPoint (GORM hooks for coords)
│   │   ├── shipment_event.go   # ShipmentEvent + Location
│   │   └── hub.go              # Hub model
│   ├── middleware/
│   │   ├── auth.go             # JWT Bearer token / cookie verification
│   │   ├── cors.go             # Manual CORS headers (origin: localhost:5173)
│   │   └── logger.go           # Request logging via zerolog
│   ├── auth/handler.go         # Register, Login, Me, Logout
│   ├── shipment/handler.go     # CRUD + status update for shipments (334 lines)
│   ├── hub/handler.go          # CRUD for logistics hubs (85 lines)
│   ├── tracking/handler.go     # Public tracking lookup (25 lines)
│   ├── analytics/handler.go    # Dashboard aggregate stats (33 lines)
│   └── websocket/
│       ├── hub.go              # Connection hub (room-based broadcast)
│       └── client.go           # WebSocket client + upgrade handler
├── pkg/utils/
│   ├── hash.go                 # bcrypt hash/compare helpers
│   └── response.go             # Standard JSON response writers (Success/Error/SuccessWithPagination)
├── Dockerfile                  # Multi-stage build (golang:1.24-alpine → alpine:3.19)
├── docker-compose.yml          # Postgres 16 + Redis 7 + backend
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
| POST   | `/api/auth/register`              | No    | `auth.Register`              | Register a new user                  |
| POST   | `/api/auth/login`                 | No    | `auth.Login`                 | Login, returns JWT + sets cookie     |
| GET    | `/api/auth/me`                    | JWT   | `auth.Me`                    | Get current user profile             |
| POST   | `/api/auth/logout`                | No    | `auth.Logout`                | Clear auth cookie                    |
| GET    | `/api/shipments`                  | No    | `shipment.List`              | List shipments (pagination, search, filter) |
| GET    | `/api/shipments/:orderId`         | No    | `shipment.GetByID`           | Get shipment by order ID             |
| POST   | `/api/shipments`                  | JWT   | `shipment.Create`            | Create a new shipment                |
| PATCH  | `/api/shipments/:orderId/status`  | JWT   | `shipment.UpdateStatus`      | Update shipment status + log event   |
| PUT    | `/api/shipments/:orderId`         | JWT   | `shipment.Update`            | Update shipment fields               |
| DELETE | `/api/shipments/:orderId`         | JWT   | `shipment.Delete`            | Delete a shipment + its events       |
| GET    | `/api/track/:trackingNumber`      | No    | `tracking.Track`             | Public tracking lookup               |
| GET    | `/api/hubs`                       | No    | `hub.List`                   | List all hubs                        |
| GET    | `/api/hubs/:id`                   | No    | `hub.GetByID`                | Get hub by ID                        |
| POST   | `/api/hubs`                       | JWT   | `hub.Create`                 | Create a new hub                     |
| PUT    | `/api/hubs/:id`                   | JWT   | `hub.Update`                 | Update hub fields                    |
| DELETE | `/api/hubs/:id`                   | JWT   | `hub.Delete`                 | Delete a hub                         |
| GET    | `/api/analytics/overview`         | JWT   | `analytics.Overview`         | Dashboard aggregate stats            |
| GET    | `/ws/tracking/:trackingNumber`    | No    | `websocket.HandleWebSocket`  | Real-time tracking updates           |
| GET    | `/ws/admin`                       | No    | `websocket.HandleWebSocket`  | Admin WebSocket (room "global")      |

**Auth notes:** Read operations (GET shipments, hubs, tracking) are **public**. Write operations (POST, PUT, PATCH, DELETE) require JWT. WebSocket endpoints are public.

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

- **Handler-centric:** Each domain package has a single `handler.go`. Handlers call GORM directly via the global `database.DB` — no service/repository layers.
- **Global state:** `database.DB`, `database.Redis`, and `config.App` are package-level globals (no dependency injection).
- **Standardized responses:** All endpoints return `{"success": true/false, "data": ...}` or `{"success": false, "error": "..."}`. Paginated endpoints include a `pagination` block.
- **JWT auth:** HS256 bearer tokens with `user_id` and `role` claims. Token TTL is hardcoded to 24h. Auth middleware checks `Authorization: Bearer` header first, then falls back to the `jwt` cookie.
- **Seed data:** On first startup, `main.go` runs `seed.SeedHubs()` and `seed.SeedShipments()` after AutoMigrate. Both skip if their table already has rows (idempotent).
- **WebSocket pub/sub:** Room-based hub broadcasts to clients by tracking number or "global" room. Infrastructure is wired up; actual push events are not yet implemented.
- **Redis is initialized but not used** by any handler yet.
- **No tests** exist in the codebase.
- **Middlewares run in order:** CORS → Logger → Auth (only on protected routes).

---

## Configuration

| Variable       | Default                                                      | Description                |
| -------------- | ------------------------------------------------------------ | -------------------------- |
| `PORT`         | `8080`                                                       | Server listen port         |
| `DATABASE_URL` | `postgres://user:pass@localhost:5432/shipments`              | PostgreSQL DSN             |
| `REDIS_URL`    | `redis://localhost:6379`                                     | Redis connection URL       |
| `JWT_SECRET`   | `change-me`                                                  | HMAC secret for JWT        |

---

## Running Locally

```bash
# With Docker Compose (Postgres + Redis + backend)
cd backend
docker compose up

# Without Docker
cd backend
go run .            # requires separate Postgres + Redis
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
| Analytics overview | Done |
| WebSocket infrastructure | Wired, no push triggers yet |
| Redis integration | Initialized, not yet in use |
| WebSocket push on status change | Not implemented |
| `/ws/driver` route | Registered in docs but not in main.go |
| Tests | None |
