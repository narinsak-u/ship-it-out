import { api } from "@/lib/api/client";
import {
  mapShipmentToOrder,
  mapEventToTrackingEvent,
  type BackendShipment,
  type BackendShipmentEvent,
} from "@/lib/api/mappers";
import type {
  Order,
  ShipmentStatus,
  TrackingEvent,
  ContactInfo,
  PaginatedResponse,
} from "@/lib/orders";

export interface OrderFormData {
  customer: ContactInfo;
  receiver: ContactInfo;
  carrier: string;
  weight: string;
  items: number;
  estimatedDelivery: string;
  status?: ShipmentStatus;
}

export interface ShipmentsQuery {
  page?: number;
  limit?: number;
  search?: string;
  status?: string;
  exclude_status?: string;
}

export async function fetchActiveDeliveries(): Promise<Order[]> {
  const result = await api.get<BackendShipment[]>("/shipments?limit=-1&exclude_status=delivered");
  if (result.error) throw new Error(result.error);
  return result.data!.map(mapShipmentToOrder);
}

export async function fetchOrdersPaginated(
  query: ShipmentsQuery = {},
): Promise<PaginatedResponse<Order>> {
  const params = new URLSearchParams();
  if (query.page) params.set("page", String(query.page));
  if (query.limit) params.set("limit", String(query.limit));
  if (query.search) params.set("search", query.search);
  if (query.status && query.status !== "all") params.set("status", query.status);
  if (query.exclude_status) params.set("exclude_status", query.exclude_status);
  const qs = params.toString();
  const result = await api.getRaw<{
    data: BackendShipment[];
    pagination: PaginatedResponse<Order>["pagination"];
  }>(`/shipments${qs ? `?${qs}` : ""}`);
  if (result.error) throw new Error(result.error);
  return {
    data: result.data!.data.map(mapShipmentToOrder),
    pagination: result.data!.pagination,
  };
}

export async function updateShipmentStatus(
  orderId: string,
  status: ShipmentStatus,
  hubId?: string,
): Promise<Order> {
  const result = await api.patch<BackendShipment>(`/shipments/${orderId}/status`, {
    status,
    ...(hubId ? { hubId } : {}),
  });
  if (result.error) throw new Error(result.error);
  return mapShipmentToOrder(result.data!);
}

export async function createOrder(data: OrderFormData): Promise<Order> {
  const { estimatedDelivery: _, ...body } = data;
  const result = await api.post<BackendShipment>("/shipments", body);
  if (result.error) throw new Error(result.error);
  return mapShipmentToOrder(result.data!);
}

export async function updateOrder(id: string, data: Partial<OrderFormData>): Promise<Order> {
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

export async function fetchOrderEvents(trackingNumber: string): Promise<TrackingEvent[]> {
  const result = await api.get<{
    shipment: BackendShipment;
    events: BackendShipmentEvent[];
  }>(`/track/${trackingNumber}`);
  if (result.error) throw new Error(result.error);
  return result.data!.events.map(mapEventToTrackingEvent);
}
