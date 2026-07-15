+++
title = "Versioned Sites"
weight = 75
+++

The versioned builder publishes the configured default branch as `Latest` and optional branches or tags below
`versions/<name>/`. This documentation site demonstrates the result with the archived `v0.4.0` tag.

## Select published refs

```toml
[params.versioning]
  enable = true
  defaultBranch = "main"
  menuName = "Versions"

  [params.versioning.latest]
    enable = true
    label = "Latest"

  [params.versioning.branches]
    refs = []
    patterns = []
    all = false

  [params.versioning.tags]
    refs = ["v0.4.0"]
    patterns = []
    all = false
```

Keep `defaultBranch` aligned with the repository's default branch in GitHub. The deployment workflow discovers the
GitHub setting at runtime, while local versioned builds use this Hugo configuration value.

Use `refs` for an explicit list, shell-style `patterns = ["v*"]` or `["release/*"]` for matching refs, or
`all = true` for every local ref of that kind. A missing explicitly named ref fails the build.

Published URLs are:

- `https://<account>.github.io/<repo>/` for Latest
- `https://<account>.github.io/<repo>/versions/<name>/` for a branch or tag

Each ref is built with its own Hugo configuration and content. The builder only injects the common version menu
and that target's `baseURL`, so historical theme and feature settings remain accurate.

## Local and pull-request builds

Normal authoring still uses the current tree only:

```bash
hugo server
```

To validate the current checkout while also building configured archived refs:

```bash
python3 scripts/build-versioned-site.py \
  --use-current-checkout \
  --base-url / \
  --destination .cache/versioned-site \
  --no-minify
```

`--use-current-checkout` matters in CI and detached pull-request checkouts: it builds the checked-out changes
as Latest instead of silently rebuilding `main`.

## Destination safety

The builder only replaces an empty directory or output containing its
`.hugo-styles-versioned-output` managed marker. For a non-empty directory created by another tool, review the
path and add `--force`. Structural protections still reject the source,
repository root, filesystem root, their ancestors, and destinations that escape through symlinks.

Destinations outside the site root require a separate acknowledgement:

```bash
python3 scripts/build-versioned-site.py \
  --destination /tmp/hugo-styles-preview \
  --allow-external-destination
```

`--force` never bypasses these structural or external-path checks. Every target is built in staging; the
destination is swapped only after all targets succeed. A failed build leaves the previous published output intact.

## Rendered-link validation

```bash
python3 scripts/build-versioned-site.py \
  --use-current-checkout \
  --base-url / \
  --destination .cache/linkcheck-site \
  --no-minify
lychee --cache --config lychee.toml --no-progress \
  --root-dir .cache/linkcheck-site '.cache/linkcheck-site/**/*.html'
```

For template-based lessons, `scripts/build-versioned-site.py` and its workflows are managed files refreshed from
the exact `hugo-styles` version pinned in `go.mod`.
