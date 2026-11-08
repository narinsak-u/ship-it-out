import { test, expect } from "@playwright/test";

test.describe("Orders List", () => {
  test("loads and shows page title", async ({ page }) => {
    await page.goto("/orders");
    await expect(page.locator("text=Shipment manifest")).toBeVisible({ timeout: 15000 });
  });

  test("customer sees shipment reads without write controls", async ({ page }) => {
    await page.route("**/api/auth/me", (route) =>
      route.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify({
          data: { id: 2, name: "Customer", email: "customer@test.com", role: "customer" },
        }),
      }),
    );
    await page.goto("/orders");
    await expect(page.locator("text=Shipment manifest")).toBeVisible({ timeout: 15000 });
    await expect(page.getByRole("button", { name: "New Order" })).toHaveCount(0);
    await expect(page.getByRole("button", { name: "Edit shipment" })).toHaveCount(0);
    await page.goto("/orders/create");
    await expect(page).toHaveURL(/\/orders$/);
  });
});
