import Manifest from "@/features/shipments/Manifest";
import { getShipmentsPage } from "@/features/shipments/server-api";
import type { ShipmentPage } from "@/features/shipments/types";

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
