import { chromium, expect, test, type Locator } from "@playwright/test";

const episodePath = "/episodes/01-introduction/";
const longEpisodePath = "/episodes/04-physics-doc-features/";

const scrollToHeading = (target: Locator) => target.evaluate((anchor) => {
  const heading = anchor.closest("h2, h3, h4, h5, h6")!;
  const top = window.scrollY + heading.getBoundingClientRect().top - 80;
  window.scrollTo({ top, behavior: "instant" });
  // Give the scroll handler and its animation frame time to save the reading position.
  return new Promise<void>((resolve) => requestAnimationFrame(() => requestAnimationFrame(() => resolve())));
});

test("footer edit links include the repository content path", async ({ page }, testInfo) => {
  test.skip(testInfo.project.name !== "desktop-chromium", "Link generation is viewport-independent.");

  await page.goto(episodePath);

  await expect(page.getByRole("link", { name: "Edit this page" })).toHaveAttribute(
    "href",
    "https://github.com/oer-particle-physics/hugo-styles/edit/main/content/episodes/01-introduction/index.md",
  );
});

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

test("episode progress advances in both menus and persists as a checkmark", async ({ page }, testInfo) => {
  await page.goto(longEpisodePath);
  const progress = page.locator(`[data-episode-progress="${longEpisodePath}"]`);
  const sidebar = page.locator(".hextra-sidebar-container").locator(progress);
  const navbar = page.locator(".hextra-nav-menu-items").locator(progress);
  await expect(navbar).toHaveAttribute("aria-label", /^[1-9]\d?% reached$/);
  const initial = parseInt((await navbar.getAttribute("aria-label"))!, 10);

  await page.locator("[data-episode-content]").evaluate((content) => {
    const { top, height } = content.getBoundingClientRect();
    window.scrollTo({ top: window.scrollY + top + height / 2, behavior: "instant" });
  });
  await expect.poll(async () => parseInt((await navbar.getAttribute("aria-label"))!, 10)).toBeGreaterThan(initial);
  const furthest = (await navbar.getAttribute("aria-label"))!;
  await expect(furthest).toMatch(/% reached$/);
  for (const menu of await sidebar.all()) {
    await expect(menu).toHaveAttribute("aria-label", furthest);
  }

  await page.evaluate(() => window.scrollTo({ top: 0, behavior: "instant" }));
  await page.reload();
  await expect(navbar).toHaveAttribute("aria-label", furthest);

  await page.locator("[data-episode-content]").evaluate((content) => {
    const { bottom } = content.getBoundingClientRect();
    window.scrollTo({ top: Math.ceil(window.scrollY + bottom - window.innerHeight), behavior: "instant" });
  });
  await expect(navbar).toHaveAttribute("aria-label", "Complete");
  await page.reload();

  if (testInfo.project.name === "mobile-chromium") {
    await page.getByRole("button", { name: "Menu", exact: true }).click();
    await expect(sidebar.locator(".lesson-episode-progress-check:visible")).toBeInViewport();
  } else {
    await expect(sidebar.locator(".lesson-episode-progress-check:visible")).toBeVisible();
    await page.getByRole("button", { name: "Examples", exact: true }).click();
    await expect(navbar.locator(".lesson-episode-progress-check")).toBeVisible();
  }
});

test("Continue returns to the last-read heading after scrolling back", async ({ page }) => {
  await page.goto("/episodes/");
  await expect(page.locator("[data-episode-resume]")).toBeHidden();

  await page.goto(longEpisodePath);
  const progress = page.locator(`.hextra-nav-menu-items [data-episode-progress="${longEpisodePath}"]`);
  await expect(progress).toHaveAttribute("aria-label", /^[1-9]\d?% reached$/);
  const laterHeading = "synced-tabs-for-command-variants";
  const previousHeading = "mermaid-for-analysis-flow";
  for (const id of [laterHeading, previousHeading]) {
    await scrollToHeading(page.locator(`[id="${id}"]`));
  }

  const destination = `${longEpisodePath}#${previousHeading}`;
  await page.goto("/");
  await expect(page.getByRole("link", { name: "Continue where you left off" })).toHaveAttribute("href", destination);
  await page.goto("/episodes/");
  await page.getByRole("link", { name: "Continue where you left off" }).click();
  await expect(page).toHaveURL(destination);
  await expect(page.getByRole("heading", { name: "Mermaid for analysis flow" })).toBeInViewport();
});

