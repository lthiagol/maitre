Reviewer: T11 done. Verdict: **changes-needed**. Filed **1** blocker — **F-145** (medium): the template's project file is rewritten through a second copy of the plan's schema, so the suite's fixture can drift from the model silently.

**The store.** `internal/ui/seed_test.go` — `TestMain` makes a directory the suite owns (`os.MkdirTemp("", "maitre-ui-templates")`) and seeds **two** templates into it: the ordinary plan at `<dir>/demo/kitchen`, the pathological one at `<dir>/stress/demo/kitchen`. The seed happens **once per process**, in `TestMain` — never per test. A test draws a plan with `seededPlan(t)` / `stressPlan(t)`: it copies the template into `<t.TempDir()>/demo/kitchen` and opens it. Nine seeding sites in the package now draw a copy.

**The copy.** Byte for byte. `copyTree` walks the template and writes each file's bytes as they are — nothing is parsed, rewritten or re-serialised on the way. The witness is `TestTheCopyIsTheTemplate`: two copies made at different times hold the same content, a copy holds the same content as the template, and the template holds exactly the seeder's **fifteen** files, named in the test.

**The names.** The copy lands in a directory named `demo`, so the plan names itself `demo` — the name every golden holds. **Every golden is unchanged**: `git diff -- internal/ui/testdata/` is empty over the slice, and `make goldens-check` regenerated all five and the tree stayed clean.

**The stress plan.** A **second** template, seeded on its own with `core.SeedStress` — not a copy of the first. `TestTheStressPlanIsItsOwnTemplate` proves it is a different plan and holds the `STRESS:` recipes.

**Count two.** The seed count does not move with the runs — `seeds == 2`, one per template — while the copies do (`TestTheCopiesGrowAndTheSeedDoesNot` draws its own copies, so it does not depend on where it lands under `-shuffle=on`). A write in one copy reaches neither another copy nor the template (`TestTheCopiesDoNotSeeEachOther`). The cleanup is `TestMain`'s: it removes the template directory **and fails the suite if it survived** — proved by mutation (*"the template directory survived: …"*).

**The cost.** Measured by me, on this machine, the same way both times — the test binary, run from the package.

| | before (the brief) | before (the receipt) | **after (mine)** |
|---|---|---|---|
| `internal/ui`, plain | 163.8s | 182.2s | **10.19s** |
| `internal/ui`, `-race -shuffle=on -count=2` | 561.3s | 535.2s | **37.9s** |

go test's limit is 600s per package: **89% → 6%**, a **14×** margin where the slice asked for 3×. The runner's figures (**10.5s** / **38.3s**) stand — mine are within a second of them.

**The judgement call — the template names nothing. It is right, and here is why.** The name does not come from the file: `core.projectName` (`internal/core/store.go:141`) answers an empty name with the **project's directory**, and `TestProjectNameHasOneAnswer` already holds that rule in the domain. The template naming nothing therefore uses a rule the store **already owns** instead of adding a second one, and a copy in `demo/` says `demo` — the name every golden holds — while a copy in `we\nird` or in a hundred `p`s says what it said before. The alternative (a named template) would force one name on every copy and lose the two odd-path tests. **What a reader trips over:** nothing reads the template's own name — the template is read only by `copyTree` and by the witness — so the emptiness breaks no reader. The receipt's claim *"the only difference from the seeder's output is the name itself"* is **true today** (proved by command: the template's project file is byte-identical to a freshly seeded plan's, apart from the name and the created stamp). But **nothing holds it**, and the mechanism that makes it true is a **second copy of the schema** — see F-145.

**F-145.** `unnameTheTemplate` unmarshals `project.json` into its **own** struct (`Name`, `Created`, `Charter`, `Settings`) and re-marshals it, so it holds **its own copy** of `core.projectFile`. Any field the domain adds to the project file is therefore **dropped from the template**, and the suite's plan stops being the plan the seeder writes. Proved by command, and reverted: adding an `owner` field to `core.projectFile` → **`internal/ui` green** (the fixture kept the old shape and said nothing) while `internal/core` goes red in `TestDomainShapeIsVersioned`; and a probe on the rewrite showed the field **gone** from the file afterwards. The byte-for-byte witness cannot see it — it compares a copy to the template, and both come from the same drifted source — and neither can the goldens. Nothing is wrong today; this is a hole in the rule *"the fixture cannot drift from the model"* (AGENTS.md, slice 3b). **The remedy is one line**: rewrite the name in the **raw bytes** — `strings.Replace(raw, `"name": "demo"`, `"name": ""`, 1)`, the way `internal/core/store_test.go:44` already does — so the template holds whatever the seeder wrote and no field can be dropped; and the replace is already guarded, because a template that kept its name turns `TestTheCopyIsNamedAfterItsDirectory` red.

