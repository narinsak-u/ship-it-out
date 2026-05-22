# Carriers Page Design

**Date:** 2026-05-22
**Project:** ship-simple / Harbor Ops
**Status:** Approved

## Overview

Add a Carriers management dashboard page at `/carriers` for internal Harbor Ops administrators to manage carrier operations — drivers, hubs, analytics, and active deliveries. Built with Vue Query (simulated API), modular tab panels, and the existing shadcn-vue / Tailwind design system.

## Route

```
/carriers  →  CarriersView.vue  (lazy-loaded)
```

Add "Carriers" nav link to `SiteHeader.vue` between Orders and existing links.

## Architecture

### Approach: Modular Tab Panels (selected)

A single `CarriersView.vue` container with tab navigation. Each tab is a separate async component. No nested routing — tab state is a local `ref`. This matches the project's clean separation style while keeping the page fast to load.

### Component Tree

```
CarriersView.vue
├── SiteHeader (shared)
├── Page hero section (title, summary stats)
├── Tab bar (4 tabs: Drivers, Hubs, Analytics, Active Deliveries)
├── DriversPanel.vue           (async)
│   ├── Stats cards (total, available, on_delivery, off_duty)
│   ├── Driver table (search, filter by carrier/status)
│   └── AssignDriverModal.vue  (select driver + target order)
├── HubsPanel.vue              (async)
│   ├── Stats cards
│   ├── Hub table/list
│   └── HubFormModal.vue       (add/edit hub)
├── AnalyticsPanel.vue         (async)
│   ├── KPI cards
│   ├── Carrier performance (CSS bar charts)
│   └── Status distribution
└── DeliveriesPanel.vue        (async)
    ├── Active shipments table (filters to non-delivered)
    ├── Inline status update (dropdown per row)
    ├── Inline driver assignment (quick action per row)
    └── ShipmentMap (lite, async, reuses existing component)
```

## Data Model

### New Types

All defined in `src/lib/carriers.ts`:

```typescript
interface Carrier {
  id: string;
  name: string;
  contactEmail: string;
  phone: string;
  status: "active" | "inactive";
  fleetSize: number;
  totalDrivers: number;
  totalHubs: number;
  createdAt: string;
}

interface Driver {
  id: string;
  name: string;
  carrierId: string;
  phone: string;
  email: string;
  status: "available" | "on_delivery" | "off_duty";
  vehicleInfo: string;
}

interface Hub {
  id: string;
  name: string;
  carrierId: string;
  address: string;
  coords: GeoPoint;
  capacity: number;
  currentUtilization: number;
  status: "active" | "maintenance" | "closed";
}
```

### Extended Order

Add optional `driverId?: string` field to the existing `Order` type in `src/lib/orders.ts`.

### File Structure

```
src/lib/
  carriers.ts         # Types + mock data (Carrier, Driver, Hub)
  orders.ts           # Extended: driverId on Order
  api/
    carriers.ts       # Simulated: fetchCarriers, fetchDrivers, fetchHubs,
                      #   assignDriverToOrder, createHub, updateHub
    orders.ts         # Simulated: fetchActiveDeliveries, updateShipmentStatus
  utils.ts            # unchanged

src/hooks/queries/
  useCarriers.ts      # useQuery + useMutation (carriers list, drivers list, etc.)
  useDrivers.ts       # useQuery + useMutation (driver CRUD, assignment)
  useHubs.ts          # useQuery + useMutation (hub CRUD)
  useDeliveries.ts    # useQuery + useMutation (active deliveries, status updates)
```

## Data Flow

- `CarriersView.vue` orchestrates all Vue Query hooks for the 4 domains
- Each tab panel receives data via props
- Mutations use `useMutation` with `queryClient.invalidateQueries` on success
- Active deliveries query uses `refetchInterval: 15000` for live-ish polling
- No Pinia needed — all state is local refs + Vue Query cache

## Tab Details

### Drivers Tab

- Summary cards: total drivers, available, on_delivery, off_duty
- Table with columns: Driver name, Carrier, Status, Vehicle, Actions
- Actions column: "Assign to Shipment" button opens modal
- Modal: select a driver (filtered to available) + select a target order (searchable) → assign
- Search/filter: by name, carrier, status

### Hubs Tab

- Summary cards: total hubs, active, maintenance, closed
- Table with columns: Hub name, Carrier, Address, Capacity (utilization bar), Status
- "Add Hub" button opens form modal
- Each row has edit/delete actions
- Form: name, carrier (dropdown), address, capacity, status

### Analytics Tab

- KPI cards: total shipments, on-time rate, active carriers, avg delivery time
- Carrier performance bars: each carrier's shipment count + on-time % as CSS bars
- Status distribution: count per status across all carriers (simple stacked bar)
- No external chart library — all CSS-based bars to avoid bundle bloat

### Active Deliveries Tab

- Table of all non-delivered orders (from existing `orders` data)
- Columns: Order ID, Tracking, Customer, Carrier, Driver (shows assigned driver or "Unassigned"), Status, ETA, Actions
- Actions: "Update Status" dropdown (changes status inline) and "Assign Driver" (reuses shared logic)
- ShipmentMap (lite) — auto-fetches from existing `ShipmentMap.vue` to show live positions
- Refresh indicator showing seconds since last update

## State Management

| State | Mechanism |
|---|---|
| Tab selection | Local `ref<'drivers'\|'hubs'\|'analytics'\|'deliveries'>` |
| Search/filter | Local `ref<string>` per tab |
| Modal open/close | Local `ref<boolean>` |
| Server data | Vue Query hooks (useQuery) |
| Mutations | Vue Query hooks (useMutation) |
| Live polling | refetchInterval on deliveries query |

## Error Handling

| Scenario | Behavior |
|---|---|
| Loading | Skeleton placeholders (existing shadcn Skeleton) |
| API error | Inline error message: "Failed to load [resource]" with retry button |
| Empty state | "No [resource] found" with contextual suggestion |
| Mutation error | Inline error next to action button |
| Not found (bad route) | N/A — single route, no dynamic params |

Follows existing patterns from `OrdersView.vue` and `OrderDetailView.vue`.

## Testing

No test framework installed yet. When added:
- `useCarriers.ts` / `useDrivers.ts` — test with mock query client
- `CarriersView.vue` — mount with mocked Vue Query data, test tab switching
- `AssignDriverModal.vue` — test form submission, validation, driver filtering

## Non-Goals

- No real-time WebSocket (polling only for live feel)
- No external chart library (CSS bars only)
- No authentication/authz (none exists yet in frontend)
- No backend API integration (simulated only, ready to swap)
- No i18n (hardcoded English, matching existing app)
