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

Before tagging a new `hugo-styles` release:

1. run `go test ./...` in `cmd/hugo-styles-migrate`
2. run `npm run check:flexsearch`
3. run `hugo --gc --minify`
4. update `CHANGELOG.md`
5. tag and publish a release from `hugo-styles`

After release, downstream lesson repositories consume updates via:

- `hugo-styles-template` vendored refresh PRs (`go.mod`, `go.sum`, `_vendor/`)
- direct-module Dependabot PRs or manual `hugo mod get -u ...`
