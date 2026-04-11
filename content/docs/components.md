+++
title = "Components"
weight = 40
+++

These are the supported lesson-specific components on top of Hextra.

## Challenge, hint, and solution

```text
{{</* challenge title="Warm-up exercise" */>}}
Prompt the learner to do something concrete.

{{</* hint */>}}
Offer a small nudge.
{{</* /hint */>}}

{{</* solution */>}}
Provide the worked answer.
{{</* /solution */>}}
{{</* /challenge */>}}
```

Use these for real learning activities, not just decorative emphasis.

## Callouts

Use the `callout` shortcode for the common Carpentries-style boxes:

- `note`
- `prereq`
- `checklist`
- `testimonial`
- `discussion`
- `warning`
- `caution`

```text
{{</* callout type="discussion" title="Compare approaches" */>}}
Ask learners to compare trade-offs before you reveal an answer.
{{</* /callout */>}}
```

## Audience-aware content

```text
{{</* learner */>}}
Learner-facing support text.
{{</* /learner */>}}

{{</* instructor */>}}
Facilitation guidance for instructors.
{{</* /instructor */>}}
```

The top-bar learner/instructor switch controls these blocks.

## Details blocks

```text
{{</* details title="Why this matters" */>}}
Expanded explanation or optional background.
{{</* /details */>}}
```

## Tabs

Use Hextra-native tabs for OS, shell, or package-manager variants:

```text
{{</* tabs */>}}
{{</* tab name="bash" selected=true */>}}
```bash
echo hello
```
{{</* /tab */>}}
{{</* tab name="fish" */>}}
```fish
echo hello
```
{{</* /tab */>}}
{{</* /tabs */>}}
```

Tabs sync by default. Disable sync per page with:

```toml
+++
[tabs]
  sync = false
+++
```

## Videos

Use Hugo-native shortcodes when supported:

```text
{{</* youtube aqz-KE-bpKQ */>}}
```

Avoid raw iframe HTML for providers that Hugo already supports directly.

## Hextra home components

For landing pages, overview sections, or documentation marketing surfaces, you can also use Hextra's own home components directly:

- `hextra/hero-badge`
- `hextra/hero-headline`
- `hextra/hero-subtitle`
- `hextra/hero-button`
- `hextra/feature-grid`
- `hextra/feature-card`

These are a good fit for non-pedagogy UI such as homepage highlights, getting-started links, or overview cards. The homepage of `hugo-styles` now uses them as the default pattern.

If you want the same overall homepage structure that Hextra uses, set:

```toml
+++
layout = "hextra-home"
+++
```

and compose the page from `hextra/hero-badge`, `hextra/hero-headline`, `hextra/hero-subtitle`, `hextra/hero-button`, and `hextra/feature-grid`.
