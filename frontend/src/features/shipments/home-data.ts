import type { Shipment } from "@/lib/api/preview";

interface ApiEnvelope<T> {
  success: boolean;
  data?: T;
}

interface AnalyticsOverview {
  total: number;
  active: number;
  delivered: number;
}

function apiBase() {
  const base = process.env.API_BASE_URL;
  if (!base) throw new Error("API_BASE_URL is required");
  return base.replace(/\/$/, "");
}

async function get<T>(path: string): Promise<T> {
  const response = await fetch(`${apiBase()}/api/${path}`, { cache: "no-store" });
  if (!response.ok) throw new Error(`Home API returned ${response.status}`);
  const envelope = (await response.json()) as ApiEnvelope<T>;
  if (!envelope.success || envelope.data === undefined)
    throw new Error("Home API response was invalid");
  return envelope.data;
}

export async function getHomeData() {
  const [analytics, shipments] = await Promise.allSettled([
    get<AnalyticsOverview>("analytics/overview"),
    get<Shipment[]>("shipments?limit=3&exclude_status=delivered"),
  ]);

  return {
    analytics: analytics.status === "fulfilled" ? analytics.value : null,
    shipments: shipments.status === "fulfilled" ? shipments.value : null,
  };
}
