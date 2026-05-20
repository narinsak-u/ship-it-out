import { createFileRoute, Link, notFound } from "@tanstack/react-router";
import { lazy, Suspense, useEffect, useState } from "react";
import { ArrowLeft, MapPin, Package, Truck, Calendar, Hash, User, Weight, Maximize2 } from "lucide-react";
import { SiteHeader } from "@/components/SiteHeader";
import { StatusBadge } from "@/components/StatusBadge";
import { getOrder, type TrackingEvent } from "@/lib/orders";

const ShipmentMap = lazy(() =>
  import("@/components/ShipmentMap").then((m) => ({ default: m.ShipmentMap })),
);

export const Route = createFileRoute("/orders/$orderId")({
  ssr: false,
  head: ({ params }) => ({
    meta: [
      { title: `${params.orderId} — Shipment Detail · Harbor Ops` },
      { name: "description", content: `Live tracking timeline and shipment details for ${params.orderId}.` },
      { property: "og:title", content: `Shipment ${params.orderId} — Harbor Ops` },
      { property: "og:description", content: `Live tracking and route detail for ${params.orderId}.` },
    ],
  }),
  loader: ({ params }) => {
    const order = getOrder(params.orderId);
    if (!order) throw notFound();
    return { order };
  },
  notFoundComponent: () => (
    <div className="min-h-screen">
      <SiteHeader />
      <div className="mx-auto max-w-2xl px-6 py-32 text-center">
        <h1 className="font-mono text-4xl">404</h1>
        <p className="mt-3 text-muted-foreground">Shipment not found.</p>
        <Link to="/orders" className="mt-6 inline-block font-mono text-sm text-primary">← Back to orders</Link>
      </div>
    </div>
  ),
  errorComponent: ({ error }) => (
    <div className="min-h-screen">
      <SiteHeader />
      <div className="mx-auto max-w-2xl px-6 py-32 text-center">
        <h1 className="font-mono text-2xl">Something broke</h1>
        <p className="mt-3 text-muted-foreground">{error.message}</p>
      </div>
    </div>
  ),
  component: OrderDetail,
});

