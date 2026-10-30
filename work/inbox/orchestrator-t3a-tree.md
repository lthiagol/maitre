Orchestrator: the change you found is the **reviewer's** — it is its mutation #1 from the
review brief ("make the promotion of an idea succeed"), and it is **in flight right
now**. The reviewer is watching what the console's footer and banner do when the
domain's answer is taken away.

**So: do not touch the file.** The reviewer owns that line and will revert it when
its experiment ends; a second revert would lose the experiment and hide its evidence.

**Your receipt stands.** `make check` green, `make check-race` green, tree clean — all
true at `8ac6459`, and your commits are intact and unchanged. The tree stopped being
clean because of the reviewer's mutation, not because of anything you did.

**What you did was exactly right**, and it is now written into the round's protocol:
`work/tui/README.md`'s *The tree is shared* says what a lane does when it finds a
stranger's edit — it does not touch the file, it reports what changed and what it
breaks, it says whose it probably is, and it stands by.

**Stand by.** The reviewer's verdict is what comes next, and it may still be running.
When it lands I will tell you what it says, and whether T3a is closed or needs a fix
cycle. Do not start anything, and do not re-run the gate while the reviewer's mutation
is in the tree — it would fail, and the failure would be the reviewer's experiment,
not yours.

The receipt you wrote is in `work/inbox/runner-tui-t3a-tree.md`, and it will be part of
T3a's record — a shared tree needs a lane that stops and says so.
