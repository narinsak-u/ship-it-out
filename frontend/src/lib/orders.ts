export type ShipmentStatus =
  | "pending"
  | "picked_up"
  | "departed"
  | "in_transit"
  | "out_for_delivery"
  | "delivered"
  | "delayed";

export interface Location {
  name: string;
  lat: number;
  lng: number;
}

export interface TrackingEvent {
  timestamp: string;
  location: Location;
  status: string;
  description: string;
}

export interface GeoPoint {
  lat: number;
  lng: number;
}

export interface ContactInfo {
  name: string;
  zipcode: string;
  subDistrict: string;
  district: string;
  province: string;
  coords: GeoPoint;
}

export interface Order {
  id: string;
  trackingNumber: string;
  customer: ContactInfo;
  receiver: ContactInfo;
  origin: string;
  destination: string;
  currentCoords: GeoPoint;
  status: ShipmentStatus;
  carrier: string;
  hubId?: string;
  weight: string;
  items: number;
  estimatedDelivery: string;
  createdAt: string;
  progress: number;
  events: TrackingEvent[];
}

export const orders: Order[] = [
  {
    id: "ORD-10245",
    trackingNumber: "TRK-9F2A-44B1",
    customer: {
      name: "Aria Nakamura",
      zipcode: "3011",
      subDistrict: "Stadsdriehoek",
      district: "Centrum",
      province: "Zuid-Holland",
      coords: { lat: 51.9244, lng: 4.4777 },
    },
    receiver: {
      name: "James Mitchell",
      zipcode: "11201",
      subDistrict: "DUMBO",
      district: "Brooklyn",
      province: "New York",
      coords: { lat: 40.6782, lng: -73.9442 },
    },
    origin: "Stadsdriehoek, Centrum, Zuid-Holland",
    destination: "DUMBO, Brooklyn, New York",
    currentCoords: { lat: 46.5, lng: -34.5 },
    status: "in_transit",
    carrier: "Pacific Freight",
    weight: "12.4 kg",
    items: 3,
    estimatedDelivery: "May 24, 2026",
    createdAt: "May 18, 2026",
    progress: 62,
    events: [
      {
        timestamp: "May 21, 10:00",
        location: { name: "North Sea", lat: 62.5, lng: -2.5 },
        status: "Delayed",
        description: "Weather hold, 36h estimate.",
      },
      {
        timestamp: "May 18, 22:14",
        location: { name: "Oslo Port", lat: 59.9, lng: 10.75 },
        status: "Departed",
        description: "",
      },
      {
        timestamp: "May 16, 09:02",
        location: { name: "Oslo Warehouse", lat: 59.9139, lng: 10.7522 },
        status: "Picked up",
        description: "",
      },
    ],
  },
  {
    id: "ORD-10249",
    trackingNumber: "TRK-FF02-1188",
    customer: {
      name: "Priya Anand",
      zipcode: "400001",
      subDistrict: "Fort",
      district: "Mumbai City",
      province: "Maharashtra",
      coords: { lat: 19.076, lng: 72.8777 },
    },
    receiver: {
      name: "Ahmed Al-Rashid",
      zipcode: "00000",
      subDistrict: "Downtown",
      district: "Dubai",
      province: "Dubai",
      coords: { lat: 25.2048, lng: 55.2708 },
    },
    origin: "Fort, Mumbai City, Maharashtra",
    destination: "Downtown, Dubai, Dubai",
    currentCoords: { lat: 19.076, lng: 72.8777 },
    status: "pending",
    carrier: "Gulf Logistics",
    weight: "0.9 kg",
    items: 1,
    estimatedDelivery: "May 25, 2026",
    createdAt: "May 20, 2026",
    progress: 8,
    events: [
      {
        timestamp: "May 20, 16:40",
        location: { name: "Mumbai Warehouse", lat: 19.076, lng: 72.8777 },
        status: "Label created",
        description: "Awaiting pickup.",
      },
    ],
  },
  {
    id: "ORD-10250",
    trackingNumber: "TRK-5E73-220B",
    customer: {
      name: "Hugo Martín",
      zipcode: "08001",
      subDistrict: "El Raval",
      district: "Ciutat Vella",
      province: "Barcelona",
      coords: { lat: 41.3851, lng: 2.1734 },
    },
    receiver: {
      name: "Camille Dubois",
      zipcode: "13001",
      subDistrict: "Vieux-Port",
      district: "Marseille",
      province: "Provence-Alpes",
      coords: { lat: 43.2965, lng: 5.3698 },
    },
    origin: "El Raval, Ciutat Vella, Barcelona",
    destination: "Vieux-Port, Marseille, Provence-Alpes",
    currentCoords: { lat: 42.6986, lng: 2.8954 },
    status: "in_transit",
    carrier: "Mediterranean Freight",
    weight: "7.3 kg",
    items: 4,
    estimatedDelivery: "May 22, 2026",
    createdAt: "May 19, 2026",
    progress: 55,
    events: [
      {
        timestamp: "May 20, 13:22",
        location: { name: "Perpignan", lat: 42.69, lng: 2.89 },
        status: "In transit",
        description: "Crossing border.",
      },
      {
        timestamp: "May 19, 18:00",
        location: { name: "Barcelona Hub", lat: 41.3851, lng: 2.1734 },
        status: "Departed",
        description: "",
      },
    ],
  },
];

export const statusLabels: Record<ShipmentStatus, string> = {
  pending: "Pending",
  picked_up: "Picked Up",
  departed: "Departed",
  in_transit: "In Transit",
  out_for_delivery: "Out for Delivery",
  delivered: "Delivered",
  delayed: "Delayed",
};

export function getOrder(id: string) {
  return orders.find((o) => o.id === id);
}
