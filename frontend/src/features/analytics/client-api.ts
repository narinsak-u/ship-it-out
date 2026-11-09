import { toAnalyticsData, type AnalyticsData } from "./data";

export async function getAnalyticsData(signal: AbortSignal): Promise<AnalyticsData> {
  const [overviewResponse, seriesResponse] = await Promise.all([
    fetch("/api/analytics/overview", { signal }),
    fetch("/api/analytics/timeseries", { signal }),
  ]);
  if (!overviewResponse.ok || !seriesResponse.ok) throw new Error("Analytics request failed");
  const [overviewResult, seriesResult] = await Promise.all([
    overviewResponse.json(),
    seriesResponse.json(),
  ]);
  return toAnalyticsData(overviewResult, seriesResult);
}
