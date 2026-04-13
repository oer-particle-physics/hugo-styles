+++
title = "Versioned Sites"
weight = 75
+++

Hextra's public docs site does not use a special built-in versioning feature.
It combines two ordinary Hugo features:

1. a navbar dropdown created from nested `menus.main` entries
2. a deployment step that builds multiple copies of the site under stable subpaths

Upstream references:

- [Hextra docs menu config](https://github.com/imfing/hextra/blob/main/docs/hugo.yaml)
- [Hextra multi-version build script](https://github.com/imfing/hextra/blob/main/build.sh)
- [Hextra GitHub Pages workflow](https://github.com/imfing/hextra/blob/main/.github/workflows/pages.yml)

## What Hextra does

Hextra's docs config defines a top-level `Versions` menu item with children that point at published URLs such as:

- `/versions/latest/`
- `/versions/v0.11/`
- `/versions/v0.10/`

The theme renders that as a dropdown because child menu entries are already supported in the navbar.
There is no automatic tag discovery and no theme-specific version registry.
The links are configured manually in site config.

## Can this use a branch or a tag?

Yes.
Hextra's `build.sh` stores each version as a Git ref and then runs `git checkout <ref>` before building that copy of the site.
That means the ref can be:

- a branch such as `main`
- a tag such as `v0.11.3`

Upstream uses `main` for the `latest` docs build and release tags for archived versions.

## What this means for `hugo-styles`

`hugo-styles` already preserves nested navbar menus, so no layout or module change is required to support a version dropdown.
To use the pattern in a lesson repository, you need to do two things:

1. publish separate builds at predictable URLs
2. add matching `Versions` links in `hugo.toml`

If you do not need archived lesson releases, keep the default single-version deployment from `hugo-styles-template`.

## Minimal menu example

```toml
[[menus.main]]
  name = "Versions"
  identifier = "versions"
  weight = 70

[[menus.main]]
  name = "Development ↗"
  parent = "versions"
  url = "https://example.github.io/example-lesson/versions/latest/"
  weight = 10

[[menus.main]]
  name = "v1.0 ↗"
  parent = "versions"
  url = "https://example.github.io/example-lesson/versions/v1.0/"
  weight = 20
```

Use absolute URLs here because each entry typically points at a separately published site copy.

## Minimal build idea

Hextra's upstream script builds each version into `public/versions/<name>/`.
The same idea works for lesson sites:

```bash
VERSIONS=(
  "main:latest"
  "v1.0.0:v1.0"
)

for VERSION in "${VERSIONS[@]}"; do
  IFS=':' read -r REF NAME <<< "$VERSION"
  git checkout "$REF"
  hugo --baseURL "$BASE_URL/versions/$NAME/" --destination "../public/versions/$NAME"
done
```

In practice, you would replace the template's single `hugo --gc --minify` build step with a small wrapper script like this in GitHub Actions.
If you want a cleaner local workflow than repeated `git checkout`, prefer separate worktrees for each ref.

## Tradeoffs

- Good fit for documentation that needs archived releases.
- More CI complexity and longer build times than a single Pages deployment.
- Menu entries stay manual, so adding a new release means updating both the build script and the menu links.
