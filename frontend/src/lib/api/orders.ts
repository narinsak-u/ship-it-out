import { orders, type Order, type ShipmentStatus, type ContactInfo, type Location, statusLabels } from "@/lib/orders";

function delay(ms = 200): Promise<void> {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

export interface OrderFormData {
  customer: ContactInfo;
  receiver: ContactInfo;
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
  const hex = () =>
    Math.floor(Math.random() * 0x10000)
      .toString(16)
      .toUpperCase()
      .padStart(4, "0");
  return `TRK-${hex()}-${hex()}`;
}

function today(): string {
  const d = new Date();
  const months = [
    "Jan",
    "Feb",
    "Mar",
    "Apr",
    "May",
    "Jun",
    "Jul",
    "Aug",
    "Sep",
    "Oct",
    "Nov",
    "Dec",
  ];
  return `${months[d.getMonth()]} ${d.getDate()}, ${d.getFullYear()}`;
}

function composeAddress(info: ContactInfo): string {
  return `${info.subDistrict}, ${info.district}, ${info.province}`;
}

export async function fetchActiveDeliveries() {
  await delay();
  return orders.filter((o) => o.status !== "delivered");
}

export async function updateShipmentStatus(orderId: string, status: ShipmentStatus) {
  await delay(100);
  const order = orders.find((o) => o.id === orderId);
  if (!order) throw new Error("Order not found");

  const now = today();
  let location: Location;

  if (status === "delivered") {
    location = {
      name: order.receiver.name,
      lat: order.receiver.coords.lat,
      lng: order.receiver.coords.lng,
    };
    order.progress = 100;
  } else if (status === "out_for_delivery") {
    location = {
      name: order.destination,
      lat: order.receiver.coords.lat,
      lng: order.receiver.coords.lng,
    };
    order.progress = Math.max(order.progress, 70);
  } else {
    const last = order.events[0]?.location;
    location = {
      name: last?.name ?? "Unknown",
      lat: last?.lat ?? 0,
      lng: last?.lng ?? 0,
    };
    if (status === "in_transit") order.progress = Math.max(order.progress, 30);
  }

  order.status = status;
  order.events.unshift({
    timestamp: now,
    location: { ...location },
    status: statusLabels[status],
    description: `Status updated to ${statusLabels[status]}.`,
  });
  order.currentCoords = { lat: location.lat, lng: location.lng };

  return { ...order };
}

export async function createOrder(data: OrderFormData): Promise<Order> {
  await delay(200);
  const now = today();
  const pickupLocation: Location = {
    name: composeAddress(data.customer),
    lat: data.customer.coords.lat,
    lng: data.customer.coords.lng,
  };
  const order: Order = {
    id: generateId(),
    trackingNumber: generateTrackingNumber(),
    customer: { ...data.customer },
    receiver: { ...data.receiver },
    origin: composeAddress(data.customer),
    destination: composeAddress(data.receiver),
    currentCoords: { lat: pickupLocation.lat, lng: pickupLocation.lng },
    status: "pending",
    carrier: data.carrier,
    weight: data.weight,
    items: data.items,
    estimatedDelivery: data.estimatedDelivery,
    createdAt: now,
    progress: 0,
    events: [
      {
        timestamp: now,
        location: pickupLocation,
        status: "Created",
        description: "Order created.",
      },
    ],
  };
  orders.unshift(order);
  return { ...order };
}

export async function updateOrder(id: string, data: Partial<OrderFormData>): Promise<Order> {
  await delay(200);
  const idx = orders.findIndex((o) => o.id === id);
  if (idx === -1) throw new Error("Order not found");
  const order = orders[idx];
  if (data.customer !== undefined) {
    order.customer = { ...data.customer };
    order.origin = composeAddress(data.customer);
  }
  if (data.receiver !== undefined) {
    order.receiver = { ...data.receiver };
    order.destination = composeAddress(data.receiver);
  }
  if (data.carrier !== undefined) order.carrier = data.carrier;
  if (data.weight !== undefined) order.weight = data.weight;
  if (data.items !== undefined) order.items = data.items;
  if (data.estimatedDelivery !== undefined) order.estimatedDelivery = data.estimatedDelivery;
  if (data.status !== undefined) order.status = data.status;
  return { ...order };
}
