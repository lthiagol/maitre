You are a **fresh runner**. You did not write Q9. The repository is
`/Users/thiago/Code/Pessoal/maitre`. Read `AGENTS.md`, then
[work/tui/plan.md](../tui/plan.md) § 5, then this brief. This brief wins where they
differ. The orchestrator's pane id is `w1H:pJ`; use it verbatim in the notify command.

This slice is **T5**. The Path tab becomes **Ready**. Q3 already ruled: the tab
shows both blocks, and a small task's row says *promote → sign off → run*.

## The question the tab answers

> A recipe is on the Ready tab iff an agent could claim it right now.

That is `core.Adoptable`. It is a view: the `// view` marker, a registration in
`core.DerivedViews`, and a row in `docs/03` § 8 whose name comes from the runtime.
The wire fails both ways.

**Ask the claim, do not restate it.** A recipe is adoptable iff `EventStart` would
be accepted. Ask that transition on a **copy** — `Transition` mutates. A gate added
to the claim is a gate this view sees. The meaning, so a reader can check it:

| In | Out |
|---|---|
| approved, dependencies closed, not blocked, not cancelled | in |
| approved, an unmet dependency | out |
| blocked, cancelled, already executing, still being shaped | out |

## D-4 — one answer

`ReadinessOf` today answers *"can this be claimed?"* and **ignores `depends_on`**.
An approved recipe that waits on an open one reads as ready. `Path` knows the
dependency. Two answers.

**`ReadinessOf.Ready` is true iff the recipe is adoptable.** The reason names the
dependency when that is why. Do not leave a second predicate. Callers that wanted
"the lifecycle is approved" regardless of dependencies were the disagreement — fix
them.

`Path` stays. The Board and `maitre-cli path` still ask it. Do not retarget either
at `Adoptable`.

## The tab

The label is **Ready**. The count stays **7**:

`Home · Recipes · Backlog · Ready · Board · Log · Settings`

`Path` is not a tab label. The view `path` and the CLI command `path` keep their
names. Update every line that calls the **tab** Path — `docs/05`, `AGENTS.md`, the
Home tab's own test.

| Block | The rows | The move |
|---|---|---|
| **ADOPTABLE** | `core.Adoptable`, and nothing else | none. `space` writes **nothing**. The footer does not offer *claim*. The console does not claim (ADR-0003) |
| **SMALL TASKS** | open backlog items of kind `feature` or `bug`, from `core.BacklogGroups` | *promote*. The row says **promote → sign off → run**. The footer asks the domain for the move |

An idea is not on this tab. A promoted or dismissed item is not on it. An empty
block is an empty list, present, not a missing section.

An adoptable row opens the recipe's detail, the same detail the Recipes tab shows.
`space` on a small task promotes, through the move the backlog already uses. The
cursor follows the item: a promotion does not leave the cursor naming the id that
left, and it does not switch tabs (Q10).

The Board does not become this list. A recipe the Board still shows and Ready does
not — a draft, or an approved recipe with an open dependency — is the proof the
overlap shrank.

## The acceptance criteria

| AC | What |
|--:|---|
| **AC-01** | `core.Adoptable` is the claim's answer, registered, with a row in `docs/03` § 8. The wire fails both ways. `ReadinessOf.Ready` agrees with it, including the unmet-dependency case |
| **AC-02** | the tab is Ready, the seven names are right, the two blocks are the two views, the small task's row names the two further steps, `space` claims nothing and promotes a small task |
| **AC-03** | the Board and `maitre-cli path` still answer `Path`. The goldens that moved are claimed, and you say why. `DomainVersion` is **7** |
| **AC-04** | `work/tui/README.md`'s T5 row names this slice's commits |

One commit per criterion. The state row's own commit is not one of the ids it names.

## The invariants

| | |
|---|---|
| **`core.DomainVersion`** | stays **7**. No format change |
| **The domain** | `internal/core` and `internal/cli` stay stdlib-only. No dependency is added |
| **The tab count** | stays **7** |
| **The console never claims a recipe** | ADR-0003 |
| **Q8** | the free edit stays unbuilt |
| **Q11, Q12, Q13** | still open. Do not answer them |
| **The gate** | `make check`, then `make check-race` before you report |

## The scope — what not to touch

| | |
|---|---|
| **The Board's drawing** | T6. This slice does not redraw the tree |
| **The filter** | T7 |
| **The CLI's verbs** | no new command, no rename of `path` |
| **`kitchen/`** | do not rewrite the repo's plan |
| **`docs/audit.md`** | the reviewer's |
| **`work/tui/plan.md`** | the orchestrator's |
| **`work/`** | the orchestrator's, except the T5 state row and your receipt |

## The report

```bash
mkdir -p work/inbox
cat > work/inbox/runner-tui-t5.md <<EOF
Runner: T5 done. Commits: <n>. Adoptable: <what is in, what is out, and that it asks the claim>. Readiness: <that it agrees, including the unmet dependency>. The tab: <the seven names, the two blocks, the row's words>. The moves: <space claims nothing, space promotes a small task, the cursor does not follow the promotion off the tab>. The Board: <a recipe it still shows that Ready does not>. The goldens: <what moved, and why>. DomainVersion: <7>. The gate: <make check, make check-race>. Notes: <anything the reviewer should know>.
EOF
git add work/inbox/runner-tui-t5.md && git commit -q -m "inbox: runner reports tui-t5"
herdr agent prompt "w1H:pJ" "Runner: T5 done. Commits: <n>. Adoptable: <what is in, what is out, and that it asks the claim>. Readiness: <that it agrees, including the unmet dependency>. The tab: <the seven names, the two blocks, the row's words>. The moves: <space claims nothing, space promotes a small task, the cursor does not follow the promotion off the tab>. The Board: <a recipe it still shows that Ready does not>. The goldens: <what moved, and why>. DomainVersion: <7>. The gate: <make check, make check-race>. Notes: <anything the reviewer should know>."
```

Write the receipt, commit it, and then nudge. If a file you did not write has
changed, stop and report it. Stage by name. Then stop.
