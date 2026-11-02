Runner: T10 fix 3 done. Commits: 4 — `3fb3701` (AC-01), `4bf7d93` (AC-02), `a8d2bab` (AC-03), `d0f0f15` (AC-04).

**The third copy.** The startup task's copy of the format is now **one constant** (`rosterExample` in `internal/ui/souschef.go`) — the copy the console hands the orchestrator at spawn, in the same breath as the work — and it is **wired to the parser** the way the skill's copy is. The wire test (`internal/ui/crew_wire_test.go`) extracts **the startup task's own example** from the rendered task and feeds it to the console's `parseRoster` and `parseProgress`; a second test holds that the task's copy is **one of the skill's own lines**, not a variant of it.

**Your probe now fails.** I ran it by hand: the startup task's separator drifts from *" · "* to *" / "*, every word intact, and **two** tests go red — *the roster's line names the job "orchestrator / round 2 / 6 interactions", which is not one of the crew's* and *the startup task's copy is not the skill's own line — the skill's table is the specification*. Nothing red became two red.

**The specification.** It is said in three places, where a reader will find it: the **skill's table** says so in its own section — *"This table is the specification. The console's parser and the startup task's example follow it: when the format changes, this table changes first, and the wire test holds the other two to it"* — and the **skill's own witness** holds that claim; the console's **constant** says it in its doc comment; and the **wire's** comment says it for the reader who comes looking for the wire. The startup task's doc says that its copy is **an example**, not the specification.

**The flag.** It names the allowance that ran out. The budget is **two** allowances — the rounds and the interactions — and *"the budget is spent"* was fired by both, so *round 2 of 3 · 12 of 12 interactions · the budget is spent* said the budget was spent though the round allowance was not. The stamp now reads **"the rounds are spent"**, **"the interactions are spent"**, or **"the rounds and the interactions are spent"** — one for each case, and the witness holds all three.

**DomainVersion:** 7. The tab count stays **8**; `internal/core` does not change; the budget is **28 of 30**.

**The gate:** `make check` — **pass** (internal/ui 197.7s). `make check-race` — **pass** (internal/ui **529.1s**, internal/core 106.2s, internal/cli 13.1s).

Notes for the reviewer:

1. **The startup task's copy is now a constant, not a literal.** The task's rendered text is unchanged — `herdr tab rename <TAB_ID> "maitre: orchestrator · round <N> · <M> interactions"` — but it is one literal in the package instead of a retyped copy, so the console's side of the format lives in **one** place: the constant and the parser, joined by the wire to the skill's table.
2. **The equality test is the "one answer" check**: the startup task's copy must be *one of the skill's own lines*. It is a second, independent wire from the parse check — a drift that still parses but is not the skill's line is caught too.
3. **The flag's wording is mine**: *"the rounds are spent"* · *"the interactions are spent"* · *"the rounds and the interactions are spent"*. The requirement was yours — the reader must not have to guess which allowance ran out.
4. **The seeder and the render path are untouched** — not this cycle's, as ruled. The gate's headroom moved from 523.0s to **529.1s** of go test's 600s default under `-race -shuffle=on -count=2` (**88%**); this cycle's four tests cost about what the ones they replace cost, and the **shared seeded store** is the next slice and the only thing that buys the margin.
5. **Nothing else moved**: `internal/core` untouched, `core.DomainVersion` **7**, the tab count **8**, the budget **28 of 30**, and no file outside this cycle's scope was touched.
