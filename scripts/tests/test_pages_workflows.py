from __future__ import annotations

import unittest
from pathlib import Path


REPO_ROOT = Path(__file__).parents[2]
PAGES_WORKFLOWS = (
    REPO_ROOT / ".github" / "workflows" / "pages.yml",
    REPO_ROOT / "dist" / "template-files" / ".github" / "workflows" / "pages.yml",
)


class PagesWorkflowTests(unittest.TestCase):
    def test_pushes_are_gated_by_the_repository_default_branch(self) -> None:
        for workflow_path in PAGES_WORKFLOWS:
            with self.subTest(workflow=workflow_path.relative_to(REPO_ROOT)):
                workflow = workflow_path.read_text(encoding="utf-8")

                self.assertIn('branches: ["**"]', workflow)
                self.assertNotIn('branches: ["main"]', workflow)
                self.assertIn("github.event_name == 'pull_request' ||", workflow)
                self.assertIn(
                    "github.ref == format('refs/heads/{0}', "
                    "github.event.repository.default_branch)",
                    workflow,
                )


if __name__ == "__main__":
    unittest.main()
