"use client";

import { useEffect, useRef } from "react";
import type { Map } from "leaflet";
import type { Hub } from "@/features/hubs/types";
import type { Shipment } from "@/features/shipments/types";
import "leaflet/dist/leaflet.css";
export function DeliveryMap({ order, hub }: Readonly<{ order: Shipment; hub?: Hub }>) {
  const container = useRef<HTMLDivElement>(null);
  const cartoKey = process.env.NEXT_PUBLIC_CARTO_API_KEY;

  useEffect(() => {
    let map: Map | undefined;
    let disposed = false;
    void import("leaflet").then((L) => {
      if (disposed || !container.current) return;
      const origin: [number, number] = [order.customer.coords.lat, order.customer.coords.lng];
      const current: [number, number] = [order.currentCoords.lat, order.currentCoords.lng];
      const destination: [number, number] = [order.receiver.coords.lat, order.receiver.coords.lng];
      map = L.map(container.current, {
        center: current,
        zoom: 6,
        scrollWheelZoom: false,
        zoomControl: false,
      });
      if (cartoKey) {
        L.tileLayer(
          `https://{s}.basemaps.cartocdn.com/dark_all/{z}/{x}/{y}{r}.png?key=${encodeURIComponent(cartoKey)}`,
          {
            attribution:
              '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> &copy; CARTO',
          },
        ).addTo(map);
      }
      L.polyline([origin, current], {
        color: "oklch(0.78 0.11 195)",
        weight: 3,
        opacity: 0.9,
      }).addTo(map);
      L.polyline([current, destination], {
        color: "oklch(0.78 0.11 195)",
        weight: 2,
        opacity: 0.6,
        dashArray: "6 8",
      }).addTo(map);
      const originIcon = L.divIcon({
        className: "",
        html: '<div style="width:14px;height:14px;border-radius:9999px;background:oklch(0.78 0.11 195);box-shadow:0 0 0 4px oklch(0.78 0.11 195 / 0.25),0 0 12px oklch(0.78 0.11 195 / 0.6)"></div>',
        iconSize: [14, 14],
        iconAnchor: [7, 7],
      });
      const destinationIcon = L.divIcon({
        className: "",
        html: '<div style="width:14px;height:14px;border-radius:9999px;background:oklch(0.65 0.12 210);box-shadow:0 0 0 4px oklch(0.65 0.12 210 / 0.25)"></div>',
        iconSize: [14, 14],
        iconAnchor: [7, 7],
      });
      const currentIcon = L.divIcon({
        className: "",
        html: '<div style="position:relative;width:28px;height:28px"><div style="position:absolute;inset:0;border-radius:9999px;background:oklch(0.78 0.15 75 / 0.3);animation:delivery-pulse 2s ease-out infinite"></div><div style="position:absolute;top:7px;left:7px;width:14px;height:14px;border-radius:9999px;background:oklch(0.78 0.15 75);border:2px solid white"></div></div><style>@keyframes delivery-pulse{0%{transform:scale(.6);opacity:1}100%{transform:scale(1.8);opacity:0}}</style>',
        iconSize: [28, 28],
        iconAnchor: [14, 14],
      });
      L.marker(origin, { title: `Origin: ${order.origin}`, icon: originIcon }).addTo(map);
      L.marker(destination, {
        title: `Destination: ${order.destination}`,
        icon: destinationIcon,
      }).addTo(map);
      L.marker(current, {
        title: `Current location: ${current[0]}, ${current[1]}`,
        icon: currentIcon,
      }).addTo(map);
      map.fitBounds(L.latLngBounds([origin, current, destination]), { padding: [40, 40] });
    });
    return () => {
      disposed = true;
      map?.remove();
    };
  }, [
    order.id,
    order.origin,
    order.destination,
    order.carrier,
    order.customer.coords.lat,
    order.customer.coords.lng,
    order.currentCoords.lat,
    order.currentCoords.lng,
    order.receiver.coords.lat,
    order.receiver.coords.lng,
    hub?.name,
    cartoKey,
  ]);

  return (
    <div className="relative mt-6 overflow-hidden rounded-lg border border-border">
      <div
        ref={container}
        className="leaflet-map h-[300px]"
        role="img"
        aria-label={`Delivery route from ${order.origin} to ${order.destination}; current location ${hub?.name ?? "in transit"}`}
      />
      {!cartoKey && (
        <p className="map-key-warning" role="status">
          Map tiles unavailable. Configure NEXT_PUBLIC_CARTO_API_KEY.
        </p>
      )}
    </div>
  );
}