function OrderDetail() {
  const { order } = Route.useLoaderData();
  const [mounted, setMounted] = useState(false);
  useEffect(() => setMounted(true), []);

  const meta = [
    { icon: Hash, label: "Tracking #", value: order.trackingNumber },
    { icon: User, label: "Customer", value: order.customer },
    { icon: Truck, label: "Carrier", value: order.carrier },
    { icon: Package, label: "Items", value: `${order.items}` },
    { icon: Weight, label: "Weight", value: order.weight },
    { icon: Calendar, label: "Created", value: order.createdAt },
  ];

  return (
    <div className="min-h-screen">
      <SiteHeader />

      <div className="mx-auto grid max-w-[1600px] gap-0 lg:grid-cols-[minmax(0,1fr)_minmax(0,1.1fr)]">
        {/* LEFT: details */}
        <div className="border-r border-border">
          <div className="px-6 py-8 lg:px-10 lg:py-10">
            <Link to="/orders" className="group inline-flex items-center gap-1.5 font-mono text-xs uppercase tracking-widest text-muted-foreground hover:text-foreground">
              <ArrowLeft className="h-3.5 w-3.5 transition-transform group-hover:-translate-x-1" />
              All orders
            </Link>

            <div className="mt-6 flex flex-wrap items-center gap-3">
              <h1 className="font-mono text-4xl font-semibold tracking-tight">{order.trackingNumber}</h1>
              <StatusBadge status={order.status} />
            </div>
            <div className="mt-2 font-mono text-sm text-muted-foreground">
              Order ID <span className="text-primary">{order.id}</span> · {order.customer}
            </div>

            {/* Route summary card */}
            <div className="mt-8 rounded-xl border border-border bg-card p-5 shadow-elegant">
              <div className="flex items-start gap-4">
                <div className="mt-1 flex flex-col items-center gap-1">
                  <span className="h-2.5 w-2.5 rounded-full bg-primary" />
                  <span className="h-10 w-px border-l border-dashed border-border" />
                  <span className="h-2.5 w-2.5 rounded-full bg-accent" />
                </div>
                <div className="flex-1 space-y-3">
                  <div>
                    <div className="font-mono text-[11px] uppercase tracking-widest text-muted-foreground">From</div>
                    <div className="font-mono text-sm">{order.origin}</div>
                  </div>
                  <div>
                    <div className="font-mono text-[11px] uppercase tracking-widest text-muted-foreground">To</div>
                    <div className="font-mono text-sm">{order.destination}</div>
                  </div>
                </div>
                <div className="rounded-md border border-border bg-secondary px-3 py-1.5 font-mono text-xs text-primary">
                  {order.carrier}
                </div>
              </div>

              <div className="mt-5 grid grid-cols-2 gap-4 border-t border-border pt-5 sm:grid-cols-3">
                <div>
                  <div className="font-mono text-[11px] uppercase tracking-widest text-muted-foreground">Progress</div>
                  <div className="font-mono text-sm">{order.progress}%</div>
                </div>
                <div>
                  <div className="font-mono text-[11px] uppercase tracking-widest text-muted-foreground">Created</div>
                  <div className="font-mono text-sm">{order.createdAt}</div>
                </div>
                <div>
                  <div className="font-mono text-[11px] uppercase tracking-widest text-muted-foreground">ETA</div>
                  <div className="font-mono text-sm text-primary">{order.estimatedDelivery}</div>
                </div>
              </div>

              <div className="mt-4 h-1.5 overflow-hidden rounded-full bg-secondary">
                <div className="h-full bg-gradient-accent transition-all" style={{ width: `${order.progress}%` }} />
              </div>
            </div>

            {/* Metadata grid */}
            <div className="mt-6 grid grid-cols-2 gap-3 sm:grid-cols-3">
              {meta.map((m) => (
                <div key={m.label} className="rounded-lg border border-border bg-card p-4">
                  <div className="flex items-center gap-2 font-mono text-[10px] uppercase tracking-wider text-muted-foreground">
                    <m.icon className="h-3.5 w-3.5 text-primary" />
                    {m.label}
                  </div>
                  <div className="mt-1.5 font-mono text-sm">{m.value}</div>
                </div>
              ))}
            </div>

            {/* Timeline */}
            <div className="mt-10">
              <h2 className="font-mono text-sm uppercase tracking-widest text-muted-foreground">Shipment status</h2>
              <ol className="relative mt-5 space-y-5 border-l border-border pl-6">
                {order.events.map((e: TrackingEvent, i: number) => (
                  <li key={i} className="relative">
                    <span className={`absolute -left-[31px] flex h-4 w-4 items-center justify-center rounded-full ring-4 ring-background ${i === 0 ? "bg-primary shadow-glow" : "bg-muted-foreground/40"}`}>
                      {i === 0 && <span className="h-1.5 w-1.5 animate-pulse rounded-full bg-primary-foreground" />}
                    </span>
                    <div className="flex items-center justify-between">
                      <span className="font-mono text-sm font-medium">{e.status}</span>
                      <span className="font-mono text-xs text-muted-foreground">{e.timestamp}</span>
                    </div>
                    <div className="mt-1 flex items-center gap-1.5 font-mono text-xs text-primary">
                      <MapPin className="h-3 w-3" />
                      {e.location}
                    </div>
                    {e.description && (
                      <p className="mt-1 text-sm text-muted-foreground">{e.description}</p>
                    )}
                  </li>
                ))}
              </ol>
            </div>
          </div>
        </div>

        {/* RIGHT: map */}
        <div className="relative bg-secondary/40 lg:sticky lg:top-16 lg:h-[calc(100vh-4rem)]">
          <div className="h-[420px] w-full lg:h-full">
            {mounted ? (
              <Suspense fallback={<MapSkeleton />}>
                <ShipmentMap
                  origin={order.originCoords}
                  destination={order.destinationCoords}
                  current={order.currentCoords}
                  originLabel={order.origin}
                  destinationLabel={order.destination}
                  carrier={order.carrier}
                  status={order.status}
                />
              </Suspense>
            ) : (
              <MapSkeleton />
            )}
          </div>

          {/* Floating telemetry card */}
          <div className="pointer-events-none absolute left-16 top-4 z-[400] rounded-lg border border-border bg-card/95 px-4 py-3 font-mono text-xs shadow-elegant backdrop-blur">
            <div className="flex items-center gap-2 text-muted-foreground">
              <span className="h-1.5 w-1.5 animate-pulse rounded-full bg-primary" />
              LIVE TELEMETRY
            </div>
            <div className="mt-1 text-sm text-foreground">{order.carrier}</div>
            <div className="text-muted-foreground">
              {order.currentCoords.lat.toFixed(2)}°, {order.currentCoords.lng.toFixed(2)}°
            </div>
          </div>

          <div className="pointer-events-none absolute right-4 top-4 z-[400] rounded-lg border border-border bg-card/95 px-3 py-2 font-mono text-[10px] uppercase tracking-widest text-muted-foreground shadow-elegant backdrop-blur">
            <Maximize2 className="mr-1 inline h-3 w-3" />
            Geo route
          </div>
        </div>
      </div>
    </div>
  );
}

function MapSkeleton() {
  return (
    <div className="flex h-full w-full items-center justify-center bg-gradient-hero">
      <div className="font-mono text-xs uppercase tracking-widest text-muted-foreground">
        Loading geo telemetry…
      </div>
    </div>
  );
}
