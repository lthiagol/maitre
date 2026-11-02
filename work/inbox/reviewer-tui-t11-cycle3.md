Reviewer: T11 cycle 3 done. Verdict: **ok**. Filed **0** blockers (none). **F-146 is closed**, and with it **T11's findings are all closed** — the slice's closing row is written.

**The witness asks the DOMAIN.** `core.Seed` calls `store.Init` (`internal/core/seed.go:33`), so `core.Open(dir).Init()` writes the domain's **own** project file **without a seed** — one `Init()`, no `core.Seed`. `TestTheTemplateIsTheDomainOutput` parses the template's file and that one, drops `name` and `created` from both, and compares with `reflect.DeepEqual` — the **field set and the values**, nested included. **No list of names anywhere in the witness** (`grep` for `projectFields` finds nothing). It asks the domain for **both** templates — the ordinary one and the pathological one.

**My two cases are red, reproduced by hand and reverted.**

| the case | the result |
|---|---|
| the rewrite parses and re-serialises again **with** the domain's field | **2 red** — `TestTheTemplateIsTheDomainOutput`: *the template's project file is not the domain's own* (got without the field, want with it), and `TestTheRewriteBlanksOnlyTheTopLevelName` |
| a **nested** name that gets blanked | **2 red** — the same witness, plus the rewrite's own |

Last cycle both of these were **green**. That was the finding; it is closed.

**The rewrite blanks only the top-level name, and the comment is true.** `topLevelValue` walks the document's **top level** with `json.Decoder` and locates the value through the **decoder's offsets**; `blankTheName` replaces that value's bytes with `""` and copies the rest verbatim. Proved by command:

- a **nested** name survives — with the domain's `settings` carrying `{"theme":{"name":"dark"}}`, the template holds `"name": "dark"` ✓, and the whole ui suite stays green ✓ (the fixture equals the model's);
- the template's project file is **byte-identical** to a freshly seeded plan's **apart from the name value and the created stamp** ✓ — so *\"every other byte is left as the domain wrote it\"* is **true**, and a doc that contradicts the code is not a bug here.

**The wire counts CALLS, not mentions.** The sources are parsed with `go/parser` and `go/ast`, and the calls to `Seed` / `SeedStress` are counted by name. Proved by command: a comment saying `core.Seed(core.Open(\"kitchen\"))` is **not** counted (the wire is green with it in the package); an **aliased** call **is** counted (*zz_mutation_test.go calls the seeder: Seed()*); `seededPlan(` is not a call; and the sources' count still has to agree with the counter — a seeder call smuggled into the template's own file → *the sources call the seeder 3 times, but the suite seeded 2*.

**The cleanup and the removal.** `cleanupTemplates` reports the **actual** reason — *is still there* · *could not check* · *could not remove* — and `whyNotRemoved` tells the two apart. The removal witness drives **the function `TestMain` calls** (`cleanupTemplates`), and **the call itself** is held by reading the source: take `cleanupTemplates` out of `TestMain` and `TestMainRemovesTheTemplates` goes red — *seed_test.go: TestMain does not call cleanupTemplates — the templates would survive the suite* (proved, reverted; the litter the mutation left was cleaned up).

**The note — confirmed: a domain that ADDS a field is red in `internal/core` and green in `internal/ui`, and that is right there.** Proved by command — the domain's project file gains an `owner` field:

| | |
|---|---|
| **`internal/core`** | **red** — `TestDomainShapeIsVersioned`: the domain is fenced by its own version wire |
| **`internal/ui`** | **green** — and the green is a *true* statement: the fixture **keeps** the field. The template holds `"owner": "the-chef"`, and the witness compares it against the domain's own output and is satisfied |

`internal/ui` holds the **fixture's fidelity to the model**; `internal/core` holds the **model's move**. Each fact is in its own place, and neither needs the other — the old list made the ui suite a **second** witness of the domain's shape, which is the copy F-145 was about. Nothing is lost: a field the console renders moves a **golden** (red), and a field nobody reads needs nothing in ui. The green is the fixture's, not a silence.

**Mutations — four from the brief plus one of my own, every one reverted.**

| # | The mutation | What went red |
|--:|---|---|
| 1 | the rewrite **parses and re-serialises** again, **with** the domain's field | **2** — *the template's project file is not the domain's own* |
| 2 | the rewrite blanks a **nested** name | **2** — the same witness, plus the rewrite's own |
| 3 | the wire counts **mentions** again | the comment counts → *zz_mutation_test.go calls the seeder: Seed()* |
| 4 | the removal call taken out of **`TestMain`** | `TestMainRemovesTheTemplates` — *TestMain does not call cleanupTemplates* |
| 5 | **mine**: an **aliased** call in a test file | the wire — *zz_mutation_test.go calls the seeder: Seed()* |

All reverted; `git status --short` is empty except the ledger commit and this receipt.

**The questions.** `work/questions.md` — all **nine** rows are struck. **Nothing answered, nothing missed, nothing newly raised.**

**DomainVersion:** 7 (`internal/core/lifecycle.go:225`). `internal/core` did **not** change — `git diff cbc6ace..HEAD -- internal/core/` is empty — `docs/03` § 8 did not move, and the budget is unchanged: **0 / 0 / 0 / 28**.

**The goldens:** unchanged — `git diff` over `internal/ui/testdata/` is empty, and `make goldens-check` regenerated all five with the tree clean.

**The cost**, measured by me: **9.38s** plain and **38.25s** under `-race -shuffle=on -count=2` — the runner's **9.8s** / **38.6s** stand, and the store's numbers hold (about 10s plain, under 40s under the race gate).

**The gate.** `make check` **green** and `make check-race` **green**, confirmed by me on the committed tree.

**LOW notes.**
- The cleanup's third answer — *could not remove* — is **reachable** (proved: a path under a file) but **no test witnesses it**; the other two answers are witnessed. A note, not a finding.
- The witness drops exactly **two** names, `name` and `created` — the two that are not facts about the plan's shape. If the domain added a **second** timestamp field, the witness would compare two different moments and go **red** — loud, not silent, and the message shows the shape. Worth knowing before the plan's file grows a second clock.
- The wire counts calls by **name** in any receiver, so an unrelated `rand.Seed(...)` in a ui test would count against the suite. There is none today; over-inclusive, and loud.
- A seeder reached through a **function value** (`f := core.Seed; f(store)`) is not a *call* in the source and would not be counted. The old wire missed it too, and the cost it would hide shows in the suite's own wall time — the thing the slice protects.
- The witness compares parsed **values**, not bytes: a rewrite that re-serialised with the same values would be green. The bytes are held by the rewrite's own witness and by construction, and nothing but the domain reads the fixture's bytes — a note.

**The render path.** Not this slice's, and not mine. `work/briefs/tui-t12.md` goes to the chef first: it must achieve **no process spawn in a snapshot**, the crew read in **one** spawn, the **goldens** rendering with **herdr off the PATH**, and the live console still asking **asynchronously** — and it must not break what the Sous-Chef tab shows, a single golden, **F-140's** witness, or the **poll's frequency**.

**What I did not touch:** the implementation (mutated and reverted), `work/` (except this receipt), the runner (not prompted), and the render path.
