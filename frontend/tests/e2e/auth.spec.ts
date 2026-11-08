import { test, expect, type Page } from "@playwright/test";

const demoUser = {
  id: 2,
  name: "Ship It Out Demo",
  email: "demo@shipitout.example",
  role: "customer",
  created_at: "2026-01-01T00:00:00Z",
};

async function mockPublicShipments(page: Page) {
  await page.route("**/api/shipments**", (route) =>
    route.fulfill({
      status: 200,
      contentType: "application/json",
      body: JSON.stringify({
        data: [],
        pagination: { page: 1, limit: 10, total: 0, totalPages: 0 },
      }),
    }),
  );
}

test.describe("Demo authentication", () => {
  test("shows credentials, restores session, logs out, and keeps unsigned browsing public", async ({
    page,
  }) => {
    let authenticated = false;
    await page.route("**/api/auth/me", async (route) => {
      if (authenticated) expect(route.request().headers().cookie).toContain("jwt=demo-session");
      await route.fulfill({
        status: authenticated ? 200 : 401,
        contentType: "application/json",
        body: JSON.stringify(authenticated ? { data: demoUser } : { error: "unauthorized" }),
      });
    });
    await page.route("**/api/auth/login", async (route) => {
      expect(route.request().postDataJSON()).toEqual({
        email: "demo@shipitout.example",
        password: "demo-read-only-2026",
      });
      authenticated = true;
      await route.fulfill({
        status: 200,
        contentType: "application/json",
        headers: { "set-cookie": "jwt=demo-session; Path=/; HttpOnly" },
        body: JSON.stringify({ data: { user: demoUser } }),
      });
    });
    await page.route("**/api/auth/logout", async (route) => {
      authenticated = false;
      await route.fulfill({
        status: 200,
        contentType: "application/json",
        headers: { "set-cookie": "jwt=; Path=/; Max-Age=0; HttpOnly" },
        body: JSON.stringify({ success: true }),
      });
    });
    await mockPublicShipments(page);

    await page.goto("/orders");
    await expect(page.getByRole("heading", { name: "Shipment manifest" })).toBeVisible();
    await expect(page.getByRole("button", { name: /guest/i })).toHaveCount(0);
    await page.getByRole("button", { name: "Sign in" }).click();
    await expect(page.getByText("Read-only demo account")).toBeVisible();
    await expect(page.getByText("demo@shipitout.example", { exact: false })).toBeVisible();
    await expect(page.getByText("demo-read-only-2026", { exact: false })).toBeVisible();
    const password = page.locator('input[type="password"]');
    await expect(password).toHaveCount(1);
    await password.fill("arbitrary-user-password");
    await expect(password).toHaveAttribute("type", "password");
    await password.fill("demo-read-only-2026");
    await page.locator("form").getByRole("button", { name: "Sign In" }).click();
    await expect(page.getByText("customer (Ship It Out Demo)")).toBeVisible();
    await expect(page.getByRole("button", { name: "New Order" })).toHaveCount(0);

    await page.reload();
    await expect(page.getByText("customer (Ship It Out Demo)")).toBeVisible();
    await page.getByRole("button", { name: "Sign out" }).click();
    await expect(page.getByRole("button", { name: "Sign in" })).toBeVisible();
    await expect(page.getByRole("heading", { name: "Shipment manifest" })).toBeVisible();
    await page.reload();
    await expect(page.getByRole("button", { name: "Sign in" })).toBeVisible();
  });

  test("staff sees write controls while demo customer does not", async ({ page }) => {
    await page.route("http://localhost:8080/api/hubs**", (route) =>
      route.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify({ data: [] }),
      }),
    );
    let role = "customer";
    await page.route("**/api/auth/me", (route) =>
      route.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify({ data: { ...demoUser, role } }),
      }),
    );
    await mockPublicShipments(page);
    await page.goto("/orders");
    await expect(page.getByRole("button", { name: "New Order" })).toHaveCount(0);
    await page.goto("/carriers");
    await page.getByRole("button", { name: "Hubs" }).click();
    await expect(page.getByText("Total Hubs")).toBeVisible({ timeout: 15000 });
    await expect(page.getByRole("button", { name: "Add Hub" })).toHaveCount(0);

    role = "staff";
    await page.reload();
    await page.goto("/orders");
    await expect(page.getByRole("button", { name: "New Order" })).toBeVisible();
    await page.goto("/carriers");
    await page.getByRole("button", { name: "Hubs" }).click();
    await expect(page.getByText("Total Hubs")).toBeVisible({ timeout: 15000 });
    await expect(page.getByRole("button", { name: "Add Hub" })).toBeVisible();

    role = "admin";
    await page.reload();
    await page.goto("/orders");
    await expect(page.getByRole("button", { name: "New Order" })).toBeVisible();
  });
});
