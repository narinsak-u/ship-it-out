# Order Create/Edit Form Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add order creation and editing with a dedicated form page at `/orders/create` and `/orders/:orderId/edit`.

**Architecture:** `OrderFormView.vue` detects create vs edit from the route and renders a reusable `OrderForm.vue` component. New API functions and Vue Query mutations handle persistence via the existing in-memory data layer.

**Tech Stack:** Vue 3 Composition API, shadcn-vue (Input, Button, Skeleton), lucide-vue-next, TanStack Vue Query

---

### Task 1: Extend orders API layer

**Files:**
- Modify: `frontend/src/lib/api/orders.ts`

- [ ] **Step 1: Add `OrderFormData` type and `createOrder`/`updateOrder` functions**

Read `frontend/src/lib/api/orders.ts` and replace its content with:

```typescript
import { orders, type Order, type ShipmentStatus } from "@/lib/orders";

function delay(ms = 200): Promise<void> {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

export interface OrderFormData {
  customer: string;
  origin: string;
  destination: string;
  carrier: string;
  weight: string;
  items: number;
  estimatedDelivery: string;
  status?: ShipmentStatus;
}

function generateId(): string {
  const num = orders.length + 10251;
  return `ORD-${num}`;
}

function generateTrackingNumber(): string {
  const hex = () => Math.floor(Math.random() * 0x10000).toString(16).toUpperCase().padStart(4, "0");
  return `TRK-${hex()}-${hex()}`;
}

function today(): string {
  const d = new Date();
  const months = ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"];
  return `${months[d.getMonth()]} ${d.getDate()}, ${d.getFullYear()}`;
}

export async function fetchActiveDeliveries() {
  await delay();
  return orders.filter((o) => o.status !== "delivered");
}

export async function updateShipmentStatus(orderId: string, status: ShipmentStatus) {
  await delay(100);
  const order = orders.find((o) => o.id === orderId);
  if (!order) throw new Error("Order not found");
  order.status = status;
  if (status === "delivered") order.progress = 100;
  if (status === "in_transit") order.progress = Math.max(order.progress, 30);
  if (status === "out_for_delivery") order.progress = Math.max(order.progress, 70);
  return { ...order };
}

export async function createOrder(data: OrderFormData): Promise<Order> {
  await delay(200);
  const now = today();
  const order: Order = {
    id: generateId(),
    trackingNumber: generateTrackingNumber(),
    customer: data.customer,
    origin: data.origin,
    destination: data.destination,
    originCoords: { lat: 0, lng: 0 },
    destinationCoords: { lat: 0, lng: 0 },
    currentCoords: { lat: 0, lng: 0 },
    status: "pending",
    carrier: data.carrier,
    weight: data.weight,
    items: data.items,
    estimatedDelivery: data.estimatedDelivery,
    createdAt: now,
    progress: 0,
    events: [{ timestamp: now, location: data.origin, status: "Created", description: "Order created." }],
  };
  orders.unshift(order);
  return { ...order };
}

export async function updateOrder(id: string, data: Partial<OrderFormData>): Promise<Order> {
  await delay(200);
  const idx = orders.findIndex((o) => o.id === id);
  if (idx === -1) throw new Error("Order not found");
  const order = orders[idx];
  if (data.customer !== undefined) order.customer = data.customer;
  if (data.origin !== undefined) order.origin = data.origin;
  if (data.destination !== undefined) order.destination = data.destination;
  if (data.carrier !== undefined) order.carrier = data.carrier;
  if (data.weight !== undefined) order.weight = data.weight;
  if (data.items !== undefined) order.items = data.items;
  if (data.estimatedDelivery !== undefined) order.estimatedDelivery = data.estimatedDelivery;
  if (data.status !== undefined) order.status = data.status;
  return { ...order };
}
```

- [ ] **Step 2: Commit**

```bash
git add frontend/src/lib/api/orders.ts
git commit -m "feat: add createOrder and updateOrder API functions"
```

---

### Task 2: Create order mutation hooks

**File:**
- Create: `frontend/src/hooks/useOrders.ts`

- [ ] **Step 1: Create `useOrders.ts` with `useCreateOrder` and `useUpdateOrder`**

Create `frontend/src/hooks/useOrders.ts`:

```typescript
import { useMutation, useQueryClient } from "@tanstack/vue-query";
import { createOrder, updateOrder, type OrderFormData } from "@/lib/api/orders";

export function useCreateOrder() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: OrderFormData) => createOrder(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["deliveries"] });
    },
  });
}

export function useUpdateOrder() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<OrderFormData> }) =>
      updateOrder(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["deliveries"] });
    },
  });
}
```

- [ ] **Step 2: Commit**

```bash
git add frontend/src/hooks/useOrders.ts
git commit -m "feat: add useCreateOrder and useUpdateOrder hooks"
```

