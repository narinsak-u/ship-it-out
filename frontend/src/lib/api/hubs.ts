import type { Hub } from "@/lib/hubs";
import { api } from "@/lib/api/client";
import { mapBackendHubToHub, type BackendHub } from "@/lib/api/mappers";

export async function fetchHubs(): Promise<Hub[]> {
  const result = await api.get<BackendHub[]>("/hubs");
  if (result.error) throw new Error(result.error);
  return result.data!.map(mapBackendHubToHub);
}

export async function createHub(data: Omit<Hub, "id">): Promise<Hub> {
  const result = await api.post<BackendHub>("/hubs", data);
  if (result.error) throw new Error(result.error);
  return mapBackendHubToHub(result.data!);
}

export async function updateHub(id: string, data: Partial<Hub>): Promise<Hub> {
  const result = await api.put<BackendHub>(`/hubs/${id}`, data);
  if (result.error) throw new Error(result.error);
  return mapBackendHubToHub(result.data!);
}

export async function deleteHub(id: string): Promise<void> {
  const result = await api.del(`/hubs/${id}`);
  if (result.error) throw new Error(result.error);
}
