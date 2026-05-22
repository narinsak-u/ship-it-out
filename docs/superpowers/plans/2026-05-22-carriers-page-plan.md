# Carriers Page Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a Carriers management dashboard at `/carriers` with tabbed panels for drivers, hubs, analytics, and active deliveries.

**Architecture:** Modular tab panels inside `CarriersView.vue`. Each tab is a separate async component. Data lives in a new `src/lib/carriers.ts` module with types + mock data. TanStack Vue Query hooks provide data access with simulated API functions. No test framework installed yet.

**Tech Stack:** Vue 3 (Composition API, `<script setup lang="ts">`), TanStack Vue Query, shadcn-vue (Card, Table, Skeleton, Input, Button, Badge), Tailwind CSS v4, lucide-vue-next icons.

---

### Task 1: Extend Order type + create carriers data module

**Files:**
- Modify: `frontend/src/lib/orders.ts`
- Create: `frontend/src/lib/carriers.ts`

- [ ] **Step 1: Add `driverId` to Order type in `orders.ts`**

In `frontend/src/lib/orders.ts`, add `driverId?: string` to the `Order` interface:

```typescript
export interface Order {
  id: string;
  trackingNumber: string;
  customer: string;
  destination: string;
  origin: string;
  originCoords: GeoPoint;
  destinationCoords: GeoPoint;
  currentCoords: GeoPoint;
  status: ShipmentStatus;
  carrier: string;
  driverId?: string;
  weight: string;
  items: number;
  estimatedDelivery: string;
  createdAt: string;
  progress: number;
  events: TrackingEvent[];
}
```

- [ ] **Step 2: Create `carriers.ts` with types and mock data**

Create `frontend/src/lib/carriers.ts`:

```typescript
import type { GeoPoint } from "@/lib/orders";

export type CarrierStatus = "active" | "inactive";
export type DriverStatus = "available" | "on_delivery" | "off_duty";
export type HubStatus = "active" | "maintenance" | "closed";

export interface Carrier {
  id: string;
  name: string;
  contactEmail: string;
  phone: string;
  status: CarrierStatus;
  fleetSize: number;
  totalDrivers: number;
  totalHubs: number;
  createdAt: string;
}

export interface Driver {
  id: string;
  name: string;
  carrierId: string;
  phone: string;
  email: string;
  status: DriverStatus;
  vehicleInfo: string;
}

export interface Hub {
  id: string;
  name: string;
  carrierId: string;
  address: string;
  coords: GeoPoint;
  capacity: number;
  currentUtilization: number;
  status: HubStatus;
}

export const carriers: Carrier[] = [
  { id: "CAR-001", name: "Pacific Freight", contactEmail: "ops@pacificfreight.com", phone: "+1-555-0101", status: "active", fleetSize: 42, totalDrivers: 18, totalHubs: 6, createdAt: "2024-03-15" },
  { id: "CAR-002", name: "Skyline Express", contactEmail: "dispatch@skylineexpress.io", phone: "+1-555-0102", status: "active", fleetSize: 28, totalDrivers: 12, totalHubs: 4, createdAt: "2024-06-01" },
  { id: "CAR-003", name: "Trans-Atlantic Cargo", contactEmail: "support@tacargo.com", phone: "+1-555-0103", status: "active", fleetSize: 35, totalDrivers: 15, totalHubs: 5, createdAt: "2024-01-20" },
  { id: "CAR-004", name: "Nordic Lines", contactEmail: "ops@nordiclines.no", phone: "+47-555-0104", status: "active", fleetSize: 15, totalDrivers: 6, totalHubs: 3, createdAt: "2024-08-10" },
  { id: "CAR-005", name: "Gulf Logistics", contactEmail: "info@gulflogistics.ae", phone: "+971-555-0105", status: "active", fleetSize: 20, totalDrivers: 8, totalHubs: 3, createdAt: "2025-01-05" },
  { id: "CAR-006", name: "Mediterranean Freight", contactEmail: "ops@medfreight.it", phone: "+39-555-0106", status: "inactive", fleetSize: 10, totalDrivers: 4, totalHubs: 2, createdAt: "2025-02-18" },
];

export const drivers: Driver[] = [
  { id: "DRV-001", name: "Elena Voss", carrierId: "CAR-001", phone: "+1-555-1001", email: "elena.voss@pacificfreight.com", status: "available", vehicleInfo: "VAN-442" },
  { id: "DRV-002", name: "Marco Reyes", carrierId: "CAR-002", phone: "+1-555-1002", email: "marco.reyes@skylineexpress.io", status: "on_delivery", vehicleInfo: "TRK-991" },
  { id: "DRV-003", name: "Wei Chen", carrierId: "CAR-005", phone: "+1-555-1003", email: "wei.chen@gulflogistics.ae", status: "off_duty", vehicleInfo: "BOX-207" },
  { id: "DRV-004", name: "Sofia Lambert", carrierId: "CAR-001", phone: "+1-555-1004", email: "sofia.lambert@pacificfreight.com", status: "available", vehicleInfo: "VAN-443" },
  { id: "DRV-005", name: "Lars Johansson", carrierId: "CAR-004", phone: "+47-555-1005", email: "lars.j@nordiclines.no", status: "on_delivery", vehicleInfo: "TRK-007" },
  { id: "DRV-006", name: "Amara Obi", carrierId: "CAR-003", phone: "+1-555-1006", email: "amara.obi@tacargo.com", status: "available", vehicleInfo: "VAN-881" },
  { id: "DRV-007", name: "Dimitri Pavlov", carrierId: "CAR-002", phone: "+1-555-1007", email: "dimitri.pavlov@skylineexpress.io", status: "available", vehicleInfo: "TRK-992" },
  { id: "DRV-008", name: "Camila Rojas", carrierId: "CAR-006", phone: "+39-555-1008", email: "camila.rojas@medfreight.it", status: "available", vehicleInfo: "BOX-104" },
];

export const hubs: Hub[] = [
  { id: "HUB-001", name: "Rotterdam Hub", carrierId: "CAR-001", address: "Haven 12, Rotterdam, NL", coords: { lat: 51.9244, lng: 4.4777 }, capacity: 5000, currentUtilization: 3400, status: "active" },
  { id: "HUB-002", name: "Berlin Hub", carrierId: "CAR-002", address: "Industriestr. 45, Berlin, DE", coords: { lat: 52.52, lng: 13.405 }, capacity: 3000, currentUtilization: 2100, status: "active" },
  { id: "HUB-003", name: "Lisbon Hub", carrierId: "CAR-003", address: "Av. da Marina 88, Lisbon, PT", coords: { lat: 38.7223, lng: -9.1393 }, capacity: 4000, currentUtilization: 2800, status: "active" },
  { id: "HUB-004", name: "Oslo Warehouse", carrierId: "CAR-004", address: "Havnegata 9, Oslo, NO", coords: { lat: 59.9139, lng: 10.7522 }, capacity: 2000, currentUtilization: 820, status: "active" },
  { id: "HUB-005", name: "Mumbai Warehouse", carrierId: "CAR-005", address: "Port Rd, Mumbai, IN", coords: { lat: 19.076, lng: 72.8777 }, capacity: 2500, currentUtilization: 210, status: "active" },
  { id: "HUB-006", name: "Barcelona Hub", carrierId: "CAR-006", address: "Moll d'Espanya 3, Barcelona, ES", coords: { lat: 41.3851, lng: 2.1734 }, capacity: 1800, currentUtilization: 990, status: "maintenance" },
];

export const driverStatusLabels: Record<DriverStatus, string> = {
  available: "Available",
  on_delivery: "On Delivery",
  off_duty: "Off Duty",
};

export const hubStatusLabels: Record<HubStatus, string> = {
  active: "Active",
  maintenance: "Maintenance",
  closed: "Closed",
};

export const carrierStatusLabels: Record<CarrierStatus, string> = {
  active: "Active",
  inactive: "Inactive",
};

export function getCarrier(id: string) {
  return carriers.find((c) => c.id === id);
}

export function getCarrierByName(name: string) {
  return carriers.find((c) => c.name === name);
}

export function getDriversByCarrier(carrierId: string) {
  return drivers.filter((d) => d.carrierId === carrierId);
}

export function getHubsByCarrier(carrierId: string) {
  return hubs.filter((h) => h.carrierId === carrierId);
}

export function getAvailableDrivers() {
  return drivers.filter((d) => d.status === "available");
}
```

