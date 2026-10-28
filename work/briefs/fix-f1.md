You are the **runner** — the same lane that wrote slice F. The orchestrator's pane id is
`w1Z:p1`; use it verbatim in the notify command below.

## The fix cycle for slice F — four findings

The reviewer filed **F-87**, **F-88**, **F-89** and **F-90**. The orchestrator ruled on
all four: **accepted**. Four findings, four commits — one each.

**Two of them are the kind that matter, and one of them is the orchestrator's own
briefing that caused the hole.** Read F-87 first.

### F-87 (high) — a retired name can *buy* a gate

> `line-cook` claims, `runner` reviews — and it passes `review.reviewer_independent`.

**The cause is the orchestrator's brief.** It said: *"the retired names must never be
able to sign a recipe off"* — the **approval** path. It did not say the retired names
must be refused **everywhere an actor names itself**, and the **independence
comparison** is where the hole is.

With the resolution gone, `line-cook` and `runner` are two different names, so the
gate reads the reviewer as independent of the executor. The executor claims as
`line-cook`, the runner reviews as `runner`, and the recipe closes with nobody
independent having looked at it.

**The fix — a retired name is not a name any more.**

The retired names are refused **wherever an actor names itself**: claim, review,
approve, and every other write. Not as a *special case for the approval*, but as
**the rule**: `line-cook` is a name the tool does not accept, and a name the tool
does not accept cannot appear in a live record either.

**This closes the hole without restoring the resolver**, which the ruling removed.
There is no `line-cook` → `runner` mapping any more: there is simply no `line-cook`.

**The repository's own plan moves with it.** `kitchen/recipes/RCP01-*.json` carries
`claimed_by: line-cook` and three rounds `by: line-cook` — a **live record** carrying
a retired name, which is exactly what the ruling forbids. Move it to the current
name, **consistently**: the claim and the rounds were the same agent, so they must
carry the same name, or the self-review the gate refuses would stop being refused.
**The journal's history keeps its old names** — it is history.

**Prove it:** claim as `line-cook`, review as `runner`, and confirm the gate refuses.

### F-88 (medium) — the operations doc still describes the resolution

`docs/04-operations.md` still says the retired names **name the role they stood
for**. The resolution is gone: a retired name is not "the runner", it is a name the
tool refuses.

**The fix:** the doc says what the code says.

### F-89 (medium) — the digest's encoding changed, and neither fence moved

> a plan from the previous build passes the version guard and is refused as
> tampered, and the repo's plan needed a hand-migration.

**The same class of bug as F-78, recreated by the fix for it.** The canonical payload
is a new encoding, so a stored digest from the previous build is uncomparable — and
the version guard passed it, because the plan format itself did not change. So the
digest said **tampered**, which is a lie: nothing was tampered with, the *encoding*
moved.

**The rule:** the tamper message must mean **content** and only content. An
encoding change must move the **version**, so the version guard gives the honest
message — *"this plan was written under version 6 and this tool speaks 7"* — instead
of the digest calling a plan tampered.

**The fix:**

1. **The digest's payload shape is a format fact.** The facts list must cover it, so
   that a change to the payload moves the shape, and the shape moving **must** move
   the version. This is the same lesson as the root file names and the journal
   envelope: *a fact that is not in the list is a fact the version stamp cannot see.*
2. **The version moves** (6 → 7), and the repo's own plan is stamped with it — the
   journal records it as an **append**, and the history is untouched.
3. **The repo's own plan's stored digest is re-stamped** to the new encoding, and
   the journal says so.
4. **The prose follows**: `docs/03` and the `DomainVersion` comment must say that a
   plan from another version is refused **by the version**, and that the digest's
   tamper message means content.

**Prove it:** a plan stamped **6** must be refused by the **version** guard, with the
honest message — not as tampering. And a plan stamped **7** must be written.

### F-90 (medium) — the "one format" wire has holes

> the wire scans `internal/core` only, exempts `model.go` wholesale, and omits the
> retired id literals (`R01`/`S01`/`N01`).

A wire that scans one package, exempts a whole file, and lists tokens it forgot is a
wire with three doors open. The reviewer's mutations: `'specials'` in `internal/core`
fails the wire; in `internal/ui`, in `model.go`, or as `R01`/`LegacyThing` — it
**passes**.

**The fix:**

1. **The wire scans the whole tree**, not one package.
2. **The exemption is narrowed** to where the retired names legitimately live — the
   **detection that refuses them** and the **deny-list that refuses them as names**.
   Not a whole file: the exemption must be **per place**, so a retired name cannot
   hide in `model.go` behind the deny-list's exemption.
3. **The token list is complete**: the retired *id* literals (`R01`, `S01`, `N01`)
   join it, and so does every other retired token the code could grow back.
4. **The token choice must be unambiguous.** The bare word `evidence` cannot be a
   token — it is a **current** word (the record's name) and the retired stamp field
   was called that too. Pick a token that names the retired thing and nothing else.
5. **The wire can fail** — prove it, on a token in each of the packages it now scans.

## The LOW notes — the rulings

| Note | Ruling |
|---|---|
| `activity.go`'s version-zero comment is stale | **fix it** — a comment that no longer matches the code is a bug. |
| **the file name's id form is not checked** | **fix it, and fold it into F-90.** A file named `R01-….json` in the current format is a retired id form on disk, and the lookup normalises it into the current record. **One format means the file's name is in the current format too.** Decide what the tool does with a retired form in a file name — refuse it, or say what it is — and make it a rule with a test. |
| **a plan with a retired artifact can still be written to** | **fix it.** The plan is not supported: *not supported* must mean no writes, not only "the read says so". The detection must gate the write path. |
| `GLOSSARY.md:97` names `coordinator` in prose | **fix it** — the glossary's rule is that a retired word is named **once**, in the `_Avoid_` line. |
| `docs/audit.md` was untouched by the runner | correct — it is the reviewer's. |
| `logo/` untouched | correct — another session's. |

## The tree is shared

**Another session is drawing the mascot in this repository** and has files in the tree
(`logo/`). Stage **by name** — never `git add -A`, never `git add .`. Leave anything
you did not write alone. If a file you did not write has changed, **stop and report
it**.

## Not yours — do not touch

| | |
|---|---|
| `docs/audit.md` | the reviewer's to write. It files F-87..F-90 and closes them. |
| `docs/07-vocabulary.md` | the decision of record — the orchestrator's to write. |
| `work/` | the orchestrator's. |
| The gate table, the plan format, the command count (42) | a fix may touch the **version** and the **facts list** (F-89). Nothing else. If a fix would touch anything else, **stop and report it**. |

## The gate

```bash
make lint && gofmt -l . && make test && make regen-goldens && git status --short
```

Green.

## Four commits, one per finding

```
Per-AC evidence:
- F-87: a retired name cannot claim, review or approve; go test ./internal/core -run 'Retired|Independent' -v exit 0
```

## The report — the receipt comes first

```bash
mkdir -p work/inbox
cat > work/inbox/runner-slice-f-fix-1.md <<EOF
Runner: slice F fix 1 done. Commits: <n>. Fixed: F-87..F-90 and the notes. Notes: <anything the reviewer should know>.
EOF
git add work/inbox && git commit -q -m "inbox: runner reports slice-f fix 1"
herdr agent prompt "w1Z:p1" "Runner: slice F fix 1 done. Commits: <n>. Fixed: F-87..F-90 and the notes. Notes: <anything the reviewer should know>."
```

**One write per file, and commit as soon as the file is coherent.**
