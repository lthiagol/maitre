#!/bin/sh
# The slice numbers are one number.
#
# The Sous-Chef tab's slice is named in the plan's own files: the state table,
# the table that says what each ruling lands on, the console's doc, and the
# design. They must agree — **T10 is the next slice, and its writer reads
# these files**. A reader who takes one of them at its word must not be misled.
#
# The check reads them and fails when they disagree. It says **where**: the
# places and the numbers. A check that says "mismatch" without saying where is
# a nuisance, not a witness.
#
# And it refuses to pass having found nothing. A number that cannot be read
# fails the check, the way a test that matched no test fails the gate: a check
# that proves nothing must not be green.
set -u

ROOT=$(git rev-parse --show-toplevel)
cd "$ROOT" || exit 1

README=work/tui/README.md
CONSOLE=docs/05-human-interface.md
DESIGN=work/tui/plan-automation.md

# The tab's own questions — the rulings that land on the tab's slice.
#
# They are the questions the plan asks about the tab: its name and its place
# (Q16, Q17), its move (Q18), the state it reads (Q19), the journal (Q20), the
# startup task (Q21), the stop (Q22), and the ticket (Q26, Q27).
TAB_QUESTIONS="Q16 Q17 Q18 Q19 Q20 Q21 Q22 Q26 Q27"

places=""

# add records what a place says.
add() {
	places="$places$2 $1
"
}

# fail_missing fails the check: the place does not say anything the check can
# read, so the check has proved nothing.
fail_missing() {
	echo "the slice check found no number in $1." >&2
	echo "The file changed shape — the check must be taught the new shape, or it proves nothing." >&2
	exit 1
}

# state_slice is the number the state table gives the tab.
state_slice() {
	sed -nE 's/^\| \*\*T([0-9]+)\*\* \| the Sous-Chef tab \|.*/\1/p' "$README"
}

# ruling_slice is the number a ruling lands on.
ruling_slice() {
	sed -nE "s/^\\| $1 \\| \\*\\*T([0-9]+)\\*\\*.*/\\1/p" "$README"
}

# console_slice is the number the console's doc gives the tab.
console_slice() {
	sed -nE 's/^\| \*\*Sous-Chef\*\* \|.*\(T([0-9]+) — planned\).*/\1/p' "$CONSOLE"
}

# design_slice is the number the design's question table gives the tab.
design_slice() {
	sed -nE 's/^\| \*\*Q24\*\* \|.*\*\*T([0-9]+)\*\* the Sous-Chef tab.*/\1/p' "$DESIGN"
}

# slices_slice is the number the design's slices table gives the tab.
slices_slice() {
	sed -nE 's/^\| \*\*T([0-9]+) — the Sous-Chef tab\*\* \|.*/\1/p' "$DESIGN"
}

# The state table: the plan's own account of which slice is which.
n=$(state_slice)
[ -n "$n" ] || fail_missing "$README (the state table)"
add "the state table in $README" "T$n"

# The rulings: every ruling about the tab lands on the tab's slice.
for q in $TAB_QUESTIONS; do
	n=$(ruling_slice "$q")
	[ -n "$n" ] || fail_missing "$README (the ruling $q)"
	add "the ruling $q in $README" "T$n"
done

# The console's doc: the tab table.
n=$(console_slice)
[ -n "$n" ] || fail_missing "$CONSOLE (the tab table)"
add "the tab table in $CONSOLE" "T$n"

# The design: the question that ruled the slices, and the table of them.
n=$(design_slice)
[ -n "$n" ] || fail_missing "$DESIGN (the Q24 ruling)"
add "the Q24 ruling in $DESIGN" "T$n"

n=$(slices_slice)
[ -n "$n" ] || fail_missing "$DESIGN (the slices table)"
add "the slices table in $DESIGN" "T$n"

# One number, or the check fails.
distinct=$(printf '%s' "$places" | cut -d' ' -f1 | sort -u | grep -c .)
if [ "$distinct" -ne 1 ]; then
	echo "the slice numbers disagree — the Sous-Chef tab's slice is not one number:" >&2
	printf '%s' "$places" | sed 's/^/  /' >&2
	echo >&2
	echo "Fix every place to the same number, or the plan's own files will not agree." >&2
	exit 1
fi

printf '%s' "$places" | sed 's/^/  /'
echo "the Sous-Chef tab's slice is one number: $(printf '%s' "$places" | cut -d' ' -f1 | head -1)"
