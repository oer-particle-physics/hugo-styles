package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	toml "github.com/pelletier/go-toml/v2"
	"gopkg.in/yaml.v3"
)

type finding struct {
	Path    string
	Kind    string
	Message string
}

type legacyPattern struct {
	Kind    string
	Regex   *regexp.Regexp
	Message string
}

var (
	includePattern         = regexp.MustCompile(`\{%\s*include\s+[^%]+%\}`)
	rawOpenPattern         = regexp.MustCompile(`\{%\s*raw\s*%\}`)
	rawClosePattern        = regexp.MustCompile(`\{%\s*endraw\s*%\}`)
	commentBlockPattern    = regexp.MustCompile(`(?s)\{%\s*comment\s*%\}.*?\{%\s*endcomment\s*%\}`)
	baseURLPattern         = regexp.MustCompile(`\{\{\s*site\.baseurl\s*\}\}/?`)
	pageRootLinkPattern    = regexp.MustCompile(`\{\{\s*page\.root\s*\}\}\{% link _episodes/([^.]+)\.md %\}`)
	relativeFigPattern     = regexp.MustCompile(`(?:\.\./)+fig/`)
	relativeFigHTMLPattern = regexp.MustCompile(`(<img\b[^>]*\bsrc=")(?:\.\./)+fig/`)
	relativeFigMDPattern   = regexp.MustCompile(`(!\[[^\]]*\]\()((?:\.\./)+)fig/`)
	imageTagPattern        = regexp.MustCompile(`(?is)<img\b[^>]*>`)
	attributePattern       = regexp.MustCompile(`(?i)\b([a-z0-9:-]+)\s*=\s*(?:"([^"]*)"|'([^']*)'|([^\s"'>]+))`)
	headingPattern         = regexp.MustCompile(`^(\s{0,3})(#{1,5})(\s+.*)$`)
	attrPattern            = regexp.MustCompile(`^\{:\s*\.?([a-zA-Z0-9_-]+)\s*\}$`)
	titleHeadingPattern    = regexp.MustCompile(`^##\s+(.+)$`)
	htmlCommentPattern     = regexp.MustCompile(`(?s)<!--.*?-->`)
	htmlAttrCommentPattern = regexp.MustCompile(`(?m)^<!--\s*\{:\s*\.[a-zA-Z0-9_-]+\s*\}\s*-->\s*$`)
	youtubeIframePattern   = regexp.MustCompile(`(?is)<iframe[^>]+src="https?://(?:www\.)?youtube\.com/embed/([^"?/&]+)[^"]*"[^>]*></iframe>`)
	vimeoIframePattern     = regexp.MustCompile(`(?is)<iframe[^>]+src="https?://player\.vimeo\.com/video/([^"?/&]+)[^"]*"[^>]*></iframe>`)
	fencedCodeBlockPattern = regexp.MustCompile("(?s)```.*?```|~~~.*?~~~")
	inlineCodeSpanPattern  = regexp.MustCompile("`[^`\n]+`")
	legacyTitlePattern     = regexp.MustCompile(`(?m)^title:\s*(.+?)\s*$`)
)

var legacyPatterns = []legacyPattern{
	{"liquid-include", includePattern, "remove or replace Liquid includes"},
	{"liquid-raw", rawOpenPattern, "strip Liquid raw tags after converting embedded code"},
	{"legacy-baseurl", baseURLPattern, "rewrite legacy baseurl asset links"},
	{"legacy-link", regexp.MustCompile(`\{% link [^%]+%\}`), "replace Jekyll link tags with relref or normal links"},
	{"legacy-attr-block", regexp.MustCompile(`\{:\s*\.[a-zA-Z0-9_-]+\s*\}`), "convert Carpentries fenced-attribute blocks"},
	{"legacy-gh-variables", regexp.MustCompile(`gh_variables\.html`), "remove Jekyll GitHub variable includes"},
	{"liquid-variable", regexp.MustCompile(`\{\{\s*site\.[^}]+\}\}`), "replace site-scoped Liquid variables that are not handled automatically"},
	{"legacy-auto-ids", regexp.MustCompile(`\{:\s*auto_ids\s*\}`), "rewrite glossary auto_ids syntax"},
	{"iframe-embed", regexp.MustCompile(`(?i)<iframe\b`), "convert supported iframe embeds to Hugo shortcodes"},
	{"legacy-fig-path", relativeFigPattern, "rewrite lesson-local figure paths to site-root /fig/"},
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	switch os.Args[1] {
	case "check":
		if err := runCheck(os.Args[2:]); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case "migrate":
		if err := runMigrate(os.Args[2:]); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	default:
		usage()
		os.Exit(2)
	}
}

