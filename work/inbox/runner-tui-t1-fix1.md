Runner: T1 fix 1 done. Commits: 5 (5a79638, 67b455d, b4a3052, 55bdc81, e218ec3). Fixed: F-99, F-100, F-101, F-102.

The name is sanitised: a plan whose directory basename carries a newline — maitre-cli init writes "we\nird" — renders the header " ⋈ maitre · we ird" as ONE row (sed -n l shows a single line, no continuation), and kitchen/project.json still holds "we\nird". The file on disk is untouched; the sanitising is a rendering concern. The window title gets the same name, so it is sanitised too.

The name has one answer: a plan whose project.json names it empty renders " ⋈ maitre · f100x", and the same plan with project.json removed renders " ⋈ maitre · f100x" — the same answer, and a name, not a path. Before the fix they answered "/tmp/f100x/kitchen" and "f100x". The answer now lives in the domain: core.Store.ProjectName() is the one answer, and core.LoadProject answers the same for an empty name as for a missing file. The console asks the domain (model.projectName() -> sanitize(store.ProjectName())) and no longer computes a name of its own — grep GetProject internal/ui/ finds no call site.

The root row shows: the plan directory — "/tmp/f100x/kitchen", and humanised "~/maitre-t1-evidence/kitchen" when the plan sits under the home directory. The themes are still listed (Kitchen, Catppuccin Latte, Frappe, Macchiato, Mocha, Dracula, Alucard).

DomainVersion: 7, and TestDomainShapeIsVersioned passes. The change is in the reader only — core.LoadProject's empty-name case and the new core.Store.ProjectName(); the writer core.Init is untouched and the plan format on disk does not move.

The goldens: unchanged. The fixture's project is "demo" and the header's words did not move; make goldens-check is green with no diff.

The gate: make check green, make check-race green, tree clean. The budget is 0 / 0 / 28 against 0 / 0 / 30.

Notes:
- The new core.Store.ProjectName is inside the declared "Store.*" family in core.NonViews, so it needs no row in docs/03 § 8 and the specification does not move.
- The CLI did not change: maitre-cli plan prints the name the domain now answers ("f100b"), which is the whole of the change.
- The Settings row's test change landed in 67b455d with the other UI test changes; the row's own change is b4a3052.
- T1's brief was corrected in 55bdc81 where it was wrong — the root's example and the name's source — and each correction says so in the brief.
- docs/05-human-interface.md gained one paragraph (F-102): the window title carries the project, and it is not the screen's row 0. The screen sketch is unchanged.
