# Design: Frontend Test File Reorganization

## Goal

Move all 28 Vitest `.spec.ts` files from co-located positions within `src/` into a dedicated `tests/unit/` directory that mirrors the `src/` directory structure, with updated imports and config so everything continues to work identically.

## Target Structure

```
src/                          tests/unit/
├── views/*.vue               ├── views/*.spec.ts     (5 files)
├── components/*.vue          ├── components/*.spec.ts (7 files)
├── stores/auth.ts            ├── stores/auth.spec.ts
├── hooks/*.ts                ├── hooks/*.spec.ts      (3 files)
├── composables/*.ts          ├── composables/*.spec.ts (2 files)
├── lib/*.ts                  ├── lib/*.spec.ts         (3 files)
├── lib/api/*.ts              ├── lib/api/*.spec.ts     (7 files)
```

Existing `tests/e2e/`, `tests/msw/`, and `tests/setup.ts` remain unchanged.

## Import Rewriting

Every spec file currently uses `./` relative imports to reference its source-under-test (e.g., `"./OrdersView.vue"`). These are rewritten to use the `@/` alias:

| Current | Rewritten |
|---|---|
| `"./OrdersView.vue"` | `"@/views/OrdersView.vue"` |
| `"./auth"` | `"@/stores/auth"` |
| `"./utils"` | `"@/lib/utils"` |
| `"./useOrders"` | `"@/hooks/useOrders"` |
| `"./StatusBadge.vue"` | `"@/components/StatusBadge.vue"` |

**Special case:** `src/lib/api/client.spec.ts` imports `"../../../tests/msw/server"`. After moving to `tests/unit/lib/api/client.spec.ts`, rewritten to `"../../../msw/server"`.

## Config Changes

**`vite.config.ts` (test section):**
- `include: ["src/**/*.spec.ts"]` → `include: ["tests/unit/**/*.spec.ts"]`
- `setupFiles: ["tests/setup.ts"]` — unchanged (already correct)

**`tsconfig.json`:**
- Add `"tests/unit/**/*.ts"` to `include` so `vue-tsc` type-checks test files during `npm run build`

## Verification

1. `npm run test` — all 28 Vitest tests pass
2. `npm run build` — `vue-tsc` + `vite build` succeed (including new test type-checking)
3. `npm run lint` — no new ESLint issues