---

### Task 3: Create OrderForm component

**File:**
- Create: `frontend/src/components/OrderForm.vue`

- [ ] **Step 1: Create `OrderForm.vue`**

Create `frontend/src/components/OrderForm.vue`:

```typescript
<script setup lang="ts">
import { ref } from 'vue'
import { carriers as carrierList } from '@/lib/carriers'
import { statusLabels, type ShipmentStatus } from '@/lib/orders'
import type { OrderFormData } from '@/lib/api/orders'
import Input from '@/components/ui/Input.vue'
import Button from '@/components/ui/Button.vue'

const props = defineProps<{
  initial?: Partial<OrderFormData & { status?: ShipmentStatus }>
  isEditing?: boolean
  pending?: boolean
}>()

const emit = defineEmits<{
  submit: [data: OrderFormData & { status?: ShipmentStatus }]
}>()

const customer = ref(props.initial?.customer ?? '')
const origin = ref(props.initial?.origin ?? '')
const destination = ref(props.initial?.destination ?? '')
const carrier = ref(props.initial?.carrier ?? carrierList[0]?.name ?? '')
const weight = ref(props.initial?.weight ?? '')
const items = ref(props.initial?.items ?? 1)
const estimatedDelivery = ref(props.initial?.estimatedDelivery ?? '')
const status = ref<ShipmentStatus>(props.initial?.status ?? 'pending')

const errors = ref<Record<string, string>>({})

function validate(): boolean {
  const e: Record<string, string> = {}
  if (!customer.value.trim()) e.customer = 'Required'
  if (!origin.value.trim()) e.origin = 'Required'
  if (!destination.value.trim()) e.destination = 'Required'
  if (!carrier.value.trim()) e.carrier = 'Required'
  if (!weight.value.trim()) e.weight = 'Required'
  if (!items.value || items.value < 1) e.items = 'Must be at least 1'
  if (!estimatedDelivery.value.trim()) e.estimatedDelivery = 'Required'
  errors.value = e
  return Object.keys(e).length === 0
}

function handleSubmit() {
  if (!validate()) return
  emit('submit', {
    customer: customer.value,
    origin: origin.value,
    destination: destination.value,
    carrier: carrier.value,
    weight: weight.value,
    items: items.value,
    estimatedDelivery: estimatedDelivery.value,
    ...(props.isEditing ? { status: status.value } : {}),
  })
}
</script>

<template>
  <form @submit.prevent="handleSubmit" class="space-y-5">
    <div class="grid gap-5 md:grid-cols-2">
      <!-- Customer -->
      <div>
        <label class="font-mono text-xs uppercase tracking-widest text-muted-foreground">Customer</label>
        <Input v-model="customer" class="mt-1.5 font-mono text-sm" placeholder="e.g. Aria Nakamura" />
        <p v-if="errors.customer" class="mt-1 font-mono text-xs text-destructive">{{ errors.customer }}</p>
      </div>

      <!-- Carrier -->
      <div>
        <label class="font-mono text-xs uppercase tracking-widest text-muted-foreground">Carrier</label>
        <select
          v-model="carrier"
          class="mt-1.5 flex h-10 w-full rounded-lg border border-border bg-background px-3 font-mono text-sm"
        >
          <option v-for="c in carrierList" :key="c.id" :value="c.name">{{ c.name }}</option>
        </select>
        <p v-if="errors.carrier" class="mt-1 font-mono text-xs text-destructive">{{ errors.carrier }}</p>
      </div>

      <!-- Origin -->
      <div>
        <label class="font-mono text-xs uppercase tracking-widest text-muted-foreground">Origin</label>
        <Input v-model="origin" class="mt-1.5 font-mono text-sm" placeholder="e.g. Rotterdam, NL" />
        <p v-if="errors.origin" class="mt-1 font-mono text-xs text-destructive">{{ errors.origin }}</p>
      </div>

      <!-- Destination -->
      <div>
        <label class="font-mono text-xs uppercase tracking-widest text-muted-foreground">Destination</label>
        <Input v-model="destination" class="mt-1.5 font-mono text-sm" placeholder="e.g. Brooklyn, NY" />
        <p v-if="errors.destination" class="mt-1 font-mono text-xs text-destructive">{{ errors.destination }}</p>
      </div>

      <!-- Weight -->
      <div>
        <label class="font-mono text-xs uppercase tracking-widest text-muted-foreground">Weight</label>
        <Input v-model="weight" class="mt-1.5 font-mono text-sm" placeholder="e.g. 12.4 kg" />
        <p v-if="errors.weight" class="mt-1 font-mono text-xs text-destructive">{{ errors.weight }}</p>
      </div>

      <!-- Items -->
      <div>
        <label class="font-mono text-xs uppercase tracking-widest text-muted-foreground">Items</label>
        <Input v-model.number="items" type="number" min="1" class="mt-1.5 font-mono text-sm" />
        <p v-if="errors.items" class="mt-1 font-mono text-xs text-destructive">{{ errors.items }}</p>
      </div>

      <!-- Estimated Delivery -->
      <div>
        <label class="font-mono text-xs uppercase tracking-widest text-muted-foreground">Estimated Delivery</label>
        <Input v-model="estimatedDelivery" class="mt-1.5 font-mono text-sm" placeholder="e.g. May 25, 2026" />
        <p v-if="errors.estimatedDelivery" class="mt-1 font-mono text-xs text-destructive">{{ errors.estimatedDelivery }}</p>
      </div>

      <!-- Status (edit only) -->
      <div v-if="isEditing">
        <label class="font-mono text-xs uppercase tracking-widest text-muted-foreground">Status</label>
        <select
          v-model="status"
          class="mt-1.5 flex h-10 w-full rounded-lg border border-border bg-background px-3 font-mono text-sm"
        >
          <option v-for="(label, key) in statusLabels" :key="key" :value="key">{{ label }}</option>
        </select>
      </div>
    </div>

    <div class="flex justify-end gap-3 pt-4 border-t border-border">
      <Button variant="outline" type="button" @click="$router.back()">Cancel</Button>
      <Button type="submit" :disabled="pending">
        {{ pending ? 'Saving…' : isEditing ? 'Save Changes' : 'Create Order' }}
      </Button>
    </div>
  </form>
</template>
```

