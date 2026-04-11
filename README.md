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

## What it provides

- Episode layouts with automatic rendering of `questions`, `objectives`, `keypoints`, and teaching/exercise time.
- Custom shortcodes for `challenge`, `solution`, `hint`, `learner`, `instructor`, glossary references, and learner profile references.
- Hextra-native tabs with synced variants enabled by default.
- Aggregated lesson pages for:
  - Key Points
  - All-in-One
  - Extract All Images
- A learner/instructor view toggle.
- Documentation for authors and maintainers.
- A Go-based `hugo-styles-migrate` command for migration checks and common conversions.

## Documentation map

The example site in this repository doubles as the public documentation for the module.

- `docs/quickstart`: first-run setup
- `docs/authoring`: lesson-writing model
- `docs/frontmatter`: episode metadata contract
- `docs/components`: shortcode and component reference
- `docs/glossary-profiles`: glossary/profile authoring
- `docs/deployment`: GitHub Pages workflow
- `docs/troubleshooting`: common failures and fixes
- `docs/migration`: legacy Carpentries migration flow
- `docs/updates`: downstream update and release workflow

## Update model

Downstream lessons should **not** copy layouts, assets, or shortcodes out of this repository. Instead they should import a released version of `hugo-styles` as a Hugo Module.

Typical downstream update flow:

```bash
hugo mod get -u github.com/oer-particle-physics/hugo-styles@latest
hugo mod tidy
hugo
```

For a smoother maintenance experience, downstream lesson repositories should enable Dependabot for `gomod` updates so module bumps arrive as pull requests.

## Local development

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

Before cutting a new `hugo-styles` release:

```bash
cd cmd/hugo-styles-migrate && go test ./...
cd ../..
npm run check:flexsearch
hugo --gc --minify
```

Then update `CHANGELOG.md`, create the release tag, and let downstream lesson repositories pick it up through `hugo mod get -u ...` or Dependabot PRs.
