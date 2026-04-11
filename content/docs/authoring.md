+++
title = "Authoring Guide"
weight = 20
+++

Use this guide for the overall lesson-writing model. The detailed field and shortcode references live in:

- [Front Matter]({{< relref "/docs/frontmatter" >}})
- [Components]({{< relref "/docs/components" >}})
- [Glossary & Profiles]({{< relref "/docs/glossary-profiles" >}})

For generic Hextra theme behavior and deeper Hugo mechanics, use the
[Reference]({{< relref "/reference" >}}) page to jump to the upstream docs.

## Lesson structure

Use Hugo sections to match the lesson model:

- `content/episodes/` for the main teaching flow
- `content/learners/` for learner-facing supporting pages
- `content/instructors/` for instructor-only notes
- `content/glossary/` for reusable definitions
- `content/profiles/` for learner or maintainer personas
- `content/reference.md` for a broader reference page when you need one
- `content/key-points.md`, `content/all-in-one.md`, and `content/extract-all-images.md` for generated resource surfaces

Episodes should stay short, ordered, and self-contained. Supplemental setup, reference, or facilitation material should move into the supporting sections instead of bloating the main narrative.

## Pedagogy blocks

The goal is to preserve the Carpentries teaching vocabulary while letting Hextra handle the overall visual language. The boxes below are intentional teaching signals, not generic decoration.

### Challenge, hint, and solution

```text
{{</* challenge title="Warm-up exercise" */>}}
Prompt the learner to do something concrete.

{{</* hint */>}}
Offer a small nudge when learners get stuck.
{{</* /hint */>}}

{{</* solution */>}}
Provide a worked answer or expected approach.
{{</* /solution */>}}
{{</* /challenge */>}}
```

This renders as a warm challenge panel with nested collapsible hint and solution blocks.

### Callout box types

The `callout` shortcode covers the common Carpentries box families and keeps them visually aligned with Hextra:

- `note`
- `prereq`
- `checklist`
- `testimonial`
- `discussion`
- `warning`
- `caution`

```text
{{</* callout type="note" title="Teaching note" */>}}
Context or a side comment.
{{</* /callout */>}}

{{</* callout type="discussion" title="Discuss" */>}}
Prompt learners to compare approaches.
{{</* /callout */>}}

{{</* callout type="warning" title="Watch out" */>}}
Flag a likely mistake or risky step.
{{</* /callout */>}}
```

Challenge and solution remain their own shortcodes because they carry different authoring and interaction behaviour than generic callouts.

### Audience-specific content

```text
{{</* learner */>}}
Learner-facing note.
{{</* /learner */>}}

{{</* instructor */>}}
Facilitation advice for instructors.
{{</* /instructor */>}}
```

The learner and instructor switch in the top navigation changes which of these blocks are visible, following the same idea as Workbench while staying inside the Hextra shell.

### Generic callouts and spoilers

```text
{{</* callout type="prereq" title="Before you start" */>}}
Explain assumptions or required setup.
{{</* /callout */>}}

{{</* details title="Why this matters" */>}}
Expandable extra context.
{{</* /details */>}}
```

## Tabs for setup variants

Use Hextra's native `tabs` and `tab` shortcodes for repeated setup variants such as operating systems, shells, or package managers. This keeps the lesson source compact and avoids long repeated sections.

### Operating system tabs

~~~text
{{</* tabs */>}}
{{</* tab name="macOS" selected=true */>}}
```bash
brew install hugo go
```
{{</* /tab */>}}
{{</* tab name="Linux" */>}}
```bash
sudo apt install hugo golang
```
{{</* /tab */>}}
{{</* tab name="Windows" */>}}
```powershell
winget install Hugo.Hugo.Extended
```
{{</* /tab */>}}
{{</* /tabs */>}}
~~~

### Shell tabs

~~~text
{{</* tabs */>}}
{{</* tab name="bash" selected=true */>}}
```bash
source ~/.bashrc
```
{{</* /tab */>}}
{{</* tab name="zsh" */>}}
```zsh
source ~/.zshrc
```
{{</* /tab */>}}
{{</* tab name="fish" */>}}
```fish
source ~/.config/fish/config.fish
```
{{</* /tab */>}}
{{</* /tabs */>}}
~~~

### Sync behavior

`hugo-styles` enables synced tabs by default with Hextra's built-in tab sync. That means repeated tab groups with the same labels on a page, such as `bash` / `zsh` / `fish`, stay aligned automatically.

Use synced tabs when:

- the same choice repeats across a page
- learners should not have to reselect their environment each time
- the tab names mean the same thing in each group

Opt out on a page when the tab groups are only locally meaningful:

```toml
+++
[tabs]
  sync = false
+++
```

That disables syncing for all tab groups on that page while keeping the same shortcode syntax.

### Recommended setup structure

- Use tabs for short, parallel variants of the same step.
- Use separate pages or sections when the platform workflows genuinely diverge.
- Keep tab labels short and familiar: `macOS`, `Linux`, `Windows`, `bash`, `zsh`, `fish`.
- Prefer one decision axis per tab group. Do not mix OS and shell in the same tab set.
- When multiple tab groups on a page should sync, keep the labels identical and in the same order.

## Where to see the examples

The example lesson intentionally demonstrates the lesson-specific UI:

- [Designing a Hugo-first lesson]({{< relref "/episodes/01-introduction" >}})
- [Using challenge and solution blocks]({{< relref "/episodes/02-facilitating-activity" >}})
- [Setup choices, profiles, and native embeds]({{< relref "/episodes/03-setup-choices" >}})
- [Learner setup page]({{< relref "/learners/setup" >}})

## Videos

Use Hugo's built-in provider shortcodes when they exist:

```text
{{</* youtube aqz-KE-bpKQ */>}}
```

That keeps the source cleaner than raw iframe HTML and makes migration rules easier to reason about.

## Quality checks

The shared validator can check both legacy lessons and Hugo-native lesson repos:

```bash
go run github.com/oer-particle-physics/hugo-styles/cmd/hugo-styles-migrate@latest check .
```

For Hugo lessons, the validator currently checks:

- required episode metadata
- duplicate episode weights
- unresolved glossary references
- unresolved profile references
- missing image alt text
- unsupported leftover legacy syntax
