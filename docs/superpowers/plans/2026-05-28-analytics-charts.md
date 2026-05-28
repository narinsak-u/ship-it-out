# Analytics Charts Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add 4 Unovis-based charts (area, bar, line, pie) to the AnalyticsPanel, powered by a new backend time series endpoint.

**Architecture:** Backend `/analytics/timeseries` endpoint aggregates `created_at` dates into monthly and day-of-week counts. Frontend installs `@unovis/vue` + shadcn-vue chart components, creates 4 chart components in `frontend/src/components/charts/`, and integrates them into `AnalyticsPanel.vue` with computed helpers.

**Tech Stack:** Go (Fiber, GORM), Vue 3 (Composition API), shadcn-vue chart wrappers, Unovis

---

### Task 1: Install Unovis + shadcn-vue chart components

**Files:**
- Modify: `frontend/package.json`
- (no Create — `bun install` + shadcn CLI)

- [ ] **Step 1: Add shadcn-vue chart component**

```bash
bunx shadcn-vue@latest add chart
```

This installs:
- `ChartContainer.vue`, `ChartTooltip.vue`, `ChartTooltipContent.vue`, `ChartLegendContent.vue`, `ChartCrosshair.vue`, `componentToString.ts` under `frontend/src/components/ui/chart/`
- `@unovis/vue` as a dependency in `package.json`

- [ ] **Step 2: Verify build still passes**

```bash
npm run build
```

Expected: Build succeeds (no type errors).

- [ ] **Step 3: Commit**

```bash
git add -A && git commit -m "feat: add shadcn-vue chart components (Unovis)"
```

---

### Task 2: Backend — add `/analytics/timeseries` endpoint

**Files:**
- Modify: `backend/internal/analytics/handler.go`
- Modify: `backend/cmd/server/main.go`

- [ ] **Step 1: Add time series types and handler to handler.go**

Append after the `Overview` function in `backend/internal/analytics/handler.go`:

```go
type monthCount struct {
	Month string `json:"month"`
	Count int64  `json:"count"`
}

type dayCount struct {
	Day   string `json:"day"`
	Count int64  `json:"count"`
}

func TimeSeries(c *fiber.Ctx) error {
	var byMonth []monthCount
	database.DB.Model(&models.Shipment{}).
		Select("to_char(created_at, 'YYYY-MM') as month, count(*) as count").
		Group("month").
		Order("month").
		Scan(&byMonth)

	var byDay []dayCount
	database.DB.Model(&models.Shipment{}).
		Select("trim(to_char(created_at, 'Day')) as day, count(*) as count").
		Group("day").
		Order("min(extract(dow from created_at))").
		Scan(&byDay)

	return utils.Success(c, fiber.Map{
		"by_month":        byMonth,
		"by_day_of_week":  byDay,
	})
}
```

- [ ] **Step 2: Add route in main.go**

In `backend/cmd/server/main.go`, after the analytics overview route:

```go
api.Get("/analytics/timeseries", middleware.AuthRequired(), analytics.TimeSeries)
```

- [ ] **Step 3: Verify Go code compiles**

```bash
cd backend && go build ./...
```

Expected: No errors.

- [ ] **Step 4: Commit**

```bash
git add -A && git commit -m "feat(analytics): add /analytics/timeseries endpoint"
```

---

### Task 3: Frontend — add API types and hook

**Files:**
- Modify: `frontend/src/lib/api/analytics.ts`
- Create: `frontend/src/hooks/useTimeSeries.ts`

- [ ] **Step 1: Add types and fetch function to analytics.ts**

Append to `frontend/src/lib/api/analytics.ts`:

```ts
export interface MonthlyCount {
  month: string;
  count: number;
}

export interface DayOfWeekCount {
  day: string;
  count: number;
}

export interface TimeSeriesData {
  by_month: MonthlyCount[];
  by_day_of_week: DayOfWeekCount[];
}

export async function fetchTimeSeries(): Promise<TimeSeriesData> {
  const result = await api.get<TimeSeriesData>("/analytics/timeseries");
  if (result.error) throw new Error(result.error);
  return result.data!;
}
```

