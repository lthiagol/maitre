package core

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// Journal kinds. A closed set — an unknown kind is an error, never silently ignored.
const (
	KindCreated   = "created"
	KindLifecycle = "lifecycle"
	KindEvidence  = "evidence"
	KindFinding   = "finding"
	KindComment   = "comment"
	KindDeviation = "deviation"
	KindDecision  = "decision"
	KindPromotion = "promotion"
	KindOverlay   = "overlay"
	KindControl   = "control"
	KindCriterion = "criterion"
	KindStep      = "step"
	KindRefusal   = "refusal"
)

// JournalKinds is every valid kind.
var JournalKinds = []string{
	KindCreated, KindLifecycle, KindEvidence, KindFinding, KindComment,
	KindDeviation, KindDecision, KindPromotion, KindOverlay, KindControl,
	KindCriterion, KindStep, KindRefusal,
}

// validKind reports whether k is a known journal kind.
func validKind(k string) bool {
	for _, v := range JournalKinds {
		if v == k {
			return true
		}
	}
	return false
}

// JournalEntry is one event in the plan's history.
//
// The journal is the **only** event log. A recipe's own audit holds its deviations and
// decisions — authored content — not a copy of its events. One source of truth per fact.
type JournalEntry struct {
	Seq     int             `json:"seq"`
	At      string          `json:"at"`
	Kind    string          `json:"kind"`
	Recipe  string          `json:"recipe,omitempty"`
	Actor   string          `json:"actor"`
	Session string          `json:"session,omitempty"`
	Role    string          `json:"role,omitempty"`
	Data    json.RawMessage `json:"data,omitempty"`
	// DomainVersion is the version of the domain's rules the write was made under.
	//
	// The journal is written on every mutation, so the newest entry carries the
	// version the plan was last written under. That is where the stamp costs no
	// extra write.
	DomainVersion int `json:"domain_version,omitempty"`
}

// Journal is the plan's append-only event log.
type Journal struct {
	Schema int            `json:"schema"`
	Seq    int            `json:"seq"`
	Events []JournalEntry `json:"events"`
}

// JournalSchema is the current journal format version.
const JournalSchema = 1

// PlanDomainVersion is the version of the rules the plan was last written under.
//
// The journal is written on every mutation, so the newest entry carries the version.
// A plan with no journal is a plan that does not exist yet, and it reads as version
// zero — as does a journal written before the stamp existed. **Neither is writable
// by a tool that speaks another version**: the guard refuses every plan whose version
// is not the tool's version, and only a plan with no journal at all is a plan the
// tool may create.
func PlanDomainVersion(j Journal) int {
	if len(j.Events) == 0 {
		return 0
	}
	return j.Events[len(j.Events)-1].DomainVersion
}

// LoadJournal reads the plan's journal. A missing journal is an empty one, not an error.
func (s *Store) LoadJournal() (Journal, error) {
	j := Journal{Schema: JournalSchema}
	raw, err := os.ReadFile(filepath.Join(s.Root, FileActivity))
	if err != nil {
		if os.IsNotExist(err) {
			return j, nil
		}
		return j, fmt.Errorf("read journal: %w", err)
	}
	if err := json.Unmarshal(raw, &j); err != nil {
		return j, fmt.Errorf("parse journal: %w", err)
	}
	if j.Schema != JournalSchema {
		return j, fmt.Errorf("journal schema %d is not supported (want %d)", j.Schema, JournalSchema)
	}
	return j, nil
}

// AppendJournalEvent adds an entry and returns the stored journal.
//
// The sequence number is assigned here, so it is never a caller's problem.
func (s *Store) AppendJournalEvent(e JournalEntry) (Journal, error) {
	unlock, err := s.lock()
	if err != nil {
		return Journal{}, err
	}
	defer unlock()
	return s.appendJournalEventLocked(e)
}

// appendJournalEventLocked adds an entry. The caller holds the plan lock.
func (s *Store) appendJournalEventLocked(e JournalEntry) (Journal, error) {
	if !validKind(e.Kind) {
		return Journal{}, fmt.Errorf("unknown journal kind %q", e.Kind)
	}
	j, err := s.LoadJournal()
	if err != nil {
		return Journal{}, err
	}
	e.Seq = j.Seq + 1
	if e.At == "" {
		e.At = time.Now().UTC().Format(time.RFC3339)
	}
	// Every write records the version of the rules it was made under. The newest
	// entry is the plan's stamp.
	//
	// A refusal keeps the plan's own version: it writes no artifact, so the rules
	// the plan was written under did not change. A newer tool refusing a write in
	// an older plan must not lock the older tool out of it.
	e.DomainVersion = DomainVersion
	if e.Kind == KindRefusal {
		e.DomainVersion = PlanDomainVersion(j)
	}
	j.Seq = e.Seq
	j.Events = append(j.Events, e)

	if err := writeJSONAtomic(filepath.Join(s.Root, FileActivity), j); err != nil {
		return Journal{}, err
	}
	return j, nil
}

// Since returns the entries after the given sequence number, oldest first.
//
// This is how an agent notices what changed: ask for what is newer than the last
// sequence number it saw.
//
// The list is never nil: no newer event is an empty list, and the CLI marshals an
// empty list as `[]` and a nil slice as `null`. A client cannot iterate `null`.
func (j Journal) Since(seq int) []JournalEntry {
	out := make([]JournalEntry, 0, len(j.Events))
	for _, e := range j.Events {
		if e.Seq > seq {
			out = append(out, e)
		}
	}
	sort.Slice(out, func(i, k int) bool { return out[i].Seq < out[k].Seq })
	return out
}

// Recent returns the last n entries, newest first.
func (j Journal) Recent(n int) []JournalEntry {
	out := make([]JournalEntry, 0, n)
	for i := len(j.Events) - 1; i >= 0 && len(out) < n; i-- {
		out = append(out, j.Events[i])
	}
	return out
}
