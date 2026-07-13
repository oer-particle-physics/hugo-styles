+++
title = "Authoring Guide"
weight = 20
+++

This is the shortest path from an empty template to a checked lesson. Keep the
[Front Matter]({{< relref "/docs/frontmatter" >}}) and
[Components]({{< relref "/docs/components" >}}) references open when you need field or shortcode details.

## 1. Create an episode

```bash
hugo new content --kind episode episodes/05-your-topic/index.md
```

Put the teaching sequence in `content/episodes/`. Use `content/learners/` for learner setup or reference material and
`content/instructors/` for facilitation notes. Keep glossary entries and audience profiles in their matching
`content/glossary/` and `content/profiles/` sections.

## 2. Complete the episode metadata

Every episode needs a non-empty string `title`, an integer `weight`, and non-empty lists of non-empty
`questions`, `objectives`, and `keypoints`. Optional `teaching` and `exercises` values are
non-negative integer minutes.

See [Front Matter]({{< relref "/docs/frontmatter" >}}) for a copyable example and TOML math guidance.

## 3. Add teaching components

Use project-owned components for pedagogy and audience-specific content, then Hextra components for general
documentation UI. The live references place copyable source beside every example:

- [Components]({{< relref "/docs/components" >}}): challenge, hint, solution, callouts, audience blocks, links, images, and lesson metadata
- [Hextra Features]({{< relref "/docs/hextra-features" >}}): math, Mermaid, tabs, details, steps, cards, code, PDF, Jupyter, search, theme, and other theme features
- [Glossary & Profiles]({{< relref "/docs/glossary-profiles" >}}): reusable glossary and profile pages

Episode and All-in-One pages add the learner/instructor selector automatically. On another page that contains
`learner` or `instructor` blocks, place `{{</* lesson/audience-toggle */>}}` above them.

## 4. Preview both audiences

```bash
hugo server
```

Open an episode and the All-in-One page, switch between learner and instructor, and check any local links or
page-bundle images. Add `?view=learner` or `?view=instructor` to test a direct audience URL.

## 5. Validate and build

```bash
go run github.com/oer-particle-physics/hugo-styles/cmd/hugo-styles-migrate@latest check .
hugo --gc --minify --panicOnWarning
```

The validator checks episode metadata types and values, unique weights, glossary/profile targets, image alt text,
and unsupported legacy syntax. The Hugo build catches rendering warnings. Maintainers can run the full browser
and rendered-link suite from [hugo-styles Maintenance]({{< relref "/docs/hugo-styles-maintenance" >}}).

## Live lesson examples

- [Designing a Hugo-first lesson]({{< relref "/episodes/01-introduction" >}})
- [Using challenge and solution blocks]({{< relref "/episodes/02-facilitating-activity" >}})
- [Setup choices, profiles, and native embeds]({{< relref "/episodes/03-setup-choices" >}})
- [Physics notation, diagrams, and richer docs]({{< relref "/episodes/04-physics-doc-features" >}})
- [Learner setup page]({{< relref "/learners/setup" >}})
