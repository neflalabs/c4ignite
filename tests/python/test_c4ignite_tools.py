import argparse
import io
import json
import os
import tempfile
import unittest
from pathlib import Path
from unittest.mock import patch

from scripts.python import c4ignite_tools as tools


class C4IgniteToolsTest(unittest.TestCase):
    def test_remove_setup_block(self):
        with tempfile.TemporaryDirectory() as tmpdir:
            target = Path(tmpdir) / "shell.sh"
            target.write_text(
                "line 1\n# >>> c4ignite setup >>>\n# snippet\n# <<< c4ignite setup <<<\nline 2\n",
                encoding="utf8",
            )
            args = argparse.Namespace(
                file=str(target),
                start="# >>> c4ignite setup >>>",
                end="# <<< c4ignite setup <<<",
            )
            tools.remove_setup_block(args)
            self.assertEqual(target.read_text(encoding="utf8"), "line 1\nline 2\n")

    def test_ensure_lint_script_adds_entry(self):
        with tempfile.TemporaryDirectory() as tmpdir:
            composer = Path(tmpdir) / "composer.json"
            composer.write_text("{}", encoding="utf8")
            args = argparse.Namespace(composer_file=str(composer))
            buffer = io.StringIO()
            with patch("sys.stdout", buffer):
                tools.ensure_lint_script(args)
            self.assertEqual(json.loads(composer.read_text(encoding="utf8"))["scripts"]["lint"], "vendor/bin/phpcs app")
            self.assertIn("ADDED", buffer.getvalue())

    def test_lint_script_exists(self):
        with tempfile.TemporaryDirectory() as tmpdir:
            composer = Path(tmpdir) / "composer.json"
            composer.write_text(json.dumps({"scripts": {"lint": "run me"}}), encoding="utf8")
            args = argparse.Namespace(composer_file=str(composer))
            buffer = io.StringIO()
            with patch("sys.stdout", buffer):
                tools.lint_script_exists(args)
            self.assertEqual(buffer.getvalue().strip(), "yes")

    def test_create_backup_metadata(self):
        with tempfile.TemporaryDirectory() as tmpdir:
            target = Path(tmpdir) / "meta.json"
            args = argparse.Namespace(
                output=str(target),
                includes_vendor="yes",
                includes_writable="no",
                exclude_env="yes",
                encrypted="no",
            )
            env = {
                "C4IGNITE_NOW": "2024-10-01T00:00:00Z",
                "USER": "tester",
                "PROJECT_ROOT": "/workspace",
                "C4IGNITE_GIT_COMMIT": "abc123",
            }
            with patch.dict(os.environ, env, clear=False):
                tools.create_backup_metadata(args)

            data = json.loads(target.read_text(encoding="utf8"))
            self.assertTrue(data["includes_vendor"])
            self.assertFalse(data["includes_writable"])
            self.assertTrue(data["exclude_env"])
            self.assertFalse(data["encrypted"])
            self.assertEqual(data["git_commit"], "abc123")

    def test_apply_env_defaults_appends_value(self):
        with tempfile.TemporaryDirectory() as tmpdir:
            env_file = Path(tmpdir) / ".env"
            env_file.write_text("APP_ENV=production\n", encoding="utf8")
            args = argparse.Namespace(env_file=str(env_file))
            tools.apply_env_defaults(args)
            contents = env_file.read_text(encoding="utf8")
            self.assertIn("CI_ENVIRONMENT = development", contents)

    def test_parse_release_outputs_tag(self):
        payload = {"tag_name": "v4.6.3", "tarball_url": "https://example.com/archive.tar.gz"}
        args = argparse.Namespace(field="tag")
        buffer = io.StringIO()
        stdin = io.StringIO(json.dumps(payload))
        with patch("sys.stdout", buffer), patch("sys.stdin", stdin):
            tools.parse_release(args)
        self.assertEqual(buffer.getvalue().strip(), "v4.6.3")


if __name__ == "__main__":
    unittest.main()
