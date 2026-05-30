<script setup lang="ts">
import { ref, computed, watch } from "vue";
import { useRouter } from "vue-router";
import { useQuery, useMutation, useQueryClient } from "@tanstack/vue-query";
import { toast } from "vue-sonner";
import { Plus } from "lucide-vue-next";
import { statusLabels, type ShipmentStatus } from "@/lib/orders";
import { fetchOrdersPaginated, deleteOrder } from "@/lib/api/orders";
import { orderKeys, deliveryKeys } from "@/lib/api/queryKeys";
import Button from "@/components/ui/Button.vue";
import {
  Table,
  TableHeader,
  TableBody,
  TableHead,
  TableCell,
  TableRow,
} from "@/components/ui/table";
import { useAuthStore } from "@/stores/auth";
import AuthModal from "@/components/AuthModal.vue";
import Skeleton from "@/components/ui/Skeleton.vue";
import Pagination from "@/components/Pagination.vue";
import ConfirmDialog from "@/components/ConfirmDialog.vue";
import ShipmentFilters from "@/components/ShipmentFilters.vue";
import OrderTableRow from "@/components/OrderTableRow.vue";

const authStore = useAuthStore();
const showAuthModal = ref(false);
const router = useRouter();
const queryClient = useQueryClient();

const deleteTarget = ref<string | null>(null);
const currentPage = ref(1);
const searchInput = ref("");
const debouncedSearch = ref("");
const filter = ref<ShipmentStatus | "all">("all");

watch(searchInput, (v, _old, onCleanup) => {
  const timer = setTimeout(() => {
    debouncedSearch.value = v;
    currentPage.value = 1;
  }, 300);
  onCleanup(() => clearTimeout(timer));
});

watch(filter, () => {
  currentPage.value = 1;
});

const queryKey = computed(() =>
  orderKeys.list({ page: currentPage.value, search: debouncedSearch.value, status: filter.value }),
);

const { data: pageData, isLoading } = useQuery({
  queryKey,
  queryFn: () =>
    fetchOrdersPaginated({
      page: currentPage.value,
      limit: 10,
      search: debouncedSearch.value || undefined,
      status: filter.value === "all" ? undefined : filter.value,
    }),
  staleTime: 60_000,
});

const pageItems = computed(() => pageData.value?.data ?? []);
const totalPages = computed(() => pageData.value?.pagination.totalPages ?? 1);
const totalItems = computed(() => pageData.value?.pagination.total ?? 0);

function setPage(page: number) {
  currentPage.value = page;
}

const deleteMutation = useMutation({
  mutationFn: (id: string) => deleteOrder(id),
  onSuccess: () => {
    toast.success("Order deleted");
    deleteTarget.value = null;
    queryClient.invalidateQueries({ queryKey: orderKeys.all });
    queryClient.invalidateQueries({ queryKey: deliveryKeys.all });
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

function onAuthenticated() {
  showAuthModal.value = false;
  router.push({ name: "order-create" });
}

function onGuest() {
  showAuthModal.value = false;
  router.push({ name: "order-create" });
}

function onDeleteOrder(id: string) {
  deleteTarget.value = id;
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
                {{ totalItems }} total shipments tracked across all carriers.
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
        <ShipmentFilters
          :search="searchInput"
          :filter="filter"
          :filters="FILTERS"
          @update:search="searchInput = $event"
          @update:filter="filter = $event"
        />

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
                <TableHead class="hidden md:table-cell">ETA</TableHead>
                <TableHead class="hidden md:table-cell">Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              <OrderTableRow
                v-for="o in pageItems"
                :key="o.id"
                :order="o"
                :is-authenticated="!!authStore.user"
                @edit="router.push({ name: 'order-edit', params: { orderId: $event } })"
                @delete="onDeleteOrder($event)"
              />
            </TableBody>
          </Table>

          <Pagination
            :current-page="currentPage"
            :total-pages="totalPages"
            :total-items="totalItems"
            :page-size="10"
            @update:current-page="setPage"
          />

          <div
            v-if="pageItems.length === 0"
            class="px-6 py-16 text-center font-mono text-sm text-muted-foreground"
          >
            No shipments match your filters.
          </div>
        </div>

        <div class="mt-4 font-mono text-xs text-muted-foreground">
          Status: {{ filter === "all" ? "All" : statusLabels[filter] }}
        </div>
      </section>

      <ConfirmDialog
        :open="!!deleteTarget"
        title="Delete Order"
        description="Are you sure you want to delete this order? This action cannot be undone."
        :pending="deleteMutation.isPending.value"
        @confirm="deleteTarget && deleteMutation.mutate(deleteTarget)"
        @cancel="deleteTarget = null"
      />

      <AuthModal
        :open="showAuthModal"
        @close="showAuthModal = false"
        @authenticated="onAuthenticated"
        @guest="onGuest"
      />
    </template>
  </div>
</template>
