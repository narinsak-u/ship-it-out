import { describe, it, expect } from "vitest";
import { orderKeys, deliveryKeys, hubKeys, analyticsKeys } from "./queryKeys";

describe("orderKeys", () => {
  it("all returns root key", () => {
    expect(orderKeys.all).toEqual(["orders"]);
  });
  it("detail returns scoped key", () => {
    expect(orderKeys.detail("ORD-001")).toEqual(["orders", "detail", "ORD-001"]);
  });
  it("list returns scoped key", () => {
    expect(orderKeys.list({ page: 1 })).toEqual(["orders", "list", { page: 1 }]);
  });
});

describe("deliveryKeys", () => {
  it("all returns root", () => {
    expect(deliveryKeys.all).toEqual(["deliveries"]);
  });
  it("active returns scoped", () => {
    expect(deliveryKeys.active()).toEqual(["deliveries", "active"]);
  });
});

describe("hubKeys", () => {
  it("all returns root", () => {
    expect(hubKeys.all).toEqual(["hubs"]);
  });
});

describe("analyticsKeys", () => {
  it("all returns root", () => {
    expect(analyticsKeys.all).toEqual(["analytics"]);
  });
  it("timeseries returns scoped", () => {
    expect(analyticsKeys.timeseries()).toEqual(["analytics", "timeseries"]);
  });
});
