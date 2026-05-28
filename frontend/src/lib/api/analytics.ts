import { api } from "@/lib/api/client";

export interface RegionCount {
  name: string;
  total: number;
}

export interface StatusCount {
  status: string;
  count: number;
}

export interface AnalyticsOverview {
  total: number;
  active: number;
  delivered: number;
  by_status: StatusCount[];
  by_region: RegionCount[];
}

export async function fetchAnalytics(): Promise<AnalyticsOverview> {
  const result = await api.get<AnalyticsOverview>("/analytics/overview");
  if (result.error) throw new Error(result.error);
  return result.data!;
}

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
