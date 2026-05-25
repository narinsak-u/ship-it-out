<script setup lang="ts">
import { ref, computed } from "vue";
import { Search, UserPlus } from "lucide-vue-next";
import { useDrivers } from "@/hooks/useDrivers";
import { getCarrier, type DriverStatus, driverStatusLabels } from "@/lib/carriers";
import { cn } from "@/lib/utils";
import Input from "@/components/ui/Input.vue";
import Skeleton from "@/components/ui/Skeleton.vue";
import Button from "@/components/ui/Button.vue";
import AssignDriverModal from "@/components/AssignDriverModal.vue";

const { data: drivers, isLoading, isError, refetch } = useDrivers();

const query = ref("");
const statusFilter = ref<DriverStatus | "all">("all");

const filtered = computed(() => {
  if (!drivers.value) return [];
  const q = query.value.trim().toLowerCase();
  return drivers.value.filter((d) => {
    if (statusFilter.value !== "all" && d.status !== statusFilter.value) return false;
    if (!q) return true;
    const carrier = getCarrier(d.carrierId);
    return (
      d.name.toLowerCase().includes(q) ||
      carrier?.name.toLowerCase().includes(q) ||
      d.vehicleInfo.toLowerCase().includes(q)
    );
  });
});

const statusCounts = computed(() => {
  if (!drivers.value) return { total: 0, available: 0, on_delivery: 0, off_duty: 0 };
  return {
    total: drivers.value.length,
    available: drivers.value.filter((d) => d.status === "available").length,
    on_delivery: drivers.value.filter((d) => d.status === "on_delivery").length,
    off_duty: drivers.value.filter((d) => d.status === "off_duty").length,
  };
});

const showAssignModal = ref(false);
</script>

<template>
  <div v-if="isLoading" class="space-y-4">
    <div class="grid grid-cols-4 gap-4">
      <Skeleton v-for="i in 4" :key="i" class="h-24 rounded-xl" />
    </div>
    <Skeleton class="h-64 rounded-xl" />
  </div>

  <div v-else-if="isError" class="py-12 text-center">
    <p class="font-mono text-sm text-destructive">Failed to load drivers.</p>
    <Button variant="outline" class="mt-4" @click="refetch()">Retry</Button>
  </div>

  <div v-else>
    <!-- Stats -->
    <div class="grid grid-cols-2 gap-4 md:grid-cols-4">
      <div class="rounded-lg border border-border bg-secondary/50 p-4">
        <div class="font-mono text-[11px] uppercase tracking-widest text-muted-foreground">
          Total Drivers
        </div>
        <div class="mt-1 font-mono text-3xl font-semibold">{{ statusCounts.total }}</div>
      </div>
      <div class="rounded-lg border border-border bg-secondary/50 p-4">
        <div class="font-mono text-[11px] uppercase tracking-widest text-muted-foreground">
          Available
        </div>
        <div class="mt-1 font-mono text-3xl font-semibold text-success">
          {{ statusCounts.available }}
        </div>
      </div>
      <div class="rounded-lg border border-border bg-secondary/50 p-4">
        <div class="font-mono text-[11px] uppercase tracking-widest text-muted-foreground">
          On Delivery
        </div>
        <div class="mt-1 font-mono text-3xl font-semibold text-info">
          {{ statusCounts.on_delivery }}
        </div>
      </div>
      <div class="rounded-lg border border-border bg-secondary/50 p-4">
        <div class="font-mono text-[11px] uppercase tracking-widest text-muted-foreground">
          Off Duty
        </div>
        <div class="mt-1 font-mono text-3xl font-semibold text-muted-foreground">
          {{ statusCounts.off_duty }}
        </div>
      </div>
    </div>

    <!-- Controls -->
    <div class="mt-6 flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
      <div class="flex items-center gap-2 rounded-lg border border-border bg-card px-3 md:w-72">
        <Search class="h-4 w-4 text-muted-foreground" />
        <Input
          v-model="query"
          placeholder="Search drivers..."
          class="h-10 border-0 bg-transparent font-mono text-sm shadow-none focus-visible:ring-0"
        />
      </div>
      <div class="flex items-center gap-2">
        <button
          v-for="s in ['all', 'available', 'on_delivery', 'off_duty'] as const"
          :key="s"
          @click="statusFilter = s"
          :class="
            cn(
              'rounded-full border px-3 py-1.5 font-mono text-xs uppercase tracking-wider transition-colors',
              statusFilter === s
                ? 'border-primary bg-primary/15 text-primary'
                : 'border-border text-muted-foreground hover:text-foreground',
            )
          "
        >
          {{ s === "all" ? "All" : driverStatusLabels[s] }}
        </button>
        <Button size="sm" class="gap-2 ml-2" @click="showAssignModal = true">
          <UserPlus class="h-4 w-4" /> Assign
        </Button>
      </div>
    </div>

    <!-- Table -->
    <div class="mt-4 overflow-hidden rounded-xl border border-border">
      <div
        class="hidden grid-cols-[1.5fr_1.5fr_1fr_1fr_0.8fr] gap-4 border-b border-border bg-secondary/50 px-6 py-3 font-mono text-[11px] uppercase tracking-widest text-muted-foreground md:grid"
      >
        <span>Driver</span>
        <span>Carrier</span>
        <span>Status</span>
        <span>Vehicle</span>
        <span class="text-right">Actions</span>
      </div>

      <div
        v-if="filtered.length === 0"
        class="px-6 py-12 text-center font-mono text-sm text-muted-foreground"
      >
        No drivers match your filters.
      </div>

      <div
        v-for="d in filtered"
        :key="d.id"
        class="group grid grid-cols-1 gap-2 border-b border-border px-6 py-4 transition-colors last:border-0 hover:bg-secondary/40 md:grid-cols-[1.5fr_1.5fr_1fr_1fr_0.8fr] md:items-center"
      >
        <div>
          <div class="font-mono text-sm">{{ d.name }}</div>
          <div class="font-mono text-xs text-muted-foreground">{{ d.email }}</div>
        </div>
        <div class="font-mono text-sm text-muted-foreground">
          {{ getCarrier(d.carrierId)?.name ?? d.carrierId }}
        </div>
        <div>
          <span
            :class="
              cn(
                'inline-flex items-center gap-1.5 rounded-full border px-2.5 py-0.5 text-xs font-mono uppercase tracking-wider',
                d.status === 'available' ? 'bg-success/15 text-success border-success/30' : '',
                d.status === 'on_delivery' ? 'bg-info/15 text-info border-info/30' : '',
                d.status === 'off_duty' ? 'bg-muted text-muted-foreground border-border' : '',
              )
            "
          >
            <span class="h-1.5 w-1.5 rounded-full bg-current" />
            {{ driverStatusLabels[d.status] }}
          </span>
        </div>
        <div class="font-mono text-sm text-muted-foreground">{{ d.vehicleInfo }}</div>
        <div class="text-right">
          <Button
            v-if="d.status === 'available'"
            variant="outline"
            size="sm"
            @click="showAssignModal = true"
          >
            Assign
          </Button>
        </div>
      </div>
    </div>

    <div class="mt-4 font-mono text-xs text-muted-foreground">
      Showing {{ filtered.length }} of {{ statusCounts.total }} drivers
    </div>

    <AssignDriverModal v-if="showAssignModal" @close="showAssignModal = false" />
  </div>
</template>
