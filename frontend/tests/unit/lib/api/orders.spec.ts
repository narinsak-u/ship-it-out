import { describe, it, expect } from "vitest";
import {
  fetchOrdersPaginated,
  fetchOrder,
  createOrder,
  updateOrder,
  deleteOrder,
  fetchActiveDeliveries,
  updateShipmentStatus,
  fetchOrderEvents,
} from "@/lib/api/orders";
import type { OrderFormData } from "@/lib/api/orders";

describe("orders API", () => {
  it("fetchOrdersPaginated returns paginated response", async () => {
    const result = await fetchOrdersPaginated({ page: 1, limit: 10 });
    expect(result.data).toBeInstanceOf(Array);
    expect(result.pagination).toBeDefined();
  });
  it("fetchOrdersPaginated with search param", async () => {
    const result = await fetchOrdersPaginated({ search: "test", status: "in_transit" });
    expect(result.data).toBeDefined();
  });
  it("fetchOrder returns single order", async () => {
    const result = await fetchOrder("ORD-001");
    expect(result.id).toBe("ORD-001");
  });
  it("createOrder sends POST and returns order", async () => {
    const data: OrderFormData = {
      customer: {
        name: "Test",
        zipcode: "10100",
        subDistrict: "A",
        district: "B",
        province: "C",
        coords: { lat: 1, lng: 2 },
      },
      receiver: {
        name: "Test2",
        zipcode: "50000",
        subDistrict: "D",
        district: "E",
        province: "F",
        coords: { lat: 3, lng: 4 },
      },
      carrier: "Test Carrier",
      weight: 10,
      items: 2,
      estimatedDelivery: "2026-06-01T00:00:00Z",
    };
    const result = await createOrder(data);
    expect(result.carrier).toBe("Test Carrier");
  });
  it("updateOrder sends PUT and returns updated order", async () => {
    const result = await updateOrder("ORD-001", { weight: 15 });
    expect(result).toBeDefined();
  });
  it("deleteOrder sends DELETE", async () => {
    await expect(deleteOrder("ORD-001")).resolves.toBeUndefined();
  });
  it("fetchActiveDeliveries returns non-delivered shipments", async () => {
    const result = await fetchActiveDeliveries();
    expect(result).toBeInstanceOf(Array);
  });
  it("updateShipmentStatus sends PATCH", async () => {
    const result = await updateShipmentStatus("ORD-001", "delivered", "hub-1");
    expect(result).toBeDefined();
  });
  it("fetchOrderEvents returns event array", async () => {
    const result = await fetchOrderEvents("TH202600001");
    expect(result).toBeInstanceOf(Array);
  });
});
