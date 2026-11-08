import { expect, test } from "@playwright/test";

test("demo session uses proxy cookie, restores on refresh, and cannot write", async ({ page }) => {
  await page.goto("/");
  await expect(page.getByRole("heading", { name: "Recent shipments" })).toBeVisible();
  await expect(page.getByRole("button", { name: "Sign in" })).toBeVisible();
  await page.getByRole("button", { name: "Sign in" }).click();

  const dialog = page.getByRole("dialog");
  await expect(
    dialog.getByText(
      "Read-only demo account. Email: demo@shipitout.example · Password: demo-read-only-2026",
    ),
  ).toBeVisible();
  await expect(dialog.getByLabel("Password")).toHaveAttribute("type", "password");
  await dialog.getByRole("button", { name: "Sign in", exact: true }).last().click();

  await expect(page.getByText("customer (Ship It Out Demo) · read-only")).toBeVisible();
  const cookie = (await page.context().cookies()).find(({ name }) => name === "jwt");
  expect(cookie?.httpOnly).toBe(true);
  expect(cookie?.secure).toBe(true);
  expect(
    await page.evaluate(async () => (await fetch("/api/shipments", { method: "POST" })).status),
  ).toBe(403);

  await page.reload();
  await expect(page.getByText("customer (Ship It Out Demo) · read-only")).toBeVisible();
  await page.getByRole("button", { name: "Sign out" }).click();
  await expect(page.getByRole("button", { name: "Sign in" })).toBeVisible();
  expect((await page.context().cookies()).some(({ name }) => name === "jwt")).toBe(false);
});

test("public signup cannot select a staff role", async ({ page }) => {
  await page.goto("/");
  await page.getByRole("button", { name: "Sign in" }).click();
  await page.getByRole("button", { name: "Sign up" }).click();
  await page.getByLabel("Name").fill("Preview Customer");
  await page.getByLabel("Email").fill(`customer-${Date.now()}@example.test`);
  await page.getByLabel("Password").fill("customer-preview-password");
  await page.getByRole("button", { name: "Create account" }).click();
  await expect(page.getByText("customer (Preview Customer) · read-only")).toBeVisible();

  await page.getByRole("button", { name: "Sign out" }).click();
  const result = await page.evaluate(async () => {
    const response = await fetch("/api/auth/register", {
      method: "POST",
      headers: { "content-type": "application/json" },
      body: JSON.stringify({
        name: "Role Probe",
        email: `role-${Date.now()}@example.test`,
        password: "role-probe-password",
        role: "admin",
      }),
    });
    return { status: response.status, body: await response.json() };
  });
  expect(result.status).toBe(200);
  expect(result.body.data.user.role).toBe("customer");
});

test("staff role comes from session and logout clears it", async ({ page }) => {
  await page.goto("/");
  await page.getByRole("button", { name: "Sign in" }).click();
  await page.getByLabel("Email").fill("staff@shipitout.example");
  await page.getByLabel("Password").fill("preview-staff-password");
  await page.getByRole("button", { name: "Sign in", exact: true }).last().click();
  await expect(page.getByText("staff (Preview Staff) · staff access")).toBeVisible();
  expect(
    await page.evaluate(async () => (await fetch("/api/shipments", { method: "POST" })).status),
  ).toBe(201);

  await page.getByRole("button", { name: "Sign out" }).click();
  await expect(page.getByRole("button", { name: "Sign in" })).toBeVisible();
  expect((await page.context().cookies()).some(({ name }) => name === "jwt")).toBe(false);
});

test("active navigation follows pathname, not query string", async ({ page }) => {
  await page.goto("/");
  await expect(page.getByRole("link", { name: "Home" })).toHaveAttribute("aria-current", "page");
  await page.goto("/carriers?tab=analytics");
  await expect(page.getByRole("link", { name: "Carriers" })).toHaveAttribute(
    "aria-current",
    "page",
  );
  await expect(page.getByRole("link", { name: "Orders" })).not.toHaveAttribute(
    "aria-current",
    "page",
  );
});
