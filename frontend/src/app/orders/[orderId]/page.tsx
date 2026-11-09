import { notFound } from "next/navigation";
import { getShipment } from "@/features/shipments/server-api";
import ShipmentDetail from "@/features/shipments/ShipmentDetail";
import { getTrackingEvents } from "@/features/tracking/server-api";

export default async function ShipmentPage({
  params,
  searchParams,
}: {
  params: Promise<{ orderId: string }>;
  searchParams: Promise<{ updated?: string }>;
}) {
  const [{ orderId }, { updated }] = await Promise.all([params, searchParams]);
  const shipment = await getShipment(orderId);
  if (!shipment) notFound();
  const events = await getTrackingEvents(shipment.trackingNumber);

  return <ShipmentDetail shipment={shipment} events={events} updated={updated} />;
}
