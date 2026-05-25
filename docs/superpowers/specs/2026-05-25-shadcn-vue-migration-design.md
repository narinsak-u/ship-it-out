# Shadcn-vue UI Component Migration

**Date:** 2026-05-25
**Status:** Draft

## Goal

Replace all custom UI patterns in the frontend with shadcn-vue component equivalents while preserving the existing Ocean Deep dark theme exactly — no visual changes.

## Scope

6 work areas, executed component-by-component (each independently verifiable with `npm run build`):

1. Install missing shadcn-vue components (`dialog`, `select`)
2. Replace custom modals with Dialog
3. Replace native `<select>` elements with Select
4. Refactor CSS-grid tables to Table
5. Extend Badge variants and use everywhere
6. Replace manual card patterns and separators

## Non-Goals

- No visual changes to the theme
- No changes to `styles.css` (Ocean Deep dark theme already maps correctly)
- No changes to `StatusBadge.vue`'s colored dot indicator
- No changes to existing Button, Input, or Skeleton usage (already using shadcn-vue)
- No changes to logic, data flow, or component APIs

## New Components to Install

```bash
bunx shadcn-vue@latest add dialog
bunx shadcn-vue@latest add select
```

`dialog` wraps radix-vue's Dialog primitive (DialogRoot, DialogTrigger, DialogPortal, DialogContent, DialogHeader, DialogTitle, etc.).
`select` wraps radix-vue's Select primitive (SelectRoot, SelectTrigger, SelectValue, SelectContent, SelectItem, etc.).

These will be installed to `src/components/ui/dialog/` and `src/components/ui/select/` with the project's existing new-york style, neutral base, and Ocean Deep tokens.

## Component-by-Component Migration

### 1. Dialog — Replace Custom Modals

**Files:**
- `src/components/AuthModal.vue`
- `src/components/HubFormModal.vue`
- `src/components/SiteHeader.vue` (trigger for AuthModal)
- `src/components/HubsPanel.vue` (trigger for HubFormModal)

**Current pattern (both files):**
```vue
<div class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm"
     @click.self="emit('close')">
  <div class="w-full max-w-md rounded-xl border border-border bg-card p-6 shadow-elegant">
    <!-- form content -->
  </div>
</div>
```

**Target pattern:**
```vue
<DialogRoot v-model:open="open">
  <DialogTrigger>...</DialogTrigger>
  <DialogPortal>
    <DialogOverlay />
    <DialogContent>
      <DialogHeader>
        <DialogTitle>Sign In</DialogTitle>
      </DialogHeader>
      <!-- form content -->
    </DialogContent>
  </DialogPortal>
</DialogRoot>
```

**Details:**
- `AuthModal.vue`: The existing tabbed sign-in/sign-up forms stay. Replace the backdrop wrapper with Dialog components. The modal's `open` state is controlled by `showAuthModal` ref in `SiteHeader.vue` — use `DialogRoot v-model:open` to bind it. Remove the custom `@click.self` handler (Dialog provides escape-key + overlay-click dismiss).
- `HubFormModal.vue`: Same pattern. The "Add Hub" button in `HubsPanel.vue` becomes a `DialogTrigger`. The form content lives inside `DialogContent`. Remove custom backdrop.
- No CSS changes needed — Dialog uses `--color-background`, `--color-border`, `--color-foreground` tokens matching the theme.

### 2. Select — Replace Native `<select>` Elements

**Files:**
- `src/components/DeliveriesPanel.vue` (2 selects: status, hub per row)
- `src/components/ThaiAddressGroup.vue` (1 select: sub-district)
- `src/components/HubFormModal.vue` (1 select: hub status)

**Current pattern:**
```vue
<select v-model="value"
        class="rounded-lg border border-border bg-background px-2 py-1 font-mono text-xs">
  <option value="">Select...</option>
  <option v-for="opt in options" :key="opt" :value="opt">{{ opt }}</option>
</select>
```

