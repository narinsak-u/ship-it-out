<script setup lang="ts">
import { ref, computed, defineAsyncComponent } from 'vue'
import { Truck, Warehouse, BarChart3, Package } from 'lucide-vue-next'
import { cn } from '@/lib/utils'

const DriversPanel = defineAsyncComponent(() => import('@/components/DriversPanel.vue'))
const HubsPanel = defineAsyncComponent(() => import('@/components/HubsPanel.vue'))
const AnalyticsPanel = defineAsyncComponent(() => import('@/components/AnalyticsPanel.vue'))
const DeliveriesPanel = defineAsyncComponent(() => import('@/components/DeliveriesPanel.vue'))

type Tab = 'drivers' | 'hubs' | 'analytics' | 'deliveries'

const activeTab = ref<Tab>('drivers')

const tabs: Array<{ key: Tab; label: string; icon: typeof Truck }> = [
  { key: 'drivers', label: 'Drivers', icon: Truck },
  { key: 'hubs', label: 'Hubs', icon: Warehouse },
  { key: 'analytics', label: 'Analytics', icon: BarChart3 },
  { key: 'deliveries', label: 'Active Deliveries', icon: Package },
]
</script>

<template>
  <div>

    <section class="border-b border-border bg-gradient-hero">
      <div class="mx-auto max-w-7xl px-6 py-14">
        <span class="font-mono text-xs uppercase tracking-widest text-primary">/ carriers</span>
        <h1 class="mt-3 text-4xl font-semibold tracking-tight md:text-5xl">Carrier operations</h1>
        <p class="mt-3 max-w-2xl text-muted-foreground">
          Manage carrier drivers, hubs, and monitor active deliveries across the fleet.
        </p>
      </div>
    </section>

    <section class="mx-auto max-w-7xl px-6 pt-0">
      <!-- Tab bar -->
      <div class="flex gap-1 -mt-px">
        <button
          v-for="t in tabs"
          :key="t.key"
          @click="activeTab = t.key"
          :class="cn(
            'flex items-center gap-2 rounded-t-lg px-5 py-3 font-mono text-sm transition-colors border border-border',
            activeTab === t.key
              ? 'bg-card text-foreground border-b-card -mb-px'
              : 'bg-transparent text-muted-foreground hover:text-foreground border-transparent hover:border-border',
          )"
        >
          <component :is="t.icon" class="h-4 w-4" />
          {{ t.label }}
        </button>
      </div>

      <!-- Tab content -->
      <div class="rounded-b-xl rounded-tr-xl border border-border bg-card p-6">
        <DriversPanel v-if="activeTab === 'drivers'" />
        <HubsPanel v-else-if="activeTab === 'hubs'" />
        <AnalyticsPanel v-else-if="activeTab === 'analytics'" />
        <DeliveriesPanel v-else-if="activeTab === 'deliveries'" />
      </div>
    </section>
  </div>
</template>
