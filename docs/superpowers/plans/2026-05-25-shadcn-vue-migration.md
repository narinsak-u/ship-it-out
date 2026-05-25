# Shadcn-vue UI Component Migration Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace all custom UI patterns with shadcn-vue components while preserving the Ocean Deep dark theme visually.

**Architecture:** Component-by-component migration — each of 6 work areas is independently testable. No logic changes, only template/component swaps. Theme stays unchanged in `styles.css`.

**Tech Stack:** Vue 3, shadcn-vue (new-york style, radix-vue), Tailwind CSS v4, Bun

---

## File Structure

### Files to Create (via shadcn-vue CLI):
- `src/components/ui/dialog/` — shadcn-vue Dialog (DialogRoot, DialogContent, DialogOverlay, etc.)
- `src/components/ui/select/` — shadcn-vue Select (SelectRoot, SelectTrigger, SelectItem, etc.)

### Files to Modify:

#### Dialog migration:
- `src/components/AuthModal.vue` — replace custom backdrop with Dialog components
- `src/components/SiteHeader.vue` — wrap sign-in button with DialogTrigger
- `src/components/HubFormModal.vue` — replace custom backdrop with Dialog components
- `src/components/HubsPanel.vue` — wrap "Add Hub" button with DialogTrigger

#### Select migration:
- `src/components/DeliveriesPanel.vue` — 2 native `<select>` to Select
- `src/components/ThaiAddressGroup.vue` — 1 native `<select>` to Select
- `src/components/HubFormModal.vue` — 1 native `<select>` to Select

#### Table migration:
- `src/components/DeliveriesPanel.vue` — CSS-grid table to Table components
- `src/components/HubsPanel.vue` — CSS-grid table to Table components
- `src/components/OrdersView.vue` — CSS-grid table to Table components

#### Badge migration:
- `src/components/ui/Badge.vue` — extend `badgeVariants` with 7 shipment status variants
- `src/components/StatusBadge.vue` — use shadcn Badge internally
- `src/components/HubsPanel.vue` — replace inline badge classes with Badge

#### Card migration:
- All files with manual card patterns (~15 locations)

#### Separator migration:
- Files with manual `border-t/border-b` dividers

---

### Task 1: Install shadcn-vue Dialog and Select

**Files:**
- Create: `src/components/ui/dialog/` (auto-generated)
- Create: `src/components/ui/select/` (auto-generated)

- [ ] **Step 1: Install Dialog and Select**

Run: `bunx shadcn-vue@latest add dialog select`
Workdir: `frontend/`
Expected: Creates `src/components/ui/dialog/Dialog.vue`, `src/components/ui/dialog/index.ts`, `src/components/ui/select/Select.vue`, `src/components/ui/select/index.ts` with supporting sub-components.

- [ ] **Step 2: Verify build**

Run: `npm run build`
Workdir: `frontend/`
Expected: Build passes with zero errors.

- [ ] **Step 3: Commit**

```bash
git add frontend/src/components/ui/dialog/ frontend/src/components/ui/select/
git commit -m "feat: add shadcn-vue dialog and select components"
```

---

### Task 2: Migrate AuthModal to Dialog

**Files:**
- Modify: `src/components/AuthModal.vue`
- Modify: `src/components/SiteHeader.vue`

**Context:** `AuthModal.vue` currently opens via `showAuthModal` ref in `SiteHeader.vue`. The modal shows a custom backdrop with tabbed sign-in/sign-up forms.

- [ ] **Step 1: Update AuthModal.vue**

Replace the custom backdrop wrapper with Dialog components. The form content inside stays identical.

Current wrapper (lines ~69-78):
```vue
<div class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm"
     @click.self="emit('close')">
  <div class="w-full max-w-md rounded-xl border border-border bg-card p-6 shadow-elegant">
    <!-- tabs, forms, etc. -->
  </div>
</div>
```

