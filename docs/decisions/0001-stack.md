# 0001. Implementation stack

- **Status:** 🔄 proposed
- **Date:** 2026-10-25
- **Deciders:** TBD

> **Revised** after the premises changed. The first draft recommended **Rust** on reusing
> ~30k LOC of v1 and on agent-CLI startup time. Both are void: **no code is carried over**
> (premise 5) and **the CLI is secondary** (premise 3). The premise that does decide it:
> **the human interface is the product, so the UI framework's component library decides.**

## The question

> Which stack gives the human interface the fastest path to a complete, good-looking,
> keyboard-driven plan manager — given everything is written from scratch?

Needs: a table of hundreds of recipes, a tree, markdown rendering (CPD artifacts are
prose), diff rendering (findings anchor to code), text input, tabs/sidebar, modals,
keybindings, live updates.

## The candidates

| | **Charm** (Go) | **Ratatui** (Rust) | **OpenTUI** (TS/Bun/Zig) | **MyGo** (Go) |
|---|---|---|---|---|
| Latest | bubbletea **v2.1.0** | v0.30.2 | **v0.6.8** | v0.6.4 |
| Version line | v2 (since 2026-02-24) | 0.x | 0.x, weekly breaks | 0.x, 67 releases/30d |
| First-party widgets | ✅ **bubbles v2** (list, table, tree, viewport, textinput, textarea, spinner, progress, help, key) | ❌ none — primitives only | ⚠️ partial (table, diff, markdown, code, input, select, scrollbox) | ✅ rich, but **GUI not TUI** |
| Styling | ✅ **lipgloss v2** | manual | ❌ no theme system | ✅ own `ui.Theme` |
| Markdown | ✅ **glamour v2** | 3rd-party | ✅ 2,401-LOC native | ✅ own CommonMark+GFM parser |
| Diff | ✅ `charmbracelet/git` | 3rd-party | ✅ 1,360-LOC native | ✅ native diff viewer |
| Missing for us | — | **the entire component kit** | **tree, modal, tabs container, list, virtualization, theme, focus traversal** | **kanban, gantt, editable grid, docking** |
| SSH serving | ✅ **wish v2** | manual | ✅ `@opentui/ssh` | n/a |
| Known perf risk | — | — | ⚠️ **#1339: idle spinner = 23% of a core** (open 3mo) | ✅ published benchmarks |
| Runtime | none | none | Bun ≥1.3.14 *or* Node ≥26.4 + `--experimental-ffi` | none |
| Binary | ~5–10 MB | ~2–5 MB | ~50–90 MB (Bun) | 8.8 MB |
| Ecosystem | 57 repos, all pushed this week; `crush` (28.8k★) is Charm's own agent TUI | mature | 13.6k★, **78–94% one maintainer** | 2.9k★, **29 days old, 0 importers** |

### MyGo is not a TUI framework

There is **no `tui` directory** in MyGo's entire git history. It is a **desktop
application framework** — native windows via Metal/D3D/OpenGL, plus Go→WASM→WebGL2 for
the browser. Its terminal code is `plugins/terminal`, a **Ghostty VT emulator** that hosts
other programs, including a bubbletea app.

So MyGo is not an alternative to Charm; it is a potential **later host** for it. It is also
29 days old with zero importers, no stability policy, and 67 releases in 30 days — **a v3+
option, not a v2.0 dependency.**

## Decision

**Go + Charm.**

```
bubbletea v2        the application loop (Model / Update / View)
lipgloss  v2        styling, layout, borders, adaptive color
bubbles  v2         components: table, tree, list, viewport, textinput, spinner, help, key
glamour  v2         markdown rendering (specs, notes, findings)
charmbracelet/git   diff rendering
wish      v2        serve the same TUI over SSH (deferred)
```

### Rationale

