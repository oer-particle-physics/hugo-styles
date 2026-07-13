package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	toml "github.com/pelletier/go-toml/v2"
)

const hugoStylesModule = "github.com/oer-particle-physics/hugo-styles"

var migratedSections = []string{"episodes", "learners", "instructors", "glossary", "profiles"}
var migratedStaticDirs = []string{"fig", "files", "data", "code"}

type migrationChange struct {
	RelPath string
	Source  string
	Target  string
}

func migrateIntoTemplate(source, dest string, dryRun bool) error {
	sourceRoot, destRoot, err := validateMigrationRoots(source, dest)
	if err != nil {
		return err
	}
	if err := validateTemplateDestination(destRoot); err != nil {
		return err
	}

	staging, err := os.MkdirTemp(filepath.Dir(destRoot), ".hugo-styles-migration-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(staging)

	if err := generateMigration(sourceRoot, staging); err != nil {
		return err
	}
	changes, err := planMigrationChanges(staging, destRoot)
	if err != nil {
		return err
	}
	preserved, err := preservedMigrationFiles(destRoot, changes)
	if err != nil {
		return err
	}
	printMigrationPlan(changes, preserved, dryRun)
	if dryRun {
		return nil
	}
	if err := applyMigrationChanges(changes, staging, destRoot); err != nil {
		return err
	}

	fmt.Printf("migration completed: %s -> %s\n", sourceRoot, destRoot)
	fmt.Println("next steps:")
	fmt.Println("  1. Update lesson metadata and repository URLs in hugo.toml.")
	fmt.Println("  2. Run hugo-styles-migrate check on the template repository.")
	fmt.Println("  3. Build with Hugo and review custom legacy content manually.")
	return nil
}

func validateMigrationRoots(source, dest string) (string, string, error) {
	sourceRoot, err := filepath.Abs(source)
	if err != nil {
		return "", "", err
	}
	destRoot, err := filepath.Abs(dest)
	if err != nil {
		return "", "", err
	}
	sourceRoot, err = filepath.EvalSymlinks(sourceRoot)
	if err != nil {
		return "", "", fmt.Errorf("resolve source: %w", err)
	}
	destRoot, err = filepath.EvalSymlinks(destRoot)
	if err != nil {
		return "", "", fmt.Errorf("resolve destination: %w", err)
	}

	for label, path := range map[string]string{"source": sourceRoot, "destination": destRoot} {
		stat, statErr := os.Stat(path)
		if statErr != nil {
			return "", "", fmt.Errorf("%s is not accessible: %w", label, statErr)
		}
		if !stat.IsDir() {
			return "", "", fmt.Errorf("%s is not a directory: %s", label, path)
		}
	}
	if pathsOverlap(sourceRoot, destRoot) {
		return "", "", errors.New("source and destination directories must not overlap")
	}
	if stat, statErr := os.Stat(filepath.Join(sourceRoot, "_episodes")); statErr != nil || !stat.IsDir() {
		return "", "", fmt.Errorf("source does not contain a legacy _episodes directory: %s", sourceRoot)
	}
	return sourceRoot, destRoot, nil
}

func pathsOverlap(first, second string) bool {
	return pathWithin(first, second) || pathWithin(second, first)
}

func pathWithin(path, parent string) bool {
	rel, err := filepath.Rel(parent, path)
	if err != nil {
		return false
	}
	return rel == "." || (rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)))
}

