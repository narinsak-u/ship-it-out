# ship-it-out Frontend

Vue 3 SPA for Thun-u-der Express — a shipment tracking dashboard with interactive maps, real-time status updates, and logistics analytics.

## Tech Stack

| Category        | Technology                                            |
| --------------- | ----------------------------------------------------- |
| Framework       | Vue 3.5 (Composition API, `<script setup lang="ts">`) |
| Build           | Vite 6                                                |
| Language        | TypeScript 5.7 (strict)                               |
| Routing         | Vue Router 4 (6 lazy-loaded routes)                   |
| State           | Pinia (auth) + TanStack Vue Query 5 (server state)    |
| UI              | shadcn-vue (New York) on reka-ui                      |
| Styling         | Tailwind CSS v4                                       |
| Maps            | Leaflet 1.9 + CARTO dark tiles                        |
| Geocoding       | OpenCage via `opencage-api-client`                    |
| Icons           | lucide-vue-next                                       |
| Toast           | vue-sonner                                            |
| Package manager | Bun                                                   |

## Quick Start

```bash
cd frontend
bun install
cp .env.example .env   # Add VITE_OPENCAGE_API_KEY and VITE_CARTO_API_KEY
npm run dev            # http://localhost:5173
```

## Next.js preview

The migration preview serves `/`, `/orders/:orderId`, and carrier operations on `/carriers` while Vue remains the default app. `/carriers` opens Hubs; `/carriers?tab=analytics` opens API-backed analytics. Start the Go API first, then run:

```bash
API_BASE_URL=http://localhost:8080 bun run dev:next
```

`API_BASE_URL` is the server-only Go origin, without `/api`. Browser `/api` requests use the Next.js origin. Set `VITE_CARTO_API_KEY` in `.env` for dark map tiles. Run `bun run build:next` and `bun run test:e2e:next` to check the preview.

## Scripts

| Command                 | Purpose                                |
| ----------------------- | -------------------------------------- |
| `npm run dev`           | Start Vite dev server                  |
| `npm run build`         | `vue-tsc` typecheck + `vite build`     |
| `npm run preview`       | Preview production build               |
| `npm run lint`          | ESLint (flat config + Prettier)        |
| `npm run format`        | Prettier auto-format                   |
| `npm run test`          | Vitest (unit + integration, 138 tests) |
| `npm run test:coverage` | Vitest with V8 coverage report         |
| `npm run test:e2e`      | Playwright E2E (8 tests, Chromium)     |

## Project Structure

```
src/
├── components/        # Shared Vue components
│   └── ui/            # shadcn-vue primitives (auto-generated)
├── views/             # Page-level route components (lazy-loaded)
├── lib/               # Types, API client, utilities, seed data
│   └── api/           # Endpoint functions + response mappers
├── hooks/             # TanStack Vue Query hooks (useQuery/useMutation)
├── stores/            # Pinia store (auth.ts)
├── composables/       # usePagination, useSearchFilter
├── router/            # Vue Router config
├── App.vue            # Root layout + auth init
├── main.ts            # App entry point
└── styles.css         # Tailwind entry + Ocean Deep theme
tests/
├── msw/               # Mock Service Worker handlers + server
└── e2e/               # Playwright E2E tests
```

## Routes

| Path                    | View            | Description                                             |
| ----------------------- | --------------- | ------------------------------------------------------- |
| `/`                     | HomeView        | Dashboard: tracking search, KPI stats, recent shipments |
| `/orders`               | OrdersView      | Paginated, searchable, filterable orders table          |
| `/orders/create`        | OrderFormView   | Create order (auth required)                            |
| `/orders/:orderId`      | OrderDetailView | Shipment detail with Leaflet map & timeline             |
| `/orders/:orderId/edit` | OrderFormView   | Edit order (auth required)                              |
| `/carriers`             | CarriersView    | Tabs: Hubs, Analytics, Active Deliveries                |

## Key Features

- **Tracking search** — look up by order ID or tracking number from home
- **Server-side pagination** — orders table with 300ms debounced search & status filters
- **Leaflet maps** — origin → current → destination polylines with custom markers
- **Geocoding** — addresses resolved to real lat/lng via OpenCage on form submit
- **Hub-aware status updates** — status changes to hub-based states update map position to hub coordinates
- **Auth** — HTTP-only cookie sessions, signup, and visible read-only demo login; public pages work without a session
- **Dark theme** — Ocean Deep OKLCH color palette, fully responsive

## Architecture

```
Views/Components → TanStack Vue Query hooks → lib/api/ → Go backend (localhost:8080/api)
                              ↕
                    Pinia store (auth session)
```

Pagination: server-side (`OrdersView`) sends `page`/`limit`/`search` as query params; client-side (`HubsPanel`, `DeliveriesPanel`) fetches all and paginates via `usePagination`.

## Environment

| Variable                | Required | Description                                                                                                                                            |
| ----------------------- | -------- | ------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `VITE_OPENCAGE_API_KEY` | Yes      | OpenCage Geocoding API key                                                                                                                             |
| `VITE_CARTO_API_KEY`    | Yes      | CARTO Basemaps key for dark map tiles; set in Vercel build environment for Next previews. The Next config exposes this public tile key to the browser. |

## Conventions

- **SFC:** `<script setup lang="ts">` only. No Options API. Order: `<script>` → `<template>` → no `<style>`.
- **Imports:** `@/` alias for project files. Vue/external first, then `@/`, then relative. `import type` for types.
- **State:** TanStack Vue Query for server state, Pinia (Composition API) for client state.
- **Naming:** PascalCase components (`StatusBadge.vue`). Views suffixed `View`. Hooks/composables `use`-prefixed. Constants `UPPER_SNAKE_CASE`.
- **CSS:** Tailwind utility classes only. No `<style scoped>`. Ocean Deep `oklch()` dark theme.
- **Formatting:** Prettier (semicolons, double quotes, trailing commas, 100 print width).

## Testing

**Vitest** (unit + integration): 138 tests across 30 files. Tests colocated next to source (`*.spec.ts`). API calls mocked via MSW. Coverage: 92% lines, 85% statements, 91% functions, 74% branches.

**Playwright** (E2E): 8 tests across 6 files. Runs against Vite dev server. Critical path flow uses `page.route()` for API interception (no backend needed).

| Layer       | Tests | Approach                                         |
| ----------- | ----- | ------------------------------------------------ |
| Composables | 12    | Pure logic, edge cases                           |
| Stores      | 9     | Pinia, mocked API client                         |
| API         | 42    | MSW handlers for all endpoints                   |
| Hooks       | 6     | TanStack Query, mocked fetch                     |
| Components  | 30    | Vue Test Utils, stubbed children                 |
| Views       | 7     | Route + data flow smoke tests                    |
| E2E         | 8     | Playwright, page.route API mocks                 |
| Lib         | 24    | Pure functions (geocode, analytics-utils, utils) |

See `docs/OVERVIEW.md` for detailed architecture reference.
