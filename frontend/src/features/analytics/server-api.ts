import "server-only";
import { getApiOrigin } from "@/lib/api/server";
import { toAnalyticsData, type AnalyticsData } from "./data";

export async function getAnalyticsData(): Promise<AnalyticsData> {
  const origin = getApiOrigin();
  const [overviewResponse, seriesResponse] = await Promise.all([
    fetch(`${origin}/api/analytics/overview`, { cache: "no-store" }),
    fetch(`${origin}/api/analytics/timeseries`, { cache: "no-store" }),
  ]);
  if (!overviewResponse.ok || !seriesResponse.ok) throw new Error("Analytics request failed");
  const [overviewResult, seriesResult] = await Promise.all([
    overviewResponse.json(),
    seriesResponse.json(),
  ]);
  return toAnalyticsData(overviewResult, seriesResult);
}
