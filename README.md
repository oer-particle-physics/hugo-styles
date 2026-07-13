# Hugo Styles

`hugo-styles` is the reusable core for a
[Hugo](https://gohugo.io/)-based
lesson stack that preserves the
[Carpentries](https://carpentries.org/)
pedagogy layer of
[styles](https://github.com/carpentries/styles)
and
[workbench-template-md](https://github.com/carpentries/workbench-template-md)
while keeping the authoring and deployment flow of Hugo.

This repository serves two roles:

1. A versioned Hugo Module that downstream lessons import.
2. A self-documenting example site
   that demonstrates the supported lesson features.

## Start here

- Creating a new lesson repository? Start with
  [`hugo-styles-template`](https://github.com/oer-particle-physics/hugo-styles-template).
- Looking for the shared module docs? Start with the published
  [Quickstart](https://oer-particle-physics.github.io/hugo-styles/docs/quickstart/),
  then the
  [Authoring Guide](https://oer-particle-physics.github.io/hugo-styles/docs/authoring/).
- Need general theme or framework background? Use the upstream
  [Hextra docs](https://imfing.github.io/hextra/docs/getting-started/)
  and
  [Hugo docs](https://gohugo.io/installation/).

## What it provides

- Episode layouts with automatic rendering of `questions`, `objectives`, `keypoints`, and teaching/exercise time.
- Custom shortcodes for `challenge`, `solution`, `hint`, `learner`, `instructor`, glossary references, learner profile references, and lesson homepage snippets for overview, schedule, and authors from `AUTHORS`.
- Hextra-native tabs with synced variants enabled by default.
- Aggregated lesson pages for:
  - Key Points
  - All-in-One
  - External Links
  - All Images
- A learner/instructor view toggle.
- Documentation for authors and maintainers.
- A Go-based `hugo-styles-migrate` command for migration checks and common conversions.

## Documentation map

The example site in this repository doubles as the public documentation for the module.

- [Quickstart](https://oer-particle-physics.github.io/hugo-styles/docs/quickstart/): first-run setup
- [Authoring Guide](https://oer-particle-physics.github.io/hugo-styles/docs/authoring/): lesson-writing model
- [Front Matter](https://oer-particle-physics.github.io/hugo-styles/docs/frontmatter/): episode metadata contract
- [Components](https://oer-particle-physics.github.io/hugo-styles/docs/components/): shortcode and component reference
- [Hextra Features](https://oer-particle-physics.github.io/hugo-styles/docs/hextra-features/): live theme-feature examples and configuration
- [Glossary and Profiles](https://oer-particle-physics.github.io/hugo-styles/docs/glossary-profiles/): glossary/profile authoring
- [Deployment](https://oer-particle-physics.github.io/hugo-styles/docs/deployment/): GitHub Pages workflow
- [Versioned Sites](https://oer-particle-physics.github.io/hugo-styles/docs/versioned-sites/): archived refs and safe build output
- [Troubleshooting](https://oer-particle-physics.github.io/hugo-styles/docs/troubleshooting/): common failures and fixes
- [Migration Guide](https://oer-particle-physics.github.io/hugo-styles/docs/migration/): legacy Carpentries migration flow
- [Update Guide](https://oer-particle-physics.github.io/hugo-styles/docs/updates/): downstream update and release workflow
- [hugo-styles Maintenance](https://oer-particle-physics.github.io/hugo-styles/docs/hugo-styles-maintenance/): shared-module test and release checklist
- [Reference](https://oer-particle-physics.github.io/hugo-styles/reference/): further reading for Hextra and Hugo

## Update model

Downstream lessons should **not** copy layouts, assets, or shortcodes out of this repository. Instead they should import a released version of `hugo-styles` as a Hugo Module.

The `hugo-styles-template` repository commits `_vendor/` so lesson authors can run local builds with Hugo Extended alone.
Lesson repositories receive released module updates through the scheduled **Refresh vendored Hugo modules** workflow.
It updates the pinned module, refreshes `_vendor/` and the managed maintainer files, and opens a pull request for review.

The shared sync currently manages:

- `scripts/build-versioned-site.py`
- `scripts/sync-template-files.sh`
- `lychee.toml`
- `.github/workflows/pages.yml`
- `.github/workflows/refresh-vendored-modules.yml`
- `.github/workflows/reusable-pages.yml`
- `.github/workflows/reusable-refresh-vendored-modules.yml`

## Local development

For downstream lesson authors, the practical prerequisites are:

- [Hugo Extended](https://gohugo.io/installation/)
- [Go](https://go.dev/doc/install) (optional for template-based authoring with committed `_vendor/`; required for module maintenance and migration checks)
- [lychee](https://lychee.cli.rs/guides/getting-started/) (optional for local rendered-site link checks)

Node.js is only needed in this repository when maintainers refresh vendored frontend assets or run browser tests.

```bash
hugo server
```

## Validation and tests

The shared checker can validate both legacy Carpentries lessons and Hugo-native lesson repositories:

```bash
(cd cmd/hugo-styles-migrate && go run . check ../..)
```

The Hugo-native checks currently cover:

- strict required episode metadata types and non-empty values
- duplicate episode weights
- unresolved glossary references
- unresolved profile references
- missing image alt text
- leftover unsupported legacy syntax

Regression tests for the checker and migrator live under `cmd/hugo-styles-migrate/testdata/`.

```bash
(cd cmd/hugo-styles-migrate && go test ./... && go vet ./...)
python3 -m unittest discover -s scripts/tests -v
npm run test:browser
```

Rendered-site link checks use `lychee` against a local build that mirrors the GitHub Actions workflow:

```bash
python3 scripts/build-versioned-site.py --use-current-checkout \
  --base-url / --destination .cache/linkcheck-site --no-minify
lychee --cache --config lychee.toml --no-progress --root-dir .cache/linkcheck-site '.cache/linkcheck-site/**/*.html'
```

The workflow uses `--base-url /` for this validation build so internal links are checked against local files
instead of a future GitHub Pages URL.
Keep ordinary broken links as content fixes, but if a trusted upstream URL is known to fail automated validation
for reasons such as an invalid TLS chain, add a narrow exact-match entry to the downstream repository's
`.lycheeignore` instead of enabling global insecure mode. The shared template can include a comment-only placeholder,
but the actual ignore entries are repository-specific.

## Vendored frontend asset maintenance

Search uses a vendored FlexSearch bundle, and image zoom uses a vendored Medium Zoom bundle,
so local and GitHub Pages builds do not depend on a CDN.
Most lesson authors never need Node.js for this repository, but maintainers do need it when refreshing
vendored frontend assets after a Dependabot bump.

```bash
npm ci
npm run vendor:flexsearch
npm run vendor:medium-zoom
```

The `npm run check:flexsearch` and `npm run check:medium-zoom` commands are used in CI to
confirm the committed bundles still match the pinned package versions.

## Migration tool

Run the checker or migration helper directly from this repository:

```bash
cd cmd/hugo-styles-migrate
go run . check ../..
go run . migrate --source ../old-training --dest ../clean-template-training --dry-run
go run . migrate --source ../old-training --dest ../clean-template-training
go run . check ../clean-template-training
cd ../..
```

The destination must be a clean Git worktree created from `hugo-styles-template`; migration preserves its
configuration, workflows, branding, generated-resource pages, and section indexes.

Or from another repository:

```bash
go run github.com/oer-particle-physics/hugo-styles/cmd/hugo-styles-migrate@latest check .
```

## Release workflow

`hugo-styles` now uses `release-please` for release PRs and changelog updates.
The normal maintainer flow is:

```bash
prek install --hook-type commit-msg
```

Then work with conventional commits locally. CI also runs `cz check --rev-range ...` on pull requests, so the hook is a faster local guardrail rather than the only enforcement point.

Before merging a release PR or sanity-checking a release candidate:

```bash
(cd cmd/hugo-styles-migrate && go test ./... && go vet ./...)
python3 -m unittest discover -s scripts/tests -v
npm run check:flexsearch
npm run check:medium-zoom
npm run test:browser
hugo --gc --minify --panicOnWarning
python3 scripts/build-versioned-site.py --use-current-checkout \
  --base-url / --destination .cache/linkcheck-site --no-minify
lychee --cache --config lychee.toml --no-progress --root-dir .cache/linkcheck-site '.cache/linkcheck-site/**/*.html'
```

The `release-please` workflow expects a `RELEASE_PLEASE_TOKEN` secret so the generated release PRs and tags can trigger follow-up GitHub Actions runs normally.
Once the release PR is merged, downstream lesson repositories pick up the new version through the **Refresh vendored Hugo modules** workflow.