- [ ] **Step 3: Commit**

```bash
git add frontend/src/lib/orders.ts frontend/src/lib/carriers.ts
git commit -m "feat: add carriers data model with mock data"
```

---

### Task 2: Simulated API layer

**Files:**
- Create: `frontend/src/lib/api/carriers.ts`
- Create: `frontend/src/lib/api/orders.ts`

- [ ] **Step 1: Create simulated carriers API**

Create `frontend/src/lib/api/carriers.ts`:

```typescript
import type { Driver, DriverStatus, Hub, HubStatus } from "@/lib/carriers";
import { carriers, drivers, hubs } from "@/lib/carriers";
import { orders, type Order, type ShipmentStatus } from "@/lib/orders";

function delay(ms = 200): Promise<void> {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

// --- Carriers ---

export async function fetchCarriers() {
  await delay();
  return [...carriers];
}

// --- Drivers ---

export async function fetchDrivers() {
  await delay();
  return [...drivers];
}

export async function assignDriverToOrder(driverId: string, orderId: string) {
  await delay(100);
  const driver = drivers.find((d) => d.id === driverId);
  if (!driver) throw new Error("Driver not found");
  const order = orders.find((o) => o.id === orderId);
  if (!order) throw new Error("Order not found");
  driver.status = "on_delivery" as DriverStatus;
  order.driverId = driverId;
  return { driver: { ...driver }, order: { ...order } };
}

// --- Hubs ---

export async function fetchHubs() {
  await delay();
  return [...hubs];
}

export async function createHub(data: Omit<Hub, "id">) {
  await delay(150);
  const id = `HUB-${String(hubs.length + 1).padStart(3, "0")}`;
  const hub: Hub = { id, ...data };
  hubs.push(hub);
  return { ...hub };
}

export async function updateHub(id: string, data: Partial<Hub>) {
  await delay(150);
  const idx = hubs.findIndex((h) => h.id === id);
  if (idx === -1) throw new Error("Hub not found");
  hubs[idx] = { ...hubs[idx], ...data };
  return { ...hubs[idx] };
}

export async function deleteHub(id: string) {
  await delay(100);
  const idx = hubs.findIndex((h) => h.id === id);
  if (idx === -1) throw new Error("Hub not found");
  hubs.splice(idx, 1);
}
```

- [ ] **Step 2: Create simulated orders API**

Create `frontend/src/lib/api/orders.ts`:

```typescript
import { orders, type Order, type ShipmentStatus } from "@/lib/orders";

function delay(ms = 200): Promise<void> {
  return new Promise((resolve) => setTimeout(resolve, ms));
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
```

- [ ] **Step 3: Commit**

```bash
git add frontend/src/lib/api/
git commit -m "feat: add simulated API layer for carriers and orders"
```

---

### Task 3: Vue Query hooks

**Files:**
- Create: `frontend/src/hooks/useCarriers.ts`
- Create: `frontend/src/hooks/useDrivers.ts`
- Create: `frontend/src/hooks/useHubs.ts`
- Create: `frontend/src/hooks/useDeliveries.ts`

- [ ] **Step 1: Create carriers hook**

Create `frontend/src/hooks/useCarriers.ts`:

```typescript
import { useQuery } from "@tanstack/vue-query";
import { fetchCarriers } from "@/lib/api/carriers";

export function useCarriers() {
  return useQuery({
    queryKey: ["carriers"],
    queryFn: fetchCarriers,
  });
}
```

- [ ] **Step 2: Create drivers hook**

Create `frontend/src/hooks/useDrivers.ts`:

```typescript
import { useQuery, useMutation, useQueryClient } from "@tanstack/vue-query";
import { fetchDrivers, assignDriverToOrder } from "@/lib/api/carriers";

export function useDrivers() {
  return useQuery({
    queryKey: ["drivers"],
    queryFn: fetchDrivers,
  });
}

export function useAssignDriver() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ driverId, orderId }: { driverId: string; orderId: string }) =>
      assignDriverToOrder(driverId, orderId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["drivers"] });
      queryClient.invalidateQueries({ queryKey: ["deliveries"] });
    },
  });
}
```

- [ ] **Step 3: Create hubs hook**

Create `frontend/src/hooks/useHubs.ts`:

```typescript
import { useQuery, useMutation, useQueryClient } from "@tanstack/vue-query";
import { fetchHubs, createHub, updateHub, deleteHub } from "@/lib/api/carriers";
import type { Hub } from "@/lib/carriers";

export function useHubs() {
  return useQuery({
    queryKey: ["hubs"],
    queryFn: fetchHubs,
  });
}

export function useCreateHub() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: Omit<Hub, "id">) => createHub(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["hubs"] });
    },
  });
}

export function useUpdateHub() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<Hub> }) => updateHub(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["hubs"] });
    },
  });
}

export function useDeleteHub() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => deleteHub(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["hubs"] });
    },
  });
}
```

- [ ] **Step 4: Create deliveries hook**

Create `frontend/src/hooks/useDeliveries.ts`:

```typescript
import { useQuery, useMutation, useQueryClient } from "@tanstack/vue-query";
import { fetchActiveDeliveries, updateShipmentStatus } from "@/lib/api/orders";
import type { ShipmentStatus } from "@/lib/orders";

export function useActiveDeliveries() {
  return useQuery({
    queryKey: ["deliveries"],
    queryFn: fetchActiveDeliveries,
    refetchInterval: 15_000,
  });
}

export function useUpdateShipmentStatus() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ orderId, status }: { orderId: string; status: ShipmentStatus }) =>
      updateShipmentStatus(orderId, status),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["deliveries"] });
    },
  });
}
```

- [ ] **Step 5: Commit**

```bash
git add frontend/src/hooks/
git commit -m "feat: add Vue Query hooks for carriers domain"
```

---

