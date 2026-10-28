You are the **reviewer** — the independent reviewer lane, fresh for slice F. The
orchestrator's pane id is `w1Z:p1` — it is already substituted in this file; wherever
this brief says `$ORCH_PANE`, it means `w1Z:p1`. Use it verbatim in the notify command
below. The repository is `/Users/thiago/Code/Pessoal/maitre`.

## Your slice — the review of slice F

The chef ruled that **this project is not brownfield**: no migrations, no
compatibility with previous versions, no upgrade paths. The retired formats are
**gone**. Slice F took the tolerance machinery out.

```bash
cd /Users/thiago/Code/Pessoal/maitre
git log --oneline -12
git diff <base>..HEAD --stat      # the base is the commit before slice F's first commit
```

**Cycle 1: blockers only.** File **HIGH or MEDIUM** blockers. No nits, no cosmetics,
no style, no LOW notes filed as rows. **Do not re-file F-56, F-59, F-60, F-61, P-8**
(the chef's).

### The line that decides this slice

**The journal is history; the plan's live records are current-format only.** The
journal's events keep their old names and ids forever — it is an append-only log —
and the reader still resolves what the journal names. The **live** records (recipes,
backlog items) are read in the current format and nothing else.

**Check every change against that line.** A change that makes the *journal* unreadable
is a blocker; a change that makes a *live* record unreadable in a retired format is
the slice's whole point.

### 1. Verify each acceptance criterion with a command or a mutation

**AC-01 — the reader stops reading the retired formats.**

- `internal/core/legacy.go`'s reading of the retired directories, the retired field
  names and the retired ids is **gone**, and the console's third group with it.
- **A plan carrying a retired artifact is refused, plainly** — not read as empty, and
  not silently ignored. Find the message and judge it: a chef must know in one line
  that the plan is not supported.
- `Init` does not create the retired directories (it did not before either).

**AC-02 — the deny-list keeps the retired names; the resolution goes.**

- A retired name — `line-cook`, `taster`, `expo`, `coordinator` — must still be
  **refused** as an approver. **Prove it:** try to approve a recipe as `line-cook`
  and confirm the gate refuses. (Today the refusal came *by* resolving the name to a
  role; without the resolution the refusal has to come from somewhere else.)
- The **resolution** is gone: a retired name is no longer "the runner".
- The tripwire follows: `docs/03` names every role the code refuses **and** every
  retired name. Check the wire says what it should now that the two sets are joined.

**AC-03 — the digest fences on content, and a rename never moves it.**

- The digest is taken over a **canonical payload with fixed field names** — the
  code's own shape, not the struct's JSON tags. **The proof that matters: rename a
  field and the digest must not move.** Try it.
- The payload covers **exactly** the frozen sections: intent, scope, each criterion's
  declared shape. A criterion's **evidence** and its **deferral** are excluded — the
  records grow after the approval.
- **The content fence still holds:** a recipe whose content changed after its
  sign-off is still refused. **Prove it by mutation.**
- There is **one** encoding and no candidate set.

**AC-04 — the version fence is one rule, with no upgrade path.**

- A tool writes a plan **only if** the plan's version **is** the tool's version.
  **Both directions** are refused, and each refusal says something a chef can act on.
- **Prove it:** a plan stamped one below the tool's version is refused, and one
  stamped one above is refused. A plan stamped the tool's version is written.
- The **fingerprint's job is unchanged** and is not the same job: it is the code-side
  tripwire that says *the format's facts changed*. The version is the plan-side fact.
  The prose in `docs/03` and `internal/core/lifecycle.go` must say exactly that.

**AC-05 — the tree moved to the current format.**

- **The eight deleted tests.** For **each** one: what property went with it, and **is
  that property still held by a test in the current format?** The chef said the
  retired-format tests may go; the chef did **not** say a live property may go with
  them. **Three to look at hard:**
  - `TestASelfReviewSurvivesTheRename` — the **self-review gate** must still refuse a
    recipe that its own executor reviewed. If the only test of that gate was the
    retired-format one, the property is now untested.
  - `TestDomainVersionStampsAnOlderPlan` — the fence's new rule (**a plan's version
    must be the tool's version**) must be tested, not just implemented.
  - `TestALegacyApprovalIsRestampedNotRefused` and
    `TestANormalisedLegacyApprovalIsRestamped` — the digest's *content* fence must
    still be tested.
- **The journal's history is untouched** — 37 events, unchanged, and the 38th is the
  new one. Check with `git diff` on `kitchen/activity.json`.
- **The journal's old ids still resolve.** That is the history, not a migration —
  the test survives, and it is not in the deleted list.

**AC-06 — the tripwires still hold, and one is new.**

```bash
go test ./internal/core ./internal/cli ./internal/ui -run 'Doc' -v
```

The derived-views table, the gate table, the plan-format fields, the command count
(**42**) and the roles the code refuses all agree with the code.

**The new wire** — `TestTheRetiredNamesAreOnlyWhereTheyBelong` — holds *"one format"*:
the retired names may live in exactly two places, the **detection that refuses them**
and the **deny-list that refuses them as names**. **Check the wire can fail** (a
mutation), and check its token list is complete: it should name every retired token
the code could grow back. A wire that scans for tokens it forgot to list is a wire
with a hole.

### 2. The judgement to check — "a lookup is tolerant but a record is not"

The runner's line: the **id normaliser strips any leading letters**, so `RCP01`,
`RCP1`, `1` and — as a side effect — `R01` all resolve to the same recipe, while the
**record on disk** is in the current format.

**The orchestrator's ruling:** the **lookup** is ergonomic — the number is the
identity, the directory carries the kind, and the letters are a kind hint — while the
**record** is one format. So `R01` resolving is not "the retired id is supported"; it
is the number answering to a spelling.

**But the sweep's own list names `R01`-style ids as old names that must be gone, and
the ruling says *one format*.** Judge it: is a retired id form resolving to a live
record a **second answer** to "what ids exist", or is it the same answer written
differently? Say which, and say whether the new wire's token list should name it.

### 3. Scope

Read the whole diff of the six commits. Did a change touch something the ruling did
not name? The **gate table**, the **plan format**, the **command count (42)**,
`core.DomainVersion` (**6**) and the **role names** must not have moved.

**Two things are not the runner's, and it left them:**

| | |
|---|---|
| `docs/07-vocabulary.md` | the decision of record — the orchestrator's. It has been corrected in the same spirit: the third group and the one-sided guard are gone from it. **Check it agrees with the code.** |
| `docs/audit.md` | **yours**. It is the ledger of the night and is never swept. |

**And the mascot** — another session is drawing it in this tree (`logo/`). Check the
runner committed nothing of it.

### Verdict

- No HIGH or MEDIUM findings → `verdict: ok`.
- Any HIGH or MEDIUM → `verdict: changes-needed`.
- LOW findings: **note them in your message, do not file them.**

### The ledger

`docs/audit.md` is **yours** to write. The runner did not touch it — confirm that with
`git diff`. Add any NEW blocker as a new row with an id continuing from **F-86**. Do
not re-file what is closed; the ledger is not swept. **Add a row for slice F itself.**
Commit the ledger.

### Hard rules

- Do **not** edit implementation code. Read-only for everything except
  `docs/audit.md`.
- Revert every mutation you make. The tree must be exactly as the runner left it,
  except your ledger commit.
- No `git commit --no-verify`. No push. Do not touch other panes. Do not spawn agents.
- **Do not touch `logo/`** — it is another session's work.

### Notify the orchestrator

```bash
mkdir -p work/inbox
cat > work/inbox/reviewer-slice-f.md <<EOF
Reviewer: slice F review done. Verdict: <ok|changes-needed>. Closed: <rows>. Filed <N> new blockers (<F-87.. or none>). Mutation checks: <what you mutated and reverted>. Deleted tests: <which live properties are still held, and by what>. LOW notes: <any>.
EOF
git add work/inbox && git commit -q -m "inbox: reviewer reports slice-f"
herdr agent prompt "w1Z:p1" "Reviewer: slice F review done. Verdict: <ok|changes-needed>. Closed: <rows>. Filed <N> new blockers (<F-87.. or none>). Mutation checks: <what you mutated and reverted>. Deleted tests: <which live properties are still held, and by what>. LOW notes: <any>."
```
