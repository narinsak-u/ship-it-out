import type { GeoPoint } from "@/lib/orders";

export type HubStatus = "active" | "maintenance" | "closed" | "full";

export interface Hub {
  id: string;
  name: string;
  carrierId: string;
  address: string;
  coords: GeoPoint;
  capacity: number;
  currentUtilization: number;
  status: HubStatus;
}

export const hubStatusLabels: Record<HubStatus, string> = {
  active: "Active",
  maintenance: "Maintenance",
  closed: "Closed",
  full: "Full",
};
