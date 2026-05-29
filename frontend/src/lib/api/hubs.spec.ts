import { describe, it, expect } from "vitest";
import { fetchHubs, createHub, updateHub, deleteHub } from "./hubs";

describe("hubs API", () => {
  it("fetchHubs returns array", async () => {
    const result = await fetchHubs();
    expect(result).toBeInstanceOf(Array);
  });
  it("createHub sends POST", async () => {
    const result = await createHub({
      name: "Test Hub",
      carrierId: "c1",
      address: "123 St",
      coords: { lat: 0, lng: 0 },
      capacity: 500,
      currentUtilization: 0,
      status: "active",
    });
    expect(result.name).toBe("Test Hub");
  });
  it("updateHub sends PUT", async () => {
    const result = await updateHub("hub-1", { name: "Updated Hub" });
    expect(result).toBeDefined();
  });
  it("deleteHub sends DELETE", async () => {
    await expect(deleteHub("hub-1")).resolves.toBeUndefined();
  });
});
