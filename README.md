# Harbor Ops — Shipment Tracking Dashboard

Real-time shipment tracking dashboard with a Vue 3 frontend and Go backend. Monitor cargo globally with interactive maps, status filters, and timeline tracking.

## Tech Stack

### Frontend

- **Framework:** Vue 3 (Composition API, `<script setup lang="ts">`)
- **Build tool:** Vite 6
- **Language:** TypeScript (strict mode)
- **Routing:** Vue Router 4 (lazy-loaded routes)
- **State management:** Pinia (client state) + TanStack Vue Query (server/cache state)
- **UI components:** shadcn-vue (New York style) built on Radix Vue
- **Styling:** Tailwind CSS v4
- **Icons:** Lucide Vue Next
- **Validation:** Zod
- **Maps:** Leaflet with CARTO dark tiles
- **Hosting:** Cloudflare Workers
- **Package manager:** Bun

### Backend

- **Language:** Go 1.24+
- **API:** REST / GraphQL (TBD)

## Getting Started

### Frontend

All commands must be run from the `frontend/` directory:

```bash
cd frontend
bun install
npm run dev
npm run build
npm run preview
```

### Backend

```bash
cd backend
go run .
```

## Available Scripts (Frontend)

| Command           | Description                            |
| ----------------- | -------------------------------------- |
| `npm run dev`     | Start Vite dev server                  |
| `npm run build`   | `vue-tsc` type-checking + `vite build` |
| `npm run preview` | Preview production build locally       |
| `npm run lint`    | ESLint check (flat config)             |
| `npm run format`  | Prettier auto-format                   |

## Project Structure

```
ship-simple/
├── frontend/            # Vue 3 SPA
│   ├── src/
│   │   ├── components/  # Shared Vue components
│   │   │   └── ui/      # shadcn-vue primitives (auto-generated)
│   │   ├── views/       # Page-level route components
│   │   ├── lib/         # Utilities, types, data
│   │   ├── router/      # Vue Router configuration
│   │   ├── App.vue      # Root component
│   │   ├── main.ts      # App entry point
│   │   └── styles.css   # Tailwind entry + theme tokens
│   ├── package.json
│   ├── vite.config.ts
│   └── ...
├── backend/             # Go API server
│   ├── go.mod
│   └── ...
└── README.md
```

## Features

- **Tracking search** — find shipments by ID or tracking number from the home page
- **Filterable manifest** — filter orders by status (Pending, In Transit, Out for Delivery, Delivered, Delayed) with text search
- **Route visualization** — Leaflet map showing origin → current → destination with styled polylines and custom markers
- **Timeline** — chronological event history per shipment with status indicators
- **Live telemetry overlay** — floating card on map showing carrier name and coordinates
- **Dark theme** — Ocean Deep color palette with CSS custom properties
- **Responsive** — mobile-first grid layout adapting from single column to multi-column

## shadcn-vue Components

Components in `frontend/src/components/ui/` are shadcn-vue primitives. Do not edit them directly. To add or update:

```bash
cd frontend
bunx shadcn-vue@latest add <component-name>
```

Available components: Badge, Button, Card (and subcomponents), Input, Separator, Skeleton, Table (and subcomponents).

## Deployment

The app is configured for Cloudflare Workers via `wrangler.jsonc`:

```bash
cd frontend
npm run build
npx wrangler deploy
```

## Design Tokens

CSS custom properties are defined in `frontend/src/styles.css` under the `:root` block:

| Token                 | Purpose                    |
| --------------------- | -------------------------- |
| `--color-background`  | Page background            |
| `--color-primary`     | Accent/action color (cyan) |
| `--color-success`     | Delivered status           |
| `--color-warning`     | Warning states             |
| `--color-destructive` | Error/delayed states       |
| `--color-info`        | In-transit status          |
| `--gradient-hero`     | Hero section gradient      |
| `--shadow-glow`       | Glowing accent shadow      |

## License

Private — internal project.
