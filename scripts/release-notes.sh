#!/bin/sh
# Print the changelog section for a release; missing or empty notes are an error.
set -eu
if [ "$#" -lt 1 ] || [ "$#" -gt 2 ]; then
    printf '%s\n' 'usage: release-notes.sh VERSION [CHANGELOG]' >&2
    exit 2
fi
version=$1
changelog=${2:-CHANGELOG.md}
awk -v version="$version" '
    /^## / {
        if (active) exit
        heading = "## [" version "]"
        if (index($0, heading) == 1 && (length($0) == length(heading) || substr($0, length(heading) + 1, 1) == " ")) {
            found = 1
            active = 1
        }
        next
    }
    active && /^\[[^]]+\]:[[:space:]]/ { exit }
    active {
        print
        if ($0 ~ /[^[:space:]]/ && $0 !~ /^### /) content = 1
    }
    END {
        if (!found || !content) {
            print "Missing or empty changelog section for " version > "/dev/stderr"
            exit 1
        }
    }
' "$changelog"
