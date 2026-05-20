# Harbor Ops — Shipment Tracking Dashboard

Real-time shipment tracking dashboard built with Vue 3 and deployed on Cloudflare Workers. Monitor cargo globally with interactive maps, status filters, and timeline tracking.

## Tech Stack

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

## Getting Started

```bash
# Install dependencies
bun install

# Start dev server
npm run dev

# Type-check and build for production
npm run build

# Preview production build
npm run preview
```

## Available Scripts

| Command           | Description                            |
| ----------------- | -------------------------------------- |
| `npm run dev`     | Start Vite dev server                  |
| `npm run build`   | `vue-tsc` type-checking + `vite build` |
| `npm run preview` | Preview production build locally       |
| `npm run lint`    | ESLint check (flat config)             |
| `npm run format`  | Prettier auto-format                   |

## Project Structure

```
src/
├── components/          # Shared Vue components
│   ├── ui/              # shadcn-vue primitives (auto-generated)
│   ├── ShipmentMap.vue  # Leaflet map with route polylines
│   ├── SiteHeader.vue   # Sticky navigation header
│   └── StatusBadge.vue  # Color-coded shipment status pill
├── views/               # Page-level route components
│   ├── HomeView.vue     # Landing page with tracking search + stats
│   ├── OrdersView.vue   # Filterable/searchable shipment table
│   └── OrderDetailView.vue  # Shipment detail + live map
├── lib/
│   ├── orders.ts        # In-memory order data + types
│   └── utils.ts         # cn() helper (clsx + tailwind-merge)
├── router/
│   └── index.ts         # Vue Router configuration
├── App.vue              # Root component
├── main.ts              # App entry point (Pinia, router, Vue Query)
└── styles.css           # Tailwind entry, theme tokens, custom utilities
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

Components in `src/components/ui/` are shadcn-vue primitives. Do not edit them directly. To add or update:

```bash
bunx shadcn-vue@latest add <component-name>
```

Available components: Badge, Button, Card (and subcomponents), Input, Separator, Skeleton, Table (and subcomponents).

## Deployment

The app is configured for Cloudflare Workers via `wrangler.jsonc`:

```bash
npm run build
npx wrangler deploy
```

## Design Tokens

CSS custom properties are defined in `src/styles.css` under the `:root` block:

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
