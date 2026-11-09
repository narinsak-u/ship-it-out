import { getHomeAnalytics } from "@/features/analytics/server-api";
import { getShipments } from "./server-api";


export async function getHomeData() {
  const [analytics, shipments] = await Promise.allSettled([
    getHomeAnalytics(),
    getShipments("shipments?limit=3&exclude_status=delivered"),
  ]);

  return {
    analytics: analytics.status === "fulfilled" ? analytics.value : null,
    shipments: shipments.status === "fulfilled" ? shipments.value : null,
  };
}
