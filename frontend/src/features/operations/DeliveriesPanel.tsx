"use client";

import { useCallback, useEffect, useMemo, useState } from "react";
import { mapShipmentToOrder, type BackendShipment } from "@/lib/api/mappers";
import type { Hub } from "@/lib/hubs";
import { statusLabels, type Order, type ShipmentStatus } from "@/lib/orders";
import { useAuth } from "@/features/auth/AuthProvider";
import { DeliveryMap } from "@/features/operations/DeliveryMap";
import styles from "@/features/operations/DeliveriesPanel.module.css";

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

async function fetchDeliveries(): Promise<Order[]> {
  const response = await fetch("/api/shipments?limit=-1&exclude_status=delivered", {
    cache: "no-store",
  });
  const result = await response.json();
  if (!response.ok) throw new Error(result.error || `Request failed (${response.status})`);
  return (result.data as BackendShipment[]).map(mapShipmentToOrder);
}

export function DeliveriesPanel({
  initialDeliveries,
  hubs,
  initialError,
}: Readonly<{ initialDeliveries: Order[]; hubs: Hub[]; initialError: boolean }>) {
  const { isStaff } = useAuth();
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
    selected && filtered.some((delivery) => delivery.id === selected.id) ? selected : filtered[0];
  const activeHubs = useMemo(() => hubs.filter((hub) => hub.status === "active"), [hubs]);

  async function updateStatus(order: Order) {
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
      const updated = mapShipmentToOrder(result.data as BackendShipment);
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
        <p className={styles.readOnly}>
          Public delivery view. Sign in as staff to change delivery status.
        </p>
      )}
      <div className={styles.controls}>
        <label className={styles.searchLabel}>
          <span>Search active deliveries</span>
          <input
            aria-label="Search active deliveries"
            type="search"
            value={search}
            onChange={(event) => {
              setSearch(event.target.value);
              setPage(1);
            }}
            placeholder="Order, tracking, customer, carrier"
          />
        </label>
        <div className={styles.updated}>
          <span>
            {deliveries.length} active · updated {Math.max(0, Math.floor((now - updatedAt) / 1000))}
            s ago
          </span>
          <button
            type="button"
            onClick={refresh}
            disabled={loading}
            aria-label="Refresh active deliveries"
          >
            {loading ? "Refreshing…" : "Refresh"}
          </button>
        </div>
      </div>
      {error && (
        <div className={styles.error} role="alert">
          {error}
          <button type="button" onClick={refresh}>
            Retry
          </button>
        </div>
      )}
      {feedback && (
        <p className={styles.feedback} role={feedback.includes("updated") ? "status" : "alert"}>
          {feedback}
        </p>
      )}
      <div className={styles.tableWrap}>
        <table className={styles.table}>
          <thead>
            <tr>
              <th>Order ID</th>
              <th>Tracking</th>
              <th>Customer</th>
              <th>Carrier</th>
              <th>Status</th>
              <th>Current hub</th>
              <th>Current location</th>
              <th>ETA</th>
              {isStaff && <th>Actions</th>}
            </tr>
          </thead>
          <tbody>
            {pageItems.map((order) => {
              const hub = hubs.find((item) => item.id === order.hubId);
              return (
                <tr
                  key={order.id}
                  onClick={() => setSelectedId(order.id)}
                  aria-selected={mapDelivery?.id === order.id}
                >
                  <td>{order.id}</td>
                  <td>{order.trackingNumber}</td>
                  <td>{order.customer.name}</td>
                  <td>{order.carrier}</td>
                  <td>
                    {isStaff ? (
                      <select
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
                      statusLabels[order.status].toLowerCase()
                    )}
                  </td>
                  <td>
                    {isStaff && hubRequiredStatuses[draftStatus[order.id] ?? order.status] ? (
                      <select
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
                      (hub?.name ?? "—")
                    )}
                  </td>
                  <td>
                    {order.currentCoords.lat.toFixed(4)}, {order.currentCoords.lng.toFixed(4)}
                  </td>
                  <td>{order.estimatedDelivery || "—"}</td>
                  {isStaff && (
                    <td>
                      <button
                        type="button"
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
                        Save
                      </button>
                    </td>
                  )}
                </tr>
              );
            })}
          </tbody>
        </table>
        {filtered.length === 0 && (
          <p className={styles.empty}>No active deliveries match your filters.</p>
        )}
      </div>
      <nav className={styles.pagination} aria-label="Active delivery pages">
        <span>
          {filtered.length} deliveries · Page {page} of {totalPages}
        </span>
        <div>
          <button
            type="button"
            aria-label="Previous page"
            disabled={page <= 1}
            onClick={() => setPage((current) => current - 1)}
          >
            Previous
          </button>
          <button
            type="button"
            aria-label="Next page"
            disabled={page >= totalPages}
            onClick={() => setPage((current) => current + 1)}
          >
            Next
          </button>
        </div>
      </nav>
      {mapDelivery && (
        <DeliveryMap
          key={mapDelivery.id}
          order={mapDelivery}
          hub={hubs.find((item) => item.id === mapDelivery.hubId)}
        />
      )}
    </section>
  );
}
