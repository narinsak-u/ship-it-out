import type { GeoPoint } from "@/lib/geo";

export type ShipmentStatus =
  | "pending"
  | "picked_up"
  | "departed"
  | "in_transit"
  | "out_for_delivery"
  | "delivered"
  | "delayed";

export const statusLabels: Record<ShipmentStatus, string> = {
  pending: "Pending",
  picked_up: "Picked Up",
  departed: "Departed",
  in_transit: "In Transit",
  out_for_delivery: "Out for Delivery",
  delivered: "Delivered",
  delayed: "Delayed",
};

export interface ShipmentContact {
  name: string;
  zipcode: string;
  subDistrict: string;
  district: string;
  province: string;
  coords: GeoPoint;
}


export interface ApiEnvelope<T> {
  success: boolean;
  data: T;
  pagination?: ShipmentPage["pagination"];
  error?: string;
}

export interface Shipment {
  id: string;
  trackingNumber: string;
  customer: ShipmentContact;
  receiver: ShipmentContact;
  origin: string;
  destination: string;
  currentCoords: GeoPoint;
  status: ShipmentStatus;
  carrier: string;
  hubId?: string;
  weight: number;
  items: number;
  estimatedDelivery: string;
  estimatedDeliveryRaw?: string;
  createdAt: string;
  progress: number;
}

export interface ShipmentDTO {
  id: string;
  trackingNumber: string;
  customer: ShipmentContact;
  receiver: ShipmentContact;
  currentCoords: GeoPoint;
  origin: string;
  destination: string;
  status: string;
  carrier: string;
  hubId?: string;
  weight: number;
  items: number;
  estimatedDelivery: string;
  createdAt: string;
  progress: number;
}

export interface ShipmentPage {
  data: Shipment[];
  pagination: { page: number; limit: number; total: number; totalPages: number };
}
