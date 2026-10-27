package core

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestAVerificationDoesNotHoldThePlanLock covers the lock the verification used to
// hold.
//
// The verification ran *inside* the plan lock, and the lock went stale after 30s
// while a verification may run for ten minutes. A second writer stole the lock of a
// live process, and the verification's save then clobbered the second writer's
// change: the write was silently lost and the journal and the artifact disagreed.
//
// The fix: no command runs while the plan lock is held. The lock is taken only to
// append the evidence, against the recipe as it is at that moment.
//
// The stale threshold is lowered rather than waited out: what matters is the
// relationship between the threshold and the hold. The hold must never reach the
// threshold, so a live lock can never be stolen.
func TestVerificationDoesNotHoldThePlanLock(t *testing.T) {
	ops, store := newOps(t)

	restore := staleLockAfter
	staleLockAfter = 100 * time.Millisecond
	defer func() { staleLockAfter = restore }()

	// The command marks that it started, then keeps running for a second — longer
	// than the stale window.
	marker := filepath.Join(t.TempDir(), "command-ran")
	s := spec()
	s.Criteria = []CriterionSpec{{
		Description:  "a slow check",
		Verification: fmt.Sprintf("touch %q; sleep 1", marker),
	}}
	r, err := ops.CreateRecipe(s)
	if err != nil {
		t.Fatal(err)
	}

	done := make(chan error, 1)
	go func() {
		_, err := ops.VerifyCriterion(r.ID, "AC-01")
		done <- err
	}()

	waitForFile(t, marker)

	// The plan lock is free: no command runs while it is held. A lock held here
	// would be stale by now, and the next writer would steal it.
	if _, err := os.Stat(filepath.Join(store.Root, LockName)); !os.IsNotExist(err) {
		t.Fatal("the plan lock is held while the verification command runs")
	}

	// The second writer finds the plan unlocked, so it is not let in by a stale
	// lock — and its write is not lost.
	if _, err := ops.BlockRecipe(r.ID, "the chef blocked this"); err != nil {
		t.Fatalf("a concurrent write should not be refused: %v", err)
	}

	if err := <-done; err != nil {
		t.Fatalf("verify: %v", err)
	}

	got, err := ops.GetRecipe(r.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Overlays.Blocked == nil {
		t.Fatal("the concurrent write was clobbered by the verification's save")
	}
	if len(got.Acceptance[0].Evidence) != 1 {
		t.Fatalf("the evidence was not recorded: %+v", got.Acceptance[0].Evidence)
	}
}

// waitForFile waits until the file exists, or the test fails.
func waitForFile(t *testing.T, path string) {
	t.Helper()
	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		if _, err := os.Stat(path); err == nil {
			return
		}
		time.Sleep(5 * time.Millisecond)
	}
	t.Fatalf("the verification command never ran: %s", path)
}
