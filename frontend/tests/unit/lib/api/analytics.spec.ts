import { describe, it, expect } from "vitest";
import { fetchAnalytics, fetchTimeSeries } from "@/lib/api/analytics";

describe("analytics API", () => {
  it("fetchAnalytics returns overview", async () => {
    const result = await fetchAnalytics();
    expect(result.total).toBeGreaterThan(0);
    expect(result.active).toBeDefined();
  });
  it("fetchTimeSeries returns monthly data", async () => {
    const result = await fetchTimeSeries();
    expect(result.by_month).toBeInstanceOf(Array);
    expect(result.by_day_of_week).toBeInstanceOf(Array);
  });
});
