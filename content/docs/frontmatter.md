+++
title = "Front Matter"
weight = 30
+++

Use episode front matter to describe both the teaching flow and the navigation model.

## Required episode fields

- `title`: a non-empty string shown in navigation and episode headers
- `weight`: an integer order in the lesson
- `questions`: a non-empty list of non-empty learner-facing strings
- `objectives`: a non-empty list of non-empty learning-goal strings
- `keypoints`: a non-empty list of non-empty recap strings

## Optional episode fields

- `teaching`: a non-negative integer teaching time in minutes
- `exercises`: a non-negative integer exercise time in minutes
- `summary`: override for episode card summaries
- `draft`: hide unfinished content from production builds
- `imageZoom = false`: disable click-to-zoom behavior for Markdown images on a page
- `[tabs] sync = false`: disable synced Hextra tabs on a page

## Homepage metadata

The lesson homepage usually lives in `content/_index.md` and keeps `layout = "hextra-home"`.
That page can hold the homepage content blocks in [Components]({{< relref "/docs/components" >}}).
Shared lesson metadata such as `params.lesson.title`, `params.lesson.tagline`, and
`params.lesson.description` belongs in `hugo.toml`, and `lesson/meta` can reuse those values
inside Markdown body content when you want homepage copy to stay aligned with the config.

If you want an authors block on the homepage, add a root-level `AUTHORS` file. The `lesson/authors`
shortcode reads that file directly and renders the contributors there.

## Example episode front matter

```toml
+++
title = "Using challenge and solution blocks"
weight = 20
teaching = 15
exercises = 10
questions = ["How should active-learning blocks behave in a Hugo-based lesson?"]
objectives = ["Use challenge, hint, and solution shortcodes naturally."]
keypoints = ["Hints and solutions should stay collapsible."]
+++
```

## Math in TOML front matter

If you use TOML front matter, prefer literal strings for LaTeX so backslashes are not treated as TOML escapes:

```toml
objectives = [
  'Explain the \(t\bar{t}\gamma\) workflow.'
]
```

Double-quoted TOML strings also work, but every backslash must be escaped, for example `\\(` and `\\bar`.

## Section pages

Section index pages such as `content/episodes/_index.md` or `content/glossary/_index.md` usually only need:

- `title`
- `weight`
- optional descriptive body text

## Validator expectations

The shared `check` command enforces the types and non-empty values listed above and requires every episode
weight to be unique. Values such as `weight = 2.5`, `questions = "one string"`, empty list items, and
negative times fail validation.

If you intentionally want a draft episode excluded from normal builds, still give it complete front matter. That keeps previews and future publication simpler.