func usage() {
	fmt.Println("usage:")
	fmt.Println("  hugo-styles-migrate check [path]")
	fmt.Println("  hugo-styles-migrate migrate --source <legacy repo> --dest <output dir>")
}

func runCheck(args []string) error {
	fs := flag.NewFlagSet("check", flag.ContinueOnError)
	if err := fs.Parse(args); err != nil {
		return err
	}
	root := "."
	if fs.NArg() > 0 {
		root = fs.Arg(0)
	}

	findings, err := collectFindings(root)
	if err != nil {
		return err
	}
	if len(findings) == 0 {
		fmt.Println("no migration blockers found")
		return nil
	}

	for _, f := range findings {
		fmt.Printf("%s [%s] %s\n", f.Path, f.Kind, f.Message)
	}
	return fmt.Errorf("found %d migration issues", len(findings))
}

func runMigrate(args []string) error {
	fs := flag.NewFlagSet("migrate", flag.ContinueOnError)
	source := fs.String("source", "", "legacy lesson repository")
	dest := fs.String("dest", "", "output directory")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *source == "" || *dest == "" {
		return errors.New("both --source and --dest are required")
	}

	if err := os.MkdirAll(filepath.Join(*dest, "content", "episodes"), 0o755); err != nil {
		return err
	}
	for _, dir := range []string{"learners", "instructors", "glossary", "profiles"} {
		if err := os.MkdirAll(filepath.Join(*dest, "content", dir), 0o755); err != nil {
			return err
		}
	}

	for _, dir := range []string{"fig", "files", "data", "code"} {
		srcDir := filepath.Join(*source, dir)
		if stat, err := os.Stat(srcDir); err == nil && stat.IsDir() {
			if err := copyTree(srcDir, filepath.Join(*dest, "static", dir)); err != nil {
				return err
			}
		}
	}

	if err := migrateLessonHomePage(*source, *dest); err != nil {
		return err
	}
	if err := migrateRootPage(*source, *dest, "reference.md", filepath.Join("content", "reference.md"), ""); err != nil {
		return err
	}
	if err := migrateRootPage(*source, *dest, "setup.md", filepath.Join("content", "learners", "setup.md"), ""); err != nil {
		return err
	}

	if err := migrateExtras(*source, *dest); err != nil {
		return err
	}
	if err := migrateEpisodes(*source, *dest); err != nil {
		return err
	}

	fmt.Printf("migrated %s -> %s\n", *source, *dest)
	return nil
}

func collectLegacyFindings(root string) ([]finding, error) {
	var findings []finding
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		name := d.Name()
		if d.IsDir() {
			switch name {
			case ".git", "node_modules", "public", "resources", "testdata":
				return filepath.SkipDir
			}
			return nil
		}
		if filepath.Ext(path) != ".md" {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		text := string(data)
		scanText := stripFencedCodeBlocks(text)
		for _, pattern := range legacyPatterns {
			if pattern.Regex.MatchString(scanText) {
				findings = append(findings, finding{
					Path:    path,
					Kind:    pattern.Kind,
					Message: pattern.Message,
				})
			}
		}
		if strings.Contains(path, string(filepath.Separator)+"_episodes"+string(filepath.Separator)) {
			for _, key := range []string{"questions:", "objectives:", "keypoints:"} {
				if !strings.Contains(text, "\n"+key) && !strings.HasPrefix(text, key) {
					findings = append(findings, finding{
						Path:    path,
						Kind:    "metadata",
						Message: "episode is missing " + strings.TrimSuffix(key, ":"),
					})
				}
			}
		}
		return nil
	})
	return findings, err
}

