Runner: T11 done. Commits: 5 — `825bb82` (AC-01), `90774a4` (AC-02), `f090206` (AC-03), `2fd2c13` (AC-04), `e9df136` (AC-05). The receipt is the sixth.

**The store:** `internal/ui/seed_test.go`. `TestMain` makes a directory the suite owns (`os.MkdirTemp("", "maitre-ui-templates")`) and seeds the templates into it — the ordinary plan at `<dir>/demo/kitchen`, the pathological one at `<dir>/stress/demo/kitchen`. The seed happens **once per process**, in `TestMain`, not per test. A test draws a plan with `seededPlan(t)` / `stressPlan(t)`: it copies the template into `<t.TempDir()>/demo/kitchen` and opens it. The nine seeding sites in the package (the twelve `core.Seed` calls behind them — `newModel`, `goldenStore`, `identityStore`, `settingsStore`, `stressStore` and the four direct ones) now draw a copy.

**The copy:** byte for byte. `copyTree` walks the template and writes each file's bytes as they are — nothing is parsed, rewritten or re-serialised. The witness is `TestTheCopyIsTheTemplate`: two copies made at different times hold the same content, a copy holds the same content as the template, and the template holds exactly the seeder's fifteen files, named in the test.

**The names:** the copy lands in a directory named `demo`, so the plan's name is `demo` — the name every golden holds. The template's `project.json` names **nothing**, so a copy's name follows the directory it lands in: `demo` for the ordinary copies, and the odd names for the two tests that need an odd path (`we\nird`, 100 `p`s), exactly as a plan seeded there said before. The witness is `TestTheCopyIsNamedAfterItsDirectory`; it was probed and fails when the template's own name comes back. Every golden is unchanged — `make goldens-check` regenerated all five and the tree is clean.

**The stress plan:** a **second** template. `core.SeedStress` is seeded on its own at `<dir>/stress/demo/kitchen`; it is not a copy of the first. `TestTheStressPlanIsItsOwnTemplate` proves it is a different plan and holds the `STRESS:` recipes.

**Count two:** the seed count does not move with the runs — `seeds == 2`, one per template — while the copies do. `TestTheCopiesGrowAndTheSeedDoesNot` draws its own copies, so it does not depend on where it lands in the order under `-shuffle=on`. A write in one copy reaches neither another copy nor the template (`TestTheCopiesDoNotSeeEachOther`) — the template is shared by every test in the run, so a write that reached it would show up in the next test to draw a plan. The cleanup is `TestMain`'s: it removes the template directory, and the suite **fails** if the directory survived.

**The cost:** measured on this machine, the same way both times — the test binary, run from the package.

| | before | after |
|---|---|---|
| `internal/ui`, plain | 182.2s | 10.5s |
| `internal/ui`, `-race -shuffle=on -count=2` | 535.2s | 38.3s |

go test's limit is 600s per package: the package was at **89%** and is now at **6%** — a margin of 14×, where the slice asked for 3×. The brief's baseline was 163.8s and 561.3s; this machine measures the plain run a little slower and the race run a little faster than the brief's figures.

**The render path:** not started, as instructed. Its wait is `internal/ui/render.go:41` — `snapshot()` calls `m.readCrew()`; `readCrew` (`internal/ui/crew.go:54`) asks the runtime for the agents and the tabs, and each ask is an `exec.Command("herdr", …)` (`internal/ui/herdr.go:76`). Two spawns a render. That is the render path's cost, and it belongs to whoever next touches the render path.

**DomainVersion:** 7 (`internal/core/lifecycle.go:225`). `internal/core` did **not** change — the copy helper is the ui package's own test code, so the domain's declared-helper wire did not move, `docs/03` § 8 did not move, and the budget is unchanged (0 / 0 / 0 / 28).

**The gate:** `make check` green, `make check-race` green. The slow gate now takes 1m52s in total, with `internal/ui` at 47.8s of it.

**Notes:**

- **One deliberate difference from the seeder's output:** the template's `project.json` carries an empty name — one field, written the way the store writes it. The seeder writes the name it took from the directory it seeded into, so a *named* template would force the name `demo` onto every copy, and the two identity tests that need an odd path would lose the name they were about. A plan that names nothing answers with its directory — the rule the store already follows — so a copy is indistinguishable from a plan seeded in that directory. This is the one judgement call in the slice; the reviewer may want it the other way, and it is a two-line change (`unnameTheTemplate`).
- **The witnesses were probed, not assumed.** The naming test fails when the template's name comes back; the "no test seeds its own plan" wire fails when a test file calls the seeder. The first draft of the copy-count witness depended on other tests having run first and would have failed under `-shuffle=on`; it now draws its own copies.
- `work/tui/README.md`'s T11 row names the four criterion commits, and the row's own commit is not among them. The "next work" row said `internal/ui` was at 88% of the 600s default; it now says 6% and that the next slice is not briefed yet — the orchestrator's call, not mine.
- The suite's own plan (`kitchen/`), `docs/audit.md` and the console's behaviour were not touched. No golden moved.
