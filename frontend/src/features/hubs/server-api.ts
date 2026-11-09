import { getApiOrigin } from "@/lib/api/server";
import { mapHubs, type HubDTO } from "./mappers";
import type { Hub } from "./types";

export async function getHubs(): Promise<Hub[]> {
  const response = await fetch(`${getApiOrigin()}/api/hubs`, { cache: "no-store" });
  if (!response.ok) throw new Error(`Hub API returned ${response.status}`);
  const result = (await response.json()) as { data: HubDTO[] };
  return mapHubs(result.data);
}
