# AGENTS.md — ship-simple

## Project Overview

Vue 3 app deployed on Cloudflare Workers. Shipment tracking dashboard using Tailwind CSS v4, shadcn-vue (New York), TanStack Vue Query, Pinia, and Vue Router.

**Package manager:** Bun (`bun.lock`, `bunfig.toml`). Use `bun add` / `bun remove` / `bun install`.

---

## Build / Lint / Test Commands

| Command           | Purpose                        |
| ----------------- | ------------------------------ |
| `npm run dev`     | Start dev server (Vite)        |
| `npm run build`   | vue-tsc typecheck + Vite build |
| `npm run preview` | Preview production build       |
| `npm run lint`    | ESLint check (entire project)  |
| `npm run format`  | Prettier auto-format           |

**No test framework is currently installed.** Do not assume vitest, jest, or playwright exist. If tests need to be added, use `bun add -d vitest` with `@vue/test-utils` and set up `vitest.config.ts` alongside the existing `vite.config.ts`.

ESLint config: `eslint.config.js` (flat config, typescript-eslint + eslint-plugin-vue + Prettier).
Prettier config: `.prettierrc` (printWidth 100, semi, double quotes, trailingComma all).

---

## Code Style Guidelines

### Vue SFC Patterns

- Always use `<script setup lang="ts">` (Composition API, no `setup()` function).
- Props defined with `defineProps<{ propName: Type }>()` — use interfaces or inline type.
- Local state with `ref()`, `computed()`, `watch()` from `vue`.
- Lazy-load heavy components via `defineAsyncComponent(() => import('...'))`.
- Named slots via `<slot />` with optional `#slotName` syntax in consumers.
- Always use `import type { ... }` for type-only imports.
- Do NOT use Options API (`export default { ... }`) in new code.
- Keep `<script>` first, then `<template>`, then `<style scoped>` (though scoped styles are not used in this project — all styling is via Tailwind classes).

### Imports

- Use `@/` path alias for all project imports (maps to `./src/` via vite alias & tsconfig paths).
- Group imports: Vue / external packages first, then `@/` project imports, then relative imports.
- Import icons from `lucide-vue-next` directly by name (tree-shaken).
- Use `import type { ... }` for type-only imports.

### Formatting

- Semicolons required.
- Double quotes (single quotes for strings that contain double quotes).
- Trailing commas everywhere.
- Print width: 100.
- Format with: `npm run format` (Prettier).

### Types

- Define shared types in `src/lib/*.ts` files using `export type` or `export interface`.
- Co-locate component-local types within the component's `<script setup>` block.
- Use `Record<UnionType, ValueType>` for lookup maps (e.g., `statusLabels` in `orders.ts`).
- Prefer `interface` for object shapes, `type` for unions and utility types.
- Use `as const` sparingly; prefer union types.

### Naming Conventions

- **Components:** PascalCase, multi-word names, file name matches component name (e.g., `StatusBadge.vue` → `<StatusBadge />`).
- **Views:** PascalCase, suffixed with `View` (e.g., `HomeView.vue`, `OrdersView.vue`).
- **Lib/utils:** camelCase for functions, PascalCase for types/interfaces.
- **Route names:** kebab-case (`order-detail`).
- **CSS classes:** Tailwind utility classes only (no CSS modules, no `<style scoped>`).
- **Constants:** UPPER_SNAKE_CASE for module-level constants (e.g., `FILTERS`, `MOBILE_BREAKPOINT`).
- **Directory names:** lowercase, plural for collections (`components/`, `views/`, `lib/`, `router/`).

### Component Patterns

- Use `defineProps<Props>()` with dedicated `Props` interface.
- Destructure reactive state in `<script setup>` — no `.value` needed in `<template>`.
- Use `computed()` for derived values (e.g., `filtered` in `OrdersView.vue`).
- Optional `class` prop typed as `HTMLAttributes['class']` for style customization via `cn()`.
- Use `RouterLink` for internal navigation, `<component :is="..." />` for dynamic icons.

### State Management

- **Server/cache state:** TanStack Vue Query (`VueQueryPlugin` in `main.ts`).
- **Client state:** Pinia stores (`createPinia()` in `main.ts`).
- **Local component state:** `ref()`, `computed()`, `watch()` from Vue 3.
- **No Zustand, Redux, or raw Context** — Pinia handles client state; Vue Query handles async/server state.

### Routing

- Vue Router with `createWebHistory()`.
- Routes defined in `src/router/index.ts` using `{ path, name, component }` shape.
- All route components are lazy-loaded via `() => import('@/views/...')`.
- Route params accessed via `useRoute().params`.
- Navigation with `RouterLink` component or `useRouter().push({ name: '...', params: {...} })`.
- Route names: `'home'`, `'orders'`, `'order-detail'`.

### Error Handling

- **Not found:** Manual 404 rendering (see `OrderDetailView.vue` — render fallback when data not found).
- **Async component loading:** `<Suspense>` with `#fallback` template for lazy-loaded components (see map in OrderDetailView).
- **No global error boundary** — errors propagate to console. Consider adding `app.config.errorHandler` if needed.
- **Build errors:** `vue-tsc` catches type errors during `npm run build`.

### CSS & Styling

- Tailwind CSS v4 exclusively (no CSS-in-JS, no `<style scoped>`).
- Utility: `cn()` from `@/lib/utils` (clsx + tailwind-merge).
- shadcn-vue theme tokens: `bg-background`, `text-foreground`, `text-muted-foreground`, `border-border`, `bg-card`, `bg-secondary`, `bg-primary`, `bg-destructive`, `bg-success`, `bg-warning`, `bg-info`.
- Custom CSS utilities: `bg-gradient-hero`, `bg-gradient-accent`, `shadow-elegant`, `shadow-glow` (defined in `styles.css` `@layer utilities`).
- Font: `font-mono` ("JetBrains Mono") for data-heavy text, `font-sans` ("Work Sans") for prose.
- All shadcn-vue components under `src/components/ui/` — do not edit these directly; use `bunx shadcn-vue@latest add <component>`.

### Icons

- Use `lucide-vue-next` for all icons. Import icon components by name: `import { Search, ArrowRight } from 'lucide-vue-next'`.
- Render dynamically via `<component :is="iconComponent" />` when icon is chosen at runtime.

### File Structure

```
src/
├── components/       # Shared Vue components
│   └── ui/           # shadcn-vue primitives (auto-generated)
├── views/            # Page-level views (consumed by router)
├── lib/              # Utilities, types, data, error handling
├── router/           # Vue Router setup
├── App.vue           # Root component
├── main.ts           # App entry point (creates app, installs plugins)
└── styles.css        # Tailwind entry point + theme tokens
```

### Feature Conventions

- **No tests currently exist.** When adding tests, use `vitest` with `@vue/test-utils`.
- **No i18n** — hardcoded English strings throughout.
- **Leaflet maps** with `leaflet` + custom `L.divIcon` — rendered client-side inside `<Suspense>` via `defineAsyncComponent`, guarded with `onMounted` flag.
- **In-memory data** — all order data lives in `src/lib/orders.ts` (no API calls). Future API integration should go through TanStack Vue Query.
