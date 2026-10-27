Reviewer: slice C review done. Verdict: changes-needed. Closed: F-73 (the repo's own
plan now names the project `maitre`; the console reads `maitre · maitre`). Filed 4 new
blockers (F-74 high, F-75/F-76/F-77 medium).

F-74 (high) — the live console never shows an old plan's specials and prep notes.
`reload()` reads them into `reloadMsg`, but `Update`'s `reloadMsg` branch assigns
`m.recipes`, `m.items` and `m.journal` and NOT `m.legacy`, so the field is always
empty in the running console. `--render` shows them because `snapshot()` assigns the
field itself; the live path never does. So the console lies about the plan it shows:
an old plan's backlog reads as empty. The runner's own `TestAnOldPlanStillRenders`
passes only because `loadModel` (the test helper) copies `loaded.legacy` by hand —
the path the chef actually uses is untested. Proof: mutating `Update` to assign
`m.legacy = msg.legacy` makes the legacy group appear.

F-75 (medium) — a legacy plan's dangling recipe pointer no longer reserves its id.
`reservedRecipeIDs` reads `PromotedTo` from the backlog items only, so the pointers in
`kitchen/prep-notes/*.json` are not scanned, and a new recipe can be handed the id a
dangling pointer names. Base: the same plan reserved `RCP10`. Head: it hands out
`RCP04` and the pointer is unprotected — the failure `NextRecipeID`'s own comment says
must never happen. `core.LegacyItem` does not carry a `PromotedTo` field either.

F-76 (medium) — `backlog update` silently wipes the item's tags.
`UpdateBacklogItem` guards the tag list with `if tags != nil` while the CLI's
`StringList` returns a non-nil empty slice when the flag is absent, so
`backlog update BL01 --title "X"` drops the tags with no warning. The scope list,
one line away, uses `len(scope) > 0` and IS left alone — so the two list fields of
one operation behave differently, against the rule the doc states two lines under the
row this slice added.

F-77 (medium) — `docs/03` § 1.1's small-task row says a promoted small task's
recipe is "falsifiable from the start", while the paragraph below it says the
criterion carries no verification and cannot pass `ac.falsifiable`. The code follows
the paragraph.

What holds, verified:
- AC-01 — one artifact, one directory. `BL01`/`bl01`/`BL1`/`B01`/`1` all resolve;
  a file named `B01-<slug>.json` reads as `BL01`. Mutating the normaliser to strip a
  fixed two letters fails 6 tests (recipe and backlog). `Init` creates
  `kitchen/backlog/` and not the old two.
- AC-02 — one noun. `special`/`prepnote` are gone from the dispatch table (unknown
  command). `create --kind` defaults to `idea`; `update --kind` changes the kind and
  keeps the id; `dismiss` works on every kind. The count is 46 -> 42 and the doc and
  the counting tripwire agree. Proof: mutating AGENTS.md to 41 fires the wire;
  mutating the code to drop one verb fires it too (both directions checked).
- AC-03 — `promote` creates a recipe in `draft`; a small task gets exactly one
  criterion, an idea none; the item becomes `promoted` and records `promoted_to`;
  promotion is one-way. No gate weakened: `gates` is byte-identical to the base (only
  the `DomainVersion` constant moved), no new plan-format field, no special case in a
  rule. Proof: giving the promoted criterion a verification fails
  `TestThePromotionWeakensNoGate`.
- AC-04 — the groups are derived from the kind and never stored (`BacklogGroups` is
  the only place the rule lives); a row shows the kind. An old plan renders and the
  console is honest about what it found — but see F-74: only in the snapshot.
- AC-05 — the seed and the chaos seed create backlog items; `demo` still renders.
  `docs/03` § 6 shows `kitchen/backlog/` and drops the two old directories. GLOSSARY
  has one Backlog item entry. F-73 closed.
- AC-06 — `go test ./internal/core ./internal/cli ./internal/ui -run 'Doc' -v` exit 0.

Beyond the ACs:
- The journal is append-only — `git diff` on `kitchen/activity.json` shows one
  addition (seq 34, a control/migration event), never a rewrite. The plan's stamp
  moved to 5.
- The old ids still resolve: `S01`, `N01`, `BL01` and `1` all name the same item.
- The fence stands in its designed direction: a plan stamped 6 is refused, a plan
  stamped 4 or 5 is accepted and re-stamped to 5.
- The id is kind-neutral: `BL01` stays `BL01` when an idea is shaped into a feature.
- Scope is clean — the roles, the tabs, the recipe ids, the module path, the binaries
  and `Proof` are untouched.
- Invariants hold: `internal/core` and `internal/cli` are stdlib-only, `internal/ui`
  keeps 28 external modules, the domain layer has zero UI imports.

LOW notes (not filed, per the brief):
- `backlog update` on a dismissed or promoted item still refuses with a plain
  "already dismissed" — the reason is not in the refusal text, so the console's
  banner cannot say why.
- `core.LegacyItem.Text` returns the prep note's `intent` under a field the item
  calls `body`; the two legacy artifacts keep one prose field each, but the console
  labels both "Why".
- The `GLOSSARY.md` CPD mapping row collapses `idea / note` and `brief` into one
  entry; the split is in the glossary's `_Avoid_` line rather than in the mapping.
