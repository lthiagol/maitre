You are a **fresh runner**. You did not write T6. The repository is
`/Users/thiago/Code/Pessoal/maitre`. Read `AGENTS.md`, then this brief. The
orchestrator's pane id is `w1H:pJ`.

This is **T6's fix cycle 1**. Two findings, both ruled in. The slice's shape held:
the two drawings, the click, the full-width board, the orphan, the cycle.

## F-126 (medium) — `space` on the Board signs a recipe off

**The target:** the Board's footer names no `space`, so `space` writes nothing.

**The hole.** `primaryAction` handles the Board with `actOnRecipe`. On the seeded
plan, `down` lands on RCP02 (`ready`). `space` appends a lifecycle event, `ready` →
`approved`, and the status line says the recipe was signed off. The journal grew.
A draft on the same tab writes nothing, and so does the approved fork: the key does
not claim. The write is the chef's sign-off, on a tab that does not offer it.

**What lands:** `space` on the Board leaves `activity.json` byte for byte the same,
on a recipe the Recipes tab would sign off. The footer stays without `space`. The
Recipes tab's sign-off does not move.

**The evidence:** RCP02, `space`, the journal unchanged.

## F-127 (medium) — at 80 columns a Board row wraps the footer off

**The target:** a row that does not fit is one row, and the footer is the last row.
The same rule as F-122.

**The hole.** The wide drawing truncates. The narrow drawing does not. At 80 columns
the seeded branch rows are 82 columns, and the shape fixture's cycle row is 114.
The view's last line is the footer, and a terminal wraps the long line, so the
footer is not the last row on the screen. At 100 the same rows fit, which is why
`console-board-narrow.golden` is green.

**What lands:** at 80 columns every Board row is one line, truncated, and the footer
is the last row. The 100-column golden stays unless a row there was already past
the width — if it moves, say why.

**The evidence:** the seeded plan at 80, and the cycle fixture at 80. Say the widths.

## The acceptance criteria

| AC | What |
|--:|---|
| **AC-01** | F-126. `space` on the Board writes nothing, including on a recipe that is ready to sign off. The Recipes tab still signs off |
| **AC-02** | F-127. At 80 columns a Board row is one line and the footer is the last row, on the seeded plan and on the cycle fixture |
| **AC-03** | `work/tui/README.md`'s T6 row names the fix's commits. `DomainVersion` is **7** |

One commit per criterion. The state row's own commit is not one of the ids it names.

## The invariants

| | |
|---|---|
| **The two drawings** | below 160 the connectors stay. At 160 the trunk stays left of the fork. The click still opens the detail |
| **The wide board** | still takes the full width. The detail still opens on `enter` and on a click |
| **`core.DomainVersion`** | stays **7** |
| **The gate** | `make check`, then `make check-race` before you report |

## The scope — what not to touch

| | |
|---|---|
| **`core.TreeOf`** | the tree's membership |
| **Ready, Home, the filter** | not this cycle |
| **`docs/audit.md`** | the reviewer's. You close nothing in the ledger |
| **`work/tui/plan.md`** | the orchestrator's |
| **`work/`** | the orchestrator's, except the T6 state row and your receipt |

## The report

```bash
mkdir -p work/inbox
cat > work/inbox/runner-tui-t6-fix1.md <<EOF
Runner: T6 fix 1 done. Commits: <n>. F-126: <space on RCP02, the journal unchanged, and that Recipes still signs off>. F-127: <the widths, one line, the footer last>. The drawings: <the narrow connectors, the wide columns, the click>. The goldens: <what moved, and why>. DomainVersion: <7>. The gate: <make check, make check-race>. Notes: <anything the reviewer should know>.
EOF
git add work/inbox/runner-tui-t6-fix1.md && git commit -q -m "inbox: runner reports tui-t6 fix 1"
herdr agent prompt "w1H:pJ" "Runner: T6 fix 1 done. Commits: <n>. F-126: <space on RCP02, the journal unchanged, and that Recipes still signs off>. F-127: <the widths, one line, the footer last>. The drawings: <the narrow connectors, the wide columns, the click>. The goldens: <what moved, and why>. DomainVersion: <7>. The gate: <make check, make check-race>. Notes: <anything the reviewer should know>."
```

Write the receipt, commit it, and then nudge. If a file you did not write has
changed, stop and report it. Stage by name. Then stop.
