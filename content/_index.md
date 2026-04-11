+++
title = "Hugo Styles"
layout = "hextra-home"
+++

{{< hextra/hero-badge link="docs/quickstart/" >}}
Shared Module + Example Site {{< icon name="arrow-circle-right" attributes="height=14" >}}
{{< /hextra/hero-badge >}}

<div class="hx:mt-6 hx:mb-6">
{{< hextra/hero-headline >}}
Build Carpentries-style lessons&nbsp;
<br class="hx:xl:block hx:hidden" />with Hugo and Hextra
{{< /hextra/hero-headline >}}
</div>

<div class="hx:mb-12">
{{< hextra/hero-subtitle >}}
Keep the teaching structure people already know,&nbsp;
<br class="hx:lg:block hx:hidden" />with a Hugo-first stack that is easier to update and teach from.
{{< /hextra/hero-subtitle >}}
</div>

<div class="hx:mb-6">
{{< hextra/hero-button text="Get Started" link="docs/quickstart/" >}}
</div>

<div class="hx:mt-6"></div>

{{< hextra/feature-grid cols="3" >}}
{{< hextra/feature-card
  title="Carpentries pedagogy"
  subtitle="Challenge blocks, collapsible solutions, learner and instructor views, glossary links, profiles, and aggregated teaching resources are built into the module."
  icon="academic-cap"
  link="docs/components/"
  class="hx:min-h-[220px]"
>}}
{{< hextra/feature-card
  title="Thin downstream repos"
  subtitle="Lesson repositories import this module instead of copying layouts and shortcodes, which keeps updates small, predictable, and easier to review."
  icon="book-open"
  link="episodes/01-introduction/"
  class="hx:min-h-[220px]"
>}}
{{< hextra/feature-card
  title="GitHub Pages ready"
  subtitle="The starter template deploys with GitHub Actions, validates lesson content in CI, and can receive module updates through Dependabot pull requests."
  icon="sparkles"
  link="docs/updates/"
  class="hx:min-h-[220px]"
>}}
{{< /hextra/feature-grid >}}
