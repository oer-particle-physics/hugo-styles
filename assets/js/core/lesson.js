(() => {
  const root = document.documentElement;
  const viewStorageKey = "hugo-styles:view";
  const sidebarStorageKey = "hugo-styles:sidebar-collapsed";
  const allowedViews = new Set(["learner", "instructor"]);
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

  const applySidebarState = (collapsed) => {
    const next = collapsed ? "true" : "false";
    root.dataset.sidebarCollapsed = next;
    document.querySelectorAll("[data-sidebar-toggle]").forEach((button) => {
      button.setAttribute("aria-pressed", next);
    });
    document.querySelectorAll("[data-sidebar-toggle-label]").forEach((label) => {
      label.textContent = collapsed ? "Show Nav" : "Hide Nav";
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

  const sidebarFromStorage = window.localStorage.getItem(sidebarStorageKey);
  applySidebarState(sidebarFromStorage === "true");

  document.querySelectorAll("[data-sidebar-toggle]").forEach((button) => {
    button.addEventListener("click", () => {
      const collapsed = root.dataset.sidebarCollapsed !== "true";
      window.localStorage.setItem(sidebarStorageKey, String(collapsed));
      applySidebarState(collapsed);
    });
  });
})();
