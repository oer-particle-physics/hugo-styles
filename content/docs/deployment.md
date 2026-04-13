+++
title = "Deployment"
weight = 60
+++

`hugo-styles-template` is set up for GitHub Pages via GitHub Actions rather than a committed `gh-pages` branch.
For hosting patterns beyond this shared setup, see
[Hextra Deploy Site](https://imfing.github.io/hextra/docs/guide/deploy-site/).

## Recommended deployment flow

1. create the lesson repository from `hugo-styles-template`
2. set `baseURL` in `hugo.toml` to `https://<account>.github.io/<repo>/`
3. in GitHub, enable Pages and choose `GitHub Actions` as the source
4. push `main`
5. let the included workflow build and publish the site

The included workflow deploys on pushes to `main`.
Enable Pages before the first push so that initial deploy run already has a configured publishing target.

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

## Optional: versioned lesson archives

If you want a Hextra-style version switcher, keep the default single-version deployment only for the current site and add a custom build step for archived copies.
The navbar support is already present in `hugo-styles`; the missing piece is publishing extra builds under stable URLs such as `/versions/latest/` and `/versions/v1.0/`.

See [Versioned Sites]({{< relref "/docs/versioned-sites" >}}) for the pattern Hextra uses upstream and how to adapt it for lesson repositories.
