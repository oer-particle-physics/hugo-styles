+++
title = "Updating Downstream Lessons"
weight = 100
+++

This page is for lesson maintainers updating lesson repositories.
If you maintain the shared `hugo-styles` module itself, use
[hugo-styles Maintenance]({{< relref "/docs/hugo-styles-maintenance" >}}).

## v0.6 author information

Version v0.6.0 removed support for `AUTHORS`. The `lesson/authors` shortcode now
reads `CITATION.cff` only. A lesson with no citation file can lose its authors table
while its build still passes. Follow
[Move author information to CITATION.cff]({{< relref "/docs/upgrades/authors-citation" >}})
before merging an upgrade from an older version.

## v0.5 compatibility notes

The v0.5 safety changes are intentional:

- the versioned-site builder protects marked output, requires `--force` for an unmarked non-empty destination,
  and requires `--allow-external-destination` outside the site root
- migration writes only into a clean, recognized `hugo-styles-template` Git worktree
- `lesson/image` requires both `src` and a non-empty `alt`
- episode metadata now uses strict string, integer, and non-empty-list validation
- callout titles now render trusted inline Markdown; Goldmark typography and permitted inline HTML therefore apply,
  while block or media output falls back to the original literal source

Fix validation errors rather than bypassing them. Existing content URLs and shortcode names remain compatible.

## Update through the scheduled refresh workflow

For repositories created from `hugo-styles-template`, the intended update flow is:

- keep `_vendor/` committed so lesson authors can build with Hugo Extended only
- use the **Update hugo-styles** GitHub Actions workflow
- review and merge the PR when it updates `go.mod`, `go.sum`, the managed workflows and scripts, `lychee.toml`, and `_vendor/`
- keep lesson-specific overrides in the lesson repository (`content/`, config, and selected overrides)

This avoids requiring local Go for normal lesson authoring.

### What the upgrade PR tells you

The title identifies the installed and target hugo-styles versions. The description
contains release and comparison links, upgrade notes for every release crossed,
the actual changed files, and the outcome of the configured build verification.
Notes are read from the exact target module version, together with its changelog.
If the hugo-styles version is unchanged, the title identifies a refresh at that
version; dependency or managed-file updates may still be included.

Resolve required actions and inspect the rendered lesson before merging. A passing
build does not prove that content or contributor information was preserved. Checks
for detectable unfinished migrations run again on later updates, even if the
release that introduced them has already been installed. The author check looks
for `AUTHORS` and `lesson/authors` in standard `content/` index files; custom content
locations and layout overrides need manual review. It does not compare contributors
or validate the CFF schema.

If verification fails, the updater stops before creating a PR, but the Actions job
summary retains the upgrade guidance and failed verification status. Non-stable
versions, missing notes in older modules, and unavailable changelog ranges are
identified explicitly rather than treated as evidence that no action is needed.

Custom reusable-workflow callers may still set `pr-title`. Their `pr-body` is now
additional context appended to the generated report. Keep the managed caller's
defaults to receive automatic version-aware titles. The workflow filenames and the
`chore/refresh-vendored-hugo-modules` PR branch are unchanged.

### Initial rollout of upgrade reports

An existing lesson runs its committed workflow. Updating workflow files during a
run cannot change the workflow already executing. The first PR that installs the
new updater may therefore still be titled **chore: refresh vendored Hugo modules**
and have the old static description. Its workflow may still appear as **Refresh
vendored Hugo modules** in Actions until that PR is merged.

Review the target release notes and upgrade guides manually for this first PR,
especially the v0.6 author migration. If the upgrade needs lesson edits, add them to
that PR and validate the result before merging. The next update run uses the new
reporting automatically. Re-running the old workflow before merging does not enable
the new reporting. A maintainer can also generate a report locally using the
[downstream report command]({{< relref "/docs/hugo-styles-maintenance" >}}).

### Configure `WORKFLOW_SYNC_TOKEN` once per lesson repository

For the automated refresh workflow to work reliably, add a repository Actions secret named `WORKFLOW_SYNC_TOKEN`.
The managed refresh workflow will use it automatically when it is present.

Use either:

- a fine-grained personal access token scoped to the lesson repository with `Contents: Read and write`, `Pull requests: Read and write`, and `Workflows: Read and write`
- or a GitHub App installation token with the same repository permissions

The refresh workflow updates managed files under `.github/workflows/` when `hugo-styles` ships workflow changes.
GitHub's default `GITHUB_TOKEN` can open pull requests, but GitHub rejects pushes that modify workflow files unless the token also has workflow write permission.

If `WORKFLOW_SYNC_TOKEN` is missing and no managed workflow files changed, the refresh emits a warning and continues with
the default `GITHUB_TOKEN`. If the refresh changed files under `.github/workflows/`, it stops before the pull-request step
with an actionable error explaining how to configure the secret.

### What `_vendor/` is for

`_vendor/` is a committed snapshot of Hugo module dependencies pinned by `go.mod` and `go.sum`.
Committing it keeps lesson builds reproducible and lets authors run `hugo server` without local Go.

The sync helper copies the managed maintainer files from the exact pinned `hugo-styles`
module version. That currently includes:

- `scripts/build-versioned-site.py`
- `scripts/upgrade-report.py`
- `scripts/sync-template-files.sh`
- `lychee.toml`
- `.github/workflows/cffconvert.yml`
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
