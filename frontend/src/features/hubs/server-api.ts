import { getApiOrigin } from "@/lib/api/server";
import type { Hub } from "./types";

export async function getHubs(): Promise<Hub[]> {
  const response = await fetch(`${getApiOrigin()}/api/hubs`, { cache: "no-store" });
  if (!response.ok) throw new Error(`Hub API returned ${response.status}`);
  const result = await response.json();
  return result.data as Hub[];
}
