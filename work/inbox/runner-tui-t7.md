Runner: T7 done. Commits: 9 (7a76e66, 0f264ac, 0dce0c1, c542916, 5ab1660,
1983b1a, af9deb1, 65fdbba, 74070ef). The row names the first eight. The line: one
row under the summary, between the summary and the list; the list gives up exactly
one row while the line is open (`listAreaHeight` drops by one), and the footer is the
last row of the screen — the new golden pins it. The keys: the list filters **live**
as the chef types, no second keypress to apply; `enter` keeps the filter and closes
the line; `esc` clears the filter and closes it, and the list is whole again. The
match: one function, `matchFilter` — case-insensitive, a substring — over a recipe's
id · title · intent (outcome and problem) · scope (in and out), and over a backlog
item's id · title · body · scope · tags. A recipe and a backlog item share the one
rule, and a test holds it for both kinds. The summary: `filter: <text> · N of M` — N
matches of M items — on the summary row, and the count comes from the same
computation the list comes from. The section headers count the items **shown**, and a
section with no match is not drawn at all. The selection: it follows the item — it
stays on the same recipe when that recipe still matches (by id, not by index), it
goes to the first match when it does not, and the list scrolls to it; the filter
survives a tab switch and applies to the new tab's list. The other tabs: Board, Log,
Home and Settings say "the filter does not apply here", and nothing else changes —
one message, in one place, and the open detail says it too, because the filter is a
list thing. The scripts: `/` opens the line; the single-character steps type one
character each; `text:<string>` types a whole word in one step, in the colon style
the mouse steps use — and it is a **word**, not a run of key presses, so outside the
filter line it does nothing and a script that meant to name a filter cannot quit the
console; `esc` clears; `enter` keeps. docs/05: a new **The filter** section — what it
is, the match, the count, the selection, the tabs it reaches, the tabs it does not
reach and why, and the script steps; `/` is in the interaction table and in the
Settings tab's key list. The goldens: the four existing goldens did **not** move —
the filter is off by default; a new golden, `console-filter.golden`, pins the
filter's screen and is claimed by a test. DomainVersion: 7. The gate: `make check`,
`make check-race` — both green. Notes: nine commits — one per criterion, plus one
that carries AC-01's missing Ready-tab test and the new golden, plus one that closes
two holes in the line's own invariant (a pasted newline is sanitised on the way in,
and a filter longer than the terminal is capped to the room the summary row has
left). The **footer does not name the filter key**: the goldens pin the footer and
the brief says they do not move, so the key is named in the Settings tab's key list
instead. While the line is open, `↑`/`↓` still move the selection and `tab` still
switches tabs (the line closes on a tab the filter does not reach, and the text is
kept). "The first match" is the first row of the filtered list in the tab's own
order. Nothing outside `internal/ui` changed: `internal/core` and `internal/cli` do
not move, and the plan's format does not move.