### Task 4: Route + SiteHeader updates

**Files:**
- Modify: `frontend/src/router/index.ts`
- Modify: `frontend/src/components/SiteHeader.vue`

- [ ] **Step 1: Add carriers route**

In `frontend/src/router/index.ts`, add after the `order-detail` route:

```typescript
    {
      path: '/carriers',
      name: 'carriers',
      component: () => import('@/views/CarriersView.vue'),
    },
```

- [ ] **Step 2: Add Carriers nav link to SiteHeader**

In `frontend/src/components/SiteHeader.vue`, add after the Orders nav link:

```typescript
        <RouterLink
          to="/carriers"
          class="rounded-md px-3 py-1.5 transition-colors"
          :class="route.path.startsWith('/carriers') ? 'bg-secondary text-foreground' : 'text-muted-foreground hover:text-foreground'"
        >
          Carriers
        </RouterLink>
```

- [ ] **Step 3: Commit**

```bash
git add frontend/src/router/index.ts frontend/src/components/SiteHeader.vue
git commit -m "feat: add /carriers route and nav link"
```

---

### Task 5: CarriersView shell

**File:**
- Create: `frontend/src/views/CarriersView.vue`

- [ ] **Step 1: Create CarriersView with tab navigation**

Create `frontend/src/views/CarriersView.vue`:

```typescript
<script setup lang="ts">
import { ref, computed, defineAsyncComponent } from 'vue'
import { Truck, Warehouse, BarChart3, Package } from 'lucide-vue-next'
import SiteHeader from '@/components/SiteHeader.vue'
import { cn } from '@/lib/utils'

const DriversPanel = defineAsyncComponent(() => import('@/components/DriversPanel.vue'))
const HubsPanel = defineAsyncComponent(() => import('@/components/HubsPanel.vue'))
const AnalyticsPanel = defineAsyncComponent(() => import('@/components/AnalyticsPanel.vue'))
const DeliveriesPanel = defineAsyncComponent(() => import('@/components/DeliveriesPanel.vue'))

type Tab = 'drivers' | 'hubs' | 'analytics' | 'deliveries'

const activeTab = ref<Tab>('drivers')

const tabs: Array<{ key: Tab; label: string; icon: typeof Truck }> = [
  { key: 'drivers', label: 'Drivers', icon: Truck },
  { key: 'hubs', label: 'Hubs', icon: Warehouse },
  { key: 'analytics', label: 'Analytics', icon: BarChart3 },
  { key: 'deliveries', label: 'Active Deliveries', icon: Package },
]
</script>

<template>
  <div class="min-h-screen">
    <SiteHeader />

    <section class="border-b border-border bg-gradient-hero">
      <div class="mx-auto max-w-7xl px-6 py-14">
        <span class="font-mono text-xs uppercase tracking-widest text-primary">/ carriers</span>
        <h1 class="mt-3 text-4xl font-semibold tracking-tight md:text-5xl">Carrier operations</h1>
        <p class="mt-3 max-w-2xl text-muted-foreground">
          Manage carrier drivers, hubs, and monitor active deliveries across the fleet.
        </p>
      </div>
    </section>

    <section class="mx-auto max-w-7xl px-6 pt-0">
      <!-- Tab bar -->
      <div class="flex gap-1 -mt-px">
        <button
          v-for="t in tabs"
          :key="t.key"
          @click="activeTab = t.key"
          :class="cn(
            'flex items-center gap-2 rounded-t-lg px-5 py-3 font-mono text-sm transition-colors border border-border',
            activeTab === t.key
              ? 'bg-card text-foreground border-b-card -mb-px'
              : 'bg-transparent text-muted-foreground hover:text-foreground border-transparent hover:border-border',
          )"
        >
          <component :is="t.icon" class="h-4 w-4" />
          {{ t.label }}
        </button>
      </div>

      <!-- Tab content -->
      <div class="rounded-b-xl rounded-tr-xl border border-border bg-card p-6">
        <DriversPanel v-if="activeTab === 'drivers'" />
        <HubsPanel v-else-if="activeTab === 'hubs'" />
        <AnalyticsPanel v-else-if="activeTab === 'analytics'" />
        <DeliveriesPanel v-else-if="activeTab === 'deliveries'" />
      </div>
    </section>
  </div>
</template>
```

- [ ] **Step 2: Commit**

```bash
git add frontend/src/views/CarriersView.vue
git commit -m "feat: add CarriersView with tab navigation shell"
```

---

### Task 6: DriversPanel + AssignDriverModal

**Files:**
- Create: `frontend/src/components/DriversPanel.vue`
- Create: `frontend/src/components/AssignDriverModal.vue`

- [ ] **Step 1: Create DriversPanel**

Create `frontend/src/components/DriversPanel.vue`:

