#!/bin/sh
# Run the tests, and refuse a green that proves nothing.
#
# usage: go-test.sh "<run-flags>" <pkg-pattern>
#
# `go test` exits 0 when the filter matches no test at all, and says so with
# "[no tests to run]". That is a green that proves nothing, so it fails here
# — v1's lesson: zero-match verification.
set -u

run=${1:-}
pkg=${2:-./...}
go_bin=${GO:-go}

tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT

# `set -f` keeps the flags from being globbed; the shell splits them into
# words, and a `|` in a -run pattern stays a character.
#
# The package comes first: the flags that are not `go test`'s own — `-update`
# for the goldens — belong to the test binary, and `go test` only passes
# them on after the package list.
set -f
if [ -n "$run" ]; then
	{ $go_bin test "$pkg" $run 2>&1; echo $? > "$tmp/status"; } | tee "$tmp/out"
else
	{ $go_bin test "$pkg" 2>&1; echo $? > "$tmp/status"; } | tee "$tmp/out"
fi
status=$(cat "$tmp/status")

if [ "$status" -ne 0 ]; then
	exit "$status"
fi

# A package that ran its tests says "ok"; a package the filter emptied says
# "ok ... [no tests to run]".
ran=$(grep '^ok' "$tmp/out" | grep -v 'no tests to run' | wc -l | tr -d ' ')
if [ "$ran" -eq 0 ]; then
	echo "the filter matched no test — a green that proves nothing." >&2
	exit 1
fi

exit 0
