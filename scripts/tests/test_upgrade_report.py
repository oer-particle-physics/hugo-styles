from __future__ import annotations

import importlib.util
import json
import re
import shutil
import subprocess
import sys
import tempfile
import unittest
from pathlib import Path


REPO_ROOT = Path(__file__).parents[2]
SCRIPT = REPO_ROOT / "scripts/upgrade-report.py"
SPEC = importlib.util.spec_from_file_location("upgrade_report", SCRIPT)
report = importlib.util.module_from_spec(SPEC)
sys.modules[SPEC.name] = report
SPEC.loader.exec_module(report)


class UpgradeFixture(unittest.TestCase):
    def setUp(self) -> None:
        self.temp = tempfile.TemporaryDirectory()
        self.addCleanup(self.temp.cleanup)
        self.root = Path(self.temp.name)
        self.site = self.root / "lesson"
        self.site.mkdir()
        self.module = self.root / "module"
        self.notes = self.module / "content/docs/upgrades"
        self.notes.mkdir(parents=True)
        self.authors_note = self.notes / "authors-citation.md"
        self.authors_note.write_text(
            (REPO_ROOT / "content/docs/upgrades/authors-citation.md").read_text(),
            encoding="utf-8",
        )
        self.git("init", "-q")
        (self.site / "go.mod").write_text("module example.org/lesson\n", encoding="utf-8")
        self.git("add", ".")
        self.git("-c", "user.name=Test", "-c", "user.email=test@example.org",
                 "-c", "commit.gpgsign=false", "commit", "-qm", "fixture")

    def git(self, *args: str) -> str:
        return subprocess.run(
            ["git", *args], cwd=self.site, check=True, capture_output=True, text=True
        ).stdout

    def add_note(self, name: str, version: str) -> Path:
        path = self.notes / f"{name}.md"
        path.write_text(
            f'+++\ntitle = "{name}"\n[params.upgrade]\nrelease = "{version}"\n'
            'severity = "action"\n+++\n\n'
            '## Who is affected\n\nAll lessons.\n\n'
            '## What changes\n\nA setting changes.\n\n'
            '## What to do\n\nUpdate the setting.\n\n'
            '## How to verify\n\nInspect the rendered lesson.\n',
            encoding="utf-8",
        )
        return path

    def authors(self, citation: str | None = None, *, shortcode: bool = True) -> None:
        (self.site / "AUTHORS").write_text("Example Author <example>\n", encoding="utf-8")
        (self.site / "content").mkdir(exist_ok=True)
        (self.site / "content/_index.md").write_text(
            '{{< lesson/authors >}}' if shortcode else 'A lesson without an authors table.',
            encoding="utf-8",
        )
        if citation is not None:
            (self.site / "CITATION.cff").write_text(citation, encoding="utf-8")

    def build(self, old: str = "v0.5.7", new: str = "v0.6.0", **kwargs) -> tuple[str, str]:
        return report.build_report(self.site, self.module, old, new, **kwargs)


