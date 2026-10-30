Reviewer: T1 cycle 2 done. Verdict: changes-needed. Filed 1 blocker (F-103).

F-99..F-102 verified: each by command and each by a mutation, all reverted.
- F-99 (the name is sanitised): a plan under a directory named "we\nird" renders the header on ONE row ("⋈ maitre · we ird") while project.json still holds "we\nird". Mutation - the sanitize call taken out of internal/ui/identity.go - fails TestTheNameIsSanitised three times over (the name still carries a control character; the header does not carry the sanitised name; the header is 41 rows, want 40).
- F-100 (one answer): the empty-name plan and the plan with no project.json both show "⋈ maitre · demo", and "maitre-cli plan" prints the same name - a name, never a path. The console asks: grep gives one call site, m.store.ProjectName() in internal/ui/identity.go. The writer is untouched (core.Init still writes the directory basename) and core.DomainVersion is still 7.
- F-101 (the root row): the Settings tab shows "root  ~/Code/Pessoal/maitre/kitchen" - the plan directory, humanised - and the themes are still listed. Mutation - the row put back on the parent - fails TestTheSettingsTabShowsThePlanRoot twice over.
- F-102 (docs/05): the paragraph is there and says what the window title carries and that it is not the screen's row 0; T1's brief was corrected in the same commit and says so. Nothing holds the paragraph: deleting it leaves the doc check green - a LOW note.

F-103 (medium, filed): the domain's one-answer rule for the plan's name is written TWICE and only one copy is pinned. core.Store.LoadProject answers an empty name with the project's directory name, and core.Store.ProjectName answers the same question again for itself. No test pins LoadProject - the only test that reads it reads a plan that names itself. Mutation: taking the empty-name branch out of LoadProject leaves THE WHOLE SUITE GREEN (every package), and then the two shipped surfaces disagree - "maitre-cli plan" prints "Name": "" while the console shows "demo2". Remedy: pin LoadProject's answer in internal/core/store_test.go for the same two cases ProjectName is pinned for, or drop ProjectName's redundant fallback so the regression shows through it.

The path row: a note, not a blocker. The path is not the plan's content, the unsanitised path predates the slice (the row showed the absolute path before T1), and the consequence is one extra visual row in a settings list - the header and the window title both survive the same directory. It is in the ledger's notes table; sanitize(humanPath(...)) would close it.

The state row: a note, not a blocker - cycle 1's LOW stands, the fix cycle repeated it (the row says 12 commits and names nine; the three it does not name are the three commits that wrote the row). The shape that cannot drift is a count that counts what the list names, or a row that says the count excludes the commits that wrote it.

The name in docs/03: it does NOT belong in § 8. That table is the derived views - "computed, never stored" - and the plan's name is a STORED fact the domain reads, like LoadProject and ListRecipes, which the table does not carry either; the coverage wire agrees (Store.ProjectName is a declared non-view under the Store.* family, and the table and the wire check each other both ways). No MEDIUM filed on that question. A LOW note instead: the rule itself ("a plan that names nothing answers with the project's directory name") is nowhere in docs/03.

The goldens: unchanged - make goldens leaves git diff empty. The scope held: internal/core and internal/cli stdlib-only, budget 0 / 0 / 28 against 0 / 0 / 30, seven tabs with the same names, the glyph only in internal/brand, docs/audit.md untouched by the runner.

The gate: make check green, make check-race green. The tree was left exactly as the runner left it - every mutation reverted, git status and git diff empty apart from this review's ledger commit.

LOW notes: the path row is not sanitised; docs/05's window-title paragraph is held by nothing; the plan's name rule is nowhere in the spec; the state row's count and list disagree.
