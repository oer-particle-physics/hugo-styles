from __future__ import annotations

import importlib.util
import json
import os
import subprocess
import sys
import tempfile
import unittest
from pathlib import Path
from unittest import mock


SCRIPT_PATH = Path(__file__).parents[1] / "build-versioned-site.py"
SPEC = importlib.util.spec_from_file_location("build_versioned_site", SCRIPT_PATH)
assert SPEC and SPEC.loader
builder = importlib.util.module_from_spec(SPEC)
sys.modules[SPEC.name] = builder
SPEC.loader.exec_module(builder)


class URLTests(unittest.TestCase):
    def test_normalize_and_join_base_url(self) -> None:
        self.assertEqual(
            builder.normalize_base_url("https://example.org/lesson"),
            "https://example.org/lesson/",
        )
        self.assertEqual(
            builder.join_url("https://example.org/lesson/", "versions", "v1.2.0"),
            "https://example.org/lesson/versions/v1.2.0/",
        )


class ConfigTests(unittest.TestCase):
    def test_missing_secondary_config_is_only_allowed_for_historical_builds(self) -> None:
        with tempfile.TemporaryDirectory() as temp_name:
            site = Path(temp_name)
            (site / "hugo.toml").write_text("title = 'Example'\n", encoding="utf-8")

            with self.assertRaisesRegex(builder.BuildError, "hugo-docs.toml"):
                builder.load_hugo_config(
                    site,
                    "hugo.toml,hugo-docs.toml",
                    "hugo",
                    None,
                    site / "cache",
                )

            config = builder.load_hugo_config(
                site,
                "hugo.toml,hugo-docs.toml",
                "hugo",
                None,
                site / "cache",
                allow_missing_config_files=True,
            )
            self.assertEqual(config["title"], "Example")

    def test_importing_module_does_not_inherit_docs_banner(self) -> None:
        repo_root = SCRIPT_PATH.parents[1]
        with tempfile.TemporaryDirectory() as temp_name:
            site = Path(temp_name)
            (site / "go.mod").write_text(
                "\n".join(
                    [
                        "module example.org/downstream-lesson",
                        "",
                        "go 1.26",
                        "",
                        "require github.com/oer-particle-physics/hugo-styles v0.0.0",
                        "",
                        "replace github.com/oer-particle-physics/hugo-styles "
                        f"=> {repo_root}",
                        "",
                    ]
                ),
                encoding="utf-8",
            )
            (site / "hugo.toml").write_text(
                "\n".join(
                    [
                        'baseURL = "https://example.org/lesson/"',
                        'title = "Downstream lesson"',
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

            completed = subprocess.run(
                ["hugo", "config", "--quiet", "--format", "json"],
                cwd=site,
                check=True,
                capture_output=True,
                text=True,
            )
            config = json.loads(completed.stdout)
            self.assertNotIn("banner", config.get("params", {}))


class DestinationTests(unittest.TestCase):
    def test_rejects_site_root_and_ancestor(self) -> None:
        with tempfile.TemporaryDirectory() as temp_name:
            repo = Path(temp_name).resolve()
            site = repo / "site"
            site.mkdir()

            for destination in (site, repo, repo.parent):
                with self.subTest(destination=destination):
                    with self.assertRaises(builder.BuildError):
                        builder.validate_destination(
                            site_root=site,
                            repo_root=repo,
                            destination=destination,
                            allow_external=True,
                            force=True,
                        )

    def test_requires_explicit_external_destination(self) -> None:
        with tempfile.TemporaryDirectory() as temp_name:
            root = Path(temp_name).resolve()
            site = root / "repo" / "site"
            external = root / "external"
            site.mkdir(parents=True)

            with self.assertRaisesRegex(builder.BuildError, "outside the site root"):
                builder.validate_destination(
                    site_root=site,
                    repo_root=site.parent,
                    destination=external,
                    allow_external=False,
                    force=False,
                )

            builder.validate_destination(
                site_root=site,
                repo_root=site.parent,
                destination=external,
                allow_external=True,
                force=False,
            )

    def test_rejects_symlink_escape_even_with_force(self) -> None:
        with tempfile.TemporaryDirectory() as temp_name:
            root = Path(temp_name).resolve()
            site = root / "site"
            external = root / "external"
            site.mkdir()
            external.mkdir()
            (site / "linked-output").symlink_to(external, target_is_directory=True)

            with self.assertRaisesRegex(builder.BuildError, "symbolic link"):
                builder.validate_destination(
                    site_root=site,
                    repo_root=site,
                    destination=site / "linked-output" / "public",
                    allow_external=True,
                    force=True,
                )

    def test_requires_force_for_unmarked_nonempty_destination(self) -> None:
        with tempfile.TemporaryDirectory() as temp_name:
            site = Path(temp_name).resolve()
            destination = site / "public"
            destination.mkdir()
            (destination / "index.html").write_text("old", encoding="utf-8")

            with self.assertRaisesRegex(builder.BuildError, "--force"):
                builder.validate_destination(
                    site_root=site,
                    repo_root=site,
                    destination=destination,
                    allow_external=False,
                    force=False,
                )

            builder.validate_destination(
                site_root=site,
                repo_root=site,
                destination=destination,
                allow_external=False,
                force=True,
            )
            (destination / builder.OUTPUT_MARKER).write_text("managed", encoding="utf-8")
            builder.validate_destination(
                site_root=site,
                repo_root=site,
                destination=destination,
                allow_external=False,
                force=False,
            )

    def test_publish_replaces_output_and_keeps_marker(self) -> None:
        with tempfile.TemporaryDirectory() as temp_name:
            root = Path(temp_name)
            destination = root / "public"
            staging = root / "staging"
            destination.mkdir()
            staging.mkdir()
            (destination / "old.html").write_text("old", encoding="utf-8")
            (staging / "new.html").write_text("new", encoding="utf-8")
            (staging / builder.OUTPUT_MARKER).write_text("managed", encoding="utf-8")

            builder.publish_staged_output(staging, destination)

            self.assertFalse(staging.exists())
            self.assertFalse((destination / "old.html").exists())
            self.assertEqual((destination / "new.html").read_text(), "new")
            self.assertTrue((destination / builder.OUTPUT_MARKER).is_file())

    def test_publish_restores_previous_output_when_swap_fails(self) -> None:
        with tempfile.TemporaryDirectory() as temp_name:
            root = Path(temp_name)
            destination = root / "public"
            staging = root / "staging"
            destination.mkdir()
            staging.mkdir()
            (destination / "old.html").write_text("old", encoding="utf-8")

            real_replace = os.replace
            calls = 0

            def fail_second_replace(source: Path, target: Path) -> None:
                nonlocal calls
                calls += 1
                if calls == 2:
                    raise OSError("simulated failure")
                real_replace(source, target)

            with mock.patch.object(builder.os, "replace", side_effect=fail_second_replace):
                with self.assertRaises(builder.BuildError):
                    builder.publish_staged_output(staging, destination)

            self.assertEqual((destination / "old.html").read_text(), "old")
            self.assertTrue(staging.exists())


class TargetTests(unittest.TestCase):
    def test_project_base_url_is_used_for_version_menu_links(self) -> None:
        versioning = {
            "enable": True,
            "defaultbranch": "main",
            "latest": {"enable": True, "label": "Latest"},
            "branches": {},
            "tags": {"refs": ["v1.2.0"]},
        }
        with (
            mock.patch.object(
                builder,
                "list_branch_refs",
                return_value=(["main"], {"main": "refs/heads/main"}),
            ),
            mock.patch.object(
                builder,
                "list_tag_refs",
                return_value=(["v1.2.0"], {"v1.2.0": "refs/tags/v1.2.0"}),
            ),
        ):
            targets = builder.resolve_targets(
                versioning=versioning,
                base_url="https://example.org/lesson/",
                destination_root=Path("/tmp/output"),
                repo_root=Path("/tmp/repo"),
                current_branch="main",
            )

        self.assertEqual(targets[0].menu_path, "/lesson/")
        self.assertEqual(targets[1].menu_path, "/lesson/versions/v1.2.0/")
        root_menu = builder.apply_version_menu({}, versioning, targets)
        tag_menu = builder.apply_version_menu({}, versioning, targets)
        placeholders = [builder.version_menu_placeholder(target) for target in targets]
        self.assertEqual(
            [item.get("url") for item in root_menu["menus"]["main"][1:]],
            placeholders,
        )
        self.assertEqual(
            [item.get("url") for item in tag_menu["menus"]["main"][1:]],
            placeholders,
        )

        with tempfile.TemporaryDirectory() as temp_name:
            output = Path(temp_name)
            page = output / "docs" / "page" / "index.html"
            page.parent.mkdir(parents=True)
            page.write_text(
                f'<a href="{placeholders[0]}">Latest</a>'
                f'<a href="{placeholders[1]}">v1.2.0 quoted</a>'
                f'<a href={placeholders[1]}>v1.2.0 minified</a>',
                encoding="utf-8",
            )
            builder.rewrite_version_menu_links(output, targets)
            self.assertEqual(
                page.read_text(encoding="utf-8"),
                '<a href="/lesson/">Latest</a>'
                '<a href="/lesson/versions/v1.2.0/">v1.2.0 quoted</a>'
                '<a href=/lesson/versions/v1.2.0/>v1.2.0 minified</a>',
            )

    def test_current_checkout_avoids_primary_ref_in_detached_head(self) -> None:
        versioning = {
            "enable": True,
            "defaultbranch": "main",
            "latest": {"enable": True, "label": "Latest"},
            "branches": {},
            "tags": {},
        }
        with (
            mock.patch.object(
                builder,
                "list_branch_refs",
                return_value=(["main"], {"main": "refs/remotes/origin/main"}),
            ),
            mock.patch.object(builder, "list_tag_refs", return_value=([], {})),
        ):
            configured = builder.resolve_targets(
                versioning=versioning,
                base_url="/",
                destination_root=Path("/tmp/output"),
                repo_root=Path("/tmp/repo"),
                current_branch=None,
                use_current_checkout=False,
            )
            current = builder.resolve_targets(
                versioning=versioning,
                base_url="/",
                destination_root=Path("/tmp/output"),
                repo_root=Path("/tmp/repo"),
                current_branch=None,
                use_current_checkout=True,
            )

        self.assertEqual(configured[0].git_ref, "refs/remotes/origin/main")
        self.assertIsNone(current[0].git_ref)

    def test_target_config_starts_from_each_worktree_config(self) -> None:
        target = builder.BuildTarget(
            name="latest",
            label="Latest",
            kind="root",
            git_ref=None,
            base_url="/",
            menu_path="/",
            destination=Path("/tmp/output"),
            include_in_menu=True,
        )
        per_ref = {"title": "Historical title", "menus": {"main": []}}
        with mock.patch.object(builder, "load_hugo_config", return_value=per_ref) as load:
            config = builder.target_config(
                build_root=Path("/tmp/worktree"),
                config_name="hugo.toml",
                hugo_bin="hugo",
                cache_dir=Path("/tmp/cache"),
                versioning={},
                menu_targets=[target],
            )

        self.assertEqual(config["title"], "Historical title")
        self.assertEqual(load.call_args.args[0], Path("/tmp/worktree"))
        self.assertEqual(config["menus"]["main"][0]["name"], "Versions")


if __name__ == "__main__":
    unittest.main()
