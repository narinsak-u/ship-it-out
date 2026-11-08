import { expect, test } from "@playwright/test";

test.describe("Active deliveries", () => {
  test("searches and pages public deliveries, then refreshes and polls while open", async ({
    page,
  }) => {
    await page.request.post("http://127.0.0.1:8081/__test/deliveries?reset=1");
    await page.goto("/carriers");
    await expect(page.getByRole("link", { name: "Active Deliveries" })).toHaveAttribute(
      "aria-current",
      "page",
    );
    await expect(page.getByText("ORD-TEST1", { exact: true })).toBeVisible();
    await expect(page.getByText("TH-TEST-1", { exact: true })).toBeVisible();
    await expect(page.locator(".leaflet-map")).toBeVisible();
    await expect(page.getByText(/11 deliveries · Page 1 of 2/)).toBeVisible();
    await page.getByRole("button", { name: "Next page" }).click();
    await expect(page.getByText("ORD-TEST11")).toBeVisible();
    await page.getByRole("button", { name: "Previous page" }).click();

    const search = page.getByRole("searchbox", { name: "Search active deliveries" });
    await search.fill("no matching delivery");
    await expect(page.getByText("No active deliveries match your filters.")).toBeVisible();
    await search.fill("");
    await page.getByRole("button", { name: "Refresh active deliveries" }).click();
    await expect(page.getByText(/updated 0s ago/)).toBeVisible();
    const response = await page.request.get("http://127.0.0.1:8081/__test/delivery-reads");
    const beforePoll = (await response.json()).reads as number;
    await expect
      .poll(
        async () => {
          const result = await page.request.get("http://127.0.0.1:8081/__test/delivery-reads");
          return (await result.json()).reads as number;
        },
        { timeout: 20_000 },
      )
      .toBeGreaterThan(beforePoll);
  });

  test("staff status changes require a hub and move delivery map and table", async ({ page }) => {
    await page.request.post("http://127.0.0.1:8081/__test/deliveries?reset=1");
    await page.goto("/carriers");
    await page.getByRole("button", { name: "Sign in" }).click();
    await page.locator('input[type="email"]').fill("staff@shipitout.example");
    await page.locator('input[type="password"]').fill("preview-staff-password");
    await page.locator("form").getByRole("button", { name: "Sign in" }).click();
    await expect(page.getByText("staff (Preview Staff)")).toBeVisible();

    const row = page.getByRole("row", { name: /^ORD-TEST1\b/ });
    await row.getByLabel("Status for ORD-TEST1", { exact: true }).selectOption("out_for_delivery");
    await row.getByRole("button", { name: "Save ORD-TEST1 status" }).click();
    await expect(
      page.getByText("Select an active hub for this status.", { exact: true }),
    ).toBeVisible();

    await row.getByLabel("Hub for ORD-TEST1", { exact: true }).selectOption("HUB-001");
    await row.getByRole("button", { name: "Save ORD-TEST1 status" }).click();
    await expect(page.getByText("Delivery status updated.", { exact: true })).toBeVisible();
    await expect(row.getByLabel("Status for ORD-TEST1", { exact: true })).toHaveValue(
      "out_for_delivery",
    );
    await expect(page.getByRole("img", { name: /Bangkok Hub/ })).toBeVisible();
    await expect(page.locator('.leaflet-marker-icon[title*="13.7563"]').first()).toBeVisible();
  });

  test("demo account can view deliveries but direct status writes remain forbidden", async ({
    page,
  }) => {
    await page.request.post("http://127.0.0.1:8081/__test/deliveries?reset=1");
    await page.goto("/carriers");
    await page.getByRole("button", { name: "Sign in" }).click();
    await page.locator('input[type="email"]').fill("demo@shipitout.example");
    await page.locator('input[type="password"]').fill("demo-read-only-2026");
    await page.locator("form").getByRole("button", { name: "Sign in" }).click();
    await expect(page.getByText("customer (Ship It Out Demo)")).toBeVisible();
    await expect(page.getByText("ORD-TEST1", { exact: true })).toBeVisible();
    await expect(page.getByLabel("Status for ORD-TEST1", { exact: true })).toHaveCount(0);
    await expect(page.getByLabel("Hub for ORD-TEST1", { exact: true })).toHaveCount(0);
    const response = await page.request.patch("/api/shipments/ORD-TEST1/status", {
      data: { status: "out_for_delivery", hubId: "HUB-001" },
    });
    expect(response.status()).toBe(403);
    const row = page.getByRole("row", { name: /^ORD-TEST1\b/ });
    await expect(row.getByRole("cell").nth(4)).toHaveText("in transit");
  });
});
