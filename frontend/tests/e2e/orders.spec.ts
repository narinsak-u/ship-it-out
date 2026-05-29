import { test, expect } from "@playwright/test";

test.describe("Orders List", () => {
  test("loads and shows page title", async ({ page }) => {
    await page.goto("/orders");
    await expect(page.locator("text=Shipment manifest")).toBeVisible({ timeout: 15000 });
  });
});