**Target pattern:**
```vue
<SelectRoot v-model="value">
  <SelectTrigger class="font-mono text-xs">
    <SelectValue placeholder="Select..." />
  </SelectTrigger>
  <SelectContent>
    <SelectItem v-for="opt in options" :key="opt" :value="opt">
      {{ opt }}
    </SelectItem>
  </SelectContent>
</SelectRoot>
```

**Details:**
- `DeliveriesPanel.vue`: The `draftStatus`/`draftHubId` reactive state and `@change` handlers stay identical. Only the template markup changes from `<select>` to Select components. The `canUpdate` computed and `handleUpdate` function remain untouched.
- `ThaiAddressGroup.vue`: The `onSubDistrictSelected` handler and `lookupStatus` logic stay. The disabled state (`!modelValue.zipcode`) maps to `:disabled` on `SelectTrigger`.
- `HubFormModal.vue`: The `status` v-model stays. Replace the native `<select>` in the form body.
- Select components automatically handle keyboard navigation, ARIA attributes, and focus management.

### 3. Table — Refactor CSS-Grid Tables to Native `<table>`

**Files:**
- `src/components/OrdersView.vue`
- `src/components/DeliveriesPanel.vue`
- `src/components/HubsPanel.vue`

**Current pattern (CSS grid):**
```vue
<div class="hidden grid-cols-[0.8fr_1fr_1fr_0.8fr_0.5fr] gap-4 border-b border-border bg-secondary/50 px-6 py-3 font-mono text-[11px] uppercase tracking-widest text-muted-foreground md:grid">
  <span>Order ID</span>
  <span>Tracking</span>
  <!-- ... -->
</div>
<div v-for="row in data" :key="row.id"
     class="group grid grid-cols-[0.8fr_1fr_1fr_0.8fr_0.5fr] gap-4 border-b border-border px-6 py-4 md:items-center">
  <div>...</div>
  <div>...</div>
  <!-- ... -->
</div>
```

**Target pattern:**
```vue
<Table>
  <TableHeader>
    <TableRow>
      <TableHead>Order ID</TableHead>
      <TableHead>Tracking</TableHead>
      <!-- ... -->
    </TableRow>
  </TableHeader>
  <TableBody>
    <TableRow v-for="row in data" :key="row.id">
      <TableCell>...</TableCell>
      <TableCell>...</TableCell>
      <!-- ... -->
    </TableRow>
  </TableBody>
</Table>
```

**Details:**
- The CSS grid column templates (`grid-cols-[0.8fr_1fr_...]`) are replaced by the Table component's default layout. Column widths may shift slightly — use `<TableHead class="w-[xx]">` if specific widths are needed to match the current appearance. Inspect each table after migration and adjust widths per-column only where noticeably different.
- The `hidden md:grid` responsive pattern (mobile: stacked cards, desktop: table) in the current CSS-grid tables should be preserved. The Table component renders as a native `<table>`, which naturally stacks on small screens when wrapped in a horizontally scrollable container. Add `overflow-x-auto` wrapper on mobile if needed.
- Existing sorting, filtering, and row-action logic stays unchanged.
- shadcn-vue's Table styling uses the same `--color-border`, `--color-secondary`, `--color-muted-foreground` tokens already in use.
- **`DeliveriesPanel.vue`** has the most complex table with inline `<select>` elements inside cells — these will have already been migrated to Select in step 2, so the Table step uses the Select components naturally.

### 4. Badge — Extend Variants and Use Everywhere

**File:**
- `src/components/ui/Badge.vue` (extend variants)
- `src/components/StatusBadge.vue` (refactor to use Badge)
- `src/components/HubsPanel.vue` (replace inline badge classes)

**Current shadcn Badge variants:**
```
default, secondary, destructive, outline
```

**Extended variants** (add to `badgeVariants` in `Badge.vue`):
```
pending:       "bg-warning/15 text-warning border-warning/30"
picked_up:     "bg-info/15 text-info border-info/30"
departed:      "bg-accent/15 text-accent border-accent/30"
in_transit:    "bg-primary/15 text-primary border-primary/30"
out_for_delivery: "bg-muted text-muted-foreground border-border"
delivered:     "bg-success/15 text-success border-success/30"
delayed:       "bg-destructive/15 text-destructive border-destructive/30"
```

