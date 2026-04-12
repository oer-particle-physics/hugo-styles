package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCollectFindingsValidHugoFixture(t *testing.T) {
	root := filepath.Join("testdata", "valid-hugo")
	findings, err := collectFindings(root)
	if err != nil {
		t.Fatalf("collectFindings returned error: %v", err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %#v", findings)
	}
}

func TestCollectFindingsInvalidHugoFixture(t *testing.T) {
	root := filepath.Join("testdata", "invalid-hugo")
	findings, err := collectFindings(root)
	if err != nil {
		t.Fatalf("collectFindings returned error: %v", err)
	}
	if len(findings) == 0 {
		t.Fatal("expected findings for invalid fixture")
	}

	kinds := map[string]bool{}
	for _, finding := range findings {
		kinds[finding.Kind] = true
	}

	for _, kind := range []string{"metadata", "duplicate-weight", "glossary-ref", "profile-ref", "image-alt", "legacy-attr-block"} {
		if !kinds[kind] {
			t.Fatalf("expected finding kind %q in %#v", kind, findings)
		}
	}
}

func TestCollectFindingsLegacyFixture(t *testing.T) {
	root := filepath.Join("testdata", "legacy-styles")
	findings, err := collectFindings(root)
	if err != nil {
		t.Fatalf("collectFindings returned error: %v", err)
	}
	if len(findings) == 0 {
		t.Fatal("expected legacy findings")
	}

	kinds := map[string]bool{}
	for _, finding := range findings {
		kinds[finding.Kind] = true
	}

	for _, kind := range []string{"liquid-include", "legacy-attr-block", "metadata"} {
		if !kinds[kind] {
			t.Fatalf("expected legacy finding kind %q in %#v", kind, findings)
		}
	}
}

func TestTransformMarkdownConvertsCommonPatterns(t *testing.T) {
	input := strings.TrimSpace(`
## Exercise heading
> Prompt text
{: .challenge}

<iframe src="https://www.youtube.com/embed/aqz-KE-bpKQ"></iframe>
![Diagram](../fig/diagram.svg)
<img src="../fig/example.png" alt="example">
`) + "\n"

	output := transformMarkdown(input)

	if !strings.Contains(output, `{{< challenge >}}`) {
		t.Fatalf("expected challenge shortcode conversion, got:\n%s", output)
	}
	if !strings.Contains(output, `Prompt text`) {
		t.Fatalf("expected challenge body conversion, got:\n%s", output)
	}
	if !strings.Contains(output, `{{< youtube aqz-KE-bpKQ >}}`) {
		t.Fatalf("expected youtube shortcode conversion, got:\n%s", output)
	}
	if !strings.Contains(output, `/fig/diagram.svg`) {
		t.Fatalf("expected markdown figure path rewrite, got:\n%s", output)
	}
	if !strings.Contains(output, `/fig/example.png`) {
		t.Fatalf("expected figure path rewrite, got:\n%s", output)
	}
}

func TestPromoteHeadingLevels(t *testing.T) {
	input := "# First\n## Second\n### Third\n```bash\n# untouched\n```\n"

	output := promoteHeadingLevels(input)

	if !strings.Contains(output, "## First") {
		t.Fatalf("expected level-1 heading promotion, got:\n%s", output)
	}
	if !strings.Contains(output, "### Second") {
		t.Fatalf("expected level-2 heading promotion, got:\n%s", output)
	}
	if !strings.Contains(output, "#### Third") {
		t.Fatalf("expected level-3 heading promotion, got:\n%s", output)
	}
	if !strings.Contains(output, "# untouched") {
		t.Fatalf("expected fenced code block to remain unchanged, got:\n%s", output)
	}
}

func TestTransformAndWriteNormalizesFrontMatter(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "source.md")
	dest := filepath.Join(dir, "out.md")
	err := os.WriteFile(src, []byte(strings.TrimSpace(`---
title: Lesson Title
layout: lesson
root: .
permalink: index.html
---
# Body
`)+"\n"), 0o644)
	if err != nil {
		t.Fatalf("write source: %v", err)
	}

	if err := transformAndWrite(src, dest, 10, true, true, false, "hextra-home"); err != nil {
		t.Fatalf("transformAndWrite returned error: %v", err)
	}

	output, err := os.ReadFile(dest)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}
	out := string(output)
	if !strings.HasPrefix(out, "+++\n") {
		t.Fatalf("expected TOML front matter, got:\n%s", out)
	}
	if !strings.Contains(out, `layout = 'hextra-home'`) {
		t.Fatalf("expected layout override, got:\n%s", out)
	}
	if !strings.Contains(out, "weight = 10") {
		t.Fatalf("expected weight, got:\n%s", out)
	}
	if !strings.Contains(out, "draft = true") {
		t.Fatalf("expected draft, got:\n%s", out)
	}
	if strings.Contains(out, "root") || strings.Contains(out, "permalink") {
		t.Fatalf("expected legacy keys removed, got:\n%s", out)
	}
}

