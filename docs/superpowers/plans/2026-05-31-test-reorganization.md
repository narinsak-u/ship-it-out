# Frontend Test File Reorganization — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Move all 28 Vitest `.spec.ts` files from `src/` into `tests/unit/` mirroring the source tree structure, with rewritten imports and updated config.

**Architecture:** Files move from co-located positions in `src/` into `tests/unit/<mirrored-path>/`. Every `./` relative import to the source-under-test gets rewritten to use the `@/` alias. No runtime behavior changes — purely structural.

**Tech Stack:** Vitest, TypeScript, Vite, `@/` path alias

---

### Task 1: Move `src/views/*.spec.ts` — 5 files

**Files:**
- Move: `src/views/OrdersView.spec.ts` → `tests/unit/views/OrdersView.spec.ts`
- Move: `src/views/OrderFormView.spec.ts` → `tests/unit/views/OrderFormView.spec.ts`
- Move: `src/views/OrderDetailView.spec.ts` → `tests/unit/views/OrderDetailView.spec.ts`
- Move: `src/views/HomeView.spec.ts` → `tests/unit/views/HomeView.spec.ts`
- Move: `src/views/CarriersView.spec.ts` → `tests/unit/views/CarriersView.spec.ts`

- [ ] **Step 1: Create target directory**

```bash
mkdir -p tests/unit/views
```

- [ ] **Step 2: Move first file and rewrite import**

Read `src/views/OrdersView.spec.ts`, then write to `tests/unit/views/OrdersView.spec.ts` with `"./OrdersView.vue"` → `"@/views/OrdersView.vue"` in both the `vi.mock` call and the dynamic `import()` on the line referencing `OrdersView.vue`.

- [ ] **Step 3: Move remaining 4 view spec files**

Apply the same transform to each: replace `"./<ViewName>.vue"` with `"@/views/<ViewName>.vue"`.

- `OrderFormView.spec.ts`: `"./OrderFormView.vue"` → `"@/views/OrderFormView.vue"`
- `OrderDetailView.spec.ts`: `"./OrderDetailView.vue"` → `"@/views/OrderDetailView.vue"`
- `HomeView.spec.ts`: `"./HomeView.vue"` → `"@/views/HomeView.vue"`
- `CarriersView.spec.ts`: `"./CarriersView.vue"` → `"@/views/CarriersView.vue"`

- [ ] **Step 4: Remove old files**

```bash
git rm src/views/OrdersView.spec.ts src/views/OrderFormView.spec.ts src/views/OrderDetailView.spec.ts src/views/HomeView.spec.ts src/views/CarriersView.spec.ts
```

- [ ] **Step 5: Run tests to verify views pass**

```bash
npx vitest run tests/unit/views
```

Expected: 5 test files found, all passing.

- [ ] **Step 6: Commit**

```bash
git add tests/unit/views && git commit -m "test: move view spec files to tests/unit/views"
```

---

### Task 2: Move `src/components/*.spec.ts` — 8 files

**Files:**
- Move: `src/components/ThaiAddressGroup.spec.ts` → `tests/unit/components/ThaiAddressGroup.spec.ts`
- Move: `src/components/StatusBadge.spec.ts` → `tests/unit/components/StatusBadge.spec.ts`
- Move: `src/components/SiteHeader.spec.ts` → `tests/unit/components/SiteHeader.spec.ts`
- Move: `src/components/SiteFooter.spec.ts` → `tests/unit/components/SiteFooter.spec.ts`
- Move: `src/components/Pagination.spec.ts` → `tests/unit/components/Pagination.spec.ts`
- Move: `src/components/OrderForm.spec.ts` → `tests/unit/components/OrderForm.spec.ts`
- Move: `src/components/ConfirmDialog.spec.ts` → `tests/unit/components/ConfirmDialog.spec.ts`
- Move: `src/components/AuthModal.spec.ts` → `tests/unit/components/AuthModal.spec.ts`

Transform pattern for each: replace `"./<ComponentName>.vue"` with `"@/components/<ComponentName>.vue"`.

- [ ] **Step 1: Create target directory**

```bash
mkdir -p tests/unit/components
```

- [ ] **Step 2: Move all 8 files with import rewrite**

