import { defineConfig } from "@playwright/test";

export default defineConfig({
  testDir: "./tests/e2e",
  testMatch: "next-*.spec.ts",
  timeout: 30_000,
  // shortcut: analytics fixture state is shared; isolate it per worker before increasing concurrency.
  workers: 1,
  use: { baseURL: "http://localhost:3000", trace: "retain-on-failure" },
  projects: [{ name: "chromium-next-preview", use: { browserName: "chromium" } }],
  webServer: [
    { command: "node tests/e2e/preview-api.mjs", port: 8081, reuseExistingServer: true },
    {
      command: "bun run dev",
      port: 3000,
      reuseExistingServer: true,
      env: { API_BASE_URL: "http://127.0.0.1:8081" },
    },
  ],
});
