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
├── cmd/server/main.go           # Server process entry point
├── internal/
│   ├── bootstrap/               # Startup, initialization, dependency and route composition, shutdown
│   ├── server/routes.go         # Route registration for prebuilt feature handlers
│   ├── platform/config/         # Environment config and validation
│   ├── platform/database/       # PostgreSQL connection
│   ├── seed/                    # Admin/demo account, hubs, and shipment bootstrap
│   ├── models/                  # GORM models and coordinate hooks
│   ├── middleware/              # Auth, staff authorization, CORS, logging, rate limits
│   ├── auth/                    # Handler, service, repository interface + GORM adapter
│   ├── shipment/               # Handler, service, repository interface + GORM adapter
│   ├── hub/                    # Handler, service, repository interface + GORM adapter
│   ├── tracking/               # Focused public tracking handler and service
│   ├── analytics/              # Focused analytics handler, service, repository + GORM adapter
│   └── health/                 # Health endpoint
├── pkg/utils/                   # Shared response and password helpers
├── Dockerfile / docker-compose.yml
└── docs/                        # Overview, workflows, data, and logging references
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
| POST   | `/api/shipments`                  | Staff | `shipment.Handler.Create`    | Create a new shipment                |
| PATCH  | `/api/shipments/:orderId/status`  | Staff | `shipment.Handler.UpdateStatus` | Update shipment status + log event |
| PUT    | `/api/shipments/:orderId`         | Staff | `shipment.Handler.Update`    | Update shipment fields               |
| DELETE | `/api/shipments/:orderId`         | Staff | `shipment.Handler.Delete`    | Delete a shipment + its events       |
| GET    | `/api/track/:trackingNumber`      | No    | `tracking.Handler.Track`     | Public tracking lookup               |
| GET    | `/api/hubs`                       | No    | `hub.Handler.List`           | List all hubs                        |
| GET    | `/api/hubs/:id`                   | No    | `hub.Handler.GetByID`        | Get hub by ID                        |
| POST   | `/api/hubs`                       | Staff | `hub.Handler.Create`         | Create a new hub                     |
| PUT    | `/api/hubs/:id`                   | Staff | `hub.Handler.Update`         | Update hub fields                    |
| DELETE | `/api/hubs/:id`                   | Staff | `hub.Handler.Delete`         | Delete a hub                         |
| GET    | `/api/analytics/overview`         | JWT   | `analytics.Handler.Overview` | Dashboard aggregate stats            |
| GET    | `/api/analytics/timeseries`       | JWT   | `analytics.Handler.TimeSeries` | Monthly + day-of-week trends       |

\* Rate-limited to 5 requests/minute per IP.

**Auth notes:** Shipment, hub, tracking, and analytics reads are public. Every shipment and hub write requires a current database role of `admin` or `staff`; a signed-in customer receives 403.

The Go API provisions the public read-only account `demo@shipitout.example` / `demo-read-only-2026` at startup. Its role is `customer`; the frontend shows these public credentials. Signup always creates a customer, even when a role is included in the request. Never log passwords or expose deployment-provisioned admin credentials.

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

- Feature packages keep HTTP handlers, Fiber-independent services, focused repository interfaces, and GORM adapters together.
- `cmd/server` delegates startup to `internal/bootstrap`, which initializes runtime state, composes feature dependencies and routes, and handles shutdown. `internal/server` registers routes around supplied handlers.
- Tracking and analytics use dedicated service/repository interfaces rather than shipment persistence APIs.
- Shipment creation, status changes, and deletion run their related database writes in transactions. PostgreSQL constraints and database-backed ID allocation protect writes across server instances.
- Coordinate hooks preserve the nested API representation while storing flat latitude/longitude columns.
- Staff writes require a fresh database role check; public reads, registration role assignment, cookie/JWT behavior, and route contracts remain unchanged.
- Seed/bootstrap functions run after AutoMigrate and preserve existing data; no refactor migration is introduced.
- Middleware order: CORS → security headers → logger; auth-specific rate limiting and auth/staff checks are attached to routes.

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
go run ./cmd/server # requires separate Postgres
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
