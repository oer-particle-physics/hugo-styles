+++
title = "Updating Downstream Lessons"
weight = 100
+++

This page is for lesson maintainers updating lesson repositories.
If you maintain the shared `hugo-styles` module itself, use
[hugo-styles Maintenance]({{< relref "/docs/hugo-styles-maintenance" >}}).

## v0.5 compatibility notes

The v0.5 safety changes are intentional:

- the versioned-site builder protects marked output, requires `--force` for an unmarked non-empty destination,
  and requires `--allow-external-destination` outside the site root
- migration writes only into a clean, recognized `hugo-styles-template` Git worktree
- `lesson/image` requires both `src` and a non-empty `alt`
- episode metadata now uses strict string, integer, and non-empty-list validation

Fix validation errors rather than bypassing them. Existing content URLs and shortcode names remain compatible.

## Update through the scheduled refresh workflow

For repositories created from `hugo-styles-template`, the intended update flow is:

- keep `_vendor/` committed so lesson authors can build with Hugo Extended only
- use the **Refresh vendored Hugo modules** GitHub Actions workflow
- review and merge the PR when it updates `go.mod`, `go.sum`, the managed workflow files, `scripts/build-versioned-site.py`, `scripts/sync-template-files.sh`, `lychee.toml`, and `_vendor/`
- keep lesson-specific overrides in the lesson repository (`content/`, config, and selected overrides)

This avoids requiring local Go for normal lesson authoring.

### Configure `WORKFLOW_SYNC_TOKEN` once per lesson repository

For the automated refresh workflow to work reliably, add a repository Actions secret named `WORKFLOW_SYNC_TOKEN`.
The managed refresh workflow will use it automatically when it is present.

Use either:

- a fine-grained personal access token scoped to the lesson repository with `Contents: Read and write`, `Pull requests: Read and write`, and `Workflows: Read and write`
- or a GitHub App installation token with the same repository permissions

The refresh workflow updates managed files under `.github/workflows/` when `hugo-styles` ships workflow changes.
GitHub's default `GITHUB_TOKEN` can open pull requests, but GitHub rejects pushes that modify workflow files unless the token also has workflow write permission.

If `WORKFLOW_SYNC_TOKEN` is missing, the refresh workflow may still succeed for releases that do not touch managed workflow files,
but it can fail during the pull-request step once an upstream release updates `.github/workflows/*`.

### What `_vendor/` is for

`_vendor/` is a committed snapshot of Hugo module dependencies pinned by `go.mod` and `go.sum`.
Committing it keeps lesson builds reproducible and lets authors run `hugo server` without local Go.

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
Review the rendered preview and workflow checks before merging the refresh pull request, especially when the changelog mentions a breaking change.

## Override strategy

Hugo's normal precedence rules let downstream lessons override the module safely:

- local `layouts/` override module layouts
- local `assets/` override module assets
- local `archetypes/` override module archetypes
- local `hugo.toml` values override module defaults

Use that for branding, navigation changes, or lesson-specific extras without forking the shared infrastructure.