Replace with:
```vue
<DialogRoot :open="open" @update:open="emit('close')">
  <DialogPortal>
    <DialogOverlay class="bg-black/60 backdrop-blur-sm" />
    <DialogContent class="w-full max-w-md rounded-xl border border-border bg-card p-6 shadow-elegant sm:max-w-md">
      <DialogHeader>
        <DialogTitle class="sr-only">Sign In / Sign Up</DialogTitle>
      </DialogHeader>
      <!-- tabs, forms, etc. (unchanged) -->
    </DialogContent>
  </DialogPortal>
</DialogRoot>
```

Update the props — change from emitting `close` on backdrop click to accepting `open` as a prop:
```ts
const props = defineProps<{
  open?: boolean;
}>();
const emit = defineEmits<{
  "update:open": [value: boolean];
}>();
```

Remove the `@click.self` event and the old backdrop div. Remove `emit('close')` calls — replace with `emit("update:open", false)`.

- [ ] **Step 2: Update SiteHeader.vue**

Current (sign-in button):
```vue
<button @click="showAuthModal = true" class="...">
  <LogIn class="h-3.5 w-3.5" /> Sign in
</button>
<AuthModal @close="showAuthModal = false" v-if="showAuthModal" />
```

Replace with:
```vue
<AuthModal :open="showAuthModal" @update:open="showAuthModal = $event" />
```

And the sign-in button is now just a trigger (no need for `@click` to open — Dialog handles it via `v-model:open`). Wait, actually since `AuthModal` wraps its own `DialogRoot`, the trigger needs to be inside the `DialogRoot`. But the trigger button is in `SiteHeader` while the modal content is in `AuthModal`.

The cleanest approach: keep the `showAuthModal` state in `SiteHeader`, pass it as `:open` to `AuthModal`, and let the `@update:open` handler toggle it. Remove `v-if` since Dialog handles show/hide via `:open`.

Actually, since the trigger is in SiteHeader but the DialogRoot is in AuthModal, the v-model approach needs adjustment. Simplest: keep using `showAuthModal` boolean + `@update:open` binding.

```vue
<button @click="showAuthModal = true" class="...">
  <LogIn class="h-3.5 w-3.5" /> Sign in
</button>
<AuthModal :open="showAuthModal" @update:open="showAuthModal = $event" />
```

Remove the old `v-if="showAuthModal"` — Dialog handles visibility internally.

- [ ] **Step 3: Verify build**

Run: `npm run build`
Workdir: `frontend/`
Expected: Build passes.

- [ ] **Step 4: Commit**

```bash
git add frontend/src/components/AuthModal.vue frontend/src/components/SiteHeader.vue
git commit -m "feat: migrate AuthModal to shadcn-vue Dialog"
```

---

### Task 3: Migrate HubFormModal to Dialog

**Files:**
- Modify: `src/components/HubFormModal.vue`
- Modify: `src/components/HubsPanel.vue`

**Context:** Same pattern as AuthModal — custom backdrop with form content.

- [ ] **Step 1: Update HubFormModal.vue**

Replace custom backdrop with Dialog components (same pattern as AuthModal):

```vue
<DialogRoot :open="open" @update:open="emit('close')">
  <DialogPortal>
    <DialogOverlay class="bg-black/60 backdrop-blur-sm" />
    <DialogContent class="w-full max-w-md rounded-xl border border-border bg-card p-6 shadow-elegant sm:max-w-md">
      <DialogHeader>
        <DialogTitle>{{ isEditing ? 'Edit Hub' : 'Add Hub' }}</DialogTitle>
      </DialogHeader>
      <!-- form content (unchanged) -->
    </DialogContent>
  </DialogPortal>
</DialogRoot>
```

Update props:
```ts
const props = defineProps<{
  open?: boolean;
  hub?: Hub | null;
}>();
const emit = defineEmits<{
  close: [];
  saved: [];
}>();
```

Replace `emit('close')` calls with `emit("close")` — the `@update:open` handler in the parent will also call close.

- [ ] **Step 2: Update HubsPanel.vue**

Find the HubFormModal usage and `showForm` ref. Change from:
```vue
<HubFormModal v-if="showForm" :hub="editingHub" @close="showForm = false" @saved="..." />
```

