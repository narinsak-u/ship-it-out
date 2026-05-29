# Testing Setup Design — ship-simple Frontend

**Date:** 2026-05-29
**Goal:** Add comprehensive unit, integration, and end-to-end tests targeting 70-90% coverage on business logic layers.

---

## 1. Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| API mocking | Fully mocked (no backend) | Tests run standalone, no docker/go required |
| Mocking libraries | MSW + vi.mock | MSW for API layer, vi.mock for component isolation |
| E2E scope | Critical paths (~12 tests) | Portfolio-grade, not exhaustive |
| Component testing | Vue Test Utils | Standard Vue ecosystem choice |
| Coverage provider | V8 | Fast, native to Vitest |
| Test organization | Colocated (`*.spec.ts` next to source) | Easy to find, zero refactor |
| Implementation order | Bottom-up | Fastest feedback loop |

---

## 2. Infrastructure

### Dependencies

```bash
# Unit/integration testing
bun add -d vitest @vue/test-utils happy-dom @vueuse/testing
bun add -d @vitest/coverage-v8
bun add -d msw

# Playwright E2E
bun add -d @playwright/test
```

### Vitest Configuration

Add `test` block to `vite.config.ts`:

```typescript
import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";
import tailwindcss from "@tailwindcss/vite";
import path from "path";

export default defineConfig({
  plugins: [vue(), tailwindcss()],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
  },
  test: {
    environment: "happy-dom",
    globals: true,
    include: ["src/**/*.spec.ts"],
    setupFiles: ["tests/setup.ts"],
    coverage: {
      provider: "v8",
      include: ["src/composables/**", "src/stores/**", "src/lib/**", "src/hooks/**"],
      exclude: ["src/components/ui/**", "src/main.ts", "src/vite-env.d.ts"],
      reporter: ["text", "html", "lcov"],
    },
  },
});
```

Key settings:
- `environment: 'happy-dom'` — lightweight DOM simulation (faster than jsdom)
- `globals: true` — `describe/it/expect` available without imports
- `coverage.include` — targets business logic layers, excludes shadcn-vue generated UI components
- `coverage.exclude` — skips auto-generated `components/ui/`, entry points

### Test Setup (`tests/setup.ts`)

```typescript
import { afterEach, beforeAll, afterAll } from "vitest";
import { server } from "./msw/server";

beforeAll(() => server.listen({ onUnhandledRequest: "error" }));
afterEach(() => server.resetHandlers());
afterAll(() => server.close());
```

### MSW Server (`tests/msw/server.ts`)

```typescript
import { setupServer } from "msw/node";
import { handlers } from "./handlers";

export const server = setupServer(...handlers);
```

### MSW Handlers (`tests/msw/handlers.ts`)

Shared mock handlers for all backend API endpoints. Each handler returns realistic data matching the backend response format (`{ data: ... }` or `{ error: ... }`).

Handlers cover:
- Auth: `GET /auth/me`, `POST /auth/login`, `POST /auth/register`, `POST /auth/logout`
- Shipments: `GET /shipments`, `POST /shipments`, `GET /shipments/:id`, `PUT /shipments/:id`, `DELETE /shipments/:id`, `PATCH /shipments/:id/status`
- Tracking: `GET /track/:trackingNumber`
- Hubs: `GET /hubs`, `POST /hubs`, `PUT /hubs/:id`, `DELETE /hubs/:id`
- Analytics: `GET /analytics/overview`, `GET /analytics/timeseries`

### MSW Browser (`tests/msw/browser.ts`)

Service worker setup for Playwright E2E tests. Used only in the browser context during E2E tests.

### Playwright Configuration (`playwright.config.ts`)

```typescript
import { defineConfig } from "@playwright/test";

export default defineConfig({
  testDir: "./tests/e2e",
  timeout: 30_000,
  retries: 1,
  use: {
    baseURL: "http://localhost:5173",
    trace: "on-first-retry",
  },
  projects: [{ name: "chromium", use: { browserName: "chromium" } }],
  webServer: {
    command: "npm run dev",
    port: 5173,
    reuseExistingServer: true,
  },
});
```

### npm Scripts

Add to `package.json`:

```json
{
  "scripts": {
    "test": "vitest run",
    "test:watch": "vitest",
    "test:coverage": "vitest run --coverage",
    "test:e2e": "playwright test",
    "test:e2e:ui": "playwright test --ui",
    "test:all": "vitest run && playwright test"
  }
}
```

---

## 3. Test Coverage Plan

### Layer 1: Composables (~12 tests)

| File | Tests |
|------|-------|
| `usePagination.spec.ts` | Page calculation from items length, empty array handling, single page, boundary (page 1, last page), `setPage` clamping, `nextPage`/`prevPage`, auto-reset to page 1 when items change |
| `useSearchFilter.spec.ts` | Case-insensitive matching, empty query returns all, no matches returns empty, multiple fields searched, null/undefined field values skipped |

