import type { Shipment, ShipmentDTO, ShipmentStatus } from "./types";

const months = ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"];

function formatDate(iso: string): string {
  const date = new Date(iso);
  return `${months[date.getMonth()]} ${date.getDate()}, ${date.getFullYear()}`;
}

export function mapShipment(dto: ShipmentDTO): Shipment {
  return {
    ...dto,
    status: dto.status as ShipmentStatus,
    estimatedDelivery: formatDate(dto.estimatedDelivery),
    estimatedDeliveryRaw: dto.estimatedDelivery,
    createdAt: formatDate(dto.createdAt),
  };
}
export function mapShipments(dtos: ShipmentDTO[]): Shipment[] {
  return dtos.map(mapShipment);
}
