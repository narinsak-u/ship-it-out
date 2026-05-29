import { api } from "@/lib/api/client";

export interface TrackResult {
  shipment: {
    id: string;
  };
}

export async function trackShipment(trackingNumber: string): Promise<TrackResult> {
  const result = await api.get<TrackResult>(`/track/${trackingNumber}`);
  if (result.error) throw new Error(result.error);
  return result.data!;
}
