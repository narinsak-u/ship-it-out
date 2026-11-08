export interface AnalyticsOverview {
  total: number;
  active: number;
  delivered: number;
  by_status: { status: string; count: number }[];
  by_region: { name: string; total: number }[];
}

export interface AnalyticsSeries {
  by_month: { month: string; count: number }[];
  by_day_of_week: { day: string; count: number }[];
}

export interface AnalyticsData {
  overview: AnalyticsOverview;
  series: AnalyticsSeries;
}

interface AnalyticsEnvelope<T> {
  success: boolean;
  data: T;
}

type OverviewPayload = Omit<AnalyticsOverview, "by_status" | "by_region"> &
  Partial<Pick<AnalyticsOverview, "by_status" | "by_region">>;

export function toAnalyticsData(
  overviewResult: AnalyticsEnvelope<OverviewPayload>,
  seriesResult: AnalyticsEnvelope<Partial<AnalyticsSeries>>,
): AnalyticsData {
  if (!overviewResult.success || !seriesResult.success) throw new Error("Analytics request failed");
  return {
    overview: {
      ...overviewResult.data,
      by_status: overviewResult.data.by_status ?? [],
      by_region: overviewResult.data.by_region ?? [],
    },
    series: {
      by_month: seriesResult.data.by_month ?? [],
      by_day_of_week: seriesResult.data.by_day_of_week ?? [],
    },
  };
}
