"use client";

import { useEffect, useMemo, useRef, useState } from "react";
import { useAuthStore } from "@/features/auth/auth-store";
import type { Hub, HubStatus } from "@/lib/hubs";
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
      hubs.reduce<Record<string, number>>(
        (result, hub) => {
          result.total++;
          result[hub.status]++;
          return result;
        },
        { total: 0, active: 0, maintenance: 0, closed: 0, full: 0 },
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
      <div className="grid grid-cols-2 gap-3 sm:grid-cols-5" aria-label="Hub status counts">
        <div className="grid gap-2.5 rounded-lg border border-border bg-card p-4">
          <span>Total hubs</span>
          <strong className="font-mono text-2xl">{counts.total}</strong>
        </div>
        {statuses.map((status) => (
          <div key={status} className="grid gap-2.5 rounded-lg border border-border bg-card p-4">
            <span className="capitalize text-muted-foreground">{status}</span>
            <strong
              className={`font-mono text-2xl ${status === "active" ? "text-emerald-400" : status === "closed" ? "text-destructive" : "text-amber-300"}`}
            >
              {counts[status]}
            </strong>
          </div>
        ))}
      </div>
      <div className="my-5 grid max-w-sm gap-2">
        <label htmlFor="hub-search">Search hubs</label>
        <input
          className="h-11 rounded-lg border border-border bg-card px-3 focus-visible:outline-2 focus-visible:outline-primary"
          type="search"
          placeholder="Search hubs..."
          value={query}
          onChange={(event) => {
            setQuery(event.target.value);
            setPage(1);
          }}
        />
      </div>
      {isStaff && (
        <div className="my-4 flex justify-end">
          <Button type="button" onClick={() => setEditing(null)}>
            Add Hub
          </Button>
        </div>
      )}
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
            <table className="w-full border-collapse text-left">
              <thead>
                <tr>
                  <th className="bg-secondary p-3 font-mono text-[10px] uppercase tracking-wider text-muted-foreground">
                    ID
                  </th>
                  <th className="bg-secondary p-3 font-mono text-[10px] uppercase tracking-wider text-muted-foreground">
                    Name
                  </th>
                  <th className="bg-secondary p-3 font-mono text-[10px] uppercase tracking-wider text-muted-foreground">
                    Carrier
                  </th>
                  <th className="bg-secondary p-3 font-mono text-[10px] uppercase tracking-wider text-muted-foreground">
                    Address
                  </th>
                  <th className="bg-secondary p-3 font-mono text-[10px] uppercase tracking-wider text-muted-foreground">
                    Capacity
                  </th>
                  <th className="bg-secondary p-3 font-mono text-[10px] uppercase tracking-wider text-muted-foreground">
                    Status
                  </th>
                  {isStaff && <th>Actions</th>}
                </tr>
              </thead>
              <tbody>
                {pageHubs.map((hub) => {
                  const utilization = hub.capacity
                    ? Math.round((hub.currentUtilization / hub.capacity) * 100)
                    : 0;
                  return (
                    <tr key={hub.id} className="border-b border-border">
                      <td data-label="ID">{hub.id}</td>
                      <td data-label="Name">{hub.name}</td>
                      <td data-label="Carrier">{hub.carrierId}</td>
                      <td data-label="Address">{hub.address}</td>
                      <td data-label="Capacity">
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
                          <span className="min-w-8 text-right text-xs text-muted-foreground">
                            {utilization}%
                          </span>
                        </div>
                      </td>
                      <td data-label="Status">
                        <span
                          className={`inline-flex rounded-full border px-2 py-1 capitalize ${hub.status === "active" ? "border-emerald-400/30 bg-emerald-400/10 text-emerald-300" : hub.status === "closed" ? "border-destructive/30 bg-destructive/10 text-destructive" : "border-amber-300/30 bg-amber-300/10 text-amber-200"}`}
                        >
                          {hub.status}
                        </span>
                      </td>
                      {isStaff && (
                        <td data-label="Actions">
                          <Button
                            type="button"
                            aria-label={`Edit ${hub.name}`}
                            onClick={() => setEditing(hub)}
                          >
                            Edit
                          </Button>
                          <Button
                            type="button"
                            aria-label={`Delete ${hub.name}`}
                            onClick={() => setDeleteTarget(hub)}
                          >
                            Delete
                          </Button>
                        </td>
                      )}
                    </tr>
                  );
                })}
              </tbody>
            </table>
          </div>
          <div className="flex items-center justify-between gap-3 py-4 font-mono text-xs text-muted-foreground max-sm:flex-col max-sm:items-stretch">
            <span>
              {filtered.length} hubs · Page {page} of {pages}
            </span>
            <div>
              <Button
                aria-label="Previous page"
                disabled={page <= 1}
                onClick={() => setPage(page - 1)}
              >
                Previous
              </Button>
              <Button
                aria-label="Next page"
                disabled={page >= pages}
                onClick={() => setPage(page + 1)}
              >
                Next
              </Button>
            </div>
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
