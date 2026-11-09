"use client";

import { useCallback, useEffect, useMemo, useState } from "react";
import { mapShipment } from "@/features/shipments/mappers";
import { statusLabels, type Shipment, type ShipmentDTO, type ShipmentStatus } from "@/features/shipments/types";
import type { Hub } from "@/features/hubs/types";
import { useAuthStore } from "@/features/auth/store/auth-store";
import { DeliveryMap } from "@/features/operations/DeliveryMap";
import { Check, ChevronLeft, ChevronRight, Eye, RefreshCw, Search } from "lucide-react";
import { Button } from "@/components/ui/button";

const pageSize = 10;
const statuses: ShipmentStatus[] = [
  "pending",
  "picked_up",
  "departed",
  "in_transit",
  "out_for_delivery",
  "delayed",
  "delivered",
];
const hubRequiredStatuses: Record<ShipmentStatus, boolean> = {
  pending: false,
  picked_up: false,
  departed: true,
  in_transit: true,
  out_for_delivery: true,
  delayed: true,
  delivered: false,
};

async function fetchDeliveries(): Promise<Shipment[]> {
  const response = await fetch("/api/shipments?limit=-1&exclude_status=delivered", {
    cache: "no-store",
  });
  const result = await response.json();
  return (result.data as ShipmentDTO[]).map(mapShipment);
}