For each file, read the source, replace the `./` component import with `@/components/` equivalent, write to `tests/unit/components/`.

Then remove originals:

```bash
git rm src/components/ThaiAddressGroup.spec.ts src/components/StatusBadge.spec.ts src/components/SiteHeader.spec.ts src/components/SiteFooter.spec.ts src/components/Pagination.spec.ts src/components/OrderForm.spec.ts src/components/ConfirmDialog.spec.ts src/components/AuthModal.spec.ts
```

- [ ] **Step 3: Run tests to verify components pass**

```bash
npx vitest run tests/unit/components
```

Expected: 8 test files found, all passing.

- [ ] **Step 4: Commit**

```bash
git add tests/unit/components && git commit -m "test: move component spec files to tests/unit/components"
```

---

### Task 3: Move `src/stores/` + `src/composables/` specs — 3 files

**Files:**
- Move: `src/stores/auth.spec.ts` → `tests/unit/stores/auth.spec.ts`
- Move: `src/composables/useSearchFilter.spec.ts` → `tests/unit/composables/useSearchFilter.spec.ts`
- Move: `src/composables/usePagination.spec.ts` → `tests/unit/composables/usePagination.spec.ts`

Transform:
- `stores/auth.spec.ts`: `"./auth"` → `"@/stores/auth"`
- `composables/*.spec.ts`: `"./useXxx"` → `"@/composables/useXxx"`

- [ ] **Step 1: Create target directories**

```bash
mkdir -p tests/unit/stores tests/unit/composables
```

- [ ] **Step 2: Move and rewrite all 3 files**

- [ ] **Step 3: Remove originals**

```bash
git rm src/stores/auth.spec.ts src/composables/useSearchFilter.spec.ts src/composables/usePagination.spec.ts
```

- [ ] **Step 4: Run tests**

```bash
npx vitest run tests/unit/stores tests/unit/composables
```

Expected: 3 test files found, all passing.

- [ ] **Step 5: Commit**

```bash
git add tests/unit/stores tests/unit/composables && git commit -m "test: move store and composable spec files to tests/unit"
```

---

### Task 4: Move `src/hooks/*.spec.ts` — 3 files

**Files:**
- Move: `src/hooks/useOrders.spec.ts` → `tests/unit/hooks/useOrders.spec.ts`
- Move: `src/hooks/useHubs.spec.ts` → `tests/unit/hooks/useHubs.spec.ts`
- Move: `src/hooks/useDeliveries.spec.ts` → `tests/unit/hooks/useDeliveries.spec.ts`

Transform: `"./useXxx"` → `"@/hooks/useXxx"`

- [ ] **Step 1: Create target directory**

```bash
mkdir -p tests/unit/hooks
```

- [ ] **Step 2: Move and rewrite all 3 files**

- [ ] **Step 3: Remove originals**

```bash
git rm src/hooks/useOrders.spec.ts src/hooks/useHubs.spec.ts src/hooks/useDeliveries.spec.ts
```

- [ ] **Step 4: Run tests**

```bash
npx vitest run tests/unit/hooks
```

Expected: 3 test files found, all passing.

- [ ] **Step 5: Commit**

```bash
git add tests/unit/hooks && git commit -m "test: move hook spec files to tests/unit/hooks"
```

---

### Task 5: Move `src/lib/*.spec.ts` — 3 files

**Files:**
- Move: `src/lib/utils.spec.ts` → `tests/unit/lib/utils.spec.ts`
- Move: `src/lib/geocode.spec.ts` → `tests/unit/lib/geocode.spec.ts`
- Move: `src/lib/analytics-utils.spec.ts` → `tests/unit/lib/analytics-utils.spec.ts`

Transform: `"./xxx"` → `"@/lib/xxx"`

- [ ] **Step 1: Create target directory**

```bash
mkdir -p tests/unit/lib
```

- [ ] **Step 2: Move and rewrite all 3 files**

- [ ] **Step 3: Remove originals**

```bash
git rm src/lib/utils.spec.ts src/lib/geocode.spec.ts src/lib/analytics-utils.spec.ts
```

- [ ] **Step 4: Run tests**

```bash
npx vitest run tests/unit/lib
```

