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
let originMarker: L.Marker | null = null;
let destMarker: L.Marker | null = null;
let currentMarker: L.Marker | null = null;
let lineA: L.Polyline | null = null;
let lineB: L.Polyline | null = null;

function point(p: GeoPoint): [number, number] {
  return [p.lat, p.lng];
}

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

  const ori = point(props.origin);
  const cur = point(props.current);
  const dst = point(props.destination);

  lineA = L.polyline([ori, cur], {
    color: "oklch(0.78 0.11 195)",
    weight: 3,
    opacity: 0.9,
  }).addTo(map);

  lineB = L.polyline([cur, dst], {
    color: "oklch(0.78 0.11 195)",
    weight: 2,
    opacity: 0.6,
    dashArray: "6 8",
  }).addTo(map);

  originMarker = L.marker(ori, { icon: originIcon })
    .bindPopup(`Origin · ${props.originLabel}`)
    .addTo(map);

  destMarker = L.marker(dst, { icon: destIcon })
    .bindPopup(`Destination · ${props.destinationLabel}`)
    .addTo(map);

  currentMarker = L.marker(cur, { icon: currentIcon })
    .bindPopup(`<strong>${props.carrier}</strong><br>Current location`)
    .addTo(map);

  map.fitBounds(L.latLngBounds([ori, cur, dst]), { padding: [50, 50] });
});

onUnmounted(() => {
  if (map) {
    map.remove();
    map = null;
  }
  originMarker = null;
  destMarker = null;
  currentMarker = null;
  lineA = null;
  lineB = null;
});

watch(
  () => props.current,
  (cur) => {
    if (!map) return;
    const curLatLng = L.latLng(cur.lat, cur.lng);
    currentMarker?.setLatLng(curLatLng);
    lineA?.setLatLngs([point(props.origin), curLatLng]);
    lineB?.setLatLngs([curLatLng, point(props.destination)]);
    map.setView(curLatLng);
  },
  { deep: true },
);
</script>

<template>
  <div ref="mapContainer" class="h-full w-full bg-[#2e3440]" />
</template>
