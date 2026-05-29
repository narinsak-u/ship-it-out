import { test, expect } from "@playwright/test";

test.describe("Order Detail", () => {
  test("shows not found for non-existent order", async ({ page }) => {
    await page.goto("/orders/NONEXISTENT");
    await expect(page.locator("text=Shipment not found")).toBeVisible({ timeout: 10000 });
  });
});