The color values are copied directly from `StatusBadge.vue`'s `statusStyles` to ensure zero visual change.

**StatusBadge.vue refactor:**
```vue
<template>
  <Badge :variant="status" class="gap-1.5">
    <span class="h-1.5 w-1.5 rounded-full bg-current" />
    {{ label }}
  </Badge>
</template>
```
The `font-mono`, `uppercase`, `tracking-wider` classes are moved into the badge variant entries. The colored dot stays.

**HubsPanel.vue:** Replace inline:
```vue
<span class="inline-flex items-center gap-1.5 rounded-full border px-2.5 py-0.5 text-xs font-mono uppercase tracking-wider"
      :class="h.status === 'active' ? 'bg-success/15 text-success border-success/30' : ...">
```
with:
```vue
<Badge variant="outline" :class="h.status === 'active' ? 'text-success border-success/30 bg-success/15' : ...">
```

### 5. Card — Replace Manual Card Patterns

**Files (representative — ~15 locations):**
- `src/components/AnalyticsPanel.vue`
- `src/components/OrdersView.vue`
- `src/components/OrderDetailView.vue`
- `src/components/HomeView.vue`
- `src/components/CarriersView.vue`
- `src/components/AuthModal.vue` (after Dialog migration)
- `src/components/HubFormModal.vue` (after Dialog migration)
- `src/components/OrderForm.vue`

**Current pattern:**
```vue
<div class="rounded-xl border border-border bg-card p-6 shadow-elegant">
  <h3 class="font-mono text-xs uppercase tracking-widest text-muted-foreground">Title</h3>
  <!-- content -->
</div>
```

**Target patterns:**

For content sections:
```vue
<Card class="shadow-elegant">
  <CardHeader>
    <CardTitle>Title</CardTitle>
  </CardHeader>
  <CardContent>
    <!-- content -->
  </CardContent>
</Card>
```

For KPI stat cards (AnalyticsPanel, HubsPanel):
```vue
<Card>
  <CardHeader class="flex flex-row items-center justify-between pb-2">
    <CardTitle class="text-sm font-medium text-muted-foreground">Label</CardTitle>
    <component :is="icon" class="h-4 w-4 text-muted-foreground" />
  </CardHeader>
  <CardContent>
    <div class="text-3xl font-bold text-success">{{ value }}</div>
  </CardContent>
</Card>
```

**Visual preservation:** shadcn Card uses `bg-card text-card-foreground` with `border-border` — same tokens already used manually. The `shadow-elegant` custom utility is passed via the `class` prop to Card.

### 6. Separator — Replace Manual Border Dividers

**Files:** `AuthModal.vue` (after Dialog migration), `OrderDetailView.vue`, and any file with `border-t border-border` or `border-b border-border` used as section dividers.

**Current pattern:**
```vue
<div class="border-t border-border pt-4" />
```

**Target pattern:**
```vue
<Separator class="my-4" />
```

## Verification

After each component migration step:
1. Run `npm run build` — must pass with zero errors
2. Manually verify the page renders with correct colors (theme tokens unchanged)
3. Check interactive behavior (modals open/close, selects work, table rows clickable)

Full verification checklist before final commit:
- [ ] `npm run build` passes
- [ ] `npm run lint` passes
- [ ] Modals open/close with proper focus trapping
- [ ] All select dropdowns functional
- [ ] Tables render with correct data
- [ ] Badge colors match pre-migration appearance
- [ ] Card layouts unchanged
- [ ] No regressions in any view

## Out of Scope

- `StatusBadge.vue`'s colored dot indicator (stays as-is)
- `SiteHeader.vue` navigation buttons (simple `<button>` is fine)
- `ShipmentMap.vue` (Leaflet-only, no UI components to replace)
- `SiteFooter.vue` (minimal HTML, no components needed)
