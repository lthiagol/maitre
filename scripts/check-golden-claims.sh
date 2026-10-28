#!/bin/sh
# The golden files' claims.
#
# usage: check-golden-claims.sh mark <dir> | report <dir> <filter>
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
# The scope is the directory it is given: widen it where the goldens live.
#
# A failure is loud and names the file, the cause and the remedy. Two
# causes, two remedies:
#
#   - no test source mentions the file — write a test that pins it, or
#     delete it: a golden nothing pins is a layout nothing pins;
#   - a test source mentions it, but the golden run's filter misses the
#     test that writes it — name the test so the filter matches it, or
#     extend the filter so it matches the test.
set -u

# A timestamp in the past, so "not newer than this" means "not written".
MARK=200001010000

phase=${1:-}
dir=${2:-}
filter=${3:-Golden}

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
		for f in $orphans; do
			name=$(basename "$f")
			mentioned=$(find . -name '*_test.go' -not -path './.git/*' \
				-exec grep -lF -- "$name" {} + 2>/dev/null || true)
			if [ -n "$mentioned" ]; then
				echo "$f: a test source mentions it, but the golden run's filter (-run $filter) does not match the test that writes it." >&2
				echo "  Name the test so the filter matches it — its name must carry '$filter' — or extend the filter so it matches the test." >&2
			else
				echo "$f: no test source mentions it, so no test claims it." >&2
				echo "  Write a test that pins it, or delete the file: a golden nothing pins is a layout nothing pins." >&2
			fi
		done
		exit 1
	fi
	echo "  every golden file is claimed by a test"
	;;
*)
	echo "usage: check-golden-claims.sh mark <dir> | report <dir> <filter>" >&2
	exit 2
	;;
esac
