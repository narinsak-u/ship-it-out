import { describe, it, expect } from "vitest";
import {
  computeKpis,
  computeRegionPerformance,
  computeStatusDistribution,
  computeStatusPieData,
  computeCumulativeData,
} from "./analytics-utils";
import type { AnalyticsOverview, TimeSeriesData } from "@/lib/api/analytics";

const mockAnalytics: AnalyticsOverview = {
  total: 100,
  active: 45,
  delivered: 80,
  by_status: [
    { status: "delivered", count: 80 },
    { status: "in_transit", count: 15 },
    { status: "pending", count: 5 },
  ],
  by_region: [
    { name: "Bangkok", total: 60 },
    { name: "Chiang Mai", total: 25 },
    { name: "Phuket", total: 15 },
  ],
};

const mockTimeSeries: TimeSeriesData = {
  by_month: [
    { month: "2026-01", count: 30 },
    { month: "2026-02", count: 25 },
    { month: "2026-03", count: 45 },
  ],
  by_day_of_week: [
    { day: "Mon", count: 10 },
    { day: "Tue", count: 15 },
  ],
};

describe("computeKpis", () => {
  it("calculates total and on-time rate", () => {
    const result = computeKpis(mockAnalytics);
    expect(result.total).toBe(100);
    expect(result.onTime).toBe(80);
    expect(result.regions).toBe(3);
    expect(result.avgDeliveryTime).toBe("3.2 days");
  });

  it("falls back to defaults when data is null", () => {
    const result = computeKpis(null);
    expect(result.total).toBe(0);
    expect(result.onTime).toBe(99.9);
    expect(result.regions).toBe(0);
  });

  it("handles zero total gracefully", () => {
    const zeroData: AnalyticsOverview = {
      total: 0,
      active: 0,
      delivered: 0,
      by_status: [],
      by_region: [],
    };
    const result = computeKpis(zeroData);
    expect(result.onTime).toBe(99.9);
  });
});

describe("computeRegionPerformance", () => {
  it("sorts by total descending", () => {
    const result = computeRegionPerformance(mockAnalytics);
    expect(result[0].name).toBe("Bangkok");
    expect(result[0].pct).toBe(60);
    expect(result[1].name).toBe("Chiang Mai");
    expect(result[2].name).toBe("Phuket");
  });

  it("returns empty array when data is null", () => {
    expect(computeRegionPerformance(null)).toEqual([]);
  });
});

describe("computeStatusDistribution", () => {
  it("maps status counts with percentages and labels", () => {
    const result = computeStatusDistribution(mockAnalytics);
    expect(result).toHaveLength(3);
    expect(result[0]).toMatchObject({ status: "delivered", count: 80, pct: 80 });
    expect(result[1]).toMatchObject({ status: "in_transit", count: 15, pct: 15 });
  });
});

describe("computeStatusPieData", () => {
  it("adds fill colors to status distribution", () => {
    const result = computeStatusPieData(mockAnalytics);
    expect(result[0].fill).toBe("var(--color-success)");
    expect(result[1].fill).toBe("var(--color-info)");
    expect(result[2].fill).toBe("var(--color-muted-foreground)");
  });

  it("uses fallback color for unknown statuses", () => {
    const customData: AnalyticsOverview = {
      total: 1,
      active: 1,
      delivered: 0,
      by_status: [{ status: "unknown", count: 1 }],
      by_region: [],
    };
    const result = computeStatusPieData(customData);
    expect(result[0].fill).toBe("hsl(var(--muted-foreground))");
  });
});

describe("computeCumulativeData", () => {
  it("computes running total from monthly counts", () => {
    const result = computeCumulativeData(mockTimeSeries);
    expect(result).toEqual([
      { month: "2026-01", count: 30 },
      { month: "2026-02", count: 55 },
      { month: "2026-03", count: 100 },
    ]);
  });

  it("returns empty array when data is null", () => {
    expect(computeCumulativeData(null)).toEqual([]);
  });
});
