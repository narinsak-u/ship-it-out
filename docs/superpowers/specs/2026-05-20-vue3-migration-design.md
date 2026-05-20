# Design Spec: Vue 3 Migration (SPA)

**Date:** 2026-05-20
**Status:** Draft (Pending Review)
**Topic:** Migrating from React + TanStack Start to Vue 3 (Composition API) SPA.

## 1. Overview
The goal of this project is to migrate the "Harbor Ops" shipment tracking application from its current React-based TanStack Start architecture to a modern Vue 3 Single Page Application (SPA).

### Key Objectives
- Switch from React to Vue 3 (Composition API).
- Move from TanStack Start (SSR) to a standard Client-Side Rendered (CSR) SPA.
- Maintain the visual design using `shadcn-vue`.
- Preserive server state management using TanStack Query (Vue version).

## 2. Tech Stack
- **Framework:** Vue 3.5+ (Composition API, `<script setup>`).
- **Build Tool:** Vite 6+.
- **Routing:** Vue Router 4.
- **State Management:** 
  - **Pinia:** Global UI state (if needed).
  - **@tanstack/vue-query:** Server state/caching.
- **UI & Styling:**
  - **Tailwind CSS v4** (existing).
  - **shadcn-vue:** Ported components from `radix-vue`.
  - **Lucide Vue Next:** Icons.
- **Maps:** Leaflet (via standard JS or a Vue-specific wrapper).

## 3. Architecture
The app will be converted from a file-based SSR router (TanStack Start) to a standard Vue SPA structure.

### Project Structure Mapping
- `src/main.ts`: Entry point for Vue app mounting.
- `src/App.vue`: Root container with `<RouterView>` and global providers (QueryClient).
- `src/router/index.ts`: Centralized Vue Router configuration.
- `src/views/`: Contains page-level components (formerly in `src/routes/`).
- `src/components/ui/`: Ported `.vue` components for the design system.

## 4. Migration Strategy: "Big Bang" Reconstruction
We will perform a surgical replacement of dependencies and then port files systematically.

### Phase 1: Dependency Overhaul
1. Remove React/TanStack Start dependencies.
2. Install Vue 3, Vue Router, Pinia, `@tanstack/vue-query`, and `radix-vue`.
3. Update `vite.config.ts` to use `@vitejs/plugin-vue`.

### Phase 2: Foundation
1. Create `src/main.ts` and `src/App.vue`.
2. Configure Vue Router with the 3 main routes: `/`, `/orders`, and `/orders/:orderId`.
3. Set up TanStack Query provider in `App.vue`.

### Phase 3: UI Component Porting
1. Port basic UI components (Button, Badge, Card, etc.) to `shadcn-vue`.
2. Port `SiteHeader.vue` and `StatusBadge.vue`.
3. Port `ShipmentMap.vue` (replacing `react-leaflet` with `leaflet`).

### Phase 4: Views & Logic
1. Port `HomeView.vue`.
2. Port `OrdersView.vue` using `useQuery` to fetch from `src/lib/orders.ts`.
3. Port `OrderDetailView.vue` using route params and `useQuery`.

## 5. Success Criteria
- The app builds successfully with Vite.
- All routes functional and visual parity with the React version.
- "Live Telemetry" on the map works correctly.
- No React dependencies remaining in `package.json`.
