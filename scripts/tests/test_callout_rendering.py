from __future__ import annotations

import subprocess
import tempfile
import unittest
from html.parser import HTMLParser
from pathlib import Path


REPO_ROOT = Path(__file__).parents[2]


class CalloutTitleParser(HTMLParser):
    VOID_ELEMENTS = {
        "area",
        "base",
        "br",
        "col",
        "embed",
        "hr",
        "img",
        "input",
        "link",
        "meta",
        "param",
        "source",
        "track",
        "wbr",
    }

    def __init__(self) -> None:
        super().__init__(convert_charrefs=True)
        self.titles: list[dict[str, object]] = []
        self._current: dict[str, object] | None = None
        self._depth = 0

    def handle_starttag(self, tag: str, attrs: list[tuple[str, str | None]]) -> None:
        attributes = dict(attrs)
        classes = (attributes.get("class") or "").split()
        if self._current is None and "lesson-callout-title" in classes:
            self._current = {"text": [], "elements": []}
            self._depth = 1
            return
        if self._current is not None:
            elements = self._current["elements"]
            assert isinstance(elements, list)
            elements.append((tag, attributes))
            if tag not in self.VOID_ELEMENTS:
                self._depth += 1

    def handle_startendtag(self, tag: str, attrs: list[tuple[str, str | None]]) -> None:
        if self._current is not None:
            elements = self._current["elements"]
            assert isinstance(elements, list)
            elements.append((tag, dict(attrs)))

    def handle_endtag(self, tag: str) -> None:
        if self._current is None:
            return
        self._depth -= 1
        if self._depth == 0:
            self.titles.append(self._current)
            self._current = None

    def handle_data(self, data: str) -> None:
        if self._current is not None:
            text = self._current["text"]
            assert isinstance(text, list)
            text.append(data)


