# Testing Gap Fill Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Address the 3 remaining gaps from TEST-ANALYSIS.md: ThaiAddressGroup watcher tests, critical-path E2E flow, and AnalyticsPanel logic extraction + unit tests.

**Architecture:** Each gap is independent — parallel execution. ThaiAddressGroup tests mock the `thai-data` package. The E2E flow uses Playwright `page.route()` to intercept API calls (no MSW needed). AnalyticsPanel extraction moves computed logic to pure functions in a new utility file.

**Tech Stack:** Vitest + Vue Test Utils + Playwright (with page.route API mocking)

---

## File Map

| File | Action | Responsibility |
|------|--------|---------------|
| `src/components/ThaiAddressGroup.spec.ts` | Modify | Add 6 tests for zipcode watcher, sub-district select, auto-populate |
| `tests/e2e/fulfillment.spec.ts` | Create | Playwright test: auth → create order → verify in list → view detail |
| `src/lib/analytics-utils.ts` | Create | Pure functions: computeKpis, computeRegionPerformance, computeStatusDistribution, computeStatusPieData, computeCumulativeData |
| `src/lib/analytics-utils.spec.ts` | Create | Unit tests for all 5 exported functions |
| `src/components/AnalyticsPanel.vue` | Modify | Import from analytics-utils.ts instead of inline computeds |

---

## Task 1: Strengthen ThaiAddressGroup.spec.ts

**Files:**
- Modify: `src/components/ThaiAddressGroup.spec.ts`

### Step 1: Write tests for the 3 lookup status states

```ts
// Add to top of existing file
import { describe, it, expect, vi, beforeEach } from "vitest";

const mockGetSubDistrictNames = vi.fn();
const mockGetDistrictNames = vi.fn();
const mockGetProvinceName = vi.fn();

vi.mock("thai-data", () => ({
  getSubDistrictNames: mockGetSubDistrictNames,
  getDistrictNames: mockGetDistrictNames,
  getProvinceName: mockGetProvinceName,
}));
```

Placeholder addresses used in tests:
- "10100" → 2 sub-districts: ["Phra Nakhon", "San Chao"]
- "10200" → 0 results
- "99999" → no results

```ts
describe("ThaiAddressGroup zipcode watcher", () => {
  beforeEach(() => {
    mockGetSubDistrictNames.mockReset();
    mockGetDistrictNames.mockReset();
    mockGetProvinceName.mockReset();
  });

  it("shows 'Enter zipcode first' when zipcode is shorter than 5 digits", () => {
    const wrapper = mount(ThaiAddressGroup, {
      props: {
        label: "Sender",
        modelValue: { name: "", zipcode: "123", subDistrict: "", district: "", province: "" },
        errors: {},
      },
      global: { stubs },
    });
    expect(wrapper.text()).toContain("Enter zipcode first");
  });

  it("shows 'No sub-districts found' when 5-digit zipcode has no results", () => {
    mockGetSubDistrictNames.mockReturnValue([]);
    const wrapper = mount(ThaiAddressGroup, {
      props: {
        label: "Sender",
        modelValue: { name: "", zipcode: "99999", subDistrict: "", district: "", province: "" },
        errors: {},
      },
      global: { stubs },
    });
    expect(wrapper.text()).toContain("No sub-districts found");
  });

  it("renders Select with sub-districts when zipcode has results", () => {
    mockGetSubDistrictNames.mockReturnValue(["Phra Nakhon", "San Chao"]);
    const wrapper = mount(ThaiAddressGroup, {
      props: {
        label: "Sender",
        modelValue: { name: "", zipcode: "10100", subDistrict: "", district: "", province: "" },
        errors: {},
      },
      global: { stubs },
    });
    expect(wrapper.text()).toContain("Phra Nakhon");
    expect(wrapper.text()).toContain("San Chao");
  });

  it("auto-selects sub-district when only one result", async () => {
    mockGetSubDistrictNames.mockReturnValue(["Phra Nakhon"]);
    mockGetDistrictNames.mockReturnValue(["Phra Nakhon District"]);
    mockGetProvinceName.mockReturnValue("Bangkok");
    const wrapper = mount(ThaiAddressGroup, {
      props: {
        label: "Sender",
        modelValue: { name: "", zipcode: "10100", subDistrict: "", district: "", province: "" },
        errors: {},
      },
      global: { stubs },
    });
    await wrapper.vm.$nextTick();
    expect(wrapper.emitted("update:modelValue")).toBeTruthy();
    const emitted = wrapper.emitted("update:modelValue")!;
    const last = emitted[emitted.length - 1][0] as Record<string, string>;
    expect(last.subDistrict).toBe("Phra Nakhon");
    expect(last.district).toBe("Phra Nakhon District");
    expect(last.province).toBe("Bangkok");
  });

  it("clears address fields when zipcode changes to a different 5-digit zip", async () => {
    mockGetSubDistrictNames.mockReturnValue(["Dusit"]);
    mockGetDistrictNames.mockReturnValue(["Dusit District"]);
    mockGetProvinceName.mockReturnValue("Bangkok");
    const wrapper = mount(ThaiAddressGroup, {
      props: {
        label: "Sender",
        modelValue: {
          name: "", zipcode: "10200", subDistrict: "Old", district: "Old District", province: "Old Province",
        },
        errors: {},
      },
      global: { stubs },
    });
    // Simulate typing a new zipcode by updating the prop
    await wrapper.setProps({
      modelValue: {
        name: "", zipcode: "10300", subDistrict: "Old", district: "Old District", province: "Old Province",
      },
    });
    await wrapper.vm.$nextTick();
    expect(wrapper.emitted("update:modelValue")).toBeTruthy();
    const emitted = wrapper.emitted("update:modelValue")!;
    const last = emitted[emitted.length - 1][0] as Record<string, string>;
    expect(last.subDistrict).toBe("");
    expect(last.district).toBe("");
    expect(last.province).toBe("");
  });

  it("clears address fields when zipcode goes from 5-digit to shorter", async () => {
    mockGetSubDistrictNames.mockReturnValue(["Dusit"]);
    const wrapper = mount(ThaiAddressGroup, {
      props: {
        label: "Sender",
        modelValue: {
          name: "", zipcode: "10200", subDistrict: "Dusit", district: "Dusit District", province: "Bangkok",
        },
        errors: {},
      },
      global: { stubs },
    });
    await wrapper.setProps({
      modelValue: {
        name: "", zipcode: "10", subDistrict: "Dusit", district: "Dusit District", province: "Bangkok",
      },
    });
    await wrapper.vm.$nextTick();
    const emitted = wrapper.emitted("update:modelValue")!;
    const last = emitted[emitted.length - 1][0] as Record<string, string>;
    expect(last.subDistrict).toBe("");
    expect(last.district).toBe("");
    expect(last.province).toBe("");
  });
});
```