func migrateLessonHomePage(source, dest string) error {
	srcPath := filepath.Join(source, "index.md")
	if _, err := os.Stat(srcPath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}

	data, err := os.ReadFile(srcPath)
	if err != nil {
		return err
	}
	meta, body, err := parseFrontMatter(string(data))
	if err != nil {
		return err
	}

	text := transformMarkdown(body)
	delete(meta, "layout")
	delete(meta, "root")
	delete(meta, "permalink")
	meta["layout"] = "hextra-home"

	if stringValue(meta["title"]) == "" {
		title, err := readLegacyLessonTitle(source)
		if err != nil {
			return err
		}
		if title != "" {
			meta["title"] = title
		}
	}

	text = buildLessonHomePage(stringValue(meta["title"]), firstLessonLink(source), text)
	text, err = renderFrontMatter(meta, text)
	if err != nil {
		return err
	}

	destPath := filepath.Join(dest, "content", "_index.md")
	if err := os.MkdirAll(filepath.Dir(destPath), 0o755); err != nil {
		return err
	}
	return os.WriteFile(destPath, []byte(text), 0o644)
}

func stripFencedCodeBlocks(text string) string {
	return inlineCodeSpanPattern.ReplaceAllString(fencedCodeBlockPattern.ReplaceAllString(text, ""), "")
}

func migrateRootPage(source, dest, srcRel, destRel, layoutOverride string) error {
	srcPath := filepath.Join(source, srcRel)
	if _, err := os.Stat(srcPath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	return transformAndWrite(srcPath, filepath.Join(dest, destRel), 0, false, true, false, layoutOverride)
}

func migrateExtras(source, dest string) error {
	extrasDir := filepath.Join(source, "_extras")
	entries, err := os.ReadDir(extrasDir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".md" {
			continue
		}
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		srcPath := filepath.Join(extrasDir, entry.Name())
		base := strings.TrimSuffix(entry.Name(), ".md")
		switch base {
		case "guide":
			if err := transformAndWrite(srcPath, filepath.Join(dest, "content", "instructors", "instructor-notes.md"), 0, false, false, false, ""); err != nil {
				return err
			}
		case "figures":
			continue
		default:
			if err := transformAndWrite(srcPath, filepath.Join(dest, "content", "learners", base+".md"), 0, false, false, false, ""); err != nil {
				return err
			}
		}
	}
	return nil
}

func migrateEpisodes(source, dest string) error {
	episodesDir := filepath.Join(source, "_episodes")
	entries, err := os.ReadDir(episodesDir)
	if err != nil {
		return err
	}
	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".md" || entry.Name() == ".gitkeep" {
			continue
		}
		names = append(names, entry.Name())
	}
	sort.Strings(names)
	for index, name := range names {
		srcPath := filepath.Join(episodesDir, name)
		draft := strings.HasPrefix(name, ".")
		slug := strings.TrimPrefix(strings.TrimSuffix(name, ".md"), ".")
		destPath := filepath.Join(dest, "content", "episodes", slug, "index.md")
		if err := transformAndWrite(srcPath, destPath, (index+1)*10, draft, false, true, ""); err != nil {
			return err
		}
	}
	return nil
}

func transformAndWrite(sourcePath, destPath string, weight int, draft bool, stripJekyllKeys bool, promoteHeadings bool, layoutOverride string) error {
	data, err := os.ReadFile(sourcePath)
	if err != nil {
		return err
	}
	meta, body, err := parseFrontMatter(string(data))
	if err != nil {
		return err
	}
	text := transformMarkdown(body)
	if promoteHeadings {
		text = promoteHeadingLevels(text)
	}
	if stripJekyllKeys {
		delete(meta, "layout")
		delete(meta, "root")
		delete(meta, "permalink")
	}
	if layoutOverride != "" {
		meta["layout"] = layoutOverride
	}
	if weight > 0 {
		meta["weight"] = weight
	}
	if draft {
		meta["draft"] = true
	}
	text, err = renderFrontMatter(meta, text)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(destPath), 0o755); err != nil {
		return err
	}
	return os.WriteFile(destPath, []byte(text), 0o644)
}

