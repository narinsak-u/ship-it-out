import "server-only";
import { getApiOrigin } from "@/lib/api/server";
import { mapShipment, mapShipments } from "./mappers";
import type { ApiEnvelope, Shipment, ShipmentDTO, ShipmentPage } from "./types";

async function get<T>(path: string): Promise<ApiEnvelope<T>> {
  const response = await fetch(`${getApiOrigin()}/api/${path}`, { cache: "no-store" });
  if (!response.ok) throw new Error(`Shipment API returned ${response.status}`);
  const result = (await response.json()) as ApiEnvelope<T>;
  if (!result.success) throw new Error(result.error || "Failed to load shipments");
  return result;
}

export async function getShipmentsPage(page = 1, limit = 10): Promise<ShipmentPage> {
  const result = await get<ShipmentDTO[]>(`shipments?page=${page}&limit=${limit}`);
  return { data: mapShipments(result.data), pagination: result.pagination! };
}

export async function getShipments(path: string): Promise<Shipment[]> {
  const result = await get<ShipmentDTO[]>(path);
  return mapShipments(result.data);
}

export async function getActiveShipments(): Promise<Shipment[]> {
  return getShipments("shipments?limit=-1&exclude_status=delivered");
}

export async function getShipment(id: string): Promise<Shipment | null> {
  const response = await fetch(
    `${getApiOrigin()}/api/shipments/${encodeURIComponent(id)}`,
    { cache: "no-store" },
  );
  if (response.status === 404) return null;
  if (!response.ok) throw new Error(`Shipment API returned ${response.status}`);
  const result = (await response.json()) as ApiEnvelope<ShipmentDTO>;
  if (!result.success) throw new Error(result.error || "Failed to load shipment");
  return mapShipment(result.data);
}

