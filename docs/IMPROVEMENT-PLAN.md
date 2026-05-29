# Code Improvement Plan

**Generated:** 2026-05-29  
**Review Status:** Pending  
**Source:** Code review using `code-simplification`, `improve-codebase-architecture`, and `code-review-and-quality` skills

---

## Executive Summary

Three comprehensive reviews identified **40 improvement opportunities** across the codebase:

| Category | Count | Priority |
|----------|-------|----------|
| Critical Issues | 8 | Fix immediately |
| Important Issues | 14 | Fix soon |
| Suggestions | 18 | Nice to have |

**Estimated Effort:**
- Phase 1 (Critical): 1-2 days
- Phase 2 (Security): 2-3 days
- Phase 3 (Simplification): 3-4 days
- Phase 4 (Architecture): 1-2 weeks

---

## Critical Issues (Fix Immediately)

### Backend

| # | File | Issue | Impact |
|---|------|-------|--------|
| 1 | `internal/shipment/handler.go:253-262` | Duplicate hub lookup with wrong error check (`err.Error == nil`) | Bug + potential security issue |
| 2 | `internal/shipment/handler.go:83-95` | Race condition in `generateOrderID()` - concurrent requests can collide | Data integrity |
| 3 | `internal/shipment/handler.go:29-31` | Tracking number collision risk (`time.Now().UnixMilli()%100000`) | Duplicate tracking numbers |
| 4 | `internal/config/config.go:28-32` | Default JWT secret is `"change-me"` | Session hijacking if not changed |
| 5 | `internal/shipment/handler.go:98-138` | No validation on nested `ContactInfo` fields | Malformed data accepted |

### Frontend

| # | File | Issue | Impact |
|---|------|-------|--------|
| 6 | `src/components/AuthModal.vue:21-22` | Hardcoded default credentials in UI | Security anti-pattern |
| 7 | `src/lib/geocode.ts:9` | No validation API key exists before use | Runtime crashes |
| 8 | `src/stores/auth.ts:28-32` | Silent failure on auth init | Users see logged-out state incorrectly |

---

## Phased Improvement Plan

> **Note:** Tasks marked *(SKIPPED - Development)* are deferred until production readiness.

### Phase 1: Critical Bug Fixes (1-2 days)

#### 1.1 Fix UpdateStatus Hub Lookup Bug
- **File:** `backend/internal/shipment/handler.go:246-262`
- **Changes:**
  - Remove duplicate hub lookup (lines 250-251)
  - Fix error check from `err.Error == nil` to `err != nil`
  - Remove duplicate `Lng` assignment (line 262)
- **Verification:** Test status update with valid/invalid hub IDs

#### 1.2 Fix ID Generation Race Conditions
- **Files:** `backend/internal/shipment/handler.go:83-95`, `backend/internal/hub/handler.go:15-29`
- **Changes:**
  - Add database unique constraint on `order_id` and `hub_id`
  - Use database sequence or UUID instead of scan-and-increment
- **Verification:** Concurrent request test

#### 1.3 Add Input Validation
- **File:** `backend/internal/shipment/handler.go`
- **Changes:**
  - Add validation for `CreateRequest` nested fields (sender, receiver, weight)
  - Return 400 for invalid input
- **Verification:** Submit malformed requests, verify 400 responses

#### 1.4 Remove Dead Seed Data *(SKIPPED - Development)*
- **File:** `frontend/src/lib/orders.ts:67-776`
- **Status:** Deferred until production readiness
- **Note:** Seed data useful for development testing

#### 1.5 Remove Hardcoded Credentials *(SKIPPED - Development)*
- **File:** `frontend/src/components/AuthModal.vue:21-22`
- **Status:** Deferred until production readiness
- **Note:** Hardcoded creds convenient for dev testing

---

### Phase 2: Security Hardening (2-3 days)

#### 2.1 JWT Cookie Security
- **File:** `backend/internal/auth/handler.go:38-46`
- **Changes:**
  - Add `Secure: true` flag to cookie
  - Add `SameSite: http.Strict`
- **Verification:** Inspect cookie attributes in browser dev tools

#### 2.2 Add Rate Limiting
- **Files:** `backend/internal/auth/handler.go`
- **Changes:**
  - Add rate limiting middleware for `/auth/login` and `/auth/register`
  - Use Redis-backed rate limiter (once Redis is actually used)
- **Verification:** Rapid login attempts return 429

#### 2.3 Add Route Guards *(SKIPPED - Development)*
- **File:** `frontend/src/router/index.ts`
- **Status:** Deferred until production readiness
- **Note:** Open access convenient for development testing

#### 2.4 Environment-Based API URL
- **File:** `frontend/src/lib/api/client.ts:1`
- **Changes:**
  - Replace hardcoded `localhost:8080` with `import.meta.env.VITE_API_URL`
  - Add `.env.example` with default value
- **Verification:** App works with custom `VITE_API_URL`

---

### Phase 3: Code Simplification (3-4 days)

#### Frontend Simplifications

| Task | Files | Effort |
|------|-------|--------|
| Extract seed data to separate file | `lib/orders.ts` → `lib/seed-data.ts` | Low |
| Consolidate `request()` and `requestRaw()` | `lib/api/client.ts:15-55` | Low |
| Create `useSearchFilter` composable | `composables/useSearchFilter.ts` | Medium |
| Extract form submission logic | `components/AuthModal.vue` | Low |
| Single-pass status counting | `components/HubsPanel.vue:57-67` | Low |
| Remove unused Redis connection | `backend/internal/database/redis.go` | Low |

