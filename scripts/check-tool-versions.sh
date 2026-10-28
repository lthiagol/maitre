#!/bin/sh
# Check the tool versions — the binary the gate runs, and the doc that
# names them.
#
# usage:
#   check-tool-versions.sh tool <tool> <pin>
#   check-tool-versions.sh doc  <doc> <tool> <pin>
#
# The pin lives in the Makefile; the doc's copy is checked against it, so
# the doc cannot drift from the gate — the same shape as the dependency
# budget: the gate holds the number, the doc is the checked copy.
#
# A gate that runs a tool it did not pin is a gate nobody knows what it
# ran.
#
# The scripts call each other by a relative path: they are make's
# implementation, and make runs its recipes from the repo root.
set -u

mode=${1:-}
shift || true

go_bin=${GO:-go}

case "$mode" in
tool)
	tool=${1:-}
	pin=${2:-}
	out=$(scripts/require-tool.sh "$tool" -version 2>&1) || {
		printf '%s\n' "$out" >&2
		exit 1
	}
	# staticcheck prints "staticcheck 2026.2.2 (0.8.2)"; govulncheck prints
	# "Scanner: govulncheck@v1.9.0". The parenthesised token wins; the
	# @-suffix is the fallback.
	got=$(printf '%s\n' "$out" | sed -n 's/.*(\([^()]*\)).*/\1/p' | head -1)
	if [ -z "$got" ]; then
		got=$(printf '%s\n' "$out" | sed -n 's/.*@\([^[:space:]]*\).*/\1/p' | head -1)
	fi
	want=${pin#v}
	got=${got#v}
	if [ "$got" != "$want" ]; then
		gopath=$($go_bin env GOPATH)
		echo "$tool is version '${got:-unknown}'; the pin is $pin." >&2
		echo "Run 'make tools' to install the pinned version — it lands in $gopath/bin." >&2
		echo "Put $gopath/bin first on your PATH: an older $tool earlier there still fails." >&2
		exit 1
	fi
	echo "  $tool $got matches the pin $pin"
	;;
doc)
	doc=${1:-}
	tool=${2:-}
	pin=${3:-}
	want=${pin#v}
	# Every mention that names a version is checked — not only the first.
	# Only the pin's shape counts: the tool reports two versions (the
	# release and the module), the pin is the module one, and a line that
	# carries the other version must not fail the gate.
	versions=$(grep -F "$tool" "$doc" | grep -E 'v[0-9]+\.[0-9]+\.[0-9]+' | grep -oE 'v[0-9]+\.[0-9]+\.[0-9]+' || true)
	if [ -z "$versions" ]; then
		echo "$doc names no version for $tool." >&2
		echo "Change $doc to name $pin." >&2
		exit 1
	fi
	bad=""
	for v in $versions; do
		v=${v#v}
		if [ "$v" != "$want" ]; then
			if [ -z "$bad" ]; then
				bad="$v"
			else
				bad="$bad, $v"
			fi
		fi
	done
	if [ -n "$bad" ]; then
		echo "$doc names $tool $bad; the pin is $pin." >&2
		echo "Change $doc to name $pin." >&2
		exit 1
	fi
	echo "  $doc names $tool $want, as pinned"
	;;
*)
	echo "usage: check-tool-versions.sh tool <tool> <pin> | doc <doc> <tool> <pin>" >&2
	exit 2
	;;
esac
