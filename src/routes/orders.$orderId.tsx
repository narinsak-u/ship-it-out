import { createFileRoute, Link, notFound } from "@tanstack/react-router";
import { ArrowLeft, MapPin, Package, Truck, Calendar, Hash, User, Weight } from "lucide-react";
import { SiteHeader } from "@/components/SiteHeader";
import { StatusBadge } from "@/components/StatusBadge";
import { getOrder } from "@/lib/orders";

export const Route = createFileRoute("/orders/$orderId")({
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

      <section className="border-b border-border bg-gradient-hero">
        <div className="mx-auto max-w-7xl px-6 py-12">
          <Link to="/orders" className="group inline-flex items-center gap-1.5 font-mono text-xs uppercase tracking-widest text-muted-foreground hover:text-foreground">
            <ArrowLeft className="h-3.5 w-3.5 transition-transform group-hover:-translate-x-1" />
            All orders
          </Link>
          <div className="mt-6 flex flex-wrap items-start justify-between gap-6">
            <div>
              <div className="flex items-center gap-3">
                <span className="font-mono text-sm text-primary">{order.id}</span>
                <StatusBadge status={order.status} />
              </div>
              <h1 className="mt-3 font-mono text-4xl font-semibold tracking-tight md:text-5xl">{order.trackingNumber}</h1>
              <div className="mt-4 flex items-center gap-3 font-mono text-sm text-muted-foreground">
                <MapPin className="h-4 w-4 text-primary" />
                <span>{order.origin}</span>
                <span className="text-primary">→</span>
                <span>{order.destination}</span>
              </div>
            </div>
            <div className="rounded-xl border border-border bg-card px-6 py-4 shadow-elegant">
              <div className="font-mono text-xs uppercase tracking-widest text-muted-foreground">Estimated delivery</div>
              <div className="mt-1 font-mono text-2xl">{order.estimatedDelivery}</div>
            </div>
          </div>
        </div>
      </section>

      <section className="mx-auto grid max-w-7xl gap-8 px-6 py-12 lg:grid-cols-[1fr_360px]">
        {/* Timeline */}
        <div>
          <div className="mb-6 flex items-center justify-between">
            <h2 className="text-2xl font-semibold tracking-tight">Tracking timeline</h2>
            <span className="font-mono text-xs text-muted-foreground">{order.progress}% complete</span>
          </div>

          <div className="mb-8 h-1.5 overflow-hidden rounded-full bg-secondary">
            <div className="h-full bg-gradient-accent transition-all" style={{ width: `${order.progress}%` }} />
          </div>

          <ol className="relative space-y-6 border-l border-border pl-6">
            {order.events.map((e, i) => (
              <li key={i} className="relative">
                <span className={`absolute -left-[31px] flex h-4 w-4 items-center justify-center rounded-full ring-4 ring-background ${i === 0 ? "bg-primary shadow-glow" : "bg-muted-foreground/50"}`}>
                  {i === 0 && <span className="h-1.5 w-1.5 animate-pulse rounded-full bg-primary-foreground" />}
                </span>
                <div className="rounded-lg border border-border bg-card p-4">
                  <div className="flex items-center justify-between">
                    <span className="font-mono text-sm font-medium">{e.status}</span>
                    <span className="font-mono text-xs text-muted-foreground">{e.timestamp}</span>
                  </div>
                  <div className="mt-1 flex items-center gap-1.5 font-mono text-xs text-primary">
                    <MapPin className="h-3 w-3" />
                    {e.location}
                  </div>
                  {e.description && (
                    <p className="mt-2 text-sm text-muted-foreground">{e.description}</p>
                  )}
                </div>
              </li>
            ))}
          </ol>
        </div>

        {/* Sidebar */}
        <aside className="space-y-6">
          <div className="rounded-xl border border-border bg-card p-6 shadow-elegant">
            <h3 className="font-mono text-xs uppercase tracking-widest text-muted-foreground">Shipment details</h3>
            <dl className="mt-4 space-y-4">
              {meta.map((m) => (
                <div key={m.label} className="flex items-start gap-3">
                  <m.icon className="mt-0.5 h-4 w-4 text-primary" />
                  <div className="flex-1">
                    <dt className="font-mono text-[11px] uppercase tracking-wider text-muted-foreground">{m.label}</dt>
                    <dd className="mt-0.5 font-mono text-sm">{m.value}</dd>
                  </div>
                </div>
              ))}
            </dl>
          </div>

          <div className="rounded-xl border border-border bg-gradient-hero p-6">
            <h3 className="font-mono text-xs uppercase tracking-widest text-primary">Route</h3>
            <div className="mt-4 space-y-3 font-mono text-sm">
              <div className="flex items-center gap-3">
                <div className="h-2 w-2 rounded-full bg-primary" />
                <div>
                  <div className="text-xs text-muted-foreground">Origin</div>
                  <div>{order.origin}</div>
                </div>
              </div>
              <div className="ml-1 h-6 border-l-2 border-dashed border-primary/40" />
              <div className="flex items-center gap-3">
                <div className="h-2 w-2 rounded-full bg-accent" />
                <div>
                  <div className="text-xs text-muted-foreground">Destination</div>
                  <div>{order.destination}</div>
                </div>
              </div>
            </div>
          </div>
        </aside>
      </section>
    </div>
  );
}
