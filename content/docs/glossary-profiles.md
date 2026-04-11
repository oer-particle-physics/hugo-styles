+++
title = "Glossary & Profiles"
weight = 50
+++

Glossary terms and learner profiles are normal Hugo content pages, not hidden data files.

## Glossary entries

Create entries in `content/glossary/`:

```bash
hugo new glossary/new-term.md
```

Minimal example:

```toml
+++
title = "Formative Assessment"
summary = "A quick way to learn what participants understand while teaching is still in progress."
weight = 10
+++
```

Reference a glossary term inline with:

```text
{{</* glossary formative-assessment */>}}
```

## Profiles

Create profiles in `content/profiles/`:

```bash
hugo new profile/workshop-host.md
```

Reference a profile inline with:

```text
{{</* profile workshop-host */>}}
```

## Why this content model works well

- entries have their own URLs and summaries
- glossary and profile list pages come for free
- links are validated by the shared checker
- styling stays consistent with the rest of the lesson

## Validation rules

The shared validator reports broken glossary and profile references, so typos in shortcode targets do not quietly ship.
