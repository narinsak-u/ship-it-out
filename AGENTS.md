# AGENTS.md — ship-simple

## Project Overview

TanStack Start (React 19) app deployed on Cloudflare Workers. Shipment tracking dashboard using Tailwind CSS v4, shadcn/ui (New York), TanStack React Query, and TanStack Router.

**Package manager:** Bun (`bun.lock`, `bunfig.toml`). Use `bun add` / `bun remove` / `bun install`.

---

## Build / Lint / Test Commands

| Command             | Purpose                       |
| ------------------- | ----------------------------- |
| `npm run dev`       | Start dev server (Vite)       |
| `npm run build`     | Production build (Vite)       |
| `npm run build:dev` | Dev-mode build                |
| `npm run preview`   | Preview production build      |
| `npm run lint`      | ESLint check (entire project) |
| `npm run format`    | Prettier auto-format          |

**No test framework is currently installed.** Do not assume vitest, jest, or playwright exist. If tests need to be added, use `bun add -d vitest` and set up `vitest.config.ts` alongside the existing `vite.config.ts`.

ESLint config: `eslint.config.js` (flat config, typescript-eslint + react-hooks + Prettier).
Prettier config: `.prettierrc` (printWidth 100, semi, double quotes, trailingComma all).

---

## Code Style Guidelines

### Imports

- Use `@/` path alias for all project imports (maps to `./src/` via tsconfig paths & vite-tsconfig-paths).
- Group imports: external packages first, then `@/` project imports.
- Do NOT import `server-only` (blocked by ESLint rule; TanStack Start uses `*.server.ts` convention instead).
- Use `import type { ... }` for type-only imports.
- Lazy-load heavy components (e.g., map) via `React.lazy(() => import(...))` with `Suspense`.

### Formatting

- Semicolons required.
- Double quotes (single quotes for strings that contain double quotes).
- Trailing commas everywhere.
- Print width: 100.
- Format with: `npm run format` (Prettier).

### Types

- Define shared types in `src/lib/*.ts` files using `export type` or `export interface`.
- Co-locate component-local types within the component file.
- Use `Record<UnionType, ValueType>` for lookup maps (e.g., `statusLabels`).
- Prefer `interface` for object shapes, `type` for unions and utility types.
- Use `as const` sparingly; prefer union types.

### Naming Conventions

- **Components:** PascalCase, named exports (`export function ComponentName()`), file name matches component name (e.g., `StatusBadge.tsx` → `StatusBadge`).
- **Hooks:** camelCase, `use*` prefix, file name matches hook name (e.g., `use-mobile.tsx` → `useIsMobile`).
- **Lib/utils:** camelCase for functions, PascalCase for types/interfaces.
- **Route files:** kebab-case, TanStack Router file conventions (`orders.$orderId.tsx`, `orders.index.tsx`).
- **CSS classes:** Tailwind utility classes only (no CSS modules, no styled-components).
- **Constants:** UPPER_SNAKE_CASE for module-level constants (e.g., `MOBILE_BREAKPOINT`, `TTL_MS`, `FILTERS`).

### Component Patterns

- Always use `function ComponentName()` syntax (not arrow functions).
- Always use named exports (no `export default`).
- Route component defined as `export const Route = createFileRoute(...)({ component: PageComponent });`.
- Route data loading via `loader` function, accessed with `Route.useLoaderData()`.
- Route params via `route.useParams()` or `loader` params argument.
- Destructure props inline, optionally adding `className?: string`.

### react-hook-form / Zod

- Form schemas defined with Zod (`z.object({...})`).
- Use `@hookform/resolvers` with `zodResolver` for validation.
- Components use `useForm` and shadcn form components.

### Error Handling

- **Route errors:** Use `notFound()` from `@tanstack/react-router` with `notFoundComponent`.
- **Render errors:** `errorComponent` on routes catches rendering exceptions.
- **SSR errors:** Captured in `error-capture.ts` (global `error`/`unhandledrejection` listeners), recovered in `server.ts` when h3 swallows throws.
- **Middleware errors:** `start.ts` wraps all requests in error middleware that renders a styled error page.
- **No try/catch for h3-thrown errors** (they're swallowed as normal 500 responses); handled via `normalizeCatastrophicSsrResponse`.
- Client errors: `console.error` (no toast framework integrated).

### State Management

- **Server state:** TanStack React Query (`QueryClient`, `QueryClientProvider`).
- **Local state:** `useState`, `useMemo`, `useEffect` from React 19.
- **No Zustand, Redux, or Context for global state** — React Query handles server state; client state is local.
- **Router context** carries `queryClient` via `createRootRouteWithContext`.

### CSS & Styling

- Tailwind CSS v4 exclusively (no CSS-in-JS, no CSS modules).
- Utility: `cn()` from `@/lib/utils` (clsx + tailwind-merge).
- shadcn/ui theme tokens: `bg-background`, `text-foreground`, `text-muted-foreground`, `border-border`, `bg-card`, `bg-secondary`, `bg-primary`, `bg-destructive`, `bg-success`, `bg-info`.
- Custom gradients: `bg-gradient-hero`, `bg-gradient-accent`.
- Shadows: `shadow-elegant`, `shadow-glow`.
- Font: `font-mono` for data-heavy text, `font-sans` for prose.
- All shadcn/ui components under `src/components/ui/` — do not edit these directly; use `bunx shadcn@latest add <component>`.

### File Structure

```
src/
├── components/       # Shared UI components
│   └── ui/           # shadcn/ui primitives (auto-generated)
├── hooks/            # Custom React hooks
├── lib/              # Utilities, types, data, error handling
├── routes/           # TanStack Router file-based routes
├── router.tsx        # Router setup with QueryClient
├── server.ts         # Cloudflare Workers entry
├── start.ts          # TanStack Start instance + error middleware
└── styles.css        # Tailwind entry point
```

### Feature Conventions

- **No tests currently exist.** When adding tests, use `vitest` with `@testing-library/react`.
- **No i18n** — hardcoded English strings throughout.
- **Leaflet maps** with `react-leaflet` — only render client-side (guard with `useEffect` + mounted state).
