Reviewer: T1 cycle 3 done. Verdict: ok. Filed 0 blockers (none).

F-103 verified: the rule now lives in one private function - core.Store.projectName - and both core.Store.LoadProject and core.Store.ProjectName ask it; defaultProjectName has one caller. Mutation: I took the rule out of that one copy (projectName returning doc.Name raw) and ran the FULL suite: RED in both packages, in three tests as the runner claimed - internal/core TestProjectNameHasOneAnswer and TestLoadProjectAnswersTheSameName, internal/ui TestTheNameHasOneAnswer. Reverted.

The two surfaces: maitre-cli plan and maitre --render agree in all three states - the plan names itself (demo2 / ⋈ maitre · demo2), the file names it empty (demo2 / ⋈ maitre · demo2), the file is missing (demo2 / ⋈ maitre · demo2). The CLI's behaviour did change, as the finding required: it used to print "" where the console printed the project's directory name. The agreement is not a fixture coincidence - the mutation above turns both packages red, and the CLI reads the name through LoadProject while the console reads it through ProjectName.

F-104: humanPath sanitises the path, and it is the only path the console renders (one call site, in the Settings tab's root row), so it is by construction and not by a call site's care. Proof by command: a plan under a directory named "we\nird" shows "root  .../we ird/kitchen" on ONE row. Mutation - the sanitize call taken out of humanPath - fails TestTheRootRowIsOneRow twice over. Reverted.

F-105: the row says "11 commits" and names 11 ids - the count and the list agree - and every one of the eleven is a real commit of the slice in the log. The commits that wrote the row (8c25ed8, 53ba171, e218ec3, 4e67b31) are not among them, as the shape requires. The protocol section carries the rule in one line: "The count is the number of ids the row names, and the commits that wrote the row are not among them - a commit cannot carry its own hash."

The rule in the spec: NOT written down anywhere in the spec. It is written down in the CODE - the doc comment on core.Store.projectName ("the plan's own name, and the project's directory name when the plan names nothing: a missing file, an unreadable file and an empty name are the same case"). docs/03 does not say it: § 6 names project.json as "project meta · charter · settings" and nothing about the name's fallback; docs/05 does not either. I wrote nothing into docs/03 - that is the chef's call, and the orchestrator has it.

The goldens: unchanged - make goldens leaves git diff empty. The scope held: docs/03 untouched in this cycle, docs/05's paragraph stays, internal/cli and cmd/ untouched, budget 0 / 0 / 28 against 0 / 0 / 30, seven tabs with the same names, core.DomainVersion still 7, the glyph only in internal/brand, docs/audit.md untouched by the runner.

The gate: make check green, make check-race green. The tree was left exactly as the runner left it - every mutation reverted, git status and git diff empty apart from this review's ledger commit.

LOW notes: TestTheRootRowIsOneRow splits the rows on \n and then checks each line for a \n - an assertion that cannot fail; the test still bites through its other two assertions, so the rule is held. Still open from earlier rounds: docs/05's window-title paragraph is held by nothing (deleting it leaves the doc check green), and the plan's name rule is nowhere in the spec.
