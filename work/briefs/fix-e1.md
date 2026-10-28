You are the **runner** — the same lane that wrote slice E. The orchestrator's pane id is
`w1Z:p1`; use it verbatim in the notify command below.

## The fix cycle for slice E — five findings

The reviewer filed **F-81**, **F-82**, **F-83**, **F-84** and **F-85**. The
orchestrator ruled on all five: **accepted**. Five findings, five commits — one each.

**Two of them are about the *glossary*, and the glossary is the single source for the
mapping.** F-81, F-82 and F-83 are all about the same file, and they are all about
one thing: the glossary is telling a reader what to call things *now*, and it must do
that without tripping over the names it retired.

### F-81 (medium) — the console's retired name is nowhere, and the proof did not reproduce

Two halves, and both matter.

**The name.** Every other renamed term carries its old name **once**, in its
`_Avoid_` line — `Line cook`, `Taster`, `Expo`, `Proof`, `Evidence`, `Special`,
`Prep note`. The console's old name, *maître d'*, is in **no line of the file**, so
the single source cannot tell a reader what the console was once called — and ADR
0009 § 2 records that rename.

**The proof.** Your AC-02 commit message offers a `grep -c` over `GLOSSARY.md` for
*"line cook"*, *"taster"*, *"expo"* and *"maître d'"* as the evidence that each
retired name hits only an `_Avoid_` line. **As written, that command returns 0 for
all four** — the first three because the file capitalises them, the fourth because
the word is absent.

> A commit that offers a command as its evidence must offer a command that
> **reproduces**. A command that returns nothing proves nothing, and offering it
> anyway is worse than offering none: it looks like a check.

**The fix:** the console's entry carries its retired name in its `_Avoid_` line, like
every other retired name. And the proof is a command that **works** — and the commit
message shows its actual output, not a claim about it.

### F-82 (medium) — the console is filed among the crew

The file defines the crew as *"the people and agents working it"*, its scope table
lists the jobs as chef · sous-chef · runner · reviewer · orchestrator, and **"### The
crew" has six entries — the sixth being the console.**

The console is a **surface**, not a person or an agent: ADR-0003 § 4 says *"the
console is not an agent"*, ADR 0009 § 1 gives the jobs as those five, and ADR 0009
§ 4 puts the console's names under **"What is strings only"**. F-72 closed on the
console **never being a job** — this entry crosses the line the file's own new tables
draw.

**And the entry's header is `**Maitre**`** — a spelling the vocabulary does not use.
The console is `maitre` (`:32`).

**The fix:** the console's entry leaves the crew. It belongs where ADR 0009 puts it,
with the things that are only strings. And its name is spelled the way the code spells
it.

### F-83 (medium) — two retired words are named twice

`:59` says the backlog item *"replaces the two file types the plan used to hold — a
*special* and a *prep note* —"* and `:63` lists *"Special, prep note"* again.

The rule this slice set is that a retired word is named **once**, in the `_Avoid_`
line that warns the reader off it. The prose is where a reader meets a retired word
with nothing warning them.

**The fix:** the prose stops naming the retired words. The `_Avoid_` line keeps them —
that is where the warning lives.

### F-84 (medium) — the spec lists the fingerprint's facts twice, and the two lists disagree

`docs/03` `:451` says *"the id prefixes, the directories, the root files and the field
names on disk"*, and `:550` says *"the id prefixes, the directories, the field names
on disk"*. The code hashes the root file names as well (`internal/core/format.go`).

Your own commit updated one of the two lists and not the other, so the document now
says two different things about what the fingerprint covers — and **nothing checks the
prose**, which is why it survived your gate.

**The fix:** the two lists say the same thing, and they say what the code does.

### F-85 (medium) — the spec still states the coupling this slice split

`docs/03` `:548` says **"What the stamp covers is exactly what the fingerprint
hashes"**, and `:563` says *"the fingerprint hashes it, so a format change moves the
version on its own instead of by hand"*. The `DomainVersion` comment in
`internal/core/lifecycle.go` says the same.

**The code now says the opposite** — and this slice's own commit is the counterexample:

> *"If the change is in what the tool writes, bump DomainVersion as well — **the
> version fences the writer, and the shape alone fences nothing.**"*

Your judgement was **right**: the version stayed at 6 because the writer did not
change. But the prose still teaches the belief the split removed, and a reader of it
would bump the version every time the fingerprint moves — which is the mistake the
split exists to prevent.

**The fix:** `docs/03` and the `DomainVersion` comment say what the code says. **The
version fences the writer; the fingerprint is the tripwire that says the facts
moved.** A change to the *facts list* moves the fingerprint and leaves the version
alone; a change to *what the tool writes* moves both.

> `docs/07-vocabulary.md` is the orchestrator's to write, and it has been corrected in
> the same spirit. Yours are `docs/03` and the comment.

## The two LOW notes

| Note | Ruling |
|---|---|
| `README.md`'s decisions table has no row for 0009 | **fixed by the orchestrator** — `README.md` and `work/README.md` are the orchestrator's, and they are corrected. Do not touch them. |
| `work/README.md` still said *"the taster"*, *"the line cook"*, version 4 and *"slice E is queued"* | **fixed by the orchestrator**, the same commit. Those were the last old names in the tree. |

## The tree is shared — read this before you commit

**Another session is drawing the mascot in this repository** and has files in the tree
(`logo/`).

| | |
|---|---|
| **Stage by name** | name your files explicitly. **Never `git add -A`**, never `git add .` |
| **`logo/`, and anything you did not write** | leave it alone |
| **`kitchen/`** | re-read before you commit; the other session may write to the plan |

**If a file changed that you did not change, stop and report it.**

## Not yours — do not touch

| | |
|---|---|
| `docs/audit.md` | the reviewer's to write. It files F-81..F-85 and closes them. |
| `docs/07-vocabulary.md` | the decision of record — the orchestrator's to write. |
| `README.md`, `work/README.md` | the orchestrator's. Already corrected. |
| `GLOSSARY.md`'s **crew** section heading and the console's entry | **yours** — that is F-82. Everything else in the glossary is yours too. |
| The gate table, the plan format, the command count (42), `core.DomainVersion` (6) | none of the five findings touches them. If a fix would, **stop and report it**. |

## The gate

```bash
make lint && gofmt -l . && make test && git status --short
```

Green.

## Five commits, one per finding

```
Per-AC evidence:
- F-81: the console's retired name is in its _Avoid_ line, and the proof command reproduces; go test ./internal/core -run 'Doc' -v exit 0
```

## The report — the receipt comes first

```bash
mkdir -p work/inbox
cat > work/inbox/runner-slice-e-fix-1.md <<EOF
Runner: slice E fix 1 done. Commits: <n>. Fixed: F-81..F-85. Notes: <anything the reviewer should know>.
EOF
git add work/inbox && git commit -q -m "inbox: runner reports slice-e fix 1"
herdr agent prompt "w1Z:p1" "Runner: slice E fix 1 done. Commits: <n>. Fixed: F-81..F-85. Notes: <anything the reviewer should know>."
```

**One write per file, and commit as soon as the file is coherent.**