func TestReadLegacyLessonTitle(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "_config.yml")
	err := os.WriteFile(configPath, []byte("title: Legacy Lesson\n"), 0o644)
	if err != nil {
		t.Fatalf("write config: %v", err)
	}

	title, err := readLegacyLessonTitle(dir)
	if err != nil {
		t.Fatalf("readLegacyLessonTitle returned error: %v", err)
	}
	if title != "Legacy Lesson" {
		t.Fatalf("expected legacy title, got %q", title)
	}
}

func TestReadLegacyLessonTitleFallsBackWhenConfigHasDuplicateKeys(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "_config.yml")
	err := os.WriteFile(configPath, []byte(strings.TrimSpace(`
title: Legacy Lesson
swc_site: https://example.com/one
swc_site: https://example.com/two
`)+"\n"), 0o644)
	if err != nil {
		t.Fatalf("write config: %v", err)
	}

	title, err := readLegacyLessonTitle(dir)
	if err != nil {
		t.Fatalf("readLegacyLessonTitle returned error: %v", err)
	}
	if title != "Legacy Lesson" {
		t.Fatalf("expected legacy title, got %q", title)
	}
}

func TestMigrateLessonHomePageBuildsStructuredLandingPage(t *testing.T) {
	source := t.TempDir()
	dest := t.TempDir()

	if err := os.MkdirAll(filepath.Join(source, "_episodes"), 0o755); err != nil {
		t.Fatalf("mkdir episodes: %v", err)
	}

	err := os.WriteFile(filepath.Join(source, "_config.yml"), []byte("title: Legacy Lesson\n"), 0o644)
	if err != nil {
		t.Fatalf("write config: %v", err)
	}
	err = os.WriteFile(filepath.Join(source, "_episodes", "01-first.md"), []byte("---\ntitle: First\n---\n"), 0o644)
	if err != nil {
		t.Fatalf("write episode: %v", err)
	}
	err = os.WriteFile(filepath.Join(source, "index.md"), []byte(strings.TrimSpace(`---
layout: lesson
root: .
permalink: index.html
---
Legacy lesson introduction.

Another lead paragraph.

> ## Prerequisites
>
> Familiarity with the shell.
{: .prereq}
`)+"\n"), 0o644)
	if err != nil {
		t.Fatalf("write homepage: %v", err)
	}

	if err := migrateLessonHomePage(source, dest); err != nil {
		t.Fatalf("migrateLessonHomePage returned error: %v", err)
	}

	outputPath := filepath.Join(dest, "content", "_index.md")
	output, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}

	meta, body, err := parseFrontMatter(string(output))
	if err != nil {
		t.Fatalf("parseFrontMatter returned error: %v", err)
	}
	if got := meta["title"]; got != "Legacy Lesson" {
		t.Fatalf("expected legacy title in front matter, got %#v", got)
	}
	if got := meta["layout"]; got != "hextra-home" {
		t.Fatalf("expected hextra-home layout, got %#v", got)
	}

	for _, snippet := range []string{
		"{{< hextra/hero-headline >}}",
		"Legacy Lesson",
		"{{< /hextra/hero-headline >}}",
		`{{< hextra/hero-button text="Start Lesson" link="episodes/01-first/" >}}`,
		"{{< lesson/overview >}}",
		`{{< lesson/schedule title="Schedule" >}}`,
		`{{< lesson/authors title="Authors and Contributors" >}}`,
		"## About This Lesson",
		"Legacy lesson introduction.",
		"Another lead paragraph.",
		`{{< callout type="prereq" title="Prerequisites" >}}`,
	} {
		if !strings.Contains(body, snippet) {
			t.Fatalf("expected homepage body to contain %q, got:\n%s", snippet, body)
		}
	}
}