**Approach:** Pure function tests using `ref()` for reactive inputs. No mocking needed.

### Layer 2: Store (~10 tests)

| File | Tests |
|------|-------|
| `auth.spec.ts` | `init()` restores user from `/auth/me`, `init()` with guest mode skips API, `login()` success sets user, `login()` error sets error message, `signup()` success sets user, `signup()` error sets error, `logout()` clears user and sessionStorage, `enterGuestMode()` sets sessionStorage flag, `isAuthenticated` computed reflects state |

**Approach:** Create Pinia test instance, vi.mock `@/lib/api/client` to control API responses.

### Layer 3: API Layer (~35 tests)

| File | Tests |
|------|-------|
| `client.spec.ts` | GET returns data, POST sends body, error response returns `{ error }`, network error returns error message, credentials included in requests |
| `mappers.spec.ts` | `mapShipmentToOrder` transforms all fields, `formatDate` formats ISO string, `formatTimestamp` includes time, `mapEventToTrackingEvent` handles description, `mapBackendHubToHub` transforms hub |
| `orders.spec.ts` | `fetchOrdersPaginated` with params, `fetchOrder` returns single order, `createOrder` sends POST, `updateOrder` sends PUT, `deleteOrder` sends DELETE, `fetchActiveDeliveries` excludes delivered, `updateShipmentStatus` sends PATCH with hubId, `fetchOrderEvents` returns events |
| `hubs.spec.ts` | `fetchHubs` returns array, `createHub` sends POST, `updateHub` sends PUT, `deleteHub` sends DELETE |
| `analytics.spec.ts` | `fetchAnalytics` returns overview, `fetchTimeSeries` returns monthly+daily data |
| `tracking.spec.ts` | `trackShipment` returns shipment data |
| `queryKeys.spec.ts` | `orderKeys` structure, `deliveryKeys` structure, `hubKeys` structure |

**Approach:** MSW intercepts fetch calls. Pure function tests for mappers/queryKeys.

### Layer 4: Hooks (~10 tests)

| File | Tests |
|------|-------|
| `useOrders.spec.ts` | `useCreateOrder` mutation calls API and invalidates cache, `useUpdateOrder` mutation calls API and invalidates cache |
| `useHubs.spec.ts` | `useHubs` query fetches data, `useCreateHub` invalidates cache, `useDeleteHub` shows toast on success |
| `useDeliveries.spec.ts` | `useActiveDeliveries` has 15s refetch interval, `useUpdateShipmentStatus` mutation calls API |

**Approach:** Mount composable in test component with Vue Query + Pinia providers. vi.mock API functions.

### Layer 5: Components (~35 tests)

| Component | Tests |
|-----------|-------|
| `StatusBadge.spec.ts` | Renders correct label for each of 7 shipment statuses |
| `Pagination.spec.ts` | Shows correct page range, ellipsis for large page counts, prev/next buttons disabled at boundaries, emits `update:currentPage` on click |
| `ConfirmDialog.spec.ts` | Renders title/description, confirm emits `confirm`, cancel emits `cancel`, pending state disables confirm button |
| `AuthModal.spec.ts` | Tab switching between login/signup, login form validation, login submit calls store, signup form validation, signup submit calls store, guest button calls enterGuestMode, error message display |
| `OrderForm.spec.ts` | Empty form shows validation errors, filled form enables submit, geocode errors display, cancel emits, submit emits with correct data |
| `ThaiAddressGroup.spec.ts` | Renders all 5 fields, v-model binding works, error messages display |
| `SiteHeader.spec.ts` | Shows auth state (logged in vs guest), nav links present |
| `SiteFooter.spec.ts` | Renders copyright text |

**Approach:** Vue Test Utils `mount()`/`shallowMount()`. Mock child components where needed. Test props, emits, and user interactions.

### Layer 6: Views (~15 tests)

| View | Tests |
|------|-------|
| `HomeView.spec.ts` | Renders hero section, displays stats, shows recent shipments, tracking form submits and navigates |
| `OrdersView.spec.ts` | Shows loading skeleton, renders table rows, filter buttons work, search input works, pagination works |
| `OrderDetailView.spec.ts` | Shows loading state, renders order details when loaded, shows timeline, 404 state when order not found |
| `OrderFormView.spec.ts` | Create mode renders form, edit mode loads order data, error banner on save failure |
| `CarriersView.spec.ts` | Tab switching renders correct panel |

**Approach:** Mount views with Vue Query + Pinia providers. Mock child components (e.g., `ShipmentMap`) and API calls via MSW or vi.mock.

### Layer 7: Playwright E2E (~12 tests)

