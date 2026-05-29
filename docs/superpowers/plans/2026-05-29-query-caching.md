# Query Caching Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) for syntax tracking.

**Goal:** Configure TanStack Vue Query with `staleTime` defaults and a query key factory to eliminate unnecessary refetches on page navigation, while keeping mutation-driven real-time refresh.

**Architecture:** One `QueryClient` config in `main.ts` sets global defaults. A `queryKeys.ts` factory provides hierarchical keys for precise cache invalidation. Each view/component overrides `staleTime` per query. Mutations invalidate via factory keys — active queries refetch immediately, inactive ones refetch on next mount.

**Tech Stack:** Vue 3, TanStack Vue Query v5, TypeScript

---

### Task 1: Configure QueryClient + Create Query Key Factory

**Files:**
- Modify: `frontend/src/main.ts`
- Create: `frontend/src/lib/api/queryKeys.ts`

- [ ] **Step 1: Add configured QueryClient to main.ts**

Edit `frontend/src/main.ts`. Replace the bare `app.use(VueQueryPlugin)` with a configured `QueryClient`:

```
import { createApp } from "vue";
import { createPinia } from "pinia";
-import { VueQueryPlugin } from "@tanstack/vue-query";
+import { VueQueryPlugin, QueryClient } from "@tanstack/vue-query";
import App from "./App.vue";
import router from "./router";
import "./styles.css";

const app = createApp(App);

app.use(createPinia());
app.use(router);
-app.use(VueQueryPlugin);
+
+const queryClient = new QueryClient({
+  defaultOptions: {
+    queries: {
+      staleTime: 30_000,
+      gcTime: 5 * 60_000,
+      retry: 1,
+      refetchOnWindowFocus: true,
+    },
+  },
+});
+
+app.use(VueQueryPlugin, { queryClient });

app.mount("#app");
```

- [ ] **Step 2: Create queryKeys.ts factory**

Create `frontend/src/lib/api/queryKeys.ts`:

```ts
export const orderKeys = {
  all: ["orders"] as const,
  lists: () => [...orderKeys.all, "list"] as const,
  list: (filters: Record<string, unknown>) => [...orderKeys.lists(), filters] as const,
  details: () => [...orderKeys.all, "detail"] as const,
  detail: (id: string) => [...orderKeys.details(), id] as const,
};

export const deliveryKeys = {
  all: ["deliveries"] as const,
  active: () => [...deliveryKeys.all, "active"] as const,
};

export const hubKeys = {
  all: ["hubs"] as const,
};

export const analyticsKeys = {
  all: ["analytics"] as const,
  timeseries: () => [...analyticsKeys.all, "timeseries"] as const,
};
```

- [ ] **Step 3: Verify build**

Run: `cd frontend && npm run build` — must pass cleanly.

