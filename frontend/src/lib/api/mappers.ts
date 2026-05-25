import type { Order, TrackingEvent, ShipmentStatus, GeoPoint, ContactInfo } from "@/lib/orders";
import type { Hub, HubStatus } from "@/lib/carriers";

// Backend response types (mirrors Go struct JSON tags)

export interface BackendShipment {
  id: string;
  trackingNumber: string;
  customer: ContactInfo;
  receiver: ContactInfo;
  currentCoords: GeoPoint;
  origin: string;
  destination: string;
  status: string;
  carrier: string;
  weight: string;
  items: number;
  estimatedDelivery: string;
  createdAt: string;
  progress: number;
}

export interface BackendShipmentEvent {
  id: number;
  shipmentId: number;
  status: string;
  location: { name: string; lat: number; lng: number };
  description?: string;
  timestamp: string;
}

export interface BackendHub {
  id: string;
  name: string;
  carrierId: string;
  address: string;
  coords: GeoPoint;
  capacity: number;
  currentUtilization: number;
  status: string;
  createdAt: string;
}

// Date formatters

const MONTHS = ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"];

export function formatDate(iso: string): string {
  const d = new Date(iso);
  return `${MONTHS[d.getMonth()]} ${d.getDate()}, ${d.getFullYear()}`;
}

export function formatTimestamp(iso: string): string {
  const d = new Date(iso);
  const hh = String(d.getHours()).padStart(2, "0");
  const mm = String(d.getMinutes()).padStart(2, "0");
  return `${MONTHS[d.getMonth()]} ${d.getDate()}, ${hh}:${mm}`;
}

// Mappers

export function mapShipmentToOrder(s: BackendShipment): Order {
  return {
    id: s.id,
    trackingNumber: s.trackingNumber,
    customer: s.customer,
    receiver: s.receiver,
    origin: s.origin,
    destination: s.destination,
    currentCoords: s.currentCoords,
    status: s.status as ShipmentStatus,
    carrier: s.carrier,
    weight: s.weight,
    items: s.items,
    estimatedDelivery: formatDate(s.estimatedDelivery),
    createdAt: formatDate(s.createdAt),
    progress: s.progress,
    events: [],
  };
}

export function mapEventToTrackingEvent(e: BackendShipmentEvent): TrackingEvent {
  return {
    timestamp: formatTimestamp(e.timestamp),
    location: e.location,
    status: e.status,
    description: e.description ?? "",
  };
}

export function mapBackendHubToHub(h: BackendHub): Hub {
  return {
    id: h.id,
    name: h.name,
    carrierId: h.carrierId,
    address: h.address,
    coords: h.coords,
    capacity: h.capacity,
    currentUtilization: h.currentUtilization,
    status: h.status as HubStatus,
  };
}
