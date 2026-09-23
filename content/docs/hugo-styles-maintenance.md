+++
title = "hugo-styles Maintenance"
weight = 110
+++

This page is for maintainers of the shared `hugo-styles` repository.
Lesson maintainers should start with [Updating Downstream Lessons]({{< relref "/docs/updates" >}}).

## Vendored frontend asset maintenance

Search uses a vendored FlexSearch bundle, and image zoom uses a vendored Medium Zoom bundle,
so lesson builds do not depend on a CDN.
When Dependabot bumps `flexsearch` or `medium-zoom` in `hugo-styles`, refresh the committed bundles:

```bash
npm ci
npm run vendor:flexsearch
npm run vendor:medium-zoom
```

CI runs `npm run check:flexsearch` and `npm run check:medium-zoom` to ensure the committed
bundles match the pinned package versions.

## Shared module release checklist

`hugo-styles` uses `release-please` for release PRs, generated changelog updates, and GitHub releases.
The expected local setup is:

```bash
prek install --hook-type commit-msg
```

That installs the standard `.pre-commit-config.yaml` hook using `prek`, which is the preferred local runner because it is faster than `pre-commit`.
CI remains authoritative and runs `cz check --rev-range ...` on pull requests.

### Communicate upgrade impact in the change PR

Every change PR should state whether lesson maintainers need to take action. Normal
fixes can use the generated changelog. Changes to supported content, configuration,
metadata, workflows, or rendered behavior that need attention require an upgrade
note in `content/docs/upgrades/`. Add the note in the same PR as the change.

Create a note using the shared archetype (run from this repository):

```bash
hugo new content --kind upgrade-note docs/upgrades/rename-setting.md
```

Choose a descriptive, stable filename. The archetype supplies this metadata and
the four required sections; replace every `TODO:` instruction with guidance:

```toml
+++
title = "Describe the action for lesson maintainers"

[params.upgrade]
release = "unreleased"
severity = "breaking" # breaking, deprecation, or action
# check = "legacy-authors" # optional; only registered checks are accepted
+++
```

Write four sections: `## Who is affected`, `## What changes`, `## What to do`, and
`## How to verify`. Explain the consequence of doing nothing, provide concrete file
edits or examples, and describe how to verify the result. Use plain Markdown and
absolute links: the same text becomes a documentation page
and appears in downstream upgrade PRs. The upgrade-notes index lists the pages
automatically.

Executable Hugo shortcodes are not supported in notes. To show literal shortcode
examples, use Hugo's escape comments in the source, including inside code fences:

```markdown
{{</*/* lesson/authors */*/>}}
```

The docs renderer and report both display the literal call without the escape
comments. Ordinary code examples, including Markdown headings, are preserved.

Keep `release = "unreleased"` during development. Do not guess the next version or
set Hugo's `draft` field on the page. Preview the completed guidance without a
lesson checkout or version arguments:

```bash
python3 scripts/upgrade-report.py preview
```

Use `--output /tmp/upgrade-preview.md` to save the Markdown, `--site-root /path/to/lesson`
to exercise applicability checks, or `--release v0.6.0` to also include notes already
assigned to a particular version. Without a lesson path, the preview explicitly
leaves applicability unchecked. It never claims that an upgrade or build succeeded.

The **Check upgrade guidance** PR workflow validates notes and publishes the preview
in its job summary and the `upgrade-guidance-preview` artifact. It runs again when
the PR title or description changes. A conventional breaking-change marker in the
title, description, or a commit requires a new note with `severity = "breaking"`;
updating an old guide does not satisfy that requirement. Reviewers must still catch
changes whose impact was not declared. Make this check required in branch protection
to enforce it before merging.

An optional `check` names a function in `scripts/upgrade-report.py`. Add a check
only when repository files can identify a useful action reliably. Cover affected,
migrated, and intentionally absent configurations with tests. Without a check, a
note is shown to everyone crossing its release. With a check, detected unresolved
work is also shown on later upgrades. Document a check's limits; a passing build or
the absence of a finding must not be presented as proof of a completed migration.

For removals, normally ship a deprecation warning and retain the old behavior for
at least one published release before removing it. State the planned removal
version and migration steps. If a transition period is impractical, explain why in
the change PR and call out the immediate impact prominently. Mark breaking changes
in the conventional commit and review the version proposed by release-please,
including the intended policy while the project is below v1.0.

