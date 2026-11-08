import { expect, test, type Page } from "@playwright/test";

async function resetManifest(page: Page) {
  await page.request.post("http://127.0.0.1:8081/__test/manifest?reset=1");
}

async function signIn(page: Page, email: string, password: string) {
  await page.getByRole("button", { name: "Sign in" }).click();
  await page.getByLabel("Email").fill(email);
  await page.getByLabel("Password").fill(password);
  await page.getByRole("button", { name: "Sign in", exact: true }).last().click();
}

test.describe("shipment manifest", () => {
  test.beforeEach(async ({ page }) => resetManifest(page));

  test("search, status filters, and server pagination narrow persisted shipments", async ({
    page,
  }) => {
    await page.goto("/orders");
    const rows = page.getByRole("row");
    await expect(rows).toHaveCount(11);
    await expect(rows.getByText("ORD-MANIFEST-01")).toBeVisible();
    await page.getByRole("button", { name: "Next page" }).click();
    await expect(rows.getByText("ORD-MANIFEST-11")).toBeVisible();
    await expect(rows.getByText("ORD-MANIFEST-01")).toHaveCount(0);

    await page.getByRole("button", { name: "In Transit" }).click();
    await expect(rows.getByText("ORD-MANIFEST-01")).toBeVisible();
    await expect(page.getByText("Status: In Transit")).toBeVisible();

    await page.getByRole("searchbox", { name: "Search shipments" }).fill("Receiver 10");
    await expect(rows.getByText("ORD-MANIFEST-10")).toBeVisible();
    await expect(rows.getByText("ORD-MANIFEST-02")).toHaveCount(0);
    await page.getByRole("searchbox", { name: "Search shipments" }).fill("no matching shipment");
    await expect(page.getByText("No shipments match your filters.")).toBeVisible();
  });

  test("staff confirms deletion and sees refreshed persisted manifest", async ({ page }) => {
    await page.goto("/orders");
    await signIn(page, "staff@shipitout.example", "preview-staff-password");
    await expect(
      page.getByRole("button", { name: "Delete shipment ORD-MANIFEST-01" }),
    ).toBeVisible();
    await page.getByRole("button", { name: "Delete shipment ORD-MANIFEST-01" }).click();
    await expect(page.getByRole("dialog", { name: "Delete shipment?" })).toBeVisible();
    await page.getByRole("dialog").getByRole("button", { name: "Delete shipment" }).click();

    await expect(page.getByRole("status")).toHaveText("Shipment deleted.");
    await expect(page.getByRole("row").getByText("ORD-MANIFEST-01")).toHaveCount(0);
    const response = await page.request.get("http://127.0.0.1:8081/api/shipments?page=1&limit=10");
    const result = await response.json();
    expect(result.data.some((shipment: { id: string }) => shipment.id === "ORD-MANIFEST-01")).toBe(
      false,
    );
  });

  test("staff sees delete failure without losing shipment", async ({ page }) => {
    await page.request.post("http://127.0.0.1:8081/__test/manifest?delete=failure");
    await page.goto("/orders");
    await signIn(page, "staff@shipitout.example", "preview-staff-password");
    await page.getByRole("button", { name: "Delete shipment ORD-MANIFEST-01" }).click();
    await page.getByRole("dialog").getByRole("button", { name: "Delete shipment" }).click();

    await expect(
      page.getByRole("alert").filter({ hasText: "shipment delete unavailable" }),
    ).toBeVisible();
    await expect(page.getByRole("row").getByText("ORD-MANIFEST-01")).toBeVisible();
  });

  test("demo account has no delete control and direct delete stays denied", async ({ page }) => {
    await page.goto("/orders");
    await signIn(page, "demo@shipitout.example", "demo-read-only-2026");
    await expect(page.getByText("customer (Ship It Out Demo) · read-only")).toBeVisible();
    await expect(
      page.getByRole("row").getByRole("button", { name: /delete shipment/i }),
    ).toHaveCount(0);
    const directDelete = await page.evaluate(async () => {
      const response = await fetch("/api/shipments/ORD-MANIFEST-01", { method: "DELETE" });
      return response.status;
    });
    expect(directDelete).toBe(403);
    await page.reload();
    await expect(page.getByRole("row").getByText("ORD-MANIFEST-01")).toBeVisible();
  });
  test("keyboard keeps focus inside the delete dialog and Escape cancels", async ({ page }) => {
    await resetManifest(page);
    await page.goto("/orders");
    await signIn(page, "staff@shipitout.example", "preview-staff-password");
    await page.getByRole("button", { name: "Delete shipment ORD-MANIFEST-01" }).click();
    const dialog = page.getByRole("dialog", { name: "Delete shipment?" });
    await expect(dialog).toBeVisible();
    await page.keyboard.press("Tab");
    expect(await dialog.evaluate((element) => element.contains(document.activeElement))).toBe(true);
    await page.keyboard.press("Escape");
    await expect(dialog).toBeHidden();
    await expect(page.getByRole("button", { name: "Delete shipment ORD-MANIFEST-01" })).toBeFocused();
    await expect(page.getByRole("row").getByText("ORD-MANIFEST-01")).toBeVisible();
  });

});
