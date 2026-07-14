import { defineConfig, devices } from "@playwright/test";

export default defineConfig({
  testDir: "./tests/browser",
  fullyParallel: false,
  forbidOnly: Boolean(process.env.CI),
  retries: process.env.CI ? 1 : 0,
  reporter: process.env.CI ? "github" : "list",
  use: {
    baseURL: "http://127.0.0.1:13131",
    trace: "retain-on-failure",
  },
  projects: [
    {
      name: "desktop-chromium",
      use: { ...devices["Desktop Chrome"], viewport: { width: 1440, height: 900 } },
    },
    {
      name: "mobile-chromium",
      use: { ...devices["Desktop Chrome"], viewport: { width: 390, height: 844 } },
    },
  ],
  webServer: {
    command:
      "hugo server --config hugo.toml,hugo-docs.toml --disableFastRender --bind 127.0.0.1 --port 13131 --baseURL http://127.0.0.1:13131/",
    url: "http://127.0.0.1:13131/",
    reuseExistingServer: !process.env.CI,
    timeout: 120_000,
  },
});