test("episode progress is saved while a slow image is still loading", async ({ page }) => {
  let releaseImage!: () => void;
  const imagePending = new Promise<void>((resolve) => { releaseImage = resolve; });
  const imagePath = "/slow-reading-image.png";
  await page.route(`**${imagePath}`, async (route) => {
    await imagePending;
    await route.abort().catch(() => {});
  });
  await page.route(`**${longEpisodePath}`, async (route) => {
    const response = await route.fetch();
    const html = (await response.text()).replace("</body>", `<img src="${imagePath}" alt="" hidden></body>`);
    await route.fulfill({ response, body: html });
  });

  const heading = "mermaid-for-analysis-flow";
  try {
    await page.goto(longEpisodePath, { waitUntil: "domcontentloaded" });
    await scrollToHeading(page.locator(`[id="${heading}"]`));
    expect(await page.evaluate(() => document.readyState)).toBe("interactive");

    await page.goto("/episodes/");
    await expect(page.getByRole("link", { name: "Continue where you left off" })).toHaveAttribute(
      "href",
      `${longEpisodePath}#${heading}`,
    );
  } finally {
    releaseImage();
  }
});

test("Continue places a generated heading below the sticky navbar", async ({ page }) => {
  const id = "04-physics-doc-features-objectives";
  const heading = page.locator(`[id="${id}"]`);
  await page.goto(longEpisodePath);
  await scrollToHeading(heading);

  await page.goto("/episodes/");
  const continueLink = page.getByRole("link", { name: "Continue where you left off" });
  await expect(continueLink).toHaveAttribute("href", `${longEpisodePath}#${id}`);
  await continueLink.click();

  await expect.poll(() => heading.evaluate((element) => {
    const navbar = document.querySelector(".hextra-nav-container")!;
    return element.getBoundingClientRect().top - navbar.getBoundingClientRect().bottom;
  })).toBeGreaterThanOrEqual(0);
  await expect(heading).toBeInViewport();
});

test("Continue remembers the episode restored from the back-forward cache", async ({ baseURL, viewport }) => {
  const browser = await chromium.launch({ channel: "chromium", ignoreDefaultArgs: ["--disable-back-forward-cache"] });
  try {
    const page = await browser.newPage({ baseURL, viewport });
    await page.addInitScript(() => window.addEventListener("pageshow", (event) => {
      document.documentElement.dataset.restoredFromCache = String(event.persisted);
    }));
    const firstHeading = "mermaid-for-analysis-flow";
    await page.goto(longEpisodePath);
    await scrollToHeading(page.locator(`[id="${firstHeading}"]`));
    await page.goto("/episodes/02-facilitating-activity/");
    await scrollToHeading(page.locator('[id="callout-gallery"]'));

    await page.goBack({ waitUntil: "commit" });
    await expect(page.locator("html")).toHaveAttribute("data-restored-from-cache", "true");
    await page.goto("/episodes/");
    await expect(page.getByRole("link", { name: "Continue where you left off" })).toHaveAttribute(
      "href",
      `${longEpisodePath}#${firstHeading}`,
    );
  } finally {
    await browser.close();
  }
});

test("Continue offers the next unfinished episode after completion", async ({ page }) => {
  await page.goto(episodePath);
  await page.evaluate(() => window.scrollTo({ top: document.body.scrollHeight, behavior: "instant" }));
  const marker = page.locator(`.hextra-nav-menu-items [data-episode-progress="${episodePath}"]`);
  await expect(marker).toHaveAttribute("aria-label", "Complete");

  await page.goto("/episodes/");
  await expect(page.getByRole("link", { name: "Continue where you left off" })).toHaveAttribute(
    "href",
    "/episodes/02-facilitating-activity/",
  );
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

test("lesson controls tolerate unavailable local storage", async ({ browser }, testInfo) => {
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
  await page.goto(`http://127.0.0.1:13131${longEpisodePath}`);
  await expect(page.locator(`.hextra-nav-menu-items [data-episode-progress="${longEpisodePath}"]`)).toHaveAttribute(
    "aria-label",
    /^[1-9]\d?% reached$/,
  );
  await page.goto("http://127.0.0.1:13131/episodes/");
  await expect(page.locator("[data-episode-resume]")).toBeHidden();
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