| Flow | Tests |
|------|-------|
| **Homepage** | Page loads with hero text, stats section visible, recent shipments render, tracking form submits |
| **Orders list** | Page loads table, search filters results, status filter works, pagination navigates |
| **Create order** | Click "New Order" → auth modal → login → form loads → fill → submit → redirects to orders |
| **Order detail** | Click order ID → detail page loads with tracking number, timeline renders |
| **Navigation** | Header links navigate to all routes, active state shown |

**Approach:** Playwright against Vite dev server. E2E tests register MSW as a service worker in the browser to mock API responses — no real backend needed. The MSW browser setup (`tests/msw/browser.ts`) initializes the service worker before tests run.

---

## 4. File Structure

```
frontend/
├── src/
│   ├── composables/
│   │   ├── usePagination.ts
│   │   ├── usePagination.spec.ts
│   │   ├── useSearchFilter.ts
│   │   └── useSearchFilter.spec.ts
│   ├── stores/
│   │   ├── auth.ts
│   │   └── auth.spec.ts
│   ├── lib/
│   │   ├── utils.ts
│   │   ├── utils.spec.ts
│   │   └── api/
│   │       ├── client.ts
│   │       ├── client.spec.ts
│   │       ├── mappers.ts
│   │       ├── mappers.spec.ts
│   │       ├── orders.ts
│   │       ├── orders.spec.ts
│   │       ├── hubs.ts
│   │       ├── hubs.spec.ts
│   │       ├── analytics.ts
│   │       ├── analytics.spec.ts
│   │       ├── tracking.ts
│   │       ├── tracking.spec.ts
│   │       ├── queryKeys.ts
│   │       └── queryKeys.spec.ts
│   ├── hooks/
│   │   ├── useOrders.ts
│   │   ├── useOrders.spec.ts
│   │   ├── useHubs.ts
│   │   ├── useHubs.spec.ts
│   │   ├── useDeliveries.ts
│   │   └── useDeliveries.spec.ts
│   ├── components/
│   │   ├── StatusBadge.vue
│   │   ├── StatusBadge.spec.ts
│   │   ├── Pagination.vue
│   │   ├── Pagination.spec.ts
│   │   ├── ConfirmDialog.vue
│   │   ├── ConfirmDialog.spec.ts
│   │   ├── AuthModal.vue
│   │   ├── AuthModal.spec.ts
│   │   ├── OrderForm.vue
│   │   ├── OrderForm.spec.ts
│   │   ├── ThaiAddressGroup.vue
│   │   └── ThaiAddressGroup.spec.ts
│   └── views/
│       ├── HomeView.spec.ts
│       ├── OrdersView.spec.ts
│       ├── OrderDetailView.spec.ts
│       ├── OrderFormView.spec.ts
│       └── CarriersView.spec.ts
├── tests/
│   ├── setup.ts
│   ├── msw/
│   │   ├── server.ts
│   │   ├── handlers.ts
│   │   └── browser.ts
│   └── e2e/
│       ├── home.spec.ts
│       ├── orders.spec.ts
│       ├── order-detail.spec.ts
│       ├── create-order.spec.ts
│       └── navigation.spec.ts
├── playwright.config.ts
└── package.json
```

---

## 5. Implementation Order

1. **Install dependencies** — vitest, @vue/test-utils, happy-dom, @vitest/coverage-v8, msw, @playwright/test
2. **Configure Vitest** — update vite.config.ts, create tests/setup.ts
3. **Configure Playwright** — create playwright.config.ts
4. **Create MSW handlers** — tests/msw/handlers.ts with all API mocks
5. **Layer 1: Composables** — usePagination.spec.ts, useSearchFilter.spec.ts
6. **Layer 2: Store** — auth.spec.ts
7. **Layer 3: API** — client.spec.ts, mappers.spec.ts, orders.spec.ts, hubs.spec.ts, analytics.spec.ts, tracking.spec.ts, queryKeys.spec.ts
8. **Layer 4: Hooks** — useOrders.spec.ts, useHubs.spec.ts, useDeliveries.spec.ts
9. **Layer 5: Components** — StatusBadge, Pagination, ConfirmDialog, AuthModal, OrderForm, ThaiAddressGroup, SiteHeader, SiteFooter
10. **Layer 6: Views** — HomeView, OrdersView, OrderDetailView, OrderFormView, CarriersView
11. **Layer 7: Playwright** — e2e tests for critical paths
12. **Run coverage** — verify 70-90% on business logic layers

---

## 6. Success Criteria

- `bun run test` passes all unit/integration tests
- `bun run test:coverage` shows ≥70% line coverage on composables, store, API, hooks, mappers
- `bun run test:e2e` passes all Playwright critical path tests
- No backend required to run any tests
- All tests are colocated next to their source files
- MSW handlers cover all API endpoints used by the app
