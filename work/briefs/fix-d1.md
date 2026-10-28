You are the **runner** — the same lane that wrote slice D. The orchestrator's pane id is
`w1Z:p1`; use it verbatim in the notify command below.

## The fix cycle for slice D — three findings

The reviewer filed **F-78 (high)**, **F-79** and **F-80**. The orchestrator ruled on
all three: **accepted**. Three findings, three commits — one each.

### F-78 (high) — the approval digest did not travel with the rename

The digest hashes the criterion's JSON, and the JSON carries the field **name**. The
pre-rename encoding hashes `{"id","description","verification","proofs"}`; the
version-6 encoding hashes the same with `evidence`. The orchestrator reproduced it:

```
the repo's own RCP01 stores   8813a266eaac833a280c99016d53bbceb85b95cf1e3b0121354f330de4072edc
the version-6 tool recomputes 4984098d94858a1ec142f89c987e41d5c232bc2ad85c7699893ef9b20fd75fe1
the digest of the OLD shape   8813a266eaac833a280c99016d53bbceb85b95cf1e3b0121354f330de4072edc
```

The stored digest **is** the old field name's digest. And the `immutability.respected`
gate reads the difference as **tampering** — so it refuses every gated transition of
a recipe nobody touched, with *"frozen sections changed since approval"*.

**The tolerance and the digest have to move together.** A reader that tolerates the
old names must also recognise the old digest, **the way `restampLegacyApproval`
already recognises the old lists**.

**The fix:**

1. **`legacyFingerprints` gains the old-name candidate.** It already enumerates the
   encodings an older build could have written — the `null` vs `[]` list forms. Add
   the *field name* form: the digest of the criterion encoded with `proofs` instead of
   `evidence`. The set must stay **complete** — the comment above the function says
   why the current set is complete, and that reasoning has to be rewritten, not
   extended in silence.
2. **The repository's own plan is migrated**, so its stored digest matches the
   canonical recomputation. The journal records the migration as an **append** — the
   historical events keep their old ids and their old field names.
3. **A test holds the standing discipline:** *the repository's own plan's stored
   approval digest must still match the canonical recomputation.* No test in the tree
   holds it now, which is why the finding was possible.

**Prove it by mutation.** Remove the old-name candidate and the test must fail.

> **A note on what a digest is for.** The gate's whole purpose is to catch a
> *material edit*. A rename is not an edit — the content did not change, only the
> encoding. A digest that cannot tell those two apart is not a guard; it is a
> tripwire pointed at the wrong thing.

### F-79 (medium) — the glossary's entry is narrower than the record it names

The **Evidence** entry reads *"A verification result a machine produced"* and *"An
evidence that matched nothing is not evidence"*. But a **human sign-off** is stored in
the same field — `Criterion.Evidence{Match: "signed-off", Matched: 0, SignedOff:
true}` — and `ac.all_proven` takes it as a pass.

So the **single source** of the vocabulary says a sign-off is not evidence, while the
code stores it as evidence and the gate accepts it. And `docs/03` § 5 says the
**opposite** of the glossary — *"a sign-off is not machine evidence"*.

**The fix:** the entry carries the exception, or **Evidence** is defined as *a
verification result — a machine's, or the chef's*. The two docs must say the same
thing, and the sentence must be **true**: a sign-off claims **no match** and is still
evidence, so *"an evidence that matched nothing is not evidence"* is false as it
stands. Carry the exception in the line itself.

**Also fold in the reviewer's grammar note:** *"an evidence"* reads badly. The entry
should read like English.

> The wording came from the orchestrator's slice D brief, which quoted the old
> **Proof** entry's shape. The old entry could say "a proof that matched nothing is
> not a proof" because a proof was only ever a machine's. The new word covers both,
> and the sentence has to change with it.

### F-80 (medium) — `docs/03` § 6 contradicts itself

> *"A version-5 tool would write `proofs` where the version-6 reader looks for
> `evidence`, and the records it wrote would be records the reader does not see"*

