You are a **fresh runner**. The repository is `/Users/thiago/Code/Pessoal/maitre`.
Read `AGENTS.md` first, then this brief. This brief wins where the design and the
docs differ. The orchestrator's pane id is `w1H:p2`.

This slice is **T11 — the shared seeded store**. It changes **no behaviour**. It makes
the test suite pay for the plan it builds **once** instead of two hundred and
nineteen times.

**Read this before anything else.** The gate is at **89%** of go test's 600s
default timeout. `internal/ui` takes **532.6s** under `-race -shuffle=on -count=2`.
The next slice of console work lands at ~95%, and the one after it **breaks the
gate**. This slice is not optional, and it is not a console slice — it is the
slice that makes the next console slice safe.

## The cost, measured

| | |
|---|---|
| **the suite builds a plan in** | **219** places |
| **`core.Seed` costs** | **0.608s** — it writes **15** files |
| **so the seeder costs** | about **133s** of the **195s** non-race run |
| **a copy of a seeded plan costs** | **0.007s** — **87× cheaper** |
| **the store makes it** | 0.6s once, plus 219 copies at 0.007s → about **65s** non-race, and **under 200s** under the race gate |

## The shape

**One template, seeded once. Every test gets a copy.**

| | |
|---|---|
| **the template** | seeded **once**, into a directory the suite owns |
| **the copy** | **byte-for-byte**. The goldens pin the seeded plan's content, so **nothing inside is touched** |
| **the project's name** | comes from the **parent directory's name** (`store.Init` → `filepath.Base(filepath.Dir(s.Root))`), so the copy lands in a directory named **`demo`** and the name follows |
| **the stress plan** | a **second** template. `core.SeedStress` is a different plan and does not come from the first one |
| **the cleanup** | `TestMain`, not a user's tmp |

**The seed happens once, not per test.** The suite runs **twice** in one process
under `-count=2`, so:

| | |
|---|---|
| **`-count=2`** | the template is seeded **once** and copied **twice**. A per-test `core.Seed` is exactly what this slice removes |
| **`sync.Once` or `TestMain`** | one of the two. Not a per-test call |
| **the template survives** | copying it twice must produce two identical copies |
| **no fresh timestamps** | no test asserts a fresh one. The reviewer checked, and none does — but **you** check it too, and say what you found |

## What must not break

| | |
|---|---|
| **isolation** | each test its **own** copy. Under `-shuffle=on -count=2` the writers must not see each other |
| **the goldens** | every golden file must be **unchanged** by this slice. The name comes from the directory, so the copy lands in a directory named `demo` — if a golden moves, the copy landed in the wrong place |
| **the stress plan** | a **second** template, not the same one |
| **the cleanup** | the template's directory is removed by `TestMain` — not left in a user's tmp |
| **the behaviour** | **nothing changes**. The console reads the same plans and renders the same screens |

## The acceptance criteria

| AC | What |
|--:|---|
| **AC-01** | **one template, seeded once.** Every test that needs a plan gets a **copy**, and the copy is **byte-for-byte** the template. A test holds it: two copies made at different times have the same content, and the template still holds exactly the files the seeder wrote |
| **AC-02** | the copy lands in a directory whose **name** gives the project its name — `demo` — so **every golden file is unchanged**. The stress plan is a **second** template |
| **AC-03** | it survives **`-count=2`**: the suite runs twice in one process, the template is seeded **once** and copied twice, and the copies do not see each other's writes under `-shuffle=on`. The template's directory is removed by **`TestMain`** |
| **AC-04** | the cost is **held**: the package's tests no longer seed their own plans, and the suite's time is reported — the non-race run and the race run, before and after. The target is **well under 200s** under the race gate: a margin of at least **3×** |
| **AC-05** | `work/tui/README.md`'s T11 row names this slice's commits. `DomainVersion` is **7** |

One commit per criterion. The state row's own commit is not one of the ids it names.

## The invariants

| | |
|---|---|
| **`core.DomainVersion`** | stays **7**. Nothing about the plan changes |
| **`internal/core`** | does **not** change, unless the copy helper belongs there — and if it does, say why in the receipt |
| **`internal/ui`**'s behaviour | **unchanged**. The same screens, the same goldens |
| **The budget** | unchanged: 0 / 0 / 0 / 28 |
| **No new view** | nothing is registered in `core.DerivedViews`, so `docs/03` § 8 does not move |
| **The gate** | `make check`, then `make check-race` before you report |

## The scope — what not to touch

| | |
|---|---|
| **The render path's wait** | **not** this slice. `render.go`'s `snapshot()` → `readCrew()` waits on a process, and the poll spends two spawns a second. That is the **render path's** cost, not the seeder's — it is a cheap fix of its own, and it belongs to whoever next touches the **render path**. Name it in your receipt; do not start it |
| **`internal/core`** | does not change, unless the copy helper belongs there (see the invariants) |
| **The console's behaviour** | nothing changes. This slice moves **no** screen and **no** golden |
| **`kitchen/`** | do not rewrite the repo's plan |
| **`docs/audit.md`** | the reviewer's |
| **`work/`** | the orchestrator's, except the T11 state row and your receipt |

## The report

```bash
mkdir -p work/inbox
cat > work/inbox/runner-tui-t11.md <<'EOF'
Runner: T11 done. Commits: <n>. The store: <where the template lives, how the seed happens once, and how a test gets its copy>. The copy: <that it is byte-for-byte, and the witness>. The names: <how the copy lands in a directory named demo, and that every golden is unchanged>. The stress plan: <the second template>. Count two: <what survives -count=2, and the cleanup>. The cost: <the non-race and race times, before and after, and the margin>. The render path: <where its wait is named, and that you did not start it>. DomainVersion: <7>. The gate: <make check, make check-race>. Notes: <anything the reviewer should know>.
EOF
git add work/inbox/runner-tui-t11.md && git commit -q -m "inbox: runner reports tui-t11"
herdr pane send-text "w1H:p2" "<the receipt, as above>"
herdr pane send-keys "w1H:p2" enter
```

**The nudge is `pane send-text`, not `agent prompt`.** `herdr agent prompt` refuses a
pane that reports through the `mcode-herdr` bridge — it answers `agent_not_ready`,
because the bridge registers the pane's *state* and not a named session the prompt
path can address. The receipt file is the record; the nudge is only the bell.

Write the receipt, commit it, and then nudge. If a file you did not write has
changed, stop and report it. Stage by name. Then stop.
