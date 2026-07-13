import { expect, test } from "@playwright/test";

const episodePath = "/episodes/01-introduction/";

test("episode controls respect desktop and mobile breakpoints", async ({ page }, testInfo) => {
  await page.goto(episodePath);

  const sidebarToggle = page.locator("[data-lesson-sidebar-toggle]");
  const tocToggle = page.locator("[data-lesson-toc-toggle]");
  const audienceToggle = page.locator(".lesson-audience-toggle");
  const title = page.locator("main h1");

  if (testInfo.project.name === "mobile-chromium") {
    await expect(sidebarToggle).toBeHidden();
    await expect(tocToggle).toBeHidden();

    const audienceBox = await audienceToggle.boundingBox();
    const titleBox = await title.boundingBox();
    expect(audienceBox).not.toBeNull();
    expect(titleBox).not.toBeNull();
    expect(audienceBox!.y + audienceBox!.height).toBeLessThanOrEqual(titleBox!.y);
  } else {
    await expect(sidebarToggle).toBeVisible();
    await expect(tocToggle).toBeVisible();
    await expect(page.locator(".hextra-sidebar-container")).toBeVisible();
    await expect(page.locator(".hextra-toc")).toBeVisible();
  }
});

test("audience selection persists and URL state wins", async ({ page }, testInfo) => {
  test.skip(testInfo.project.name !== "desktop-chromium", "Persistence is viewport-independent.");

  await page.goto("/all-in-one/");
  const instructorContent = page.locator('[data-audience="instructor"]').first();
  const instructorButton = page.getByRole("button", { name: "Instructor", exact: true });

  await expect(instructorContent).toBeHidden();
  await instructorButton.click();
  await expect(instructorContent).toBeVisible();
  await expect(instructorButton).toHaveAttribute("aria-pressed", "true");
  await expect(page).toHaveURL(/\?view=instructor$/);
  expect(await page.evaluate(() => localStorage.getItem("hugo-styles:view"))).toBe("instructor");

  await page.reload();
  await expect(instructorContent).toBeVisible();

  await page.goto("/all-in-one/?view=learner");
  await expect(instructorContent).toBeHidden();
  await expect(page.getByRole("button", { name: "Learner", exact: true })).toHaveAttribute(
    "aria-pressed",
    "true",
  );
});

test("audience controls tolerate unavailable local storage", async ({ browser }, testInfo) => {
  test.skip(testInfo.project.name !== "desktop-chromium", "Storage behavior is viewport-independent.");

  const context = await browser.newContext();
  await context.addInitScript(() => {
    Object.defineProperty(window, "localStorage", {
      configurable: true,
      get() {
        throw new DOMException("Storage disabled", "SecurityError");
      },
    });
  });
  const page = await context.newPage();
  await page.goto("http://127.0.0.1:13131/all-in-one/");
  await page.getByRole("button", { name: "Instructor", exact: true }).click();
  await expect(page.locator('[data-audience="instructor"]').first()).toBeVisible();
  await context.close();
});

test("All-in-One headings and table of contents have valid unique targets", async ({ page }, testInfo) => {
  test.skip(testInfo.project.name !== "desktop-chromium", "Document structure is viewport-independent.");

  await page.goto("/all-in-one/");
  await expect(page.locator("main h1")).toHaveCount(1);

  const h2s = page.locator("[data-aio-episode] > h2");
  await expect(h2s).toHaveCount(4);
  expect(await page.locator("[data-aio-episode] h1").count()).toBe(0);
  expect(await page.locator("[data-aio-episode] h2").count()).toBe(4);

  const duplicateIds = await page.locator("[id]").evaluateAll((nodes) => {
    const ids = nodes.map((node) => node.id).filter(Boolean);
    return ids.filter((id, index) => ids.indexOf(id) !== index);
  });
  expect(duplicateIds).toEqual([]);

  const targets = await page.locator('.hextra-toc a[href^="#"]').evaluateAll((links) =>
    links.map((link) => (link as HTMLAnchorElement).hash.slice(1)),
  );
  for (const target of targets) {
    await expect(page.locator(`[id="${target}"]`)).toHaveCount(1);
  }
});

test("back-to-top supports keyboard activation and reduced motion", async ({ page }, testInfo) => {
  test.skip(testInfo.project.name !== "desktop-chromium", "Interaction is viewport-independent.");

  await page.emulateMedia({ reducedMotion: "reduce" });
  await page.goto("/all-in-one/");
  await page.evaluate(() => window.scrollTo(0, document.body.scrollHeight));

  const button = page.locator("#backToTop");
  await expect(button).not.toHaveAttribute("tabindex", "-1");
  await page.evaluate(() => {
    (window as typeof window & { __lastScroll?: ScrollToOptions }).scroll = (options) => {
      window.__lastScroll = options as ScrollToOptions;
    };
  });
  await button.focus();
  await expect(button).toBeFocused();
  await button.press("Enter");
  await expect.poll(() => page.evaluate(() => window.__lastScroll?.behavior)).toBe("auto");
});
