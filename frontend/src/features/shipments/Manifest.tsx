"use client";

import { useEffect, useRef, useState } from "react";
import {
  ArrowRight,
  ChevronLeft,
  ChevronRight,
  Eye,
  Pencil,
  Plus,
  Search,
  Trash2,
} from "lucide-react";
import { useAuth } from "@/features/auth/AuthProvider";
import type { Shipment, ShipmentPage } from "@/lib/api/preview";
import styles from "./manifest.module.css";

const pageSize = 10;
const filters = [
  ["all", "All"],
  ["pending", "Pending"],
  ["in_transit", "In Transit"],
  ["out_for_delivery", "Out for Delivery"],
  ["delivered", "Delivered"],
  ["delayed", "Delayed"],
] as const;
const statusLabels: Record<string, string> = {
  pending: "Pending",
  picked_up: "Picked Up",
  departed: "Departed",
  in_transit: "In Transit",
  out_for_delivery: "Out for Delivery",
  delivered: "Delivered",
  delayed: "Delayed",
};

async function getManifest(
  page: number,
  search: string,
  status: string,
  signal?: AbortSignal,
): Promise<ShipmentPage> {
  const params = new URLSearchParams({ page: String(page), limit: String(pageSize) });
  if (search) params.set("search", search);
  if (status !== "all") params.set("status", status);
  const response = await fetch(`/api/shipments?${params}`, { cache: "no-store", signal });
  const result = await response.json();
  if (!response.ok || !result.success) {
    throw new Error(result.error || `Shipment API returned ${response.status}`);
  }
  return { data: result.data, pagination: result.pagination };
}

function statusLabel(status: string) {
  return statusLabels[status] ?? status.replaceAll("_", " ");
}

function formatDate(value: string) {
  return new Date(value).toLocaleDateString("en-US", {
    month: "short",
    day: "numeric",
    year: "numeric",
  });
}

