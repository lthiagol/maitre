Reviewer: slice G cycle 3 done. Verdict: ok. F-97 closed: yes. Filed 0 new blockers (none).

Mutation checks (every one reverted; the tree is clean apart from this receipt and the ledger commit):

F-97 — the exemption is the declaration, not the file.
- A copy of the full drawing pasted into `internal/brand/mark.go` BESIDE the declaration -> the wire now FAILS ("carries a line of the full drawing — the drawing is declared once, in internal/brand"), and the whole suite fails with it. This was the one that slipped through in cycle 2.
- A copy in each of the seven categories the wire reads -> all FAIL the wire: `internal/brand/zz_probe.go`, `internal/ui/testdata/zz_probe.golden`, `docs/zz_probe.md`, `kitchen/zz_probe.json`, `cmd/maitre/zz_probe.go`, `images/logo/zz_probe.svg`, `internal/ui/testdata/zz_probe.txt`.
- The legitimate copies still PASS: the `.txt` sources (`images/logo/ascii/bow-tie-full.txt` and `-half.txt`), the declaration itself, and — for the glyph — the rendered output (the two goldens), the wallpaper SVG, the brand sheet and the ledger's prose.
- A copy in a comment inside `mark.go` -> caught. A copy inside the wire's own test file -> caught. A copy in a raw string literal (backticks), a re-indented copy, and an escaped-newline copy -> all caught. The half drawing inside the `Full` literal -> caught.
- A duplicate element INSIDE the `Full` literal passes the scope wire (the literal is masked) but is caught by `TestTheDrawingsMatchTheSource` — a second net under the first.
- The declaration renamed `Full` -> `FullForm` fails closed: the wire cannot find the declaration and says so.

The glyph's scope is a category.
- MUST FAIL, and all fail: `internal/ui/zz_probe.go`, `scripts/zz_probe.py`, `scripts/zz_probe.sh`, `Makefile`.
- MUST PASS, and all pass: a golden, an SVG under `images/`, a doc, a brief. The real copies still pass: goldens 1 each, wallpaper 1, brand sheet 1, declaration 1.
- The four exclusions are still excluded: `work/briefs`, `work/inbox`, `work/transcripts`, `skills/`, `bin/` — a copy in each is not read, for the reason the comment gives (records or build output).

The slip: clean. `git hash-object scripts/console-shot.py` == `git rev-parse HEAD:scripts/console-shot.py` (`c3b62ec54cba8f3f12dd5560956e1f0a06c9b106`); `git status --short` empty; `git diff HEAD` empty; `git diff --cached` empty. The runner reported it instead of letting the next reader find it.

Scope: `core.DomainVersion` is 7 and did not move (`internal/core/lifecycle.go:225`). The command count is 42 (36 verbs in 6 groups + 6 workspace) and `TestTheDocsCommandCountMatchesTheCode` passes; the gate table, the plan format and the derived-views table are green. The budgets hold: core 0, cli 0, ui 28, brand 0. The card's copy stays in the card — `EVIDENCE REQUIRED`, `agentic work` and `taking orders` are in `images/logo/9-plan-behind.svg` and in no `.go` file. The runner's commit `b5dda82` touched `internal/brand/mark_test.go` only.

How it looks: `maitre --render --plain` — the header is one line, `⋈ maitre · <project>`. `maitre --render --all --plain` — the Sous-Chef tab carries the full form above the section, `Not built yet.` is still there, the vital signs read `7 / 7 / they agree`, and the last event and its row are whole. `maitre-cli` with no arguments — the full form above the title line. At 80x24 (pty capture) the tab is **byte-for-byte what it was before the mark**: `LAST EVENT` and the event row on the screen and NO mark; the header's compact glyph is the only difference from the baseline. At 80x40 the mark is there and the content is whole. The height rule is sound at every height from 10 to 45 — the mark never costs the content, the content is whole from h=20 (the baseline's own threshold), and the mark appears from h=26. Both width boundaries are pinned: mutating the full-form threshold or the half-form threshold fails the suite.

LOW notes (not filed):
- The scope wire searches by drawing LINE, so a copy split across a string concatenation (`"a" + "b"`) is not caught — the whole suite stays green. A deliberate evasion of a line-based search; the drawings are long and the same rule catches every whole-line form (literal, raw string, comment, re-indented, escaped).
- The glyph's exemption is the declaring FILE, not the declaration: a second glyph const pasted beside `Compact` in `mark.go` passes. The drawings are per-declaration (masked spans) but the glyph is per-file (`glyphDrift` skips `decl.file`). One character, and the file is the brand package's own — a note, not a blocker.
- The declaration reformatted onto one line is not gofmt-clean, so `gofmt -l` catches it; a gofmt-clean one-line form is still caught when a copy sits beside it.
