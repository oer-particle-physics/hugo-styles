+++
title = "Front Matter"
weight = 30
+++

Use episode front matter to describe both the teaching flow and the navigation model.

## Required episode fields

- `title`: the page title shown in navigation and episode headers
- `weight`: numeric order in the lesson
- `questions`: list of learner-facing questions rendered near the top
- `objectives`: list of learning goals rendered near the top
- `keypoints`: list of recap points rendered near the bottom

## Optional episode fields

- `teaching`: teaching time in minutes
- `exercises`: exercise time in minutes
- `summary`: override for episode card summaries
- `draft`: hide unfinished content from production builds
- `[tabs] sync = false`: disable synced Hextra tabs on a page

## Homepage metadata

The lesson homepage usually lives in `content/_index.md` and keeps `layout = "hextra-home"`.
That page can hold the homepage content blocks in [Components]({{< relref "/docs/components" >}}).

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

## Section pages

Section index pages such as `content/episodes/_index.md` or `content/glossary/_index.md` usually only need:

- `title`
- `weight`
- optional descriptive body text

## Validator expectations

The shared `check` command currently expects:

- every episode to define all required fields
- every episode weight to be unique
- weight values to be integers

If you intentionally want a draft episode excluded from normal builds, still give it complete front matter. That keeps previews and future publication simpler.
