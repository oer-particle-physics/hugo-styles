(() => {
  const root = document.documentElement;
  const viewStorageKey = "hugo-styles:view";
  const lessonSidebarStorageKey = "hugo-styles:lesson-sidebar-collapsed";
  const lessonTocStorageKey = "hugo-styles:lesson-toc-collapsed";
  const allowedViews = new Set(["learner", "instructor"]);
  const lessonShells = Array.from(document.querySelectorAll("[data-lesson-shell]"));
  const syncViewUrl = (view) => {
    const url = new URL(window.location.href);
    if (view === "learner") {
      url.searchParams.delete("view");
    } else {
      url.searchParams.set("view", view);
    }
    window.history.replaceState({}, "", `${url.pathname}${url.search}${url.hash}`);
  };

  const applyView = (view) => {
    const next = allowedViews.has(view) ? view : "learner";
    root.dataset.view = next;
    document.querySelectorAll("[data-view-toggle]").forEach((button) => {
      button.setAttribute("aria-pressed", String(button.dataset.viewToggle === next));
    });
  };

  const urlParams = new URLSearchParams(window.location.search);
  const fromUrl = urlParams.get("view");
  const fromStorage = window.localStorage.getItem(viewStorageKey);
  applyView(fromUrl || fromStorage || root.dataset.view || "learner");

  document.querySelectorAll("[data-view-toggle]").forEach((button) => {
    button.addEventListener("click", () => {
      const next = button.dataset.viewToggle || "learner";
      window.localStorage.setItem(viewStorageKey, next);
      applyView(next);
      syncViewUrl(next);
    });
  });

  document.querySelectorAll("[data-expand-solutions]").forEach((button) => {
    button.addEventListener("click", () => {
      const details = document.querySelectorAll("[data-lesson-disclosure='solution'], [data-lesson-disclosure='hint']");
      const shouldOpen = button.dataset.expanded !== "true";
      details.forEach((item) => {
        item.open = shouldOpen;
      });
      button.dataset.expanded = String(shouldOpen);
      button.textContent = shouldOpen ? "Collapse Hints and Solutions" : "Expand Hints and Solutions";
    });
  });

  const applyLessonToggleState = (kind, collapsed) => {
    const next = collapsed ? "true" : "false";
    const attribute = kind === "sidebar" ? "sidebarCollapsed" : "tocCollapsed";
    const buttonSelector = kind === "sidebar" ? "[data-lesson-sidebar-toggle]" : "[data-lesson-toc-toggle]";
    const collapsedLabel = kind === "sidebar" ? "Show Nav" : "Show TOC";
    const expandedLabel = kind === "sidebar" ? "Hide Nav" : "Hide TOC";
    const collapsedAria = kind === "sidebar" ? "Show lesson navigation" : "Show page table of contents";
    const expandedAria = kind === "sidebar" ? "Hide lesson navigation" : "Hide page table of contents";

    lessonShells.forEach((shell) => {
      shell.dataset[attribute] = next;
    });

    document.querySelectorAll(buttonSelector).forEach((button) => {
      button.setAttribute("aria-pressed", next);
      button.setAttribute("aria-label", collapsed ? collapsedAria : expandedAria);
      button.setAttribute("title", collapsed ? collapsedLabel : expandedLabel);
    });
  };

  document.querySelectorAll("[data-aio-search]").forEach((input) => {
    input.addEventListener("input", () => {
      const query = input.value.trim().toLowerCase();
      document.querySelectorAll("[data-aio-episode]").forEach((episode) => {
        const haystack = `${episode.dataset.title || ""} ${episode.textContent || ""}`.toLowerCase();
        episode.style.display = !query || haystack.includes(query) ? "" : "none";
      });
    });
  });

  if (lessonShells.length > 0) {
    const sidebarFromStorage = window.localStorage.getItem(lessonSidebarStorageKey);
    const tocFromStorage = window.localStorage.getItem(lessonTocStorageKey);
    applyLessonToggleState("sidebar", sidebarFromStorage === "true");
    applyLessonToggleState("toc", tocFromStorage === "true");

    document.querySelectorAll("[data-lesson-sidebar-toggle]").forEach((button) => {
      button.addEventListener("click", () => {
        const shell = button.closest("[data-lesson-shell]");
        const collapsed = shell?.dataset.sidebarCollapsed !== "true";
        window.localStorage.setItem(lessonSidebarStorageKey, String(collapsed));
        applyLessonToggleState("sidebar", collapsed);
      });
    });

    document.querySelectorAll("[data-lesson-toc-toggle]").forEach((button) => {
      button.addEventListener("click", () => {
        const shell = button.closest("[data-lesson-shell]");
        const collapsed = shell?.dataset.tocCollapsed !== "true";
        window.localStorage.setItem(lessonTocStorageKey, String(collapsed));
        applyLessonToggleState("toc", collapsed);
      });
    });
  }
})();
