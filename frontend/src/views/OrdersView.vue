<script setup lang="ts">
import { ref, computed } from "vue";
import { useRouter } from "vue-router";
import { useQuery, useMutation, useQueryClient } from "@tanstack/vue-query";
import { Search, Filter, ArrowRight, Plus, Pencil, Trash2 } from "lucide-vue-next";
import Input from "@/components/ui/Input.vue";
import StatusBadge from "@/components/StatusBadge.vue";
import { statusLabels, type ShipmentStatus } from "@/lib/orders";
import { fetchActiveDeliveries, deleteOrder } from "@/lib/api/orders";
import Button from "@/components/ui/Button.vue";
import {
  Table,
  TableHeader,
  TableBody,
  TableRow,
  TableHead,
  TableCell,
} from "@/components/ui/table";
import { cn } from "@/lib/utils";
import { useAuthStore } from "@/stores/auth";
import AuthModal from "@/components/AuthModal.vue";
import Skeleton from "@/components/ui/Skeleton.vue";

const authStore = useAuthStore();
const showAuthModal = ref(false);
const router = useRouter();
const queryClient = useQueryClient();

const { data: orders, isLoading } = useQuery({
  queryKey: ["orders"],
  queryFn: fetchActiveDeliveries,
});

const deleteMutation = useMutation({
  mutationFn: (id: string) => deleteOrder(id),
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ["orders"] });
  },
});

const FILTERS: Array<{ key: ShipmentStatus | "all"; label: string }> = [
  { key: "all", label: "All" },
  { key: "pending", label: "Pending" },
  { key: "in_transit", label: "In Transit" },
  { key: "out_for_delivery", label: "Out for Delivery" },
  { key: "delivered", label: "Delivered" },
  { key: "delayed", label: "Delayed" },
];

const filter = ref<ShipmentStatus | "all">("all");
const query = ref("");

const filtered = computed(() => {
  const q = query.value.trim().toLowerCase();
  return (orders.value ?? []).filter((o) => {
    if (filter.value !== "all" && o.status !== filter.value) return false;
    if (!q) return true;
    return (
      o.id.toLowerCase().includes(q) ||
      o.trackingNumber.toLowerCase().includes(q) ||
      o.customer.name.toLowerCase().includes(q) ||
      o.destination.toLowerCase().includes(q)
    );
  });
});

function onAuthenticated() {
  showAuthModal.value = false;
  router.push({ name: "order-create" });
}

function onGuest() {
  showAuthModal.value = false;
  router.push({ name: "order-create" });
}
</script>