Keep the same pattern since HubFormModal manages its own DialogRoot. The overlay and backdrop are handled by Dialog internally.

But update the "Add Hub" button and edit buttons — no need to change their click handlers since they set `showForm = true`. Dialog handles the open state via the `:open` prop.

Actually, looking at the HubFormModal's existing code, it uses `v-if="showForm"` which destroys and recreates the component. With Dialog, we should pass `:open="showForm"` and let Dialog manage visibility:

```vue
<HubFormModal :open="showForm" :hub="editingHub" @close="showForm = false" @saved="..." />
```

Remove `v-if`.

- [ ] **Step 3: Verify build**

Run: `npm run build`
Workdir: `frontend/`
Expected: Build passes.

- [ ] **Step 4: Commit**

```bash
git add frontend/src/components/HubFormModal.vue frontend/src/components/HubsPanel.vue
git commit -m "feat: migrate HubFormModal to shadcn-vue Dialog"
```

---

### Task 4: Migrate Native Selects — DeliveriesPanel

**Files:**
- Modify: `src/components/DeliveriesPanel.vue`

**Context:** Two native `<select>` elements per row — status selector and hub selector.

- [ ] **Step 1: Replace status `<select>` with Select**

Current:
```vue
<select
  :value="draftStatus[o.id] ?? o.status"
  @change="draftStatus[o.id] = ($event.target as HTMLSelectElement).value as ShipmentStatus"
  :disabled="!auth.isAuthenticated"
  class="rounded-lg border border-border bg-background px-2 py-1 font-mono text-xs disabled:opacity-40"
>
  <option v-for="(label, key) in statusLabels" :key="key" :value="key">
    {{ label }}
  </option>
</select>
```

Replace with:
```vue
<SelectRoot
  :model-value="draftStatus[o.id] ?? o.status"
  @update:model-value="(v) => draftStatus[o.id] = v as ShipmentStatus"
  :disabled="!auth.isAuthenticated"
>
  <SelectTrigger class="h-7 rounded-lg border border-border bg-background px-2 font-mono text-xs disabled:opacity-40">
    <SelectValue />
  </SelectTrigger>
  <SelectContent>
    <SelectGroup>
      <SelectItem v-for="(label, key) in statusLabels" :key="key" :value="key">
        {{ label }}
      </SelectItem>
    </SelectGroup>
  </SelectContent>
</SelectRoot>
```

- [ ] **Step 2: Replace hub `<select>` with Select**

Current:
```vue
<select
  v-if="usesHubSelector(draftStatus[o.id] ?? o.status)"
  :value="draftHubId[o.id] ?? o.hubId ?? ''"
  @change="draftHubId[o.id] = ($event.target as HTMLSelectElement).value"
  :disabled="!auth.isAuthenticated"
  class="w-full rounded-lg border border-border bg-background px-2 py-1 font-mono text-xs disabled:opacity-40"
>
  <option disabled value="">Select hub...</option>
  <option v-for="h in hubOptions" :key="h.id" :value="h.id">{{ h.name }}</option>
</select>
```

Replace with:
```vue
<SelectRoot
  v-if="usesHubSelector(draftStatus[o.id] ?? o.status)"
  :model-value="draftHubId[o.id] ?? o.hubId ?? ''"
  @update:model-value="(v) => draftHubId[o.id] = v"
  :disabled="!auth.isAuthenticated"
>
  <SelectTrigger class="h-7 w-full rounded-lg border border-border bg-background px-2 font-mono text-xs disabled:opacity-40">
    <SelectValue placeholder="Select hub..." />
  </SelectTrigger>
  <SelectContent>
    <SelectGroup>
      <SelectItem v-for="h in hubOptions" :key="h.id" :value="h.id">
        {{ h.name }}
      </SelectItem>
    </SelectGroup>
  </SelectContent>
</SelectRoot>
```

- [ ] **Step 3: Verify build**

Run: `npm run build`
Workdir: `frontend/`
Expected: Build passes.

