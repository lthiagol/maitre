#!/bin/sh
# The dependency budget.
#
# This file holds the numbers; AGENTS.md's dependency table is checked
# against them, so the doc cannot drift from the gate. One place holds the
# numbers, one place reads them — and the gate is that reader.
#
# The count is computed the way AGENTS.md says to compute it:
#
#   go list -deps -f '{{if not .Standard}}{{.Module.Path}}{{end}}' ./PKG \
#     | sort -u | grep -v lthiagol/maitre | wc -l
#
# The gate holds the numbers because the gate enforces them: a ceiling that
# lives in prose can be edited to fit the code.
set -u

ROOT=$(git rev-parse --show-toplevel)
cd "$ROOT" || exit 1

DOC=AGENTS.md
MODULE=lthiagol/maitre
LIST_TEMPLATE='{{if not .Standard}}{{.Module.Path}}{{end}}'

# package, expected count, ceiling
BUDGET='internal/core 0 0
internal/cli 0 0
internal/ui 28 30'

status=0
summary=""
failures=""

report() {
	status=1
	failures="$failures
  $1
  $2
"
}

# The doc must still state the counting method the gate uses.
if ! grep -qF "$LIST_TEMPLATE" "$DOC" || ! grep -qF "grep -v $MODULE" "$DOC"; then
	report "$DOC no longer states the counting command this gate uses." \
		"Change $DOC to state it again — the gate counts the way the doc says."
fi

while read -r pkg expected ceiling; do
	[ -n "$pkg" ] || continue

	live=$(go list -deps -f "$LIST_TEMPLATE" "./$pkg" | sort -u | grep -v "$MODULE" | wc -l | tr -d ' ')

	row=$(grep -F "| \`$pkg\` |" "$DOC" | head -1)
	if [ -z "$row" ]; then
		report "$pkg has no row in $DOC's dependency table." \
			"Add the row to $DOC — the table is the checked copy of the gate's numbers."
		continue
	fi
	doc_ceiling=$(printf '%s\n' "$row" | awk -F'|' '{gsub(/[ \t]/, "", $4); print $4}')
	doc_current=$(printf '%s\n' "$row" | awk -F'|' '{gsub(/[ \t]/, "", $5); print $5}')

	if [ "$live" -gt "$ceiling" ]; then
		report "$pkg: live count $live, ceiling $ceiling." \
			"Remove a dependency, or raise the ceiling in scripts/dep-budget.sh and $DOC — a dependency needs a written reason."
	fi

	if [ "$live" -ne "$expected" ]; then
		report "$pkg: live count $live, expected $expected (ceiling $ceiling)." \
			"Revert the dependency that moved the count, or update scripts/dep-budget.sh and $DOC — a dependency needs a written reason."
	fi

	if [ "$doc_ceiling" != "$ceiling" ]; then
		report "$DOC's dependency table states a ceiling of $doc_ceiling for $pkg; the gate holds $ceiling." \
			"Change $DOC's dependency table to match the gate."
	fi

	if [ "$doc_current" != "$expected" ]; then
		report "$DOC's dependency table states $doc_current external modules for $pkg; the gate holds $expected." \
			"Change $DOC's dependency table to match the gate."
	fi

	summary="$summary  $pkg: $live external modules (ceiling $ceiling)
"
done <<EOF
$BUDGET
EOF

if [ "$status" -ne 0 ]; then
	echo "dependency budget failed:" >&2
	printf '%s' "$failures" >&2
	echo >&2
	exit 1
fi

printf '%s' "$summary"