1. **`bubbles` decides it.** Ratatui is a rendering substrate with no first-party
   component library; from scratch you'd assemble third-party crates and hand-roll layout
   math. OpenTUI's equivalent is missing exactly the widgets we need. Charm ships all of
   it, first-party, v2, with a consistent design language.
2. **`glamour` is a direct hit.** CPD's artifact is prose — intent, problem, scope, ACs,
   audit trail. A stylesheet-driven CommonMark renderer is table stakes.
3. **Proven at our scale.** `crush`, Charm's own agent TUI, is 28.8k★ on bubbletea v2.
4. **Go's iteration speed suits a large from-scratch UI.** Measured on v1: Rust's
   incremental `cargo check` is 1.3–1.5 s, but a full rebuild of a large UI package tree
   is ~17 s. Go rebuilds in seconds.
5. **The agent-CLI startup argument is void.** The CLI is secondary.
6. **Go keeps the MyGo door open.** A native/web front end later is a same-language
   option — hosting the bubbletea TUI via `plugins/terminal`, or a second shell over the
   same domain package. This requires **the domain layer to have zero UI imports**
   (principle P12).

## Consequences

### Positive
- Fastest path to a complete, polished human interface.
- All the boring-but-essential components exist and are first-party.
- Markdown + diff rendering solved on day one.
- Single static binary; trivial distribution; one brew formula.
- The domain package is UI-free, so the UI is replaceable.

### Negative
- **Charm's v2 line is ~8 months old** (v2.0.0 was 2026-02-24). The stack is
  consistently v2, but churn is possible. Mitigate: pin exact versions; keep the UI
  behind a thin adapter.
- **Elm-architecture discipline required.** A big app's `Model` grows fast. Mitigate:
  organize by feature, not widget; keep view logic out of `Update`.
- **We give up OpenTUI's visual ceiling** — images, audio, 3D, Yoga flexbox, the best
  grapheme handling. Acceptable for a plan manager.
- Go's generics are adequate but not elegant for the projection machinery.

### Measured (2026-10-25, slice 2)

| Package | External deps | Policy |
|---|---:|---|
| `internal/core` | **0** | stdlib only |
| `internal/cli` | **0** | stdlib only |
| `internal/ui` | 22 | ≤ 30 |

The core uses only the standard library — JSON, atomic writes and fsync are all stdlib.
The UI carries the Charm v2 tree.

### Module paths moved in v2

Charm v2 publishes under `charm.land/<name>/v2`, not `github.com/charmbracelet/<name>/v2`:

```
charm.land/bubbletea/v2   charm.land/lipgloss/v2
charm.land/bubbles/v2     charm.land/glamour/v2
```

`go get github.com/charmbracelet/bubbletea/v2` fails with *"module declares its path as:
charm.land/bubbletea/v2"*. Use the `charm.land` paths.

### The web door stays open

Because the domain is UI-free, a **local web UI** served by the same binary is a v3+
option (`--ui web` vs `--ui tui`) with no core rewrite. A web shell gets unlimited
rendering, shareable URLs, real diff/markdown libraries and accessibility for free — at
the cost of a server and a browser. Revisit once the TUI proves the model.

## Resolved

- **Dependency budget.** Core and CLI: **0 external deps, stdlib only**. UI: **≤ 30**
  external (currently 22). A dependency needs a written reason; the core may never gain
  one for convenience.
- **Single binary or two?** **Two binaries, one brew formula.** `maitre-cli` for agents,
  `maitre` for humans. An agent cannot launch the TUI by accident — `maitre` refuses
  without a terminal, which is the structural guard. Two *formulas* was v1's bug; two
  *binaries* in one formula is fine.
- **Naming** → [ADR-0004](0004-naming.md).

## Open questions

- The tool's module path is `github.com/lthiagol/kitchen-brigade`. Confirm the hosting
  account before the first release.
- `go.mod` requires Go 1.26 (raised by the Charm v2 modules). Confirm that is an
  acceptable floor for adopters.

