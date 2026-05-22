import { orders, type ShipmentStatus } from "@/lib/orders";

function delay(ms = 200): Promise<void> {
  return new Promise((resolve) => setTimeout(resolve, ms));
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
