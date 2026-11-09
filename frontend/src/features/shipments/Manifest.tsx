"use client";

import Link from "next/link";
import { useEffect, useRef, useState } from "react";
import {
  ArrowRight,
  ChevronLeft,
  ChevronRight,
  Eye,
  Filter,
  Pencil,
  Plus,
  Search,
  Trash2,
} from "lucide-react";
import { useAuthStore } from "@/features/auth/store/auth-store";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { statusLabels, type Shipment, type ShipmentPage, type ShipmentStatus } from "./types";
import { formatShipmentDate } from "./mappers";
import { deleteShipment, getManifestPage } from "./browser-api";

const pageSize = 10;
const filters = [
  ["all", "All"],
  ["pending", "Pending"],
  ["in_transit", "In Transit"],
  ["out_for_delivery", "Out for Delivery"],
  ["delivered", "Delivered"],
  ["delayed", "Delayed"],
] as const;

function statusLabel(status: string) {
  return statusLabels[status as ShipmentStatus] ?? status.replaceAll("_", " ");
}


const statusStyles: Record<ShipmentStatus, string> = {
  pending: "border-amber-400/30 bg-amber-400/10 text-amber-400",
  picked_up: "border-sky-400/30 bg-sky-400/10 text-sky-400",
  departed: "border-cyan-400/30 bg-cyan-400/10 text-cyan-400",
  in_transit: "border-teal-400/30 bg-teal-400/10 text-teal-300",
  out_for_delivery: "border-slate-400/20 bg-slate-400/10 text-slate-300",
  delivered: "border-emerald-400/30 bg-emerald-400/10 text-emerald-400",
  delayed: "border-rose-400/30 bg-rose-400/10 text-rose-400",
};

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
    getManifestPage(page, debouncedSearch, status, controller.signal)
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
      await deleteShipment(deleteTarget.id);
      setDeleteTarget(null);
      setMessage("Shipment deleted.");
      const refreshed = await getManifestPage(page, debouncedSearch, status);
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
      <section className="mx-auto max-w-screen-2xl space-y-5 px-4 py-8 sm:px-6">
        <div className="flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between">
          <label className="relative block w-full max-w-[30rem]">
            <Search
              className="absolute left-4 top-1/2 size-4 -translate-y-1/2 text-muted-foreground"
              aria-hidden="true"
            />
            <Input
              className="h-14 border-border bg-background/40 pl-14 font-mono text-base placeholder:text-muted-foreground"
              type="search"
              aria-label="Search shipments"
              placeholder="Search by ID, tracking, customer, destination"
              value={search}
              onChange={(event) => changeSearch(event.target.value)}
            />
          </label>
          <div
            className="flex flex-wrap items-center gap-2"
            role="group"
            aria-label="Filter by status"
          >
            <Filter className="mr-1 size-5 shrink-0 text-muted-foreground" aria-hidden="true" />
            {filters.map(([key, label]) => (
              <Button
                key={key}
                type="button"
                variant="outline"
                size="sm"
                className={`h-9 rounded-full border-border bg-transparent px-4 font-mono text-xs uppercase tracking-wider text-muted-foreground hover:border-primary/50 hover:bg-primary/5 hover:text-foreground ${status === key ? "border-primary bg-primary/10 text-primary" : ""}`}
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
            <table className="w-full min-w-[1100px] table-fixed text-left text-sm">
              <thead className="border-b border-border bg-muted/50 font-mono text-xs uppercase tracking-widest text-muted-foreground">
                <tr>
                  <th className="w-[11%] px-2.5 py-4">Order ID</th>
                  <th className="w-[11%] px-2.5 py-4">Tracking</th>
                  <th className="w-[10%] px-2.5 py-4">Customer</th>
                  <th className="w-[35%] px-2.5 py-4">Route</th>
                  <th className="w-[15%] px-2.5 py-4">Status</th>
                  <th className="w-[9%] px-2.5 py-4">ETA</th>
                  <th className="w-[9%] px-2.5 py-4">Actions</th>
                </tr>
              </thead>
              <tbody>
                {pageData.map((shipment) => {
                  const statusStyle =
                    statusStyles[shipment.status as ShipmentStatus] ??
                    "border-primary/30 bg-primary/10 text-primary";
                  return (
                    <tr
                      className="border-b border-border last:border-0 hover:bg-muted/20"
                      key={shipment.id}
                    >
                      <td className="whitespace-nowrap px-2.5 py-3 font-mono" data-label="Order ID">
                        <Link
                          className="block truncate font-medium text-primary hover:underline"
                          href={`/orders/${encodeURIComponent(shipment.id)}`}
                        >
                          {shipment.id}
                        </Link>
                      </td>
                      <td
                        className="whitespace-nowrap px-2.5 py-3 font-mono text-sm text-muted-foreground"
                        data-label="Tracking"
                      >
                        {shipment.trackingNumber}
                      </td>
                      <td className="truncate px-2.5 py-3 font-medium" data-label="Customer">
                        {shipment.customer.name}
                      </td>
                      <td className="px-2.5 py-3" data-label="Route">
                        <span className="flex min-w-0 items-center gap-2 font-mono text-xs text-muted-foreground">
                          <span className="truncate">{shipment.origin}</span>
                          <ArrowRight className="size-3 shrink-0 text-primary" aria-hidden="true" />
                          <span className="truncate">{shipment.destination}</span>
                        </span>
                      </td>
                      <td className="px-2.5 py-3" data-label="Status">
                        <span
                          className={`inline-flex items-center gap-1.5 whitespace-nowrap rounded-md border px-2.5 py-1 font-mono text-xs font-semibold uppercase tracking-wide ${statusStyle}`}
                        >
                          <i className="size-2 rounded-full bg-current" aria-hidden="true" />
                          {statusLabel(shipment.status)}
                        </span>
                      </td>
                      <td
                        className="whitespace-nowrap px-2.5 py-3 font-mono text-sm text-muted-foreground"
                        data-label="ETA"
                      >
                        {formatShipmentDate(shipment.estimatedDelivery)}
                      </td>
                      <td className="px-2.5 py-3" data-label="Actions">
                        <div className="flex items-center gap-0.5">
                          <Link
                            className="inline-flex size-8 items-center justify-center rounded-md text-muted-foreground hover:bg-accent hover:text-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
                            href={`/orders/${encodeURIComponent(shipment.id)}`}
                            aria-label={`View shipment ${shipment.id}`}
                            title="View details"
                          >
                            <Eye size={18} />
                          </Link>
                          {isStaff && (
                            <Link
                              className="inline-flex size-8 items-center justify-center rounded-md text-muted-foreground hover:bg-accent hover:text-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
                              href={`/orders/${encodeURIComponent(shipment.id)}/edit`}
                              aria-label={`Edit shipment ${shipment.id}`}
                              title="Edit shipment"
                            >
                              <Pencil size={18} />
                            </Link>
                          )}
                          {isStaff && (
                            <button
                              className="inline-flex size-8 items-center justify-center rounded-md text-muted-foreground hover:bg-accent hover:text-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
                              type="button"
                              aria-label={`Delete shipment ${shipment.id}`}
                              title="Delete shipment"
                              onClick={() => {
                                setError(null);
                                setDeleteTarget(shipment);
                              }}
                            >
                              <Trash2 size={18} />
                            </button>
                          )}
                        </div>
                      </td>
                    </tr>
                  );
                })}
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
              className="flex flex-wrap items-center justify-between gap-4 border-t border-border px-4 py-3 font-mono text-sm text-muted-foreground"
              aria-label="Shipment pages"
            >
              <span>
                Showing {Math.min((page - 1) * pageSize + 1, total)}–
                {Math.min(page * pageSize, total)} of {total}
              </span>
              <div className="flex items-center gap-1">
                <Button
                  type="button"
                  variant="ghost"
                  size="icon"
                  className="size-9 text-muted-foreground"
                  aria-label="Previous page"
                  disabled={page <= 1}
                  onClick={() => setPage((current) => current - 1)}
                >
                  <ChevronLeft size={18} />
                </Button>
                {Array.from(
                  { length: Math.min(totalPages, 5) },
                  (_, index) => Math.max(1, Math.min(page - 2, totalPages - 4)) + index,
                ).map((pageNumber) => (
                  <Button
                    key={pageNumber}
                    type="button"
                    variant="ghost"
                    className={`size-9 font-mono ${page === pageNumber ? "border border-primary text-primary" : "text-muted-foreground"}`}
                    aria-label={`Page ${pageNumber}`}
                    aria-current={page === pageNumber ? "page" : undefined}
                    onClick={() => setPage(pageNumber)}
                  >
                    {pageNumber}
                  </Button>
                ))}
                <Button
                  type="button"
                  variant="ghost"
                  size="icon"
                  className="size-9 text-muted-foreground"
                  aria-label="Next page"
                  disabled={page >= totalPages}
                  onClick={() => setPage((current) => current + 1)}
                >
                  <ChevronRight size={18} />
                </Button>
              </div>
            </nav>
          )}
        </div>
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
