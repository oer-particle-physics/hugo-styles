+++
title = "Migration Guide"
weight = 90
+++

The migration helper is designed to cover the common Carpentries lesson path, not every custom Jekyll feature ever added to a lesson repository.

## Recommended migration sequence

1. Run the checker on the legacy repository.
2. Run the migrator into a clean output directory.
3. Run the checker again on the migrated Hugo lesson.
4. Build with Hugo and do a quick manual pass for layout, links, and custom branding.

```bash
go run github.com/oer-particle-physics/hugo-styles/cmd/hugo-styles-migrate@latest check .
go run github.com/oer-particle-physics/hugo-styles/cmd/hugo-styles-migrate@latest migrate --source ../legacy-lesson --dest /tmp/converted-lesson
go run github.com/oer-particle-physics/hugo-styles/cmd/hugo-styles-migrate@latest check /tmp/converted-lesson
```

## What it converts well

- `_episodes/*.md` into Hugo episode content
- common metadata fields such as `questions`, `objectives`, `keypoints`, `teaching`, and `exercises`
- legacy Carpentries exercise blocks
- common callout classes
- `{{site.baseurl}}` asset paths
- common `links.md` include cleanup
- common YouTube and Vimeo iframe embeds

## What usually needs a manual pass

- custom Jekyll includes or layouts
- repo-specific branding and footer logic
- unusual Liquid expressions
- site-specific navigation widgets
- custom glossary or reference structures
- uncommon iframe providers or hand-written embeds

That means migrations like `hsf-training-docker` and `hsf-training-cicd` should be mostly systematic, while more customised repositories such as `gitlab-cms` should expect a short manual clean-up phase.

## What the checker validates after migration

When you run `check` on a Hugo lesson, it validates:

- `title`, `weight`, `questions`, `objectives`, and `keypoints` on episodes
- unique episode weights
- glossary/profile shortcode targets
- image alt text
- unsupported leftover legacy syntax
