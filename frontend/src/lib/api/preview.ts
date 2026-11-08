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

const apiBase = () => {
  const base = process.env.API_BASE_URL;
  if (!base) throw new Error("API_BASE_URL is required");
  return base.replace(/\/$/, "");
};

export async function getShipment(orderId: string): Promise<Shipment | null> {
  const response = await fetch(`${apiBase()}/api/shipments/${encodeURIComponent(orderId)}`, {
    cache: "no-store",
  });
  if (response.status === 404) return null;
  if (!response.ok) throw new Error(`Shipment API returned ${response.status}`);
  const result = await response.json();
  return result.data as Shipment;
}

export async function getShipmentEvents(trackingNumber: string): Promise<ShipmentEvent[]> {
  const response = await fetch(`${apiBase()}/api/track/${encodeURIComponent(trackingNumber)}`, {
    cache: "no-store",
  });
  if (!response.ok) throw new Error(`Tracking API returned ${response.status}`);
  const result = await response.json();
  return (result.data?.events ?? []) as ShipmentEvent[];
}
