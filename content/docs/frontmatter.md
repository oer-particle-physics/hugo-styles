+++
title = "Front Matter"
weight = 30
+++

Use episode front matter to describe both the teaching flow and the navigation model.

## Required episode fields

- `title`: a non-empty string shown in navigation and episode headers
- `weight`: an integer order in the lesson
- `objectives`: a non-empty list of non-empty learning-goal strings
- `keypoints`: a non-empty list of non-empty recap strings

## Optional episode fields

- `questions`: a list of non-empty learner-facing strings, omitted for episodes that pose none
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

If you want an authors block on the homepage, list the authors in a root-level `CITATION.cff`.
The `lesson/authors` shortcode renders each author with a link to their `orcid`, and to their
GitHub profile when `alias` holds their GitHub handle. The ORCID and GitHub columns only appear
when at least one author has a value for them.

Each person's displayed name joins `given-names`, `name-particle`, `family-names`, and
`name-suffix` in that order, with spaces, skipping any that are missing. For an organisation or
group listed as an author, give a single `name` instead. An author with neither is listed under
their `alias`.

```yaml
authors:
  - given-names: Alexander   # shown as "Alexander von Humboldt III"
    name-particle: von
    family-names: Humboldt
    name-suffix: III
    orcid: https://orcid.org/0000-0002-1825-0097
  - name: HEP Software Foundation   # shown as "HEP Software Foundation"
  - alias: octocat                  # shown as "octocat", linked to github.com/octocat
```

An empty `CITATION.cff` renders no table. Otherwise the build fails unless the file has a
non-empty `authors` list and every author has a name or `alias`. To check the rest of the
file against the CFF specification, run a validator such as
[cffconvert](https://github.com/citation-file-format/cffconvert) (`cffconvert --validate`).

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