— two paragraphs below:

> *"a plan written before the rename keeps `proofs` and `evidence`, and the reader
> reads them."*

The reader **does** read them; `internal/core/legacy.go` reads them and
`legacy_test.go` holds it. **The version bump is still right** — the *writer* changed,
and that is what the version fences. **The consequence the paragraph states is
wrong.** Fix the sentence.

**The same sentence sits in `internal/core/lifecycle.go`'s version-6 note.** It moves
with the paragraph.

**And fold in the reviewer's related note (g):** the version-5 sentence says the
reader *"does not look"* in `kitchen/specials/`, and `legacy.go` reads that too. Same
shape, same fix.

### The LOW notes — the rulings

| Note | Ruling |
|---|---|
| *"an evidence"* — grammar | **folded into F-79.** The entry is rewritten anyway. |
| `docs/07`'s alternatives table says *"work that closes with no proof"* | **keep it.** It names the *idea* of a proof, not the record — the same distinction that kept `ac.all_proven` and the match states. `docs/07` is the orchestrator's to write. |
| `research/01-cpd-foundation.md` still says *"a proof that matched nothing is not a proof"* | **leave it.** `research/` is about v1 — history, not drift. |
| `GLOSSARY.md`'s brigade section is stale from slice A | **leave it.** slice E rewrites the glossary. |
| `scripts/stamp-plan`'s usage line does not mention the field-name migrations | **fix it, in this cycle.** The script just recorded two field-name migrations; its usage line should say so. One line of text. |
| `AGENTS.md` and `work/README.md` carry the orchestrator's post-mortem | **leave them.** They are the orchestrator's, and they are right. Do not revert them. |

## The tree is shared — read this before you commit

**Another session is working in this repository right now.** It is drawing the
mascot, and it has files in the tree that are **not yours** (`logo/`).

| | |
|---|---|
| **Stage only your own files** | name them explicitly. **Never `git add -A`**, never `git add .` — a sweep would take the other session's work into your commit |
| **`logo/` and anything you did not write** | leave it alone. Do not delete it, do not commit it, do not rename it |
| **`kitchen/`** | the other session may write to the plan too. **Re-read the files you are about to touch**, and check `git status` immediately before you commit |
| **The journal** | append-only. If a command you run writes to the plan, that is an event — and an event nobody meant is noise in the record |

**If you find a file changed that you did not change, stop and report it** rather
than committing around it.

## Not yours — do not touch

| | |
|---|---|
| `docs/audit.md` | the reviewer's to write. It files F-78..F-80 and closes them. |
| `docs/07-vocabulary.md` | the decision of record — the orchestrator's to write. |
| `GLOSSARY.md`'s brigade section and its CPD rows | slice E's rewrite. |
| `research/` | history. |
| `logo/`, and the mascot's files | another session's. |
| The gate table, the plan format, the command count (42) | none of the three findings touches them. If a fix would, **stop and report it**. |

## The gate

```bash
make lint && gofmt -l . && make test && git status --short
```

Green. `git status` will show the other session's files — that is expected, and they
are **not** yours to commit.

## Three commits, one per finding

```
Per-AC evidence:
- F-78: the old field name is a recognised encoding; go test ./internal/core -run 'Fingerprint|Legacy' -v exit 0
```

## The report — the receipt comes first

```bash
mkdir -p work/inbox
cat > work/inbox/runner-slice-d-fix-1.md <<EOF
Runner: slice D fix 1 done. Commits: <n>. Fixed: F-78, F-79, F-80. Notes: <anything the reviewer should know>.
EOF
git add work/inbox && git commit -q -m "inbox: runner reports slice-d fix 1"
herdr agent prompt "w1Z:p1" "Runner: slice D fix 1 done. Commits: <n>. Fixed: F-78, F-79, F-80. Notes: <anything the reviewer should know>."
```

**Commit the receipt before you nudge.** And **one write per file** — the rule that
would have kept your last turn alive.
