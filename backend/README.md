# Backend — Thun-u-der Express

Go/Fiber API server for the shipment tracking platform. PostgreSQL + JWT auth.

## Quick Start

```bash
# With Docker (Postgres + backend); first copy .env.example to .env and set
# JWT_SECRET, ADMIN_EMAIL, and ADMIN_PASSWORD. Set ADMIN_ROLE_CUTOVER=true
# only for the one-time cutover after a verified backup.
cp .env.example .env
docker compose up

# Without Docker (requires separate Postgres)
cp .env.example .env
go run ./cmd/server
```

## Directory Structure

```
backend/
├── cmd/server/main.go         # Server process entry point
├── internal/
│   ├── bootstrap/             # Startup, initialization, dependency and route composition, shutdown
│   ├── server/routes.go       # Route registration for prebuilt feature handlers
│   ├── platform/config/        # Environment config and validation
│   ├── platform/database/      # PostgreSQL connection
│   ├── seed/                  # Admin/demo, hub, and shipment bootstrap
│   ├── models/                # GORM models and coordinate hooks
│   ├── middleware/            # Auth, staff checks, CORS, security, logging, rate limits
│   ├── auth/                  # Handler, service, repository, GORM adapter
│   ├── shipment/              # Handler, service, repository, GORM adapter, PostgreSQL tests
│   ├── hub/                   # Handler, service, repository, GORM adapter
│   ├── tracking/              # Focused public tracking handler and service
│   ├── analytics/             # Handler, service, repository, GORM adapter
│   └── health/                # Health endpoint
├── pkg/utils/                 # Response and password helpers
├── Dockerfile
├── docker-compose.yml
└── docs/                      # Overview, workflows, data, and logging references
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
| POST   | `/api/shipments`              | Staff | Create shipment                |
| PATCH  | `/api/shipments/:orderId/status` | Staff | Update status + log event     |
| PUT    | `/api/shipments/:orderId`     | Staff | Update shipment fields         |
| DELETE | `/api/shipments/:orderId`     | Staff | Delete shipment + events       |
| GET    | `/api/hubs`                   | No   | List all hubs                  |
| GET    | `/api/hubs/:id`               | No   | Get hub by ID                  |
| POST   | `/api/hubs`                   | Staff | Create hub                     |
| PUT    | `/api/hubs/:id`               | Staff | Update hub                     |
| DELETE | `/api/hubs/:id`               | Staff | Delete hub                     |
| GET    | `/api/track/:trackingNumber`  | No   | Public tracking lookup         |
| GET    | `/api/analytics/overview`     | No   | Dashboard aggregate stats      |
| GET    | `/api/analytics/timeseries`   | No   | Monthly + day-of-week trends   |

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
go run ./cmd/server
```


## Admin bootstrap and one-time role cutover

The server always requires `ADMIN_EMAIL` and `ADMIN_PASSWORD`, supplied through the deployment secret manager. `ADMIN_NAME` is optional and defaults to `Administrator`. The password is hashed before storage and is never logged. `ADMIN_ROLE_CUTOVER=true` is required only for the first cutover; remove or set it to `false` after successful startup.

Before the first production deployment with these settings, back up the target database and verify the backup:

```bash
export DATABASE_URL='postgres://user:password@db.example:5432/shipments?sslmode=require'
export BACKUP_FILE="shipments-before-admin-cutover-$(date -u +%Y%m%dT%H%M%SZ).dump"
pg_dump --format=custom --file="$BACKUP_FILE" "$DATABASE_URL"
pg_restore --list "$BACKUP_FILE" >/dev/null
```

Set `ADMIN_EMAIL`, `ADMIN_PASSWORD`, and optionally `ADMIN_NAME` in the deployment secret manager. Then set `ADMIN_ROLE_CUTOVER=true` and deploy. First startup atomically rejects an existing email collision, demotes every current `admin` account, creates the configured admin, and records completion. A collision aborts startup without changing roles or credentials; choose an unused email and retry. Once startup succeeds, remove `ADMIN_ROLE_CUTOVER` or set it to `false`; the durable marker keeps later startups from demoting roles or resetting credentials. Do not remove or change the seeded account: a completed marker with a missing configured admin intentionally fails startup.

## Public read-only demo account

