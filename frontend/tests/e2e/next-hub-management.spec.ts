import { expect, test, type Page } from "@playwright/test";

async function signIn(page: Page, email: string, password: string) {
  await page.goto("/carriers?tab=hubs");
  await page.getByRole("button", { name: "Sign in" }).click();
  await page.getByLabel("Email").fill(email);
  await page.getByLabel("Password").fill(password);
  await page.getByRole("button", { name: "Sign in", exact: true }).last().click();
}

async function mockHubApi(page: Page) {
  let hubs = [
    {
      id: "HUB-001",
      name: "Bangkok Hub",
      carrierId: "THUN",
      address: "Bang Rak, Bangkok",
      coords: { lat: 13.7563, lng: 100.5018 },
      capacity: 1200,
      currentUtilization: 840,
      status: "active",
    },
  ];
  await page.route("**/api/geocode", (route) =>
    route.fulfill({
      status: 200,
      contentType: "application/json",
      body: JSON.stringify({ data: { lat: 13.75, lng: 100.5 } }),
    }),
  );
  await page.route("**/api/hubs**", async (route) => {
    const request = route.request();
    const url = new URL(request.url());
    if (request.method() === "GET") {
      return route.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify({ success: true, data: hubs }),
      });
    }
    if (request.method() === "POST") {
      const body = request.postDataJSON();
      hubs = [...hubs, { ...body, id: "HUB-TEST" }];
      return route.fulfill({
        status: 201,
        contentType: "application/json",
        body: JSON.stringify({ success: true, data: hubs.at(-1) }),
      });
    }
    if (request.method() === "PUT") {
      const body = request.postDataJSON();
      const id = decodeURIComponent(url.pathname.split("/").at(-1)!);
      hubs = hubs.map((hub) => (hub.id === id ? { ...hub, ...body } : hub));
      return route.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify({
          success: true,
          data: hubs.find(({ id: hubId }) => hubId === "HUB-TEST"),
        }),
      });
    }
    if (request.method() === "DELETE") {
      const id = decodeURIComponent(url.pathname.split("/").at(-1)!);
      hubs = hubs.filter((hub) => hub.id !== id);
      return route.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify({ success: true, data: { message: "hub deleted" } }),
      });
    }
    return route.fulfill({ status: 405, body: "unsupported" });
  });
}

test("staff can create, edit, and confirm deletion of a hub", async ({ page }) => {
  await mockHubApi(page);
  await signIn(page, "staff@shipitout.example", "preview-staff-password");
  await expect(page.getByRole("button", { name: "Add Hub" })).toBeVisible();
  await page.getByRole("button", { name: "Add Hub" }).click();
  const dialog = page.getByRole("dialog");
  await dialog.getByLabel("Name").fill("Test Facility");
  await dialog.getByLabel("Province").selectOption("Sakon Nakhon");
  await dialog.getByLabel("Capacity (units)").fill("200");
  await dialog.getByLabel("Current Utilization").fill("40");
  await dialog.getByRole("button", { name: "Create Hub" }).click();
  await expect(page.getByText("Test Facility", { exact: true })).toBeVisible();

  const row = page.getByRole("row", { name: /Test Facility/ });
  await row.getByRole("button", { name: "Edit Test Facility" }).click();
  await page.getByRole("dialog").getByLabel("Name").fill("Updated Facility");
  await page.getByRole("dialog").getByRole("button", { name: "Update Hub" }).click();
  await expect(page.getByText("Updated Facility", { exact: true })).toBeVisible();

  await page
    .getByRole("row", { name: /Updated Facility/ })
    .getByRole("button", { name: "Delete Updated Facility" })
    .click();
  const confirmation = page.getByRole("dialog", { name: "Delete Hub" });
  await expect(confirmation).toContainText("This action cannot be undone");
  await confirmation.getByRole("button", { name: "Delete Hub" }).click();
  await expect(page.getByText("Updated Facility", { exact: true })).toHaveCount(0);
});

test("read-only demo cannot see hub management controls", async ({ page }) => {
  await mockHubApi(page);
  await signIn(page, "demo@shipitout.example", "demo-read-only-2026");
  await expect(page.getByText("customer (Ship It Out Demo) · read-only")).toBeVisible();
  await expect(page.getByRole("button", { name: "Add Hub" })).toHaveCount(0);
  await expect(page.getByRole("button", { name: /^Edit / })).toHaveCount(0);
  await expect(page.getByRole("button", { name: /^Delete / })).toHaveCount(0);
  const geocode = await page.request.post("/api/geocode", { data: { province: "Sakon Nakhon" } });
  expect(geocode.status()).toBe(403);
});

test("hub form shows geocoding failures inline and does not save fake coordinates", async ({
  page,
}) => {
  await mockHubApi(page);
  await page.route("**/api/geocode", (route) =>
    route.fulfill({
      status: 422,
      contentType: "application/json",
      body: JSON.stringify({
        error: "Could not resolve this address. Check the fields and try again.",
      }),
    }),
  );
  await signIn(page, "staff@shipitout.example", "preview-staff-password");
  await page.getByRole("button", { name: "Add Hub" }).click();
  const dialog = page.getByRole("dialog");
  await dialog.getByLabel("Name").fill("Unresolved Facility");
  await dialog.getByLabel("Province").selectOption("Sakon Nakhon");
  await dialog.getByRole("button", { name: "Create Hub" }).click();
  await expect(dialog.getByRole("alert")).toContainText("Could not resolve this address");
  await expect(dialog).toBeVisible();
});
