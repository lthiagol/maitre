You are a **fresh runner**. You did not write T6. The repository is
`/Users/thiago/Code/Pessoal/maitre`. Read `AGENTS.md`, then
[work/tui/plan.md](../tui/plan.md) § 6, then this brief. This brief wins where they
differ. The orchestrator's pane id is `w1H:p2`.

This slice is **T7**. The filter. `/` filters the list of the current tab.

## The filter

| | |
|---|---|
| **Opens** | `/` on **Recipes**, **Backlog** and **Ready** — the tabs that list plan items |
| **The line** | one line under the summary, while the chef is typing in it |
| **Live** | the list filters as the chef types. No second keypress to apply |
| **`enter`** | keeps the filter and closes the line |
| **`esc`** | clears the filter and closes the line. The list is whole again |
| **The match** | case-insensitive, over a **recipe's** id · title · intent (outcome and problem) · scope (in and out), and over a **backlog item's** id · title · body · scope · tags |
| **The count** | the summary says the filter is on and on what: `filter: <text> · N of M` — N matches of M items |
| **The selection** | stays on the same item when it still matches, goes to the first match when it does not, and the list scrolls to it |
| **The tabs** | the filter **survives a tab switch** and applies to the new tab's list |

**One rule, in one place.** The match is a single function in `internal/ui`, used by
every tab. A recipe and a backlog item are matched by the same case-insensitive
substring rule — two implementations of "does this match" is the house's standing
failure, and the count in the summary must come from the same computation as the
list it counts.

**The filter line takes exactly one row.** The list loses one row while the line is
open, and the footer stays the last row of the screen. A line that pushes the footer
off is a line that broke the screen.

**The section headers count what is shown.** A triage section with two matching
recipes says `2`, not the whole section's count. A section with no match is not
drawn at all — an empty section header is noise that says the plan has something it
does not have.

## The tabs it does not apply to

**Board**, **Log**, **Home** and **Settings**. `/` there **says the filter does not
apply here** and changes nothing else.

The Board is a tree: a node is a relationship, and hiding one misstates what waits
on what. The Log holds journal events, not plan items. Home holds numbers and
Settings holds the plan's identity.

**A key that does nothing must say so.** The console does not promise a filter it
will not honour.

## The scripts

The filter is drivable, like every other key:

| Step | What |
|---|---|
| `/` | opens the line |
| `t`, `u`, `i` | single characters, one step each — the way `KeyFor` already works |
| `text:tui` | **new** — a whole word in one step, in the colon style the mouse steps already use |
| `esc` | clears it |
| `enter` | keeps it |

`text:` is a small addition to `internal/ui/keys.go`. The single-character steps
keep working.

## The acceptance criteria

| AC | What |
|--:|---|
| **AC-01** | on Recipes, Backlog and Ready, `/` opens a filter line under the summary; typing filters the list live; `enter` keeps the filter and closes the line; `esc` clears the filter and closes it. The list loses exactly one row while the line is open, and the footer is the last row |
| **AC-02** | the match is one rule, case-insensitive, over the recipe's id · title · intent · scope and the backlog item's id · title · body · scope · tags. A test holds the rule for both kinds |
| **AC-03** | the summary line says the filter is on and on what — `filter: <text> · N of M` — and the section headers count the items **shown**. A section with no match is not drawn |
| **AC-04** | the selection stays on the same item when it still matches, goes to the first match when it does not, and the list scrolls to it. The filter survives a tab switch and applies to the new tab |
| **AC-05** | on Board, Log, Home and Settings, `/` says the filter does not apply here, and nothing else changes |
| **AC-06** | the filter is drivable from a script: `/`, the single-character steps, and `text:<string>` |
| **AC-07** | `work/tui/README.md`'s T7 row names this slice's commits. `docs/05` says what the filter is and which tabs it covers. `DomainVersion` is **7** |

One commit per criterion. The state row's own commit is not one of the ids it names.

## The invariants

| | |
|---|---|
| **`core.DomainVersion`** | stays **7**. No format change |
| **No new view** | the match rule is the console's answer — plan § 6, option B. Nothing is registered in `core.DerivedViews`, so `docs/03` § 8 does not move |
| **The domain** | stdlib only. `internal/core` and `internal/cli` do not change at all |
| **The tab count** | stays **7**. No tab is added, renamed or removed |
| **The goldens** | the filter is off by default, so **the existing goldens do not move**. A new golden is allowed, and it is claimed by a test and committed |
| **Q11, Q12, Q13** | still open. Do not answer them |
| **The gate** | `make check`, then `make check-race` before you report |

## The scope — what not to touch

| | |
|---|---|
| **`internal/core`** | the filter is a console feature. If you think the domain needs a change, **report it** — do not make it |
| **The Board, the Log, Home, Settings** | the filter does not reach them. Do not filter them "as a bonus" |
| **Tags on the recipe** | the plan's option A. Later, with the CLI. A recipe has no tags today, and this slice does not give it one |
| **`docs/03`, `docs/06`, `docs/07`** | not this slice |
| **`kitchen/`** | do not rewrite the repo's plan |
| **`docs/audit.md`** | the reviewer's |
| **`work/tui/plan.md`** | the orchestrator's |
| **`work/`** | the orchestrator's, except the T7 state row and your receipt |

## The report

```bash
mkdir -p work/inbox
cat > work/inbox/runner-tui-t7.md <<'EOF'
Runner: T7 done. Commits: <n>. The line: <where it sits, the row the list lost, and that the footer is the last row>. The keys: <the live filter, enter keeps, esc clears>. The match: <the one rule, the fields, and that a recipe and an item share it>. The summary: <the filter line, N of M, and the section counts>. The selection: <it stays, it moves, it scrolls, and the tab switch>. The other tabs: <what Board, Log, Home and Settings say>. The scripts: <the steps you drove>. docs/05: <what it says>. The goldens: <what moved, and why>. DomainVersion: <7>. The gate: <make check, make check-race>. Notes: <anything the reviewer should know>.
EOF
git add work/inbox/runner-tui-t7.md && git commit -q -m "inbox: runner reports tui-t7"
herdr agent prompt "w1H:p2" "<the receipt, as above>"
```

Write the receipt, commit it, and then nudge. If a file you did not write has
changed, stop and report it. Stage by name. Then stop.