Every startup provisions `demo@shipitout.example` with password `demo-read-only-2026` and role `customer`. The frontend displays these public credentials; they are not an admin secret. The seed records its account identity and verifies the account remains a customer with the expected credential on restart. If the email already belongs to another account, startup fails without changing that account. Never log passwords. Public shipment and hub reads need no session; writes require a current `admin` or `staff` database role, so the demo account receives 403 for direct write requests.

To restore the pre-cutover database after an unsuccessful deployment, stop API instances first, then restore the verified backup to the intended database:

```bash
pg_restore --clean --if-exists --no-owner --dbname="$DATABASE_URL" "$BACKUP_FILE"
psql "$DATABASE_URL" -v ON_ERROR_STOP=1 -c 'DROP TABLE IF EXISTS admin_bootstrap_states'
```

Restore overwrites database objects. Confirm the database target and stop writers before running it. The explicit marker-table drop removes cutover completion; it is required before another cutover. Correct deployment settings before restarting. Never run this against a persistent database unless restoring it is intended.

The optional integration test uses only an explicitly supplied PostgreSQL test database and creates/drops a private schema:

```bash
export ADMIN_BOOTSTRAP_TEST_DATABASE_URL='postgres://user:pass@127.0.0.1:55432/shipments_ticket14?sslmode=disable'
go test ./internal/seed -run TestAdminCutoverHTTPAndPersistence -count=1
```

## Configuration

| Variable            | Default                                         | Description                                      |
| ------------------- | ----------------------------------------------- | ------------------------------------------------ |
| `PORT`              | `8080`                                          | Server listen port                               |
| `DATABASE_URL`      | `postgres://user:pass@localhost:5432/shipments` | PostgreSQL DSN                                   |
| `CORS_ORIGIN`       | `http://localhost:3000`                         | Allowed direct-browser origin; not authorization |
| `JWT_SECRET`        | `change-me`                                     | HMAC secret for JWT                              |
| `ADMIN_NAME`        | `Administrator`                                | Seeded staff display name                        |
| `ADMIN_EMAIL`       | none                                            | Required deployment-provided staff email         |
| `ADMIN_PASSWORD`    | none                                            | Required deployment-provided staff password      |
| `ADMIN_ROLE_CUTOVER`| `false`                                         | Set `true` for first cutover after verified backup; then disable |

## Testing

Tests use [testify](https://github.com/stretchr/testify) (assert, require, mock). Auth service tests use fake repositories; Fiber handler tests cover HTTP responses and error mapping without a database.

### Run all tests

```bash
go test ./... -count=1
```

### Run with race detection

```bash
go test ./... -count=1 -race
```

### Run a specific package

```bash
go test ./internal/auth/ -v -count=1
go test ./internal/shipment/ -v -count=1
```

### Run a specific test

```bash
go test ./internal/auth/ -v -count=1 -run TestRegister_Success
```

### Coverage

```bash
go test ./... -count=1 -coverprofile=coverage.out
go tool cover -func=coverage.out          # Per-function breakdown
go tool cover -html=coverage.out          # HTML report
```

### Test structure

| Package | Type | Coverage | Notes |
|---------|------|----------|-------|
| `pkg/utils` | Unit | 100% | Hash/response helpers |
| `internal/platform/config` | Unit | 100% | Env loading |
| `internal/models` | Unit | 100% | GORM BeforeSave/AfterFind hooks |
| `internal/tracking` | HTTP + mock | 89% | Handler + tracking service |
| `internal/analytics` | HTTP + mock | 81% | Handler + analytics service |
| `internal/auth` | HTTP + mock | 74% | Handler + repository |
| `internal/middleware` | Unit | 61% | Rate limiter, JWT auth middleware |
| `internal/shipment` | HTTP + mock | 54% | Handler + statusToEvent (100%) |
| `internal/hub` | HTTP + mock | 49% | Handler + repository |
| `internal/bootstrap` | HTTP smoke | — | Route composition, health, and access compatibility |
Shipment repository integration coverage requires an expendable PostgreSQL database. It creates a private schema and drops it on cleanup:

```bash
export SHIPMENT_TEST_DATABASE_URL='postgres://user:pass@127.0.0.1:5432/shipments_test?sslmode=disable'
go test ./internal/shipment -run TestGormRepositoryTransactionsAndConcurrentIDs -count=1
```

The test verifies creation, status-update, and deletion rollback behavior, plus concurrent order/tracking ID uniqueness across separate database connections. Repository mappings and unique-constraint behavior remain covered by database schema/hooks and must be verified against PostgreSQL when running the suite.