export function DeliveriesPanel({
  initialDeliveries,
  hubs,
  initialError,
}: Readonly<{ initialDeliveries: Shipment[]; hubs: Hub[]; initialError: boolean }>) {
  const isStaff = useAuthStore((state) => state.isStaff);
  const [deliveries, setDeliveries] = useState(initialDeliveries);
  const [search, setSearch] = useState("");
  const [page, setPage] = useState(1);
  const [selectedId, setSelectedId] = useState<string | null>(null);
  const [draftStatus, setDraftStatus] = useState<Record<string, ShipmentStatus>>({});
  const [draftHubId, setDraftHubId] = useState<Record<string, string>>({});
  const [updatedAt, setUpdatedAt] = useState(Date.now());
  const [loading, setLoading] = useState(initialError);
  const [error, setError] = useState(initialError ? "Failed to load active deliveries." : "");
  const [feedback, setFeedback] = useState("");
  const [now, setNow] = useState(Date.now());

  const refresh = useCallback(async () => {
    setLoading(true);
    try {
      setDeliveries(await fetchDeliveries());
      setUpdatedAt(Date.now());
      setError("");
    } catch {
      setError("Failed to load active deliveries.");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    const interval = window.setInterval(refresh, 15_000);
    return () => window.clearInterval(interval);
  }, [refresh]);
  useEffect(() => {
    const interval = window.setInterval(() => setNow(Date.now()), 1000);
    return () => window.clearInterval(interval);
  }, []);

  const filtered = useMemo(() => {
    const query = search.trim().toLowerCase();
    if (!query) return deliveries;
    return deliveries.filter((delivery) =>
      [delivery.id, delivery.trackingNumber, delivery.customer.name, delivery.carrier].some(
        (value) => value.toLowerCase().includes(query),
      ),
    );
  }, [deliveries, search]);
  const totalPages = Math.max(1, Math.ceil(filtered.length / pageSize));
  const pageItems = filtered.slice((page - 1) * pageSize, page * pageSize);
  useEffect(() => {
    setPage((current) => Math.min(current, totalPages));
  }, [totalPages]);
  const selected = deliveries.find((delivery) => delivery.id === selectedId);
  const mapDelivery =
    selected && filtered.some((delivery) => delivery.id === selected.id) ? selected : null;
  const activeHubs = useMemo(() => hubs.filter((hub) => hub.status === "active"), [hubs]);

  async function updateStatus(order: Shipment) {
    const status = draftStatus[order.id] ?? order.status;
    const hubId = draftHubId[order.id] ?? order.hubId ?? "";
    if (hubRequiredStatuses[status] && !hubId) {
      setFeedback("Select an active hub for this status.");
      return;
    }
    setFeedback("");
    try {
      const response = await fetch(`/api/shipments/${encodeURIComponent(order.id)}/status`, {
        method: "PATCH",
        headers: { "content-type": "application/json" },
        body: JSON.stringify({ status, ...(hubId ? { hubId } : {}) }),
      });
      const result = await response.json();
      if (!response.ok) throw new Error(result.error || `Request failed (${response.status})`);
      const updated = mapShipment(result.data as ShipmentDTO);
      setDeliveries((current) =>
        updated.status === "delivered"
          ? current.filter((delivery) => delivery.id !== updated.id)
          : current.map((delivery) => (delivery.id === updated.id ? updated : delivery)),
      );
      setDraftStatus((current) => {
        const next = { ...current };
        delete next[order.id];
        return next;
      });
      setDraftHubId((current) => {
        const next = { ...current };
        delete next[order.id];
        return next;
      });
      setUpdatedAt(Date.now());
      setFeedback("Delivery status updated.");
    } catch (cause) {
      setFeedback(cause instanceof Error ? cause.message : "Failed to update delivery status.");
    }
  }

  return (
    <section aria-label="Active deliveries">
      {!isStaff && (
        <p className="mb-3.5 font-mono text-xs text-muted-foreground">
          Public delivery view. Sign in as staff to change delivery status.
        </p>
      )}
      <div className="mb-5 flex items-center justify-between gap-4 max-sm:flex-col max-sm:items-stretch">
        <label className="relative block w-full max-w-sm">
          <span className="sr-only">Search active deliveries</span>
          <Search
            className="pointer-events-none absolute top-1/2 left-3 size-4 -translate-y-1/2 text-muted-foreground"
            aria-hidden="true"
          />
          <input
            className="h-11 w-full rounded-lg border border-border bg-background pl-10 pr-3 text-foreground placeholder:text-muted-foreground focus-visible:outline-2 focus-visible:outline-primary"
            aria-label="Search active deliveries"
            type="search"
            value={search}
            onChange={(event) => {
              setSearch(event.target.value);
              setPage(1);
            }}
            placeholder="Search active shipments..."
          />
        </label>
        <div className="flex items-center gap-3 font-mono text-xs text-muted-foreground">
          <span>
            {deliveries.length} active · updated {Math.max(0, Math.floor((now - updatedAt) / 1000))}
            s ago
          </span>
          <Button
            type="button"
            variant="ghost"
            size="icon-sm"
            onClick={refresh}
            disabled={loading}
            aria-label="Refresh active deliveries"
          >
            <RefreshCw className={loading ? "size-4 animate-spin" : "size-4"} aria-hidden="true" />
          </Button>
        </div>
      </div>
      {error && (
        <div
          className="flex items-center gap-3 rounded-md border border-destructive/40 bg-destructive/10 p-3 text-destructive"
          role="alert"
        >
          {error}
          <Button type="button" variant="outline" onClick={refresh}>
            Retry
          </Button>
        </div>
      )}
      {feedback && (
        <p
          className={`mb-3.5 font-mono text-xs ${feedback.includes("updated") ? "text-emerald-400" : "text-destructive"}`}
          role={feedback.includes("updated") ? "status" : "alert"}
        >
          {feedback}
        </p>
      )}
      <div className="overflow-x-auto rounded-lg border border-border">
        <table className="w-full min-w-[1050px] border-collapse text-left text-sm">
          <thead>
            <tr>
              <th className="bg-secondary px-3 py-3 font-mono text-[10px] uppercase tracking-wider text-muted-foreground">
                Order ID
              </th>
              <th className="bg-secondary px-3 py-3 font-mono text-[10px] uppercase tracking-wider text-muted-foreground">
                Tracking
              </th>
              <th className="bg-secondary px-3 py-3 font-mono text-[10px] uppercase tracking-wider text-muted-foreground">
                Customer
              </th>
              <th className="bg-secondary px-3 py-3 font-mono text-[10px] uppercase tracking-wider text-muted-foreground">
                Carrier
              </th>
              <th className="bg-secondary px-3 py-3 font-mono text-[10px] uppercase tracking-wider text-muted-foreground">
                Status
              </th>
              <th className="bg-secondary px-3 py-3 font-mono text-[10px] uppercase tracking-wider text-muted-foreground">
                Hub
              </th>
              <th className="bg-secondary px-3 py-3 font-mono text-[10px] uppercase tracking-wider text-muted-foreground">
                ETA
              </th>
              <th className="bg-secondary px-3 py-3 font-mono text-[10px] uppercase tracking-wider text-muted-foreground">
                Actions
              </th>
            </tr>
          </thead>
          <tbody>
            {pageItems.map((order) => {
              const hub = hubs.find((item) => item.id === order.hubId);
              return (
                <tr
                  className="cursor-pointer border-b border-border last:border-0 hover:bg-primary/5 aria-selected:bg-primary/5"
                  key={order.id}
                  onClick={() => setSelectedId(order.id)}
                  aria-selected={mapDelivery?.id === order.id}
                >
                  <td className="whitespace-nowrap px-3 py-3.5 font-mono font-semibold text-primary">
                    {order.id}
                  </td>
                  <td className="whitespace-nowrap px-3 py-3.5 font-mono text-muted-foreground">
                    {order.trackingNumber}
                  </td>
                  <td className="whitespace-nowrap px-3 py-3.5 font-mono font-semibold text-foreground">
                    {order.customer.name}
                  </td>
                  <td className="whitespace-nowrap px-3 py-3.5 font-mono text-muted-foreground">
                    {order.carrier}
                  </td>
                  <td className="px-3 py-2.5">
                    {isStaff ? (
                      <select
                        className="h-11 max-w-[180px] rounded-lg border border-border bg-background px-2.5 font-mono text-sm text-foreground"
                        aria-label={`Status for ${order.id}`}
                        value={draftStatus[order.id] ?? order.status}
                        onChange={(event) =>
                          setDraftStatus((current) => ({
                            ...current,
                            [order.id]: event.target.value as ShipmentStatus,
                          }))
                        }
                      >
                        {statuses.map((status) => (
                          <option key={status} value={status}>
                            {statusLabels[status]}
                          </option>
                        ))}
                      </select>
                    ) : (
                      <span className="font-mono">{statusLabels[order.status]}</span>
                    )}
                  </td>
                  <td className="px-3 py-2.5">
                    {isStaff && hubRequiredStatuses[draftStatus[order.id] ?? order.status] ? (
                      <select
                        className="h-11 min-w-40 max-w-[190px] rounded-lg border border-border bg-background px-2.5 font-mono text-sm text-foreground"
                        aria-label={`Hub for ${order.id}`}
                        value={draftHubId[order.id] ?? order.hubId ?? ""}
                        onChange={(event) =>
                          setDraftHubId((current) => ({
                            ...current,
                            [order.id]: event.target.value,
                          }))
                        }
                      >
                        <option value="">Select hub</option>
                        {activeHubs.map((item) => (
                          <option key={item.id} value={item.id}>
                            {item.name}
                          </option>
                        ))}
                      </select>
                    ) : (
                      <span className="font-mono text-muted-foreground">{hub?.name ?? "—"}</span>
                    )}
                  </td>
                  <td className="whitespace-nowrap px-3 py-3.5 font-mono text-muted-foreground">
                    {order.estimatedDelivery || "—"}
                  </td>
                  <td className="px-2 py-2">
                    <div className="flex items-center gap-1">
                      <Button
                        type="button"
                        variant="ghost"
                        size="icon-sm"
                        aria-label={`View ${order.id} location`}
                        onClick={(event) => {
                          event.stopPropagation();
                          setSelectedId(order.id);
                        }}
                      >
                        <Eye className="size-4" aria-hidden="true" />
                      </Button>
                      {isStaff && (
                        <Button
                          type="button"
                          variant="ghost"
                          size="icon-sm"
                          aria-label={`Save ${order.id} status`}
                          disabled={
                            (draftStatus[order.id] ?? order.status) === order.status &&
                            (draftHubId[order.id] ?? order.hubId ?? "") === (order.hubId ?? "")
                          }
                          onClick={(event) => {
                            event.stopPropagation();
                            void updateStatus(order);
                          }}
                        >
                          <Check className="size-4" aria-hidden="true" />
                        </Button>
                      )}
                    </div>
                  </td>
                </tr>
              );
            })}
          </tbody>
        </table>
        {filtered.length === 0 && (
          <p className="p-6 text-center text-muted-foreground">
            No active deliveries match your filters.
          </p>
        )}
      </div>
      <nav
        className="flex items-center justify-between gap-3 px-2 py-2 font-mono text-xs text-muted-foreground max-sm:flex-col max-sm:items-stretch"
        aria-label="Active delivery pages"
      >
        <span>
          Showing {filtered.length ? (page - 1) * pageSize + 1 : 0}–
          {Math.min(page * pageSize, filtered.length)} of {filtered.length}
        </span>
        <div className="flex items-center gap-1">
          <Button
            type="button"
            variant="ghost"
            size="icon-sm"
            aria-label="Previous page"
            disabled={page <= 1}
            onClick={() => setPage((current) => current - 1)}
          >
            <ChevronLeft className="size-4" aria-hidden="true" />
          </Button>
          <span aria-current="page" className="px-2 text-primary">
            {page}
          </span>
          <Button
            type="button"
            variant="ghost"
            size="icon-sm"
            aria-label="Next page"
            disabled={page >= totalPages}
            onClick={() => setPage((current) => current + 1)}
          >
            <ChevronRight className="size-4" aria-hidden="true" />
          </Button>
        </div>
      </nav>
      {mapDelivery && (
        <div className="mt-4 h-80 overflow-hidden rounded-lg border border-border">
          <DeliveryMap
            key={mapDelivery.id}
            order={mapDelivery}
            hub={hubs.find((item) => item.id === mapDelivery.hubId)}
          />
        </div>
      )}
    </section>
  );
}
