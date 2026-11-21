from __future__ import annotations

import subprocess
import tempfile
import unittest
from pathlib import Path


REPO_ROOT = Path(__file__).resolve().parents[2]


class NavigationRenderingTests(unittest.TestCase):
    def test_episodes_menu_falls_back_when_configured_children_are_missing(self) -> None:
        with tempfile.TemporaryDirectory() as temp_name:
            temp = Path(temp_name)
            config = temp / "hugo.toml"
            config.write_text(
                (REPO_ROOT / "hugo.toml").read_text(encoding="utf-8").replace(
                    'name = "Examples"\n', 'name = "Examples"\n    identifier = "episodes"\n', 1
                )
                + '\n[[menus.main]]\nname = "Retired episode"\nparent = "episodes"\n'
                'pageRef = "/episodes/retired"\nweight = 1\n',
                encoding="utf-8",
            )
            output = temp / "public"
            build = subprocess.run(
                [
                    "hugo", "--config", f"{config},hugo-docs.toml", "--panicOnWarning",
                    "--baseURL", "https://example.org/lesson/versions/v1/", "--destination", str(output),
                ],
                cwd=REPO_ROOT,
                capture_output=True,
                text=True,
                check=False,
            )
            self.assertEqual(build.returncode, 0, build.stdout + build.stderr)
            html = (output / "reference" / "index.html").read_text(encoding="utf-8")

        self.assertTrue('data-episode-resume-placement="navbar"' in html, "Missing wide Continue slot")
        self.assertTrue('data-episode-resume-placement="dropdown"' in html, "Missing dropdown Continue slot")
        dropdown = html.index('data-episode-resume-placement="dropdown"')
        first_episode = 'href="/lesson/versions/v1/episodes/01-introduction/"'
        self.assertIn(first_episode, html[dropdown:html.index("</ul>", dropdown)])
        self.assertNotIn("Retired episode", html)
