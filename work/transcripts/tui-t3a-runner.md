# The runner's lane — T3a, captured 2026-10-30

Dumped from pane `w1H:pA` with `herdr pane read --source recent-unwrapped --lines 5000`.
The lane that wrote T3a: the backlog's detail, the refusal that is the domain's, the
shape every writer writes - and two fix cycles after the review.
The review that judged it is in [tui-t3a-reviewer.md](tui-t3a-reviewer.md).

    ┃

    ┃
    ┃  ← Edit internal/ui/backlog_test.go
    ┃
    ┃     1   package ui                                                                                                                                          1   package ui
    ┃     2                                                                                                                                                       2
    ┃     3   import (                                                                                                                                            3   import (
    ┃                                                                                                                                                             4 +   "crypto/sha256"
    ┃                                                                                                                                                             5 +   "encoding/hex"
    ┃                                                                                                                                                             6 +   "os"
    ┃                                                                                                                                                             7 +   "path/filepath"
    ┃     4     "strings"                                                                                                                                         8     "strings"
    ┃     5     "testing"                                                                                                                                         9     "testing"
    ┃     6                                                                                                                                                      10
    ┃     7     "charm.land/bubbletea/v2"                                                                                                                        11     "charm.land/bubbletea/v2"
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && gofmt -l . ; go test ./internal/ui -run 'Backlog|Promot|Space|Footer' -v 2>&1 | grep -E "^(=== RUN|--- |ok|FAIL|    )" | head -40
    ┃
    ┃  (26 earlier lines)
    ┃  === RUN   TestTheFooterNamesTheNoteAndTheTriage
    ┃  --- PASS: TestTheFooterNamesTheNoteAndTheTriage (0.11s)
    ┃  === RUN   TestTheFooterIsPinnedToTheBottom
    ┃  --- PASS: TestTheFooterIsPinnedToTheBottom (0.69s)
    ┃  === RUN   TestTheFooterIsContextual
    ┃  --- PASS: TestTheFooterIsContextual (0.64s)
    ┃  ok    github.com/lthiagol/maitre/internal/ui  11.417s
    ┃

       + Thought · 530ms



    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && cp internal/ui/views.go /private/var/folders/5l/jx0d52nd65v87ltmtrwwtf5h0000gn/T/opencode/views.go.keep && python3 - <<'EOF'
    ┃  p = "internal/ui/views.go"…
    ┃
    ┃  (11 earlier lines)
    ┃            BL03  The chef wants a weekly digest               idea    open
    ┃            BL04  Trying the board on a real project           idea    open
    ┃           SMALL TASKS                                                                                    2
    ┃          ──────────────────────────────────────────────────────────────────────────────────────────────────
    ┃            BL05  Wire the review queue                        feature promoted
    ┃            BL06  Sharpen the evidence contract                bug     open
    ┃

       + Thought · 503ms

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && go test ./internal/ui 2>&1 | grep -c "^--- FAIL"; go test ./internal/ui 2>&1 | grep "^--- FAIL"; cp /private/var/folders/5l/jx0d52nd65v87ltmtrwwtf5h0000gn/T/opencode/views.go.keep internal/ui/views.go && git diff --stat && git status --short && echo "(tree clean —
    ┃  mutation reverted)"
    ┃
    ┃  (2 earlier lines)
    ┃   internal/ui/backlog_test.go | 79 +++++++++++++++++++++++++++++++++++++++++++--
    ┃   internal/ui/model.go        | 16 +++++++++
    ┃   2 files changed, 93 insertions(+), 2 deletions(-)
    ┃   M internal/ui/backlog_test.go
    ┃   M internal/ui/model.go
    ┃  (tree clean — mutation reverted)
    ┃

       + Thought · 954ms

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && make build 2>&1 | tail -2 && cd /private/var/folders/5l/jx0d52nd65v87ltmtrwwtf5h0000gn/T/opencode/t3a && rm -rf fix1 && /Users/thiago/Code/Pessoal/maitre/bin/maitre-cli demo --dir fix1 >/dev/null && cd fix1 && echo "=== F-115: activity.json before ===" && shasum -a
    ┃   256 kitchen/activity.json && wc -c < kitchen/activity.json
    ┃
    ┃  go build -o bin/maitre-cli ./cmd/maitre-cli
    ┃  go build -o bin/maitre ./cmd/maitre
    ┃  === F-115: activity.json before ===
    ┃  948239b69b273e9b3df3ae7c0e744b60f37c25f4cac71fd5894d4dd44527bb10  kitchen/activity.json
    ┃     16418
    ┃

    ┃
    ┃  $ cd /private/var/folders/5l/jx0d52nd65v87ltmtrwwtf5h0000gn/T/opencode/t3a/fix1 && M=/Users/thiago/Code/Pessoal/maitre/bin/maitre && echo "=== space on an idea ===" && $M --render --plain --keys "tab,space" 2>&1 | sed -n '1,10p' && echo && echo "=== activity.json after ===" && shasum -a 256 kitchen/activit
    ┃  y.json && wc -c < kitchen/activity.json
    ┃
    ┃  (9 earlier lines)
    ┃    BL01  The board needs a way to show evidence       idea    open
    ┃    BL02  A recipe with no ACs is probably an idea     idea    open
    ┃
    ┃  === activity.json after ===
    ┃  948239b69b273e9b3df3ae7c0e744b60f37c25f4cac71fd5894d4dd44527bb10  kitchen/activity.json
    ┃     16418
    ┃

       + Thought · 644ms

    ┃
    ┃  $ cd /private/var/folders/5l/jx0d52nd65v87ltmtrwwtf5h0000gn/T/opencode/t3a/fix1 && M=/Users/thiago/Code/Pessoal/maitre/bin/maitre-cli && echo "=== the CLI's refusal stands ===" && $M backlog promote BL01; echo "exit=$?"; echo "--- and the journal records the attempt ---" && python3 -c "
    ┃  import json…
    ┃
    ┃  (6 earlier lines)
    ┃    "ok": false
    ┃  }
    ┃  exit=1
    ┃  --- and the journal records the attempt ---
    ┃  58 created None {"id": "BL06", "item_kind": "bug", "kind": "backlog", "title": "Sharpen the evidence contract"}
    ┃  59 refusal BL01 {"reason": "\"idea\" carries no acceptance criteria, and a recipe without them can never be signed off \u2014 shape it i
    ┃

       + Thought · 1.6s

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && grep -n "F-114\|F-115" -A 3 docs/audit.md | head -40; echo "--- notes ---"; grep -n "Notes" -A 12 docs/audit.md | tail -30
    ┃
    ┃  No tool named "bash" is currently available. Please use a tool from the available tool list.
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && grep -n "F-114\|F-115" docs/audit.md | head -10; echo "--- notes table ---"; grep -n "^| N-" docs/audit.md | head -10
    ┃
    ┃  (1 earlier line)
    ┃  he suite** — the promotion of an idea made to succeed (core, cli and ui), the writer's normalisation taken out, `PromotionOf` answering for every kind (core, cli and ui), and the view's registration or its `docs/03` row removed (the wire, both ways); **the fifth did not**: the console deciding
    ┃  promotability **for itself** leaves the whole `internal/ui` suite **green** — **two blockers filed** (F-114, F-115). The scope held: `internal/core` and `internal/cli` stay stdlib-only, `core.DomainVersion` is still **7**, the budget is still **0 / 0 / 28**, the seven tabs kept their names, `maitre-cli`
    ┃  gained no verb, `make goldens` leaves `git diff` empty, and `docs/03` and `docs/07` no longer say an idea can be promoted straight to a recipe — **verdict: changes-needed** |
    ┃  123:| **F-114** | medium | **The footer's guard does not bite — the console may decide promotability for itself and the suite stays green.** `internal/ui/views.go` asks `core.PromotionOf` for the footer, but **nothing holds it**: the reviewer's mutation — the footer deciding by the item's kind instead of
    ┃  asking the domain — leaves the **whole `internal/ui` suite green**. The two cases the test holds (an idea, a small task) are the two where the console's own rule and the domain's answer happen to **agree**; the case where they differ — a **promoted** small task, which the domain refuses because it is not
    ┃  open — is untested, and the mutated console then offers `space promote to recipe` on it (proved by command: `BL05` promoted → the mutated build offers the move, the real build does not). *One fact, one answer* is a house rule, and nothing here holds it. | `internal/ui/views.go:523` — `primaryActionLabel`;
    ┃  `internal/ui/backlog_test.go:91` — `TestTheFooterDoesNotOfferAMoveTheDomainRefuses` |
    ┃  124:| **F-115** | medium | **A refused `space` on the Backlog writes to the plan's journal — a misclick writes to the record.** The console asks the domain by **calling the operation**, and the operation journals every refusal (`refuseLocked`), so pressing `space` on an idea appends a `refusal` event to
    ┃  `activity.json` — the plan's only event log, append-only, and there is no way to take it back — while the footer, in the same frame, offers **no move at all**. The journal's own rule is that a refusal is *"an attempt that was turned away"*: nobody attempted anything, so the record says the chef tried to
    ┃  promote an idea. The console's own precedent is the recipe path, which asks `core.ChefMoveOf` and does **not** call the operation when there is no move. The remedy is the answer the footer already reads: `core.PromotionOf` — the console shows the domain's reason and writes nothing. **The CLI's refusal
    ┃  stays**: an agent asked, and the plan said no. | `internal/ui/model.go:852` — `promoteSelected`; `internal/core/operations.go:939` — `refuseLocked` |
    ┃  --- notes table ---
    ┃

       + Thought · 511ms

    ┃
    ┃  ← Edit GLOSSARY.md
    ┃
    ┃    57   The plan's loose end — one artifact with a **kind**: an `idea` (a hunch, an                                                                        57   The plan's loose end — one artifact with a **kind**: an `idea` (a hunch, an
    ┃    58   observation, an agent discovery), a `feature` or a `bug` (a small task, ready to work).                                                            58   observation, an agent discovery), a `feature` or a `bug` (a small task, ready to work).
    ┃    59   It replaces the two file types the plan used to hold — see                                                                                         59   It replaces the two file types the plan used to hold — see
    ┃    60   [ADR 0009](docs/decisions/0009-vocabulary.md) § 3 — and the **kind** is the only truth                                                             60   [ADR 0009](docs/decisions/0009-vocabulary.md) § 3 — and the **kind** is the only truth
    ┃    61 - about its shape. It becomes a recipe when it is promoted, and shaping an idea into a                                                               61 + about its shape. A **small task** becomes a recipe when it is promoted, and shaping an
    ┃    62 - small task is an update, not a promotion: the id does not change.                                                                                  62 + idea into a small task is an update, not a promotion: the id does not change.
    ┃    63   _Avoid_: Special, prep note, note, brief, ticket                                                                                                   63   _Avoid_: Special, prep note, note, brief, ticket
    ┃    64                                                                                                                                                      64
    ┃    65   **Course**:                                                                                                                                        65   **Course**:
    ┃    66   One intent across several recipes — a themed set of dishes served together.                                                                        66   One intent across several recipes — a themed set of dishes served together.
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git diff --stat && gofmt -l . ; go vet ./... && go test ./internal/ui ./internal/core ./internal/cli 2>&1 | tail -6
    ┃
    ┃  GLOSSARY.md                 |  4 +--
    ┃   internal/ui/backlog_test.go | 79 +++++++++++++++++++++++++++++++++++++++++++--
    ┃   internal/ui/model.go        | 16 +++++++++
    ┃   3 files changed, 95 insertions(+), 4 deletions(-)
    ┃  ok    github.com/lthiagol/maitre/internal/ui  123.901s
    ┃  ok    github.com/lthiagol/maitre/internal/core  24.112s
    ┃  ok    github.com/lthiagol/maitre/internal/cli  4.940s
    ┃

       + Thought · 487ms

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git diff GLOSSARY.md
    ┃
    ┃  (8 earlier lines)
    ┃  -about its shape. It becomes a recipe when it is promoted, and shaping an idea into a
    ┃  -small task is an update, not a promotion: the id does not change.
    ┃  +about its shape. A **small task** becomes a recipe when it is promoted, and shaping an
    ┃  +idea into a small task is an update, not a promotion: the id does not change.
    ┃   _Avoid_: Special, prep note, note, brief, ticket
    ┃
    ┃   **Course**:
    ┃

       Now the commits — one per finding:

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git add internal/ui/model.go && git commit -q -m "ui: a refused space writes nothing to the plan
    ┃  …
    ┃
    ┃  aa4ba5f ui: a refused space writes nothing to the plan
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git add internal/ui/backlog_test.go && git commit -q -m "ui: a test holds that the footer asks the domain
    ┃  …
    ┃
    ┃  99a7576 docs: the glossary says a small task becomes a recipe
    ┃  da36146 ui: a test holds that the footer asks the domain
    ┃  aa4ba5f ui: a refused space writes nothing to the plan
    ┃  dacfa6f work: T3a's fix cycle 1 - the footer that is not held, and the refusal that writes
    ┃

       + Thought · 1.6s

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git show dacfa6f -- work/tui/README.md
    ┃
    ┃  (34 earlier lines)
    ┃   | **T2a** | the mouse, part 1 | the mouse is on; the render builds the hit map; the wheel scrolls | ✅ closed · 12 commits · a72ea98, 115b650, 284fb65, 0672f9f, 9f971cb, 7035ec7, e4b7061, 619b342, 47be828, f1cc1ca, a332158, d8d7b95 | [tui-t2a.md](../briefs/tui-t2a.md) · [fix 1](../briefs/fix-tui-t2a-1.md)
    ┃  · [fix 2](../briefs/fix-tui-t2a-2.md) |
    ┃   | **T2b** | the mouse, part 2 | the clicks: a tab, a row, a detail; the mouse drivable from a script | ✅ closed · 9 commits · 4dfd5a4, 187d706, f64e6a2, ac8df77, 588cf5f, 9c716ce, 0144216, 6bb7545, f8969ef | [tui-t2b.md](../briefs/tui-t2b.md) · [fix 1](../briefs/fix-tui-t2b-1.md) |
    ┃  -| **T3a** | the backlog is readable | the detail pane; `space` stops being a trap on an idea; D-3, D-5 | ✅ closed · 5 commits · 4dc8079, 1b550c6, 1804ccc, b873821, 5edae71 | [tui-t3a.md](../briefs/tui-t3a.md) |
    ┃  +| **T3a** | the backlog is readable | the detail pane; `space` stops being a trap on an idea; D-3, D-5 | 🔄 **review ok · fix cycle 1** — 5 commits · 4dc8079, 1b550c6, 1804ccc, b873821, 5edae71 | [tui-t3a.md](../briefs/tui-t3a.md) · [fix 1](../briefs/fix-tui-t3a-1.md) |
    ┃   | **T3b** | the backlog moves | dismiss · create · edit | ⏳ planned | — |
    ┃   | **T4** | the Home tab | the landing tab; the art; the Sous-Chef dissolves; D-2 | ⏳ planned | — |
    ┃   | **T5** | the Ready tab | `core.Adoptable`; the Path becomes Ready; D-1, D-4 | ⏳ planned | — |
    ┃

       + Thought · 1.6s

    ┃
    ┃  ← Edit work/tui/README.md
    ┃
    ┃    14   |--:|---|---|---|---|                                                                                                                              14   |--:|---|---|---|---|
    ┃    15   | **T1** | the identity | the window title carries the project; the header's name stops being a footnote; `--project`; the Settings tab shows the  15   | **T1** | the identity | the window title carries the project; the header's name stops being a footnote; `--project`; the Settings tab shows the
    ┃         plan's path | ✅ closed · 11 commits · 94a53ff, 8fc7438, f0b2aa5, 82289c6, 0bbcb35, 5a79638, 67b455d, b4a3052, 55bdc81, 8a6866b, cf8c05e | [tui-        plan's path | ✅ closed · 11 commits · 94a53ff, 8fc7438, f0b2aa5, 82289c6, 0bbcb35, 5a79638, 67b455d, b4a3052, 55bdc81, 8a6866b, cf8c05e | [tui-
    ┃         t1.md](../briefs/tui-t1.md) · [fix 1](../briefs/fix-tui-t1-1.md) · [fix 2](../briefs/fix-tui-t1-2.md) |                                                 t1.md](../briefs/tui-t1.md) · [fix 1](../briefs/fix-tui-t1-1.md) · [fix 2](../briefs/fix-tui-t1-2.md) |
    ┃    16   | **T2a** | the mouse, part 1 | the mouse is on; the render builds the hit map; the wheel scrolls | ✅ closed · 12 commits · a72ea98, 115b650,     16   | **T2a** | the mouse, part 1 | the mouse is on; the render builds the hit map; the wheel scrolls | ✅ closed · 12 commits · a72ea98, 115b650,
    ┃         284fb65, 0672f9f, 9f971cb, 7035ec7, e4b7061, 619b342, 47be828, f1cc1ca, a332158, d8d7b95 | [tui-t2a.md](../briefs/tui-t2a.md) · [fix 1](../briefs/      284fb65, 0672f9f, 9f971cb, 7035ec7, e4b7061, 619b342, 47be828, f1cc1ca, a332158, d8d7b95 | [tui-t2a.md](../briefs/tui-t2a.md) · [fix 1](../
    ┃         fix-tui-t2a-1.md) · [fix 2](../briefs/fix-tui-t2a-2.md) |                                                                                               briefs/fix-tui-t2a-1.md) · [fix 2](../briefs/fix-tui-t2a-2.md) |
    ┃    17   | **T2b** | the mouse, part 2 | the clicks: a tab, a row, a detail; the mouse drivable from a script | ✅ closed · 9 commits · 4dfd5a4, 187d706,   17   | **T2b** | the mouse, part 2 | the clicks: a tab, a row, a detail; the mouse drivable from a script | ✅ closed · 9 commits · 4dfd5a4, 187d706,
    ┃         f64e6a2, ac8df77, 588cf5f, 9c716ce, 0144216, 6bb7545, f8969ef | [tui-t2b.md](../briefs/tui-t2b.md) · [fix 1](../briefs/fix-tui-t2b-1.md) |              f64e6a2, ac8df77, 588cf5f, 9c716ce, 0144216, 6bb7545, f8969ef | [tui-t2b.md](../briefs/tui-t2b.md) · [fix 1](../briefs/fix-tui-t2b-1.md) |
    ┃    18 - | **T3a** | the backlog is readable | the detail pane; `space` stops being a trap on an idea; D-3, D-5 | 🔄 **review ok · fix cycle 1** — 5        18 + | **T3a** | the backlog is readable | the detail pane; `space` stops being a trap on an idea; D-3, D-5 | 🔄 **implemented · fix cycle 1** — 8
    ┃         commits · 4dc8079, 1b550c6, 1804ccc, b873821, 5edae71 | [tui-t3a.md](../briefs/tui-t3a.md) · [fix 1](../briefs/fix-tui-t3a-1.md) |                      commits · 4dc8079, 1b550c6, 1804ccc, b873821, 5edae71, aa4ba5f, da36146, 99a7576 | [tui-t3a.md](../briefs/tui-t3a.md) · [fix 1](../briefs/fix-
    ┃                                                                                                                                                                 tui-t3a-1.md) |
    ┃    19   | **T3b** | the backlog moves | dismiss · create · edit | ⏳ planned | — |                                                                         19   | **T3b** | the backlog moves | dismiss · create · edit | ⏳ planned | — |
    ┃    20   | **T4** | the Home tab | the landing tab; the art; the Sous-Chef dissolves; D-2 | ⏳ planned | — |                                                20   | **T4** | the Home tab | the landing tab; the art; the Sous-Chef dissolves; D-2 | ⏳ planned | — |
    ┃    21   | **T5** | the Ready tab | `core.Adoptable`; the Path becomes Ready; D-1, D-4 | ⏳ planned | — |                                                   21   | **T5** | the Ready tab | `core.Adoptable`; the Path becomes Ready; D-1, D-4 | ⏳ planned | — |
    ┃    22   | **T6** | the Board, as a visual place | the tree drawn left to right; clicking a node | ⏳ planned | — |                                         22   | **T6** | the Board, as a visual place | the tree drawn left to right; clicking a node | ⏳ planned | — |
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git add work/tui/README.md && git commit -q -m "work: T3a's row names the fix cycle's three commits" && git log --oneline -5 && git status --short
    ┃
    ┃  8f39289 work: T3a's row names the fix cycle's three commits
    ┃  99a7576 docs: the glossary says a small task becomes a recipe
    ┃  da36146 ui: a test holds that the footer asks the domain
    ┃  aa4ba5f ui: a refused space writes nothing to the plan
    ┃  dacfa6f work: T3a's fix cycle 1 - the footer that is not held, and the refusal that writes
    ┃

       Now the gate:

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && make check 2>&1 | tail -20 && git status --short && echo "(tree clean)"
    ┃
    ┃  (14 earlier lines)
    ┃    internal/core: 0 external modules (ceiling 0)
    ┃    internal/cli: 0 external modules (ceiling 0)
    ┃    internal/ui: 28 external modules (ceiling 30)
    ┃  /Library/Developer/CommandLineTools/usr/bin/make tidy-check
    ┃
    ┃  ==> tidy-check
    ┃  (tree clean)
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && make goldens >/dev/null 2>&1; git status --short; git diff --stat; echo "(goldens unchanged)" ; make check-race 2>&1 | tail -14
    ┃
    ┃  (8 earlier lines)
    ┃  GO=go scripts/dep-budget.sh
    ┃    internal/core: 0 external modules (ceiling 0)
    ┃    internal/cli: 0 external modules (ceiling 0)
    ┃    internal/ui: 28 external modules (ceiling 30)
    ┃  /Library/Developer/CommandLineTools/usr/bin/make tidy-check
    ┃
    ┃  ==> tidy-check
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && make check-race 2>&1 | grep -E "^(ok|FAIL|---)" | head -10; git status --short; echo "(tree clean)"
    ┃
    ┃  ok    github.com/lthiagol/maitre/internal/brand  1.521s
    ┃  ok    github.com/lthiagol/maitre/internal/cli  9.823s
    ┃  ok    github.com/lthiagol/maitre/internal/core  90.658s
    ┃  ok    github.com/lthiagol/maitre/internal/ui  318.064s
    ┃  ok    github.com/lthiagol/maitre/internal/ui  5.174s
    ┃  (tree clean)
    ┃

       + Thought · 1.3s

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && make build >/dev/null 2>&1 && cd /private/var/folders/5l/jx0d52nd65v87ltmtrwwtf5h0000gn/T/opencode/t3a && rm -rf fix1b && /Users/thiago/Code/Pessoal/maitre/bin/maitre-cli demo --dir fix1b >/dev/null && cd fix1b && echo "before:" && shasum -a 256 kitchen/activity.js
    ┃  on && stat -f "%z bytes" kitchen/activity.json && /Users/thiago/Code/Pessoal/maitre/bin/maitre --render --plain --keys "tab,space" >/dev/null 2>&1 && echo "after:" && shasum -a 256 kitchen/activity.json && stat -f "%z bytes" kitchen/activity.json && echo && echo "the reason still lands:" && /Users/t…
    ┃
    ┃  (5 earlier lines)
    ┃  16418 bytes
    ┃
    ┃  the reason still lands:
    ┃   off — shape it into a small task first
    ┃   IDEAS                                                                                          4
    ┃  ──────────────────────────────────────────────────────────────────────────────────────────────────
    ┃

       + Thought · 607ms

    ┃
    ┃  $ cd /private/var/folders/5l/jx0d52nd65v87ltmtrwwtf5h0000gn/T/opencode/t3a/fix1b && /Users/thiago/Code/Pessoal/maitre/bin/maitre --render --plain --keys "tab,space" 2>&1 | sed -n '1,8p'
    ┃
    ┃  ⋈ maitre · fix1b
    ┃  ────────────────────────────────────────────────────────────────────────────────────────────────────
    ┃    Recipes  ▸ Backlog    Path    Board    Log    Sous-Chef    Settings
    ┃   6 items
    ┃   ⛔ refused — "idea" carries no acceptance criteria, and a recipe without them can never be signed
    ┃   off — shape it into a small task first
    ┃   IDEAS                                                                                          4
    ┃  ──────────────────────────────────────────────────────────────────────────────────────────────────
    ┃

    ┃
    ┃  $ mkdir -p /Users/thiago/Code/Pessoal/maitre/work/inbox && cat > /Users/thiago/Code/Pessoal/maitre/work/inbox/runner-tui-t3a-fix1.md <<'EOF'
    ┃  Runner: T3a fix 1 done. Commits: 4 (aa4ba5f, da36146, 99a7576, 8f39289) — 3 code/doc commits and the state row. Fixed: F-114, F-115.…
    ┃
    ┃  322b280 inbox: runner reports tui-t3a fix 1
    ┃  8f39289 work: T3a's row names the fix cycle's three commits
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && herdr agent prompt "w1H:p1" "Runner: T3a fix 1 done. Commits: 4 (aa4ba5f, da36146, 99a7576, 8f39289) -- 3 code/doc commits and the state row. Fixed: F-114, F-115. The footer's guard: the test uses the PROMOTED SMALL TASK -- an item that was a feature, was promoted,
    ┃   and is no longer open. It is the case where the console's own rule ('is it a small task?') and the domain's answer differ: the kind rule says yes, core.PromotionOf says no, because the item is not open any more. The two cases the old test held -- an idea and an open small task -- are the two where …
    ┃
    ┃  {"id":"cli:agent:prompt","result":{"agent":{"agent":"opencode","agent_session":{"agent":"opencode","kind":"id","source":"herdr:opencode","value":"ses_e737d6341ffeTdoEYTgykU25M8"},"agent_status":"done","completion_seq":309,"cwd":"/Users/thiago/Code/Pessoal/maitre","focused":false,"foreground_cwd":"/Users/
    ┃  thiago/Code/Pessoal/maitre","pane_id":"w1H:p1","revision":15,"screen_detection_skipped":true,"state_change_seq":309,"tab_id":"w1H:t1","terminal_id":"term_65f0cd2aa89d26","terminal_title":"OC | Project work proposal","terminal_title_stripped":"OC | Project work proposal","workspace_id":"w1H"},"type":
    ┃  "agent_prompted"}}
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git log --oneline -8 && echo "--- tree ---" && git status --short && echo "(clean)"
    ┃
    ┃  (3 earlier lines)
    ┃  da36146 ui: a test holds that the footer asks the domain
    ┃  aa4ba5f ui: a refused space writes nothing to the plan
    ┃  dacfa6f work: T3a's fix cycle 1 - the footer that is not held, and the refusal that writes
    ┃  0955ca4 inbox: reviewer reports tui-t3a
    ┃  4b40f8c docs: T3a's review - the ledger's two blockers and the slice's row
    ┃  --- tree ---
    ┃  (clean)
    ┃

       T3a's fix cycle 1 is done — 4 commits (aa4ba5f, da36146, 99a7576, 8f39289), receipt committed (322b280), orchestrator nudged.

       F-114 — the footer's guard bites now

       The case: a promoted small task. It was a feature, it was promoted, and it is no longer open. The console's own rule ("is it a small task?") says yes; core.PromotionOf says no. The two cases the old test held — an idea and an open small task — are exactly the two where the two answers agree, which is why a
       console deciding for itself passed them.

       TestTheFooterAsksTheDomainAboutPromotability holds the differing case: promote BL05, then assert the footer offers no space promote to recipe on it. The review's mutation (the footer deciding by item.Kind == core.KindIdea) now turns the suite red, and exactly one test fails:

       --- FAIL: TestTheFooterAsksTheDomainAboutPromotability
           backlog_test.go:178: the footer offers a promotion the domain refuses — the console decided for itself

       The mutation was reverted; git diff is empty.

       F-115 — the refusal writes nothing

       The console asks core.PromotionOf — a read — and shows the domain's reason; it no longer calls the operation, so nothing is journalled.

       before: 8960774117cca6aab1abaeb5ff8c494a0f0a7fbe6d7f5ebbcfc78666d7075f8c  (16418 bytes)
       after:  8960774117cca6aab1abaeb5ff8c494a0f0a7fbe6d7f5ebbcfc78666d7075f8c  (16418 bytes)

       Where the reason lands now: the banner (⛔ refused — "idea" carries no acceptance criteria…) and the item's own detail. Neither writes to the plan.

       The CLI's refusal stands — exit 1, the same reason, and the journal records the agent's attempt (seq 59, refusal, BL01), which is right: an agent asked.

       DomainVersion 7, unmoved; internal/core and internal/cli untouched. The goldens are unchanged. make check and make check-race are green, the tree is clean.

       One thing beyond the two blockers: GLOSSARY.md said "It becomes a recipe when it is promoted", which read as though an idea could — the review's LOW note (a). It now says "A small task becomes a recipe when it is promoted". The other two LOW notes — the detail's arrow on a promoted item, and the Log tab's
       detail being unreachable in the narrow layout — I left alone and named in the receipt, since the cycle does not widen and the second is the orchestrator's call.

       Build · MiniMax-M3.1-Flash · 21m 57s · 118.1 tok/s

    ┃
    ┃  You are the **runner**, again — the same lane that wrote T3a. The orchestrator's pane id
    ┃  is `w1H:p1`; use it verbatim in the notify command below.
    ┃
    ┃  This is **T3a's fix cycle 2**, and it is **small**: two things in the seam T3a
    ┃  introduced.
    ┃
    ┃  ## Why this cycle exists
    ┃
    ┃  Your fix cycle 1 asked a question, and the orchestrator has answered it: **the Log tab's
    ┃  detail rides in this cycle.** You said it was one line in `hasDetail` and one test.
    ┃  It is a little more than that, and the reason is worth writing down.
    ┃
    ┃  ## F-116 (medium) — the key and the pane disagree on the Log tab, and the footer promises a key that does nothing
    ┃
    ┃  **The target:** the console's **one** detail pane and the key that opens it read
    ┃  **one** answer, on every tab — and the footer never promises a key that does nothing.
    ┃
    ┃  **The hole.** T3a introduced `Model.hasDetail()` so that *"the pane and the key read
    ┃  the same answer, so the key can never open something the pane does not show"*. On the
    ┃  **Log** tab they do not:
    ┃
    ┃  | | The Log tab, today |
    ┃  |---|---|
    ┃  | the **pane** | shows the journal entry — `detailLines` has a `case TabLog:` and renders `m.logDetail(width)` |
    ┃  | the **key** | `enter` opens **nothing** — `hasDetail()` has no `case TabLog:` |
    ┃  | the **footer** | says **`enter detail`** — the default case adds that key on every tab but Settings |
    ┃
    ┃  So the footer promises a key that does nothing, and the pane shows something the key
    ┃  cannot open. The same promise is wrong on the **Sous-Chef**: the pane says *"Nothing to
    ┃  open here — the Sous-Chef is not a list of recipes"*, and the footer still says
    ┃  `enter detail`.
    ┃
    ┃  **What lands — two changes, both in the seam:**
    ┃
    ┃  1. **`hasDetail()` answers for the Log tab.** The pane already renders a journal
    ┃     entry, so the key must be able to open it. One case.
    ┃  2. **The footer offers `enter detail` only where there is a detail to open** — the
    ┃     same answer the key reads. So the console never promises a key that does nothing:
    ┃     not on the Log before (1), and not on the Sous-Chef, and not on a tab with nothing
    ┃     selected.
    ┃
    ┃  **The evidence:**
    ┃
    ┃  | | |
    ┃  |---|---|
    ┃  | the **Log** | `enter` on the Log tab opens the event's detail — in the **narrow** layout, where the detail is not beside the list. Show it by command: the frame changes, and it shows the event |
    ┃  | the **footer** | on the **Log** tab the footer offers `enter detail`; on the **Sous-Chef** it does **not**; and on the **Recipes** tab with **nothing selected** it does not |
    ┃  | **one answer** | the key and the pane agree on every tab — check them **together**, tab by tab |
    ┃
    ┃  ## What does *not* ride in this cycle
    ┃
    ┃  **The detail's arrow on a promoted item** — the one that reads like an instruction
    ┃  rather than a fact — **stays for a later round.** It wants a different marker for the
    ┃  domain's answer than the recipe detail's *"why it needs the chef"* arrow, and that is a
    ┃  question about the console's two arrow shapes, not a line.
    ┃
    ┃  **The backlog item's own absent lists** stay, as ruled — the note is in the ledger
    ┃  with the corrected reason.
    ┃
    ┃  ## The scope — what not to touch
    ┃
    ┃  | | |
    ┃  |---|---|
    ┃  | **The refusals** | **closed.** F-114 and F-115 are fixed and verified: the footer asks the domain and a test holds it, and a refused `space` writes nothing to the plan |
    ┃  | **The domain** | `internal/core` is untouched by this cycle. If the seam needs a domain change, **stop and report it** — do not reach for it |
    ┃  | **The Sous-Chef's pane** | its message — *"Nothing to open here"* — stays. This cycle changes the **footer's promise**, not that message |
    ┃  | **`internal/cli`** | unchanged |
    ┃  | **The tabs** | seven, the same names |
    ┃  | **`core.DomainVersion`** | stays **7** |
    ┃  | **The goldens** | expected unchanged — but **say so with the evidence**: the footer's key list is rendered text, and the goldens hold the **Recipes** tab, where the footer's keys do not change |
    ┃  | **`docs/05`** | unchanged unless a line became false — and if the footer's promise changes, check that the doc does not promise it either |
    ┃  | **`docs/audit.md`** | the ledger is the **reviewer's**. This cycle files no finding and touches no row |
    ┃  | **`work/`** | the orchestrator's, except the state row |
    ┃
    ┃  ## The invariants
    ┃
    ┃  - `internal/core` and `internal/cli` stay **stdlib-only**; `internal/ui` keeps its
    ┃    budget. **No dependency is added.**
    ┃  - **The domain layer has zero UI imports.**
    ┃  - **The tab count is 7** and the tab names do not change.
    ┃  - **One fact, one answer** — a fact the console shows is computed in the domain.
    ┃  - The gate is `make check`. A focused check is `make test PKG=… RUN=…`.
    ┃
    ┃  ## The report — the receipt comes FIRST
    ┃
    ┃  ```bash
    ┃  mkdir -p work/inbox
    ┃  cat > work/inbox/runner-tui-t3a-fix2.md <<EOF
    ┃  Runner: T3a fix 2 done. Commits: <n>. Fixed: F-116. The Log's detail: <what enter opens there now, and the command that shows it>. The footer's promise: <what it offers on the Log, on the Sous-Chef, and on Recipes with nothing selected>. One answer: <the tabs you checked, and that the key and the pane
    ┃  agree>. DomainVersion: <7>. The goldens: <unchanged / changed, and why>. The gate: <make check, make check-race>. Notes: <anything the reviewer should know>.
    ┃  EOF
    ┃  git add work/inbox && git commit -q -m "inbox: runner reports tui-t3a fix 2"
    ┃  herdr agent prompt "w1H:p1" "Runner: T3a fix 2 done. Commits: <n>. Fixed: F-116. The Log's detail: <what enter opens there now, and the command that shows it>. The footer's promise: <what it offers on the Log, on the Sous-Chef, and on Recipes with nothing selected>. One answer: <the tabs you checked, and
    ┃  that the key and the pane agree>. DomainVersion: <7>. The goldens: <unchanged / changed, and why>. The gate: <make check, make check-race>. Notes: <anything the reviewer should know>."
    ┃  ```
    ┃
    ┃  **Write the receipt, commit it, and *then* nudge** — the nudge only reaches an idle
    ┃  receiver; the receipt is the half that cannot be lost.
    ┃
    ┃  **One write per file, and commit as soon as the file is coherent.**
    ┃
    ┃  **And if you find a file you did not write has changed: stop and report it**, and do
    ┃  not touch it. It is most likely the reviewer's mutation, in flight.
    ┃
    ┃  ## Before you report
    ┃
    ┃  ```bash
    ┃  make check && git status --short
    ┃  ```
    ┃
    ┃  Green, and the tree clean. Report the **real** numbers: the commits, the gate's
    ┃  result, and what the footer promises on each tab.
    ┃
    ┃
    ┃  ---
    ┃
    ┃  ## The notify command — verbatim, use it as it is
    ┃
    ┃  The orchestrator's pane id is **w1H:p1**. Your receipt is the record; the nudge is the
    ┃  bell. Send **both**, in this order, and the nudge **last**:
    ┃
    ┃  ```bash
    ┃  mkdir -p work/inbox
    ┃  cat > work/inbox/runner-tui-t3a-fix2.md <<EOF
    ┃  Runner: T3a fix 2 done. Commits: <n>. Fixed: F-116. The Log's detail: <what enter opens there now, and the command that shows it>. The footer's promise: <what it offers on the Log, on the Sous-Chef, and on Recipes with nothing selected>. One answer: <the tabs you checked, and that the key and the pane
    ┃  agree>. DomainVersion: <7>. The goldens: <unchanged / changed, and why>. The gate: <make check, make check-race>. Notes: <anything the reviewer should know>.
    ┃  EOF
    ┃  git add work/inbox && git commit -q -m "inbox: runner reports tui-t3a fix 2"
    ┃  herdr agent prompt "w1H:p1" "Runner: T3a fix 2 done. Commits: <n>. Fixed: F-116. The Log's detail: <what enter opens there now, and the command that shows it>. The footer's promise: <what it offers on the Log, on the Sous-Chef, and on Recipes with nothing selected>. One answer: <the tabs you checked, and
    ┃  that the key and the pane agree>. DomainVersion: <7>. The goldens: <unchanged / changed, and why>. The gate: <make check, make check-race>. Notes: <anything the reviewer should know>."
    ┃  ```
    ┃
    ┃  **If the nudge fails** — it only reaches an idle receiver — the receipt is the half
    ┃  that cannot be lost. Write the receipt anyway, and say so in it.
    ┃
    ┃  **Nothing else waits on you.** After the nudge, end your turn. Do not poll for a
    ┃  reply, do not re-dispatch yourself, do not start another slice. The orchestrator
    ┃  decides what happens next.
    ┃
    ┃  ## What the orchestrator verified before writing this brief
    ┃
    ┃  So you do not repeat the work — fix cycle 1's two fixes were both reproduced:
    ┃
    ┃  - **F-115** — `maitre --render --plain --keys "tab,space"` on an idea leaves
    ┃    `activity.json` **byte-for-byte unchanged** (the same sha256 before and after), and
    ┃    the banner carries the domain's reason.
    ┃  - **F-114** — the review's mutation, the footer deciding by the item's kind, now turns
    ┃    `internal/ui` **red** in exactly one test: `TestTheFooterAsksTheDomainAboutPromotability`.
    ┃
    ┃  **And the seam, as it stands:** the pane renders a journal entry on the Log tab
    ┃  (`detailLines` has a `case TabLog:`), the key cannot open it (`hasDetail()` has no
    ┃  `case TabLog:`), and the footer promises `enter detail` on every tab but Settings.
    ┃

       + Thought · 14.3s

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git status --short && git log --oneline -2 && grep -n "func (m \*Model) openNotes" -A 12 internal/ui/note.go
    ┃
    ┃  (8 earlier lines)
    ┃  30-}
    ┃  31-
    ┃  32-// selectedNote is the tasting note under the cursor.
    ┃  33-func (m *Model) selectedNote() (core.Finding, bool) {
    ┃  34-  notes := m.openNotes()
    ┃  35-  if m.selNote < 0 || m.selNote >= len(notes) {
    ┃  36-    return core.Finding{}, false
    ┃

       + Thought · 888ms

    ┃
    ┃  ← Edit internal/ui/model.go
    ┃
    ┃    654                                                                                                                                                     654
    ┃    655   // hasDetail reports whether the current selection has a detail to open.                                                                          655   // hasDetail reports whether the current selection has a detail to open.
    ┃    656   //                                                                                                                                                656   //
    ┃    657   // The console has **one** detail pane, and it follows the tab: a recipe on Recipes,                                                              657   // The console has **one** detail pane, and it follows the tab: a recipe on Recipes,
    ┃    658 - // Path and Board, a backlog item on Backlog. The pane and the key read the same                                                                  658 + // Path and Board, a backlog item on Backlog, a journal entry on Log. The pane and
    ┃    659 - // answer, so the key can never open something the pane does not show — and a                                                                     659 + // the key read the same answer, so the key can never open something the pane does
    ┃    660 - // backlog item is not a recipe, so it needs its own answer here.                                                                                 660 + // not show — and a backlog item is not a recipe, so it needs its own answer here.
    ┃    661   func (m *Model) hasDetail() bool {                                                                                                                661   func (m *Model) hasDetail() bool {
    ┃    662     switch m.tab {                                                                                                                                  662     switch m.tab {
    ┃    663     case TabRecipes, TabPath, TabBoard:                                                                                                             663     case TabRecipes, TabPath, TabBoard:
    ┃    664       _, ok := m.currentRecipe()                                                                                                                    664       _, ok := m.currentRecipe()
    ┃    665       return ok                                                                                                                                     665       return ok
    ┃    666     case TabBacklog:                                                                                                                                666     case TabBacklog:
    ┃    667       _, ok := m.selectedBacklogItem()                                                                                                              667       _, ok := m.selectedBacklogItem()
    ┃    668       return ok                                                                                                                                     668       return ok
    ┃                                                                                                                                                            669 +   case TabLog:
    ┃                                                                                                                                                            670 +     events := m.logWindow()
    ┃                                                                                                                                                            671 +     return m.selLog >= 0 && m.selLog < len(events)
    ┃    669     }                                                                                                                                               672     }
    ┃    670     return false                                                                                                                                    673     return false
    ┃    671   }                                                                                                                                                 674   }
    ┃    672                                                                                                                                                     675
    ┃

    ┃
    ┃  ← Edit internal/ui/views.go
    ┃
    ┃    478     case m.tab == TabSettings:                                                                                                                      478     case m.tab == TabSettings:
    ┃    479       // The Settings tab has no recipe to open: `enter` applies the theme.                                                                         479       // The Settings tab has no recipe to open: `enter` applies the theme.
    ┃    480       keys = append(keys, key{"↑↓", "move"}, key{"tab", "switch"}, key{"enter", "apply theme"})                                                     480       keys = append(keys, key{"↑↓", "move"}, key{"tab", "switch"}, key{"enter", "apply theme"})
    ┃    481     default:                                                                                                                                        481     default:
    ┃    482 -     keys = append(keys, key{"↑↓", "move"}, key{"tab", "switch"}, key{"enter", "detail"})                                                          482 +     keys = append(keys, key{"↑↓", "move"}, key{"tab", "switch"})
    ┃                                                                                                                                                            483 +     // The footer never promises a key that does nothing. `enter` opens the
    ┃                                                                                                                                                            484 +     // detail, and there is one only where the pane has something to show —
    ┃                                                                                                                                                            485 +     // the same answer the key reads.
    ┃                                                                                                                                                            486 +     if m.hasDetail() {
    ┃                                                                                                                                                            487 +       keys = append(keys, key{"enter", "detail"})
    ┃                                                                                                                                                            488 +     }
    ┃    483     }                                                                                                                                               489     }
    ┃    484     if label := m.primaryActionLabel(); label != "" {                                                                                               490     if label := m.primaryActionLabel(); label != "" {
    ┃    485       keys = append(keys, key{"space", label})                                                                                                      491       keys = append(keys, key{"space", label})
    ┃    486     }                                                                                                                                               492     }
    ┃

       Now the tests:

    ┃
    ┃  # Wrote internal/ui/seam_test.go
    ┃
    ┃     1 package ui
    ┃     2
    ┃     3 import (
    ┃     4   "path/filepath"
    ┃     5   "strings"
    ┃     6   "testing"
    ┃     7
    ┃     8   "charm.land/bubbletea/v2"
    ┃     9
    ┃    10   "github.com/lthiagol/maitre/internal/core"
    ┃    11 )
    ┃    12
    ┃    13 // TestEnterOpensTheLogEntry covers the Log tab.
    ┃    14 //
    ┃    15 // The pane already renders the journal entry — `detailLines` has a `case TabLog:` —
    ┃    16 // so the key that opens the pane must open it too.
    ┃    17 func TestEnterOpensTheLogEntry(t *testing.T) {
    ┃    18   m := newModel(t)
    ┃    19   m.tab = TabLog
    ┃    20
    ┃    21   before := plain(m)
    ┃    22   if !strings.Contains(before, "recent activity") {
    ┃    23     t.Fatalf("the Log tab's summary is missing\n%s", before)
    ┃    24   }
    ┃    25
    ┃    26   m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
    ┃    27   if m.mode != ModeDetail {
    ┃    28     t.Fatalf("enter did not open the detail on the Log tab: mode is %v", m.mode)
    ┃    29   }
    ┃    30   after := plain(m)
    ┃    31   // The event's own detail: the data the journal carries. Nothing else on
    ┃    32   // this screen carries a "Data" section.
    ┃    33   if !strings.Contains(after, "Data") {
    ┃    34     t.Errorf("the open detail does not show the event\n%s", after)
    ┃    35   }
    ┃    36 }
    ┃    37
    ┃    38 // TestTheFooterPromisesOnlyWhatItDoes covers the footer's promise.
    ┃    39 //
    ┃    40 // A key that does nothing is a promise the console cannot keep. The footer offers
    ┃    41 // `enter detail` only where there is a detail to open — the Log's entry, the
    ┃    42 // backlog item's, a recipe's — and not where there is none: not on the Sous-Chef,
    ┃    43 // and not on a tab with nothing selected.
    ┃    44 func TestTheFooterPromisesOnlyWhatItDoes(t *testing.T) {
    ┃    45   m := newModel(t)
    ┃    46
    ┃    47   // The Log: the pane shows the entry, so the footer offers the key.
    ┃    48   m.tab = TabLog
    ┃    49   if out := plain(m); !strings.Contains(out, "enter detail") {
    ┃    50     t.Errorf("the footer does not offer the Log's detail\n%s", out)
    ┃    51   }
    ┃    52
    ┃    53   // The Sous-Chef: there is nothing to open.
    ┃    54   m.tab = TabSousChef
    ┃    55   if out := plain(m); strings.Contains(out, "enter detail") {
    ┃    56     t.Errorf("the footer promises a detail the Sous-Chef does not have\n%s", out)
    ┃    57   }
    ┃    58
    ┃    59   // The Recipes tab with nothing selected: nothing to open either.
    ┃    60   empty := emptyModel(t)
    ┃    61   empty.tab = TabRecipes
    ┃    62   if out := plain(empty); strings.Contains(out, "enter detail") {
    ┃    63     t.Errorf("the footer promises a detail with nothing selected\n%s", out)
    ┃    64   }
    ┃    65 }
    ┃    66
    ┃    67 // TestTheKeyAndThePaneAgreeOnEveryTab holds the seam.
    ┃    68 //
    ┃    69 // The console has **one** detail pane, and the key that opens it reads one answer:
    ┃    70 // the footer offers `enter detail` exactly where the key opens something, and the
    ┃    71 // key opens something exactly where the footer offers it. Tab by tab.
    ┃    72 func TestTheKeyAndThePaneAgreeOnEveryTab(t *testing.T) {
    ┃    73   m := newModel(t)
    ┃    74   for i := range TabNames {
    ┃    75     tab := Tab(i)
    ┃    76     m.tab = tab
    ┃    77     m.mode = ModeSplit
    ┃    78     m.clampSelections()
    ┃    79
    ┃    80     promised := strings.Contains(plain(m), "enter detail")
    ┃    81     if promised != m.hasDetail() {
    ┃    82       t.Errorf("%s: the footer promises the detail (%v) and there is one to open (%v)",
    ┃    83         tab.Label(), promised, m.hasDetail())
    ┃    84       continue
    ┃    85     }
    ┃    86     if !promised {
    ┃    87       continue
    ┃    88     }
    ┃    89     // And the key does what the footer says.
    ┃    90     m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
    ┃    91     if m.mode != ModeDetail {
    ┃    92       t.Errorf("%s: the footer promises the detail and the key does not open it", tab.Label())
    ┃    93     }
    ┃    94   }
    ┃    95 }
    ┃    96
    ┃    97 // emptyModel builds a console over a plan with nothing in it.
    ┃    98 func emptyModel(t *testing.T) *Model {
    ┃    99   t.Helper()
    ┃   100   store := core.Open(filepath.Join(t.TempDir(), core.PlanDirName))
    ┃   101   if err := store.Init(); err != nil {
    ┃   102     t.Fatalf("init: %v", err)
    ┃   103   }
    ┃   104   return loadModel(t, store)
    ┃   105 }
    ┃   106
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && gofmt -l . ; go test ./internal/ui -run 'Log|Footer|KeyAndThePane|Seam' -v 2>&1 | grep -E "^(=== RUN|--- |ok|FAIL|    )" | head -40
    ┃
    ┃  (16 earlier lines)
    ┃  === RUN   TestEnterOpensTheLogEntry
    ┃  --- PASS: TestEnterOpensTheLogEntry (1.00s)
    ┃  === RUN   TestTheFooterPromisesOnlyWhatItDoes
    ┃  --- PASS: TestTheFooterPromisesOnlyWhatItDoes (0.82s)
    ┃  === RUN   TestTheKeyAndThePaneAgreeOnEveryTab
    ┃  --- PASS: TestTheKeyAndThePaneAgreeOnEveryTab (0.91s)
    ┃  ok    github.com/lthiagol/maitre/internal/ui  10.077s
    ┃

       + Thought · 1.0s

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && cp internal/ui/model.go /private/var/folders/5l/jx0d52nd65v87ltmtrwwtf5h0000gn/T/opencode/model.go.keep && cp internal/ui/views.go /private/var/folders/5l/jx0d52nd65v87ltmtrwwtf5h0000gn/T/opencode/views.go.keep2 && python3 - <<'EOF'
    ┃  p = "internal/ui/model.go"…
    ┃
    ┃  mutation A applied: hasDetail has no Log case
    ┃  --- FAIL: TestEnterOpensTheLogEntry (0.97s)
    ┃  --- FAIL: TestTheFooterPromisesOnlyWhatItDoes (0.88s)
    ┃  FAIL
    ┃  FAIL  github.com/lthiagol/maitre/internal/ui  125.430s
    ┃  FAIL
    ┃

       + Thought · 1.8s

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && cp /private/var/folders/5l/jx0d52nd65v87ltmtrwwtf5h0000gn/T/opencode/model.go.keep internal/ui/model.go && python3 - <<'EOF'
    ┃  p = "internal/ui/views.go"…
    ┃
    ┃  mutation B applied: the footer always promises the detail
    ┃  --- FAIL: TestTheFooterPromisesOnlyWhatItDoes (0.81s)
    ┃  --- FAIL: TestTheKeyAndThePaneAgreeOnEveryTab (0.82s)
    ┃  FAIL
    ┃  FAIL  github.com/lthiagol/maitre/internal/ui  123.677s
    ┃  FAIL
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && cp /private/var/folders/5l/jx0d52nd65v87ltmtrwwtf5h0000gn/T/opencode/views.go.keep2 internal/ui/views.go && git diff --stat && git status --short && echo "(mutations reverted)" && go test ./internal/ui 2>&1 | tail -3
    ┃
    ┃  internal/ui/model.go | 9 ++++++---
    ┃   internal/ui/views.go | 8 +++++++-
    ┃   2 files changed, 13 insertions(+), 4 deletions(-)
    ┃   M internal/ui/model.go
    ┃   M internal/ui/views.go
    ┃  ?? internal/ui/seam_test.go
    ┃  (mutations reverted)
    ┃  ok    github.com/lthiagol/maitre/internal/ui  125.744s
    ┃

       + Thought · 1.2s

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && make build >/dev/null 2>&1 && cd /private/var/folders/5l/jx0d52nd65v87ltmtrwwtf5h0000gn/T/opencode/t3a && rm -rf fix2 && /Users/thiago/Code/Pessoal/maitre/bin/maitre-cli demo --dir fix2 >/dev/null && cd fix2 && M=/Users/thiago/Code/Pessoal/maitre/bin/maitre && echo
    ┃   "=== the Log tab: enter opens the event's detail ===" && $M --render --plain --keys "tab,tab,tab,tab,enter" 2>&1 | sed -n '1,20p'
    ┃
    ┃  (12 earlier lines)
    ┃        {
    ┃                "id": "BL06",
    ┃                "item_kind": "bug",
    ┃                "kind": "backlog",
    ┃                "title": "Sharpen the evidence contract"
    ┃              }
    ┃

    ┃
    ┃  $ cd /private/var/folders/5l/jx0d52nd65v87ltmtrwwtf5h0000gn/T/opencode/t3a/fix2 && M=/Users/thiago/Code/Pessoal/maitre/bin/maitre && echo "=== the Log tab's footer ===" && $M --render --plain --keys "tab,tab,tab,tab" 2>&1 | tail -2 && echo "=== the Sous-Chef's footer ===" && $M --render --plain --keys "tab
    ┃  ,tab,tab,tab,tab" 2>&1 | tail -2 && echo "=== the Recipes tab's footer (a recipe is selected) ===" && $M --render --plain 2>&1 | tail -2 && echo "=== the Settings tab's footer ===" && $M --render --plain --keys "tab,tab,tab,tab,tab,tab" 2>&1 | tail -2
    ┃
    ┃  (4 earlier lines)
    ┃   ↑↓ move · tab switch · r reload · q quit=== the Recipes tab's footer (a recipe is selected) ===
    ┃   ↑↓ move · tab switch · enter apply theme · r reload · q quit
    ┃
    ┃  === the Settings tab's footer ===
    ┃
    ┃   ↑↓ move · tab switch · enter apply theme · r reload · q quit
    ┃

       + Thought · 1.8s

       → Explored: 1 read

       + Thought · 1.4s

    ┃
    ┃  $ cd /private/var/folders/5l/jx0d52nd65v87ltmtrwwtf5h0000gn/T/opencode/t3a/fix2 && M=/Users/thiago/Code/Pessoal/maitre/bin/maitre && echo "=== the Recipes tab's footer (a recipe is selected) ===" && $M --render --plain --keys "tab,tab,tab,tab,tab,tab,tab" 2>&1 | tail -2 && echo && echo "=== the Backlog tab
    ┃  's footer ===" && $M --render --plain --keys "tab" 2>&1 | tail -2 && echo && echo "=== the Path tab's footer ===" && $M --render --plain --keys "tab,tab" 2>&1 | tail -2 && echo && echo "=== the Board tab's footer ===" && $M --render --plain --keys "tab,tab,tab" 2>&1 | tail -2
    ┃
    ┃  (6 earlier lines)
    ┃  === the Path tab's footer ===
    ┃
    ┃   ↑↓ move · tab switch · enter detail · r reload · q quit
    ┃  === the Board tab's footer ===
    ┃
    ┃   ↑↓ move · tab switch · enter detail · r reload · q quit
    ┃

       + Thought · 876ms

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && grep -n "enter detail\|enter\`" docs/05-human-interface.md | head -20; echo "--- Log tab in docs/05 ---"; grep -n -i "log tab\|journal entry" docs/05-human-interface.md | head -20
    ┃
    ┃  250:  list takes the whole terminal and `enter` opens the detail full screen. A split at 100
    ┃  259:| `enter` | expand the detail to full screen; on **Settings**, apply the theme under the cursor |
    ┃  279:| click | on the detail pane | open the detail — the keyboard's `enter` |
    ┃  295:`enter` applies the theme.
    ┃  311:**THEME** heading; the cursor moves over them and `enter` applies the one under it. The
    ┃  --- Log tab in docs/05 ---
    ┃

       → Explored: 1 read

       + Thought · 4.3s

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && make goldens >/dev/null 2>&1; git status --short; git diff --stat; echo "(the goldens are unchanged — the Recipes tab's footer still offers 'enter detail')"
    ┃
    ┃  M internal/ui/model.go
    ┃   M internal/ui/views.go
    ┃  ?? internal/ui/seam_test.go
    ┃   internal/ui/model.go | 9 ++++++---
    ┃   internal/ui/views.go | 8 +++++++-
    ┃   2 files changed, 13 insertions(+), 4 deletions(-)
    ┃  (the goldens are unchanged — the Recipes tab's footer still offers 'enter detail')
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git add internal/ui/model.go internal/ui/views.go internal/ui/seam_test.go && git commit -q -m "ui: the key and the pane agree on every tab
    ┃  …
    ┃
    ┃  63964ce ui: the key and the pane agree on every tab
    ┃  c2e5b7a work: T3a's fix cycle 2 - the key and the pane, on every tab
    ┃

    ┃
    ┃  ← Edit work/tui/README.md
    ┃
    ┃  Could not find oldString in /Users/thiago/Code/Pessoal/maitre/work/tui/README.md. It must match exactly, including whitespace and indentation.
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && grep -n "T3a" work/tui/README.md | head -5; ls work/briefs/ | grep t3a
    ┃
    ┃  18:| **T3a** | the backlog is readable | the detail pane; `space` stops being a trap on an idea; D-3, D-5 | 🔄 **fix cycle 2** — 7 commits · 4dc8079, 1b550c6, 1804ccc, b873821, 5edae71, aa4ba5f, da36146 | [tui-t3a.md](../briefs/tui-t3a.md) · [fix 1](../briefs/fix-tui-t3a-1.md) · [fix 2](../briefs/fix-tui-
    ┃  t3a-2.md) |
    ┃  30:**The briefs that exist** are T1, T2a, T2b and T3a. The rest are written **just before**
    ┃  33:**The next session is T3b** — the backlog's moves: dismiss · create · edit. T3a left
    ┃  fix-tui-t3a-1.md
    ┃  fix-tui-t3a-2.md
    ┃  review-tui-t3a.md
    ┃  tui-t3a.md
    ┃

       + Thought · 5.2s

    ┃
    ┃  ← Edit work/tui/README.md
    ┃
    ┃    14   |--:|---|---|---|---|                                                                                                                              14   |--:|---|---|---|---|
    ┃    15   | **T1** | the identity | the window title carries the project; the header's name stops being a footnote; `--project`; the Settings tab shows the  15   | **T1** | the identity | the window title carries the project; the header's name stops being a footnote; `--project`; the Settings tab shows the
    ┃         plan's path | ✅ closed · 11 commits · 94a53ff, 8fc7438, f0b2aa5, 82289c6, 0bbcb35, 5a79638, 67b455d, b4a3052, 55bdc81, 8a6866b, cf8c05e | [tui-        plan's path | ✅ closed · 11 commits · 94a53ff, 8fc7438, f0b2aa5, 82289c6, 0bbcb35, 5a79638, 67b455d, b4a3052, 55bdc81, 8a6866b, cf8c05e | [tui-
    ┃         t1.md](../briefs/tui-t1.md) · [fix 1](../briefs/fix-tui-t1-1.md) · [fix 2](../briefs/fix-tui-t1-2.md) |                                                 t1.md](../briefs/tui-t1.md) · [fix 1](../briefs/fix-tui-t1-1.md) · [fix 2](../briefs/fix-tui-t1-2.md) |
    ┃    16   | **T2a** | the mouse, part 1 | the mouse is on; the render builds the hit map; the wheel scrolls | ✅ closed · 12 commits · a72ea98, 115b650,     16   | **T2a** | the mouse, part 1 | the mouse is on; the render builds the hit map; the wheel scrolls | ✅ closed · 12 commits · a72ea98, 115b650,
    ┃         284fb65, 0672f9f, 9f971cb, 7035ec7, e4b7061, 619b342, 47be828, f1cc1ca, a332158, d8d7b95 | [tui-t2a.md](../briefs/tui-t2a.md) · [fix 1](../briefs/      284fb65, 0672f9f, 9f971cb, 7035ec7, e4b7061, 619b342, 47be828, f1cc1ca, a332158, d8d7b95 | [tui-t2a.md](../briefs/tui-t2a.md) · [fix 1](../
    ┃         fix-tui-t2a-1.md) · [fix 2](../briefs/fix-tui-t2a-2.md) |                                                                                               briefs/fix-tui-t2a-1.md) · [fix 2](../briefs/fix-tui-t2a-2.md) |
    ┃    17   | **T2b** | the mouse, part 2 | the clicks: a tab, a row, a detail; the mouse drivable from a script | ✅ closed · 9 commits · 4dfd5a4, 187d706,   17   | **T2b** | the mouse, part 2 | the clicks: a tab, a row, a detail; the mouse drivable from a script | ✅ closed · 9 commits · 4dfd5a4, 187d706,
    ┃         f64e6a2, ac8df77, 588cf5f, 9c716ce, 0144216, 6bb7545, f8969ef | [tui-t2b.md](../briefs/tui-t2b.md) · [fix 1](../briefs/fix-tui-t2b-1.md) |              f64e6a2, ac8df77, 588cf5f, 9c716ce, 0144216, 6bb7545, f8969ef | [tui-t2b.md](../briefs/tui-t2b.md) · [fix 1](../briefs/fix-tui-t2b-1.md) |
    ┃    18 - | **T3a** | the backlog is readable | the detail pane; `space` stops being a trap on an idea; D-3, D-5 | 🔄 **fix cycle 2** — 7 commits · 4dc8079, 18 + | **T3a** | the backlog is readable | the detail pane; `space` stops being a trap on an idea; D-3, D-5 | 🔄 **fix cycle 2** — 8 commits ·
    ┃         1b550c6, 1804ccc, b873821, 5edae71, aa4ba5f, da36146 | [tui-t3a.md](../briefs/tui-t3a.md) · [fix 1](../briefs/fix-tui-t3a-1.md) · [fix 2](../           4dc8079, 1b550c6, 1804ccc, b873821, 5edae71, aa4ba5f, da36146, 63964ce | [tui-t3a.md](../briefs/tui-t3a.md) · [fix 1](../briefs/fix-tui-t3a-1.md)
    ┃         briefs/fix-tui-t3a-2.md) |                                                                                                                              · [fix 2](../briefs/fix-tui-t3a-2.md) |
    ┃    19   | **T3b** | the backlog moves | dismiss · create · edit | ⏳ planned | — |                                                                         19   | **T3b** | the backlog moves | dismiss · create · edit | ⏳ planned | — |
    ┃    20   | **T4** | the Home tab | the landing tab; the art; the Sous-Chef dissolves; D-2 | ⏳ planned | — |                                                20   | **T4** | the Home tab | the landing tab; the art; the Sous-Chef dissolves; D-2 | ⏳ planned | — |
    ┃    21   | **T5** | the Ready tab | `core.Adoptable`; the Path becomes Ready; D-1, D-4 | ⏳ planned | — |                                                   21   | **T5** | the Ready tab | `core.Adoptable`; the Path becomes Ready; D-1, D-4 | ⏳ planned | — |
    ┃    22   | **T6** | the Board, as a visual place | the tree drawn left to right; clicking a node | ⏳ planned | — |                                         22   | **T6** | the Board, as a visual place | the tree drawn left to right; clicking a node | ⏳ planned | — |
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git add work/tui/README.md && git commit -q -m "work: T3a's row names the fix cycle 2 commit" && git log --oneline -3 && git status --short
    ┃
    ┃  c47728e work: T3a's row names the fix cycle 2 commit
    ┃  63964ce ui: the key and the pane agree on every tab
    ┃  c2e5b7a work: T3a's fix cycle 2 - the key and the pane, on every tab
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && make check 2>&1 | tail -14 && git status --short && echo "(tree clean)"
    ┃
    ┃  (8 earlier lines)
    ┃    internal/core: 0 external modules (ceiling 0)
    ┃    internal/cli: 0 external modules (ceiling 0)
    ┃    internal/ui: 28 external modules (ceiling 30)
    ┃  /Library/Developer/CommandLineTools/usr/bin/make tidy-check
    ┃
    ┃  ==> tidy-check
    ┃  (tree clean)
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && make check-race 2>&1 | grep -E "^(ok|FAIL|---)" ; git status --short; echo "(tree clean)"
    ┃
    ┃  ok    github.com/lthiagol/maitre/internal/brand  1.488s
    ┃  ok    github.com/lthiagol/maitre/internal/cli  9.355s
    ┃  ok    github.com/lthiagol/maitre/internal/core  88.986s
    ┃  ok    github.com/lthiagol/maitre/internal/ui  301.336s
    ┃  ok    github.com/lthiagol/maitre/internal/ui  4.269s
    ┃  (tree clean)
    ┃

       + Thought · 4.2s

    ┃
    ┃  $ mkdir -p /Users/thiago/Code/Pessoal/maitre/work/inbox && cat > /Users/thiago/Code/Pessoal/maitre/work/inbox/runner-tui-t3a-fix2.md <<'EOF'
    ┃  Runner: T3a fix 2 done. Commits: 2 (63964ce, c47728e) — the seam and the state row. Fixed: F-116.…
    ┃
    ┃  134ff14 inbox: runner reports tui-t3a fix 2
    ┃  c47728e work: T3a's row names the fix cycle 2 commit
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && herdr agent prompt "w1H:p1" "Runner: T3a fix 2 done. Commits: 2 (63964ce, c47728e) -- the seam and the state row. Fixed: F-116. The Log's detail: enter on the Log tab now opens the event's detail -- the same pane the recipe gets, the same geometry, the same scroll,
    ┃   the same esc. hasDetail() answers for the Log, so the key and the pane read one answer. By command: 'maitre --render --plain --keys \"tab,tab,tab,tab,enter\"' shows 'created', the timestamp, 'chef - chef' and the Data section with the event's JSON. The footer's promise: the footer offers 'enter det…
    ┃
    ┃  {"id":"cli:agent:prompt","result":{"agent":{"agent":"opencode","agent_session":{"agent":"opencode","kind":"id","source":"herdr:opencode","value":"ses_e737d6341ffeTdoEYTgykU25M8"},"agent_status":"done","completion_seq":323,"cwd":"/Users/thiago/Code/Pessoal/maitre","focused":false,"foreground_cwd":"/Users/
    ┃  thiago/Code/Pessoal/maitre","pane_id":"w1H:p1","revision":15,"screen_detection_skipped":true,"state_change_seq":323,"tab_id":"w1H:t1","terminal_id":"term_65f0cd2aa89d26","terminal_title":"OC | Project work proposal","terminal_title_stripped":"OC | Project work proposal","workspace_id":"w1H"},"type":
    ┃  "agent_prompted"}}
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git log --oneline -6 && echo "--- tree ---" && git status --short && echo "(clean)"
    ┃
    ┃  134ff14 inbox: runner reports tui-t3a fix 2
    ┃  c47728e work: T3a's row names the fix cycle 2 commit
    ┃  63964ce ui: the key and the pane agree on every tab
    ┃  c2e5b7a work: T3a's fix cycle 2 - the key and the pane, on every tab
    ┃  322b280 inbox: runner reports tui-t3a fix 1
    ┃  8f39289 work: T3a's row names the fix cycle's three commits
    ┃  --- tree ---
    ┃  (clean)
    ┃

       T3a's fix cycle 2 is done — 2 commits (63964ce, c47728e), receipt committed (134ff14), orchestrator nudged.

       F-116 — the key and the pane agree on every tab

       Two changes, both in the seam:

       1. hasDetail() answers for the Log. The pane already rendered the journal entry; the key could not open it. Now enter on the Log opens the event's detail — same pane, same geometry, same scroll, same esc:
          $ maitre --render --plain --keys "tab,tab,tab,tab,enter"
           created
           2026-10-30T22:03:09Z
           chef · chef
          Data
               { "id": "BL06", "item_kind": "bug", … }
       2. The footer offers enter detail only where there is a detail to open — the same answer the key reads.

       ┌────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┬────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┐
       │ Tab                                                                                                                                                │ The footer's promise                                                                                                                                       │
       ├────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┼────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
       │ Recipes (a recipe selected)                                                                                                                        │ enter detail ✓                                                                                                                                             │
       ├────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┼────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
       │ Backlog · Path · Board                                                                                                                             │ enter detail ✓                                                                                                                                             │
       ├────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┼────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
       │ Log                                                                                                                                                │ enter detail ✓ — and the entry opens                                                                                                                       │
       ├────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┼────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
       │ Sous-Chef                                                                                                                                          │ no enter detail                                                                                                                                            │
       ├────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┼────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
       │ Recipes, nothing selected                                                                                                                          │ no enter detail                                                                                                                                            │
       ├────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┼────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
       │ Settings                                                                                                                                           │ enter apply theme                                                                                                                                          │
       └────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┴────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┘

       One answer, tab by tab: TestTheKeyAndThePaneAgreeOnEveryTab walks all seven tabs and holds that the footer offers enter detail exactly where hasDetail() is true — and where it offers it, the key opens the detail. Two mutations, both reverted, both red: dropping the Log case fails TestEnterOpensTheLogEntry
       + TestTheFooterPromisesOnlyWhatItDoes; making the footer promise unconditionally fails TestTheFooterPromisesOnlyWhatItDoes + TestTheKeyAndThePaneAgreeOnEveryTab.

       DomainVersion 7, unmoved — internal/core and internal/cli untouched; the seam needed no domain change. The goldens are unchanged (make goldens leaves git diff empty): the footer's key list is rendered text, and the goldens hold the Recipes tab, where the keys do not change. make check and make check-race
       green; tree clean.

       docs/05 is unchanged — no line of it became false: the doc never promised enter detail on every tab, and the interaction table's enter row is still true. The Sous-Chef's own pane message stays, as ruled.

       Two things I named in the receipt for the reviewer: the agreement test holds the agreement, while the two facts (the Log has a detail; the Sous-Chef and an empty list have none) are held by the other two tests — the split is deliberate, since the agreement test alone cannot catch "the Log has no detail at
       all". And the Log's detail was already drawn in the split layout before this cycle; what was missing was the key and the footer's promise.

       Build · MiniMax-M3.1-Flash · 16m 12s · 135.3 tok/s

    ┃
    ┃
    ┃
    ┃  Build · MiniMax-M3.1-Flash MiniMax Token Plan (minimax.io) · medium
    ╹▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀
    ~/Code/Pessoal/maitre:master                                                                                                                                                                                                                                                             281.0K (28%)  ctrl+p commands
```
