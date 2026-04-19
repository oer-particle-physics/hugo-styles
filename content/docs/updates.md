+++
title = "Updating Downstream Lessons"
weight = 100
+++

This page is for lesson maintainers updating lesson repositories.
If you maintain the shared `hugo-styles` module itself, use
[hugo-styles Maintenance]({{< relref "/docs/hugo-styles-maintenance" >}}).

## Recommended path: template + automated vendor refresh (no local Go)

For repositories created from `hugo-styles-template`, the intended update flow is:

- keep `_vendor/` committed so lesson authors can build with Hugo Extended only
- use the **Refresh vendored Hugo modules** GitHub Actions workflow
- review and merge the PR when it updates `go.mod`, `go.sum`, the managed workflow files, `scripts/build-versioned-site.py`, `scripts/sync-template-files.sh`, `lychee.toml`, and `_vendor/`
- keep lesson-specific overrides in the lesson repository (`content/`, config, and selected overrides)

This avoids requiring local Go for normal lesson authoring.

### What `_vendor/` is for

`_vendor/` is a committed snapshot of Hugo module dependencies pinned by `go.mod` and `go.sum`.
Committing it keeps lesson builds reproducible and lets authors run `hugo server` without local Go.

### Refresh locally (optional, if Go is available)

If you do have Go locally and want to refresh manually:

```bash
hugo mod get -u github.com/oer-particle-physics/hugo-styles@latest
hugo mod tidy
./scripts/sync-template-files.sh
hugo mod vendor
python3 scripts/build-versioned-site.py
```

The sync helper copies the managed maintainer files from the exact pinned `hugo-styles`
module version. That currently includes:

- `scripts/build-versioned-site.py`
- `scripts/sync-template-files.sh`
- `lychee.toml`
- `.github/workflows/pages.yml`
- `.github/workflows/refresh-vendored-modules.yml`
- `.github/workflows/reusable-pages.yml`
- `.github/workflows/reusable-refresh-vendored-modules.yml`

That keeps the committed maintainer files aligned with `go.mod` rather than downloading an unrelated head revision.

## Direct module mode (without `_vendor/`)

If a lesson repository imports `github.com/oer-particle-physics/hugo-styles` directly (without committed `_vendor/`),
enable Dependabot for `gomod` so updates arrive as pull requests.

Manual fallback:

```bash
hugo mod get -u github.com/oer-particle-physics/hugo-styles@latest
hugo mod tidy
hugo mod graph
hugo
```

Review the rendered preview before merging a module bump, especially when the changelog mentions a breaking change.

## Override strategy

Hugo's normal precedence rules let downstream lessons override the module safely:

- local `layouts/` override module layouts
- local `assets/` override module assets
- local `archetypes/` override module archetypes
- local `hugo.toml` values override module defaults

Use that for branding, navigation changes, or lesson-specific extras without forking the shared infrastructure.
