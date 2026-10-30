Runner: T3a fix 1 done. Commits: 4 (aa4ba5f, da36146, 99a7576, 8f39289) — 3 code/doc commits and the state row. Fixed: F-114, F-115.

**F-114 — the footer's guard bites now.** The test uses the **promoted small task**: an item that was a `feature`, was promoted, and is no longer open. It is the case where the console's own rule ("is it a small task?") and the domain's answer **differ** — the kind rule says yes, `core.PromotionOf` says no, because the item is no longer open. The two cases the old test held (an idea, an open small task) are the two where the two answers **agree**, so a console that decided for itself passed them. `TestTheFooterAsksTheDomainAboutPromotability` (internal/ui/backlog_test.go) holds the differing case: promote BL05, then assert the footer offers **no** `space promote to recipe` on it. The reviewer's mutation — the footer deciding by the item's kind (`item.Kind == core.KindIdea`) instead of asking the domain — now turns the suite **red**, and **exactly one** test fails:

```
$ go test ./internal/ui        # with the mutation applied
--- FAIL: TestTheFooterAsksTheDomainAboutPromotability
    backlog_test.go:178: the footer offers a promotion the domain refuses — the console decided for itself
1 test failed
```

The mutation was reverted; `git diff` is empty.

**F-115 — the refusal writes nothing.** The console now asks `core.PromotionOf` — a **read** — and shows the domain's reason; it does not call the operation, so nothing is journalled. `activity.json` is **byte-for-byte unchanged**:

```
$ shasum -a 256 kitchen/activity.json      # before
8960774117cca6aab1abaeb5ff8c494a0f0a7fbe6d7f5ebbcfc78666d7075f8c  kitchen/activity.json   (16418 bytes)
$ maitre --render --plain --keys "tab,space"
$ shasum -a 256 kitchen/activity.json      # after
8960774117cca6aab1abaeb5ff8c494a0f0a7fbe6d7f5ebbcfc78666d7075f8c  kitchen/activity.json   (16418 bytes)
```

**Where the reason lands now:** in the **banner**, as before — `⛔ refused — "idea" carries no acceptance criteria, and a recipe without them can never be signed off — shape it into a small task first` — and in the **item's own detail**, on the line the domain's answer puts there. Neither writes to the plan. The console's precedent is the recipe path's: it asks `core.ChefMoveOf` and does not call the operation when there is no move.

**The CLI's refusal:** it stands, with the reason and exit 1 — an agent asked, and the plan said no. The journal records that attempt, as it should:

```
$ maitre-cli backlog promote BL01
{ "details": ["\"idea\" carries no acceptance criteria, and a recipe without them can never be signed off — shape it into a small task first"], "error": "failed", "ok": false }
exit=1
$ … kitchen/activity.json
59 refusal BL01 {"reason": "\"idea\" carries no acceptance criteria …"}
```

**DomainVersion:** 7 — unmoved. `internal/core` and `internal/cli` are untouched by this cycle; no domain change was needed.

**The goldens:** unchanged — `make goldens` leaves `git diff` empty.

**The gate:** `make check` green, `make check-race` green (core 90 s, ui 318 s), gofmt/vet/staticcheck clean, the budget unmoved (0 / 0 / 28 of 30). The tree is clean.

**Notes for the reviewer:**
- **`internal/core` was not touched.** The remedy was already in the domain: `core.PromotionOf` is a read, and the console asks it.
- **`GLOSSARY.md` was corrected** — the review's LOW note (a): "It becomes a recipe when it is promoted" read as though an idea could; it now says "**A small task** becomes a recipe when it is promoted". One sentence, and it was the only line left saying the old thing.
- The review's LOW note (b) — the detail's `→` on a promoted item ("→ backlog item BL05 is promoted") reads like an instruction rather than a fact — is **not** fixed here: it is not a blocker and this cycle does not widen. It wants a different marker for the domain's answer than the recipe detail's "why it needs the chef" arrow.
- The review's LOW note (c) — the Log tab's detail is still unreachable with `enter` in the narrow layout (`hasDetail` covers Recipes/Path/Board/Backlog) — is **not** fixed here either. The ledger names it as T3a's territory; it is one line in `hasDetail` and one test, and it is the orchestrator's call whether it rides in this cycle or the next.
- The **backlog item's own absent lists** were left alone, as ruled.
