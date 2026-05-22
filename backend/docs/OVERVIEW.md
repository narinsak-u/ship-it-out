# Backend Overview

Go API server for a shipment tracking platform. Built with **Fiber v2**, **GORM** (PostgreSQL), **go-redis**, and **JWT** authentication.

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
├── cmd/server/main.go          # Entry point: config load, DB init, route setup, server start
├── internal/
│   ├── config/config.go        # Global Config struct (Port, DatabaseURL, RedisURL, JWTSecret)
│   ├── database/
│   │   ├── postgres.go         # GORM connection + auto-migration
│   │   └── redis.go            # go-redis client (initialized but unused)
│   ├── models/
│   │   ├── user.go             # User model (ID, Name, Email, Password, Role, CreatedAt)
│   │   ├── shipment.go         # Shipment model (tracking number, addresses, weight, status, dates)
│   │   └── shipment_event.go   # ShipmentEvent model (FK to Shipment, status, location, description)
│   ├── middleware/
│   │   ├── auth.go             # JWT Bearer token verification
│   │   ├── cors.go             # Manual CORS headers
│   │   └── logger.go           # Request logging via zerolog
│   ├── auth/handler.go         # POST /api/auth/register, POST /api/auth/login
│   ├── shipment/handler.go     # CRUD + status update for shipments
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
└── docs/OVERVIEW.md
```

---

## API Endpoints

| Method | Path                              | Auth     | Handler                    | Description                        |
| ------ | --------------------------------- | -------- | -------------------------- | ---------------------------------- |
| POST   | `/api/auth/register`              | No       | `auth.Register`            | Register a new user                |
| POST   | `/api/auth/login`                 | No       | `auth.Login`               | Login, returns JWT token           |
| GET    | `/api/shipments/`                 | JWT      | `shipment.List`            | List all shipments                 |
| POST   | `/api/shipments/`                 | JWT      | `shipment.Create`          | Create a new shipment              |
| GET    | `/api/shipments/:id`              | JWT      | `shipment.GetByID`         | Get shipment by ID                 |
| PATCH  | `/api/shipments/:id/status`       | JWT      | `shipment.UpdateStatus`    | Update shipment status + log event |
| GET    | `/api/track/:trackingNumber`      | No       | `tracking.Track`           | Public tracking lookup             |
| GET    | `/api/analytics/overview`         | JWT      | `analytics.Overview`       | Dashboard aggregate stats          |
| GET    | `/ws/tracking/:trackingNumber`    | No       | `websocket.HandleWebSocket`| Real-time tracking updates         |
| GET    | `/ws/admin`                       | No       | `websocket.HandleWebSocket`| Admin WebSocket (room "global")    |
| GET    | `/ws/driver`                      | No       | `websocket.HandleWebSocket`| Driver WebSocket (room "global")   |

---

## Data Model

```
User (standalone)
  - ID, Name, Email, Password (bcrypt, hidden from JSON), Role, CreatedAt

Shipment 1──N ShipmentEvent
  - Shipment: ID, TrackingNumber (unique), SenderName, ReceiverName,
              OriginAddress, DestinationAddress, Weight, Status, EstimatedDelivery, CreatedAt
  - ShipmentEvent: ID, ShipmentID (FK), Status, Location, Description, CreatedAt
```

- GORM auto-migration creates all tables at startup.
- Tracking number format: `TH<YYYY><5-digit>` (e.g., `TH202612345`).

---

## Architecture Notes

- **Handler-centric:** Each domain package has a single `handler.go`. Handlers call GORM directly via the global `database.DB` — no service/repository layers.
- **Global state:** `database.DB`, `database.Redis`, and `config.App` are package-level globals (no dependency injection).
- **Standardized responses:** All endpoints return `{"success": true/false, "data": ...}` or `{"success": false, "error": "..."}`.
- **JWT auth:** HS256 bearer tokens with `user_id` and `role` claims. Token TTL is hardcoded to 24h.
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
