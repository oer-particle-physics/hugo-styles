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

Before merging a release PR or sanity-checking a release candidate:

1. run `go test ./...` and `go vet ./...` in `cmd/hugo-styles-migrate`
2. run `python3 -m unittest discover -s scripts/tests -v`
3. run `npm ci`, both `npm run check:*` asset checks, and `npm run test:browser`
4. run `hugo --gc --minify --panicOnWarning`
5. build the current checkout and configured archived refs with `--use-current-checkout`
6. run `lychee` against that versioned output
7. confirm the `RELEASE_PLEASE_TOKEN` secret is available, then merge the release PR

```bash
(cd cmd/hugo-styles-migrate && go test ./... && go vet ./...)
python3 -m unittest discover -s scripts/tests -v
npm ci
npm run check:flexsearch
npm run check:medium-zoom
npm run test:browser
hugo --gc --minify --panicOnWarning
python3 scripts/build-versioned-site.py --use-current-checkout \
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

After release, downstream lesson repositories consume updates via:

- `hugo-styles-template` vendored refresh PRs (`go.mod`, `go.sum`, managed files, `_vendor/`)
- direct-module Dependabot PRs or manual `hugo mod get -u ...`
