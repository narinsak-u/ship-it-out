<script setup lang="ts">
import { ref, computed, onMounted, defineAsyncComponent } from "vue";
import { useRoute, RouterLink } from "vue-router";
import { useQuery } from "@tanstack/vue-query";
import { ArrowLeft, MapPin, Truck, Calendar, Hash, User, Weight, Maximize2 } from "lucide-vue-next";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import StatusBadge from "@/components/StatusBadge.vue";
import Skeleton from "@/components/ui/Skeleton.vue";
import { fetchOrder, fetchOrderEvents } from "@/lib/api/orders";

const ShipmentMap = defineAsyncComponent(() => import("@/components/ShipmentMap.vue"));

const route = useRoute();
const orderId = route.params.orderId as string;

const { data: order, isLoading } = useQuery({
  queryKey: ["order", orderId],
  queryFn: () => fetchOrder(orderId),
});

const { data: events } = useQuery({
  queryKey: ["order-events", orderId],
  queryFn: () => {
    if (!order.value) return [];
    return fetchOrderEvents(order.value.trackingNumber);
  },
  enabled: computed(() => !!order.value),
});

const timeline = computed(() => [...(events.value ?? [])].reverse());

const mounted = ref(false);
onMounted(() => {
  mounted.value = true;
});

const meta = computed(() => {
  const o = order.value;
  if (!o) return [];
  return [
    { icon: Hash, label: "Tracking #", value: o.trackingNumber },
    { icon: User, label: "Customer", value: o.customer.name },
    { icon: Truck, label: "Carrier", value: o.carrier },
    { icon: Weight, label: "Weight/Kg", value: `${o.weight} kg` },
    { icon: Calendar, label: "Created", value: o.createdAt },
  ];
});
</script>