- [ ] **Step 2: Create useTimeSeries hook**

New file `frontend/src/hooks/useTimeSeries.ts`:

```ts
import { useQuery } from "@tanstack/vue-query";
import { fetchTimeSeries } from "@/lib/api/analytics";

export function useTimeSeries() {
  return useQuery({
    queryKey: ["analytics", "timeseries"],
    queryFn: fetchTimeSeries,
  });
}
```

- [ ] **Step 3: Verify build passes**

```bash
npm run build
```

Expected: Build succeeds.

- [ ] **Step 4: Commit**

```bash
git add -A && git commit -m "feat(analytics): add useTimeSeries hook and API client"
```

---

### Task 4: Create ShipmentsAreaChart component

**Files:**
- Create: `frontend/src/components/charts/ShipmentsAreaChart.vue`

- [ ] **Step 1: Create the component**

```vue
<script setup lang="ts">
import type { HTMLAttributes } from "vue";
import { VisArea, VisXYContainer, VisAxis } from "@unovis/vue";
import type { ChartConfig } from "@/components/ui/chart";
import {
  ChartContainer,
  ChartTooltip,
  ChartCrosshair,
  ChartTooltipContent,
  componentToString,
} from "@/components/ui/chart";

export interface CumulativeEntry {
  month: string;
  count: number;
}

const props = defineProps<{
  data: CumulativeEntry[];
  class?: HTMLAttributes["class"];
}>();

type Data = CumulativeEntry;

const chartConfig = {
  count: {
    label: "Total Shipments",
    color: "hsl(var(--primary))",
  },
} satisfies ChartConfig;
</script>

<template>
  <div v-if="data.length === 0" class="flex h-64 items-center justify-center text-sm text-muted-foreground">
    No data
  </div>
  <ChartContainer
    v-else
    :config="chartConfig"
    :class="props.class"
  >
    <VisXYContainer :data="data">
      <VisArea
        :x="(d: Data) => d.month"
        :y="(d: Data) => d.count"
        :color="[chartConfig.count.color]"
      />
      <VisAxis
        type="x"
        :x="(d: Data) => d.month"
        :tick-line="false"
        :domain-line="false"
        :grid-line="false"
        :tick-format="(v: string) => {
          const date = new Date(v + '-01');
          return date.toLocaleDateString('en-US', { month: 'short', year: '2-digit' });
        }"
        :tick-values="data.map(d => d.month)"
      />
      <VisAxis
        type="y"
        :tick-line="false"
        :domain-line="false"
        :grid-line="true"
      />
      <ChartTooltip />
      <ChartCrosshair
        :template="componentToString(chartConfig, ChartTooltipContent)"
        :color="[chartConfig.count.color]"
      />
    </VisXYContainer>
  </ChartContainer>
</template>
```

- [ ] **Step 2: Verify build passes**

```bash
npm run build
```

Expected: Build succeeds.

- [ ] **Step 3: Commit**

```bash
git add -A && git commit -m "feat(analytics): add ShipmentsAreaChart component"
```

---

### Task 5: Create ShipmentsBarChart component

**Files:**
- Create: `frontend/src/components/charts/ShipmentsBarChart.vue`

- [ ] **Step 1: Create the component**

