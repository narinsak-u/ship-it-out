import type { Driver, DriverStatus, Hub, HubStatus } from "@/lib/carriers";
import { carriers, drivers, hubs } from "@/lib/carriers";
import { orders } from "@/lib/orders";
import { api } from "@/lib/api/client";
import { mapBackendHubToHub, type BackendHub } from "@/lib/api/mappers";

function delay(ms = 200): Promise<void> {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

// --- Carriers ---

export async function fetchCarriers() {
  await delay();
  return [...carriers];
}

// --- Drivers ---

export async function fetchDrivers() {
  await delay();
  return [...drivers];
}

export async function assignDriverToOrder(driverId: string, orderId: string) {
  await delay(100);
  const driver = drivers.find((d) => d.id === driverId);
  if (!driver) throw new Error("Driver not found");
  const order = orders.find((o) => o.id === orderId);
  if (!order) throw new Error("Order not found");
  driver.status = "on_delivery" as DriverStatus;
  order.driverId = driverId;
  return { driver: { ...driver }, order: { ...order } };
}

// --- Hubs ---

export async function fetchHubs(): Promise<Hub[]> {
  const result = await api.get<BackendHub[]>("/hubs");
  if (result.error) throw new Error(result.error);
  return result.data!.map(mapBackendHubToHub);
}

export async function createHub(data: Omit<Hub, "id">): Promise<Hub> {
  const result = await api.post<BackendHub>("/hubs", data);
  if (result.error) throw new Error(result.error);
  return mapBackendHubToHub(result.data!);
}

export async function updateHub(id: string, data: Partial<Hub>): Promise<Hub> {
  const result = await api.put<BackendHub>(`/hubs/${id}`, data);
  if (result.error) throw new Error(result.error);
  return mapBackendHubToHub(result.data!);
}

export async function deleteHub(id: string): Promise<void> {
  const result = await api.del(`/hubs/${id}`);
  if (result.error) throw new Error(result.error);
}
