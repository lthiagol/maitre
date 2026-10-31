Runner: T4 done. Commits: 4 (219b331, 79c7fb5, 36a1a35, dcaac47) + the state row (1814c2a). Per-AC evidence: AC-01..AC-06 stamped. Pre-notify gate: pass.

**The tab.** The console opens on **Home** — it is where the chef lands. The tab count is **7** and the names are `Home · Recipes · Backlog · Path · Board · Log · Settings`: the **Sous-Chef dissolved into Home**, and its content is Home's **automation** block (the plan's domain version · this build's · whether they agree · the last journal event). The Sous-Chef's "Not built yet" note is gone with the tab — the orchestrator's tab comes back when the orchestrator does, and what the chef needs from the automation today is its vital signs.

Evidence (AC-01): `internal/ui/home_test.go` — `TestTheConsoleLandsOnHome` (the landing tab, the count, the seven names) and `TestTheHomeTabShowsTheAutomation` (the vital signs, the last event). A pty capture at 120×40 (`scripts/console-shot.py`) shows the bar and the automation block:

```
▸ Home    Recipes    Backlog    Path    Board    Log    Settings
 THE AUTOMATION
  plan written under 7 · this build speaks 7 · they agree
  last event created     ·     chef      01:21:21
```

**The numbers.** Every row is on the screen, and **the console computes none of them**:

| Block | The numbers | Where the number comes from |
|---|---|---|
| The plan | name · path | `core.Project` (the plan's own name) and the plan's root — the console reads them, it does not count them |
| The recipes | total 7 · need you 2 · on the line 2 · served 1 | the triage's sections, counted by the domain (`core.TriageSections`) |
| The evidence | criteria 12 · proven 6 · failed 2 · rejected 0 · unproven 0 | `core.EvidenceSummaryOf`, summed over the recipes |
| The tasting | open notes 1 · open findings 1 | `core.OpenFindingsOf` — the notes of any phase, and the ones the **external** tasting filed |
| The backlog | ideas 4 · small tasks 2 | `core.BacklogGroups` — the item's kind puts it in one group or the other |
| The automation | plan written under 7 · this build speaks 7 · they agree | the journal (`core.PlanDomainVersion`) and `core.DomainVersion` |

Evidence (AC-02): `TestTheHomeTabShowsTheNumbers` asks `core.PlanStats` for the plan's numbers and looks for each one on the screen; `internal/core/stats_test.go` — `TestPlanStatsCountsThePlan` counts them in the domain.

**The domain's answer.** A **new view**: `core.PlanStats` (the type) and **`core.PlanStatsOf`** (the function) in `internal/core/stats.go` — the `// view` marker, a registration in `core.DerivedViews` under the name `plan_stats`, and a row in **`docs/03` § 8** whose name comes from the runtime. The wire checks both directions and **fails when either side goes**: `TestThePlanStatsViewIsHeldInBothDirections` proves it (the row relabelled → the table check trips; the registration removed → the coverage check trips). The plan's **name and path are not in the view** — they are stored facts, read from `core.Project` and the plan's root; a view holds what is computed.

**No items.** The invariant — *the Home tab shows numbers; it shows no items* — is held by `TestTheHomeTabShowsNoItems` (no row holds an item, the tab's item count is zero, the cursor has nowhere to land), `TestSpaceDoesNothingOnHome` (`space` writes **nothing** to the journal — the file is byte-for-byte the same) and `TestTheHomeFooterPromisesNoMove` (the footer promises no `space`, no `enter detail`, no cursor — only `tab`, `r`, `q`).

The mutation (AC-03): making a row on Home selectable — `itemIdx: -1` → `itemIdx: 0` in `internal/ui/home.go` — turns the suite **red**:

```
--- FAIL: TestTheHomeTabShowsNoItems (0.89s)
    home_test.go:152: row 0 holds the item 0 — Home shows numbers, not items
    ... (18 rows)
FAIL	github.com/lthiagol/maitre/internal/ui
```

**The art.** The mark moved from the Sous-Chef to Home and sits **beside** the numbers — the point of moving it is that it can share the space. The decoration yields to the content in both dimensions: the numbers keep their column and the mark is drawn only in what is left (the **full** form, or none), and the numbers must be whole on the screen.

At **80×24** (pty capture) the numbers are whole and **the mark is absent** — the tab is what it would be without the art:

```
 THE EVIDENCE
  criteria 12 · proven 6 · failed 2 · rejected 0 · unproven 0
 THE TASTING
  open notes 1 · open findings 1
```

At **120×40** the mark is drawn beside them, on their rows, and the numbers are not pushed down:

```
 THE EVIDENCE                                                  ⣼⣿⣿⣶⣶⣦⣤⣀⡀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣀⣤⣴⣶⣶⣿⣿⣧
  criteria 12 · proven 6 · failed 2 · rejected 0 · unproven 0  ⢸⣿⣿⣿⣿⣏⠙⠛⠻⠇⢰⣿⣿⣿⣿⣿⣿⡆⠸⠟⠛⠋⣹⣿⣿⣿⣿⡇
```

Evidence (AC-04): `TestTheHomeKeepsItsNumbersAt80x24`, `TestTheHomeDrawsTheMarkBesideTheNumbers`, `TestTheHomeArtNeedsBothDimensions` (the width and the height, each on its own).

**The goldens.** `make goldens` leaves a diff, and `make goldens-check` is green with the change claimed:

| File | What moved |
|---|---|
| `console-recipes.golden` | the tab bar only — `Home` is first, the Sous-Chef is gone |
| `console-stress.golden` | the tab bar only |
| `console-home.golden` | **new** — the Home tab at 80 columns, the classic default: the numbers whole and the mark absent |

Two things in the new golden are the machine's own and are not pinned: the plan's path (the fixture lives under the home directory, so the golden holds the console's short form, `~/demo/kitchen`) and the clock (`HH:MM:SS`, as wide as the time it replaces).

**The gate.** `make check` green, `make check-race` green, the tree clean. `core.DomainVersion` is **7** — no format, no field, no directory. No dependency added: `internal/core` 0, `internal/cli` 0, `internal/ui` 28 (ceiling 30).

**Notes for T5.** The tab order is `Home · Recipes · Backlog · Path · Board · Log · Settings` and the count is 7 — T5 renames **Path → Ready** and the count does not move. The console's landing tab is Home, so **every test that means another tab must say so**: `internal/ui/console_test.go` has `recipesModel(t)` for the Recipes tab, and the model builders in the other files switch explicitly. T5's `core.Adoptable` is a view like `core.PlanStats` — the `// view` marker, a registration in `core.DerivedViews`, a row in `docs/03` § 8, and the wire checked both ways. The Home tab renders **full width** even at 160 columns (`renderLayout` gives it the single layout): it has no detail, and a pane beside it would hold nothing.
