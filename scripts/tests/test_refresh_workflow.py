from __future__ import annotations

import unittest
from pathlib import Path


REPO_ROOT = Path(__file__).parents[2]
REFRESH_WORKFLOW = (
    REPO_ROOT / ".github" / "workflows" / "reusable-refresh-vendored-modules.yml"
)


class RefreshWorkflowTests(unittest.TestCase):
    @classmethod
    def setUpClass(cls) -> None:
        cls.workflow = REFRESH_WORKFLOW.read_text(encoding="utf-8")

    def test_missing_token_check_runs_before_pull_request_creation(self) -> None:
        credential_check = self.workflow.index(
            "- name: Check workflow update credentials"
        )
        create_pull_request = self.workflow.index("- name: Create pull request")

        self.assertLess(credential_check, create_pull_request)
        self.assertIn(
            "WORKFLOW_SYNC_TOKEN: ${{ secrets.WORKFLOW_SYNC_TOKEN }}",
            self.workflow,
        )

    def test_missing_token_check_detects_workflow_changes(self) -> None:
        self.assertIn(
            'git -C "${GITHUB_WORKSPACE}" status --porcelain '
            "--untracked-files=all -- .github/workflows",
            self.workflow,
        )
        self.assertIn(
            "::error title=WORKFLOW_SYNC_TOKEN required::", self.workflow
        )
        self.assertIn(
            "::warning title=WORKFLOW_SYNC_TOKEN not configured::", self.workflow
        )

    def test_default_token_remains_the_fallback(self) -> None:
        self.assertIn(
            "token: ${{ secrets.WORKFLOW_SYNC_TOKEN || github.token }}",
            self.workflow,
        )

    def test_report_survives_build_failure_and_precedes_pr_creation(self) -> None:
        steps = [self.workflow.index(f"- name: {name}") for name in (
            "Record installed hugo-styles version",
            "Refresh module metadata and managed files",
            "Prepare upgrade report",
            "Verify refreshed site build",
            "Publish upgrade summary",
            "Create pull request",
        )]
        self.assertEqual(steps, sorted(steps))
        self.assertIn("always() && steps.report.outcome == 'success'", self.workflow)
        self.assertIn("VALIDATION_OUTCOME: ${{ steps.verify.outcome }}", self.workflow)
        self.assertIn("title: ${{ steps.report.outputs.title }}", self.workflow)
        self.assertIn("body-path: ${{ runner.temp }}/hugo-styles-upgrade.md", self.workflow)
        self.assertNotIn("continue-on-error", self.workflow)

    def test_managed_caller_uses_generated_title_and_keeps_existing_branch(self) -> None:
        caller = (REPO_ROOT / "dist/template-files/.github/workflows/refresh-vendored-modules.yml").read_text()
        self.assertIn("name: Update hugo-styles", caller)
        self.assertNotIn("pr-title:", caller)
        self.assertNotIn("pr-body:", caller)
        self.assertIn('default: "chore/refresh-vendored-hugo-modules"', self.workflow)

    def test_report_helper_is_distributed_and_checked_for_drift(self) -> None:
        manifest = (REPO_ROOT / "dist/template-files/manifest.tsv").read_text()
        self.assertIn("scripts/upgrade-report.py|scripts/upgrade-report.py|0755", manifest)
        pages = (REPO_ROOT / "dist/template-files/.github/workflows/pages.yml").read_text()
        self.assertIn("        scripts/upgrade-report.py\n", pages)

    def test_release_publication_is_gated_before_release_please(self) -> None:
        workflow = (REPO_ROOT / ".github/workflows/release-please.yml").read_text()
        self.assertLess(workflow.index("check-publication"),
                        workflow.index("uses: googleapis/release-please-action"))
        self.assertIn("fetch-depth: 0", workflow)
        self.assertIn("skip-github-release: ${{ steps.upgrade_policy.outputs.publication-ready != 'true' }}",
                      workflow)
        self.assertNotIn("continue-on-error", workflow)

    def test_pr_preview_is_available_even_when_release_policy_fails(self) -> None:
        workflow = (REPO_ROOT / ".github/workflows/upgrade-notes.yml").read_text()
        self.assertIn("types: [opened, synchronize, reopened, edited]", workflow)
        self.assertIn("--base-ref", workflow)
        self.assertIn("id: preview\n        if: ${{ !cancelled() }}", workflow)
        self.assertIn("steps.preview.outcome == 'success'", workflow)
        self.assertIn('args+=(--release "${RELEASE_VERSION}")', workflow)
        self.assertIn("actions/upload-artifact@", workflow)
        self.assertNotIn("continue-on-error", workflow)


if __name__ == "__main__":
    unittest.main()
