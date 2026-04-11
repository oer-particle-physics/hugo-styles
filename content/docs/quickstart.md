+++
title = "Quickstart"
weight = 10
+++

## Create a new lesson

1. Create a new repository from `hugo-styles-template`.
2. Update the site metadata in `hugo.toml`.
   In particular, `params.lesson.repo` is used for the source links and the GitHub button in the top navigation.
3. Run `hugo server`.
4. Add episodes with `hugo new --kind episode episodes/my-episode/index.md`.
5. Customize the landing page by editing `content/_index.md` while keeping `layout = "hextra-home"`.
   See [Components]({{< relref "/docs/components" >}}) for the recommended homepage pattern.

## Build locally

```bash
hugo server
```

## Run the shared checks

```bash
go run github.com/oer-particle-physics/hugo-styles/cmd/hugo-styles-migrate@latest check .
```

## Deploy on GitHub Pages

The thin template repository includes a GitHub Actions workflow that builds the site and publishes the generated artifact to GitHub Pages.
See the [Deployment]({{< relref "/docs/deployment" >}}) guide for the exact steps and the [Troubleshooting]({{< relref "/docs/troubleshooting" >}}) guide for common failures.

## Use the migration helper

```bash
go run github.com/oer-particle-physics/hugo-styles/cmd/hugo-styles-migrate@latest check .
```

That command flags legacy Carpentries syntax, unsupported Jekyll constructs, and missing lesson metadata before you start a migration.
