<script setup lang="ts">
import { ref, computed } from "vue";
import { Search, Plus, Pencil, Trash2 } from "lucide-vue-next";
import { useHubs, useDeleteHub } from "@/hooks/useHubs";
import { hubStatusLabels } from "@/lib/hubs";
import { useAuthStore } from "@/stores/auth";
import { cn } from "@/lib/utils";
import Badge from "@/components/ui/Badge.vue";
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
import HubStatsCards from "@/components/HubStatsCards.vue";
import Pagination from "@/components/Pagination.vue";
import { usePagination } from "@/composables/usePagination";
import ConfirmDialog from "@/components/ConfirmDialog.vue";

const auth = useAuthStore();

const { data: hubs, isLoading, isError, refetch } = useHubs();
const deleteHub = useDeleteHub();
const deleteTarget = ref<string | null>(null);

const query = ref("");

const filtered = computed(() => {
  if (!hubs.value) return [];
  const q = query.value.trim().toLowerCase();
  return hubs.value.filter((h) => {
    if (!q) return true;
    return h.name.toLowerCase().includes(q) || h.address.toLowerCase().includes(q);
  });
});

const { currentPage, totalPages, pageItems, setPage } = usePagination(filtered, 10);

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

interface HubCounts {
  total: number;
  active: number;
  maintenance: number;
  closed: number;
  full: number;
}

const hubStatusCounts = computed(() => {
  if (!hubs.value) return { total: 0, active: 0, maintenance: 0, closed: 0, full: 0 };
  return hubs.value.reduce<HubCounts>(
    (acc, h) => {
      acc.total++;
      if (h.status in acc) acc[h.status as keyof HubCounts]++;
      return acc;
    },
    { total: 0, active: 0, maintenance: 0, closed: 0, full: 0 },
  );
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
    <HubStatsCards :counts="hubStatusCounts" />

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
          <TableRow
            class="border-b border-border bg-secondary/50 font-mono text-[11px] uppercase tracking-widest text-muted-foreground hover:bg-secondary/50"
          >
            <TableHead class="hidden md:table-cell">ID</TableHead>
            <TableHead class="hidden md:table-cell">Name</TableHead>
            <TableHead class="hidden md:table-cell">Carrier</TableHead>
            <TableHead class="hidden md:table-cell">Address</TableHead>
            <TableHead class="hidden md:table-cell">Capacity</TableHead>
            <TableHead class="hidden md:table-cell">Status</TableHead>
            <TableHead v-if="auth.isAuthenticated" class="hidden md:table-cell">
              Actions
            </TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow
            v-for="h in pageItems"
            :key="h.id"
            class="border-b border-border transition-colors hover:bg-secondary/40"
          >
            <TableCell class="font-mono text-xs text-muted-foreground">{{ h.id }}</TableCell>
            <TableCell class="font-mono text-sm">{{ h.name }}</TableCell>
            <TableCell class="font-mono text-sm text-muted-foreground">
              {{ h.carrierId }}
            </TableCell>
            <TableCell class="font-mono text-xs text-muted-foreground">{{ h.address }}</TableCell>
            <TableCell>
              <div class="flex items-center gap-2">
                <div class="h-2 flex-1 overflow-hidden rounded-full bg-secondary">
                  <div
                    class="h-full rounded-full transition-all"
                    :class="h.currentUtilization / h.capacity > 0.8 ? 'bg-warning' : 'bg-primary'"
                    :style="{
                      width: `${Math.min(100, (h.currentUtilization / h.capacity) * 100)}%`,
                    }"
                  />
                </div>
                <span class="font-mono text-xs text-muted-foreground">
                  {{ Math.round((h.currentUtilization / h.capacity) * 100) }}%
                </span>
              </div>
            </TableCell>
            <TableCell>
              <Badge
                variant="outline"
                class="gap-1.5 font-mono text-xs uppercase tracking-wider"
                :class="
                  h.status === 'active'
                    ? 'border-success/30 bg-success/15 text-success'
                    : h.status === 'maintenance'
                      ? 'border-warning/30 bg-warning/15 text-warning'
                      : 'border-destructive/30 bg-destructive/15 text-destructive'
                "
              >
                <span class="h-1.5 w-1.5 rounded-full bg-current" />
                {{ hubStatusLabels[h.status] }}
              </Badge>
            </TableCell>
            <TableCell v-if="auth.isAuthenticated">
              <button
                @click="openEdit(h.id)"
                class="rounded cursor-pointer p-1.5 text-muted-foreground hover:text-foreground"
              >
                <Pencil class="h-4 w-4" />
              </button>
              <button
                @click="deleteTarget = h.id"
                class="rounded cursor-pointer p-1.5 text-muted-foreground hover:text-destructive"
              >
                <Trash2 class="h-4 w-4" />
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
        No hubs match your filters.
      </div>
    </div>

    <ConfirmDialog
      :open="!!deleteTarget"
      title="Delete Hub"
      description="Are you sure you want to delete this hub? This action cannot be undone."
      :pending="deleteHub.isPending.value"
      @confirm="
        deleteTarget &&
        deleteHub.mutate(deleteTarget, {
          onSuccess: () => {
            deleteTarget = null;
          },
        })
      "
      @cancel="deleteTarget = null"
    />

    <HubFormModal :open="showForm" :hub-id="editingHubId" @close="showForm = false" />
  </div>
</template>
