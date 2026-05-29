<script setup lang="ts">
import { ref, computed } from "vue";
import { useRouter } from "vue-router";
import { useQuery } from "@tanstack/vue-query";
import { ArrowRight, Boxes, Search, Truck, Globe2, Activity } from "lucide-vue-next";
import Button from "@/components/ui/Button.vue";
import Input from "@/components/ui/Input.vue";
import StatusBadge from "@/components/StatusBadge.vue";
import { fetchActiveDeliveries } from "@/lib/api/orders";
import { fetchAnalytics } from "@/lib/api/analytics";
import { trackShipment } from "@/lib/api/tracking";

const router = useRouter();

const { data: orders } = useQuery({
  queryKey: ["orders"],
  queryFn: fetchActiveDeliveries,
});

const { data: analytics } = useQuery({
  queryKey: ["analytics"],
  queryFn: fetchAnalytics,
});

const query = ref("");

const onTrack = async (e: Event) => {
  e.preventDefault();
  const q = query.value.trim();
  if (!q) return;

  const result = await trackShipment(q).catch(() => null);
  if (result) {
    router.push({ name: "order-detail", params: { orderId: result.shipment.id } });
  } else {
    router.push({ name: "orders" });
  }
};

const stats = computed(() => {
  const total = analytics.value?.total ?? 0;
  const delivered = analytics.value?.delivered ?? 0;
  return [
    {
      label: "Active shipments",
      value: total - delivered,
      icon: Truck,
    },
    {
      label: "Delivered",
      value: delivered,
      icon: Boxes,
    },
    { label: "On-time rate", value: "99.9%", icon: Activity }, // hard coded for now
    { label: "Provinces served", value: 77, icon: Globe2 }, // hard coded for now
  ];
});

const recent = computed(() => (orders.value ?? []).slice(0, 3));
</script>

<template>
  <div>
    <!-- Hero -->
    <section class="relative overflow-hidden bg-gradient-hero">
      <div
        class="absolute inset-0 opacity-[0.04]"
        :style="{
          backgroundImage:
            'linear-gradient(to right, currentColor 1px, transparent 1px), linear-gradient(to bottom, currentColor 1px, transparent 1px)',
          backgroundSize: '48px 48px',
        }"
      />
      <div class="relative mx-auto max-w-7xl px-6 py-24 md:py-32">
        <div class="max-w-3xl">
          <span
            class="inline-flex items-center gap-2 rounded-full border border-border bg-background/40 px-3 py-1 font-mono text-xs uppercase tracking-widest text-primary"
          >
            <span class="h-1.5 w-1.5 animate-pulse rounded-full bg-primary" />
            Live ops console
          </span>
          <h1
            class="mt-6 gap-1 flex flex-col text-5xl font-semibold leading-[1.05] tracking-tight md:text-7xl"
          >
            <div class="flex items-center gap-6">
              <span>Move fast</span>
              <Truck class="w-20 h-20 animate-running" />
            </div>
            <span class="bg-gradient-accent bg-clip-text w-fit">Break nothing.</span>
          </h1>
          <p class="mt-6 max-w-xl font-sans text-lg text-muted-foreground">
            Trace every parcel from origin warehouse to doorstep. Realtime telemetry, port-side
            timelines, and exception alerts in one console.
          </p>

          <form
            @submit="onTrack"
            class="mt-10 flex max-w-xl gap-2 rounded-xl border border-border bg-card p-2 shadow-elegant"
            aria-label="Track a shipment"
          >
            <div class="flex flex-1 items-center gap-2 px-3">
              <Search class="h-4 w-4 text-muted-foreground" />
              <Input
                v-model="query"
                placeholder="Enter tracking number"
                class="h-10 border-0 bg-transparent font-mono text-sm shadow-none focus-visible:ring-0"
              />
            </div>
            <Button type="submit" class="h-10 gap-2"> Track <ArrowRight class="h-4 w-4" /> </Button>
          </form>
        </div>
      </div>
    </section>

    <!-- Stats grid -->
    <section class="border-y border-border bg-card/40" aria-label="Key metrics">
      <div class="mx-auto grid max-w-7xl grid-cols-2 divide-border md:grid-cols-4 md:divide-x">
        <div v-for="s in stats" :key="s.label" class="flex items-center gap-4 px-6 py-8">
          <div
            class="flex h-11 w-11 items-center justify-center rounded-lg bg-secondary text-primary"
          >
            <component :is="s.icon" class="h-5 w-5" />
          </div>
          <div>
            <div class="font-mono text-2xl font-semibold">{{ s.value }}</div>
            <div class="text-xs uppercase tracking-wider text-muted-foreground">{{ s.label }}</div>
          </div>
        </div>
      </div>
    </section>

    <!-- Recent shipments -->
    <section class="mx-auto max-w-7xl px-6 py-20">
      <div class="flex items-end justify-between">
        <div>
          <h2 class="text-3xl font-semibold tracking-tight">Recent shipments</h2>
          <p class="mt-2 text-sm text-muted-foreground">Latest activity from the fleet.</p>
        </div>
        <RouterLink
          to="/orders"
          class="group flex items-center gap-1 font-mono text-sm text-primary"
        >
          View all <ArrowRight class="h-4 w-4 transition-transform group-hover:translate-x-1" />
        </RouterLink>
      </div>

      <div class="mt-8 grid gap-4 md:grid-cols-3">
        <RouterLink
          v-for="o in recent"
          :key="o.id"
          :to="{ name: 'order-detail', params: { orderId: o.id } }"
          class="group rounded-xl border border-border bg-card p-6 transition-all hover:border-primary/40 hover:shadow-elegant"
        >
          <div class="flex items-center justify-between">
            <span class="font-mono text-xs text-muted-foreground">{{ o.id }}</span>
            <StatusBadge :status="o.status" />
          </div>
          <div class="mt-4 font-mono text-lg">{{ o.trackingNumber }}</div>
          <div class="mt-1 text-sm text-muted-foreground">{{ o.customer.name }}</div>
          <div class="mt-6 flex items-center gap-2 font-mono text-xs">
            <span>{{ o.origin }}</span>
            <ArrowRight class="h-3 w-3 text-primary" />
            <span>{{ o.destination }}</span>
          </div>
          <div class="mt-4 h-1 overflow-hidden rounded-full bg-secondary">
            <div
              class="h-full bg-gradient-accent transition-all"
              :style="{ width: `${o.progress}%` }"
            />
          </div>
        </RouterLink>
      </div>
    </section>
  </div>
</template>
