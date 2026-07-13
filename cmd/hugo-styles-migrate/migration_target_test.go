package main

import (
	"errors"
	"flag"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidateEpisodeDocRejectsMalformedMetadata(t *testing.T) {
	doc := contentDoc{
		Path: "episode.md",
		Meta: map[string]any{
			"title":      42,
			"weight":     10.5,
			"questions":  []any{""},
			"objectives": []any{"valid", 2},
			"keypoints":  []any{},
			"teaching":   -1,
			"exercises":  "5",
		},
	}

	findings, _, hasWeight := validateEpisodeDoc(doc)
	if hasWeight {
		t.Fatal("a floating-point weight must not be accepted as an integer")
	}
	if len(findings) != 7 {
		t.Fatalf("expected seven metadata findings, got %#v", findings)
	}
}

func TestLessonImageRequiresSourceAndAltText(t *testing.T) {
	doc := contentDoc{
		Path: "page.md",
		Body: `{{< lesson/image alt="" >}}
{{< lesson/image src="/fig/example.svg" >}}`,
	}

	findings := missingAltTextFindings(doc)
	if len(findings) != 3 {
		t.Fatalf("expected missing source and alt findings, got %#v", findings)
	}
}

func TestCommandHelpReturnsFlagErrHelp(t *testing.T) {
	if err := runCheck([]string{"--help"}); !errors.Is(err, flag.ErrHelp) {
		t.Fatalf("check --help returned %v", err)
	}
	if err := runMigrate([]string{"--help"}); !errors.Is(err, flag.ErrHelp) {
		t.Fatalf("migrate --help returned %v", err)
	}
}

func TestMigrationDryRunDoesNotChangeTemplate(t *testing.T) {
	source, dest := migrationFixture(t)
	before := mustReadTestFile(t, filepath.Join(dest, "content", "_index.md"))

	if err := migrateIntoTemplate(source, dest, true); err != nil {
		t.Fatalf("dry run returned error: %v", err)
	}
	if after := mustReadTestFile(t, filepath.Join(dest, "content", "_index.md")); after != before {
		t.Fatalf("dry run changed the destination homepage: %q", after)
	}
	if _, err := os.Stat(filepath.Join(dest, "content", "episodes", "01-first", "index.md")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("dry run created migrated content: %v", err)
	}
}

func TestMigrationReplacesManagedContentAndPreservesInfrastructure(t *testing.T) {
	source, dest := migrationFixture(t)

	if err := migrateIntoTemplate(source, dest, false); err != nil {
		t.Fatalf("migration returned error: %v", err)
	}

	for _, path := range []string{
		"content/episodes/01-first/index.md",
		"content/learners/setup.md",
		"content/instructors/instructor-notes.md",
		"content/reference.md",
		"static/fig/new.svg",
		"AUTHORS",
	} {
		if _, err := os.Stat(filepath.Join(dest, filepath.FromSlash(path))); err != nil {
			t.Errorf("expected migrated path %s: %v", path, err)
		}
	}
	for _, path := range []string{
		"content/episodes/old/index.md",
		"content/learners/old.md",
		"static/fig/old.svg",
	} {
		if _, err := os.Stat(filepath.Join(dest, filepath.FromSlash(path))); !errors.Is(err, os.ErrNotExist) {
			t.Errorf("expected managed path %s to be removed, got %v", path, err)
		}
	}
	for path, expected := range map[string]string{
		"content/episodes/_index.md":       "episode index\n",
		"content/learners/_index.md":       "learner index\n",
		"content/all-in-one.md":            "generated resource\n",
		".github/workflows/pages.yml":      "workflow\n",
		"static/favicon.svg":               "branding\n",
		"content/docs/template-guide.md":   "template documentation\n",
		"layouts/shortcodes/template.html": "template layout\n",
	} {
		if actual := mustReadTestFile(t, filepath.Join(dest, filepath.FromSlash(path))); actual != expected {
			t.Errorf("preserved path %s changed: %q", path, actual)
		}
	}
}

func TestMigrationRejectsDirtyUnrecognizedAndOverlappingDestinations(t *testing.T) {
	source, dest := migrationFixture(t)
	writeTestFile(t, filepath.Join(dest, "dirty.txt"), "dirty\n")
	if err := migrateIntoTemplate(source, dest, true); err == nil || !strings.Contains(err.Error(), "must be clean") {
		t.Fatalf("expected dirty-worktree error, got %v", err)
	}

	other := t.TempDir()
	writeTestFile(t, filepath.Join(other, "_episodes", "episode.md"), "---\ntitle: Test\n---\n")
	if _, _, err := validateMigrationRoots(other, filepath.Join(other, "_episodes")); err == nil || !strings.Contains(err.Error(), "must not overlap") {
		t.Fatalf("expected overlap error, got %v", err)
	}

	unrecognized := t.TempDir()
	writeTestFile(t, filepath.Join(unrecognized, "hugo.toml"), "title = 'Other'\n")
	writeTestFile(t, filepath.Join(unrecognized, "go.mod"), "module example.test/other\n")
	initTestGitRepo(t, unrecognized)
	if err := validateTemplateDestination(unrecognized); err == nil || !strings.Contains(err.Error(), "not recognized") {
		t.Fatalf("expected template-recognition error, got %v", err)
	}
}

func TestTemplateRecognitionIgnoresCommentOnlyModuleMentions(t *testing.T) {
	dest := t.TempDir()
	writeTestFile(t, filepath.Join(dest, "hugo.toml"), "# "+hugoStylesModule+"\n[module]\n")
	writeTestFile(t, filepath.Join(dest, "go.mod"), "module example.test/lesson\n\n// require "+hugoStylesModule+" v0.4.0\n")
	initTestGitRepo(t, dest)

	if err := validateTemplateDestination(dest); err == nil || !strings.Contains(err.Error(), "not recognized") {
		t.Fatalf("expected comment-only module mentions to be rejected, got %v", err)
	}
}

func TestApplyMigrationChangesRestoresDestinationOnFailure(t *testing.T) {
	root := t.TempDir()
	staging := filepath.Join(root, "staging")
	dest := filepath.Join(root, "dest")
	writeTestFile(t, filepath.Join(staging, "one.txt"), "new\n")
	writeTestFile(t, filepath.Join(staging, "blocked", "two.txt"), "new two\n")
	writeTestFile(t, filepath.Join(dest, "one.txt"), "old\n")
	writeTestFile(t, filepath.Join(dest, "blocked"), "not a directory\n")

	changes := []migrationChange{
		{RelPath: "one.txt", Source: filepath.Join(staging, "one.txt"), Target: filepath.Join(dest, "one.txt")},
		{RelPath: filepath.Join("blocked", "two.txt"), Source: filepath.Join(staging, "blocked", "two.txt"), Target: filepath.Join(dest, "blocked", "two.txt")},
	}
	if err := applyMigrationChanges(changes, staging, dest); err == nil {
		t.Fatal("expected transactional install failure")
	}
	if got := mustReadTestFile(t, filepath.Join(dest, "one.txt")); got != "old\n" {
		t.Fatalf("destination was not restored: %q", got)
	}
	if got := mustReadTestFile(t, filepath.Join(dest, "blocked")); got != "not a directory\n" {
		t.Fatalf("unmanaged blocker changed: %q", got)
	}
}

func migrationFixture(t *testing.T) (string, string) {
	t.Helper()
	root := t.TempDir()
	source := filepath.Join(root, "legacy")
	dest := filepath.Join(root, "template")

	writeTestFile(t, filepath.Join(source, "_config.yml"), "title: Migrated lesson\n")
	writeTestFile(t, filepath.Join(source, "index.md"), "---\nlayout: lesson\n---\nLegacy introduction.\n")
	writeTestFile(t, filepath.Join(source, "_episodes", "01-first.md"), strings.TrimSpace(`---
title: First
questions:
  - What is migrated?
objectives:
  - Migrate safely.
keypoints:
  - Templates remain intact.
teaching: 5
exercises: 5
---
# First topic
`)+"\n")
	writeTestFile(t, filepath.Join(source, "setup.md"), "---\ntitle: Setup\n---\nSetup.\n")
	writeTestFile(t, filepath.Join(source, "reference.md"), "---\ntitle: Reference\n---\nReference.\n")
	writeTestFile(t, filepath.Join(source, "_extras", "guide.md"), "---\ntitle: Guide\n---\nGuide.\n")
	writeTestFile(t, filepath.Join(source, "fig", "new.svg"), "new figure\n")
	writeTestFile(t, filepath.Join(source, "AUTHORS"), "Legacy Author\n")

	writeTestFile(t, filepath.Join(dest, "hugo.toml"), "[module]\n[[module.imports]]\npath = '"+hugoStylesModule+"'\n")
	writeTestFile(t, filepath.Join(dest, "go.mod"), "module example.test/lesson\n\nrequire "+hugoStylesModule+" v0.4.0\n")
	writeTestFile(t, filepath.Join(dest, "content", "_index.md"), "old homepage\n")
	writeTestFile(t, filepath.Join(dest, "content", "episodes", "_index.md"), "episode index\n")
	writeTestFile(t, filepath.Join(dest, "content", "episodes", "old", "index.md"), "old episode\n")
	writeTestFile(t, filepath.Join(dest, "content", "learners", "_index.md"), "learner index\n")
	writeTestFile(t, filepath.Join(dest, "content", "learners", "old.md"), "old learner page\n")
	writeTestFile(t, filepath.Join(dest, "content", "instructors", "_index.md"), "instructor index\n")
	writeTestFile(t, filepath.Join(dest, "content", "glossary", "_index.md"), "glossary index\n")
	writeTestFile(t, filepath.Join(dest, "content", "profiles", "_index.md"), "profile index\n")
	writeTestFile(t, filepath.Join(dest, "content", "all-in-one.md"), "generated resource\n")
	writeTestFile(t, filepath.Join(dest, "content", "docs", "template-guide.md"), "template documentation\n")
	writeTestFile(t, filepath.Join(dest, "layouts", "shortcodes", "template.html"), "template layout\n")
	writeTestFile(t, filepath.Join(dest, ".github", "workflows", "pages.yml"), "workflow\n")
	writeTestFile(t, filepath.Join(dest, "static", "favicon.svg"), "branding\n")
	writeTestFile(t, filepath.Join(dest, "static", "fig", "old.svg"), "old figure\n")
	initTestGitRepo(t, dest)
	return source, dest
}

func writeTestFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("create parent for %s: %v", path, err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func mustReadTestFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(data)
}

func initTestGitRepo(t *testing.T, root string) {
	t.Helper()
	for _, args := range [][]string{
		{"init", "--quiet", root},
		{"-C", root, "config", "user.email", "tests@example.test"},
		{"-C", root, "config", "user.name", "Tests"},
		{"-C", root, "add", "."},
		{"-C", root, "commit", "--quiet", "-m", "fixture"},
	} {
		if output, err := exec.Command("git", args...).CombinedOutput(); err != nil {
			t.Fatalf("git %s: %v: %s", strings.Join(args, " "), err, output)
		}
	}
}
