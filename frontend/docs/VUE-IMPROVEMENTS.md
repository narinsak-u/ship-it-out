# Vue Best Practices Improvement Plan

Audit conducted against `vue-best-practices` skill: reactivity, SFC structure, component boundaries, data flow, composables.

---

## 🔴 Phase 1 — Component Splitting (high impact)

### 1.1 Split `OrdersView.vue` (279 lines)

| Current | Extract to | Responsibility |
|---------|-----------|----------------|
| Search input + filter pills inline | `components/ShipmentFilters.vue` | Search query input, status filter chip bar. Props: `modelValue` (search), `modelValue` (filter). Emits: `update:modelValue`. |
| `<TableRow>` block (8 columns + actions) | `components/OrderTableRow.vue` | Row display for a single order. Props: `order: Order`. Emits: `delete`, `edit`. |
| Loading skeleton | `components/OrdersSkeleton.vue` | Skeleton placeholder matching the table layout. |

Keep `OrdersView` as a thin composition surface: hero header, `ShipmentFilters`, `OrderTableRow` usage in `<Table>`, `Pagination`, modals.

### 1.2 Split `DeliveriesPanel.vue` (307 lines)

| Current | Extract to | Responsibility |
|---------|-----------|----------------|
| Table row with inline Select (status + hub) | `components/DeliveryTableRow.vue` | Row with status/hub Select, update button. Props: `order`, `hubs`, `draftStatus`, `draftHubId`. Emits: `update:status`, `update:hub`, `update`. |
| Ping timer "Xs ago" label | `components/PingTimer.vue` or inline computed | Derived display of seconds since last refetch. |

Keep `DeliveriesPanel` as composition surface: search, table container, map.

### 1.3 Split `AuthModal.vue` (235 lines)

| Current | Extract to | Responsibility |
|---------|-----------|----------------|
| Login form block | `components/LoginForm.vue` | Email + password fields, submit, validation, guest/switch links. |
| Signup form block | `components/SignupForm.vue` | Name + email + password + confirm, submit, validation, switch link. |
| Form validation + submission | `composables/useAuthForm.ts` | Shared login/signup logic, error handling, submitting state. |

Keep `AuthModal` as shell: dialog wrapper, tab switcher, delegates to `LoginForm` / `SignupForm`.

### 1.4 Extract `HubsPanel` stats cards section

Move the 5 stat cards (Total, Active, Maintenance, Full, Closed) into `components/HubStatsCards.vue`. Props: `counts: HubStatusCounts`.

---

## 🟡 Phase 2 — Reactivity & lifecycle fixes (medium impact)

### 2.1 Fix debounce timer leak in `OrdersView.vue`

Replace `let debounceTimer` at setup scope with watcher `onCleanup`:

```ts
watch(searchInput, (v, _old, onCleanup) => {
  const timer = setTimeout(() => {
    debouncedSearch.value = v;
    currentPage.value = 1;
  }, 300);
  onCleanup(() => clearTimeout(timer));
});
```

File: `src/views/OrdersView.vue:39-46`

### 2.2 Extract geocoding logic from `OrderForm.vue`

Move the `geocodeAddress` + `Promise.allSettled` pattern into `composables/useGeocodeSubmit.ts`:

```ts
function useGeocodeSubmit() {
  const geocoding = ref(false);
  const geocodeErrors = ref<Record<string, string>>({});
  async function geocodeBoth(sender, receiver) { /* ... */ }
  return { geocoding, geocodeErrors, geocodeBoth };
}
```

### 2.3 Fix HubFormModal watcher side effects

Replace the `watch([currentUtilization, capacity])` with:
- `utilizationError` as `computed`
- Separate `watch` for auto-setting status to "full" only when util === cap

```ts
const utilizationError = computed(() => {
  if (currentUtilization.value > capacity.value)
    return `Cannot exceed capacity (${capacity.value})`;
  return "";
});
```

---

## 🟢 Phase 3 — Minor improvements (low impact)

### 3.1 Replace `mounted` workaround in OrderDetailView

`OrderDetailView.vue:35-38`: the `mounted` ref + `v-if="mounted"` on `<Suspense>` is a hydration safety pattern. Replace with a simpler approach — wrap the async component trigger in a client-only expression, or remove the workaround entirely (SSR is not in use).

### 3.2 Use typed query key factory for order-events

`OrderDetailView.vue:24`: replace `["order-events", orderId]` with an entry in `queryKeys.ts`.

```ts
export const eventKeys = {
  all: ["events"] as const,
  byTracking: (tn: string) => [...eventKeys.all, tn] as const,
};
```

### 3.3 Tighten `StatusPieChart.tooltipFn` type

`StatusPieChart.vue:33`: change `Record<string, any>` to `{ data?: StatusPieEntry }`.

### 3.4 Normalize import ordering

Consistently apply: `vue` → `vue-router` → `vue-query/sonner/lucide` → `@/components` → `@/lib` → `@/hooks`/`@/composables`/`@/stores` → relative

### 3.5 `shallowRef` for primitive state

Per reactivity reference: use `shallowRef()` instead of `ref()` for primitive values (string, number, boolean). Apply across all files.

---

## Execution order

1. Phase 1 (component splits) — largest behavioral change, touch many files
2. Phase 2 (reactivity/lifecycle) — targeted fixes, no new files
3. Phase 3 (minor) — lint-level improvements, safe to batch
