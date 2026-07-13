+++
title = "Components"
weight = 40
+++

This is the live reference for features owned by `hugo-styles`. Each example is rendered here; open the adjacent **Use it** block to copy its source. Theme-level features live in the [Hextra feature guide]({{< relref "/docs/hextra-features" >}}).

## Challenge, hint, and solution

{{< challenge title="Choose a reproducible cut" subtitle="Keep the first pass small." >}}
Pick one event-selection threshold and explain what it removes.

{{< hint >}}
Start with a transverse-momentum threshold.
{{< /hint >}}

{{< solution >}}
For this toy sample, use \(p_T > 25\,\mathrm{GeV}\), then record the value next to the result.
{{< /solution >}}
{{< /challenge >}}

{{< details title="Use it" closed="true" >}}
{{< codeblock lang="text" >}}
{{</* challenge title="Choose a reproducible cut" subtitle="Keep the first pass small." */>}}
Pick one event-selection threshold and explain what it removes.

{{</* hint */>}}Start with a transverse-momentum threshold.{{</* /hint */>}}
{{</* solution */>}}Use `p_T > 25 GeV` and record it.{{</* /solution */>}}
{{</* /challenge */>}}
{{< /codeblock >}}
{{< /details >}}

## Callout palette

{{< callout type="note" title="Note" >}}Supporting context that belongs near the main path.{{< /callout >}}
{{< callout type="prereq" title="Prerequisite" >}}What learners need before starting.{{< /callout >}}
{{< callout type="checklist" title="Checklist" >}}A short set of checks before moving on.{{< /callout >}}
{{< callout type="discussion" title="Discussion" >}}A prompt for comparing approaches.{{< /callout >}}
{{< callout type="testimonial" title="Instructor perspective" >}}A brief experience or facilitation observation.{{< /callout >}}
{{< callout type="warning" title="Warning" >}}A likely mistake with meaningful consequences.{{< /callout >}}
{{< callout type="caution" title="Caution" >}}A migration or workflow edge case worth checking.{{< /callout >}}

{{< details title="Use it" closed="true" >}}
{{< codeblock lang="text" >}}
{{</* callout type="note" title="Note" */>}}Supporting context.{{</* /callout */>}}
{{</* callout type="prereq" title="Prerequisite" */>}}Required setup.{{</* /callout */>}}
{{</* callout type="checklist" title="Checklist" */>}}Checks to complete.{{</* /callout */>}}
{{</* callout type="discussion" title="Discussion" */>}}Compare approaches.{{</* /callout */>}}
{{</* callout type="testimonial" title="Instructor perspective" */>}}An experience.{{</* /callout */>}}
{{</* callout type="warning" title="Warning" */>}}A likely mistake.{{</* /callout */>}}
{{</* callout type="caution" title="Caution" */>}}An edge case.{{</* /callout */>}}
{{< /codeblock >}}
{{< /details >}}

## Audience-aware content

Supporting pages can place the selector exactly where readers need it. Episode and All-in-One pages include it automatically.

{{< lesson/audience-toggle >}}

{{< learner >}}
Learners see the task, expected inputs, and enough context to proceed.
{{< /learner >}}

{{< instructor >}}
Instructors see pacing notes, debrief prompts, and facilitation context.
{{< /instructor >}}

{{< details title="Use it" closed="true" >}}
{{< codeblock lang="text" >}}
{{</* lesson/audience-toggle */>}}

{{</* learner */>}}Learner-facing guidance.{{</* /learner */>}}
{{</* instructor */>}}Facilitation guidance.{{</* /instructor */>}}
{{< /codeblock >}}

Set the initial view in `hugo.toml`:

```toml
[params.lesson]
  defaultView = "learner" # or "instructor"
```

The precedence is a valid `?view=` value, a saved choice, this configured default, then `learner`.
{{< /details >}}

## Glossary and profile links

Inline references resolve to normal content pages: {{< glossary formative-assessment >}} and {{< profile workshop-host >}}.

{{< details title="Use it" closed="true" >}}
{{< codeblock lang="text" >}}
{{</* glossary formative-assessment */>}}
{{</* profile workshop-host */>}}
{{< /codeblock >}}
{{< /details >}}

## Lesson image

`lesson/image` accepts page-bundle resources and the shared `static/{fig,files,data,code}` paths. It requires `src` and meaningful `alt` text and follows the global or per-page image-zoom setting.

{{< lesson/image src="/fig/lesson-flow.svg" alt="Lesson content flowing into facilitation and reusable resources." width="560" class="hx:rounded-xl" >}}

{{< details title="Use it" closed="true" >}}
{{< codeblock lang="text" >}}
{{</* lesson/image
  src="/fig/lesson-flow.svg"
  alt="Lesson content flowing into facilitation and reusable resources."
  width="560"
  class="hx:rounded-xl"
*/>}}
{{< /codeblock >}}

`style` remains available for legacy content, but prefer `width` and `class`. Disable zoom for one page with `imageZoom = false` in its front matter.
{{< /details >}}

## Live lesson metadata

The configured lesson title is **{{< lesson/meta "title" >}}**. This value comes directly from `params.lesson`, so content does not need to duplicate it.

{{< details title="Use it" closed="true" >}}
{{< codeblock lang="text" >}}
{{</* lesson/meta "title" */>}}
{{</* lesson/meta "tagline" */>}}
{{</* lesson/meta "description" */>}}
{{</* lesson/meta "siteTitle" */>}}
{{< /codeblock >}}
{{< /details >}}

## Generated lesson resources

The homepage already renders the [episode overview, schedule, and authors]({{< relref "/" >}}). The same episode metadata also drives:

- [All-in-One]({{< relref "/all-in-one" >}})
- [Key Points]({{< relref "/key-points" >}})
- [External Links]({{< relref "/external-links" >}})
- [All Images]({{< relref "/extract-all-images" >}})

{{< details title="Use it on a homepage" closed="true" >}}
{{< codeblock lang="text" >}}
{{</* lesson/overview */>}}
{{</* lesson/schedule title="Schedule" */>}}
{{</* lesson/authors title="Authors and Contributors" */>}}
{{< /codeblock >}}
{{< /details >}}
