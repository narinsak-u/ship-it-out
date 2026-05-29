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
│   ├── config/config.go     # Env-based config (Port, DatabaseURL, JWTSecret)
│   ├── database/
│   │   └── postgres.go      # GORM connection
│   ├── seed/
│   │   ├── hubs.go          # Demo hub data (6 hubs, Eastern Thailand)
│   │   └── shipments.go     # Demo shipments + tracking events (3 orders, Eastern Thailand)
│   ├── models/
│   │   ├── user.go          # User model
│   │   ├── shipment.go      # Shipment + ContactInfo + GeoPoint
│   │   ├── shipment_event.go# ShipmentEvent + Location
│   │   └── hub.go           # Hub model
│   ├── middleware/
│   │   ├── auth.go          # JWT Bearer/cookie verification
│   │   ├── cors.go          # CORS headers for frontend dev server
│   │   └── logger.go        # Request logging via zerolog
│   ├── auth/handler.go      # Register, Login, Me, Logout
│   ├── shipment/handler.go  # Shipment CRUD + status update
│   ├── hub/handler.go       # Hub CRUD
│   ├── tracking/handler.go  # Public tracking lookup
│   ├── analytics/handler.go # Dashboard aggregate stats
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

| Method | Path                          | Auth | Description                |
| ------ | ----------------------------- | ---- | -------------------------- |
| POST   | `/api/auth/register`          | No   | Create account             |
| POST   | `/api/auth/login`             | No   | Sign in, get JWT           |
| GET    | `/api/auth/me`                | JWT  | Current user profile       |
| POST   | `/api/auth/logout`            | No   | Clear auth cookie          |
| GET    | `/api/shipments`              | JWT  | List all shipments         |
| POST   | `/api/shipments`              | JWT  | Create shipment            |
| GET    | `/api/shipments/:id`          | JWT  | Get shipment by ID         |
| PATCH  | `/api/shipments/:id/status`   | JWT  | Update status + log event  |
| GET    | `/api/hubs`                   | JWT  | List all hubs              |
| POST   | `/api/hubs`                   | JWT  | Create hub                 |
| GET    | `/api/hubs/:id`               | JWT  | Get hub by ID              |
| PUT    | `/api/hubs/:id`               | JWT  | Update hub                 |
| DELETE | `/api/hubs/:id`               | JWT  | Delete hub                 |
| GET    | `/api/track/:trackingNumber`  | No   | Public tracking lookup     |
| GET    | `/api/analytics/overview`     | JWT  | Dashboard stats            |

## Seed Data

Demo data is seeded automatically on first server start (right after `AutoMigrate`):

- **6 hubs** across Eastern Thailand (Laem Chabang Port, Pattaya, Rayong, Chanthaburi, Chachoengsao, Trat)
- **3 shipments** with tracking events — all within Eastern Thailand (Chonburi ↔ Chanthaburi, Chachoengsao ↔ Pattaya, Rayong ↔ Trat)

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

| Variable       | Default                                                      |
| -------------- | ------------------------------------------------------------ |
| `PORT`         | `8080`                                                       |
| `DATABASE_URL` | `postgres://user:pass@localhost:5432/shipments`              |
| `JWT_SECRET`   | `change-me`                                                  |
