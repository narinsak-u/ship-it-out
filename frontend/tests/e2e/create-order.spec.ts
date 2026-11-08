import { test, expect } from "@playwright/test";

test.describe("Create Order", () => {
  test("unsigned visitors cannot open the create order page directly", async ({ page }) => {
    await page.goto("/orders/create");
    await expect(page).toHaveURL(/\/orders$/);
    await expect(page.getByRole("heading", { name: "Shipment manifest" })).toBeVisible();
  });

  test("staff can open the create order page directly", async ({ page }) => {
    await page.route("**/api/auth/me", (route) =>
      route.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify({ data: { id: 1, name: "Staff", role: "admin" } }),
      }),
    );
    await page.goto("/orders/create");
    await expect(page.getByRole("heading", { name: "Create Order" })).toBeVisible();
  });
});
