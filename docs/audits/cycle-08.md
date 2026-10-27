# Cycle 8 — the tasting notes

The per-cycle review report from the audit of 2026-10-26, covering fix cycle 7,
tasting-note triage (T3) and the two rounds the fixes needed.

**It survived by luck.** The nightly scratch was wiped mid-cycle, so every earlier
per-cycle report is gone; this one carried the state forward — and it lived in
`/tmp`, which is volatile. It is rescued here so the reasoning behind the ledger's rows
is not lost with the next cleanup.

The findings themselves are in the ledger: [../audit.md](../audit.md).

---

# Tasting notes — cycle 8 (fix cycle 7 + T3 + the two rounds the fixes needed)

**Taster:** independent session (the "Taster" pane). Read the code, attacked every fix,
and drove the shipped binaries. The fixes were applied by the line cook; a second,
independent verifier agent (its own session) re-checked them twice.

**Verdict: the fixes hold.** Every note filed for this cycle is closed. The one thing
worth the next cycle's attention is the line cook's own P-1, which is now closed in code
but not in the *record* — see §4.

> This file is the surviving record. The nightly scratch (`/tmp/kb-audit`,
> `/tmp/herdr-orc`) was wiped mid-cycle, so the earlier per-cycle reports are gone; this
> one carries the state forward. It is in /tmp, which is volatile — say the word and I
> will put it in the repo instead (I did not want to invent a home for a review record in
> the project's tree without your say-so).

---

## 1. What holds

| # | Note | Verdict | How it was checked |
|---|---|---|---|
| F-14 | `scale recipe update --depends-on` | **holds** | `depends_on: ["R02"]` lands in the file; the list is replaced, not appended; the flag alone changes nothing; the dependency gate then refuses the claim; `scale path` says `needs R02`. |
| F-19 | a refusal is journalled | **holds** | A gate refusal +1 event (with the event, the gates and the reason); an illegal transition +1; a frozen-section edit +1; "no such recipe", "a claim needs an agent", a missing argument, an unknown event, "that event has its own verb", a bad flag → **+0**. One dismissal → +1, never two. The journal's kind enum is closed and the doc lists every kind. |
| F-45 | the diagram's `close` edge | **holds** | It lists all five gates the code runs. A shared check now compares **every** edge against the transition table, both ways. |
| F-46 | the `request_changes` edge | **holds** | The edge is there with its gate; the same shared check covers it. |
| F-47 | the command-count wire | **holds** | It now checks **every** occurrence of "N commands" in every file that states it. Changing the *first* of AGENTS.md's two counts fails it: *"the file says "45 commands" in one place and "46 commands" in another"*. |
| F-48 | a refusal is never dropped | **holds** | The whole reason is on screen at **40, 60, 80, 120 and 200** columns — as banner rows when it does not fit, on the summary line when it does. The screen still draws exactly `height` rows. |
| T3 | the console triages a tasting note | **holds** | The chef selects a note with `n`/`N` (wrapping, over the **open** notes only); the note wears the selected-row style; `space` opens a prompt offering **dismiss** and **comment** only; a dismissal writes the reason and journals exactly one event; a comment lands on the note, authored by the chef, in its thread; `esc` leaves the plan untouched; the footer names `n/N note`; `grep -n 'knownGap\|excepted'` returns **nothing** — the named exception is deleted from the invariant test, not replaced. |
| — | the console cannot resolve a note | **holds** | `grep -rn ResolveFinding internal/ui` → nothing. The prompt's choices are dismiss and comment. The console's every write is approve, close, promote ×2, dismiss, comment, record-round, pass-review. **The fixer can still resolve through the CLI** — `scale finding resolve R01 F-01 --note "…"` → status `fixed`, the resolution written, the journal +1. |
| — | the cursor and the gate see the same notes | **holds** | The cursor moves over `core.OpenFindingsOf(rec)` — open findings of **any** phase; `review.clean` and `review.no_open_findings` count that same list. A recipe with one external and one self open note: the cursor sees 2, the gate counts 2. `review.findings_exist` counts the **external** subset only, and deliberately so — it answers "has the other party's complaint been answered?", which is the question of entering and leaving remediation. The cursor answers "what still blocks serving this?" — the right one. |
| F-49 | a title change must not leave two files | **holds** (round 2) | One file after a rename; `show` returns the new title; `list` counts the recipe once; the console draws it once; the dependency the gate reads is the one the chef wrote. Break tests: two recipes with the same slug, a child id, a special, a prep note, a failing save — all hold. |
| F-50 | a note's outcome is written once | **holds** (round 2) | A dismissal of a fixed note is refused ("note F-01 is fixed — a dismissal applies to an open note") and the fixer's record survives; a resolution of a dismissed note is refused likewise; both refusals are journalled; the console's triage cannot overwrite a note the agent closed. |
| F-51 | the version guard must see the rules it protects | **holds** (round 3) | `core.DomainVersion` is **2** and the fingerprint moved (`0b55ac7113d0d99c`). The chef moves and the note rule are now visible to it, so an older tool is refused rather than allowed to apply weaker rules. |
| F-52 | the removal must not delete a file of another kind | **holds** (round 3) | The copy is matched on the **whole id** (`strings.EqualFold(base, id)`), not the normalised number. A rogue `C01-rogue.json` in the recipes directory now **survives** a rename. |
| F-53 | the reopen must not erase the fixer's record | **holds** (round 3) | The reopen carries the erased resolution into the comment it already writes: *"it still drops findings (it was closed as "fixed it properly in the loader")"*. The substance survives, as the reopen's own doc always said it would. |
| F-54 | a failed removal must not be a silent partial write | **holds** (round 3) | The removal failure is **non-fatal**: the operation now succeeds, the journal grows, the leftover is cleaned by the next write, and the reader shows the **new** title (the tie-break is no longer the name order). |

## 2. Did the fixes break something? — no

`make clean && make build` clean · `make test` green (core 26.7s · cli 9.5s · ui 104.0s) ·
`make lint` clean · `gofmt -l .` empty · `go test -race` on all three packages green ·
`go test -count=2 ./...` green · the demo plan seeds and renders (7 recipes, 2 need you) ·
`scale demo --chaos` seeds and renders (17 recipes, 3 need you) · the repo's own
`kitchen/` parses, renders, and its stored approval digest still equals the canonical
recomputation (`8813a266eaac`) · the goldens have not moved (both still Oct 25 22:16) ·
`TestDomainShapeIsVersioned` passes with the new fingerprint.

## 3. The two things I checked with the shipped binary, in the orchestrator's four questions

| Question | Answer |
|---|---|
| Can the console really not resolve a note? | **Confirmed it cannot**, and the fixer's move is untouched: `scale finding resolve` still works (status `fixed`, the resolution written, the journal +1). |
| Does the note cursor and the gate see the same notes? | **The same set** — the cursor moves over exactly the list `review.clean` counts, open findings of any phase. |
| Are the refusals journalled in the right places, and only there? | **Yes** — the plan's refusals (gate, illegal transition, frozen section) are journalled, exactly one event each; a bad request is not. |
| Is the refusal visible at every width? | **Yes** — 40, 60, 80, 120, 200. |

## 4. What is still open

Nothing that I would call a blocker. Two things, both for the record rather than the code:

1. **The line cook's pending file is gone** (`/tmp/herdr-orc/linecook-pending.md` — wiped
   with the scratch). It held eight items it had *filed and not fixed*: P-1 (now fixed
   ✓), P-2 (overlay refusals are not journalled — the console's banner and the journal
   use "refusal" for two different sets), P-3 (`refusalText` uses a type assertion, not
   `errors.As` — a wrapped `GateError` would print agent-written prose on the summary
   path), P-4 (a comment on a finding that does not exist is accepted and no surface
   renders it), P-5 (a CLI comment without `--author` writes `"author": ""`), P-6 (the
   dependency count in AGENTS.md — **now fixed**, it reads 28 ✓), P-7 (`m.detail` is
   assigned and never read), P-8 (`--depends-on=` writes `depends_on: [""]` — inert, the
   same as `--in-scope=`). The line cook is mid-flight on P-4/P-5 as I write this.
