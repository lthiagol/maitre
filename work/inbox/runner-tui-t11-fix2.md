Runner: T11 fix 2 done. Commits: 4 — `d0ab68d` (AC-01), `55368d4` (AC-02), `ecebfe4` (AC-03), `d22cc8e` (AC-04). The receipt is the fifth.

**The witness** asks the **domain**, not a list of names. `core.Seed` calls `store.Init` (`internal/core/seed.go:33`), so `core.Open(dir).Init()` writes the domain's own project file — no seed, one `Init()`. The witness parses the template's file and that one, drops `name` and `created` from both, and compares: the field set and the values, nested included. **No list of names anywhere** — the hand-written list is gone.

Both of the reviewer's cases are red, proved by command and reverted:

- the rewrite parses and re-serialises again **and** the domain's project file gains a field → *the template's project file is not the domain's own*
- the domain's project file carries a nested name (`settings.theme.name`) and the rewrite blanks it → the same

**The rewrite** blanks **only the top-level name**. The value is located in the domain's own JSON — the top level, through the decoder's offsets — and replaced by an empty string; every other byte is the domain's. A nested name survives, and the code's comment is **true**: the template's file is byte-identical to the seeder's apart from the name's value (proved by dumping it). `TestTheRewriteBlanksOnlyTheTopLevelName` holds the rewrite itself, whatever the domain's data is.

**The wire** counts **calls**, not mentions. The sources are parsed (`go/parser`, `go/ast`) and the wire counts the calls to `Seed`/`SeedStress`: a comment that says `core.Seed(` is not a call, an alias and a newline are the same call, and `seededPlan(` is not a call. Probed: an aliased call is still caught, a comment is not. The count must still agree with the counter — the sources' count and what the suite actually did.

**The cleanup** reports the **actual** reason: *is still there* for a directory that survived, *could not check* for one that cannot be read, and *could not remove* for a removal that failed.

**The removal witness** drives the function `TestMain` calls (`cleanupTemplates`), and the **call itself** is held by reading the source: a `TestMain` that stopped removing the templates would leave them behind and the suite would stay green. Probed: take the call out of `TestMain` and the witness goes red.

**The goldens:** unchanged — `make goldens-check` regenerated all five and the tree stayed clean.

**The cost:** the store's numbers hold.

| | before (the reviewer) | after (this cycle) |
|---|---|---|
| `internal/ui`, plain | 10.06s | **9.8s** |
| `internal/ui`, `-race -shuffle=on -count=2` | 38.06s | **38.6s** |

The slow gate is 1m45s in total, with `internal/ui` at 44.7s of it.

**DomainVersion:** 7 (`internal/core/lifecycle.go:225`). `internal/core` did not change — `git diff` over it is empty.

**The gate:** `make check` green, `make check-race` green.

**Notes:**

- **The list of field names is gone.** The reviewer said it "stays, for the tripwire it is: a domain that moves must turn the ui suite red"; this cycle's brief says no hand-written list of names anywhere in the witness, and the brief wins. The consequence to know: a domain that **adds** a field is red in `internal/core` (`TestDomainShapeIsVersioned` — the domain is fenced) and **green** in `internal/ui` — and green is right there, because the fixture keeps the field. What `internal/ui` now holds is the *fixture's* fidelity to the model, which is what F-146 was about.
- The rewrite parses the domain's JSON to **locate** the top-level value; it does not re-serialise it, and every other byte is the domain's ✓.
- **Not touched:** the render path, the poll's frequency, and the SIGKILL leftover (a note, not a finding). No stale directory is pruned at `TestMain`'s start — two lanes would race each other's pruning.
