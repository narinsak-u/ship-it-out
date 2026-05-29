<script setup lang="ts">
import { computed } from "vue";
import { useQuery } from "@tanstack/vue-query";
import { fetchAnalytics, fetchTimeSeries } from "@/lib/api/analytics";
import { analyticsKeys } from "@/lib/api/queryKeys";
import {
  computeKpis,
  computeRegionPerformance,
  computeStatusPieData,
  computeCumulativeData,
} from "@/lib/analytics-utils";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import Skeleton from "@/components/ui/Skeleton.vue";
import Button from "@/components/ui/Button.vue";
import ShipmentsAreaChart from "@/components/charts/ShipmentsAreaChart.vue";
import ShipmentsBarChart from "@/components/charts/ShipmentsBarChart.vue";
import ShipmentsLineChart from "@/components/charts/ShipmentsLineChart.vue";
import StatusPieChart from "@/components/charts/StatusPieChart.vue";

const {
  data: analytics,
  isLoading,
  isError,
  refetch,
} = useQuery({
  queryKey: analyticsKeys.all,
  queryFn: fetchAnalytics,
  staleTime: 5 * 60_000,
});

const { data: timeSeries } = useQuery({
  queryKey: analyticsKeys.timeseries(),
  queryFn: fetchTimeSeries,
  staleTime: 5 * 60_000,
});

const kpis = computed(() => computeKpis(analytics.value ?? null));
const regionPerformance = computed(() => computeRegionPerformance(analytics.value ?? null));
const statusPieData = computed(() => computeStatusPieData(analytics.value ?? null));
const cumulativeData = computed(() => computeCumulativeData(timeSeries.value ?? null));
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
