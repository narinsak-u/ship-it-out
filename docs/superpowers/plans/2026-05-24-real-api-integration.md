# Real API Integration Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development or superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Replace all mock/in-memory API functions for orders and hubs with real HTTP calls to the Go/Fiber backend.

**Architecture:** The frontend has a mock API layer (`src/lib/api/orders.ts`, `src/lib/api/carriers.ts`) that operates on in-memory arrays. A hooks layer wraps them with TanStack VueQuery. An existing `api` client in `client.ts` (native fetch, `credentials: 'include'`) is already used by the auth store. We replace the mock functions with HTTP calls via the `api` client, add mapper functions to convert backend response shapes to frontend types, and update components that consume in-memory data directly.

**Tech Stack:** Vue 3, TanStack VueQuery, native fetch, Go/Fiber backend

**Key constraints:**
- Keep existing hook signatures unchanged — functions return the same `Order[]` / `Hub[]` / `Order` types
- Backend Shipment `id` is `uint` (number in JSON) → frontend `Order.id` is `string`; map via `String()`
- Backend dates are RFC3339 → frontend uses formatted strings like "May 24, 2026"; map in a shared formatter
- Carrier/driver mock data stays as-is (no backend endpoints for those yet)
- All auth-gated routes use cookie-based auth already handled by `client.ts`

---

### File Inventory

**Backend — modify:**
- `cmd/server/main.go` — register new PUT + DELETE shipment routes
- `internal/shipment/handler.go` — add `Update` and `Delete` handlers

**Frontend — create:**
- `src/lib/api/mappers.ts` — mapper functions (backend → frontend type conversion + date formatting)

**Frontend — modify:**
- `src/lib/api/client.ts` — add `put` and `patch` methods
- `src/lib/api/orders.ts` — rewrite all functions to use `api` client + mapper
- `src/lib/api/carriers.ts` — rewrite hub functions to use `api` client; keep carrier/driver mock functions
- `src/views/OrdersView.vue` — replace direct in-memory `orders` with VueQuery
- `src/views/OrderDetailView.vue` — fetch shipment + events from API
- `src/views/HomeView.vue` — fetch shipments from API
- `src/views/OrderFormView.vue` — fetch order from API for edit mode
- `src/components/DeliveriesPanel.vue` — fix `handleAssign` to use query data instead of in-memory `orders`

---

### Task 1: Backend — Add PUT and DELETE endpoints for shipments

**Files:**
- Modify: `backend/internal/shipment/handler.go`
- Modify: `backend/cmd/server/main.go`

**Why:** The frontend's `updateOrder()` and `deleteOrder()` need real endpoints. The backend currently has List, Create, GetByID, and UpdateStatus — but no general update or delete.

**Changes to `internal/shipment/handler.go`:**

Add an `UpdateRequest` struct (reuses same fields as `CreateRequest` plus optional `estimatedDelivery` and `status`) and two new handlers at the bottom of the file:

```go
// UpdateRequest is the JSON body for updating an existing shipment.
// Same fields as CreateRequest plus optional estimatedDelivery and status.
type UpdateRequest struct {
	Customer          *models.ContactInfo `json:"customer,omitempty"`
	Receiver          *models.ContactInfo `json:"receiver,omitempty"`
	Carrier           *string             `json:"carrier,omitempty"`
	Weight            *string             `json:"weight,omitempty"`
	Items             *int                `json:"items,omitempty"`
	EstimatedDelivery *time.Time          `json:"estimatedDelivery,omitempty"`
	Status            *string             `json:"status,omitempty"`
}

// Update modifies an existing shipment's fields.
func Update(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return utils.Error(c, 400, "invalid id")
	}

	var req UpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, 400, "invalid request body")
	}

	var shipment models.Shipment
	if result := database.DB.First(&shipment, id); result.Error != nil {
		return utils.Error(c, 404, "shipment not found")
	}

	if req.Customer != nil {
		shipment.Customer = *req.Customer
		shipment.Origin = composeAddress(*req.Customer)
	}
	if req.Receiver != nil {
		shipment.Receiver = *req.Receiver
		shipment.Destination = composeAddress(*req.Receiver)
	}
	if req.Carrier != nil {
		shipment.Carrier = *req.Carrier
	}
	if req.Weight != nil {
		shipment.Weight = *req.Weight
	}
	if req.Items != nil {
		shipment.Items = *req.Items
	}
	if req.EstimatedDelivery != nil {
		shipment.EstimatedDelivery = *req.EstimatedDelivery
	}
	if req.Status != nil {
		shipment.Status = *req.Status
	}

	database.DB.Save(&shipment)
	return utils.Success(c, shipment)
}

// Delete removes a shipment and its events from the database.
func Delete(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return utils.Error(c, 400, "invalid id")
	}

	var shipment models.Shipment
	if result := database.DB.First(&shipment, id); result.Error != nil {
		return utils.Error(c, 404, "shipment not found")
	}

	// Delete associated events first (FK constraint)
	database.DB.Where("shipment_id = ?", shipment.ID).Delete(&models.ShipmentEvent{})
	database.DB.Delete(&shipment)

	return utils.Success(c, fiber.Map{"message": "shipment deleted"})
}
```

