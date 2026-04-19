+++
title = "hugo-styles Maintenance"
weight = 110
+++

This page is for maintainers of the shared `hugo-styles` repository.
Lesson maintainers should start with [Updating Downstream Lessons]({{< relref "/docs/updates" >}}).

## Search bundle maintenance

Search uses a vendored FlexSearch bundle so lesson builds do not depend on a CDN.
When Dependabot bumps `flexsearch` in `hugo-styles`, refresh the committed bundle:

```bash
npm ci
npm run vendor:flexsearch
```

CI runs `npm run check:flexsearch` to ensure the committed bundle matches the pinned package version.

## Shared module release checklist

`hugo-styles` uses `release-please` for release PRs, generated changelog updates, and GitHub releases.
The expected local setup is:

```bash
prek install --hook-type commit-msg
```

That installs the standard `.pre-commit-config.yaml` hook using `prek`, which is the preferred local runner because it is faster than `pre-commit`.
CI remains authoritative and runs `cz check --rev-range ...` on pull requests.

Before merging a release PR or sanity-checking a release candidate:

1. run `go test ./...` in `cmd/hugo-styles-migrate`
2. run `npm run check:flexsearch`
3. run `python3 scripts/build-versioned-site.py --base-url / --destination .cache/linkcheck-site --no-minify`
4. run `lychee --cache --config lychee.toml --no-progress --root-dir .cache/linkcheck-site '.cache/linkcheck-site/**/*.html'`
5. run `hugo --gc --minify`
6. confirm the `RELEASE_PLEASE_TOKEN` secret is available to the workflow
7. merge the release PR that `release-please` opens

After release, downstream lesson repositories consume updates via:

- `hugo-styles-template` vendored refresh PRs (`go.mod`, `go.sum`, managed files, `_vendor/`)
- direct-module Dependabot PRs or manual `hugo mod get -u ...`
