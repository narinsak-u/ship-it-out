# ship-simple Frontend Overview

Vue 3 SPA — shipment tracking dashboard with Tailwind CSS v4, shadcn-vue (New York style), TanStack Vue Query, Pinia, and Vue Router.

**Tech stack:** Vue 3.5 (Composition API), TypeScript 5.7, Vite 6, Tailwind CSS v4, shadcn-vue + reka-ui, TanStack Vue Query 5, Pinia 2, Vue Router 4, Leaflet 1.9, vue-sonner, opencage-api-client 2.1, lucide-vue-next

**Package manager:** Bun (see `bun.lock`, `bunfig.toml`).

**Build:** `vue-tsc` typecheck → `vite build`. All dev/build/preview/lint/format commands run from `frontend/` (see `package.json` scripts).

---

## Project structure

```
frontend/
├── src/
│   ├── components/        # Shared Vue components
│   │   └── ui/            # shadcn-vue primitives (auto-generated)
│   ├── views/             # Page-level views (consumed by router)
│   ├── lib/               # Types, API client, utilities, data
│   │   └── api/           # API functions + mappers
│   ├── hooks/             # TanStack Vue Query hooks
│   ├── stores/            # Pinia stores
│   ├── composables/       # Reusable composition functions
│   ├── router/            # Vue Router setup
│   ├── App.vue            # Root component
│   ├── main.ts            # Entry point
│   └── styles.css         # Tailwind entry + theme tokens
├── docs/                  # Documentation
├── package.json
├── vite.config.ts
├── tsconfig.json
├── components.json        # shadcn-vue config
├── eslint.config.js       # Flat config (typescript-eslint + Prettier)
├── .prettierrc            # printWidth: 100, semi, double quotes, trailingComma all
└── .env                   # VITE_OPENCAGE_API_KEY
```

---

## Entry points

### `main.ts`

Creates the Vue app, installs Pinia, Vue Router, and the TanStack Vue Query plugin. Imports `styles.css` (Tailwind entry + theme). Mounts to `#app`.

### `App.vue`

Root layout: `SiteHeader` → `RouterView` → `SiteFooter` + `Toaster` (vue-sonner). Calls `authStore.init()` on mount to restore auth session.

### `styles.css`

Tailwind CSS v4 entry (`@import "tailwindcss"`). Defines the "Ocean Deep" dark theme using OKLCH color tokens. Custom utilities: `.bg-gradient-hero`, `.bg-gradient-accent`, `.shadow-elegant`, `.shadow-glow`. Imports `tw-animate-css` and Google Fonts (JetBrains Mono, Work Sans).

---

## Router (`src/router/index.ts`)

6 routes, all lazy-loaded via dynamic `import()`:

| Path | Name | View |
|------|------|------|
| `/` | `home` | `HomeView` |
| `/orders` | `orders` | `OrdersView` |
| `/orders/create` | `order-create` | `OrderFormView` |
| `/orders/:orderId/edit` | `order-edit` | `OrderFormView` |
| `/orders/:orderId` | `order-detail` | `OrderDetailView` |
| `/carriers` | `carriers` | `CarriersView` |

Uses `createWebHistory()`.

---

## Views

### `HomeView.vue`

Landing page / dashboard. Shows hero section with tracking search (looks up by order ID or tracking number), 4 KPI stat cards (Active shipments, Delivered 30d, On-time rate, Countries served), and recent shipments grid. Fetches active deliveries via `useQuery(["orders"], fetchActiveDeliveries)`.

### `OrdersView.vue`

Paginated, searchable, filterable orders table. Uses **server-side pagination** (not the client-side composable) — `currentPage` ref + `useQuery` with key `["orders", page, search, filter]`. Search debounced at 300ms. Filter pills for each status + "all". CRUD: create (via AuthModal), edit, delete (via ConfirmDialog). Auth-guarded actions.

### `OrderDetailView.vue`

