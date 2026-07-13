package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	toml "github.com/pelletier/go-toml/v2"
	"gopkg.in/yaml.v3"
)

type contentDoc struct {
	Path     string
	RelPath  string
	Section  string
	BaseName string
	Body     string
	Meta     map[string]any
}

var (
	markdownImagePattern  = regexp.MustCompile(`!\[([^\]]*)\]\(([^)]+)\)`)
	htmlImagePattern      = regexp.MustCompile(`(?is)<img\b([^>]*?)>`)
	lessonImagePattern    = regexp.MustCompile(`(?is)\{\{[<%]\s*lesson/image\b([^}%]*?)[>%]\}\}`)
	glossaryShortcodeExpr = regexp.MustCompile(`\{\{<\s*glossary\s+(?:"([^"]+)"|([A-Za-z0-9._/-]+))`)
	profileShortcodeExpr  = regexp.MustCompile(`\{\{<\s*profile\s+(?:"([^"]+)"|([A-Za-z0-9._/-]+))`)
)

func collectFindings(root string) ([]finding, error) {
	findings, err := collectLegacyFindings(root)
	if err != nil {
		return nil, err
	}

	contentFindings, err := collectContentFindings(root)
	if err != nil {
		return nil, err
	}
	findings = append(findings, contentFindings...)

	sort.Slice(findings, func(i, j int) bool {
		if findings[i].Path == findings[j].Path {
			if findings[i].Kind == findings[j].Kind {
				return findings[i].Message < findings[j].Message
			}
			return findings[i].Kind < findings[j].Kind
		}
		return findings[i].Path < findings[j].Path
	})
	return findings, nil
}

func collectContentFindings(root string) ([]finding, error) {
	contentRoot := filepath.Join(root, "content")
	if stat, err := os.Stat(contentRoot); err != nil || !stat.IsDir() {
		if errors.Is(err, os.ErrNotExist) || err == nil {
			return nil, nil
		}
		return nil, err
	}

	docs, findings, err := loadContentDocs(contentRoot)
	if err != nil {
		return nil, err
	}

	glossarySlugs := collectSectionSlugs(docs, "glossary")
	profileSlugs := collectSectionSlugs(docs, "profiles")
	episodeWeights := map[int][]string{}

	for _, doc := range docs {
		findings = append(findings, unresolvedShortcodeFindings(doc, glossarySlugs, profileSlugs)...)
		findings = append(findings, missingAltTextFindings(doc)...)

		if !isEpisodeDoc(doc) {
			continue
		}
		episodeFindings, weight, hasWeight := validateEpisodeDoc(doc)
		findings = append(findings, episodeFindings...)
		if hasWeight {
			episodeWeights[weight] = append(episodeWeights[weight], doc.Path)
		}
	}

	for weight, paths := range episodeWeights {
		if len(paths) < 2 {
			continue
		}
		sort.Strings(paths)
		for _, path := range paths {
			findings = append(findings, finding{
				Path:    path,
				Kind:    "duplicate-weight",
				Message: fmt.Sprintf("episode weight %d is used multiple times", weight),
			})
		}
	}

	return findings, nil
}

func loadContentDocs(contentRoot string) ([]contentDoc, []finding, error) {
	var docs []contentDoc
	var findings []finding

	err := filepath.WalkDir(contentRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			switch d.Name() {
			case ".git", "node_modules", "public", "resources":
				return filepath.SkipDir
			}
			return nil
		}
		if filepath.Ext(path) != ".md" {
			return nil
		}

		relPath, err := filepath.Rel(contentRoot, path)
		if err != nil {
			return err
		}
		relPath = filepath.ToSlash(relPath)

		raw, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		meta, body, err := parseFrontMatter(string(raw))
		if err != nil {
			findings = append(findings, finding{
				Path:    path,
				Kind:    "frontmatter",
				Message: err.Error(),
			})
			meta = map[string]any{}
			body = string(raw)
		}

		section := ""
		parts := strings.Split(relPath, "/")
		if len(parts) > 0 {
			section = parts[0]
		}
		docs = append(docs, contentDoc{
			Path:     path,
			RelPath:  relPath,
			Section:  section,
			BaseName: filepath.Base(path),
			Body:     body,
			Meta:     meta,
		})
		return nil
	})
	return docs, findings, err
}

