package settings

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// planDir makes a directory to hold the settings file.
func planDir(t *testing.T) string {
	t.Helper()
	return t.TempDir()
}

// TestAMissingFileYieldsTheDefaults covers the read of a plan nobody has set up.
//
// A missing file is not an error: the chef has not chosen anything yet, and the
// defaults are the answer until they do.
func TestAMissingFileYieldsTheDefaults(t *testing.T) {
	path := Path(planDir(t))

	rec, err := Read(path)
	if err != nil {
		t.Fatalf("a missing file should not be an error: %v", err)
	}
	want := Defaults()
	if rec != want {
		t.Fatalf("the defaults are %+v, want %+v", rec, want)
	}
	if rec.View.Theme != DefaultTheme || rec.View.Order != DefaultOrder || !rec.View.ShowServed {
		t.Fatalf("the defaults are %+v, want the theme %q, the order %q and served shown",
			rec.View, DefaultTheme, DefaultOrder)
	}
}

// TestTheFileHasItsOwnVersion covers the shape of the file.
//
// The settings carry a version of their own, the way the plan's records do. It is
// the schema's version, not the plan's — the plan's format does not move.
func TestTheFileHasItsOwnVersion(t *testing.T) {
	dir := planDir(t)
	path := Path(dir)

	if err := Write(path, Defaults()); err != nil {
		t.Fatalf("write: %v", err)
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var doc map[string]any
	if err := json.Unmarshal(raw, &doc); err != nil {
		t.Fatalf("the file is not a JSON document: %v\n%s", err, raw)
	}
	if doc["version"] != float64(Version) {
		t.Fatalf("the file's version is %v, want %d", doc["version"], Version)
	}
	view, ok := doc["view"].(map[string]any)
	if !ok {
		t.Fatalf("the file has no view section: %s", raw)
	}
	for _, key := range []string{"theme", "order", "showServed"} {
		if _, ok := view[key]; !ok {
			t.Errorf("the view section has no %q: %s", key, raw)
		}
	}
}

// TestTheWriterStampsItsOwnVersion covers the version on the way out.
//
// The writer knows the schema it writes, so the version on the file is this
// build's — a record that carries another version is not written as it was.
func TestTheWriterStampsItsOwnVersion(t *testing.T) {
	path := Path(planDir(t))

	rec := Defaults()
	rec.Version = 99
	if err := Write(path, rec); err != nil {
		t.Fatalf("write: %v", err)
	}

	got, err := Read(path)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if got.Version != Version {
		t.Fatalf("the file's version is %d, want %d", got.Version, Version)
	}
}

// TestAnUnknownVersionIsRefusedByName covers the version fence.
//
// The tool does not guess what it cannot read: a safe unknown value disabled a
// whole feature with green tests. The refusal names the version it refused.
func TestAnUnknownVersionIsRefusedByName(t *testing.T) {
	path := Path(planDir(t))
	raw := `{"version": 99, "view": {"theme": "Kitchen", "order": "triage", "showServed": true}}`
	if err := os.WriteFile(path, []byte(raw), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := Read(path)
	if err == nil {
		t.Fatal("an unknown version should be refused")
	}
	if !strings.Contains(err.Error(), "99") {
		t.Fatalf("the refusal should name the version it refused: %v", err)
	}
	if !strings.Contains(err.Error(), "not one this build has") {
		t.Fatalf("the refusal should say why: %v", err)
	}
}

// TestAMissingVersionIsRefused covers a file that names no version.
//
// The file's version is part of its schema. A file without one is a file this
// build cannot read — and a safe guess would be the silent fallback the house
// learned to refuse.
func TestAMissingVersionIsRefused(t *testing.T) {
	path := Path(planDir(t))
	if err := os.WriteFile(path, []byte(`{"view": {}}`), 0o644); err != nil {
		t.Fatal(err)
	}

	if _, err := Read(path); err == nil {
		t.Fatal("a file with no version should be refused")
	}
}

// TestUnknownKeysAreIgnored covers the forward compatibility.
//
// A file written by a newer build, with keys this build does not know, still
// reads: the unknown keys are ignored, and the ones this build knows are read.
func TestUnknownKeysAreIgnored(t *testing.T) {
	path := Path(planDir(t))
	raw := `{
	  "version": 1,
	  "crew": {"harness": "codex", "count": 3},
	  "view": {"theme": "Dracula", "order": "lane", "showServed": false, "density": "compact"}
	}`
	if err := os.WriteFile(path, []byte(raw), 0o644); err != nil {
		t.Fatal(err)
	}

	rec, err := Read(path)
	if err != nil {
		t.Fatalf("a file with unknown keys should read: %v", err)
	}
	if rec.View.Theme != "Dracula" || rec.View.Order != "lane" || rec.View.ShowServed {
		t.Fatalf("the record is %+v, want the theme Dracula, the order lane and served hidden", rec.View)
	}
}

// TestTheWriteIsAtomic covers the write.
//
// The bytes go to a temporary file and are renamed over the old one, so a
// half-written file cannot lose the chef's choices. The write leaves nothing
// behind, and a write that cannot happen leaves the old file alone.
func TestTheWriteIsAtomic(t *testing.T) {
	dir := planDir(t)
	path := Path(dir)

	if err := Write(path, Defaults()); err != nil {
		t.Fatalf("write: %v", err)
	}
	rec := Defaults()
	rec.View.Theme = "Dracula"
	if err := Write(path, rec); err != nil {
		t.Fatalf("write: %v", err)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 || entries[0].Name() != FileName {
		names := make([]string, 0, len(entries))
		for _, e := range entries {
			names = append(names, e.Name())
		}
		t.Fatalf("the directory holds %v, want only %s — a write leaves nothing behind", names, FileName)
	}

	// A write that cannot happen leaves the file as it was: the old bytes are
	// still the whole file, not half of it.
	missing := Path(filepath.Join(dir, "nowhere", FileName))
	if err := Write(missing, rec); err == nil {
		t.Fatal("a write into a directory that is not there should fail")
	}
	got, err := Read(path)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if got.View.Theme != "Dracula" {
		t.Fatalf("the file is %+v, want the theme Dracula", got.View)
	}
}

// TestTheRoundTrip covers the record's own life: written, read back, the same.
func TestTheRoundTrip(t *testing.T) {
	path := Path(planDir(t))

	rec := Record{View: View{Theme: "Alucard", Order: "plain", ShowServed: false}}
	if err := Write(path, rec); err != nil {
		t.Fatalf("write: %v", err)
	}

	got, err := Read(path)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if got.View != rec.View {
		t.Fatalf("the record is %+v, want %+v", got.View, rec.View)
	}
}