```typescript
<script setup lang="ts">
import { ref, computed } from 'vue'
import { Search, UserPlus } from 'lucide-vue-next'
import { useDrivers } from '@/hooks/useDrivers'
import { getCarrier, type DriverStatus, driverStatusLabels, drivers as allDrivers } from '@/lib/carriers'
import { cn } from '@/lib/utils'
import Input from '@/components/ui/Input.vue'
import Skeleton from '@/components/ui/Skeleton.vue'
import Button from '@/components/ui/Button.vue'
import AssignDriverModal from '@/components/AssignDriverModal.vue'

const { data: drivers, isLoading, isError, refetch } = useDrivers()

const query = ref('')
const statusFilter = ref<DriverStatus | 'all'>('all')

const filtered = computed(() => {
  if (!drivers.value) return []
  const q = query.value.trim().toLowerCase()
  return drivers.value.filter((d) => {
    if (statusFilter.value !== 'all' && d.status !== statusFilter.value) return false
    if (!q) return true
    const carrier = getCarrier(d.carrierId)
    return (
      d.name.toLowerCase().includes(q) ||
      carrier?.name.toLowerCase().includes(q) ||
      d.vehicleInfo.toLowerCase().includes(q)
    )
  })
})

const statusCounts = computed(() => {
  if (!drivers.value) return { total: 0, available: 0, on_delivery: 0, off_duty: 0 }
  return {
    total: drivers.value.length,
    available: drivers.value.filter((d) => d.status === 'available').length,
    on_delivery: drivers.value.filter((d) => d.status === 'on_delivery').length,
    off_duty: drivers.value.filter((d) => d.status === 'off_duty').length,
  }
})

const showAssignModal = ref(false)
</script>

<template>
  <div v-if="isLoading" class="space-y-4">
    <div class="grid grid-cols-4 gap-4">
      <Skeleton v-for="i in 4" :key="i" class="h-24 rounded-xl" />
    </div>
    <Skeleton class="h-64 rounded-xl" />
  </div>

  <div v-else-if="isError" class="py-12 text-center">
    <p class="font-mono text-sm text-destructive">Failed to load drivers.</p>
    <Button variant="outline" class="mt-4" @click="refetch()">Retry</Button>
  </div>

  <div v-else>
    <!-- Stats -->
    <div class="grid grid-cols-2 gap-4 md:grid-cols-4">
      <div class="rounded-lg border border-border bg-secondary/50 p-4">
        <div class="font-mono text-[11px] uppercase tracking-widest text-muted-foreground">Total Drivers</div>
        <div class="mt-1 font-mono text-3xl font-semibold">{{ statusCounts.total }}</div>
      </div>
      <div class="rounded-lg border border-border bg-secondary/50 p-4">
        <div class="font-mono text-[11px] uppercase tracking-widest text-muted-foreground">Available</div>
        <div class="mt-1 font-mono text-3xl font-semibold text-success">{{ statusCounts.available }}</div>
      </div>
      <div class="rounded-lg border border-border bg-secondary/50 p-4">
        <div class="font-mono text-[11px] uppercase tracking-widest text-muted-foreground">On Delivery</div>
        <div class="mt-1 font-mono text-3xl font-semibold text-info">{{ statusCounts.on_delivery }}</div>
      </div>
      <div class="rounded-lg border border-border bg-secondary/50 p-4">
        <div class="font-mono text-[11px] uppercase tracking-widest text-muted-foreground">Off Duty</div>
        <div class="mt-1 font-mono text-3xl font-semibold text-muted-foreground">{{ statusCounts.off_duty }}</div>
      </div>
    </div>

    <!-- Controls -->
    <div class="mt-6 flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
      <div class="flex items-center gap-2 rounded-lg border border-border bg-card px-3 md:w-72">
        <Search class="h-4 w-4 text-muted-foreground" />
        <Input
          v-model="query"
          placeholder="Search drivers..."
          class="h-10 border-0 bg-transparent font-mono text-sm shadow-none focus-visible:ring-0"
        />
      </div>
      <div class="flex items-center gap-2">
        <button
          v-for="s in (['all', 'available', 'on_delivery', 'off_duty'] as const)"
          :key="s"
          @click="statusFilter = s"
          :class="cn(
            'rounded-full border px-3 py-1.5 font-mono text-xs uppercase tracking-wider transition-colors',
            statusFilter === s
              ? 'border-primary bg-primary/15 text-primary'
              : 'border-border text-muted-foreground hover:text-foreground',
          )"
        >
          {{ s === 'all' ? 'All' : driverStatusLabels[s] }}
        </button>
        <Button size="sm" class="gap-2 ml-2" @click="showAssignModal = true">
          <UserPlus class="h-4 w-4" /> Assign
        </Button>
      </div>
    </div>

    <!-- Table -->
    <div class="mt-4 overflow-hidden rounded-xl border border-border">
      <div class="hidden grid-cols-[1.5fr_1.5fr_1fr_1fr_0.8fr] gap-4 border-b border-border bg-secondary/50 px-6 py-3 font-mono text-[11px] uppercase tracking-widest text-muted-foreground md:grid">
        <span>Driver</span>
        <span>Carrier</span>
        <span>Status</span>
        <span>Vehicle</span>
        <span class="text-right">Actions</span>
      </div>

      <div v-if="filtered.length === 0" class="px-6 py-12 text-center font-mono text-sm text-muted-foreground">
        No drivers match your filters.
      </div>

      <div v-for="d in filtered" :key="d.id" class="group grid grid-cols-1 gap-2 border-b border-border px-6 py-4 transition-colors last:border-0 hover:bg-secondary/40 md:grid-cols-[1.5fr_1.5fr_1fr_1fr_0.8fr] md:items-center">
        <div>
          <div class="font-mono text-sm">{{ d.name }}</div>
          <div class="font-mono text-xs text-muted-foreground">{{ d.email }}</div>
        </div>
        <div class="font-mono text-sm text-muted-foreground">{{ getCarrier(d.carrierId)?.name ?? d.carrierId }}</div>
        <div>
          <span
            :class="cn(
              'inline-flex items-center gap-1.5 rounded-full border px-2.5 py-0.5 text-xs font-mono uppercase tracking-wider',
              d.status === 'available' ? 'bg-success/15 text-success border-success/30' : '',
              d.status === 'on_delivery' ? 'bg-info/15 text-info border-info/30' : '',
              d.status === 'off_duty' ? 'bg-muted text-muted-foreground border-border' : '',
            )"
          >
            <span class="h-1.5 w-1.5 rounded-full bg-current" />
            {{ driverStatusLabels[d.status] }}
          </span>
        </div>
        <div class="font-mono text-sm text-muted-foreground">{{ d.vehicleInfo }}</div>
        <div class="text-right">
          <Button
            v-if="d.status === 'available'"
            variant="outline"
            size="sm"
            @click="showAssignModal = true"
          >
            Assign
          </Button>
        </div>
      </div>
    </div>

    <div class="mt-4 font-mono text-xs text-muted-foreground">
      Showing {{ filtered.length }} of {{ statusCounts.total }} drivers
    </div>

    <AssignDriverModal v-if="showAssignModal" @close="showAssignModal = false" />
  </div>
</template>
```

- [ ] **Step 2: Create AssignDriverModal**

Create `frontend/src/components/AssignDriverModal.vue`:

