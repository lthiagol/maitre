Reviewer: T10 cycle 4 done. Verdict: ok. Filed 0 blockers (none).

**The third copy.** The startup task's copy of the format is now **one constant** — `rosterExample` in `internal/ui/souschef.go` — and the rendered task's text is **unchanged**: `herdr tab rename <TAB_ID> "maitre: orchestrator · round <N> · <M> interactions"`. I checked it by command, and the task carries the constant verbatim.

**It is wired to the parser, twice and independently.** `internal/ui/crew_wire_test.go` extracts **the startup task's own example** from the rendered task and feeds it to the console's `parseRoster` and `parseProgress` (*`TestTheStartupTaskExampleIsTheFormatTheConsoleReads`*); a **second** wire holds that the copy is **one of the skill's own lines**, not a variant of it (*`TestTheStartupTaskCopiesTheSkillIsTheSpecification`*).

**My probe now turns TWO tests red — reproduced by hand.** The separator drifts from *" · "* to *" / "*, every word intact, and **two** tests go red: *"the roster's line names the job 'orchestrator / round 2 / 6 interactions', which is not one of the crew's"* and *"the startup task's copy is not the skill's own line — the skill's table is the specification"*. Nothing red became two red, and the two wires are **independent**: with the parse wire broken, only **one** goes red (the equality wire still holds it), and the other way round.

**The specification** is said where a reader will find it — the **skill's table** (*"This table is the specification. The console's parser and the startup task's example follow it — when the format changes, this table changes first, and the wire test holds the other two to it"*), held by the **skill's own witness** (*`TestTheSkillSaysTheOrchestratorKeepsTheLabelsCurrent`* now asserts that sentence); the **console's constant**'s doc comment; and the **wire's** comment. The startup task's doc says its copy is **an example**, not the specification.

**The flag** names the allowance that ran out — **"the rounds are spent"** · **"the interactions are spent"** · **"the rounds and the interactions are spent"**. Read by command: *"round 3 of 3 · 6 of 12 interactions · the rounds are spent"*; *"round 2 of 3 · 12 of 12 interactions · the interactions are spent"*; *"round 3 of 3 · 12 of 12 interactions · the rounds and the interactions are spent"*; and no flag at all when neither has run out. A reader never has to guess which allowance ran out.

**The four notes.**
1. **The copy is a constant, not a retyped literal — holds.** The console's side of the format lives in **one** place: the constant and the parser, joined by the wire to the skill's table. I checked that the rendered text did not move.
2. **The equality test is the "one answer" check — holds, and it is independent.** The task's copy must be *one of the skill's own lines*; my mutation shows it catches what the parse wire does not, and vice versa. One further thing I probed and it holds: **the copy cannot become a fixed round and six** — the constant adopting the skill's *bash-block* line (*"maitre: orchestrator · round 2 · 6 interactions"*, a filled form that would fix the progress at 2 · 6 forever) turns the **substring witness red**, because it holds the placeholders. So the three wires together hold the copy to the **specification**, not merely to "a line the skill shows".
3. **The flag's wording is the runner's; the requirement was the reader's — holds.** The requirement was that a reader must not have to guess which allowance ran out, and the three wordings deliver it.
4. **The seeder and the render path are untouched, as ruled — holds.** The headroom moved from 523.0s to **529.1s**; I measured **526.4s** — **88%** of go test's 600s default under `-race -shuffle=on -count=2`. The shared seeded store is the next slice and the only thing that buys the margin.

**The next slice — the shared seeded store. I would hold the same line**, with the conditions from the last cycle carried forward, because this cycle did not move the clock.