**Changes to `cmd/server/main.go`:**

Add two new routes in the shipment group, after the existing routes:

```go
shipmentGroup.Put("/:id", shipment.Update)    // PUT /api/shipments/:id
shipmentGroup.Delete("/:id", shipment.Delete) // DELETE /api/shipments/:id
```

---

### Task 2: Frontend client — Add `put` and `patch` methods

**File:** `frontend/src/lib/api/client.ts`

**Why:** The current `api` object only has `get`, `post`, and `del`. We need `put` and `patch` for the update endpoints.

```ts
export const api = {
  get: <T = unknown>(path: string) => request<T>(path),
  post: <T = unknown>(path: string, body?: unknown) =>
    request<T>(path, { method: 'POST', body: body ? JSON.stringify(body) : undefined }),
  put: <T = unknown>(path: string, body?: unknown) =>
    request<T>(path, { method: 'PUT', body: body ? JSON.stringify(body) : undefined }),
  patch: <T = unknown>(path: string, body?: unknown) =>
    request<T>(path, { method: 'PATCH', body: body ? JSON.stringify(body) : undefined }),
  del: <T = unknown>(path: string) => request<T>(path, { method: 'DELETE' }),
};
```

---

### Task 3: Create mapper utilities

**File:** `frontend/src/lib/api/mappers.ts`

**Why:** The backend returns `Shipment` objects with numeric IDs and RFC3339 dates. The frontend `Order` type uses string IDs and formatted date strings. We need mappers to convert between them.

```ts
import type { Order, TrackingEvent, ShipmentStatus, GeoPoint, ContactInfo } from '@/lib/orders'
import type { Hub, HubStatus } from '@/lib/carriers'

// ---- Backend response types (mirrors Go structs) ----

export interface BackendShipment {
  id: number
  trackingNumber: string
  customer: ContactInfo
  receiver: ContactInfo
  currentCoords: GeoPoint
  origin: string
  destination: string
  status: string
  carrier: string
  driverId?: string
  weight: string
  items: number
  estimatedDelivery: string
  createdAt: string
  progress: number
}

export interface BackendShipmentEvent {
  id: number
  shipmentId: number
  status: string
  location: { name: string; lat: number; lng: number }
  description?: string
  timestamp: string
}

export interface BackendHub {
  id: string
  name: string
  carrierId: string
  address: string
  coords: GeoPoint
  capacity: number
  currentUtilization: number
  status: string
  createdAt: string
}

// ---- Date formatters ----

const MONTHS = ['Jan','Feb','Mar','Apr','May','Jun','Jul','Aug','Sep','Oct','Nov','Dec']

export function formatDate(iso: string): string {
  const d = new Date(iso)
  return `${MONTHS[d.getMonth()]} ${d.getDate()}, ${d.getFullYear()}`
}

export function formatTimestamp(iso: string): string {
  const d = new Date(iso)
  const hh = String(d.getHours()).padStart(2, '0')
  const mm = String(d.getMinutes()).padStart(2, '0')
  return `${MONTHS[d.getMonth()]} ${d.getDate()}, ${hh}:${mm}`
}

// ---- Mappers ----

export function mapShipmentToOrder(s: BackendShipment): Order {
  return {
    id: String(s.id),
    trackingNumber: s.trackingNumber,
    customer: s.customer,
    receiver: s.receiver,
    origin: s.origin,
    destination: s.destination,
    currentCoords: s.currentCoords,
    status: s.status as ShipmentStatus,
    carrier: s.carrier,
    driverId: s.driverId,
    weight: s.weight,
    items: s.items,
    estimatedDelivery: formatDate(s.estimatedDelivery),
    createdAt: formatDate(s.createdAt),
    progress: s.progress,
    events: [],
  }
}

export function mapEventToTrackingEvent(e: BackendShipmentEvent): TrackingEvent {
  return {
    timestamp: formatTimestamp(e.timestamp),
    location: e.location,
    status: e.status,
    description: e.description ?? '',
  }
}

export function mapBackendHubToHub(h: BackendHub): Hub {
  return {
    id: h.id,
    name: h.name,
    carrierId: h.carrierId,
    address: h.address,
    coords: h.coords,
    capacity: h.capacity,
    currentUtilization: h.currentUtilization,
    status: h.status as HubStatus,
  }
}
```