```typescript
<script setup lang="ts">
import { ref, computed } from 'vue'
import { Search, X } from 'lucide-vue-next'
import { useDrivers, useAssignDriver } from '@/hooks/useDrivers'
import { orders } from '@/lib/orders'
import { getCarrier } from '@/lib/carriers'
import Input from '@/components/ui/Input.vue'
import Button from '@/components/ui/Button.vue'

const emit = defineEmits<{ close: [] }>()

const { data: drivers } = useDrivers()
const assignMutation = useAssignDriver()

const driverQuery = ref('')
const selectedDriverId = ref('')
const selectedOrderId = ref('')

const availableDrivers = computed(() => {
  if (!drivers.value) return []
  const q = driverQuery.value.trim().toLowerCase()
  return drivers.value.filter((d) => {
    if (d.status !== 'available') return false
    if (!q) return true
    return d.name.toLowerCase().includes(q) || d.vehicleInfo.toLowerCase().includes(q)
  })
})

const unfinishedOrders = computed(() =>
  orders.filter((o) => o.status !== 'delivered' && !o.driverId)
)

const canSubmit = computed(() => selectedDriverId.value && selectedOrderId.value)

async function handleAssign() {
  if (!canSubmit.value) return
  await assignMutation.mutateAsync({
    driverId: selectedDriverId.value,
    orderId: selectedOrderId.value,
  })
  emit('close')
}
</script>

<template>
  <div class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm" @click.self="emit('close')">
    <div class="w-full max-w-lg rounded-xl border border-border bg-card p-6 shadow-elegant">
      <div class="flex items-center justify-between">
        <h2 class="font-mono text-lg font-semibold">Assign Driver to Shipment</h2>
        <button @click="emit('close')" class="text-muted-foreground hover:text-foreground">
          <X class="h-5 w-5" />
        </button>
      </div>

      <!-- Driver selection -->
      <div class="mt-5">
        <label class="font-mono text-xs uppercase tracking-widest text-muted-foreground">Select Driver</label>
        <div class="mt-2 flex items-center gap-2 rounded-lg border border-border bg-background px-3">
          <Search class="h-4 w-4 text-muted-foreground" />
          <Input
            v-model="driverQuery"
            placeholder="Search available drivers..."
            class="h-10 border-0 bg-transparent font-mono text-sm shadow-none focus-visible:ring-0"
          />
        </div>
        <div class="mt-2 max-h-40 overflow-y-auto space-y-1">
          <button
            v-for="d in availableDrivers"
            :key="d.id"
            @click="selectedDriverId = d.id"
            class="w-full rounded-lg px-3 py-2 text-left font-mono text-sm transition-colors"
            :class="selectedDriverId === d.id ? 'bg-primary/15 text-primary' : 'hover:bg-secondary'"
          >
            <div>{{ d.name }}</div>
            <div class="text-xs text-muted-foreground">{{ d.vehicleInfo }} · {{ getCarrier(d.carrierId)?.name }}</div>
          </button>
          <div v-if="availableDrivers.length === 0" class="py-4 text-center font-mono text-xs text-muted-foreground">
            No available drivers found.
          </div>
        </div>
      </div>

      <!-- Order selection -->
      <div class="mt-5">
        <label class="font-mono text-xs uppercase tracking-widest text-muted-foreground">Target Shipment</label>
        <div class="mt-2 space-y-1 max-h-40 overflow-y-auto">
          <button
            v-for="o in unfinishedOrders"
            :key="o.id"
            @click="selectedOrderId = o.id"
            class="w-full rounded-lg px-3 py-2 text-left font-mono text-sm transition-colors"
            :class="selectedOrderId === o.id ? 'bg-primary/15 text-primary' : 'hover:bg-secondary'"
          >
            <div>{{ o.id }} — {{ o.customer }}</div>
            <div class="text-xs text-muted-foreground">{{ o.origin }} → {{ o.destination }}</div>
          </button>
        </div>
      </div>

      <div v-if="assignMutation.isError" class="mt-3 font-mono text-xs text-destructive">
        Failed to assign driver. Please try again.
      </div>

      <div class="mt-6 flex justify-end gap-3">
        <Button variant="outline" @click="emit('close')">Cancel</Button>
        <Button :disabled="!canSubmit || assignMutation.isPending" @click="handleAssign">
          {{ assignMutation.isPending ? 'Assigning…' : 'Assign Driver' }}
        </Button>
      </div>
    </div>
  </div>
</template>
```

- [ ] **Step 3: Commit**

```bash
git add frontend/src/components/DriversPanel.vue frontend/src/components/AssignDriverModal.vue
git commit -m "feat: add DriversPanel and AssignDriverModal"
```

---

### Task 7: HubsPanel + HubFormModal

**Files:**
- Create: `frontend/src/components/HubsPanel.vue`
- Create: `frontend/src/components/HubFormModal.vue`

- [ ] **Step 1: Create HubsPanel**

Create `frontend/src/components/HubsPanel.vue`:

```typescript
<script setup lang="ts">
import { ref, computed } from 'vue'
import { Search, Plus, Pencil, Trash2 } from 'lucide-vue-next'
import { useHubs, useDeleteHub } from '@/hooks/useHubs'
import { getCarrier, hubStatusLabels } from '@/lib/carriers'
import { cn } from '@/lib/utils'
import Input from '@/components/ui/Input.vue'
import Skeleton from '@/components/ui/Skeleton.vue'
import Button from '@/components/ui/Button.vue'
import HubFormModal from '@/components/HubFormModal.vue'

const { data: hubs, isLoading, isError, refetch } = useHubs()
const deleteHub = useDeleteHub()

const query = ref('')

const filtered = computed(() => {
  if (!hubs.value) return []
  const q = query.value.trim().toLowerCase()
  return hubs.value.filter((h) => {
    if (!q) return true
    return (
      h.name.toLowerCase().includes(q) ||
      h.address.toLowerCase().includes(q) ||
      getCarrier(h.carrierId)?.name.toLowerCase().includes(q)
    )
  })
})

const showForm = ref(false)
const editingHubId = ref<string | null>(null)

function openAdd() {
  editingHubId.value = null
  showForm.value = true
}

function openEdit(id: string) {
  editingHubId.value = id
  showForm.value = true
}

const hubStatusCounts = computed(() => {
  if (!hubs.value) return { total: 0, active: 0, maintenance: 0, closed: 0 }
  return {
    total: hubs.value.length,
    active: hubs.value.filter((h) => h.status === 'active').length,
    maintenance: hubs.value.filter((h) => h.status === 'maintenance').length,
    closed: hubs.value.filter((h) => h.status === 'closed').length,
  }
})
</script>

<template>
  <div v-if="isLoading" class="space-y-4">
    <div class="grid grid-cols-4 gap-4">
      <Skeleton v-for="i in 4" :key="i" class="h-24 rounded-xl" />
    </div>
    <Skeleton class="h-64 rounded-xl" />
  </div>

  <div v-else-if="isError" class="py-12 text-center">
    <p class="font-mono text-sm text-destructive">Failed to load hubs.</p>
    <Button variant="outline" class="mt-4" @click="refetch()">Retry</Button>
  </div>

  <div v-else>
    <div class="grid grid-cols-2 gap-4 md:grid-cols-4">
      <div class="rounded-lg border border-border bg-secondary/50 p-4">
        <div class="font-mono text-[11px] uppercase tracking-widest text-muted-foreground">Total Hubs</div>
        <div class="mt-1 font-mono text-3xl font-semibold">{{ hubStatusCounts.total }}</div>
      </div>
      <div class="rounded-lg border border-border bg-secondary/50 p-4">
        <div class="font-mono text-[11px] uppercase tracking-widest text-muted-foreground">Active</div>
        <div class="mt-1 font-mono text-3xl font-semibold text-success">{{ hubStatusCounts.active }}</div>
      </div>
      <div class="rounded-lg border border-border bg-secondary/50 p-4">
        <div class="font-mono text-[11px] uppercase tracking-widest text-muted-foreground">Maintenance</div>
        <div class="mt-1 font-mono text-3xl font-semibold text-warning">{{ hubStatusCounts.maintenance }}</div>
      </div>
      <div class="rounded-lg border border-border bg-secondary/50 p-4">
        <div class="font-mono text-[11px] uppercase tracking-widest text-muted-foreground">Closed</div>
        <div class="mt-1 font-mono text-3xl font-semibold text-destructive">{{ hubStatusCounts.closed }}</div>
      </div>
    </div>

    <div class="mt-6 flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
      <div class="flex items-center gap-2 rounded-lg border border-border bg-card px-3 md:w-72">
        <Search class="h-4 w-4 text-muted-foreground" />
        <Input
          v-model="query"
          placeholder="Search hubs..."
          class="h-10 border-0 bg-transparent font-mono text-sm shadow-none focus-visible:ring-0"
        />
      </div>
      <Button size="sm" class="gap-2" @click="openAdd">
        <Plus class="h-4 w-4" /> Add Hub
      </Button>
    </div>

    <div class="mt-4 overflow-hidden rounded-xl border border-border">
      <div class="hidden grid-cols-[1.3fr_1.3fr_1.8fr_1fr_1fr_0.6fr] gap-4 border-b border-border bg-secondary/50 px-6 py-3 font-mono text-[11px] uppercase tracking-widest text-muted-foreground md:grid">
        <span>Name</span>
        <span>Carrier</span>
        <span>Address</span>
        <span>Capacity</span>
        <span>Status</span>
        <span class="text-right">Actions</span>
      </div>

      <div v-if="filtered.length === 0" class="px-6 py-12 text-center font-mono text-sm text-muted-foreground">
        No hubs match your filters.
      </div>

      <div v-for="h in filtered" :key="h.id" class="group grid grid-cols-1 gap-2 border-b border-border px-6 py-4 transition-colors last:border-0 hover:bg-secondary/40 md:grid-cols-[1.3fr_1.3fr_1.8fr_1fr_1fr_0.6fr] md:items-center">
        <div class="font-mono text-sm">{{ h.name }}</div>
        <div class="font-mono text-sm text-muted-foreground">{{ getCarrier(h.carrierId)?.name ?? h.carrierId }}</div>
        <div class="font-mono text-xs text-muted-foreground">{{ h.address }}</div>
        <div class="flex items-center gap-2">
          <div class="h-2 flex-1 overflow-hidden rounded-full bg-secondary">
            <div
              class="h-full rounded-full transition-all"
              :class="(h.currentUtilization / h.capacity) > 0.8 ? 'bg-warning' : 'bg-primary'"
              :style="{ width: `${Math.min(100, (h.currentUtilization / h.capacity) * 100)}%` }"
            />
          </div>
          <span class="font-mono text-xs text-muted-foreground">{{ Math.round((h.currentUtilization / h.capacity) * 100) }}%</span>
        </div>
        <div>
          <span
            :class="cn(
              'inline-flex items-center gap-1.5 rounded-full border px-2.5 py-0.5 text-xs font-mono uppercase tracking-wider',
              h.status === 'active' ? 'bg-success/15 text-success border-success/30' : '',
              h.status === 'maintenance' ? 'bg-warning/15 text-warning border-warning/30' : '',
              h.status === 'closed' ? 'bg-destructive/15 text-destructive border-destructive/30' : '',
            )"
          >
            <span class="h-1.5 w-1.5 rounded-full bg-current" />
            {{ hubStatusLabels[h.status] }}
          </span>
        </div>
        <div class="flex justify-end gap-1">
          <button @click="openEdit(h.id)" class="rounded p-1.5 text-muted-foreground hover:text-foreground">
            <Pencil class="h-4 w-4" />
          </button>
          <button
            @click="deleteHub.mutate(h.id)"
            class="rounded p-1.5 text-muted-foreground hover:text-destructive"
          >
            <Trash2 class="h-4 w-4" />
          </button>
        </div>
      </div>
    </div>

    <HubFormModal
      v-if="showForm"
      :hub-id="editingHubId"
      @close="showForm = false"
    />
  </div>
</template>
```

