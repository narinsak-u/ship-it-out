import { getApiOrigin } from "@/lib/api/server";

export interface ShipmentContact {
  name: string;
  zipcode: string;
  subDistrict: string;
  district: string;
  province: string;
  coords: { lat: number; lng: number };
}

export interface Shipment {
  id: string;
  trackingNumber: string;
  customer: ShipmentContact;
  receiver: ShipmentContact;
  currentCoords: { lat: number; lng: number };
  origin: string;
  destination: string;
  status: string;
  carrier: string;
  weight: number;
  items: number;
  estimatedDelivery: string;
  createdAt: string;
  progress: number;
}

export interface ShipmentEvent {
  status: string;
  location: { name: string; lat: number; lng: number };
  description?: string;
  timestamp: string;
}

export interface ShipmentPage {
  data: Shipment[];
  pagination: { page: number; limit: number; total: number; totalPages: number };
}

export async function getShipmentsPage(page = 1, limit = 10): Promise<ShipmentPage> {
  const response = await fetch(`${getApiOrigin()}/api/shipments?page=${page}&limit=${limit}`, {
    cache: "no-store",
  });
  if (!response.ok) throw new Error(`Shipment API returned ${response.status}`);
  const result = await response.json();
  if (!result.success) throw new Error(result.error || "Failed to load shipments");
  return { data: result.data, pagination: result.pagination };
}

export async function getShipment(orderId: string): Promise<Shipment | null> {
  const response = await fetch(`${getApiOrigin()}/api/shipments/${encodeURIComponent(orderId)}`, {
    cache: "no-store",
  });
  if (response.status === 404) return null;
  if (!response.ok) throw new Error(`Shipment API returned ${response.status}`);
  const result = await response.json();
  return result.data as Shipment;
}

export async function getShipmentEvents(trackingNumber: string): Promise<ShipmentEvent[]> {
  const response = await fetch(
    `${getApiOrigin()}/api/track/${encodeURIComponent(trackingNumber)}`,
    {
      cache: "no-store",
    },
  );
  if (!response.ok) throw new Error(`Tracking API returned ${response.status}`);
  const result = await response.json();
  return (result.data?.events ?? []) as ShipmentEvent[];
}
