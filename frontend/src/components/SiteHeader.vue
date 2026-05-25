<script setup lang="ts">
import { ref } from "vue";
import { Package, LogIn, LogOut } from "lucide-vue-next";
import { RouterLink, useRoute, useRouter } from "vue-router";
import { useAuthStore } from "@/stores/auth";
import AuthModal from "@/components/AuthModal.vue";

const route = useRoute();
const authStore = useAuthStore();
const router = useRouter();

const showAuthModal = ref(false);
</script>

<template>
  <header class="sticky top-0 z-40 border-b border-border bg-background/80 backdrop-blur-md">
    <div class="mx-auto flex h-16 max-w-7xl items-center justify-between px-6">
      <RouterLink to="/" class="flex items-center gap-2">
        <div
          class="flex h-8 w-8 items-center justify-center rounded-md bg-gradient-accent text-primary-foreground shadow-glow"
        >
          <Package class="h-4 w-4" />
        </div>
        <span class="font-mono text-sm font-semibold tracking-tight uppercase">
          thun-u-der/express
        </span>
      </RouterLink>
      <nav class="flex items-center gap-1 font-mono text-sm">
        <RouterLink
          to="/"
          class="rounded-md px-3 py-1.5 transition-colors"
          :class="
            route.path === '/'
              ? 'bg-secondary text-foreground'
              : 'text-muted-foreground hover:text-foreground'
          "
        >
          Home
        </RouterLink>
        <RouterLink
          to="/orders"
          class="rounded-md px-3 py-1.5 transition-colors"
          :class="
            route.path.startsWith('/orders')
              ? 'bg-secondary text-foreground'
              : 'text-muted-foreground hover:text-foreground'
          "
        >
          Orders
        </RouterLink>
        <RouterLink
          to="/carriers"
          class="rounded-md px-3 py-1.5 transition-colors"
          :class="
            route.path.startsWith('/carriers')
              ? 'bg-secondary text-foreground'
              : 'text-muted-foreground hover:text-foreground'
          "
        >
          Carriers
        </RouterLink>
        <div class="ml-4 flex items-center gap-2 border-l border-border pl-4">
          <template v-if="authStore.loading">
            <span class="font-mono text-xs text-muted-foreground">...</span>
          </template>
          <template v-else-if="authStore.user">
            <span class="font-mono text-xs text-muted-foreground"
              >Admin ({{ authStore.user.name }})</span
            >
            <button
              @click="
                authStore.logout();
                router.push({ name: 'home' });
              "
              class="flex cursor-pointer items-center gap-1.5 rounded-md px-3 py-1.5 font-mono text-xs text-muted-foreground transition-colors hover:text-foreground"
            >
              <LogOut class="h-3.5 w-3.5" /> Sign out
            </button>
          </template>
          <template v-else-if="authStore.isGuest">
            <span class="font-mono text-xs text-muted-foreground">Guest</span>
            <button
              @click="showAuthModal = true"
              class="flex cursor-pointer items-center gap-1.5 rounded-md px-3 py-1.5 font-mono text-xs transition-colors text-primary hover:text-foreground"
            >
              <LogIn class="h-3.5 w-3.5" /> Sign in
            </button>
          </template>
          <button
            v-else
            @click="showAuthModal = true"
            class="flex cursor-pointer items-center gap-1.5 rounded-md px-3 py-1.5 font-mono text-xs transition-colors text-primary hover:text-foreground"
          >
            <LogIn class="h-3.5 w-3.5" /> Sign in
          </button>
        </div>
      </nav>
    </div>
  </header>

  <AuthModal
    v-if="showAuthModal"
    @close="showAuthModal = false"
    @authenticated="showAuthModal = false"
    @guest="showAuthModal = false"
  />
</template>
