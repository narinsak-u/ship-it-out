# Harbor Ops — Shipment Tracking Dashboard (GEMINI.md)

Comprehensive instructional context for AI agents working on the Harbor Ops shipment tracking platform.

## Project Overview
Harbor Ops is a real-time shipment tracking dashboard featuring a modern Vue 3 SPA frontend and a high-performance Go API server. It is designed to provide a portfolio-grade logistics platform with real-time updates, interactive maps, and analytics.

- **Frontend:** Vue 3.5+ (Composition API, Vite 6, Tailwind CSS v4, shadcn-vue).
- **Backend:** Go 1.24+ (Fiber v2, GORM, PostgreSQL, WebSockets).

---

## Tech Stack Details

### Frontend
- **Framework:** Vue 3 (Composition API, `<script setup lang="ts">`)
- **Build Tool:** Vite 6
- **Language:** TypeScript (Strict Mode)
- **Routing:** Vue Router 4 (Lazy-loaded routes)
- **State Management:** Pinia (Client state) + TanStack Vue Query 5 (Server state)
- **Styling:** Tailwind CSS v4 (Utility-first, using `@tailwindcss/vite`)
- **UI Components:** shadcn-vue (Radix Vue)
- **Maps:** Leaflet (with CARTO dark tiles)
- **Package Manager:** Bun (preferred)

### Backend
- **Language:** Go 1.24+
- **Framework:** Fiber v2 (Fast HTTP framework)
- **ORM:** GORM with PostgreSQL (pgx driver)
- **Rate Limiting:** In-memory sliding window (5 req/min per IP on auth endpoints)
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
- `docker-compose up -d` — Start PostgreSQL (required for local development).

---

## Project Structure

```text
ship-it-out/
├── frontend/          # Vue 3 SPA
│   ├── src/
│   │   ├── components/  # Shared Vue components
│   │   │   └── ui/      # shadcn-vue primitives (do not edit directly)
│   │   ├── views/       # Page-level route components
│   │   ├── lib/         # Utilities, types, API clients, and mock data
│   │   ├── stores/      # Pinia state management
│   │   ├── router/      # Vue Router configuration
│   │   ├── App.vue      # Root component
│   │   ├── main.ts      # App entry point
│   │   └── styles.css   # Tailwind v4 entry point & theme tokens
│   ├── package.json
│   ├── vite.config.ts
│   └── ...
├── backend/           # Go API server
│   ├── cmd/server/    # Application entry point (main.go)
│   ├── internal/      # Domain logic (auth, shipment, tracking, etc.)
│   ├── pkg/           # Shared utilities
│   ├── go.mod
│   └── docker-compose.yml
├── docs/              # Additional documentation
├── specs/             # Project requirements, designs, and PROJECT-PLAN.md
├── AGENTS.md          # Exhaustive coding conventions and style guide
└── README.md          # High-level project overview
```

---

## Development Conventions

### General
- **Contextual Precedence:** This `GEMINI.md` file and `AGENTS.md` take precedence over default AI behaviors.
- **Surgical Edits:** Use `replace` for targeted code changes. Avoid full-file rewrites unless necessary.

### Frontend (Vue 3)
- **Patterns:** Always use `<script setup lang="ts">`. Avoid Options API.
- **State:** Use `useQuery` for all data fetching. Use Pinia for cross-component UI state.
- **Styling:** Use Tailwind CSS v4 utility classes. Avoid custom CSS in SFCs.
- **UI Primitives:** Do not edit `src/components/ui/` files directly. Use `bunx shadcn-vue@latest add <component>` to update.
- **Detailed Guidance:** Refer to `AGENTS.md` for exhaustive rules on component structure, imports, and naming.

### Backend (Go)
- **Package Layout:** Follow the `internal/` pattern to keep domain logic private.
- **Error Handling:** Use explicit error checking; do not ignore errors.
- **Logging:** Use `zerolog` for structured logging throughout the application.
- **API Versioning:** Prefix all routes with `/api`.

---

## Key Files to Watch
- `frontend/src/lib/api/`: API client and data fetching logic.
- `backend/internal/models/`: Database schema definitions (GORM models).
- `backend/cmd/server/main.go`: Main server setup and route registration.
- `specs/PROJECT-PLAN.md`: The roadmap for feature implementation and status.
- `AGENTS.md`: The "Bible" for code style and architectural patterns in this project.

---

## Feature Roadmap (Summary)
1. **Phase 1 (MVP):** Auth (JWT), Shipment CRUD, Public Tracking Page, Real-time status via WebSockets.
2. **Phase 2 (Advanced):** Driver System (Dashboard & Navigation), Live Map Tracking (GPS), ETA Prediction, Hub/Warehouse System.
3. **Phase 3 (WOW):** Event-Driven Architecture (message queue), Background Workers, Optimistic UI, Barcode/QR scanning.