export default function Manifest({
  initialData,
  initialError,
  initialMessage,
}: {
  initialData: ShipmentPage | null;
  initialError: string | null;
  initialMessage?: string;
}) {
  const { isStaff } = useAuth();
  const [data, setData] = useState(initialData);
  const [page, setPage] = useState(1);
  const [search, setSearch] = useState("");
  const [debouncedSearch, setDebouncedSearch] = useState("");
  const [status, setStatus] = useState("all");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(initialError);
  const [message, setMessage] = useState(initialMessage ?? "");
  const [deleteTarget, setDeleteTarget] = useState<Shipment | null>(null);
  const [deleting, setDeleting] = useState(false);
  const dialog = useRef<HTMLDialogElement>(null);
  const firstFetch = useRef(true);

  useEffect(() => {
    const timer = window.setTimeout(() => {
      setDebouncedSearch(search);
    }, 300);
    return () => window.clearTimeout(timer);
  }, [search]);
  useEffect(() => {
    if (initialMessage) setMessage(initialMessage);
  }, [initialMessage]);

  useEffect(() => {
    if (firstFetch.current && initialData && !initialError) {
      firstFetch.current = false;
      return;
    }
    firstFetch.current = false;
    const controller = new AbortController();
    setLoading(true);
    setError(null);
    getManifest(page, debouncedSearch, status, controller.signal)
      .then(setData)
      .catch((reason: unknown) => {
        if (!controller.signal.aborted)
          setError(reason instanceof Error ? reason.message : "Failed to load shipments");
      })
      .finally(() => {
        if (!controller.signal.aborted) setLoading(false);
      });
    return () => controller.abort();
  }, [page, debouncedSearch, status, initialData, initialError]);

  useEffect(() => {
    if (deleteTarget && !dialog.current?.open) dialog.current?.showModal();
    else if (!deleteTarget && dialog.current?.open) dialog.current.close();
  }, [deleteTarget]);

  function changeSearch(value: string) {
    setSearch(value);
    setPage(1);
  }

  function changeStatus(value: string) {
    setStatus(value);
    setPage(1);
  }

  async function confirmDelete() {
    if (!deleteTarget || !isStaff) return;
    setDeleting(true);
    setError(null);
    setMessage("");
    try {
      const response = await fetch(`/api/shipments/${encodeURIComponent(deleteTarget.id)}`, {
        method: "DELETE",
        credentials: "same-origin",
        cache: "no-store",
      });
      const result = await response.json();
      if (!response.ok || result.success === false) {
        throw new Error(result.error || `Delete failed (${response.status})`);
      }
      setDeleteTarget(null);
      setMessage("Shipment deleted.");
      const refreshed = await getManifest(page, debouncedSearch, status);
      setData(refreshed);
      if (page > refreshed.pagination.totalPages)
        setPage(Math.max(1, refreshed.pagination.totalPages));
    } catch (reason) {
      setError(reason instanceof Error ? reason.message : "Failed to delete shipment");
    } finally {
      setDeleting(false);
    }
  }

  const pageData = data?.data ?? [];
  const pagination = data?.pagination;
  const total = pagination?.total ?? 0;
  const totalPages = pagination?.totalPages ?? 1;

  return (
    <main>
      <section className={styles.hero}>
        <div className={styles.container}>
          <div>
            <span className={styles.kicker}>/ orders</span>
            <h1>Shipment manifest</h1>
            <p>{total} total shipments tracked across all carriers.</p>
          </div>
          {isStaff && (
            <a className={styles.newOrder} href="/orders/create">
              <Plus size={16} /> New Order
            </a>
          )}
        </div>
      </section>
      <section className={styles.container}>
        <div className={styles.filters}>
          <label className={styles.search}>
            <Search size={16} aria-hidden="true" />
            <input
              aria-label="Search shipments"
              type="search"
              placeholder="Search by ID, tracking, customer, destination"
              value={search}
              onChange={(event) => changeSearch(event.target.value)}
            />
          </label>
          <div className={styles.statusFilters} aria-label="Filter by status">
            {filters.map(([key, label]) => (
              <button
                key={key}
                type="button"
                aria-pressed={status === key}
                onClick={() => changeStatus(key)}
              >
                {label}
              </button>
            ))}
          </div>
        </div>

        {error && !deleteTarget && (
          <p className={styles.error} role="alert">
            {error}
          </p>
        )}
        {message && (
          <p className={styles.loading} role="status">
            {message}
          </p>
        )}
        {loading && (
          <p className={styles.loading} role="status">
            Loading shipments…
          </p>
        )}
        <div className={styles.tableCard} aria-busy={loading}>
          <div className={styles.tableScroll}>
            <table>
              <thead>
                <tr>
                  <th>Order ID</th>
                  <th>Tracking</th>
                  <th>Customer</th>
                  <th>Route</th>
                  <th>Status</th>
                  <th>ETA</th>
                  <th>Actions</th>
                </tr>
              </thead>
              <tbody>
                {pageData.map((shipment) => (
                  <tr key={shipment.id}>
                    <td data-label="Order ID">
                      <a
                        className={styles.orderId}
                        href={`/orders/${encodeURIComponent(shipment.id)}`}
                      >
                        {shipment.id}
                      </a>
                    </td>
                    <td data-label="Tracking" className={styles.mono}>
                      {shipment.trackingNumber}
                    </td>
                    <td data-label="Customer">{shipment.customer.name}</td>
                    <td data-label="Route">
                      <span className={styles.route}>
                        {shipment.origin}
                        <ArrowRight size={13} aria-hidden="true" />
                        {shipment.destination}
                      </span>
                    </td>
                    <td data-label="Status">
                      <span
                        className={`${styles.badge} ${styles[`status_${shipment.status}`] ?? ""}`}
                      >
                        <i aria-hidden="true" />
                        {statusLabel(shipment.status)}
                      </span>
                    </td>
                    <td data-label="ETA" className={styles.mono}>
                      {formatDate(shipment.estimatedDelivery)}
                    </td>
                    <td data-label="Actions">
                      <div className={styles.actions}>
                        <a
                          href={`/orders/${encodeURIComponent(shipment.id)}`}
                          aria-label={`View shipment ${shipment.id}`}
                          title="View details"
                        >
                          <Eye size={16} />
                        </a>
                        {isStaff && (
                          <a
                            href={`/orders/${encodeURIComponent(shipment.id)}/edit`}
                            aria-label={`Edit shipment ${shipment.id}`}
                            title="Edit shipment"
                          >
                            <Pencil size={16} />
                          </a>
                        )}
                        {isStaff && (
                          <button
                            type="button"
                            aria-label={`Delete shipment ${shipment.id}`}
                            title="Delete shipment"
                            onClick={() => {
                              setError(null);
                              setDeleteTarget(shipment);
                            }}
                          >
                            <Trash2 size={16} />
                          </button>
                        )}
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
            {!loading && !error && pageData.length === 0 && (
              <p className={styles.empty}>No shipments match your filters.</p>
            )}
          </div>
          {totalPages > 1 && (
            <nav className={styles.pagination} aria-label="Shipment pages">
              <span>
                Showing {Math.min((page - 1) * pageSize + 1, total)}–
                {Math.min(page * pageSize, total)} of {total}
              </span>
              <div>
                <button
                  type="button"
                  aria-label="Previous page"
                  disabled={page <= 1}
                  onClick={() => setPage((current) => current - 1)}
                >
                  <ChevronLeft size={16} />
                </button>
                <span aria-current="page">
                  {page} / {totalPages}
                </span>
                <button
                  type="button"
                  aria-label="Next page"
                  disabled={page >= totalPages}
                  onClick={() => setPage((current) => current + 1)}
                >
                  <ChevronRight size={16} />
                </button>
              </div>
            </nav>
          )}
        </div>
        <p className={styles.currentStatus}>
          Status: {status === "all" ? "All" : statusLabel(status)}
        </p>
      </section>
      <dialog
        ref={dialog}
        className={styles.dialog}
        aria-labelledby="delete-heading"
        onCancel={(event) => {
          event.preventDefault();
          setDeleteTarget(null);
          setError(null);
        }}
      >
        <h2 id="delete-heading">Delete shipment?</h2>
        {deleteTarget && error && (
          <p className={styles.error} role="alert">
            {error}
          </p>
        )}
        <div>
          <button
            type="button"
            disabled={deleting}
            onClick={() => {
              setDeleteTarget(null);
              setError(null);
            }}
          >
            Cancel
          </button>
          <button type="button" disabled={deleting} onClick={confirmDelete}>
            {deleting ? "Deleting…" : "Delete shipment"}
          </button>
        </div>
      </dialog>
    </main>
  );
}
