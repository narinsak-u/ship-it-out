<script setup lang="ts">
import { onMounted, onUnmounted, ref, watch } from "vue";
import L from "leaflet";
import "leaflet/dist/leaflet.css";
import type { GeoPoint, ShipmentStatus } from "@/lib/orders";

interface Props {
  origin: GeoPoint;
  destination: GeoPoint;
  current: GeoPoint;
  originLabel: string;
  destinationLabel: string;
  carrier: string;
  status: ShipmentStatus;
}

const props = defineProps<Props>();
const mapContainer = ref<HTMLElement | null>(null);
let map: L.Map | null = null;

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

onMounted(() => {
  if (!mapContainer.value) return;

  map = L.map(mapContainer.value, {
    center: [props.current.lat, props.current.lng],
    zoom: 3,
    scrollWheelZoom: false,
    zoomControl: false,
  });

  L.tileLayer("https://{s}.basemaps.cartocdn.com/dark_all/{z}/{x}/{y}{r}.png", {
    attribution:
      '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> &copy; CARTO',
  }).addTo(map);

  const points = [
    [props.origin.lat, props.origin.lng] as [number, number],
    [props.current.lat, props.current.lng] as [number, number],
    [props.destination.lat, props.destination.lng] as [number, number],
  ];

  L.polyline([points[0], points[1]], {
    color: "oklch(0.78 0.11 195)",
    weight: 3,
    opacity: 0.9,
  }).addTo(map);

  L.polyline([points[1], points[2]], {
    color: "oklch(0.78 0.11 195)",
    weight: 2,
    opacity: 0.6,
    dashArray: "6 8",
  }).addTo(map);

  L.marker([props.origin.lat, props.origin.lng], { icon: originIcon })
    .bindPopup(`Origin · ${props.originLabel}`)
    .addTo(map);

  L.marker([props.destination.lat, props.destination.lng], { icon: destIcon })
    .bindPopup(`Destination · ${props.destinationLabel}`)
    .addTo(map);

  L.marker([props.current.lat, props.current.lng], { icon: currentIcon })
    .bindPopup(`<strong>${props.carrier}</strong><br>Current location`)
    .addTo(map);

  const bounds = L.latLngBounds(points);
  map.fitBounds(bounds, { padding: [50, 50] });
});

onUnmounted(() => {
  if (map) {
    map.remove();
    map = null;
  }
});

watch(
  () => props.current,
  (newVal) => {
    if (map) {
      map.setView([newVal.lat, newVal.lng]);
    }
  },
  { deep: true },
);
</script>

<template>
  <div ref="mapContainer" class="h-full w-full bg-[#2e3440]" />
</template>
