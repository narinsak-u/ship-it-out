import { expect, test, type Page } from "@playwright/test";

async function signIn(page: Page, email: string, password: string) {
  await page.goto("/orders/create");
  await page.getByRole("button", { name: "Sign in" }).click();
  await page.getByLabel("Email").fill(email);
  await page.getByLabel("Password").fill(password);
  await page.getByRole("dialog").locator("form").getByRole("button", { name: "Sign in" }).click();
}

test("staff creates a geocoded Shipment and returns to the updated manifest", async ({ page }) => {
  await page.request.post("http://127.0.0.1:8081/__test/manifest?reset=1");
  await page.route("**/api/geocode", (route) =>
    route.fulfill({
      status: 200,
      contentType: "application/json",
      body: JSON.stringify({ data: { lat: 13.75, lng: 100.5 } }),
    }),
  );
  await signIn(page, "staff@shipitout.example", "preview-staff-password");
  await expect(page.getByText("staff (Preview Staff) · staff access")).toBeVisible();

  await page.getByLabel("Sender Info name").fill("Sender One");
  await page.getByLabel("Sender Info zipcode").fill("10900");
  await page.getByLabel("Sender Info sub-district").selectOption({ index: 1 });
  await page.getByLabel("Receiver Info name").fill("Receiver One");
  await page.getByLabel("Receiver Info zipcode").fill("10800");
  await page.getByLabel("Receiver Info sub-district").selectOption({ index: 1 });
  await page.getByLabel("Weight/Kg").fill("2.5");
  await page.getByLabel("Items").fill("3");
  await page.getByRole("button", { name: "Create Shipment" }).click();

  await expect(page).toHaveURL(/\/orders\?created=1$/);
  await expect(page.getByText("Shipment created.", { exact: true })).toBeVisible();
  await expect(page.getByRole("row").filter({ hasText: "Sender One" })).toContainText(
    /ORD-CREATED-\d+/,
  );
});

test("invalid Thai address stays on form and does not geocode", async ({ page }) => {
  let geocodeCalls = 0;
  await page.route("**/api/geocode", async (route) => {
    geocodeCalls++;
    await route.fulfill({
      status: 200,
      contentType: "application/json",
      body: JSON.stringify({ data: { lat: 13.75, lng: 100.5 } }),
    });
  });
  await signIn(page, "staff@shipitout.example", "preview-staff-password");
  await page.getByLabel("Sender Info name").fill("Sender One");
  await page.getByLabel("Sender Info zipcode").fill("10900");
  await page.getByLabel("Sender Info sub-district").selectOption({ index: 1 });
  await page.getByLabel("Sender Info zipcode").fill("123");
  await expect(page.getByLabel("Sender Info sub-district")).toBeDisabled();
  await expect(page.getByLabel("Sender Info district")).toHaveValue("");
  await expect(page.getByLabel("Sender Info province")).toHaveValue("");
  await page.getByLabel("Receiver Info name").fill("Receiver One");
  await page.getByLabel("Receiver Info zipcode").fill("10800");
  await page.getByLabel("Receiver Info sub-district").selectOption({ index: 1 });
  await page.getByLabel("Weight/Kg").fill("2");
  await page.getByLabel("Items").fill("1");
  await page.getByRole("button", { name: "Create Shipment" }).click();

  await expect(page).toHaveURL(/\/orders\/create$/);
  await expect(page.getByRole("alert").first()).toContainText("5-digit Thai zipcode");
  expect(geocodeCalls).toBe(0);
});

test("demo sees no create form and Go rejects direct Shipment writes", async ({ page }) => {
  await signIn(page, "demo@shipitout.example", "demo-read-only-2026");
  await expect(page.getByText("customer (Ship It Out Demo) · read-only")).toBeVisible();
  await expect(page.getByRole("button", { name: "Create Shipment" })).toHaveCount(0);
  const response = await page.request.post("/api/shipments", { data: {} });
  expect(response.status()).toBe(403);
});
