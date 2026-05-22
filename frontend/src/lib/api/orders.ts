import { orders, type Order, type ShipmentStatus } from "@/lib/orders";

function delay(ms = 200): Promise<void> {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

export interface OrderFormData {
  customer: string;
  origin: string;
  destination: string;
  carrier: string;
  weight: string;
  items: number;
  estimatedDelivery: string;
  status?: ShipmentStatus;
}

function generateId(): string {
  const num = orders.length + 10251;
  return `ORD-${num}`;
}

function generateTrackingNumber(): string {
  const hex = () => Math.floor(Math.random() * 0x10000).toString(16).toUpperCase().padStart(4, "0");
  return `TRK-${hex()}-${hex()}`;
}

function today(): string {
  const d = new Date();
  const months = ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"];
  return `${months[d.getMonth()]} ${d.getDate()}, ${d.getFullYear()}`;
}

export async function fetchActiveDeliveries() {
  await delay();
  return orders.filter((o) => o.status !== "delivered");
}

export async function updateShipmentStatus(orderId: string, status: ShipmentStatus) {
  await delay(100);
  const order = orders.find((o) => o.id === orderId);
  if (!order) throw new Error("Order not found");
  order.status = status;
  if (status === "delivered") order.progress = 100;
  if (status === "in_transit") order.progress = Math.max(order.progress, 30);
  if (status === "out_for_delivery") order.progress = Math.max(order.progress, 70);
  return { ...order };
}

export async function createOrder(data: OrderFormData): Promise<Order> {
  await delay(200);
  const now = today();
  const order: Order = {
    id: generateId(),
    trackingNumber: generateTrackingNumber(),
    customer: data.customer,
    origin: data.origin,
    destination: data.destination,
    originCoords: { lat: 0, lng: 0 },
    destinationCoords: { lat: 0, lng: 0 },
    currentCoords: { lat: 0, lng: 0 },
    status: "pending",
    carrier: data.carrier,
    weight: data.weight,
    items: data.items,
    estimatedDelivery: data.estimatedDelivery,
    createdAt: now,
    progress: 0,
    events: [{ timestamp: now, location: data.origin, status: "Created", description: "Order created." }],
  };
  orders.unshift(order);
  return { ...order };
}

export async function updateOrder(id: string, data: Partial<OrderFormData>): Promise<Order> {
  await delay(200);
  const idx = orders.findIndex((o) => o.id === id);
  if (idx === -1) throw new Error("Order not found");
  const order = orders[idx];
  if (data.customer !== undefined) order.customer = data.customer;
  if (data.origin !== undefined) order.origin = data.origin;
  if (data.destination !== undefined) order.destination = data.destination;
  if (data.carrier !== undefined) order.carrier = data.carrier;
  if (data.weight !== undefined) order.weight = data.weight;
  if (data.items !== undefined) order.items = data.items;
  if (data.estimatedDelivery !== undefined) order.estimatedDelivery = data.estimatedDelivery;
  if (data.status !== undefined) order.status = data.status;
  return { ...order };
}
