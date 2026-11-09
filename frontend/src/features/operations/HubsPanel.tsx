"use client";

import { useEffect, useMemo, useRef, useState } from "react";
import { useAuthStore } from "@/features/auth/store/auth-store";
import type { Hub, HubStatus } from "@/lib/hubs";
import { ChevronLeft, ChevronRight, Pencil, Plus, Search, Trash2 } from "lucide-react";
import { HubForm } from "./HubForm";
import { Button } from "@/components/ui/button";

const pageSize = 10;
const statuses: HubStatus[] = ["active", "maintenance", "full", "closed"];

export function HubsPanel({
  hubs: initialHubs,
  initialError,
}: {
  hubs: Hub[];
  initialError: boolean;
}) {
  const [hubs, setHubs] = useState(initialHubs);
  const [error, setError] = useState(initialError);
  const [loading, setLoading] = useState(false);
  const [query, setQuery] = useState("");
  const [page, setPage] = useState(1);
  const [editing, setEditing] = useState<Hub | null | undefined>(undefined);
  const [actionError, setActionError] = useState("");
  const [deleteTarget, setDeleteTarget] = useState<Hub | null>(null);
  const confirmRef = useRef<HTMLDialogElement>(null);
  const isStaff = useAuthStore((state) => state.isStaff);
  const [notice, setNotice] = useState("");

  useEffect(() => {
    if (deleteTarget && !confirmRef.current?.open) confirmRef.current?.showModal();
    else if (!deleteTarget && confirmRef.current?.open) confirmRef.current.close();
  }, [deleteTarget]);
  async function refresh() {
    const response = await fetch("/api/hubs", { cache: "no-store" });
    const result = await response.json();
    if (!response.ok) throw new Error(result.error || `Hub API returned ${response.status}`);
    setHubs(result.data as Hub[]);
  }

  async function removeHub(hub: Hub) {
    setActionError("");
    setNotice("");
    try {
      const response = await fetch(`/api/hubs/${encodeURIComponent(hub.id)}`, { method: "DELETE" });
      const result = await response.json();
      if (!response.ok) throw new Error(result.error || `Hub API returned ${response.status}`);
      setHubs((current) => current.filter((item) => item.id !== hub.id));
      setNotice("Hub deleted.");
      setDeleteTarget(null);
    } catch (cause) {
      setActionError(cause instanceof Error ? cause.message : "Failed to delete hub.");
      setDeleteTarget(null);
    }
  }
  const counts = useMemo(
    () =>
      hubs.reduce<Record<HubStatus | "total", number>>(
        (result, hub) => {
          result.total++;
          result[hub.status]++;
          return result;
        },
        { total: 0, active: 0, maintenance: 0, full: 0, closed: 0 },
      ),
    [hubs],
  );
  const filtered = useMemo(() => {
    const search = query.trim().toLowerCase();
    return hubs.filter(
      (hub) =>
        !search ||
        `${hub.name} ${hub.address} ${hub.id} ${hub.carrierId}`.toLowerCase().includes(search),
    );
  }, [hubs, query]);
  const pages = Math.ceil(filtered.length / pageSize);
  const pageHubs = filtered.slice((page - 1) * pageSize, page * pageSize);

  async function retry() {
    setLoading(true);
    try {
      await refresh();
      setError(false);
    } catch {
      setError(true);
    } finally {
      setLoading(false);
    }
  }

  if (loading)
    return (
      <p className="py-8 text-center text-muted-foreground" role="status">
        Loading hubs…
      </p>
    );
  if (error)
    return (
      <div className="grid justify-items-center gap-3 py-8 text-muted-foreground" role="alert">
        <p>Failed to load hubs.</p>
        <Button onClick={retry}>Retry</Button>
      </div>
    );

  return (
    <>
      <div className="mb-5 grid grid-cols-2 gap-3 sm:grid-cols-5" aria-label="Hub status counts">
        <div className="grid gap-2 rounded-lg border border-border bg-secondary/40 p-5">
          <span className="font-mono text-[11px] uppercase tracking-[0.14em] text-muted-foreground">
            Total hubs
          </span>
          <strong className="font-mono text-4xl font-semibold leading-none">{counts.total}</strong>
        </div>
        {statuses.map((status) => (
          <div
            key={status}
            className="grid gap-2 rounded-lg border border-border bg-secondary/40 p-5"
          >
            <span className="font-mono text-[11px] uppercase tracking-[0.14em] text-muted-foreground">
              {status}
            </span>
            <strong
              className={`font-mono text-4xl font-semibold leading-none ${status === "active" ? "text-emerald-400" : status === "maintenance" ? "text-amber-400" : "text-destructive"}`}
            >
              {counts[status]}
            </strong>
          </div>
        ))}
      </div>
      <div className="mb-5 flex items-center justify-between gap-4 max-sm:flex-col max-sm:items-stretch">
        <label className="relative block w-full max-w-sm">
          <span className="sr-only">Search hubs</span>
          <Search
            className="pointer-events-none absolute top-1/2 left-3 size-4 -translate-y-1/2 text-muted-foreground"
            aria-hidden="true"
          />
          <input
            className="h-11 w-full rounded-lg border border-border bg-background pl-10 pr-3 text-foreground placeholder:text-muted-foreground focus-visible:outline-2 focus-visible:outline-primary"
            id="hub-search"
            aria-label="Search hubs"
            type="search"
            placeholder="Search hubs..."
            value={query}
            onChange={(event) => {
              setQuery(event.target.value);
              setPage(1);
            }}
          />
        </label>
        {isStaff && (
          <Button type="button" className="max-sm:w-full" onClick={() => setEditing(null)}>
            <Plus className="size-4" aria-hidden="true" />
            Add Hub
          </Button>
        )}
      </div>
      {actionError && (
        <p role="alert" className="py-3 text-destructive">
          {actionError}
        </p>
      )}
      {notice && (
        <p role="status" className="py-3 text-emerald-400">
          {notice}
        </p>
      )}
      {filtered.length ? (
        <>
          <div className="overflow-x-auto rounded-lg border border-border">
            <table className="w-full min-w-[900px] border-collapse text-left text-sm">
              <thead>
                <tr>
                  <th className="bg-secondary px-3 py-3 font-mono text-[10px] uppercase tracking-wider text-muted-foreground">
                    ID
                  </th>
                  <th className="bg-secondary px-3 py-3 font-mono text-[10px] uppercase tracking-wider text-muted-foreground">
                    Name
                  </th>
                  <th className="bg-secondary px-3 py-3 font-mono text-[10px] uppercase tracking-wider text-muted-foreground">
                    Carrier
                  </th>
                  <th className="bg-secondary px-3 py-3 font-mono text-[10px] uppercase tracking-wider text-muted-foreground">
                    Address
                  </th>
                  <th className="bg-secondary px-3 py-3 font-mono text-[10px] uppercase tracking-wider text-muted-foreground">
                    Capacity
                  </th>
                  <th className="bg-secondary px-3 py-3 font-mono text-[10px] uppercase tracking-wider text-muted-foreground">
                    Status
                  </th>
                  {isStaff && (
                    <th className="bg-secondary px-3 py-3 font-mono text-[10px] uppercase tracking-wider text-muted-foreground">
                      Actions
                    </th>
                  )}
                </tr>
              </thead>
              <tbody>
                {pageHubs.map((hub) => {
                  const utilization = hub.capacity
                    ? Math.round((hub.currentUtilization / hub.capacity) * 100)
                    : 0;
                  return (
                    <tr
                      key={hub.id}
                      className="border-b border-border last:border-0 hover:bg-primary/5"
                    >
                      <td className="whitespace-nowrap px-3 py-3.5 font-mono text-muted-foreground">
                        {hub.id}
                      </td>
                      <td className="whitespace-nowrap px-3 py-3.5 font-mono font-semibold text-foreground">
                        {hub.name}
                      </td>
                      <td className="whitespace-nowrap px-3 py-3.5 font-mono text-muted-foreground">
                        {hub.carrierId}
                      </td>
                      <td className="whitespace-nowrap px-3 py-3.5 font-mono text-muted-foreground">
                        {hub.address}
                      </td>
                      <td className="px-3 py-3.5">
                        <div className="flex min-w-28 items-center gap-2">
                          <div
                            className="h-2 min-w-12 flex-1 overflow-hidden rounded-full bg-secondary [&>span]:block [&>span]:h-full [&>span]:rounded-full [&>span]:bg-primary"
                            role="progressbar"
                            aria-label={`${hub.name} capacity`}
                            aria-valuemin={0}
                            aria-valuemax={100}
                            aria-valuenow={utilization}
                          >
                            <span
                              style={{ width: `${Math.min(100, Math.max(0, utilization))}%` }}
                            />
                          </div>
                          <span className="min-w-8 text-right font-mono text-xs text-muted-foreground">
                            {utilization}%
                          </span>
                        </div>
                      </td>
                      <td className="px-3 py-3.5">
                        <span
                          className={`inline-flex items-center gap-2 rounded-md border px-2 py-1 font-mono text-xs font-semibold uppercase tracking-wide ${hub.status === "active" ? "border-emerald-400/30 bg-emerald-400/10 text-emerald-400" : hub.status === "closed" ? "border-destructive/30 bg-destructive/10 text-destructive" : "border-amber-300/30 bg-amber-300/10 text-amber-300"}`}
                        >
                          <span className="size-2 rounded-full bg-current" aria-hidden="true" />
                          {hub.status}
                        </span>
                      </td>
                      {isStaff && (
                        <td className="px-2 py-2">
                          <div className="flex items-center gap-1">
                            <Button
                              type="button"
                              variant="ghost"
                              size="icon-sm"
                              className="text-muted-foreground hover:text-primary"
                              aria-label={`Edit ${hub.name}`}
                              onClick={() => setEditing(hub)}
                            >
                              <Pencil className="size-4" aria-hidden="true" />
                            </Button>
                            <Button
                              type="button"
                              variant="ghost"
                              size="icon-sm"
                              className="text-muted-foreground hover:text-destructive"
                              aria-label={`Delete ${hub.name}`}
                              onClick={() => setDeleteTarget(hub)}
                            >
                              <Trash2 className="size-4" aria-hidden="true" />
                            </Button>
                          </div>
                        </td>
                      )}
                    </tr>
                  );
                })}
              </tbody>
            </table>
          </div>
          <div className="flex items-center justify-between gap-3 rounded-b-lg border border-t-0 border-border px-3 py-2 font-mono text-xs text-muted-foreground max-sm:flex-col max-sm:items-stretch">
            <span>
              Showing {(page - 1) * pageSize + 1}–{Math.min(page * pageSize, filtered.length)} of{" "}
              {filtered.length}
            </span>
            <nav className="flex items-center gap-1" aria-label="Hub pages">
              <Button
                variant="ghost"
                size="icon-sm"
                aria-label="Previous page"
                disabled={page <= 1}
                onClick={() => setPage(page - 1)}
              >
                <ChevronLeft className="size-4" aria-hidden="true" />
              </Button>
              {Array.from({ length: pages }, (_, index) => index + 1).map((pageNumber) => (
                <Button
                  key={pageNumber}
                  variant={pageNumber === page ? "outline" : "ghost"}
                  size="icon-sm"
                  aria-label={`Page ${pageNumber}`}
                  aria-current={pageNumber === page ? "page" : undefined}
                  onClick={() => setPage(pageNumber)}
                >
                  {pageNumber}
                </Button>
              ))}
              <Button
                variant="ghost"
                size="icon-sm"
                aria-label="Next page"
                disabled={page >= pages}
                onClick={() => setPage(page + 1)}
              >
                <ChevronRight className="size-4" aria-hidden="true" />
              </Button>
            </nav>
          </div>
        </>
      ) : (
        <p className="p-8 text-center text-muted-foreground">No hubs match your filters.</p>
      )}
      <dialog
        ref={confirmRef}
        aria-labelledby="hub-delete-title"
        className="fixed inset-0 z-50 m-auto max-h-[calc(100dvh-24px)] w-[min(100%-24px,480px)] overflow-auto rounded-xl border border-border bg-card p-6 text-foreground backdrop:bg-black/65"
        onCancel={() => setDeleteTarget(null)}
      >
        <h2 id="hub-delete-title">Delete Hub</h2>
        <p>Are you sure you want to delete {deleteTarget?.name}? This action cannot be undone.</p>
        <div className="flex justify-end gap-2">
          <button type="button" onClick={() => setDeleteTarget(null)}>
            Cancel
          </button>
          <button type="button" onClick={() => deleteTarget && void removeHub(deleteTarget)}>
            Delete Hub
          </button>
        </div>
      </dialog>
      {editing !== undefined && (
        <HubForm
          hub={editing}
          onClose={() => setEditing(undefined)}
          onSaved={async () => {
            setEditing(undefined);
            setNotice(editing ? "Hub updated." : "Hub created.");
            try {
              await refresh();
            } catch {
              setActionError("Hub saved, but the listing could not refresh.");
            }
          }}
        />
      )}
    </>
  );
}
