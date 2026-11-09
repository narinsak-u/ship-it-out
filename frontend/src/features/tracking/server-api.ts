import "server-only";
import { getApiOrigin } from "@/lib/api/server";
import type { TrackingEvent } from "./types";

export async function getTrackingEvents(trackingNumber: string): Promise<TrackingEvent[]> {
  const response = await fetch(
    `${getApiOrigin()}/api/track/${encodeURIComponent(trackingNumber)}`,
    { cache: "no-store" },
  );
  if (!response.ok) throw new Error(`Tracking API returned ${response.status}`);
  const result = (await response.json()) as { data?: { events?: TrackingEvent[] } };
  return result.data?.events ?? [];
}
