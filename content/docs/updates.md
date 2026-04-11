+++
title = "Updating Downstream Lessons"
weight = 100
+++

The update path is built around the shared module, not around repeatedly syncing a template repository.

## Recommended downstream setup

- keep lesson-specific content in the lesson repository
- import `github.com/oer-particle-physics/hugo-styles` as a Hugo Module
- override only what is genuinely repo-specific
- let Dependabot open module-update PRs

## Manual update flow

```bash
hugo mod get -u github.com/oer-particle-physics/hugo-styles@latest
hugo mod tidy
hugo mod graph
hugo
```

Review the rendered preview before merging the module bump, especially if the changelog mentions a breaking change.

## Search bundle maintenance

The shared module vendors the FlexSearch runtime locally so lesson builds stay self-contained.
That bundle is maintained in the `hugo-styles` repository rather than in downstream lessons.

When Dependabot opens an npm update for `flexsearch`, refresh the committed bundle in the module repo with:

```bash
npm ci
npm run vendor:flexsearch
```

CI also runs `npm run check:flexsearch` to catch version bumps where the vendored bundle was not refreshed yet.

## Override strategy

Hugo's normal precedence rules let downstream lessons override the module safely:

- local `layouts/` override module layouts
- local `assets/` override module assets
- local `archetypes/` override module archetypes
- local `hugo.toml` values override module defaults

Use that for branding, navigation changes, or lesson-specific extras without forking the shared infrastructure.

## Maintainer release checklist

When updating the shared module itself:

1. run `go test ./...` in `cmd/hugo-styles-migrate`
2. run `npm run check:flexsearch`
3. run `hugo --gc --minify`
4. update `CHANGELOG.md`
5. tag and publish a release from the `hugo-styles` repository
