<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { Search, Filter, ArrowRight, Plus, Pencil, Trash2 } from 'lucide-vue-next'
import Input from '@/components/ui/Input.vue'
import StatusBadge from '@/components/StatusBadge.vue'
import { orders, statusLabels, type ShipmentStatus } from '@/lib/orders'
import Button from '@/components/ui/Button.vue'
import { cn } from '@/lib/utils'
import { useAuthStore } from '@/stores/auth'
import AuthModal from '@/components/AuthModal.vue'

const authStore = useAuthStore()
const showAuthModal = ref(false)

const router = useRouter()
const deletedIds = ref(new Set<string>())

const FILTERS: Array<{ key: ShipmentStatus | 'all'; label: string }> = [
  { key: 'all', label: 'All' },
  { key: 'pending', label: 'Pending' },
  { key: 'in_transit', label: 'In Transit' },
  { key: 'out_for_delivery', label: 'Out for Delivery' },
  { key: 'delivered', label: 'Delivered' },
  { key: 'delayed', label: 'Delayed' },
]

const filter = ref<ShipmentStatus | 'all'>('all')
const query = ref('')

const filtered = computed(() => {
  const q = query.value.trim().toLowerCase()
  return orders.filter((o) => {
    if (deletedIds.value.has(o.id)) return false
    if (filter.value !== 'all' && o.status !== filter.value) return false
    if (!q) return true
    return (
      o.id.toLowerCase().includes(q) ||
      o.trackingNumber.toLowerCase().includes(q) ||
      o.customer.toLowerCase().includes(q) ||
      o.destination.toLowerCase().includes(q)
    )
  })
})

function deleteOrder(id: string) {
  deletedIds.value.add(id)
  const idx = orders.findIndex((o) => o.id === id)
  if (idx !== -1) orders.splice(idx, 1)
}

function onAuthenticated() {
  showAuthModal.value = false
  router.push({ name: 'order-create' })
}

function onGuest() {
  showAuthModal.value = false
  router.push({ name: 'order-create' })
}
</script>

<template>
  <div>

    <section class="border-b border-border bg-gradient-hero">
      <div class="mx-auto max-w-7xl px-6 py-14">
        <div class="flex items-start justify-between">
          <div>
            <span class="font-mono text-xs uppercase tracking-widest text-primary">/ orders</span>
            <h1 class="mt-3 text-4xl font-semibold tracking-tight md:text-5xl">Shipment manifest</h1>
            <p class="mt-3 max-w-2xl text-muted-foreground">
              {{ orders.length }} total shipments tracked across all carriers.
            </p>
          </div>
          <div v-if="authStore.isAuthenticated" class="shrink-0">
            <RouterLink :to="{ name: 'order-create' }">
              <Button class="gap-2">
                <Plus class="h-4 w-4" /> New Order
              </Button>
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
            :class="cn(
              'rounded-full border px-3 py-1.5 font-mono text-xs uppercase tracking-wider transition-colors',
              filter === f.key
                ? 'border-primary bg-primary/15 text-primary'
                : 'border-border text-muted-foreground hover:text-foreground',
            )"
          >
            {{ f.label }}
          </button>
        </div>
      </div>

      <!-- Table -->
      <div class="mt-8 overflow-hidden rounded-xl border border-border bg-card shadow-elegant">
        <div class="hidden grid-cols-[1.1fr_1.4fr_1.6fr_2fr_1.2fr_0.6fr_0.6fr] gap-4 border-b border-border bg-secondary/50 px-6 py-3 font-mono text-[11px] uppercase tracking-widest text-muted-foreground md:grid">
          <span>Order ID</span>
          <span>Tracking</span>
          <span>Customer</span>
          <span>Route</span>
          <span>Status</span>
          <span class="text-right">ETA</span>
          <span v-if="authStore.isAuthenticated" class="text-right">Actions</span>
        </div>

        <div v-if="filtered.length === 0" class="px-6 py-16 text-center font-mono text-sm text-muted-foreground">
          No shipments match your filters.
        </div>
        <template v-else>
          <div
            v-for="o in filtered"
            :key="o.id"
            class="group grid grid-cols-1 gap-2 border-b border-border px-6 py-4 transition-colors last:border-0 hover:bg-secondary/40 md:grid-cols-[1.1fr_1.4fr_1.6fr_2fr_1.2fr_0.6fr_0.6fr] md:items-center md:gap-4"
          >
            <RouterLink :to="{ name: 'order-detail', params: { orderId: o.id } }" class="font-mono text-sm text-primary">
              {{ o.id }}
            </RouterLink>
            <span class="font-mono text-sm text-muted-foreground">{{ o.trackingNumber }}</span>
            <span class="text-sm">{{ o.customer }}</span>
            <span class="flex items-center gap-2 font-mono text-xs text-muted-foreground">
              <span>{{ o.origin }}</span>
              <ArrowRight class="h-3 w-3 text-primary" />
              <span>{{ o.destination }}</span>
            </span>
            <span><StatusBadge :status="o.status" /></span>
            <span class="font-mono text-xs text-muted-foreground md:text-right">{{ o.estimatedDelivery }}</span>
            <div v-if="authStore.isAuthenticated" class="flex justify-end gap-1">
              <button
                @click.stop="router.push({ name: 'order-edit', params: { orderId: o.id } })"
                class="rounded p-1.5 text-muted-foreground hover:text-primary"
              >
                <Pencil class="h-4 w-4" />
              </button>
              <button
                @click.stop="deleteOrder(o.id)"
                class="rounded p-1.5 text-muted-foreground hover:text-destructive"
              >
                <Trash2 class="h-4 w-4" />
              </button>
            </div>
          </div>
        </template>
      </div>

      <div class="mt-4 font-mono text-xs text-muted-foreground">
        Showing {{ filtered.length }} of {{ orders.length }} · Status: {{ filter === 'all' ? 'All' : statusLabels[filter] }}
      </div>
    </section>

    <AuthModal
      v-if="showAuthModal"
      @close="showAuthModal = false"
      @authenticated="onAuthenticated"
      @guest="onGuest"
    />
  </div>
</template>