### Step 2: Run tests to verify they pass

```bash
bun run test -- src/components/ThaiAddressGroup.spec.ts
```
Expected: All 8 tests pass (2 existing + 6 new).

---

## Task 2: Implement Critical Path E2E Flow

**Files:**
- Create: `tests/e2e/fulfillment.spec.ts`

### Step 1: Write the E2E test

The test uses `page.route()` to intercept all backend API calls. Key assumptions:
- Backend runs at `http://localhost:8080` (production-like) — but Playwright tests against the Vite dev server only
- The Vite dev server proxies `/api/*` to the backend
- Since there's no backend running, `page.route()` intercepts ALL `**/api/**` requests

The flow:
1. Navigate to `/orders` (orders page)
2. App boots, calls `GET /api/auth/me` → page.route returns user data (simulating logged-in state)
3. "New Order" button should be visible (auth.state has user)
4. Click "New Order" → navigate to `/orders/create`
5. Fill form fields and submit
6. `POST /api/shipments` → page.route returns mock order
7. Navigate to `/orders` — verify new order appears in table
8. Click order link → navigate to `/orders/ORD-E2E`
9. `GET /api/shipments/ORD-E2E` → page.route returns detail
10. Verify detail page shows tracking number

```ts
import { test, expect } from "@playwright/test";

const MOCK_USER = {
  id: 1,
  name: "Test Admin",
  email: "admin@test.com",
  role: "admin",
  created_at: "2026-01-01T00:00:00Z",
};

const MOCK_ORDER = {
  id: "ORD-E2E",
  tracking_number: "TH2026E2E",
  status: "pending",
  progress: 0,
  customer: { name: "E2E Customer", zipcode: "10100", sub_district: "Phra Nakhon", district: "Bangkok", province: "Bangkok" },
  receiver: { name: "E2E Receiver", zipcode: "10200", sub_district: "Dusit", district: "Bangkok", province: "Bangkok" },
  origin: "Bangkok",
  destination: "Bangkok",
  carrier: "Thun-u-der Express",
  weight: 5.0,
  items: 2,
  estimated_delivery: null,
  current_coords: null,
  created_at: new Date().toISOString(),
  updated_at: new Date().toISOString(),
  events: [],
};

test.describe("Fulfillment Critical Path", () => {
  test("logs in, creates order, finds it in list, views detail", async ({ page }) => {
    // Intercept auth
    await page.route("**/api/auth/me", async (route) => {
      await route.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify({ data: MOCK_USER }),
      });
    });

    await page.route("**/api/auth/login", async (route) => {
      await route.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify({ data: { user: MOCK_USER } }),
      });
    });

    // Intercept shipment creation
    await page.route("**/api/shipments", async (route) => {
      if (route.request().method() === "POST") {
        await route.fulfill({
          status: 201,
          contentType: "application/json",
          body: JSON.stringify({ data: MOCK_ORDER }),
        });
      } else {
        await route.fulfill({
          status: 200,
          contentType: "application/json",
          body: JSON.stringify({
            data: [MOCK_ORDER],
            pagination: { page: 1, limit: 10, total: 1, totalPages: 1 },
          }),
        });
      }
    });

    // Intercept single order fetch
    await page.route("**/api/shipments/*", async (route) => {
      await route.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify({ data: MOCK_ORDER }),
      });
    });

    // Intercept tracking (prevent 404 noise)
    await page.route("**/api/track/*", async (route) => {
      await route.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify({ data: { shipment: MOCK_ORDER, events: [] } }),
      });
    });

    // Step 1: Navigate to orders page
    await page.goto("/orders");
    await expect(page.locator("text=Shipment manifest")).toBeVisible({ timeout: 15000 });

    // Step 2: Click "New Order" (user is authed, so navigates directly)
    const newOrderBtn = page.getByRole("button", { name: /new order/i });
    await expect(newOrderBtn).toBeVisible();
    await newOrderBtn.click();
    await expect(page).toHaveURL(/\/orders\/create/, { timeout: 10000 });

    // Step 3: Fill sender info (inputs are within the Sender Info fieldset)
    // The form uses Input stubs via shadcn - inputs are found by placeholder
    // Wait for form to render
    await expect(page.locator("text=Sender Info")).toBeVisible();

    // Fill sender name
    const senderNameInput = page.getByPlaceholder("e.g. ประวิทย์ ใจดี");
    await senderNameInput.fill("E2E Customer");

    // Fill sender zipcode
    const senderZipInput = page.getByPlaceholder("e.g. 10200");
    await senderZipInput.fill("10100");

    // Fill receiver name
    const receiverNameInput = page.getByPlaceholder("e.g. ประวิทย์ ใจดี").nth(1);
    await receiverNameInput.fill("E2E Receiver");

    // Fill receiver zipcode
    const receiverZipInput = page.getByPlaceholder("e.g. 10200").nth(1);
    await receiverZipInput.fill("10200");

    // Fill weight
    const weightInput = page.getByPlaceholder("e.g. 12.4");
    await weightInput.fill("5.0");

    // Fill items (by number input)
    const itemsInput = page.locator("input[type='number']").last();
    await itemsInput.fill("2");

    // Step 4: Submit the form
    const submitBtn = page.getByRole("button", { name: /create order/i });
    await submitBtn.click();

    // Step 5: After successful creation, should redirect to orders list
    // The OrderFormView handles redirect after mutation success
    await expect(page).toHaveURL(/\/orders/, { timeout: 10000 });

    // Step 6: Verify the new order appears in the table
    await expect(page.locator("text=ORD-E2E")).toBeVisible({ timeout: 10000 });
    await expect(page.locator("text=TH2026E2E")).toBeVisible();

    // Step 7: Click on the order ID to view details
    await page.locator("a:has-text('ORD-E2E')").first().click();
    await expect(page).toHaveURL(/\/orders\/ORD-E2E/, { timeout: 10000 });
  });
});
```