class UpgradeReportTests(UpgradeFixture):
    def test_title_links_and_changelog_cover_actual_versions(self) -> None:
        (self.module / "CHANGELOG.md").write_text(
            '# Changelog\n\n## [0.6.0](https://example.org/6)\n\nNew authors.\n'
            '\n## [0.5.7](https://example.org/5)\n\nAlready installed.\n',
            encoding="utf-8",
        )
        title, body = self.build()
        self.assertEqual(title, "chore: update hugo-styles from v0.5.7 to v0.6.0")
        self.assertIn("/compare/v0.5.7...v0.6.0", body)
        self.assertIn("/releases/tag/v0.6.0", body)
        self.assertIn("New authors.", body)
        self.assertNotIn("Already installed.", body)
        self.assertIn("Build verification has not completed.", body)

    def test_skipped_releases_are_included_and_future_notes_excluded(self) -> None:
        self.add_note("Earlier migration", "v0.5.0")
        self.add_note("Current migration", "v0.7.0")
        self.add_note("Future migration", "v0.8.0")
        _, body = self.build("v0.4.0", "v0.7.0")
        self.assertIn("Earlier migration", body)
        self.assertIn("Move author information", body)
        self.assertIn("Current migration", body)
        self.assertNotIn("Future migration", body)

    def test_resolved_older_notes_are_not_repeated(self) -> None:
        _, body = self.build("v0.6.0", "v0.7.0")
        self.assertNotIn("Move author information", body)
        self.assertIn("No applicable upgrade notes", body)

    def test_unresolved_migration_is_retained_after_its_release(self) -> None:
        self.authors()
        for old, new in (("v0.6.0", "v0.7.0"), ("v0.6.0", "v0.6.0")):
            with self.subTest(old=old, new=new):
                _, body = self.build(old, new)
                self.assertIn("Action required: migrate author information", body)
                self.assertIn("authors table will disappear", body)

    def test_migration_is_not_applied_before_its_release(self) -> None:
        self.authors()
        _, body = self.build("v0.4.0", "v0.5.7")
        self.assertNotIn("Move author information", body)
        self.assertNotIn("Action required: migrate", body)

    def test_missing_empty_and_comment_only_citations_need_migration(self) -> None:
        for citation in (None, "", "\n# A comment\n"):
            with self.subTest(citation=citation):
                self.authors(citation)
                self.assertEqual(report.legacy_authors(self.site).title,
                                 "Action required: migrate author information")

    def test_both_files_require_review_without_claiming_data_is_missing(self) -> None:
        self.authors("authors:\n  - name: Existing author\n")
        _, body = self.build()
        self.assertIn("Review required: reconcile contributor information", body)
        self.assertNotIn("The authors table will disappear", body)
        self.assertIn("does not compare names", body)

    def test_intentionally_absent_table_and_migrated_lesson_have_no_finding(self) -> None:
        self.assertIsNone(report.legacy_authors(self.site))
        self.authors(shortcode=False)
        self.assertIsNone(report.legacy_authors(self.site))
        self.authors("authors:\n  - name: Existing author\n")
        (self.site / "AUTHORS").unlink()
        self.assertIsNone(report.legacy_authors(self.site))

    def test_commented_shortcode_is_ignored_and_multilingual_homepage_checked(self) -> None:
        self.authors(shortcode=False)
        (self.site / "content/_index.md").write_text('<!-- {{< lesson/authors >}} -->')
        self.assertIsNone(report.legacy_authors(self.site))
        (self.site / "content/de").mkdir()
        (self.site / "content/de/_index.md").write_text('{{< lesson/authors title="Autoren" >}}')
        self.assertIsNotNone(report.legacy_authors(self.site))

    def test_same_version_refresh_and_noop_file_list(self) -> None:
        title, body = self.build("v0.6.0", "v0.6.0")
        self.assertEqual(title, "chore: refresh hugo-styles v0.6.0")
        self.assertNotIn("/compare/", body)
        self.assertIn("No file changes detected.", body)

    def test_unknown_pseudo_and_prerelease_versions_are_explicit(self) -> None:
        self.add_note("Future migration", "v0.8.0")
        for old in ("unknown", "v0.5.8-0.20260820120000-123456abcdef", "v0.6.0-rc.1"):
            with self.subTest(old=old):
                _, body = self.build(old, "v0.7.0")
                self.assertIn("Version range could not be determined", body)
                self.assertIn("Move author information", body)
                self.assertNotIn("Future migration", body)
                self.assertNotIn("/compare/", body)
        _, body = self.build("v0.6.0", "v0.7.0-rc.1")
        self.assertIn("Version range could not be determined", body)
        self.assertNotIn("No applicable", body)

    def test_missing_notes_and_changelog_do_not_claim_no_action_needed(self) -> None:
        self.authors_note.unlink()
        self.notes.rmdir()
        _, body = self.build()
        self.assertIn("does not provide structured upgrade notes", body)
        self.assertIn("Release highlights could not be selected", body)
        self.assertNotIn("No applicable", body)

    def test_pending_notes_must_be_assigned_before_release(self) -> None:
        self.add_note("Pending migration", "unreleased")
        self.assertEqual(len(report.read_notes(self.module)), 2)
        with self.assertRaisesRegex(ValueError, "assign the release version"):
            report.read_notes(self.module, require_released=True)
        _, body = self.build()
        self.assertIn("unassigned upgrade notes", body)
        self.assertNotIn("### Pending migration", body)

    def test_bad_metadata_or_missing_guidance_is_rejected(self) -> None:
        original = self.authors_note.read_text()
        cases = (
            ('severity = "breaking"', 'severity = "informational"'),
            ('release = "v0.6.0"', 'release = "tomorrow"'),
            ('check = "legacy-authors"', 'check = "typo"'),
            ('## What to do', '## Missing instructions'),
            ('[params.upgrade]', 'draft = true\n[params.upgrade]'),
            ('[params.upgrade]', 'params = "bad"\n[other]'),
        )
        for before, after in cases:
            with self.subTest(after=after):
                self.authors_note.write_text(original.replace(before, after))
                with self.assertRaises(ValueError):
                    report.read_notes(self.module)

    def test_changed_files_include_staged_deleted_and_new_files_with_vendor_summary(self) -> None:
        (self.site / "go.mod").unlink()
        (self.site / "CITATION.cff").write_text("title: Lesson\n")
        self.git("add", "CITATION.cff")
        (self.site / "new file.md").write_text("New content")
        vendor = self.site / "_vendor/module"
        vendor.mkdir(parents=True)
        (vendor / "one.html").write_text("one")
        (vendor / "two.html").write_text("two")
        _, body = self.build()
        self.assertIn("- `go.mod`", body)
        self.assertIn("- `CITATION.cff`", body)
        self.assertIn("- `new file.md`", body)
        self.assertIn("- `_vendor/` (2 changed files)", body)
        self.assertNotIn("_vendor/module/one.html", body)

    def test_custom_context_and_title_preserve_generated_guidance(self) -> None:
        self.authors()
        title, body = self.build(title="chore: custom update", extra_body="Additional lesson checks.")
        self.assertEqual(title, "chore: custom update")
        self.assertIn("Additional lesson checks.", body)
        self.assertIn("Action required", body)
        with self.assertRaises(ValueError):
            self.build(title="title\noutput=injection")
        with self.assertRaises(ValueError):
            self.build(old="v0.5.7\noutput=injection")

    def test_downgrades_are_rejected(self) -> None:
        with self.assertRaisesRegex(ValueError, "older than"):
            self.build("v0.7.0", "v0.6.0")

    def test_cli_and_validation_summary_keep_report_outside_lesson(self) -> None:
        self.authors()
        output = self.root / "report.md"
        github_output = self.root / "github-output"
        completed = subprocess.run(
            [sys.executable, str(SCRIPT), "report", "--site-root", str(self.site),
             "--module-root", str(self.module), "--old-version", "v0.5.7",
             "--new-version", "v0.6.0", "--output", str(output),
             "--github-output", str(github_output)],
            capture_output=True, text=True, check=True,
        )
        self.assertIn("chore: update hugo-styles", completed.stdout)
        self.assertEqual(github_output.read_text(),
                         "title=chore: update hugo-styles from v0.5.7 to v0.6.0\n")
        before = self.git("status", "--porcelain")
        for outcome, expected in (
            ("failure", "Build verification failed"),
            ("success", "configured build verification passed"),
            ("skipped", "Build verification was not run"),
            ("cancelled", "Build verification was cancelled"),
        ):
            with self.subTest(outcome=outcome):
                summary = self.root / f"summary-{outcome}.md"
                subprocess.run(
                    [sys.executable, str(SCRIPT), "finalize", "--report", str(output),
                     "--validation", outcome, "--summary", str(summary)], check=True,
                )
                self.assertEqual(output.read_text(), summary.read_text().rstrip() + "\n")
                self.assertIn(expected, summary.read_text())
                self.assertIn("Action required", summary.read_text())
                self.assertEqual(output.read_text().count("## Validation"), 1)
        self.assertEqual(before, self.git("status", "--porcelain"))


