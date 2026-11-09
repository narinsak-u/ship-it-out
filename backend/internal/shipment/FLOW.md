# Shipment request flow

`bootstrap.NewApp` constructs Shipment infrastructure, application service, Hub lookup adapter, and Fiber handler. `server.RegisterRoutes` keeps public reads and staff-only writes on their existing routes.

```text
Fiber route → presentation.Handler → application.Service → application ports → infrastructure.GormRepository → PostgreSQL
```

## Layer ownership

- `domain/`: shipment and shipment-event entities, status values, address/coordinate value objects, and domain errors. No Fiber, JSON, or GORM dependencies.
- `application/`: list, read, create, update, status update, and delete workflows; validation; repository and focused Hub lookup ports.
- `infrastructure/`: GORM records, domain mappers, query filters, tracking read adapter, and PostgreSQL transactions.
- `presentation/`: Fiber handlers, request/response DTOs, and HTTP error mapping. DTOs preserve the existing JSON names.

`internal/models` remains in use by Tracking, Analytics, and seed code. Shipment's persistence adapter owns its records and schema migration models; the database schema and existing rows remain unchanged.

## Persistence guarantees

- Create and initial `Label Created` event commit or roll back together.
- Status update and its mapped event commit or roll back together.
- Shipment and its events delete in one transaction.
- ID allocation holds PostgreSQL transaction advisory lock `30030`, scans existing `ORD-` and current-year `TH` IDs, and keeps the existing formats.

The PostgreSQL integration test creates an isolated schema and checks mapping, rollback, and concurrent unique IDs. Run it against an expendable database:

```bash
SHIPMENT_TEST_DATABASE_URL='postgres://user:pass@127.0.0.1:5432/shipments_test?sslmode=disable' \
  go test ./internal/shipment/infrastructure -run TestGormRepositoryTransactionsAndConcurrentIDs -count=1
```
