import "server-only";
import { getApiOrigin } from "@/lib/api/server";
import { toAnalyticsData, type AnalyticsData } from "./data";


export async function getHomeAnalytics(): Promise<
  Pick<AnalyticsData["overview"], "total" | "active" | "delivered">
> {
  const response = await fetch(`${getApiOrigin()}/api/analytics/overview`, { cache: "no-store" });
  if (!response.ok) throw new Error("Analytics request failed");
  const result = (await response.json()) as {
    success: boolean;
    data?: AnalyticsData["overview"];
  };
  if (!result.success || !result.data) throw new Error("Analytics request failed");
  const { total, active, delivered } = result.data;
  return { total, active, delivered };
}
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
