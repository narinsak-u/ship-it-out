<script setup lang="ts">
import { ref, computed } from "vue";
import { Search, Plus, Pencil, Trash2 } from "lucide-vue-next";
import { useHubs, useDeleteHub } from "@/hooks/useHubs";
import { getCarrier, hubStatusLabels } from "@/lib/carriers";
import { useAuthStore } from "@/stores/auth";
import { cn } from "@/lib/utils";
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
import HubFormModal from "@/components/HubFormModal.vue";

const auth = useAuthStore();

const { data: hubs, isLoading, isError, refetch } = useHubs();
const deleteHub = useDeleteHub();

const query = ref("");

const filtered = computed(() => {
  if (!hubs.value) return [];
  const q = query.value.trim().toLowerCase();
  return hubs.value.filter((h) => {
    if (!q) return true;
    return (
      h.name.toLowerCase().includes(q) ||
      h.address.toLowerCase().includes(q) ||
      getCarrier(h.carrierId)?.name.toLowerCase().includes(q)
    );
  });
});

const showForm = ref(false);
const editingHubId = ref<string | null>(null);

function openAdd() {
  editingHubId.value = null;
  showForm.value = true;
}

function openEdit(id: string) {
  editingHubId.value = id;
  showForm.value = true;
}

const hubStatusCounts = computed(() => {
  if (!hubs.value) return { total: 0, active: 0, maintenance: 0, closed: 0 };
  return {
    total: hubs.value.length,
    active: hubs.value.filter((h) => h.status === "active").length,
    maintenance: hubs.value.filter((h) => h.status === "maintenance").length,
    closed: hubs.value.filter((h) => h.status === "closed").length,
  };
});
</script>

<template>
  <div v-if="isLoading" class="space-y-4">
    <div class="grid grid-cols-4 gap-4">
      <Skeleton v-for="i in 4" :key="i" class="h-24 rounded-xl" />
    </div>
    <Skeleton class="h-64 rounded-xl" />
  </div>

  <div v-else-if="isError" class="py-12 text-center">
    <p class="font-mono text-sm text-destructive">Failed to load hubs.</p>
    <Button variant="outline" class="mt-4" @click="refetch()">Retry</Button>
  </div>

  <div v-else>
    <div class="grid grid-cols-2 gap-4 md:grid-cols-4">
      <div class="rounded-lg border border-border bg-secondary/50 p-4">
        <div class="font-mono text-[11px] uppercase tracking-widest text-muted-foreground">
          Total Hubs
        </div>
        <div class="mt-1 font-mono text-3xl font-semibold">{{ hubStatusCounts.total }}</div>
      </div>
      <div class="rounded-lg border border-border bg-secondary/50 p-4">
        <div class="font-mono text-[11px] uppercase tracking-widest text-muted-foreground">
          Active
        </div>
        <div class="mt-1 font-mono text-3xl font-semibold text-success">
          {{ hubStatusCounts.active }}
        </div>
      </div>
      <div class="rounded-lg border border-border bg-secondary/50 p-4">
        <div class="font-mono text-[11px] uppercase tracking-widest text-muted-foreground">
          Maintenance
        </div>
        <div class="mt-1 font-mono text-3xl font-semibold text-warning">
          {{ hubStatusCounts.maintenance }}
        </div>
      </div>
      <div class="rounded-lg border border-border bg-secondary/50 p-4">
        <div class="font-mono text-[11px] uppercase tracking-widest text-muted-foreground">
          Closed
        </div>
        <div class="mt-1 font-mono text-3xl font-semibold text-destructive">
          {{ hubStatusCounts.closed }}
        </div>
      </div>
    </div>

    <div class="mt-6 flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
      <div class="flex items-center gap-2 rounded-lg border border-border bg-card px-3 md:w-72">
        <Search class="h-4 w-4 text-muted-foreground" />
        <Input
          v-model="query"
          placeholder="Search hubs..."
          class="h-10 border-0 bg-transparent font-mono text-sm shadow-none focus-visible:ring-0"
        />
      </div>
      <Button v-if="auth.isAuthenticated" size="sm" class="gap-2" @click="openAdd">
        <Plus class="h-4 w-4" /> Add Hub
      </Button>
    </div>

    <!-- table -->
    <div class="mt-4 overflow-hidden rounded-xl border border-border">
      <Table>
        <TableHeader>
          <TableRow class="border-b border-border bg-secondary/50 font-mono text-[11px] uppercase tracking-widest text-muted-foreground hover:bg-secondary/50">
            <TableHead class="hidden md:table-cell">Name</TableHead>
            <TableHead class="hidden md:table-cell">Carrier</TableHead>
            <TableHead class="hidden md:table-cell">Address</TableHead>
            <TableHead class="hidden md:table-cell">Capacity</TableHead>
            <TableHead class="hidden md:table-cell">Status</TableHead>
            <TableHead v-if="auth.isAuthenticated" class="hidden md:table-cell text-right">Actions</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow
            v-for="h in filtered"
            :key="h.id"
            class="border-b border-border transition-colors hover:bg-secondary/40"
          >
            <TableCell class="font-mono text-sm">{{ h.name }}</TableCell>
            <TableCell class="font-mono text-sm text-muted-foreground">
              {{ getCarrier(h.carrierId)?.name ?? h.carrierId }}
            </TableCell>
            <TableCell class="font-mono text-xs text-muted-foreground">{{ h.address }}</TableCell>
            <TableCell>
              <div class="flex items-center gap-2">
                <div class="h-2 flex-1 overflow-hidden rounded-full bg-secondary">
                  <div
                    class="h-full rounded-full transition-all"
                    :class="h.currentUtilization / h.capacity > 0.8 ? 'bg-warning' : 'bg-primary'"
                    :style="{ width: `${Math.min(100, (h.currentUtilization / h.capacity) * 100)}%` }"
                  />
                </div>
                <span class="font-mono text-xs text-muted-foreground"
                  >{{ Math.round((h.currentUtilization / h.capacity) * 100) }}%</span
                >
              </div>
            </TableCell>
            <TableCell>
              <span
                :class="
                  cn(
                    'inline-flex items-center gap-1.5 rounded-full border px-2.5 py-0.5 text-xs font-mono uppercase tracking-wider',
                    h.status === 'active' ? 'bg-success/15 text-success border-success/30' : '',
                    h.status === 'maintenance' ? 'bg-warning/15 text-warning border-warning/30' : '',
                    h.status === 'closed'
                      ? 'bg-destructive/15 text-destructive border-destructive/30'
                      : '',
                  )
                "
              >
                <span class="h-1.5 w-1.5 rounded-full bg-current" />
                {{ hubStatusLabels[h.status] }}
              </span>
            </TableCell>
            <TableCell v-if="auth.isAuthenticated" class="text-right">
              <button
                @click="openEdit(h.id)"
                class="rounded p-1.5 text-muted-foreground hover:text-foreground"
              >
                <Pencil class="h-4 w-4" />
              </button>
              <button
                @click="deleteHub.mutate(h.id)"
                class="rounded p-1.5 text-muted-foreground hover:text-destructive"
              >
                <Trash2 class="h-4 w-4" />
              </button>
            </TableCell>
          </TableRow>
        </TableBody>
      </Table>

      <div
        v-if="filtered.length === 0"
        class="px-6 py-12 text-center font-mono text-sm text-muted-foreground"
      >
        No hubs match your filters.
      </div>
    </div>

    <HubFormModal :open="showForm" :hub-id="editingHubId" @close="showForm = false" />
  </div>
</template>
