import { expect, test } from "@playwright/test";

const featurePath = "/docs/hextra-features/";

function desktopOnly(projectName: string) {
  test.skip(projectName !== "desktop-chromium", "Feature behavior is viewport-independent.");
}

test("banner dismissal, search, and theme controls work", async ({ page }, testInfo) => {
  desktopOnly(testInfo.project.name);

  await page.goto(featurePath);

  const banner = page.locator(".hextra-banner");
  await expect(banner).toBeVisible();
  const bannerFeatureLink = banner.getByRole("link", { name: "feature guide" });
  await expect(bannerFeatureLink).toHaveAttribute("href", "/docs/hextra-features/");
  expect(await bannerFeatureLink.getAttribute("target")).toBeNull();
  await page.getByRole("button", { name: "Close banner" }).click();
  expect(await page.evaluate(() => localStorage.getItem("hugo-styles-feature-demo-v1"))).toBe("0");
  await page.reload();
  await expect(banner).toBeHidden();

  const search = page.locator(".hextra-search-input").first();
  await expect(search).toBeVisible();
  await search.pressSequentially("migration", { delay: 20 });
  await expect(page.locator(".hextra-search-results").first()).toContainText("Migration Guide");

  const themeToggle = page.getByRole("button", { name: "Change theme" });
  await themeToggle.click();
  const themeMenu = page.locator(".hextra-theme-toggle-options");
  await expect(themeMenu.getByRole("menuitemradio", { name: "Light" })).toBeVisible();
  await expect(themeMenu.getByRole("menuitemradio", { name: "Dark" })).toBeVisible();
  await expect(themeMenu.getByRole("menuitemradio", { name: "System" })).toBeVisible();
  await themeMenu.getByRole("menuitemradio", { name: "Dark" }).click();
  await expect(page.locator("html")).toHaveClass(/dark/);
});

test("synced tabs, code copy, and image zoom work", async ({ context, page }, testInfo) => {
  desktopOnly(testInfo.project.name);
  await context.grantPermissions(["clipboard-read", "clipboard-write"]);
  await page.goto(featurePath);

  const tabLists = page.locator('[role="tablist"][data-tab-group="bash,fish"]');
  await expect(tabLists).toHaveCount(2);
  await tabLists.nth(0).getByRole("tab", { name: "fish" }).click();
  await expect(tabLists.nth(1).getByRole("tab", { name: "fish" })).toHaveAttribute(
    "aria-selected",
    "true",
  );

  const codeBlock = page
    .locator("#syntax-highlighting-and-code-copy")
    .locator("xpath=../following-sibling::div[contains(@class, 'hextra-code-block')][1]");
  const copyButton = codeBlock.getByRole("button", { name: "Copy code" });
  await copyButton.click();
  await expect.poll(() => page.evaluate(() => navigator.clipboard.readText())).toContain(
    'selected = [event for event in events if event.pt > 25]',
  );

  const zoomImage = page.locator('img[data-zoomable][src*="analysis-flow.svg"]').first();
  await zoomImage.click();
  const openedZoomImage = page.locator(
    'img[data-zoomable][src*="analysis-flow.svg"].medium-zoom-image--opened',
  );
  await expect(openedZoomImage).toBeVisible();
});

test("Markdown, PDF, and notebook outputs are local and rendered", async ({ page }, testInfo) => {
  desktopOnly(testInfo.project.name);
  await page.route("https://**", (route) => route.abort());
  await page.goto(featurePath);

  const markdownURL = await page.locator(".hextra-page-context-menu-copy").getAttribute("data-url");
  expect(markdownURL).not.toBeNull();
  const localMarkdownURL = new URL(new URL(markdownURL!).pathname, page.url()).toString();
  const markdownResponse = await page.request.get(localMarkdownURL);
  expect(markdownResponse.ok()).toBeTruthy();
  const markdown = await markdownResponse.text();
  expect(markdown).toContain("# Hextra Feature Guide");
  expect(markdown).toContain("![A four-stage toy analysis");

  const pdfResponse = await page.request.get(
    new URL("particle-analysis-handout.pdf", page.url()).toString(),
  );
  expect(pdfResponse.ok()).toBeTruthy();
  expect(pdfResponse.headers()["content-type"]).toContain("application/pdf");
  await expect(page.locator('iframe[title="PDF viewer"]')).toBeVisible();

  const notebookResponse = await page.request.get(
    new URL("particle-analysis.ipynb", page.url()).toString(),
  );
  expect(notebookResponse.ok()).toBeTruthy();
  expect((await notebookResponse.json()).nbformat).toBe(4);
  await expect(page.locator(".hextra-jupyter-code-cell-outputs")).toContainText(
    "selected=9 signal=4 control=2",
  );
});

test("callout titles render inline Markdown", async ({ page }, testInfo) => {
  desktopOnly(testInfo.project.name);
  await page.goto("/docs/components/");

  const title = page.locator(".lesson-callout-title", {
    hasText: "Run jobs before continuing",
  });
  await expect(title).toHaveCount(1);
  await expect(title).toBeVisible();
  await expect(title.locator("code")).toHaveText("jobs");
});