---

### Task 4: Rewrite `src/lib/api/orders.ts` with real API calls

**File:** `frontend/src/lib/api/orders.ts`

**Why:** Replace all mock functions with real HTTP calls. Keep identical export signatures so hooks don't need to change.

```ts
import { api } from '@/lib/api/client'
import { mapShipmentToOrder, mapEventToTrackingEvent, type BackendShipment, type BackendShipmentEvent } from '@/lib/api/mappers'
import type { Order, ShipmentStatus, ContactInfo } from '@/lib/orders'

export interface OrderFormData {
  customer: ContactInfo
  receiver: ContactInfo
  carrier: string
  weight: string
  items: number
  estimatedDelivery: string
  status?: ShipmentStatus
}

export async function fetchActiveDeliveries(): Promise<Order[]> {
  const result = await api.get<BackendShipment[]>('/shipments')
  if (result.error) throw new Error(result.error)
  return result.data
    .filter((s) => s.status !== 'delivered')
    .map(mapShipmentToOrder)
}

export async function updateShipmentStatus(orderId: string, status: ShipmentStatus): Promise<Order> {
  const result = await api.patch<BackendShipment>(`/shipments/${orderId}/status`, { status })
  if (result.error) throw new Error(result.error)
  return mapShipmentToOrder(result.data)
}

export async function createOrder(data: OrderFormData): Promise<Order> {
  // Backend CreateRequest doesn't accept estimatedDelivery (it's auto-generated)
  const { estimatedDelivery: _, ...body } = data
  const result = await api.post<BackendShipment>('/shipments', body)
  if (result.error) throw new Error(result.error)
  return mapShipmentToOrder(result.data)
}

export async function updateOrder(id: string, data: Partial<OrderFormData>): Promise<Order> {
  const result = await api.put<BackendShipment>(`/shipments/${id}`, data)
  if (result.error) throw new Error(result.error)
  return mapShipmentToOrder(result.data)
}

export async function deleteOrder(id: string): Promise<void> {
  const result = await api.del(`/shipments/${id}`)
  if (result.error) throw new Error(result.error)
}

export async function fetchOrder(id: string): Promise<Order> {
  const result = await api.get<BackendShipment>(`/shipments/${id}`)
  if (result.error) throw new Error(result.error)
  return mapShipmentToOrder(result.data)
}

export async function fetchOrderEvents(trackingNumber: string): Promise<TrackingEvent[]> {
  const result = await api.get<{ shipment: BackendShipment; events: BackendShipmentEvent[] }>(`/track/${trackingNumber}`)
  if (result.error) throw new Error(result.error)
  return result.data.events.map(mapEventToTrackingEvent)
}
```

Also export `TrackingEvent` type at the top:
```ts
import type { Order, ShipmentStatus, TrackingEvent, ContactInfo } from '@/lib/orders'
```

---

### Task 5: Rewrite hub functions in `src/lib/api/carriers.ts` with real API calls

**File:** `frontend/src/lib/api/carriers.ts`

**Why:** Replace hub mock functions with real HTTP calls. Keep carrier/driver mock functions unchanged. Keep all existing exports.

Import and add at the top:
```ts
import { api } from '@/lib/api/client'
import { mapBackendHubToHub, type BackendHub } from '@/lib/api/mappers'
```

Replace the `fetchHubs`, `createHub`, `updateHub`, `deleteHub` functions while keeping the old carrier/driver mock functions untouched:

```ts
// ---- Hub functions (real API) ----

export async function fetchHubs(): Promise<Hub[]> {
  const result = await api.get<BackendHub[]>('/hubs')
  if (result.error) throw new Error(result.error)
  return result.data.map(mapBackendHubToHub)
}

export async function createHub(data: Omit<Hub, 'id'>): Promise<Hub> {
  const result = await api.post<BackendHub>('/hubs', data)
  if (result.error) throw new Error(result.error)
  return mapBackendHubToHub(result.data)
}

export async function updateHub(id: string, data: Partial<Hub>): Promise<Hub> {
  const result = await api.put<BackendHub>(`/hubs/${id}`, data)
  if (result.error) throw new Error(result.error)
  return mapBackendHubToHub(result.data)
}

export async function deleteHub(id: string): Promise<void> {
  const result = await api.del(`/hubs/${id}`)
  if (result.error) throw new Error(result.error)
}
```

