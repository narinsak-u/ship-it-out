# Order Create/Edit Form Design

**Date:** 2026-05-22
**Project:** ship-it-out / Harbor Ops
**Status:** Approved

## Overview

Add order creation and editing capability to the Harbor Ops dashboard. A dedicated form page at `/orders/create` and `/orders/:orderId/edit` allows admins to create new shipments and edit existing ones. "Create new order" buttons are added to the Orders list and Carriers page hero sections.

## Routes

```
/orders/create          → OrderFormView.vue (create mode — blank form)
/orders/:orderId/edit   → OrderFormView.vue (edit mode — pre-populated form)
```

The existing `order-detail` route (`/orders/:orderId`) is not affected — the `/edit` segment makes these distinct routes.

## Architecture

A single `OrderFormView.vue` page determines create vs edit mode from the route. It renders a reusable `OrderForm.vue` component. New simulated API functions and Vue Query mutations handle the persistence.

### Component Tree

```
OrderFormView.vue
├── Page hero (title: "Create Order" or "Edit Order")
└── OrderForm.vue
    ├── Customer input
    ├── Origin input
    ├── Destination input
    ├── Carrier dropdown (from carriers data)
    ├── Weight input
    ├── Items input
    ├── Estimated delivery input
    ├── Status dropdown (edit only)
    └── Submit / Cancel buttons
```

### File Structure

```
src/
  views/
    OrderFormView.vue         # Page, detects mode from route
  components/
    OrderForm.vue             # Reusable form component
  lib/
    api/
      orders.ts               # Extended: createOrder, updateOrder
  hooks/
    useOrders.ts              # New: useCreateOrder, useUpdateOrder
```

### Modified Files

- `src/router/index.ts` — add create and edit routes
- `src/views/OrdersView.vue` — add "Create new order" button in hero right side
- `src/views/CarriersView.vue` — add "Create new order" button in hero right side

## Form Fields

### Create Mode

| Field | Type | Source | Required |
|---|---|---|---|
| Customer | text input | manual | yes |
| Origin | text input | manual | yes |
| Destination | text input | manual | yes |
| Carrier | dropdown | carriers list (from `carriers.ts`) | yes |
| Weight | text input | manual (e.g. "12.4 kg") | yes |
| Items | number input | manual | yes |
| Estimated Delivery | text input | manual (e.g. "May 25, 2026") | yes |

**Auto-generated:** id (ORD-10251... auto-incrementing), trackingNumber (TRK-XXXX-XXXX format with random hex characters, e.g. TRK-9F2A-44B1), createdAt (current date formatted like "May 22, 2026"), status ("pending"), progress (0), events ([]), originCoords/destinationCoords/currentCoords (default { lat: 0, lng: 0 }).

### Edit Mode

Same 7 fields as create, pre-populated from the existing order data. Plus:

| Field | Type | Source | Required |
|---|---|---|---|
| Status | dropdown | ShipmentStatus enum | yes |

Statuses available: pending, in_transit, out_for_delivery, delivered, delayed.

**Not editable:** id, trackingNumber, createdAt, progress, events, coords.

### Layout

Two-column grid on desktop (label left, field right), single column on mobile. Uses existing shadcn `Input` components and native `<select>` for dropdowns. Consistent with the project's dark theme styling.

## Data Flow

### Create

1. Admin clicks "Create new order" → `router.push({ name: 'order-create' })`
2. `OrderFormView` renders `OrderForm` with no `order` prop
3. Admin fills form → clicks "Create Order"
4. `OrderForm` validates required fields → emits `submit` with `OrderFormData`
5. `OrderFormView` calls `useCreateOrder().mutateAsync(data)`
6. `createOrder()` in API layer generates id, trackingNumber, createdAt, sets defaults, pushes to orders array
7. On success: invalidates queries → redirects to `/orders`
8. On error: inline error banner "Failed to create order"

### Edit

1. Admin clicks edit from order-detail page → `router.push({ name: 'order-edit', params: { orderId } })`
2. `OrderFormView` fetches order data, renders `OrderForm` with pre-populated values
3. Admin modifies fields → clicks "Save Changes"
4. `OrderForm` validates → emits `submit` with updated `OrderFormData`
5. `OrderFormView` calls `useUpdateOrder().mutateAsync({ id, data })`
6. `updateOrder()` patches the order in the array
7. On success: invalidates queries → redirects to `/orders/:orderId`
8. On error: inline error banner
9. If order not found for edit: 404 message with back link

## State & Error Handling

| Scenario | Behavior |
|---|---|
| Loading edit data | Skeleton placeholder |
| Order not found (edit) | 404 message + "Back to orders" link |
| Missing required field | Inline validation on submit attempt |
| Mutation error | Inline banner: "Failed to save order" |
| Create success | Redirect to /orders |
| Edit success | Redirect to /orders/:orderId |

## API Layer (extensions to existing)

New functions in `frontend/src/lib/api/orders.ts`:

```typescript
interface OrderFormData {
  customer: string;
  origin: string;
  destination: string;
  carrier: string;
  weight: string;
  items: number;
  estimatedDelivery: string;
  status?: ShipmentStatus;  // edit only
}

createOrder(data: OrderFormData): Promise<Order>
updateOrder(id: string, data: Partial<OrderFormData>): Promise<Order>
```

## Vue Query Hooks

New file `frontend/src/hooks/useOrders.ts`:

```typescript
useCreateOrder() — useMutation, invalidates ["orders"]
useUpdateOrder() — useMutation, invalidates ["orders"]
```

## Non-Goals

- No map picker for origin/destination coordinates (uses defaults)
- No multi-step wizard (single form page)
- No event/timeline editing on create (starts empty)
- No duplicate detection
- No batch creation
