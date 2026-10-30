Runner: T3a — a file I did not write has changed in the shared tree. Reporting, not touching.

**What changed.** `internal/core/backlog.go`, one line, uncommitted in the working tree:

```
-	if p := PromotionOf(*v); !p.Allowed {
+	if p := PromotionOf(*v); false && !p.Allowed {
```

The file's mtime is 18:05, after my receipt (8ac6459, committed with a clean tree) and during my `make check-race` run. I did not make the edit and I have not touched the file.

**What it does.** It disables the domain's refusal of an idea — the rule AC-03 landed — so an idea promotes again. The core tests now fail on it:

```
$ go test ./internal/core -run 'Promotion|Promot'
--- FAIL: TestAnIdeaIsNotPromoted        an idea must not be promoted
--- FAIL: TestThePromotionAnswerIsOnePlace
--- FAIL: TestAPromotedItemMayNotBePromotedAgain
--- FAIL: TestAPromotionIsOneWay
FAIL
```

**What I did.** Nothing to the file: it is another lane's working tree, and the rule is to stop and report. The runner's commits are intact and unchanged — `4dc8079`, `1b550c6`, `1804ccc`, `b873821`, `5edae71` (slice) and `a0ca384`, `051081c`, `8ac6459` (work and receipt). `git diff` shows the edit is **uncommitted**, so restoring it is one `git checkout -- internal/core/backlog.go` away whenever the lane that owns it is done.

**Why it may be deliberate.** If the reviewer is checking that the console's refusal comes from the **domain** — by removing the domain's answer and watching the console's footer and banner change with it — then this is the experiment, and it wants the file back the way it was when the experiment ends.

**What the runner's receipt said at the time:** `make check` green, `make check-race` green, the tree clean — all true at 8ac6459. The tree is no longer clean, and it is not the runner's doing.
