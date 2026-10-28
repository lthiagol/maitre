You are the **reviewer** — the same lane that reviewed slice E. The orchestrator's pane
id is `w1Z:p1`; use it verbatim in the notify command below.

## Cycle 2 — the fixes to slice E

The runner fixed the five findings you filed. **Cycle 2 is blockers only**: file
**HIGH or MEDIUM** blockers. No nits, no cosmetics, no style.

**Do not re-file:** F-56, F-59, F-60, F-61, P-8 (the chef's), or your own LOW notes
from cycle 1. (`README.md`'s decisions table and `work/README.md`'s old role names
were the orchestrator's and are fixed.)

### 1. Are the five findings actually closed?

Verify each with a **command**, not by reading the diff and believing it.

**F-81 — the console's retired name is carried, and the proof reproduces.**

- The console's entry carries *maître d'* in its `_Avoid_` line, like every other
  retired name.
- **The proof must reproduce.** The runner offers
  ```bash
  grep -n -i -E "line cook|taster|expo|proof|special|prep note|maître d'" GLOSSARY.md
  ```
  and says *seven retired names, seven hits, every one of them the line that warns,
  and 0 hits outside an `_Avoid_` line*. **Run it.** The count is the claim: **every
  hit must be on an `_Avoid_` line, and the number of retired names must be the
  number of hits** — a retired name named twice is F-83 again, and a retired name
  named nowhere is F-81 again.

**F-82 — the console is not one of the crew.**

- "### The crew" holds **five** entries, and the console is under a separate section
  (the strings), spelled `maitre` — the way the code spells it.
- **Check the console's entry does not now claim to be a string.** The section it
  moved into must describe it honestly: the console is a **surface**, and the entry
  should say what it carries (the window title, the tab bar, the CLI's help) without
  filing the console itself as a string.

**F-83 — the prose stops naming the retired words.**

- The prose that explained what the backlog item replaced no longer names *special*
  or *prep note*; the `_Avoid_` line keeps them.
- **Check the prose still says something useful.** A line that lost its content to
  satisfy a rule is not a fix — the reader still has to learn what the backlog item
  replaced, and the fix points at ADR 0009 § 3. **Follow the pointer and check it
  says the thing.**

**F-84 — the spec says the same thing twice, and it says what the code does.**

- The two lists of the fingerprint's facts in `docs/03` agree — and so does the
  **third** list the runner found in the `DomainVersion` comment.
- **Check the claim against the code**: `internal/core/format.go` hashes the id
  prefixes, the directories, the root files and the field names. All three lists say
  that, and no list says anything the code does not do.

**F-85 — the version fences the writer, and the fingerprint is the tripwire.**

- `docs/03` (§ 6, § 6.1 and the version-4 history), the `DomainVersion` comment and
  its version-4 paragraph now say what the code says. `docs/07-vocabulary.md` was the
  orchestrator's and is corrected in the same spirit — **check it too.**
- **Check the code is unchanged where it must be**: `core.DomainVersion` is still
  **6**, and the fence still stands in its designed direction — a plan stamped **7**
  is refused, a plan stamped **6** is accepted.
- **Check the prose does not now overclaim the other way.** The rule is *the version
  fences the writer; the fingerprint says the facts moved* — a change to the **facts
  list** moves the fingerprint and leaves the version alone; a change to **what the
  tool writes** moves both. If the new prose says the fingerprint fences everything
  on its own, that is the same error in the opposite direction.

### 2. The sweep's count moved — check the movement

The runner reports the sweep's count went from **78 to 77**, because F-83 removed the
one prose line, and notes that **the command does not match *maître d'***.

**Run the sweep command yourself.** Check:

- that the remaining lines each name **something that still exists** — the legacy
  reader, the retired-role aliases, the tests of the old ids and the migration
  events, the dated version notes, the tripwires' own mutations;
- that the count is explained: 78 → 77 is one line removed, plus whatever the
  glossary's new `_Avoid_` lines added. **Do the arithmetic** and say whether the
  count reconciles;
- and whether the command **should** match *maître d'* now — the name is in the
  glossary's `_Avoid_` line, so a sweep command that looks for old names must either
  find it there or exclude the glossary's `_Avoid_` lines. **If the sweep command
  would now flag the glossary's own `_Avoid_` lines as drift, the command needs to
  know the difference** — and that is a finding about the command, not about the file.

### 3. The one LOW the runner left

`docs/03` § 6 links to `#the-version-guard`, and there is no such heading. The runner
called it pre-existing and left it. **Judge it:** a link to a heading that does not
exist is a broken pointer in the spec — a reader following it finds nothing. File it
if you judge it a blocker; note it if you do not.

### 4. The mascot and the shared tree

**The mascot did not land.** The art is in `logo/` (untracked), and the console
header, the CLI's help, the README and the Sous-Chef tab carry no drawing — **AC-01
moves to the drawing session.** Check that this slice committed **nothing** of the
other session's work, and that the goldens moved with the sweep alone.

### 5. Scope

Read the whole diff of the five commits. Did a fix touch anything the findings did
not name? The gate table, the plan format, the command count (**42**),
`core.DomainVersion` (**6**) and `kitchen/activity.json` must not have moved.

### 6. The ledger

Move **F-81, F-82, F-83, F-84, F-85** from open to closed, with the file and the
evidence that now holds each. **Add a row for the slice itself** — this is the last
slice of the rename, and the ledger's round table should say so. Commit the ledger.
The ledger is not swept: keep the findings' wording as it was filed.

### Verdict

- No HIGH or MEDIUM findings → `verdict: ok`.
- Any HIGH or MEDIUM → `verdict: changes-needed`.
- LOW findings: **note them in your message, do not file them.**

### Hard rules

- Do **not** edit implementation code. Read-only for everything except
  `docs/audit.md`.
- Revert every mutation you make. The tree must be exactly as the runner left it,
  except your ledger commit.
- No `git commit --no-verify`. No push. Do not touch other panes. Do not spawn agents.
- **Do not touch `logo/`** — it is another session's work.

### Notify the orchestrator

```bash
mkdir -p work/inbox
cat > work/inbox/reviewer-slice-e-cycle-2.md <<EOF
Reviewer: slice E cycle 2 done. Verdict: <ok|changes-needed>. F-81..F-85 closed: <yes|no>. Filed <N> new blockers (<F-86.. or none>). Sweep count: <the number, and whether it reconciles>. Mascot: <did it land? yes/no>. LOW notes: <any>.
EOF
git add work/inbox && git commit -q -m "inbox: reviewer reports slice-e cycle 2"
herdr agent prompt "w1Z:p1" "Reviewer: slice E cycle 2 done. Verdict: <ok|changes-needed>. F-81..F-85 closed: <yes|no>. Filed <N> new blockers (<F-86.. or none>). Sweep count: <the number, and whether it reconciles>. Mascot: <did it land? yes/no>. LOW notes: <any>."
```
