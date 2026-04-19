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
  - Extract All Images
- A learner/instructor view toggle.
- Documentation for authors and maintainers.
- A Go-based `hugo-styles-migrate` command for migration checks and common conversions.

## Documentation map

The example site in this repository doubles as the public documentation for the module.

- [Quickstart](https://oer-particle-physics.github.io/hugo-styles/docs/quickstart/): first-run setup
- [Authoring Guide](https://oer-particle-physics.github.io/hugo-styles/docs/authoring/): lesson-writing model
- [Front Matter](https://oer-particle-physics.github.io/hugo-styles/docs/frontmatter/): episode metadata contract
- [Components](https://oer-particle-physics.github.io/hugo-styles/docs/components/): shortcode and component reference
- [Glossary and Profiles](https://oer-particle-physics.github.io/hugo-styles/docs/glossary-profiles/): glossary/profile authoring
- [Deployment](https://oer-particle-physics.github.io/hugo-styles/docs/deployment/): GitHub Pages workflow
- [Troubleshooting](https://oer-particle-physics.github.io/hugo-styles/docs/troubleshooting/): common failures and fixes
- [Migration Guide](https://oer-particle-physics.github.io/hugo-styles/docs/migration/): legacy Carpentries migration flow
- [Update Guide](https://oer-particle-physics.github.io/hugo-styles/docs/updates/): downstream update and release workflow
- [Reference](https://oer-particle-physics.github.io/hugo-styles/reference/): further reading for Hextra and Hugo

## Update model

Downstream lessons should **not** copy layouts, assets, or shortcodes out of this repository. Instead they should import a released version of `hugo-styles` as a Hugo Module.

Typical downstream update flow:

```bash
hugo mod get -u github.com/oer-particle-physics/hugo-styles@latest
hugo mod tidy
hugo
```

For a smoother maintenance experience, downstream lesson repositories should enable Dependabot for `gomod` updates so module bumps arrive as pull requests.

The `hugo-styles-template` repository commits `_vendor/` so lesson authors can run local builds with Hugo Extended alone.
Template maintainers still use Go when refreshing `go.mod`/`go.sum`, managed maintainer files, and `_vendor/`,
either through the **Refresh vendored Hugo modules** workflow or locally with:

```bash
hugo mod get -u github.com/oer-particle-physics/hugo-styles@latest
hugo mod tidy
./scripts/sync-template-files.sh
hugo mod vendor
```

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

Node.js is only needed in this repository when maintainers refresh the vendored search bundle.

```bash
hugo server
```

## Validation and tests

The shared checker can validate both legacy Carpentries lessons and Hugo-native lesson repositories:

```bash
(cd cmd/hugo-styles-migrate && go run . check ../..)
```

The Hugo-native checks currently cover:

- required episode metadata
- duplicate episode weights
- unresolved glossary references
- unresolved profile references
- missing image alt text
- leftover unsupported legacy syntax

Regression tests for the checker and migrator live under `cmd/hugo-styles-migrate/testdata/`.

```bash
(cd cmd/hugo-styles-migrate && go test ./...)
```

Rendered-site link checks use `lychee` against a local build that mirrors the GitHub Actions workflow:

```bash
python3 scripts/build-versioned-site.py --base-url / --destination .cache/linkcheck-site --no-minify
lychee --cache --config lychee.toml --no-progress --root-dir .cache/linkcheck-site '.cache/linkcheck-site/**/*.html'
```

The workflow uses `--base-url /` for this validation build so internal links are checked against local files
instead of a future GitHub Pages URL.

## Search bundle maintenance

Search uses a vendored FlexSearch bundle so local and GitHub Pages builds do not depend on a CDN.
Most lesson authors never need Node.js for this repository, but maintainers do need it when refreshing the
vendored search asset after a Dependabot bump.

```bash
npm ci
npm run vendor:flexsearch
```

The `npm run check:flexsearch` command is used in CI to confirm the committed bundle still matches the pinned
package version.

## Migration tool

Run the checker or migration helper directly from this repository:

```bash
cd cmd/hugo-styles-migrate
go run . check ../..
go run . migrate --source ../old-training --dest /tmp/converted-training
go run . check /tmp/converted-training
cd ../..
```

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
cd cmd/hugo-styles-migrate && go test ./...
cd ../..
npm run check:flexsearch
python3 scripts/build-versioned-site.py --base-url / --destination .cache/linkcheck-site --no-minify
lychee --cache --config lychee.toml --no-progress --root-dir .cache/linkcheck-site '.cache/linkcheck-site/**/*.html'
hugo --gc --minify
```

The `release-please` workflow expects a `RELEASE_PLEASE_TOKEN` secret so the generated release PRs and tags can trigger follow-up GitHub Actions runs normally.
Once the release PR is merged, downstream lesson repositories pick up the new version through `hugo mod get -u ...`, the **Refresh vendored Hugo modules** workflow, or Dependabot PRs.
