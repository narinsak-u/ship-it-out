# Analytics Charts — Design Spec

## Overview

Add 4 Recharts-powered charts to the AnalyticsPanel: an area chart (cumulative shipments per month, full width), a bar chart (shipments per day of week), a line chart (shipments per month count), and a pie chart (status distribution replacing the current horizontal bars).

## Backend — New `/analytics/timeseries` Endpoint

**File:** `backend/internal/analytics/handler.go` — add `TimeSeries` handler.

**Route:** `GET /analytics/timeseries`

**Response shape:**

```json
{
  "by_month": [
    { "month": "2025-12", "count": 4 },
    { "month": "2026-01", "count": 7 }
  ],
  "by_day_of_week": [
    { "day": "Monday", "count": 5 },
    { "day": "Tuesday", "count": 8 }
  ]
}
```

**SQL queries:**
- `by_month`: `SELECT to_char(created_at, 'YYYY-MM') as month, COUNT(*) as count FROM shipments GROUP BY month ORDER BY month`
- `by_day_of_week`: `SELECT TRIM(to_char(created_at, 'Day')) as day, COUNT(*) as count FROM shipments GROUP BY day ORDER BY MIN(EXTRACT(DOW FROM created_at))`
  - `TRIM()` removes Postgres padding to 9 chars; PostgreSQL locale set to English.

**Route registration:** Add alongside existing `GET /analytics/overview` in router.

## Frontend — Dependencies

- Add `recharts` npm package.
- Run `bunx shadcn-vue@latest add chart` to install shadcn-vue chart wrapper components (`ChartContainer`, `ChartTooltip`, `ChartLegend`, etc.).

## Frontend — Data Layer

### API client (`frontend/src/lib/api/analytics.ts`)

Add types:

```ts
export interface MonthlyCount {
  month: string;  // "2026-05"
  count: number;
}
export interface DayOfWeekCount {
  day: string;    // "Monday"
  count: number;
}
export interface TimeSeriesData {
  by_month: MonthlyCount[];
  by_day_of_week: DayOfWeekCount[];
}

export async function fetchTimeSeries(): Promise<TimeSeriesData>;
```

### Vue Query hook (`frontend/src/hooks/useTimeSeries.ts`)

```ts
export function useTimeSeries() {
  return useQuery({ queryKey: ["analytics", "timeseries"], queryFn: fetchTimeSeries });
}
```

### Computed helpers (in AnalyticsPanel.vue)

- `cumulativeData: CumulativeEntry[]` — walk `by_month`, accumulate running total for area chart
  - `interface CumulativeEntry { month: string; count: number; }`
- `monthlyData: MonthlyCount[]` — pass `by_month` directly to line chart
- `dailyData: DayOfWeekCount[]` — pass `by_day_of_week` directly to bar chart
- `statusPieData: StatusPieEntry[]` — reformat existing `statusDistribution` for pie chart (add `fill` color per status)
  - `interface StatusPieEntry { status: string; label: string; count: number; pct: number; fill: string; }`

## Frontend — Components

4 new chart components in `frontend/src/components/charts/`:

| Component | Chart Type | Data | Key Props |
|---|---|---|---|
| `ShipmentsAreaChart.vue` | AreaChart (filled) | `cumulativeData` | `data: CumulativeEntry[]`, `class?` |
| `ShipmentsBarChart.vue` | BarChart | `dailyData` | `data: DayOfWeekCount[]`, `class?` |
| `ShipmentsLineChart.vue` | LineChart | `monthlyData` | `data: MonthlyCount[]`, `class?` |
| `StatusPieChart.vue` | PieChart | `statusPieData` | `data: StatusPieEntry[]`, `class?` |

Each component:
- Uses shadcn-vue `<ChartContainer>`, `<ChartTooltip>`, `<ChartLegend>`
- Receives typed `data` prop + optional `class` prop typed as `HTMLAttributes['class']`
- Renders empty state if data is empty
- No Card wrapper (caller provides layout)

### AnalyticsPanel.vue layout updates

```
[KPI cards - 4 columns]                     ← unchanged
[Shipments by Region]                       ← unchanged
[ShipmentsAreaChart - full width Card]      ← new
[ShipmentsBarChart | ShipmentsLineChart | StatusPieChart] ← 3-column grid
```

- Remove the existing "Shipment Status Distribution" horizontal bar section.
- Add `useTimeSeries()` alongside existing `useAnalytics()` / `useCarriers()`.
- Loading/error states already handled by existing skeleton/error block.

### StatusPieChart color mapping

| Status | Color |
|---|---|
| delivered | `hsl(var(--success))` |
| delayed | `hsl(var(--destructive))` |
| in_transit | `hsl(var(--info))` |
| out_for_delivery | `hsl(var(--primary))` |
| pending | `hsl(var(--muted-foreground))` |
| picked_up | `hsl(var(--warning))` |
| departed | `hsl(var(--secondary))` |

### Empty state

If any data array is empty, the chart component renders a centered "No data" message in `text-muted-foreground`. No separate empty-state skeleton needed.

## Files changed

**New files:**
- `frontend/src/hooks/useTimeSeries.ts`
- `frontend/src/components/charts/ShipmentsAreaChart.vue`
- `frontend/src/components/charts/ShipmentsBarChart.vue`
- `frontend/src/components/charts/ShipmentsLineChart.vue`
- `frontend/src/components/charts/StatusPieChart.vue`

**Modified files:**
- `frontend/package.json` (add `recharts`)
- `frontend/src/lib/api/analytics.ts` (add `TimeSeriesData` + `fetchTimeSeries`)
- `frontend/src/components/AnalyticsPanel.vue` (integrate charts, remove status bars)
- `backend/internal/analytics/handler.go` (add `TimeSeries` handler)
- Backend router (register `/analytics/timeseries` route)

## Not in scope

- No new loading/error states — inherits existing skeleton/error from AnalyticsPanel.
- No test framework — no tests added (none exist in project).
- No animation configuration — Recharts defaults suffice.
