import { mapHubs, type HubDTO } from "./mappers";
import type { Hub } from "./types";

export type HubInput = Omit<Hub, "id">;

export async function getHubs(): Promise<Hub[]> {
  const response = await fetch("/api/hubs", { cache: "no-store" });
  const result = (await response.json()) as { data: HubDTO[]; error?: string };
  if (!response.ok) throw new Error(result.error || `Hub API returned ${response.status}`);
  return mapHubs(result.data);
}

export async function saveHub(hub: HubInput, id?: string): Promise<void> {
  const response = await fetch(id ? `/api/hubs/${encodeURIComponent(id)}` : "/api/hubs", {
    method: id ? "PUT" : "POST",
    headers: { "content-type": "application/json" },
    body: JSON.stringify(hub),
  });
  const result = await response.json();
  if (!response.ok) throw new Error(result.error || `Failed to save hub (${response.status}).`);
}

export async function deleteHub(id: string): Promise<void> {
  const response = await fetch(`/api/hubs/${encodeURIComponent(id)}`, { method: "DELETE" });
  const result = await response.json();
  if (!response.ok) throw new Error(result.error || `Hub API returned ${response.status}`);
}