| | |
|:--:|---|
| **it must achieve** | the seeder runs **once** per test **process** — a `sync.Once` or a `TestMain` — so it survives `-count=2` without rebuilding; every test gets a **copy**; and `internal/ui` falls from **526s** to **well under 200s** under the race gate — a margin of at least **3×**, not 1.13×. The arithmetic: the suite builds a plan in **219** places (`newModel` 65 · `recipesModel` 88 · `sousModel` 34 · `loadModel` 32), `core.Seed` costs **0.608s** and writes **15 files**, so ~**133s** of the ~**195s** non-race run is the seeder rebuilding the *same* plan; a copy costs **0.007s**, so the store makes it 0.6s + 219 × 0.007s ≈ **2.1s** |
| **isolation** | each test its **own** copy, or the writers see each other under `-shuffle=on -count=2` |
| **the goldens** | the project's name comes from the **parent directory's** name, so the copy lands in a directory named **`demo`**; and the copy must be **byte-for-byte**, because the goldens pin the seeded plan's *content* |
| **the stress plan** | a **second** template (`core.SeedStress`), not the same one |
| **the cleanup** | `TestMain`, not a user's tmp |
| **timestamps** | no test asserts a fresh one — I checked twice, and none does, so a fixed `created` stamp is safe and makes the goldens *more* deterministic |

**It must land before any slice that adds tests.** At 88% the next slice of this size lands at ~95%, and the one after it breaks the gate. And the **render path's wait** — `render.go`'s `snapshot()` → `readCrew()`, two `herdr` calls — is *not* the seeder's cost and should be named where it is fixed; the **poll's** two spawns a second are the price of asking a runtime on every tick (5ms a call, measured) and belong to whoever next touches the poll.

**Mutations — four from the brief plus one of my own, each proved by command, each reverted; `git status --short` is empty except my ledger commit.**

| | the mutation | what went red |
|--:|---|---|
| 1 | the startup task's **separator** drifts, every word intact | **2** — `TestTheStartupTaskExampleIsTheFormatTheConsoleReads`, `TestTheStartupTaskCopiesTheSkillIsTheSpecification` |
| 2 | the wire reads a **different** example than the task's own | **1** — `TestTheStartupTaskCopiesTheSkillIsTheSpecification` (the parse wire is gone; the equality wire still holds) |
| 3 | the equality test compares against a **variant** of the skill's line | **2** — `TestTheStartupTaskTellsTheOrchestratorToKeepItsLabelCurrent`, `TestTheStartupTaskCopiesTheSkillIsTheSpecification` |
| 4 | the flag says *"the budget is spent"* again | **4** — `TestTheStampPrintsNoNumberItCannotCompute`, `TestTheStampReadsInAllThreeShapes`, `TestTheBudgetIsSpentWhenTheReportReachesTheCap`, `TestTheFlagNamesTheAllowanceThatRanOut` |
| 5 | **mine** — the constant adopts the skill's **bash-block** line, with its numbers fixed | **1** — `TestTheStartupTaskTellsTheOrchestratorToKeepItsLabelCurrent` (the substring witness holds the placeholders, so the copy cannot become a fixed round and six) |

**The questions.** `work/questions.md` was **not touched** by the cycle, and every row in it is struck — Q11–Q13 and Q25 closed, Q16–Q30 ruled. **Nothing was newly raised.**

**DomainVersion:** 7. `internal/core` did not change; the tab count is **8**; the budget is **28 of 30**. The T10 row names **18** commits and says *"18 commits"*; **every one is an ancestor of HEAD** (checked one by one with `git merge-base --is-ancestor`), and the row's own commit (`d0f0f15`) is **not** among them.

**The gate:** `make check` — **pass**. `make check-race` — **pass** (`internal/ui` 526.376s; `internal/core` 105.9s; `internal/cli` 11.3s).

**LOW notes** (in the ledger's notes table, not findings): the **skill shows the roster's line twice** — the table's placeholder form (the specification) and the bash block's filled illustration — and the filled one is the more tempting to copy **verbatim**, which would hand the orchestrator a fixed round and six; the skill could mark the numbers as the orchestrator's own. And the two carried forward from the earlier cycles — the **render path** waits on a process, and the **poll** spends two process spawns a second while the tab is visible.

**The slice is closed.** T10's five findings — F-140, F-141, F-142, F-143, F-144 — are all closed, and the ledger now carries the row that closes the slice. The next work is the **shared seeded store**, and it must land before the next console slice.
