You are a **fresh runner**. You did not write T5. The repository is
`/Users/thiago/Code/Pessoal/maitre`. Read `AGENTS.md`, then
[work/tui/plan.md](../tui/plan.md) § 5, the Board paragraph, then this brief. This
brief wins where they differ. The orchestrator's pane id is `w1H:pJ`.

This slice is **T6**. The Board becomes a visual place. The tree's **membership
does not change**. `core.TreeOf` is still the tree. The Board still shows what
`core.Path` shows, not `core.Adoptable`.

## The two drawings

The width is the console's width.

| Width | The drawing |
|---|---|
| **< 160** | today's vertical tree, **byte for byte**. The connectors stay `├─` and `└─`. The narrow golden does not move |
| **≥ 160** | left to right. The trunk is the left column. A fork sits to the **right** of the recipe it waits on. A node shows the recipe's **title**. The dependency it waits on is still readable |

One tree, two drawings. The render chooses. A second function that guesses where the nodes are is a second answer — the hit map is emitted **by the render**, the way the other tabs already do it.

A node that does not fit is **truncated**, not wrapped. The footer is the last row at both widths.

An orphan and a recipe a cycle cannot reach are still on the board. The wide drawing does not drop them.

## The click

A click on a node selects that recipe and opens its detail — the same detail `enter` opens. It works in **both** drawings. `ui.Drive` can express the click, the way T2b's clicks are driven. The keyboard still moves down the tree and `enter` still opens the detail, at both widths.

`space` does not gain a move. The console does not claim.

## The shot

`scripts/console-shot.py` at **160×40** and at **80×24**, on a seeded plan. Quote the wide capture's trunk and one fork in the receipt. Do not commit the capture.

## The acceptance criteria

| AC | What |
|--:|---|
| **AC-01** | below 160 columns the vertical tree is unchanged, and the narrow golden does not move |
| **AC-02** | at 160 columns and above, the trunk's column is left of its fork's column, the node shows the title, an orphan and a cycle are still drawn, and the footer is the last row. A test holds the columns |
| **AC-03** | a click on a node, wide and narrow, selects that recipe and opens its detail. The click is drivable. The hit map comes from the render |
| **AC-04** | `work/tui/README.md`'s T6 row names this slice's commits. `docs/05` says what the wide board is. `DomainVersion` is **7** |

One commit per criterion. The state row's own commit is not one of the ids it names.

## The invariants

| | |
|---|---|
| **`core.DomainVersion`** | stays **7**. No format change. No new view |
| **The domain** | stdlib only. No dependency is added. `pyte` is the script's, not a module |
| **The tab count** | stays **7**. The tab is still **Board** |
| **The console never claims** | ADR-0003 |
| **Q11, Q12, Q13** | still open. Do not answer them |
| **The gate** | `make check`, then `make check-race` before you report |

## The scope — what not to touch

| | |
|---|---|
| **`core.TreeOf`, `core.Path`, `core.Adoptable`** | what the tree contains. This slice draws it |
| **Ready, Home, the filter** | not this slice |
| **`docs/07`** | the orchestrator corrected the order-tab row. Do not touch it |
| **`kitchen/`** | do not rewrite the repo's plan |
| **`docs/audit.md`** | the reviewer's |
| **`work/tui/plan.md`** | the orchestrator's |
| **`work/`** | the orchestrator's, except the T6 state row and your receipt |

## The report

```bash
mkdir -p work/inbox
cat > work/inbox/runner-tui-t6.md <<EOF
Runner: T6 done. Commits: <n>. The narrow tree: <unchanged, and the golden>. The wide tree: <the trunk's column, the fork's column, the title, the orphan, the cycle, the footer>. The click: <wide and narrow, the recipe selected, the detail open, and that the map comes from the render>. The shot: <the wide lines you quote>. docs/05: <what it says>. The goldens: <what moved, and why>. DomainVersion: <7>. The gate: <make check, make check-race>. Notes: <anything the reviewer should know>.
EOF
git add work/inbox/runner-tui-t6.md && git commit -q -m "inbox: runner reports tui-t6"
herdr agent prompt "w1H:pJ" "Runner: T6 done. Commits: <n>. The narrow tree: <unchanged, and the golden>. The wide tree: <the trunk's column, the fork's column, the title, the orphan, the cycle, the footer>. The click: <wide and narrow, the recipe selected, the detail open, and that the map comes from the render>. The shot: <the wide lines you quote>. docs/05: <what it says>. The goldens: <what moved, and why>. DomainVersion: <7>. The gate: <make check, make check-race>. Notes: <anything the reviewer should know>."
```

Write the receipt, commit it, and then nudge. If a file you did not write has
changed, stop and report it. Stage by name. Then stop.
