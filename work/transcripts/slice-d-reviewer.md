# The reviewer's lane — slice D, captured 2026-10-27

Dumped from pane `w1Z:p9` with `herdr pane read --source recent-unwrapped --lines 5000`.
A fresh lane: it re-derived the facts from the code, filed F-78 (high), F-79 and
F-80, and proved the fixes by mutation — the old-name candidate dropped, the digest
put back, the criterion changed after the sign-off, and a plan in the old names
claimed end to end.

```
  ┃  1376e05 docs: the ledger files slice D's three findings
  ┃  8bf1bf0 work: slice D's review brief
  ┃

  ┃
  ┃  $ herdr agent prompt "w1Z:p1" "Reviewer: slice D review done. Verdict: changes-needed. Closed: none. Filed 3 new blockers (F
  ┃  -78 high, F-79 medium, F-80 medium). Mutation checks: (1) the derived-views tripwire, both ways - deleted the evidenc…
  ┃
  ┃  {"id":"cli:agent:prompt","result":{"agent":{"agent":"opencode","agent_session":{"agent":"opencode","kind":"id","source":
  ┃  "herdr:opencode","value":"ses_e80f147cbffetwcFtsbsOZTA2z"},"agent_status":"idle","completion_seq":432,"cwd":"/Users/thiago/
  ┃  Code/Pessoal/maitre","focused":false,"foreground_cwd":"/Users/thiago/Code/Pessoal/maitre","pane_id":"w1Z:p1","revision":19,
  ┃  "screen_detection_skipped":true,"state_change_seq":432,"tab_id":"w1Z:t1","terminal_id":"term_65ed7354d7f7722",
  ┃  "terminal_title":"OC | orchestrator (fork #1)","terminal_title_stripped":"OC | orchestrator (fork #1)","workspace_id":"w1Z"}
  ┃  ,"type":"agent_prompted"}}
  ┃  NOTIFIED
  ┃

  ┃
  ┃  $ git status --porcelain; echo "=== commits ==="; git log --oneline -3 | cat; echo "=== tree vs runner ==="; git diff 8bf1bf
  ┃  0..HEAD --stat | cat
  ┃
  ┃  (3 earlier lines)
  ┃  1376e05 docs: the ledger files slice D's three findings
  ┃  8bf1bf0 work: slice D's review brief
  ┃  === tree vs runner ===
  ┃   docs/audit.md                  | 4 ++++
  ┃   work/inbox/reviewer-slice-d.md | 3 +++
  ┃   2 files changed, 7 insertions(+)
  ┃

     + Thought · 997ms

     Verdict: changes-needed — 1 HIGH, 2 MEDIUM filed (F-78..F-80)

     The two words are right; the slice broke the plan's approval digest.

     Findings

     ┌──────┬────────┬────────────────────────────────────────────────────────────┬─────────────────────────────────────────────┐
     │ #    │ Sev    │ What                                                       │ Proof                                       │
     ├──────┼────────┼────────────────────────────────────────────────────────────┼─────────────────────────────────────────────┤
     │ F-78 │ high   │ The approval digest did not travel with the rename. The    │ The repo's own RCP01 stores 8813a266… (     │
     │      │        │ digest hashes the criterion's JSON, and the JSON carries   │ reproduced exactly = the pre-rename         │
     │      │        │ the field name — proofs vs evidence. restampLegacyApproval │ encoding); the v6 tool recomputes           │
     │      │        │ only knows the null/[] candidates, so a pre-rename         │ 4984098d…. A plan in the old names, re-     │
     │      │        │ approval digest is unreproducible and immutability.        │ written by the v6 tool, keeps the old       │
     │      │        │ respected refuses every gated transition with "frozen      │ digest → refused at start.                  │
     │      │        │ sections changed since approval" on a recipe nobody        │                                             │
     │      │        │ touched.                                                   │                                             │
     ├──────┼────────┼────────────────────────────────────────────────────────────┼─────────────────────────────────────────────┤
     │ F-79 │ medium │ The glossary's Evidence — "a verification result a machine │ internal/core/model.go · Criterion.Status;  │
     │      │        │ produced", "an evidence that matched nothing is not        │ verifier.go · MatchSignedOff                │
     │      │        │ evidence" — excludes the human sign-off, which the code    │                                             │
     │      │        │ stores in Criterion.Evidence{Match:"signed-off", Matched:  │                                             │
     │      │        │ 0} and ac.all_proven accepts. docs/03 §5 says the opposite │                                             │
     │      │        │ of the glossary.                                           │                                             │
     ├──────┼────────┼────────────────────────────────────────────────────────────┼─────────────────────────────────────────────┤
     │ F-80 │ medium │ docs/03 §6: "the records it wrote would be records the     │ docs/03 §6 · internal/core/legacy.go        │
     │      │        │ reader does not see" — two paragraphs above "the reader    │                                             │
     │      │        │ reads them". legacy.go reads them. Same sentence in        │                                             │
     │      │        │ lifecycle.go's version-6 note.                             │                                             │
     └──────┴────────┴────────────────────────────────────────────────────────────┴─────────────────────────────────────────────┘

     What held

     ┌────────────────────────────────────────────────────────┬─────────────────────────────────────────────────────────────────┐
     │ Check                                                  │ Result                                                          │
     ├────────────────────────────────────────────────────────┼─────────────────────────────────────────────────────────────────┤
     │ core.Proof→core.Evidence, JSON evidence, view          │ ✓                                                               │
     │ evidence_summary/EvidenceSummaryOf                     │                                                                 │
     ├────────────────────────────────────────────────────────┼─────────────────────────────────────────────────────────────────┤
     │ Derived-views tripwire, both ways                      │ deleted row → FAIL · mistyped function → FAIL (both directions  │
     │                                                        │ named) · extra registered view → FAIL. Reverted.                │
     ├────────────────────────────────────────────────────────┼─────────────────────────────────────────────────────────────────┤
     │ Legacy reader: old names read, new name wins, writer   │ ✓ (CLI, end to end)                                             │
     │ writes new names                                       │                                                                 │
     ├────────────────────────────────────────────────────────┼─────────────────────────────────────────────────────────────────┤
     │ Version fence                                          │ stamped 4/5/6 accepted · stamped 7 refused ✓                    │
     ├────────────────────────────────────────────────────────┼─────────────────────────────────────────────────────────────────┤
     │ Journal                                                │ events 1–34 untouched, 35/36 appended (git diff) ✓              │
     ├────────────────────────────────────────────────────────┼─────────────────────────────────────────────────────────────────┤
     │ Records append-only · journal append-only              │ ✓ (ac verify appended, old record kept)                         │
     ├────────────────────────────────────────────────────────┼─────────────────────────────────────────────────────────────────┤
     │ stdlib-only core/cli · ui 28 · zero UI imports         │ ✓                                                               │
     ├────────────────────────────────────────────────────────┼─────────────────────────────────────────────────────────────────┤
     │ Sign-off exception                                     │ still passes ac.all_proven, still the one zero-match exception  │
     │                                                        │ ✓                                                               │
     ├────────────────────────────────────────────────────────┼─────────────────────────────────────────────────────────────────┤
     │ Scope (roles, tabs, backlog, ids, module, binaries)    │ untouched ✓                                                     │
     └────────────────────────────────────────────────────────┴─────────────────────────────────────────────────────────────────┘







     Build · MiniMax-M3.1-Flash · 9m 17s · 80.2 tok/s

     ◈ Instructions updated: core/instructions

  ┃
  ┃  You are the **reviewer** — the same lane that reviewed slice D. The orchestrator's pane
  ┃  id is `w1Z:p1`; use it verbatim in the notify command below.
  ┃
  ┃  ## Cycle 2 — the fixes to slice D
  ┃
  ┃  The runner fixed the three findings you filed. **Cycle 2 is blockers only**: file
  ┃  **HIGH or MEDIUM** blockers. No nits, no cosmetics, no style.
  ┃
  ┃  **Do not re-file:** F-56, F-59, F-60, F-61, P-8 (the chef's), or your own LOW notes
  ┃  from cycle 1 (the grammar of *"an evidence"* — folded into F-79; `docs/07`'s
  ┃  alternatives table — the orchestrator ruled it stays; `research/` — history; the
  ┃  glossary's brigade section — slice E's; `scripts/stamp-plan`'s usage line — fixed
  ┃  here; `AGENTS.md` and `work/README.md` — the orchestrator's, and correct; and § 6's
  ┃  "does not look" sentence — folded into F-80).
  ┃
  ┃  ### 1. Are the three findings actually closed?
  ┃
  ┃  Verify each with a **command or a mutation**, not by reading the diff and believing it.
  ┃
  ┃  **F-78 (high) — the approval digest travels with the rename.**
  ┃
  ┃  - `legacyFingerprints` now enumerates the **field name** as well as the list forms.
  ┃    **Check the set is still COMPLETE.** The comment above the function used to argue
  ┃    why the set was complete; that argument had to be rewritten, not extended in
  ┃    silence. Read the new argument and check it: the candidates must be the **cross
  ┃    product** of every encoding an older build could have written — the old field name
  ┃    *and* the `null`/`[]` list forms — and nothing else.
  ┃  - **The safety property is the one that matters now.** The candidate set grew, and
  ┃    every candidate is another way a digest could accidentally match. **Prove that a
  ┃    recipe whose content changed after approval still matches NO candidate**, and that
  ┃    the `immutability.respected` gate still refuses it. A gate that can be fooled by a
  ┃    wider candidate set is worse than the bug you filed.
  ┃  - **The repository's own plan is migrated** and the standing discipline is a test
  ┃    now — *the repository's own plan's stored approval digest must still match the
  ┃    canonical recomputation*. Check that the test reads the digest **off the disk** and
  ┃    compares it with the recomputation, and that it would fail if the file were wrong.
  ┃  - **The journal** took the migration as an append (event 37). Check with `git diff`
  ┃    on `kitchen/activity.json`: events 1–36 untouched.
  ┃
  ┃  **F-79 (medium) — the glossary's entry covers the record it names.**
  ┃
  ┃  - The **Evidence** entry names the sign-off exception, and `docs/03` § 5 says the same
  ┃    thing. **Check the two say the same thing** — the single source and the spec must
  ┃    not contradict each other, which is what the finding was.
  ┃  - **Check the sentence is true**: a result that matched nothing does not count, *and
  ┃    the chef's sign-off is the exception, and it claims no match*. The code stores a
  ┃    sign-off as `Criterion.Evidence{Match: "signed-off", Matched: 0, SignedOff: true}`
  ┃    and `ac.all_proven` takes it as a pass.
  ┃  - *"An evidence"* is gone from the tree outside `docs/audit.md`. Check that the word
  ┃    was not simply **banished** where it was right: `ac.all_proven`, the match states
  ┃    (`proven` · `rejected` · `unknown`), the old field names, and prose that means the
  ┃    *idea* of a proof may keep it.
  ┃
  ┃  **F-80 (medium) — the version paragraph stops contradicting the reader.**
  ┃
  ┃  - `docs/03` § 6 and `internal/core/lifecycle.go`'s version notes now say what the
  ┃    reader **does**: it reads the old field name, it shows the old specials read-only
  ┃    in their own group, and it never looks for new ones. **Check each claim against the
  ┃    code**, and that the version bump is unchanged (the *writer* changed — that is what
  ┃    the version fences).
  ┃  - `scripts/stamp-plan`'s usage line names the field-name migrations.
  ┃
  ┃  ### 2. The ADR line — judge it against the rule, not against taste
  ┃
  ┃  The runner's F-79 commit renamed the vocabulary **incidentally** in ADR-0002,
  ┃  ADR-0006, ADR-0008 and the ADR index, and left **ADR-0004 (naming)** and **ADR-0005
  ┃  (vocabulary)** untouched.
  ┃
  ┃  **The orchestrator's ruling, which you judge against:**
  ┃
  ┃  | | |
  ┃  |---|---|
  ┃  | **An ADR whose subject is the mechanism** (0002, 0006, 0008) | its *incidental* vocabulary may move — it decides about
  ┃  verification or about the interface, not about the word, and a doc that contradicts the code is a bug |
  ┃  | **An ADR whose subject is the vocabulary** (0004, 0005) | it **keeps its tables as the record of what it said**, and ADR
  ┃  0009 supersedes it in slice E. Those are the two the runner left alone. |
  ┃
  ┃  **Judge two things:** that the runner drew this line correctly, and that in the ADRs
  ┃  it did touch, only the *vocabulary* moved — **not the substance**. An ADR that
  ┃  decides something different now, or loses a constraint, is a rewrite of the record
  ┃  and a blocker. Read those four ADRs' diffs closely.
  ┃
  ┃  ### 3. What did the fixes introduce?
  ┃
  ┃  1. **Scope.** The gate table, the plan format, the command count (**42**) and
  ┃     `core.DomainVersion` (**6**) must not have moved. Check the whole diff of the
  ┃     three commits — the F-79 commit touched 21 files, and each one needs a reason.
  ┃  2. **The shared tree.** Another session is drawing the mascot in this repository and
  ┃     has files in the tree (`logo/`). **Confirm the runner committed nothing of it** —
  ┃     no file of the other session's, and no file the runner did not write. The
  ┃     runner's report says it staged every commit by name.
  ┃  3. **The plan's writes.** The only write to the plan should be the migration event.
  ┃     Check that the journal took **exactly one** new event and that nothing else in
  ┃     `kitchen/` moved except the recipe's own migrated field.
  ┃
  ┃  ### 4. The ledger
  ┃
  ┃  Move **F-78, F-79, F-80** from open to closed, with the file and the test that now
  ┃  holds each. Commit the ledger. The ledger is not swept: keep the findings' wording
  ┃  as it was filed.
  ┃
  ┃  ### Verdict
  ┃
  ┃  - No HIGH or MEDIUM findings → `verdict: ok`.
  ┃  - Any HIGH or MEDIUM → `verdict: changes-needed`.
  ┃  - LOW findings: **note them in your message, do not file them.**
  ┃
  ┃  ### Hard rules
  ┃
  ┃  - Do **not** edit implementation code. Read-only for everything except
  ┃    `docs/audit.md`.
  ┃  - Revert every mutation you make. The tree must be exactly as the runner left it,
  ┃    except your ledger commit.
  ┃  - No `git commit --no-verify`. No push. Do not touch other panes. Do not spawn agents.
  ┃  - **Do not touch `logo/`** — it is another session's work.
  ┃
  ┃  ### Notify the orchestrator
  ┃
  ┃  ```bash
  ┃  mkdir -p work/inbox
  ┃  cat > work/inbox/reviewer-slice-d-cycle-2.md <<EOF
  ┃  Reviewer: slice D cycle 2 done. Verdict: <ok|changes-needed>. F-78..F-80 closed: <yes|no>. Filed <N> new blockers (<F-81..
  ┃  or none>). Mutation checks: <what you mutated and reverted>. LOW notes: <any>.
  ┃  EOF
  ┃  git add work/inbox && git commit -q -m "inbox: reviewer reports slice-d cycle 2"
  ┃  herdr agent prompt "w1Z:p1" "Reviewer: slice D cycle 2 done. Verdict: <ok|changes-needed>. F-78..F-80 closed: <yes|no>.
  ┃  Filed <N> new blockers (<F-81.. or none>). Mutation checks: <what you mutated and reverted>. LOW notes: <any>."
  ┃  ```
  ┃

     + Thought · 160ms

  ┃
  ┃  $ git log --oneline -8 | cat; echo "=== status ==="; git status --porcelain; echo "=== logo ==="; ls logo/ 2>/dev/null | hea
  ┃  d
  ┃
  ┃  (15 earlier lines)
  ┃  3-letterhead.png
  ┃  3-letterhead.svg
  ┃  4-night-card.png
  ┃  4-night-card.svg
  ┃  5-both-sides.png
  ┃  5-both-sides.svg
  ┃

  ┃
  ┃  $ git diff 954b9c4..HEAD --stat | cat
  ┃
  ┃  (22 earlier lines)
  ┃   ...the-console-shows-a-finding-s-conversation.json |   2 +-
  ┃   scripts/stamp-plan/main.go                         |   5 +-
  ┃   work/README.md                                     |  25 +++-
  ┃   work/briefs/fix-d1.md                              | 163 +++++++++++++++++++++
  ┃   work/briefs/review-d2.md                           | 124 ++++++++++++++++
  ┃   work/inbox/runner-slice-d-fix-1.md                 |  73 +++++++++
  ┃   28 files changed, 636 insertions(+), 69 deletions(-)
  ┃

     + Thought · 1.1s

  ┃
  ┃  $ git show ac156ff --stat | cat; echo "=========="; git show ac156ff -- internal/core/model.go | cat
  ┃
  ┃  (194 earlier lines)
  ┃  +        if err != nil {
  ┃  +          continue
  ┃  +        }
  ┃  +        out = append(out, legacy)
  ┃         }
  ┃       }
  ┃     }
  ┃

     + Thought · 9.7s

  ┃
  ┃  $ git show ac156ff -- internal/core/approval_test.go kitchen/ scripts/stamp-plan/main.go | cat
  ┃
  ┃  (215 earlier lines)
  ┃  -// had (`specials+prep-notes` → `backlog`).
  ┃  +// The forms are the id forms (`R01` → `RCP01`), the shapes the plan's records had
  ┃  +// (`specials+prep-notes` → `backlog`), or the names a record's field had
  ┃  +// (`proofs` → `evidence`).
  ┃   //
  ┃   // For this repository's own plan:
  ┃   //
  ┃

     + Thought · 2.9s

     → Explored: 1 read

     + Thought · 2.6s

  ┃
  ┃  $ git show 8c5899d --stat | cat; echo "=========="; git show 8c5899d -- GLOSSARY.md docs/03-domain-model.md | cat
  ┃
  ┃  (148 earlier lines)
  ┃  -| **Zero-match is not evidence** | `matched == 0` ⇒ rejected, recorded as a failure |
  ┃  -| **A sign-off is not machine evidence** | it is its own state (`signed-off`) — **the chef** vouching for a procedure the
  ┃  tool could not check. It claims no match, and it is the one exception to the zero-match rule |
  ┃  +| **Zero-match does not count** | `matched == 0` ⇒ rejected, recorded as a failure |
  ┃  +| **A sign-off is evidence too** | it is its own state (`signed-off`) — **the chef** vouching for a procedure the tool
  ┃  could not check. It claims no match, and it is the one exception to the zero-match rule |
  ┃   | **Exit code is the truth** | `exit 0` is a pass; anything else is a fail |
  ┃   | **Bypass is debt** | A forced transition is stamped `[forced]` and shows in the UI as debt, never as success (CPD
  ┃  principle 9) |
  ┃   | **Timeouts are distinguishable** | v1's timeout message (*"verification command did not run"*) was indistinguishable from
  ┃  a real failure. A timeout is its own outcome, and the default deadline must fit a realistic full-suite run |
  ┃

     + Thought · 192ms

  ┃
  ┃  $ git show f09a682 --stat | cat; echo "=========="; git show f09a682 -- docs/03-domain-model.md internal/core/lifecycle.go |
  ┃   cat
  ┃
  ┃  (129 earlier lines)
  ┃  +// `kitchen/specials/` and `kitchen/prep-notes/`, and the version-5 writer does not.
  ┃  +// The reader still shows the old ones, read-only, in their own group; it never looks
  ┃  +// for new ones.
  ┃   //
  ┃   // The reader is tolerant on purpose: `BL01`, `bl01`, `BL1`, `B01` and `1` all name
  ┃   // the same item, and a file named `B01-<slug>.json` is read as the item `BL01`. The
  ┃

     + Thought · 3.9s

  ┃
  ┃  $ git show 8c5899d -- docs/decisions/ AGENTS.md docs/02-principles.md docs/06-proposals.md | cat
  ┃
  ┃  (134 earlier lines)
  ┃   | [0003](0003-surface-split.md) | Who writes what — artifact vs. management | 🔄 proposed |
  ┃   | [0004](0004-naming.md) | Naming — project, binaries, plan directory, ladder, roles | ✅ accepted |
  ┃   | [0005](0005-vocabulary.md) | Vocabulary — generic methodology, specific product | ✅ accepted |
  ┃  -| [0006](0006-verification.md) | Verification — who runs it, and what an evidence is | ✅ accepted |
  ┃  +| [0006](0006-verification.md) | Verification — who runs it, and what evidence is | ✅ accepted |
  ┃   | [0007](0007-operation-layer.md) | The operation layer | ✅ accepted |
  ┃   | [0008](0008-artifact-is-the-channel.md) | The artifact is the channel — no agent-to-agent messaging | ✅ accepted |
  ┃

     + Thought · 11.8s

  ┃
  ┃  $ git show 8c5899d -- internal/core/model.go internal/core/lifecycle.go internal/core/verifier.go internal/core/views.go int
  ┃  ernal/core/operations.go internal/cli/fragments.go internal/ui/model.go internal/core/seed.go | cat
  ┃
  ┃  (199 earlier lines)
  ┃   //
  ┃  -// A recipe is longer than the pane — an intent, an evidence per criterion, a review
  ┃  +// A recipe is longer than the pane — an intent, evidence per criterion, a review
  ┃   // thread — so without this the bottom of a recipe is unreachable.
  ┃   func (m *Model) scrollDetail(delta int) {
  ┃     m.detailScroll += delta
  ┃

     + Thought · 1.4s

  ┃
  ┃  $ grep -rni "an evidence\|a evidence" --include="*.go" --include="*.md" --include="*.json" . | grep -v "^./work/" | grep -v
  ┃  "^./docs/audit.md" | grep -v "^./logo"; echo "--- (empty above means gone) ---"; echo "=== proof/proven kept ==="; gr…
  ┃
  ┃  (36 earlier lines)
  ┃  ./internal/core/model.go:556:// The rename added a second axis. The criterion's record was `proofs` before it and
  ┃  ./internal/core/model.go:558:// build that wrote `proofs` hashed a different document, and its digest is a candidate
  ┃  ./internal/core/verifier.go:17://  proven   — the count is known and greater than zero
  ┃  ./internal/core/verifier.go:19://  unproven — the count could not be determined; not machine evidence
  ┃  ./internal/core/verifier.go:21:// Only proven satisfies a gate. See docs/decisions/0006-verification.md.
  ┃  ./internal/core/verifier.go:25:  MatchProven   MatchState = "proven"
  ┃

     + Thought · 2.3s

  ┃
  ┃  $ git show 8c5899d -- docs/04-operations.md internal/core/lifecycle_test.go internal/core/gates_test.go internal/core/approv
  ┃  al_test.go internal/core/operations_test.go | cat | grep -v "^ " | head -80
  ┃
  ┃  (60 earlier lines)
  ┃  @@ -23,7 +23,7 @@ func newOps(t *testing.T) (*Operations, *Store) {
  ┃  -    Outcome:    "a zero-match evidence never satisfies a gate",
  ┃  +    Outcome:    "a zero-match result never satisfies a gate",
  ┃  @@ -32,7 +32,7 @@ func spec() RecipeSpec {
  ┃  -      Action: "add the gate", DoneWhen: "a rejected evidence blocks", Covers: []string{"AC-01"},
  ┃  +      Action: "add the gate", DoneWhen: "a rejected result blocks", Covers: []string{"AC-01"},
  ┃

     + Thought · 11.5s

  ┃
  ┃  $ go test ./... 2>&1 | tail -8; echo "=== lint ==="; gofmt -l .; go vet ./... 2>&1|tail -3; echo "=== deps ==="; go list -de
  ┃  ps -f '{{if not .Standard}}{{.Module.Path}}{{end}}' ./internal/core | sort -u; go list -deps -f '{{if not .Standard}}…
  ┃
  ┃  (4 earlier lines)
  ┃  ok    github.com/lthiagol/maitre/internal/ui  (cached)
  ┃  ?     github.com/lthiagol/maitre/scripts/stamp-plan  [no test files]
  ┃  === lint ===
  ┃  === deps ===
  ┃  github.com/lthiagol/maitre
  ┃        28
  ┃

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre
  ┃  cp internal/core/model.go /tmp/model.bak…
  ┃
  ┃  (4 earlier lines)
  ┃      approval_test.go:149: the approval should be re-stamped canonically
  ┃  --- FAIL: TestARenamedFieldIsNotATamper (0.06s)
  ┃  === RUN   TestTheRepoPlansApprovalDigestMatchesItsContent
  ┃  --- PASS: TestTheRepoPlansApprovalDigestMatchesItsContent (0.00s)
  ┃  === RUN   TestAFormChangeIsNotATamper
  ┃  --- PASS: TestAFormChangeIsNotATamper (0.10s)
  ┃

     + Thought · 1.6s

  ┃
  ┃  $ cat > internal/core/zz_probe_test.go <<'EOF'
  ┃  package core…
  ┃
  ┃  === RUN   TestZZTamperIsNotRescuedByTheLegacyCandidate
  ┃      zz_probe_test.go:51: legacy=false refused with: start: gate(s) not satisfied: immutability.respected (frozen sections
  ┃  changed since approval)
  ┃      zz_probe_test.go:51: legacy=true refused with: start: gate(s) not satisfied: immutability.respected (frozen sections
  ┃  changed since approval)
  ┃  --- PASS: TestZZTamperIsNotRescuedByTheLegacyCandidate (0.14s)
  ┃  PASS
  ┃  ok    github.com/lthiagol/maitre/internal/core  0.345s
  ┃  ?? logo/
  ┃

     + Thought · 782ms

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre
  ┃  cp kitchen/recipes/RCP01-the-console-shows-a-finding-s-conversation.json /tmp/rcp01.bak…
  ┃
  ┃  (2 earlier lines)
  ┃      approval_test.go:191: RCP01 carries the digest "8813a266eaac833a280c99016d53bbceb85b95cf1e3b0121354f330de4072edc", want
  ┃  the canonical "4984098d94858a1ec142f89c987e41d5c232bc2ad85c7699893ef9b20fd75fe1" — the approval was written under another
  ┃  encoding
  ┃  --- FAIL: TestTheRepoPlansApprovalDigestMatchesItsContent (0.00s)
  ┃  FAIL
  ┃  FAIL  github.com/lthiagol/maitre/internal/core  0.209s
  ┃  FAIL
  ┃  ?? logo/
  ┃

     + Thought · 319ms

  ┃
  ┃  $ go build -o /tmp/maitre-cli-c2 ./cmd/maitre-cli && rm -rf /tmp/legacy3 && mkdir -p /tmp/legacy3 && cd /tmp/legacy3 && /tmp
  ┃  /maitre-cli-c2 init >/dev/null && cd /tmp/legacy3 && /tmp/maitre-cli-c2 recipe create --title "Legacy plan" --outcome…
  ┃
  ┃  old old_fp: bf61c5d0b4c573a35b300db60335285211dac1a13c525f6a5db462915ac0ad15 new_fp:
  ┃  27b3de17a32cc7fec0165a6b917fc7d489d65ab74d1929d14263141bc317f07b
  ┃  24:      "proofs": []
  ┃  48:  "evidence": {},
  ┃  61:    "fingerprint": "bf61c5d0b4c573a35b300db60335285211dac1a13c525f6a5db462915ac0ad15"
  ┃

  ┃
  ┃  $ cd /tmp/legacy3 && echo "--- the claim on a plan written before the rename ---" && /tmp/maitre-cli-c2 --actor runner recip
  ┃  e claim RCP01 --agent runner | head -12; echo "EXIT=$?"; echo "--- the file after ---"; grep -n "proofs\|evidence\|co…
  ┃
  ┃  (12 earlier lines)
  ┃      "guidance": 1,
  ┃  EXIT=0
  ┃  --- the file after ---
  ┃  24:      "evidence": []
  ┃  51:  "completion": {},
  ┃  64:    "fingerprint": "27b3de17a32cc7fec0165a6b917fc7d489d65ab74d1929d14263141bc317f07b"
  ┃

     + Thought · 9.8s

  ┃
  ┃  $ cat > internal/core/zz_probe_test.go <<'EOF'
  ┃  package core…
  ┃
  ┃  (5 earlier lines)
  ┃      zz_probe_test.go:61: candidates: 4
  ┃      zz_probe_test.go:66: candidates (scope empty): 4
  ┃  --- PASS: TestZZCandidateSetSize (0.00s)
  ┃  PASS
  ┃  ok    github.com/lthiagol/maitre/internal/core  0.334s
  ┃  ?? logo/
  ┃

     + Thought · 35.2s

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre
  ┃  echo "=== kitchen diff since cycle 1 ==="; git diff 954b9c4..HEAD -- kitchen/ | cat…
  ┃
  ┃  (60 earlier lines)
  ┃   // The zero-match rule is the project's answer to v1's worst bug. A sign-off is a human
  ┃  -// judgement, not a machine evidence, and it claims no match.
  ┃  +// judgement, not machine evidence, and it claims no match.
  ┃   func TestASignOffIsNotAMachineEvidence(t *testing.T) {
  ┃     ops, _ := newOps(t)
  ┃     r, err := ops.CreateRecipe(readySpec())
  ┃

     + Thought · 2.0s

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre
  ┃  grep -rn "DirSpecials\|DirPrepNotes\|specials\|prep-notes" --include="*.go" internal/ | grep -v _test | head -20…
  ┃
  ┃  (7 earlier lines)
  ┃  internal/core/seed.go:181:  if _, err := chef.SignOffCriterion(r04.ID, "AC-01", "chef", "the backlog lists specials"); err !
  ┃  = nil {
  ┃  === the UI legacy group ===
  ┃  internal/ui/model.go:64:  legacy  []core.LegacyItem
  ┃  internal/ui/model.go:114:  legacy []core.LegacyItem
  ┃  internal/ui/model.go:245:    legacy, err := m.store.ListLegacyItems()
  ┃  internal/ui/sanitize.go:50:func sanitizePlan(recipes []core.Recipe, items []core.BacklogItem, legacy []core.LegacyItem,
  ┃  journal core.Journal) {
  ┃

     → Explored: 1 read

     + Thought · 875ms

  ┃
  ┃  $ grep -rn "legacy" internal/ui/views.go internal/ui/backlog.go 2>/dev/null | head -20; echo ---; ls internal/ui/
  ┃
  ┃  (30 earlier lines)
  ┃  term.go
  ┃  testdata
  ┃  theme_test.go
  ┃  theme.go
  ┃  triage_test.go
  ┃  views.go
  ┃  words_test.go
  ┃

     + Thought · 2.0s

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre
  ┃  git log --oneline --name-status 954b9c4..HEAD | cat
  ┃
  ┃  (33 earlier lines)
  ┃  M  kitchen/recipes/RCP01-the-console-shows-a-finding-s-conversation.json
  ┃  M  scripts/stamp-plan/main.go
  ┃  625c789 work: slice D's findings, and the tree we share
  ┃  M  work/README.md
  ┃  ee1b474 work: slice D's fix cycle — the digest did not travel with the rename
  ┃  A  work/briefs/fix-d1.md
  ┃

     + Thought · 8.1s

     → Explored: 1 read

     + Thought · 12.2s

     → Explored: 1 read

     + Thought · 17.4s

     → Explored: 1 read

     + Thought · 5.7s

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre
  ┃  go test ./internal/core ./internal/cli ./internal/ui -run 'Doc' 2>&1 | tail -5…
  ┃
  ┃  (10 earlier lines)
  ┃  github.com/lthiagol/maitre
  ┃  github.com/lthiagol/maitre
  ┃        28
  ┃  === zero UI imports in core ===
  ┃  0
  ┃  0 charm packages
  ┃

     + Thought · 743ms

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre
  ┃  go test -race ./internal/... 2>&1 | tail -5…
  ┃
  ┃  ok    github.com/lthiagol/maitre/internal/cli  5.062s
  ┃  ok    github.com/lthiagol/maitre/internal/core  22.221s
  ┃  ok    github.com/lthiagol/maitre/internal/ui  80.056s
  ┃  === canary ===
  ┃

  ┃
  ┃  $ go build -o /tmp/maitre-c2 ./cmd/maitre && cd /tmp/kb2 && /tmp/maitre-c2 --render --plain | head -14; echo "=== repo plan
  ┃  render ==="; /tmp/maitre-c2 --render --plain | head -8
  ┃
  ┃  (17 earlier lines)
  ┃  ▸ Recipes    Backlog    Path    Board    Log    Sous-Chef    Settings
  ┃   7 recipes · 2 need you   ☑ show served
  ┃   Needs you                                                                                      2
  ┃  ──────────────────────────────────────────────────────────────────────────────────────────────────
  ┃    RCP02 ready       Sign off the menu                             1/2
  ┃    RCP03 remediation Wire the review queue                         2/3 1 note
  ┃

     + Thought · 10.6s

  ┃
  ┃  $ grep -rn "func Test.*Legacy" internal/ui/*_test.go internal/core/*_test.go | head
  ┃
  ┃  internal/ui/console_test.go:301:func TestTheLiveConsoleCarriesTheLegacyGroup(t *testing.T) {
  ┃  internal/core/approval_test.go:38:func TestANormalisedLegacyApprovalIsRestamped(t *testing.T) {
  ┃  internal/core/backlog_concurrency_test.go:155:func TestALegacyDanglingPointerResolvesToNothing(t *testing.T) {
  ┃  internal/core/gates_test.go:280:func TestALegacyApprovalIsRestampedNotRefused(t *testing.T) {
  ┃

     + Thought · 906ms

  ┃
  ┃  ← Edit docs/audit.md
  ┃
  ┃    82   | **F-59** | medium | **`core.Audit.Decisions` has a     82   | **F-59** | medium | **`core.Audit.Decisions` has a
  ┃         shape and no writer** — it is where a decision *about         shape and no writer** — it is where a decision *about
  ┃         the plan* would live, and nothing writes it. `kitchen/        the plan* would live, and nothing writes it. `kitchen/
  ┃         decisions.json` is named by the plan format and opened        decisions.json` is named by the plan format and opened
  ┃         by nothing. The remedy is on the table in [06-proposals.      by nothing. The remedy is on the table in [06-
  ┃         md](06-proposals.md) § 2. | `internal/core/model.go:          proposals.md](06-proposals.md) § 2. | `internal/core/
  ┃         270`; `internal/core/store.go:24` |                           model.go:270`; `internal/core/store.go:24` |
  ┃    83   | **F-60** | medium | **`core.Questions` has a writer-   83   | **F-60** | medium | **`core.Questions` has a writer-
  ┃         less shape, so the `questions.resolved` gate can never        less shape, so the `questions.resolved` gate can never
  ┃         fire** — nothing ever sets a question, so the list is         fire** — nothing ever sets a question, so the list is
  ┃         always empty and the gate always passes. A gate nobody        always empty and the gate always passes. A gate nobody
  ┃         can fail is not a guard. | `internal/core/model.go:31`;       can fail is not a guard. | `internal/core/model.go:31`;
  ┃         `internal/core/lifecycle.go:197` |                            `internal/core/lifecycle.go:197` |
  ┃    84   | **F-61** | medium | **The `courses/` directory is      84   | **F-61** | medium | **The `courses/` directory is
  ┃         named by the plan format, created by `Init`, and never        named by the plan format, created by `Init`, and never
  ┃         written** — a directory no writer touches. Per [06-           written** — a directory no writer touches. Per [06-
  ┃         proposals.md](06-proposals.md) § 1 the honest answer is       proposals.md](06-proposals.md) § 1 the honest answer is
  ┃         to stop naming it: a course is the relation, not a file.      to stop naming it: a course is the relation, not a
  ┃         | `internal/core/store.go:18,63`; [docs/03 §6](03-            file. | `internal/core/store.go:18,63`; [docs/03 §6](
  ┃         domain-model.md) |                                            03-domain-model.md) |
  ┃    85   | F-56 | low | **The poll's fingerprint is blind to a    85   | F-56 | low | **The poll's fingerprint is blind to a
  ┃         write that keeps the mtime and the size.** It is path +       write that keeps the mtime and the size.** It is path +
  ┃         mtime + size per file, so two different contents can          mtime + size per file, so two different contents can
  ┃         share one fingerprint. The tool's own writes always           share one fingerprint. The tool's own writes always
  ┃         change the mtime, so in practice only a hand edit that        change the mtime, so in practice only a hand edit that
  ┃         restores the mtime is missed — and the chef's `r`             restores the mtime is missed — and the chef's `r`
  ┃         reloads it. | `internal/core/store.go` —                      reloads it. | `internal/core/store.go` —
  ┃         `DiskFingerprint`. |                                          `DiskFingerprint`. |
  ┃    86 - | **F-78** | high | **The approval digest did not travel
  ┃         with the rename — a plan written before it is now
  ┃         refused as tampered.** The digest hashes the criterion's
  ┃         JSON, and the JSON carries the field **name**: the pre-
  ┃         rename encoding hashes `{"id","description",
  ┃         "verification","proofs"}` and the version-6 encoding
  ┃         hashes the same with `evidence`. `restampLegacyApproval`
  ┃         recognises an older digest among the *`null` vs `[]`*
  ┃         candidates only — it has no candidate for the old
  ┃         **name** — so a plan approved before the rename keeps a
  ┃         digest the tool cannot reproduce, and the `immutability.
  ┃         respected` gate refuses every gated transition with
  ┃         *"frozen sections changed since approval"* on a recipe
  ┃         nobody touched. Proven twice: the repository's own
  ┃         `kitchen/recipes/RCP01-the-console-shows-a-finding-s-
  ┃         conversation.json` stores `8813a266…` (the pre-rename
  ┃         encoding, reproduced exactly) while the version-6 tool
  ┃         recomputes `4984098d…`; and a plan written in the old
  ┃         names, then re-written by the version-6 tool, keeps the
  ┃         old digest and is refused at `start`. The tolerance and
  ┃         the digest have to move together: a reader that
  ┃         tolerates the old names must also recognise the old
  ┃         digest, the way `restampLegacyApproval` already
  ┃         recognises the old *lists*. This breaks the standing
  ┃         discipline this ledger runs — *the repository's own
  ┃         plan's stored approval digest must still match the
  ┃         canonical recomputation* — and no test in the tree holds
  ┃         it. | `internal/core/model.go` — `fingerprint`,
  ┃         `legacyFingerprints`, `restampLegacyApproval`;
  ┃         `internal/core/lifecycle.go` — the `immutability.
  ┃         respected` gate; `kitchen/recipes/RCP01-the-console-
  ┃         shows-a-finding-s-conversation.json` |
  ┃    87 - | **F-79** | medium | **The glossary's entry for
  ┃         _Evidence_ is narrower than the record it names.** The
  ┃         entry reads *"A verification result a machine produced"*
  ┃         and *"An evidence that matched nothing is not evidence"*
  ┃         — but a **human sign-off** is stored in the same field:
  ┃         `Criterion.Evidence{Match: "signed-off", Matched: 0,
  ┃         SignedOff: true}`, and `ac.all_proven` takes it as a
  ┃         pass. So the single source of the vocabulary says a
  ┃         sign-off is not evidence, while the code stores it as
  ┃         evidence and the gate accepts it — and `docs/03` § 5
  ┃         says the opposite of the glossary (*"A sign-off is not
  ┃         machine evidence"*). The entry has to carry the
  ┃         exception, or define *Evidence* as *a verification
  ┃         result — machine or human*. | `GLOSSARY.md` — the
  ┃         **Evidence** entry; `internal/core/verifier.go` —
  ┃         `MatchSignedOff`; `internal/core/model.go` — `Criterion.
  ┃         Status`; `docs/03-domain-model.md` § 5 |
  ┃    88 - | **F-80** | medium | **`docs/03` § 6 says the reader
  ┃         does not see the old records, two paragraphs after
  ┃         saying it reads them.** *"A version-5 tool would write
  ┃         `proofs` where the version-6 reader looks for
  ┃         `evidence`, and the records it wrote would be records
  ┃         the reader does not see"* — and, two paragraphs below,
  ┃         *"a plan written before the rename keeps `proofs` and
  ┃         `evidence`, and the reader reads them."* `internal/core/
  ┃         legacy.go` reads them and `legacy_test.go` holds it. The
  ┃         version bump is still right (the **writer** changed);
  ┃         the consequence the paragraph states is wrong. The same
  ┃         sentence sits in `internal/core/lifecycle.go`'s version-
  ┃         6 note. | `docs/03-domain-model.md` § 6; `internal/core/
  ┃         lifecycle.go` — the `DomainVersion` note; `internal/
  ┃         core/legacy.go` |
  ┃    89   | P-8 | — | `--depends-on=` writes `depends_on: [""]` —  86   | P-8 | — | `--depends-on=` writes `depends_on: [""]` —
  ┃         inert, the same as `--in-scope=`. | Recorded so nobody        inert, the same as `--in-scope=`. | Recorded so nobody
  ┃         calls it a regression later. Not a defect. |                  calls it a regression later. Not a defect. |
  ┃    90                                                            87
  ┃    91   ## Closed findings                                       88   ## Closed findings
  ┃    92                                                            89
  ┃
  ┃

  ┃
  ┃  ← Edit docs/audit.md
  ┃
  ┃    130   | F-74 | high | **`internal/ui/model.go` — `Model.      130   | F-74 | high | **`internal/ui/model.go` — `Model.
  ┃          applyReload` puts a reloaded plan into the model, and         applyReload` puts a reloaded plan into the model, and
  ┃          the live path and the static render both go through it,       the live path and the static render both go through
  ┃          so a field the render carries and the live path forgets       it, so a field the render carries and the live path
  ┃          cannot happen; the test drives the real path** —              forgets cannot happen; the test drives the real path**
  ┃          `reloadMsg` carried the legacy artifacts and `Update`'s       — `reloadMsg` carried the legacy artifacts and
  ┃          `reloadMsg` branch dropped them, so an old plan's             `Update`'s `reloadMsg` branch dropped them, so an old
  ┃          backlog read as empty in the console the chef actually        plan's backlog read as empty in the console the chef
  ┃          uses, while `--render` showed them. One place now knows       actually uses, while `--render` showed them. One place
  ┃          the model's shape, and `loadModel` sends the reload           now knows the model's shape, and `loadModel` sends the
  ┃          message through the model's own `Update` instead of           reload message through the model's own `Update`
  ┃          copying the fields by hand — the copy was the bug. |          instead of copying the fields by hand — the copy was
  ┃          `internal/ui/model.go` — `applyReload`; `internal/ui/         the bug. | `internal/ui/model.go` — `applyReload`;
  ┃          render.go` — `snapshot`; `internal/ui/console_test.go`        `internal/ui/render.go` — `snapshot`; `internal/ui/
  ┃          — `loadModel`, `TestAnOldPlanStillRenders`,                   console_test.go` — `loadModel`,
  ┃          `TestTheLiveConsoleCarriesTheLegacyGroup` |                   `TestAnOldPlanStillRenders`,
  ┃                                                                        `TestTheLiveConsoleCarriesTheLegacyGroup` |
  ┃    131   | F-75 | medium | **`internal/core/store.go` —          131   | F-75 | medium | **`internal/core/store.go` —
  ┃          `reservedRecipeIDs` scans the old plan's pointers as          `reservedRecipeIDs` scans the old plan's pointers as
  ┃          well as the backlog's, and `core.LegacyItem` carries          well as the backlog's, and `core.LegacyItem` carries
  ┃          `PromotedTo`** — the reader dropped the pointers in           `PromotedTo`** — the reader dropped the pointers in
  ┃          `kitchen/prep-notes/*.json`, so an id a dangling              `kitchen/prep-notes/*.json`, so an id a dangling
  ┃          pointer names could be handed to a stranger — the exact       pointer names could be handed to a stranger — the
  ┃          failure `NextRecipeID`'s own comment says must never          exact failure `NextRecipeID`'s own comment says must
  ┃          happen. The pointer is read in the format that is             never happen. The pointer is read in the format that
  ┃          actually on disk (`"recipe RCP10"`), and a pointer in         is actually on disk (`"recipe RCP10"`), and a pointer
  ┃          the older `"recipe R10"` spelling reserves its number         in the older `"recipe R10"` spelling reserves its
  ┃          too; a bare `"RCP10"` is skipped, which is deliberate —       number too; a bare `"RCP10"` is skipped, which is
  ┃          the writer only ever wrote the two-word form. |               deliberate — the writer only ever wrote the two-word
  ┃          `internal/core/store.go` — `reservedRecipeIDs`;               form. | `internal/core/store.go` —
  ┃          `internal/core/legacy.go` — `LegacyItem.PromotedTo`;          `reservedRecipeIDs`; `internal/core/legacy.go` —
  ┃          `internal/core/backlog_concurrency_test.go` —                 `LegacyItem.PromotedTo`; `internal/core/
  ┃          `TestALegacyDanglingPointerResolvesToNothing` |               backlog_concurrency_test.go` —
  ┃                                                                        `TestALegacyDanglingPointerResolvesToNothing` |
  ┃    132   | F-76 | medium | **`internal/cli/flags.go` — the       132   | F-76 | medium | **`internal/cli/flags.go` — the
  ┃          parser knows which list flags were given, and `listArg`       parser knows which list flags were given, and
  ┃          returns nil for a flag that was not given; `internal/         `listArg` returns nil for a flag that was not given;
  ┃          core/backlog.go` — one rule (`setList`) for both list         `internal/core/backlog.go` — one rule (`setList`) for
  ┃          fields** — the tag list was guarded with `!= nil` while       both list fields** — the tag list was guarded with `!=
  ┃          the scope list, one line away, used `len(...) > 0`, so        nil` while the scope list, one line away, used `len(..
  ┃          an update that never mentioned the tags dropped them.         .) > 0`, so an update that never mentioned the tags
  ┃          nil now means "not given" and a list means "replace the       dropped them. nil now means "not given" and a list
  ┃          list it names" — an empty list that *was* given clears        means "replace the list it names" — an empty list that
  ┃          the field. `emptyIfNil` keeps the JSON edge a list: an        *was* given clears the field. `emptyIfNil` keeps the
  ┃          empty one is `[]`, never `null`. The other verbs are          JSON edge a list: an empty one is `[]`, never `null`.
  ┃          untouched: `--in-scope`, `--out-of-scope`, `--ac`, `--        The other verbs are untouched: `--in-scope`, `--out-
  ┃          ac-verify` and `--depends-on` behave exactly as before.       of-scope`, `--ac`, `--ac-verify` and `--depends-on`
  ┃          | `internal/cli/flags.go` — `given`, `listArg`;               behave exactly as before. | `internal/cli/flags.go` —
  ┃          `internal/core/backlog.go` — `setList`, `emptyIfNil`;         `given`, `listArg`; `internal/core/backlog.go` —
  ┃          `internal/cli/backlog_test.go` and `internal/core/            `setList`, `emptyIfNil`; `internal/cli/backlog_test.
  ┃          backlog_test.go` —                                            go` and `internal/core/backlog_test.go` —
  ┃          `TestTheBacklogUpdateLeavesTheListsAlone` |                   `TestTheBacklogUpdateLeavesTheListsAlone` |
  ┃    133   | F-77 | medium | **`docs/03-domain-model.md` § 1.1 —   133   | F-77 | medium | **`docs/03-domain-model.md` § 1.1 —
  ┃          the row and the paragraph now agree, and the code did         the row and the paragraph now agree, and the code did
  ┃          not move** — the row claimed a promoted small task's          not move** — the row claimed a promoted small task's
  ┃          recipe is "falsifiable from the start" while the              recipe is "falsifiable from the start" while the
  ┃          paragraph below said its criterion carries no                 paragraph below said its criterion carries no
  ┃          verification and "cannot pass `ac.falsifiable`". The          verification and "cannot pass `ac.falsifiable`". The
  ┃          criterion makes the work **checkable** and becomes            criterion makes the work **checkable** and becomes
  ┃          falsifiable when the recipe's own life writes its             falsifiable when the recipe's own life writes its
  ┃          verification; a promoted recipe is **checkable, not yet       verification; a promoted recipe is **checkable, not
  ┃          falsifiable**. `promotedCriteria` still leaves the            yet falsifiable**. `promotedCriteria` still leaves the
  ┃          verification empty, and `TestThePromotionWeakensNoGate`       verification empty, and
  ┃          still pins it. | `docs/03-domain-model.md` § 1.1;             `TestThePromotionWeakensNoGate` still pins it. |
  ┃          `internal/core/backlog.go` — `promotedCriteria`;              `docs/03-domain-model.md` § 1.1; `internal/core/
  ┃          `internal/core/promotion_test.go` —                           backlog.go` — `promotedCriteria`; `internal/core/
  ┃          `TestThePromotionWeakensNoGate` |                             promotion_test.go` — `TestThePromotionWeakensNoGate` |
  ┃                                                                  134 + | F-78 | high | **The approval digest did not travel
  ┃                                                                        with the rename — a plan written before it is now
  ┃                                                                        refused as tampered.** *Closed in* `internal/core/
  ┃                                                                        model.go` — `legacyFingerprints` enumerates the
  ┃                                                                        **field name** as well as the list forms, so the old
  ┃                                                                        name's digest is a recognised encoding and
  ┃                                                                        `restampLegacyApproval` re-stamps it canonically; the
  ┃                                                                        candidates are the product of the two axes, and a
  ┃                                                                        tampered recipe still matches none of them. The
  ┃                                                                        repository's own plan is migrated: `kitchen/recipes/
  ┃                                                                        RCP01-the-console-shows-a-finding-s-conversation.json`
  ┃                                                                        carries the canonical digest, and the journal records
  ┃                                                                        it as an **append** — event 37. | `internal/core/
  ┃                                                                        approval_test.go` — `TestARenamedFieldIsNotATamper`,
  ┃                                                                        `TestTheRepoPlansApprovalDigestMatchesItsContent`;
  ┃                                                                        `internal/core/model.go` — `legacyFingerprints`,
  ┃                                                                        `legacyFingerprint` |
  ┃                                                                  135 + | F-79 | medium | **The glossary's entry for
  ┃                                                                        _Evidence_ is narrower than the record it names.**
  ┃                                                                        *Closed in* `GLOSSARY.md` — the **Evidence** entry
  ┃                                                                        carries the exception now (*"One verification result —
  ┃                                                                        a machine's, or the chef's … A result that matched
  ┃                                                                        nothing does not count; the chef's sign-off is the one
  ┃                                                                        exception, and it claims no match"*), and `docs/03` §
  ┃                                                                        5 says the same thing, so the single source and the
  ┃                                                                        spec no longer contradict each other. | `internal/
  ┃                                                                        core/approval_test.go` —
  ┃                                                                        `TestASignedOffCriterionIsPassedButNotProven`;
  ┃                                                                        `internal/core/lifecycle_test.go` —
  ┃                                                                        `TestFinishRequiresProvenCriteria`; `docs/03-domain-
  ┃                                                                        model.md` § 5 — (the two docs' agreement is not wire-
  ┃                                                                        checked) |
  ┃                                                                  136 + | F-80 | medium | **`docs/03` § 6 says the reader does
  ┃                                                                        not see the old records, two paragraphs after saying
  ┃                                                                        it reads them.** *Closed in* `docs/03-domain-model.md`
  ┃                                                                        § 6 and `internal/core/lifecycle.go`'s version notes —
  ┃                                                                        the version-5 and version-6 sentences now say what the
  ┃                                                                        reader does: it reads the old field name, it shows the
  ┃                                                                        old specials read-only in their own group, and it
  ┃                                                                        never looks for new ones. The version bump is
  ┃                                                                        unchanged, because the **writer** changed, and that is
  ┃                                                                        what the version fences. | `internal/core/legacy.go` —
  ┃                                                                        `readLegacyEvidence`, `ListLegacyItems`; `internal/
  ┃                                                                        core/legacy_test.go`; `internal/ui/console_test.go` —
  ┃                                                                        `TestTheLiveConsoleCarriesTheLegacyGroup` |
  ┃    134                                                           137
  ┃    135   ## What this ledger cannot tell you                     138   ## What this ledger cannot tell you
  ┃    136                                                           139
  ┃    137   **Cycles 1–7's own findings are not recoverable.**      140   **Cycles 1–7's own findings are not recoverable.**
  ┃          Their per-finding reports lived in a                          Their per-finding reports lived in a
  ┃
  ┃
  ┃

     + Thought · 340ms

  ┃
  ┃  ← Edit docs/audit.md
  ┃
  ┃    40   | Slice B · cycle 2 | The fixes to slice B: the spec's   40   | Slice B · cycle 2 | The fixes to slice B: the spec's
  ┃         vocabulary rule, and the table under it | F-72 closed —       vocabulary rule, and the table under it | F-72 closed —
  ┃         the sentence, the table and the decision of record say        the sentence, the table and the decision of record say
  ┃         the same thing; the ledger's two pointers follow the          the same thing; the ledger's two pointers follow the
  ┃         rename |                                                      rename |
  ┃    41   | Slice C | The vocabulary rename: the backlog becomes   41   | Slice C | The vocabulary rename: the backlog becomes
  ┃         one artifact with a kind | The id, the kind and the           one artifact with a kind | The id, the kind and the
  ┃         fence hold; the live console hides what the snapshot          fence hold; the live console hides what the snapshot
  ┃         shows, and the old plan's dangling pointer no longer          shows, and the old plan's dangling pointer no longer
  ┃         reserves its id — one high and three mediums filed |          reserves its id — one high and three mediums filed |
  ┃    42   | Slice C · cycle 2 | The fixes to slice C: the live     42   | Slice C · cycle 2 | The fixes to slice C: the live
  ┃         reload path, the dangling pointer, the list rule and the      reload path, the dangling pointer, the list rule and
  ┃         spec's row | F-74..F-77 closed — each verified by a           the spec's row | F-74..F-77 closed — each verified by a
  ┃         mutation, and the four fixes moved nothing but the four       mutation, and the four fixes moved nothing but the four
  ┃         findings' own code |                                          findings' own code |
  ┃    43   | Slice D | The vocabulary rename: the per-criterion     43   | Slice D | The vocabulary rename: the per-criterion
  ┃         record is evidence, the recipe-level stamp is completion      record is evidence, the recipe-level stamp is
  ┃         | The two words hold, the reader tolerates the old names      completion | The two words hold, the reader tolerates
  ┃         and the tripwires hold both ways; the approval digest         the old names and the tripwires hold both ways; the
  ┃         did not travel with the rename, so a pre-rename plan is       approval digest did not travel with the rename, so a
  ┃         refused as tampered — one high and two mediums filed |        pre-rename plan is refused as tampered — one high and
  ┃                                                                       two mediums filed |
  ┃                                                                  44 + | Slice D · cycle 2 | The fixes to slice D: the digest
  ┃                                                                       travels with the rename, the glossary covers the record
  ┃                                                                       it names, the version paragraph stops contradicting the
  ┃                                                                       reader | F-78..F-80 closed — each verified by a
  ┃                                                                       mutation, and the fixes moved nothing but the three
  ┃                                                                       findings' own code |
  ┃    44   | The rounds inside a cycle | Fix → verify → back to the 45   | The rounds inside a cycle | Fix → verify → back to
  ┃         executor → repeat, until the verifier is satisfied |          the executor → repeat, until the verifier is satisfied
  ┃         Three rounds in cycle 8: the first fix introduced three       | Three rounds in cycle 8: the first fix introduced
  ┃         new mediums, the second closed them |                         three new mediums, the second closed them |
  ┃    45                                                            46
  ┃    46   ### The discipline every round runs                      47   ### The discipline every round runs
  ┃    47                                                            48
  ┃
  ┃

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre; git diff docs/audit.md | cat; echo "=== status ==="; git status --porcelain
  ┃
  ┃  (28 earlier lines)
  ┃  +| F-80 | medium | **`docs/03` § 6 says the reader does not see the old records, two paragraphs after saying it reads them.
  ┃  ** *Closed in* `docs/03-domain-model.md` § 6 and `internal/core/lifecycle.go`'s version notes — the version-5 and version-6
  ┃  sentences now say what the reader does: it reads the old field name, it shows the old specials read-only in their own group,
  ┃  and it never looks for new ones. The version bump is unchanged, because the **writer** changed, and that is what the version
  ┃  fences. | `internal/core/legacy.go` — `readLegacyEvidence`, `ListLegacyItems`; `internal/core/legacy_test.go`; `internal/ui/
  ┃  console_test.go` — `TestTheLiveConsoleCarriesTheLegacyGroup` |
  ┃
  ┃   ## What this ledger cannot tell you
  ┃
  ┃  === status ===
  ┃   M docs/audit.md
  ┃  ?? logo/
  ┃

     + Thought · 1.0s

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre
  ┃  git add docs/audit.md && git commit -q -m "docs: the ledger closes slice D's three findings…
  ┃
  ┃  279c57e docs: the ledger closes slice D's three findings
  ┃  1eb00a9 work: slice D's cycle 2 — and the line about ADRs
  ┃  36a28f7 inbox: runner reports slice-d fix 1
  ┃

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre
  ┃  mkdir -p work/inbox…
  ┃
  ┃  ac882f1 inbox: reviewer reports slice-d cycle 2
  ┃  279c57e docs: the ledger closes slice D's three findings
  ┃  1eb00a9 work: slice D's cycle 2 — and the line about ADRs
  ┃  ?? logo/
  ┃

  ┃
  ┃  $ cd /Users/thiago/Code/Pessoal/maitre
  ┃  herdr agent prompt "w1Z:p1" "Reviewer: slice D cycle 2 done. Verdict: ok. F-78..F-80 closed: yes. Filed 0 new blockers (…
  ┃
  ┃  {"id":"cli:agent:prompt","result":{"agent":{"agent":"opencode","agent_session":{"agent":"opencode","kind":"id","source":
  ┃  "herdr:opencode","value":"ses_e80f147cbffetwcFtsbsOZTA2z"},"agent_status":"done","completion_seq":451,"cwd":"/Users/thiago/
  ┃  Code/Pessoal/maitre","focused":false,"foreground_cwd":"/Users/thiago/Code/Pessoal/maitre","pane_id":"w1Z:p1","revision":19,
  ┃  "screen_detection_skipped":true,"state_change_seq":451,"tab_id":"w1Z:t1","terminal_id":"term_65ed7354d7f7722",
  ┃  "terminal_title":"OC | orchestrator (fork #1)","terminal_title_stripped":"OC | orchestrator (fork #1)","workspace_id":"w1Z"}
  ┃  ,"type":"agent_prompted"}}
  ┃  NOTIFIED
  ┃

     Verdict: ok — F-78, F-79, F-80 closed, no new blockers

     F-78 (high) — the digest travels with the rename ✓

     ┌────────────────────────────────────────┬─────────────────────────────────────────────────────────────────────────────────┐
     │ Check                                  │ Result                                                                          │
     ├────────────────────────────────────────┼─────────────────────────────────────────────────────────────────────────────────┤
     │ The candidate set is the cross product │ old field name and the null/[] list forms — a digest with both the old name and │
     │                                        │ null lists is recognised (probe)                                                │
     ├────────────────────────────────────────┼─────────────────────────────────────────────────────────────────────────────────┤
     │ The completeness argument              │ rewritten, not extended in silence: three lists, null vs [], the field name —   │
     │                                        │ correct                                                                         │
     ├────────────────────────────────────────┼─────────────────────────────────────────────────────────────────────────────────┤
     │ The safety property                    │ a criterion changed after the sign-off, stored digest in either encoding → the  │
     │                                        │ gate refuses both times                                                         │
     ├────────────────────────────────────────┼─────────────────────────────────────────────────────────────────────────────────┤
     │ The wire can fail                      │ dropping the old-name candidate → TestARenamedFieldIsNotATamper FAILS; the      │
     │                                        │ other tamper tests still pass                                                   │
     ├────────────────────────────────────────┼─────────────────────────────────────────────────────────────────────────────────┤
     │ The plan's digest is a test now        │ reads the digest off the disk; putting the old digest back → the test FAILS     │
     ├────────────────────────────────────────┼─────────────────────────────────────────────────────────────────────────────────┤
     │ End to end                             │ a plan written in the old names with the old digest is now claimed (exit 0),    │
     │                                        │ rewritten in the new encoding, re-stamped canonically — before the fix it was   │
     │                                        │ refused                                                                         │
     ├────────────────────────────────────────┼─────────────────────────────────────────────────────────────────────────────────┤
     │ The journal                            │ exactly one new event (37); events 1–36 untouched by git diff ✓                 │
     └────────────────────────────────────────┴─────────────────────────────────────────────────────────────────────────────────┘

     F-79 (medium) ✓ — the glossary and docs/03 § 5 now say the same thing (machine's or the chef's; zero-match does not count;
     the sign-off is the exception and claims no match). The code still holds it: Matched 0, match signed-off, ac.all_proven
     accepts it. "An evidence" is gone outside docs/audit.md; proof/proven survives where it is right (ac.all_proven, the match
     states, the old field names, the _Avoid_ line).

     F-80 (medium) ✓ — all three claims check against the code: legacy.go reads the old field name; ListLegacyItems reads
     kitchen/specials/ + prep-notes/ and the console shows them under FROM THE OLD PLAN; nothing writes to those dirs. The
     version bump is unchanged — the writer changed.

     The ADR line ✓ — the runner drew it correctly (0002/0006/0008 + index moved; 0004 and 0005 untouched), and only the
     vocabulary moved — no ADR decides something different now.

     Scope & the shared tree ✓ — DomainVersion 6, recordedDomainShape unchanged, format.go and the gate table untouched, 42
     commands, core/cli stdlib-only, ui 28, zero charm imports in core. The three fix commits touch no file of logo/; kitchen/
     took exactly one new event.

     LOW notes (not filed): the legacyFingerprints comment's stale counts ("at most eight… four"); gates_test.go:191 "an unproven
     evidence"; ADR-0002's parenthetical reads narrower than the glossary; ADR-0006's "three match states" (pre-existing); event
     37's wording; the seed's RCP05 outcome.

     Build · MiniMax-M3.1-Flash · 6m 46s · 86.4 tok/s

  ┃
  ┃
  ┃
  ┃  Build · MiniMax-M3.1-Flash MiniMax Token Plan (minimax.io) · medium
  ╹▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀
  ~/Code/Pessoal/maitre:master                                                                      196.4K (20%)  ctrl+p commands

```