---

### Task 6: Update OrdersView to use VueQuery

**File:** `frontend/src/views/OrdersView.vue`

**Why:** Currently uses direct in-memory `orders` import. Need to fetch from API.

Replace the `<script>` block:

```vue
<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query'
import { Search, Filter, ArrowRight, Plus, Pencil, Trash2 } from 'lucide-vue-next'
import Input from '@/components/ui/Input.vue'
import StatusBadge from '@/components/StatusBadge.vue'
import { statusLabels, type ShipmentStatus } from '@/lib/orders'
import { fetchActiveDeliveries, deleteOrder } from '@/lib/api/orders'
import Button from '@/components/ui/Button.vue'
import { cn } from '@/lib/utils'
import { useAuthStore } from '@/stores/auth'
import AuthModal from '@/components/AuthModal.vue'

const authStore = useAuthStore()
const showAuthModal = ref(false)
const router = useRouter()
const queryClient = useQueryClient()

const { data: orders, isLoading } = useQuery({
  queryKey: ['orders'],
  queryFn: fetchActiveDeliveries,
})

const deleteMutation = useMutation({
  mutationFn: (id: string) => deleteOrder(id),
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['orders'] })
  },
})

const FILTERS: Array<{ key: ShipmentStatus | 'all'; label: string }> = [
  { key: 'all', label: 'All' },
  { key: 'pending', label: 'Pending' },
  { key: 'in_transit', label: 'In Transit' },
  { key: 'out_for_delivery', label: 'Out for Delivery' },
  { key: 'delivered', label: 'Delivered' },
  { key: 'delayed', label: 'Delayed' },
]

const filter = ref<ShipmentStatus | 'all'>('all')
const query = ref('')

const filtered = computed(() => {
  const q = query.value.trim().toLowerCase()
  return (orders.value ?? []).filter((o) => {
    if (filter.value !== 'all' && o.status !== filter.value) return false
    if (!q) return true
    return (
      o.id.toLowerCase().includes(q) ||
      o.trackingNumber.toLowerCase().includes(q) ||
      o.customer.name.toLowerCase().includes(q) ||
      o.destination.toLowerCase().includes(q)
    )
  })
})

function onAuthenticated() {
  showAuthModal.value = false
  router.push({ name: 'order-create' })
}

function onGuest() {
  showAuthModal.value = false
  router.push({ name: 'order-create' })
}
</script>
```

Template changes needed:
1. The heading count: `orders.length` → `orders?.length ?? 0`
2. The "Showing X of Y" line at bottom
3. Delete button handler: `deleteOrder(o.id)` → `deleteMutation.mutate(o.id)`
4. Add loading state

Replace:
- `{{ orders.length }} total shipments` → `{{ orders?.length ?? 0 }} total shipments`
- `@click.stop="deleteOrder(o.id)"` → `@click.stop="deleteMutation.mutate(o.id)"`
- `Showing {{ filtered.length }} of {{ orders.length }}` → `Showing {{ filtered.length }} of {{ orders?.length ?? 0 }}`
- Add `v-if="isLoading"` skeleton state before the main content

---

### Task 7: Update OrderDetailView to fetch from API

**File:** `frontend/src/views/OrderDetailView.vue`

**Why:** Currently reads from in-memory `getOrder(orderId)`. Need to fetch shipment and events from backend.

Replace the `<script>` block:

```vue
<script setup lang="ts">
import { ref, computed, onMounted, defineAsyncComponent } from 'vue'
import { useRoute, RouterLink } from 'vue-router'
import { useQuery } from '@tanstack/vue-query'
import {
  ArrowLeft, MapPin, Truck, Calendar, Hash, User, Weight, Maximize2
} from 'lucide-vue-next'
import StatusBadge from '@/components/StatusBadge.vue'
import Skeleton from '@/components/ui/Skeleton.vue'
import { fetchOrder, fetchOrderEvents } from '@/lib/api/orders'

const ShipmentMap = defineAsyncComponent(() => import('@/components/ShipmentMap.vue'))

const route = useRoute()
const orderId = route.params.orderId as string

const { data: order, isLoading } = useQuery({
  queryKey: ['order', orderId],
  queryFn: () => fetchOrder(orderId),
})

const { data: events } = useQuery({
  queryKey: ['order-events', orderId],
  queryFn: () => {
    if (!order.value) return []
    return fetchOrderEvents(order.value.trackingNumber)
  },
  enabled: computed(() => !!order.value),
})

const mounted = ref(false)
onMounted(() => {
  mounted.value = true
})

const meta = computed(() => {
  const o = order.value
  if (!o) return []
  return [
    { icon: Hash, label: 'Tracking #', value: o.trackingNumber },
    { icon: User, label: 'Customer', value: o.customer.name },
    { icon: Truck, label: 'Carrier', value: o.carrier },
    { icon: Weight, label: 'Weight', value: o.weight },
    { icon: Calendar, label: 'Created', value: o.createdAt },
  ]
})
</script>
```