Note: The Cancel button uses `$router.back()` which requires `useRouter()` in the parent or `defineEmits`. Actually, since `OrderForm` is used inside `OrderFormView`, let's emit a `cancel` event instead.

Replace the Cancel button line with:

```typescript
      <Button variant="outline" type="button" @click="emit('cancel')">Cancel</Button>
```

And add to the emits:

```typescript
const emit = defineEmits<{
  submit: [data: OrderFormData & { status?: ShipmentStatus }]
  cancel: []
}>()
```

- [ ] **Step 2: Commit**

```bash
git add frontend/src/components/OrderForm.vue
git commit -m "feat: add reusable OrderForm component"
```

---

### Task 4: Create OrderFormView page

**File:**
- Create: `frontend/src/views/OrderFormView.vue`

- [ ] **Step 1: Create `OrderFormView.vue`**

Create `frontend/src/views/OrderFormView.vue`:

```typescript
<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getOrder } from '@/lib/orders'
import { useCreateOrder, useUpdateOrder } from '@/hooks/useOrders'
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
const order = computed(() => (orderId.value ? getOrder(orderId.value) : null))
const isPending = computed(() => createOrder.isPending.value || updateOrder.isPending.value)

async function handleSubmit(data: OrderFormData & { status?: ShipmentStatus }) {
  if (isEditing.value && orderId.value) {
    await updateOrder.mutateAsync({ id: orderId.value, data })
    router.push({ name: 'order-detail', params: { orderId: orderId.value } })
  } else {
    const created = await createOrder.mutateAsync(data)
    router.push({ name: 'orders' })
  }
}

function handleCancel() {
  if (isEditing.value && orderId.value) {
    router.push({ name: 'order-detail', params: { orderId: orderId.value } })
  } else {
    router.push({ name: 'orders' })
  }
}
</script>

<template>
  <div>
    <section class="border-b border-border bg-gradient-hero">
      <div class="mx-auto max-w-7xl px-6 py-14">
        <span class="font-mono text-xs uppercase tracking-widest text-primary">/ orders</span>
        <h1 class="mt-3 text-4xl font-semibold tracking-tight md:text-5xl">
          {{ isEditing ? 'Edit Order' : 'Create Order' }}
        </h1>
        <p class="mt-3 max-w-2xl text-muted-foreground">
          {{ isEditing ? 'Update shipment details and status.' : 'Register a new shipment in the system.' }}
        </p>
      </div>
    </section>

    <section class="mx-auto max-w-3xl px-6 py-10">
      <div v-if="isEditing && !order" class="py-16 text-center">
        <p class="font-mono text-lg">Order not found</p>
        <p class="mt-2 text-sm text-muted-foreground">No order matches ID "{{ orderId }}".</p>
        <button
          @click="router.push({ name: 'orders' })"
          class="mt-6 font-mono text-sm text-primary hover:underline"
        >
          ← Back to orders
        </button>
      </div>

      <div v-else-if="isEditing && !order">
        <Skeleton class="h-96 rounded-xl" />
      </div>

      <div v-else class="rounded-xl border border-border bg-card p-6 shadow-elegant">
        <OrderForm
          :initial="order ?? undefined"
          :is-editing="isEditing"
          :pending="isPending"
          @submit="handleSubmit"
          @cancel="handleCancel"
        />
      </div>

      <div v-if="createOrder.isError || updateOrder.isError" class="mt-4 rounded-lg bg-destructive/15 px-4 py-3 font-mono text-sm text-destructive">
        Failed to save order. Please try again.
      </div>
    </section>
  </div>
</template>
```

