Runner: T1 done. Commits: 7 (94a53ff, 8fc7438, f0b2aa5, 82289c6, 8c25ed8, 0bbcb35, 53ba171). Per-AC evidence: AC-01..AC-06 stamped. Pre-notify gate: pass (make check green, make check-race green, tree clean). The window title now reads: "maitre · <project>" — e.g. "maitre · demo" for the seeded plan, "maitre · maitre" for the repo's own kitchen. The goldens: unchanged (the header's words did not move; only its style did, and the Settings tab is not in the golden).

Notes:
- The name has one source: `model.projectName()` in `internal/ui/identity.go`. The header and the window title both read it; it answers the plan's project.json name, the `--project` override, then the plan's root. `grep -rn GetProject internal/ui/` returns one call site.
- `--project` is display only and reaches all three constructors through `ui.WithProject(name)`; `--render` and `--keys` honour it (they build the same model). `project.json` and `activity.json` are byte-for-byte unchanged after a run with the flag.
- The brief said the Settings tab had no path row. It had one, labelled "root", showing the plan directory (<project>/kitchen). It now shows the project's directory — the one the chef can open — with the home directory written as "~". The themes are still listed.
- The brand's wire holds the compact glyph to internal/brand: it may ride the rendered output but not be pasted into code. The header's test asks for `brand.Compact`.
- docs/05-human-interface.md needed no fix: its screen sketch already said "row 0  maitre - <project>", which was aspirational until now and is true now.
- The orchestrator's commits (2d068fa, 0e3b8cc) landed in the middle of this slice; no file of theirs was touched by these commits.
