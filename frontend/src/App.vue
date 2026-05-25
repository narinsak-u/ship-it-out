<script setup lang="ts">
import { onMounted } from "vue";
import { RouterView } from "vue-router";
import { useAuthStore } from "@/stores/auth";
import SiteHeader from "@/components/SiteHeader.vue";
import SiteFooter from "@/components/SiteFooter.vue";

const authStore = useAuthStore();

onMounted(() => {
  authStore.init();
});
</script>

<template>
  <div class="flex min-h-screen flex-col">
    <SiteHeader />
    <main class="flex-1">
      <RouterView />
    </main>
    <SiteFooter />
  </div>
</template>
