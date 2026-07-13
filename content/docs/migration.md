+++
title = "Migration Guide"
weight = 90
+++

The migrator moves lesson content into an existing
[hugo-styles-template](https://github.com/oer-particle-physics/hugo-styles-template) repository. It deliberately
keeps the template's configuration and automation instead of creating a second, incomplete site.

## Before you start

The destination must:

- be the root of a Git worktree with no tracked or untracked changes
- import and require `github.com/oer-particle-physics/hugo-styles` in `hugo.toml` and `go.mod`
- be separate from the legacy source tree

Commit or stash destination changes first. There is no force option for bypassing these checks.

## Preview, then migrate

```bash
go run github.com/oer-particle-physics/hugo-styles/cmd/hugo-styles-migrate@latest check ../legacy-lesson
go run github.com/oer-particle-physics/hugo-styles/cmd/hugo-styles-migrate@latest migrate \
  --source ../legacy-lesson \
  --dest ../new-template-lesson \
  --dry-run
go run github.com/oer-particle-physics/hugo-styles/cmd/hugo-styles-migrate@latest migrate \
  --source ../legacy-lesson \
  --dest ../new-template-lesson
```

The dry run prints every absolute path it would add, replace, remove, or preserve. The real migration stages all
changes first and restores the original destination if applying them fails.

## Replaced content

- the homepage and regular content below `episodes`, `learners`, `instructors`, `glossary`, and `profiles`
- optional `content/reference.md` and root `AUTHORS`
- optional `static/fig`, `static/files`, `static/data`, and `static/code`

Section `_index.md` files are preserved. So are template configuration, branding, workflows, generated-resource
pages, and all other infrastructure outside those managed paths.

## What it converts

The common path covers legacy episodes and metadata, Carpentries exercise/callout blocks,
`{{site.baseurl}}` asset paths, common link includes, and YouTube/Vimeo embeds. Custom Jekyll includes,
unusual Liquid, repository-specific layouts, branding, and uncommon embeds still need review.

## Finish the migration

```bash
cd ../new-template-lesson
go run github.com/oer-particle-physics/hugo-styles/cmd/hugo-styles-migrate@latest check .
hugo --gc --minify --panicOnWarning
```

Then edit lesson metadata and repository URLs in `hugo.toml`, preview both audience views, and review any
custom legacy constructs reported by the checker.