```vue
<script setup lang="ts">
import type { HTMLAttributes } from "vue";
import { VisGroupedBar, VisXYContainer, VisAxis } from "@unovis/vue";
import type { ChartConfig } from "@/components/ui/chart";
import {
  ChartContainer,
  ChartTooltip,
  ChartCrosshair,
  ChartTooltipContent,
  componentToString,
} from "@/components/ui/chart";
import type { DayOfWeekCount } from "@/lib/api/analytics";

const props = defineProps<{
  data: DayOfWeekCount[];
  class?: HTMLAttributes["class"];
}>();

type Data = DayOfWeekCount;

const chartConfig = {
  count: {
    label: "Shipments",
    color: "hsl(var(--primary))",
  },
} satisfies ChartConfig;
</script>

<template>
  <div v-if="data.length === 0" class="flex h-48 items-center justify-center text-sm text-muted-foreground">
    No data
  </div>
  <ChartContainer
    v-else
    :config="chartConfig"
    :class="props.class"
  >
    <VisXYContainer :data="data">
      <VisGroupedBar
        :x="(d: Data) => d.day"
        :y="(d: Data) => d.count"
        :color="[chartConfig.count.color]"
        :rounded-corners="4"
        bar-padding="0.1"
        group-padding="0"
      />
      <VisAxis
        type="x"
        :x="(d: Data) => d.day"
        :tick-line="false"
        :domain-line="false"
        :grid-line="false"
        :tick-format="(v: string) => v.slice(0, 3)"
        :tick-values="data.map(d => d.day)"
      />
      <VisAxis
        type="y"
        :tick-line="false"
        :domain-line="false"
        :grid-line="true"
      />
      <ChartTooltip />
      <ChartCrosshair
        :template="componentToString(chartConfig, ChartTooltipContent)"
        :color="[chartConfig.count.color]"
      />
    </VisXYContainer>
  </ChartContainer>
</template>
```

- [ ] **Step 2: Verify build passes**

```bash
npm run build
```

Expected: Build succeeds.

- [ ] **Step 3: Commit**

```bash
git add -A && git commit -m "feat(analytics): add ShipmentsBarChart component"
```

---

### Task 6: Create ShipmentsLineChart component

**Files:**
- Create: `frontend/src/components/charts/ShipmentsLineChart.vue`

- [ ] **Step 1: Create the component**

```vue
<script setup lang="ts">
import type { HTMLAttributes } from "vue";
import { VisLine, VisXYContainer, VisAxis } from "@unovis/vue";
import type { ChartConfig } from "@/components/ui/chart";
import {
  ChartContainer,
  ChartTooltip,
  ChartCrosshair,
  ChartTooltipContent,
  componentToString,
} from "@/components/ui/chart";
import type { MonthlyCount } from "@/lib/api/analytics";

const props = defineProps<{
  data: MonthlyCount[];
  class?: HTMLAttributes["class"];
}>();

type Data = MonthlyCount;

const chartConfig = {
  count: {
    label: "Shipments",
    color: "hsl(var(--primary))",
  },
} satisfies ChartConfig;
</script>

<template>
  <div v-if="data.length === 0" class="flex h-48 items-center justify-center text-sm text-muted-foreground">
    No data
  </div>
  <ChartContainer
    v-else
    :config="chartConfig"
    :class="props.class"
  >
    <VisXYContainer :data="data">
      <VisLine
        :x="(d: Data) => d.month"
        :y="(d: Data) => d.count"
        :color="[chartConfig.count.color]"
      />
      <VisAxis
        type="x"
        :x="(d: Data) => d.month"
        :tick-line="false"
        :domain-line="false"
        :grid-line="false"
        :tick-format="(v: string) => {
          const date = new Date(v + '-01');
          return date.toLocaleDateString('en-US', { month: 'short', year: '2-digit' });
        }"
        :tick-values="data.map(d => d.month)"
      />
      <VisAxis
        type="y"
        :tick-line="false"
        :domain-line="false"
        :grid-line="true"
      />
      <ChartTooltip />
      <ChartCrosshair
        :template="componentToString(chartConfig, ChartTooltipContent)"
        :color="[chartConfig.count.color]"
      />
    </VisXYContainer>
  </ChartContainer>
</template>
```

- [ ] **Step 2: Verify build passes**

```bash
npm run build
```

Expected: Build succeeds.

- [ ] **Step 3: Commit**

