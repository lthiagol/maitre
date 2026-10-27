Runner: slice D done. Commits: 5 (e9710af, 1d378b2, 399608c, a9f9fc4, b4ff937).
Per-AC evidence: AC-01..AC-04 stamped. Pre-notify gate: pass (make lint clean,
gofmt -l . empty, make test green, make regen-goldens wrote nothing new, git status
clean). The journal is a clean append: events 1–34 untouched, 35 and 36 added.

**AC-01 — the per-criterion record is evidence** (`1d378b2`). `core.Proof` is
`core.Evidence`; `Criterion.Evidence []Evidence` is the JSON `evidence`; the view is
`evidence_summary`, computed by `EvidenceSummaryOf`, and `docs/03` § 8 names it in the
same commit (the table is tripwire-checked both ways, and it failed until the doc
moved). The verifier, the gates, the operations, the CLI's JSON and the console
follow. The journal's event kind stays `evidence` — it was already right. The records
stay append-only and the criterion's status still derives from its latest one.

**AC-02 — the recipe-level stamp is completion** (`e9710af`). `core.Completion`, JSON
`completion`, `r.Completion.CompletedAt`, set on `finish` exactly where it was. The
domain version moved 5 → 6 with it (`docs/03` § 6 says 6; the recorded fingerprint
moved with the format).

**AC-03 — the vocabulary in the docs and the glossary** (`399608c`, `a9f9fc4`,
`b4ff937`). `docs/03` § 5 is **Evidence**, the anatomy and the plan format say
evidence, § 8's row is `evidence_summary`. `docs/04`'s verbs say evidence. The
glossary's **Proof** is **Evidence** ("an evidence that matched nothing is not
evidence"), its `_Avoid_` line drops Evidence and gains Proof, the mapping row is
`| evidence | **evidence** |`, and a dash row names the **completion** with a line
under the table saying CPD has no term for it. The seeded demo plan and the two
console goldens say evidence too (`make regen-goldens`).

**AC-04 — the tripwires still hold** (evidence here, no code moved).
`go test ./internal/core ./internal/cli ./internal/ui -run 'Doc' -v` exits 0: the
derived-views table, the gate table, the lifecycles, the journal kinds, the operation
marks, the plan-format version and fields, the roles the code refuses, and the
command count (**42 commands: 36 verbs in 6 groups, plus 6 workspace commands**).
`make test` is green. `internal/core` and `internal/cli` pull **0** external modules,
`internal/ui` pulls **28**, and the domain imports no UI.

## The prose that kept the word "proof"

| Kept | Where | Why |
|---|---|---|
| `ac.all_proven` | the gate's name, in the code and in the gate table | a gate's name is a *rule*; renaming it is a rule change, not a rename |
| `proven` · `rejected` · `unknown` | the match states, and `EvidenceSummary.Proven` | they name the **state a criterion is in**, which the record derives — not the record |
| "a criterion is **proved** by an existing suite" | `docs/06`'s brownfield proposal | the verb is the act, not the record |
| "Proof" in the glossary's `_Avoid_` line | `GLOSSARY.md` | the point of the line is to warn the reader off the word the code dropped |
| `proofs` / `evidence` as **old field names** | the legacy reader, `docs/03`'s version paragraph, the journal's two migration events | they name the format *before* the rename, which is the point |

One place could not keep the old wording: the sign-off was "a sign-off is not a
proof", and under the new words that would have been **false** — a sign-off *is* an
evidence (one that claims no match, with its own `signed-off` state). It now says
"a sign-off is not **machine** evidence".

## Notes for the reviewer

1. **AC-02 landed before AC-01.** The name `Evidence` had to be free for the
   per-criterion record, so the recipe-level stamp moved first. Each commit is a
   whole AC and the suite is green at every one of them; only the order differs from
   the brief's numbering.

2. **The version moved in AC-02 and the fingerprint moved twice.** The slice changes
   the format twice (the stamp, then the record), and `core.DomainVersion` is one
   number. It went 5 → 6 in the commit that changed the format first, so a stale
   version-5 tool is fenced off from the start; the recorded fingerprint moved again
   in AC-01 (`cff4269d589690ce`) because the format moved again. The version stays
   **6**, as the brief says.

3. **The reader tolerates the old field names** (`internal/core/legacy.go`,
   `internal/core/legacy_test.go`): a plan written before the rename keeps `proofs`
   and `evidence` for the stamp, and the reader reads them — the same tolerance the
   ids and the backlog's two directories already had. A file carrying both names is
   a broken plan, and the **new** name wins. The repo's own plan was migrated by
   hand and the journal says so (`evidence -> completion`, `proofs -> evidence`).

4. **A collision the rename created, and how it was settled.** `docs/02`'s
   principle template has a section labelled **Evidence** — the failure that
   motivated the principle. Harmless under the old words; under the new ones it
   would have been *one word for two things*, which is what the glossary exists to
   prevent. The label is now **The failure**, in all thirteen principles. Say if
   you would rather have another word.

5. **The sweep went wider than `docs/03` and `docs/04`** — AC-03 asks for *any*
   prose that says "proof" where it means the record, so it touched `docs/01`
   (G1), `docs/02` (P5 and the label), `docs/05` (four mentions), `docs/06` (the
   brownfield proposal), ADR-0006 / 0002 / 0008 and the decisions index,
   `AGENTS.md` and `prototypes/README.md`. **`AGENTS.md` was wrong about the
   code**: it named `ProofSummary` as a derived view, which has not existed since
   this slice; it names `EvidenceSummary` now. The ADRs' titles moved with their
   bodies, and `AGENTS.md`'s decisions row and the index follow.

6. **Two things I did not touch, and they are yours.**
   - **`GLOSSARY.md`'s brigade section is stale from slice A**: it still says *line
     cook*, *taster* and *expo*, and its mapping rows still map CPD's runner to
     *line cook*. The code refuses `line-cook` and the roles are `runner`,
     `reviewer`, `orchestrator`. The brief says slice D does not touch the roles, so
     I did not — but the glossary is the single source, and it is the one doc that
     still contradicts the code on something other than this slice.
   - **`docs/07-vocabulary.md`** still says the plan format went "2 → 3, then 4, then
     5". The version is 6 now. Slice C's runner left the same two docs to the
     orchestrator, so I left them again.

7. **The seeded demo plan's text moved, so the console's goldens moved.**
   "Sharpen the proof contract" is now "Sharpen the evidence contract", and the two
   goldens say so. `make regen-goldens` is clean on a second run.

LOW notes: none.
