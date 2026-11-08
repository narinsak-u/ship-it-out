"use client";

import { useEffect, useRef } from "react";
import "leaflet/dist/leaflet.css";
import type { Map as LeafletMap } from "leaflet";

interface Props {
  origin: { lat: number; lng: number };
  destination: { lat: number; lng: number };
  current: { lat: number; lng: number };
  originLabel: string;
  destinationLabel: string;
  carrier: string;
}

export default function RouteMap(props: Props) {
  const element = useRef<HTMLDivElement>(null);
  const cartoKey = process.env.NEXT_PUBLIC_CARTO_API_KEY;
  useEffect(() => {
    let map: LeafletMap | undefined;
    let disposed = false;
    void import("leaflet").then((L) => {
      if (disposed || !element.current) return;
      map = L.map(element.current, {
        center: [props.current.lat, props.current.lng],
        zoom: 3,
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
      const points: [number, number][] = [
        [props.origin.lat, props.origin.lng],
        [props.current.lat, props.current.lng],
        [props.destination.lat, props.destination.lng],
      ];
      L.polyline(points.slice(0, 2), {
        color: "oklch(0.78 0.11 195)",
        weight: 3,
        opacity: 0.9,
      }).addTo(map);
      L.polyline(points.slice(1), {
        color: "oklch(0.78 0.11 195)",
        weight: 2,
        opacity: 0.6,
        dashArray: "6 8",
      }).addTo(map);
      const marker = (color: string) =>
        L.divIcon({
          className: "",
          html: `<span style="display:block;width:14px;height:14px;border-radius:50%;background:${color};box-shadow:0 0 0 4px ${color}40"></span>`,
          iconSize: [14, 14],
          iconAnchor: [7, 7],
        });
      const popup = (label: string) => {
        const node = document.createElement("span");
        node.textContent = label;
        return node;
      };
      L.marker(points[0], { icon: marker("oklch(0.78 0.11 195)") })
        .bindPopup(popup(`Origin · ${props.originLabel}`))
        .addTo(map);
      L.marker(points[2], { icon: marker("oklch(0.65 0.12 210)") })
        .bindPopup(popup(`Destination · ${props.destinationLabel}`))
        .addTo(map);
      L.marker(points[1], { icon: marker("oklch(0.78 0.15 75)") })
        .bindPopup(popup(`${props.carrier} · Current location`))
        .addTo(map);
      map.fitBounds(L.latLngBounds(points), { padding: [50, 50] });
    });
    return () => {
      disposed = true;
      map?.remove();
    };
  }, [props, cartoKey]);
  return (
    <div className="leaflet-map" ref={element} role="region" aria-label="Shipment route map">
      {!cartoKey && (
        <p className="map-key-warning" role="status">
          Map tiles unavailable. Configure NEXT_PUBLIC_CARTO_API_KEY.
        </p>
      )}
    </div>
  );
}
