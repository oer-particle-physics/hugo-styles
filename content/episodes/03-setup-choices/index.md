+++
title = "Setup choices, profiles, and native embeds"
weight = 30
teaching = 15
exercises = 5
questions = ["How do we document setup and learner variation without overloading the main narrative?"]
objectives = ["Show how profile pages and glossary pages support lesson navigation.", "Use Hugo-native embed syntax instead of hand-written iframe HTML when a provider shortcode exists."]
keypoints = ["Profiles and glossary entries are easier to maintain as real content pages.", "Supported video providers should use Hugo shortcodes rather than raw iframe HTML."]
[tabs]
  sync = false
+++

Profiles make it easier to explain why a lesson exists for different audiences. In this example, the {{< profile workshop-host >}} profile cares most about pacing, while the {{< profile self-paced-maintainer >}} profile cares more about update flow and repo ownership.

For terminology, content can link to short definitions like {{< glossary audience-toggle >}} without interrupting the main narrative.

![A facilitation sketch showing where setup decisions and audience mode intersect.](fig/facilitation-map.svg)

## Tabs for local setup choices

This page opts out of synced tabs so the examples can stay local to this episode. That is useful when one tab set is about operating systems and another is about shells, but you do not want changing one to affect the other.

{{< tabs >}}
{{< tab name="macOS" selected=true >}}
Tell learners where a tool lives in Finder or which package manager you expect them to use.
{{< /tab >}}
{{< tab name="Linux" >}}
Tell learners which package name or distribution-specific prerequisite matters.
{{< /tab >}}
{{< tab name="Windows" >}}
Tell learners whether PowerShell, WSL, or a native installer is the supported route.
{{< /tab >}}
{{< /tabs >}}

{{< tabs >}}
{{< tab name="bash" selected=true >}}
```bash
export LESSON_ENV=workshop
```
{{< /tab >}}
{{< tab name="zsh" >}}
```zsh
export LESSON_ENV=workshop
```
{{< /tab >}}
{{< tab name="fish" >}}
```fish
set -x LESSON_ENV workshop
```
{{< /tab >}}
{{< /tabs >}}

Because this page disables sync, the second shell tab group below starts independently instead of following the choice above.

{{< tabs >}}
{{< tab name="bash" selected=true >}}
```bash
echo $LESSON_ENV
```
{{< /tab >}}
{{< tab name="zsh" >}}
```zsh
echo $LESSON_ENV
```
{{< /tab >}}
{{< tab name="fish" >}}
```fish
echo $LESSON_ENV
```
{{< /tab >}}
{{< /tabs >}}

Below is a native Hugo video embed example:

{{< youtube aqz-KE-bpKQ >}}

{{< callout type="testimonial" title="Migration expectation" >}}
The migration helper should cover the common path well, but highly customised legacy sites should still expect a short manual pass.
{{< /callout >}}
