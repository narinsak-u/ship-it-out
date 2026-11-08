"use client";

import { useState, type FormEvent } from "react";
import { ArrowRight, Search } from "lucide-react";
import { useRouter } from "next/navigation";
import styles from "./home.module.css";

interface ApiResult {
  success: boolean;
  data?: { shipment?: { id: string }; id?: string };
}

async function lookup(path: string): Promise<string | null> {
  const response = await fetch(path);
  if (!response.ok) return null;
  const result = (await response.json()) as ApiResult;
  return result.success ? (result.data?.shipment?.id ?? result.data?.id ?? null) : null;
}

export default function TrackingSearch() {
  const router = useRouter();
  const [query, setQuery] = useState("");
  const [loading, setLoading] = useState(false);

  async function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const value = query.trim();
    if (!value || loading) return;

    setLoading(true);
    try {
      const encoded = encodeURIComponent(value);
      const tracking = await lookup(`/api/track/${encoded}`).catch(() => null);
      const orderId = tracking ?? (await lookup(`/api/shipments/${encoded}`).catch(() => null));
      if (orderId) {
        router.push(`/orders/${encodeURIComponent(orderId)}`);
        return;
      }
      router.push("/orders");
    } finally {
      setLoading(false);
    }
  }

  return (
    <form className={styles.search} onSubmit={submit} aria-label="Track a shipment">
      <label className={styles.searchField}>
        <span className={styles.srOnly}>Tracking number or order ID</span>
        <Search className={styles.searchIcon} size={18} aria-hidden="true" />
        <input
          value={query}
          onChange={(event) => setQuery(event.target.value)}
          placeholder="Enter tracking number or order ID"
          autoComplete="off"
        />
      </label>
      <button type="submit" disabled={loading || !query.trim()}>
        {loading ? (
          "Searching…"
        ) : (
          <>
            Track <ArrowRight size={16} aria-hidden="true" />
          </>
        )}
      </button>
    </form>
  );
}
