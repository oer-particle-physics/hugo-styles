+++
title = "Move author information to CITATION.cff"

[params.upgrade]
release = "v0.6.0"
severity = "breaking"
check = "legacy-authors"
+++

## Who is affected

Lessons using the `lesson/authors` shortcode to display contributors on the homepage.
Lessons that do not display this table do not need to add a citation file for this feature.

## What changes

Starting with v0.6.0, `lesson/authors` reads only the root-level `CITATION.cff`.
Support for `AUTHORS` was removed without a transition period. If the citation file
is missing or empty, the table disappears even though the Hugo build can pass.
If both files exist, the table may change because it now uses the citation authors.

## What to do

1. Transfer or reconcile all intended contributors from `AUTHORS` into the `authors`
   list in `CITATION.cff`. Preserve existing citation metadata and authors.
2. Use `given-names` and `family-names` for people, `name` for organisations, and
   `alias` for GitHub handles. Put ORCID URLs in `orcid`, not GitHub handles.
3. After checking the result, remove `AUTHORS` if nothing else in the repository
   needs it. The updater never deletes or converts contributor information automatically.

A minimal example (replace the title and author with the lesson's actual information):

```yaml
cff-version: 1.2.0
message: "Please cite this lesson using the metadata in this file."
title: "Example lesson"
authors:
  - given-names: Ada
    family-names: Lovelace
    alias: example-handle
```

See the [author metadata documentation](https://oer-particle-physics.github.io/hugo-styles/docs/frontmatter/)
for name particles, suffixes, organisation names, and ORCID links.

## How to verify

Validate the file with `cffconvert --validate` or the managed **Validate CITATION.cff**
workflow. Build the lesson, then inspect the homepage: every intended contributor
should be present, with the correct GitHub and ORCID links. A passing Hugo build
does not establish that the contributor list is complete.