Single order detail page. Two-column layout: left has route summary (origin → destination), metadata cards (tracking #, customer, carrier, weight, created), progress bar, event timeline; right has a sticky Leaflet map (`ShipmentMap` via `<Suspense>` + `defineAsyncComponent`). Fetches order by ID and its tracking events.

### `OrderFormView.vue`

Create/edit order. Determines mode from `route.params.orderId`. Renders the `OrderForm` component. Uses `useCreateOrder` / `useUpdateOrder` mutations. Shows save error banner on failure.

### `CarriersView.vue`

Tabbed view with 3 panels: Hubs (`HubsPanel`), Analytics (`AnalyticsPanel`), Active Deliveries (`DeliveriesPanel`). All panels are lazy-loaded via `defineAsyncComponent`. Tab state tracked with `activeTab` ref.

---

## Components

### Project components

| Component | Purpose | Key behavior |
|-----------|---------|--------------|
| `StatusBadge` | Shipment status indicator | Renders `<Badge>` with variant matching status + colored dot |
| `Pagination` | Page navigation | Smart ellipsis, sliding window (max 7 pages), "Showing X-Y of Z" |
| `ShipmentMap` | Leaflet route map | CartoDB dark tiles, custom divIcon markers, polylines, animated pulse on current position, auto-fit bounds |
| `SiteHeader` | Top navigation | Logo, nav links (Home/Orders/Carriers), auth state display, sign in/out |
| `SiteFooter` | Footer | Copyright text |
| `AuthModal` | Login/signup dialog | Tabbed login/signup, client-side validation, guest mode, toast on success |
| `ConfirmDialog` | Delete confirmation | Generic confirm/cancel dialog with pending state |
| `OrderForm` | Order create/edit form | Two `ThaiAddressGroup` (sender/receiver) + parcel info. Geocodes both addresses via OpenCage on submit before emitting. Shows geocode errors inline. Button has "Resolving addresses…" state. |
| `ThaiAddressGroup` | Thai address form group | 5 fields (name, zipcode, subDistrict, district, province). Auto-lookup from `thai-data` package by zipcode. Sub-district shown as `<Select>` when results available. |
| `HubFormModal` | Hub create/edit dialog | Geocodes address on submit before mutation; inline error display |
| `HubsPanel` | Hubs listing tab | Stat cards, search, Table with CRUD, pagination |
| `AnalyticsPanel` | Analytics tab | KPI cards, carrier performance bars, status distribution. Uses in-memory `orders` data. |
| `DeliveriesPanel` | Active deliveries tab | Table with inline status/hub selectors, auto-refetch (15s), mini map |

### UI components (`src/components/ui/`)

Generated by `shadcn-vue`. Primitives built on `reka-ui` (formerly radix-vue):

| Package | Components |
|---------|-----------|
| `Badge` | Badge with CVA variants (default, secondary, destructive, outline + 7 shipment status variants) |
| `Button` | Button with CVA variants (default, destructive, outline, secondary, ghost, link) and sizes |
| `Card` | Card, CardHeader, CardTitle, CardDescription, CardContent, CardFooter |
| `Dialog` | Dialog, DialogTrigger, DialogContent, DialogScrollContent, DialogOverlay, DialogHeader, DialogFooter, DialogTitle, DialogDescription, DialogClose |
| `Input` | Styled input with v-model |
| `Select` | Select, SelectTrigger, SelectContent, SelectItem, SelectValue, SelectGroup, SelectItemText, SelectLabel, SelectSeparator, SelectScrollDownButton, SelectScrollUpButton |
| `Separator` | Horizontal/vertical separator |
| `Skeleton` | Pulsing placeholder |
| `Table` | Table, TableHeader, TableBody, TableRow, TableHead, TableCell, TableCaption |
| `Sonner` | Toaster wrapper (custom theme via CSS variables, icon slots) |

---

## Libraries (`src/lib/`)

### `lib/orders.ts`

Defines core types: `ShipmentStatus` (union), `GeoPoint`, `ContactInfo`, `Location`, `TrackingEvent`, `Order`, `PaginatedResponse<T>`. Contains an in-memory `orders` array (14 hardcoded orders with realistic data) and `statusLabels` lookup map. The `getOrder(id)` helper finds by ID.

### `lib/carriers.ts`

Defines `Carrier` and `Hub` types, plus `CarrierStatus` / `HubStatus` unions. Contains in-memory data: 6 carriers, 6 hubs. Exports `hubStatusLabels`, `carrierStatusLabels`, and lookup helpers (`getCarrier`, `getCarrierByName`, `getHubsByCarrier`).

### `lib/geocode.ts`

Single exported function: `async geocodeAddress(subDistrict, district, province)` → `Promise<{lat, lng}>`. Calls the OpenCage Geocoding API via `opencage-api-client`. Reads key from `import.meta.env.VITE_OPENCAGE_API_KEY`. Throws on: missing key, network error, no results. Used by `OrderForm` and `HubFormModal`.

### `lib/utils.ts`

`cn()` — utility combining `clsx` + `tailwind-merge` for conditional Tailwind class merging.

### `lib/api/client.ts`

Thin fetch wrapper around `http://localhost:8080/api`. Exports `api` object with methods: `get`, `getRaw`, `post`, `del`, `put`, `patch`. Uses `credentials: "include"` for cookie-based auth. Returns discriminated union `{ data: T } | { error: string }`.

### `lib/api/orders.ts`

8 API functions mapping to backend endpoints:
- `fetchActiveDeliveries()` → GET `/shipments?limit=-1&exclude_status=delivered`
- `fetchOrdersPaginated(query)` → GET `/shipments?page&limit&search&status&exclude_status`
- `fetchOrder(id)` → GET `/shipments/:id`
- `createOrder(data)` → POST `/shipments`
- `updateOrder(id, data)` → PUT `/shipments/:id`
- `deleteOrder(id)` → DELETE `/shipments/:id`
- `updateShipmentStatus(orderId, status, hubId?)` → PATCH `/shipments/:id/status`
- `fetchOrderEvents(trackingNumber)` → GET `/track/:trackingNumber`

### `lib/api/carriers.ts`

5 functions: `fetchCarriers()` (in-memory, no API), `fetchHubs()`, `createHub()`, `updateHub()`, `deleteHub()`.

### `lib/api/mappers.ts`

Backend response types (`BackendShipment`, `BackendShipmentEvent`, `BackendHub`). Date formatters (`formatDate`, `formatTimestamp`). Mapper functions (`mapShipmentToOrder`, `mapEventToTrackingEvent`, `mapBackendHubToHub`).

---

## Hooks (`src/hooks/`)

TanStack Vue Query hooks wrapping API functions:

| Hook file | Exports |
|-----------|---------|
| `useOrders.ts` | `useCreateOrder()`, `useUpdateOrder()` — mutations, invalidate `["deliveries"]` |
| `useDeliveries.ts` | `useActiveDeliveries()` — query with 15s refetchInterval; `useUpdateShipmentStatus()` — mutation |
| `useHubs.ts` | `useHubs()` — query; `useCreateHub()`, `useUpdateHub()`, `useDeleteHub()` — mutations |
| `useCarriers.ts` | `useCarriers()` — query (in-memory data) |

Mutation hooks use `useQueryClient().invalidateQueries()` on success to trigger refetch. `useDeleteHub` shows a toast.

---

## Stores (`src/stores/`)

### `auth.ts`

Pinia store for authentication. State: `user`, `loading`, `error`, `isGuest`. Actions: `init()` (restore session), `login()`, `signup()`, `logout()`, `enterGuestMode()`. Guest mode persisted in `sessionStorage("harborops_guest")`. Authenticated session via HTTP-only cookie (`credentials: "include"`). Backend endpoints: `POST /auth/login`, `POST /auth/register`, `GET /auth/me`, `POST /auth/logout`.

---

## Composables (`src/composables/`)

### `usePagination.ts`

Generic **client-side** pagination. `usePagination<T>(items: Ref<T[]>, pageSize = 10)` returns `{ currentPage, totalPages, pageItems, setPage, nextPage, prevPage }`. Resets to page 1 when items array changes. Used by `HubsPanel` and `DeliveriesPanel`. Note: `OrdersView` uses server-side pagination instead (separate `currentPage` ref + `useQuery`).

---

## Data architecture

```
                   ┌──────────────────┐
                   │  Go backend API   │
                   │ localhost:8080/api│
                   └────────┬─────────┘
                            │ HTTP (credentials: include)
                   ┌────────▼─────────┐
                   │  api/client.ts    │
                   │  fetch wrapper    │
                   └────────┬─────────┘
                            │
          ┌─────────────────┼──────────────────┐
          │                 │                  │
  ┌───────▼────────┐ ┌─────▼──────┐ ┌─────────▼───────┐
  │ lib/api/orders │ │lib/api/    │ │ stores/auth.ts  │
  │                │ │carriers.ts │ │ POST /auth/*    │
  │ GET/POST/PUT   │ │            │ │ GET /auth/me    │
  │ PATCH/DELETE   │ │ GET/POST   │ │                 │
  │ /shipments/*   │ │ PUT/DELETE │ └─────────────────┘
  └───────┬────────┘ │ /hubs/*    │
          │          └────────────┘
          │
  ┌───────▼────────────────────────────────────────────┐
  │               TanStack Vue Query                   │
  │  hooks/useOrders, useDeliveries, useHubs            │
  │  (useQuery → loading/data/error; useMutation)      │
  └──────────────────────┬─────────────────────────────┘
                         │
          ┌──────────────┼──────────────┐
          │              │              │
  ┌───────▼──────┐ ┌────▼──────┐ ┌─────▼──────────┐
  │   Views      │ │Components │ │ Stores/Pinia   │
  │ (pages)      │ │  (reusable)│ │ auth.ts        │
  └──────────────┘ └───────────┘ └────────────────┘
```

**Two pagination strategies:**
- **Server-side** (`OrdersView`): page/search/filter sent as query params, backend returns `PaginatedResponse`
- **Client-side** (`HubsPanel`, `DeliveriesPanel`): fetch all, paginate in browser via `usePagination` composable

**Geocoding flow:**
Form submit → `geocodeAddress()` → OpenCage API → real `{lat, lng}` → backend persists → Leaflet map reads coordinates via `currentCoords`, `customer.coords`, `receiver.coords`.

**Status update flow:**
DeliveriesPanel selects status + hub → `PATCH /shipments/:id/status` → backend updates status + hub assignment + sets `currentCoords` to hub's location → frontend refetches → map reflects new position.
