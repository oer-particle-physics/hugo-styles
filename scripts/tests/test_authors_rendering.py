from __future__ import annotations

import subprocess
import tempfile
import textwrap
import unittest
from dataclasses import dataclass, field
from html.parser import HTMLParser
from pathlib import Path


REPO_ROOT = Path(__file__).parents[2]

HEADER = "cff-version: 1.2.0\nmessage: Please cite this lesson.\ntitle: Example lesson\n"


@dataclass
class Cell:
    text: str = ""
    links: list[str] = field(default_factory=list)


class AuthorsTableParser(HTMLParser):
    def __init__(self) -> None:
        super().__init__(convert_charrefs=True)
        self.tables = 0
        self.headers: list[str] = []
        self.rows: list[list[Cell]] = []
        self._cell: Cell | None = None

    def handle_starttag(self, tag: str, attrs: list[tuple[str, str | None]]) -> None:
        if tag == "table":
            self.tables += 1
        elif tag == "tr" and self.tables:
            self._row: list[Cell] = []
        elif tag in {"th", "td"}:
            self._cell = Cell()
        elif tag == "a" and self._cell is not None:
            self._cell.links.append(dict(attrs).get("href") or "")

    def handle_endtag(self, tag: str) -> None:
        if tag in {"th", "td"} and self._cell is not None:
            self._cell.text = " ".join(self._cell.text.split())
            if tag == "th":
                self.headers.append(self._cell.text)
            else:
                self._row.append(self._cell)
            self._cell = None
        elif tag == "tr" and self.tables and self._row:
            self.rows.append(self._row)

    def handle_data(self, data: str) -> None:
        if self._cell is not None:
            self._cell.text += data


@dataclass
class Build:
    returncode: int
    output: str
    table: AuthorsTableParser | None

    @property
    def names(self) -> list[str]:
        assert self.table is not None
        return [row[0].text for row in self.table.rows]


