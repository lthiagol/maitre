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
set -u

mode=${1:-}
shift || true

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
		echo "$tool is version '${got:-unknown}'; the pin is $pin." >&2
		echo "Run 'make tools' to install the pinned version." >&2
		exit 1
	fi
	echo "  $tool $got matches the pin $pin"
	;;
doc)
	doc=${1:-}
	tool=${2:-}
	pin=${3:-}
	row=$(grep -F "\`$tool\`" "$doc" | grep -E 'v?[0-9]+\.[0-9]+\.[0-9]+' | head -1)
	if [ -z "$row" ]; then
		echo "$doc names no version for $tool." >&2
		echo "Change $doc to name $pin." >&2
		exit 1
	fi
	stated=$(printf '%s\n' "$row" | grep -oE 'v?[0-9]+\.[0-9]+\.[0-9]+' | head -1)
	want=${pin#v}
	stated=${stated#v}
	if [ "$stated" != "$want" ]; then
		echo "$doc names $tool $stated; the pin is $pin." >&2
		echo "Change $doc to name $pin." >&2
		exit 1
	fi
	echo "  $doc names $tool $stated, as pinned"
	;;
*)
	echo "usage: check-tool-versions.sh tool <tool> <pin> | doc <doc> <tool> <pin>" >&2
	exit 2
	;;
esac
