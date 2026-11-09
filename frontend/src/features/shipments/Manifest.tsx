"use client";

import Link from "next/link";
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
import { useAuthStore } from "@/features/auth/auth-store";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { statusLabels, type ShipmentStatus } from "@/lib/orders";
import type { Shipment, ShipmentPage } from "@/lib/api/preview";

const pageSize = 10;
const filters = [
  ["all", "All"],
  ["pending", "Pending"],
  ["in_transit", "In Transit"],
  ["out_for_delivery", "Out for Delivery"],
  ["delivered", "Delivered"],
  ["delayed", "Delayed"],
] as const;

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
  return statusLabels[status as ShipmentStatus] ?? status.replaceAll("_", " ");
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
  const isStaff = useAuthStore((state) => state.isStaff);
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
      <section className="border-b border-border bg-card/50">
        <div className="mx-auto flex max-w-7xl flex-wrap items-center justify-between gap-6 px-4 py-10 sm:px-6">
          <div>
            <span className="font-mono text-sm text-primary">/ orders</span>
            <h1 className="mt-2 text-3xl font-bold tracking-tight">Shipment manifest</h1>
            <p className="mt-2 text-sm text-muted-foreground">
              {total} total shipments tracked across all carriers.
            </p>
          </div>
          {isStaff && (
            <Link
              className="inline-flex h-10 items-center gap-2 rounded-md bg-primary px-4 text-sm font-medium text-primary-foreground hover:bg-primary/90 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
              href="/orders/create"
            >
              <Plus size={16} /> New Order
            </Link>
          )}
        </div>
      </section>
      <section className="mx-auto max-w-7xl space-y-5 px-4 py-8 sm:px-6">
        <div className="flex flex-col gap-4">
          <label className="relative block max-w-md">
            <Search
              className="absolute left-3 top-1/2 size-4 -translate-y-1/2 text-muted-foreground"
              aria-hidden="true"
            />
            <Input
              className="pl-9"
              type="search"
              aria-label="Search shipments"
              placeholder="Search by ID, tracking, customer, destination"
              value={search}
              onChange={(event) => changeSearch(event.target.value)}
            />
          </label>
          <div className="flex flex-wrap gap-2" aria-label="Filter by status">
            {filters.map(([key, label]) => (
              <Button
                key={key}
                type="button"
                variant={status === key ? "default" : "outline"}
                size="sm"
                aria-pressed={status === key}
                onClick={() => changeStatus(key)}
              >
                {label}
              </Button>
            ))}
          </div>
        </div>

        {error && !deleteTarget && (
          <p
            className="rounded-md border border-destructive/40 bg-destructive/10 p-3 text-sm text-destructive"
            role="alert"
          >
            {error}
          </p>
        )}
        {message && (
          <p className="rounded-md border border-border bg-muted p-3 text-sm" role="status">
            {message}
          </p>
        )}
        {loading && (
          <p className="rounded-md border border-border bg-muted p-3 text-sm" role="status">
            Loading shipments…
          </p>
        )}
        <div
          className="overflow-hidden rounded-xl border border-border bg-card shadow-sm"
          aria-busy={loading}
        >
          <div className="overflow-x-auto">
            <table className="w-full text-left text-sm">
              <thead className="border-b border-border bg-muted/50 text-xs uppercase text-muted-foreground">
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
                  <tr className="border-b border-border last:border-0" key={shipment.id}>
                    <td className="whitespace-nowrap px-4 py-3" data-label="Order ID">
                      <Link
                        className="font-medium text-primary hover:underline"
                        href={`/orders/${encodeURIComponent(shipment.id)}`}
                      >
                        {shipment.id}
                      </Link>
                    </td>
                    <td
                      className="whitespace-nowrap px-4 py-3 font-mono text-xs"
                      data-label="Tracking"
                    >
                      {shipment.trackingNumber}
                    </td>
                    <td className="px-4 py-3" data-label="Customer">
                      {shipment.customer.name}
                    </td>
                    <td className="px-4 py-3" data-label="Route">
                      <span className="flex min-w-64 items-center gap-2 text-xs text-muted-foreground">
                        {shipment.origin}
                        <ArrowRight size={13} aria-hidden="true" />
                        {shipment.destination}
                      </span>
                    </td>
                    <td data-label="Status">
                      <span className="inline-flex items-center gap-1.5 rounded-full bg-primary/10 px-2.5 py-1 text-xs font-medium capitalize text-primary">
                        <i aria-hidden="true" />
                        {statusLabel(shipment.status)}
                      </span>
                    </td>
                    <td className="whitespace-nowrap px-4 py-3 font-mono text-xs" data-label="ETA">
                      {formatDate(shipment.estimatedDelivery)}
                    </td>
                    <td data-label="Actions">
                      <div className="flex items-center gap-1">
                        <Link
                          className="inline-flex size-9 items-center justify-center rounded-md text-muted-foreground hover:bg-accent hover:text-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
                          href={`/orders/${encodeURIComponent(shipment.id)}`}
                          aria-label={`View shipment ${shipment.id}`}
                          title="View details"
                        >
                          <Eye size={16} />
                        </Link>
                        {isStaff && (
                          <Link
                            className="inline-flex size-9 items-center justify-center rounded-md text-muted-foreground hover:bg-accent hover:text-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
                            href={`/orders/${encodeURIComponent(shipment.id)}/edit`}
                            aria-label={`Edit shipment ${shipment.id}`}
                            title="Edit shipment"
                          >
                            <Pencil size={16} />
                          </Link>
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
              <p className="p-8 text-center text-sm text-muted-foreground">
                No shipments match your filters.
              </p>
            )}
          </div>
          {totalPages > 1 && (
            <nav
              className="flex flex-wrap items-center justify-between gap-4 border-t border-border px-4 py-3 text-sm text-muted-foreground"
              aria-label="Shipment pages"
            >
              <span>
                Showing {Math.min((page - 1) * pageSize + 1, total)}–
                {Math.min(page * pageSize, total)} of {total}
              </span>
              <div className="flex items-center gap-2">
                <Button
                  type="button"
                  variant="outline"
                  size="icon"
                  aria-label="Previous page"
                  disabled={page <= 1}
                  onClick={() => setPage((current) => current - 1)}
                >
                  <ChevronLeft size={16} />
                </Button>
                <span aria-current="page" className="min-w-14 text-center">
                  {page} / {totalPages}
                </span>
                <Button
                  type="button"
                  variant="outline"
                  size="icon"
                  aria-label="Next page"
                  disabled={page >= totalPages}
                  onClick={() => setPage((current) => current + 1)}
                >
                  <ChevronRight size={16} />
                </Button>
              </div>
            </nav>
          )}
        </div>
        <p className="text-xs text-muted-foreground">
          Status: {status === "all" ? "All" : statusLabel(status)}
        </p>
      </section>
      <dialog
        ref={dialog}
        className="fixed inset-0 m-auto w-[min(28rem,calc(100%-2rem))] rounded-xl border border-border bg-card p-6 text-card-foreground shadow-xl backdrop:bg-black/60"
        aria-labelledby="delete-heading"
        onCancel={(event) => {
          event.preventDefault();
          setDeleteTarget(null);
          setError(null);
        }}
      >
        <h2 id="delete-heading" className="text-lg font-semibold">
          Delete shipment?
        </h2>
        {deleteTarget && error && (
          <p
            className="mb-4 rounded-md border border-destructive/40 bg-destructive/10 p-3 text-sm text-destructive"
            role="alert"
          >
            {error}
          </p>
        )}
        <div className="mt-6 flex justify-end gap-3">
          <Button
            type="button"
            variant="outline"
            disabled={deleting}
            onClick={() => {
              setDeleteTarget(null);
              setError(null);
            }}
          >
            Cancel
          </Button>
          <Button type="button" variant="destructive" disabled={deleting} onClick={confirmDelete}>
            {deleting ? "Deleting…" : "Delete shipment"}
          </Button>
        </div>
      </dialog>
    </main>
  );
}
