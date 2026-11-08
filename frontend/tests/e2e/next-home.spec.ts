import { expect, test } from "@playwright/test";

test("direct home URL renders public metrics and recent shipment in server HTML", async ({
  page,
}) => {
  const response = await page.goto("/");
  const html = await response?.text();
  expect(html).toContain("Move fast");
  expect(html).toContain("ORD-TEST1");
  expect(html).toContain("TH-TEST-1");
  expect(html).toContain(">9<");
  await expect(page.getByRole("heading", { name: "Recent shipments" })).toBeVisible();
});

test("tracking number and order ID open shipment detail", async ({ page }) => {
  await page.goto("/");
  const search = page.getByRole("textbox", { name: "Tracking number or order ID" });
  await search.fill("TH-TEST-1");
  await page.getByRole("button", { name: "Track" }).click();
  await expect(page).toHaveURL(/\/orders\/ORD-TEST1$/);

  await page.goto("/");
  await page.getByRole("textbox", { name: "Tracking number or order ID" }).fill("ORD-TEST1");
  await page.getByRole("button", { name: "Track" }).click();
  await expect(page).toHaveURL(/\/orders\/ORD-TEST1$/);
});

test("failed search follows orders fallback", async ({ page }) => {
  await page.goto("/");
  await page.getByRole("textbox", { name: "Tracking number or order ID" }).fill("not-found");
  await page.getByRole("button", { name: "Track" }).click();
  await expect(page).toHaveURL(/\/orders$/);
});
