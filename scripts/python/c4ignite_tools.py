#!/usr/bin/env python3

import argparse
import json
import os
import sys
from pathlib import Path

__all__ = [
    "remove_setup_block",
    "ensure_lint_script",
    "lint_script_exists",
    "create_backup_metadata",
    "apply_env_defaults",
    "parse_release",
    "main",
]


def _load_text_lines(path: Path):
    try:
        with path.open("r", encoding="utf8") as handle:
            return handle.readlines()
    except FileNotFoundError:
        return None


def remove_setup_block(args: argparse.Namespace) -> int:
    path = Path(args.file)
    lines = _load_text_lines(path)
    if lines is None:
        return 0

    start = args.start.strip()
    end = args.end.strip()

    output = []
    skipping = False
    for line in lines:
        stripped = line.strip()
        if stripped == start:
            skipping = True
            continue
        if skipping and stripped == end:
            skipping = False
            continue
        if not skipping:
            output.append(line)

    with path.open("w", encoding="utf8") as handle:
        handle.writelines(output)
    return 0


def ensure_lint_script(args: argparse.Namespace) -> int:
    path = Path(args.composer_file)
    with path.open("r", encoding="utf8") as handle:
        data = json.load(handle)

    scripts = data.setdefault("scripts", {})
    changed = False

    if "lint" not in scripts:
        scripts["lint"] = "vendor/bin/phpcs app"
        changed = True

    if changed:
        with path.open("w", encoding="utf8") as handle:
            json.dump(data, handle, indent=2)
            handle.write("\n")
        print("ADDED")
    else:
        print("EXISTS")
    return 0


def lint_script_exists(args: argparse.Namespace) -> int:
    path = Path(args.composer_file)
    if not path.is_file():
        print("no")
        return 0

    with path.open("r", encoding="utf8") as handle:
        data = json.load(handle)

    scripts = data.get("scripts", {})
    print("yes" if "lint" in scripts else "no")
    return 0


def _str_to_bool(value: str) -> bool:
    return value.lower() in {"1", "true", "yes", "y", "on"}


def create_backup_metadata(args: argparse.Namespace) -> int:
    payload = {
        "created_at": os.environ.get("C4IGNITE_NOW", ""),
        "hostname": os.uname().nodename,
        "user": os.environ.get("USER", ""),
        "project_path": os.environ.get("PROJECT_ROOT", ""),
        "includes_vendor": _str_to_bool(args.includes_vendor),
        "includes_writable": _str_to_bool(args.includes_writable),
        "exclude_env": _str_to_bool(args.exclude_env),
        "encrypted": _str_to_bool(args.encrypted),
        "git_commit": os.environ.get("C4IGNITE_GIT_COMMIT", ""),
    }

    path = Path(args.output)
    with path.open("w", encoding="utf8") as handle:
        json.dump(payload, handle, indent=2)
        handle.write("\n")
    return 0


def apply_env_defaults(args: argparse.Namespace) -> int:
    path = Path(args.env_file)
    lines = _load_text_lines(path)
    if lines is None:
        return 0

    found = False
    result = []
    for line in lines:
        stripped = line.lstrip()
        if stripped.startswith("CI_ENVIRONMENT") or stripped.startswith("# CI_ENVIRONMENT"):
            if not found:
                result.append("CI_ENVIRONMENT = development\n")
                found = True
            continue
        result.append(line)

    if not found:
        if result and not result[-1].endswith("\n"):
            result[-1] = result[-1] + "\n"
        result.append("\nCI_ENVIRONMENT = development\n")

    with path.open("w", encoding="utf8") as handle:
        handle.writelines(result)
    return 0


def parse_release(args: argparse.Namespace) -> int:
    try:
        data = json.load(sys.stdin)
    except json.JSONDecodeError as exc:
        raise SystemExit(f"Failed to parse release JSON: {exc}") from exc

    if args.field == "tag":
        value = data.get("tag_name") or data.get("name") or ""
    elif args.field == "tarball":
        value = data.get("tarball_url") or ""
    else:
        raise SystemExit(f"Unsupported release field: {args.field}")

    print(value)
    return 0


def main(argv: list[str] | None = None) -> int:
    parser = argparse.ArgumentParser(prog="c4ignite-tools")
    subparsers = parser.add_subparsers(dest="command", required=True)

    remove_parser = subparsers.add_parser("remove-setup-block")
    remove_parser.add_argument("file")
    remove_parser.add_argument("start")
    remove_parser.add_argument("end")
    remove_parser.set_defaults(func=remove_setup_block)

    ensure_parser = subparsers.add_parser("ensure-lint-script")
    ensure_parser.add_argument("composer_file")
    ensure_parser.set_defaults(func=ensure_lint_script)

    exists_parser = subparsers.add_parser("lint-script-exists")
    exists_parser.add_argument("composer_file")
    exists_parser.set_defaults(func=lint_script_exists)

    metadata_parser = subparsers.add_parser("create-backup-metadata")
    metadata_parser.add_argument("output")
    metadata_parser.add_argument("includes_vendor")
    metadata_parser.add_argument("includes_writable")
    metadata_parser.add_argument("exclude_env")
    metadata_parser.add_argument("encrypted")
    metadata_parser.set_defaults(func=create_backup_metadata)

    env_parser = subparsers.add_parser("apply-env-defaults")
    env_parser.add_argument("env_file")
    env_parser.set_defaults(func=apply_env_defaults)

    release_parser = subparsers.add_parser("parse-release")
    release_parser.add_argument("field", choices=("tag", "tarball"))
    release_parser.set_defaults(func=parse_release)

    args = parser.parse_args(argv)
    return args.func(args)


if __name__ == "__main__":
    sys.exit(main())
