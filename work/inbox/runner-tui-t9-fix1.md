Runner: T9 fix 1 done. Commits: 3 — `b40eedf` (AC-01, the witness), `97f1de9` (AC-02, the doc), `8a6bf57` (AC-03, the row's own).

**The atomicity.** The write now has **two seams** in `internal/settings` — `writeTemp` (the bytes into the temporary file) and `renameTemp` (the rename over the old one) — so a fault can be injected at **either** stage of the atomic write, and two witnesses hold the property the brief names: **the file that was there before is whole**.

| Witness | The fault | What it holds |
|---|---|---|
| `TestAHalfWrittenFileCannotLoseTheChefsChoices` | the bytes go **half way** into the temporary file and the write fails | the old file is byte for byte whole, and the half-written temporary file is gone |
| `TestARenameThatFailsCannotLoseTheChefsChoices` | the bytes are staged and the **rename** fails | the old file is byte for byte whole, and the staged bytes went with the fault |

The witnesses are **property** witnesses, not mechanism watchers: the property is *"the chef's file survives a write that fails"*. A write that replaces the old file in place has neither stage, so the fault cannot be injected at all — the write succeeds, and the witness says so on the line *"the write should have failed — the fault was injected into it"*.

**The mutation result, proved by command and reverted:**

| Mutation | What went red |
|---|---|
| the temporary file and the rename replaced by a straight `os.WriteFile` (the reviewer's) | **both** witnesses red; `internal/cli` and `internal/ui` stay **green** |
| the temporary file kept, the rename replaced by writing over the target | the **rename** witness red |

The second mutation is why there are two seams: with a fault only at the write stage, that mutation left the suite green — so the first version of this fix was not enough, and the commits were rebuilt. The implementation is back to the temp file and the rename; `git status` is clean.

**The doc.** `docs/05`'s Settings prose now says what the code does: the theme, the order and show-served are persisted in `kitchen/settings.json` — the **tool's** file, in the plan's folder, **not a plan record**; the plan's format does not move, `core.DomainVersion` stays **7**, and the plan's readers ignore the file; the console is the only writer, so a theme or an order it does not know is refused by name and the console says so. My commit makes the file's **own summary table agree with that prose**: the Theme row said *"persisted in the plan's `settings.json`"* — which reads as the plan owning the file — and now says *"persisted in `kitchen/settings.json` — the **tool's** file, in the plan's folder"*, beside the Settings row that says the same. One more stale label in the same file: the tab table marked the Sous-Chef tab *"(T9 — planned)"*, and it is **T10** — `work/tui/README.md`'s table, the T9 brief and `plan-automation.md` § 14 all say so. Fixed, and I say so here because it is outside the Settings prose.

**The orchestrator's corrections.** Found **already fixed** in `babf240`, not redone: `plan-automation.md` § 15 now states that the **version fence wins** (a file at a version this build does not know is refused by name, before its keys are read; unknown keys are ignored only within the version it knows); § 15's CLI verb is **`settings show` only**, with the reason (the console is the only writer); and `work/questions.md`'s Q25 is **struck**. `docs/05`'s prose had also already been rewritten in that commit — I did not redo it; my commit makes the *table* agree with the prose it had written.

DomainVersion: **7**. The count did not move — still **46 commands, 40 verbs in 8 groups, plus 6 workspace commands**. `internal/core` did not change (`git diff babf240..HEAD -- internal/core` is empty), `internal/settings` stays **stdlib only** (0 external modules), and the record's shape did not move.

The gate: `make check` — pass; `make check-race` — pass.

Notes for the reviewer:

1. **The seams are the slice's only new production surface** — two package-level function variables in `internal/settings`, restored by `defer` in each witness, so `-shuffle=on` and `-count=2` are safe and the race gate is green.
2. **The two witnesses are not twins.** The first injects the fault in the *bytes*, the second in the *last step*. Between them they cover a failure at either stage of the write, and each mutation above turns a different one red.
3. **The existing `TestTheWriteIsAtomic` stayed** — it proves the weaker things (no file left behind; a write into a folder that is not there leaves the old file whole). The new witnesses are the ones that bite.
4. **The commits were rebuilt.** The first AC-01 attempt (`f482661`, since replaced) had one seam, and the second mutation above left it green — so the three commits were re-made over the same file sets: `b40eedf`, `97f1de9`, `8a6bf57`. The row in `work/tui/README.md` names the final ids.
5. **Nothing outside the slice was touched** — `docs/05-human-interface.md`, `internal/settings/` and `work/tui/README.md` are the only files the fix cycle changed.
