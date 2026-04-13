+++
title = "Physics notation, diagrams, and richer docs"
weight = 40
teaching = 20
exercises = 10
questions = ["How can we communicate a particle-physics analysis workflow clearly without custom frontend code?"]
objectives = ["Use LaTeX and Mermaid to express analysis logic and notation directly in lesson content.", "Use Hextra components to keep setup, navigation, and optional depth structured and maintainable."]
keypoints = ["Prefer built-in Hextra features over custom templates for diagrams and notation.", "Use tabs, details, steps, cards, and filetree to separate core path from optional depth."]
+++

This episode demonstrates documentation-focused Hextra features in a particle-physics context.
The goal is to keep source files readable for authors while giving learners clear visual structure.

## LaTeX for physics notation

Inline notation is useful for compact expressions like \(p_T > 25\,\mathrm{GeV}\), \(|\eta| < 2.4\), and \(\Delta R < 0.4\).

A standalone equation block can communicate selection logic and uncertainty estimates more clearly:

$$
\begin{aligned}
N_{\text{sig}} &= N_{\text{obs}} - N_{\text{bkg}} \\
Z &\approx \frac{N_{\text{sig}}}{\sqrt{N_{\text{bkg}} + (\delta N_{\text{bkg}})^2}}
\end{aligned}
$$

{{< callout type="note" title="Notation scope" >}}
Keep first-pass equations close to the learning objective.
Move detailed derivations to optional `details` blocks.
{{< /callout >}}

## Mermaid for analysis flow

```mermaid
flowchart LR
  A[Detector data] --> B[Reconstruction]
  B --> C[Quality filters]
  C --> D[Object selection]
  D --> E[Signal region]
  D --> F[Control region]
  E --> G[Fit and inference]
  F --> G
```

```mermaid
sequenceDiagram
  participant C as Collaborator
  participant CI as CI checks
  participant R as Reviewers
  C->>CI: Push selection update
  CI-->>C: Validation and build report
  C->>R: Open PR with plots and notes
  R-->>C: Request uncertainty clarification
```

## Synced tabs for command variants

Use repeated tab labels to keep shell choices synced across blocks.

{{< tabs >}}
{{< tab name="bash" selected=true >}}
```bash
hugo server
```
{{< /tab >}}
{{< tab name="zsh" >}}
```zsh
hugo server
```
{{< /tab >}}
{{< tab name="fish" >}}
```fish
hugo server
```
{{< /tab >}}
{{< /tabs >}}

{{< tabs >}}
{{< tab name="bash" selected=true >}}
```bash
go run github.com/oer-particle-physics/hugo-styles/cmd/hugo-styles-migrate@latest check .
```
{{< /tab >}}
{{< tab name="zsh" >}}
```zsh
go run github.com/oer-particle-physics/hugo-styles/cmd/hugo-styles-migrate@latest check .
```
{{< /tab >}}
{{< tab name="fish" >}}
```fish
go run github.com/oer-particle-physics/hugo-styles/cmd/hugo-styles-migrate@latest check .
```
{{< /tab >}}
{{< /tabs >}}

## FileTree and Cards for orientation

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
{{< card link="/docs/authoring" title="Authoring Guide" subtitle="Structure and pedagogy conventions." icon="book-open" tag="core" >}}
{{< card link="/docs/hextra-features" title="Hextra Features" subtitle="Curated feature usage for lesson authors." icon="cog" tag="core" >}}
{{< card link="/docs/updates" title="Update workflow" subtitle="How to pull shared module updates safely." icon="academic-cap" tag="maintenance" tagColor="blue" >}}
{{< card link="/reference" title="Reference links" subtitle="Upstream docs for deeper options." icon="book-open" tag="upstream" tagColor="gray" >}}
{{< /cards >}}

## Steps plus optional deep dive

{{% steps %}}

### Write the question and objective first

Anchor the episode around one concrete analysis skill.

### Add one core equation and one core diagram

Use LaTeX and Mermaid to keep explanation compact and precise.

### Keep setup variants in tabs

Avoid long repeated command blocks in linear prose.

### Hide advanced context in details blocks

Keep the default reading path short, then let motivated learners expand.

{{% /steps %}}

{{< details title="Optional: invariant-mass reminder" closed="true" >}}
For two reconstructed leptons, \(m_{\ell\ell}^2 = (E_1 + E_2)^2 - \lVert \vec{p}_1 + \vec{p}_2 \rVert^2\).  
Use this as context when explaining mass-window selection.
{{< /details >}}

{{< challenge title="Document a toy signal-selection flow" >}}
Add one Mermaid flowchart and one LaTeX equation to a new episode page that explains a toy signal-vs-background workflow.

{{< hint >}}
Start with only four stages: reconstruction, quality filters, signal region, control region.
{{< /hint >}}

{{< solution >}}
Keep the first version minimal: one flowchart, one equation for \(Z\) or a counting estimate, and one optional `details` block for background theory.
{{< /solution >}}
{{< /challenge >}}