Template changes needed:
- `v-if="order"` stays the same (it's now reactive via `useQuery`)
- `v-for="e in order.events"` → `v-for="e in events ?? []"`
- Add `v-if="isLoading"` skeleton state

---

### Task 8: Update HomeView to fetch from API

**File:** `frontend/src/views/HomeView.vue`

**Why:** Currently imports `orders` directly from `@/lib/orders`.

Replace the direct import with a VueQuery fetch:

```vue
<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useQuery } from '@tanstack/vue-query'
import { ArrowRight, Boxes, Search, Truck, Globe2, Activity } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
import Input from '@/components/ui/Input.vue'
import StatusBadge from '@/components/StatusBadge.vue'
import { fetchActiveDeliveries } from '@/lib/api/orders'
import Skeleton from '@/components/ui/Skeleton.vue'

const router = useRouter()

const { data: orders } = useQuery({
  queryKey: ['orders'],
  queryFn: fetchActiveDeliveries,
})

const query = ref('')

const onTrack = (e: Event) => {
  e.preventDefault()
  if (!orders.value) return
  const q = query.value.trim().toLowerCase()
  const match = orders.value.find(
    (o) => o.id.toLowerCase() === q || o.trackingNumber.toLowerCase() === q,
  )
  if (match) {
    router.push({ name: 'order-detail', params: { orderId: match.id } })
  } else {
    router.push({ name: 'orders' })
  }
}

const stats = computed(() => {
  const all = orders.value ?? []
  return [
    { label: 'Active shipments', value: all.filter((o) => o.status !== 'delivered').length, icon: Truck },
    { label: 'Delivered (30d)', value: 184, icon: Boxes },
    { label: 'On-time rate', value: '97.4%', icon: Activity },
    { label: 'Countries served', value: 42, icon: Globe2 },
  ]
})

const recent = computed(() => (orders.value ?? []).slice(0, 3))
</script>
```

Template changes: the tracking placeholder hint needs updating since the seeded tracking numbers are now different:
- `placeholder="Enter tracking # (try TRK-9F2A-44B1)"` → `placeholder="Enter tracking # (e.g. TH2026...)"`

---

### Task 9: Update OrderFormView to fetch from API

**File:** `frontend/src/views/OrderFormView.vue`

**Why:** Currently uses `getOrder(orderId)` from in-memory data for edit mode.

Replace the direct import with a `useQuery`:

```vue
<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useQuery } from '@tanstack/vue-query'
import { useCreateOrder, useUpdateOrder } from '@/hooks/useOrders'
import { fetchOrder } from '@/lib/api/orders'
import type { OrderFormData } from '@/lib/api/orders'
import type { ShipmentStatus } from '@/lib/orders'
import OrderForm from '@/components/OrderForm.vue'
import Skeleton from '@/components/ui/Skeleton.vue'

const route = useRoute()
const router = useRouter()
const createOrder = useCreateOrder()
const updateOrder = useUpdateOrder()

const isEditing = computed(() => !!route.params.orderId)
const orderId = computed(() => route.params.orderId as string | undefined)

const { data: order } = useQuery({
  queryKey: ['order', orderId.value],
  queryFn: () => fetchOrder(orderId.value!),
  enabled: isEditing,
})
</script>
```

---

### Task 10: Fix DeliveriesPanel to use query data instead of in-memory `orders`

**File:** `frontend/src/components/DeliveriesPanel.vue`

**Why:** `handleAssign` uses the in-memory `orders` array from `@/lib/orders` to find the order. Should use `deliveries.value` (the query data) instead.

Remove: `import { orders, statusLabels, type ShipmentStatus } from '@/lib/orders'`
Add: `import { statusLabels, type ShipmentStatus } from '@/lib/orders'`

Change the `handleAssign` function:
```ts
function handleAssign(orderId: string) {
  const order = deliveries.value?.find((o) => o.id === orderId)
  if (!order) return
  const available = driverData.filter((d) => d.status === 'available')
  if (available.length === 0) return
  assignDriver.mutate({ driverId: available[0].id, orderId })
}
```