#### Backend Simplifications

| Task | Files | Effort |
|------|-------|--------|
| Extract `composeAddress` to utility | `shipment/handler.go`, `seed/shipments.go` | Low |
| Extract `getLocationForStatus()` helper | `shipment/handler.go:156-221` | Medium |
| Extract ID generation to utility | `pkg/utils/idgen.go` | Medium |
| Move province-region map to data package | `analytics/handler.go` → `internal/data/regions.go` | Low |
| Split routes to module files | `cmd/server/main.go` → `internal/*/routes.go` | Medium |

---

### Phase 4: Architecture Improvements (1-2 weeks)

#### 4.1 Introduce Repository Interfaces

```go
type ShipmentRepository interface {
    Create(ctx context.Context, s *models.Shipment) error
    FindByID(ctx context.Context, id string) (*models.Shipment, error)
    List(ctx context.Context, limit, offset int) ([]*models.Shipment, int, error)
    UpdateStatus(ctx context.Context, id string, status string) error
    Delete(ctx context.Context, id string) error
}
```

- **Benefits:**
  - Enables unit testing without real database
  - Allows mocking for integration tests
  - Decouples handlers from GORM directly
- **Files Affected:** All handler files, models, database package

#### 4.2 Replace Global State with Dependency Injection

- **Current:**
  ```go
  var DB *gorm.DB  // global
  var App Config   // global singleton
  ```
- **Target:**
  ```go
  type App struct {
      DB *gorm.DB
      Config *config.Config
  }
  func NewApp(cfg *config.Config) *App { ... }
  ```
- **Benefits:**
  - Testable with mock database
  - Configurable per environment
  - Clear dependency graph

#### 4.3 Consolidate Shallow Modules

| Module | Action | Reason |
|--------|--------|--------|
| `hooks/useAnalytics.ts` | Merge or delete | Single `useQuery` call |
| `hooks/useTimeSeries.ts` | Merge or delete | Single `useQuery` call |
| `composables/usePagination.ts` | Merge into components | Pagination is UI concern |

#### 4.4 Add Database Indexes

```sql
CREATE INDEX idx_shipments_status ON shipments(status);
CREATE INDEX idx_shipments_created_at ON shipments(created_at);
CREATE INDEX idx_shipments_tracking_number ON shipments(tracking_number);
CREATE INDEX idx_hubs_status ON hubs(status);
```

---

## File Change Summary

### Files to Delete

| File | Reason |
|------|--------|
| `frontend/src/lib/orders.ts` (lines 67-790) | Seed data + unused function |
| `backend/internal/database/redis.go` | Unused connection |

### Files to Create

| File | Purpose |
|------|---------|
| `frontend/src/lib/seed-data.ts` | Extracted seed data (if needed for dev) |
| `frontend/src/composables/useSearchFilter.ts` | Reusable search filter |
| `backend/pkg/utils/idgen.go` | Shared ID generation |
| `backend/internal/data/regions.go` | Province-region mapping |
| `backend/internal/auth/routes.go` | Module-specific routes |
| `backend/internal/shipment/routes.go` | Module-specific routes |
| `backend/internal/hub/routes.go` | Module-specific routes |
| `backend/internal/tracking/routes.go` | Module-specific routes |
| `backend/internal/analytics/routes.go` | Module-specific routes |

### Files to Modify (High Priority)

| File | Changes |
|------|---------|
| `backend/internal/shipment/handler.go` | Fix bug, extract helpers, add validation |
| `backend/internal/hub/handler.go` | Fix ID generation |
| `backend/internal/auth/handler.go` | Add rate limiting, secure cookies |
| `frontend/src/lib/api/client.ts` | Environment-based URL, consolidate requests |
| `frontend/src/components/AuthModal.vue` | Remove hardcoded creds, extract logic |
| `frontend/src/router/index.ts` | Add auth guards |
| `frontend/src/lib/orders.ts` | Remove seed data |

---

## Verification Checklist

After each phase:

- [ ] All backend tests pass (`cd backend && go test ./...`)
- [ ] All frontend builds succeed (`cd frontend && npm run build`)
- [ ] No new linting errors (`npm run lint`, `go vet ./...`)
- [ ] Manual verification of fixed bugs
- [ ] Security issues validated with testing

---

## Risk Assessment

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Breaking existing functionality | Medium | Run full test suite after each change |
| Database migration issues | Low | Use GORM auto-migrate, backup before changes |
| Environment variable misconfiguration | Low | Document in `.env.example`, validate on startup |
| Race condition fixes introduce new bugs | Medium | Add concurrent integration tests |

---

## Recommended Next Steps

1. **Start with Phase 1** - Fix critical bugs first (especially the `UpdateStatus` bug)
2. **Create GitHub issues** for each phase using this plan
3. **Tackle simplification tasks** in parallel with security fixes (they're independent)
4. **Save architecture improvements** for last - they require more design discussion

---

## Related Documents

- `specs/PROJECT-PLAN.md` - Feature roadmap
- `backend/docs/OVERVIEW.md` - Backend architecture
- `frontend/docs/OVERVIEW.md` - Frontend architecture
- `AGENTS.md` - Project structure and conventions
