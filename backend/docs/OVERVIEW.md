# Backend Overview

Go API server for the Thun-u-der Express shipment tracking platform. Built with **Fiber v2**, **GORM** (PostgreSQL), **go-redis**, and **JWT** authentication.

---

## Tech Stack

| Component        | Library                                      |
| ---------------- | -------------------------------------------- |
| HTTP Framework   | `gofiber/fiber/v2` (FastHTTP)                |
| ORM              | `gorm.io/gorm` + `gorm.io/driver/postgres`   |
| Cache            | `redis/go-redis/v9`                          |
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
│   │   ├── postgres.go         # GORM connection (global DB var)
│   │   └── redis.go            # go-redis client (global Redis var)
│   ├── seed/
│   │   ├── hubs.go             # SeedHubs — inserts 6 demo hubs (skips if table non-empty)
│   │   └── shipments.go        # SeedShipments — inserts 3 demo shipments + events (skips if table non-empty)
│   ├── models/
│   │   ├── user.go             # User model (ID, Name, Email, Password, Role, CreatedAt)
│   │   ├── shipment.go         # Shipment + ContactInfo + GeoPoint (embedded, GORM hooks for coords)
│   │   ├── shipment_event.go   # ShipmentEvent + Location (FK to Shipment, status, description, timestamp)
│   │   └── hub.go              # Hub model (ID, Name, CarrierID, Address, Coords, Capacity, Status)
│   ├── middleware/
│   │   ├── auth.go             # JWT Bearer token / cookie verification
│   │   ├── cors.go             # Manual CORS headers
│   │   └── logger.go           # Request logging via zerolog
│   ├── auth/handler.go         # POST /api/auth/register, POST /api/auth/login, GET /api/auth/me, POST /api/auth/logout
│   ├── shipment/handler.go     # CRUD + status update for shipments
│   ├── hub/handler.go          # CRUD for logistics hubs
│   ├── tracking/handler.go     # GET /api/track/:trackingNumber (public)
│   ├── analytics/handler.go    # GET /api/analytics/overview (dashboard stats)
│   └── websocket/
│       ├── hub.go              # Connection hub (room-based broadcast)
│       └── client.go           # WebSocket client + upgrade handler
├── pkg/utils/
│   ├── hash.go                 # bcrypt hash/compare helpers
│   └── response.go             # Standard JSON response writers (Success/Error)
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

| Method | Path                              | Auth     | Handler                      | Description                          |
| ------ | --------------------------------- | -------- | ---------------------------- | ------------------------------------ |
| POST   | `/api/auth/register`              | No       | `auth.Register`              | Register a new user                  |
| POST   | `/api/auth/login`                 | No       | `auth.Login`                 | Login, returns JWT + sets cookie     |
| GET    | `/api/auth/me`                    | JWT      | `auth.Me`                    | Get current user profile             |
| POST   | `/api/auth/logout`                | No       | `auth.Logout`                | Clear auth cookie                    |
| GET    | `/api/shipments`                  | JWT      | `shipment.List`              | List all shipments                   |
| POST   | `/api/shipments`                  | JWT      | `shipment.Create`            | Create a new shipment                |
| GET    | `/api/shipments/:id`              | JWT      | `shipment.GetByID`           | Get shipment by ID                   |
| PATCH  | `/api/shipments/:id/status`       | JWT      | `shipment.UpdateStatus`      | Update shipment status + log event   |
| GET    | `/api/hubs`                       | JWT      | `hub.List`                   | List all hubs                        |
| POST   | `/api/hubs`                       | JWT      | `hub.Create`                 | Create a new hub                     |
| GET    | `/api/hubs/:id`                   | JWT      | `hub.GetByID`                | Get hub by ID                        |
| PUT    | `/api/hubs/:id`                   | JWT      | `hub.Update`                 | Update hub fields                    |
| DELETE | `/api/hubs/:id`                   | JWT      | `hub.Delete`                 | Delete a hub                         |
| GET    | `/api/track/:trackingNumber`      | No       | `tracking.Track`             | Public tracking lookup               |
| GET    | `/api/analytics/overview`         | JWT      | `analytics.Overview`         | Dashboard aggregate stats            |
| GET    | `/ws/tracking/:trackingNumber`    | No       | `websocket.HandleWebSocket`  | Real-time tracking updates           |
| GET    | `/ws/admin`                       | No       | `websocket.HandleWebSocket`  | Admin WebSocket (room "global")      |
| GET    | `/ws/driver`                      | No       | `websocket.HandleWebSocket`  | Driver WebSocket (room "global")     |

---

## Data Model

```
User (standalone)
  - ID, Name, Email, Password (bcrypt, hidden from JSON), Role, CreatedAt

Shipment 1──N ShipmentEvent
  - Shipment: ID, TrackingNumber (unique), Customer (ContactInfo embedded),
    Receiver (ContactInfo embedded), Origin, Destination, CurrentCoords,
    Status, Carrier, DriverID (optional), Weight, Items, EstimatedDelivery,
    CreatedAt, Progress
  - ShipmentEvent: ID, ShipmentID (FK), Status, Location (embedded),
    Description (optional), CreatedAt

Hub (standalone)
  - ID (string PK), Name, CarrierID, Address, Coords (GeoPoint),
    Capacity, CurrentUtilization, Status, CreatedAt
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

This gives us clean JSON with nested objects while keeping normal DB columns for querying.

- GORM auto-migration creates all tables at startup.
- Tracking number format: `TH<YYYY><5-digit>` (e.g., `TH202612345`).

---

## Architecture Notes

- **Handler-centric:** Each domain package has a single `handler.go`. Handlers call GORM directly via the global `database.DB` — no service/repository layers.
- **Global state:** `database.DB`, `database.Redis`, and `config.App` are package-level globals (no dependency injection).
- **Standardized responses:** All endpoints return `{"success": true/false, "data": ...}` or `{"success": false, "error": "..."}`.
- **JWT auth:** HS256 bearer tokens with `user_id` and `role` claims. Token TTL is hardcoded to 24h. Auth middleware checks `Authorization: Bearer` header first, then falls back to the `jwt` cookie.
- **Seed data:** On first startup, `main.go` runs `seed.SeedHubs()` and `seed.SeedShipments()` after AutoMigrate. Both skip if their table already has rows (idempotent).
- **WebSocket pub/sub:** Room-based hub broadcasts to clients by tracking number or "global" room. Infrastructure is wired up; actual push events are not yet implemented.
- **Redis is initialized but not used** by any handler yet.
- **No tests** exist in the codebase.
- **Middlewares run in order:** CORS → Logger → Auth (selectively applied).

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
