You are the **runner** — the same lane that wrote slice C. The orchestrator's pane id is
`w1Z:p1`; use it verbatim in the notify command below.

## The fix cycle for slice C — four findings

The reviewer filed **F-74 (high)**, **F-75**, **F-76** and **F-77**. The orchestrator
ruled on all four: **accepted**. Four findings, four commits — one each, so the
reviewer can judge them apart.

### F-74 (high) — the live console never shows an old plan's specials

`internal/ui/model.go`'s `case reloadMsg:` assigns `m.recipes`, `m.items` and
`m.journal` — **not** `m.legacy`. The reload message carries `legacy` (line 256) and the
branch drops it. So `m.legacy` is always empty in the running console, and an old
plan's backlog reads as empty.

`--render` does not have the bug: `internal/ui/render.go` assigns the field itself.
**The console lies to the chef** — and the console is this project's main surface.

And the test that was supposed to hold it does not:

- `internal/ui/console_test.go`'s `loadModel` helper copies `loaded.legacy` **by hand**,
  so the test exercised the helper, not the path the chef uses.
- That is the standing lesson: a test that sets up the state it is testing is not a
  test.

**The fix:**

1. The live path assigns the field: `case reloadMsg:` must carry `m.legacy` the way it
   carries `m.items`.
2. **The test must drive the real path.** If the helper hand-copies the field, the
   helper is the bug. Make the test go through the reload — the same `reloadMsg` the
   live console gets — and let it fail before the fix.
3. **Prove the fix by mutation**: change the assignment back and watch the test fail.
   A test that cannot fail is not a test.

**One fact, one answer.** If the reload path and the `--render` path each assign the
field themselves, there are two places that know the model's shape. If the fix makes
them two, that is the finding again in a smaller coat.

### F-75 (medium) — a dangling pointer no longer reserves its id

`internal/core/store.go`'s `reservedRecipeIDs` scans `ListBacklogItems()` only, so the
`promoted_to` pointers in `kitchen/prep-notes/*.json` are not scanned. `core.LegacyItem`
does not even carry a `PromotedTo` field, so the legacy reader drops the pointer on the
floor.

The code's own comment on `NextRecipeID` says what must never happen:

> A promotion records where it went before the target is written, so a crash can leave
> a pointer to an id that does not exist — and that pointer must resolve to nothing,
> never to a stranger.

The same plan reserved `RCP10` before this slice; now it hands out `RCP04` and the
pointer is unprotected. (`len()+1` id allocation is a standing v1 lesson — it produced
`[AC-01, AC-03, AC-03]`.)

**The fix:**

- `core.LegacyItem` carries `PromotedTo` (JSON `promoted_to`), and the legacy reader
  reads it.
- `reservedRecipeIDs` scans the legacy items' pointers as well as the backlog's.
- A test: a plan with a legacy prep note pointing at `RCP10` must not hand out
  `RCP10` to anything else.

**Do not** add a migration verb or a "fix the dangling pointer" command. That is the
open question the orchestrator put to the chef, and it is not this cycle's business.

### F-76 (medium) — `backlog update` silently wipes the tags

`UpdateBacklogItem` guards the tags with `if tags != nil`, and the CLI's `StringList`
returns a **non-nil empty slice** when the flag is absent. So
`backlog update BL01 --title X` drops the tags, with no warning.

The scope list, one line away, uses `len(scope) > 0` and **is** left alone — so the two
list fields of one operation behave differently.

The rule the doc states (`docs/04-operations.md`, *Every field of the model has a
writer*):

> A list flag replaces the list it names; a flag that is not given leaves its field
> alone.

**The fix:** both list fields obey that rule. A flag that is not given leaves the field
alone; a flag that is given replaces the list. One mechanism, not two guards.

**Two constraints, both of them the project's own rules:**

| | |
|---|---|
| **No `null` in JSON** | `docs/04-operations.md` also says *"A list is always a list: an empty one is `[]`, never `null`."* If the fix makes an absent flag read as nil internally, make sure **nil never reaches a JSON document**. The JSON edge is `[]`. |
| **The other verbs are untouched** | `StringList` is used by `--in-scope`, `--out-of-scope`, `--ac` and `--ac-verify` too. Whatever you change, those verbs must behave exactly as before. |

**One question I am not deciding for you, and you must not decide it either:** the
parser refuses a list flag with no value (`--tag needs a value`), so a list can be
replaced but never emptied. The doc's rule says a list flag replaces the list it
names — which implies `--tag` with no values should *clear* the tags. Changing the
parser's rule would touch every list flag, so it is out of this cycle. **Say so in
your report**; the orchestrator will put the question to the chef.

### F-77 (medium) — the spec contradicts itself about a promoted small task

`docs/03-domain-model.md` § 1.1's table says a small task's promoted recipe is *"a
recipe with a criterion is falsifiable from the start"*, and the paragraph below says
the criterion's verification is the recipe's own life, and *"a criterion with no
verification cannot pass `ac.falsifiable`"*. The code follows the paragraph.

**The paragraph is right and the code is right** — a promoted recipe is not falsifiable
until its criterion has its verification, which is the recipe's own life after
grooming.

**The fix:** the table's row. Say what is true — the criterion makes the work
checkable, and it becomes falsifiable when the verification is written.

> The orchestrator's slice C brief said *"so the work is falsifiable from the start"*.
> The runner put it in the doc because the brief said it. **A brief is not a
> specification**, and this is the cost of writing a phrase in a brief and not reading
> the code behind it. The brief is corrected too.

### Not yours — do not touch

| | |
|---|---|
| `docs/audit.md` | the reviewer's to write. It files F-74..F-77 and closes them. |
| `docs/07-vocabulary.md` | the decision of record — the orchestrator's to write. |
| `GLOSSARY.md`'s CPD mapping | slice E's. The reviewer's note about two CPD terms sharing one row goes to slice E. |
| The parser's "needs a value" rule | the chef's question. Report it; do not change it. |
| The gate table, the plan format, the command count (42) | none of the four findings touches them. If a fix would, **stop and report it**. |

### The gate

```bash
make lint && gofmt -l . && make test && git status --short
```

Green, and the tree clean apart from the orchestrator's own `work/` files.

### Four commits, one per finding

```
Per-AC evidence:
- F-74: the live reload path carries the legacy group; go test ./internal/ui -run 'Legacy|OldPlan' -v exit 0
```

### Notify the orchestrator

```bash
mkdir -p work/inbox
cat > work/inbox/runner-slice-c-fix-1.md <<EOF
Runner: slice C fix 1 done. Commits: <n>. Fixed: F-74, F-75, F-76, F-77. Notes: <anything the reviewer should know — especially the list-clearing question>.
EOF
git add work/inbox && git commit -q -m "inbox: runner reports slice-c fix 1"
herdr agent prompt "w1Z:p1" "Runner: slice C fix 1 done. Commits: <n>. Fixed: F-74, F-75, F-76, F-77. Notes: <anything the reviewer should know — especially the list-clearing question>."
```
