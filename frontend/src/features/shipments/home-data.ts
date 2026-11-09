import { getApiOrigin } from "@/lib/api/server";
import { getShipments } from "./server-api";
import type { AnalyticsOverview } from "@/features/analytics/data";

interface ApiEnvelope<T> {
  success: boolean;
  data?: T;
}

type HomeAnalyticsOverview = Pick<AnalyticsOverview, "total" | "active" | "delivered">;

async function get<T>(path: string): Promise<T> {
  const response = await fetch(`${getApiOrigin()}/api/${path}`, { cache: "no-store" });
  if (!response.ok) throw new Error(`Home API returned ${response.status}`);
  const envelope = (await response.json()) as ApiEnvelope<T>;
  if (!envelope.success || envelope.data === undefined)
    throw new Error("Home API response was invalid");
  return envelope.data;
}

export async function getHomeData() {
  const [analytics, shipments] = await Promise.allSettled([
    get<HomeAnalyticsOverview>("analytics/overview"),
    getShipments("shipments?limit=3&exclude_status=delivered"),
  ]);

  return {
    analytics: analytics.status === "fulfilled" ? analytics.value : null,
    shipments: shipments.status === "fulfilled" ? shipments.value : null,
  };
}