func transformMarkdown(text string) string {
	text = commentBlockPattern.ReplaceAllString(text, "")
	text = includePattern.ReplaceAllString(text, "")
	text = rawOpenPattern.ReplaceAllString(text, "")
	text = rawClosePattern.ReplaceAllString(text, "")
	text = htmlAttrCommentPattern.ReplaceAllString(text, "")
	text = baseURLPattern.ReplaceAllString(text, "")
	text = pageRootLinkPattern.ReplaceAllString(text, `{{< relref "/episodes/$1" >}}`)
	text = relativeFigPattern.ReplaceAllString(text, "/fig/")
	text = relativeFigHTMLPattern.ReplaceAllString(text, `${1}/fig/`)
	text = relativeFigMDPattern.ReplaceAllString(text, `${1}/fig/`)
	text = convertImageTags(text)
	text = youtubeIframePattern.ReplaceAllString(text, "{{< youtube $1 >}}")
	text = vimeoIframePattern.ReplaceAllString(text, "{{< vimeo $1 >}}")
	text = strings.ReplaceAll(text, "{:auto_ids}", "")
	text = convertAttrBlocks(text)
	text = convertFenceAttributes(text)
	text = cleanupSpacing(text)
	return text
}

func promoteHeadingLevels(text string) string {
	var out []string
	scanner := bufio.NewScanner(strings.NewReader(text))
	inFence := false
	var fenceMarker string
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimLeft(line, " ")
		switch {
		case !inFence && (strings.HasPrefix(trimmed, "```") || strings.HasPrefix(trimmed, "~~~")):
			fenceMarker = strings.Repeat(string(trimmed[0]), len(trimmed)-len(strings.TrimLeft(trimmed, string(trimmed[0]))))
			inFence = true
			out = append(out, line)
			continue
		case inFence && strings.HasPrefix(trimmed, fenceMarker):
			inFence = false
			fenceMarker = ""
			out = append(out, line)
			continue
		}

		if !inFence {
			if match := headingPattern.FindStringSubmatch(line); len(match) == 4 {
				prefix := match[1]
				hashes := match[2]
				if len(hashes) < 6 {
					line = prefix + strings.Repeat("#", len(hashes)+1) + match[3]
				}
			}
		}
		out = append(out, line)
	}
	return strings.Join(out, "\n")
}

func renderFrontMatter(meta map[string]any, body string) (string, error) {
	if len(meta) == 0 {
		return body, nil
	}
	data, err := toml.Marshal(meta)
	if err != nil {
		return "", err
	}
	return "+++\n" + strings.TrimSpace(string(data)) + "\n+++\n" + body, nil
}

func convertAttrBlocks(text string) string {
	lines := strings.Split(text, "\n")
	for i := 0; i < len(lines); i++ {
		class := parseAttrClass(lines[i])
		if class == "" {
			continue
		}
		if isFenceAttrClass(class) {
			if !followsFence(lines, i) {
				lines = spliceLines(lines, i, i+1, nil)
				i--
			}
			continue
		}
		attrDepth := quoteDepth(lines[i])
		targetDepth := attrDepth + 1
		start := findAttrBlockStart(lines, i, targetDepth)
		if start >= i {
			continue
		}
		block := strings.Join(unquoteLevels(lines[start:i], targetDepth), "\n")
		block = convertAttrBlocks(block)
		replacement := wrapAttrBlock(class, block)
		replacement = applyQuoteDepth(replacement, attrDepth)
		replLines := strings.Split(replacement, "\n")
		lines = spliceLines(lines, start, i+1, replLines)
		i = start + len(replLines) - 1
	}
	return strings.Join(lines, "\n")
}

func findAttrBlockStart(lines []string, index, targetDepth int) int {
	for i := index - 1; i >= 0; i-- {
		if quoteDepth(lines[i]) == targetDepth {
			trimmed := strings.TrimSpace(unquoteLevels([]string{lines[i]}, targetDepth)[0])
			if titleHeadingPattern.MatchString(trimmed) {
				return i
			}
		}
	}
	start := index - 1
	for start >= 0 && quoteDepth(lines[start]) >= targetDepth {
		start--
	}
	return start + 1
}

func followsFence(lines []string, index int) bool {
	for i := index - 1; i >= 0; i-- {
		trimmed := strings.TrimSpace(lines[i])
		if trimmed == "" {
			continue
		}
		return trimmed == "~~~"
	}
	return false
}

func wrapAttrBlock(class, block string) string {
	title, body := extractHeading(block)
	switch class {
	case "challenge":
		return wrapShortcode("challenge", title, body)
	case "solution":
		return wrapShortcode("solution", title, body)
	case "hint":
		return wrapShortcode("hint", title, body)
	case "callout":
		return wrapCallout("note", title, body)
	case "prereq", "checklist", "testimonial", "caution", "warning", "discussion":
		return wrapCallout(class, title, body)
	default:
		return block
	}
}

