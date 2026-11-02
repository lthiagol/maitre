Runner: T12 done. Commits: 4 — `d197c41` (AC-01), `b0040e0` (AC-02), `ade00a2` (AC-03), `3bfdb85` (the read says what it holds). The row's own commit is the fifth.

**The snapshot:** it asks nothing. `internal/ui/render.go`'s `snapshot()` no longer reads the crew. The runtime is asked by a **command** — the poll, and the switch that brings the tab on screen — and the answer arrives as a message, so a frame draws what the console already knows. The witness is **process level**: a fake `herdr` sits on the PATH and records every call, so the count is the claim. `TestASnapshotSpawnsNothing` holds the snapshot at **zero** calls (it draws the Sous-Chef tab twice, and the fake is never run), and `TestTheSnapshotDoesNotAskTheRuntime` holds the seam — no crew, no error, never asked.

**One spawn:** the read asks the runtime **once**. The seam is one call — `crew() (herdrCrew, error)` — and the runtime answers with `herdr api snapshot`, which carries the **agents** and the **tabs** in one document: the state from the agent's `agent_status`, the job from the tab's label, matched by `tab_id`. `TestTheCrewReadCostsOneSpawn` holds the read at **one** fork, asserts the call is `api snapshot`, and asserts the one answer carries both questions (the member's state *and* its job).

**The live console:** the poll asks the runtime while the tab is visible, and the tab shows **what it last heard** — the answer is held on the model and the frame reads it. `TestThePollAsksTheRuntimeInOneProcess` holds the poll at one fork; `TestTheTabShowsWhatItLastHeard` holds that the five frames after the answer fork **nothing**; `TestTheClickOnTheTabReachesTheRuntime` holds F-140's witness at the process level — the mouse's switch reaches the runtime once, the tab does not say *the automation is not running* while it is, and the footer never offers to start a second crew. The filter line's switch holds too (`TestTheFilterLineSwitchAsksTheRuntime`, unchanged).

**The goldens:** **not one moved.** `make goldens-check` regenerated all five and the tree is clean. `TestTheGoldensRenderWithHerdrOffThePath` renders all five with **herdr off the PATH** and compares them, byte for byte, with the files in the tree.

**The poll:** its **frequency is unchanged** — `PollInterval` is still a second. What changed is what the poll asks for, not how often. The poll's frequency stays with the poll.

**DomainVersion:** 7 (`internal/core/lifecycle.go:225`). `internal/core` did not change, `docs/03` § 8 did not move, and the budget is unchanged — 0 / 0 / 0 / 28.

**The gate:** `make check` green, `make check-race` green (`-race -shuffle=on -count=2`).

**Notes:**

- **The tab has a third state, and it is the one judgement call in the slice.** A console that has **not** asked the runtime says *"the crew has not been asked"* and offers **no move**. It does not say *the automation is not running*: that is a guess about a crew it has not heard about, and a guess there is a move that can start a second crew — the hazard F-140 is about. The state shows in the static render (`maitre --render --all`) and for the frame between a tab switch and the answer. The three states are: asked and running (the crew's rows) · asked and gone (*the automation is not running*) · **not asked** (*the crew has not been asked*) · asked and failed (*herdr is not here*).
- **The tests that fake the crew now say so.** `withCrew` gives the console an answer (the crew's rows *and* the fact that it has one), and the two tests that moved a crew — start and stop — **ask the runtime first**, because a move offered on a guess is a move that can start a second crew. The claims under test did not move: the stamp, the progress, the vital signs and the blocked agent read as they did.
- **`space` in the not-asked state asks, and does not move.** The console does not know whether a crew runs, so it asks — in a command, so the frame does not wait — and the move comes back with the answer. In a `--keys` script the switch's own ask has already landed, so `space` starts the crew as it always did.
- **`docs/05-human-interface.md`** named the two commands the read used (`herdr agent list`, `herdr tab list`). It now names the one it uses, carries *one read, one process* and *a frame never waits on a process*, and lists the new state in the failures table.
- **`work/tui/README.md`'s T12 row** names the four criterion commits, and the row's own commit is not among them.
- The suite's own plan (`kitchen/`), `docs/audit.md`, the seeder and the console's behaviour the chef sees did not move.
