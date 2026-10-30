# The reviewer's lane — T2a, captured 2026-10-30

Dumped from pane `w1H:p5` with `herdr pane read --source recent-unwrapped --lines 5000`.
The lane that judged T2a: three cycles, and the findings they closed —
F-106..F-109.
The lane that wrote it is in [tui-t2a-runner.md](tui-t2a-runner.md).

    ┃                                                                               67   | The rounds inside a cycle | Fix → verify → back to the executor
    ┃                                                                                    → repeat, until the verifier is satisfied | Three rounds in cycle
    ┃     66   | The rounds inside a cycle | Fix → verify → back to the executor →       8: the first fix introduced three new mediums, the second closed
    ┃          repeat, until the verifier is satisfied | Three rounds in cycle 8:        them |
    ┃          the first fix introduced three new mediums, the second closed them   68
    ┃          |                                                                    69   ### The discipline every round runs
    ┃                                                                               70
    ┃     67
    ┃     68   ### The discipline every round runs
    ┃     69
    ┃
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git diff --stat && grep -n "F-10[6-9]" docs/audit.md | cut -c1-70 && make check 2>&1 | tail -6
    ┃
    ┃  (7 earlier lines)
    ┃  214:| F-108 | low | **The zero box was a one-cell region at the origin
    ┃    internal/core: 0 external modules (ceiling 0)
    ┃    internal/cli: 0 external modules (ceiling 0)
    ┃    internal/ui: 28 external modules (ceiling 30)
    ┃  /Library/Developer/CommandLineTools/usr/bin/make tidy-check
    ┃
    ┃  ==> tidy-check
    ┃

       + Thought · 6.8s

    ┃
    ┃  ← Edit docs/audit.md
    ┃
    ┃     62   | Slice T1 · cycle 1 | The review of the identity: four mutations,   62   | Slice T1 · cycle 1 | The review of the identity: four mutations,
    ┃          the write-nothing proof, and the gate | Each mutation fails a test        the write-nothing proof, and the gate | Each mutation fails a test
    ┃          and was reverted, and the tree was left exactly as the runner left        and was reverted, and the tree was left exactly as the runner left
    ┃          it: the hardcoded `"maitre"` fails **four** window-title tests; a         it: the hardcoded `"maitre"` fails **four** window-title tests; a
    ┃          `--project` that writes the name into `project.json` fails *"the          `--project` that writes the name into `project.json` fails *"the
    ┃          override writes nothing"* and *"the plan keeps its own name"*; the        override writes nothing"* and *"the plan keeps its own name"*; the
    ┃          name back in `styFaint` fails *"the header's name is not a                name back in `styFaint` fails *"the header's name is not a
    ┃          footnote"*; a title that computes its own name fails *"the override       footnote"*; a title that computes its own name fails *"the
    ┃          reaches the console"*. The write-nothing claim was proved **by            override reaches the console"*. The write-nothing claim was proved
    ┃          command** and not by the diff: a copy of a plan, `maitre --render -       **by command** and not by the diff: a copy of a plan, `maitre --
    ┃          -project …` over it, and `cmp` on both `project.json` and                 render --project …` over it, and `cmp` on both `project.json` and
    ┃          `activity.json`. `make check` and `make check-race` are green,            `activity.json`. `make check` and `make check-race` are green,
    ┃          `make goldens` leaves the tree clean, and the state row is **in**         `make goldens` leaves the tree clean, and the state row is **in**
    ┃          the slice's commits — **no blockers filed** |                             the slice's commits — **no blockers filed** |
    ┃     63   | Slice T1 · cycle 2 | The review of the fix cycle: the four         63   | Slice T1 · cycle 2 | The review of the fix cycle: the four
    ┃          findings, the two observations, and the `docs/03` question | F-99..       findings, the two observations, and the `docs/03` question | F-99.
    ┃          F-102 are closed, and each was proved by a mutation: the `sanitize`       .F-102 are closed, and each was proved by a mutation: the
    ┃          call taken out of the name's path fails `TestTheNameIsSanitised`          `sanitize` call taken out of the name's path fails
    ┃          three times over; the `root` row put back on the parent fails             `TestTheNameIsSanitised` three times over; the `root` row put back
    ┃          `TestTheSettingsTabShowsThePlanRoot`; the empty-name branch taken         on the parent fails `TestTheSettingsTabShowsThePlanRoot`; the
    ┃          out of `LoadProject` leaves **the whole suite green** — and that is       empty-name branch taken out of `LoadProject` leaves **the whole
    ┃          **F-103**, the one blocker filed. The name has one answer in the          suite green** — and that is **F-103**, the one blocker filed. The
    ┃          code, but the rule is written **twice** and only one copy is              name has one answer in the code, but the rule is written **twice**
    ┃          pinned: `maitre-cli plan` then prints `"Name": ""` while the              and only one copy is pinned: `maitre-cli plan` then prints
    ┃          console shows the project's directory name. The **path row** is           `"Name": ""` while the console shows the project's directory name.
    ┃          judged a note and not a blocker: the path is not the plan's               The **path row** is judged a note and not a blocker: the path is
    ┃          content, the unsanitised path **predates** the slice, and the             not the plan's content, the unsanitised path **predates** the
    ┃          consequence is one extra visual row in a settings list — the header       slice, and the consequence is one extra visual row in a settings
    ┃          and the window title both survive. The **state row** is judged a          list — the header and the window title both survive. The **state
    ┃          note as well: the count should count what the list names, or the          row** is judged a note as well: the count should count what the
    ┃          row should say the count excludes the commits that wrote it. The          list names, or the row should say the count excludes the commits
    ┃          plan's name does **not** belong in `docs/03` § 8: that table is the       that wrote it. The plan's name does **not** belong in `docs/03` §
    ┃          derived views — *computed, never stored* — and the name is a              8: that table is the derived views — *computed, never stored* —
    ┃          **stored** fact the domain reads, like `LoadProject` and                  and the name is a **stored** fact the domain reads, like
    ┃          `ListRecipes`, which the table does not carry either. `make check`        `LoadProject` and `ListRecipes`, which the table does not carry
    ┃          and `make check-race` are green, the goldens did not move, and the        either. `make check` and `make check-race` are green, the goldens
    ┃          tree was left exactly as the runner left it — **one blocker filed**       did not move, and the tree was left exactly as the runner left it
    ┃          (F-103) |                                                                 — **one blocker filed** (F-103) |
    ┃     64   | Slice T1 · cycle 3 | The **last** cycle of T1: the blocker, the    64   | Slice T1 · cycle 3 | The **last** cycle of T1: the blocker, the
    ┃          two notes, and the question that is the chef's | F-103..F-105 are         two notes, and the question that is the chef's | F-103..F-105 are
    ┃          closed, each proved by a mutation and each reverted: the rule taken       closed, each proved by a mutation and each reverted: the rule
    ┃          out of the one copy turns the suite **red** in both `internal/core`       taken out of the one copy turns the suite **red** in both
    ┃          and `internal/ui` — **three** tests, as the runner claimed; the           `internal/core` and `internal/ui` — **three** tests, as the runner
    ┃          `sanitize` call taken out of `humanPath` fails the root-row test          claimed; the `sanitize` call taken out of `humanPath` fails the
    ┃          twice over; the state row's count is the number of ids the row            root-row test twice over; the state row's count is the number of
    ┃          names, and the **protocol section** carries the shape in one line.        ids the row names, and the **protocol section** carries the shape
    ┃          The two surfaces agree in **all three states** — the plan names           in one line. The two surfaces agree in **all three states** — the
    ┃          itself, the file names it empty, the file is missing — and the            plan names itself, the file names it empty, the file is missing —
    ┃          CLI's behaviour **changed** to match the console's, which is the          and the CLI's behaviour **changed** to match the console's, which
    ┃          point of F-103. The rule is written down in the **code** (`core.          is the point of F-103. The rule is written down in the **code** (
    ┃          Store.projectName`'s doc comment) and **nowhere in the spec**:            `core.Store.projectName`'s doc comment) and **nowhere in the
    ┃          `docs/03` § 6 names `project.json` as *project meta · charter ·           spec**: `docs/03` § 6 names `project.json` as *project meta ·
    ┃          settings* and says nothing about the name's fallback, and `docs/05`       charter · settings* and says nothing about the name's fallback,
    ┃          does not either — where it belongs is the **chef's** call, and            and `docs/05` does not either — where it belongs is the **chef's**
    ┃          **nothing was written into `docs/03`**. The goldens did not move,         call, and **nothing was written into `docs/03`**. The goldens did
    ┃          the budget is still **0 / 0 / 28**, `core.DomainVersion` is still         not move, the budget is still **0 / 0 / 28**, `core.DomainVersion`
    ┃          **7**, the seven tabs kept their names, the CLI's surface did not         is still **7**, the seven tabs kept their names, the CLI's surface
    ┃          change, and the tree was left exactly as the runner left it — **no        did not change, and the tree was left exactly as the runner left
    ┃          blockers filed** |                                                        it — **no blockers filed** |
    ┃     65   | Slice T2a | The mouse, part 1: the mouse is on; the render builds  65   | Slice T2a | The mouse, part 1: the mouse is on; the render
    ┃          the hit map; the wheel scrolls the list and the detail | The mouse        builds the hit map; the wheel scrolls the list and the detail |
    ┃          is on **every** view the console draws — the console's, the error's       The mouse is on **every** view the console draws — the console's,
    ┃          and the opening one's — and the error and opening views **clear the       the error's and the opening one's — and the error and opening
    ┃          map**, so no event can land on a tab that is not on screen. The map       views **clear the map**, so no event can land on a tab that is not
    ┃          is built **by the render**: `tabBar` emits the tab boxes as it            on screen. The map is built **by the render**: `tabBar` emits the
    ┃          draws them, `hitMapOf` is called from the three render functions          tab boxes as it draws them, `hitMapOf` is called from the three
    ┃          and from nowhere else, the banner pushes the list's rows down in          render functions and from nowhere else, the banner pushes the
    ┃          the map as it does on the screen, and the geometry constants (            list's rows down in the map as it does on the screen, and the
    ┃          `headerRows`, `footerRows`, `splitThreshold`) are read where they         geometry constants (`headerRows`, `footerRows`, `splitThreshold`)
    ┃          are written. The wheel over the list drives `m.move` — the                are read where they are written. The wheel over the list drives
    ┃          keyboard's own move — and a scripted `wheel:down` and three `down`s       `m.move` — the keyboard's own move — and a scripted `wheel:down`
    ┃          land on **byte-identical** frames. The wheel over the detail drives       and three `down`s land on **byte-identical** frames. The wheel
    ┃          `m.scrollDetail`, but **the split layout never renders that               over the detail drives `m.scrollDetail`, but **the split layout
    ┃          offset**: at ≥160 columns — where the console draws the list and          never renders that offset**: at ≥160 columns — where the console
    ┃          the detail side by side, and `ModeSplit` is the default — the pane        draws the list and the detail side by side, and `ModeSplit` is the
    ┃          does not move, and the clamp measures the detail at a different           default — the pane does not move, and the clamp measures the
    ┃          width than the pane draws it. The scripted mouse's **position             detail at a different width than the pane draws it. The scripted
    ┃          form** cannot be reached through `maitre --keys`, and the doc             mouse's **position form** cannot be reached through `maitre --
    ┃          prints the command line that fails. The scope held: **no click            keys`, and the doc prints the command line that fails. The scope
    ┃          exists**, on purpose; `internal/core` and `internal/cli` are              held: **no click exists**, on purpose; `internal/core` and
    ┃          untouched; `core.DomainVersion` is still **7**; the budget is still       `internal/cli` are untouched; `core.DomainVersion` is still **7**;
    ┃          **0 / 0 / 28**; the seven tabs kept their names; the goldens did          the budget is still **0 / 0 / 28**; the seven tabs kept their
    ┃          not move — **two blockers filed** (F-106, F-107) |                        names; the goldens did not move — **two blockers filed** (F-106,
    ┃                                                                                    F-107) |
    ┃                                                                               66 + | Slice T2a · cycle 1 | The review of the slice: the map, the
    ┃                                                                                    wheel, and six mutations | The slice's shape held: the map **is**
    ┃                                                                                    built by the render (`tabBar` emits the tab boxes, `hitMapOf` is
    ┃                                                                                    called from the three render functions and nowhere else, the
    ┃                                                                                    banner pushes the list's rows down in the map, and the geometry
    ┃                                                                                    constants are read where they are written), the wheel **is** the
    ┃                                                                                    keyboard's scroll (a scripted `wheel:down` and three `down`s land
    ┃                                                                                    on byte-identical frames), and all **six** mutations failed the
    ┃                                                                                    suite — the mouse mode off the error view, the render not emitting
    ┃                                                                                    the map, a second scroll variable, the wheel over the header, the
    ┃                                                                                    wheel under an open prompt, and the scripted step acting without
    ┃                                                                                    drawing first. Two blockers filed: **F-106** (the wheel over the
    ┃                                                                                    detail does nothing in the split layout) and **F-107** (the
    ┃                                                                                    scripted position form is unreachable through `maitre --keys`, and
    ┃                                                                                    the doc prints the command that fails). The zero box's one cell at
    ┃                                                                                    the origin was noted as **F-108** and ruled into the fix cycle.
    ┃                                                                                    The goldens did not move; `internal/core` was untouched; `core.
    ┃                                                                                    DomainVersion` was still **7**; the budget was still **0 / 0 /
    ┃                                                                                    28** — **verdict: changes-needed** |
    ┃     66   | Slice T2a · cycle 2 | The review of the fix cycle: the split       67   | Slice T2a · cycle 2 | The review of the fix cycle: the split
    ┃          detail, the position form, the absent region, and the **reversed**        detail, the position form, the absent region, and the **reversed**
    ┃          test | F-106, F-107 and F-108 are closed, each proved by a                test | F-106, F-107 and F-108 are closed, each proved by a
    ┃          mutation: taking the offset out of the split layout fails                 mutation: taking the offset out of the split layout fails
    ┃          `TestTheWheelScrollsTheSplitDetail`; a splitter that cuts a               `TestTheWheelScrollsTheSplitDetail`; a splitter that cuts a
    ┃          bracketed position at its comma fails the two position tests;             bracketed position at its comma fails the two position tests;
    ┃          taking the `there` guard out of `hitBox` fails three tests. The           taking the `there` guard out of `hitBox` fails three tests. The
    ┃          offset is applied in **one** place and the pane's width has **one**       offset is applied in **one** place and the pane's width has
    ┃          answer (`detailPaneWidth`), and the wheel at the origin over the          **one** answer (`detailPaneWidth`), and the wheel at the origin
    ┃          list scrolls the list. **But the bound's guard does not bite**: the       over the list scrolls the list. **But the bound's guard does not
    ┃          bound put back on the full-screen width leaves the **whole                bite**: the bound put back on the full-screen width leaves the
    ┃          `internal/ui` suite green** — **one blocker filed** (F-109). The          **whole `internal/ui` suite green** — **one blocker filed** (F-
    ┃          **reversed** test is judged **sound**: the rule is one rule — the         109). The **reversed** test is judged **sound**: the rule is one
    ┃          detail answers over the detail, the list answers the rest — stated        rule — the detail answers over the detail, the list answers the
    ┃          in one place in the code and in the doc; the full-screen detail           rest — stated in one place in the code and in the doc; the full-
    ┃          keeps the *no list, nothing happens* answer, pinned by                    screen detail keeps the *no list, nothing happens* answer, pinned
    ┃          `TestTheWheelOverTheFullScreenHeaderDoesNothing`; and the old rule        by `TestTheWheelOverTheFullScreenHeaderDoesNothing`; and the old
    ┃          and the fix cycle's own F-108 evidence (*a wheel at the origin over       rule and the fix cycle's own F-108 evidence (*a wheel at the
    ┃          the list scrolls the list*) cannot both hold, because the origin          origin over the list scrolls the list*) cannot both hold, because
    ┃          **is** the header's first cell. The scope held: **no click                the origin **is** the header's first cell. The scope held: **no
    ┃          exists**, on purpose; `internal/core`, `internal/cli` and `docs/03`       click exists**, on purpose; `internal/core`, `internal/cli` and
    ┃          are untouched; `core.DomainVersion` is still **7**; the budget is         `docs/03` are untouched; `core.DomainVersion` is still **7**; the
    ┃          still **0 / 0 / 28**; the seven tabs kept their names; the goldens        budget is still **0 / 0 / 28**; the seven tabs kept their names;
    ┃          did not move; `make check` and `make check-race` are green. |             the goldens did not move; `make check` and `make check-race` are
    ┃                                                                                    green. |
    ┃     67   | The rounds inside a cycle | Fix → verify → back to the executor →  68   | The rounds inside a cycle | Fix → verify → back to the executor
    ┃          repeat, until the verifier is satisfied | Three rounds in cycle 8:        → repeat, until the verifier is satisfied | Three rounds in cycle
    ┃          the first fix introduced three new mediums, the second closed them        8: the first fix introduced three new mediums, the second closed
    ┃          |                                                                         them |
    ┃     68                                                                        69
    ┃     69   ### The discipline every round runs                                  70   ### The discipline every round runs
    ┃
    ┃
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git diff --stat && git status --short && echo "--- only the ledger changed ---" && git diff --name-only
    ┃
    ┃  docs/audit.md | 8 ++++++--
    ┃   1 file changed, 6 insertions(+), 2 deletions(-)
    ┃   M docs/audit.md
    ┃  --- only the ledger changed ---
    ┃  docs/audit.md
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git add docs/audit.md && git commit -q -m "docs: T2a's fix cycle — the split detail and the position form
    ┃  are closed, the bound's guard is not…
    ┃
    ┃  97779e5 docs: T2a's fix cycle — the split detail and the position form are closed, the bound's guard is not
    ┃  98aebc9 work: T2a's cycle 2 review brief - and the reversed test
    ┃  CLEAN
    ┃

       + Thought · 1.6s

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && mkdir -p work/inbox && cat > work/inbox/reviewer-tui-t2a-cycle2.md <<'EOF'
    ┃  Reviewer: T2a cycle 2 done. Verdict: changes-needed. Filed 1 blocker (F-109).…
    ┃
    ┃  d939ab1 inbox: reviewer reports tui-t2a cycle 2
    ┃  97779e5 docs: T2a's fix cycle — the split detail and the position form are closed, the bound's guard is not
    ┃  CLEAN
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && herdr agent prompt "w1H:p1" "Reviewer: T2a cycle 2 done. Verdict: changes-needed. Filed 1 blocker (F-109).
    ┃   F-106: proved by command — maitre --render --width 200 --keys \"enter,esc\" against \"enter,esc,wheel:down@(150,6)\": the detail column chan…
    ┃
    ┃  {"id":"cli:agent:prompt","result":{"agent":{"agent":"opencode","agent_session":{"agent":"opencode","kind":"id","source":"herdr:opencode","value":
    ┃  "ses_e737d6341ffeTdoEYTgykU25M8"},"agent_status":"done","completion_seq":145,"cwd":"/Users/thiago/Code/Pessoal/maitre","focused":false,
    ┃  "foreground_cwd":"/Users/thiago/Code/Pessoal/maitre","pane_id":"w1H:p1","revision":15,"screen_detection_skipped":true,"state_change_seq":145,
    ┃  "tab_id":"w1H:t1","terminal_id":"term_65f0cd2aa89d26","terminal_title":"OC | Project work proposal","terminal_title_stripped":"OC | Project work
    ┃  proposal","workspace_id":"w1H"},"type":"agent_prompted"}}
    ┃

       verdict: changes-needed — 1 blocker filed (F-109).

       What I proved, by command

       ┌─────────────┬────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┐
       │ Check       │ Result                                                                                                                             │
       ├─────────────┼────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
       │ F-106       │ --width 200 --keys "enter,esc" vs "enter,esc,wheel:down@(150,6)": the detail column changed by 6 lines (the recipe's first three   │
       │             │ rows gone), the list column is byte-identical. The other way: a wheel over the list moves the list and leaves the detail's scroll  │
       │             │ alone                                                                                                                              │
       ├─────────────┼────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
       │ The bound   │ counted at the width the pane draws with — detailPaneWidth(): width-4 full screen, width-listColumn()-3 in the split layout;       │
       │             │ renderSplit, renderDetailFull and detailLineCount all ask that one function; the offset is applied in one place (detailScrollRows) │
       │             │ — but nothing holds it                                                                                                             │
       ├─────────────┼────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
       │ F-107       │ --keys "wheel:down@(12,5)" exits 0; "tab,down,space" still exits 0; the old spelling says malformed mouse step "wheel:down@12";    │
       │             │ docs/05 prints the bracketed form                                                                                                  │
       ├─────────────┼────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
       │ One grammar │ the splitter ui.SplitSteps and the parser parseMouseStep are in the same file; cmd/maitre calls ui.SplitSteps and nothing else     │
       ├─────────────┼────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
       │ F-108       │ --keys "wheel:down@(0,0)" is byte-identical to "down,down,down"; taking the there guard out turns three tests red                  │
       ├─────────────┼────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
       │ The goldens │ unchanged                                                                                                                          │
       ├─────────────┼────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
       │ The gate    │ make check and make check-race green; no click; internal/core, internal/cli, docs/03 untouched; DomainVersion 7; budget 0/0/28     │
       └─────────────┴────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┘

       The blocker

       F-109 (medium) — TestTheScrollsBoundIsThePanesWidth asks whether the bound equals the pane's row count, but the seeded recipe's detail has the same
       count at both widths, so the assertion holds whatever the bound is. My mutation — detailLineCount() back on the full-screen width, the pre-fix
       defect — leaves the whole internal/ui suite green, and the consequence is real: five rows of the split pane become unreachable.

       The reversed test — judged sound, not filed

       The rule is one rule (the detail answers over the detail, the list answers the rest when it is on screen), in one place, stated once in docs/05. The
       old rule and the fix cycle's own F-108 evidence cannot both hold — the origin is the header's first cell — so the reversal was forced by the ruling,
       not chosen. The full-screen detail still ignores the wheel over the header, pinned by a test and proved by command ("enter" and "enter,wheel:down@(
       150,0)" are byte-identical).

       LOW notes: docs/05's table says "anywhere else → scroll the list", which the full-screen detail's header does not do (one more line would make it
       whole); detailScrollRows writes m.detailScroll during a render; the wheel over the list in the split layout moves the selection, so the pane shows
       another recipe from its top.

       Ledger 97779e5, receipt d939ab1, nudge sent. The tree is exactly as the runner left it.

       Build · MiniMax-M3.1-Flash · 14m 40s · 136.8 tok/s

    ┃
    ┃  You are the **reviewer** for **T2a** — cycle 3. The same lane that reviewed cycles 1
    ┃  and 2. The orchestrator's pane id is `w1H:p1`; use it verbatim in the notify command
    ┃  below.
    ┃
    ┃  Read [work/tui/README.md](../tui/README.md) for the state and the protocol, and the
    ┃  slice's briefs: [tui-t2a.md](tui-t2a.md) · [fix-tui-t2a-1.md](fix-tui-t2a-1.md) ·
    ┃  [fix-tui-t2a-2.md](fix-tui-t2a-2.md).
    ┃
    ┃  ## Cycle 3 is the **last** cycle
    ┃
    ┃  **This is the last cycle for T2a.** Anything still open after it is a **note for the
    ┃  ledger**, not another round: the orchestrator escalates to the chef.
    ┃
    ┃  File **HIGH** blockers only. A MEDIUM finding at this point is a **note for the
    ┃  ledger** — the cycle cap says the slice ends, and a note is how a finding ends without
    ┃  a fourth round.
    ┃
    ┃  ## What cycles 1 and 2 decided
    ┃
    ┃  | Cycle | Verdict | What it closed |
    ┃  |--:|---|---|
    ┃  | 1 | `changes-needed` | **F-106** (the split layout's detail never scrolled) · **F-107** (the position form could not survive the command line) ·
    ┃  **F-108** (a region that is not drawn answered for one cell) |
    ┃  | 2 | `changes-needed` | **F-109** — the bound's width was held by nothing, because the test compared two values that **both** went through the same
    ┃  function, and the fixture's detail was too short for the two widths to differ |
    ┃
    ┃  **The wheel's fallback is closed.** Cycle 2 judged the reversed test **sound**: the
    ┃  list is the console's one scrollable thing, the wheel scrolls it everywhere except
    ┃  over the detail, and the rule is stated once. **Do not reopen it.**
    ┃
    ┃  The fix cycle landed in three commits: `a332158` · `d8d7b95` · `22b9de6`.
    ┃
    ┃  ## What to verify — each by command, not by reading
    ┃
    ┃  ### 1. F-109 — the bound's width is held by a test that can fail
    ┃
    ┃  - The bound is counted at the width the pane **draws with**, and the scroll reaches
    ┃    the pane's **last row**.
    ┃  - **The tripwire must bite.** Count the bound at the full-screen width again — the
    ┃    suite must go **red**, and it must fail **for the bound's width** and not for some
    ┃    unrelated number.
    ┃  - **The test must not ask the code under test for the answer.** The old test compared
    ┃    two values that both went through `detailPaneWidth()`. The new one computes the
    ┃    pane's width **in the test**. Check that.
    ┃  - **The fixture must make the two widths differ**, and the test must say so if they
    ┃    ever stop differing — a fixture that silently agreed again would make the tripwire
    ┃    vacuous a second time.
    ┃
    ┃  ### 2. The document's wheel table reads whole
    ┃
    ┃  - Every row of the table is **true**, including the row that says the wheel scrolls
    ┃    the list *"anywhere else, while the list is on screen"*.
    ┃  - The prose and the table agree — check that the **full-screen detail** is the case
    ┃    the table carves out, and that the console really does what the table says there.
    ┃
    ┃  ### 3. The scope held
    ┃
    ┃  | | |
    ┃  |---|---|
    ┃  | **No click** | the click is **T2b**. There must be **no** click handler. `grep` for it |
    ┃  | **The wheel's fallback** | unchanged since cycle 2 |
    ┃  | **`maitre-cli`** | no verb is added and no command changes |
    ┃  | **The tabs** | seven, the same names |
    ┃  | **The domain** | `internal/core` untouched |
    ┃  | **`core.DomainVersion`** | still **7** |
    ┃  | **The budget** | **0 / 0 / 28** against **0 / 0 / 30** |
    ┃  | **The goldens** | `make goldens` leaves `git diff` empty |
    ┃  | **`docs/03`** | untouched |
    ┃  | **`docs/audit.md`** | the ledger is **yours** |
    ┃
    ┃  ## The notes the fix cycle left for a later round
    ┃
    ┃  **Do not treat these as findings.** The fix cycle was told to leave them, and a cycle
    ┃  that widens drifts:
    ┃
    ┃  - `detailScrollRows` writes `m.detailScroll` during a **render** — the clamping. It
    ┃    moved there from `renderDetailFull`; it is not new.
    ┃  - The wheel over the list in the split layout moves the **selection**, so the pane
    ┃    shows another recipe from its top and the detail does not scroll. That is the
    ┃    console's model.
    ┃
    ┃  **If you think either is a finding**, say so in your message — the ledger's notes table
    ┃  is where a note goes, and it is the orchestrator who decides that.
    ┃
    ┃  ## The ledger — `docs/audit.md` is yours
    ┃
    ┃  - The **slice T2a** row and the **cycle 1** and **cycle 2** rows are there. Add the
    ┃    **cycle 3** row.
    ┃  - Anything you file is **F-108** onward — F-106 and F-107 were closed in cycle 1's
    ┃    fix, and F-109 was closed in cycle 2's fix.
    ┃  - **Never sweep the ledger.**
    ┃
    ┃  ## The verdict
    ┃
    ┃  | | |
    ┃  |---|---|
    ┃  | No HIGH finding | `verdict: ok` |
    ┃  | Any HIGH finding | `verdict: changes-needed` |
    ┃
    ┃  **Do not invent reasons to fail.** A slice that does what it said is `ok`.
    ┃
    ┃  ## Hard rules
    ┃
    ┃  - Do **not** edit implementation code — **except** `docs/audit.md`, which is yours.
    ┃  - Revert every mutation you make, and **confirm the revert**.
    ┃  - No `git commit --no-verify`. No push. Do not touch other panes. Do not spawn agents.
    ┃  - **One commit** for the ledger. One commit for the receipt.
    ┃
    ┃  ## Notify the orchestrator
    ┃
    ┃  **The receipt comes FIRST** — the file cannot be lost, the nudge is only the bell.
    ┃
    ┃  ```bash
    ┃  mkdir -p work/inbox
    ┃  cat > work/inbox/reviewer-tui-t2a-cycle3.md <<EOF
    ┃  Reviewer: T2a cycle 3 done. Verdict: <ok|changes-needed>. Filed <N> HIGH blockers (<F-108.. or none>). F-109: <the mutation you took, what failed,
    ┃  and that it failed for the bound's width>. The fixture: <how the two widths differ, and the row counts>. The document: <what the table says now, and
    ┃  whether every row is true>. The scope: <what you checked>. The goldens: <unchanged / changed, and why>. The gate: <make check, and make check-race
    ┃  if you ran it>. Notes for the ledger: <any>.
    ┃  EOF
    ┃  git add work/inbox && git commit -q -m "inbox: reviewer reports tui-t2a cycle 3"
    ┃  herdr agent prompt "w1H:p1" "Reviewer: T2a cycle 3 done. Verdict: <ok|changes-needed>. Filed <N> HIGH blockers (<F-108.. or none>). F-109: <the
    ┃  mutation you took, what failed, and that it failed for the bound's width>. The fixture: <how the two widths differ, and the row counts>. The
    ┃  document: <what the table says now, and whether every row is true>. The scope: <what you checked>. The goldens: <unchanged / changed, and why>. The
    ┃  gate: <make check, and make check-race if you ran it>. Notes for the ledger: <any>."
    ┃  ```
    ┃
    ┃  If the nudge lands on a busy receiver, the receipt is the record — **write it anyway**.
    ┃
    ┃  **Nothing else waits on you.** After the nudge, end your turn.
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
    ┃  cat > work/inbox/reviewer-tui-t2a-cycle3.md <<EOF
    ┃  Reviewer: T2a cycle 3 done. Verdict: <ok|changes-needed>. Filed <N> HIGH blockers (<F-108.. or none>). F-109: <the mutation you took, what failed,
    ┃  and that it failed for the bound's width>. The fixture: <how the two widths differ, and the row counts>. The document: <what the table says now, and
    ┃  whether every row is true>. The scope: <what you checked>. The goldens: <unchanged / changed, and why>. The gate: <make check, and make check-race
    ┃  if you ran it>. Notes for the ledger: <any>.
    ┃  EOF
    ┃  git add work/inbox && git commit -q -m "inbox: reviewer reports tui-t2a cycle 3"
    ┃  herdr agent prompt "w1H:p1" "Reviewer: T2a cycle 3 done. Verdict: <ok|changes-needed>. Filed <N> HIGH blockers (<F-108.. or none>). F-109: <the
    ┃  mutation you took, what failed, and that it failed for the bound's width>. The fixture: <how the two widths differ, and the row counts>. The
    ┃  document: <what the table says now, and whether every row is true>. The scope: <what you checked>. The goldens: <unchanged / changed, and why>. The
    ┃  gate: <make check, and make check-race if you ran it>. Notes for the ledger: <any>."
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
    ┃  So you do not repeat the work:
    ┃
    ┃  - **F-109** — the reviewer's mutation (count the bound at the full-screen width)
    ┃    turns `internal/ui` **red**, with three messages: *the pane has 29 rows at 107
    ┃    columns; the full screen has 27 at 196* · *the bound is 27 rows; the pane has 29 rows
    ┃    at 107 columns* · *the scroll stopped at 26, want the pane's last row 28*. The
    ┃    tripwire fails **for the bound's width**.
    ┃  - **The document** — the wheel table's second row now reads *anywhere else, **while the
    ┃    list is on screen***, and the prose keeps the full-screen detail as the case the
    ┃    table carves out.
    ┃  - The gate: `make check` and `make check-race` green, the goldens unchanged, the tree
    ┃    clean.
    ┃

       + Thought · 154ms

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git log --oneline -8 && git status --short && echo CLEAN && git show --stat a332158 d8d7b95 22b9de6 | head
    ┃   -50
    ┃
    ┃  (42 earlier lines)
    ┃
    ┃      inbox: runner reports tui-t2a fix 2
    ┃
    ┃   work/inbox/runner-tui-t2a-fix2.md | 45 +++++++++++++++++++++++++++++++++++++++
    ┃   work/tui/README.md                |  2 +-
    ┃   2 files changed, 46 insertions(+), 1 deletion(-)
    ┃

       → Explored: 1 read

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git show d8d7b95 && git show a332158 -- internal/ui/mouse_test.go
    ┃
    ┃  (75 earlier lines)
    ┃  -    t.Errorf("the scroll stopped at %d, want %d", m.detailScroll, want)
    ┃  -  }
    ┃  -}
    ┃  -
    ┃   // TestThePositionFormSurvivesTheCommandLine covers the spelling on the command
    ┃   // line.
    ┃   //
    ┃

       + Thought · 3.2s

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && go test ./internal/ui -run 'TestTheScrollsBoundIsThePanesWidth' -v 2>&1 | grep -E "bound_test|^(---|ok|FAI
    ┃  L|PASS)"
    ┃
    ┃  bound_test.go:77: the pane has 29 rows at 107 columns; the full screen has 27 at 196
    ┃  --- PASS: TestTheScrollsBoundIsThePanesWidth (0.17s)
    ┃  PASS
    ┃  ok    github.com/lthiagol/maitre/internal/ui  (cached)
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && python3 - <<'EOF'
    ┃  p='internal/ui/model.go'…
    ┃
    ┃  (2 earlier lines)
    ┃      bound_test.go:81: the bound is 27 rows; the pane has 29 rows at 107 columns
    ┃      bound_test.go:90: the scroll stopped at 26, want the pane's last row 28
    ┃  FAIL
    ┃  FAIL  github.com/lthiagol/maitre/internal/ui  90.096s
    ┃  FAIL
    ┃  CLEAN
    ┃

       + Thought · 4.5s

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && grep -rn "wheel\|mouse" docs/*.md | grep -v "^docs/05-human-interface.md:2[6-9]" | cut -c1-160
    ┃
    ┃  docs/audit.md:65:| Slice T2a | The mouse, part 1: the mouse is on; the render builds the hit map; the wheel scrolls the list and the detail | The
    ┃  mouse is on **
    ┃  docs/audit.md:66:| Slice T2a · cycle 1 | The review of the slice: the map, the wheel, and six mutations | The slice's shape held: the map **is**
    ┃  built by the re
    ┃  docs/audit.md:67:| Slice T2a · cycle 2 | The review of the fix cycle: the split detail, the position form, the absent region, and the **reversed**
    ┃  test | F-106,
    ┃  docs/audit.md:119:| **F-109** | medium | **The bound's fix is not held by any test — the suite is green with the bound back at the full-screen
    ┃  width.** F-106's
    ┃  docs/audit.md:213:| F-106 | medium | **The wheel over the detail pane did not scroll the detail in the split layout, and the detail's row count had
    ┃  two answers.
    ┃  docs/audit.md:214:| F-107 | medium | **The scripted mouse's position form was unreachable through `maitre --keys`, and the doc printed the command
    ┃  that failed.*
    ┃  docs/audit.md:215:| F-108 | low | **The zero box was a one-cell region at the origin**, so the top-left corner of a screen with no detail pane
    ┃  looked like it ha
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && sed -n '262,295p' docs/05-human-interface.md
    ┃
    ┃  (27 earlier lines)
    ┃  The coordinates are the terminal's own — zero-based, with (0, 0) at the upper-left
    ┃  corner. The position is **bracketed**, because the steps are separated by commas and
    ┃  the comma inside a position must not be one.
    ┃
    ┃  **The theme is the chef's to pick.** The Settings tab lists the palettes under a
    ┃  **THEME** heading; the cursor moves over them and `enter` applies the one under it. The
    ┃  active theme is marked, so the chef can see which one is in force. The choice is for the
    ┃

       + Thought · 1.1s

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && make build 2>&1 | tail -1 && cd /tmp/rev2 && M=/Users/thiago/Code/Pessoal/maitre/bin/maitre
    ┃  $M --render --width 200 --keys "enter,esc" --plain > /tmp/q0.txt…
    ┃
    ┃  (4 earlier lines)
    ┃
    ┃  == the full-screen detail, the wheel over the header ==
    ┃  IDENTICAL — nothing happens (no list on screen)
    ┃
    ┃  == the full-screen detail, the wheel over the detail ==
    ┃  CHANGED — the detail scrolls
    ┃

       + Thought · 1.8s

    ┃
    ┃  $ cd /tmp && diff /tmp/Rq0.txt /tmp/Rq1.txt | head -6; echo "..."; echo "== the first detail row, before and after =="; sed -n '5p' /tmp/Rq0.txt; se
    ┃  d -n '5p' /tmp/Rq1.txt
    ┃
    ┃  (4 earlier lines)
    ┃  ---
    ┃  >   RCP05 · Sharpen the evidence contract
    ┃  ...
    ┃  == the first detail row, before and after ==
    ┃    RCP02 · Sign off the menu
    ┃    RCP05 · Sharpen the evidence contract
    ┃

       + Thought · 5.8s

       → Explored: 1 read

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && sed -n '13,18p' work/tui/README.md | cut -c1-260 && git diff fddaa0d..HEAD --stat
    ┃
    ┃  (6 earlier lines)
    ┃   docs/05-human-interface.md           |   6 +-
    ┃   internal/ui/bound_test.go            |  92 ++++++++++++++++++++++++++
    ┃   internal/ui/mouse_test.go            |  29 ---------
    ┃   work/briefs/review-tui-t2a-cycle3.md | 122 +++++++++++++++++++++++++++++++++++
    ┃   work/inbox/runner-tui-t2a-fix2.md    |  45 +++++++++++++
    ┃   work/tui/README.md                   |   2 +-
    ┃   6 files changed, 263 insertions(+), 33 deletions(-)
    ┃

       + Thought · 2.5s

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && make check 2>&1 | tail -8 && make goldens 2>&1 | tail -2 && git status --short && echo "GOLDENS: no diff"
    ┃  && echo "=== the click ===" && grep -rn "click\|MouseButton" internal/ cmd/ --include="*.go" | grep -v "_test.go"; echo "=== core / cli / doc…
    ┃
    ┃  (10 earlier lines)
    ┃  GOLDENS: no diff
    ┃  === the click ===
    ┃  internal/ui/keys.go:50:// clicks against a map from a frame that is not on screen clicks the wrong thing.
    ┃  internal/ui/mouse.go:22:// The console answers the wheel and leaves the rest for the click: a click that
    ┃  === core / cli / docs03 since cycle 1 ===
    ┃  (empty = untouched)
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && make check-race 2>&1 | tail -6 && git status --short && echo CLEAN && grep -n "DomainVersion = " internal/
    ┃  core/lifecycle.go && grep -c "" internal/ui/testdata/*.golden
    ┃
    ┃  (4 earlier lines)
    ┃
    ┃  ==> tidy-check
    ┃  CLEAN
    ┃  225:const DomainVersion = 7
    ┃  internal/ui/testdata/console-recipes.golden:40
    ┃  internal/ui/testdata/console-stress.golden:40
    ┃

       + Thought · 9.0s

    ┃
    ┃  ← Edit docs/audit.md
    ┃
    ┃    115   |---|---|---|---|                                                   115   |---|---|---|---|
    ┃    116   | **F-59** | medium | **`core.Audit.Decisions` has a shape and no   116   | **F-59** | medium | **`core.Audit.Decisions` has a shape and no
    ┃          writer** — it is where a decision *about the plan* would live, and        writer** — it is where a decision *about the plan* would live, and
    ┃          nothing writes it. `kitchen/decisions.json` is named by the plan          nothing writes it. `kitchen/decisions.json` is named by the plan
    ┃          format and opened by nothing. The remedy is on the table in [06-          format and opened by nothing. The remedy is on the table in [06-
    ┃          proposals.md](06-proposals.md) § 2. | `internal/core/model.go:270`;       proposals.md](06-proposals.md) § 2. | `internal/core/model.go:
    ┃          `internal/core/store.go:24` |                                             270`; `internal/core/store.go:24` |
    ┃    117   | **F-60** | medium | **`core.Questions` has a writer-less shape,   117   | **F-60** | medium | **`core.Questions` has a writer-less shape,
    ┃          so the `questions.resolved` gate can never fire** — nothing ever          so the `questions.resolved` gate can never fire** — nothing ever
    ┃          sets a question, so the list is always empty and the gate always          sets a question, so the list is always empty and the gate always
    ┃          passes. A gate nobody can fail is not a guard. | `internal/core/          passes. A gate nobody can fail is not a guard. | `internal/core/
    ┃          model.go:31`; `internal/core/lifecycle.go:197` |                          model.go:31`; `internal/core/lifecycle.go:197` |
    ┃    118   | **F-61** | medium | **The `courses/` directory is named by the    118   | **F-61** | medium | **The `courses/` directory is named by the
    ┃          plan format, created by `Init`, and never written** — a directory         plan format, created by `Init`, and never written** — a directory
    ┃          no writer touches. Per [06-proposals.md](06-proposals.md) § 1 the         no writer touches. Per [06-proposals.md](06-proposals.md) § 1 the
    ┃          honest answer is to stop naming it: a course is the relation, not a       honest answer is to stop naming it: a course is the relation, not
    ┃          file. | `internal/core/store.go:18,63`; [docs/03 §6](03-domain-           a file. | `internal/core/store.go:18,63`; [docs/03 §6](03-domain-
    ┃          model.md) |                                                               model.md) |
    ┃    119 - | **F-109** | medium | **The bound's fix is not held by any test —
    ┃          the suite is green with the bound back at the full-screen width.**
    ┃          F-106's second half asked for **one** answer for the detail's row
    ┃          count, and the code has it: `detailLineCount()` counts at
    ┃          `detailPaneWidth()`. But `TestTheScrollsBoundIsThePanesWidth` asks
    ┃          whether the bound equals the pane's row count, and the **seeded**
    ┃          recipe's detail has the same row count at both widths — so the
    ┃          assertion holds whatever the bound is. The reviewer's mutation —
    ┃          `detailLineCount()` put back on `m.width - 4`, the **full-screen**
    ┃          width, which is the pre-fix defect — leaves the **whole `internal/
    ┃          ui` suite green**. The consequence is real: with the bound at the
    ┃          full-screen width the split pane's last rows are unreachable (a
    ┃          probe at 200×30 with wrapping prose: the pane has 41 rows, the
    ┃          bound 36 — five rows the scroll can never reach; with the fix, 41
    ┃          and 41). The guard needs a fixture whose detail wraps, or an
    ┃          assertion that fails when the two widths give different counts. |
    ┃          `internal/ui/mouse_test.go` — `TestTheScrollsBoundIsThePanesWidth`;
    ┃          `internal/ui/model.go` — `detailLineCount` |
    ┃    120   | F-56 | low | **The poll's fingerprint is blind to a write that    119   | F-56 | low | **The poll's fingerprint is blind to a write that
    ┃          keeps the mtime and the size.** It is path + mtime + size per file,       keeps the mtime and the size.** It is path + mtime + size per
    ┃          so two different contents can share one fingerprint. The tool's own       file, so two different contents can share one fingerprint. The
    ┃          writes always change the mtime, so in practice only a hand edit           tool's own writes always change the mtime, so in practice only a
    ┃          that restores the mtime is missed — and the chef's `r` reloads it.        hand edit that restores the mtime is missed — and the chef's `r`
    ┃          | `internal/core/store.go` — `DiskFingerprint`. |                         reloads it. | `internal/core/store.go` — `DiskFingerprint`. |
    ┃    121   | P-8 | — | `--depends-on=` writes `depends_on: [""]` — inert, the  120   | P-8 | — | `--depends-on=` writes `depends_on: [""]` — inert, the
    ┃          same as `--in-scope=`. | Recorded so nobody calls it a regression         same as `--in-scope=`. | Recorded so nobody calls it a regression
    ┃          later. Not a defect. |                                                    later. Not a defect. |
    ┃    122                                                                       121
    ┃    123   ### Notes left for a later round                                    122   ### Notes left for a later round
    ┃
    ┃

    ┃
    ┃  ← Edit docs/audit.md
    ┃
    ┃    211   | F-105 | low | **The state row's count and its list disagree** —   211   | F-105 | low | **The state row's count and its list disagree** —
    ┃          the row said *12 commits* and named nine, and the three it did not        the row said *12 commits* and named nine, and the three it did not
    ┃          name were the three commits that wrote the row. *Closed in* `work/        name were the three commits that wrote the row. *Closed in* `work/
    ┃          tui/README.md` — the shape is one that cannot drift: **the count is       tui/README.md` — the shape is one that cannot drift: **the count
    ┃          the number of ids the row names**, and the commits that wrote the         is the number of ids the row names**, and the commits that wrote
    ┃          row are not among them (a commit cannot carry its own hash). The          the row are not among them (a commit cannot carry its own hash).
    ┃          row now says *11 commits* and names **eleven real ids** — every one       The row now says *11 commits* and names **eleven real ids** —
    ┃          of them in the log — and the **protocol section** carries the shape       every one of them in the log — and the **protocol section**
    ┃          in one line, so the next writer knows it. | `work/tui/README.md` —        carries the shape in one line, so the next writer knows it. |
    ┃          the T1 row and the protocol section |                                     `work/tui/README.md` — the T1 row and the protocol section |
    ┃    212   | F-106 | medium | **The wheel over the detail pane did not scroll  212   | F-106 | medium | **The wheel over the detail pane did not scroll
    ┃          the detail in the split layout, and the detail's row count had two        the detail in the split layout, and the detail's row count had two
    ┃          answers.** The offset was applied in two places and the split             answers.** The offset was applied in two places and the split
    ┃          layout read neither; the bound was counted at the full-screen width       layout read neither; the bound was counted at the full-screen
    ┃          while the split pane draws at a narrower one, so the scroll stopped       width while the split pane draws at a narrower one, so the scroll
    ┃          early. *Closed in* `internal/ui/detail.go` — the offset is applied        stopped early. *Closed in* `internal/ui/detail.go` — the offset is
    ┃          **once**, where the detail is drawn (`detailScrollRows`), and             applied **once**, where the detail is drawn (`detailScrollRows`),
    ┃          `detailPaneWidth` is the **one** answer for the pane's width, which       and `detailPaneWidth` is the **one** answer for the pane's width,
    ┃          `renderSplit`, `renderDetailFull` and `detailLineCount` all ask.          which `renderSplit`, `renderDetailFull` and `detailLineCount` all
    ┃          *The bound's code is fixed, but the test that holds it does not           ask. *The bound's code is fixed, but the test that holds it does
    ┃          bite — that is F-109.* | `internal/ui/detail.go` —                        not bite — that is F-109.* | `internal/ui/detail.go` —
    ┃          `detailScrollRows`, `detailPaneWidth`; `internal/ui/views.go` —           `detailScrollRows`, `detailPaneWidth`; `internal/ui/views.go` —
    ┃          `renderSplit`, `renderDetailFull`; `internal/ui/model.go` —               `renderSplit`, `renderDetailFull`; `internal/ui/model.go` —
    ┃          `detailLineCount`; `internal/ui/mouse_test.go` —                          `detailLineCount`; `internal/ui/mouse_test.go` —
    ┃          `TestTheWheelScrollsTheSplitDetail` (the reviewer's mutation — the        `TestTheWheelScrollsTheSplitDetail` (the reviewer's mutation — the
    ┃          offset taken out of the split layout — fails it) |                        offset taken out of the split layout — fails it) |
    ┃    213   | F-107 | medium | **The scripted mouse's position form was         213   | F-107 | medium | **The scripted mouse's position form was
    ┃          unreachable through `maitre --keys`, and the doc printed the              unreachable through `maitre --keys`, and the doc printed the
    ┃          command that failed.** *Closed in* `internal/ui/mouse.go` — the           command that failed.** *Closed in* `internal/ui/mouse.go` — the
    ┃          position is **bracketed** (`wheel:down@(12,5)`), and `SplitSteps` —       position is **bracketed** (`wheel:down@(12,5)`), and `SplitSteps`
    ┃          the splitter that knows the grammar — lives **beside** the parser         — the splitter that knows the grammar — lives **beside** the
    ┃          that reads it, so the command line and the driver cannot disagree         parser that reads it, so the command line and the driver cannot
    ┃          about what a step is; `cmd/maitre/main.go` — the flag calls `ui.          disagree about what a step is; `cmd/maitre/main.go` — the flag
    ┃          SplitSteps` and nothing else; `docs/05-human-interface.md` — the          calls `ui.SplitSteps` and nothing else; `docs/05-human-interface.
    ┃          doc prints the form that works. The old spelling is refused with a        md` — the doc prints the form that works. The old spelling is
    ┃          message a reader can act on (*malformed mouse step "wheel:                refused with a message a reader can act on (*malformed mouse step
    ┃          down@12"*), and `--keys "tab,down,space"` still works. | `internal/       "wheel:down@12"*), and `--keys "tab,down,space"` still works. |
    ┃          ui/mouse.go` — `SplitSteps`, `parseCoords`, `isMouseStep`; `cmd/          `internal/ui/mouse.go` — `SplitSteps`, `parseCoords`,
    ┃          maitre/main.go` — the `--keys` split; `docs/05-human-interface.md`        `isMouseStep`; `cmd/maitre/main.go` — the `--keys` split; `docs/
    ┃          — the scripted-mouse block; `internal/ui/mouse_test.go` —                 05-human-interface.md` — the scripted-mouse block; `internal/ui/
    ┃          `TestThePositionFormSurvivesTheCommandLine`,                              mouse_test.go` — `TestThePositionFormSurvivesTheCommandLine`,
    ┃          `TestThePositionFormIsDrivable` (the reviewer's mutation — a              `TestThePositionFormIsDrivable` (the reviewer's mutation — a
    ┃          splitter that cuts at the comma — fails both) |                           splitter that cuts at the comma — fails both) |
    ┃    214   | F-108 | low | **The zero box was a one-cell region at the         214   | F-108 | low | **The zero box was a one-cell region at the
    ┃          origin**, so the top-left corner of a screen with no detail pane          origin**, so the top-left corner of a screen with no detail pane
    ┃          looked like it had one, and the cleared map answered for the list.        looked like it had one, and the cleared map answered for the list.
    ┃          *Closed in* `internal/ui/hit.go` — `hitBox` carries a `there` flag,       *Closed in* `internal/ui/hit.go` — `hitBox` carries a `there`
    ┃          `newHitBox` builds a region that is on screen, and the zero box           flag, `newHitBox` builds a region that is on screen, and the zero
    ┃          answers for nothing. The wheel at the origin over the list now            box answers for nothing. The wheel at the origin over the list now
    ┃          scrolls the list. | `internal/ui/hit.go` — `hitBox`, `newHitBox`;         scrolls the list. | `internal/ui/hit.go` — `hitBox`, `newHitBox`;
    ┃          `internal/ui/mouse_test.go` — `TestTheOriginIsNotAPhantomDetail`,         `internal/ui/mouse_test.go` — `TestTheOriginIsNotAPhantomDetail`,
    ┃          `TestAClearedMapAnswersNothing` (the reviewer's mutation — the            `TestAClearedMapAnswersNothing` (the reviewer's mutation — the
    ┃          `there` guard taken out — fails three tests) |                            `there` guard taken out — fails three tests) |
    ┃                                                                              215 + | F-109 | medium | **The bound's width was held by nothing** — the
    ┃                                                                                    test compared two values that **both** went through
    ┃                                                                                    `detailPaneWidth()`, and the fixture's detail was too short for
    ┃                                                                                    the two widths to differ, so the bound could be counted at the
    ┃                                                                                    full screen's width and the comparison would still hold. *Closed
    ┃                                                                                    in* `internal/ui/bound_test.go` — a **new** test: the pane's width
    ┃                                                                                    is computed **in the test** (`m.width - m.listColumn() - 3`) and
    ┃                                                                                    the pane's rows are counted **at that width**, not through the
    ┃                                                                                    console's answer — a tripwire that asks the code under test for
    ┃                                                                                    the answer is not a tripwire. The fixture is a recipe whose detail
    ┃                                                                                    **wraps** — **29 rows at the pane's 107 columns, 27 at the full
    ┃                                                                                    screen's 196** — and the test says so out loud if the two widths
    ┃                                                                                    ever agree again. The reviewer's mutation — the bound counted at
    ┃                                                                                    the full-screen width again — turns `internal/ui` **red** in that
    ┃                                                                                    one test, and **for the bound's width**: *the bound is 27 rows;
    ┃                                                                                    the pane has 29 rows at 107 columns* · *the scroll stopped at 26,
    ┃                                                                                    want the pane's last row 28*. | `internal/ui/bound_test.go` —
    ┃                                                                                    `TestTheScrollsBoundIsThePanesWidth`, `longDetailModel`;
    ┃                                                                                    `internal/ui/model.go` — `detailLineCount` |
    ┃    215                                                                       216
    ┃    216   **Cycles 1–7's own findings are not recoverable.** Their per-       217   **Cycles 1–7's own findings are not recoverable.** Their per-
    ┃          finding reports lived in a                                                finding reports lived in a
    ┃    217   scratch directory that was wiped mid-cycle. The *fixes* are all in  218   scratch directory that was wiped mid-cycle. The *fixes* are all in
    ┃          the repository, and                                                       the repository, and
    ┃    218   cycle 8's report re-verified fifteen of them with evidence — those  219   cycle 8's report re-verified fifteen of them with evidence — those
    ┃          are the rows above.                                                       are the rows above.
    ┃

    ┃
    ┃  ← Edit docs/audit.md
    ┃
    ┃     64   | Slice T1 · cycle 3 | The **last** cycle of T1: the blocker, the    64   | Slice T1 · cycle 3 | The **last** cycle of T1: the blocker, the
    ┃          two notes, and the question that is the chef's | F-103..F-105 are         two notes, and the question that is the chef's | F-103..F-105 are
    ┃          closed, each proved by a mutation and each reverted: the rule taken       closed, each proved by a mutation and each reverted: the rule
    ┃          out of the one copy turns the suite **red** in both `internal/core`       taken out of the one copy turns the suite **red** in both
    ┃          and `internal/ui` — **three** tests, as the runner claimed; the           `internal/core` and `internal/ui` — **three** tests, as the runner
    ┃          `sanitize` call taken out of `humanPath` fails the root-row test          claimed; the `sanitize` call taken out of `humanPath` fails the
    ┃          twice over; the state row's count is the number of ids the row            root-row test twice over; the state row's count is the number of
    ┃          names, and the **protocol section** carries the shape in one line.        ids the row names, and the **protocol section** carries the shape
    ┃          The two surfaces agree in **all three states** — the plan names           in one line. The two surfaces agree in **all three states** — the
    ┃          itself, the file names it empty, the file is missing — and the            plan names itself, the file names it empty, the file is missing —
    ┃          CLI's behaviour **changed** to match the console's, which is the          and the CLI's behaviour **changed** to match the console's, which
    ┃          point of F-103. The rule is written down in the **code** (`core.          is the point of F-103. The rule is written down in the **code** (
    ┃          Store.projectName`'s doc comment) and **nowhere in the spec**:            `core.Store.projectName`'s doc comment) and **nowhere in the
    ┃          `docs/03` § 6 names `project.json` as *project meta · charter ·           spec**: `docs/03` § 6 names `project.json` as *project meta ·
    ┃          settings* and says nothing about the name's fallback, and `docs/05`       charter · settings* and says nothing about the name's fallback,
    ┃          does not either — where it belongs is the **chef's** call, and            and `docs/05` does not either — where it belongs is the **chef's**
    ┃          **nothing was written into `docs/03`**. The goldens did not move,         call, and **nothing was written into `docs/03`**. The goldens did
    ┃          the budget is still **0 / 0 / 28**, `core.DomainVersion` is still         not move, the budget is still **0 / 0 / 28**, `core.DomainVersion`
    ┃          **7**, the seven tabs kept their names, the CLI's surface did not         is still **7**, the seven tabs kept their names, the CLI's surface
    ┃          change, and the tree was left exactly as the runner left it — **no        did not change, and the tree was left exactly as the runner left
    ┃          blockers filed** |                                                        it — **no blockers filed** |
    ┃     65   | Slice T2a | The mouse, part 1: the mouse is on; the render builds  65   | Slice T2a | The mouse, part 1: the mouse is on; the render
    ┃          the hit map; the wheel scrolls the list and the detail | The mouse        builds the hit map; the wheel scrolls the list and the detail |
    ┃          is on **every** view the console draws — the console's, the error's       The mouse is on **every** view the console draws — the console's,
    ┃          and the opening one's — and the error and opening views **clear the       the error's and the opening one's — and the error and opening
    ┃          map**, so no event can land on a tab that is not on screen. The map       views **clear the map**, so no event can land on a tab that is not
    ┃          is built **by the render**: `tabBar` emits the tab boxes as it            on screen. The map is built **by the render**: `tabBar` emits the
    ┃          draws them, `hitMapOf` is called from the three render functions          tab boxes as it draws them, `hitMapOf` is called from the three
    ┃          and from nowhere else, the banner pushes the list's rows down in          render functions and from nowhere else, the banner pushes the
    ┃          the map as it does on the screen, and the geometry constants (            list's rows down in the map as it does on the screen, and the
    ┃          `headerRows`, `footerRows`, `splitThreshold`) are read where they         geometry constants (`headerRows`, `footerRows`, `splitThreshold`)
    ┃          are written. The wheel over the list drives `m.move` — the                are read where they are written. The wheel over the list drives
    ┃          keyboard's own move — and a scripted `wheel:down` and three `down`s       `m.move` — the keyboard's own move — and a scripted `wheel:down`
    ┃          land on **byte-identical** frames. The wheel over the detail drives       and three `down`s land on **byte-identical** frames. The wheel
    ┃          `m.scrollDetail`, but **the split layout never renders that               over the detail drives `m.scrollDetail`, but **the split layout
    ┃          offset**: at ≥160 columns — where the console draws the list and          never renders that offset**: at ≥160 columns — where the console
    ┃          the detail side by side, and `ModeSplit` is the default — the pane        draws the list and the detail side by side, and `ModeSplit` is the
    ┃          does not move, and the clamp measures the detail at a different           default — the pane does not move, and the clamp measures the
    ┃          width than the pane draws it. The scripted mouse's **position             detail at a different width than the pane draws it. The scripted
    ┃          form** cannot be reached through `maitre --keys`, and the doc             mouse's **position form** cannot be reached through `maitre --
    ┃          prints the command line that fails. The scope held: **no click            keys`, and the doc prints the command line that fails. The scope
    ┃          exists**, on purpose; `internal/core` and `internal/cli` are              held: **no click exists**, on purpose; `internal/core` and
    ┃          untouched; `core.DomainVersion` is still **7**; the budget is still       `internal/cli` are untouched; `core.DomainVersion` is still **7**;
    ┃          **0 / 0 / 28**; the seven tabs kept their names; the goldens did          the budget is still **0 / 0 / 28**; the seven tabs kept their
    ┃          not move — **two blockers filed** (F-106, F-107) |                        names; the goldens did not move — **two blockers filed** (F-106,
    ┃                                                                                    F-107) |
    ┃     66   | Slice T2a · cycle 1 | The review of the slice: the map, the        66   | Slice T2a · cycle 1 | The review of the slice: the map, the
    ┃          wheel, and six mutations | The slice's shape held: the map **is**         wheel, and six mutations | The slice's shape held: the map **is**
    ┃          built by the render (`tabBar` emits the tab boxes, `hitMapOf` is          built by the render (`tabBar` emits the tab boxes, `hitMapOf` is
    ┃          called from the three render functions and nowhere else, the banner       called from the three render functions and nowhere else, the
    ┃          pushes the list's rows down in the map, and the geometry constants        banner pushes the list's rows down in the map, and the geometry
    ┃          are read where they are written), the wheel **is** the keyboard's         constants are read where they are written), the wheel **is** the
    ┃          scroll (a scripted `wheel:down` and three `down`s land on byte-           keyboard's scroll (a scripted `wheel:down` and three `down`s land
    ┃          identical frames), and all **six** mutations failed the suite — the       on byte-identical frames), and all **six** mutations failed the
    ┃          mouse mode off the error view, the render not emitting the map, a         suite — the mouse mode off the error view, the render not emitting
    ┃          second scroll variable, the wheel over the header, the wheel under        the map, a second scroll variable, the wheel over the header, the
    ┃          an open prompt, and the scripted step acting without drawing first.       wheel under an open prompt, and the scripted step acting without
    ┃          Two blockers filed: **F-106** (the wheel over the detail does             drawing first. Two blockers filed: **F-106** (the wheel over the
    ┃          nothing in the split layout) and **F-107** (the scripted position         detail does nothing in the split layout) and **F-107** (the
    ┃          form is unreachable through `maitre --keys`, and the doc prints the       scripted position form is unreachable through `maitre --keys`, and
    ┃          command that fails). The zero box's one cell at the origin was            the doc prints the command that fails). The zero box's one cell at
    ┃          noted as **F-108** and ruled into the fix cycle. The goldens did          the origin was noted as **F-108** and ruled into the fix cycle.
    ┃          not move; `internal/core` was untouched; `core.DomainVersion` was         The goldens did not move; `internal/core` was untouched; `core.
    ┃          still **7**; the budget was still **0 / 0 / 28** — **verdict:             DomainVersion` was still **7**; the budget was still **0 / 0 /
    ┃          changes-needed** |                                                        28** — **verdict: changes-needed** |
    ┃     67   | Slice T2a · cycle 2 | The review of the fix cycle: the split       67   | Slice T2a · cycle 2 | The review of the fix cycle: the split
    ┃          detail, the position form, the absent region, and the **reversed**        detail, the position form, the absent region, and the **reversed**
    ┃          test | F-106, F-107 and F-108 are closed, each proved by a                test | F-106, F-107 and F-108 are closed, each proved by a
    ┃          mutation: taking the offset out of the split layout fails                 mutation: taking the offset out of the split layout fails
    ┃          `TestTheWheelScrollsTheSplitDetail`; a splitter that cuts a               `TestTheWheelScrollsTheSplitDetail`; a splitter that cuts a
    ┃          bracketed position at its comma fails the two position tests;             bracketed position at its comma fails the two position tests;
    ┃          taking the `there` guard out of `hitBox` fails three tests. The           taking the `there` guard out of `hitBox` fails three tests. The
    ┃          offset is applied in **one** place and the pane's width has **one**       offset is applied in **one** place and the pane's width has
    ┃          answer (`detailPaneWidth`), and the wheel at the origin over the          **one** answer (`detailPaneWidth`), and the wheel at the origin
    ┃          list scrolls the list. **But the bound's guard does not bite**: the       over the list scrolls the list. **But the bound's guard does not
    ┃          bound put back on the full-screen width leaves the **whole                bite**: the bound put back on the full-screen width leaves the
    ┃          `internal/ui` suite green** — **one blocker filed** (F-109). The          **whole `internal/ui` suite green** — **one blocker filed** (F-
    ┃          **reversed** test is judged **sound**: the rule is one rule — the         109). The **reversed** test is judged **sound**: the rule is one
    ┃          detail answers over the detail, the list answers the rest — stated        rule — the detail answers over the detail, the list answers the
    ┃          in one place in the code and in the doc; the full-screen detail           rest — stated in one place in the code and in the doc; the full-
    ┃          keeps the *no list, nothing happens* answer, pinned by                    screen detail keeps the *no list, nothing happens* answer, pinned
    ┃          `TestTheWheelOverTheFullScreenHeaderDoesNothing`; and the old rule        by `TestTheWheelOverTheFullScreenHeaderDoesNothing`; and the old
    ┃          and the fix cycle's own F-108 evidence (*a wheel at the origin over       rule and the fix cycle's own F-108 evidence (*a wheel at the
    ┃          the list scrolls the list*) cannot both hold, because the origin          origin over the list scrolls the list*) cannot both hold, because
    ┃          **is** the header's first cell. The scope held: **no click                the origin **is** the header's first cell. The scope held: **no
    ┃          exists**, on purpose; `internal/core`, `internal/cli` and `docs/03`       click exists**, on purpose; `internal/core`, `internal/cli` and
    ┃          are untouched; `core.DomainVersion` is still **7**; the budget is         `docs/03` are untouched; `core.DomainVersion` is still **7**; the
    ┃          still **0 / 0 / 28**; the seven tabs kept their names; the goldens        budget is still **0 / 0 / 28**; the seven tabs kept their names;
    ┃          did not move; `make check` and `make check-race` are green. |             the goldens did not move; `make check` and `make check-race` are
    ┃                                                                                    green. |
    ┃                                                                               68 + | Slice T2a · cycle 3 | The **last** cycle of T2a: the bound's
    ┃                                                                                    tripwire, the wheel table, and the scope | F-109 is closed: the
    ┃                                                                                    tripwire **bites** — the reviewer's mutation, the bound counted at
    ┃                                                                                    the full-screen width again, turns `internal/ui` **red** in that
    ┃                                                                                    one test and **for the bound's width** (*the bound is 27 rows; the
    ┃                                                                                    pane has 29 rows at 107 columns* · *the scroll stopped at 26, want
    ┃                                                                                    the pane's last row 28*). The new test does not ask the code under
    ┃                                                                                    test for the answer — the pane's width is computed **in the test**
    ┃                                                                                    — and its fixture **makes the two widths differ** (29 rows at 107
    ┃                                                                                    columns, 27 at 196) and says so out loud if they ever agree again.
    ┃                                                                                    The document's wheel table reads whole: the second row carries
    ┃                                                                                    *"while the list is on screen"*, and the prose names the case the
    ┃                                                                                    table carves out; every row is true, proved by command — over the
    ┃                                                                                    header in the split layout the list scrolls, and in the full-
    ┃                                                                                    screen detail the wheel over the header does **nothing** (`enter`
    ┃                                                                                    and `enter,wheel:down@(150,0)` are byte-identical) while the wheel
    ┃                                                                                    over the detail scrolls it. The wheel's **fallback is untouched**
    ┃                                                                                    since cycle 2. The scope held: **no click**; `internal/core`,
    ┃                                                                                    `internal/cli` and `docs/03` untouched; `core.DomainVersion` still
    ┃                                                                                    **7**; the budget still **0 / 0 / 28**; the seven tabs kept their
    ┃                                                                                    names; the goldens did not move; `make check` and `make check-
    ┃                                                                                    race` green — **no blockers filed** |
    ┃     68   | The rounds inside a cycle | Fix → verify → back to the executor →  69   | The rounds inside a cycle | Fix → verify → back to the executor
    ┃          repeat, until the verifier is satisfied | Three rounds in cycle 8:        → repeat, until the verifier is satisfied | Three rounds in cycle
    ┃          the first fix introduced three new mediums, the second closed them        8: the first fix introduced three new mediums, the second closed
    ┃          |                                                                         them |
    ┃     69                                                                        70
    ┃     70   ### The discipline every round runs                                  71   ### The discipline every round runs
    ┃     71                                                                        72
    ┃
    ┃
    ┃

       Now the two notes in the ledger's notes table:

    ┃
    ┃  ← Edit docs/audit.md
    ┃
    ┃    134   | A golden whose name carries a **space** | The scripts walk the    134   | A golden whose name carries a **space** | The scripts walk the
    ┃          goldens by word, so a name with a space in it would be mishandled.        goldens by word, so a name with a space in it would be mishandled.
    ┃          There is none today. |                                                    There is none today. |
    ┃    135   | `docs/05`'s window-title paragraph is **held by nothing** | The   135   | `docs/05`'s window-title paragraph is **held by nothing** | The
    ┃          paragraph that says what the window title carries (F-102) has no          paragraph that says what the window title carries (F-102) has no
    ┃          wire: deleting it leaves the doc check green. A doc fix with              wire: deleting it leaves the doc check green. A doc fix with
    ┃          nothing holding it is a doc fix that can be undone by the next            nothing holding it is a doc fix that can be undone by the next
    ┃          reader. |                                                                 reader. |
    ┃    136   | The plan's name rule is **nowhere in the spec** | **Closed by the 136   | The plan's name rule is **nowhere in the spec** | **Closed by
    ┃          chef's ruling (2026-10-30, Q7):** the rule lives in `docs/03` § 6,        the chef's ruling (2026-10-30, Q7):** the rule lives in `docs/03`
    ┃          beside the plan format — it is about what the reader answers from         § 6, beside the plan format — it is about what the reader answers
    ┃          `project.json`, and § 6 is where the reader looks. It is **not** in       from `project.json`, and § 6 is where the reader looks. It is
    ┃          § 8: the name is a *stored* fact, not a derived view. A wire (            **not** in § 8: the name is a *stored* fact, not a derived view. A
    ┃          `TestThePlanFormatStatesTheNameRule`) holds the sentence, and             wire (`TestThePlanFormatStatesTheNameRule`) holds the sentence,
    ┃          `TestTheNameRuleTripwireFailsOnRealDrift` proves the wire can fail.       and `TestTheNameRuleTripwireFailsOnRealDrift` proves the wire can
    ┃          |                                                                         fail. |
    ┃    137   | A test assertion that **cannot fail** | `TestTheRootRowIsOneRow`  137   | A test assertion that **cannot fail** | `TestTheRootRowIsOneRow`
    ┃          splits the rows on `\n` and then checks each line for a `\n` — a          splits the rows on `\n` and then checks each line for a `\n` — a
    ┃          line that came out of a split cannot carry one. The test still            line that came out of a split cannot carry one. The test still
    ┃          bites through its other two assertions (*the root row does not            bites through its other two assertions (*the root row does not
    ┃          carry the sanitised path*, *the root row carries a newline*), so          carry the sanitised path*, *the root row carries a newline*), so
    ┃          the rule is held; the loop is decoration. |                               the rule is held; the loop is decoration. |
    ┃                                                                              138 + | The render **writes** the scroll's offset | `detailScrollRows`
    ┃                                                                                    clamps `m.detailScroll` — and writes it — while the frame is being
    ┃                                                                                    drawn. It moved there from `renderDetailFull`, which clamped it
    ┃                                                                                    inline; it is **not new**, and the render already writes `m.hits`.
    ┃                                                                                    Worth knowing before T2b adds more to the render path. |
    ┃                                                                              139 + | The wheel over the list moves the **selection** | In the split
    ┃                                                                                    layout the wheel over the list is the keyboard's `down`, so the
    ┃                                                                                    selection moves and the detail pane shows **another recipe from
    ┃                                                                                    its top** — the detail does not scroll. That is the console's
    ┃                                                                                    model (one answer, not two), but a reader of the pane may expect
    ┃                                                                                    the recipe to stay put. |
    ┃    138                                                                       140
    ┃    139   ## Closed findings                                                  141   ## Closed findings
    ┃    140                                                                       142
    ┃    141   | # | Severity | Closed in | Held by |                              143   | # | Severity | Closed in | Held by |
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git diff --stat && grep -n "F-10[6-9]" docs/audit.md | cut -c1-60 && make check 2>&1 | tail -4 && git stat
    ┃  us --short
    ┃
    ┃  (9 earlier lines)
    ┃  218:| F-109 | medium | **The bound's width was held by nothi
    ┃    internal/ui: 28 external modules (ceiling 30)
    ┃  /Library/Developer/CommandLineTools/usr/bin/make tidy-check
    ┃
    ┃  ==> tidy-check
    ┃   M docs/audit.md
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && git add docs/audit.md && git commit -q -m "docs: T2a's last cycle — the bound's tripwire bites, and the sl
    ┃  ice is done…
    ┃
    ┃  ac1bc3c docs: T2a's last cycle — the bound's tripwire bites, and the slice is done
    ┃  26a94e3 work: T2a's cycle 3 review brief - the last cycle
    ┃  CLEAN
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && mkdir -p work/inbox && cat > work/inbox/reviewer-tui-t2a-cycle3.md <<'EOF'
    ┃  Reviewer: T2a cycle 3 done. Verdict: ok. Filed 0 HIGH blockers (none).…
    ┃
    ┃  3d5b773 inbox: reviewer reports tui-t2a cycle 3
    ┃  ac1bc3c docs: T2a's last cycle — the bound's tripwire bites, and the slice is done
    ┃  CLEAN
    ┃

    ┃
    ┃  $ cd /Users/thiago/Code/Pessoal/maitre && herdr agent prompt "w1H:p1" "Reviewer: T2a cycle 3 done. Verdict: ok. Filed 0 HIGH blockers (none). F-109:
    ┃   the mutation — detailLineCount() counted at the full-screen width again (m.width - 4 instead of m.detailPaneWidth()), the pre-fix defect — t…
    ┃
    ┃  {"id":"cli:agent:prompt","result":{"agent":{"agent":"opencode","agent_session":{"agent":"opencode","kind":"id","source":"herdr:opencode","value":
    ┃  "ses_e737d6341ffeTdoEYTgykU25M8"},"agent_status":"idle","completion_seq":157,"cwd":"/Users/thiago/Code/Pessoal/maitre","focused":false,
    ┃  "foreground_cwd":"/Users/thiago/Code/Pessoal/maitre","pane_id":"w1H:p1","revision":15,"screen_detection_skipped":true,"state_change_seq":157,
    ┃  "tab_id":"w1H:t1","terminal_id":"term_65f0cd2aa89d26","terminal_title":"OC | Project work proposal","terminal_title_stripped":"OC | Project work
    ┃  proposal","workspace_id":"w1H"},"type":"agent_prompted"}}
    ┃

       verdict: ok — no HIGH blockers filed; none needed.

       What I verified, by command

       ┌──────────────────────────────────────┬───────────────────────────────────────────────────────────────────────────────────────────────────────────┐
       │ Check                                │ Result                                                                                                    │
       ├──────────────────────────────────────┼───────────────────────────────────────────────────────────────────────────────────────────────────────────┤
       │ F-109 — the tripwire bites           │ My mutation (detailLineCount() back on the full-screen width — the pre-fix defect) turns the whole        │
       │                                      │ internal/ui suite red, and it fails for the bound's width: the pane has 29 rows at 107 columns; the full  │
       │                                      │ screen has 27 at 196 · the bound is 27 rows; the pane has 29 rows · the scroll stopped at 26, want the    │
       │                                      │ pane's last row 28. One test, nothing else                                                                │
       ├──────────────────────────────────────┼───────────────────────────────────────────────────────────────────────────────────────────────────────────┤
       │ The test does not ask the code under │ The pane's width is computed in the test (m.width - m.listColumn() - 3) and the rows are counted at that  │
       │ test                                 │ width — the old test compared two values that both went through detailPaneWidth()                         │
       ├──────────────────────────────────────┼───────────────────────────────────────────────────────────────────────────────────────────────────────────┤
       │ The fixture                          │ longDetailModel builds a plan through the core operations; a long outcome wraps into 29 rows at 107       │
       │                                      │ columns, 27 at 196, and the test fails out loud if the two widths ever agree again                        │
       ├──────────────────────────────────────┼───────────────────────────────────────────────────────────────────────────────────────────────────────────┤
       │ The document reads whole             │ the table's second row now says "anywhere else, while the list is on screen"; the prose names the carved- │
       │                                      │ out case. Every row true by command: split layout + header → the list scrolls; full-screen detail +       │
       │                                      │ header → enter and enter,wheel:down@(150,0) are byte-identical (nothing); full-screen detail + detail →   │
       │                                      │ the frame changes                                                                                         │
       ├──────────────────────────────────────┼───────────────────────────────────────────────────────────────────────────────────────────────────────────┤
       │ The scope                            │ no click handler; the wheel's fallback untouched since cycle 2; maitre-cli unchanged; seven tabs;         │
       │                                      │ internal/core, internal/cli, docs/03 untouched; DomainVersion 7; budget 0/0/28                            │
       ├──────────────────────────────────────┼───────────────────────────────────────────────────────────────────────────────────────────────────────────┤
       │ The goldens                          │ unchanged                                                                                                 │
       ├──────────────────────────────────────┼───────────────────────────────────────────────────────────────────────────────────────────────────────────┤
       │ The gate                             │ make check and make check-race green                                                                      │
       └──────────────────────────────────────┴───────────────────────────────────────────────────────────────────────────────────────────────────────────┘

       The ledger

       F-109 closed (mutation-proven), the cycle-3 row added, and the two notes the fix cycle left went to the notes table — I do not think either is a
       finding: detailScrollRows writes the offset during a render (moved from renderDetailFull, not new), and the wheel over the list moves the selection
       so the pane shows another recipe from its top (the console's model).

       Nothing is open: F-106, F-107, F-108 and F-109 are all closed, and the open table holds no T2a finding. Ledger ac1bc3c, receipt 3d5b773, nudge sent.
       The tree is exactly as the runner left it.

       Build · MiniMax-M3.1-Flash · 10m 47s · 117.5 tok/s

    ┃
    ┃
    ┃
    ┃  Build · MiniMax-M3.1-Flash MiniMax Token Plan (minimax.io) · medium
    ╹▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀
    ~/Code/Pessoal/maitre:master                                                                                              203.0K (20%)  ctrl+p commands
```