- [ ] **Step 2: Create HubFormModal**

Create `frontend/src/components/HubFormModal.vue`:

```typescript
<script setup lang="ts">
import { ref, computed } from 'vue'
import { X } from 'lucide-vue-next'
import { useHubs, useCreateHub, useUpdateHub } from '@/hooks/useHubs'
import { carriers, hubs, hubStatusLabels, type HubStatus } from '@/lib/carriers'
import Input from '@/components/ui/Input.vue'
import Button from '@/components/ui/Button.vue'

const props = defineProps<{ hubId?: string | null }>()
const emit = defineEmits<{ close: [] }>()

const { data: hubsData } = useHubs()
const createHub = useCreateHub()
const updateHub = useUpdateHub()

const existing = computed(() => {
  if (!props.hubId || !hubsData.value) return null
  return hubsData.value.find((h) => h.id === props.hubId) ?? null
})

const name = ref(existing.value?.name ?? '')
const carrierId = ref(existing.value?.carrierId ?? carriers[0].id)
const address = ref(existing.value?.address ?? '')
const capacity = ref(existing.value?.capacity ?? 1000)
const status = ref<HubStatus>(existing.value?.status ?? 'active')

const isEditing = computed(() => !!props.hubId)

async function handleSubmit() {
  const data = {
    name: name.value,
    carrierId: carrierId.value,
    address: address.value,
    coords: { lat: 0, lng: 0 },
    capacity: capacity.value,
    currentUtilization: existing.value?.currentUtilization ?? 0,
    status: status.value,
  }

  if (isEditing.value && props.hubId) {
    await updateHub.mutateAsync({ id: props.hubId, data })
  } else {
    await createHub.mutateAsync(data)
  }
  emit('close')
}
</script>

<template>
  <div class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm" @click.self="emit('close')">
    <div class="w-full max-w-md rounded-xl border border-border bg-card p-6 shadow-elegant">
      <div class="flex items-center justify-between">
        <h2 class="font-mono text-lg font-semibold">{{ isEditing ? 'Edit Hub' : 'Add Hub' }}</h2>
        <button @click="emit('close')" class="text-muted-foreground hover:text-foreground">
          <X class="h-5 w-5" />
        </button>
      </div>

      <div class="mt-5 space-y-4">
        <div>
          <label class="font-mono text-xs uppercase tracking-widest text-muted-foreground">Name</label>
          <Input v-model="name" class="mt-1.5 font-mono text-sm" placeholder="Hub name" />
        </div>
        <div>
          <label class="font-mono text-xs uppercase tracking-widest text-muted-foreground">Carrier</label>
          <select
            v-model="carrierId"
            class="mt-1.5 flex h-10 w-full rounded-lg border border-border bg-background px-3 font-mono text-sm"
          >
            <option v-for="c in carriers" :key="c.id" :value="c.id">{{ c.name }}</option>
          </select>
        </div>
        <div>
          <label class="font-mono text-xs uppercase tracking-widest text-muted-foreground">Address</label>
          <Input v-model="address" class="mt-1.5 font-mono text-sm" placeholder="Full address" />
        </div>
        <div class="grid grid-cols-2 gap-4">
          <div>
            <label class="font-mono text-xs uppercase tracking-widest text-muted-foreground">Capacity (units)</label>
            <Input v-model.number="capacity" type="number" class="mt-1.5 font-mono text-sm" />
          </div>
          <div>
            <label class="font-mono text-xs uppercase tracking-widest text-muted-foreground">Status</label>
            <select
              v-model="status"
              class="mt-1.5 flex h-10 w-full rounded-lg border border-border bg-background px-3 font-mono text-sm"
            >
              <option v-for="(label, key) in hubStatusLabels" :key="key" :value="key">{{ label }}</option>
            </select>
          </div>
        </div>
      </div>

      <div v-if="createHub.isError || updateHub.isError" class="mt-3 font-mono text-xs text-destructive">
        Failed to save hub. Please try again.
      </div>

      <div class="mt-6 flex justify-end gap-3">
        <Button variant="outline" @click="emit('close')">Cancel</Button>
        <Button
          :disabled="!name || createHub.isPending || updateHub.isPending"
          @click="handleSubmit"
        >
          {{ createHub.isPending || updateHub.isPending ? 'Saving…' : isEditing ? 'Update Hub' : 'Create Hub' }}
        </Button>
      </div>
    </div>
  </div>
</template>
```

