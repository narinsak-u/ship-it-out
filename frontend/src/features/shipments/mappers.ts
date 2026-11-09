import type { Shipment, ShipmentDTO, ShipmentStatus } from "./types";

export function mapShipment(dto: ShipmentDTO): Shipment {
  return {
    ...dto,
    status: dto.status as ShipmentStatus,
    estimatedDeliveryRaw: dto.estimatedDelivery,
    events: [],
  };
}

export function mapShipments(dtos: ShipmentDTO[]): Shipment[] {
  return dtos.map(mapShipment);
}
