Runner: T10 fix 2 done. Commits: 5 — `7edc649` (AC-01), `9003efb` (AC-02), `7d58dc8` (AC-03), `7597466` (AC-04), `99412f6` (AC-05).

**The label.** It is sanitised **where it enters** — in `readCrew`, before it is anywhere near the screen (`internal/ui/crew.go`). The label is agent-written text: the console writes the first one and the orchestrator writes it from then on, so a newline in it put a line of the agent's text where the crew's rows are a table, and an escape sequence reached the terminal. The witness is on **this** tab — `TestTheCrewLabelIsSanitisedOnTheSousChefTab` injects a label with a newline and an escape, renders the Sous-Chef tab, and holds that no control character reaches the render, that the label's escape is not in the raw render, and that the injected text does not become a line of its own. **The witness bites**: taking the `sanitize` out turns it red.

**One format.** The wire is the **skill's own example**, read from `internal/cli/skills/orchestrator-role/SKILL.md` and fed to the console's parser (`internal/ui/crew_wire_test.go`). Every roster's line the skill shows is extracted, the placeholders are filled in, and each one must parse: a roster's line, a job of the crew's, and — for the orchestrator's — a round and an interaction count the console reads. **Your probe now fails**: the skill drifts to *"round #<N> · <M> interactions"*, the skill's own witness is updated, **`internal/cli` stays green and `internal/ui` goes red** — *the console reads no round from the skill's example*. I ran that probe by hand: two green packages was the finding, one green package and one red is the fix.

**The stamp.** Three shapes, and the heading reads as two facts with the vital signs:

- a progress reported → `THE STAMP · round 2 of 3 · 6 of 12 interactions · plan 7 · build 7 · they agree`
- **no** progress → `THE STAMP · plan 7 · build 7 · they agree` — the separator is there whether a progress was reported or not (it used to read *"THE STAMPplan 7"*).
- a progress **over** the cap → `THE STAMP · round 7 · 39 interactions · the budget is spent · plan 7 · …` — the numbers as reported, never *"7 of 3"*.

And two more: the flag fires when the report **reaches** the cap (`round 3 of 3 · 12 of 12 interactions · the budget is spent`) — it used to wait for the next round to go over; and a **zeroed** report (`round 0 · 0 interactions`) is a report, not an absence — the console now carries a **reported** flag beside each number, so a zero reads as a zero and an absence stays an absence.

**The pane's own tab id.** The startup task and `orchestrator-role` point at **`herdr pane current`**, which carries the pane and the tab that pane is in — not `herdr pane list`, where the orchestrator had to find itself. The lanes' tab ids are still in `herdr tab list`, and the skill's mechanics table says so.

**DomainVersion:** 7. The tab count stays **8**; `internal/core` does not change; the budget is **28 of 30**.

**The gate:** `make check` — **pass** (internal/ui 195.8s). `make check-race` — **pass** (internal/ui **523.0s**, internal/core 107.3s, internal/cli 11.2s).

Notes for the reviewer:

1. **The wire lives in `internal/ui` and reads the skill's file by path** (`../cli/skills/orchestrator-role/SKILL.md`). The console does not import the CLI — the path *is* the wire, and a skill that moves breaks the test loudly rather than silently. I took your first option: a test that takes the skill's own example and feeds it to the console's parser.
2. **The zeroed report's shape is mine**: `round 0 of 3 · 0 of 12 interactions` — the orchestrator's numbers, echoed against the cap. A round is 1-based, so a *reported* zero round means "no round yet"; the console prints what was reported and does not read it as an absence.
3. **The over-budget wording is unchanged** — *"the budget is spent"* — and it now fires **at** the cap, not one round later.
4. **The notes I did not fix, and where they belong.** The render path still waits on a process (`render.go`'s `snapshot()` → `readCrew()`, two `herdr` calls) and the poll still spends two process spawns a second while the tab is visible. Both belong to the **seeder/render slice that is next**: the render path's wait is the same seeder cost the suite pays 219 times over, and the poll's two spawns are the price of asking a runtime on every tick. I did not start the seeder — it is the next slice, as ruled.
5. **The gate's headroom did not move**: internal/ui is **523.0s** of go test's 600s default under `-race -shuffle=on -count=2` — **87%**, the same number as your measurement, because this cycle's five tests cost about what the ones they replace cost. The seeder slice is not optional now: the next slice of this size lands at ~90%, and the one after it breaks the gate.
6. **The five commits are one per criterion and each is green** — I checked `make lint` and the focused tests on each, and nothing else in the tree moved.
