#!/bin/sh
# The golden files' claims.
#
# usage: check-golden-claims.sh mark <dir> | report <dir>
#
# A golden file that no test writes is a layout nothing pins — and the check
# that says "the goldens are current" would pass without touching it.
#
# The claim is derived, not declared: `mark` moves every golden's timestamp
# into the past, the golden run rewrites the ones a test claims, and
# `report` names the ones nobody wrote. No manifest of test names, and no
# scan for a path in a test file — the paths live in constants, and a scan
# would be a guess.
#
# A failure is loud and names the file. Two causes, one message: no test
# claims the file, or the golden run's filter misses the test that does —
# and the second is a real hole too, because `make goldens` would not
# regenerate the file either.
set -u

# A timestamp in the past, so "not newer than this" means "not written".
MARK=200001010000

phase=${1:-}
dir=${2:-}

case "$phase" in
mark)
	find "$dir" -name '*.golden' -type f -exec touch -t "$MARK" {} +
	;;
report)
	ref=$(mktemp)
	touch -t "$MARK" "$ref"
	orphans=$(find "$dir" -name '*.golden' -type f ! -newer "$ref")
	rm -f "$ref"
	if [ -n "$orphans" ]; then
		echo "these golden files are claimed by no test:" >&2
		printf '%s\n' "$orphans" >&2
		echo "No test in the golden run writes them: either no test claims them, or the run's filter misses the test that does." >&2
		exit 1
	fi
	echo "  every golden file is claimed by a test"
	;;
*)
	echo "usage: check-golden-claims.sh mark <dir> | report <dir>" >&2
	exit 2
	;;
esac
