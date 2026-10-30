Runner: T3a done. Commits: 5 slice commits (4dc8079, 1b550c6, 1804ccc, b873821, 5edae71) + 2 work commits (a0ca384, 051081c). Per-AC evidence: AC-01..AC-06 stamped below. Pre-notify gate: pass.

**AC-01 — a backlog item has a detail pane.** The seam is `Model.hasDetail()` (internal/ui/model.go): the pane follows the tab, and the key and the pane read one answer, so `enter` opens whatever the pane shows. The item has its own rendering (`backlogDetail`), not the recipe's. The detail shows: the id and title, the kind and status, the dates (`created … · updated …`), the body's own words under **Why**, the **Scope**, the **Tags**, and **Promoted to recipe RCP01** when the item was promoted.

```
$ maitre --render --plain --keys "tab,enter"          # the item's detail, full screen
 BL01 · The board needs a way to show evidence
 idea · open
 created 2026-10-30T20:55:10Z · updated 2026-10-30T20:55:10Z

 → "idea" carries no acceptance criteria, and a recipe without them can never be signed off — shape it into a small task first

Why

    Evidence the chef cannot read is a claim, not evidence.

Tags
  board

$ maitre --render --plain --keys "tab,click@(5,6)"     # the click opens it too
 BL01 · The board needs a way to show evidence
```

**AC-02 — the console offers no move the domain refuses.** The console asks the domain (`core.PromotionOf`) and renders the answer; it does not decide by itself. On an **idea** the footer has **no** `space` key at all — `↑↓ move · tab switch · enter detail · r reload · q quit` — and the item's detail carries the domain's reason. On a **small task** the footer still reads `space promote to recipe`. Pressing `space` on an idea asks the domain to promote and shows the domain's reason in the refusal banner:

```
$ maitre --render --plain --keys "tab,space"
 ⛔ refused — "idea" carries no acceptance criteria, and a recipe without them can never be signed
 off — shape it into a small task first
```

**AC-03 — D-3: an idea cannot be promoted into a dead end.** The refusal is the domain's, and the rule lives in one place: `core.PromotionOf(item)` (internal/core/backlog.go) — a **derived view**, registered in `core.DerivedViews` and given a row in `docs/03` § 8. `PromoteBacklogItem` reads it before it writes and refuses with its reason; the console's footer reads the same answer. A small task still promotes, with its one criterion:

```
$ maitre-cli backlog promote BL01          # BL01 is an idea
{ "details": ["\"idea\" carries no acceptance criteria, and a recipe without them can never be
   signed off — shape it into a small task first"], "error": "failed", "ok": false }   (exit 1)

$ maitre-cli backlog promote BL02          # BL02 is a small task
{ "ok": true, "recipe": { "acceptance": [ { "description": "Wire the review queue",
  "id": "AC-01", … } ], "id": "RCP01", "lifecycle": "draft" … } }                        (exit 0)
```

**AC-04 — D-5: every writer writes the same shape.** The rule is in the **writer**: `Store.SaveRecipe` calls `r.normaliseLists()` before it writes, so a recipe written *first* time says the same thing as one read and written back. Before: `"scope": {"in_scope": null, "out_of_scope": null}, "acceptance": null`. Now:

```
kitchen/recipes/RCP02-a-plain-small-task.json   (the promotion)
  scope     : {"in_scope": [], "out_of_scope": []}
  acceptance: [{"id": "AC-01", "description": "A plain small task", …}]
kitchen/recipes/RCP03-a-plain-recipe.json        (recipe create)
  scope     : {"in_scope": [], "out_of_scope": []}
  acceptance: []
```

The **digest does not move**: the approval's payload canonicalises its lists, so `null` and `[]` hash the same. `TestTheRepoPlansApprovalDigestMatchesItsContent` (the repository's own plan, read from the file's own bytes) and the new `TestTheApprovalDigestStillVerifies` (approve → write back → read from disk → recompute) both pass.

**AC-05 — the goldens.** Unchanged. `make goldens` leaves `git diff` empty: the goldens hold the **Recipes** tab, and this slice changed the **Backlog** tab and the shape a recipe is *written* in — not what a recipe reads as.

**AC-06 — the state row.** `work/tui/README.md`: T3a is `✅ closed · 5 commits`, and **T3b** — the backlog's moves: dismiss · create · edit — is named as the next session.

**The refusal:** the reason comes from the domain, from `core.PromotionOf` — `"idea" carries no acceptance criteria, and a recipe without them can never be signed off — shape it into a small task first`. It names the remedy. The console renders it; it does not decide it.

**The footer:** on an idea — no `space` key at all (there is no move to offer), and the detail carries the reason. On a small task — `space promote to recipe`.

**The shape:** the promotion writes `"scope": {"in_scope": [], "out_of_scope": []}` and `"acceptance": []` (or the small task's one criterion); `recipe create` writes the same shape.

**The digest:** an existing plan's approval still verifies — the repository's own plan (`TestTheRepoPlansApprovalDigestMatchesItsContent`) and a fresh approve → write → read → recompute (`TestTheApprovalDigestStillVerifies`).

**The view:** `core.PromotionOf` — a **derived view**, with the `// view` marker, registered in `core.DerivedViews` as `promotion`, and a row in `docs/03` § 8 (the wire checks both directions). It was not declared in `core.NonViews`: it is computed and never stored, so it is a view and not a helper.

**DomainVersion:** 7 — unmoved. `TestDomainShapeIsVersioned` passes; "a list is `[]`, never `null`" is the format's own rule, and the writer now keeps it.

**The gate:** `make check` green, `make check-race` green, `gofmt` clean, `go vet` clean, staticcheck clean, the dependency budget unmoved (0 / 0 / 28 of 30). The tree is clean.

**Notes for T3b:**
- The **kind change** is still missing, and the refusal points at it: *"shape it into a small task first"*. `UpdateBacklogItem` already changes the kind (`backlog update --kind feature --scope …`), so the console's edit surface is the only thing missing. `core.PromotionOf` is the answer the console must ask after the edit — the footer's `space promote to recipe` will appear by itself, because the footer reads the domain.
- The refusal **writes a journal event** (`refuseLocked`), so a refused promotion is auditable — and the console's `space` on an idea therefore leaves a `refusal` event in the plan. Worth knowing before the console's own moves start writing.
- `maitre-cli backlog promote` on an idea exits **1**, like every other `RefusalError` (a *gate* refusal is 2). No verb changed.
- The backlog item's file shape was **not** touched: `scope` and `tags` are `omitempty`, so they are absent rather than `null`. Making them always-present would be a format change and a `DomainVersion` bump — out of this slice.
- `docs/03` § 1.1 and `docs/07` were corrected: they said an idea could be promoted straight to a recipe, which AC-03 made false.