If the import of `QueryClient` causes tree-shaking issues (it shouldn't with `@tanstack/vue-query` v5), check the error. This step validates the foundation is sound before dependent tasks.

- [ ] **Step 4: Commit**

```bash
cd frontend
git add src/main.ts src/lib/api/queryKeys.ts
git commit -m "feat: configure QueryClient staleTime and add query key factory"
```

---

### Task 2: Update Mutation Hooks to Use Factory Keys

**Files:**
- Modify: `frontend/src/hooks/useOrders.ts`
- Modify: `frontend/src/hooks/useDeliveries.ts`
- Modify: `frontend/src/hooks/useHubs.ts`

- [ ] **Step 1: Update useOrders.ts**

Edit `frontend/src/hooks/useOrders.ts`:

```
 import { useMutation, useQueryClient } from "@tanstack/vue-query";
+import { orderKeys, deliveryKeys } from "@/lib/api/queryKeys";
 import { createOrder, updateOrder, type OrderFormData } from "@/lib/api/orders";

 export function useCreateOrder() {
   const queryClient = useQueryClient();
   return useMutation({
     mutationFn: (data: OrderFormData) => createOrder(data),
     onSuccess: () => {
-      queryClient.invalidateQueries({ queryKey: ["deliveries"] });
+      queryClient.invalidateQueries({ queryKey: orderKeys.all });
+      queryClient.invalidateQueries({ queryKey: deliveryKeys.all });
     },
   });
 }

 export function useUpdateOrder() {
   const queryClient = useQueryClient();
   return useMutation({
     mutationFn: ({ id, data }: { id: string; data: Partial<OrderFormData> }) =>
       updateOrder(id, data),
     onSuccess: () => {
-      queryClient.invalidateQueries({ queryKey: ["deliveries"] });
+      queryClient.invalidateQueries({ queryKey: orderKeys.all });
+      queryClient.invalidateQueries({ queryKey: deliveryKeys.all });
     },
   });
 }
```

- [ ] **Step 2: Update useDeliveries.ts**

Edit `frontend/src/hooks/useDeliveries.ts`:

```
 import { useQuery, useMutation, useQueryClient } from "@tanstack/vue-query";
+import { deliveryKeys } from "@/lib/api/queryKeys";
 import { fetchActiveDeliveries, updateShipmentStatus } from "@/lib/api/orders";
 import type { ShipmentStatus } from "@/lib/orders";

 export function useActiveDeliveries() {
   return useQuery({
-    queryKey: ["deliveries"],
+    queryKey: deliveryKeys.active(),
     queryFn: fetchActiveDeliveries,
     refetchInterval: 15_000,
   });
 }

 export function useUpdateShipmentStatus() {
   const queryClient = useQueryClient();
   return useMutation({
     mutationFn: ({ orderId, status, hubId }: { orderId: string; status: ShipmentStatus; hubId?: string }) =>
       updateShipmentStatus(orderId, status, hubId),
     onSuccess: () => {
-      queryClient.invalidateQueries({ queryKey: ["deliveries"] });
+      queryClient.invalidateQueries({ queryKey: deliveryKeys.all });
     },
   });
 }
```

- [ ] **Step 3: Update useHubs.ts**

Edit `frontend/src/hooks/useHubs.ts`:

```
 import { toast } from "vue-sonner";
 import { useQuery, useMutation, useQueryClient } from "@tanstack/vue-query";
+import { hubKeys } from "@/lib/api/queryKeys";
 import { fetchHubs, createHub, updateHub, deleteHub } from "@/lib/api/hubs";
 import type { Hub } from "@/lib/hubs";

 export function useHubs() {
   return useQuery({
-    queryKey: ["hubs"],
+    queryKey: hubKeys.all,
     queryFn: fetchHubs,
   });
 }

 export function useCreateHub() {
   const queryClient = useQueryClient();
   return useMutation({
     mutationFn: (data: Omit<Hub, "id">) => createHub(data),
     onSuccess: () => {
-      queryClient.invalidateQueries({ queryKey: ["hubs"] });
+      queryClient.invalidateQueries({ queryKey: hubKeys.all });
     },
   });
 }

 export function useUpdateHub() {
   const queryClient = useQueryClient();
   return useMutation({
     mutationFn: ({ id, data }: { id: string; data: Partial<Hub> }) => updateHub(id, data),
     onSuccess: () => {
-      queryClient.invalidateQueries({ queryKey: ["hubs"] });
+      queryClient.invalidateQueries({ queryKey: hubKeys.all });
     },
   });
 }

 export function useDeleteHub() {
   const queryClient = useQueryClient();
   return useMutation({
     mutationFn: (id: string) => deleteHub(id),
     onSuccess: () => {
       toast.success("Hub deleted");
-      queryClient.invalidateQueries({ queryKey: ["hubs"] });
+      queryClient.invalidateQueries({ queryKey: hubKeys.all });
     },
   });
 }
```

- [ ] **Step 4: Verify build**

Run: `cd frontend && npm run build` — must pass cleanly. This validates all imports and type references.

- [ ] **Step 5: Commit**

```bash
cd frontend
git add src/hooks/useOrders.ts src/hooks/useDeliveries.ts src/hooks/useHubs.ts
git commit -m "feat: update mutation hooks to use query key factory"
```

---

### Task 3: Update Views — HomeView + OrdersView

**Files:**
- Modify: `frontend/src/views/HomeView.vue`
- Modify: `frontend/src/views/OrdersView.vue`

- [ ] **Step 1: Update HomeView.vue**

Edit `frontend/src/views/HomeView.vue`. Two changes: import factory + swap keys + add staleTime.

Add import:
```
+import { deliveryKeys, analyticsKeys } from "@/lib/api/queryKeys";
 import { fetchActiveDeliveries } from "@/lib/api/orders";
```

Swap query keys and add staleTime:
```
 const { data: orders } = useQuery({
-  queryKey: ["orders"],
+  queryKey: deliveryKeys.active(),
   queryFn: fetchActiveDeliveries,
+  staleTime: 10_000,
 });

 const { data: analytics } = useQuery({
-  queryKey: ["analytics"],
+  queryKey: analyticsKeys.all,
   queryFn: fetchAnalytics,
+  staleTime: 5 * 60_000,
 });
```

- [ ] **Step 2: Update OrdersView.vue**

Edit `frontend/src/views/OrdersView.vue`. Three changes: import factory, swap query key, update inline delete mutation.

Add import:
```
+import { orderKeys, deliveryKeys } from "@/lib/api/queryKeys";
 import { fetchOrdersPaginated, deleteOrder } from "@/lib/api/orders";
```

Swap query key and add staleTime:
```
 const { data: pageData, isLoading } = useQuery({
-  queryKey: ["orders", currentPage, debouncedSearch, filter],
+  queryKey: orderKeys.list({ page: currentPage.value, search: debouncedSearch.value, status: filter.value }),
   queryFn: () =>
     fetchOrdersPaginated({
       page: currentPage.value,
       limit: 10,
       search: debouncedSearch.value || undefined,
       status: filter.value === "all" ? undefined : filter.value,
     }),
+  staleTime: 60_000,
 });
```

Update inline delete mutation invalidation:
```
   onSuccess: () => {
     toast.success("Order deleted");
     deleteTarget.value = null;
-    queryClient.invalidateQueries({ queryKey: ["orders"] });
+    queryClient.invalidateQueries({ queryKey: orderKeys.all });
+    queryClient.invalidateQueries({ queryKey: deliveryKeys.all });
   },
```

- [ ] **Step 3: Verify build**

Run: `cd frontend && npm run build` — must pass cleanly.

- [ ] **Step 4: Commit**

```bash
cd frontend
git add src/views/HomeView.vue src/views/OrdersView.vue
git commit -m "feat: update HomeView and OrdersView with factory keys and staleTime"
```

---

### Task 4: Update Views — OrderDetailView + OrderFormView

**Files:**
- Modify: `frontend/src/views/OrderDetailView.vue`
- Modify: `frontend/src/views/OrderFormView.vue`

- [ ] **Step 1: Update OrderDetailView.vue**

Edit `frontend/src/views/OrderDetailView.vue`. Add import, swap keys, add staleTime.

Add import:
```
+import { orderKeys } from "@/lib/api/queryKeys";
 import { fetchOrder, fetchOrderEvents } from "@/lib/api/orders";
```

Update order query:
```
 const { data: order, isLoading } = useQuery({
-  queryKey: ["order", orderId],
+  queryKey: orderKeys.detail(orderId),
   queryFn: () => fetchOrder(orderId),
+  staleTime: 60_000,
 });
```

Update events query (keeps its own key — separate resource from order detail):
```
 const { data: events } = useQuery({
   queryKey: ["order-events", orderId],
   queryFn: () => {
     if (!order.value) return [];
     return fetchOrderEvents(order.value.trackingNumber);
   },
   enabled: computed(() => !!order.value),
+  staleTime: 60_000,
 });
```

- [ ] **Step 2: Update OrderFormView.vue**

Edit `frontend/src/views/OrderFormView.vue`. Add import and swap query key.

Add import:
```
+import { orderKeys } from "@/lib/api/queryKeys";
 import { fetchOrder } from "@/lib/api/orders";
```

Swap query key:
```
 const { data: order } = useQuery({
-  queryKey: ["order", orderId.value],
+  queryKey: orderKeys.detail(orderId.value),
   queryFn: () => fetchOrder(orderId.value!),
   enabled: isEditing,
 });
```

- [ ] **Step 3: Verify build**

Run: `cd frontend && npm run build` — must pass cleanly.

- [ ] **Step 4: Commit**

```bash
cd frontend
git add src/views/OrderDetailView.vue src/views/OrderFormView.vue
git commit -m "feat: update OrderDetailView and OrderFormView with factory keys and staleTime"
```

---

### Task 5: Update Components — AnalyticsPanel + DeliveriesPanel + HubsPanel

**Files:**
- Modify: `frontend/src/components/AnalyticsPanel.vue`
- Modify: `frontend/src/components/DeliveriesPanel.vue`
- Modify: `frontend/src/components/HubsPanel.vue`

- [ ] **Step 1: Update AnalyticsPanel.vue**

Edit `frontend/src/components/AnalyticsPanel.vue`. Add import, swap keys, add staleTime.

Add import:
```
+import { analyticsKeys } from "@/lib/api/queryKeys";
 import { fetchAnalytics, fetchTimeSeries } from "@/lib/api/analytics";
```

Update analytics query:
```
 const {
   data: analytics,
   isLoading,
   isError,
   refetch,
 } = useQuery({
-  queryKey: ["analytics"],
+  queryKey: analyticsKeys.all,
   queryFn: fetchAnalytics,
+  staleTime: 5 * 60_000,
 });
```

Update timeseries query:
```
 const { data: timeSeries } = useQuery({
-  queryKey: ["analytics", "timeseries"],
+  queryKey: analyticsKeys.timeseries(),
   queryFn: fetchTimeSeries,
+  staleTime: 5 * 60_000,
 });
```

- [ ] **Step 2: Update DeliveriesPanel.vue**

Edit `frontend/src/components/DeliveriesPanel.vue`. Add import, swap keys, add staleTime.

Add import:
```
+import { deliveryKeys, hubKeys } from "@/lib/api/queryKeys";
 import { useActiveDeliveries, useUpdateShipmentStatus } from "@/hooks/useDeliveries";
```

Update hubs query:
```
 const { data: hubs } = useQuery({
-  queryKey: ["hubs"],
+  queryKey: hubKeys.all,
   queryFn: fetchHubs,
+  staleTime: 2 * 60_000,
 });
```

- [ ] **Step 3: Update HubsPanel.vue**

Edit `frontend/src/components/HubsPanel.vue`. No import needed (useHubs already uses factory keys from the hook file change in Task 2).

No changes needed — the hook already returns the data, and the component doesn't define its own query keys.

Actually wait — let me reconsider. `HubsPanel.vue` uses `useHubs()` from hooks, not an inline `useQuery`. The query key is defined inside `useHubs()` which was already updated in Task 2. So no changes needed.

- [ ] **Step 4: Verify build**

Run: `cd frontend && npm run build` — must pass cleanly.

- [ ] **Step 5: Commit**

```bash
cd frontend
git add src/components/AnalyticsPanel.vue src/components/DeliveriesPanel.vue
git commit -m "feat: update AnalyticsPanel and DeliveriesPanel with factory keys and staleTime"
```

---

### Task 6: Final Build Verification

- [ ] **Step 1: Run full build**

Run: `cd frontend && npm run build`

Expected: clean exit, no TypeScript errors, no lint errors.

If build fails, fix any issues (likely an import path or key type mismatch).

- [ ] **Step 2: Squash commits (optional)**

If the user prefers a single clean commit, squash the 5 commits:

```bash
git reset --soft HEAD~5 && git commit -m "feat: configure query caching with staleTime and key factory"
```

Otherwise the 5 incremental commits are fine as-is.