- [ ] **Step 2: Commit**

```bash
git add frontend/src/views/OrderFormView.vue
git commit -m "feat: add OrderFormView page for create and edit"
```

---

### Task 5: Routes + "Create new order" buttons

**Files:**
- Modify: `frontend/src/router/index.ts`
- Modify: `frontend/src/views/OrdersView.vue`
- Modify: `frontend/src/views/CarriersView.vue`

- [ ] **Step 1: Add create and edit routes**

In `frontend/src/router/index.ts`, add before the `order-detail` route (order matters — literal `create` must come before the `:orderId` param):

```typescript
    {
      path: '/orders/create',
      name: 'order-create',
      component: () => import('@/views/OrderFormView.vue'),
    },
    {
      path: '/orders/:orderId/edit',
      name: 'order-edit',
      component: () => import('@/views/OrderFormView.vue'),
    },
```

Add them BEFORE the `order-detail` route (at line 16-17) so `/orders/create` matches before `/orders/:orderId`.

- [ ] **Step 2: Add "Create new order" button to OrdersView**

In `frontend/src/views/OrdersView.vue`, add `Plus` to the lucide imports:

```typescript
import { Search, Filter, ArrowRight, Plus } from 'lucide-vue-next'
```

Add `Button` to the imports:

```typescript
import Button from '@/components/ui/Button.vue'
```

In the hero section, after the description `<p>` and before closing the hero `<section>`, add a button on the right side. Change the hero div from:

```html
      <div class="mx-auto max-w-7xl px-6 py-14">
        <span class="font-mono text-xs uppercase tracking-widest text-primary">/ orders</span>
        <h1 class="mt-3 text-4xl font-semibold tracking-tight md:text-5xl">Shipment manifest</h1>
        <p class="mt-3 max-w-2xl text-muted-foreground">
          {{ orders.length }} total shipments tracked across all carriers.
        </p>
      </div>
```

To:

```html
      <div class="mx-auto max-w-7xl px-6 py-14">
        <div class="flex items-start justify-between">
          <div>
            <span class="font-mono text-xs uppercase tracking-widest text-primary">/ orders</span>
            <h1 class="mt-3 text-4xl font-semibold tracking-tight md:text-5xl">Shipment manifest</h1>
            <p class="mt-3 max-w-2xl text-muted-foreground">
              {{ orders.length }} total shipments tracked across all carriers.
            </p>
          </div>
          <RouterLink
            :to="{ name: 'order-create' }"
            class="hidden shrink-0 md:block"
          >
            <Button class="gap-2">
              <Plus class="h-4 w-4" /> New Order
            </Button>
          </RouterLink>
        </div>
      </div>
```

- [ ] **Step 3: Add "Create new order" button to CarriersView**

In `frontend/src/views/CarriersView.vue`, add `Plus` and `RouterLink` to the imports:

```typescript
import { Truck, Warehouse, BarChart3, Package, Plus } from 'lucide-vue-next'
import { RouterLink } from 'vue-router'
```

Import `Button`:

```typescript
import Button from '@/components/ui/Button.vue'
```

In the hero section, wrap the content in a flex row and add the button:

```html
      <div class="mx-auto max-w-7xl px-6 py-14">
        <div class="flex items-start justify-between">
          <div>
            <span class="font-mono text-xs uppercase tracking-widest text-primary">/ carriers</span>
            <h1 class="mt-3 text-4xl font-semibold tracking-tight md:text-5xl">Carrier operations</h1>
            <p class="mt-3 max-w-2xl text-muted-foreground">
              Manage carrier drivers, hubs, and monitor active deliveries across the fleet.
            </p>
          </div>
          <RouterLink
            :to="{ name: 'order-create' }"
            class="hidden shrink-0 md:block"
          >
            <Button class="gap-2">
              <Plus class="h-4 w-4" /> New Order
            </Button>
          </RouterLink>
        </div>
      </div>
```

- [ ] **Step 4: Commit**

```bash
git add frontend/src/router/index.ts frontend/src/views/OrdersView.vue frontend/src/views/CarriersView.vue
git commit -m "feat: add order routes and create order buttons"
```

---

### Task 6: Build verification

- [ ] **Step 1: Run typecheck and build**

```bash
cd frontend && npm run build
```

Expected: clean build with no errors. If there are type issues, fix them and re-run.

- [ ] **Step 2: If build passes, commit any remaining fixes**

```bash
git add -A
git commit -m "chore: fix type issues from order form build"
```
