+++
title = "Quickstart"
weight = 10
+++

This page is for lesson authors creating a new repository from
[`hugo-styles-template`](https://github.com/oer-particle-physics/hugo-styles-template).
If you want to work on the shared module itself, start with the repository
README and the [hugo-styles Maintenance]({{< relref "/docs/hugo-styles-maintenance" >}}) page instead.

Use this order in a fresh lesson repository:

1. run local setup checks
2. update `hugo.toml`
3. replace the sample content
4. come back here for deeper authoring, deployment, migration, and update guidance

## Before you start

- Install [Hugo Extended](https://gohugo.io/installation/).
- For normal authoring from `hugo-styles-template`, local Go is optional because the template commits `_vendor/`.
- Install [Go](https://go.dev/doc/install) only if you plan to maintain module versions locally.
- You only need GitHub Pages access if you plan to deploy with the template workflow.
- You do not need Node.js for normal lesson authoring. It is only used in the `hugo-styles` module repository for
  frontend asset maintenance and browser tests.

## Create a new lesson

1. Create a new repository from
   [`hugo-styles-template`](https://github.com/oer-particle-physics/hugo-styles-template)
   by clicking **Use this template** at the top right of the repository page.
2. Run `hugo version` and `hugo server` to confirm the local toolchain before editing lesson content.
3. Update the site metadata in `hugo.toml`.
   Change both the site-level `title` and `params.lesson.title`; the former controls Hugo's site title,
   while the latter is reused by lesson components and metadata snippets.
   Update `params.lesson.repo` and `params.lesson.docsRepo` because they are used for the source, edit, and citation links in the page footer.
   The top-nav GitHub icon is configured separately in `[[menus.main]]`, so update that menu URL too.
   If the repository's default branch is not `main`, set both `params.lesson.editBranch` and
   `params.versioning.defaultBranch` to its name.
   If you plan to deploy on GitHub Pages, set `baseURL` to `https://<account>.github.io/<repo>/`.
   In the repository settings, enable Pages and select `GitHub Actions` as the source before the first push to the default branch.
4. Add episodes with `hugo new content --kind episode episodes/my-episode/index.md`.
   Episode order is controlled by front matter `weight` (not by file or folder name).
   Use unique integer weights such as `10`, `20`, `30`; numbered slugs like `01-intro` are optional.
5. Customize the landing page by editing `content/_index.md` while keeping `layout = "hextra-home"`.
   See [Components]({{< relref "/docs/components" >}}) for the recommended homepage pattern.
   If you want authors rendered on the homepage, add an `AUTHORS` file in the repository root
   and list one GitHub handle per line.

## Deploy on GitHub Pages

The thin template repository includes a GitHub Actions workflow that builds the site and publishes the generated artifact to GitHub Pages.
See the [Deployment]({{< relref "/docs/deployment" >}}) guide for the exact steps and the [Troubleshooting]({{< relref "/docs/troubleshooting" >}}) guide for common failures.

## Recommended repository secret for automated upstream refreshes

If you want the **Refresh vendored Hugo modules** workflow to keep opening update PRs without manual intervention,
add a repository Actions secret named `WORKFLOW_SYNC_TOKEN`.
The managed refresh workflow will use it automatically when it is present.

Use either:

- a fine-grained personal access token scoped to this repository with `Contents: Read and write`, `Pull requests: Read and write`, and `Workflows: Read and write`
- or a GitHub App installation token with the same repository permissions

This extra token is needed because upstream refreshes can update the managed workflow files under `.github/workflows/`.
GitHub's default `GITHUB_TOKEN` can create the PR, but GitHub rejects pushes that modify workflow files unless the token also has workflow write permission.

If `WORKFLOW_SYNC_TOKEN` is not configured, the refresh workflow can still work for releases that only change `go.mod`, `go.sum`, scripts, or `_vendor/`,
but it may fail on releases that also update the managed workflow wrappers.

## Migrating existing lessons

If you are migrating from a legacy Carpentries-style lesson, use the
[Migration Guide]({{< relref "/docs/migration" >}}).
That guide covers when to run migration checks and conversion commands.

## Further reading

- [hugo-styles-template](https://github.com/oer-particle-physics/hugo-styles-template)
- [Hextra Getting Started](https://imfing.github.io/hextra/docs/getting-started/)
- [Hextra Tabs](https://imfing.github.io/hextra/docs/guide/shortcodes/tabs/)
- [Hextra Deploy Site](https://imfing.github.io/hextra/docs/guide/deploy-site/)
- [Hugo Installation](https://gohugo.io/installation/)
- [Hugo Modules](https://gohugo.io/hugo-modules/)
- [Hugo Embedded Shortcodes](https://gohugo.io/content-management/shortcodes/#embedded)
