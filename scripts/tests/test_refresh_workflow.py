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


if __name__ == "__main__":
    unittest.main()
