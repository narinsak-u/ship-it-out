import { test, expect } from "@playwright/test";

test.describe("Navigation", () => {
  test("navigates between pages via header links", async ({ page }) => {
    await page.goto("/");

    await page.locator("a", { hasText: "Orders" }).click();
    await expect(page).toHaveURL(/\/orders/);

    await page.locator("a", { hasText: "Carriers" }).click();
    await expect(page).toHaveURL(/\/carriers/);

    await page.locator("a", { hasText: "Home" }).click();
    await expect(page).toHaveURL(/\/$/);
  });
});
