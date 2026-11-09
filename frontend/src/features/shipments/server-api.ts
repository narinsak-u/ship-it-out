import "server-only";
import { getApiOrigin } from "@/lib/api/server";
import { mapShipment, mapShipments } from "./mappers";
import type { Shipment, ShipmentDTO, ShipmentPage } from "./types";

interface Envelope<T> {
  success: boolean;
  data: T;
  pagination?: ShipmentPage["pagination"];
  error?: string;
}

async function get<T>(path: string): Promise<Envelope<T>> {
  const response = await fetch(`${getApiOrigin()}/api/${path}`, { cache: "no-store" });
  if (!response.ok) throw new Error(`Shipment API returned ${response.status}`);
  const result = (await response.json()) as Envelope<T>;
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

export async function getShipment(id: string): Promise<Shipment | null> {
  const response = await fetch(
    `${getApiOrigin()}/api/shipments/${encodeURIComponent(id)}`,
    { cache: "no-store" },
  );
  if (response.status === 404) return null;
  if (!response.ok) throw new Error(`Shipment API returned ${response.status}`);
  const result = (await response.json()) as Envelope<ShipmentDTO>;
  if (!result.success) throw new Error(result.error || "Failed to load shipment");
  return mapShipment(result.data);
}

