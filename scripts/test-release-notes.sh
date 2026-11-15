#!/bin/sh
set -eu
script_dir=$(CDPATH='' cd -- "$(dirname -- "$0")" && pwd)
test_dir=$(mktemp -d)
trap 'rm -rf "$test_dir"' EXIT HUP INT TERM
cat > "$test_dir/changelog.md" <<'CHANGELOG'
# Changelog

## [Unreleased]

## [1.2.30] - 2026-11-14

### Added

- Different version.

## [1.2.3] - 2026-11-13

### Fixed

- Correct release.

[Unreleased]: https://example.com/commits
CHANGELOG
"$script_dir/release-notes.sh" 1.2.3 "$test_dir/changelog.md" > "$test_dir/actual"
printf '\n### Fixed\n\n- Correct release.\n\n' > "$test_dir/expected"
cmp "$test_dir/expected" "$test_dir/actual"
if "$script_dir/release-notes.sh" 2.0.0 "$test_dir/changelog.md" > /dev/null 2>&1; then
    printf '%s\n' 'missing version unexpectedly accepted' >&2
    exit 1
fi
if "$script_dir/release-notes.sh" Unreleased "$test_dir/changelog.md" > /dev/null 2>&1; then
    printf '%s\n' 'empty section unexpectedly accepted' >&2
    exit 1
fi
