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
// writes a plan only if the plan's version is the tool's version, rather than
// applying rules the plan was not written under.
// lock takes the plan lock and refuses a plan written under other rules.
func (s *Store) lock() (func(), error) {
	return s.lockUnder(s.checkWritable)
}

// checkWritable refuses a write to a plan the tool does not support.
//
// Two fences, and both hold:
//
//   - the plan's **version** must be the tool's version — a stale tool does not
//     write a plan it does not understand; and
//   - the plan must carry **no retired artifact** — a plan in a retired format is
//     not supported, and *not supported* means **no writes**, not only a message at
//     read time.
//
// The second fence was missing: the reader refused such a plan and the writer wrote
// to it anyway, so the plan kept taking records in a format the tool does not
// speak. A message is not a fence.
func (s *Store) checkWritable() error {
	if err := s.checkDomainVersion(); err != nil {
		return err
	}
	return s.checkRetiredPlan()
}

// lockUnder takes the plan lock, with the given guard.
//
// The guard runs once the lock is held: the plan's version cannot change under the
// lock, so the check and the write see the same plan.
func (s *Store) lockUnder(guard func() error) (func(), error) {
	path := filepath.Join(s.Root, LockName)
	deadline := time.Now().Add(5 * time.Second)

	for {
		f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
		if err == nil {
			fmt.Fprintf(f, "%d\n", os.Getpid())
			f.Close()
			unlock := func() { os.Remove(path) }
			if guard != nil {
				if err := guard(); err != nil {
					unlock()
					return nil, err
				}
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

// checkDomainVersion refuses a write to a plan written under other rules.
//
// **One rule:** a tool writes a plan only if the plan's version **is** the tool's
// version. A plan stamped with an older version is not supported, and there is no
// upgrade path — a newer tool does not write an older plan, because it would apply
// rules the plan was not written under.
//
// The refusal says what the chef can do, and it promises no migration: a plan under
// another version is read, and written by the tool that speaks its version.
//
// Reading is always allowed: the chef can still look at a plan.
//
// A plan with **no journal** is a plan that does not exist yet. There is nothing to
// be incompatible with, and the first write creates it.
func (s *Store) checkDomainVersion() error {
	j, err := s.LoadJournal()
	if err != nil {
		return err
	}
	if len(j.Events) == 0 {
		return nil
	}
	plan := PlanDomainVersion(j)
	if plan == DomainVersion {
		return nil
	}
	if plan > DomainVersion {
		return fmt.Errorf("this plan was written under domain version %d and this tool speaks version %d: upgrade the tool, or read the plan without writing to it", plan, DomainVersion)
	}
	return fmt.Errorf("this plan was written under domain version %d and this tool speaks version %d: this tool does not write plans written under another version — read it without writing to it, or write the plan with the tool that speaks its version", plan, DomainVersion)
}

// checkDeclarationForward refuses a declaration that would move a plan backwards.
//
// The declaration is the one write that crosses the version fence, and it may only
// move the plan **forward**: a tool may not stamp a plan with a version it does not
// speak, and a plan under a newer version is not this tool's to move.
//
// The retired fence holds here too: a plan that still carries a retired artifact is
// not supported, and the declaration records a plan that has already moved — so the
// plan is migrated before it is declared.
func (s *Store) checkDeclarationForward() error {
	j, err := s.LoadJournal()
	if err != nil {
		return err
	}
	plan := PlanDomainVersion(j)
	if plan > DomainVersion {
		return fmt.Errorf("this plan is under domain version %d and this tool speaks %d: a declaration only moves a plan forward", plan, DomainVersion)
	}
	return s.checkRetiredPlan()
}