CI validates upgrade-note metadata, non-empty sections, scaffold placeholders, and
registered check names. Reviewers still need to assess whether the instructions are
correct and complete. Most changes need only a Markdown note; no detector or workflow
code changes are required.

### Review the release and its upgrade report

Before merging a release PR or sanity-checking a release candidate:

1. check out the release PR branch, fetch tags, and assign pending notes using its
   proposed manifest version:

   ```bash
   git fetch --tags origin
   python3 scripts/upgrade-report.py assign-release
   python3 scripts/upgrade-report.py check-notes --require-released
   ```

   Commit the changed notes to the release PR. The command edits only pending note
   versions, preserves published notes and their content, and refuses to assign an
   already tagged version. Re-running it after assignment makes no further edits.
   Release Please can regenerate its branch when more changes land on `main`,
   replacing manual commits. After a bot refresh, update your checkout and rerun
   assignment before merging; the release check catches restored pending notes.
2. review required actions alongside the generated changelog; include a link to the
   relevant upgrade guide in breaking-change release notes
3. run `go test ./...` and `go vet ./...` in `cmd/hugo-styles-migrate`
4. run `python3 -m unittest discover -s scripts/tests -v`
5. run `npm ci`, both `npm run check:*` asset checks, and `npm run test:browser`
6. run `hugo --config hugo.toml,hugo-docs.toml --gc --minify --panicOnWarning`
7. build the current checkout and configured archived refs with `--use-current-checkout`
8. run `lychee` against that versioned output
9. test the upgrade on a disposable lesson checkout and inspect its generated report,
   including an upgrade that skips a release and any unresolved migration findings
10. confirm the `RELEASE_PLEASE_TOKEN` secret is available, then merge the release PR

Release PRs that change `.release-please-manifest.json` must have no pending notes
and no notes assigned beyond the proposed version. Their CI preview includes notes
for that version as well as any pending guidance, so it remains useful when the
release check fails.

The release workflow also checks before Release Please can publish. If the manifest
version already has a Git tag, this is a development run: pending notes are allowed
and Release Please may update release PRs but cannot publish a release. For a new
manifest version, incomplete, unassigned, or future-dated notes stop publication.
Published notes stay in the repository so lessons skipping releases receive all
relevant instructions.

```bash
(cd cmd/hugo-styles-migrate && go test ./... && go vet ./...)
python3 -m unittest discover -s scripts/tests -v
npm ci
npm run check:flexsearch
npm run check:medium-zoom
npm run test:browser
hugo --config hugo.toml,hugo-docs.toml --gc --minify --panicOnWarning
python3 scripts/build-versioned-site.py --use-current-checkout \
  --config hugo.toml,hugo-docs.toml \
  --base-url / --destination .cache/versioned-site --no-minify
lychee --cache --config lychee.toml --no-progress \
  --root-dir .cache/versioned-site '.cache/versioned-site/**/*.html'
```

The release workflow creates the root `vX.Y.Z` tag and idempotently mirrors it as
`cmd/hugo-styles-migrate/vX.Y.Z`. After release, confirm both tags exist and that:

```bash
go list -m github.com/oer-particle-physics/hugo-styles/cmd/hugo-styles-migrate@latest
```

resolves to the semantic release rather than a pseudo-version.

To exercise the full downstream report, pass the lesson checkout and the target
module checkout explicitly (Python 3.11 or later is required):

```bash
python3 scripts/upgrade-report.py report \
  --site-root /path/to/disposable-lesson \
  --module-root . \
  --old-version v0.5.7 --new-version v0.6.0 \
  --output /tmp/hugo-styles-upgrade.md
```

Use the actual installed and proposed versions. Upgrade the disposable lesson
before generating the report so the file list describes the proposed changes.
The command does not run a build or modify the lesson. The workflow later records
whether its configured verification succeeded, failed, or was skipped; do not
describe a report preview as build validation.

After release, run **Update hugo-styles** in `hugo-styles-template`, inspect the PR,
and check the rendered lesson before merging. New lesson repositories then inherit
the current updater. Existing lessons receive managed files through their own update
PRs. Follow the [initial rollout guidance]({{< relref "/docs/updates#initial-rollout-of-upgrade-reports" >}})
for lessons whose installed workflow still uses the old static PR description.
