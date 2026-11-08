import { test, expect } from "@playwright/test";

test.describe("Next carrier hubs", () => {
  test("renders public hubs on a direct carrier URL", async ({ page }) => {
    const response = await page.goto("/carriers?tab=hubs");
    expect(response).not.toBeNull();
    const html = await response!.text();
    expect(html).toContain("Bangkok Hub");
    expect(html).toContain("HUB-001");

    await expect(page.getByRole("heading", { name: "Carrier operations" })).toBeVisible();
    await expect(page.getByText("Bangkok Hub")).toBeVisible();
    await expect(page.getByText("Chonburi Hub")).toBeVisible();
  });

  test("paginates hub rows, filters hubs, and shows status counts and empty state", async ({
    page,
  }) => {
    await page.goto("/carriers?tab=hubs");
    const counts = page.getByLabel("Hub status counts");
    await expect(counts).toContainText("active");
    await expect(counts).toContainText("maintenance");

    await expect(page.getByText("HUB-012")).toHaveCount(0);
    await page.getByRole("button", { name: "Next page" }).click();
    await expect(page.getByText("HUB-012")).toBeVisible();
    await page.getByRole("button", { name: "Previous page" }).click();
    await expect(page.getByText("HUB-012")).toHaveCount(0);

    const search = page.getByRole("searchbox", { name: "Search hubs" });
    await search.fill("Chonburi");
    await expect(page.getByText("Chonburi Hub")).toBeVisible();
    await expect(page.getByText("Bangkok Hub")).toHaveCount(0);
    await expect(page.getByText(/hubs · Page/)).toBeVisible();

    await search.fill("no matching hub");
    await expect(page.getByText("No hubs match your filters.")).toBeVisible();
  });
});