<template>
  <div v-if="isLoading" class="mx-auto max-w-7xl px-6 py-32">
    <div class="grid grid-cols-1 gap-8 lg:grid-cols-2">
      <div class="space-y-6">
        <Skeleton class="h-8 w-48" />
        <Skeleton class="h-12 w-96" />
        <Skeleton class="h-48 rounded-xl" />
        <Skeleton class="h-64 rounded-xl" />
      </div>
      <Skeleton class="h-105 rounded-xl lg:h-full" />
    </div>
  </div>
  <div v-else-if="order">
    <div class="mx-auto grid max-w-400 gap-0 lg:grid-cols-[minmax(0,1fr)_minmax(0,1.1fr)]">
      <!-- LEFT: details -->
      <div class="border-r border-border">
        <div class="px-6 py-8 lg:px-10 lg:py-10">
          <RouterLink
            to="/orders"
            class="group inline-flex items-center gap-1.5 font-mono text-xs uppercase tracking-widest text-muted-foreground hover:text-foreground"
          >
            <ArrowLeft class="h-3.5 w-3.5 transition-transform group-hover:-translate-x-1" />
            All orders
          </RouterLink>

          <div class="mt-6 flex flex-wrap items-center gap-3">
            <h1 class="font-mono text-4xl font-semibold tracking-tight">
              {{ order.trackingNumber }}
            </h1>
            <StatusBadge :status="order.status" />
          </div>
          <div class="mt-2 font-mono text-sm text-muted-foreground">
            Order ID <span class="text-primary">{{ order.id }}</span> · {{ order.customer.name }}
          </div>

          <!-- Route summary card -->
          <Card class="mt-8 shadow-elegant">
            <CardContent class="p-5">
              <div class="flex items-start gap-4">
                <div class="mt-1 flex flex-col items-center gap-1">
                  <span class="h-2.5 w-2.5 rounded-full bg-primary" />
                  <span class="h-10 w-px border-l border-dashed border-border" />
                  <span class="h-2.5 w-2.5 rounded-full bg-accent" />
                </div>
                <div class="flex-1 space-y-3">
                  <div>
                    <div
                      class="font-mono text-[11px] uppercase tracking-widest text-muted-foreground"
                    >
                      From
                    </div>
                    <div class="font-mono text-sm">{{ order.origin }}</div>
                  </div>
                  <div>
                    <div
                      class="font-mono text-[11px] uppercase tracking-widest text-muted-foreground"
                    >
                      To
                    </div>
                    <div class="font-mono text-sm">{{ order.destination }}</div>
                  </div>
                </div>
                <div
                  class="rounded-md border border-border bg-secondary px-3 py-1.5 font-mono text-xs text-primary"
                >
                  {{ order.carrier }}
                </div>
              </div>

              <div class="mt-5 grid grid-cols-2 gap-4 border-t border-border pt-5 sm:grid-cols-3">
                <div>
                  <div
                    class="font-mono text-[11px] uppercase tracking-widest text-muted-foreground"
                  >
                    Progress
                  </div>
                  <div class="font-mono text-sm">{{ order.progress }}%</div>
                </div>
                <div>
                  <div
                    class="font-mono text-[11px] uppercase tracking-widest text-muted-foreground"
                  >
                    Created
                  </div>
                  <div class="font-mono text-sm">{{ order.createdAt }}</div>
                </div>
                <div>
                  <div
                    class="font-mono text-[11px] uppercase tracking-widest text-muted-foreground"
                  >
                    ETA
                  </div>
                  <div class="font-mono text-sm text-primary">{{ order.estimatedDelivery }}</div>
                </div>
              </div>

              <div class="mt-4 h-1.5 overflow-hidden rounded-full bg-secondary">
                <div
                  class="h-full bg-gradient-accent transition-all"
                  :style="{ width: `${order.progress}%` }"
                />
              </div>
            </CardContent>
          </Card>

          <!-- Metadata grid -->
          <div class="mt-6 grid grid-cols-2 gap-3 sm:grid-cols-3">
            <Card v-for="m in meta" :key="m.label" class="rounded-lg shadow-elegant">
              <CardHeader class="flex flex-row items-center gap-2 p-4 pb-0">
                <component :is="m.icon" class="h-3.5 w-3.5 text-primary" />
                <CardTitle
                  class="font-mono text-[10px] uppercase tracking-wider text-muted-foreground"
                >
                  {{ m.label }}
                </CardTitle>
              </CardHeader>
              <CardContent class="p-4 pt-1.5">
                <div class="font-mono text-sm">{{ m.value }}</div>
              </CardContent>
            </Card>
          </div>

          <!-- Timeline -->
          <div class="mt-10">
            <h2 class="font-mono text-sm uppercase tracking-widest text-muted-foreground">
              Shipment status
            </h2>
            <ol class="relative mt-5 space-y-5 border-l border-border pl-6">
              <li v-for="(e, i) in timeline" :key="i" class="relative">
                <span
                  :class="[
                    'absolute -left-7.75 flex h-4 w-4 items-center justify-center rounded-full ring-4 ring-background',
                    i === 0 ? 'bg-primary shadow-glow' : 'bg-muted-foreground/40',
                  ]"
                >
                  <span
                    v-if="i === 0"
                    class="h-1.5 w-1.5 animate-pulse rounded-full bg-primary-foreground"
                  />
                </span>
                <div class="flex items-center justify-between">
                  <span class="font-mono text-sm font-medium">{{ e.status }}</span>
                  <span class="font-mono text-xs text-muted-foreground">{{ e.timestamp }}</span>
                </div>
                <div class="mt-1 flex items-center gap-1.5 font-mono text-xs text-primary">
                  <MapPin class="h-3 w-3" />
                  {{ e.location.name }}
                </div>
                <p v-if="e.description" class="mt-1 text-sm text-muted-foreground">
                  {{ e.description }}
                </p>
              </li>
            </ol>
          </div>
        </div>
      </div>

      <!-- RIGHT: map -->
      <div class="relative bg-secondary/40 lg:sticky lg:top-16 lg:h-[calc(100vh-4rem)]">
        <div class="h-105 w-full lg:h-full">
          <Suspense v-if="mounted">
            <ShipmentMap
              :origin="order.customer.coords"
              :destination="order.receiver.coords"
              :current="order.currentCoords"
              :origin-label="order.origin"
              :destination-label="order.destination"
              :carrier="order.carrier"
              :status="order.status"
            />
            <template #fallback>
              <div class="flex h-full w-full items-center justify-center bg-gradient-hero">
                <div class="font-mono text-xs uppercase tracking-widest text-muted-foreground">
                  Loading geo telemetry…
                </div>
              </div>
            </template>
          </Suspense>
          <div v-else class="flex h-full w-full items-center justify-center bg-gradient-hero">
            <div class="font-mono text-xs uppercase tracking-widest text-muted-foreground">
              Loading geo telemetry…
            </div>
          </div>
        </div>

        <!-- Floating telemetry card -->
        <div
          class="pointer-events-none absolute left-16 top-4 z-400 rounded-lg border border-border bg-card/95 px-4 py-3 font-mono text-xs shadow-elegant backdrop-blur"
        >
          <div class="flex items-center gap-2 text-muted-foreground">
            <span class="h-1.5 w-1.5 animate-pulse rounded-full bg-primary" />
            LIVE TELEMETRY
          </div>
          <div class="mt-1 text-sm text-foreground">{{ order.carrier }}</div>
          <div class="text-muted-foreground">
            {{ order.currentCoords.lat.toFixed(2) }}°, {{ order.currentCoords.lng.toFixed(2) }}°
          </div>
        </div>

        <div
          class="pointer-events-none absolute right-4 top-4 z-400 rounded-lg border border-border bg-card/95 px-3 py-2 font-mono text-[10px] uppercase tracking-widest text-muted-foreground shadow-elegant backdrop-blur"
        >
          <Maximize2 class="mr-1 inline h-3 w-3" />
          Geo route
        </div>
      </div>
    </div>
  </div>
  <div v-else>
    <div class="mx-auto max-w-2xl px-6 py-32 text-center">
      <h1 class="font-mono text-4xl">404</h1>
      <p class="mt-3 text-muted-foreground">Shipment not found.</p>
      <RouterLink to="/orders" class="mt-6 inline-block font-mono text-sm text-primary"
        >← Back to orders</RouterLink
      >
    </div>
  </div>
</template>
