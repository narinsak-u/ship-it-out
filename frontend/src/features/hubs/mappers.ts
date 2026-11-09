import type { GeoPoint } from "@/lib/geo";
import type { Hub, HubStatus } from "./types";

export interface HubDTO {
  id: string;
  name: string;
  carrierId: string;
  address: string;
  coords: GeoPoint;
  capacity: number;
  currentUtilization: number;
  status: string;
  createdAt: string;
}

export function mapHub(dto: HubDTO): Hub {
  return {
    id: dto.id,
    name: dto.name,
    carrierId: dto.carrierId,
    address: dto.address,
    coords: dto.coords,
    capacity: dto.capacity,
    currentUtilization: dto.currentUtilization,
    status: dto.status as HubStatus,
  };
}

export function mapHubs(dtos: HubDTO[]): Hub[] {
  return dtos.map(mapHub);
}
