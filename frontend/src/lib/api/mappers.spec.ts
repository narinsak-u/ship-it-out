import { describe, it, expect } from "vitest";
import {
  mapShipmentToOrder,
  mapEventToTrackingEvent,
  mapBackendHubToHub,
  formatDate,
  formatTimestamp,
} from "./mappers";
import type { BackendShipment, BackendShipmentEvent, BackendHub } from "./mappers";

describe("mapShipmentToOrder", () => {
  const backendShipment: BackendShipment = {
    id: "ORD-001",
    trackingNumber: "TH202600001",
    customer: {
      name: "John",
      zipcode: "10100",
      subDistrict: "A",
      district: "B",
      province: "C",
      coords: { lat: 1, lng: 2 },
    },
    receiver: {
      name: "Jane",
      zipcode: "50000",
      subDistrict: "D",
      district: "E",
      province: "F",
      coords: { lat: 3, lng: 4 },
    },
    currentCoords: { lat: 1, lng: 2 },
    origin: "Bangkok",
    destination: "Chiang Mai",
    status: "in_transit",
    carrier: "Test Carrier",
    weight: 10,
    items: 2,
    estimatedDelivery: "2026-06-01T00:00:00Z",
    createdAt: "2026-05-28T10:30:00Z",
    progress: 50,
  };
  it("maps all fields correctly", () => {
    const order = mapShipmentToOrder(backendShipment);
    expect(order.id).toBe("ORD-001");
    expect(order.trackingNumber).toBe("TH202600001");
    expect(order.status).toBe("in_transit");
    expect(order.weight).toBe(10);
    expect(order.items).toBe(2);
  });
  it("formats dates", () => {
    const order = mapShipmentToOrder(backendShipment);
    expect(order.estimatedDelivery).toContain("Jun");
    expect(order.createdAt).toContain("May");
  });
  it("includes hubId when present", () => {
    const order = mapShipmentToOrder({ ...backendShipment, hubId: "hub-1" });
    expect(order.hubId).toBe("hub-1");
  });
});

describe("formatDate", () => {
  it("formats ISO string to readable date", () => {
    expect(formatDate("2026-06-01T00:00:00Z")).toBe("Jun 1, 2026");
  });
});

describe("formatTimestamp", () => {
  it("formats ISO string with time", () => {
    const result = formatTimestamp("2026-06-01T14:30:00Z");
    expect(result).toContain("Jun 1");
    expect(result).toMatch(/^\w{3} \d{1,2}, \d{2}:\d{2}$/);
  });
});

describe("mapEventToTrackingEvent", () => {
  it("maps with description", () => {
    const be: BackendShipmentEvent = {
      id: 1,
      shipmentId: 1,
      status: "In Transit",
      location: { name: "Bangkok", lat: 13.75, lng: 100.5 },
      description: "Moving",
      timestamp: "2026-06-01T10:00:00Z",
    };
    const ev = mapEventToTrackingEvent(be);
    expect(ev.status).toBe("In Transit");
    expect(ev.description).toBe("Moving");
  });
  it("handles missing description", () => {
    const be: BackendShipmentEvent = {
      id: 2,
      shipmentId: 1,
      status: "Delivered",
      location: { name: "Chiang Mai", lat: 18.78, lng: 98.98 },
      timestamp: "2026-06-02T12:00:00Z",
    };
    const ev = mapEventToTrackingEvent(be);
    expect(ev.description).toBe("");
  });
});

describe("mapBackendHubToHub", () => {
  it("transforms hub correctly", () => {
    const bh: BackendHub = {
      id: "h1",
      name: "Bangkok Hub",
      carrierId: "c1",
      address: "123 Rd",
      coords: { lat: 13.75, lng: 100.5 },
      capacity: 1000,
      currentUtilization: 500,
      status: "active",
      createdAt: "2026-01-01T00:00:00Z",
    };
    const h = mapBackendHubToHub(bh);
    expect(h.name).toBe("Bangkok Hub");
    expect(h.status).toBe("active");
    expect(h.capacity).toBe(1000);
  });
});
