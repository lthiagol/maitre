Runner: T9 done. Commits: 6 — 38d1d48 (AC-01), 2aba9f4 (AC-02), 089ef6b (AC-03), 1ed7f3f (AC-04), 2d65fa0 (AC-05), 17cdefe (AC-06, the row's own).

The package: `internal/settings` — stdlib only, 0 external modules. It owns the record's schema (`version`, `view`), the read, the atomic write and the version check. Nothing else knows the file's shape.

The file: `kitchen/settings.json` — `{"version":1,"view":{"theme":"Kitchen","order":"triage","showServed":true}}`, in the plan's folder beside the plan's records. `internal/core` did not change (the diff is empty) and does not read the file: a test reads the core's sources and fails if one mentions `settings.json` or `internal/settings`. Another test seeds a plan, writes the record into the plan's folder, and reads the plan back — the plan's readers ignore it. `core.DomainVersion` stays **7**; `planFormatFacts()` did not move. The repo's own `kitchen/` is not stamped and holds no settings file.

The mechanics: a missing file yields the defaults, and it is not an error. The write is atomic — the bytes go to a temporary file in the same folder and are renamed over the old one, so a half-written file cannot lose the chef's choices; a write into a folder that is not there fails and leaves the old file whole. An unknown version is refused **by name** ("the settings file has version 99, which is not one this build has (version 1)") — the tool does not guess. Unknown keys are ignored, so a file with keys this build does not know still reads. The writer stamps its own version.

The verb: `maitre-cli settings show` — the record, as one JSON document, the way every list the CLI prints is printed. It is **read-only**: the console is the only writer, because the console owns the vocabulary. The count moved with it — **46 commands: 40 verbs in 8 groups, plus 6 workspace commands** — and `AGENTS.md` and `docs/04-operations.md` move with the code; `TestTheDocsCommandCountMatchesTheCode` holds them.

The console: the Settings tab's VIEW preferences — the theme, the order and show-served — are read from the record when the console starts, and written when the chef changes one of them (`enter` applies the theme, `1`/`2`/`3` set the order, `c` toggles the served recipes). The chef's choices survive a restart: a test picks a theme, an order and the toggle through the keys, then a second console over the same plan starts with them. The record is not the plan's record: the journal's event count does not move when the chef changes the theme.

The refusal: a theme or an order the console does not know is refused, **by name**, and the console **says so** on the status line — "the theme Neon is not one this build has" — and keeps its own default. It does not fall back silently, and it never writes a value it does not know: the refusal happens before the write, the way the domain's writer refuses before it writes. Each value is applied on its own, so an unknown theme does not cost the chef the order.

The budget: `internal/settings 0 0` in `scripts/dep-budget.sh`, and the row "stdlib only / 0 / 0" for `internal/settings` in `AGENTS.md`'s dependency table; `AGENTS.md`'s Layout section names the package. `work/tui/README.md`'s T9 row names the five commits above. DomainVersion: **7**.

The gate: `make check` — pass; `make check-race` — pass.

Notes for the reviewer:

1. **The exception wire grew a second group.** `settings show` never opens the plan's records, so `planlessGroups` asks the code and finds the `settings` group beside `skill`. The documents must name both, and the phrase is per group — the wire's `operationExceptionPhrase` now knows two sentences ("…over the tool's own embedded files…" and "…over the tool's own settings file…"), and the test was renamed to `TestEveryVerbOutsideTheReadOnlyGroupsMapsToAnOperation`. `AGENTS.md` and `docs/04-operations.md` carry both.

2. **Forward compatibility, precisely.** The brief asks for two rules that touch: unknown keys are ignored, and an unknown version is refused by name. The version fence wins — a file at version 2 is refused before its keys are read. The unknown-key rule therefore covers keys added **within** version 1 (the test writes a file with a `crew` section and a `density` key and reads it). `plan-automation.md` § 15 says T10 moves the schema's version to 2, so a T9 build will refuse a T10 file by version — that is the fence the design asked for, and T10 may want to say so in the brief.

3. **The record's defaults name the console's defaults.** The record says `Kitchen` and `triage`, because the CLI prints the record and the CLI does not know the console's vocabulary — so the two names are written twice. A test holds the tie (`TestTheRecordsDefaultsAreTheConsolesDefaults`). A rename of the console's default palette must come with the other name; the failure would be loud, not silent.

4. **`kitchen/settings.json` is untracked, and the console writes it.** Running the console in this repo — even `--keys "c"` — writes the file into `kitchen/`. The brief says the repo's own plan must not get one; it does not, until a console runs here. The orchestrator may want it in `.gitignore` (it is not there today). I did not touch `.gitignore`: no acceptance criterion named it.

5. **The console's write is a disk write in the key path.** `applyTheme`, the order keys and `c` write the file synchronously. The plan's fingerprint hashes every `.json` in the folder, so a write is seen by the poll and the console reloads the plan once — harmless, and the design says the poll already watches the file.

6. **`work/tui/README.md` has a second row labelled T9** (the Sous-Chef tab, which is T10). I left it alone — the brief gave me the T9 state row only.
