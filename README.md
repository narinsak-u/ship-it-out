# Thun-u-der Express — Shipment Tracking Dashboard

Real-time shipment tracking platform with a **Vue 3** frontend and **Go** backend. Monitor cargo globally with interactive maps, status updates, and timeline tracking.

---

## Tech Stack

### Frontend

| Component | Technology |
|-----------|-----------|
| Framework | Vue 3.5 (Composition API, `<script setup lang="ts">`) |
| Build tool | Vite 6 |
| Language | TypeScript 5.7 (strict mode) |
| Routing | Vue Router 4 (lazy-loaded routes) |
| State management | Pinia (client state) + TanStack Vue Query (server/cache state) |
| UI components | shadcn-vue (New York style) on reka-ui |
| Styling | Tailwind CSS v4 |
| Icons | lucide-vue-next |
| Maps | Leaflet with CARTO dark tiles |
| Geocoding | OpenCage API (`opencage-api-client`) |
| Toast | vue-sonner |
| Package manager | Bun |

### Backend

| Component | Technology |
|-----------|-----------|
| Language | Go 1.24 |
| HTTP framework | Fiber v2 (FastHTTP) |
| ORM | GORM v2 with PostgreSQL |
| Auth | JWT (HS256), HTTP-only cookies |
| Password hashing | bcrypt |
| Cache | go-redis/v9 |
| WebSocket | gorilla/websocket via fiber/contrib |
| Logging | zerolog |
| Containerization | Docker Compose (Postgres 16 + Redis 7) |

---

## Project Structure

```
ship-simple/
├── frontend/                     # Vue 3 SPA
│   ├── src/
│   │   ├── components/           # Shared Vue components
│   │   │   └── ui/               # shadcn-vue primitives (auto-generated)
│   │   ├── views/                # Page-level route components (lazy-loaded)
│   │   ├── lib/                  # Types, API client, utilities, seed data
│   │   │   └── api/              # Endpoint functions + response mappers
│   │   ├── hooks/                # TanStack Vue Query hooks
│   │   ├── stores/               # Pinia store (auth)
│   │   ├── composables/          # Reusable composition functions
│   │   ├── router/               # Vue Router config (6 routes)
│   │   ├── App.vue               # Root component
│   │   ├── main.ts               # Entry point
│   │   └── styles.css            # Tailwind entry + Ocean Deep theme
│   ├── docs/OVERVIEW.md          # Frontend architecture reference
│   ├── package.json
│   ├── vite.config.ts
│   ├── tsconfig.json
│   └── components.json           # shadcn-vue config
├── backend/                      # Go API server
│   ├── cmd/server/main.go        # Entry point: bootstrap, migrate, seed, routes
│   ├── internal/
│   │   ├── config/               # Environment-based configuration
│   │   ├── database/             # GORM + Redis connections
│   │   ├── models/               # User, Shipment, ShipmentEvent, Hub
│   │   ├── middleware/           # Auth (JWT), CORS, Logger
│   │   ├── auth/                 # Register, Login, Me, Logout
│   │   ├── shipment/             # Shipment CRUD + status updates
│   │   ├── hub/                  # Hub CRUD
│   │   ├── tracking/             # Public tracking lookup
│   │   ├── analytics/            # Dashboard aggregate stats
│   │   ├── seed/                 # Demo data (6 hubs, 12 shipments)
│   │   └── websocket/            # Real-time tracking infrastructure
│   ├── pkg/utils/                # Response writers, bcrypt helpers
│   ├── docs/                     # OVERVIEW.md, WORKFLOW.md
│   ├── Dockerfile
│   ├── docker-compose.yml
│   └── go.mod
└── README.md
```

---

## Getting Started

### Frontend

```bash
cd frontend
bun install
npm run dev        # Start Vite dev server
npm run build      # vue-tsc typecheck + vite build
npm run preview    # Preview production build
npm run lint       # ESLint check
npm run format     # Prettier auto-format
```

### Backend

```bash
# With Docker Compose (recommended)
cd backend
docker compose up

# Or manually (Postgres + Redis required)
cd backend
go run .
```

