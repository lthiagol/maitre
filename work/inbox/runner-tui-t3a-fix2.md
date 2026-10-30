Runner: T3a fix 2 done. Commits: 2 (63964ce, c47728e) — the seam and the state row. Fixed: F-116.

**The Log's detail.** `enter` on the Log tab now opens the event's detail — the same pane the recipe gets, the same geometry, the same scroll, the same `esc`. `hasDetail()` answers for the Log, so the key and the pane read one answer:

```
$ maitre --render --plain --keys "tab,tab,tab,tab,enter"
  Recipes    Backlog    Path    Board  ▸ Log    Sous-Chef    Settings
 recent activity · 58 events
 created
 2026-10-30T22:03:09Z
 chef · chef

Data
      { "id": "BL06", "item_kind": "bug", "kind": "backlog",
        "title": "Sharpen the evidence contract" }
```

**The footer's promise.** The footer offers `enter detail` only where there is a detail to open — the same answer the key reads:

| Tab | The footer's keys |
|---|---|
| **Recipes** (a recipe is selected) | `↑↓ move · tab switch · enter detail · space sign off · 1 2 3 order · c served · r reload · q quit` |
| **Backlog** · **Path** · **Board** | `↑↓ move · tab switch · enter detail · r reload · q quit` |
| **Log** | `↑↓ move · tab switch · enter detail · r reload · q quit` — the entry opens |
| **Sous-Chef** | `↑↓ move · tab switch · r reload · q quit` — **no** `enter detail` |
| **Recipes with nothing selected** | no `enter detail` |
| **Settings** | `enter apply theme` — never `enter detail` |

**One answer.** All seven tabs were checked together, by `TestTheKeyAndThePaneAgreeOnEveryTab`: for each tab the footer offers `enter detail` **exactly** where `hasDetail()` is true, and where it offers it the key opens the detail. The tab-by-tab promise: Recipes · Backlog · Path · Board · Log · Sous-Chef · Settings.

Two mutations, both reverted, both red:
- `hasDetail()` without the Log case → `TestEnterOpensTheLogEntry` and `TestTheFooterPromisesOnlyWhatItDoes` fail.
- the footer offering `enter detail` unconditionally → `TestTheFooterPromisesOnlyWhatItDoes` and `TestTheKeyAndThePaneAgreeOnEveryTab` fail.

**DomainVersion:** 7 — unmoved. `internal/core` and `internal/cli` are untouched: the seam needed no domain change.

**The goldens:** unchanged — `make goldens` leaves `git diff` empty. The footer's key list is rendered text, and the goldens hold the **Recipes** tab, where a recipe is selected and the footer's keys do not change.

**The gate:** `make check` green, `make check-race` green (core 89 s, ui 301 s), gofmt/vet/staticcheck clean, the budget unmoved (0 / 0 / 28 of 30). The tree is clean.

**Notes for the reviewer:**
- **`docs/05` is unchanged**, because no line of it became false: the doc never promised `enter detail` on every tab — the interaction table's `enter` row ("expand the detail to full screen; on **Settings**, apply the theme under the cursor") is still true, and the layout paragraph is about the split, not about a per-tab promise. The Sous-Chef's own pane message ("Nothing to open here — the Sous-Chef is not a list of recipes") is untouched, as ruled.
- The tab-by-tab test holds the **agreement** (footer ⟺ key), and the two facts — the Log *has* a detail, and the Sous-Chef and an empty list have none — are held by the two other tests. That split is deliberate: the agreement test cannot catch "the Log has no detail at all", because both sides would change together — and the Log's case is caught by the other two.
- The **Log's detail in the split layout** was already drawn before this cycle; what was missing was the key and the footer's promise. The narrow layout is where the difference shows, and that is where the evidence above was taken.
- The detail's arrow on a promoted item and the backlog item's absent lists stay for a later round, as ruled.
