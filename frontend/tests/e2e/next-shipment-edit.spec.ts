import { expect, test, type Page } from "@playwright/test";
import { signIn } from "./auth";

async function allowGeocoding(page: Page) {
  await page.route("**/api/geocode", (route) =>
    route.fulfill({
      status: 200,
      contentType: "application/json",
      body: JSON.stringify({ data: { lat: 13.75, lng: 100.5 } }),
    }),
  );
}

test("staff can open a direct edit URL and save refreshed Shipment details", async ({ page }) => {
  await page.request.post("http://127.0.0.1:8081/__test/shipment?reset=1");
  await page.goto("/orders/ORD-TEST1/edit");
  await signIn(page, "staff@shipitout.example", "preview-staff-password");
  await allowGeocoding(page);

  await expect(page.getByLabel("Sender Info name")).toHaveValue("Sender");
  await expect(page.getByLabel("Sender Info zipcode")).toHaveValue("10900");
  await expect(page.getByLabel("Sender Info sub-district")).toHaveValue("ลาดยาว");
  await expect(page.getByLabel("Sender Info district")).toHaveValue("จตุจักร");
  await expect(page.getByLabel("Sender Info province")).toHaveValue("กรุงเทพมหานคร");
  await expect(page.getByLabel("Receiver Info zipcode")).toHaveValue("10800");
  await expect(page.getByLabel("Receiver Info sub-district")).toHaveValue("วงศ์สว่าง");
  await expect(page.getByLabel("Receiver Info district")).toHaveValue("บางซื่อ");
  await expect(page.getByLabel("Weight/Kg")).toHaveValue("2");
  await expect(page.getByLabel("Items")).toHaveValue("2");
  await expect(page.getByLabel("Estimated Delivery")).toHaveValue("2026-12-01");
  await page.getByLabel("Sender Info name").fill("Edited Sender");
  await page.getByLabel("Weight/Kg").fill("4");
  await page.getByLabel("Items").fill("3");
  await page.getByLabel("Estimated Delivery").fill("2026-12-02");
  await page.getByRole("button", { name: "Save Changes" }).click();

  await expect(page).toHaveURL("/orders/ORD-TEST1?updated=1");
  await expect(page.getByText("Shipment updated.", { exact: true })).toBeVisible();
  await expect(page.getByText("Edited Sender", { exact: true })).toBeVisible();
});

test("failed Shipment update shows error and retains entered values", async ({ page }) => {
  await page.request.post("http://127.0.0.1:8081/__test/shipment?reset=1");
  await page.request.post("http://127.0.0.1:8081/__test/shipment?update=failure");
  await page.goto("/orders/ORD-TEST1/edit");
  await signIn(page, "staff@shipitout.example", "preview-staff-password");
  await allowGeocoding(page);
  await page.getByLabel("Sender Info name").fill("Retained Sender");
  await page.getByRole("button", { name: "Save Changes" }).click();

  await expect(page.getByText("shipment update unavailable", { exact: true })).toBeVisible();
  await expect(page).toHaveURL(/\/orders\/ORD-TEST1\/edit$/);
  await expect(page.getByLabel("Sender Info name")).toHaveValue("Retained Sender");
});

test("unsigned and demo users cannot edit through the UI or API", async ({ page }) => {
  await page.request.post("http://127.0.0.1:8081/__test/shipment?reset=1");
  await page.goto("/orders/ORD-TEST1/edit");
  await expect(page.getByText(/Staff access required to edit shipments/)).toBeVisible();
  await expect(page.getByRole("button", { name: "Save Changes" })).toHaveCount(0);
  const unsignedResponse = await page.request.put("/api/shipments/ORD-TEST1", { data: {} });
  expect(unsignedResponse.status()).toBe(401);

  await signIn(page, "demo@shipitout.example", "demo-read-only-2026");
  await expect(page.getByText("customer (Ship It Out Demo) · read-only")).toBeVisible();
  const demoResponse = await page.request.put("/api/shipments/ORD-TEST1", { data: {} });
  expect(demoResponse.status()).toBe(403);
});
