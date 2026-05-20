import { useEffect } from "react";
import { MapContainer, TileLayer, Marker, Polyline, Popup, useMap } from "react-leaflet";
import L from "leaflet";
import "leaflet/dist/leaflet.css";
import type { GeoPoint, ShipmentStatus } from "@/lib/orders";

const originIcon = L.divIcon({
  className: "",
  html: `<div style="width:14px;height:14px;border-radius:9999px;background:oklch(0.78 0.11 195);box-shadow:0 0 0 4px oklch(0.78 0.11 195 / 0.25),0 0 12px oklch(0.78 0.11 195 / 0.6)"></div>`,
  iconSize: [14, 14],
  iconAnchor: [7, 7],
});

const destIcon = L.divIcon({
  className: "",
  html: `<div style="width:14px;height:14px;border-radius:9999px;background:oklch(0.65 0.12 210);box-shadow:0 0 0 4px oklch(0.65 0.12 210 / 0.25)"></div>`,
  iconSize: [14, 14],
  iconAnchor: [7, 7],
});

const currentIcon = L.divIcon({
  className: "",
  html: `
    <div style="position:relative;width:28px;height:28px;">
      <div style="position:absolute;inset:0;border-radius:9999px;background:oklch(0.78 0.15 75 / 0.3);animation:pulse 2s ease-out infinite;"></div>
      <div style="position:absolute;top:7px;left:7px;width:14px;height:14px;border-radius:9999px;background:oklch(0.78 0.15 75);border:2px solid white;"></div>
    </div>
    <style>@keyframes pulse{0%{transform:scale(0.6);opacity:1}100%{transform:scale(1.8);opacity:0}}</style>
  `,
  iconSize: [28, 28],
  iconAnchor: [14, 14],
});

function FitBounds({ points }: { points: GeoPoint[] }) {
  const map = useMap();
  useEffect(() => {
    const bounds = L.latLngBounds(points.map((p) => [p.lat, p.lng] as [number, number]));
    map.fitBounds(bounds, { padding: [50, 50] });
  }, [map, points]);
  return null;
}

export interface ShipmentMapProps {
  origin: GeoPoint;
  destination: GeoPoint;
  current: GeoPoint;
  originLabel: string;
  destinationLabel: string;
  carrier: string;
  status: ShipmentStatus;
}

export function ShipmentMap({
  origin,
  destination,
  current,
  originLabel,
  destinationLabel,
  carrier,
}: ShipmentMapProps) {
  return (
    <MapContainer
      center={[current.lat, current.lng]}
      zoom={3}
      scrollWheelZoom={false}
      style={{ height: "100%", width: "100%", background: "oklch(0.18 0.04 245)" }}
    >
      <TileLayer
        attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> &copy; CARTO'
        url="https://{s}.basemaps.cartocdn.com/dark_all/{z}/{x}/{y}{r}.png"
      />
      <Polyline
        positions={[
          [origin.lat, origin.lng],
          [current.lat, current.lng],
        ]}
        pathOptions={{ color: "oklch(0.78 0.11 195)", weight: 3, opacity: 0.9 }}
      />
      <Polyline
        positions={[
          [current.lat, current.lng],
          [destination.lat, destination.lng],
        ]}
        pathOptions={{
          color: "oklch(0.78 0.11 195)",
          weight: 2,
          opacity: 0.6,
          dashArray: "6 8",
        }}
      />
      <Marker position={[origin.lat, origin.lng]} icon={originIcon}>
        <Popup>Origin · {originLabel}</Popup>
      </Marker>
      <Marker position={[destination.lat, destination.lng]} icon={destIcon}>
        <Popup>Destination · {destinationLabel}</Popup>
      </Marker>
      <Marker position={[current.lat, current.lng]} icon={currentIcon}>
        <Popup>
          <strong>{carrier}</strong>
          <br />
          Current location
        </Popup>
      </Marker>
      <FitBounds points={[origin, destination, current]} />
    </MapContainer>
  );
}
