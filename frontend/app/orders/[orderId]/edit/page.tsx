import { notFound } from "next/navigation";
import { ShipmentFormPage } from "@/features/shipments/ShipmentFormPage";
import { mapShipmentToOrder } from "@/lib/api/mappers";
import { getShipment } from "@/lib/api/preview";

export default async function EditShipmentPage({
  params,
}: {
  params: Promise<{ orderId: string }>;
}) {
  const { orderId } = await params;
  const shipment = await getShipment(orderId);
  if (!shipment) notFound();
  return <ShipmentFormPage initial={mapShipmentToOrder(shipment)} />;
}
