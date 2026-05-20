import { createFileRoute, Link } from "@tanstack/react-router";
import { useMemo, useState } from "react";
import { ArrowRight, Search, Filter } from "lucide-react";
import { SiteHeader } from "@/components/SiteHeader";
import { Input } from "@/components/ui/input";
import { orders, statusLabels, type ShipmentStatus } from "@/lib/orders";
import { StatusBadge } from "@/components/StatusBadge";
import { cn } from "@/lib/utils";

export const Route = createFileRoute("/orders")({
  head: () => ({
    meta: [
      { title: "Orders — Harbor Ops" },
      { name: "description", content: "Browse and filter all active and completed shipments across the fleet." },
      { property: "og:title", content: "Orders — Harbor Ops" },
      { property: "og:description", content: "Browse all active and completed shipments." },
    ],
  }),
  component: OrdersPage,
});

const FILTERS: Array<{ key: ShipmentStatus | "all"; label: string }> = [
  { key: "all", label: "All" },
  { key: "pending", label: "Pending" },
  { key: "in_transit", label: "In Transit" },
  { key: "out_for_delivery", label: "Out for Delivery" },
  { key: "delivered", label: "Delivered" },
  { key: "delayed", label: "Delayed" },
];

function OrdersPage() {
  const [filter, setFilter] = useState<ShipmentStatus | "all">("all");
  const [query, setQuery] = useState("");

  const filtered = useMemo(() => {
    const q = query.trim().toLowerCase();
    return orders.filter((o) => {
      if (filter !== "all" && o.status !== filter) return false;
      if (!q) return true;
      return (
        o.id.toLowerCase().includes(q) ||
        o.trackingNumber.toLowerCase().includes(q) ||
        o.customer.toLowerCase().includes(q) ||
        o.destination.toLowerCase().includes(q)
      );
    });
  }, [filter, query]);

  return (
    <div className="min-h-screen">
      <SiteHeader />

      <section className="border-b border-border bg-gradient-hero">
        <div className="mx-auto max-w-7xl px-6 py-14">
          <span className="font-mono text-xs uppercase tracking-widest text-primary">/ orders</span>
          <h1 className="mt-3 text-4xl font-semibold tracking-tight md:text-5xl">Shipment manifest</h1>
          <p className="mt-3 max-w-2xl text-muted-foreground">
            {orders.length} total shipments tracked across all carriers.
          </p>
        </div>
      </section>

      <section className="mx-auto max-w-7xl px-6 py-10">
        {/* Controls */}
        <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
          <div className="flex items-center gap-2 rounded-lg border border-border bg-card px-3 md:w-96">
            <Search className="h-4 w-4 text-muted-foreground" />
            <Input
              value={query}
              onChange={(e) => setQuery(e.target.value)}
              placeholder="Search by ID, tracking, customer, destination"
              className="h-11 border-0 bg-transparent font-mono text-sm shadow-none focus-visible:ring-0"
            />
          </div>
          <div className="flex items-center gap-2 overflow-x-auto">
            <Filter className="h-4 w-4 shrink-0 text-muted-foreground" />
            {FILTERS.map((f) => (
              <button
                key={f.key}
                onClick={() => setFilter(f.key)}
                className={cn(
                  "rounded-full border px-3 py-1.5 font-mono text-xs uppercase tracking-wider transition-colors",
                  filter === f.key
                    ? "border-primary bg-primary/15 text-primary"
                    : "border-border text-muted-foreground hover:text-foreground",
                )}
              >
                {f.label}
              </button>
            ))}
          </div>
        </div>

        {/* Table */}
        <div className="mt-8 overflow-hidden rounded-xl border border-border bg-card shadow-elegant">
          <div className="hidden grid-cols-[1.1fr_1.4fr_1.6fr_2fr_1.2fr_0.6fr] gap-4 border-b border-border bg-secondary/50 px-6 py-3 font-mono text-[11px] uppercase tracking-widest text-muted-foreground md:grid">
            <span>Order ID</span>
            <span>Tracking</span>
            <span>Customer</span>
            <span>Route</span>
            <span>Status</span>
            <span className="text-right">ETA</span>
          </div>

          {filtered.length === 0 ? (
            <div className="px-6 py-16 text-center font-mono text-sm text-muted-foreground">
              No shipments match your filters.
            </div>
          ) : (
            filtered.map((o) => (
              <Link
                key={o.id}
                to="/orders/$orderId"
                params={{ orderId: o.id }}
                className="group grid grid-cols-1 gap-2 border-b border-border px-6 py-4 transition-colors last:border-0 hover:bg-secondary/40 md:grid-cols-[1.1fr_1.4fr_1.6fr_2fr_1.2fr_0.6fr] md:items-center md:gap-4"
              >
                <span className="font-mono text-sm text-primary">{o.id}</span>
                <span className="font-mono text-sm text-muted-foreground">{o.trackingNumber}</span>
                <span className="text-sm">{o.customer}</span>
                <span className="flex items-center gap-2 font-mono text-xs text-muted-foreground">
                  <span>{o.origin}</span>
                  <ArrowRight className="h-3 w-3 text-primary" />
                  <span>{o.destination}</span>
                </span>
                <span><StatusBadge status={o.status} /></span>
                <span className="font-mono text-xs text-muted-foreground md:text-right">{o.estimatedDelivery}</span>
              </Link>
            ))
          )}
        </div>

        <div className="mt-4 font-mono text-xs text-muted-foreground">
          Showing {filtered.length} of {orders.length} · Status: {filter === "all" ? "All" : statusLabels[filter]}
        </div>
      </section>
    </div>
  );
}
