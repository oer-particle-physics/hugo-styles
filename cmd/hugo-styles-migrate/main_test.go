package main

import (
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
}