**The render path.** The wait belongs to the **render path**, and the brief's placement is right: `crew.go` states *"a frame must never wait on a process"*, and `render.go`'s `snapshot()` → `readCrew()` → `herdr.go`'s `exec.Command("herdr", …)` breaks exactly that rule — two spawns a render. The runner **naming** it in the receipt rather than leaving it quiet is **right**: a note that is not written down is lost, and the brief forbade starting it. It was already a ledger note from the T10 cycle, so nothing was lost twice. **The brief for it goes to the chef first**, and it must achieve: **(1)** a snapshot costs **no** process spawn — the `--render` and `--keys` paths must not block on `herdr`; **(2)** the goldens do **not** depend on whether `herdr` is on the `PATH` (true today — no golden renders the Sous-Chef tab — and worth holding before one does); **(3)** the live console keeps asking the runtime **asynchronously**, unchanged. **The poll's two spawns a second does *not* belong with it.** Different owner (the poll), different symptom (the *frequency* of the ask, not a frame's *blocking*), and it is **deliberate**: F-140's witness `TestThePollAsksTheRuntimeWhileTheTabIsVisible` requires the poll to ask while the Sous-Chef tab is visible **even when the plan has not changed**. It is already a ledger note, measured at ~5ms a call — ~10ms a second. What the two share is the **cause** — the crew read is two spawns — so the render path's brief should also say the crew read must cost **one** spawn, not two; the poll's *frequency* stays with the poll.

**Mutations — six from the brief plus one of my own, every one reverted.**

| # | The mutation | What went red |
|--:|---|---|
| 1 | a test seeds its **own** plan | the wire: `TestNoTestSeedsItsOwnPlan` — *"zz_mutation_test.go seeds its own plan: core.Seed( — a test takes a copy, not a seed"* |
| 2 | the copy is **not** byte-for-byte (one byte appended) | **2**: `TestTheCopyIsTheTemplate` (*"the copy is not the template … differ in activity.json"*) and `TestTheCopiesDoNotSeeEachOther`. The goldens stayed green — the byte witness is the load-bearing one |
| 3 | `TestMain` does **not** remove the template directory | the suite: `TestTheTemplatesAreRemoved` and `TestMain` — *"the template directory survived"* — and the directory was indeed left behind |
| 4 | the template **names** the plan | **1**: `TestTheCopyIsNamedAfterItsDirectory` — *"the plan's name is \"demo\", want the directory's name \"the-other-one\""* |
| 5 | the stress plan **shares** the first template | **19**, including `TestTheStressPlanIsItsOwnTemplate` and four goldens |
| 6 | **mine**: a per-test seed inside the exempted file | the counter, not the wire: *"the suite seeded 201 times, want 2"* — and the suite went to **180.6s** |
| 7 | **mine**: the domain's project file gains a field | **nothing in `internal/ui`** — that is F-145 |

All reverted; `git status --short` is empty except the ledger commit and this receipt.

**The questions.** `work/questions.md` — every row is struck (Q8–Q13, Q25, Q28, Q29). **Nothing answered, nothing missed, nothing newly raised.** This slice changed no behaviour and no vocabulary.

**DomainVersion:** 7 (`internal/core/lifecycle.go:225`). `internal/core` did **not** change — `git diff 337b7fe..HEAD -- internal/core/` is empty — and the copy helper stayed in the ui package's own test code, so the domain's declared-helper wire and `docs/03` § 8 did not move. The budget is unchanged: **0 / 0 / 0 / 28**.

**The gate.** `make check` **green** and `make check-race` **green**, confirmed by me on the committed tree (`internal/ui` **43.3s** in the slow gate). The canaries too: `maitre-cli demo` renders, and the repo's own `kitchen/` still renders.

**LOW notes.**
- `TestTheTemplatesAreRemoved`'s failure message prints the **nil** stat error — *"the template directory survived the cleanup: `<nil>`"*. The reason is `%v` of a nil error; a reader gets `<nil>` as the cause.
- The template directory lives in the user's `TMPDIR` and `TestMain` removes it — but a **killed** test binary (SIGKILL) leaves a 15-file plan behind. The brief's *"not a user's tmp"* is met by the removal and the survival check; the killed-run case is the leftover.
- The wire matches the **literal** `core.Seed(` / `core.SeedStress(`; an aliased call would slip past it. The **counter** covers the cost anyway — proved by mutation 6 (*"the suite seeded 201 times"*) — so the pair holds the rule the wire alone does not.
- The `seeds` / `copies` counters are plain ints, and no test in `internal/ui` calls `t.Parallel()` — so the race gate has nothing to catch there. Worth knowing before a parallel test lands.

**What I did not touch:** the implementation (mutated and reverted), `work/` (except this receipt), the runner (not prompted), and the render path (not this slice's, and not mine).
