from __future__ import annotations

import subprocess
import tempfile
import unittest
from pathlib import Path


REPO_ROOT = Path(__file__).resolve().parents[2]


class FooterRenderingTests(unittest.TestCase):
    def render_footers(self, pages: dict[str, str]) -> dict[str, str]:
        with tempfile.TemporaryDirectory() as temp_name:
            site = Path(temp_name)
            (site / "go.mod").write_text(
                "module example.org/footer-test\n\n"
                "go 1.26\n\n"
                "require github.com/oer-particle-physics/hugo-styles v0.0.0\n\n"
                f"replace github.com/oer-particle-physics/hugo-styles => {REPO_ROOT}\n",
                encoding="utf-8",
            )
            (site / "hugo.toml").write_text(
                'baseURL = "https://example.org/lesson/"\n'
                'disableKinds = ["taxonomy", "term", "RSS", "sitemap", "404"]\n'
                # Leave Lastmod unset on date-only pages to exercise the footer's fallback.
                '\n[frontmatter]\nlastmod = ["lastmod"]\n'
                '\n[params]\ndateFormat = "2006-01-02"\n'
                '\n[module]\n'
                '  [[module.imports]]\n'
                '    path = "github.com/oer-particle-physics/hugo-styles"\n'
                '    [[module.imports.mounts]]\n'
                '      source = "layouts"\n'
                '      target = "layouts"\n',
                encoding="utf-8",
            )
            for layout in ("home.html", "episodes/page.html", "episodes/section.html"):
                path = site / "layouts" / layout
                path.parent.mkdir(parents=True, exist_ok=True)
                path.write_text('{{ partial "footer.html" . }}\n', encoding="utf-8")

            for name, dates in pages.items():
                path = site / "content" / f"{name}.md"
                path.parent.mkdir(parents=True, exist_ok=True)
                path.write_text(f'+++\ntitle = "Lesson page"\n{dates}+++\n', encoding="utf-8")

            build = subprocess.run(
                ["hugo", "--cacheDir", str(site / "cache")],
                cwd=site,
                capture_output=True,
                text=True,
                check=False,
            )
            self.assertEqual(build.returncode, 0, build.stdout + build.stderr)
            return {
                name: (site / "public" / name.removesuffix("_index") / "index.html").read_text(
                    encoding="utf-8"
                )
                for name in pages
            }

    def test_homepage_displays_newer_episode_date(self) -> None:
        footers = self.render_footers({
            "_index": "lastmod = 2025-01-01\n",
            "episodes/older": "date = 2025-02-01\n",
            "episodes/newer": "lastmod = 2025-03-01\n",
        })
        self.assertIn("<span>Updated: 2025-03-01</span>", footers["_index"])
        self.assertIn("<span>Updated: 2025-02-01</span>", footers["episodes/older"])
        self.assertIn("<span>Updated: 2025-03-01</span>", footers["episodes/newer"])

    def test_homepage_displays_its_own_date_when_newest(self) -> None:
        footers = self.render_footers({
            "_index": "lastmod = 2025-03-01\n",
            "episodes/older": "lastmod = 2025-01-01\n",
        })
        self.assertIn("<span>Updated: 2025-03-01</span>", footers["_index"])
        self.assertIn("<span>Updated: 2025-01-01</span>", footers["episodes/older"])

    def test_undated_site_omits_updated_label(self) -> None:
        footers = self.render_footers({"_index": "", "episodes/undated": ""})
        for name, footer in footers.items():
            with self.subTest(page=name):
                self.assertIn('class="lesson-footer ', footer)
                self.assertNotIn("Updated:", footer)