func validateTemplateDestination(dest string) error {
	hugoConfig, err := os.ReadFile(filepath.Join(dest, "hugo.toml"))
	if err != nil {
		return fmt.Errorf("destination must contain hugo.toml: %w", err)
	}
	goMod, err := os.ReadFile(filepath.Join(dest, "go.mod"))
	if err != nil {
		return fmt.Errorf("destination must contain go.mod: %w", err)
	}
	var config struct {
		Module struct {
			Imports []struct {
				Path string `toml:"path"`
			} `toml:"imports"`
		} `toml:"module"`
	}
	if err := toml.Unmarshal(hugoConfig, &config); err != nil {
		return fmt.Errorf("destination hugo.toml is invalid: %w", err)
	}
	hasImport := false
	for _, item := range config.Module.Imports {
		if item.Path == hugoStylesModule {
			hasImport = true
			break
		}
	}
	if !hasImport || !goModRequires(string(goMod), hugoStylesModule) {
		return fmt.Errorf("destination is not recognized as a hugo-styles template: %s", dest)
	}

	top, err := gitCommand(dest, "rev-parse", "--show-toplevel")
	if err != nil {
		return errors.New("destination must be the root of a Git worktree")
	}
	topRoot, err := filepath.EvalSymlinks(strings.TrimSpace(top))
	if err != nil || topRoot != dest {
		return errors.New("destination must be the root of its Git worktree")
	}
	status, err := gitCommand(dest, "status", "--porcelain", "--untracked-files=all")
	if err != nil {
		return fmt.Errorf("inspect destination worktree: %w", err)
	}
	if strings.TrimSpace(status) != "" {
		return errors.New("destination Git worktree must be clean; commit or stash changes first")
	}
	return nil
}

func goModRequires(contents, modulePath string) bool {
	inRequireBlock := false
	for _, rawLine := range strings.Split(contents, "\n") {
		line := strings.TrimSpace(strings.SplitN(rawLine, "//", 2)[0])
		if line == "" {
			continue
		}
		if inRequireBlock {
			if line == ")" {
				inRequireBlock = false
				continue
			}
			fields := strings.Fields(line)
			if len(fields) > 0 && fields[0] == modulePath {
				return true
			}
			continue
		}
		if line == "require (" {
			inRequireBlock = true
			continue
		}
		if strings.HasPrefix(line, "require ") {
			fields := strings.Fields(strings.TrimPrefix(line, "require "))
			if len(fields) > 0 && fields[0] == modulePath {
				return true
			}
		}
	}
	return false
}

func gitCommand(root string, args ...string) (string, error) {
	commandArgs := append([]string{"-C", root}, args...)
	output, err := exec.Command("git", commandArgs...).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git %s: %s", strings.Join(args, " "), strings.TrimSpace(string(output)))
	}
	return string(output), nil
}

func generateMigration(source, staging string) error {
	for _, section := range migratedSections {
		if err := os.MkdirAll(filepath.Join(staging, "content", section), 0o755); err != nil {
			return err
		}
	}
	for _, dir := range migratedStaticDirs {
		sourceDir := filepath.Join(source, dir)
		if stat, err := os.Stat(sourceDir); err == nil && stat.IsDir() {
			if err := copyTree(sourceDir, filepath.Join(staging, "static", dir)); err != nil {
				return err
			}
		}
	}
	if authors := filepath.Join(source, "AUTHORS"); fileExists(authors) {
		if err := copyFile(authors, filepath.Join(staging, "AUTHORS")); err != nil {
			return err
		}
	}
	if err := migrateLessonHomePage(source, staging); err != nil {
		return err
	}
	if err := migrateRootPage(source, staging, "reference.md", filepath.Join("content", "reference.md"), ""); err != nil {
		return err
	}
	if err := migrateRootPage(source, staging, "setup.md", filepath.Join("content", "learners", "setup.md"), ""); err != nil {
		return err
	}
	if err := migrateExtras(source, staging); err != nil {
		return err
	}
	return migrateEpisodes(source, staging)
}

func fileExists(path string) bool {
	stat, err := os.Stat(path)
	return err == nil && !stat.IsDir()
}

func planMigrationChanges(staging, dest string) ([]migrationChange, error) {
	changes := map[string]migrationChange{}
	add := func(rel string) {
		rel = filepath.Clean(rel)
		changes[rel] = migrationChange{
			RelPath: rel,
			Source:  filepath.Join(staging, rel),
			Target:  filepath.Join(dest, rel),
		}
	}

	for _, rel := range []string{filepath.Join("content", "_index.md"), filepath.Join("content", "reference.md"), "AUTHORS"} {
		if fileExists(filepath.Join(staging, rel)) {
			add(rel)
		}
	}
	for _, section := range migratedSections {
		relDir := filepath.Join("content", section)
		for _, root := range []string{dest, staging} {
			entries, err := os.ReadDir(filepath.Join(root, relDir))
			if err != nil {
				if errors.Is(err, os.ErrNotExist) {
					continue
				}
				return nil, err
			}
			for _, entry := range entries {
				if entry.Name() != "_index.md" {
					add(filepath.Join(relDir, entry.Name()))
				}
			}
		}
	}
	for _, dir := range migratedStaticDirs {
		rel := filepath.Join("static", dir)
		if _, err := os.Stat(filepath.Join(dest, rel)); err == nil {
			add(rel)
		}
		if _, err := os.Stat(filepath.Join(staging, rel)); err == nil {
			add(rel)
		}
	}

	result := make([]migrationChange, 0, len(changes))
	for _, change := range changes {
		result = append(result, change)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].RelPath < result[j].RelPath })
	return result, nil
}