- [ ] **Step 4: Commit**

```bash
git add frontend/src/components/DeliveriesPanel.vue
git commit -m "feat: migrate DeliveriesPanel selects to shadcn-vue Select"
```

---

### Task 5: Migrate Native Select — ThaiAddressGroup

**Files:**
- Modify: `src/components/ThaiAddressGroup.vue`

- [ ] **Step 1: Replace sub-district `<select>` with Select**

Current:
```vue
<select
  v-if="lookupStatus === 'has-results'"
  :value="modelValue.subDistrict"
  @change="onSubDistrictSelected"
  class="mt-1.5 flex h-10 w-full rounded-lg border border-border bg-background px-3 font-mono text-sm"
>
  <option disabled value="">Select sub-district...</option>
  <option v-for="sd in subDistricts" :key="sd" :value="sd">{{ sd }}</option>
</select>
```

Replace with:
```vue
<SelectRoot
  v-if="lookupStatus === 'has-results'"
  :model-value="modelValue.subDistrict"
  @update:model-value="onSubDistrictSelected"
>
  <SelectTrigger class="mt-1.5 flex h-10 w-full rounded-lg border border-border bg-background px-3 font-mono text-sm">
    <SelectValue placeholder="Select sub-district..." />
  </SelectTrigger>
  <SelectContent>
    <SelectGroup>
      <SelectItem v-for="sd in subDistricts" :key="sd" :value="sd">{{ sd }}</SelectItem>
    </SelectGroup>
  </SelectContent>
</SelectRoot>
```

Note: `onSubDistrictSelected` currently receives a native Event. With Select's `@update:model-value`, it receives the string value directly. Update the handler:
```ts
function onSubDistrictSelected(value: string) {
  // ... existing logic using `value` instead of `(event.target as HTMLSelectElement).value`
}
```

- [ ] **Step 2: Verify build**

Run: `npm run build`
Workdir: `frontend/`
Expected: Build passes.

- [ ] **Step 3: Commit**

```bash
git add frontend/src/components/ThaiAddressGroup.vue
git commit -m "feat: migrate ThaiAddressGroup select to shadcn-vue Select"
```

---

### Task 6: Migrate Native Select — HubFormModal

**Files:**
- Modify: `src/components/HubFormModal.vue`

- [ ] **Step 1: Replace hub status `<select>` with Select**

After the Dialog migration, find the hub status select in the form body. Replace:

Current:
```vue
<select
  v-model="status"
  class="mt-1.5 flex h-10 w-full rounded-lg border border-border bg-background px-3 font-mono text-sm"
>
  <option value="active">Active</option>
  <option value="maintenance">Maintenance</option>
  <option value="closed">Closed</option>
</select>
```

Replace with:
```vue
<SelectRoot v-model="status">
  <SelectTrigger class="mt-1.5 flex h-10 w-full rounded-lg border border-border bg-background px-3 font-mono text-sm">
    <SelectValue placeholder="Select status..." />
  </SelectTrigger>
  <SelectContent>
    <SelectGroup>
      <SelectItem value="active">Active</SelectItem>
      <SelectItem value="maintenance">Maintenance</SelectItem>
      <SelectItem value="closed">Closed</SelectItem>
    </SelectGroup>
  </SelectContent>
</SelectRoot>
```

- [ ] **Step 2: Verify build**

Run: `npm run build`
Workdir: `frontend/`
Expected: Build passes.

- [ ] **Step 3: Commit**

```bash
git add frontend/src/components/HubFormModal.vue
git commit -m "feat: migrate HubFormModal select to shadcn-vue Select"
```

---

### Task 7: Migrate Table — DeliveriesPanel

**Files:**
- Modify: `src/components/DeliveriesPanel.vue`

**Context:** The active deliveries list uses a CSS-grid layout for the "table" header and rows. Replace with shadcn-vue Table components.

- [ ] **Step 1: Replace header grid with TableHeader**

