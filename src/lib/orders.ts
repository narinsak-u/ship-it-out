export type ShipmentStatus = "pending" | "in_transit" | "out_for_delivery" | "delivered" | "delayed";

export interface TrackingEvent {
  timestamp: string;
  location: string;
  status: string;
  description: string;
}

export interface GeoPoint {
  lat: number;
  lng: number;
}

export interface Order {
  id: string;
  trackingNumber: string;
  customer: string;
  destination: string;
  origin: string;
  originCoords: GeoPoint;
  destinationCoords: GeoPoint;
  currentCoords: GeoPoint;
  status: ShipmentStatus;
  carrier: string;
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
    customer: "Aria Nakamura",
    origin: "Rotterdam, NL",
    destination: "Brooklyn, NY",
    status: "in_transit",
    carrier: "Pacific Freight",
    weight: "12.4 kg",
    items: 3,
    estimatedDelivery: "May 24, 2026",
    createdAt: "May 18, 2026",
    progress: 62,
    originCoords: { lat: 51.9244, lng: 4.4777 },
    destinationCoords: { lat: 40.6782, lng: -73.9442 },
    currentCoords: { lat: 46.5, lng: -34.5 },
    events: [
      { timestamp: "May 22, 09:14", location: "Mid-Atlantic", status: "In transit", description: "Vessel underway, on schedule." },
      { timestamp: "May 20, 17:42", location: "Rotterdam Port", status: "Departed", description: "Container loaded onto MV Northstar." },
      { timestamp: "May 19, 11:08", location: "Rotterdam Hub", status: "Processed", description: "Customs cleared." },
      { timestamp: "May 18, 08:30", location: "Rotterdam Warehouse", status: "Picked up", description: "Package received." },
    ],
  },
  {
    id: "ORD-10246",
    trackingNumber: "TRK-3C81-77D2",
    customer: "Marcus Brenner",
    origin: "Shenzhen, CN",
    destination: "Berlin, DE",
    status: "out_for_delivery",
    carrier: "Skyline Express",
    weight: "2.1 kg",
    items: 1,
    estimatedDelivery: "May 20, 2026",
    createdAt: "May 14, 2026",
    progress: 92,
    originCoords: { lat: 22.5431, lng: 114.0579 },
    destinationCoords: { lat: 52.52, lng: 13.405 },
    currentCoords: { lat: 52.52, lng: 13.38 },
    events: [
      { timestamp: "May 20, 07:01", location: "Berlin Mitte", status: "Out for delivery", description: "Loaded on courier vehicle." },
      { timestamp: "May 19, 22:18", location: "Berlin Hub", status: "Arrived", description: "Sorted for last-mile." },
      { timestamp: "May 17, 04:55", location: "Frankfurt Airport", status: "Customs cleared", description: "Released for transit." },
      { timestamp: "May 14, 19:32", location: "Shenzhen", status: "Shipped", description: "Departed origin facility." },
    ],
  },
  {
    id: "ORD-10247",
    trackingNumber: "TRK-22E5-09FA",
    customer: "Lina Okafor",
    origin: "São Paulo, BR",
    destination: "Lisbon, PT",
    status: "delivered",
    carrier: "Trans-Atlantic Cargo",
    weight: "5.8 kg",
    items: 2,
    estimatedDelivery: "May 17, 2026",
    createdAt: "May 09, 2026",
    progress: 100,
    originCoords: { lat: -23.5505, lng: -46.6333 },
    destinationCoords: { lat: 38.7223, lng: -9.1393 },
    currentCoords: { lat: 38.7223, lng: -9.1393 },
    events: [
      { timestamp: "May 17, 14:21", location: "Lisbon", status: "Delivered", description: "Signed by L. Okafor." },
      { timestamp: "May 17, 09:00", location: "Lisbon Hub", status: "Out for delivery", description: "" },
      { timestamp: "May 15, 06:12", location: "Lisbon Port", status: "Arrived", description: "" },
      { timestamp: "May 09, 11:45", location: "São Paulo", status: "Shipped", description: "" },
    ],
  },
  {
    id: "ORD-10248",
    trackingNumber: "TRK-7B11-A3CC",
    customer: "Theo Lindqvist",
    origin: "Oslo, NO",
    destination: "Reykjavík, IS",
    status: "delayed",
    carrier: "Nordic Lines",
    weight: "18.9 kg",
    items: 5,
    estimatedDelivery: "May 26, 2026",
    createdAt: "May 16, 2026",
    progress: 41,
    originCoords: { lat: 59.9139, lng: 10.7522 },
    destinationCoords: { lat: 64.1466, lng: -21.9426 },
    currentCoords: { lat: 62.5, lng: -2.5 },
    events: [
      { timestamp: "May 21, 10:00", location: "North Sea", status: "Delayed", description: "Weather hold, 36h estimate." },
      { timestamp: "May 18, 22:14", location: "Oslo Port", status: "Departed", description: "" },
      { timestamp: "May 16, 09:02", location: "Oslo Warehouse", status: "Picked up", description: "" },
    ],
  },
  {
    id: "ORD-10249",
    trackingNumber: "TRK-FF02-1188",
    customer: "Priya Anand",
    origin: "Mumbai, IN",
    destination: "Dubai, AE",
    status: "pending",
    carrier: "Gulf Logistics",
    weight: "0.9 kg",
    items: 1,
    estimatedDelivery: "May 25, 2026",
    createdAt: "May 20, 2026",
    progress: 8,
    originCoords: { lat: 19.076, lng: 72.8777 },
    destinationCoords: { lat: 25.2048, lng: 55.2708 },
    currentCoords: { lat: 19.076, lng: 72.8777 },
    events: [
      { timestamp: "May 20, 16:40", location: "Mumbai Warehouse", status: "Label created", description: "Awaiting pickup." },
    ],
  },
  {
    id: "ORD-10250",
    trackingNumber: "TRK-5E73-220B",
    customer: "Hugo Martín",
    origin: "Barcelona, ES",
    destination: "Marseille, FR",
    status: "in_transit",
    carrier: "Mediterranean Freight",
    weight: "7.3 kg",
    items: 4,
    estimatedDelivery: "May 22, 2026",
    createdAt: "May 19, 2026",
    progress: 55,
    originCoords: { lat: 41.3851, lng: 2.1734 },
    destinationCoords: { lat: 43.2965, lng: 5.3698 },
    currentCoords: { lat: 42.6986, lng: 2.8954 },
    events: [
      { timestamp: "May 20, 13:22", location: "Perpignan", status: "In transit", description: "Crossing border." },
      { timestamp: "May 19, 18:00", location: "Barcelona Hub", status: "Departed", description: "" },
    ],
  },
];

export const statusLabels: Record<ShipmentStatus, string> = {
  pending: "Pending",
  in_transit: "In Transit",
  out_for_delivery: "Out for Delivery",
  delivered: "Delivered",
  delayed: "Delayed",
};

export function getOrder(id: string) {
  return orders.find((o) => o.id === id);
}
