<script setup lang="ts">
import { ref, computed, defineAsyncComponent, onMounted, onUnmounted } from "vue";
import { useQuery } from "@tanstack/vue-query";
import { toast } from "vue-sonner";
import { Search, RefreshCw, Check, Eye } from "lucide-vue-next";
import { useActiveDeliveries, useUpdateShipmentStatus } from "@/hooks/useDeliveries";
import { fetchHubs } from "@/lib/api/carriers";
import { statusLabels, type ShipmentStatus } from "@/lib/orders";
import { useAuthStore } from "@/stores/auth";
import Input from "@/components/ui/Input.vue";
import Skeleton from "@/components/ui/Skeleton.vue";
import Button from "@/components/ui/Button.vue";
import {
  Table,
  TableHeader,
  TableBody,
  TableRow,
  TableHead,
  TableCell,
} from "@/components/ui/table";
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import Pagination from "@/components/Pagination.vue";
import { usePagination } from "@/composables/usePagination";

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
const selectedOrderId = ref<string | null>(null);
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

const { currentPage, totalPages, pageItems, setPage } = usePagination(filtered, 10);

const hubOptions = computed(() => (hubs.value ?? []).filter((h) => h.status === "active"));

const selectedOrder = computed(() => {
  const id = selectedOrderId.value;
  if (!id || !deliveries.value) return null;
  return deliveries.value.find((d) => d.id === id) ?? null;
});

const mapOrder = computed(() => selectedOrder.value ?? filtered.value[0] ?? null);

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
        toast.success("Delivery status updated");
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
      <Table>
        <TableHeader>
          <TableRow
            class="border-b border-border bg-secondary/50 font-mono text-[11px] uppercase tracking-widest text-muted-foreground hover:bg-secondary/50"
          >
            <TableHead class="hidden md:table-cell">Order ID</TableHead>
            <TableHead class="hidden md:table-cell">Tracking</TableHead>
            <TableHead class="hidden md:table-cell">Customer</TableHead>
            <TableHead class="hidden md:table-cell">Carrier</TableHead>
            <TableHead class="hidden md:table-cell">Status</TableHead>
            <TableHead class="hidden md:table-cell">Hub</TableHead>
            <TableHead class="hidden md:table-cell">ETA</TableHead>
            <TableHead class="hidden md:table-cell">Actions</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow
            v-for="o in pageItems"
            :key="o.id"
            class="cursor-pointer border-b border-border transition-colors hover:bg-secondary/40"
            :class="selectedOrderId === o.id ? 'bg-secondary/60' : ''"
            @click="selectedOrderId = o.id"
          >
            <TableCell class="font-mono text-sm text-primary">{{ o.id }}</TableCell>
            <TableCell class="font-mono text-xs text-muted-foreground">{{
              o.trackingNumber
            }}</TableCell>
            <TableCell class="font-mono text-sm">{{ o.customer.name }}</TableCell>
            <TableCell class="font-mono text-sm text-muted-foreground">{{ o.carrier }}</TableCell>
            <TableCell>
              <Select
                :model-value="draftStatus[o.id] ?? o.status"
                @update:model-value="(v) => (draftStatus[o.id] = v as ShipmentStatus)"
                :disabled="!auth.isAuthenticated"
              >
                <SelectTrigger
                  class="h-7 rounded-lg border border-border bg-background px-2 font-mono text-xs disabled:opacity-40"
                >
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
            </TableCell>
            <TableCell>
              <Select
                v-if="usesHubSelector(draftStatus[o.id] ?? o.status)"
                :model-value="draftHubId[o.id] ?? o.hubId ?? ''"
                @update:model-value="(v) => (draftHubId[o.id] = (v ?? '') as string)"
                :disabled="!auth.isAuthenticated"
              >
                <SelectTrigger
                  class="h-7 w-full rounded-lg border border-border bg-background px-2 font-mono text-xs disabled:opacity-40"
                >
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
            </TableCell>
            <TableCell class="font-mono text-xs text-muted-foreground">
              {{ o.estimatedDelivery }}
            </TableCell>
            <TableCell class="flex gap-1">
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
                class="rounded cursor-pointer p-1.5 text-muted-foreground transition-colors hover:text-primary disabled:opacity-30 disabled:pointer-events-none"
                title="Update"
              >
                <Check class="h-4 w-4" />
              </button>
            </TableCell>
          </TableRow>
        </TableBody>
      </Table>

      <Pagination
        :current-page="currentPage"
        :total-pages="totalPages"
        :total-items="filtered.length"
        :page-size="10"
        @update:current-page="setPage"
      />

      <div
        v-if="filtered.length === 0"
        class="px-6 py-12 text-center font-mono text-sm text-muted-foreground"
      >
        No active deliveries match your filters.
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
            v-if="mapOrder"
            :key="mapOrder.id"
            :origin="mapOrder.customer.coords"
            :destination="mapOrder.receiver.coords"
            :current="mapOrder.currentCoords"
            :origin-label="mapOrder.origin"
            :destination-label="mapOrder.destination"
            :carrier="mapOrder.carrier"
            :status="mapOrder.status"
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
