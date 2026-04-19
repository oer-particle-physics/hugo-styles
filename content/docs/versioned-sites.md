+++
title = "Versioned Sites"
weight = 75
+++

`hugo-styles` and `hugo-styles-template` now ship with version-aware deployment support.
The default behavior is intentionally small:

- the default branch is published at the site root and labelled `Latest`
- plain `hugo server` still serves only the current checkout
- archived branches or tags are only added when you opt into them in `hugo.toml`

This keeps local authoring simple while still letting the GitHub Pages build publish multiple refs under stable URLs.

## The config model

Both repositories include this block by default:

```toml
[params.versioning]
  enable = true
  defaultBranch = "main"
  menuName = "Versions"
  menuIdentifier = "versions"
  menuWeight = 70

  [params.versioning.latest]
    enable = true
    label = "Latest"

  [params.versioning.branches]
    refs = []
    patterns = []
    all = false

  [params.versioning.tags]
    refs = []
    patterns = []
    all = false
```

## What this does by default

With the default config:

- the GitHub Pages workflow builds the configured default branch as the site root
- the navbar gets a `Versions` dropdown with one entry: `Latest`
- no extra branch or tag builds are published

The generated output uses these URL patterns:

- `https://<account>.github.io/<repo>/` for `Latest`
- `https://<account>.github.io/<repo>/versions/<name>/` for extra branches or tags

## Add specific branches or tags

To publish named refs explicitly:

```toml
[params.versioning.branches]
  refs = ["release-2026"]

[params.versioning.tags]
  refs = ["v1.0.0", "v1.1.0"]
```

These appear in the dropdown as:

- `Latest`
- `release-2026`
- `v1.0.0`
- `v1.1.0`

If `Latest` already points at the default branch, that branch is not duplicated under its raw branch name.

## Match tags or branches with wildcards

Pattern matching uses shell-style globs:

```toml
[params.versioning.tags]
  patterns = ["v*"]

[params.versioning.branches]
  patterns = ["release/*"]
```

Typical uses:

- `v*` for release tags
- `release/*` for maintenance branches
- `docs-*` for documentation branches

## Publish all branches or all tags

If you want broad discovery instead of an explicit list:

```toml
[params.versioning.branches]
  all = true

[params.versioning.tags]
  all = true
```

You can combine `all = true` with `patterns` or `refs`, but in most cases one approach is enough.

## Disable `Latest`

If you only want named version entries in the dropdown:

```toml
[params.versioning.latest]
  enable = false
```

The root site still has to exist for GitHub Pages, so the default branch continues to build at the site root.
This option only hides the `Latest` menu entry.

## Local development

For day-to-day authoring, keep using:

```bash
hugo server
```

That serves only the current checkout.
It does not try to build other tags or branches locally.

If you want the full deployment artifact on your machine, run the same helper the GitHub Pages workflow uses:

```bash
python3 scripts/build-versioned-site.py
```

You can override the production base URL while testing:

```bash
python3 scripts/build-versioned-site.py --base-url http://localhost:8000/
```

Then serve `public/` with any static file server.

For local rendered-site link checking, build a validation artifact with a root base URL and run `lychee`
against the generated HTML:

```bash
python3 scripts/build-versioned-site.py --base-url / --destination .cache/linkcheck-site --no-minify
lychee --cache --config lychee.toml --no-progress --root-dir .cache/linkcheck-site '.cache/linkcheck-site/**/*.html'
```

Using `--base-url /` keeps the generated version menu and internal links portable on disk, which makes local
link checking work even for GitHub Pages project sites that deploy under `/<repo>/`.

For repositories created from `hugo-styles-template`, that helper is committed in the lesson repository
and refreshed from the pinned `hugo-styles` module version by the template's update workflow and local
`./scripts/sync-build-versioned-site.sh` helper.

## How the build works

The helper script:

1. reads the effective Hugo config via `hugo config --format json`
2. resolves the configured branches and tags from Git
3. creates temporary Git worktrees instead of repeatedly checking out the current directory
4. generates a small per-build config file with the correct `Versions` menu entries
5. runs Hugo once for the root site and once for each extra version

This stays close to Hugo's built-in behavior:

- Hugo still handles routing, menus, `baseURL`, and output paths
- the custom logic is limited to ref discovery and menu generation

## Current limitation

The shared build helper supports one site root per repository via Hugo's normal `--source` option.
It does not currently support a different `source_dir` per ref like Hextra's own script does for its historical `docs` and `exampleSite` layouts.

That tradeoff is deliberate:

- it keeps the lesson workflow simpler
- it matches the current `hugo-styles` and template repository layouts
- it avoids adding extra per-version config unless a real downstream use case needs it

If you need to publish older refs from a different site subdirectory, extend the helper at that point instead of carrying that complexity by default.
