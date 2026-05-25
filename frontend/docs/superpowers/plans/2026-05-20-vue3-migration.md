# Vue 3 Migration Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Migrate the Harbor Ops shipment tracking app from React + TanStack Start to a Vue 3 SPA with Vite and shadcn-vue.

**Architecture:** A standard Vue 3 Single Page Application using the Composition API (`<script setup>`). It replaces TanStack Router with Vue Router and TanStack Start's SSR with client-side rendering. Server state is managed by `@tanstack/vue-query`.

**Tech Stack:** Vue 3.5, Vite 6, Vue Router 4, Pinia 2, @tanstack/vue-query 5, Radix Vue, Tailwind CSS v4, Lucide Vue Next.

---

### Task 1: Dependency Overhaul

**Files:**

- Modify: `package.json`
- Modify: `vite.config.ts`
- Delete: `src/router.tsx`, `src/routeTree.gen.ts`, `src/server.ts`, `src/start.ts`

- [ ] **Step 1: Update `package.json` dependencies**
      Remove all React and TanStack Start related packages. Add Vue, Vue Router, Pinia, and Vue Query.

```json
{
  "dependencies": {
    "vue": "^3.5.13",
    "vue-router": "^4.5.0",
    "pinia": "^2.3.0",
    "@tanstack/vue-query": "^5.62.7",
    "radix-vue": "^1.9.11",
    "lucide-vue-next": "^0.468.0",
    "leaflet": "^1.9.4",
    "clsx": "^2.1.1",
    "tailwind-merge": "^2.5.5",
    "class-variance-authority": "^0.7.1"
  },
  "devDependencies": {
    "@vitejs/plugin-vue": "^5.2.1",
    "vue-tsc": "^2.2.0"
  }
}
```

- [ ] **Step 2: Run install**
      Run: `npm install`

- [ ] **Step 3: Update `vite.config.ts`**
      Replace the React configuration with the standard Vue plugin.

```typescript
import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";
import path from "path";

export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
  },
});
```

- [ ] **Step 4: Remove React-specific files**
      Remove files used by TanStack Start/React that are no longer needed.
      Run: `rm src/router.tsx src/routeTree.gen.ts src/server.ts src/start.ts`

- [ ] **Step 5: Commit**

```bash
git add package.json vite.config.ts
git rm src/router.tsx src/routeTree.gen.ts src/server.ts src/start.ts
git commit -m "chore: migrate dependencies to Vue 3"
```

---

### Task 2: Foundation (Main & App)

**Files:**

- Create: `src/main.ts`
- Create: `src/App.vue`
- Create: `src/router/index.ts`

- [ ] **Step 1: Create Vue Router configuration**

```typescript
import { createRouter, createWebHistory } from "vue-router";

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: "/",
      name: "home",
      component: () => import("@/views/HomeView.vue"),
    },
    {
      path: "/orders",
      name: "orders",
      component: () => import("@/views/OrdersView.vue"),
    },
    {
      path: "/orders/:orderId",
      name: "order-detail",
      component: () => import("@/views/OrderDetailView.vue"),
    },
  ],
});

export default router;
```

- [ ] **Step 2: Create `App.vue` root component**
      Include the QueryClientProvider equivalent and RouterView.

```vue
<script setup lang="ts">
import { VueQueryPlugin } from "@tanstack/vue-query";
import { RouterView } from "vue-router";
</script>

<template>
  <RouterView />
</template>
```

- [ ] **Step 3: Create `src/main.ts` entry point**

```typescript
import { createApp } from "vue";
import { createPinia } from "pinia";
import { VueQueryPlugin } from "@tanstack/vue-query";
import App from "./App.vue";
import router from "./router";
import "./styles.css";

const app = createApp(App);

app.use(createPinia());
app.use(router);
app.use(VueQueryPlugin);

app.mount("#app");
```

- [ ] **Step 4: Commit**

```bash
git add src/main.ts src/App.vue src/router/index.ts
git commit -m "feat: setup Vue foundation and router"
```

---

### Task 3: Port Base UI Components (Shadcn)

**Files:**

- Create: `src/components/ui/Badge.vue`, `src/components/ui/Button.vue`, `src/components/ui/Card.vue` (etc.)

- [ ] **Step 1: Port `Badge.vue`**

```vue
<script setup lang="ts">
import type { HTMLAttributes } from "vue";
import { type BadgeVariants, badgeVariants } from ".";
import { cn } from "@/lib/utils";

interface Props {
  variant?: BadgeVariants["variant"];
  class?: HTMLAttributes["class"];
}

defineProps<Props>();
</script>

<template>
  <div :class="cn(badgeVariants({ variant }), props.class)">
    <slot />
  </div>
</template>
```

- [ ] **Step 2: Repeat for core components used in routes**
      (Follow shadcn-vue registry for Button, Card, Separator, etc.)

- [ ] **Step 3: Commit**

```bash
git add src/components/ui/*.vue
git commit -m "feat: port shadcn-vue base components"
```

---

### Task 4: Port Domain Components

**Files:**

- Create: `src/components/SiteHeader.vue`
- Create: `src/components/StatusBadge.vue`
- Create: `src/components/ShipmentMap.vue`

- [ ] **Step 1: Port `SiteHeader.vue`**
      Convert from `SiteHeader.tsx`. Use `router-link` instead of `Link`.

- [ ] **Step 2: Port `StatusBadge.vue`**
      Convert from `StatusBadge.tsx`.

- [ ] **Step 3: Port `ShipmentMap.vue`**
      Convert from `ShipmentMap.tsx`. Use standard `leaflet` in `onMounted`.

- [ ] **Step 4: Commit**

```bash
git add src/components/*.vue
git commit -m "feat: port domain components to Vue"
```

---

### Task 5: Port Views (Pages)

**Files:**

- Create: `src/views/HomeView.vue`
- Create: `src/views/OrdersView.vue`
- Create: `src/views/OrderDetailView.vue`

- [ ] **Step 1: Port `HomeView.vue`**
      Simple landing page port.

- [ ] **Step 2: Port `OrdersView.vue`**
      Use `useQuery` from `@tanstack/vue-query` to fetch `orders`.

- [ ] **Step 3: Port `OrderDetailView.vue`**
      Use `useRoute` to get `orderId`. Use `useQuery` for details. Port the timeline and telemetry logic.

- [ ] **Step 4: Final verification**
      Run `npm run dev` and verify all pages work as expected.

- [ ] **Step 5: Commit**

```bash
git add src/views/*.vue
git commit -m "feat: port all view components"
```
