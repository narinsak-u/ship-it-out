# Backend — Thun-u-der Express

Go/Fiber API server for the shipment tracking platform. PostgreSQL + JWT auth.

## Quick Start

```bash
# With Docker (Postgres + backend)
docker compose up

# Without Docker (requires separate Postgres)
cp .env.example .env
go run .
```

## Directory Structure

```
backend/
├── cmd/server/main.go       # Entry point: config → DB → migrate → seed → routes → server
├── internal/
│   ├── config/config.go     # Env-based config (Port, DatabaseURL, JWTSecret, JWTTTL)
│   ├── data/regions.go      # Static Thailand province-to-region mapping
│   ├── database/
│   │   └── postgres.go      # GORM connection
│   ├── seed/
│   │   ├── hubs.go          # Demo hub data (10 hubs, nationwide)
│   │   └── shipments.go     # Demo shipments + tracking events (20 orders, nationwide)
│   ├── models/
│   │   ├── user.go          # User model
│   │   ├── shipment.go      # Shipment + ContactInfo + GeoPoint
│   │   ├── shipment_event.go# ShipmentEvent + Location
│   │   └── hub.go           # Hub model
│   ├── middleware/
│   │   ├── auth.go          # JWT Bearer/cookie verification
│   │   ├── cors.go          # CORS headers for frontend dev server
│   │   ├── logger.go        # Request logging via zerolog
│   │   └── ratelimit.go     # In-memory rate limiter (5 req/min on auth)
│   ├── auth/
│   │   ├── handler.go       # Handler (Register, Login, Me, Logout)
│   │   ├── repository.go    # Repository interface
│   │   └── gorm_repository.go
│   ├── shipment/
│   │   ├── handler.go       # Handler (List, Create, GetByID, Update, UpdateStatus, Delete)
│   │   ├── repository.go    # Repository interface + filter/result types
│   │   └── gorm_repository.go
│   ├── hub/
│   │   ├── handler.go       # Handler (List, GetByID, Create, Update, Delete)
│   │   ├── repository.go    # Repository interface
│   │   └── gorm_repository.go
│   ├── tracking/handler.go  # Handler (Track) — uses shipment.Repository
│   ├── analytics/handler.go # Handler (Overview, TimeSeries) — uses shipment.Repository
├── pkg/utils/
│   ├── hash.go              # bcrypt helpers
│   └── response.go          # Standard JSON response writers
├── Dockerfile
├── docker-compose.yml       # Postgres 16 + backend
└── docs/
    ├── OVERVIEW.md          # Full API reference + architecture
    └── WORKFLOW.md          # Workflow descriptions
```

## API Overview

| Method | Path                          | Auth | Description                    |
| ------ | ----------------------------- | ---- | ------------------------------ |
| POST   | `/api/auth/register`          | No*  | Create account                 |
| POST   | `/api/auth/login`             | No*  | Sign in, get JWT               |
| GET    | `/api/auth/me`                | JWT  | Current user profile           |
| POST   | `/api/auth/logout`            | No   | Clear auth cookie              |
| GET    | `/api/shipments`              | No   | List shipments (paginated)     |
| GET    | `/api/shipments/:orderId`     | No   | Get shipment by order ID       |
| POST   | `/api/shipments`              | JWT  | Create shipment                |
| PATCH  | `/api/shipments/:orderId/status` | JWT | Update status + log event     |
| PUT    | `/api/shipments/:orderId`     | JWT  | Update shipment fields         |
| DELETE | `/api/shipments/:orderId`     | JWT  | Delete shipment + events       |
| GET    | `/api/hubs`                   | No   | List all hubs                  |
| GET    | `/api/hubs/:id`               | No   | Get hub by ID                  |
| POST   | `/api/hubs`                   | JWT  | Create hub                     |
| PUT    | `/api/hubs/:id`               | JWT  | Update hub                     |
| DELETE | `/api/hubs/:id`               | JWT  | Delete hub                     |
| GET    | `/api/track/:trackingNumber`  | No   | Public tracking lookup         |
| GET    | `/api/analytics/overview`     | JWT  | Dashboard aggregate stats      |
| GET    | `/api/analytics/timeseries`   | JWT  | Monthly + day-of-week trends   |

\* Rate-limited to 5 requests/minute per IP.

## Seed Data

Demo data is seeded automatically on first server start (right after `AutoMigrate`):

- **10 hubs** nationwide (Bangkok, Chonburi, Kanchanaburi, Chiang Mai, Phuket, Korat, Khon Kaen, Udon Thani, Ubon Ratchathani, Buriram)
- **20 shipments** with 1-6 tracking events each, spanning Jan–May 2026 with varied statuses (pending → delivered)

Both functions are **idempotent** — they check if the table already has rows and skip if so.

To **re-seed** from scratch (deletes existing data):

```bash
# With Docker: drop volumes and restart
docker compose down -v && docker compose up

# Without Docker: connect to Postgres and truncate tables
psql $DATABASE_URL -c "TRUNCATE hubs, shipments, shipment_events RESTART IDENTITY CASCADE;"
go run .
```

## Configuration

| Variable       | Default                                                      | Description                |
| -------------- | ------------------------------------------------------------ | -------------------------- |
| `PORT`         | `8080`                                                       | Server listen port         |
| `DATABASE_URL` | `postgres://user:pass@localhost:5432/shipments`              | PostgreSQL DSN             |
| `JWT_SECRET`   | `change-me`                                                  | HMAC secret for JWT        |

## Testing

```bash
go test ./... -count=1          # Run all tests
go test ./... -count=1 -race    # With race detection
go test ./... -v -count=1       # Verbose output
```

Tests use [testify](https://github.com/stretchr/testify) (assert, require, mock) with Fiber's `app.Test()` for HTTP integration. Handler tests mock the repository layer — no database required. Coverage: ~64 tests across 10 packages, 46% overall (handler/business logic 70-100%, GORM repo code requires a live PostgreSQL database).