class CalloutRenderingTests(unittest.TestCase):
    def render_titles(
        self,
        content: str,
        *,
        unsafe: bool,
        decorate_external_links: bool = False,
        inline_image_hook: bool = False,
    ) -> list[dict[str, object]]:
        with tempfile.TemporaryDirectory() as temp_name:
            site = Path(temp_name)
            (site / "content").mkdir()
            (site / "layouts" / "partials" / "shortcodes").mkdir(parents=True)

            (site / "go.mod").write_text(
                "\n".join(
                    [
                        "module example.org/callout-title-test",
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
                        'title = "Callout title test"',
                        'disableKinds = ["taxonomy", "term", "RSS", "sitemap", "robotsTXT", "404"]',
                        "",
                        "[params]",
                        f"  externalLinkDecoration = {str(decorate_external_links).lower()}",
                        "",
                        "[module]",
                        "  [[module.imports]]",
                        '    path = "github.com/oer-particle-physics/hugo-styles"',
                        "    [[module.imports.mounts]]",
                        '      source = "layouts"',
                        '      target = "layouts"',
                        "",
                        "[markup.goldmark.renderer]",
                        f"  unsafe = {str(unsafe).lower()}",
                        "",
                    ]
                ),
                encoding="utf-8",
            )
            (site / "layouts" / "home.html").write_text(
                "<!doctype html><html><body>{{ .Content }}</body></html>\n",
                encoding="utf-8",
            )
            (site / "layouts" / "partials" / "shortcodes" / "callout.html").write_text(
                '<div class="callout">{{ .content }}</div>\n',
                encoding="utf-8",
            )
            if inline_image_hook:
                (site / "layouts" / "_markup").mkdir()
                (site / "layouts" / "_markup" / "render-image.html").write_text(
                    '<span class="image-reference">{{ .PlainText }}</span>\n',
                    encoding="utf-8",
                )
            (site / "content" / "_index.md").write_text(
                f'+++\ntitle = "Callout title cases"\n+++\n\n{content}',
                encoding="utf-8",
            )

            destination = site / "public"
            completed = subprocess.run(
                [
                    "hugo",
                    "--quiet",
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
            self.assertEqual(completed.returncode, 0, completed.stderr)

            parser = CalloutTitleParser()
            parser.feed((destination / "index.html").read_text(encoding="utf-8"))
            return parser.titles

    def test_shortcode_wrappers_match(self) -> None:
        current = REPO_ROOT / "layouts" / "_shortcodes" / "callout.html"
        legacy = REPO_ROOT / "layouts" / "shortcodes" / "callout.html"
        self.assertEqual(
            current.read_text(encoding="utf-8"),
            legacy.read_text(encoding="utf-8"),
        )

    def test_titles_render_supported_inline_markdown_and_fallback(self) -> None:
        titles = self.render_titles(
            """{{< callout title="R&D > baseline's target" >}}Body.{{< /callout >}}
{{< callout title="Use `<div>` and **care**" >}}Body.{{< /callout >}}
{{< callout title="An *important* note" >}}Body.{{< /callout >}}
{{< callout title="[Setup](/setup/?a=1&b=2)" >}}Body.{{< /callout >}}
{{< callout title="<em>Trusted HTML</em>" >}}Body.{{< /callout >}}
{{< callout title="[External](https://example.net/)" >}}Body.{{< /callout >}}
{{< callout title="Line one<br>Line two" >}}Body.{{< /callout >}}
{{< callout title=`<span title="<div>">Inline</span>` >}}Body.{{< /callout >}}
{{< callout title=`Text <!-- <div> --> safe` >}}Body.{{< /callout >}}
{{< callout title="# Heading syntax" >}}Body.{{< /callout >}}
{{< callout title="![Image](/image.png)" >}}Body.{{< /callout >}}
{{< callout title="<div>Block HTML</div>" >}}Body.{{< /callout >}}
{{< callout title="<script>alert(1)</script>" >}}Body.{{< /callout >}}
{{< callout title="- List item" >}}Body.{{< /callout >}}
{{< callout title="> Quoted" >}}Body.{{< /callout >}}
{{< callout title="<table><tr><td>Cell</td></tr></table>" >}}Body.{{< /callout >}}
{{< callout title="---" >}}Body.{{< /callout >}}
{{< callout title="<iframe src='https://example.org/'></iframe>" >}}Body.{{< /callout >}}
{{< callout title="<video></video>" >}}Body.{{< /callout >}}
{{< callout title="<input type='text'>" >}}Body.{{< /callout >}}
{{< callout title="<style>em { color: red; }</style>" >}}Body.{{< /callout >}}
{{< callout title="<meta name='x' content='y'>" >}}Body.{{< /callout >}}
""",
            unsafe=True,
            decorate_external_links=True,
        )
        self.assertEqual(len(titles), 22)

        def title_text(index: int) -> str:
            text = titles[index]["text"]
            assert isinstance(text, list)
            return "".join(text)

        def title_elements(index: int) -> list[tuple[str, dict[str, str | None]]]:
            elements = titles[index]["elements"]
            assert isinstance(elements, list)
            return elements

        self.assertEqual(title_text(0), "R&D > baseline’s target")
        self.assertEqual(title_text(1), "Use <div> and care")
        self.assertEqual([tag for tag, _ in title_elements(1)], ["code", "strong"])
        self.assertEqual(title_text(2), "An important note")
        self.assertEqual([tag for tag, _ in title_elements(2)], ["em"])

        link = title_elements(3)
        self.assertEqual([tag for tag, _ in link], ["a"])
        self.assertEqual(link[0][1]["href"], "/lesson/setup/?a=1&b=2")

        self.assertEqual(title_text(4), "Trusted HTML")
        self.assertEqual([tag for tag, _ in title_elements(4)], ["em"])

        external_link = title_elements(5)
        self.assertEqual([tag for tag, _ in external_link], ["a", "svg", "path"])
        self.assertEqual(external_link[0][1]["href"], "https://example.net/")
        self.assertEqual(external_link[0][1]["target"], "_blank")

        self.assertEqual(title_text(6), "Line oneLine two")
        self.assertEqual([tag for tag, _ in title_elements(6)], ["br"])

        self.assertEqual(title_text(7), "Inline")
        self.assertEqual([tag for tag, _ in title_elements(7)], ["span"])
        self.assertEqual(title_elements(7)[0][1]["title"], "<div>")
        self.assertEqual(title_text(8), "Text  safe")
        self.assertEqual(title_elements(8), [])

        literal_fallbacks = [
            "# Heading syntax",
            "![Image](/image.png)",
            "<div>Block HTML</div>",
            "<script>alert(1)</script>",
            "- List item",
            "> Quoted",
            "<table><tr><td>Cell</td></tr></table>",
            "---",
            "<iframe src='https://example.org/'></iframe>",
            "<video></video>",
            "<input type='text'>",
            "<style>em { color: red; }</style>",
            "<meta name='x' content='y'>",
        ]
        for index, expected in enumerate(literal_fallbacks, start=9):
            self.assertEqual(title_text(index), expected)
            self.assertEqual(title_elements(index), [])

        disallowed = {
            "blockquote",
            "div",
            "h1",
            "h2",
            "h3",
            "h4",
            "h5",
            "h6",
            "hr",
            "iframe",
            "img",
            "input",
            "meta",
            "ol",
            "p",
            "pre",
            "script",
            "style",
            "table",
            "ul",
            "video",
        }
        for title in titles:
            elements = title["elements"]
            assert isinstance(elements, list)
            self.assertTrue(disallowed.isdisjoint(tag for tag, _ in elements))

    def test_raw_html_disabled_falls_back_to_literal_source(self) -> None:
        titles = self.render_titles(
            '{{< callout title="<em>Trusted HTML</em>" >}}Body.{{< /callout >}}\n',
            unsafe=False,
        )
        self.assertEqual(len(titles), 1)
        self.assertEqual(titles[0]["text"], ["<em>Trusted HTML</em>"])
        self.assertEqual(titles[0]["elements"], [])

    def test_downstream_render_hook_may_keep_image_syntax_inline(self) -> None:
        titles = self.render_titles(
            '{{< callout title="![Diagram](/diagram.png)" >}}Body.{{< /callout >}}\n',
            unsafe=False,
            inline_image_hook=True,
        )
        self.assertEqual(len(titles), 1)
        self.assertEqual(titles[0]["text"], ["Diagram"])
        self.assertEqual(
            titles[0]["elements"],
            [("span", {"class": "image-reference"})],
        )


if __name__ == "__main__":
    unittest.main()
