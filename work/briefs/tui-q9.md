You are a **fresh runner**. You did not write T4. The repository is
`/Users/thiago/Code/Pessoal/maitre`. Read `AGENTS.md`, then this brief. The
orchestrator's pane id is `w1H:pJ`; use it verbatim in the notify command.

This slice lands two rulings from 2026-10-31. Nothing else.

| # | Ruling |
|--:|---|
| **Q9** | an empty `scope` or `tags` on a backlog item is written **`[]`**. `DomainVersion` stays **7** |
| **Q10** | after a write, the selection **follows the item**. A selection that points at another item is a lie. The fix belongs in the model |

## Q9 — the list on disk

**The hole.** `BacklogItem.Scope` and `.Tags` are tagged `omitempty`. An empty list
is left off the file. A reader that asks for the list gets nothing. A recipe's
empty lists are written `[]`. One record, two answers.

**What lands.** A save writes `"scope": []` and `"tags": []` when the lists are
empty. A nil list and an empty list are the same record on disk. The guarantee
lives in **one place**, so a caller that passes nil cannot forget it.

A file that **lacks** the keys still loads. The lists read as empty. The next save
writes `[]`.

**`DomainVersion` stays 7.** The shape fingerprint hashes field **names**, and
`scope` and `tags` are already among them (`internal/core/format.go` takes the tag's
part before the comma). Removing `omitempty` does not move the fingerprint. If
`TestDomainShapeIsVersioned` goes red, **stop and report it**. Do not bump the
version.

**The evidence.** The raw file, not the loaded struct. A created item with no scope
and no tags contains both keys, each `[]`. An update that clears the lists does
too. A file without the keys loads, and saving it puts the keys back.

## Q10 — the selection

**The hole.** The backlog cursor is an index. A shaping moves an idea into *Small
tasks*, the index stays, and the footer now offers another item's move.

**What lands.** The cursor remembers the item's **id**, and a reload resolves that
id to its row.

| Write | The cursor |
|---|---|
| **Shape** | the same id, on its new row. An idea shaped into a small task does not leave the cursor on the idea that slid into the old index |
| **Create** | the new id |
| **Dismiss, promote** | the id has left the backlog. The cursor does not still name it. It sits on the nearest remaining row. It does **not** switch tabs, and it does not follow a promotion into the recipe |

A reload that does not remove the id leaves that id selected. The recipe cursor,
the note cursor, the scroll, the theme, the order, the served checkbox and the tab
still survive a reload — `TestTheChefsPlaceSurvivesAReload` stays green.

## The acceptance criteria

| AC | What | The commit |
|--:|---|---|
| **AC-01** | the raw file carries `"scope": []` and `"tags": []` for an empty list, including a nil that was saved, and a file that lacked the keys loads and the next save writes them | the writer |
| **AC-02** | after a shape, the cursor is on that id; after a create, on the new id; after a dismiss and a promote, not on the id that left | the selection |
| **AC-03** | `work/tui/README.md`'s Q9 row names this slice's commits | the state row |

One commit per criterion, as it goes. The state row's own commit is not one of the
ids it names.

## The invariants

| | |
|---|---|
| **`core.DomainVersion`** | stays **7** |
| **The domain** | `internal/core` and `internal/cli` stay stdlib-only. No dependency is added |
| **The tab count** | stays **7** |
| **The console never authors a recipe** | a promotion still does what it does today. The cursor does not chase the new recipe |
| **Q8** | the free edit stays unbuilt |
| **Q11, Q12, Q13** | still open. Do not answer them |
| **The gate** | `make check`, then `make check-race` before you report |

## The scope — what not to touch

| | |
|---|---|
| **Home, Path, Board, the filter** | not this slice |
| **The CLI's verbs** | no new command. The same writer changing what an empty list looks like is this slice |
| **`kitchen/`** | the repo's own plan. Do not rewrite it |
| **`docs/audit.md`** | the reviewer's |
| **`work/tui/plan.md`** | the orchestrator's |
| **`work/`** | the orchestrator's, except the Q9 state row and your receipt |

The goldens move only if the screen changed. Say which, and why. A golden that
moved for no stated reason is the slice doing more than it said.

## The report

```bash
mkdir -p work/inbox
cat > work/inbox/runner-tui-q9.md <<EOF
Runner: Q9 done. Commits: <n>. Q9: <the raw file, and that a missing key still loads>. Q10: <shape, create, dismiss, promote — the id under the cursor>. DomainVersion: <7>. The goldens: <unchanged / what moved, and why>. The gate: <make check, make check-race>. Notes: <anything the reviewer should know>.
EOF
git add work/inbox/runner-tui-q9.md && git commit -q -m "inbox: runner reports tui-q9"
herdr agent prompt "w1H:pJ" "Runner: Q9 done. Commits: <n>. Q9: <the raw file, and that a missing key still loads>. Q10: <shape, create, dismiss, promote — the id under the cursor>. DomainVersion: <7>. The goldens: <unchanged / what moved, and why>. The gate: <make check, make check-race>. Notes: <anything the reviewer should know>."
```

Write the receipt, commit it, and then nudge. If a file you did not write has
changed, stop and report it. Do not touch it. Stage by name. Then stop.
