+++
title = "Components"
weight = 40
+++

These are the supported lesson-specific components on top of Hextra.
This page focuses on the Carpentries-style layer added by `hugo-styles`.
For the broader Hextra and Hugo feature set, use the
[Reference]({{< relref "/reference" >}}) page and the linked upstream docs.

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

For the full Hextra tab behavior and additional examples, see
[Hextra Tabs](https://imfing.github.io/hextra/docs/guide/shortcodes/tabs/).

## Videos

Use Hugo-native shortcodes when supported:

```text
{{</* youtube aqz-KE-bpKQ */>}}
```

Avoid raw iframe HTML for providers that Hugo already supports directly.
For the full list of embedded providers, see
[Hugo Embedded Shortcodes](https://gohugo.io/content-management/shortcodes/#embedded).

## Hextra home components

For landing pages, overview sections, or documentation marketing surfaces, you can also use Hextra's own home components directly:

- `hextra/hero-badge`
- `hextra/hero-headline`
- `hextra/hero-subtitle`
- `hextra/hero-button`
- `hextra/feature-grid`
- `hextra/feature-card`

These are a good fit for non-pedagogy UI such as homepage highlights, getting-started links, or overview cards. The homepage of `hugo-styles` now uses them as the default pattern.
For more theme-level homepage guidance, see
[Hextra Getting Started](https://imfing.github.io/hextra/docs/getting-started/).

If you want the same overall homepage structure that Hextra uses, set:

```toml
+++
layout = "hextra-home"
+++
```

and compose the page from `hextra/hero-badge`, `hextra/hero-headline`, `hextra/hero-subtitle`, `hextra/hero-button`, and `hextra/feature-grid`.

For most downstream lesson repositories, the easiest customization path is:

1. Keep `layout = "hextra-home"` in `content/_index.md`.
2. Edit the copy, links, and cards in that file.
3. Prefer Hextra shortcodes over custom layout overrides.

A minimal starting point looks like this:

```md
+++
title = "My Lesson"
layout = "hextra-home"
+++

{{< hextra/hero-badge link="docs/quickstart/" >}}
Workshop-ready lesson materials
{{< /hextra/hero-badge >}}

<div class="hx:mt-6 hx:mb-6">
{{< hextra/hero-headline >}}
Teach particle physics with a Carpentries-style lesson
{{< /hextra/hero-headline >}}
</div>

<div class="hx:mb-12">
{{< hextra/hero-subtitle >}}
Reuse the familiar teaching structure while keeping the site easy to maintain with Hugo.
{{< /hextra/hero-subtitle >}}
</div>

<div class="hx:mb-6">
{{< hextra/hero-button text="Start Learning" link="episodes/01-introduction/" >}}
</div>

<div class="hx:mt-6"></div>

{{< hextra/feature-grid cols="3" >}}
{{< hextra/feature-card title="Episodes" subtitle="Step-by-step lesson flow." icon="book-open" link="episodes/" class="hx:min-h-[220px]" >}}
{{< hextra/feature-card title="Setup" subtitle="Environment and data requirements." icon="cog" link="learners/setup/" class="hx:min-h-[220px]" >}}
{{< hextra/feature-card title="Teaching Notes" subtitle="Support for instructors and helpers." icon="academic-cap" link="instructors/" class="hx:min-h-[220px]" >}}
{{< /hextra/feature-grid >}}
```

Recommended customization boundaries:

- Change the homepage content in `content/_index.md`.
- Change site-wide metadata such as the title, repository URL, and menu entries in `hugo.toml`.
- Add more `hextra/feature-card` entries or additional content sections if the lesson needs more orientation material.
- Avoid copying `layouts/hextra-home.html`, `navbar.html`, or other theme templates unless you need a structural change that content cannot express.
- Reach for custom CSS only after the built-in Hextra shortcodes and utilities stop being enough.

That keeps downstream lesson repositories thin and makes future Hextra or `hugo-styles` updates much easier to adopt.

For text-only homepages, adding small wrapper divs with Hextra utility classes like `hx:mt-6`, `hx:mb-6`, and `hx:mb-12` is often enough to create better breathing room between the hero badge, headline, subtitle, button, and card grid.