func parseFrontMatter(text string) (map[string]any, string, error) {
	text = strings.TrimPrefix(text, "\ufeff")
	switch {
	case strings.HasPrefix(text, "---\n"):
		front, body, ok := splitDelimitedFrontMatter(text, "---")
		if !ok {
			return nil, "", errors.New("invalid YAML front matter")
		}
		var meta map[string]any
		if err := yaml.Unmarshal([]byte(front), &meta); err != nil {
			return nil, "", fmt.Errorf("invalid YAML front matter: %w", err)
		}
		if meta == nil {
			meta = map[string]any{}
		}
		return meta, body, nil
	case strings.HasPrefix(text, "+++\n"):
		front, body, ok := splitDelimitedFrontMatter(text, "+++")
		if !ok {
			return nil, "", errors.New("invalid TOML front matter")
		}
		var meta map[string]any
		if err := toml.Unmarshal([]byte(front), &meta); err != nil {
			return nil, "", fmt.Errorf("invalid TOML front matter: %w", err)
		}
		if meta == nil {
			meta = map[string]any{}
		}
		return meta, body, nil
	default:
		return map[string]any{}, text, nil
	}
}

func splitDelimitedFrontMatter(text, delimiter string) (string, string, bool) {
	open := delimiter + "\n"
	close := "\n" + delimiter + "\n"
	if !strings.HasPrefix(text, open) {
		return "", "", false
	}
	idx := strings.Index(text[len(open):], close)
	if idx == -1 {
		return "", "", false
	}
	idx += len(open)
	return text[len(open):idx], text[idx+len(close):], true
}

func collectSectionSlugs(docs []contentDoc, section string) map[string]struct{} {
	slugs := map[string]struct{}{}
	for _, doc := range docs {
		if doc.Section != section || doc.BaseName == "_index.md" {
			continue
		}
		slug := docSlug(doc)
		if slug != "" {
			slugs[slug] = struct{}{}
		}
	}
	return slugs
}

func docSlug(doc contentDoc) string {
	if doc.BaseName == "index.md" {
		return filepath.Base(filepath.Dir(doc.Path))
	}
	return strings.TrimSuffix(doc.BaseName, filepath.Ext(doc.BaseName))
}

func isEpisodeDoc(doc contentDoc) bool {
	return doc.Section == "episodes" && doc.BaseName != "_index.md"
}

func validateEpisodeDoc(doc contentDoc) ([]finding, int, bool) {
	var findings []finding
	if title, ok := doc.Meta["title"].(string); !ok || strings.TrimSpace(title) == "" {
		findings = append(findings, metadataFinding(doc.Path, "episode title must be a non-empty string"))
	}
	for _, key := range []string{"questions", "objectives", "keypoints"} {
		if !nonEmptyStringList(doc.Meta[key]) {
			findings = append(findings, metadataFinding(
				doc.Path,
				fmt.Sprintf("episode %s must be a non-empty list of non-empty strings", key),
			))
		}
	}

	weight, ok := strictIntValue(doc.Meta["weight"])
	if !ok {
		findings = append(findings, metadataFinding(doc.Path, "episode weight must be an integer"))
	}
	for _, key := range []string{"teaching", "exercises"} {
		value, present := doc.Meta[key]
		if !present {
			continue
		}
		minutes, valid := strictIntValue(value)
		if !valid || minutes < 0 {
			findings = append(findings, metadataFinding(
				doc.Path,
				fmt.Sprintf("episode %s must be a non-negative integer", key),
			))
		}
	}
	if !ok {
		return findings, 0, false
	}
	return findings, weight, true
}

func metadataFinding(path, message string) finding {
	return finding{Path: path, Kind: "metadata", Message: message}
}

func nonEmptyStringList(value any) bool {
	switch typed := value.(type) {
	case []string:
		if len(typed) == 0 {
			return false
		}
		for _, item := range typed {
			if strings.TrimSpace(item) == "" {
				return false
			}
		}
		return true
	case []any:
		if len(typed) == 0 {
			return false
		}
		for _, item := range typed {
			text, ok := item.(string)
			if !ok || strings.TrimSpace(text) == "" {
				return false
			}
		}
		return true
	default:
		return false
	}
}

