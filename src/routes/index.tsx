import { createFileRoute, Link, useNavigate } from "@tanstack/react-router";
import { useState } from "react";
import { ArrowRight, Boxes, Search, Truck, Globe2, Activity } from "lucide-react";
import { SiteHeader } from "@/components/SiteHeader";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { orders } from "@/lib/orders";
import { StatusBadge } from "@/components/StatusBadge";

export const Route = createFileRoute("/")({
  head: () => ({
    meta: [
      { title: "Harbor Ops — Real-time Shipment Tracking" },
      { name: "description", content: "Track global shipments, monitor fleet status, and manage orders from a single ops console." },
      { property: "og:title", content: "Harbor Ops — Real-time Shipment Tracking" },
      { property: "og:description", content: "Track global shipments and manage orders from a single console." },
    ],
  }),
  component: Home,
});

function Home() {
  const navigate = useNavigate();
  const [query, setQuery] = useState("");

  const onTrack = (e: React.FormEvent) => {
    e.preventDefault();
    const q = query.trim().toLowerCase();
    const match = orders.find(
      (o) => o.id.toLowerCase() === q || o.trackingNumber.toLowerCase() === q,
    );
    if (match) navigate({ to: "/orders/$orderId", params: { orderId: match.id } });
    else navigate({ to: "/orders" });
  };

  const stats = [
    { label: "Active shipments", value: orders.filter((o) => o.status !== "delivered").length, icon: Truck },
    { label: "Delivered (30d)", value: 184, icon: Boxes },
    { label: "On-time rate", value: "97.4%", icon: Activity },
    { label: "Countries served", value: 42, icon: Globe2 },
  ];

  const recent = orders.slice(0, 3);

  return (
    <div className="min-h-screen">
      <SiteHeader />

      {/* Hero */}
      <section className="relative overflow-hidden bg-gradient-hero">
        <div className="absolute inset-0 opacity-[0.04]" style={{
          backgroundImage: "linear-gradient(to right, currentColor 1px, transparent 1px), linear-gradient(to bottom, currentColor 1px, transparent 1px)",
          backgroundSize: "48px 48px",
        }} />
        <div className="relative mx-auto max-w-7xl px-6 py-24 md:py-32">
          <div className="max-w-3xl">
            <span className="inline-flex items-center gap-2 rounded-full border border-border bg-background/40 px-3 py-1 font-mono text-xs uppercase tracking-widest text-primary">
              <span className="h-1.5 w-1.5 animate-pulse rounded-full bg-primary" />
              Live ops console
            </span>
            <h1 className="mt-6 text-5xl font-semibold leading-[1.05] tracking-tight md:text-7xl">
              Move cargo with<br />
              <span className="bg-gradient-accent bg-clip-text text-transparent">radar precision.</span>
            </h1>
            <p className="mt-6 max-w-xl font-sans text-lg text-muted-foreground">
              Trace every parcel from origin warehouse to doorstep. Realtime telemetry, port-side timelines, and exception alerts in one console.
            </p>

            <form onSubmit={onTrack} className="mt-10 flex max-w-xl gap-2 rounded-xl border border-border bg-card p-2 shadow-elegant">
              <div className="flex flex-1 items-center gap-2 px-3">
                <Search className="h-4 w-4 text-muted-foreground" />
                <Input
                  value={query}
                  onChange={(e) => setQuery(e.target.value)}
                  placeholder="Enter tracking # (try TRK-9F2A-44B1)"
                  className="h-10 border-0 bg-transparent font-mono text-sm shadow-none focus-visible:ring-0"
                />
              </div>
              <Button type="submit" className="h-10 gap-2">
                Track <ArrowRight className="h-4 w-4" />
              </Button>
            </form>
          </div>
        </div>
      </section>

      {/* Stats grid */}
      <section className="border-y border-border bg-card/40">
        <div className="mx-auto grid max-w-7xl grid-cols-2 divide-border md:grid-cols-4 md:divide-x">
          {stats.map((s) => (
            <div key={s.label} className="flex items-center gap-4 px-6 py-8">
              <div className="flex h-11 w-11 items-center justify-center rounded-lg bg-secondary text-primary">
                <s.icon className="h-5 w-5" />
              </div>
              <div>
                <div className="font-mono text-2xl font-semibold">{s.value}</div>
                <div className="text-xs uppercase tracking-wider text-muted-foreground">{s.label}</div>
              </div>
            </div>
          ))}
        </div>
      </section>

      {/* Recent shipments */}
      <section className="mx-auto max-w-7xl px-6 py-20">
        <div className="flex items-end justify-between">
          <div>
            <h2 className="text-3xl font-semibold tracking-tight">Recent shipments</h2>
            <p className="mt-2 text-sm text-muted-foreground">Latest activity from the fleet.</p>
          </div>
          <Link to="/orders" className="group flex items-center gap-1 font-mono text-sm text-primary">
            View all <ArrowRight className="h-4 w-4 transition-transform group-hover:translate-x-1" />
          </Link>
        </div>

        <div className="mt-8 grid gap-4 md:grid-cols-3">
          {recent.map((o) => (
            <Link
              key={o.id}
              to="/orders/$orderId"
              params={{ orderId: o.id }}
              className="group rounded-xl border border-border bg-card p-6 transition-all hover:border-primary/40 hover:shadow-elegant"
            >
              <div className="flex items-center justify-between">
                <span className="font-mono text-xs text-muted-foreground">{o.id}</span>
                <StatusBadge status={o.status} />
              </div>
              <div className="mt-4 font-mono text-lg">{o.trackingNumber}</div>
              <div className="mt-1 text-sm text-muted-foreground">{o.customer}</div>
              <div className="mt-6 flex items-center gap-2 font-mono text-xs">
                <span>{o.origin}</span>
                <ArrowRight className="h-3 w-3 text-primary" />
                <span>{o.destination}</span>
              </div>
              <div className="mt-4 h-1 overflow-hidden rounded-full bg-secondary">
                <div className="h-full bg-gradient-accent transition-all" style={{ width: `${o.progress}%` }} />
              </div>
            </Link>
          ))}
        </div>
      </section>

      <footer className="border-t border-border py-8">
        <div className="mx-auto max-w-7xl px-6 font-mono text-xs text-muted-foreground">
          © 2026 Harbor Ops · Logistics control plane
        </div>
      </footer>
    </div>
  );
}
