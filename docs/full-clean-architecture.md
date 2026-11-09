# Ship It Out — Feature-based Clean Architecture Refactoring Specification

## 1. Objective

Refactor the existing Ship It Out Go Fiber backend into a **Feature-based Clean Architecture**, following idiomatic Go project conventions.

Use:

- `cmd/server/` for the application entry point.
- `internal/` for application-specific feature packages.
- Four architectural layers within each feature where necessary:
  - `domain/`
  - `application/`
  - `infrastructure/`
  - `presentation/`

Preserve existing behavior, database schema, API contracts, and frontend compatibility.

## 2. Target Architecture

```text
backend/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── shipment/
│   │   ├── domain/
│   │   ├── application/
│   │   ├── infrastructure/
│   │   └── presentation/
│   ├── tracking/
│   │   ├── application/
│   │   └── presentation/
│   ├── hub/
│   │   ├── domain/
│   │   ├── application/
│   │   ├── infrastructure/
│   │   └── presentation/
│   ├── auth/
│   │   ├── domain/
│   │   ├── application/
│   │   ├── infrastructure/
│   │   └── presentation/
│   ├── analytics/
│   │   ├── application/
│   │   ├── infrastructure/
│   │   └── presentation/
│   ├── platform/
│   ├── bootstrap/
│   └── seed/
├── go.mod
└── Dockerfile
```

Do not create empty layers or unnecessary files.

## 3. Dependency Rules

### Domain

- Contains entities, business rules, value objects, and domain errors.
- Must not import Fiber, GORM, SQL drivers, or other architectural layers.
- Must not contain HTTP-specific code.

### Application

- Contains use cases, services, input/output types, and repository interfaces.
- Depends on Domain.
- Must not depend on Infrastructure or Presentation.
- Uses `context.Context` for I/O operations.
- Defines interfaces close to consumers.

### Infrastructure

- Implements repository and external service interfaces.
- Contains GORM persistence models, database queries, mappers, and external adapters.
- May depend on Domain and Application.
- Must not depend on Presentation.

### Presentation

- Contains Fiber handlers, routes, request/response DTOs, and HTTP error mapping.
- Calls Application services through interfaces or concrete application service types.
- Must not directly execute database queries or depend on GORM repositories.
- Must preserve existing HTTP contracts.

### Bootstrap

- Acts as the Composition Root.
- Constructs repositories, services, and handlers.
- Wires dependencies through constructors.
- May import all architectural layers.
- Must not contain business logic.

## 4. Feature Refactoring

### Shipment

Separate:

- Domain: Shipment, ShipmentEvent, Status, business invariants.
- Application: Create, Update, Delete, List, GetByID, UpdateStatus workflows.
- Infrastructure: GORM repositories, persistence models, mappers.
- Presentation: Fiber handlers, routes, HTTP DTOs.

Ensure transactional consistency for:

1. Shipment creation and initial event.
2. Status update and tracking event.
3. Shipment deletion and related events.

Preserve existing ID formats and implement concurrency-safe generation.

### Tracking

Use an Application service and Presentation handler.

Reuse a focused Shipment read interface for tracking lookup and event history.

Do not duplicate Shipment domain entities or create unnecessary persistence implementations.

### Hub

Separate Hub entities, application services, GORM persistence, and HTTP handlers.

Preserve all existing CRUD APIs.

### Authentication

Separate:

- Domain: User and roles.
- Application: Registration, login, user lookup, token issuance orchestration.
- Infrastructure: GORM user repository, password hashing, JWT signing.
- Presentation: Fiber handlers, cookies, authentication middleware.

Do not allow public users to assign privileged roles.

Never expose password hashes through API responses.

### Analytics

Use:

- Application services for statistics and aggregation workflows.
- Infrastructure for database aggregation queries.
- Presentation for HTTP handlers.

Reuse focused read ports rather than importing the entire Shipment repository interface.

## 5. Application Entry Point

Keep:

```text
cmd/server/main.go
```

Responsibilities:

- Initialize logging.
- Start application bootstrap.
- Start Fiber server.
- Handle startup errors and graceful shutdown.

Use manual dependency injection.

Do not introduce dependency injection frameworks.

## 6. Compatibility Requirements

Preserve:

- All existing `/api/*` endpoints.
- Existing HTTP methods.
- Request and response JSON structures.
- Query parameters and pagination.
- Authentication cookies and authorization behavior, except explicitly documented security fixes.
- PostgreSQL schema, relationships, and existing data.
- Existing frontend compatibility.
- Docker Compose and environment configuration.

Do not introduce breaking changes without explicit approval.

## 7. Testing

Maintain or add:

- Domain unit tests for business rules.
- Application unit tests using fake repositories.
- Infrastructure integration tests with PostgreSQL.
- Presentation handler tests using Fiber.
- API regression tests.
- Transaction and concurrency tests.

Run:

```bash
go test ./...
go vet ./...
go test -race ./...
```

Do not remove existing tests to make the refactor pass.

## 8. Migration Strategy

Implement incrementally:

1. Inspect the existing code and establish a test baseline.
2. Refactor Shipment into four layers.
3. Refactor Tracking and Hub.
4. Refactor Authentication.
5. Refactor Analytics.
6. Move shared infrastructure into Platform.
7. Introduce Bootstrap wiring and simplify `cmd/server/main.go`.
8. Remove obsolete code after verifying replacements.
9. Update documentation, Docker configuration, and tests.

Keep the application buildable throughout the migration.

## 9. Constraints

- Keep Go Fiber v2, GORM, and PostgreSQL.
- Do not introduce microservices.
- Do not introduce CQRS, Event Sourcing, or DI frameworks.
- Avoid one use-case struct per endpoint unless justified.
- Avoid unnecessary interfaces and abstractions.
- Do not duplicate domain entities across features.
- Do not create circular dependencies.
- Do not change frontend code unless required for compatibility.
- Do not deploy or modify production data automatically.

## 10. Definition of Done

- [ ] Backend is organized by business features.
- [ ] Four Clean Architecture layers exist where appropriate.
- [ ] `cmd/server/main.go` remains the entry point.
- [ ] Domain and Application are framework-independent.
- [ ] Infrastructure implements Application ports.
- [ ] Presentation depends on Application, not persistence implementations.
- [ ] Bootstrap wires dependencies manually.
- [ ] Shipment operations are transaction-safe.
- [ ] Existing API contracts are preserved.
- [ ] Database schema and data remain compatible.
- [ ] All relevant tests pass.
- [ ] Docker builds successfully.
- [ ] Frontend integration remains functional.
- [ ] No unused legacy code remains.

**Implementation principle:** Favor clear feature boundaries and dependency direction over unnecessary abstraction. Apply full Clean Architecture where it adds value, and keep simple features lightweight.