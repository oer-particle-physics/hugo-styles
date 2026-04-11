+++
title = "Quickstart"
weight = 10
+++

This page is for lesson authors creating a new repository from
[`hugo-styles-template`](https://github.com/oer-particle-physics/hugo-styles-template).
If you want to work on the shared module itself, start with the repository
README and the [Update Guide]({{< relref "/docs/updates" >}}) instead.

## Before you start

- Install [Hugo Extended](https://gohugo.io/installation/).
- Install [Go](https://go.dev/doc/install) so Hugo can resolve modules and run the shared checker.
- You only need GitHub Pages access if you plan to deploy with the template workflow.
- You do not need Node.js for normal lesson authoring. It is only used in the `hugo-styles` module repository for search bundle maintenance.

## Create a new lesson

1. Create a new repository from `hugo-styles-template`.
2. Update the site metadata in `hugo.toml`.
   In particular, `params.lesson.repo` is used for the source links and the GitHub button in the top navigation.
3. Run `hugo server`.
4. Add episodes with `hugo new --kind episode episodes/my-episode/index.md`.
5. Customize the landing page by editing `content/_index.md` while keeping `layout = "hextra-home"`.
   See [Components]({{< relref "/docs/components" >}}) for the recommended homepage pattern.
   If you want authors rendered on the homepage, add an `AUTHORS` file in the repository root
   and list one GitHub handle per line.

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

## Further reading

- [hugo-styles-template](https://github.com/oer-particle-physics/hugo-styles-template)
- [Hextra Getting Started](https://imfing.github.io/hextra/docs/getting-started/)
- [Hextra Tabs](https://imfing.github.io/hextra/docs/guide/shortcodes/tabs/)
- [Hextra Deploy Site](https://imfing.github.io/hextra/docs/guide/deploy-site/)
- [Hugo Installation](https://gohugo.io/installation/)
- [Hugo Modules](https://gohugo.io/hugo-modules/)
- [Hugo Embedded Shortcodes](https://gohugo.io/content-management/shortcodes/#embedded)
