import { expect, test } from "@playwright/test";

test("direct detail URL renders API-backed shipment and timeline", async ({ page }) => {
  await page.goto("/orders/ORD-TEST1");
  await expect(page.getByRole("heading", { name: "TH-TEST-1" })).toBeVisible();
  await expect(page.getByText("in transit", { exact: true }).first()).toBeVisible();
  await expect(page.getByText("Bangkok Hub")).toBeVisible();
  await expect(page.getByRole("region", { name: "Shipment route map" })).toBeVisible();
});

test("missing shipment has not-found state without a map", async ({ page }) => {
  await page.goto("/orders/missing");
  await expect(page.getByRole("heading", { name: "Shipment not found" })).toBeVisible({ timeout: 15000 });
  await expect(page.locator(".leaflet-map")).toHaveCount(0);
});
test("warns when CARTO key is missing instead of fetching unusable tiles", async ({ page }) => {
  const tileRequests: string[] = [];
  page.on("request", (request) => {
    if (request.url().includes(".basemaps.cartocdn.com/")) tileRequests.push(request.url());
  });
  await page.goto("/orders/ORD-TEST1");
  await expect(page.getByRole("region", { name: "Shipment route map" })).toBeVisible();
  await expect(page.locator(".map-key-warning")).toHaveText(
    "Map tiles unavailable. Configure NEXT_PUBLIC_CARTO_API_KEY.",
  );
  expect(tileRequests).toHaveLength(0);
});
