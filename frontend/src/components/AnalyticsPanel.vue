<script setup lang="ts">
import { computed } from "vue";
import { useCarriers } from "@/hooks/useCarriers";
import { orders, statusLabels } from "@/lib/orders";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import Skeleton from "@/components/ui/Skeleton.vue";
import Button from "@/components/ui/Button.vue";

const { data: carriersData, isLoading, isError, refetch } = useCarriers();

const kpis = computed(() => {
  const total = orders.length;
  const delivered = orders.filter((o) => o.status === "delivered").length;
  const onTime = Math.round((delivered / Math.max(total, 1)) * 100);
  const activeCarriers = carriersData.value?.filter((c) => c.status === "active").length ?? 0;
  return { total, onTime, activeCarriers, avgDeliveryTime: "3.2 days" };
});

const carrierPerformance = computed(() => {
  if (!carriersData.value) return [];
  return carriersData.value.map((c) => {
    const carrierOrders = orders.filter((o) => o.carrier === c.name);
    const delivered = carrierOrders.filter((o) => o.status === "delivered").length;
    return {
      name: c.name,
      total: carrierOrders.length,
      delivered,
      onTimeRate:
        carrierOrders.length > 0 ? Math.round((delivered / carrierOrders.length) * 100) : 0,
    };
  });
});

const statusDistribution = computed(() => {
  const counts: Record<string, number> = {};
  for (const o of orders) {
    counts[o.status] = (counts[o.status] || 0) + 1;
  }
  return Object.entries(counts).map(([status, count]) => ({
    status,
    label: statusLabels[status as keyof typeof statusLabels] ?? status,
    count,
    pct: Math.round((count / orders.length) * 100),
  }));
});

const maxCarrierOrders = computed(() =>
  Math.max(...carrierPerformance.value.map((c) => c.total), 1),
);

const maxStatusCount = computed(() => Math.max(...statusDistribution.value.map((s) => s.count), 1));
</script>

<template>
  <div v-if="isLoading" class="space-y-4">
    <div class="grid grid-cols-4 gap-4">
      <Skeleton v-for="i in 4" :key="i" class="h-24 rounded-xl" />
    </div>
    <Skeleton class="h-48 rounded-xl" />
    <Skeleton class="h-48 rounded-xl" />
  </div>

  <div v-else-if="isError" class="py-12 text-center">
    <p class="font-mono text-sm text-destructive">Failed to load analytics data.</p>
    <Button variant="outline" class="mt-4" @click="refetch()">Retry</Button>
  </div>

  <div v-else>
    <!-- KPI cards -->
    <div class="grid grid-cols-2 gap-4 md:grid-cols-4">
      <Card class="shadow-elegant">
        <CardHeader class="flex flex-row items-center justify-between pb-2">
          <CardTitle class="font-mono text-xs uppercase tracking-widest text-muted-foreground">
            Total Shipments
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div class="font-mono text-3xl font-semibold">{{ kpis.total }}</div>
        </CardContent>
      </Card>
      <Card class="shadow-elegant">
        <CardHeader class="flex flex-row items-center justify-between pb-2">
          <CardTitle class="font-mono text-xs uppercase tracking-widest text-muted-foreground">
            On-Time Rate
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div class="font-mono text-3xl font-semibold text-success">{{ kpis.onTime }}%</div>
        </CardContent>
      </Card>
      <Card class="shadow-elegant">
        <CardHeader class="flex flex-row items-center justify-between pb-2">
          <CardTitle class="font-mono text-xs uppercase tracking-widest text-muted-foreground">
            Active Carriers
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div class="font-mono text-3xl font-semibold text-info">{{ kpis.activeCarriers }}</div>
        </CardContent>
      </Card>
      <Card class="shadow-elegant">
        <CardHeader class="flex flex-row items-center justify-between pb-2">
          <CardTitle class="font-mono text-xs uppercase tracking-widest text-muted-foreground">
            Avg Delivery
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div class="font-mono text-3xl font-semibold">{{ kpis.avgDeliveryTime }}</div>
        </CardContent>
      </Card>
    </div>

    <!-- Carrier Performance -->
    <div class="mt-8">
      <h3 class="font-mono text-sm font-semibold">Carrier Performance</h3>
      <div class="mt-4 space-y-3">
        <div
          v-for="c in carrierPerformance"
          :key="c.name"
          class="rounded-lg border border-border bg-card p-4"
        >
          <div class="flex items-center justify-between">
            <span class="font-mono text-sm">{{ c.name }}</span>
            <span class="font-mono text-xs text-muted-foreground"
              >{{ c.delivered }}/{{ c.total }} delivered</span
            >
          </div>
          <div class="mt-2 flex items-center gap-3">
            <div class="h-2 flex-1 overflow-hidden rounded-full bg-secondary">
              <div
                class="h-full rounded-full bg-gradient-accent transition-all"
                :style="{ width: `${(c.total / maxCarrierOrders) * 100}%` }"
              />
            </div>
            <span
              class="font-mono text-xs"
              :class="c.onTimeRate >= 80 ? 'text-success' : 'text-warning'"
            >
              {{ c.onTimeRate }}%
            </span>
          </div>
        </div>
      </div>
    </div>

    <!-- Status Distribution -->
    <div class="mt-8">
      <h3 class="font-mono text-sm font-semibold">Shipment Status Distribution</h3>
      <div class="mt-4 space-y-2">
        <div v-for="s in statusDistribution" :key="s.status" class="flex items-center gap-3">
          <span class="w-32 font-mono text-xs text-muted-foreground">{{ s.label }}</span>
          <div class="h-5 flex-1 overflow-hidden rounded-full bg-secondary">
            <div
              class="h-full rounded-full transition-all"
              :class="
                s.status === 'delivered'
                  ? 'bg-success'
                  : s.status === 'delayed'
                    ? 'bg-destructive'
                    : s.status === 'in_transit'
                      ? 'bg-info'
                      : s.status === 'out_for_delivery'
                        ? 'bg-primary'
                        : 'bg-muted-foreground/40'
              "
              :style="{ width: `${(s.count / maxStatusCount) * 100}%` }"
            />
          </div>
          <span class="w-16 text-right font-mono text-xs text-muted-foreground"
            >{{ s.count }} ({{ s.pct }}%)</span
          >
        </div>
      </div>
    </div>
  </div>
</template>
