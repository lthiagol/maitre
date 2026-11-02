Reviewer: T11 cycle 2 done. Verdict: **changes-needed**. Filed **1** blocker — **F-146** (medium): F-145's guarantee is true by construction and held by **no witness**. **F-145 itself is closed** — the rewrite is textual, the schema copy is gone, and the probe now fails in both packages.

**The project file.** The name is rewritten **textually**: `blankTheName` reads the seeder's **own bytes** and empties the name's value where it stands — a regexp on `"name"\s*:\s*"` and the value it holds — and **nothing is parsed and nothing is re-serialised**. The copy of the plan's schema is **gone**. Proved by command: the template's project file is byte-identical to a freshly seeded plan's apart from the name's value — every other byte the domain's.

**The probe now fails, in both packages.** Adding an `owner` field to `core.projectFile`:

| | |
|---|---|
| **`internal/ui`** | red — `TestTheTemplateCarriesTheDomainFields`: *got: [charter created name owner settings], want: [charter created name settings]* |
| **`internal/core`** | red — `TestDomainShapeIsVersioned` |

Both reverted; `git diff` over `internal/core/` is empty.

**The counters.** `seeds` and `copies` are `atomic.Int64`, and `TestTheCountersAreSafeUnderAParallelTest` runs two subtests **in parallel**, each drawing two plans, with the parent's cleanup checking the count once they are done. Proved by mutation: **plain ints → the race detector says `DATA RACE`** (two races, in `copyTemplate` from the parallel subtests); the atomic counters say **nothing** under `-race -count=2`. The cleanup's failure message says **what** happened and **where** the leftover is — *the cleanup did not remove the template directory: `<path>` is still there* — and it no longer prints the nil stat error (proved by mutation: with the old message the reader gets `<nil>` as the cause).

**The wire.** It matches the **call**, not one spelling: `\bSeed(Stress)?\s*\(`. Proved by command — it catches the qualified call, an **aliased** one (`corealias.Seed(`), one with a **space** before the parenthesis, and a **newline** (`core.Seed\n(`); and it is not too wide (`seededPlan(`, `unseeded(`, `theSeeder(`, `core.Seeded(` are not calls). The **count** agrees with the counter: a seeder call smuggled into the template's own file → *the sources call the seeder 3 times, but the suite seeded 2*. The counter still covers the cost.

**F-146 — the guarantee is true by construction and held by no witness.** `TestTheTemplateCarriesTheDomainFields` compares the template's project file against a **hand-written list of names**, so the two go stale **together** and agree. The mechanism's fidelity is what makes the domain's move visible at all — the domain adds a field → the textual rewrite keeps it → the template carries a name the list does not → red. **But if the mechanism ever stops keeping it, the wire goes quiet.** *Proved by command, and reverted*: the rewrite **parsed and re-serialised again** — F-145's own mechanism — **and** the domain's project file carrying a new field → **`internal/ui` green**. The fixture silently dropped the field again and the wire said nothing. `internal/core` did go red, so the domain is fenced — but the ui fixture is wrong and the ui suite is silent about it. The same hole shows a second way, and it is **live today**: the rewrite blanks **every** `"name"` key it finds, **at any depth**, so a nested `"name"` in the domain's project file is blanked too — *proved*: with the domain's `settings` carrying `{"theme":{"name":"dark"}}`, the template held `"name": ""` while a seeded plan held `"dark"`, and **both packages were green**. The code's own comment claims *"every other byte is left as the domain wrote it"* — false in both cases, and a doc that contradicts the code is a bug. **The remedy is a witness, and it costs one `Init()` and no seed**: `core.Seed` calls `store.Init()` (`internal/core/seed.go:33`), so the domain's **own** project file is available without seeding — `core.Open(dir).Init()` writes it. Parse the template's file and that one, drop `name` and `created` from both, and compare: the witness holds the field set and the values (nested included), and it needs **no** list of names to do it. *Proved by command*: the witness is green today and goes red under **both** cases the list cannot see. The list stays, for the tripwire it is — a domain that moves must turn the ui suite red.

**The two notes.**

