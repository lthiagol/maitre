package core

import (
	"reflect"
	"testing"
)

// TestEveryFieldHasAWriterOrIsGone covers C3.
//
// The model carried fields nothing could write: a branch the plan has no repository
// for, a completion stamp nobody stamped, a step's evidence that duplicated the
// criterion's evidence, a finding's "fixed in" that only ever got cleared. A field with
// no writer is a promise the code does not keep — v1's `needs_regrooming` had a field,
// an event and a renderer, and no way to set it.
//
// The list below is the model's whole field set. Adding a field means adding it here,
// where the question "who writes it?" has to be answered.
func TestEveryFieldHasAWriterOrIsGone(t *testing.T) {
	want := map[string][]string{
		"Recipe": {"ID", "Title", "Slug", "Created", "Updated", "Lifecycle", "Intent",
			"Scope", "Acceptance", "Guidance", "WorkPackages", "Questions", "DependsOn",
			"Execution", "Completion", "Review", "Audit", "Approval", "Overlays"},
		"Intent":    {"Outcome", "Problem"},
		"Scope":     {"InScope", "OutOfScope"},
		"Criterion": {"ID", "Description", "Verification", "Evidence", "Deferral"},
		"Step": {"ID", "WorkPackage", "Order", "Action", "Files", "DependsOn", "Covers",
			"DoneWhen", "Status", "Deviation"},
		"WorkPackage": {"ID", "Name", "Goal", "Rollback"},
		"Execution":   {"ClaimedBy", "ClaimedAt"},
		"Completion":  {"CompletedAt"},
		"Review":      {"Findings", "Comments", "Rounds"},
		"Finding": {"ID", "Phase", "Severity", "Category", "Description", "Anchor",
			"Status", "Created", "Resolved", "Resolution"},
		"Anchor":    {"File", "Line", "Side"},
		"Comment":   {"ID", "Author", "Body", "At", "Finding"},
		"Round":     {"N", "Reviewer", "Verdict", "At", "Notes"},
		"Audit":     {"Deviations", "Decisions", "Documentation"},
		"Deviation": {"At", "Step", "Note", "By"},
		"Decision":  {"ID", "Summary", "Context", "At"},
		"Question":  {"ID", "Question", "Status", "Answer"},
		"Approval":  {"By", "At", "Fingerprint"},
		"Overlays":  {"Blocked", "Deferred", "Cancelled"},
		"Overlay":   {"Reason", "By", "At"},
		"Evidence": {"Command", "ExitCode", "TimedOut", "OutputDigest", "Matched",
			"MatchUnit", "Match", "SignedOff", "At", "By"},
		"Deferral": {"Reason", "FollowUp"},
	}

	got := modelFields()
	for name, fields := range want {
		gotFields, ok := got[name]
		if !ok {
			t.Errorf("the model has no %s", name)
			continue
		}
		if !sameFields(gotFields, fields) {
			t.Errorf("%s fields = %v, want %v", name, gotFields, fields)
		}
	}
	for name := range got {
		if _, ok := want[name]; !ok {
			t.Errorf("the model has a new struct %q — say who writes its fields", name)
		}
	}
}

// modelFields is every exported field of every struct in the recipe model.
func modelFields() map[string][]string {
	out := map[string][]string{}
	collectFields(reflect.TypeOf(Recipe{}), out)
	return out
}

// collectFields walks a struct and the structs it holds.
func collectFields(t reflect.Type, out map[string][]string) {
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return
	}
	names := make([]string, 0, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if !f.IsExported() {
			continue
		}
		names = append(names, f.Name)
	}
	out[t.Name()] = names
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if !f.IsExported() {
			continue
		}
		ft := f.Type
		for ft.Kind() == reflect.Slice || ft.Kind() == reflect.Pointer {
			ft = ft.Elem()
		}
		if ft.Kind() == reflect.Struct {
			collectFields(ft, out)
		}
	}
}

// sameFields compares two field lists, order aside.
func sameFields(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	seen := map[string]bool{}
	for _, f := range a {
		seen[f] = true
	}
	for _, f := range b {
		if !seen[f] {
			return false
		}
	}
	return true
}