func wrapShortcode(name, title, body string) string {
	if title != "" {
		return fmt.Sprintf("{{< %s title=%q >}}\n%s\n{{< /%s >}}", name, sanitizeShortcode(title), strings.TrimSpace(body), name)
	}
	return fmt.Sprintf("{{< %s >}}\n%s\n{{< /%s >}}", name, strings.TrimSpace(body), name)
}

func wrapCallout(kind, title, body string) string {
	if title != "" {
		return fmt.Sprintf("{{< callout type=%q title=%q >}}\n%s\n{{< /callout >}}", kind, sanitizeShortcode(title), strings.TrimSpace(body))
	}
	return fmt.Sprintf("{{< callout type=%q >}}\n%s\n{{< /callout >}}", kind, strings.TrimSpace(body))
}

func convertFenceAttributes(text string) string {
	scanner := bufio.NewScanner(strings.NewReader(text))
	var out []string
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, "~~~") {
			out = append(out, line)
			continue
		}
		fenceLang := strings.TrimSpace(strings.TrimPrefix(trimmed, "~~~"))
		var block []string
		for scanner.Scan() {
			next := scanner.Text()
			if strings.TrimSpace(next) == "~~~" {
				break
			}
			block = append(block, next)
		}
		attrLine := ""
		if scanner.Scan() {
			attrLine = scanner.Text()
		}
		class := parseAttrClass(attrLine)
		lang := fenceLanguage(fenceLang, class)
		out = append(out, "```"+lang)
		out = append(out, block...)
		out = append(out, "```")
		if attrLine != "" && class == "" {
			out = append(out, attrLine)
		}
	}
	return strings.Join(out, "\n")
}

func cleanupSpacing(text string) string {
	text = strings.ReplaceAll(text, "\n\n\n\n", "\n\n\n")
	text = strings.ReplaceAll(text, "\n\n\n\n", "\n\n\n")
	return strings.TrimSpace(text) + "\n"
}

func buildLessonHomePage(title, startLink, body string) string {
	body = strings.TrimSpace(htmlCommentPattern.ReplaceAllString(body, ""))
	lead, remainder := splitLessonHomeIntro(body)
	var out []string

	if title != "" {
		out = append(out, strings.Join([]string{
			`<div class="hx:mt-6 hx:mb-6">`,
			`{{< hextra/hero-headline >}}`,
			title,
			`{{< /hextra/hero-headline >}}`,
			`</div>`,
		}, "\n"))
	}

	if lead != "" {
		out = append(out, strings.Join([]string{
			`<div class="hx:mb-12">`,
			`{{< hextra/hero-subtitle >}}`,
			lead,
			`{{< /hextra/hero-subtitle >}}`,
			`</div>`,
		}, "\n"))
	}

	if startLink != "" {
		out = append(out, strings.Join([]string{
			`<div class="hx:mb-6">`,
			fmt.Sprintf(`{{< hextra/hero-button text=%q link=%q >}}`, "Start Lesson", startLink),
			`</div>`,
		}, "\n"))
	}

	out = append(out, strings.Join([]string{
		`<div class="hx:mt-6"></div>`,
		`{{< lesson/overview >}}`,
		`<div class="hx:mt-6"></div>`,
		`{{< lesson/schedule title="Schedule" >}}`,
		`<div class="hx:mt-6"></div>`,
		`{{< lesson/authors title="Authors and Contributors" >}}`,
	}, "\n"))

	if remainder != "" {
		out = append(out, "## About This Lesson\n\n"+remainder)
	}

	return cleanupSpacing(strings.Join(out, "\n\n"))
}

func splitLessonHomeIntro(text string) (string, string) {
	blocks := splitMarkdownBlocks(text)
	if len(blocks) == 0 {
		return "", ""
	}

	leadEnd := 0
	for leadEnd < len(blocks) && isHeroLeadBlock(blocks[leadEnd]) {
		leadEnd++
	}

	if leadEnd == 0 {
		return "", strings.TrimSpace(text)
	}

	lead := strings.Join(blocks[:leadEnd], "\n\n")
	remainder := strings.Join(blocks[leadEnd:], "\n\n")
	return strings.TrimSpace(lead), strings.TrimSpace(remainder)
}

