# GEMINI.md — Harbor Ops (ship-simple)

## Project Structure

```
ship-simple/
├── frontend/          # Vue 3 SPA
│   ├── src/           # Application source
│   ├── package.json
│   ├── vite.config.ts
│   └── ...
├── backend/           # Go API server
│   ├── go.mod
│   └── ...
```

## Project Overview

Harbor Ops is a shipment tracking dashboard prototype. The **frontend** is a Vue 3 SPA; the **backend** is a Go API server.

- **Frontend Framework:** Vue 3.5+ (Composition API, `<script setup>`)
- **Build Tool:** Vite 6+
- **Routing:** Vue Router 4
- **Server State:** TanStack Vue Query 5
- **Global State:** Pinia 2
- **Styling:** Tailwind CSS v4 (using `@tailwindcss/vite` plugin)
- **UI Components:** `shadcn-vue` (built on `radix-vue`)
- **Icons:** Lucide Vue Next
- **Maps:** Leaflet
- **Package Manager:** Bun (preferred)
- **Backend:** Go 1.24+

---

## Getting Started (Frontend)

All commands must be run from the `frontend/` directory:

| Command           | Purpose                                                    |
| :---------------- | :--------------------------------------------------------- |
| `npm run dev`     | Start the development server (Vite)                        |
| `npm run build`   | Build the project for production (`vue-tsc && vite build`) |
| `npm run preview` | Preview the production build locally                       |
| `npm run lint`    | Run ESLint check                                           |
| `npm run format`  | Run Prettier auto-format                                   |

**Note:** No automated test framework is currently configured. To add tests, consider installing `vitest`.

---

## Getting Started (Backend)

```bash
cd backend
go run .
```

---

## Development Conventions

### Architecture & Pattern

- **Composition API:** Always use `<script setup lang="ts">` for components.
- **SPA:** The project is a client-side rendered SPA.
- **Server State:** Use `useQuery` from `@tanstack/vue-query` for all data fetching logic.
- **Routing:** Centralized configuration in `frontend/src/router/index.ts`.

### Naming & Structure

- **Components:** Use PascalCase for component names and filenames (e.g., `StatusBadge.vue`).
- **Views:** Page-level components are located in `frontend/src/views/` (e.g., `OrdersView.vue`).
- **UI Primitives:** shadcn-vue components are in `frontend/src/components/ui/`.
- **Path Aliases:** Use the `@/` alias to reference the `frontend/src/` directory.

### Styling & UI

- **Tailwind CSS v4:** Use utility classes exclusively. Avoid custom CSS files unless necessary.
- **Dynamic Classes:** Use the `cn()` utility from `@/lib/utils` for merging Tailwind classes (`clsx` + `tailwind-merge`).
- **Icons:** Use `lucide-vue-next` components.

### Data Handling

- Mock data and types are defined in `frontend/src/lib/orders.ts`.
- Types should be exported and reused across components.

---

## Important Note

`AGENTS.md` contains the full code style guide. **This `GEMINI.md` file is the foundational source of truth** for project overview and development workflows.
