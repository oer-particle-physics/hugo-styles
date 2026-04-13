+++
title = "Hextra Features for Physics Lessons"
weight = 55
+++

Use this page for practical, lesson-author-focused Hextra features that are especially useful in particle-physics tutorials.
It complements upstream Hextra docs with domain-shaped examples instead of re-documenting every option.

{{< callout type="note" title="Scope" >}}
This page covers the high-value defaults for lesson authors.
For full option matrices and edge cases, use the linked upstream Hextra pages.
{{< /callout >}}

## High-value feature map

| Feature | Good fit in physics lessons | Upstream |
| --- | --- | --- |
| LaTeX | equations, symbols, uncertainty notation | [LaTeX](https://imfing.github.io/hextra/docs/guide/latex/) |
| Mermaid | selection flow and collaboration diagrams | [Diagrams](https://imfing.github.io/hextra/docs/guide/diagrams/) |
| Tabs | OS/shell/package-manager variants | [Tabs](https://imfing.github.io/hextra/docs/guide/shortcodes/tabs/) |
| Details | optional derivations and deep dives | [Details](https://imfing.github.io/hextra/docs/guide/shortcodes/details/) |
| FileTree | repository and data layout orientation | [FileTree](https://imfing.github.io/hextra/docs/guide/shortcodes/filetree/) |
| Cards | quick links to key lesson resources | [Cards](https://imfing.github.io/hextra/docs/guide/shortcodes/cards/) |
| Steps | short procedural workflows | [Steps](https://imfing.github.io/hextra/docs/guide/shortcodes/steps/) |
| Syntax highlighting | line numbers and highlighted snippets | [Syntax Highlighting](https://imfing.github.io/hextra/docs/guide/syntax-highlighting/) |

## LaTeX for analysis notation

Inline notation works well for terms like \(p_T\), \(\eta\), and \(\Delta R\).

For standalone equations:

$$
\begin{aligned}
N_{\text{sig}} &= N_{\text{obs}} - N_{\text{bkg}} \\
Z &\approx \frac{N_{\text{sig}}}{\sqrt{N_{\text{bkg}} + (\delta N_{\text{bkg}})^2}}
\end{aligned}
$$

`hugo-styles` enables Goldmark passthrough delimiters so `\(...\)` and `$$...$$` render correctly.

## Mermaid for analysis and teaching flow

```mermaid
flowchart LR
  A[Raw events] --> B[Reconstruction]
  B --> C[Object selection]
  C --> D[Signal region]
  C --> E[Control region]
  D --> F[Histogram + fit]
  E --> F
```

```mermaid
sequenceDiagram
  participant DAQ as Data acquisition
  participant RECO as Reconstruction
  participant ANA as Analysis team
  participant REV as Internal review
  DAQ->>RECO: Calibrated data stream
  RECO->>ANA: Ntuples + metadata
  ANA->>REV: Selection + uncertainty model
  REV-->>ANA: Feedback and sign-off
```

## Tabs for setup variants

Use the same labels across tab groups when you want sync behavior.

{{< tabs >}}
{{< tab name="bash" selected=true >}}
```bash
python -m venv .venv
source .venv/bin/activate
```
{{< /tab >}}
{{< tab name="zsh" >}}
```zsh
python -m venv .venv
source .venv/bin/activate
```
{{< /tab >}}
{{< tab name="fish" >}}
```fish
python -m venv .venv
source .venv/bin/activate.fish
```
{{< /tab >}}
{{< /tabs >}}

## File orientation with FileTree and Cards

{{< filetree/container >}}
  {{< filetree/folder name="content" >}}
    {{< filetree/folder name="episodes" >}}
      {{< filetree/file name="01-introduction/index.md" >}}
      {{< filetree/file name="04-physics-doc-features/index.md" >}}
    {{< /filetree/folder >}}
    {{< filetree/folder name="learners" state="closed" >}}
      {{< filetree/file name="setup.md" >}}
    {{< /filetree/folder >}}
  {{< /filetree/folder >}}
  {{< filetree/file name="hugo.toml" >}}
{{< /filetree/container >}}

{{< cards cols="2" >}}
{{< card link="/docs/authoring" title="Authoring model" subtitle="Where pedagogy conventions live." icon="book-open" tag="core" >}}
{{< card link="/docs/components" title="Component reference" subtitle="Lesson shortcodes and examples." icon="puzzle" tag="core" >}}
{{< card link="/docs/updates" title="Update workflow" subtitle="How downstream lessons consume module releases." icon="refresh" tag="maintenance" tagColor="blue" >}}
{{< card link="/reference" title="External references" subtitle="Hextra and Hugo upstream docs." icon="external-link" tag="upstream" tagColor="gray" >}}
{{< /cards >}}

## Steps and optional depth

{{% steps %}}

### Define the learning objective

Write the specific analysis skill the learner should gain.

### Show the minimal reproducible command path

Keep platform-specific command variants in short tab groups.

### Add one visual summary

Use Mermaid for workflow shape before detailed prose.

### Add one optional deep-dive block

Use `details` for derivations that are useful but not required for first pass.

{{% /steps %}}

{{< details title="Optional derivation: Poisson counting uncertainty" closed="true" >}}
For quick estimates in counting analyses, \(\sigma_N \approx \sqrt{N}\) is often enough to explain uncertainty propagation in early lesson stages.
{{< /details >}}

## Syntax-highlighted snippets

```python {filename="selection.py",linenos=table,hl_lines=[2,5],linenostart=1}
events = load_events("events.parquet")
selected = events[(events.pt > 25) & (abs(events.eta) < 2.4)]
control = selected[selected.m_ll.between(70, 110)]
signal = selected[selected.m_ll.between(110, 160)]
plot_mass(control, signal)
```

For badge/PDF/video and other utility shortcodes, see [Other Shortcodes](https://imfing.github.io/hextra/docs/guide/shortcodes/others/).
