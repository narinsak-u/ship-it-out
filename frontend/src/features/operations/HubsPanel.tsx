"use client";

import { useEffect, useMemo, useRef, useState } from "react";
import { useAuth } from "@/features/auth/AuthProvider";
import type { Hub, HubStatus } from "@/lib/hubs";
import { HubForm } from "./HubForm";
import styles from "./HubsPanel.module.css";

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
  const { isStaff } = useAuth();
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
      <p className={styles.message} role="status">
        Loading hubs…
      </p>
    );
  if (error)
    return (
      <div className={styles.message} role="alert">
        <p>Failed to load hubs.</p>
        <button onClick={retry}>Retry</button>
      </div>
    );

  return (
    <>
      <div className={styles.stats} aria-label="Hub status counts">
        <div>
          <span>Total hubs</span>
          <strong>{counts.total}</strong>
        </div>
        {statuses.map((status) => (
          <div key={status} className={styles[`stat${status[0].toUpperCase()}${status.slice(1)}`]}>
            <span>{status}</span>
            <strong>{counts[status]}</strong>
          </div>
        ))}
      </div>
      <div className={styles.controls}>
        <label htmlFor="hub-search">Search hubs</label>
        <input
          id="hub-search"
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
        <div className={styles.actions}>
          <button type="button" onClick={() => setEditing(null)}>Add Hub</button>
        </div>
      )}
      {actionError && <p role="alert" className={styles.message}>{actionError}</p>}
      {notice && <p role="status" className={styles.message}>{notice}</p>}
      {filtered.length ? (
        <>
          <div className={styles.tableWrap}>
            <table className={styles.table}>
              <thead>
                <tr>
                  <th>ID</th>
                  <th>Name</th>
                  <th>Carrier</th>
                  <th>Address</th>
                  <th>Capacity</th>
                  <th>Status</th>
                  {isStaff && <th>Actions</th>}
                </tr>
              </thead>
              <tbody>
                {pageHubs.map((hub) => {
                  const utilization = hub.capacity
                    ? Math.round((hub.currentUtilization / hub.capacity) * 100)
                    : 0;
                  return (
                    <tr key={hub.id}>
                      <td data-label="ID">{hub.id}</td>
                      <td data-label="Name">{hub.name}</td>
                      <td data-label="Carrier">{hub.carrierId}</td>
                      <td data-label="Address">{hub.address}</td>
                      <td data-label="Capacity">
                        <div className={styles.capacity}>
                          <div
                            className={styles.capacityTrack}
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
                          <span className={styles.utilization}>{utilization}%</span>
                        </div>
                      </td>
                      <td data-label="Status">
                        <span className={`${styles.status} ${styles[hub.status]}`}>
                          {hub.status}
                        </span>
                      </td>
                      {isStaff && (
                        <td data-label="Actions">
                          <button type="button" aria-label={`Edit ${hub.name}`} onClick={() => setEditing(hub)}>Edit</button>
                          <button type="button" aria-label={`Delete ${hub.name}`} onClick={() => setDeleteTarget(hub)}>Delete</button>
                        </td>
                      )}
                    </tr>
                  );
                })}
              </tbody>
            </table>
          </div>
          <div className={styles.pagination}>
            <span>
              {filtered.length} hubs · Page {page} of {pages}
            </span>
            <div>
              <button
                aria-label="Previous page"
                disabled={page <= 1}
                onClick={() => setPage(page - 1)}
              >
                Previous
              </button>
              <button
                aria-label="Next page"
                disabled={page >= pages}
                onClick={() => setPage(page + 1)}
              >
                Next
              </button>
            </div>
          </div>
        </>
      ) : (
        <p className={styles.empty}>No hubs match your filters.</p>
      )}
      <dialog ref={confirmRef} aria-labelledby="hub-delete-title" className={styles.dialog} onCancel={() => setDeleteTarget(null)}>
        <h2 id="hub-delete-title">Delete Hub</h2>
        <p>Are you sure you want to delete {deleteTarget?.name}? This action cannot be undone.</p>
        <div className={styles.formActions}>
          <button type="button" onClick={() => setDeleteTarget(null)}>Cancel</button>
          <button type="button" onClick={() => deleteTarget && void removeHub(deleteTarget)}>Delete Hub</button>
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