```bash
git add -A && git commit -m "feat(analytics): add ShipmentsLineChart component"
```

---

### Task 7: Create StatusPieChart component

**Files:**
- Create: `frontend/src/components/charts/StatusPieChart.vue`

- [ ] **Step 1: Create the component**

```vue
<script setup lang="ts">
import type { HTMLAttributes } from "vue";
import { computed } from "vue";
import { VisDonut, VisSingleContainer } from "@unovis/vue";
import type { ChartConfig } from "@/components/ui/chart";
import {
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
  componentToString,
} from "@/components/ui/chart";

export interface StatusPieEntry {
  status: string;
  label: string;
  count: number;
  pct: number;
  fill: string;
}

const props = defineProps<{
  data: StatusPieEntry[];
  class?: HTMLAttributes["class"];
}>();

type Data = StatusPieEntry;

const chartConfig = computed<ChartConfig>(() => {
  const config: ChartConfig = {};
  for (const entry of props.data) {
    config[entry.status] = { label: entry.label, color: entry.fill };
  }
  return config;
});
</script>

<template>
  <div v-if="data.length === 0" class="flex h-48 items-center justify-center text-sm text-muted-foreground">
    No data
  </div>
  <ChartContainer
    v-else
    :config="chartConfig"
    :class="props.class"
  >
    <VisSingleContainer :data="data">
      <VisDonut
        :value="(d: Data) => d.count"
        :label="(d: Data) => d.label"
        :color="(d: Data) => d.fill"
      />
      <ChartTooltip
        :template="componentToString(chartConfig, ChartTooltipContent, {
          labelKey: 'count',
          nameKey: 'status',
        })"
      />
    </VisSingleContainer>
  </ChartContainer>
</template>
```

- [ ] **Step 2: Verify build passes**

```bash
npm run build
```

Expected: Build succeeds.

- [ ] **Step 3: Commit**

```bash
git add -A && git commit -m "feat(analytics): add StatusPieChart component"
```

---

### Task 8: Integrate charts into AnalyticsPanel.vue

**Files:**
- Modify: `frontend/src/components/AnalyticsPanel.vue`

- [ ] **Step 1: Update script section**

Replace the existing `<script setup>` block with:

```vue
<script setup lang="ts">
import { computed } from "vue";
import { useCarriers } from "@/hooks/useCarriers";
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

const { data: carriersData } = useCarriers();
const { data: analytics, isLoading, isError, refetch } = useAnalytics();
const { data: timeSeries } = useTimeSeries();

const kpis = computed(() => {
  const total = analytics.value?.total ?? 0;
  const delivered = analytics.value?.delivered ?? 0;
  const onTime = Math.round((delivered / Math.max(total, 1)) * 100);
  const activeCarriers = carriersData.value?.filter((c) => c.status === "active").length ?? 0;
  return { total, onTime, activeCarriers, avgDeliveryTime: "3.2 days" };
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
    delivered: "hsl(var(--success))",
    delayed: "hsl(var(--destructive))",
    in_transit: "hsl(var(--info))",
    out_for_delivery: "hsl(var(--primary))",
    pending: "hsl(var(--muted-foreground))",
    picked_up: "hsl(var(--warning))",
    departed: "hsl(var(--secondary))",
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
```

- [ ] **Step 2: Update template section**

Replace the full template section after the `<!-- Shipments by Region -->` block. The old "TODO" comments and "Shipment Status Distribution" section should be replaced with:

```vue
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
```

Specific changes:
- Remove the existing `maxRegionOrders` and `maxStatusCount` computed properties (no longer used).
- Remove the existing "Shipment Status Distribution" section (lines 143-171 equivalent).
- Remove the TODO comment lines.
- Add the imports and code above.

- [ ] **Step 3: Verify build passes**

```bash
npm run build
```

Expected: Build succeeds.

- [ ] **Step 4: Commit**

```bash
git add -A && git commit -m "feat(analytics): integrate charts into AnalyticsPanel"
```