func unresolvedShortcodeFindings(doc contentDoc, glossarySlugs, profileSlugs map[string]struct{}) []finding {
	var findings []finding
	for _, match := range glossaryShortcodeExpr.FindAllStringSubmatch(doc.Body, -1) {
		slug := normalizeShortcodeSlug(firstNonEmpty(match[1], match[2]))
		if _, ok := glossarySlugs[slug]; !ok {
			findings = append(findings, finding{
				Path:    doc.Path,
				Kind:    "glossary-ref",
				Message: fmt.Sprintf("unresolved glossary reference %q", slug),
			})
		}
	}
	for _, match := range profileShortcodeExpr.FindAllStringSubmatch(doc.Body, -1) {
		slug := normalizeShortcodeSlug(firstNonEmpty(match[1], match[2]))
		if _, ok := profileSlugs[slug]; !ok {
			findings = append(findings, finding{
				Path:    doc.Path,
				Kind:    "profile-ref",
				Message: fmt.Sprintf("unresolved profile reference %q", slug),
			})
		}
	}
	return findings
}

func normalizeShortcodeSlug(slug string) string {
	slug = strings.Trim(slug, `"'`)
	return strings.TrimPrefix(strings.TrimSpace(slug), "/")
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

func missingAltTextFindings(doc contentDoc) []finding {
	var findings []finding
	for _, match := range markdownImagePattern.FindAllStringSubmatch(doc.Body, -1) {
		if strings.TrimSpace(match[1]) != "" {
			continue
		}
		findings = append(findings, finding{
			Path:    doc.Path,
			Kind:    "image-alt",
			Message: fmt.Sprintf("image %q is missing alt text", strings.TrimSpace(match[2])),
		})
	}
	for _, match := range htmlImagePattern.FindAllStringSubmatch(doc.Body, -1) {
		attrText := match[1]
		attrMap := parseHTMLAttrs(attrText)
		alt, ok := attrMap["alt"]
		if !ok || strings.TrimSpace(alt) == "" {
			src := strings.TrimSpace(attrMap["src"])
			findings = append(findings, finding{
				Path:    doc.Path,
				Kind:    "image-alt",
				Message: fmt.Sprintf("HTML image %q is missing alt text", src),
			})
		}
	}
	for _, match := range lessonImagePattern.FindAllStringSubmatch(doc.Body, -1) {
		attrs := parseQuotedAttrs(match[1])
		src := strings.TrimSpace(attrs["src"])
		if src == "" {
			findings = append(findings, finding{
				Path:    doc.Path,
				Kind:    "image-src",
				Message: "lesson/image is missing src",
			})
		}
		if strings.TrimSpace(attrs["alt"]) == "" {
			findings = append(findings, finding{
				Path:    doc.Path,
				Kind:    "image-alt",
				Message: fmt.Sprintf("lesson/image %q is missing alt text", src),
			})
		}
	}
	return findings
}

func parseQuotedAttrs(attrText string) map[string]string {
	attrs := map[string]string{}
	for _, match := range attributePattern.FindAllStringSubmatch(attrText, -1) {
		attrs[strings.ToLower(match[1])] = firstNonEmpty(match[2], match[3], match[4])
	}
	return attrs
}

func parseHTMLAttrs(attrText string) map[string]string {
	attrs := map[string]string{}
	for _, field := range strings.Fields(attrText) {
		parts := strings.SplitN(field, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.ToLower(strings.TrimSpace(parts[0]))
		value := strings.Trim(parts[1], `"'`)
		attrs[key] = value
	}
	return attrs
}

func strictIntValue(value any) (int, bool) {
	switch typed := value.(type) {
	case int:
		return typed, true
	case int8:
		return int(typed), true
	case int16:
		return int(typed), true
	case int64:
		return int(typed), true
	case int32:
		return int(typed), true
	case uint:
		return int(typed), true
	case uint8:
		return int(typed), true
	case uint16:
		return int(typed), true
	case uint32:
		return int(typed), true
	default:
		return 0, false
	}
}