class UpgradeAuthoringTests(UpgradeFixture):
    def setUp(self) -> None:
        super().setUp()
        self.manifest = self.module / report.MANIFEST
        self.manifest.write_text(json.dumps({".": "0.6.0"}) + "\n")
        self.module_git("init", "-q")
        self.commit("chore: published fixture")
        self.module_git("-c", "tag.gpgsign=false", "tag", "v0.6.0")
        self.base = self.module_git("rev-parse", "HEAD").strip()

    def module_git(self, *args: str) -> str:
        return report.git(self.module, *args)

    def commit(self, message: str) -> None:
        self.module_git("add", ".")
        self.module_git("-c", "user.name=Test", "-c", "user.email=test@example.org",
                        "-c", "commit.gpgsign=false", "commit", "--allow-empty", "-qm", message)

    def pending(self) -> Path:
        path = self.add_note("future-setting", "unreleased")
        path.write_text(path.read_text().replace('severity = "action"', 'severity = "breaking"'))
        return path

    def cli(self, *args: str, check: bool = True) -> subprocess.CompletedProcess:
        return subprocess.run(
            [sys.executable, str(SCRIPT), *args, "--module-root", str(self.module)],
            capture_output=True, text=True, check=check,
        )

    def test_code_examples_survive_report_and_preview_formatting(self) -> None:
        path = self.pending()
        example = '\n```markdown\n## Required lesson heading\n```\n\n~~~python\n# Keep this comment\n~~~\n'
        path.write_text(path.read_text() + example)
        preview = report.build_preview(self.module)
        self.assertIn(example.strip(), preview)
        self.assertIn("#### What to do", preview)
        self.manifest.write_text('{".": "0.7.0"}\n')
        report.assign_release(self.module)
        _, body = self.build("v0.6.0", "v0.7.0")
        self.assertIn(example.strip(), body)

    def test_fences_inside_longer_fences_and_indented_code_are_preserved(self) -> None:
        source = '## Prose\n\n````markdown\n```markdown\n## Code heading\n```\n````\n\n    ## Indented code\n'
        expected = source.replace('## Prose', '#### Prose', 1)
        self.assertEqual(report.render_note_body(source), expected)

    def test_literal_shortcode_examples_are_unescaped_but_not_executed(self) -> None:
        path = self.pending()
        path.write_text(path.read_text() + '\n```markdown\n{{</* lesson/authors */>}}\n'
                        '{{%/* challenge */%}}Text{{%/* /challenge */%}}\n```\n')
        report.read_notes(self.module)
        preview = report.build_preview(self.module)
        self.assertIn('{{< lesson/authors >}}', preview)
        self.assertIn('{{% challenge %}}Text{{% /challenge %}}', preview)
        self.assertNotIn('/*', preview)
        self.assertIn('/*', path.read_text())
        path.write_text(path.read_text() + '\n{{< lesson/authors >}}\n')
        with self.assertRaisesRegex(ValueError, "escape literal shortcode"):
            report.read_notes(self.module)

    def test_empty_comment_only_and_placeholder_sections_fail(self) -> None:
        path = self.pending()
        original = path.read_text()
        for content in ('', '<!-- Explain later. -->', 'TODO: Add migration steps.', 'TBD', '```\n```'):
            with self.subTest(content=content):
                path.write_text(original.replace('Update the setting.', content))
                with self.assertRaises(ValueError):
                    report.read_notes(self.module)

    def test_example_headings_do_not_satisfy_required_sections(self) -> None:
        path = self.pending()
        path.write_text(path.read_text().replace('## How to verify', '## Extra') +
                        '\n```markdown\n## How to verify\nA heading in an example.\n```\n')
        with self.assertRaisesRegex(ValueError, "missing '## How to verify'"):
            report.read_notes(self.module)

    def test_preview_needs_no_git_checkout_or_versions_and_does_not_claim_validation(self) -> None:
        self.pending()
        # The preview also works with an unpacked module, without a Git history.
        shutil.rmtree(self.module / '.git')
        output = self.root / 'preview.md'
        summary = self.root / 'summary.md'
        self.cli('preview', '--output', str(output), '--summary', str(summary))
        body = output.read_text()
        self.assertIn('future-setting (unreleased, breaking)', body)
        self.assertIn('Applicability has not been checked', body)
        self.assertNotIn('Move author information', body)
        self.assertNotIn('No applicable upgrade notes', body)
        self.assertNotIn('verification passed', body)
        self.assertEqual(body.strip(), summary.read_text().strip())
        self.assertIn(body, self.cli('preview').stdout)

    def test_preview_can_include_assigned_notes_and_optional_lesson_checks(self) -> None:
        self.pending()
        self.authors()
        body = self.cli('preview', '--release', 'v0.6.0', '--site-root', str(self.site)).stdout
        self.assertIn('future-setting', body)
        self.assertIn('Action required: migrate author information', body)
        self.assertNotIn('Applicability has not been checked', body)

    def test_assignment_uses_manifest_and_preserves_historical_notes_and_body(self) -> None:
        path = self.pending()
        text = path.read_text().replace('release = "unreleased"', "release = 'unreleased' # planned")
        text = text.replace('title = "future-setting"', 'title = "Support +++ in examples"')
        text += '\n```toml\nrelease = "unreleased"\n```\n'
        path.write_text(text)
        historical = self.authors_note.read_bytes()
        self.manifest.write_text('{".": "0.7.0"}\n')
        self.cli('assign-release')
        self.assertEqual(self.authors_note.read_bytes(), historical)
        self.assertEqual(path.read_text(), text.replace("release = 'unreleased' # planned", 'release = "v0.7.0" # planned', 1))
        self.assertIn('No pending notes', self.cli('assign-release').stdout)

    def test_assignment_refuses_already_tagged_version(self) -> None:
        path = self.pending()
        before = path.read_bytes()
        result = self.cli('assign-release', check=False)
        self.assertNotEqual(result.returncode, 0)
        self.assertIn('already tagged', result.stderr)
        self.assertEqual(before, path.read_bytes())

    def test_assignment_validates_all_notes_before_writing(self) -> None:
        path = self.pending()
        before = path.read_bytes()
        bad = self.add_note('invalid-note', 'unreleased')
        bad.write_text(bad.read_text().replace('Update the setting.', 'TODO: Fix me.'))
        self.manifest.write_text('{".": "0.7.0"}\n')
        with self.assertRaises(ValueError):
            report.assign_release(self.module)
        self.assertEqual(before, path.read_bytes())

    def test_publication_allows_development_but_blocks_incomplete_release(self) -> None:
        self.pending()
        output = self.root / 'outputs'
        self.cli('check-publication', '--github-output', str(output))
        self.assertIn('publication-ready=false', output.read_text())
        self.manifest.write_text('{".": "0.7.0"}\n')
        output.unlink()
        result = self.cli('check-publication', '--github-output', str(output), check=False)
        self.assertNotEqual(result.returncode, 0)
        self.assertIn('assign the release version', result.stderr)
        self.assertFalse(output.exists())
        report.assign_release(self.module)
        self.cli('check-publication', '--github-output', str(output))
        self.assertIn('publication-ready=true', output.read_text())

    def test_release_branch_refresh_requires_assignment_again(self) -> None:
        path = self.pending()
        self.commit('feat!: change config')
        self.base = self.module_git('rev-parse', 'HEAD').strip()
        pending = path.read_text()
        self.manifest.write_text('{".": "0.7.0"}\n')
        self.commit('chore: release 0.7.0')
        self.cli('assign-release')
        self.commit('chore: assign upgrade notes')
        self.assertTrue(report.check_publication(self.module))
        self.assertEqual(report.check_pr(self.module, self.base), 'v0.7.0')
        # Release Please may rebuild its branch from main, restoring pending notes.
        path.write_text(pending)
        self.commit('chore: refresh release branch')
        with self.assertRaisesRegex(ValueError, 'assign the release version'):
            report.check_publication(self.module)
        output = self.root / 'outputs'
        result = self.cli('check-pr', '--base-ref', self.base,
                          '--github-output', str(output), check=False)
        self.assertNotEqual(result.returncode, 0)
        self.assertIn('release-version=v0.7.0', output.read_text())
        self.cli('assign-release')
        self.assertTrue(report.check_publication(self.module))
        self.assertEqual(report.check_pr(self.module, self.base), 'v0.7.0')

    def test_invalid_manifest_fails_with_a_useful_error(self) -> None:
        for value in ('[]', 'null', '{".": "next"}'):
            with self.subTest(value=value):
                self.manifest.write_text(value)
                result = self.cli('check-publication', check=False)
                self.assertNotEqual(result.returncode, 0)
                self.assertIn('must contain a stable root version', result.stderr)
                self.assertNotIn('Traceback', result.stderr)

    def test_release_pr_requires_assignment_and_a_coherent_manifest(self) -> None:
        path = self.pending()
        self.manifest.write_text('{".": "0.7.0"}\n')
        self.commit('chore: release 0.7.0')
        with self.assertRaisesRegex(ValueError, 'assign the release version'):
            report.check_pr(self.module, self.base)
        path.write_text(path.read_text().replace('"unreleased"', '"v0.8.0"'))
        with self.assertRaisesRegex(ValueError, 'newer than the release manifest'):
            report.check_pr(self.module, self.base)
        path.write_text(path.read_text().replace('"v0.8.0"', '"v0.7.0"'))
        self.assertEqual(report.check_pr(self.module, self.base), 'v0.7.0')
        self.manifest.write_text('{".": "0.5.0"}\n')
        with self.assertRaisesRegex(ValueError, 'must increase'):
            report.check_pr(self.module, self.base)

    def test_declared_breaking_pr_requires_a_new_breaking_note(self) -> None:
        for title, body in (('feat!: change config', ''), ('feat: config', 'BREAKING CHANGE: Old config removed.')):
            with self.subTest(title=title, body=body):
                with self.assertRaisesRegex(ValueError, 'add a new breaking upgrade note'):
                    report.check_pr(self.module, self.base, title, body)
        path = self.pending()
        self.commit('feat: change config')
        self.assertEqual(report.check_pr(self.module, self.base, 'feat!: change config'), '')
        path.write_text(path.read_text().replace('severity = "breaking"', 'severity = "action"'))
        with self.assertRaises(ValueError):
            report.check_pr(self.module, self.base, 'feat!: change config')

    def test_commit_breaking_marker_is_checked_even_without_pr_title_marker(self) -> None:
        self.commit('feat: change config\n\nBREAKING CHANGE: A setting was removed.')
        with self.assertRaisesRegex(ValueError, 'declares a breaking change'):
            report.check_pr(self.module, self.base, 'feat: change config')

    def test_documented_markers_in_pr_body_are_not_breaking_declarations(self) -> None:
        self.assertEqual(report.check_pr(self.module, self.base, 'docs: explain commits',
                         '```text\nBREAKING CHANGE: Example\n```\n'), '')

    def test_archetype_creates_pending_note_and_requires_placeholders_to_be_filled(self) -> None:
        (self.module / 'hugo.toml').write_text('title = "Note authoring"\n')
        (self.module / 'archetypes').mkdir()
        shutil.copyfile(REPO_ROOT / 'archetypes/upgrade-note.md', self.module / 'archetypes/upgrade-note.md')
        subprocess.run(['hugo', 'new', 'content', '--kind', 'upgrade-note',
                        'docs/upgrades/new-setting.md'], cwd=self.module,
                       check=True, capture_output=True, text=True)
        path = self.notes / 'new-setting.md'
        self.assertIn('release = "unreleased"', path.read_text())
        with self.assertRaisesRegex(ValueError, 'replace scaffold placeholders'):
            report.read_notes(self.module)
        path.write_text(path.read_text().replace('TODO:', 'Guidance:'))
        self.assertIn('New Setting', self.cli('preview').stdout)

    def test_hugo_and_report_show_the_same_literal_shortcode_example(self) -> None:
        source = '{{</* lesson/authors */>}}\n{{%/* hint */%}}Help{{%/* /hint */%}}'
        # The maintainer documentation must itself show a copyable escaped call.
        instructions = (REPO_ROOT / 'content/docs/hugo-styles-maintenance.md').read_text()
        recipe = re.search(r'```markdown\n(.*?)\n```', instructions, flags=re.S)[1]
        source += '\n' + recipe
        (self.site / 'hugo.toml').write_text(
            'baseURL = "https://example.org/"\ndisableKinds = ["taxonomy", "term", "RSS", "sitemap"]\n')
        (self.site / 'content').mkdir()
        (self.site / 'layouts').mkdir()
        (self.site / 'layouts/home.html').write_text('{{ .Content }}')
        (self.site / 'content/_index.md').write_text('+++\ntitle = "Example"\n+++\n\n```markdown\n' + source + '\n```\n')
        subprocess.run(['hugo'], cwd=self.site, check=True, capture_output=True, text=True)
        from html.parser import HTMLParser

        class Text(HTMLParser):
            value = ''

            def handle_data(self, data):
                self.value += data

        html = Text()
        html.feed((self.site / 'public/index.html').read_text())
        literal = '{{< lesson/authors >}}\n{{% hint %}}Help{{% /hint %}}'
        self.assertIn(literal, html.value)
        self.assertIn('{{</* lesson/authors */>}}', html.value)
        self.assertIn(literal, report.render_note_body('```markdown\n' + source + '\n```\n'))


if __name__ == "__main__":
    unittest.main()
