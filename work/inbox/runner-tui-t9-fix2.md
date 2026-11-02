Runner: T9 fix 2 done. Commits: 3 — `2fb29f2` (AC-01, the check), `b170d0c` (AC-02, the names), `107a961` (AC-03, the row and the gate table).

**The check.** `scripts/check-slice-numbers.sh`, run by the target **`make slices-check`**, wired into `check` and `check-race` — so `make check` runs it and cannot skip it. It reads **thirteen places in the plan's three files** and requires them to be one number:

| Place | What it reads |
|---|---|
| `work/tui/README.md` — the state table | the Sous-Chef row's slice |
| `work/tui/README.md` — *"what each ruling lands on"* | the slice of the tab's **nine** rulings — Q16–Q22, Q26, Q27 |
| `docs/05-human-interface.md` — the tab table | the tab's `*(T10 — planned)*` |
| `work/tui/plan-automation.md` — § 10 | Q24's ruling, the number attached to *"the Sous-Chef tab"* |
| `work/tui/plan-automation.md` — § 11 | the slices table's *"**T10 — the Sous-Chef tab**"* row |

On a disagreement it **fails and names the places and the numbers** — the files and the two numbers, one line each, so the reader knows where to look. And it **refuses to pass having found nothing**: a place whose number cannot be read fails the check (*"the slice check found no number in docs/05-human-interface.md (the tab table). The file changed shape — the check must be taught the new shape, or it proves nothing."*), the way a test that matched no test fails the gate.

**The proof, by command, each mutation reverted:**

| Mutation | The check |
|---|---|
| `docs/05`'s tab table: T10 → T9 | **exit 1**, the places and the two numbers |
| the README's state table: T10 → T9 | **exit 1** |
| the README's ruling Q19: T10 → T9 | **exit 1** — the ruling that does not say *tab* is checked too |
| the design's § 10 Q24 ruling: T10 → T9 | **exit 1** |
| the design's § 11 table: T10 → T9 | **exit 1** |
| `docs/05`'s tab table with no number at all | **exit 1** — the check found nothing, and said so |

**The name.** `TestTheWriteIsAtomic` is gone: it held two things under one name, and its name claimed a third (the atomicity) that it did not hold — the two witnesses hold that. It is now **one claim, one test**:

| Test | What it holds | The mutation that turns **only** it red |
|---|---|---|
| `TestAWriteLeavesNothingBehind` | the **tidiness**: a successful write and a failed one both leave the folder holding the record and nothing else | the write does not remove its temporary file |
| `TestAWriteIntoAFolderThatIsNotThereIsRefused` | the **refusal**: a write into a folder that is not there fails, and the file that was there before is untouched | the write creates the folder it was asked for |
| `TestAHalfWrittenFileCannotLoseTheChefsChoices` | the **atomicity**, the fault in the bytes | the write overwrites the old file in place |
| `TestARenameThatFailsCannotLoseTheChefsChoices` | the **atomicity**, the fault in the last step | the rename replaced by a write over the target |

The two witnesses now assert **only** the property they are named for — *the file that was there before is whole* — so the tidiness assertion no longer rides inside a test whose name says nothing about it.

**The gate.** The target is `make slices-check`; `AGENTS.md`'s gate table carries the row, and `make check` and `make check-race` run it (right after `versions-check`). DomainVersion: **7**.

The gate: `make check` — pass; `make check-race` — pass.

Notes for the reviewer:

1. **The check's one written claim is the list of the tab's questions** — Q16–Q22, Q26, Q27. Those are the questions the plan asks *about the tab*, and they are named in the script. Everything else the check reads is derived from the files. A ruling about the tab that is not in that list would not be checked — a new question about the tab must be added to the list, and the list is the one thing a reader has to read to know what the check holds. I could not derive it: the rules about the tab's state (Q19), the journal (Q20), the startup task (Q21) and the stop (Q22) never say *tab*, so no reading of the text can pick them out — hence the list.
2. **The check is a shell script**, beside `check-golden-claims.sh` and `dep-budget.sh`, in the house's shape. It is a check and not a test because its subject is the plan's own files, not the code.
3. **`AGENTS.md`'s gate table is not checked against the Makefile** — I looked, and no test reads the Makefile (the brand test only reads it for the mark). The brief says the gate's list is checked against the Makefile; today that wire does not exist. The new target is in both places, by hand. A wire for that is a next round's work, not this one's.
4. **The orchestrator's labels are untouched.** The README's state table, its rulings table, `plan-automation.md` § 10 and § 11 and `docs/05` are as the orchestrator left them; the only line I touched in `work/tui/README.md` is **T9's own row**, which AC-03 asks for.
5. **The invariants held**: `internal/core` did not change, `core.DomainVersion` stays **7**, the count did not move (46 commands, 40 verbs in 8 groups), `internal/settings` stays stdlib-only (0 external modules), and the record's shape did not move.
6. **The check found the drift the orchestrator had already fixed** — it is green on the tree as it stands, and it fails on every one of the five places when its number is changed.
