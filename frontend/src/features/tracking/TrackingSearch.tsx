"use client";

import { useState, type FormEvent } from "react";
import { useRouter } from "next/navigation";
import { ArrowRight, Search } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";

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
    <form
      className="flex flex-col gap-3 sm:flex-row"
      onSubmit={submit}
      aria-label="Track a shipment"
    >
      <label className="relative flex-1">
        <span className="sr-only">Tracking number or order ID</span>
        <Search
          className="absolute left-3 top-1/2 size-[18px] -translate-y-1/2 text-muted-foreground"
          aria-hidden="true"
        />
        <Input
          className="h-12 pl-10"
          value={query}
          onChange={(event) => setQuery(event.target.value)}
          placeholder="Enter tracking number or order ID"
          autoComplete="off"
        />
      </label>
      <Button className="h-12 gap-2 px-6" type="submit" disabled={loading || !query.trim()}>
        {loading ? (
          "Searching…"
        ) : (
          <>
            Track <ArrowRight size={16} aria-hidden="true" />
          </>
        )}
      </Button>
    </form>
  );
}
