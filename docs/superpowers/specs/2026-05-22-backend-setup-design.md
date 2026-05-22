# Backend Scaffold — Initial Setup Design

## Context

Initial Go backend scaffold for the real-time shipment tracking system defined in `specs/PROJECT-PLAN.md`. This covers the Phase 1 MVP backend foundation — folder structure, dependencies, entry point, and model layer.

## Structure

```
backend/
├── cmd/server/main.go         # Entry point — config load, DB connect, Fiber start
├── internal/
│   ├── config/                # Env/config loading (godotenv)
│   ├── database/              # PostgreSQL (GORM) + Redis (go-redis) connection init
│   ├── models/                # GORM model definitions
│   ├── auth/                  # Login/register handlers + JWT logic
│   ├── shipment/              # Shipment CRUD handlers
│   ├── tracking/              # Tracking number generation + public lookup
│   ├── analytics/             # Analytics query endpoints
│   ├── websocket/             # WebSocket hub pattern + client management
│   └── middleware/            # JWT auth guard, CORS, request logging
├── pkg/utils/                 # Shared helpers (hash, response formatters)
├── Dockerfile
├── docker-compose.yml
├── go.mod / go.sum
└── .env.example
```

## Dependencies

| Package                                | Purpose                         |
| -------------------------------------- | ------------------------------- |
| `github.com/gofiber/fiber/v2`          | HTTP framework                  |
| `github.com/gofiber/contrib/websocket` | Fiber WebSocket adapter         |
| `github.com/gorilla/websocket`         | Core WebSocket lib              |
| `gorm.io/gorm`                         | ORM                             |
| `gorm.io/driver/postgres`              | Postgres driver                 |
| `github.com/redis/go-redis/v9`         | Redis client (Pub/Sub, caching) |
| `github.com/golang-jwt/jwt/v5`         | JWT creation/validation         |
| `github.com/joho/godotenv`             | `.env` file loader              |
| `github.com/rs/zerolog`                | Structured logging              |

## Entry Point (`cmd/server/main.go`)

1. Load `.env` config via godotenv
2. Connect to PostgreSQL via GORM with auto-migration of models
3. Connect to Redis
4. Create Fiber app with CORS + request-logging middleware
5. Mount route groups: `/api/auth`, `/api/shipments`, `/api/track/:trackingNumber`, `/ws`
6. Start HTTP server on `:$PORT` (default 8080)

## Models (Phase 1)

### User

- `id` (uint, PK, auto-increment)
- `name` (string)
- `email` (string, unique)
- `password` (string, bcrypt hash)
- `role` (string — admin, operator, driver, customer)
- `created_at` (time)

### Shipment

- `id` (uint, PK)
- `tracking_number` (string, unique, indexed)
- `sender_name`, `receiver_name` (string)
- `origin_address`, `destination_address` (text)
- `weight` (float)
- `status` (string — enum of statuses)
- `estimated_delivery` (time)
- `created_at` (time)

### ShipmentEvent

- `id` (uint, PK)
- `shipment_id` (uint, FK → shipments)
- `status` (string)
- `location` (string, nullable)
- `description` (text)
- `created_at` (time)

## Route Map

```
POST   /api/auth/register            # Create account
POST   /api/auth/login               # Returns JWT
POST   /api/auth/refresh             # Refresh token

GET    /api/shipments                # List (auth required)
POST   /api/shipments                # Create (auth required)
GET    /api/shipments/:id            # Get by ID (auth required)
PATCH  /api/shipments/:id/status     # Update status (auth required)

GET    /api/track/:trackingNumber    # Public tracking lookup

WS     /ws/tracking/:trackingNumber  # Realtime updates for a shipment
WS     /ws/admin                     # Admin-wide events
WS     /ws/driver                    # Driver-specific events
```

## Config (.env)

```
PORT=8080
DATABASE_URL=postgres://user:pass@db:5432/shipments
REDIS_URL=redis://redis:6379
JWT_SECRET=change-me
```

## Docker

### Dockerfile (multi-stage)

- **Build stage**: `golang:1.24-alpine`, copies `go.mod`/`go.sum`, runs `go mod download`, builds binary
- **Run stage**: `alpine:3.19`, copies binary from build stage, exposes `$PORT`, runs as non-root user

### docker-compose.yml

Three services:

| Service   | Image                    | Ports       | Depends on    |
| --------- | ------------------------ | ----------- | ------------- |
| `backend` | builds from `./backend/` | `8080:8080` | `db`, `redis` |
| `db`      | `postgres:16-alpine`     | `5432:5432` | —             |
| `redis`   | `redis:7-alpine`         | `6379:6379` | —             |

Environment variables wired through `.env` file.

## Non-Goals (Phase 2+)

- Driver/Hub models and handlers
- ETA prediction logic
- Analytics aggregation
- Background workers
- Barcode/QR generation
- Nginx reverse proxy
