import { test, expect } from "@playwright/test";

test.describe("Homepage", () => {
  test("loads and shows hero section", async ({ page }) => {
    await page.goto("/");
    await expect(page.locator("text=Move fast")).toBeVisible();
    await expect(page.locator("text=Live ops console")).toBeVisible();
  });

  test("shows tracking search form", async ({ page }) => {
    await page.goto("/");
    await expect(page.getByPlaceholder("Enter tracking number")).toBeVisible();
    await expect(page.getByRole("button", { name: /track/i })).toBeVisible();
  });

  test("shows stats section", async ({ page }) => {
    await page.goto("/");
    await expect(page.locator("text=Recent shipments")).toBeVisible();
  });
});