<template>
  <div>
    <div v-if="isLoading" class="mx-auto max-w-7xl px-6 py-32">
      <Skeleton class="h-12 w-64" />
      <Skeleton class="mt-4 h-8 w-96" />
      <Skeleton class="mt-8 h-96 rounded-xl" />
    </div>
    <template v-if="!isLoading">
      <section class="border-b border-border bg-gradient-hero">
        <div class="mx-auto max-w-7xl px-6 py-14">
          <div class="flex items-start justify-between">
            <div>
              <span class="font-mono text-xs uppercase tracking-widest text-primary">/ orders</span>
              <h1 class="mt-3 text-4xl font-semibold tracking-tight md:text-5xl">
                Shipment manifest
              </h1>
              <p class="mt-3 max-w-2xl text-muted-foreground">
                {{ orders?.length ?? 0 }} total shipments tracked across all carriers.
              </p>
            </div>
            <div v-if="authStore.user" class="shrink-0">
              <RouterLink :to="{ name: 'order-create' }">
                <Button class="gap-2"> <Plus class="h-4 w-4" /> New Order </Button>
              </RouterLink>
            </div>
            <div v-else class="shrink-0">
              <Button class="gap-2" @click="showAuthModal = true">
                <Plus class="h-4 w-4" /> New Order
              </Button>
            </div>
          </div>
        </div>
      </section>

      <section class="mx-auto max-w-7xl px-6 py-10">
        <!-- Controls -->
        <div class="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
          <div class="flex items-center gap-2 rounded-lg border border-border bg-card px-3 md:w-96">
            <Search class="h-4 w-4 text-muted-foreground" />
            <Input
              v-model="query"
              placeholder="Search by ID, tracking, customer, destination"
              class="h-11 border-0 bg-transparent font-mono text-sm shadow-none focus-visible:ring-0"
            />
          </div>
          <div class="flex items-center gap-2 overflow-x-auto">
            <Filter class="h-4 w-4 shrink-0 text-muted-foreground" />
            <button
              v-for="f in FILTERS"
              :key="f.key"
              @click="filter = f.key"
              :class="
                cn(
                  'rounded-full border px-3 py-1.5 font-mono text-xs uppercase tracking-wider transition-colors',
                  filter === f.key
                    ? 'border-primary bg-primary/15 text-primary'
                    : 'border-border text-muted-foreground hover:text-foreground',
                )
              "
            >
              {{ f.label }}
            </button>
          </div>
        </div>

        <!-- Table -->
        <div class="mt-8 overflow-hidden rounded-xl border border-border bg-card shadow-elegant">
          <Table>
            <TableHeader>
              <TableRow
                class="border-b border-border bg-secondary/50 font-mono text-[11px] uppercase tracking-widest text-muted-foreground hover:bg-secondary/50"
              >
                <TableHead class="hidden md:table-cell">Order ID</TableHead>
                <TableHead class="hidden md:table-cell">Tracking</TableHead>
                <TableHead class="hidden md:table-cell">Customer</TableHead>
                <TableHead class="hidden md:table-cell">Route</TableHead>
                <TableHead class="hidden md:table-cell">Status</TableHead>
                <TableHead class="hidden md:table-cell text-right">ETA</TableHead>
                <TableHead v-if="authStore.user" class="hidden md:table-cell text-right"
                  >Actions</TableHead
                >
              </TableRow>
            </TableHeader>
            <TableBody>
              <TableRow
                v-for="o in filtered"
                :key="o.id"
                class="border-b border-border transition-colors hover:bg-secondary/40"
              >
                <TableCell>
                  <RouterLink
                    :to="{ name: 'order-detail', params: { orderId: o.id } }"
                    class="font-mono text-sm text-primary"
                  >
                    {{ o.id }}
                  </RouterLink>
                </TableCell>
                <TableCell class="font-mono text-sm text-muted-foreground">{{
                  o.trackingNumber
                }}</TableCell>
                <TableCell class="text-sm">{{ o.customer.name }}</TableCell>
                <TableCell>
                  <span class="flex items-center gap-2 font-mono text-xs text-muted-foreground">
                    <span>{{ o.origin }}</span>
                    <ArrowRight class="h-3 w-3 text-primary" />
                    <span>{{ o.destination }}</span>
                  </span>
                </TableCell>
                <TableCell><StatusBadge :status="o.status" /></TableCell>
                <TableCell class="font-mono text-xs text-muted-foreground text-right">{{
                  o.estimatedDelivery
                }}</TableCell>
                <TableCell v-if="authStore.user" class="text-right">
                  <button
                    @click.stop="router.push({ name: 'order-edit', params: { orderId: o.id } })"
                    class="rounded p-1.5 text-muted-foreground hover:text-primary"
                  >
                    <Pencil class="h-4 w-4" />
                  </button>
                  <button
                    @click.stop="deleteMutation.mutate(o.id)"
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
            class="px-6 py-16 text-center font-mono text-sm text-muted-foreground"
          >
            No shipments match your filters.
          </div>
        </div>

        <div class="mt-4 font-mono text-xs text-muted-foreground">
          Showing {{ filtered.length }} of {{ orders?.length ?? 0 }} · Status:
          {{ filter === "all" ? "All" : statusLabels[filter] }}
        </div>
      </section>

      <AuthModal
        v-if="showAuthModal"
        @close="showAuthModal = false"
        @authenticated="onAuthenticated"
        @guest="onGuest"
      />
    </template>
  </div>
</template>
