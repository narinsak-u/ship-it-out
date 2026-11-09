import { mapShipments } from "./mappers";
import type { ShipmentDTO, ShipmentPage } from "./types";

interface Envelope<T> {
  success: boolean;
  data: T;
  pagination?: ShipmentPage["pagination"];
  error?: string;
}

export async function getManifestPage(
  page: number,
  search: string,
  status: string,
  signal?: AbortSignal,
): Promise<ShipmentPage> {
  const params = new URLSearchParams({ page: String(page), limit: "10" });
  if (search) params.set("search", search);
  if (status !== "all") params.set("status", status);
  const response = await fetch(`/api/shipments?${params}`, { cache: "no-store", signal });
  const result = (await response.json()) as Envelope<ShipmentDTO[]>;
  if (!response.ok || !result.success) {
    throw new Error(result.error || `Shipment API returned ${response.status}`);
  }
  return { data: mapShipments(result.data), pagination: result.pagination! };
}

export async function deleteShipment(id: string): Promise<void> {
  const response = await fetch(`/api/shipments/${encodeURIComponent(id)}`, {
    method: "DELETE",
    credentials: "same-origin",
    cache: "no-store",
  });
  const result = (await response.json()) as Envelope<unknown>;
  if (!response.ok || result.success === false) {
    throw new Error(result.error || `Delete failed (${response.status})`);
  }
}

export async function saveShipment(id: string | undefined, payload: object): Promise<void> {
  const response = await fetch(id ? `/api/shipments/${encodeURIComponent(id)}` : "/api/shipments", {
    method: id ? "PUT" : "POST",
    headers: { "content-type": "application/json" },
    body: JSON.stringify(payload),
  });
  const result = (await response.json()) as Envelope<unknown>;
  if (!response.ok || !result.success) {
    throw new Error(result.error || `Shipment save failed (${response.status})`);
  }
}
