import Link from "next/link";
import { notFound } from "next/navigation";
import { getShipment } from "@/features/shipments/server-api";
import { getTrackingEvents } from "@/features/tracking/server-api";
import RouteMap from "@/features/tracking/RouteMap";

function date(value: string) {
  return new Date(value).toLocaleDateString("en-US", {
    month: "short",
    day: "numeric",
    year: "numeric",
  });
}

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

  return (
    <main className="mx-auto grid max-w-7xl gap-8 px-4 py-8 sm:px-6 lg:grid-cols-[minmax(0,1fr)_minmax(20rem,0.8fr)]">
      <section className="min-w-0 space-y-6">
        {updated === "1" && (
          <p
            className="rounded-md border border-primary/30 bg-primary/10 p-3 text-sm text-primary"
            role="status"
          >
            Shipment updated.
          </p>
        )}
        <Link
          className="inline-flex items-center gap-2 text-sm text-muted-foreground hover:text-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
          href="/orders"
        >
          ← <span>All orders</span>
        </Link>
        <div>
          <div className="flex flex-wrap items-center justify-between gap-3">
            <h1 className="break-all text-3xl font-bold tracking-tight">
              {shipment.trackingNumber}
            </h1>
            <span className="inline-flex items-center gap-2 rounded-full bg-primary/10 px-3 py-1 text-sm font-medium capitalize text-primary">
              <i className="size-2 rounded-full bg-primary" aria-hidden="true" />
              {shipment.status.replaceAll("_", " ")}
            </span>
          </div>
          <p className="mt-2 text-sm text-muted-foreground">
            Order ID <strong className="font-medium text-foreground">{shipment.id}</strong> ·{" "}
            {shipment.customer.name}
          </p>
        </div>
        <article className="rounded-xl border border-border bg-card p-5 sm:p-6">
          <div className="flex items-center gap-4">
            <div className="flex flex-col items-center">
              <span className="size-3 rounded-full bg-primary" />
              <span className="h-12 w-px bg-border" />
              <span className="size-3 rounded-full border-2 border-primary" />
            </div>
            <div className="min-w-0 flex-1 space-y-3">
              <div>
                <small className="text-xs font-semibold tracking-wider text-muted-foreground">
                  FROM
                </small>
                <div className="truncate">{shipment.origin}</div>
              </div>
              <div>
                <small className="text-xs font-semibold tracking-wider text-muted-foreground">
                  TO
                </small>
                <div className="truncate">{shipment.destination}</div>
              </div>
            </div>
            <span className="hidden rounded-full bg-muted px-3 py-1 text-xs text-muted-foreground sm:inline-flex">
              {shipment.carrier}
            </span>
          </div>
          <div className="mt-6 grid grid-cols-3 gap-3 border-t border-border pt-5 text-sm">
            <div>
              <small className="text-xs text-muted-foreground">PROGRESS</small>
              <div className="mt-1 font-semibold">{shipment.progress}%</div>
            </div>
            <div>
              <small className="text-xs text-muted-foreground">CREATED</small>
              <div className="mt-1">{date(shipment.createdAt)}</div>
            </div>
            <div>
              <small className="text-xs text-muted-foreground">ETA</small>
              <div className="mt-1 font-semibold text-primary">
                {date(shipment.estimatedDelivery)}
              </div>
            </div>
          </div>
          <div className="mt-4 h-1.5 overflow-hidden rounded-full bg-muted">
            <span
              className="block h-full rounded-full bg-primary"
              style={{ width: `${shipment.progress}%` }}
            />
          </div>
        </article>
        <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
          {[
            ["Tracking #", shipment.trackingNumber],
            ["Customer", shipment.customer.name],
            ["Carrier", shipment.carrier],
            ["Weight/Kg", `${shipment.weight} kg`],
            ["Created", date(shipment.createdAt)],
          ].map(([label, value]) => (
            <article className="min-w-0 rounded-lg border border-border bg-card p-4" key={label}>
              <small className="text-xs text-muted-foreground">{label}</small>
              <div className="mt-1 truncate text-sm font-medium">{value}</div>
            </article>
          ))}
        </div>
        <section className="rounded-xl border border-border bg-card p-5 sm:p-6">
          <h2 className="text-lg font-semibold">Shipment status</h2>
          <ol className="mt-5 space-y-5 border-l border-border pl-5">
            {[...events].reverse().map((event, index) => (
              <li className="relative" key={`${event.timestamp}-${index}`}>
                <span
                  className={`absolute -left-[1.58rem] top-1 size-3 rounded-full border-2 border-background ${index === 0 ? "bg-primary" : "bg-muted-foreground"}`}
                />
                <div className="flex flex-wrap items-baseline justify-between gap-x-3 gap-y-1">
                  <strong className="text-sm">{event.status}</strong>
                  <time className="text-xs text-muted-foreground">
                    {new Date(event.timestamp).toLocaleString("en-US", {
                      month: "short",
                      day: "numeric",
                      hour: "2-digit",
                      minute: "2-digit",
                      hour12: false,
                    })}
                  </time>
                </div>
                <p className="mt-1 text-sm text-muted-foreground">⌖ {event.location.name}</p>
                {event.description && (
                  <p className="mt-1 text-sm text-muted-foreground">{event.description}</p>
                )}
              </li>
            ))}
          </ol>
        </section>
      </section>
      <aside className="relative h-[24rem] overflow-hidden rounded-xl border border-border bg-card lg:sticky lg:top-6 lg:h-[calc(100vh-3rem)]">
        <RouteMap
          origin={shipment.customer.coords}
          destination={shipment.receiver.coords}
          current={shipment.currentCoords}
          originLabel={shipment.origin}
          destinationLabel={shipment.destination}
          carrier={shipment.carrier}
        />
        <div className="absolute bottom-5 left-5 right-5 rounded-lg border border-border bg-background/95 p-4 shadow-lg backdrop-blur">
          <small className="flex items-center gap-2 text-xs font-semibold tracking-wider text-primary">
            <i className="size-2 rounded-full bg-primary" aria-hidden="true" /> LIVE TELEMETRY
          </small>
          <strong className="mt-2 block">{shipment.carrier}</strong>
          <span className="mt-1 block font-mono text-sm text-muted-foreground">
            {shipment.currentCoords.lat.toFixed(2)}°, {shipment.currentCoords.lng.toFixed(2)}°
          </span>
        </div>
        <div className="absolute right-4 top-4 rounded bg-background/90 px-2 py-1 font-mono text-xs text-muted-foreground">
          ⌗ GEO ROUTE
        </div>
      </aside>
    </main>
  );
}
