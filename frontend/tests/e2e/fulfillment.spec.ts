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
  trackingNumber: "TH2026E2E",
  customer: {
    name: "E2E Customer",
    zipcode: "10100",
    subDistrict: "Phra Nakhon",
    district: "Bangkok",
    province: "Bangkok",
  },
  receiver: {
    name: "E2E Receiver",
    zipcode: "10200",
    subDistrict: "Dusit",
    district: "Bangkok",
    province: "Bangkok",
  },
  origin: "Bangkok",
  destination: "Bangkok",
  carrier: "Thun-u-der Express",
  weight: 5.0,
  items: 2,
  status: "pending",
  progress: 0,
  estimatedDelivery: null,
  currentCoords: null,
  createdAt: new Date().toISOString(),
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

    // Intercept all /api/shipments requests, dispatch by method and path
    await page.route("**/api/shipments**", async (route) => {
      const url = new URL(route.request().url());
      const path = url.pathname;
      const method = route.request().method();
      // /api/shipments (list) vs /api/shipments/ORD-E2E (detail)
      const isDetail = path !== "/api/shipments" && path.startsWith("/api/shipments/");

      if (method === "POST") {
        await route.fulfill({
          status: 201,
          contentType: "application/json",
          body: JSON.stringify({ data: MOCK_ORDER }),
        });
      } else if (isDetail) {
        await route.fulfill({
          status: 200,
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

    // Intercept tracking
    await page.route("**/api/track/*", async (route) => {
      await route.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify({ data: { shipment: MOCK_ORDER, events: [] } }),
      });
    });

    // Intercept OpenCage geocoding API
    await page.route("https://api.opencagedata.com/geocode/v1/json*", async (route) => {
      await route.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify({ results: [{ geometry: { lat: 13.75, lng: 100.5 } }] }),
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

    // Step 3: Fill sender info
    await expect(page.locator("text=Sender Info")).toBeVisible();

    const senderGroup = page.getByRole("group", { name: "Sender Info" });
    await senderGroup.getByPlaceholder("e.g. ประวิทย์ ใจดี").fill("E2E Customer");
    await senderGroup.getByPlaceholder("e.g. 10200").fill("10100");

    // Wait for sub-district select to appear and select the first option
    await expect(senderGroup.getByRole("combobox")).toBeVisible({ timeout: 5000 });
    await senderGroup.getByRole("combobox").click();
    await expect(page.getByRole("option").first()).toBeVisible({ timeout: 5000 });
    await page.getByRole("option").first().click();

    // Step 4: Fill receiver info
    const receiverGroup = page.getByRole("group", { name: "Receiver Info" });
    await receiverGroup.getByPlaceholder("e.g. ประวิทย์ ใจดี").fill("E2E Receiver");
    await receiverGroup.getByPlaceholder("e.g. 10200").fill("10200");

    // Wait for sub-district select to appear and select the first option
    await expect(receiverGroup.getByRole("combobox")).toBeVisible({ timeout: 5000 });
    await receiverGroup.getByRole("combobox").click();
    await expect(page.getByRole("option").first()).toBeVisible({ timeout: 5000 });
    await page.getByRole("option").first().click();

    // Step 5: Fill parcel info
    const weightInput = page.getByPlaceholder("e.g. 12.4");
    await weightInput.fill("5.0");

    const itemsInput = page.locator("input[type='number']").last();
    await itemsInput.fill("2");

    // Step 6: Submit the form
    const submitBtn = page.getByRole("button", { name: /create order/i });
    await expect(submitBtn).toBeEnabled({ timeout: 5000 });
    await submitBtn.click();

    // Step 7: Should redirect to orders list after creation
    await expect(page).toHaveURL(/\/orders/, { timeout: 10000 });

    // Step 8: Verify the new order appears in the table
    await expect(page.locator("text=ORD-E2E")).toBeVisible({ timeout: 10000 });
    await expect(page.locator("text=TH2026E2E")).toBeVisible();

    // Step 9: Click on the order ID to view details
    await page.locator("a:has-text('ORD-E2E')").first().click();
    await expect(page).toHaveURL(/\/orders\/ORD-E2E/, { timeout: 10000 });
  });
});
