import type { Driver, DriverStatus, Hub, HubStatus } from "@/lib/carriers";
import { carriers, drivers, hubs } from "@/lib/carriers";
import { orders, type Order, type ShipmentStatus } from "@/lib/orders";

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

export async function fetchHubs() {
  await delay();
  return [...hubs];
}

export async function createHub(data: Omit<Hub, "id">) {
  await delay(150);
  const id = `HUB-${String(hubs.length + 1).padStart(3, "0")}`;
  const hub: Hub = { id, ...data };
  hubs.push(hub);
  return { ...hub };
}

export async function updateHub(id: string, data: Partial<Hub>) {
  await delay(150);
  const idx = hubs.findIndex((h) => h.id === id);
  if (idx === -1) throw new Error("Hub not found");
  hubs[idx] = { ...hubs[idx], ...data };
  return { ...hubs[idx] };
}

export async function deleteHub(id: string) {
  await delay(100);
  const idx = hubs.findIndex((h) => h.id === id);
  if (idx === -1) throw new Error("Hub not found");
  hubs.splice(idx, 1);
}