Expected: 3 test files found, all passing.

- [ ] **Step 5: Commit**

```bash
git add tests/unit/lib && git commit -m "test: move lib spec files to tests/unit/lib"
```

---

### Task 6: Move `src/lib/api/*.spec.ts` — 7 files

**Files:**
- Move: `src/lib/api/tracking.spec.ts` → `tests/unit/lib/api/tracking.spec.ts`
- Move: `src/lib/api/queryKeys.spec.ts` → `tests/unit/lib/api/queryKeys.spec.ts`
- Move: `src/lib/api/orders.spec.ts` → `tests/unit/lib/api/orders.spec.ts`
- Move: `src/lib/api/mappers.spec.ts` → `tests/unit/lib/api/mappers.spec.ts`
- Move: `src/lib/api/hubs.spec.ts` → `tests/unit/lib/api/hubs.spec.ts`
- Move: `src/lib/api/analytics.spec.ts` → `tests/unit/lib/api/analytics.spec.ts`
- Move: `src/lib/api/client.spec.ts` → `tests/unit/lib/api/client.spec.ts`

Transform for most files: `"./xxx"` → `"@/lib/api/xxx"`

Special case for `client.spec.ts`: also rewrite `"../../../tests/msw/server"` → `"../../../msw/server"`

- [ ] **Step 1: Create target directory**

```bash
mkdir -p tests/unit/lib/api
```

- [ ] **Step 2: Move 6 standard files with rewrite**

For each of `tracking`, `queryKeys`, `orders`, `mappers`, `hubs`, `analytics`: replace `"./xxx"` with `"@/lib/api/xxx"`.

- [ ] **Step 3: Move `client.spec.ts` with dual rewrite**

Two changes:
- `"./client"` → `"@/lib/api/client"`
- `"../../../tests/msw/server"` → `"../../../msw/server"`

- [ ] **Step 4: Remove originals**

```bash
git rm src/lib/api/tracking.spec.ts src/lib/api/queryKeys.spec.ts src/lib/api/orders.spec.ts src/lib/api/mappers.spec.ts src/lib/api/hubs.spec.ts src/lib/api/analytics.spec.ts src/lib/api/client.spec.ts
```

- [ ] **Step 5: Run tests**

```bash
npx vitest run tests/unit/lib/api
```

Expected: 7 test files found, all passing.

- [ ] **Step 6: Commit**

```bash
git add tests/unit/lib/api && git commit -m "test: move lib/api spec files to tests/unit/lib/api"
```

---

### Task 7: Update `vite.config.ts`

**Files:**
- Modify: `frontend/vite.config.ts:40`

- [ ] **Step 1: Update Vitest include pattern**

Change line 40 from:
```ts
include: ["src/**/*.spec.ts"],
```
to:
```ts
include: ["tests/unit/**/*.spec.ts"],
```

- [ ] **Step 2: Run full test suite**

```bash
npx vitest run
```

Expected: 28 test files found from `tests/unit/`, all passing.

- [ ] **Step 3: Commit**

```bash
git add vite.config.ts && git commit -m "test: update vitest include path to tests/unit"
```

---

### Task 8: Update `tsconfig.json`

**Files:**
- Modify: `frontend/tsconfig.json:19`

- [ ] **Step 1: Add tests to include**

Change line 19 from:
```json
"include": ["src/**/*.ts", "src/**/*.vue", "vite.config.ts"]
```
to:
```json
"include": ["src/**/*.ts", "src/**/*.vue", "tests/unit/**/*.ts", "vite.config.ts"]
```

- [ ] **Step 2: Run build to verify type-checking**

```bash
npm run build
```

Expected: `vue-tsc` and `vite build` both succeed.

- [ ] **Step 3: Commit**

```bash
git add tsconfig.json && git commit -m "test: include tests/unit in tsconfig for type-checking"
```

---

### Task 9: Final verification

- [ ] **Step 1: Run all checks**

```bash
npm run test
npm run build
npm run lint
```

Expected: 28 Vitest tests pass, build succeeds, lint is clean.

- [ ] **Step 2: Verify no spec files remain in src/**

```bash
find src -name "*.spec.ts"
```

Expected: no output (no `.spec.ts` files left in `src/`).
