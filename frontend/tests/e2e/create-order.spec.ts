import { test, expect } from "@playwright/test";

test.describe("Create Order", () => {
  test("create order page renders", async ({ page }) => {
    await page.goto("/orders/create");
    await expect(page.getByRole("heading", { name: "Create Order" })).toBeVisible({
      timeout: 10000,
    });
  });
});
