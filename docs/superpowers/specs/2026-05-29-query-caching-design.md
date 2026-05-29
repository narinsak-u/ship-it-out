# Query Caching & Data Fetching Performance

**Date:** 2026-05-29
**Status:** Approved
**Goal:** Eliminate unnecessary refetches when navigating between pages by configuring TanStack Vue Query with appropriate `staleTime` values and a centralized query key factory.

---

## Problem

The frontend uses TanStack Vue Query but with default configuration (`staleTime: 0`), causing every component mount to trigger a refetch even when data hasn't changed. Navigating away from a page and back always hits the network.

## Approach

**Global defaults + per-query overrides + query key factory.** One `QueryClient` config in `main.ts`, each query defines its own `staleTime`, and a key factory enables precise cache invalidation from mutations.

---

## Changes

### 1. QueryClient Configuration (`src/main.ts`)

Replace bare `app.use(VueQueryPlugin)` with a configured client:

```ts
import { VueQueryPlugin, QueryClient } from "@tanstack/vue-query";

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 30_000,       // 30s — safety net for anything without explicit staleTime
      gcTime: 5 * 60_000,      // 5min before unused cache is GC'd
      retry: 1,
      refetchOnWindowFocus: true,
    },
  },
});

app.use(VueQueryPlugin, { queryClient });
```

### 2. Query Key Factory (`src/lib/api/queryKeys.ts`)

Hierarchical keys for type-safe, precise cache invalidation:

| Factory | Key Structure | Used By |
|---------|--------------|---------|
| `orderKeys` | `['orders']` / `['orders','list',filters]` / `['orders','detail',id]` | Orders list, detail, edit form |
| `deliveryKeys` | `['deliveries']` / `['deliveries','active']` | HomeView "active deliveries", DeliveriesPanel |
| `hubKeys` | `['hubs']` | HubsPanel |
| `analyticsKeys` | `['analytics']` / `['analytics','timeseries']` | HomeView stats, AnalyticsPanel |

### 3. Per-Query staleTime

| Query | staleTime | Rationale |
|-------|-----------|-----------|
| Active deliveries (`deliveryKeys.active()`) | 10s | Real-time — already has `refetchInterval: 15000` |
| Hubs (`hubKeys.all`) | 2 min | Rarely changes |
| Orders list (`orderKeys.list(...)`) | 1 min | Paginated, searchable |
| Order detail (`orderKeys.detail(id)`) | 1 min | Shared between detail view and edit form |
| Order events (`['order-events', id]`) | 1 min | Tied to order detail (unchanged key) |
| Analytics overview (`analyticsKeys.all`) | 5 min | Aggregate stats |
| Time series (`analyticsKeys.timeseries()`) | 5 min | Historical chart data |

### 4. Mutation → Real-time Refresh Flow

When a mutation `onSuccess` fires, `invalidateQueries` marks matching cache entries as stale:

- **Active queries** (component mounted and watching the key) → refetch **immediately in the background**. The UI updates as soon as the response arrives, giving a real-time feel.
- **Inactive queries** (navigated away, no active observer) → marked stale. On next mount, `refetchOnMount` sees stale data and refetches.
- **Other users / external updates**: not covered here — that would require WebSocket push or periodic polling (already done for active deliveries via `refetchInterval: 15000`).

This means: edit an order → redirect back to orders list → the list query was invalidated, mounts as stale, refetches immediately. Same for hub edits/creates.

### 5. Mutation Invalidation Table

| Mutation | File | Current Invalidation | New Invalidation |
|----------|------|---------------------|------------------|
| `useCreateOrder` | `hooks/useOrders.ts` | `["deliveries"]` | `orderKeys.all`, `deliveryKeys.all` |
| `useUpdateOrder` | `hooks/useOrders.ts` | `["deliveries"]` | `orderKeys.all`, `deliveryKeys.all` |
| `deleteOrder` (inline mutation) | `views/OrdersView.vue` | `["orders"]` | `orderKeys.all`, `deliveryKeys.all` |
| `useUpdateShipmentStatus` | `hooks/useDeliveries.ts` | `["deliveries"]` | `deliveryKeys.all` |
| `useCreateHub` | `hooks/useHubs.ts` | `["hubs"]` | `hubKeys.all` |
| `useUpdateHub` | `hooks/useHubs.ts` | `["hubs"]` | `hubKeys.all` |
| `useDeleteHub` | `hooks/useHubs.ts` | `["hubs"]` | `hubKeys.all` |

Key changes for coverage:
- `useCreateOrder` / `useUpdateOrder` now also invalidate `orderKeys.all` so the OrdersView list refreshes (previously only deliveries panel refreshed)
- Delete order (inline in OrdersView) also invalidates both orders + deliveries
- `useUpdateShipmentStatus` invalidates deliveries (previously correct, stays correct)

### 6. View Key Swaps

| File | Old Key | New Key |
|------|---------|---------|
| `HomeView.vue` | `["orders"]` | `deliveryKeys.active()` |
| `OrdersView.vue` | `["orders", page, search, filter]` | `orderKeys.list(...)` |
| `OrderDetailView.vue` | `["order", id]` | `orderKeys.detail(id)` |
| `OrderFormView.vue` | `["order", id]` | `orderKeys.detail(id)` |
| `AnalyticsPanel.vue` | `["analytics"]` | `analyticsKeys.all` |
| `AnalyticsPanel.vue` | `["analytics","timeseries"]` | `analyticsKeys.timeseries()` |
| `DeliveriesPanel.vue` | `["deliveries"]` | `deliveryKeys.active()` |
| `DeliveriesPanel.vue` | `["hubs"]` | `hubKeys.all` |
| `HubsPanel.vue` | `["hubs"]` | `hubKeys.all` |

---

## Files Changed

| File | Change |
|------|--------|
| `src/main.ts` | Add configured QueryClient |
| `src/lib/api/queryKeys.ts` | **New** — key factory |
| `src/hooks/useDeliveries.ts` | Swap key strings for factory in mutations |
| `src/hooks/useHubs.ts` | Swap key strings for factory |
| `src/hooks/useOrders.ts` | Swap key strings, add orderKeys.all + deliveryKeys.all invalidation |
| `src/hooks/useDeliveries.ts` | Swap key strings to deliveryKeys |
| `src/views/HomeView.vue` | Swap key strings, add staleTime |
| `src/views/OrdersView.vue` | Swap key strings, add staleTime; update delete mutation invalidation |
| `src/views/OrderDetailView.vue` | Swap key strings, add staleTime |
| `src/views/OrderFormView.vue` | Swap key strings |
| `src/components/AnalyticsPanel.vue` | Swap key strings, add staleTime |
| `src/components/DeliveriesPanel.vue` | Swap key strings, add staleTime |
| `src/components/HubsPanel.vue` | Swap key strings |

---

## Verification

1. `npm run build` — TypeScript + Vite must pass
2. Navigate between pages: network tab should show **no duplicate requests** for cached data within staleTime
3. After mutation (create/edit/delete), affected queries should refetch automatically
4. Active deliveries auto-refresh every 15s (unchanged behavior)
