import type { Shipment } from "./types";
import { ShipmentForm } from "./ShipmentForm";

export function ShipmentFormPage({ initial }: { initial?: Shipment }) {
  const editing = !!initial;
  return (
    <main>
      <section className="border-b border-border bg-card/50">
        <div className="mx-auto max-w-7xl px-4 py-12 sm:px-6">
          <span className="font-mono text-sm text-primary">/ orders</span>
          <h1 className="mt-3 text-3xl font-bold tracking-tight">
            {editing ? "Edit Shipment" : "Create Shipment"}
          </h1>
          <p className="mt-2 text-muted-foreground">
            {editing ? "Update shipment details." : "Register a new shipment in the system."}
          </p>
        </div>
      </section>
      <section className="mx-auto max-w-5xl px-4 py-8 sm:px-6">
        <ShipmentForm initial={initial} />
      </section>
    </main>
  );
}
