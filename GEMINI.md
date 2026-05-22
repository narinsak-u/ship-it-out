# GEMINI.md — Harbor Ops (ship-simple)

Comprehensive instructional context for AI agents working on the Harbor Ops shipment tracking platform.

## Project Overview
Harbor Ops is a real-time shipment tracking dashboard. It features a modern Vue 3 SPA frontend and a high-performance Go API server.

- **Primary Goal:** Provide a portfolio-grade logistics platform with real-time updates, interactive maps, and analytics.
- **Frontend:** Vue 3.5+ (Composition API, Vite, Tailwind CSS v4, shadcn-vue).
- **Backend:** Go 1.24+ (Fiber v2, GORM, PostgreSQL, Redis, WebSockets).

---

## Project Structure

```text
ship-simple/
├── frontend/          # Vue 3 SPA
│   ├── src/
│   │   ├── components/  # Shared Vue components
│   │   │   └── ui/      # shadcn-vue primitives (do not edit directly)
│   │   ├── views/       # Page-level route components
│   │   ├── lib/         # Utilities, types, mock data
│   │   ├── router/      # Vue Router configuration
│   │   ├── App.vue      # Root component
│   │   └── main.ts      # App entry point
│   ├── package.json
│   ├── vite.config.ts
│   └── ...
├── backend/           # Go API server
│   ├── cmd/server/    # Application entry point (main.go)
│   ├── internal/      # Domain logic (auth, shipment, tracking, etc.)
│   ├── pkg/           # Shared utilities
│   ├── go.mod
│   └── docker-compose.yml
├── docs/              # Additional documentation & plans
└── specs/             # Project requirements and designs
```

---

## Tech Stack Details

### Frontend
- **Framework:** Vue 3 (Composition API, `<script setup lang="ts">`)
- **Build Tool:** Vite 6
- **Routing:** Vue Router 4 (Lazy-loaded routes)
- **State Management:** Pinia (Client state) + TanStack Vue Query 5 (Server state)
- **Styling:** Tailwind CSS v4 (Utility-first, using `@tailwindcss/vite`)
- **UI Components:** shadcn-vue (Radix Vue)
- **Maps:** Leaflet
- **Package Manager:** Bun (preferred)

### Backend
- **Framework:** Fiber v2 (Fast HTTP framework)
- **ORM:** GORM with PostgreSQL (pgx driver)
- **Cache/PubSub:** Redis (for real-time events)
- **Auth:** JWT (JSON Web Tokens)
- **Real-time:** Gorilla WebSockets
- **Logging:** Zerolog

---

## Building and Running

### Frontend
Run from the `frontend/` directory:
- `bun install` — Install dependencies.
- `npm run dev` — Start the development server (Vite).
- `npm run build` — Build for production (`vue-tsc && vite build`).
- `npm run lint` — Run ESLint check.
- `npm run format` — Run Prettier auto-format.

### Backend
Run from the `backend/` directory:
- `go run ./cmd/server/main.go` — Start the API server.
- `docker-compose up -d` — Start PostgreSQL and Redis (if using Docker).

---

## Development Conventions

### Frontend (Vue 3)
- **Patterns:** Always use `<script setup lang="ts">`.
- **State:** Use `useQuery` for all data fetching. Use Pinia for cross-component UI state.
- **Naming:** PascalCase for components/views (e.g., `StatusBadge.vue`). camelCase for functions.
- **Path Aliases:** Use `@/` to reference `frontend/src/`.
- **UI Primitives:** Do not edit `src/components/ui/` files directly. Use `bunx shadcn-vue@latest add <component>` to update or add.
- **Detailed Style Guide:** Refer to `AGENTS.md` for exhaustive rules on component structure, imports, and types.

### Backend (Go)
- **Package Layout:** Follow the `internal/` pattern for domain-specific logic.
- **Config:** Use `internal/config` for environment variable management (via `.env`).
- **Models:** Define database schemas in `internal/models`.
- **API Versioning:** All routes should be under `/api`.

---

## Key Files to Watch
- `frontend/src/lib/orders.ts`: Current source of mock data and shared shipment types.
- `backend/internal/database/postgres.go`: DB connection and migration logic.
- `PROJECT-PLAN.md`: Current implementation roadmap and feature list.

---

## Project Status (AI Context)
- **Authentication:** Partially implemented (JWT tokens, Register/Login routes).
- **Shipments:** Core CRUD and Status updates in place.
- **Tracking:** Public tracking route `/api/track/:trackingNumber` implemented.
- **Real-time:** WebSocket hubs for Admin, Driver, and Tracking active.
- **Frontend Views:** Home, Orders, and OrderDetail views are functional.

When suggesting changes, prioritize matching existing patterns (Composition API, Tailwind v4 utilities, and Go Fiber patterns).
