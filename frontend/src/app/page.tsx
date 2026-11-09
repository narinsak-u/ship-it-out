import { Activity, ArrowRight, Boxes, Globe2, Truck } from "lucide-react";

import Link from "next/link";
import { getHomeData } from "@/features/shipments/home-data";
import TrackingSearch from "@/features/shipments/TrackingSearch";

export const dynamic = "force-dynamic";

const metrics = [
  { label: "Active shipments", icon: Truck },
  { label: "Delivered", icon: Boxes },
  { label: "On-time rate", icon: Activity },
  { label: "Provinces served", icon: Globe2 },
];

export default async function HomePage() {
  const { analytics, shipments } = await getHomeData();
  const values = [
    analytics ? analytics.total - analytics.delivered : "—",
    analytics?.delivered ?? "—",
    "99.9%",
    "77",
  ];

  return (
    <main>
      <section className="border-b border-border bg-card/50">
        <div className="mx-auto max-w-7xl px-4 py-20 sm:px-6 lg:py-28">
          <div className="max-w-3xl">
            <span className="inline-flex items-center gap-2 rounded-full border border-primary/30 bg-primary/10 px-3 py-1 text-xs font-semibold uppercase tracking-wider text-primary">
              <i className="size-2 rounded-full bg-primary" /> Live ops console
            </span>
            <h1 className="mt-6 text-5xl font-bold tracking-tight sm:text-7xl">
              <span className="flex items-center gap-4">
                Move fast <Truck className="size-10 text-primary" aria-hidden="true" />
              </span>
              <span className="block text-muted-foreground">Break nothing.</span>
            </h1>
            <p className="mt-6 max-w-2xl text-lg text-muted-foreground">
              Trace every parcel from origin warehouse to doorstep. Realtime telemetry, port-side
              timelines, and exception alerts in one console.
            </p>
            <div className="mt-8 max-w-xl">
              <TrackingSearch />
            </div>
          </div>
        </div>
      </section>
      <section className="border-b border-border" aria-label="Key metrics">
        <div className="mx-auto grid max-w-7xl grid-cols-2 divide-x divide-border px-4 sm:grid-cols-4 sm:px-6">
          {metrics.map((metric, index) => {
            const Icon = metric.icon;
            return (
              <div className="flex items-center gap-4 px-4 py-6 first:pl-0" key={metric.label}>
                <span
                  className="flex size-11 shrink-0 items-center justify-center rounded-lg bg-primary/10 text-primary"
                  aria-hidden="true"
                >
                  <Icon size={20} />
                </span>
                <div>
                  <div className="text-2xl font-bold">{values[index]}</div>
                  <div className="text-sm text-muted-foreground">{metric.label}</div>
                </div>
              </div>
            );
          })}
        </div>
      </section>
      <section className="mx-auto max-w-7xl px-4 py-12 sm:px-6">
        <div className="mb-6 flex items-end justify-between gap-4">
          <div>
            <h2 className="text-2xl font-semibold">Recent shipments</h2>
            <p className="mt-1 text-muted-foreground">Latest activity from the fleet.</p>
          </div>
          <Link
            className="inline-flex items-center gap-2 text-sm font-medium text-primary hover:underline focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
            href="/orders"
          >
            View all <ArrowRight size={16} aria-hidden="true" />
          </Link>
        </div>
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {shipments === null ? (
            <p className="text-sm text-muted-foreground" role="status">
              Recent shipments are unavailable right now.
            </p>
          ) : shipments.length === 0 ? (
            <p className="text-sm text-muted-foreground">No active shipments yet.</p>
          ) : (
            shipments.map((shipment) => (
              <Link
                className="rounded-xl border border-border bg-card p-5 transition-colors hover:border-primary/50 hover:bg-accent focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
                href={`/orders/${encodeURIComponent(shipment.id)}`}
                key={shipment.id}
              >
                <div className="flex items-center justify-between gap-3">
                  <span className="font-semibold">{shipment.id}</span>
                  <span className="rounded-full bg-primary/10 px-2.5 py-1 text-xs font-medium capitalize text-primary">
                    {shipment.status.replaceAll("_", " ")}
                  </span>
                </div>
                <div className="mt-3 font-mono text-sm text-muted-foreground">
                  {shipment.trackingNumber}
                </div>
                <div className="mt-1 text-sm">{shipment.customer.name}</div>
                <div className="mt-4 flex items-center gap-2 text-xs text-muted-foreground">
                  <span className="truncate">{shipment.origin}</span>
                  <ArrowRight className="size-3.5 shrink-0" aria-hidden="true" />
                  <span className="truncate">{shipment.destination}</span>
                </div>
                <div
                  className="mt-4 h-1.5 overflow-hidden rounded-full bg-muted"
                  aria-label={`${shipment.progress}% complete`}
                >
                  <span
                    className="block h-full rounded-full bg-primary"
                    style={{ width: `${shipment.progress}%` }}
                  />
                </div>
              </Link>
            ))
          )}
        </div>
        {analytics === null && (
          <p className="mt-4 text-sm text-muted-foreground" role="status">
            Dashboard metrics are unavailable right now.
          </p>
        )}
      </section>
    </main>
  );
}
