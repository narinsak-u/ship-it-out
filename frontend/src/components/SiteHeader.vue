<script setup lang="ts">
import { Package } from 'lucide-vue-next';
import { RouterLink, useRoute } from 'vue-router';

const route = useRoute();
</script>

<template>
  <header class="sticky top-0 z-40 border-b border-border bg-background/80 backdrop-blur-md">
    <div class="mx-auto flex h-16 max-w-7xl items-center justify-between px-6">
      <RouterLink to="/" class="flex items-center gap-2">
        <div class="flex h-8 w-8 items-center justify-center rounded-md bg-gradient-accent text-primary-foreground shadow-glow">
          <Package class="h-4 w-4" />
        </div>
        <span class="font-mono text-sm font-semibold tracking-tight">HARBOR/OPS</span>
      </RouterLink>
      <nav class="flex items-center gap-1 font-mono text-sm">
        <RouterLink
          to="/"
          class="rounded-md px-3 py-1.5 transition-colors"
          :class="route.path === '/' ? 'bg-secondary text-foreground' : 'text-muted-foreground hover:text-foreground'"
        >
          Home
        </RouterLink>
        <RouterLink
          to="/orders"
          class="rounded-md px-3 py-1.5 transition-colors"
          :class="route.path.startsWith('/orders') ? 'bg-secondary text-foreground' : 'text-muted-foreground hover:text-foreground'"
        >
          Orders
        </RouterLink>
        <RouterLink
          to="/carriers"
          class="rounded-md px-3 py-1.5 transition-colors"
          :class="route.path.startsWith('/carriers') ? 'bg-secondary text-foreground' : 'text-muted-foreground hover:text-foreground'"
        >
          Carriers
        </RouterLink>
      </nav>
    </div>
  </header>
</template>
