import Link from "next/link";
import { AnalyticsPanel } from "@/features/analytics/AnalyticsPanel";
import { DeliveriesPanel } from "@/features/operations/DeliveriesPanel";
import { HubsPanel } from "@/features/operations/HubsPanel";
import styles from "@/features/operations/CarrierShell.module.css";
import type { Hub } from "@/lib/hubs";
import type { Order } from "@/lib/orders";
import { mapShipmentToOrder, type BackendShipment } from "@/lib/api/mappers";
import { toAnalyticsData, type AnalyticsData } from "@/features/analytics/data";
import { getApiOrigin } from "@/lib/api/server";

export const dynamic = "force-dynamic";

async function getHubs(): Promise<Hub[]> {
  const response = await fetch(`${getApiOrigin()}/api/hubs`, { cache: "no-store" });
  if (!response.ok) throw new Error(`Hub API returned ${response.status}`);
  const result = await response.json();
  return result.data as Hub[];
}

async function getActiveDeliveries(): Promise<Order[]> {
  const response = await fetch(
    `${getApiOrigin()}/api/shipments?limit=-1&exclude_status=delivered`,
    { cache: "no-store" },
  );
  if (!response.ok) throw new Error(`Shipment API returned ${response.status}`);
  const result = await response.json();
  return (result.data as BackendShipment[]).map(mapShipmentToOrder);
}

async function getAnalytics(): Promise<AnalyticsData> {
  const origin = getApiOrigin();
  const [overviewResponse, seriesResponse] = await Promise.all([
    fetch(`${origin}/api/analytics/overview`, { cache: "no-store" }),
    fetch(`${origin}/api/analytics/timeseries`, { cache: "no-store" }),
  ]);
  if (!overviewResponse.ok || !seriesResponse.ok) throw new Error("Analytics request failed");

  const [overviewResult, seriesResult] = await Promise.all([
    overviewResponse.json(),
    seriesResponse.json(),
  ]);
  return toAnalyticsData(overviewResult, seriesResult);
}

export default async function CarriersPage({
  searchParams,
}: {
  searchParams: Promise<{ tab?: string }>;
}) {
  const { tab } = await searchParams;
  const analyticsSelected = tab === "analytics";
  const hubsSelected = tab === "hubs";
  let hubs: Hub[] = [];
  let error = false;
  let deliveries: Order[] = [];
  let deliveriesError = false;
  let analytics: AnalyticsData | null = null;
  let analyticsError = false;

  if (analyticsSelected) {
    try {
      analytics = await getAnalytics();
    } catch {
      analyticsError = true;
    }
  } else if (hubsSelected) {
    try {
      hubs = await getHubs();
    } catch {
      error = true;
    }
  } else {
    try {
      [deliveries, hubs] = await Promise.all([getActiveDeliveries(), getHubs()]);
    } catch {
      deliveriesError = true;
    }
  }

  return (
    <main>
      <section className={styles.hero}>
        <div className={styles.heroInner}>
          <span className={styles.eyebrow}>/ carriers</span>
          <h1>Carrier operations</h1>
          <p>Manage carrier hubs and monitor active deliveries across the fleet.</p>
        </div>
      </section>
      <section className={styles.content}>
        <nav className={styles.tabs} aria-label="Carrier operations">
          <Link
            className={styles.tab}
            href="/carriers"
            aria-current={!analyticsSelected && !hubsSelected ? "page" : undefined}
          >
            Active Deliveries
          </Link>
          <Link
            className={styles.tab}
            href="/carriers?tab=hubs"
            aria-current={hubsSelected ? "page" : undefined}
          >
            Hubs
          </Link>
          <Link
            className={styles.tab}
            href="/carriers?tab=analytics"
            aria-current={analyticsSelected ? "page" : undefined}
          >
            Analytics
          </Link>
        </nav>
        <div className={styles.panel}>
          {analyticsSelected ? (
            <AnalyticsPanel initialData={analytics} initialError={analyticsError} />
          ) : hubsSelected ? (
            <HubsPanel hubs={hubs} initialError={error} />
          ) : (
            <DeliveriesPanel
              initialDeliveries={deliveries}
              hubs={hubs}
              initialError={deliveriesError}
            />
          )}
        </div>
      </section>
    </main>
  );
}
