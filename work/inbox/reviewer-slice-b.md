Reviewer: slice B review done. Verdict: changes-needed. Closed: none. Filed 2 new blockers (F-72, F-73 — F-72 is the only MEDIUM; F-73 is the low `kitchen/project.json` drift).

## What held

- **The sharp edge holds.** `core.agentRole()` resolves a name as written *and* as a role,
  and `core.sameAgent()` compares the executor with the reviewer through it. A plan
  written before the rename cannot close on its own executor's review — proved end to
  end on a hand-edited plan (`claimed_by: line-cook`): the executor reviewing under
  `line-cook` is refused, the same executor reviewing under the **new** name `runner`
  is refused, and an independent `reviewer` is accepted.
- **The deny-list got tighter, not looser.** Every current and every retired role, in
  every spelling (`Line Cook`, `LINE-COOK`, ` line-cook `, `Runner`), is refused the
  approval; `chef` and an arbitrary human name are not. An unknown role still fails to
  match the list, exactly as before — no silent enum fallback.
- **The tabs.** `TabBoard` / `TabSousChef` throughout; `Chopboard` and `Expo` survive
  nowhere outside the three record exclusions. The Sous-Chef tab still says the
  automation is not built yet, still shows the vital signs, and grew no crew view. The
  header is one line; the goldens moved by the tab bar and nothing else.
- **One word per fact.** The row shows the plain state word, the lane is a section
  header, and the detail pane's meta line no longer carries the lane's label beside
  the state. The blocked reason kept its own line, so nothing was lost.
- **Invariants.** `internal/core` and `internal/cli` are stdlib-only; `internal/ui` is
  at 28 external modules; the domain layer imports no UI; `core.DomainVersion` is 4;
  `kitchen/activity.json` is untouched — 33 events, the 32 old ones still `R01`, the
  migration the 33rd, a `control` event stamped 4.
- **Scope.** The roles and the tabs only. The backlog, the ids, the module path, the
  binaries, `Proof`→evidence, the glossary and the mascot are untouched. The WIP diff
  (38 files) is entirely the rename — no other change rode along.

## The blocker

**F-72 (medium)** — `docs/03-domain-model.md` § "The CPD ↔ tool mapping" now contradicts
itself. The table keeps `runner · reviewer · orchestrator` in the **Kept plain** column
(this slice moved them), but the rule two lines above still says *"the **artifact and the
roles** get themed names"*. The table is the right half — the decision of record puts the
agents' jobs in the plain vocabulary — so the sentence is what is stale. Nothing checks
it: the new deny-list tripwire asks whether the doc *names* a role, not which column it
sits in.

## Mutation checks (all reverted; the tree is as the runner left it)

| Mutation | Result |
|---|---|
| `core.agentRole` returns `""` | `TestASelfReviewSurvivesTheRename`, `TestTheRenamedRolesAreTheAgentRoles` and 13 more fail |
| drop the `line-cook` → `runner` mapping | `TestASelfReviewSurvivesTheRename`, `TestTheRenamedRolesAreTheAgentRoles` fail |
| `core.sameAgent` ignores the role | `TestASelfReviewSurvivesTheRename` fails |
| `docs/03` no longer names `runner` | the new deny-list tripwire fails on real drift |

## LOW notes (not filed)

- `kitchen/project.json` still says `"name": "new-project"`; the console's header
  renders it, so the repo's own console reads `maitre · new-project`. **Filed as F-73
  (low)** — it is plan data, it belongs to slice C's plan migration, and it is not
  slice B's to change.
- `docs/audit.md`'s closed rows T4 and T6 still name `internal/ui/chop_test.go` and
  `internal/ui/expo_test.go`, which this slice renamed. A closed row claims a test
  holds the fix, and those two point at files that no longer exist. Not swept here.
- `docs/04-operations.md` names the roles in prose as well as `docs/03` does; only
  `docs/03` is wired to the code.
- The CLI's documented example (`--actor runner`) has no tripwire; a future rename
  would silently put a retired name back in the help text.
- The detail pane's readiness reason restates the state in kitchen words ("ready" and
  "awaiting sign-off"). The lane's *label* is gone from the pane, so the invariant as
  stated holds; the sentence remains.
- `GLOSSARY.md` still maps `runner`→`line cook`, `reviewer`→`taster`,
  `orchestrator`→`expo` — slice E's rewrite, and it is the file the whole rename
  points at.
