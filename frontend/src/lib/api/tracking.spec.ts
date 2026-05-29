import { describe, it, expect } from "vitest";
import { trackShipment } from "./tracking";

describe("tracking API", () => {
  it("trackShipment returns shipment data", async () => {
    const result = await trackShipment("TH202600001");
    expect(result.shipment.id).toBe("ORD-001");
  });
});