func preservedMigrationFiles(dest string, changes []migrationChange) ([]string, error) {
	var preserved []string
	err := filepath.WalkDir(dest, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			if entry.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}
		for _, change := range changes {
			if pathWithin(path, change.Target) {
				return nil
			}
		}
		preserved = append(preserved, path)
		return nil
	})
	sort.Strings(preserved)
	return preserved, err
}

func printMigrationPlan(changes []migrationChange, preserved []string, dryRun bool) {
	label := "migration plan"
	if dryRun {
		label = "migration dry run"
	}
	fmt.Printf("%s (%d managed paths):\n", label, len(changes))
	for _, change := range changes {
		targetExists := pathExists(change.Target)
		sourceExists := pathExists(change.Source)
		action := "replace"
		switch {
		case sourceExists && !targetExists:
			action = "add"
		case targetExists && !sourceExists:
			action = "remove"
		}
		path := filepath.ToSlash(change.RelPath)
		if dryRun {
			path = change.Target
		}
		fmt.Printf("  %-7s %s\n", action, path)
	}
	if dryRun {
		fmt.Printf("preserve (%d files):\n", len(preserved))
		for _, path := range preserved {
			fmt.Printf("  preserve %s\n", path)
		}
	}
}

func pathExists(path string) bool {
	_, err := os.Lstat(path)
	return err == nil
}

func applyMigrationChanges(changes []migrationChange, staging, dest string) (returnErr error) {
	backup, err := os.MkdirTemp(filepath.Dir(dest), ".hugo-styles-migration-backup-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(backup)

	backedUp := make([]migrationChange, 0, len(changes))
	installed := make([]migrationChange, 0, len(changes))
	defer func() {
		if returnErr == nil {
			return
		}
		var rollbackErrors []error
		for index := len(installed) - 1; index >= 0; index-- {
			if err := os.RemoveAll(installed[index].Target); err != nil {
				rollbackErrors = append(rollbackErrors, fmt.Errorf("remove installed %s: %w", installed[index].RelPath, err))
			}
		}
		for index := len(backedUp) - 1; index >= 0; index-- {
			change := backedUp[index]
			backupPath := filepath.Join(backup, change.RelPath)
			if err := os.MkdirAll(filepath.Dir(change.Target), 0o755); err != nil {
				rollbackErrors = append(rollbackErrors, fmt.Errorf("prepare restore for %s: %w", change.RelPath, err))
				continue
			}
			if err := os.Rename(backupPath, change.Target); err != nil {
				rollbackErrors = append(rollbackErrors, fmt.Errorf("restore %s: %w", change.RelPath, err))
			}
		}
		if len(rollbackErrors) > 0 {
			returnErr = fmt.Errorf("%w; rollback failed: %v", returnErr, errors.Join(rollbackErrors...))
		}
	}()

	for _, change := range changes {
		if !pathExists(change.Target) {
			continue
		}
		backupPath := filepath.Join(backup, change.RelPath)
		if err := os.MkdirAll(filepath.Dir(backupPath), 0o755); err != nil {
			return err
		}
		if err := os.Rename(change.Target, backupPath); err != nil {
			return fmt.Errorf("back up %s: %w", change.RelPath, err)
		}
		backedUp = append(backedUp, change)
	}

	for _, change := range changes {
		if !pathExists(change.Source) {
			continue
		}
		if err := os.MkdirAll(filepath.Dir(change.Target), 0o755); err != nil {
			return err
		}
		if err := os.Rename(change.Source, change.Target); err != nil {
			return fmt.Errorf("install %s: %w", change.RelPath, err)
		}
		installed = append(installed, change)
	}
	return nil
}
