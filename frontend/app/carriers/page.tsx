import Link from "next/link";
import { AnalyticsPanel } from "@/features/analytics/AnalyticsPanel";
import { DeliveriesPanel } from "@/features/operations/DeliveriesPanel";
import { HubsPanel } from "@/features/operations/HubsPanel";
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
      <section className="border-b border-border bg-[linear-gradient(135deg,#102b35,#111b27_65%)]">
        <div className="mx-auto max-w-7xl px-4 py-8 sm:px-6 sm:py-12">
          <span className="font-mono text-[11px] uppercase tracking-[0.12em] text-primary">
            / carriers
          </span>
          <h1 className="mt-3 text-3xl font-semibold tracking-tight sm:text-5xl">
            Carrier operations
          </h1>
          <p className="mt-2 text-muted-foreground">
            Manage carrier hubs and monitor active deliveries across the fleet.
          </p>
        </div>
      </section>
      <section className="mx-auto max-w-7xl px-4 py-8 sm:px-6 sm:py-12">
        <nav className="mb-[-1px] flex overflow-x-auto" aria-label="Carrier operations">
          {[
            ["Active Deliveries", "/carriers", !analyticsSelected && !hubsSelected],
            ["Hubs", "/carriers?tab=hubs", hubsSelected],
            ["Analytics", "/carriers?tab=analytics", analyticsSelected],
          ].map(([label, href, active]) => (
            <Link
              key={href as string}
              className={`whitespace-nowrap rounded-t-lg border px-4 py-3 font-mono text-sm focus-visible:relative focus-visible:z-10 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary ${active ? "border-border border-b-card bg-card text-foreground" : "border-border bg-secondary text-muted-foreground"}`}
              href={href as string}
              aria-current={active ? "page" : undefined}
            >
              {label}
            </Link>
          ))}
        </nav>
        <div className="rounded-b-lg rounded-tr-lg border border-border bg-card p-4 shadow-xl sm:p-6">
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