| | the note | my answer |
|:--:|---|---|
| 1 | the wire is a list of field names in the test, not a copy of the schema in the mechanism — does it re-create the drift? | **It holds — no.** It is a copy in the **witness**, not in the mechanism: it changes nothing, so it cannot shape the fixture, and it fails **loudly** when the domain moves (proved: the probe turns `internal/ui` red). It is a second source of truth about the domain's project shape, and the price of it is that a domain change touches two files — which is the price of every wire this repo has (F-143, F-144). What it **cannot** do is see the *mechanism* regress — **F-146**. |
| 2 | the **SIGKILL** leftover is not fixed — whose is it? | **It does not matter for correctness, and it belongs to the OS's tmp, not to the tool.** Proved: a `kill -9` at 4s leaves the directory — **244K**, **43** files — in the user's `TMPDIR`. But each run makes its **own** `os.MkdirTemp` directory, so a stale one is **never read** and can never be mistaken for the current run's; the OS reclaims its tmp. The name says what it is (`maitre-ui-templates…`), so a reader who finds one knows. I would **not** prune stale directories at `TestMain`'s start: two lanes running the same suite at once would race each other's pruning. **A note, not a finding** — and the fix brief's "the cleanup is `TestMain`'s" holds. |

**The render path.** Not this cycle's, and not mine. Its brief is written (`work/briefs/tui-t12.md`) and goes to the chef first. It must achieve: **no process spawn in a snapshot**; **one** spawn for the crew read, not two; the **goldens** rendering with **herdr off the PATH**; and the live console still asking **asynchronously**. It must not break: what the Sous-Chef tab shows (the crew's rows, their states, their jobs), **not one golden** moving, **F-140's** witness (the mouse's switch and the filter line's switch still reach the runtime, and a click never leaves the tab saying *the automation is not running* while it is), and the **poll's frequency** — which is deliberate and stays with the poll; T12 changes *what* the poll asks, not *how often*.

**Mutations — four from the brief plus one of my own, every one reverted.**

| # | The mutation | What went red |
|--:|---|---|
| 1 | the template's project file **parsed and re-serialised** again, **with** the domain's field | **nothing in `internal/ui`** — the domain's field goes missing and the fixture drifts silently again. That is **F-146** (`internal/core` did go red) |
| 2 | the counters **plain ints** again | the race detector: **two `DATA RACE`s**, and `TestTheCountersAreSafeUnderAParallelTest` fails |
| 3 | the wire matching only the literal `core.Seed(` | an **aliased** call slips past and the wire is **green** — the hole; with the narrow pattern the count check also fires, because `core.SeedStress(` stops matching (*the sources call the seeder 1 times, but the suite seeded 2*) |
| 4 | the cleanup's message printing the **nil stat error** | the message ends with `<nil>`: *the cleanup did not remove the template directory: `<path>` is still there: `<nil>`* |
| 5 | **mine**: a seeder call smuggled into the template's own file | the count: *the sources call the seeder 3 times, but the suite seeded 2* |

All reverted; `git status --short` is empty except the ledger commit and this receipt.

**The questions.** `work/questions.md` — all **nine** rows are struck. **Nothing answered, nothing missed, nothing newly raised.** This cycle changed no behaviour and no vocabulary.

**DomainVersion:** 7 (`internal/core/lifecycle.go:225`). `internal/core` did **not** change — `git diff 968def2..HEAD -- internal/core/` is empty — and the budget is unchanged: **0 / 0 / 0 / 28**.

**The goldens:** unchanged — `git diff` over `internal/ui/testdata/` is empty, and `make goldens-check` regenerated all five with the tree clean.

**The cost**, measured by me, the same way as last cycle: **10.06s** plain and **38.06s** under `-race -shuffle=on -count=2` — the runner's **10.4s** / **37.8s** stand, and the store's numbers hold (about 10s plain, under 40s under the race gate).

**The gate.** `make check` **green** and `make check-race` **green**, confirmed by me on the committed tree.

**LOW notes.**
- The wire's pattern is textual, so a seeder call **in a comment** would count — e.g. a comment containing `core.Seed(` would make the count red with a message about calls that are not calls. Fail-safe (loud, not green), and the message names the file.
- The cleanup's condition is `!os.IsNotExist(err)`, so a stat error that is *not* "not exist" (a permission error, say) would also print *is still there* — which would be the wrong reason. Rare, and loud.
- `blankTheName` blanks **every** `"name"` it finds — the nested case above. That is the live half of **F-146**; the remedy that covers it is the witness, not a narrower regexp.
- `TestTheTemplatesAreRemoved` checks the removal of a directory it makes itself — it holds `removeTemplates`, not `TestMain`'s own call. `TestMain`'s call is held by the suite's exit code, which is how the mutation showed it.

**What I did not touch:** the implementation (mutated and reverted), `work/` (except this receipt), the runner (not prompted), and the render path (not this cycle's).
