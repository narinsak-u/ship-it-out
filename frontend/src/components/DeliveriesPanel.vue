<script setup lang="ts">
import { ref, computed, defineAsyncComponent, onMounted, onUnmounted } from "vue";
import { Search, RefreshCw, ArrowRight, UserPlus } from "lucide-vue-next";
import { useActiveDeliveries, useUpdateShipmentStatus } from "@/hooks/useDeliveries";
import { useAssignDriver } from "@/hooks/useDrivers";
import { drivers as driverData } from "@/lib/carriers";
import { statusLabels, type ShipmentStatus } from "@/lib/orders";
import Input from "@/components/ui/Input.vue";
import Skeleton from "@/components/ui/Skeleton.vue";
import Button from "@/components/ui/Button.vue";

const ShipmentMap = defineAsyncComponent(() => import("@/components/ShipmentMap.vue"));

const { data: deliveries, isLoading, isError, refetch, dataUpdatedAt } = useActiveDeliveries();
const updateStatus = useUpdateShipmentStatus();
const assignDriver = useAssignDriver();

const query = ref("");
const mounted = ref(false);
const secondsSinceUpdate = ref(0);
let interval: ReturnType<typeof setInterval> | undefined;

onMounted(() => {
  mounted.value = true;
  interval = setInterval(() => {
    secondsSinceUpdate.value = Math.round((Date.now() - dataUpdatedAt.value) / 1000);
  }, 1000);
});

onUnmounted(() => {
  if (interval) clearInterval(interval);
});

const filtered = computed(() => {
  if (!deliveries.value) return [];
  const q = query.value.trim().toLowerCase();
  return deliveries.value.filter((o) => {
    if (!q) return true;
    return (
      o.id.toLowerCase().includes(q) ||
      o.trackingNumber.toLowerCase().includes(q) ||
      o.customer.name.toLowerCase().includes(q) ||
      o.carrier.toLowerCase().includes(q)
    );
  });
});

function getDriverForOrder(orderId: string) {
  const order = deliveries.value?.find((o) => o.id === orderId);
  if (!order?.driverId) return null;
  return driverData.find((d) => d.id === order.driverId) ?? null;
}

function handleAssign(orderId: string) {
  const order = deliveries.value?.find((o) => o.id === orderId);
  if (!order) return;
  const available = driverData.filter((d) => d.status === "available");
  if (available.length === 0) return;
  assignDriver.mutate({ driverId: available[0].id, orderId });
}
</script>

<template>
  <div v-if="isLoading" class="space-y-4">
    <Skeleton class="h-48 rounded-xl" />
    <Skeleton class="h-64 rounded-xl" />
  </div>

  <div v-else-if="isError" class="py-12 text-center">
    <p class="font-mono text-sm text-destructive">Failed to load active deliveries.</p>
    <Button variant="outline" class="mt-4" @click="refetch()">Retry</Button>
  </div>

  <div v-else>
    <!-- Controls bar -->
    <div class="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
      <div class="flex items-center gap-2 rounded-lg border border-border bg-card px-3 md:w-72">
        <Search class="h-4 w-4 text-muted-foreground" />
        <Input
          v-model="query"
          placeholder="Search active shipments..."
          class="h-10 border-0 bg-transparent font-mono text-sm shadow-none focus-visible:ring-0"
        />
      </div>
      <div class="flex items-center gap-3">
        <span class="font-mono text-xs text-muted-foreground">
          {{ deliveries?.length ?? 0 }} active · updated {{ secondsSinceUpdate }}s ago
        </span>
        <button
          @click="refetch()"
          class="rounded p-1.5 text-muted-foreground hover:text-foreground"
        >
          <RefreshCw class="h-4 w-4" />
        </button>
      </div>
    </div>

    <!-- Table -->
    <div class="mt-4 overflow-hidden rounded-xl border border-border">
      <div
        class="hidden grid-cols-[0.8fr_1fr_1fr_1fr_1fr_1fr_1fr_0.6fr] gap-4 border-b border-border bg-secondary/50 px-6 py-3 font-mono text-[11px] uppercase tracking-widest text-muted-foreground md:grid"
      >
        <span>Order ID</span>
        <span>Tracking</span>
        <span>Customer</span>
        <span>Carrier</span>
        <span>Driver</span>
        <span>Status</span>
        <span>ETA</span>
        <span class="text-right">Actions</span>
      </div>

      <div
        v-if="filtered.length === 0"
        class="px-6 py-12 text-center font-mono text-sm text-muted-foreground"
      >
        No active deliveries match your filters.
      </div>

      <div
        v-for="o in filtered"
        :key="o.id"
        class="group grid grid-cols-1 gap-2 border-b border-border px-6 py-4 transition-colors last:border-0 hover:bg-secondary/40 md:grid-cols-[0.8fr_1fr_1fr_1fr_1fr_1fr_1fr_0.6fr] md:items-center"
      >
        <div class="font-mono text-sm text-primary">{{ o.id }}</div>
        <div class="font-mono text-xs text-muted-foreground">{{ o.trackingNumber }}</div>
        <div class="font-mono text-sm">{{ o.customer.name }}</div>
        <div class="font-mono text-sm text-muted-foreground">{{ o.carrier }}</div>
        <div>
          <div v-if="getDriverForOrder(o.id)" class="font-mono text-sm">
            {{ getDriverForOrder(o.id)?.name }}
          </div>
          <div v-else class="font-mono text-xs text-muted-foreground">Unassigned</div>
        </div>
        <div>
          <select
            :value="o.status"
            @change="
              (e) =>
                updateStatus.mutate({
                  orderId: o.id,
                  status: (e.target as HTMLSelectElement).value as ShipmentStatus,
                })
            "
            class="rounded-lg border border-border bg-background px-2 py-1 font-mono text-xs"
          >
            <option v-for="(label, key) in statusLabels" :key="key" :value="key">
              {{ label }}
            </option>
          </select>
        </div>
        <div class="font-mono text-xs text-muted-foreground">{{ o.estimatedDelivery }}</div>
        <div class="flex justify-end gap-1">
          <button
            v-if="!o.driverId"
            @click="handleAssign(o.id)"
            class="rounded p-1.5 text-muted-foreground hover:text-primary"
            title="Assign driver"
          >
            <UserPlus class="h-4 w-4" />
          </button>
          <RouterLink
            :to="{ name: 'order-detail', params: { orderId: o.id } }"
            class="rounded p-1.5 text-muted-foreground hover:text-foreground"
          >
            <ArrowRight class="h-4 w-4" />
          </RouterLink>
        </div>
      </div>
    </div>

    <!-- Mini map -->
    <div
      v-if="mounted && filtered.length > 0"
      class="mt-6 overflow-hidden rounded-xl border border-border"
    >
      <Suspense>
        <div class="h-[300px] w-full">
          <ShipmentMap
            v-if="filtered[0]"
            :origin="filtered[0].customer.coords"
            :destination="filtered[0].receiver.coords"
            :current="filtered[0].currentCoords"
            :origin-label="filtered[0].origin"
            :destination-label="filtered[0].destination"
            :carrier="filtered[0].carrier"
            :status="filtered[0].status"
          />
        </div>
        <template #fallback>
          <div class="flex h-[300px] w-full items-center justify-center bg-gradient-hero">
            <div class="font-mono text-xs uppercase tracking-widest text-muted-foreground">
              Loading geo telemetry…
            </div>
          </div>
        </template>
      </Suspense>
    </div>
  </div>
</template>
