+++
title = "Deployment"
weight = 60
+++

`hugo-styles-template` is set up for GitHub Pages via GitHub Actions rather than a committed `gh-pages` branch.
For hosting patterns beyond this shared setup, see
[Hextra Deploy Site](https://imfing.github.io/hextra/docs/guide/deploy-site/).

## Recommended deployment flow

1. create the lesson repository from `hugo-styles-template`
2. push `main`
3. in GitHub, enable Pages with the Actions source
4. let the included workflow build and publish the site

## Local production build

```bash
hugo --gc --minify
```

## Typical checks before publishing

```bash
go run github.com/oer-particle-physics/hugo-styles/cmd/hugo-styles-migrate@latest check .
hugo --gc --minify
```

## Why this uses Actions

- it works for any static site generator Hugo can produce
- it avoids the old branch-publishing coupling
- it matches the shared-module update model more naturally

The template repo already includes the workflow and Dependabot setup required for a typical lesson repository.
