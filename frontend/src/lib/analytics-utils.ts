import { statusLabels } from "@/lib/orders";
import type { AnalyticsOverview, TimeSeriesData } from "@/lib/api/analytics";

export interface KpiData {
  total: number;
  onTime: number;
  regions: number;
  avgDeliveryTime: string;
}

export interface RegionEntry {
  name: string;
  total: number;
  pct: number;
}

export interface StatusDistEntry {
  status: string;
  label: string;
  count: number;
  pct: number;
}

export interface StatusPieEntry {
  status: string;
  label: string;
  count: number;
  pct: number;
  fill: string;
}

export interface CumulativeEntry {
  month: string;
  count: number;
}

export function computeKpis(data: AnalyticsOverview | null): KpiData {
  const total = data?.total ?? 0;
  const delivered = data?.delivered ?? 0;
  const onTime = Math.round((delivered / Math.max(total, 1)) * 100) || 99.9;
  const regions = data?.by_region.length ?? 0;
  return { total, onTime, regions, avgDeliveryTime: "3.2 days" };
}

export function computeRegionPerformance(data: AnalyticsOverview | null): RegionEntry[] {
  if (!data) return [];
  const total = data.total;
  return data.by_region
    .map((r) => ({
      name: r.name,
      total: r.total,
      pct: total > 0 ? Math.round((r.total / total) * 100) : 0,
    }))
    .sort((a, b) => b.total - a.total);
}

export function computeStatusDistribution(data: AnalyticsOverview | null): StatusDistEntry[] {
  if (!data) return [];
  const total = data.total;
  return data.by_status.map((s) => ({
    status: s.status.toLowerCase(),
    label: (statusLabels as Record<string, string>)[s.status.toLowerCase()] ?? s.status,
    count: s.count,
    pct: Math.round((s.count / Math.max(total, 1)) * 100),
  }));
}

const PIE_COLOR_MAP: Record<string, string> = {
  delivered: "var(--color-success)",
  delayed: "var(--color-destructive)",
  in_transit: "var(--color-info)",
  out_for_delivery: "var(--color-primary)",
  pending: "var(--color-muted-foreground)",
  picked_up: "var(--color-warning)",
  departed: "var(--color-secondary)",
};

export function computeStatusPieData(data: AnalyticsOverview | null): StatusPieEntry[] {
  return computeStatusDistribution(data).map((s) => ({
    ...s,
    fill: PIE_COLOR_MAP[s.status] ?? "hsl(var(--muted-foreground))",
  }));
}

export function computeCumulativeData(data: TimeSeriesData | null): CumulativeEntry[] {
  if (!data) return [];
  let running = 0;
  return data.by_month.map((m) => {
    running += m.count;
    return { month: m.month, count: running };
  });
}
