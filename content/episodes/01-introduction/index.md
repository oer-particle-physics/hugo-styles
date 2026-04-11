+++
title = "Designing a Hugo-first lesson"
weight = 10
teaching = 10
exercises = 5
questions = ["What should a modern lesson template preserve from the old Carpentries stack?"]
objectives = ["Identify the teaching features that matter more than the old implementation details.", "Recognise which pieces should live in a shared module versus in a lesson repository."]
keypoints = ["Preserve pedagogy and author ergonomics, not the historical implementation.", "A thin lesson repo plus a shared module gives a much cleaner update path."]
+++

The lesson infrastructure in this repository is intentionally split in two:

- a **shared module** that owns the reusable lesson system
- a **starter template** that stays light and project-specific

That means people can keep their tutorials current without copying framework files across repositories.

{{< callout type="prereq" title="Who this is for" >}}
This stack is aimed at lesson authors who want Carpentries-style pedagogy in a Hugo-native workflow.
{{< /callout >}}

The module keeps the teaching model front and centre. For example, we can link directly to {{< glossary formative-assessment >}} practices and connect activities to a target learner profile such as {{< profile workshop-host >}}.

![An example lesson map that connects content, facilitation, and update flow.](fig/lesson-flow.svg)

{{< learner >}}
As you read the example lesson, look for the places where metadata becomes visible structure: questions, objectives, key points, and active-learning prompts.
{{< /learner >}}

{{< instructor >}}
This first episode is a good place to explain the module/template split before diving into syntax. Learners usually care less about the build system than maintainers do.
{{< /instructor >}}

