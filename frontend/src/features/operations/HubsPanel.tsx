"use client";

import { useMemo, useState } from "react";
import type { Hub, HubStatus } from "@/lib/hubs";
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
      const response = await fetch("/api/hubs", { cache: "no-store" });
      if (!response.ok) throw new Error(`Hub API returned ${response.status}`);
      const result = await response.json();
      setHubs(result.data as Hub[]);
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
    </>
  );
}
