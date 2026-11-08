import Manifest from "@/features/shipments/Manifest";
import { getShipmentsPage, type ShipmentPage } from "@/lib/api/preview";

export const dynamic = "force-dynamic";

export default async function OrdersPage({
  searchParams,
}: {
  searchParams: Promise<{ created?: string }>;
}) {
  const { created } = await searchParams;
  let initialData: ShipmentPage | null = null;
  let initialError: string | null = null;
  try {
    initialData = await getShipmentsPage();
  } catch (error) {
    initialError = error instanceof Error ? error.message : "Failed to load shipments";
  }
  return (
    <Manifest
      initialData={initialData}
      initialError={initialError}
      initialMessage={created === "1" ? "Shipment created." : ""}
    />
  );
}
