import type { Shipment, ShipmentDTO, ShipmentStatus } from "./types";

export function formatShipmentDate(iso: string): string {
  return new Date(iso).toLocaleDateString("en-US", {
    timeZone: "UTC",
    month: "short",
    day: "numeric",
    year: "numeric",
  });
}

export function mapShipment(dto: ShipmentDTO): Shipment {
  return {
    ...dto,
    status: dto.status as ShipmentStatus,
    estimatedDeliveryRaw: dto.estimatedDelivery,
  };
}

export function mapShipments(dtos: ShipmentDTO[]): Shipment[] {
  return dtos.map(mapShipment);
}