### Step 2: Run the E2E test

```bash
bun run test:e2e -- tests/e2e/fulfillment.spec.ts --headed
```
Expected: Test passes — browser shows the fulfillment flow executing.

---

## Task 3: Extract AnalyticsPanel Logic + Unit Tests

**Files:**
- Create: `src/lib/analytics-utils.ts`
- Create: `src/lib/analytics-utils.spec.ts`
- Modify: `src/components/AnalyticsPanel.vue`

### Step 1: Create analytics-utils.ts

Extract 5 pure functions from AnalyticsPanel.vue's computed properties:

```ts
import { statusLabels } from "@/lib/orders";
import type { AnalyticsOverview, TimeSeriesData } from "@/lib/api/analytics";
import type { StatusPieEntry } from "@/components/charts/StatusPieChart.vue";
import type { CumulativeEntry } from "@/components/charts/ShipmentsAreaChart.vue";

export interface KpiData {
  total: number;
  onTime: number;
  regions: number;
  avgDeliveryTime: string;
}

export interface RegionEntry {
  name: string;
  total: number;
  pct: number;
}

export interface StatusDistEntry {
  status: string;
  label: string;
  count: number;
  pct: number;
}

export function computeKpis(data: AnalyticsOverview | null): KpiData {
  const total = data?.total ?? 0;
  const delivered = data?.delivered ?? 0;
  const onTime = Math.round((delivered / Math.max(total, 1)) * 100) || 99.9;
  const regions = data?.by_region.length ?? 0;
  return { total, onTime, regions, avgDeliveryTime: "3.2 days" };
}

export function computeRegionPerformance(data: AnalyticsOverview | null): RegionEntry[] {
  if (!data) return [];
  const total = data.total;
  return data.by_region
    .map((r) => ({
      name: r.name,
      total: r.total,
      pct: total > 0 ? Math.round((r.total / total) * 100) : 0,
    }))
    .sort((a, b) => b.total - a.total);
}

export function computeStatusDistribution(data: AnalyticsOverview | null): StatusDistEntry[] {
  if (!data) return [];
  const total = data.total;
  return data.by_status.map((s) => ({
    status: s.status.toLowerCase(),
    label: (statusLabels as Record<string, string>)[s.status.toLowerCase()] ?? s.status,
    count: s.count,
    pct: Math.round((s.count / Math.max(total, 1)) * 100),
  }));
}

const PIE_COLOR_MAP: Record<string, string> = {
  delivered: "var(--color-success)",
  delayed: "var(--color-destructive)",
  in_transit: "var(--color-info)",
  out_for_delivery: "var(--color-primary)",
  pending: "var(--color-muted-foreground)",
  picked_up: "var(--color-warning)",
  departed: "var(--color-secondary)",
};

export function computeStatusPieData(data: AnalyticsOverview | null): StatusPieEntry[] {
  return computeStatusDistribution(data).map((s) => ({
    ...s,
    fill: PIE_COLOR_MAP[s.status] ?? "hsl(var(--muted-foreground))",
  }));
}

export function computeCumulativeData(data: TimeSeriesData | null): CumulativeEntry[] {
  if (!data) return [];
  let running = 0;
  return data.by_month.map((m) => {
    running += m.count;
    return { month: m.month, count: running };
  });
}
```