Current (full table wrapper + header):
```vue
<div class="mt-4 overflow-hidden rounded-xl border border-border">
  <div class="hidden grid-cols-[0.8fr_1fr_1fr_1fr_1fr_0.8fr_0.8fr_1fr_0.5fr] gap-4 border-b border-border bg-secondary/50 px-6 py-3 font-mono text-[11px] uppercase tracking-widest text-muted-foreground md:grid">
    <span>Order ID</span>
    <span>Tracking</span>
    <span>Customer</span>
    <span>Carrier</span>
    <span>Status</span>
    <span>Hub</span>
    <span>ETA</span>
    <span class="text-right">Actions</span>
  </div>
```

Replace with:
```vue
<div class="mt-4 overflow-hidden rounded-xl border border-border">
  <Table>
    <TableHeader>
      <TableRow class="border-b border-border bg-secondary/50 font-mono text-[11px] uppercase tracking-widest text-muted-foreground hover:bg-secondary/50">
        <TableHead>Order ID</TableHead>
        <TableHead>Tracking</TableHead>
        <TableHead>Customer</TableHead>
        <TableHead>Carrier</TableHead>
        <TableHead>Status</TableHead>
        <TableHead>Hub</TableHead>
        <TableHead>ETA</TableHead>
        <TableHead class="text-right">Actions</TableHead>
      </TableRow>
    </TableHeader>
```

