+++
title = "Hextra Feature Guide"
weight = 55
+++

This page is a live, offline-capable showcase of the Hextra features selected for lesson authors. Use the copy button on any code block, and use the **Copy page** menu beside the title to view this page as Markdown.

## Math

Inline notation such as \(p_T > 25\,\mathrm{GeV}\) stays readable in a sentence. Display equations can carry the central result:

$$
Z \approx \frac{N_{\mathrm{sig}}}{\sqrt{N_{\mathrm{bkg}} + (\delta N_{\mathrm{bkg}})^2}}
$$

{{< details title="Use it" closed="true" >}}
```md
Inline: \(p_T > 25\,\mathrm{GeV}\)

$$
Z \approx \frac{N_{\mathrm{sig}}}{\sqrt{N_{\mathrm{bkg}} + (\delta N_{\mathrm{bkg}})^2}}
$$
```
{{< /details >}}

## Mermaid diagrams

```mermaid
flowchart LR
  A[Toy events] --> B[pT selection]
  B --> C[Signal region]
  B --> D[Control region]
  C --> E[Count and report]
  D --> E
```

{{< details title="Use it" closed="true" >}}
````md
```mermaid
flowchart LR
  A[Toy events] --> B[pT selection]
  B --> C[Signal region]
  B --> D[Control region]
  C --> E[Count and report]
  D --> E
```
````
{{< /details >}}

## Synced tabs

Choose a shell here:

{{< tabs >}}
{{< tab name="bash" selected=true >}}`source .venv/bin/activate`{{< /tab >}}
{{< tab name="fish" >}}`source .venv/bin/activate.fish`{{< /tab >}}
{{< /tabs >}}

The matching choice follows in this second group:

{{< tabs >}}
{{< tab name="bash" selected=true >}}`python analysis.py`{{< /tab >}}
{{< tab name="fish" >}}`python analysis.py`{{< /tab >}}
{{< /tabs >}}

{{< details title="Use it" closed="true" >}}
{{< codeblock lang="text" >}}
{{</* tabs */>}}
{{</* tab name="bash" selected=true */>}}bash command{{</* /tab */>}}
{{</* tab name="fish" */>}}fish command{{</* /tab */>}}
{{</* /tabs */>}}
{{< /codeblock >}}

Enable sync globally, then opt out for a page only when its tab groups are unrelated:

```toml
# hugo.toml
[params.page.tabs]
  sync = true

# page front matter
[tabs]
  sync = false
```
{{< /details >}}

## Details and steps

{{< details title="Optional interpretation" closed="true" >}}
The toy counts demonstrate the rendering workflow, not a statistical claim.
{{< /details >}}

{{% steps %}}

### Record the input

Name the sample and the expected columns.

### Apply one explicit cut

Keep the first selection easy to reproduce.

### Report counts with context

State the signal and control windows beside the result.

{{% /steps %}}

{{< details title="Use them" closed="true" >}}
{{< codeblock lang="text" >}}
{{</* details title="Optional interpretation" closed="true" */>}}
Extra context.
{{</* /details */>}}

{{%/* steps */%}}
### Record the input
### Apply one explicit cut
### Report counts with context
{{%/* /steps */%}}
{{< /codeblock >}}
{{< /details >}}

## Cards, badges, and icons

{{< cards cols="2" >}}
{{< card link="/docs/components" title="Lesson components" subtitle="Project-owned teaching features." icon="puzzle" tag="live" tagColor="blue" >}}
{{< card link="/all-in-one" title="All-in-One" subtitle="The complete lesson with audience controls." icon="book-open" tag="generated" tagColor="green" >}}
{{< /cards >}}

{{< badge content="offline asset" color="green" icon="check" >}}
{{< badge content="copyable source" color="blue" icon="clipboard" >}}
{{< icon name="beaker" >}} analysis-ready

{{< details title="Use them" closed="true" >}}
{{< codeblock lang="text" >}}
{{</* cards cols="2" */>}}
{{</* card link="/docs/components" title="Lesson components" icon="puzzle" */>}}
{{</* card link="/all-in-one" title="All-in-One" icon="book-open" */>}}
{{</* /cards */>}}

{{</* badge content="offline asset" color="green" icon="check" */>}}
{{</* icon name="beaker" */>}}
{{< /codeblock >}}
{{< /details >}}

## File tree

