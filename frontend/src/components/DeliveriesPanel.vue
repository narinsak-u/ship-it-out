<script setup lang="ts">
import { ref, computed, defineAsyncComponent, onMounted, onUnmounted } from "vue";
import { useQuery } from "@tanstack/vue-query";
import { Search, RefreshCw, Check, Eye } from "lucide-vue-next";
import { useActiveDeliveries, useUpdateShipmentStatus } from "@/hooks/useDeliveries";
import { fetchHubs } from "@/lib/api/carriers";
import { statusLabels, type ShipmentStatus } from "@/lib/orders";
import { useAuthStore } from "@/stores/auth";
import Input from "@/components/ui/Input.vue";
import Skeleton from "@/components/ui/Skeleton.vue";
import Button from "@/components/ui/Button.vue";
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

const ShipmentMap = defineAsyncComponent(() => import("@/components/ShipmentMap.vue"));

const auth = useAuthStore();

const { data: deliveries, isLoading, isError, refetch, dataUpdatedAt } = useActiveDeliveries();
const updateStatus = useUpdateShipmentStatus();

const { data: hubs } = useQuery({
  queryKey: ["hubs"],
  queryFn: fetchHubs,
});

const query = ref("");
const mounted = ref(false);
const secondsSinceUpdate = ref(0);
let interval: ReturnType<typeof setInterval> | undefined;

const draftStatus = ref<Record<string, ShipmentStatus>>({});
const draftHubId = ref<Record<string, string>>({});

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

const hubOptions = computed(() => (hubs.value ?? []).filter((h) => h.status === "active"));

function usesHubSelector(status: ShipmentStatus) {
  return (
    status === "departed" ||
    status === "in_transit" ||
    status === "out_for_delivery" ||
    status === "delayed"
  );
}

function canUpdate(orderId: string) {
  const o = deliveries.value?.find((d) => d.id === orderId);
  if (!o) return false;
  const ds = draftStatus.value[orderId] ?? o.status;
  const dh = draftHubId.value[orderId] ?? "";
  const changed = ds !== o.status || (usesHubSelector(ds) && dh !== "");
  if (!changed) return false;
  if (usesHubSelector(ds) && !dh) return false;
  return true;
}

function handleUpdate(orderId: string) {
  const o = deliveries.value?.find((d) => d.id === orderId);
  if (!o) return;
  const status = draftStatus.value[orderId] ?? o.status;
  if (usesHubSelector(status) && !draftHubId.value[orderId]) return;
  updateStatus.mutate(
    { orderId, status, hubId: draftHubId.value[orderId] },
    {
      onSuccess: () => {
        delete draftStatus.value[orderId];
        delete draftHubId.value[orderId];
      },
    },
  );
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
        class="hidden grid-cols-[0.8fr_1fr_1fr_1fr_1fr_0.8fr_0.8fr_1fr_0.5fr] gap-4 border-b border-border bg-secondary/50 px-6 py-3 font-mono text-[11px] uppercase tracking-widest text-muted-foreground md:grid"
      >
        <span>Order ID</span>
        <span>Tracking</span>
        <span>Customer</span>
        <span>Carrier</span>
        <span>Status</span>
        <span>Hub</span>
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
        class="group grid grid-cols-1 gap-2 border-b border-border px-6 py-4 transition-colors last:border-0 hover:bg-secondary/40 md:grid-cols-[0.8fr_1fr_1fr_1fr_1fr_0.8fr_0.8fr_1fr_0.5fr] md:items-center"
      >
        <div class="font-mono text-sm text-primary">{{ o.id }}</div>
        <div class="font-mono text-xs text-muted-foreground">{{ o.trackingNumber }}</div>
        <div class="font-mono text-sm">{{ o.customer.name }}</div>
        <div class="font-mono text-sm text-muted-foreground">{{ o.carrier }}</div>
        <div>
          <Select
            :model-value="draftStatus[o.id] ?? o.status"
            @update:model-value="(v) => draftStatus[o.id] = v as ShipmentStatus"
            :disabled="!auth.isAuthenticated"
          >
            <SelectTrigger class="h-7 rounded-lg border border-border bg-background px-2 font-mono text-xs disabled:opacity-40">
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectGroup>
                <SelectItem v-for="(label, key) in statusLabels" :key="key" :value="key">
                  {{ label }}
                </SelectItem>
              </SelectGroup>
            </SelectContent>
          </Select>
        </div>
        <!-- Hub column -->
        <div>
          <Select
            v-if="usesHubSelector(draftStatus[o.id] ?? o.status)"
            :model-value="draftHubId[o.id] ?? o.hubId ?? ''"
            @update:model-value="(v) => draftHubId[o.id] = (v ?? '') as string"
            :disabled="!auth.isAuthenticated"
          >
            <SelectTrigger class="h-7 w-full rounded-lg border border-border bg-background px-2 font-mono text-xs disabled:opacity-40">
              <SelectValue placeholder="Select hub..." />
            </SelectTrigger>
            <SelectContent>
              <SelectGroup>
                <SelectItem v-for="h in hubOptions" :key="h.id" :value="h.id">
                  {{ h.name }}
                </SelectItem>
              </SelectGroup>
            </SelectContent>
          </Select>
          <span v-else class="font-mono text-xs text-muted-foreground">&mdash;</span>
        </div>
        <div class="font-mono text-xs text-muted-foreground">{{ o.estimatedDelivery }}</div>

        <!-- Actions -->
        <div class="flex justify-end gap-1">
          <RouterLink
            :to="{ name: 'order-detail', params: { orderId: o.id } }"
            class="rounded p-1.5 text-muted-foreground transition-colors hover:text-foreground"
            title="View details"
          >
            <Eye class="h-4 w-4" />
          </RouterLink>
          <button
            @click="handleUpdate(o.id)"
            :disabled="!auth.isAuthenticated || !canUpdate(o.id)"
            class="rounded p-1.5 text-muted-foreground transition-colors hover:text-primary disabled:opacity-30 disabled:pointer-events-none"
            title="Update"
          >
            <Check class="h-4 w-4" />
          </button>
        </div>
      </div>
    </div>

    <!-- Mini map -->
    <div
      v-if="mounted && filtered.length > 0"
      class="mt-6 overflow-hidden rounded-xl border border-border"
    >
      <Suspense>
        <div class="h-75 w-full">
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
          <div class="flex h-75 w-full items-center justify-center bg-gradient-hero">
            <div class="font-mono text-xs uppercase tracking-widest text-muted-foreground">
              Loading geo telemetry…
            </div>
          </div>
        </template>
      </Suspense>
    </div>
  </div>
</template>