Note: The `hidden md:grid` responsive behavior. The Table component renders as a native `<table>` — on small screens it will need horizontal scroll. Wrap the table in a scrollable container or keep the existing `overflow-hidden rounded-xl border border-border` wrapper (which doesn't scroll). The current CSS-grid design hides the header on mobile (`hidden md:grid`). With native `<table>`, we hide the `<TableHeader>` on mobile using `class="hidden md:table-header-group"`.

- [ ] **Step 2: Replace row data grid with TableBody + TableRow + TableCell**

Current data rows:
```vue
<div v-for="o in filtered" :key="o.id"
     class="group grid grid-cols-1 gap-2 border-b border-border px-6 py-4 transition-colors last:border-0 hover:bg-secondary/40 md:grid-cols-[0.8fr_1fr_1fr_1fr_1fr_0.8fr_0.8fr_1fr_0.5fr] md:items-center">
  <div class="font-mono text-sm text-primary">{{ o.id }}</div>
  <div class="font-mono text-xs text-muted-foreground">{{ o.trackingNumber }}</div>
  <div class="font-mono text-sm">{{ o.customer.name }}</div>
  <div class="font-mono text-sm text-muted-foreground">{{ o.carrier }}</div>
  <div><!-- Select for status (already migrated in Task 4) --></div>
  <div><!-- Select for hub (already migrated) or em-dash --></div>
  <div class="font-mono text-xs text-muted-foreground">{{ o.estimatedDelivery }}</div>
  <div class="flex justify-end gap-1"><!-- action buttons --></div>
</div>
```

Replace with:
```vue
<TableBody>
  <TableRow v-for="o in filtered" :key="o.id"
    class="border-b border-border transition-colors hover:bg-secondary/40">
    <TableCell class="font-mono text-sm text-primary">{{ o.id }}</TableCell>
    <TableCell class="font-mono text-xs text-muted-foreground">{{ o.trackingNumber }}</TableCell>
    <TableCell class="font-mono text-sm">{{ o.customer.name }}</TableCell>
    <TableCell class="font-mono text-sm text-muted-foreground">{{ o.carrier }}</TableCell>
    <TableCell><!-- Select for status --></TableCell>
    <TableCell><!-- Select for hub or em-dash --></TableCell>
    <TableCell class="font-mono text-xs text-muted-foreground">{{ o.estimatedDelivery }}</TableCell>
    <TableCell class="text-right"><!-- action buttons --></TableCell>
  </TableRow>
</TableBody>
```

Close the Table:
```vue
  </Table>
</div>
```

The empty state (`v-if="filtered.length === 0"`) stays as a `<div>` between `</TableBody>` and `</Table>`.

- [ ] **Step 3: Adjust column widths if needed**

Check the rendered table. If column widths differ noticeably from the original design, add width classes to `TableHead` elements:
```vue
<TableHead class="w-[0.8fr]">Order ID</TableHead>
```

Run dev server: `npm run dev`, navigate to the orders page, compare layout.

- [ ] **Step 4: Verify build**

Run: `npm run build`
Workdir: `frontend/`
Expected: Build passes.

- [ ] **Step 5: Commit**

```bash
git add frontend/src/components/DeliveriesPanel.vue
git commit -m "feat: migrate DeliveriesPanel table to shadcn-vue Table"
```

---

### Task 8: Migrate Table — HubsPanel

**Files:**
- Modify: `src/components/HubsPanel.vue`

**Context:** Same CSS-grid pattern as DeliveriesPanel. The hub list has columns for ID, Name, Address, Capacity, Utilization, Status, and Actions.

- [ ] **Step 1: Replace header and rows with Table components**

Follow the same pattern as Task 7. Current grid columns: look at the grid template in the file.

- [ ] **Step 2: Verify build**

Run: `npm run build`
Workdir: `frontend/`
Expected: Build passes.

- [ ] **Step 3: Commit**

```bash
git add frontend/src/components/HubsPanel.vue
git commit -m "feat: migrate HubsPanel table to shadcn-vue Table"
```

---

### Task 9: Migrate Table — OrdersView

**Files:**
- Modify: `src/components/OrdersView.vue`

**Context:** Same CSS-grid pattern. The orders list view.

- [ ] **Step 1: Replace header and rows with Table components**

Follow the same pattern as Task 7.

- [ ] **Step 2: Verify build**

Run: `npm run build`
Workdir: `frontend/`
Expected: Build passes.

- [ ] **Step 3: Commit**

```bash
git add frontend/src/components/OrdersView.vue
git commit -m "feat: migrate OrdersView table to shadcn-vue Table"
```

---

### Task 10: Extend Badge Variants and Migrate StatusBadge

**Files:**
- Modify: `src/components/ui/Badge.vue`
- Modify: `src/components/StatusBadge.vue`

- [ ] **Step 1: Extend `badgeVariants` in Badge.vue**

Add 7 new variants matching the shipment status colors from `StatusBadge.vue`'s `statusStyles`:

```ts
export const badgeVariants = cva(
  "inline-flex items-center rounded-full border px-2.5 py-0.5 text-xs font-semibold transition-colors focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2",
  {
    variants: {
      variant: {
        default: "border-transparent bg-primary text-primary-foreground shadow hover:bg-primary/80",
        secondary: "border-transparent bg-secondary text-secondary-foreground hover:bg-secondary/80",
        destructive: "border-transparent bg-destructive text-destructive-foreground shadow hover:bg-destructive/80",
        outline: "text-foreground",
        // Shipment status variants (Ocean Deep colors)
        pending: "border-warning/30 bg-warning/15 text-warning",
        picked_up: "border-info/30 bg-info/15 text-info",
        departed: "border-accent/30 bg-accent/15 text-accent",
        in_transit: "border-primary/30 bg-primary/15 text-primary",
        out_for_delivery: "border-border bg-muted text-muted-foreground",
        delivered: "border-success/30 bg-success/15 text-success",
        delayed: "border-destructive/30 bg-destructive/15 text-destructive",
      },
    },
    defaultVariants: {
      variant: "default",
    },
  },
);
```

- [ ] **Step 2: Update StatusBadge.vue to use Badge**

Current:
```vue
<script setup lang="ts">
import type { ShipmentStatus } from "@/lib/orders";

const props = defineProps<{
  status: ShipmentStatus;
}>();

const statusLabels: Record<ShipmentStatus, string> = {
  pending: "Pending",
  picked_up: "Picked Up",
  departed: "Departed",
  in_transit: "In Transit",
  out_for_delivery: "Out for Delivery",
  delivered: "Delivered",
  delayed: "Delayed",
};

const statusStyles: Record<ShipmentStatus, string> = {
  pending: "bg-warning/15 text-warning border-warning/30",
  picked_up: "bg-info/15 text-info border-info/30",
  departed: "bg-accent/15 text-accent border-accent/30",
  in_transit: "bg-primary/15 text-primary border-primary/30",
  out_for_delivery: "bg-muted text-muted-foreground border-border",
  delivered: "bg-success/15 text-success border-success/30",
  delayed: "bg-destructive/15 text-destructive border-destructive/30",
};
</script>

<template>
  <span
    class="inline-flex items-center gap-1.5 rounded-full border px-2.5 py-0.5 font-mono text-xs uppercase tracking-wider"
    :class="statusStyles[status]"
  >
    <span class="h-1.5 w-1.5 rounded-full bg-current" />
    {{ statusLabels[status] }}
  </span>
</template>
```

Replace with:
```vue
<script setup lang="ts">
import type { ShipmentStatus } from "@/lib/orders";
import { Badge } from "@/components/ui/Badge";

const props = defineProps<{
  status: ShipmentStatus;
}>();

const statusLabels: Record<ShipmentStatus, string> = {
  pending: "Pending",
  picked_up: "Picked Up",
  departed: "Departed",
  in_transit: "In Transit",
  out_for_delivery: "Out for Delivery",
  delivered: "Delivered",
  delayed: "Delayed",
};
</script>

<template>
  <Badge :variant="status" class="gap-1.5 font-mono text-xs uppercase tracking-wider">
    <span class="h-1.5 w-1.5 rounded-full bg-current" />
    {{ statusLabels[status] }}
  </Badge>
</template>
```

- [ ] **Step 3: Verify build**

Run: `npm run build`
Workdir: `frontend/`
Expected: Build passes.

- [ ] **Step 4: Commit**

```bash
git add frontend/src/components/ui/Badge.vue frontend/src/components/StatusBadge.vue
git commit -m "feat: extend Badge variants with shipment statuses, refactor StatusBadge"
```

---

### Task 11: Replace Inline Hub Status Badges in HubsPanel

**Files:**
- Modify: `src/components/HubsPanel.vue`

**Context:** After Table migration in Task 8, the hub status badges in the table are currently inline Tailwind classes. Replace with shadcn Badge using existing variants.

- [ ] **Step 1: Replace inline badge with Badge component**

Find the hub status cell in the table. Current inline pattern:
```vue
<span class="inline-flex items-center gap-1.5 rounded-full border px-2.5 py-0.5 text-xs font-mono uppercase tracking-wider"
      :class="h.status === 'active' ? 'bg-success/15 text-success border-success/30' :
               h.status === 'maintenance' ? 'bg-warning/15 text-warning border-warning/30' :
               'bg-destructive/15 text-destructive border-destructive/30'">
  {{ h.status }}
</span>
```

Replace with:
```vue
<Badge variant="outline"
       :class="h.status === 'active' ? 'border-success/30 bg-success/15 text-success' :
               h.status === 'maintenance' ? 'border-warning/30 bg-warning/15 text-warning' :
               'border-destructive/30 bg-destructive/15 text-destructive'"
       class="font-mono text-xs uppercase tracking-wider">
  {{ h.status }}
</Badge>
```

- [ ] **Step 2: Verify build**

Run: `npm run build`
Workdir: `frontend/`
Expected: Build passes.

- [ ] **Step 3: Commit**

```bash
git add frontend/src/components/HubsPanel.vue
git commit -m "feat: replace inline hub status badges with shadcn Badge"
```

---

### Task 12: Migrate Card Patterns

**Files:** All files with manual card sections (~15 locations):
- `src/components/AnalyticsPanel.vue`
- `src/components/HomeView.vue`
- `src/components/OrdersView.vue`
- `src/components/OrderDetailView.vue`
- `src/components/CarriersView.vue`
- `src/components/OrderForm.vue`
- (and any other file with `rounded-xl border border-border bg-card`)

**Context:** Replace `div class="rounded-xl border border-border bg-card p-6 shadow-elegant"` patterns with `<Card class="shadow-elegant">` + `<CardHeader>` + `<CardTitle>` + `<CardContent>`.

- [ ] **Step 1: Migrate AnalyticsPanel.vue cards**

Replace KPI stat cards and section cards:
```vue
<!-- Before -->
<div class="rounded-xl border border-border bg-card p-6 shadow-elegant">
  <div class="flex items-center justify-between">
    <h3 class="font-mono text-xs uppercase tracking-widest text-muted-foreground">On-Time Rate</h3>
    <Truck class="h-4 w-4 text-muted-foreground" />
  </div>
  <div class="mt-1 font-mono text-3xl font-semibold text-success">{{ kpis.onTime }}%</div>
</div>

<!-- After -->
<Card class="shadow-elegant">
  <CardHeader class="flex flex-row items-center justify-between pb-2">
    <CardTitle class="font-mono text-xs uppercase tracking-widest text-muted-foreground">On-Time Rate</CardTitle>
    <Truck class="h-4 w-4 text-muted-foreground" />
  </CardHeader>
  <CardContent>
    <div class="font-mono text-3xl font-semibold text-success">{{ kpis.onTime }}%</div>
  </CardContent>
</Card>
```

Apply same pattern to all 4 KPI cards and the recent orders section card.

- [ ] **Step 2: Migrate HomeView.vue cards**

Replace the hero section and feature cards.

- [ ] **Step 3: Migrate OrdersView.vue cards**

Replace the filter/sort bar container and any section cards.

- [ ] **Step 4: Migrate OrderDetailView.vue cards**

Replace the order info panels and timeline card.

- [ ] **Step 5: Migrate CarriersView.vue cards**

Replace the tab content containers.

- [ ] **Step 6: Migrate OrderForm.vue**

Replace the form fieldset containers. Note: OrderForm uses `<fieldset>` with `<legend>` — the Card component may not be a good fit here since `<fieldset>` has semantic meaning for form grouping. Keep the fieldset but apply card styling via the existing classes. This is one location where the Card replacement is NOT appropriate.

- [ ] **Step 7: Verify build**

Run: `npm run build`
Workdir: `frontend/`
Expected: Build passes. Check each page for visual regressions.

- [ ] **Step 8: Commit**

```bash
git add frontend/src/components/AnalyticsPanel.vue frontend/src/components/HomeView.vue frontend/src/components/OrdersView.vue frontend/src/components/OrderDetailView.vue frontend/src/components/CarriersView.vue
git commit -m "feat: migrate manual card patterns to shadcn-vue Card"
```

---

### Task 13: Migrate Separator Patterns

**Files:**
- `src/components/AuthModal.vue`
- `src/components/OrderDetailView.vue`
- (any other file with `border-t border-border` or `border-b border-border` used as dividers)

- [ ] **Step 1: Replace manual separators**

Search for `border-t border-border` and `border-b border-border` patterns used as section dividers (not as actual borders on elements).

Replace:
```vue
<div class="border-t border-border pt-4" />
<!-- or -->
<hr class="border-t border-border" />
```

With:
```vue
<Separator class="my-4" />
```

Imports needed:
```ts
import { Separator } from "@/components/ui/Separator";
```

- [ ] **Step 2: Verify build**

Run: `npm run build`
Workdir: `frontend/`
Expected: Build passes.

- [ ] **Step 3: Commit**

```bash
git add <modified files>
git commit -m "feat: migrate manual separators to shadcn-vue Separator"
```

---

### Final Verification

- [ ] **Step 1: Full build check**

Run: `npm run build`
Workdir: `frontend/`
Expected: Build passes with zero errors.

- [ ] **Step 2: Lint check**

Run: `npm run lint`
Workdir: `frontend/`
Expected: No lint errors.

- [ ] **Step 3: Format**

Run: `npm run format`
Workdir: `frontend/`
Expected: Files formatted.

- [ ] **Step 4: Final commit**

```bash
git add -A
git commit -m "chore: final verification and formatting after shadcn-vue migration"
```
