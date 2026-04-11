+++
title = "Using challenge and solution blocks"
weight = 20
teaching = 15
exercises = 10
questions = ["How should active-learning blocks behave in a Hugo-based lesson?"]
objectives = ["Use challenge, hint, and solution shortcodes in a way that feels natural to authors.", "Show how learner and instructor context can coexist in the same episode."]
keypoints = ["Challenge blocks should be easy to read in source form and easy to scan on the page.", "Hints and solutions should stay collapsible but also support bulk expansion in the all-in-one page."]
+++

Shortcodes should help the author express teaching intent, not force them to think about page mechanics.

{{< challenge title="Draft a reusable exercise block" >}}
Write a short prompt that asks contributors to identify which parts of a lesson framework should remain shared across repositories.

{{< hint >}}
Think about pieces that change rarely and are expensive to copy around.
{{< /hint >}}

{{< solution >}}
Reusable framework elements usually include layouts, shortcodes, shared CSS, small interaction scripts, and validation logic. Lesson-specific prose and branding should normally stay local to each tutorial repository.
{{< /solution >}}
{{< /challenge >}}

{{< callout type="checklist" title="Review checklist" >}}
- Does the activity ask learners to do something concrete?
- Is the expected answer small enough for a workshop?
- Could an instructor reveal the solution only when needed?
{{< /callout >}}

{{< details title="Why not keep legacy fenced-attribute syntax?" closed="true" >}}
It is harder to teach to new contributors, less portable, and more awkward to validate than a small set of explicit shortcodes.
{{< /details >}}

{{< instructor >}}
If you need to pace the room, ask learners to compare answers in pairs before expanding the solution.
{{< /instructor >}}

## Callout gallery

The older `styles` lessons relied heavily on recognisable teaching boxes. In this stack, they are still first-class teaching components, but their presentation stays closer to Hextra's overall visual system.
Expand the `Show code` toggles to see the shortcode syntax for each box without making the page much longer.

{{< callout type="note" title="Note" >}}
Use this for side context or reminders that support the main narrative.
{{< /callout >}}
{{< details title="Show code (note)" closed="true" >}}
{{< codeblock lang="text" >}}
{{</* callout type="note" title="Note" */>}}
Use this for side context or reminders that support the main narrative.
{{</* /callout */>}}
{{< /codeblock >}}
{{< /details >}}

{{< callout type="prereq" title="Prerequisite" >}}
Use this when learners need knowledge or setup in place before continuing.
{{< /callout >}}
{{< details title="Show code (prereq)" closed="true" >}}
{{< codeblock lang="text" >}}
{{</* callout type="prereq" title="Prerequisite" */>}}
Use this when learners need knowledge or setup in place before continuing.
{{</* /callout */>}}
{{< /codeblock >}}
{{< /details >}}

{{< callout type="checklist" title="Checklist" >}}
- Confirm the prompt is concrete.
- Keep the exercise small enough for a workshop.
- Make the reveal optional.
{{< /callout >}}
{{< details title="Show code (checklist)" closed="true" >}}
{{< codeblock lang="text" >}}
{{</* callout type="checklist" title="Checklist" */>}}
- Confirm the prompt is concrete.
- Keep the exercise small enough for a workshop.
- Make the reveal optional.
{{</* /callout */>}}
{{< /codeblock >}}
{{< /details >}}

{{< callout type="discussion" title="Discussion" >}}
Ask learners to compare trade-offs before showing a canonical answer.
{{< /callout >}}
{{< details title="Show code (discussion)" closed="true" >}}
{{< codeblock lang="text" >}}
{{</* callout type="discussion" title="Discussion" */>}}
Ask learners to compare trade-offs before showing a canonical answer.
{{</* /callout */>}}
{{< /codeblock >}}
{{< /details >}}

{{< callout type="testimonial" title="Instructor perspective" >}}
“This kind of box is useful when I want to signal facilitation intent without breaking the flow of the lesson.”
{{< /callout >}}
{{< details title="Show code (testimonial)" closed="true" >}}
{{< codeblock lang="text" >}}
{{</* callout type="testimonial" title="Instructor perspective" */>}}
“This kind of box is useful when I want to signal facilitation intent without breaking the flow of the lesson.”
{{</* /callout */>}}
{{< /codeblock >}}
{{< /details >}}

{{< callout type="warning" title="Common trap" >}}
Avoid using a generic callout when the content is actually an exercise or a solution. Those need their own behaviour.
{{< /callout >}}
{{< details title="Show code (warning)" closed="true" >}}
{{< codeblock lang="text" >}}
{{</* callout type="warning" title="Common trap" */>}}
Avoid using a generic callout when the content is actually an exercise or a solution. Those need their own behaviour.
{{</* /callout */>}}
{{< /codeblock >}}
{{< /details >}}

{{< callout type="caution" title="Migration caution" >}}
Legacy lessons may contain nested blockquote syntax that should be migrated automatically, but unusual customisations still deserve a quick manual pass.
{{< /callout >}}
{{< details title="Show code (caution)" closed="true" >}}
{{< codeblock lang="text" >}}
{{</* callout type="caution" title="Migration caution" */>}}
Legacy lessons may contain nested blockquote syntax that should be migrated automatically, but unusual customisations still deserve a quick manual pass.
{{</* /callout */>}}
{{< /codeblock >}}
{{< /details >}}

{{< learner >}}
As a learner, you should mostly see activities, guidance, and supporting notes.
{{< /learner >}}

{{< instructor >}}
As an instructor, this page can also surface pacing advice, debrief prompts, and workshop-specific facilitation notes without requiring a separate fork of the lesson.
{{< /instructor >}}
