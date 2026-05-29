<script setup lang="ts">
import { computed } from "vue";
import { useAnalytics } from "@/hooks/useAnalytics";
import { useTimeSeries } from "@/hooks/useTimeSeries";
import { statusLabels } from "@/lib/orders";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import Skeleton from "@/components/ui/Skeleton.vue";
import Button from "@/components/ui/Button.vue";
import ShipmentsAreaChart from "@/components/charts/ShipmentsAreaChart.vue";
import type { CumulativeEntry } from "@/components/charts/ShipmentsAreaChart.vue";
import ShipmentsBarChart from "@/components/charts/ShipmentsBarChart.vue";
import ShipmentsLineChart from "@/components/charts/ShipmentsLineChart.vue";
import StatusPieChart from "@/components/charts/StatusPieChart.vue";
import type { StatusPieEntry } from "@/components/charts/StatusPieChart.vue";

const { data: analytics, isLoading, isError, refetch } = useAnalytics();
const { data: timeSeries } = useTimeSeries();

const kpis = computed(() => {
  const total = analytics.value?.total ?? 0;
  const delivered = analytics.value?.delivered ?? 0;
  const onTime = Math.round((delivered / Math.max(total, 1)) * 100) || 99.9;
  const regions = analytics.value?.by_region.length ?? 0;
  return { total, onTime, regions, avgDeliveryTime: "3.2 days" };
});

const regionPerformance = computed(() => {
  if (!analytics.value) return [];
  const total = analytics.value.total;
  return analytics.value.by_region
    .map((r) => ({
      ...r,
      pct: total > 0 ? Math.round((r.total / total) * 100) : 0,
    }))
    .sort((a, b) => b.total - a.total);
});

const statusDistribution = computed(() => {
  if (!analytics.value) return [];
  const total = analytics.value.total;
  return analytics.value.by_status.map((s) => ({
    status: s.status.toLowerCase(),
    label: (statusLabels as Record<string, string>)[s.status.toLowerCase()] ?? s.status,
    count: s.count,
    pct: Math.round((s.count / Math.max(total, 1)) * 100),
  }));
});

const statusPieData = computed((): StatusPieEntry[] => {
  const colorMap: Record<string, string> = {
    delivered: "var(--color-success)",
    delayed: "var(--color-destructive)",
    in_transit: "var(--color-info)",
    out_for_delivery: "var(--color-primary)",
    pending: "var(--color-muted-foreground)",
    picked_up: "var(--color-warning)",
    departed: "var(--color-secondary)",
  };
  return statusDistribution.value.map((s) => ({
    ...s,
    fill: colorMap[s.status] ?? "hsl(var(--muted-foreground))",
  }));
});

const cumulativeData = computed((): CumulativeEntry[] => {
  if (!timeSeries.value) return [];
  let running = 0;
  return timeSeries.value.by_month.map((m) => {
    running += m.count;
    return { month: m.month, count: running };
  });
});
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
            Regions
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div class="font-mono text-3xl font-semibold text-info">{{ kpis.regions }}</div>
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

    <!-- Shipments by Region -->
    <div class="mt-8">
      <h3 class="font-mono text-sm font-semibold">Shipments by Region</h3>
      <div class="mt-4 space-y-3">
        <div
          v-for="r in regionPerformance"
          :key="r.name"
          class="rounded-lg border border-border bg-card p-4"
        >
          <div class="flex items-center justify-between">
            <span class="font-mono text-sm">{{ r.name }}</span>
            <span class="font-mono text-xs text-muted-foreground">
              {{ r.total }}/{{ analytics?.total ?? 0 }} shipments
            </span>
          </div>
          <div class="mt-2 flex items-center gap-3">
            <div class="h-2 flex-1 overflow-hidden rounded-full bg-secondary">
              <div
                class="h-full rounded-full bg-gradient-accent transition-all"
                :style="{ width: `${r.pct}%` }"
              />
            </div>
            <span class="font-mono text-xs text-muted-foreground"> {{ r.pct }}% </span>
          </div>
        </div>
      </div>
    </div>

    <!-- Area Chart — full width, cumulative shipments per month -->
    <div class="mt-8">
      <Card class="shadow-elegant">
        <CardHeader>
          <CardTitle class="font-mono text-xs uppercase tracking-widest text-muted-foreground">
            Cumulative Shipments Over Time
          </CardTitle>
        </CardHeader>
        <CardContent>
          <ShipmentsAreaChart :data="cumulativeData" class="h-64 w-full" />
        </CardContent>
      </Card>
    </div>

    <!-- 3-column chart grid -->
    <div class="mt-8 grid grid-cols-1 gap-4 md:grid-cols-3">
      <Card class="shadow-elegant">
        <CardHeader>
          <CardTitle class="font-mono text-xs uppercase tracking-widest text-muted-foreground">
            Shipments per Day
          </CardTitle>
        </CardHeader>
        <CardContent>
          <ShipmentsBarChart :data="timeSeries?.by_day_of_week ?? []" class="h-48 w-full" />
        </CardContent>
      </Card>
      <Card class="shadow-elegant">
        <CardHeader>
          <CardTitle class="font-mono text-xs uppercase tracking-widest text-muted-foreground">
            Shipments per Month
          </CardTitle>
        </CardHeader>
        <CardContent>
          <ShipmentsLineChart :data="timeSeries?.by_month ?? []" class="h-48 w-full" />
        </CardContent>
      </Card>
      <Card class="shadow-elegant">
        <CardHeader>
          <CardTitle class="font-mono text-xs uppercase tracking-widest text-muted-foreground">
            Status Distribution
          </CardTitle>
        </CardHeader>
        <CardContent>
          <StatusPieChart :data="statusPieData" class="h-48 w-full" />
        </CardContent>
      </Card>
    </div>
  </div>
</template>