- [ ] **Step 3: Commit**

```bash
git add frontend/src/components/HubsPanel.vue frontend/src/components/HubFormModal.vue
git commit -m "feat: add HubsPanel and HubFormModal with CRUD"
```

---

### Task 8: AnalyticsPanel

**File:**
- Create: `frontend/src/components/AnalyticsPanel.vue`

- [ ] **Step 1: Create AnalyticsPanel with CSS bar charts**

Create `frontend/src/components/AnalyticsPanel.vue`:

```typescript
<script setup lang="ts">
import { computed } from 'vue'
import { useCarriers } from '@/hooks/useCarriers'
import { orders, statusLabels } from '@/lib/orders'
import Skeleton from '@/components/ui/Skeleton.vue'

const { data: carriersData, isLoading, isError, refetch } = useCarriers()

const kpis = computed(() => {
  const total = orders.length
  const delivered = orders.filter((o) => o.status === 'delivered').length
  const onTime = Math.round((delivered / Math.max(total, 1)) * 100)
  const activeCarriers = carriersData.value?.filter((c) => c.status === 'active').length ?? 0
  return { total, onTime, activeCarriers, avgDeliveryTime: '3.2 days' }
})

const carrierPerformance = computed(() => {
  if (!carriersData.value) return []
  return carriersData.value.map((c) => {
    const carrierOrders = orders.filter((o) => o.carrier === c.name)
    const delivered = carrierOrders.filter((o) => o.status === 'delivered').length
    return {
      name: c.name,
      total: carrierOrders.length,
      delivered,
      onTimeRate: carrierOrders.length > 0 ? Math.round((delivered / carrierOrders.length) * 100) : 0,
    }
  })
})

const statusDistribution = computed(() => {
  const counts: Record<string, number> = {}
  for (const o of orders) {
    counts[o.status] = (counts[o.status] || 0) + 1
  }
  return Object.entries(counts).map(([status, count]) => ({
    status,
    label: statusLabels[status as keyof typeof statusLabels] ?? status,
    count,
    pct: Math.round((count / orders.length) * 100),
  }))
})

const maxCarrierOrders = computed(() =>
  Math.max(...carrierPerformance.value.map((c) => c.total), 1)
)

const maxStatusCount = computed(() =>
  Math.max(...statusDistribution.value.map((s) => s.count), 1)
)
</script>

<template>
  <div v-if="isLoading" class="space-y-4">
    <div class="grid grid-cols-4 gap-4">
      <Skeleton v-for="i in 4" :key="i" class="h-24 rounded-xl" />
    </div>
    <Skeleton class="h-48 rounded-xl" />
    <Skeleton class="h-48 rounded-xl" />
  </div>

  <div v-else-if="isError" class="py-12 text-center">
    <p class="font-mono text-sm text-destructive">Failed to load analytics data.</p>
    <button
      @click="refetch()"
      class="mt-4 font-mono text-xs uppercase tracking-widest text-primary hover:underline"
    >
      Retry
    </button>
  </div>

  <div v-else>
    <!-- KPI cards -->
    <div class="grid grid-cols-2 gap-4 md:grid-cols-4">
      <div class="rounded-lg border border-border bg-secondary/50 p-4">
        <div class="font-mono text-[11px] uppercase tracking-widest text-muted-foreground">Total Shipments</div>
        <div class="mt-1 font-mono text-3xl font-semibold">{{ kpis.total }}</div>
      </div>
      <div class="rounded-lg border border-border bg-secondary/50 p-4">
        <div class="font-mono text-[11px] uppercase tracking-widest text-muted-foreground">On-Time Rate</div>
        <div class="mt-1 font-mono text-3xl font-semibold text-success">{{ kpis.onTime }}%</div>
      </div>
      <div class="rounded-lg border border-border bg-secondary/50 p-4">
        <div class="font-mono text-[11px] uppercase tracking-widest text-muted-foreground">Active Carriers</div>
        <div class="mt-1 font-mono text-3xl font-semibold text-info">{{ kpis.activeCarriers }}</div>
      </div>
      <div class="rounded-lg border border-border bg-secondary/50 p-4">
        <div class="font-mono text-[11px] uppercase tracking-widest text-muted-foreground">Avg Delivery</div>
        <div class="mt-1 font-mono text-3xl font-semibold">{{ kpis.avgDeliveryTime }}</div>
      </div>
    </div>

    <!-- Carrier Performance -->
    <div class="mt-8">
      <h3 class="font-mono text-sm font-semibold">Carrier Performance</h3>
      <div class="mt-4 space-y-3">
        <div v-for="c in carrierPerformance" :key="c.name" class="rounded-lg border border-border bg-card p-4">
          <div class="flex items-center justify-between">
            <span class="font-mono text-sm">{{ c.name }}</span>
            <span class="font-mono text-xs text-muted-foreground">{{ c.delivered }}/{{ c.total }} delivered</span>
          </div>
          <div class="mt-2 flex items-center gap-3">
            <div class="h-2 flex-1 overflow-hidden rounded-full bg-secondary">
              <div
                class="h-full rounded-full bg-gradient-accent transition-all"
                :style="{ width: `${(c.total / maxCarrierOrders) * 100}%` }"
              />
            </div>
            <span class="font-mono text-xs" :class="c.onTimeRate >= 80 ? 'text-success' : 'text-warning'">
              {{ c.onTimeRate }}%
            </span>
          </div>
        </div>
      </div>
    </div>

    <!-- Status Distribution -->
    <div class="mt-8">
      <h3 class="font-mono text-sm font-semibold">Shipment Status Distribution</h3>
      <div class="mt-4 space-y-2">
        <div v-for="s in statusDistribution" :key="s.status" class="flex items-center gap-3">
          <span class="w-32 font-mono text-xs text-muted-foreground">{{ s.label }}</span>
          <div class="h-5 flex-1 overflow-hidden rounded-full bg-secondary">
            <div
              class="h-full rounded-full transition-all"
              :class="
                s.status === 'delivered' ? 'bg-success' :
                s.status === 'delayed' ? 'bg-destructive' :
                s.status === 'in_transit' ? 'bg-info' :
                s.status === 'out_for_delivery' ? 'bg-primary' :
                'bg-muted-foreground/40'
              "
              :style="{ width: `${(s.count / maxStatusCount) * 100}%` }"
            />
          </div>
          <span class="w-16 text-right font-mono text-xs text-muted-foreground">{{ s.count }} ({{ s.pct }}%)</span>
        </div>
      </div>
    </div>
  </div>
</template>
```

- [ ] **Step 2: Commit**

```bash
git add frontend/src/components/AnalyticsPanel.vue
git commit -m "feat: add AnalyticsPanel with KPI cards and CSS charts"
```

---

### Task 9: DeliveriesPanel

**File:**
- Create: `frontend/src/components/DeliveriesPanel.vue`

- [ ] **Step 1: Create DeliveriesPanel**

Create `frontend/src/components/DeliveriesPanel.vue`:

```typescript
<script setup lang="ts">
import { ref, computed, defineAsyncComponent, onMounted } from 'vue'
import { Search, RefreshCw, MapPin, ArrowRight, UserPlus } from 'lucide-vue-next'
import { useActiveDeliveries, useUpdateShipmentStatus } from '@/hooks/useDeliveries'
import { drivers as driverData, getCarrierByName } from '@/lib/carriers'
import { statusLabels, type ShipmentStatus } from '@/lib/orders'
import { cn } from '@/lib/utils'
import Input from '@/components/ui/Input.vue'
import Skeleton from '@/components/ui/Skeleton.vue'
import StatusBadge from '@/components/StatusBadge.vue'
import Button from '@/components/ui/Button.vue'

const ShipmentMap = defineAsyncComponent(() => import('@/components/ShipmentMap.vue'))

const { data: deliveries, isLoading, isError, refetch, dataUpdatedAt } = useActiveDeliveries()
const updateStatus = useUpdateShipmentStatus()

const query = ref('')
const mounted = ref(false)

onMounted(() => {
  mounted.value = true
})

const filtered = computed(() => {
  if (!deliveries.value) return []
  const q = query.value.trim().toLowerCase()
  return deliveries.value.filter((o) => {
    if (!q) return true
    return (
      o.id.toLowerCase().includes(q) ||
      o.trackingNumber.toLowerCase().includes(q) ||
      o.customer.toLowerCase().includes(q) ||
      o.carrier.toLowerCase().includes(q)
    )
  })
})

function getDriverForOrder(orderId: string) {
  const order = deliveries.value?.find((o) => o.id === orderId)
  if (!order?.driverId) return null
  return driverData.find((d) => d.id === order.driverId) ?? null
}

const formatter = new Intl.RelativeTimeFormat('en', { numeric: 'auto' })
const secondsSinceUpdate = ref(0)
let interval: ReturnType<typeof setInterval> | undefined

function updateTime() {
  secondsSinceUpdate.value = Math.round((Date.now() - dataUpdatedAt.value) / 1000)
}

onMounted(() => {
  interval = setInterval(updateTime, 1000)
})
</script>

<template>
  <div v-if="isLoading" class="space-y-4">
    <Skeleton class="h-48 rounded-xl" />
    <Skeleton class="h-64 rounded-xl" />
  </div>

  <div v-else-if="isError" class="py-12 text-center">
    <p class="font-mono text-sm text-destructive">Failed to load active deliveries.</p>
    <Button variant="outline" class="mt-4" @click="refetch()">Retry</Button>
  </div>

  <div v-else>
    <!-- Controls bar -->
    <div class="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
      <div class="flex items-center gap-2 rounded-lg border border-border bg-card px-3 md:w-72">
        <Search class="h-4 w-4 text-muted-foreground" />
        <Input
          v-model="query"
          placeholder="Search active shipments..."
          class="h-10 border-0 bg-transparent font-mono text-sm shadow-none focus-visible:ring-0"
        />
      </div>
      <div class="flex items-center gap-3">
        <span class="font-mono text-xs text-muted-foreground">
          {{ deliveries?.length ?? 0 }} active · updated {{ secondsSinceUpdate }}s ago
        </span>
        <button @click="refetch()" class="rounded p-1.5 text-muted-foreground hover:text-foreground">
          <RefreshCw class="h-4 w-4" />
        </button>
      </div>
    </div>

    <!-- Table -->
    <div class="mt-4 overflow-hidden rounded-xl border border-border">
      <div class="hidden grid-cols-[0.8fr_1fr_1fr_1fr_1fr_1fr_0.6fr] gap-4 border-b border-border bg-secondary/50 px-6 py-3 font-mono text-[11px] uppercase tracking-widest text-muted-foreground md:grid">
        <span>Order ID</span>
        <span>Customer</span>
        <span>Carrier</span>
        <span>Driver</span>
        <span>Status</span>
        <span>ETA</span>
        <span class="text-right">Actions</span>
      </div>

      <div v-if="filtered.length === 0" class="px-6 py-12 text-center font-mono text-sm text-muted-foreground">
        No active deliveries match your filters.
      </div>

      <div v-for="o in filtered" :key="o.id" class="group grid grid-cols-1 gap-2 border-b border-border px-6 py-4 transition-colors last:border-0 hover:bg-secondary/40 md:grid-cols-[0.8fr_1fr_1fr_1fr_1fr_1fr_0.6fr] md:items-center">
        <div class="font-mono text-sm text-primary">{{ o.id }}</div>
        <div class="font-mono text-sm">{{ o.customer }}</div>
        <div class="font-mono text-sm text-muted-foreground">{{ o.carrier }}</div>
        <div>
          <div v-if="getDriverForOrder(o.id)" class="font-mono text-sm">
            {{ getDriverForOrder(o.id)?.name }}
          </div>
          <div v-else class="font-mono text-xs text-muted-foreground">Unassigned</div>
        </div>
        <div>
          <select
            :value="o.status"
            @change="(e) => updateStatus.mutate({ orderId: o.id, status: (e.target as HTMLSelectElement).value as ShipmentStatus })"
            class="rounded-lg border border-border bg-background px-2 py-1 font-mono text-xs"
          >
            <option v-for="(label, key) in statusLabels" :key="key" :value="key">{{ label }}</option>
          </select>
        </div>
        <div class="font-mono text-xs text-muted-foreground">{{ o.estimatedDelivery }}</div>
        <div class="flex justify-end gap-1">
          <RouterLink
            :to="{ name: 'order-detail', params: { orderId: o.id } }"
            class="rounded p-1.5 text-muted-foreground hover:text-foreground"
          >
            <ArrowRight class="h-4 w-4" />
          </RouterLink>
        </div>
      </div>
    </div>

    <!-- Mini map -->
    <div v-if="mounted && filtered.length > 0" class="mt-6 overflow-hidden rounded-xl border border-border">
      <Suspense>
        <div class="h-[300px] w-full">
          <ShipmentMap
            v-if="filtered[0]"
            :origin="filtered[0].originCoords"
            :destination="filtered[0].destinationCoords"
            :current="filtered[0].currentCoords"
            :origin-label="filtered[0].origin"
            :destination-label="filtered[0].destination"
            :carrier="filtered[0].carrier"
            :status="filtered[0].status"
          />
        </div>
        <template #fallback>
          <div class="flex h-[300px] w-full items-center justify-center bg-gradient-hero">
            <div class="font-mono text-xs uppercase tracking-widest text-muted-foreground">
              Loading geo telemetry…
            </div>
          </div>
        </template>
      </Suspense>
    </div>
  </div>
</template>
```

- [ ] **Step 2: Commit**

```bash
git add frontend/src/components/DeliveriesPanel.vue
git commit -m "feat: add DeliveriesPanel with live status updates and map"
```

---

### Task 10: Lint check & build verification

- [ ] **Step 1: Run TypeScript check and build**

```bash
cd frontend && npm run build
```

Expected: clean build with no type errors.

If there are type errors, fix them (likely import path or type mismatch issues) and re-run.

- [ ] **Step 2: If build passes, commit any final fixes**

```bash
git add -A
git commit -m "chore: fix type issues from carriers page build"
```
