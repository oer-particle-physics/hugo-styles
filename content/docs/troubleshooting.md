+++
title = "Troubleshooting"
weight = 70
+++

## The checker reports legacy syntax in a migrated repo

Run the checker on the migrated output and search for:

- leftover fenced-attribute syntax such as `{: .challenge}`
- Liquid includes or variables
- raw iframe embeds

The checker is designed to catch exactly those leftovers so you can clean them up before publishing.

## Tabs are not syncing

Check that:

- the tab labels match exactly
- the page does not set `[tabs] sync = false`
- you are using Hextra's `tabs` and `tab` shortcodes, not a copied custom variant

## A glossary or profile link is broken

The shortcode target should match the content slug, for example:

- `content/glossary/formative-assessment.md` -> `{{</* glossary formative-assessment */>}}`
- `content/profiles/workshop-host.md` -> `{{</* profile workshop-host */>}}`

## TOML front matter reports an invalid escaped character for LaTeX

TOML double-quoted strings treat backslashes as escapes, so `\( ... \)` fails unless every backslash is doubled.

Prefer literal strings in front matter:

```toml
objectives = [
  'Clone and configure the \(t\bar{t}\gamma\) analysis repository.'
]
```

If the same expression renders in the Markdown body but not in `questions`, `objectives`, or `keypoints`, check your lesson's `hugo.toml`. When you override `[markup.goldmark.extensions]`, keep the passthrough block from `hugo-styles` so `\(...\)` and `\[...\]` still reach the math renderer.

## Search updates but results seem odd

Aggregated pages such as `All-in-One`, `Key Points`, and `External Links` are intentionally excluded from indexing so they do not crowd out the main lesson pages. If a result feels missing, check whether the relevant page is a generated resource or a real content page.

## A downstream lesson does not pick up a shared-module fix

Run:

```bash
hugo mod get -u github.com/oer-particle-physics/hugo-styles@latest
hugo mod tidy
```

Then rebuild locally. If the change still does not appear, look for a local override in `layouts/`, `assets/`, or `archetypes/`.
