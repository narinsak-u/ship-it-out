import { test, expect } from "@playwright/test";

test.describe("Next carrier analytics", () => {
  test.describe.configure({ mode: "serial" });
  test.afterEach(async ({ page }) => {
    await page.request.post("http://127.0.0.1:8081/__test/analytics?mode=normal");
  });
  test("opens from hubs, deep-links, refreshes, and shows API-backed charts", async ({ page }) => {
    await page.goto("/carriers?tab=hubs");
    await expect(page.getByText("Bangkok Hub")).toBeVisible();
    await page.getByRole("link", { name: "Analytics" }).click();
    await expect(page).toHaveURL(/\/carriers\?tab=analytics$/);
    const response = await page.goto("/carriers?tab=analytics");
    expect(response).not.toBeNull();
    const html = await response!.text();
    expect(html).toContain("126");
    expect(html).toContain("Cumulative shipments:");
    await expect(page.getByText("126", { exact: true })).toBeVisible();
    await expect(page.getByText("65%", { exact: true })).toBeVisible();
    await expect(page.getByText("4", { exact: true })).toBeVisible();
    await expect(page.getByText("3.2 days", { exact: true })).toBeVisible();
    await expect(page.getByText("Phuket", { exact: true })).toBeVisible();
    await expect(page.getByText("0/126 shipments", { exact: true })).toBeVisible();
    await expect(
      page
        .getByRole("region", { name: "Cumulative Shipments Over Time" })
        .getByText("Jun 26", { exact: true }),
    ).toBeVisible();
    await expect(page.getByText("Bangkok", { exact: true })).toBeVisible();
    await expect(
      page.getByRole("img", { name: /Cumulative shipments:.*Jun 26 18.*Oct 26 126/ }),
    ).toBeVisible();
    await expect(
      page.getByRole("img", { name: /Shipments by day: Monday 21.*Sunday 8/ }),
    ).toBeVisible();
    await expect(page.getByRole("img", { name: /Monthly shipments:.*Aug 26 31/ })).toBeVisible();
    await expect(page.getByRole("img", { name: /Shipment statuses:.*delivered 82/ })).toBeVisible();

    await page.reload();
    await page.getByRole("link", { name: "Hubs" }).click();
    await expect(page).toHaveURL(/\/carriers\?tab=hubs$/);
    await expect(page.getByText("Bangkok Hub")).toBeVisible();
  });

  test("shows clear API failure and empty states", async ({ page }) => {
    await page.request.post("http://127.0.0.1:8081/__test/analytics?mode=failure");
    await page.goto("/carriers?tab=analytics");
    await expect(
      page.getByRole("alert").filter({ hasText: "Failed to load analytics data." }),
    ).toBeVisible();

    await page.request.post("http://127.0.0.1:8081/__test/analytics?mode=empty");
    await page.reload();
    await expect(page.getByText("No analytics data available yet.", { exact: true })).toBeVisible();
  });
});
