Runner: slice H fix 1 done. Commits: 2 (plus this receipt). Fixed: F-98 and the seven notes. The tool versions live in: the `Makefile` (`STATICCHECK_VERSION`, `GOVULNCHECK_VERSION`) — the gate holds the pin, and `docs/testing.md`'s copy is the checked copy, the same shape as the dependency budget; a pin that lives in a doc can be edited to fit the tool on the `PATH`, and the doc cannot drift from the gate. Notes: `docs/audit.md` is the reviewer's — not touched.

Per-AC evidence:

- **F-98** — `make test-race` honours the filter and fails on a zero-match run; `make test-race PKG=./internal/brand RUN='-run NoSuchTest'` exits 2. One place decides how a test run is invoked: `test` and `test-race` differ in the flags they add (`TEST_FLAGS`) and in nothing else — the filter and the package always reach `scripts/go-test.sh`, so the zero-match check always applies. `make test-race PKG=./internal/brand RUN='-run TestTheDrawingsMatchTheSource'` exits 0 and runs that test (twice, under `-count=2`).
- **the tool's version is never checked** — `make staticcheck` and `make vuln` now check the version of the binary they run against the pin. A stale `staticcheck` on the `PATH` fails: *"staticcheck is version '0.8.1'; the pin is v0.8.2. / Run 'make tools' to install the pinned version."* (proved with a fake `staticcheck` earlier on the `PATH`).
- **the tool versions are written twice and only one is used** — the pins live in the `Makefile`; `docs/testing.md` names them and `make versions-check` (in `check` and `check-race`) fails when the doc and the pin disagree: *"docs/testing.md names staticcheck 0.8.1; the pin is v0.8.2. / Change docs/testing.md to name v0.8.2."*
- **`dep-budget.sh` ignores `GO`** — it honours `GO` now, like `require-tool.sh` and `go-test.sh`. `make budget` passes it.
- **`tidy-check` leaves the tree dirty if tidy itself fails** — `go.mod` and `go.sum` are put back when tidy fails. Proved with a broken `go.mod`: the gate failed and the file was exactly as it was.
- **`install-hooks` overwrites silently** — it prints the hook it replaces before replacing it.
- **`make -j` is unsafe** — `.NOTPARALLEL:` in the `Makefile`, with the reason: the gate's targets touch the goldens, `go.mod` and `go.sum`. `make -j4 check` is green and serial.
- **the ledger's discipline table** — `docs/audit.md` is the reviewer's. Not touched.

The gate: `make check` 1m14s green, `make check-race` 3m12s green, `git status --short` clean. Counts unchanged: 0 / 0 / 28. `core.DomainVersion` is still 7.