{{< filetree/container >}}
  {{< filetree/folder name="content" >}}
    {{< filetree/folder name="episodes" >}}
      {{< filetree/file name="01-introduction/index.md" >}}
    {{< /filetree/folder >}}
    {{< filetree/folder name="docs" state="open" >}}
      {{< filetree/folder name="hextra-features" state="open" >}}
        {{< filetree/file name="index.md" >}}
        {{< filetree/file name="particle-analysis.ipynb" >}}
        {{< filetree/file name="particle-analysis-handout.pdf" >}}
      {{< /filetree/folder >}}
    {{< /filetree/folder >}}
  {{< /filetree/folder >}}
{{< /filetree/container >}}

{{< details title="Use it" closed="true" >}}
{{< codeblock lang="text" >}}
{{</* filetree/container */>}}
  {{</* filetree/folder name="content" */>}}
    {{</* filetree/file name="hugo.toml" */>}}
  {{</* /filetree/folder */>}}
{{</* /filetree/container */>}}
{{< /codeblock >}}
{{< /details >}}

## Syntax highlighting and code copy

```python {filename="selection.py",linenos=table,hl_lines=[2,4]}
events = load_events("toy-events.json")
selected = [event for event in events if event.pt > 25]
signal = [event for event in selected if 80 <= event.mass < 100]
print(f"signal={len(signal)}")
```

The copy button in the upper-right corner is enabled automatically.

{{< details title="Use it" closed="true" >}}
````md
```python {filename="selection.py",linenos=table,hl_lines=[2,4]}
selected = [event for event in events if event.pt > 25]
```
````
{{< /details >}}

## Image zoom

Click this page-bundle image to enlarge it:

![A four-stage toy analysis moving from events through selection and regions to a reported result.](analysis-flow.svg)

{{< details title="Use it" closed="true" >}}
```md
![Meaningful alt text](analysis-flow.svg)
```

Enable the local vendored zoom script globally, with an optional page opt-out:

```toml
# hugo.toml
[params.imageZoom]
  enable = true
  js = "js/vendor/medium-zoom.min.js"

# page front matter
imageZoom = false
```
{{< /details >}}

## Local PDF

[Open or download the handout](particle-analysis-handout.pdf) if the embedded viewer is not convenient.

{{< pdf "particle-analysis-handout.pdf" >}}

{{< details title="Use it" closed="true" >}}
{{< codeblock lang="text" >}}
[Open or download the handout](particle-analysis-handout.pdf).

{{</* pdf "particle-analysis-handout.pdf" */>}}
{{< /codeblock >}}
{{< /details >}}

## Local Jupyter notebook

The code cell and its saved deterministic output render without a notebook server or network request.

{{% jupyter "particle-analysis.ipynb" %}}

{{< details title="Use it" closed="true" >}}
{{< codeblock lang="text" >}}
{{%/* jupyter "particle-analysis.ipynb" */%}}
{{< /codeblock >}}
Keep the `.ipynb` file in the same page bundle as `index.md`.
{{< /details >}}

## Banner, search, and theme

The dismissible banner at the top identifies this documentation as a live demo. Search is visible in the navigation, and the theme menu offers Light, Dark, and System.

{{< details title="Configure them" closed="true" >}}
```toml
[params.banner]
  key = "hugo-styles-feature-demo-v1"
  message = "This documentation is a live feature demo. See the [feature guide](/docs/hextra-features/)."

[params.search]
  enable = true
  type = "flexsearch"

[params.theme]
  default = "system"
  displayToggle = true
```
Use a root-relative path for an internal banner link so Hugo applies the site's base path and keeps navigation in the same tab. Changing the banner key makes a revised announcement visible even to readers who dismissed an older one.
{{< /details >}}

## Markdown context menu

Use **Copy page** beside this title, or open its menu to view the exact Markdown source.

{{< details title="Configure it" closed="true" >}}
```toml
[outputs]
  page = ["HTML", "Markdown"]
  section = ["HTML", "Markdown"]

[params.page.contextMenu]
  enable = true
```

Opt out for one page with `contextMenu = false` in front matter.
{{< /details >}}

## Version selector

The deployed demo builds **Latest** and the archived **v0.4.0** site. Use **Versions** in the navigation to move between them.

{{< details title="Configure it" closed="true" >}}
```toml
[params.versioning]
  enable = true
  defaultBranch = "main"

  [params.versioning.latest]
    enable = true
    label = "Latest"

  [params.versioning.tags]
    refs = ["v0.4.0"]
```
{{< /details >}}

The curated gallery intentionally leaves out comments/blog features, Hextra's separate data-file term system, include-based content reuse, and CDN-backed terminal recordings. They overlap with this lesson model or weaken the offline build contract.