func splitMarkdownBlocks(text string) []string {
	text = strings.TrimSpace(text)
	if text == "" {
		return nil
	}

	var blocks []string
	var current []string
	for _, line := range strings.Split(text, "\n") {
		if strings.TrimSpace(line) == "" {
			if len(current) > 0 {
				blocks = append(blocks, strings.TrimSpace(strings.Join(current, "\n")))
				current = nil
			}
			continue
		}
		current = append(current, line)
	}
	if len(current) > 0 {
		blocks = append(blocks, strings.TrimSpace(strings.Join(current, "\n")))
	}
	return blocks
}

func isHeroLeadBlock(block string) bool {
	block = strings.TrimSpace(block)
	if block == "" {
		return false
	}
	switch {
	case strings.HasPrefix(block, "#"):
		return false
	case strings.HasPrefix(block, "{{<"):
		return false
	case strings.HasPrefix(block, ">"):
		return false
	case strings.HasPrefix(block, "```"), strings.HasPrefix(block, "~~~"):
		return false
	case strings.HasPrefix(block, "- "), strings.HasPrefix(block, "* "):
		return false
	case orderedListItemPattern.MatchString(block):
		return false
	case strings.HasPrefix(block, "<") && !strings.HasPrefix(block, "<!--"):
		return false
	default:
		return true
	}
}

func readLegacyLessonTitle(source string) (string, error) {
	for _, name := range []string{"_config.yml", "_config.yaml"} {
		path := filepath.Join(source, name)
		data, err := os.ReadFile(path)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return "", err
		}

		var meta map[string]any
		if err := yaml.Unmarshal(data, &meta); err != nil {
			if title := extractLegacyLessonTitle(data); title != "" {
				return title, nil
			}
			return "", fmt.Errorf("parse %s: %w", path, err)
		}
		if title := stringValue(meta["title"]); title != "" {
			return title, nil
		}
	}
	return "", nil
}

func extractLegacyLessonTitle(data []byte) string {
	match := legacyTitlePattern.FindSubmatch(data)
	if len(match) != 2 {
		return ""
	}
	title := strings.TrimSpace(string(match[1]))
	title = strings.Trim(title, `"'`)
	return strings.TrimSpace(title)
}

func firstLessonLink(source string) string {
	episodesDir := filepath.Join(source, "_episodes")
	entries, err := os.ReadDir(episodesDir)
	if err != nil {
		return "episodes/"
	}

	var names []string
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || filepath.Ext(name) != ".md" || name == ".gitkeep" || strings.HasPrefix(name, ".") {
			continue
		}
		names = append(names, name)
	}
	if len(names) == 0 {
		return "episodes/"
	}

	sort.Strings(names)
	slug := strings.TrimSuffix(names[0], ".md")
	return "episodes/" + slug + "/"
}

func stringValue(value any) string {
	if text, ok := value.(string); ok {
		return strings.TrimSpace(text)
	}
	return ""
}

func convertImageTags(text string) string {
	return imageTagPattern.ReplaceAllStringFunc(text, func(tag string) string {
		attrs := parseHTMLAttributes(tag)
		src, ok := attrs["src"]
		if !ok {
			return tag
		}
		if !(strings.Contains(src, "fig/") || strings.HasPrefix(src, "/fig/") || strings.HasPrefix(src, "../fig/")) {
			return tag
		}

		src = normalizeFigureSrc(src)
		parts := []string{fmt.Sprintf(`src=%q`, src)}
		for _, key := range []string{"alt", "width", "style", "class"} {
			if value, ok := attrs[key]; ok && value != "" {
				parts = append(parts, fmt.Sprintf(`%s=%q`, key, value))
			}
		}
		return fmt.Sprintf("{{< lesson/image %s >}}", strings.Join(parts, " "))
	})
}

func parseHTMLAttributes(tag string) map[string]string {
	attrs := make(map[string]string)
	for _, m := range attributePattern.FindAllStringSubmatch(tag, -1) {
		value := m[2]
		if value == "" {
			value = m[3]
		}
		if value == "" {
			value = m[4]
		}
		attrs[strings.ToLower(m[1])] = value
	}
	return attrs
}

