package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/lthiagol/maitre/internal/settings"
)

// TestSettingsShowPrintsTheRecord covers the read.
//
// The record is the tool's file: the CLI reads it and prints it as one document,
// the way every list the CLI prints is printed.
func TestSettingsShowPrintsTheRecord(t *testing.T) {
	dir := newPlan(t)

	res := cli(t, dir, "settings", "show")
	if !res.ok(t) {
		t.Fatalf("settings show failed:\n%s", res.Stderr)
	}
	doc := res.doc(t)
	if doc["ok"] != true {
		t.Fatalf("the document is %v, want ok", doc)
	}
	rec, ok := doc["settings"].(map[string]any)
	if !ok {
		t.Fatalf("the document has no settings: %s", res.Stdout)
	}
	if rec["version"] != float64(settings.Version) {
		t.Errorf("the record's version is %v, want %d", rec["version"], settings.Version)
	}
	view, ok := rec["view"].(map[string]any)
	if !ok {
		t.Fatalf("the record has no view section: %s", res.Stdout)
	}
	if view["theme"] != settings.DefaultTheme || view["order"] != settings.DefaultOrder {
		t.Errorf("the record's view is %v, want the theme %q and the order %q",
			view, settings.DefaultTheme, settings.DefaultOrder)
	}
}

// TestSettingsShowReadsTheFile covers the record the chef chose.
//
// The console writes the file; the CLI reads what is there.
func TestSettingsShowReadsTheFile(t *testing.T) {
	dir := newPlan(t)
	rec := settings.Record{View: settings.View{Theme: "Dracula", Order: "lane", ShowServed: false}}
	if err := settings.Write(settings.Path(filepath.Join(dir, "kitchen")), rec); err != nil {
		t.Fatalf("write the settings: %v", err)
	}

	res := cli(t, dir, "settings", "show")
	if !res.ok(t) {
		t.Fatalf("settings show failed:\n%s", res.Stderr)
	}
	doc := res.doc(t)
	raw, err := json.Marshal(doc["settings"])
	if err != nil {
		t.Fatal(err)
	}
	var got settings.Record
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatal(err)
	}
	if got.View != rec.View {
		t.Fatalf("the record is %+v, want %+v", got.View, rec.View)
	}
}

// TestSettingsShowRefusesAnUnknownVersion covers the version fence.
//
// The tool does not guess what it cannot read, and it says which version it
// refused.
func TestSettingsShowRefusesAnUnknownVersion(t *testing.T) {
	dir := newPlan(t)
	path := filepath.Join(dir, "kitchen", settings.FileName)
	if err := os.WriteFile(path, []byte(`{"version": 99, "view": {}}`), 0o644); err != nil {
		t.Fatal(err)
	}

	res := cli(t, dir, "settings", "show")
	if res.Code == 0 {
		t.Fatalf("an unknown version should be refused:\n%s", res.Stdout)
	}
	if !strings.Contains(res.Stderr, "99") {
		t.Fatalf("the refusal should name the version it refused: %s", res.Stderr)
	}
}

// TestSettingsShowNeedsAPlan covers the plan's folder.
//
// The record lives in the plan's directory, so there is no record to read outside
// a plan.
func TestSettingsShowNeedsAPlan(t *testing.T) {
	res := cli(t, t.TempDir(), "settings", "show")
	if res.Code == 0 {
		t.Fatalf("settings show outside a plan should fail:\n%s", res.Stdout)
	}
	if !strings.Contains(res.Stderr, "no plan") {
		t.Fatalf("the failure should say why: %s", res.Stderr)
	}
}