class AuthorsRenderingTests(unittest.TestCase):
    def build(self, citation: str | None, *, header: bool = True) -> Build:
        """Build a site with this CITATION.cff; header adds the required keys."""
        if citation is not None:
            citation = textwrap.dedent(citation)
            if header:
                citation = HEADER + citation
        with tempfile.TemporaryDirectory() as temp_name:
            site = Path(temp_name)
            (site / "content").mkdir()
            (site / "layouts").mkdir()

            (site / "go.mod").write_text(
                "\n".join(
                    [
                        "module example.org/authors-test",
                        "",
                        "go 1.26",
                        "",
                        "require github.com/oer-particle-physics/hugo-styles v0.0.0",
                        "",
                        "replace github.com/oer-particle-physics/hugo-styles "
                        f"=> {REPO_ROOT}",
                        "",
                    ]
                ),
                encoding="utf-8",
            )
            (site / "hugo.toml").write_text(
                "\n".join(
                    [
                        'baseURL = "https://example.org/lesson/"',
                        'title = "Authors test"',
                        'disableKinds = ["taxonomy", "term", "RSS", "sitemap", "robotsTXT", "404"]',
                        "",
                        "[module]",
                        "  [[module.imports]]",
                        '    path = "github.com/oer-particle-physics/hugo-styles"',
                        "    [[module.imports.mounts]]",
                        '      source = "layouts"',
                        '      target = "layouts"',
                        "",
                    ]
                ),
                encoding="utf-8",
            )
            (site / "layouts" / "home.html").write_text(
                "<!doctype html><html><body>{{ .Content }}</body></html>\n",
                encoding="utf-8",
            )
            (site / "content" / "_index.md").write_text(
                '+++\ntitle = "Authors"\n+++\n\n{{< lesson/authors >}}\n',
                encoding="utf-8",
            )
            if citation is not None:
                (site / "CITATION.cff").write_text(citation, encoding="utf-8")

            destination = site / "public"
            completed = subprocess.run(
                [
                    "hugo",
                    "--cacheDir",
                    str(site / "cache"),
                    "--destination",
                    str(destination),
                ],
                cwd=site,
                check=False,
                capture_output=True,
                text=True,
            )
            table = None
            if completed.returncode == 0:
                table = AuthorsTableParser()
                table.feed((destination / "index.html").read_text(encoding="utf-8"))
            return Build(completed.returncode, completed.stdout + completed.stderr, table)

    def render(self, citation: str | None, *, header: bool = True) -> Build:
        build = self.build(citation, header=header)
        self.assertEqual(build.returncode, 0, build.output)
        self.assertNotIn("WARN", build.output)
        return build

    def assertBuildFails(self, citation: str, message: str, *, header: bool = True) -> None:
        build = self.build(citation, header=header)
        self.assertNotEqual(build.returncode, 0, build.output)
        self.assertIn(message, build.output)

    def assertNoTable(self, build: Build) -> None:
        assert build.table is not None
        self.assertEqual(build.table.tables, 0)

    def test_missing_file_renders_nothing(self) -> None:
        self.assertNoTable(self.render(None))

    def test_empty_file_renders_nothing(self) -> None:
        for citation in ("", "\n", "# only a comment\n"):
            with self.subTest(citation=citation):
                self.assertNoTable(self.render(citation, header=False))

    def test_malformed_file_fails_the_build(self) -> None:
        cases = {
            "unparseable": ("authors: [\n", "could not be parsed"),
            "top-level list": ("- given-names: Ada\n", "must be a mapping"),
            "top-level scalar": ("just a string\n", "must be a mapping"),
            "empty mapping": ("{}\n", "must have an authors list"),
            "empty list": ("[]\n", "must be a mapping"),
            "false": ("false\n", "must be a mapping"),
            "zero": ("0\n", "must be a mapping"),
            "empty string": ('""\n', "must be a mapping"),
        }
        for label, (citation, message) in cases.items():
            with self.subTest(label):
                self.assertBuildFails(citation, message, header=False)

    def test_bad_authors_list_fails_the_build(self) -> None:
        cases = {
            "no authors key": ("", "must have an authors list"),
            "null authors": ("authors:\n", "must have an authors list"),
            "authors is a string": ("authors: Ada Lovelace\n", "must have an authors list"),
            "authors is a mapping": ("authors:\n  given-names: Ada\n", "must have an authors list"),
            "empty authors": ("authors: []\n", "must not be empty"),
        }
        for label, (citation, message) in cases.items():
            with self.subTest(label):
                self.assertBuildFails(citation, message)

    def test_bad_author_entries_fail_the_build(self) -> None:
        cases = {
            "string entry": ("  - Ada Lovelace\n", "author 1 must be a mapping"),
            "list entry": ("  - [nested]\n", "author 1 must be a mapping"),
            "no name": ("  - orcid: https://orcid.org/0000-0002-1825-0097\n", "author 1 has no name"),
        }
        for label, (entry, message) in cases.items():
            with self.subTest(label):
                self.assertBuildFails("authors:\n" + entry, message)

    def test_full_person_name(self) -> None:
        build = self.render(
            """\
            authors:
              - given-names: Alexander
                name-particle: von
                family-names: Humboldt
                name-suffix: III
              - family-names: Curie
            """
        )
        self.assertEqual(build.names, ["Alexander von Humboldt III", "Curie"])

    def test_organisation_author(self) -> None:
        build = self.render(
            """\
            authors:
              - name: HEP Software Foundation
                alias: hsf
            """
        )
        assert build.table is not None
        self.assertEqual(build.names, ["HEP Software Foundation"])
        self.assertEqual(build.table.rows[0][1].links, ["https://github.com/hsf"])

    def test_alias_only_author(self) -> None:
        build = self.render(
            """\
            authors:
              - alias: octocat
              - alias: https://github.com/hubot/
              - alias: not a github handle
            """
        )
        assert build.table is not None
        self.assertEqual(build.names, ["octocat", "hubot", "not a github handle"])
        self.assertEqual(build.table.headers, ["Contributor", "GitHub"])
        self.assertEqual(build.table.rows[2][1].text, "-")

    def test_missing_identifiers_hide_columns(self) -> None:
        build = self.render(
            """\
            authors:
              - given-names: Ada
                family-names: Lovelace
            """
        )
        assert build.table is not None
        self.assertEqual(build.table.headers, ["Contributor"])
        self.assertEqual(build.names, ["Ada Lovelace"])

    def test_mixed_identifiers(self) -> None:
        build = self.render(
            """\
            authors:
              - given-names: Ada
                family-names: Lovelace
                orcid: https://orcid.org/0000-0002-1825-0097
              - given-names: Grace
                family-names: Hopper
                orcid: 0000-0002-9079-593X
                alias: "@ghopper"
              - given-names: Alan
                family-names: Turing
            """
        )
        assert build.table is not None
        self.assertEqual(build.table.headers, ["Contributor", "ORCID", "GitHub"])
        ada, grace, alan = build.table.rows
        self.assertEqual(ada[1].links, ["https://orcid.org/0000-0002-1825-0097"])
        self.assertEqual(ada[1].text, "0000-0002-1825-0097")
        self.assertEqual(ada[2].text, "-")
        self.assertEqual(grace[1].text, "0000-0002-9079-593X")
        self.assertEqual(grace[2].links, ["https://github.com/ghopper"])
        self.assertEqual(grace[2].text, "@ghopper")
        self.assertEqual([alan[1].text, alan[2].text], ["-", "-"])


if __name__ == "__main__":
    unittest.main()
