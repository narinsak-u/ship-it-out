import type { Page } from "@playwright/test";

export async function signIn(page: Page, email: string, password: string) {
  await page.getByRole("button", { name: "Sign in" }).click();
  await page.getByLabel("Email").fill(email);
  await page.getByLabel("Password").fill(password);
  await page
    .getByRole("dialog")
    .locator("form")
    .getByRole("button", { name: "Sign in" })
    .click();
}
