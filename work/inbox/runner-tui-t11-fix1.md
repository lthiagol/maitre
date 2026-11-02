Runner: T11 fix 1 done. Commits: 4 — `2accb50` (AC-01), `3c2be62` (AC-02), `6ccfbdc` (AC-03), `5a82933` (AC-04). The receipt is the fifth.

**The project file:** the name is rewritten **textually**. `unnameTheTemplate` reads the seeder's own bytes and empties the name's value where it stands — a regexp on `"name"\s*:\s*"` and the value it holds — and nothing else. Nothing is parsed and nothing is re-serialised: the file the domain wrote is the file the template keeps, so a field the domain adds is in the template, because the template **is** the seeder's output. The copy of the plan's schema — the fixture's own struct — is gone.

**The probe now fails.** The domain's project file gains a field, and a test goes red:

    got:  [charter created name owner settings]
    want: [charter created name settings]

The wire is `TestTheTemplateCarriesTheDomainFields`: the template's project file must carry the fields the domain writes, and the list of names fails loudly the moment the domain moves. It is a list of **names** and not a copy of the schema — it changes nothing, and a copy of the schema in the *rewrite* is what F-145 was. Probed by command: the field was added to `core.projectFile`, `internal/ui` went red on the wire, `internal/core` went red too (in `TestDomainShapeIsVersioned`, as it must). Both reverted — `git diff` over `internal/core/` is empty.

**The counters:** `seeds` and `copies` are `atomic.Int64` — shared state that a test calling `t.Parallel()` cannot race. The witness is `TestTheCountersAreSafeUnderAParallelTest`: two subtests run in parallel and draw plans at the same time, and the parent's cleanup checks the count once they are done. Probed: with plain counters the race detector says `DATA RACE`; with these it says nothing.

**The cleanup's failure** says what happened and where the leftover is — *the cleanup did not remove the template directory: `<path>` is still there* — in the test and in `TestMain`'s own message. It no longer prints the nil stat error.

**The wire** matches the **call**, not one spelling: `\bSeed(Stress)?\s*\(`, so a qualifier, an alias, a space before the parenthesis and a newline before it are all caught. Probed: an aliased `c.Seed(` and a spaced `core.Seed (` both go red. And the wire asks **how many** times: the count of seeder calls in the package's own sources must agree with the counter, which is what the suite actually did — so a seeder call smuggled into the template's own file is caught too (*the sources call the seeder 3 times, but the suite seeded 2*). The counter still covers the cost: it says the suite seeded once per template, whatever the run count.

**The goldens:** unchanged — `make goldens-check` regenerated all five and the tree stayed clean.

**The cost:** measured the same way as before — the test binary, run from the package.

| | before (the slice) | after (this cycle) |
|---|---|---|
| `internal/ui`, plain | 10.5s | **10.4s** |
| `internal/ui`, `-race -shuffle=on -count=2` | 38.3s | **37.8s** |

The gate's own numbers: `internal/ui` 44.6s in the slow gate, and the slow gate is 1m47s in total.

**DomainVersion:** 7 (`internal/core/lifecycle.go:225`). `internal/core` did not change.

**The gate:** `make check` green, `make check-race` green.

**Notes:**

- **The one judgement in this cycle:** the wire is a **list of field names** written in the test. It is not a copy of the schema in the *mechanism* — it changes nothing and fails loudly, and that is the difference F-145 was about. The alternative would be to hold no shape at all in the ui suite and lean on `internal/core`'s own version wire — but then the reviewer's probe is green in `internal/ui` again, which this cycle's AC-01 forbids. If the reviewer wants the pin somewhere else, it is one list and one test.
- The wire's pattern does not match its own source — the `\b` before `Seed` is preceded by a word character there — so the wire does not count itself. That is why the count is 2 and the counter says 2.
- The render path was **not** started and the poll's frequency was **not** touched: the brief for T12 goes to the chef first, as the fix brief says.
- The reviewer's other low note — a **killed** test binary (SIGKILL) leaves a 15-file plan in the user's tmp — is **not** fixed here. The removal and the survival check hold; the SIGKILL case does not.
