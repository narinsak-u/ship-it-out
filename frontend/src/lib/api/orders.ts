import { api } from "@/lib/api/client";
import {
  mapShipmentToOrder,
  mapEventToTrackingEvent,
  type BackendShipment,
  type BackendShipmentEvent,
} from "@/lib/api/mappers";
import type { Order, ShipmentStatus, TrackingEvent, ContactInfo } from "@/lib/orders";

export interface OrderFormData {
  customer: ContactInfo;
  receiver: ContactInfo;
  carrier: string;
  weight: string;
  items: number;
  estimatedDelivery: string;
  status?: ShipmentStatus;
}

export async function fetchActiveDeliveries(): Promise<Order[]> {
  const result = await api.get<BackendShipment[]>("/shipments");
  if (result.error) throw new Error(result.error);
  return result.data!
    .filter((s) => s.status !== "delivered")
    .map(mapShipmentToOrder);
}

export async function updateShipmentStatus(
  orderId: string,
  status: ShipmentStatus,
): Promise<Order> {
  const result = await api.patch<BackendShipment>(
    `/shipments/${orderId}/status`,
    { status },
  );
  if (result.error) throw new Error(result.error);
  return mapShipmentToOrder(result.data!);
}

export async function createOrder(data: OrderFormData): Promise<Order> {
  const { estimatedDelivery: _, ...body } = data;
  const result = await api.post<BackendShipment>("/shipments", body);
  if (result.error) throw new Error(result.error);
  return mapShipmentToOrder(result.data!);
}

export async function updateOrder(
  id: string,
  data: Partial<OrderFormData>,
): Promise<Order> {
  const result = await api.put<BackendShipment>(`/shipments/${id}`, data);
  if (result.error) throw new Error(result.error);
  return mapShipmentToOrder(result.data!);
}

export async function deleteOrder(id: string): Promise<void> {
  const result = await api.del(`/shipments/${id}`);
  if (result.error) throw new Error(result.error);
}

export async function fetchOrder(id: string): Promise<Order> {
  const result = await api.get<BackendShipment>(`/shipments/${id}`);
  if (result.error) throw new Error(result.error);
  return mapShipmentToOrder(result.data!);
}

export async function fetchOrderEvents(
  trackingNumber: string,
): Promise<TrackingEvent[]> {
  const result = await api.get<{
    shipment: BackendShipment;
    events: BackendShipmentEvent[];
  }>(`/track/${trackingNumber}`);
  if (result.error) throw new Error(result.error);
  return result.data!.events.map(mapEventToTrackingEvent);
}
