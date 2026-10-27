package core

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// LockName is the advisory lock file inside the plan directory.
const LockName = ".lock"

// staleLockAfter is how long a lock file may sit before it is assumed abandoned.
//
// It must exceed the longest hold. The lock is only ever held across a read, a write
// and a journal append — no command runs inside it — so the hold is milliseconds and
// this is generous. It is a variable so a test can shorten it.
var staleLockAfter = 30 * time.Second

// lock takes the plan directory's advisory lock.
//
// It is an exclusive-create lock file: portable, dependency-free, and good enough for the
// one place we need it — the journal's read-modify-write. The window is tiny, so the
// chance of two writers colliding is small but not zero, which is why the lock exists at
// all.
//
// The lock is advisory: it stops two `maitre-cli` processes from interleaving. It does not
// stop a hand edit — which the plan-zone rule already forbids.
//
// Every write takes this lock, so it is also where the version guard lives: a tool
// whose rules are older than the plan's refuses to write, rather than applying weaker
// rules than the plan was written under.
func (s *Store) lock() (func(), error) {
	path := filepath.Join(s.Root, LockName)
	deadline := time.Now().Add(5 * time.Second)

	for {
		f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
		if err == nil {
			fmt.Fprintf(f, "%d\n", os.Getpid())
			f.Close()
			unlock := func() { os.Remove(path) }
			if err := s.checkDomainVersion(); err != nil {
				unlock()
				return nil, err
			}
			return unlock, nil
		}
		if !os.IsExist(err) {
			return nil, fmt.Errorf("lock %s: %w", path, err)
		}

		// A lock left by a crashed process must not wedge the plan forever.
		if info, statErr := os.Stat(path); statErr == nil && time.Since(info.ModTime()) > staleLockAfter {
			os.Remove(path)
			continue
		}
		if time.Now().After(deadline) {
			return nil, fmt.Errorf("plan is locked by another process (%s); try again", path)
		}
		time.Sleep(10 * time.Millisecond)
	}
}

// checkDomainVersion refuses a write to a plan written under newer rules.
//
// A tool whose rules are older than the plan's would let through a gate the newer tool
// refuses: the plan was written under stronger rules, and the write would apply weaker
// ones. Reading is allowed — the chef can still look at the plan.
func (s *Store) checkDomainVersion() error {
	j, err := s.LoadJournal()
	if err != nil {
		return err
	}
	plan := PlanDomainVersion(j)
	if plan > DomainVersion {
		return fmt.Errorf("this plan was written under domain version %d and this tool speaks version %d: upgrade the tool, or read the plan without writing to it", plan, DomainVersion)
	}
	return nil
}
