# AGENTS.md — ship-simple

## Project Structure

```
ship-simple/
├── frontend/                 # Vue 3 SPA (Vite, Tailwind v4, shadcn-vue)
│   ├── src/
│   │   ├── components/ui/    # shadcn-vue primitives (auto-generated; don't edit)
│   │   ├── views/            # Page-level route components (lazy-loaded)
│   │   ├── lib/api/          # Types, API client, utilities, seed data
│   │   ├── hooks/            # TanStack Vue Query hooks (useOrders, useHubs, etc.)
│   │   ├── stores/           # Pinia stores (auth)
│   │   ├── composables/      # Composition functions (usePagination)
│   │   ├── router/           # Vue Router config
│   │   ├── App.vue
│   │   ├── main.ts
│   │   └── styles.css        # Tailwind v4 + Ocean Deep oklch theme
│   ├── package.json
│   └── vite.config.ts
├── backend/                   # Go API (Fiber v2, GORM, PostgreSQL)
│   ├── cmd/server/main.go
│   ├── internal/
│   │   ├── config/            # Env-based configuration
│   │   ├── database/          # GORM (PostgreSQL) connection
│   │   ├── models/            # User, Shipment, ShipmentEvent, Hub
│   │   ├── middleware/        # Auth (JWT), CORS, Logger
│   │   ├── auth/              # Register, Login, Me, Logout
│   │   ├── shipment/          # Shipment CRUD + status updates
│   │   ├── hub/               # Hub CRUD
│   │   ├── tracking/          # Public tracking lookup
│   │   ├── analytics/         # Dashboard aggregate stats
│   │   ├── seed/              # Demo data
│   │   └── websocket/         # Real-time tracking (Gorilla WebSockets)
│   ├── pkg/utils/             # Response writers, bcrypt helpers
│   ├── go.mod
│   └── docker-compose.yml     # Postgres 16
├── docs/
├── specs/PROJECT-PLAN.md      # Feature roadmap
├── GEMINI.md                  # AI context
└── README.md
```

## Commands

### Frontend (run from `frontend/`)

| Command | Purpose |
|---------|---------|
| `npm run dev` | Start Vite dev server |
| `npm run build` | `vue-tsc` + `vite build` |
| `npm run preview` | Preview production build |
| `npm run lint` | ESLint (flat config) |
| `npm run format` | Prettier auto-format |

**No test framework installed.** To add: `bun add -d vitest @vue/test-utils`. **Package manager:** Bun. **shadcn-vue:** `bunx shadcn-vue@latest add <component>` — don't edit `src/components/ui/`.

### Backend (run from `backend/`)

| Command | Purpose |
|---------|---------|
| `go run .` | Start API server (needs Postgres) |
| `docker compose up -d` | Start Postgres 16 |
| `go build ./...` | Compile all packages |
| `go vet ./...` | Static analysis |
| `go test ./...` | Run all tests |
| `go test -v -run TestName ./internal/shipment` | Single test |
| `go mod tidy` | Clean dependencies |

Backend: `http://localhost:8080`, Frontend: `http://localhost:5173`.

---

## Code Style — Frontend (Vue 3 / TypeScript)

- Always `<script setup lang="ts">`. No Options API. SFC order: `<script>` → `<template>` → no `<style>`.
- Props via `defineProps<{ propName: Type }>()`. Optional `class` typed as `HTMLAttributes['class']`.
- State: `ref()`, `computed()`, `watch()`. Lazy-load with `defineAsyncComponent(() => import('...'))`.
- `@/` alias maps to `./src/`. Import order: Vue/external → `@/` → relative. `import type` for types.
- Icons from `lucide-vue-next` by name (tree-shaken). Dynamic via `<component :is="..." />`.
- `interface` for objects, `type` for unions/utility types. `Record<Union, Value>` for maps.
- **Naming:** Components PascalCase (`StatusBadge.vue`). Views PascalCase + `View` suffix. Hooks/composables `use` prefix. Stores camelCase (`auth.ts`). Constants UPPER_SNAKE_CASE. Route names kebab-case.
- **State:** TanStack Vue Query for server state, Pinia (Composition API) for client state.
- **Routing:** `createWebHistory()`, lazy-loaded routes. Names: `home`, `orders`, `order-{create,edit,detail}`.
- **CSS:** Tailwind v4 only. No `<style scoped>`. `cn()` utility. Ocean Deep `oklch()` palette, dark mode.
- Custom: `bg-gradient-{hero,accent}`, `shadow-{elegant,glow}`. Fonts: `font-mono`, `font-sans`.

---

## Code Style — Backend (Golang)

- **Handlers:** `func(c *fiber.Ctx) error` — no custom wrappers.
- **Responses:** Always `pkg/utils` — `Success(c, data)`, `Error(c, status, msg)`, `SuccessWithPagination(c, data, page, limit, total)`.
- **Models:** GORM struct tags (`gorm:"..."`, `json:"..."`). Lifecycle hooks: `BeforeSave`, `AfterFind`.
- **Config:** Global `config.App` singleton via `config.Load()`.
- **Database:** Global `database.DB` (*gorm.DB).
- **Logging:** `zerolog` with `.Str().Err().Msg()` chaining.
- **Auth:** JWT (HS256), HTTP-only cookies, bcrypt.
- **API:** All routes under `/api/`.
- **Request structs** are local; update requests use pointer fields (`*string`, `*float64`) for partial updates.
- **Response format:** `{"success": true, "data": ...}` / `{"success": false, "error": "..."}`.
- **Packages:** lowercase, single word. Files: `snake_case.go`. Exported: PascalCase. Unexported: camelCase.
- **Error handling:** Explicit checks; `utils.Error(c, status, msg)` for handlers; `log.Fatal().Err(err)` for startup.
- **Data layer:** Auto-migrate on startup. Embedded structs with `embeddedPrefix`. `gorm:"-"` computed fields synced via `BeforeSave`/`AfterFind`.

---

## Environment Variables

### Backend (`backend/.env`)
| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | Server port |
| `DATABASE_URL` | `postgres://user:pass@localhost:5432/shipments` | PostgreSQL DSN |
| `JWT_SECRET` | `change-me` | JWT HMAC secret |

### Frontend (`frontend/.env`)
| Variable | Description |
|----------|-------------|
| `VITE_OPENCAGE_API_KEY` | OpenCage Geocoding key |

---

## Reference Docs

- `frontend/docs/OVERVIEW.md` — Frontend architecture
- `backend/docs/OVERVIEW.md` — Backend architecture
- `backend/docs/WORKFLOW.md` — API request/response flows
- `backend/docs/DATA.md` — Database schema
- `specs/PROJECT-PLAN.md` — Feature roadmap