The backend starts on `http://localhost:8080` and the frontend dev server on `http://localhost:5173`.

---

## Frontend Scripts

| Command | Description |
|---------|-------------|
| `npm run dev` | Start Vite dev server |
| `npm run build` | `vue-tsc` typecheck + `vite build` |
| `npm run preview` | Preview production build locally |
| `npm run lint` | ESLint check (flat config) |
| `npm run format` | Prettier auto-format |

---

## Features

- **Tracking search** — find shipments by order ID or tracking number from the home page
- **Filterable manifest** — paginated orders table with status filters, text search, and CRUD
- **Route visualization** — Leaflet map showing origin → current → destination with styled polylines and custom markers
- **Timeline** — chronological event history per shipment with status indicators and location
- **Live telemetry overlay** — floating card on map showing current coordinates
- **Status updates with hub awareness** — changing a shipment's status to a hub-based state (departed, in_transit) updates the map position to that hub's real coordinates
- **Geocoding** — addresses are resolved to real lat/lng via OpenCage API before creating orders or hubs
- **Hub management** — CRUD operations for logistics hubs with capacity tracking
- **Analytics dashboard** — KPI cards and carrier performance breakdowns
- **Dark theme** — Ocean Deep OKLCH color palette
- **Responsive layout** — mobile-first grid system

---

## Backend API Overview

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| POST | `/api/auth/register` | No | Register a new user |
| POST | `/api/auth/login` | No | Login, returns JWT + cookie |
| GET | `/api/auth/me` | JWT | Get current user profile |
| POST | `/api/auth/logout` | No | Clear auth cookie |
| GET | `/api/shipments` | No | List shipments (paginated, searchable, filterable) |
| GET | `/api/shipments/:orderId` | No | Get shipment by order ID |
| POST | `/api/shipments` | JWT | Create a new shipment |
| PATCH | `/api/shipments/:orderId/status` | JWT | Update status + log event |
| PUT | `/api/shipments/:orderId` | JWT | Update shipment fields |
| DELETE | `/api/shipments/:orderId` | JWT | Delete shipment + events |
| GET | `/api/track/:trackingNumber` | No | Public tracking lookup |
| GET | `/api/hubs` | No | List all hubs |
| GET | `/api/hubs/:id` | No | Get hub by ID |
| POST | `/api/hubs` | JWT | Create a hub |
| PUT | `/api/hubs/:id` | JWT | Update hub fields |
| DELETE | `/api/hubs/:id` | JWT | Delete a hub |
| GET | `/api/analytics/overview` | JWT | Dashboard aggregate stats |
| GET | `/ws/tracking/:trackingNumber` | No | Real-time tracking WebSocket |

---

## Design Tokens

CSS custom properties defined in `frontend/src/styles.css` under the `:root` block:

| Token | Purpose |
|-------|---------|
| `--color-background` | Page background |
| `--color-primary` | Accent/action color (cyan) |
| `--color-success` | Delivered status |
| `--color-warning` | Warning states |
| `--color-destructive` | Error/delayed states |
| `--color-info` | In-transit status |
| `--gradient-hero` | Hero section gradient |
| `--shadow-glow` | Glowing accent shadow |

---

## Environment Variables

### Backend (`backend/.env`)

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | Server listen port |
| `DATABASE_URL` | `postgres://user:pass@localhost:5432/shipments` | PostgreSQL DSN |
| `REDIS_URL` | `redis://localhost:6379` | Redis connection URL |
| `JWT_SECRET` | `change-me` | HMAC secret for JWT |

### Frontend (`frontend/.env`)

| Variable | Description |
|----------|-------------|
| `VITE_OPENCAGE_API_KEY` | OpenCage Geocoding API key |

---

## Documentation

- `frontend/docs/OVERVIEW.md` — Frontend architecture, components, data flow
- `backend/docs/OVERVIEW.md` — Backend architecture, routes, data model
- `backend/docs/WORKFLOW.md` — Detailed request/response flows for every endpoint
