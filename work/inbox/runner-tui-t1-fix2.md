Runner: T1 fix 2 done. Commits: 3 (8a6866b, cf8c05e, 4e67b31). Fixed: F-103, F-104, F-105.

The rule lives in: one private function — core.Store.projectName(doc projectFile) — and everything asks it. I chose the private helper over letting LoadProject answer directly because LoadProject and ProjectName are both asked by a different surface (the CLI and the console), and a rule either of them carries can be deleted without the other noticing: that was F-103. One copy, three askers, and the pin sits on the copy. core.Store.ProjectName and core.Store.LoadProject both call it; grep finds one branch that answers "the plan names nothing".

The mutation: I removed the empty-name branch from core.Store.projectName and ran the suite. It went RED: internal/core failed TestProjectNameHasOneAnswer and TestLoadProjectAnswersTheSameName, and internal/ui failed TestTheNameHasOneAnswer ("the name is \"\", want the project's directory name \"demo\""). The mutation was reverted and the tree is exactly as it was — git status is clean and the diff holds only the three commits' work.

The two surfaces: in all three states the CLI and the console answer the same. The plan names itself — cli: f2, console: " ⋈ maitre · f2". The file names it empty — cli: f2, console: " ⋈ maitre · f2". The file is missing — cli: f2, console: " ⋈ maitre · f2". Before the fix the CLI printed "" for the empty-name state while the console printed the directory name.

The path is sanitised: a plan under a directory named "we\nird" shows the root row as ONE row — "  root        /tmp/f2x/we ird/kitchen" (sed -n l prints a single line) — and the header is one row too. The sanitising lives inside humanPath, so every path the console draws is sanitised by construction and not by a call site's care.

The state row: the count counts what the list names. The row says 11 commits and names 11 ids; the commits that wrote the row (8c25ed8, 53ba171, e218ec3, 4e67b31) are not among them, because a commit cannot carry its own hash. The protocol section in work/tui/README.md now says so in one line, so the next writer knows the shape and cannot read the row two ways.

DomainVersion: 7, and TestDomainShapeIsVersioned passes. The change is in the reader only — core.Init still writes what it wrote, and the plan format on disk does not move.

The goldens: unchanged. The fixture's project is "demo", the header's words did not move, and make goldens-check is green with no diff.

The gate: make check green, make check-race green, tree clean. The budget is 0 / 0 / 28 against 0 / 0 / 30.

Notes:
- The CLI's behaviour changed, as the brief allowed: maitre-cli plan now prints the domain's answer for a plan that names nothing ("f2") instead of "". That is the point of F-103 — the two surfaces agree.
- The error strings did not change: readProjectFile keeps "read project: %w" and "parse project: %w", and LoadProject returns them as before. A missing file is now found with errors.Is(err, os.ErrNotExist), because the error is wrapped.
- LoadProject carries the name even when the file cannot be read, so a caller that ignores the error still gets the plan's name.
- docs/03 was not touched: the name is not a view, and the reviewer is right about that. Where the RULE belongs is the chef's call — my reading is that it belongs with the plan's anatomy (the format section of docs/03), since it says where the name comes from and what a plan that names nothing answers, and § 8 is the derived views. I wrote it nowhere; it is in the code, where it is enforced.
- The state row's ids are the slice's code commits; the four commits that wrote the row are deliberately not named, and the count agrees with the list.
