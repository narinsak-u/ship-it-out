import { ChartNoAxesColumnIncreasing, Package, Warehouse } from "lucide-react";
import Link from "next/link";
import type { AnalyticsData } from "@/features/analytics/data";
import { getAnalyticsData } from "@/features/analytics/server-api";
import { AnalyticsPanel } from "@/features/analytics/AnalyticsPanel";
import { DeliveriesPanel } from "@/features/operations/DeliveriesPanel";
import { HubsPanel } from "@/features/hubs/HubsPanel";
import type { Hub } from "@/features/hubs/types";
import { getHubs } from "@/features/hubs/server-api";
import type { Shipment } from "@/features/shipments/types";
import { getShipments } from "@/features/shipments/server-api";

export const dynamic = "force-dynamic";

async function getActiveDeliveries(): Promise<Shipment[]> {
  return getShipments("shipments?limit=-1&exclude_status=delivered");
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
  let deliveries: Shipment[] = [];
  let deliveriesError = false;
  let analytics: AnalyticsData | null = null;
  let analyticsError = false;
  if (analyticsSelected) {
    try {
      analytics = await getAnalyticsData();
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
    <main className="min-h-screen">
      <section className="border-b border-border bg-[linear-gradient(135deg,#102b35,#111b27_65%)]">
        <div className="mx-auto max-w-screen-2xl px-4 py-12 sm:px-6 sm:py-20">
          <span className="font-mono text-[11px] uppercase tracking-[0.12em] text-primary">
            / carriers
          </span>
          <h1 className="mt-3 font-mono text-3xl font-semibold tracking-tight sm:text-6xl">
            Carrier operations
          </h1>
          <p className="mt-2 text-muted-foreground">
            Manage carrier hubs and monitor active deliveries across the fleet.
          </p>
        </div>
      </section>
      <section className="mx-auto w-full max-w-screen-2xl px-4 py-8 font-mono text-sm sm:px-6 sm:py-8">
        <nav
          className="flex overflow-x-auto border-b border-border"
          aria-label="Carrier operations"
        >
          {[
            {
              label: "Active Deliveries",
              href: "/carriers",
              active: !analyticsSelected && !hubsSelected,
              icon: Package,
            },
            { label: "Hubs", href: "/carriers?tab=hubs", active: hubsSelected, icon: Warehouse },
            {
              label: "Analytics",
              href: "/carriers?tab=analytics",
              active: analyticsSelected,
              icon: ChartNoAxesColumnIncreasing,
            },
          ].map(({ label, href, active, icon: Icon }) => (
            <Link
              key={href}
              className={`inline-flex shrink-0 items-center gap-2 whitespace-nowrap border border-b-0 px-5 py-3 sm:py-[18px] focus-visible:relative focus-visible:z-10 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary ${active ? "mb-[-1px] rounded-t-lg border-border bg-card text-foreground" : "border-transparent text-muted-foreground hover:text-foreground"}`}
              href={href}
              aria-current={active ? "page" : undefined}
            >
              <Icon className="size-5" aria-hidden="true" />
              {label}
            </Link>
          ))}
        </nav>
        <div className="min-h-[calc(100dvh-24rem)] rounded-b-xl rounded-tr-xl border border-t-0 border-border bg-card p-4 sm:p-8">
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
