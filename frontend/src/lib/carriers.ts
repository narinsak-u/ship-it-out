import type { GeoPoint } from "@/lib/orders";

export type CarrierStatus = "active" | "inactive";
export type HubStatus = "active" | "maintenance" | "closed";

export interface Carrier {
  id: string;
  name: string;
  contactEmail: string;
  phone: string;
  status: CarrierStatus;
  fleetSize: number;
  totalHubs: number;
  createdAt: string;
}

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

export const carriers: Carrier[] = [
  {
    id: "CAR-001",
    name: "Pacific Freight",
    contactEmail: "ops@pacificfreight.com",
    phone: "+1-555-0101",
    status: "active",
    fleetSize: 42,
    totalHubs: 6,
    createdAt: "2024-03-15",
  },
  {
    id: "CAR-002",
    name: "Skyline Express",
    contactEmail: "dispatch@skylineexpress.io",
    phone: "+1-555-0102",
    status: "active",
    fleetSize: 28,
    totalHubs: 4,
    createdAt: "2024-06-01",
  },
  {
    id: "CAR-003",
    name: "Trans-Atlantic Cargo",
    contactEmail: "support@tacargo.com",
    phone: "+1-555-0103",
    status: "active",
    fleetSize: 35,
    totalHubs: 5,
    createdAt: "2024-01-20",
  },
  {
    id: "CAR-004",
    name: "Nordic Lines",
    contactEmail: "ops@nordiclines.no",
    phone: "+47-555-0104",
    status: "active",
    fleetSize: 15,
    totalHubs: 3,
    createdAt: "2024-08-10",
  },
  {
    id: "CAR-005",
    name: "Gulf Logistics",
    contactEmail: "info@gulflogistics.ae",
    phone: "+971-555-0105",
    status: "active",
    fleetSize: 20,
    totalHubs: 3,
    createdAt: "2025-01-05",
  },
  {
    id: "CAR-006",
    name: "Mediterranean Freight",
    contactEmail: "ops@medfreight.it",
    phone: "+39-555-0106",
    status: "inactive",
    fleetSize: 10,
    totalHubs: 2,
    createdAt: "2025-02-18",
  },
];

export const hubs: Hub[] = [
  {
    id: "HUB-001",
    name: "Rotterdam Hub",
    carrierId: "CAR-001",
    address: "Haven 12, Rotterdam, NL",
    coords: { lat: 51.9244, lng: 4.4777 },
    capacity: 5000,
    currentUtilization: 3400,
    status: "active",
  },
  {
    id: "HUB-002",
    name: "Berlin Hub",
    carrierId: "CAR-002",
    address: "Industriestr. 45, Berlin, DE",
    coords: { lat: 52.52, lng: 13.405 },
    capacity: 3000,
    currentUtilization: 2100,
    status: "active",
  },
  {
    id: "HUB-003",
    name: "Lisbon Hub",
    carrierId: "CAR-003",
    address: "Av. da Marina 88, Lisbon, PT",
    coords: { lat: 38.7223, lng: -9.1393 },
    capacity: 4000,
    currentUtilization: 2800,
    status: "active",
  },
  {
    id: "HUB-004",
    name: "Oslo Warehouse",
    carrierId: "CAR-004",
    address: "Havnegata 9, Oslo, NO",
    coords: { lat: 59.9139, lng: 10.7522 },
    capacity: 2000,
    currentUtilization: 820,
    status: "active",
  },
  {
    id: "HUB-005",
    name: "Mumbai Warehouse",
    carrierId: "CAR-005",
    address: "Port Rd, Mumbai, IN",
    coords: { lat: 19.076, lng: 72.8777 },
    capacity: 2500,
    currentUtilization: 210,
    status: "active",
  },
  {
    id: "HUB-006",
    name: "Barcelona Hub",
    carrierId: "CAR-006",
    address: "Moll d'Espanya 3, Barcelona, ES",
    coords: { lat: 41.3851, lng: 2.1734 },
    capacity: 1800,
    currentUtilization: 990,
    status: "maintenance",
  },
];

export const hubStatusLabels: Record<HubStatus, string> = {
  active: "Active",
  maintenance: "Maintenance",
  closed: "Closed",
};

export const carrierStatusLabels: Record<CarrierStatus, string> = {
  active: "Active",
  inactive: "Inactive",
};

export function getCarrier(id: string) {
  return carriers.find((c) => c.id === id);
}

export function getCarrierByName(name: string) {
  return carriers.find((c) => c.name === name);
}

export function getHubsByCarrier(carrierId: string) {
  return hubs.filter((h) => h.carrierId === carrierId);
}