2. **The version bump's scope.** `DomainVersion` is 2 ✓, which protects every plan the
   new build writes from now on ✓ — and *not* the plans already stamped 1 (the repo's own
   `kitchen/` is stamped 0 ✓ and is re-stamped on its next write ✓). The verifier's
   measurement stands: the bump buys the future, not the corpus ✓.

## 5. Coverage

| Area | What I did |
|---|---|
| The eleven fixes (F-14, F-19, F-45, F-46, F-47, F-48) | Read the code, and attacked each with the shipped binary. |
| T3 | Read `internal/ui/note.go`, the grown prompt in `internal/ui/review.go`, the detail pane's note rendering, `core.ChefMoveOf`, `core.OpenFindingsOf`; drove the whole flow (select → triage → record → reload → see it) with `maitre --keys`. |
| The four the orchestrator named | Verified by command, with the output quoted above ✓. |
| The three the second round found (F-52, F-53, F-54) | Reproduced by me before they were fixed ✓; re-verified after ✓. |
| The tripwires | Corrupted the real docs in a scratch copy — a gate renamed, a gate invented, a kind removed, a kind invented, a mark flipped, the count changed in either file, a diagram gate dropped, the `request_changes` edge removed, the operations heading renamed — and confirmed each one fails its wire ✓. |
| The suite | Cold build, `make test`, `make lint`, `gofmt -l .`, `go test -race` on core/cli/ui, `go test -count=2 ./...`, the two canaries, the repo's own plan, the goldens' mtimes. |