### Step 2: Write unit tests

```ts
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
      total: 0, active: 0, delivered: 0,
      by_status: [], by_region: [],
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
    expect(result[0]).toMatchObject({
      status: "delivered", count: 80, pct: 80,
    });
    expect(result[1]).toMatchObject({
      status: "in_transit", count: 15, pct: 15,
    });
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
      total: 1, active: 1, delivered: 0,
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
```

### Step 3: Update AnalyticsPanel.vue

Replace the 5 inline computed properties with imports from analytics-utils.ts.

Remove these imports at the top (no longer needed directly in the component):
```ts
import { statusLabels } from "@/lib/orders";
import type { CumulativeEntry } from "@/components/charts/ShipmentsAreaChart.vue";
import type { StatusPieEntry } from "@/components/charts/StatusPieChart.vue";
```

Add:
```ts
import {
  computeKpis,
  computeRegionPerformance,
  computeStatusPieData,
  computeCumulativeData,
} from "@/lib/analytics-utils";
```

Replace the 5 computed properties:
```ts
const kpis = computed(() => computeKpis(analytics.value));
const regionPerformance = computed(() => computeRegionPerformance(analytics.value));
const statusDistribution = computed(() => computeStatusDistribution(analytics.value));
const statusPieData = computed(() => computeStatusPieData(analytics.value));
const cumulativeData = computed(() => computeCumulativeData(timeSeries.value));
```

Wait — `statusDistribution` is no longer needed as a separate computed since `computeStatusPieData` internally calls `computeStatusDistribution`. But looking at the template, `statusDistribution` is NOT used in the template — only `statusPieData` and `cumulativeData` and `regionPerformance` and `kpis`.

Check the AnalyticsPanel template — `statusDistribution` is not referenced. So we can remove it entirely. The only computed that was using it was `statusPieData`.

### Step 4: Run tests

```bash
bun run test -- src/lib/analytics-utils.spec.ts
bun run test:coverage
```
Expected: 15+ new tests pass. Coverage should improve for branch coverage.