func normalizeFigureSrc(src string) string {
	src = strings.TrimSpace(src)
	for strings.HasPrefix(src, "../") {
		src = strings.TrimPrefix(src, "../")
	}
	src = strings.TrimPrefix(src, "./")
	src = strings.TrimPrefix(src, "/")
	if strings.HasPrefix(src, "fig/") {
		src = "/" + src
	}
	return src
}

func parseAttrClass(line string) string {
	matches := attrPattern.FindStringSubmatch(strings.TrimSpace(line))
	if len(matches) == 2 {
		return matches[1]
	}
	return ""
}

var orderedListItemPattern = regexp.MustCompile(`^\d+\.\s`)

func isFenceAttrClass(class string) bool {
	if class == "source" || class == "output" || class == "bash" {
		return true
	}
	return strings.HasPrefix(class, "language-")
}

func isQuotedLine(line string) bool {
	trimmed := strings.TrimLeft(line, " ")
	return strings.HasPrefix(trimmed, ">")
}

func unquoteLevels(lines []string, levels int) []string {
	out := make([]string, 0, len(lines))
	for _, line := range lines {
		current := line
		for level := 0; level < levels; level++ {
			trimmed := strings.TrimLeft(current, " ")
			prefixSpaces := len(current) - len(trimmed)
			if strings.HasPrefix(trimmed, "> ") {
				current = current[:prefixSpaces] + trimmed[2:]
				continue
			}
			if strings.HasPrefix(trimmed, ">") {
				current = current[:prefixSpaces] + trimmed[1:]
				continue
			}
			break
		}
		out = append(out, current)
	}
	return out
}

func applyQuoteDepth(text string, depth int) string {
	if depth == 0 {
		return text
	}
	lines := strings.Split(text, "\n")
	prefix := strings.Repeat("> ", depth)
	blankPrefix := strings.TrimRight(prefix, " ")
	for i, line := range lines {
		if line == "" {
			lines[i] = blankPrefix
			continue
		}
		lines[i] = prefix + line
	}
	return strings.Join(lines, "\n")
}

func extractHeading(block string) (string, string) {
	lines := strings.Split(strings.TrimSpace(block), "\n")
	if len(lines) == 0 {
		return "", ""
	}
	matches := titleHeadingPattern.FindStringSubmatch(strings.TrimSpace(lines[0]))
	if len(matches) != 2 {
		return "", strings.TrimSpace(block)
	}
	return matches[1], strings.TrimSpace(strings.Join(lines[1:], "\n"))
}

func sanitizeShortcode(text string) string {
	return strings.ReplaceAll(text, `"`, `'`)
}

func quoteDepth(line string) int {
	trimmed := strings.TrimLeft(line, " ")
	depth := 0
	for strings.HasPrefix(trimmed, ">") {
		depth++
		trimmed = strings.TrimLeft(strings.TrimPrefix(trimmed, ">"), " ")
	}
	return depth
}

func fenceLanguage(existing, class string) string {
	switch {
	case existing != "":
		return existing
	case class == "source" || class == "bash":
		return "bash"
	case class == "output":
		return "text"
	case strings.HasPrefix(class, "language-"):
		return strings.TrimPrefix(class, "language-")
	default:
		return "text"
	}
}

func spliceLines(lines []string, start, end int, replacement []string) []string {
	out := make([]string, 0, len(lines)-end+start+len(replacement))
	out = append(out, lines[:start]...)
	out = append(out, replacement...)
	out = append(out, lines[end:]...)
	return out
}

func splitYAMLFrontMatter(text string) (string, string, bool) {
	if !strings.HasPrefix(text, "---\n") {
		return "", "", false
	}
	closing := strings.Index(text[4:], "\n---\n")
	if closing == -1 {
		return "", "", false
	}
	closing += 4
	front := text[:closing+5]
	body := text[closing+5:]
	return strings.TrimSuffix(front, "\n---\n"), body, true
}

func copyTree(source, dest string) error {
	return filepath.WalkDir(source, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dest, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		return copyFile(path, target)
	})
}

func copyFile(source, dest string) error {
	in, err := os.Open(source)
	if err != nil {
		return err
	}
	defer in.Close()

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, in); err != nil {
		return err
	}
	return os.WriteFile(dest, buf.Bytes(), 0o644)
}
